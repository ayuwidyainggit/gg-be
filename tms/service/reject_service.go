package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"regexp"
	"scyllax-tms/config"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"
	"gorm.io/gorm"
)

type RejectService interface {
	GetReject(ctx context.Context, dataFilter entity.RejectQueryFilter) (response []entity.RejectResponse)
	GetRejectPartial(ctx context.Context, dataFilter entity.RejectQueryFilter) (response []entity.RejectPartialResponse)
	RejectAll(ctx context.Context, headers map[string]string, request entity.RejectRequest)
	RejectPartial(ctx context.Context, headers map[string]string, request entity.RejectPartialRequest)
	RejectCancel(ctx context.Context, request entity.RejectCancelRequest)
	Conversion(ctx context.Context, headers map[string]string, request entity.ConversionRequest) (entity.ConversionResponse, error)
	UpdateQtyFinalOrder(ctx context.Context, headers map[string]string, request entity.UpdateQtyFinalOrderRequest, orderNo string) error
}

type RejectServiceImpl struct {
	shipmentInvoicesRepo    repository.ShipmentInvoicesRepo
	shipmentService         ShipmentService
	shipmentOrderStatusRepo repository.ShipmentOrderStatusRepo
	validate                *validator.Validate
}

func NewRejectServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, shipmentService ShipmentService, shipmentOrderStatusRepo repository.ShipmentOrderStatusRepo, validate *validator.Validate) RejectService {
	return &RejectServiceImpl{
		shipmentInvoicesRepo:    shipmentInvoicesRepo,
		shipmentService:         shipmentService,
		shipmentOrderStatusRepo: shipmentOrderStatusRepo,
		validate:                validate,
	}
}

func (service *RejectServiceImpl) GetReject(ctx context.Context, dataFilter entity.RejectQueryFilter) (response []entity.RejectResponse) {
	result := service.shipmentInvoicesRepo.GetAllReject(ctx, dataFilter)

	re := regexp.MustCompile(`(?i)^(Reject|Reject All)$`)

	// for _, row := range result {
	// 	var res entity.RejectResponse
	// 	helper.Automapper(row, &res)
	// 	response = append(response, res)
	// }

	for _, row := range result {
		if re.MatchString(row.ProductStatus) {
			var res entity.RejectResponse
			helper.Automapper(row, &res)
			res.CtgId1 = "SMALL" // unit_id1 := small

			if row.UnitId2 == row.UnitId3 {
				// unit_id2 == unit_id3 := middle
				res.CtgId2 = "MIDDLE"
				res.CtgId3 = "MIDDLE"
			} else {
				// unit_id3 != unit_id2 -> unit_id3 := large
				res.CtgId2 = "MIDDLE"
				res.CtgId3 = "LARGE"
			}
			response = append(response, res)
		}
	}

	return response
}

func (service *RejectServiceImpl) GetRejectPartial(ctx context.Context, dataFilter entity.RejectQueryFilter) (response []entity.RejectPartialResponse) {
	result := service.shipmentInvoicesRepo.GetAllReject(ctx, dataFilter)

	productMap := make(map[int]*entity.RejectPartialResponse)

	for _, row := range result {
		if row.ReasonID != nil {
			reasonID := *row.ReasonID

			// Check if ReasonID already exists in productMap
			if _, exists := productMap[reasonID]; !exists {
				var res entity.RejectPartialResponse
				helper.Automapper(row, &res)
				res.Products = nil
				// Set CtgId1
				// res.CtgId1 = "SMALL" // unit_id1 := small

				// // Custom logic for CtgId2 and CtgId3
				// if row.UnitId2 == row.UnitId3 {
				// 	// unit_id2 == unit_id3 := middle
				// 	res.CtgId2 = "MIDDLE"
				// 	res.CtgId3 = "MIDDLE"
				// } else {
				// 	// unit_id3 != unit_id2 -> unit_id3 := large
				// 	res.CtgId2 = "MIDDLE"
				// 	res.CtgId3 = "LARGE"
				// }
				productMap[reasonID] = &res
			}

			if row.ProductID != 0 {
				var product entity.ProductMap
				helper.Automapper(row, &product)

				// product.Stock1 = helper.CalculateQty(row.Qty1, row.QtyReject1)
				// product.Stock2 = helper.CalculateQty(row.Qty2, row.QtyReject2)
				// product.Stock3 = helper.CalculateQty(row.Qty3, row.QtyReject3)
				product.Stock1 = *row.QtyReject1
				product.Stock2 = *row.QtyReject2
				product.Stock3 = *row.QtyReject3

				// Set CtgId fields based on UnitID comparison
				product.CtgId1 = "SMALL" // Default value
				if row.UnitId2 == row.UnitId3 {
					product.CtgId2 = "MIDDLE"
					product.CtgId3 = "MIDDLE"
				} else {
					product.CtgId2 = "MIDDLE"
					product.CtgId3 = "LARGE"
				}

				// Append the product to Products slice of corresponding RejectPartialResponse
				productMap[reasonID].Products = append(productMap[reasonID].Products, product)
			}
		}
	}

	for _, product := range productMap {
		response = append(response, *product)
	}
	return response
}

func (service *RejectServiceImpl) RejectAll(ctx context.Context, headers map[string]string, request entity.RejectRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ProductStatus: "Reject",
		ReasonID:      &request.ReasonID,
		ReasonName:    &request.ReasonName,
		UnloadAt:      &request.CurrentTime,
	}
	err = service.shipmentInvoicesRepo.UpdateReject(ctx, request.ID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	result := service.shipmentInvoicesRepo.GetAllOrderNo(ctx)
	for _, v := range result {
		var status string

		switch v.ProductStatus {
		case "-":
			status = "On Delivery"
		case "Receive":
			status = "Received"
		case "Reject Partial":
			status = "Partial Received"
		case "Reject All":
			status = "Cancelled"
		default:
			continue
		}

		if err := service.shipmentOrderStatusRepo.CreateOrUpdate(ctx, model.ShipmentOrderStatus{
			OrderNo:     v.OrderNo,
			StatusOrder: status,
		}); err != nil {
			panic(exception.NewInternalServerError(err.Error()))
		}
	}

	// Fetch all order numbers for the given shipment number
	orderNos := service.shipmentInvoicesRepo.FindAllOrderNoById(ctx, request.ID)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", orderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: orderNo,
						Status:  9,
					},
				},
			}
			log.Printf("Order Update: %+v", orderUpdate)
			err = service.shipmentService.MobileUpdateStatusOrder(ctx, headers, orderUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return
			}
		} else if strings.HasPrefix(orderNo, "SR") {
			log.Printf("Processing OrderNo with SR prefix: %s", orderNo)
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: orderNo,
						Status:  4,
					},
				},
			}
			log.Printf("Return Update: %+v", returnUpdate)
			err = service.shipmentService.MobileUpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating return status: %v", err)
				return
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", orderNo)
		}
	}
}

func (service *RejectServiceImpl) RejectCancel(ctx context.Context, request entity.RejectCancelRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ReasonID:   nil,
		ReasonName: nil,
	}

	err = service.shipmentInvoicesRepo.RejectCancel(ctx, request.ID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

}

func (service *RejectServiceImpl) RejectPartial(ctx context.Context, headers map[string]string, request entity.RejectPartialRequest) {
	// Validate request
	if err := service.validate.Struct(request); err != nil {
		helper.ErrorPanic(err)
	}

	// Begin transaction
	tx, err := service.shipmentInvoicesRepo.BeginTx(ctx)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
	defer service.handleTransactionRollbackOrCommit(tx, &err)

	// Process each request data
	for _, req := range request.Data {
		if err := service.validate.Struct(req); err != nil {
			helper.ErrorPanic(err)
		}

		// Process each product in the request
		for _, product := range req.Products {
			if err := service.processProduct(ctx, tx, headers, product, req); err != nil {
				tx.Rollback()
				log.Printf("Error processing product %d: %v", product.ID, err)
				panic(exception.NewInternalServerError(err.Error()))
			}
		}

		// Update shipment order status
		if err := service.updateShipmentOrderStatus(ctx, req); err != nil {
			tx.Rollback()
			log.Printf("Error updating shipment order status: %v", err)
			panic(exception.NewInternalServerError(err.Error()))
		}

		// Update order status for related orders
		if err := service.updateOrderStatusForRelatedOrders(ctx, headers, req); err != nil {
			tx.Rollback()
			log.Printf("Error updating order status for related orders: %v", err)
			panic(exception.NewInternalServerError(err.Error()))
		}
	}
}

// handleTransactionRollbackOrCommit handles transaction commit or rollback
func (service *RejectServiceImpl) handleTransactionRollbackOrCommit(tx *gorm.DB, err *error) {
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	} else if *err != nil {
		tx.Rollback()
	} else {
		*err = tx.Commit().Error
	}
}

// processProduct processes a single product in the reject request
func (service *RejectServiceImpl) processProduct(ctx context.Context, tx *gorm.DB, headers map[string]string, product entity.Product, req entity.RejectPartialBody) error {
	for i, qty := range product.Qty {
		dataset := model.ShipmentInvoices{
			ReasonID:      &req.ReasonID,
			ReasonName:    &req.ReasonName,
			OutletID:      req.OutletID,
			ProductStatus: "Reject Partial",
			ShipmentNo:    &req.ShipmentNo,
			UnloadAt:      &req.CurrentTime,
		}

		switch i {
		case 0:
			dataset.QtyReject1 = &qty.Stock
			dataset.UnitId1 = qty.UnitID
		case 1:
			dataset.QtyReject2 = &qty.Stock
			dataset.UnitId2 = qty.UnitID
		case 2:
			dataset.QtyReject3 = &qty.Stock
			dataset.UnitId3 = qty.UnitID
		}

		// Update shipment invoice with rejection data
		if err := service.shipmentInvoicesRepo.UpdateRejectPartial(ctx, product.ID, dataset, tx); err != nil {
			return fmt.Errorf("failed to update reject partial for product ID %d: %v", product.ID, err)
		}
	}

	// Handle stock and rejection quantities
	return service.handleStockAndRejectQuantities(ctx, headers, product, req)
}

// handleStockAndRejectQuantities calculates and updates final stock and reject quantities
func (service *RejectServiceImpl) handleStockAndRejectQuantities(ctx context.Context, headers map[string]string, product entity.Product, req entity.RejectPartialBody) error {
	// Fetch existing shipment invoice data
	data, err := service.shipmentInvoicesRepo.FindByOneColumn(ctx, "id", product.ID)
	if err != nil {
		log.Printf("Failed to find shipment invoice by ID %d: %v", product.ID, err)
		return err
	}

	// Convert stock and reject quantities
	stockQuantity := entity.ConversionRequest{
		ProId: int(data.ProductID),
		Qty1:  int(data.Qty1),
		Qty2:  int(data.Qty2),
		Qty3:  int(data.Qty3),
	}
	rejectQuantity := entity.ConversionRequest{
		ProId: int(data.ProductID),
		Qty1:  int(product.Qty[0].Stock),
		Qty2:  int(product.Qty[1].Stock),
		Qty3:  int(product.Qty[2].Stock),
	}

	// Get total stock and reject quantities
	totalStockQuantity, err := service.getTotalQuantity(ctx, headers, stockQuantity)
	if err != nil {
		return err
	}

	totalRejectQuantity, err := service.getTotalQuantity(ctx, headers, rejectQuantity)
	if err != nil {
		return err
	}

	// Calculate final quantity
	finalQuantity := totalStockQuantity - totalRejectQuantity
	if finalQuantity <= 0 {
		log.Printf("Total Reject must lower than stock: reject=%d, stock=%d", totalRejectQuantity, totalStockQuantity)
		return fmt.Errorf("Total Reject (%d) must lower than stock (%d)", totalRejectQuantity, totalStockQuantity)
	}

	// Final quantity conversion request
	finalQuantityRequest := entity.ConversionRequest{
		ProId: int(data.ProductID),
		Qty1:  finalQuantity,
		Qty2:  0,
		Qty3:  0,
	}

	// Get final quantity response
	finalQuantityResponse, err := service.Conversion(ctx, headers, finalQuantityRequest)
	if err != nil {
		log.Printf("Failed to find response for final quantity: %v", err)
		return err
	}

	// Update final quantity in order
	updateFinalQuantityRequest := entity.UpdateQtyFinalOrderRequest{
		DetailsFinal: struct {
			Normal []entity.ConversionRequest `json:"normal"`
		}{
			Normal: []entity.ConversionRequest{
				{
					OrderDetailID: IntPtr(data.OrderDetailID),
					ProId:         int(data.ProductID),
					Qty1:          int(finalQuantityResponse.Data.Qty1),
					Qty2:          int(finalQuantityResponse.Data.Qty2),
					Qty3:          int(finalQuantityResponse.Data.Qty3),
				},
			},
		},
	}

	// Update final order quantities
	return service.UpdateQtyFinalOrder(ctx, headers, updateFinalQuantityRequest, *data.OrderNo)
}

// getTotalQuantity converts quantity and returns total quantity
func (service *RejectServiceImpl) getTotalQuantity(ctx context.Context, headers map[string]string, quantity entity.ConversionRequest) (int, error) {
	totalQuantityResponse, err := service.Conversion(ctx, headers, quantity)
	if err != nil {
		log.Printf("Failed to find response: %v", err)
		return 0, err
	}
	return totalQuantityResponse.Data.TotalQty, nil
}

// updateShipmentOrderStatus updates the shipment order status based on the request
func (service *RejectServiceImpl) updateShipmentOrderStatus(ctx context.Context, req entity.RejectPartialBody) error {
	result := service.shipmentInvoicesRepo.GetAllOrderNo(ctx)
	for _, v := range result {
		status := service.getOrderStatus(v.ProductStatus)
		if status == "" {
			continue
		}
		if err := service.shipmentOrderStatusRepo.CreateOrUpdate(ctx, model.ShipmentOrderStatus{
			OrderNo:     v.OrderNo,
			StatusOrder: status,
		}); err != nil {
			return fmt.Errorf("failed to update shipment order status for OrderNo %s: %v", v.OrderNo, err)
		}
	}
	return nil
}

// getOrderStatus returns the status based on product status
func (service *RejectServiceImpl) getOrderStatus(productStatus string) string {
	switch productStatus {
	case "-":
		return "On Delivery"
	case "Receive":
		return "Received"
	case "Reject Partial":
		return "Partial Received"
	case "Reject All":
		return "Cancelled"
	default:
		return ""
	}
}

// updateOrderStatusForRelatedOrders updates the order status for related orders
func (service *RejectServiceImpl) updateOrderStatusForRelatedOrders(ctx context.Context, headers map[string]string, req entity.RejectPartialBody) error {
	// Fetch all order numbers for the given shipment number
	orderNos := service.shipmentInvoicesRepo.FindAllOrderNoByShipmentNo(ctx, req.ShipmentNo, req.OutletID)
	if err := service.processRelatedOrders(ctx, headers, orderNos); err != nil {
		return err
	}
	return nil
}

// processRelatedOrders processes each related order based on the order number
func (service *RejectServiceImpl) processRelatedOrders(ctx context.Context, headers map[string]string, orderNos []string) error {
	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SO") {
			err := service.processOrder(ctx, headers, orderNo)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(orderNo, "SR") {
			err := service.processReturn(ctx, headers, orderNo)
			if err != nil {
				return err
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", orderNo)
		}
	}
	return nil
}

// processOrder updates the order status for a given order number with "SO" prefix
func (service *RejectServiceImpl) processOrder(ctx context.Context, headers map[string]string, orderNo string) error {
	orderUpdate := entity.UpdateStatusOrder{
		Orders: []entity.OrderItem{
			{
				OrderNo: orderNo,
				Status:  5,
			},
		},
	}
	log.Printf("Order Update: %+v", orderUpdate)
	return service.shipmentService.MobileUpdateStatusOrder(ctx, headers, orderUpdate)
}

// processReturn updates the return status for a given order number with "SR" prefix
func (service *RejectServiceImpl) processReturn(ctx context.Context, headers map[string]string, orderNo string) error {
	returnUpdate := entity.UpdateStatusReturn{
		Returns: []entity.ReturnItem{
			{
				OrderNo: orderNo,
				Status:  4,
			},
		},
	}
	log.Printf("Return Update: %+v", returnUpdate)
	return service.shipmentService.MobileUpdateStatusReturn(ctx, headers, returnUpdate)
}


func (service *RejectServiceImpl) Conversion(ctx context.Context, headers map[string]string, request entity.ConversionRequest) (entity.ConversionResponse, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		return entity.ConversionResponse{}, exception.NewInternalServerError(err.Error())
	}

	endpointUrl := fmt.Sprintf("%s/v1/orders/conversion", config.KongUrlMobile)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return entity.ConversionResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return entity.ConversionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return entity.ConversionResponse{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return entity.ConversionResponse{}, fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	// Unmarshal the response body into ConversionResponse
	var conversionResponse entity.ConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&conversionResponse); err != nil {
		return entity.ConversionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return conversionResponse, nil
}

func (service *RejectServiceImpl) UpdateQtyFinalOrder(ctx context.Context, headers map[string]string, request entity.UpdateQtyFinalOrderRequest, orderNo string) error {
	config, err := config.LoadConfig(".")
	if err != nil {
		return exception.NewInternalServerError(err.Error())
	}

	endpointUrl := fmt.Sprintf("%s/v1/orders/final/%s", config.KongUrlMobile, orderNo)
	log.Printf("endpoint: %s", endpointUrl)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	log.Printf("req: %s", req)

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	return nil
}

func IntPtr(i int) *int {
    return &i
}

func convertPointerToInt(ptr *int64) int {
	if ptr == nil {
		return 0 // Default value if pointer is nil
	}
	return int(*ptr)
}