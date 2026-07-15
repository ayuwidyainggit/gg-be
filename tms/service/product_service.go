package service

import (
	"context"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/repository"

	"github.com/go-playground/validator/v10"
)

type ProductService interface {
	GetProduct(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) (response []entity.ProductResponse)
}

type ProductServiceImpl struct {
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	validate             *validator.Validate
}

func NewProductServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		validate:             validate,
	}
}

func (service *ProductServiceImpl) GetProduct(ctx context.Context, dataFilter entity.ShipmentInvoicesQueryFilter) (response []entity.ProductResponse) {
	result := service.shipmentInvoicesRepo.FindAll(ctx, dataFilter)
	//var data []entity.ProductResponse
	for _, row := range result {
		var res entity.ProductResponse
		helper.Automapper(row, &res)
		if row.Shipment != nil {
			res.DriverID = row.Shipment.DriverID
		}
		if row.OrderNo != nil {
			res.ReturnNo = *row.OrderNo
		}
		res.Salesman = row.SalesmanName
		res.ReturnDate = row.DeliveryDate.Format("2006-01-02")
		res.CtgId1 = "SMALL" // unit_id1 := small

		if row.UnitId2 == row.UnitId3 {
			// unit_id2 == unit_id3 := middle
			res.CtgId2 = "MIDDLE"
			res.CtgId3 = "MIDDLE"
		} else if row.UnitId1 == row.UnitId2{
			res.CtgId2 = "SMALL"
			res.CtgId3 = "MIDDLE"
		} else {
			// unit_id3 != unit_id2 -> unit_id3 := large
			res.CtgId2 = "MIDDLE"
			res.CtgId3 = "LARGE"
		}
		response = append(response, res)
	}
	return response
}
