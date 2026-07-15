package service

import (
	"mobile/entity"
	"mobile/pkg/conversion"
	"mobile/pkg/structs"
	"mobile/repository"
)

type StockService interface {
	ListGudangUtama(dataFilter entity.StockQueryFilter, EmpId int64, custID string) (response entity.StockGudangUtamaListResponse, err error)
	ListGudangCanvas(dataFilter entity.StockQueryFilter, EmpId int64, custID string) (response entity.StockGudangCanvasistResponse, err error)
	// Detail(RoNo string, custID string) (response entity.OrderResponse, err error)
}

func NewStockService(stockRepository repository.StockRepository, transaction repository.Dbtransaction) *stockServiceImpl {
	return &stockServiceImpl{
		StockRepository: stockRepository,
		Transaction:     transaction,
	}
}

type stockServiceImpl struct {
	StockRepository repository.StockRepository
	Transaction     repository.Dbtransaction
}

func (service *stockServiceImpl) ListGudangUtama(dataFilter entity.StockQueryFilter, EmpId int64, custID string) (response entity.StockGudangUtamaListResponse, err error) {
	ro, err := service.StockRepository.FindByEmpId(EmpId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.StockRepository.FindDetailProductGudangUtama(dataFilter, response.WhId)
	if err != nil {
		return response, err
	}
	for _, row := range details {
		var vResp entity.DetilsGudangUtamaProduct

		qtyUnit := &conversion.Qty{
			Qty:       int(row.Qty),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyUnitConversion conversion.QtyConversionResult
		if row.Qty > 0 {
			qtyUnitConversion = qtyUnit.ConvToQtyConversion()
		}

		qtyUnitOrder := &conversion.Qty{
			Qty:       int(row.QtyOrder),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyUnitOrderConversion conversion.QtyConversionResult
		if row.Qty > 0 {
			qtyUnitOrderConversion = qtyUnitOrder.ConvToQtyConversion()
		}

		qtyIncOrder := &conversion.Qty{
			Qty:       int(row.QtyOrder) + int(row.QtyOrder),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyIncOrderConversion conversion.QtyConversionResult
		if row.Qty > 0 {
			qtyIncOrderConversion = qtyIncOrder.ConvToQtyConversion()
		}

		structs.Automapper(row, &vResp)
		vResp.TotalQty = row.Qty
		vResp.Qty1 = float64(qtyUnitConversion.Qty1)
		vResp.Qty2 = float64(qtyUnitConversion.Qty2)
		vResp.Qty3 = float64(qtyUnitConversion.Qty3)

		vResp.TotalQtyOrder = row.QtyOrder
		vResp.QtyOrder1 = float64(qtyUnitOrderConversion.Qty1)
		vResp.QtyOrder2 = float64(qtyUnitOrderConversion.Qty2)
		vResp.QtyOrder3 = float64(qtyUnitOrderConversion.Qty3)

		vResp.TotalQtyIncOnOrder = row.QtyOrder + row.QtyOrder
		vResp.QtyIncOnOrder1 = float64(qtyIncOrderConversion.Qty1)
		vResp.QtyIncOnOrder2 = float64(qtyIncOrderConversion.Qty2)
		vResp.QtyIncOnOrder3 = float64(qtyIncOrderConversion.Qty3)

		response.DetailsProduct = append(response.DetailsProduct, vResp)

	}

	return response, nil
}

func (service *stockServiceImpl) ListGudangCanvas(dataFilter entity.StockQueryFilter, EmpId int64, custID string) (response entity.StockGudangCanvasistResponse, err error) {
	ro, err := service.StockRepository.FindByEmpIdCanvas(EmpId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.StockRepository.FindDetailProductGudangCanvas(dataFilter, response.WhId)
	if err != nil {
		return response, err
	}
	for _, row := range details {
		var vResp entity.DetilsGudangCanvasProduct

		qtyUnit := &conversion.Qty{
			Qty:       int(row.Qty),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyUnitConversion conversion.QtyConversionResult
		if row.Qty > 0 {
			qtyUnitConversion = qtyUnit.ConvToQtyConversion()
		}

		// var qtyStock conversion.QtyConversionResult
		// if row.QtyStock > 0 {
		// 	qtyStock = qtyUnit.ConvToQtyConversion()
		// }

		// qtyUnitOrder := &conversion.Qty{
		// 	Qty:       int(row.QtyOrder),
		// 	ConvUnit2: int(row.ConvUnit2),
		// 	ConvUnit3: int(row.ConvUnit3),
		// }

		// var qtyUnitOrderConversion conversion.QtyConversionResult
		// if row.Qty > 0 {
		// 	qtyUnitOrderConversion = qtyUnitOrder.ConvToQtyConversion()
		// }

		qtyIncOrder := &conversion.Qty{
			Qty:       int(row.Qty) + int(row.QtyOrder),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyIncOrderConversion conversion.QtyConversionResult
		if row.Qty >= 0 {
			qtyIncOrderConversion = qtyIncOrder.ConvToQtyConversion()
		}

		structs.Automapper(row, &vResp)
		vResp.TotalQtyAvailable = row.Qty
		vResp.Qty1Available = float64(qtyUnitConversion.Qty1)
		vResp.Qty2Available = float64(qtyUnitConversion.Qty2)
		vResp.Qty3Available = float64(qtyUnitConversion.Qty3)

		vResp.TotalQtyStock = row.QtyStock
		vResp.Qty1Stock = float64(qtyIncOrderConversion.Qty1)
		vResp.Qty2Stock = float64(qtyIncOrderConversion.Qty2)
		vResp.Qty3Stock = float64(qtyIncOrderConversion.Qty3)

		response.DetailsProduct = append(response.DetailsProduct, vResp)

	}

	return response, nil
}
