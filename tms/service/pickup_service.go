package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"
)

type PickUpService interface {
	PickUpAll(ctx context.Context, headers map[string]string, request entity.PickUpRequest)
	SkipPickUp(ctx context.Context, request entity.SkipPickUpRequest)
	PickUpPartial(ctx context.Context, headers map[string]string, request entity.PickupPartialRequest)
}

type PickUpServiceImpl struct {
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	shipmentService      ShipmentService
	validate             *validator.Validate
}

func NewPickUpServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, shipmentService ShipmentService, validate *validator.Validate) PickUpService {
	return &PickUpServiceImpl{
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		shipmentService:      shipmentService,
		validate:             validate,
	}
}

func (service *PickUpServiceImpl) PickUpAll(ctx context.Context, headers map[string]string, request entity.PickUpRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ProductStatus: "Pick Up",
		PickupAt:      &request.CurrentTime,
	}

	err = service.shipmentInvoicesRepo.UpdatePickUp(ctx, request.ID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	orderNos := service.shipmentInvoicesRepo.FindAllOrderNoById(ctx, request.ID)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SR") {
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: orderNo,
						Status:  5,
					},
				},
			}

			err = service.shipmentService.MobileUpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return
			}
		}
	}
}

func (service *PickUpServiceImpl) SkipPickUp(ctx context.Context, request entity.SkipPickUpRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	for _, req := range request.Data {
		dataset := model.ShipmentInvoices{
			ProductStatus: "Skip",
			ReasonID:      &req.ReasonID,
			ReasonName:    &req.ReasonName,
		}
		err := service.shipmentInvoicesRepo.UpdatePickUpPartial(ctx, req.ID, dataset)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}
	}
}

func (service *PickUpServiceImpl) PickUpPartial(ctx context.Context, headers map[string]string, request entity.PickupPartialRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	// Iterate over the products in request.Data
	for _, product := range request.Data.Products {
		dataset := model.ShipmentInvoices{
			OutletID:      request.Data.OutletID,
			ProductStatus: "Pick Up",
			ShipmentNo:    &request.Data.ShipmentNo,
			PickupAt:      &request.Data.CurrentTime,
		}
		// Handle multiple quantities using switch
		for i, qty := range product.Qty {
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
		}
		err := service.shipmentInvoicesRepo.UpdatePickUpPartial(ctx, product.ID, dataset)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		orderNos := service.shipmentInvoicesRepo.FindAllOrderNoByShipmentNo(ctx, request.Data.ShipmentNo, request.Data.OutletID)
		if err != nil {
			panic(exception.NewInternalServerError(err.Error()))
		}

		for _, orderNo := range orderNos {
			if strings.HasPrefix(orderNo, "SR") {
				returnUpdate := entity.UpdateStatusReturn{
					Returns: []entity.ReturnItem{
						{
							OrderNo: orderNo,
							Status:  5,
						},
					},
				}

				err = service.shipmentService.MobileUpdateStatusReturn(ctx, headers, returnUpdate)
				if err != nil {
					log.Printf("Error updating order status: %v", err)
					return
				}
			}
		}
	}
}
