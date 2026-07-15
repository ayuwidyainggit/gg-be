package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/config/env"
	"inventory/pkg/constant"
	"inventory/pkg/conversion"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/pkg/validation"
	"inventory/repository"
	"log"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StockDisposalService interface {
	Store(request entity.CreateStockDisposalBody) (response entity.StockDisposalResponse, err error)
	Detail(sdID int64, custID, parentCustId string) (response entity.StockDisposalDetailResponse, err error)
	List(dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string) (data []entity.StockDisposalListResponse, total int64, lastPage int, err error)
	ProductLookup(dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string) (data []entity.StockDisposalProductLookupResponse, total int64, lastPage int, err error)
}

func NewStockDisposalService(
	stockDisposalRepository repository.StockDisposalRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction,
	validator *validation.Validate,
	config env.ConfigEnv) *stockDisposalServiceImpl {
	return &stockDisposalServiceImpl{
		StockDisposalRepository: stockDisposalRepository,
		StockRepository:         stockRepository,
		Transaction:             transaction,
		Validator:               validator,
		Config:                  config,
	}
}

type stockDisposalServiceImpl struct {
	StockDisposalRepository repository.StockDisposalRepository
	StockRepository         repository.StockRepository
	Transaction             repository.Dbtransaction
	Validator               *validation.Validate
	Config                  env.ConfigEnv
}

// roundPrice rounds a float64 value with custom logic: <= 5 round down, >5 round up
func roundPrice(value float64) float64 {
	// Get integer part
	intPart := math.Floor(value)
	// Get decimal part (first 2 digits after decimal point)
	decimalPart := (value - intPart) * 100
	// Get first digit after decimal point
	firstDigit := math.Floor(decimalPart / 10)

	// If first digit <= 5, round down (floor), else round up (ceil)
	if firstDigit <= 5 {
		return math.Floor(value)
	}
	return math.Ceil(value)
}

// calculateTotals calculates total subtotal and VAT value from products
func calculateTotals(products []entity.CreateStockDisposalProductBody) (totalSubTotal, totalVatValue float64) {
	for _, product := range products {
		subTotal := roundPrice(product.GrossPrice)
		vatValue := roundPrice(product.VatValue)
		totalSubTotal += subTotal
		totalVatValue += vatValue
	}
	return totalSubTotal, totalVatValue
}

// buildQtyUnit creates QtyUnit from product and productModel
func buildQtyUnit(product entity.CreateStockDisposalProductBody, productModel *model.Product) *conversion.QtyUnit {
	return &conversion.QtyUnit{
		Qty1:      int(product.Qty1),
		Qty2:      int(product.Qty2),
		Qty3:      int(product.Qty3),
		ConvUnit2: int(productModel.ConvUnit2),
		ConvUnit3: int(productModel.ConvUnit3),
	}
}

func isRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows)
}

func productValidationError(productIndex string, err error) string {
	if isRecordNotFoundError(err) {
		return fmt.Sprintf("%s: product not found", productIndex)
	}

	return fmt.Sprintf("%s: failed to validate product: %v", productIndex, err)
}

// mapSupplierWarehouseInfo maps supplier and warehouse info from model to Store response
func mapSupplierWarehouseInfo(stockDisposalList *model.StockDisposalList, response *entity.StockDisposalResponse) {
	if stockDisposalList.SupCode != nil {
		response.SupCode = *stockDisposalList.SupCode
	}
	if stockDisposalList.SupName != nil {
		response.SupName = *stockDisposalList.SupName
	}
	if stockDisposalList.WhCode != nil {
		response.WhCode = *stockDisposalList.WhCode
	}
	if stockDisposalList.WhName != nil {
		response.WhName = *stockDisposalList.WhName
	}
}

// mapSupplierWarehouseInfoToList maps supplier and warehouse info from model to List response
func mapSupplierWarehouseInfoToList(stockDisposal *model.StockDisposalList, listResponse *entity.StockDisposalListResponse) {
	if stockDisposal.WhCode != nil {
		listResponse.WhCode = *stockDisposal.WhCode
	}
	if stockDisposal.WhName != nil {
		listResponse.WhName = *stockDisposal.WhName
	}
	if stockDisposal.SupCode != nil {
		listResponse.SupCode = *stockDisposal.SupCode
	}
	if stockDisposal.SupName != nil {
		listResponse.SupName = *stockDisposal.SupName
	}
}

// mapSupplierWarehouseInfoToDetail maps supplier and warehouse info from model to Detail response
func mapSupplierWarehouseInfoToDetail(stockDisposalList *model.StockDisposalList, response *entity.StockDisposalDetailResponse) {
	if stockDisposalList.SupCode != nil {
		response.SupCode = *stockDisposalList.SupCode
	}
	if stockDisposalList.SupName != nil {
		response.SupName = *stockDisposalList.SupName
	}
	if stockDisposalList.WhCode != nil {
		response.WhCode = *stockDisposalList.WhCode
	}
	if stockDisposalList.WhName != nil {
		response.WhName = *stockDisposalList.WhName
	}
}

// validateFileSize validates file sizes for all products and returns validation errors
// Note: This is used by the Store method. File URLs are now passed directly from FE.
// If FileUrl is provided, file size validation is skipped (assumes FE already validated before upload).
func validateFileSize(products []entity.CreateStockDisposalProductBody) []string {
	var validationErrors []string
	for i, product := range products {
		if product.UploadFile != nil {
			// Skip file size validation if FileUrl is provided (URL-based upload)
			// Frontend already validated file size before uploading to OBS
			if product.UploadFile.FileUrl != "" {
				continue
			}

			// Validate file size only for direct file uploads (if any)
			if product.UploadFile.FileSize <= 0 {
				fileName := getFileNameOrDefault(product.UploadFile)
				validationErrors = append(validationErrors,
					fmt.Sprintf("product index %d: file %s has invalid size: %d bytes", i, fileName, product.UploadFile.FileSize))
			}
			if product.UploadFile.FileSize > constant.MAX_FILE_SIZE_BYTES {
				fileName := getFileNameOrDefault(product.UploadFile)
				validationErrors = append(validationErrors,
					fmt.Sprintf("product index %d: file %s exceeds maximum size: %d bytes (max: %d bytes)",
						i, fileName, product.UploadFile.FileSize, constant.MAX_FILE_SIZE_BYTES))
			}
		}
	}
	return validationErrors
}

// validateAndParseDate parses date string from YYYY-MM-DD format to time.Time
func validateAndParseDate(dateStr string) (time.Time, error) {
	disposalDate, err := str.DateStrToRfc3339String(dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	disposalDateParsed, err := time.Parse(time.RFC3339, disposalDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}
	return disposalDateParsed, nil
}

// formatFileSize converts bytes to human-readable format (e.g., "1mb")
func formatFileSize(sizeBytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if sizeBytes >= GB {
		return fmt.Sprintf("%.2fgb", float64(sizeBytes)/float64(GB))
	} else if sizeBytes >= MB {
		return fmt.Sprintf("%.2fmb", float64(sizeBytes)/float64(MB))
	} else if sizeBytes >= KB {
		return fmt.Sprintf("%.2fkb", float64(sizeBytes)/float64(KB))
	}
	return fmt.Sprintf("%db", sizeBytes)
}

// getFileNameOrDefault returns file name or "unknown" if empty
func getFileNameOrDefault(uploadFile *entity.CreateStockDisposalFileBody) string {
	if uploadFile == nil || uploadFile.FileName == "" {
		return "unknown"
	}
	return uploadFile.FileName
}

// buildStockDisposalModel creates stock disposal model from request
func buildStockDisposalModel(request entity.CreateStockDisposalBody, warehouse *model.WarehouseStockWhList, disposalDate time.Time, totalSubTotal, totalVatValue float64) model.StockDisposal {
	return model.StockDisposal{
		CustID:       warehouse.CustID,
		SupID:        request.SupID,
		WhID:         request.WhID,
		StockType:    warehouse.StockType,
		GrNo:         request.GrNo,
		Note:         request.Note,
		DisposalDate: &disposalDate,
		TrCode:       constant.TR_CODE_STOCK_DISPOSAL,
		SubTotal:     roundPrice(totalSubTotal),
		VatValue:     roundPrice(totalVatValue),
		Total:        roundPrice(totalSubTotal + totalVatValue),
		CreatedBy:    request.CreatedBy,
		UpdatedBy:    &request.CreatedBy,
	}
}

// createStockDisposalDetails creates stock disposal details and returns stock update entities
func (service *stockDisposalServiceImpl) createStockDisposalDetails(request entity.CreateStockDisposalBody, stockDisposalModel model.StockDisposal, productMap map[int64]*model.Product, disposalDate time.Time, txCtx context.Context) ([]*entity.StockUpdate, error) {
	var stockUpdateEntities []*entity.StockUpdate

	for _, product := range request.Products {
		productModel := productMap[product.ProID]
		if productModel == nil {
			return nil, fmt.Errorf("product model not found for product %d", product.ProID)
		}

		qtyUnit := buildQtyUnit(product, productModel)
		totalQty, err := qtyUnit.ToTotalQuantity()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total quantity for product %d: %w", product.ProID, err)
		}

		// Calculate prices with rounding
		grossPrice := roundPrice(product.GrossPrice)
		vatValue := roundPrice(product.VatValue)
		subTotal := roundPrice(product.GrossPrice)

		// Create detail model
		var detailModel model.StockDisposalDetail
		detailModel.CustID = stockDisposalModel.CustID
		detailModel.SdID = stockDisposalModel.SdID
		detailModel.ProID = product.ProID
		detailModel.UnitID1 = product.UnitID1
		detailModel.UnitID2 = product.UnitID2
		detailModel.UnitID3 = product.UnitID3
		detailModel.Qty1 = product.Qty1
		detailModel.Qty2 = product.Qty2
		detailModel.Qty3 = product.Qty3
		detailModel.PurchPrice1 = product.PurchPrice1
		detailModel.PurchPrice2 = product.PurchPrice2
		detailModel.PurchPrice3 = product.PurchPrice3
		detailModel.GrossPrice = grossPrice
		detailModel.Vat = product.Vat
		detailModel.VatValue = vatValue
		detailModel.SubTotal = subTotal

		// Set file fields only if upload_file is provided
		if product.UploadFile != nil {
			detailModel.FileName = &product.UploadFile.FileName
			detailModel.FileType = &product.UploadFile.FileType
			detailModel.MediaCategory = &product.UploadFile.MediaCategory
			if product.UploadFile.FileUrl != "" {
				detailModel.FileUrl = &product.UploadFile.FileUrl
			}
			detailModel.FileSize = &product.UploadFile.FileSize
		}

		detailModel.CreatedBy = request.CreatedBy

		detailCreated, err := service.StockDisposalRepository.CreateDetail(txCtx, &detailModel)
		if err != nil {
			return nil, err
		}

		// Create stock update entity
		stockUpdateEntity := entity.StockUpdate{
			CustID:    stockDisposalModel.CustID,
			WhID:      request.WhID,
			ProID:     product.ProID,
			StockDate: disposalDate,
			TrCode:    constant.TR_CODE_STOCK_DISPOSAL,
			TrNo:      stockDisposalModel.SdNumber,
			ItemCdn:   1,
			QtyOut:    float64(totalQty),
			QtyIn:     0,
			UnitPrice: product.PurchPrice1,
			RefDetId:  detailCreated.SdDetailID,
		}

		stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
	}

	return stockUpdateEntities, nil
}

// buildStockDisposalResponse builds response from stock disposal list model
func buildStockDisposalResponse(stockDisposalList *model.StockDisposalList) entity.StockDisposalResponse {
	var response entity.StockDisposalResponse
	structs.Automapper(stockDisposalList, &response)

	// Format date
	if stockDisposalList.DisposalDate != nil {
		response.DisposalDate = stockDisposalList.DisposalDate.Format(constant.DATE_FORMAT_DISPLAY)
	}

	// Map supplier/warehouse info
	mapSupplierWarehouseInfo(stockDisposalList, &response)

	// Set additional fields
	response.SdNumber = stockDisposalList.SdNumber
	response.TrCode = stockDisposalList.TrCode
	response.WhID = stockDisposalList.WhID
	response.StockType = stockDisposalList.StockType
	response.SubTotal = roundPrice(stockDisposalList.SubTotal)
	response.VatValue = roundPrice(stockDisposalList.VatValue)
	response.Total = roundPrice(stockDisposalList.SubTotal + stockDisposalList.VatValue)

	return response
}

// validateProducts validates products: qty, existence, and stock availability
func (service *stockDisposalServiceImpl) validateProducts(request entity.CreateStockDisposalBody, effectiveCustID string, c context.Context) (map[int64]*model.Product, []string) {
	var validationErrors []string
	productMap := make(map[int64]*model.Product)

	for i, product := range request.Products {
		productIndex := fmt.Sprintf("product index %d (pro_id: %d)", i, product.ProID)

		// Validate quantity
		if product.Qty1 == 0 && product.Qty2 == 0 && product.Qty3 == 0 {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s: please input qty (qty1, qty2, or qty3 must be greater than 0)", productIndex))
			continue
		}

		// Validate product existence
		productModel, err := service.StockDisposalRepository.FindProductByID(c, product.ProID, request.CustID, request.ParentCustID)
		if err != nil {
			validationErrors = append(validationErrors, productValidationError(productIndex, err))
			continue
		}
		productMap[product.ProID] = productModel

		// Validate quantity calculation
		qtyUnit := buildQtyUnit(product, productModel)
		totalQty, err := qtyUnit.ToTotalQuantity()
		if err != nil {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s: failed to calculate total quantity: %v", productIndex, err))
			continue
		}

		// Validate stock availability
		availableQty, err := service.StockDisposalRepository.GetAvailableStock(c, effectiveCustID, request.WhID, product.ProID)
		if err != nil {
			availableQty = 0
		}

		if float64(totalQty) > availableQty {
			validationErrors = append(validationErrors,
				fmt.Sprintf("%s: insufficient stock - available %.2f, requested %.2f",
					productIndex, availableQty, float64(totalQty)))
		}
	}

	return productMap, validationErrors
}

func (service *stockDisposalServiceImpl) Store(request entity.CreateStockDisposalBody) (response entity.StockDisposalResponse, err error) {
	c := context.Background()

	// Validate minimal 1 product required
	if len(request.Products) == 0 {
		return response, errors.New("products is required, minimum 1 product")
	}

	// Collect all validation errors
	var validationErrors []string

	// Validate file size
	fileErrors := validateFileSize(request.Products)
	validationErrors = append(validationErrors, fileErrors...)

	// Parse and validate date
	disposalDateParsed, err := validateAndParseDate(request.Date)
	if err != nil {
		return response, err
	}

	// Get warehouse to get stock_type
	warehouse, err := service.StockDisposalRepository.FindWarehouseByID(c, request.WhID, request.CustID)
	if err != nil {
		if isRecordNotFoundError(err) {
			return response, errors.New("warehouse not found or inactive")
		}

		return response, fmt.Errorf("failed to validate warehouse: %w", err)
	}

	// Validate supplier
	_, err = service.StockDisposalRepository.FindSupplierByID(c, request.SupID, request.ParentCustID)
	if err != nil {
		if isRecordNotFoundError(err) {
			return response, errors.New("supplier not found or inactive")
		}

		return response, fmt.Errorf("failed to validate supplier: %w", err)
	}

	// Validate products: qty, existence, and stock availability
	productMap, productErrors := service.validateProducts(request, warehouse.CustID, c)
	validationErrors = append(validationErrors, productErrors...)

	// Return all validation errors at once
	if len(validationErrors) > 0 {
		return response, fmt.Errorf("validation failed (%d error(s)):\n%s",
			len(validationErrors), strings.Join(validationErrors, "\n"))
	}

	totalSubTotal, totalVatValue := calculateTotals(request.Products)

	// Create stock disposal model
	stockDisposalModel := buildStockDisposalModel(request, warehouse, disposalDateParsed, totalSubTotal, totalVatValue)

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// Save header (SD number will be auto-generated via BeforeCreate hook)
		err := service.StockDisposalRepository.Store(txCtx, &stockDisposalModel)
		if err != nil {
			return err
		}

		// Create details and get stock update entities
		stockUpdateEntities, err := service.createStockDisposalDetails(request, stockDisposalModel, productMap, disposalDateParsed, txCtx)
		if err != nil {
			return err
		}

		// Update stock with batch update
		if len(stockUpdateEntities) > 0 {
			log.Printf("StockDisposalService: Updating stock with %d entities for SD %s", len(stockUpdateEntities), stockDisposalModel.SdNumber)
			err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			if err != nil {
				log.Printf("StockDisposalService: Stock update failed for SD %s: %v", stockDisposalModel.SdNumber, err)
				return fmt.Errorf("failed to update stock: %w", err)
			}
			log.Printf("StockDisposalService: Stock update successful for SD %s", stockDisposalModel.SdNumber)
		} else {
			log.Printf("StockDisposalService: No stock updates to process for SD %s", stockDisposalModel.SdNumber)
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	// Get created stock disposal with joins for response
	stockDisposalList, err := service.StockDisposalRepository.FindByNumber(c, stockDisposalModel.SdNumber, stockDisposalModel.CustID, request.ParentCustID, warehouse.CustID)
	if err != nil {
		return response, err
	}

	// Build response
	response = buildStockDisposalResponse(stockDisposalList)

	return response, nil
}

func (service *stockDisposalServiceImpl) Detail(sdID int64, custID, parentCustId string) (response entity.StockDisposalDetailResponse, err error) {
	c := context.Background()

	stockDisposal, err := service.StockDisposalRepository.FindByID(c, sdID, custID)
	if err != nil {
		return response, err
	}

	stockDisposalList, err := service.StockDisposalRepository.FindByNumber(c, stockDisposal.SdNumber, custID, parentCustId, stockDisposal.CustID)
	if err != nil {
		return response, err
	}

	// Get all details with joins
	details, err := service.StockDisposalRepository.FindDetail(c, sdID, custID)
	if err != nil {
		return response, err
	}

	response.SdID = stockDisposalList.SdID
	response.SdNumber = stockDisposalList.SdNumber
	if stockDisposalList.DisposalDate != nil {
		response.Date = stockDisposalList.DisposalDate.Format(constant.DATE_FORMAT_DETAIL)
	}
	response.SupID = stockDisposalList.SupID
	response.WhID = stockDisposalList.WhID
	response.StockType = stockDisposalList.StockType
	response.GrNo = stockDisposalList.GrNo
	response.Note = stockDisposalList.Note

	mapSupplierWarehouseInfoToDetail(stockDisposalList, &response)

	// Map details to response
	var totalSubTotal float64
	var totalVatValue float64
	for _, detail := range details {
		var productResponse entity.StockDisposalProductResponse
		productResponse.ProID = detail.ProID
		productResponse.ProCode = detail.ProCode
		productResponse.ProName = detail.ProName
		productResponse.UnitID1 = detail.UnitID1
		productResponse.UnitID2 = detail.UnitID2
		productResponse.UnitID3 = detail.UnitID3
		productResponse.Qty1 = detail.Qty1
		productResponse.Qty2 = detail.Qty2
		productResponse.Qty3 = detail.Qty3
		productResponse.PurchPrice1 = roundPrice(detail.PurchPrice1)
		productResponse.PurchPrice2 = roundPrice(detail.PurchPrice2)
		productResponse.PurchPrice3 = roundPrice(detail.PurchPrice3)
		productResponse.GrossPrice = roundPrice(detail.GrossPrice)
		productResponse.Vat = detail.Vat
		productResponse.VatValue = roundPrice(detail.VatValue)
		productResponse.SubTotal = roundPrice(detail.GrossPrice)

		// Set upload_file only if file fields are present
		// Prefer FileUrl over FileBase64 for new data
		if detail.FileName != nil && detail.FileType != nil && detail.MediaCategory != nil && detail.FileSize != nil {
			// Only set UploadFile if we have FileUrl (new data) or FileBase64 (old data for backward compatibility)
			hasFileUrl := detail.FileUrl != nil && *detail.FileUrl != ""
			hasFileBase64 := detail.FileBase64 != nil && *detail.FileBase64 != ""

			if hasFileUrl || hasFileBase64 {
				fileResponse := entity.StockDisposalFileResponse{
					FileName:      *detail.FileName,
					FileType:      *detail.FileType,
					MediaCategory: *detail.MediaCategory,
					FileSize:      formatFileSize(*detail.FileSize),
				}

				// Use FileUrl if available, otherwise leave empty (old data with FileBase64 won't show URL)
				if hasFileUrl {
					fileResponse.FileUrl = *detail.FileUrl
				}

				productResponse.UploadFile = &fileResponse
			}
		}

		response.DataProducts = append(response.DataProducts, productResponse)

		totalSubTotal += roundPrice(detail.GrossPrice)
		totalVatValue += roundPrice(detail.VatValue)
	}

	response.Subtotal = roundPrice(totalSubTotal)
	response.Vat = roundPrice(totalVatValue)
	response.Total = roundPrice(totalSubTotal + totalVatValue)

	return response, nil
}

func (service *stockDisposalServiceImpl) List(dataFilter entity.StockDisposalQueryFilter, custId, parentCustId string) (data []entity.StockDisposalListResponse, total int64, lastPage int, err error) {
	c := context.Background()

	// Set cust_id and parent_cust_id in filter
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId

	// Query with filters
	stockDisposals, total, lastPage, err := service.StockDisposalRepository.FindAllByCustId(c, dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	// Map to response
	for _, stockDisposal := range stockDisposals {
		var listResponse entity.StockDisposalListResponse
		listResponse.SdID = stockDisposal.SdID
		if stockDisposal.DisposalDate != nil {
			listResponse.Date = stockDisposal.DisposalDate.Format(constant.DATE_FORMAT_DETAIL)
		}
		listResponse.SdNumber = stockDisposal.SdNumber
		listResponse.WhID = stockDisposal.WhID
		listResponse.SupID = stockDisposal.SupID
		mapSupplierWarehouseInfoToList(&stockDisposal, &listResponse)
		calculatedSubtotal := stockDisposal.CalculatedSubtotal
		listResponse.Subtotal = roundPrice(calculatedSubtotal)
		listResponse.Vat = roundPrice(stockDisposal.VatValue)
		listResponse.VatValue = roundPrice(stockDisposal.VatValue)
		listResponse.Total = roundPrice(calculatedSubtotal + stockDisposal.VatValue)

		data = append(data, listResponse)
	}

	return data, total, lastPage, nil
}

func (service *stockDisposalServiceImpl) ProductLookup(dataFilter entity.StockDisposalProductLookupQueryFilter, custId, parentCustId string) (data []entity.StockDisposalProductLookupResponse, total int64, lastPage int, err error) {
	c := context.Background()

	// Query products dengan stock breakdown
	products, total, lastPage, err := service.StockDisposalRepository.FindProductsForLookup(c, dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	// Map to response
	for _, product := range products {
		var response entity.StockDisposalProductLookupResponse
		response.ProID = product.ProID
		response.ProCode = product.ProCode
		response.ProName = product.ProName
		response.Vat = int(product.Vat) // Convert float to int
		response.ConvUnit2 = product.ConvUnit2
		response.ConvUnit3 = product.ConvUnit3
		response.UnitID1 = product.UnitID1
		response.UnitID2 = product.UnitID2
		response.UnitID3 = product.UnitID3
		response.PurchPrice1 = product.PurchPrice1
		response.PurchPrice2 = product.PurchPrice2
		response.PurchPrice3 = product.PurchPrice3
		response.MinStockQty = product.MinStockQty
		response.SafStockQty = product.SafStockQty
		response.Qty1 = product.Qty1
		response.Qty2 = product.Qty2
		response.Qty3 = product.Qty3
		response.TotalQty = product.TotalQty
		response.InTransitStock1 = product.InTransitStock1
		response.InTransitStock2 = product.InTransitStock2
		response.InTransitStock3 = product.InTransitStock3

		data = append(data, response)
	}

	return data, total, lastPage, nil
}
