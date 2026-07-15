package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
)

type ArriveService interface {
	Arrive(ctx context.Context, request entity.ArriveRequest)
}

type ArriveServiceImpl struct {
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	validate             *validator.Validate
}

func NewArriveServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) ArriveService {
	return &ArriveServiceImpl{
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		validate:             validate,
	}
}

func (service *ArriveServiceImpl) Arrive(ctx context.Context, request entity.ArriveRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ArriveAt:     &request.CurrentTime,
		OutletStatus: "On Progress",
		ShipmentNo:   &request.ShipmentNo,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAt(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}
