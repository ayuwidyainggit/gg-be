package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"master/entity"
	"master/model"
	"master/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

type AutomaticReplenishmentProductService interface {
	List(filter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]entity.AutomaticReplenishmentProductResponse, int, int, error)
	Export(filter entity.AutomaticReplenishmentProductQueryFilter, custId string) (*bytes.Buffer, string, string, error)
	DownloadTemplate(format string) (*bytes.Buffer, string, string, error)
	Import(req entity.AutomaticReplenishmentProductImportRequest, custId string, createdBy int64) (entity.AutomaticReplenishmentProductImportResponse, error)
	Detail(params entity.DetailAutomaticReplenishmentProductParams) (entity.AutomaticReplenishmentProductDetailResponse, error)
	Create(request []*entity.CreateAutomaticReplenishmentProductRequest, custId string, createdBy int64) ([]*entity.AutomaticReplenishmentProductResponse, error)
	Update(id int64, request entity.UpdateAutomaticReplenishmentProductRequest, custId string, updatedBy int64) error
	Delete(custId string, id int64, deletedBy int64) error
}

type automaticReplenishmentProductService struct {
	repo                  repository.AutomaticReplenishmentProductRepository
	productRepository     repository.ProductRepository
	distributorRepository repository.DistributorRepository
}

func NewAutomaticReplenishmentProductService(
	repo repository.AutomaticReplenishmentProductRepository,
	productRepository repository.ProductRepository,
	distributorRepository repository.DistributorRepository,
) AutomaticReplenishmentProductService {
	return &automaticReplenishmentProductService{
		repo:                  repo,
		productRepository:     productRepository,
		distributorRepository: distributorRepository,
	}
}

func (s *automaticReplenishmentProductService) List(filter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]entity.AutomaticReplenishmentProductResponse, int, int, error) {
	products, total, lastPage, err := s.repo.FindAll(filter, custId)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, List error:", err.Error())
		return nil, 0, 0, err
	}

	return mapAutomaticReplenishmentProductResponses(products), total, lastPage, nil
}

func (s *automaticReplenishmentProductService) Export(filter entity.AutomaticReplenishmentProductQueryFilter, custId string) (*bytes.Buffer, string, string, error) {
	products, err := s.repo.FindAllExport(filter, custId)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Export error:", err.Error())
		return nil, "", "", err
	}

	responses := mapAutomaticReplenishmentProductResponses(products)
	switch strings.ToLower(strings.TrimSpace(filter.Format)) {
	case "csv":
		buffer, err := createAutomaticReplenishmentProductExportCSV(responses)
		if err != nil {
			log.Error("AutomaticReplenishmentProductService, Export create csv error:", err.Error())
			return nil, "", "", err
		}
		return buffer, "text/csv", "automatic_replenishment_products.csv", nil
	case "xls":
		buffer, err := createAutomaticReplenishmentProductExportWorkbook(responses)
		if err != nil {
			log.Error("AutomaticReplenishmentProductService, Export create xls error:", err.Error())
			return nil, "", "", err
		}
		return buffer, "application/vnd.ms-excel", "automatic_replenishment_products.xls", nil
	default:
		buffer, err := createAutomaticReplenishmentProductExportWorkbook(responses)
		if err != nil {
			log.Error("AutomaticReplenishmentProductService, Export create xlsx error:", err.Error())
			return nil, "", "", err
		}
		return buffer, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "automatic_replenishment_products.xlsx", nil
	}
}

func (s *automaticReplenishmentProductService) DownloadTemplate(format string) (*bytes.Buffer, string, string, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "csv":
		buffer, err := createAutomaticReplenishmentProductTemplateCSV()
		if err != nil {
			return nil, "", "", err
		}
		return buffer, "text/csv", "automatic_replenishment_product_template.csv", nil
	case "xls":
		buffer, err := createAutomaticReplenishmentProductTemplateWorkbook()
		if err != nil {
			return nil, "", "", err
		}
		return buffer, "application/vnd.ms-excel", "automatic_replenishment_product_template.xls", nil
	default:
		buffer, err := createAutomaticReplenishmentProductTemplateWorkbook()
		if err != nil {
			return nil, "", "", err
		}
		return buffer, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "automatic_replenishment_product_template.xlsx", nil
	}
}

func (s *automaticReplenishmentProductService) Import(req entity.AutomaticReplenishmentProductImportRequest, custId string, createdBy int64) (entity.AutomaticReplenishmentProductImportResponse, error) {
	response := entity.AutomaticReplenishmentProductImportResponse{
		FileURL:     req.FileURL,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	rows, filename, err := downloadAutomaticReplenishmentImportRows(req.FileURL)
	if err != nil {
		return response, err
	}
	response.FileName = filename
	if len(rows) <= 1 {
		return response, fmt.Errorf("template does not contain data rows")
	}
	if err := validateAutomaticReplenishmentImportHeader(rows[0]); err != nil {
		return response, err
	}

	type importRow struct {
		model   model.AutomaticReplenishmentProduct
		rowText int
	}

	// Validate the whole file first so we can return row-level feedback
	// without persisting a partial import.
	pending := make([]importRow, 0, len(rows)-1)
	failedReasons := make([]string, 0)
	seenCombinations := make(map[string]int)
	now := time.Now()

	for rowIndex, row := range rows[1:] {
		actualRow := rowIndex + 2
		if isAutomaticReplenishmentImportRowEmpty(row) {
			continue
		}

		response.TotalRow++

		distributorCode := automaticReplenishmentImportCell(row, 0)
		productCode := automaticReplenishmentImportCell(row, 1)
		limitAction, err := parseAutomaticReplenishmentLimitAction(automaticReplenishmentImportCell(row, 2))
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		maxOrderQty, err := parseAutomaticReplenishmentPositiveInt(automaticReplenishmentImportCell(row, 3), limitAction == "RESTRICTED", "maximum order quantity")
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		minStockQty, err := parseAutomaticReplenishmentPositiveInt(automaticReplenishmentImportCell(row, 4), false, "minimum stock")
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		safetyStockQty, err := parseAutomaticReplenishmentPositiveInt(automaticReplenishmentImportCell(row, 5), true, "safety stock")
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		minOrderQty, err := parseAutomaticReplenishmentPositiveInt(automaticReplenishmentImportCell(row, 6), true, "minimum order quantity")
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		if distributorCode == "" {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: distributor code is required", actualRow))
			continue
		}
		if productCode == "" {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: product code is required", actualRow))
			continue
		}

		distributor, err := s.distributorRepository.FindOneByParentCustIdAndDistributorCode(custId, distributorCode)
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: distributor %s not found", actualRow, distributorCode))
			continue
		}

		product, err := s.productRepository.FindOneByProductCodeAndCustId(productCode, distributor.DistributorCustID)
		if err != nil || product.ProductId == 0 {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: product %s not found for distributor %s", actualRow, productCode, distributorCode))
			continue
		}

		// Track duplicates inside the same import file separately from duplicates
		// that already exist in the database.
		combinationKey := automaticReplenishmentProductCombinationKey(int64(product.ProductId), distributor.DistributorId)
		if existingRow, exists := seenCombinations[combinationKey]; exists {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: setup already exists in import file for product %s and distributor %s (first declared at row %d)", actualRow, productCode, distributorCode, existingRow))
			continue
		}

		exists, err := s.repo.IsExistsByProductAndDistributor(custId, int64(product.ProductId), distributor.DistributorId)
		if err != nil {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: failed to check existing setup: %v", actualRow, err))
			continue
		}
		if exists {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: setup already exists for product %s and distributor %s", actualRow, productCode, distributorCode))
			continue
		}
		seenCombinations[combinationKey] = actualRow

		pending = append(pending, importRow{
			rowText: actualRow,
			model: model.AutomaticReplenishmentProduct{
				CustId:          custId,
				ProId:           int64(product.ProductId),
				DistributorId:   distributor.DistributorId,
				LimitAction:     limitAction,
				MaxOrderQty:     maxOrderQty,
				MaxOrderType:    "L",
				MinStockQty:     minStockQty,
				MinStockType:    "L",
				SafetyStockQty:  safetyStockQty,
				SafetyStockType: "L",
				MinOrderQty:     minOrderQty,
				MinOrderType:    "L",
				IsActive:        boolPtr(true),
				CreatedBy:       createdBy,
				CreatedAt:       now,
				UpdatedBy:       &createdBy,
				UpdatedAt:       &now,
			},
		})
	}

	response.FailedReasons = failedReasons
	response.FailedRow = len(failedReasons)
	if len(failedReasons) > 0 {
		return response, fmt.Errorf("import validation failed")
	}

	// Only write after all rows pass validation to keep the import atomic
	// from the API consumer's perspective.
	for _, item := range pending {
		if _, err := s.repo.Create(&item.model); err != nil {
			response.FailedReasons = append(response.FailedReasons, fmt.Sprintf("row %d: failed to save data: %v", item.rowText, err))
		}
	}

	response.FailedRow = len(response.FailedReasons)
	response.SuccessRow = len(pending) - response.FailedRow
	if response.FailedRow > 0 {
		return response, fmt.Errorf("import failed while saving data")
	}

	response.SuccessRow = len(pending)
	return response, nil
}

func (s *automaticReplenishmentProductService) Detail(params entity.DetailAutomaticReplenishmentProductParams) (entity.AutomaticReplenishmentProductDetailResponse, error) {
	product, err := s.repo.FindOne(params)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Detail error:", err.Error())
		return entity.AutomaticReplenishmentProductDetailResponse{}, err
	}

	updatedAt := product.UpdatedAt.Format("2006-01-02 15:04:05")
	response := entity.AutomaticReplenishmentProductDetailResponse{
		CustId:          product.CustId,
		Id:              product.Id,
		ProId:           product.ProId,
		DistributorId:   product.DistributorId,
		LimitAction:     product.LimitAction,
		MaxOrderQty:     product.MaxOrderQty,
		MaxOrderType:    product.MaxOrderType,
		MinStockQty:     product.MinStockQty,
		MinStockType:    product.MinStockType,
		SafetyStockQty:  product.SafetyStockQty,
		SafetyStockType: product.SafetyStockType,
		MinOrderQty:     product.MinOrderQty,
		MinOrderType:    product.MinOrderType,
		IsActive:        product.IsActive,
		CreatedBy:       product.CreatedBy,
		CreatedByName:   product.CreatedByName,
		CreatedAt:       product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:       product.UpdatedBy,
		UpdatedByName:   product.UpdatedByName,
		UpdatedAt:       &updatedAt,
	}

	if product.ProCode != "" {
		response.ProCode = product.ProCode
	}
	if product.ProName != "" {
		response.ProName = product.ProName
	}
	if product.DistributorCode != "" {
		response.DistributorCode = product.DistributorCode
	}
	if product.DistributorName != "" {
		response.DistributorName = product.DistributorName
	}

	if product.UpdatedBy != nil {
		response.UpdatedBy = product.UpdatedBy
		response.UpdatedAt = stringPtr(product.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
	if product.DeletedBy != nil {
		response.DeletedBy = product.DeletedBy
		response.DeletedAt = stringPtr(product.DeletedAt.Format("2006-01-02 15:04:05"))
	}
	if product.IsDel != nil {
		response.IsDel = product.IsDel
	}

	return response, nil
}

func (s *automaticReplenishmentProductService) Create(request []*entity.CreateAutomaticReplenishmentProductRequest, custId string, createdBy int64) (responses []*entity.AutomaticReplenishmentProductResponse, err error) {
	now := time.Now()
	seenCombinations := make(map[string]struct{}, len(request))
	for _, req := range request {
		// Reject duplicates both within the current payload and against the
		// existing customer setup.
		combinationKey := automaticReplenishmentProductCombinationKey(req.ProId, req.DistributorId)
		if _, exists := seenCombinations[combinationKey]; exists {
			return nil, fmt.Errorf("product_id: %d and distributor_id: %d combination is already exists in request", req.ProId, req.DistributorId)
		}

		var exists bool
		exists, err = s.repo.IsExistsByProductAndDistributor(custId, req.ProId, req.DistributorId)
		if err != nil {
			log.Error("AutomaticReplenishmentProductService, Create, IsExistsByProductAndDistributor error:", err.Error())
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("product_id: %d and distributor_id: %d combination is already exists", req.ProId, req.DistributorId)
		}

		product := model.AutomaticReplenishmentProduct{
			CustId:          custId,
			ProId:           req.ProId,
			DistributorId:   req.DistributorId,
			LimitAction:     req.LimitAction,
			MaxOrderQty:     req.MaxOrderQty,
			MaxOrderType:    req.MaxOrderType,
			MinStockQty:     req.MinStockQty,
			MinStockType:    req.MinStockType,
			SafetyStockQty:  req.SafetyStockQty,
			SafetyStockType: req.SafetyStockType,
			MinOrderQty:     req.MinOrderQty,
			MinOrderType:    req.MinOrderType,
			IsActive:        boolPtr(true),
			CreatedBy:       createdBy,
			CreatedAt:       now,
			UpdatedBy:       &createdBy,
			UpdatedAt:       &now,
		}

		_, err = s.repo.Create(&product)
		if err != nil {
			log.Error("AutomaticReplenishmentProductService, Create error:", err.Error())
			return
		}

		seenCombinations[combinationKey] = struct{}{}
	}

	return
}

func automaticReplenishmentProductCombinationKey(proId, distributorId int64) string {
	return strconv.FormatInt(proId, 10) + ":" + strconv.FormatInt(distributorId, 10)
}

func (s *automaticReplenishmentProductService) Update(id int64, request entity.UpdateAutomaticReplenishmentProductRequest, custId string, updatedBy int64) error {
	// Check if exists
	exists, err := s.repo.IsExists(id, custId)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Update, IsExists error:", err.Error())
		return err
	}
	if !exists {
		return fmt.Errorf("automatic replenishment product not found")
	}

	now := time.Now()
	product := model.AutomaticReplenishmentProduct{
		CustId:          custId,
		ProId:           request.ProId,
		DistributorId:   request.DistributorId,
		LimitAction:     request.LimitAction,
		MaxOrderQty:     request.MaxOrderQty,
		MaxOrderType:    request.MaxOrderType,
		MinStockQty:     request.MinStockQty,
		MinStockType:    request.MinStockType,
		SafetyStockQty:  request.SafetyStockQty,
		SafetyStockType: request.SafetyStockType,
		MinOrderQty:     request.MinOrderQty,
		MinOrderType:    request.MinOrderType,
		UpdatedBy:       &updatedBy,
		UpdatedAt:       &now,
	}

	err = s.repo.Update(id, &product)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Update error:", err.Error())
		return err
	}

	return nil
}

func (s *automaticReplenishmentProductService) Delete(custId string, id int64, deletedBy int64) error {
	// Check if exists
	exists, err := s.repo.IsExists(id, custId)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Delete, IsExists error:", err.Error())
		return err
	}
	if !exists {
		return fmt.Errorf("automatic replenishment product not found")
	}

	err = s.repo.Delete(custId, id, deletedBy)
	if err != nil {
		log.Error("AutomaticReplenishmentProductService, Delete error:", err.Error())
		return err
	}

	return nil
}

func mapAutomaticReplenishmentProductResponses(products []*model.AutomaticReplenishmentProduct) []entity.AutomaticReplenishmentProductResponse {
	responses := make([]entity.AutomaticReplenishmentProductResponse, 0, len(products))
	for _, product := range products {
		var updatedAt *string
		if product.UpdatedAt != nil {
			formatted := product.UpdatedAt.Format("2006-01-02 15:04:05")
			updatedAt = &formatted
		}

		response := entity.AutomaticReplenishmentProductResponse{
			CustId:          product.CustId,
			Id:              product.Id,
			ProId:           product.ProId,
			ProCode:         product.ProCode,
			ProName:         product.ProName,
			DistributorId:   product.DistributorId,
			DistributorCode: product.DistributorCode,
			DistributorName: product.DistributorName,
			LimitAction:     product.LimitAction,
			MaxOrderQty:     product.MaxOrderQty,
			MaxOrderType:    product.MaxOrderType,
			MinStockQty:     product.MinStockQty,
			MinStockType:    product.MinStockType,
			SafetyStockQty:  product.SafetyStockQty,
			SafetyStockType: product.SafetyStockType,
			MinOrderQty:     product.MinOrderQty,
			MinOrderType:    product.MinOrderType,
			CreatedBy:       product.CreatedBy,
			CreatedByName:   product.CreatedByName,
			CreatedAt:       product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedBy:       product.UpdatedBy,
			UpdatedByName:   product.UpdatedByName,
			UpdatedAt:       updatedAt,
		}

		responses = append(responses, response)
	}

	return responses
}

func automaticReplenishmentExportHeaders() []string {
	return []string{
		"cust_id",
		"id",
		"pro_id",
		"pro_code",
		"pro_name",
		"distributor_id",
		"distributor_code",
		"distributor_name",
		"limit_action",
		"max_order_qty",
		"max_order_type",
		"min_stock_qty",
		"min_stock_type",
		"safety_stock_qty",
		"safety_stock_type",
		"min_order_qty",
		"min_order_type",
		"created_by",
		"created_by_name",
		"created_at",
		"updated_by",
		"updated_by_name",
		"updated_at",
	}
}

func automaticReplenishmentTemplateHeaders() []string {
	return []string{
		"Distributor Code",
		"Product Code",
		"Limit Action (Maximum Order Quantity)",
		"Maximum Order Quantity",
		"Minimum Stock (Largest Unit)",
		"Safety Stock (Largest Unit)",
		"Minimum order Quantity (Largest Unit)",
	}
}

func automaticReplenishmentTemplateInstructions() [][]string {
	return [][]string{
		{"Column Name (Template)", "Instruction"},
		{"Distributor Code", "Required. Use distributor code registered under the current customer."},
		{"Product Code", "Required. Product must exist for the selected distributor."},
		{"Limit Action (Maximum Order Quantity)", "Required. Accepted values: 0/1/2 or RESTRICTED/UNRESTRICTED/WARNING."},
		{"Maximum Order Quantity", "Required when limit action is RESTRICTED. Must be an integer >= 1."},
		{"Minimum Stock (Largest Unit)", "Optional. Must be an integer >= 1 when filled."},
		{"Safety Stock (Largest Unit)", "Required. Must be an integer >= 1."},
		{"Minimum order Quantity (Largest Unit)", "Required. Must be an integer >= 1."},
	}
}

func createAutomaticReplenishmentProductExportWorkbook(products []entity.AutomaticReplenishmentProductResponse) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Replenishment Products"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := automaticReplenishmentExportHeaders()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for rowIndex, product := range products {
		row := []interface{}{
			product.CustId,
			product.Id,
			product.ProId,
			product.ProCode,
			product.ProName,
			product.DistributorId,
			product.DistributorCode,
			product.DistributorName,
			product.LimitAction,
			product.MaxOrderQty,
			product.MaxOrderType,
			product.MinStockQty,
			product.MinStockType,
			product.SafetyStockQty,
			product.SafetyStockType,
			product.MinOrderQty,
			product.MinOrderType,
			product.CreatedBy,
			product.CreatedByName,
			product.CreatedAt,
			product.UpdatedBy,
			product.UpdatedByName,
			product.UpdatedAt,
		}

		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Use a fixed width for each export column so the sheet is readable
	// immediately after download.
	for i := range headers {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			continue
		}
		_ = f.SetColWidth(sheetName, colName, colName, 18)
	}

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	return f.WriteToBuffer()
}

func createAutomaticReplenishmentProductExportCSV(products []entity.AutomaticReplenishmentProductResponse) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	if err := writer.Write(automaticReplenishmentExportHeaders()); err != nil {
		return nil, err
	}

	for _, product := range products {
		record := []string{
			product.CustId,
			strconv.FormatInt(product.Id, 10),
			strconv.FormatInt(product.ProId, 10),
			product.ProCode,
			product.ProName,
			strconv.FormatInt(product.DistributorId, 10),
			product.DistributorCode,
			product.DistributorName,
			product.LimitAction,
			strconv.Itoa(product.MaxOrderQty),
			product.MaxOrderType,
			strconv.Itoa(product.MinStockQty),
			product.MinStockType,
			strconv.Itoa(product.SafetyStockQty),
			product.SafetyStockType,
			strconv.Itoa(product.MinOrderQty),
			product.MinOrderType,
			strconv.FormatInt(product.CreatedBy, 10),
			product.CreatedByName,
			product.CreatedAt,
			formatAutomaticReplenishmentNullableInt64(product.UpdatedBy),
			formatAutomaticReplenishmentNullableString(product.UpdatedByName),
			formatAutomaticReplenishmentNullableString(product.UpdatedAt),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer, nil
}

func createAutomaticReplenishmentProductTemplateWorkbook() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Template"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := automaticReplenishmentTemplateHeaders()
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Keep usage guidance in a separate sheet so the template stays clean for data entry.
	instructionSheet := "Instructions"
	if _, err := f.NewSheet(instructionSheet); err != nil {
		return nil, err
	}
	for rowIndex, row := range automaticReplenishmentTemplateInstructions() {
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			f.SetCellValue(instructionSheet, cell, value)
			if rowIndex == 0 {
				f.SetCellStyle(instructionSheet, cell, cell, headerStyle)
			}
		}
	}

	for i := range headers {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			continue
		}
		_ = f.SetColWidth(sheetName, colName, colName, 28)
	}
	_ = f.SetColWidth(instructionSheet, "A", "A", 36)
	_ = f.SetColWidth(instructionSheet, "B", "B", 80)

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	return f.WriteToBuffer()
}

func createAutomaticReplenishmentProductTemplateCSV() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	if err := writer.Write(automaticReplenishmentTemplateHeaders()); err != nil {
		return nil, err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer, nil
}

func formatAutomaticReplenishmentNullableInt64(value *int64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatInt(*value, 10)
}

func formatAutomaticReplenishmentNullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func downloadAutomaticReplenishmentImportRows(fileURL string) ([][]string, string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	filename := fileURL
	if parts := strings.Split(fileURL, "/"); len(parts) > 0 {
		filename = parts[len(parts)-1]
	}

	if strings.EqualFold(strings.TrimSpace(filename[strings.LastIndex(filename, ".")+1:]), "csv") {
		reader := csv.NewReader(bytes.NewReader(content))
		rows, err := reader.ReadAll()
		if err != nil {
			return nil, "", err
		}
		return rows, filename, nil
	}

	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		// Some uploaded files are CSVs with the wrong extension, so fall back
		// to CSV parsing before returning the Excel error.
		reader := csv.NewReader(bytes.NewReader(content))
		rows, csvErr := reader.ReadAll()
		if csvErr == nil {
			return rows, filename, nil
		}
		return nil, "", err
	}
	defer func() { _ = f.Close() }()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, "", fmt.Errorf("excel file has no sheet")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, "", err
	}

	return rows, filename, nil
}

func validateAutomaticReplenishmentImportHeader(header []string) error {
	expected := []string{
		"Distributor Code",
		"Product Code",
		"Limit Action (Maximum Order Quantity)",
		"Maximum Order Quantity",
		"Minimum Stock (Largest Unit)",
		"Safety Stock (Largest Unit)",
		"Minimum order Quantity (Largest Unit)",
	}

	for i, want := range expected {
		if automaticReplenishmentImportCell(header, i) != want {
			return fmt.Errorf("invalid template header at column %d: expected %q", i+1, want)
		}
	}

	return nil
}

func automaticReplenishmentImportCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func isAutomaticReplenishmentImportRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func parseAutomaticReplenishmentLimitAction(value string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	switch normalized {
	case "0", "RESTRICTED":
		return "RESTRICTED", nil
	case "1", "UNRESTRICTED":
		return "UNRESTRICTED", nil
	case "2", "WARNING":
		return "WARNING", nil
	default:
		return "", fmt.Errorf("limit action must be 0/1/2 or RESTRICTED/UNRESTRICTED/WARNING")
	}
}

func parseAutomaticReplenishmentPositiveInt(value string, required bool, field string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		if required {
			return 0, fmt.Errorf("%s is required", field)
		}
		return 0, nil
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", field)
	}
	if number < 1 {
		return 0, fmt.Errorf("%s must be greater than or equal to 1", field)
	}

	return number, nil
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}
