package service

import (
	"context"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UnloadService interface {
	TravelList(ctx context.Context, params entity.TravelListParams) (response entity.TravelListResponse)
	Unload(ctx context.Context, headers map[string]string, request entity.UnloadRequest)
	Resume(ctx context.Context, request entity.UnloadRequest)
	Onhold(ctx context.Context, request entity.UnloadRequest)
}

type UnloadServiceImpl struct {
	shipmentInvoicesRepo    repository.ShipmentInvoicesRepo
	shipmentService         ShipmentService
	shipmentOrderStatusRepo repository.ShipmentOrderStatusRepo
	validate                *validator.Validate
}

func NewUnloadServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, shipmentService ShipmentService, shipmentOrderStatusRepo repository.ShipmentOrderStatusRepo, validate *validator.Validate) UnloadService {
	return &UnloadServiceImpl{
		shipmentInvoicesRepo:    shipmentInvoicesRepo,
		shipmentService:         shipmentService,
		shipmentOrderStatusRepo: shipmentOrderStatusRepo,
		validate:                validate,
	}
}

func (service *UnloadServiceImpl) TravelList(ctx context.Context, params entity.TravelListParams) (response entity.TravelListResponse) {
	data, err := service.shipmentInvoicesRepo.FindTodoList(ctx, params.OutletId, params.ShipmentNo)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	helper.Automapper(data, &response)

	return response

}

func (service *UnloadServiceImpl) Unload(ctx context.Context, headers map[string]string, request entity.UnloadRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		UnloadAt:      &request.CurrentTime,
		OutletStatus:  "On Progress", // todo change Receive All to be on progress for unload
		ProductStatus: "Receive",
		ShipmentNo:    &request.ShipmentNo,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAtUnload(ctx, request.OutletID, dataset)
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
	orderNos := service.shipmentInvoicesRepo.FindAllOrderNoByShipmentNo(ctx, request.ShipmentNo, request.OutletID)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// Process each order number
	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", orderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: orderNo,
						Status:  4,
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

func (service *UnloadServiceImpl) Resume(ctx context.Context, request entity.UnloadRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ResumeAt:     &request.CurrentTime,
		OutletStatus: "On Resume",
		ShipmentNo:   &request.ShipmentNo,
		// OnHold:       nil,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAt(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.shipmentInvoicesRepo.UpdateColumnAt(ctx, "on_hold", nil, request.OutletID, dataset)
}

func (service *UnloadServiceImpl) Onhold(ctx context.Context, request entity.UnloadRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		OnHold:       &request.CurrentTime,
		OutletStatus: "On Hold",
		ShipmentNo:   &request.ShipmentNo,
		// ResumeAt:     nil,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAt(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
	service.shipmentInvoicesRepo.UpdateColumnAt(ctx, "resume_at", nil, request.OutletID, dataset)
}
