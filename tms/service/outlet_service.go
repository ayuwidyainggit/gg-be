package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/repository"
)

type OutletService interface {
	GetOutlet(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) (response []entity.OutletResponse)
	GetOutletByParams(ctx context.Context, params entity.OutletParams) (response []entity.OutletResponse)
}

type OutletServiceImpl struct {
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	validate             *validator.Validate
}

func NewOutletServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) OutletService {
	return &OutletServiceImpl{
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		validate:             validate,
	}
}

func (service *OutletServiceImpl) GetOutlet(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) (response []entity.OutletResponse) {
	data := service.shipmentInvoicesRepo.FindAll(ctx, dataFilter)

	seenOutletId := make(map[int]bool)
	totalProduct := make(map[int]int)

	for _, row := range data {
		totalProduct[row.OutletID]++
	}

	for _, row := range data {
		if !seenOutletId[row.OutletID] {
			var res entity.OutletResponse
			helper.Automapper(row, &res)
			res.TotalProduct = totalProduct[row.OutletID]
			response = append(response, res)

			seenOutletId[row.OutletID] = true
		}
	}

	return response
}

func (service *OutletServiceImpl) GetOutletByParams(ctx context.Context, params entity.OutletParams) (response []entity.OutletResponse) {
	paramSlice := []any{params.DriverId, params.OutletId, params.ShipmentNo}
	data := service.shipmentInvoicesRepo.FindByOutletId(ctx, paramSlice)

	for _, row := range data {
		var res entity.OutletResponse
		helper.Automapper(row, &res)
		response = append(response, res)
	}
	return response
}
