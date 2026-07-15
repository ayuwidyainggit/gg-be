package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/pkg/conversion"
	"mobile/pkg/structs"
	"mobile/repository"
)

type ProductService interface {
	List(dataFilter entity.ProductsQueryFilter, custId, parentCustId string, EmpId int64, IsActiveGudangCanvas, IsActiveGudangUtama bool) (data []entity.ProductsResp, total int64, lastPage int, err error)
	Detail(entity.DetailProductParams) (entity.ProductDetailResponse, error)
}

type productServiceImpl struct {
	Config                 env.ConfigEnv
	ProductRepository      repository.ProductRepository
	MProductDistRepository repository.MProductDistRepository
	Transaction            repository.Dbtransaction
}

func NewProductService(
	config env.ConfigEnv,
	productRepository repository.ProductRepository,
	mProductDistRepository repository.MProductDistRepository,
	transaction repository.Dbtransaction,
) *productServiceImpl {
	return &productServiceImpl{
		Config:                 config,
		ProductRepository:      productRepository,
		MProductDistRepository: mProductDistRepository,
		Transaction:            transaction,
	}
}

func (service *productServiceImpl) List(dataFilter entity.ProductsQueryFilter, custId, parentCustId string, EmpId int64, IsActiveGudangCanvas, IsActiveGudangUtama bool) (data []entity.ProductsResp, total int64, lastPage int, err error) {
	salesmanTakingOrder, _ := service.ProductRepository.GetSalesmanTakingOrder(custId, EmpId)
	salesmanCanvas, _ := service.ProductRepository.GetSalesmanCanvas(custId, EmpId)
	if salesmanTakingOrder.EmpId == 0 && salesmanCanvas.EmpId == 0 {
		return nil, 0, 0, entity.ErrUserNotSalesman
	}

	if salesmanCanvas.EmpId == 0 && dataFilter.Mode == "canvas" {
		return nil, 0, 0, entity.ErrUserNotSalesmanCanvas
	}
	if salesmanTakingOrder.EmpId == 0 && dataFilter.Mode != "canvas" {
		return nil, 0, 0, entity.ErrUserNotSalesmanTakingOrder
	}

	// If the salesman is not taking order, use the canvas warehouse
	tmpWhId := int64(0)

	if IsActiveGudangCanvas && IsActiveGudangUtama {
		if dataFilter.Mode == "canvas" && salesmanCanvas.EmpId > 0 {
			tmpWhId = int64(salesmanCanvas.WhId)
		} else if salesmanTakingOrder.EmpId > 0 {
			tmpWhId = salesmanTakingOrder.WhId
		}
	} else {
		if !salesmanTakingOrder.IsTakingOrder {
			tmpWhId = int64(salesmanCanvas.WhId)
		}
		// fmt.Println("salesmanTakingOrder.IsTakingOrder >>>", salesmanTakingOrder.IsTakingOrder)
		if salesmanTakingOrder.IsTakingOrder && salesmanTakingOrder.WhId > 0 {
			dataFilter.WhId = &salesmanTakingOrder.WhId
			tmpWhId = salesmanTakingOrder.WhId
			// fmt.Println("IS TAKING ORDER")
		}
	}

	products, total, lastPage, err := service.ProductRepository.FindAllByCustId(dataFilter, custId, parentCustId, tmpWhId, salesmanTakingOrder.IsTakingOrder)
	if err != nil {
		return data, total, lastPage, err
	}

	data = make([]entity.ProductsResp, 0)
	for _, row := range products {
		var vResp entity.ProductsResp
		structs.Automapper(row, &vResp)
		switch row.StatusMMP {
		case 1:
			vResp.StatusManageMinimumPriceName = "Submit"
		case 2:
			vResp.StatusManageMinimumPriceName = "Active"
		default:
			vResp.StatusManageMinimumPriceName = "Non Active"
		}

		switch row.LimitAction {
		case 1:
			vResp.LimitActionName = "Warning"
		case 2:
			vResp.LimitActionName = "Restricted"
		default:
			vResp.LimitActionName = "-"
		}

		qtyUnit := &conversion.Qty{
			Qty:       int(row.Qty),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}
		var qtyUnitConversion conversion.QtyConversionResult
		if row.Qty >= 0 {
			qtyUnitConversion = qtyUnit.ConvToQtyConversion()
		}

		structs.Automapper(row, &vResp)
		// vResp.TotalQty = row.Qty
		vResp.Qty1 = float64(qtyUnitConversion.Qty1)
		vResp.Qty2 = float64(qtyUnitConversion.Qty2)
		vResp.Qty3 = float64(qtyUnitConversion.Qty3)

		vResp.Stock1 = float64(row.Qty)
		if row.ConvUnit2 > 0 {
			vResp.Stock2 = float64(row.Qty) / float64(row.ConvUnit2)
		}
		if row.ConvUnit2 > 0 && row.ConvUnit3 > 0 {
			vResp.Stock3 = float64(row.Qty) / float64(row.ConvUnit2*row.ConvUnit3)
		}
		vResp.Stock4 = 0
		vResp.Stock5 = 0
		vResp.CtgId1 = "SMALL"
		vResp.CtgId2 = "MIDDLE"
		vResp.CtgId3 = "LARGE"

		if row.UnitId1 == row.UnitId2 {
			vResp.CtgId2 = "SMALL"
			vResp.CtgId3 = "LARGE"

		} else if row.UnitId2 == row.UnitId3 {
			vResp.CtgId2 = "LARGE"
			vResp.CtgId3 = "LARGE"

		} else if row.UnitId1 == row.UnitId2 && row.UnitId2 == row.UnitId3 {
			vResp.CtgId1 = "LARGE"
			vResp.CtgId2 = "LARGE"
			vResp.CtgId3 = "LARGE"
		}

		if vResp.LimitAction == 1 {
			vResp.LimitActionName = "Warning"
		}
		if vResp.LimitAction == 2 {
			vResp.LimitActionName = "Restricted"
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) Detail(params entity.DetailProductParams) (response entity.ProductDetailResponse, err error) {
	product, err := service.ProductRepository.FindOneByProductIdAndCustId(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(product, &response)
	if err != nil {
		return response, err
	}

	return response, err
}
