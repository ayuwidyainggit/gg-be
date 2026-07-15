package service

import (
	"context"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/errmsg"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/pkg/validation"
	"inventory/repository"
	"log"
)

type StockService interface {
	List(dataFilter entity.StockQueryFilter, custID, parentCustID string) (data []entity.StockList, total int64, lastPage int, err error)
	Report(dataFilter entity.StockReportQueryFilter) (data []entity.StockReport, total int64, lastPage int, err error)
	Store(request entity.CreateStock) (err error)
	StoreBulk(request entity.CreateBulkStock) (err error)
	OpnameLookup(dataFilter entity.StockOpnameLookupQueryFilter) (data []entity.StockOpnameLookup, total int64, lastPage int, err error)
}

func NewStockService(stockRepository repository.StockRepository,
	transaction repository.Dbtransaction,
	validator *validation.Validate) *stockServiceImpl {
	return &stockServiceImpl{
		StockRepository: stockRepository,
		Transaction:     transaction,
		Validator:       validator,
	}
}

type stockServiceImpl struct {
	StockRepository repository.StockRepository
	Transaction     repository.Dbtransaction
	Validator       *validation.Validate
}

func (service *stockServiceImpl) List(dataFilter entity.StockQueryFilter, custId, parentCustId string) (data []entity.StockList, total int64, lastPage int, err error) {

	stocks, total, lastPage, err := service.StockRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	// log.Println("stocks:", structs.StructToJson(stocks))
	for _, row := range stocks {
		var vResp entity.StockList
		structs.Automapper(row, &vResp)
		stockDate := row.StockDate.Format("2006-01-02")
		vResp.StockDate = stockDate
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *stockServiceImpl) Store(request entity.CreateStock) (err error) {
	c := context.Background()

	// request.ItemCdn = 4
	//validate & format input data
	request, err = service.validateAndFormat(request)
	if err != nil {
		return err
	}

	var stockModel model.Stock
	err = structs.Automapper(request, &stockModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockRepository.Store(txCtx, &stockModel)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *stockServiceImpl) StoreBulk(request entity.CreateBulkStock) (err error) {
	c := context.Background()

	// request.ItemCdn = 4
	//validate & format input data
	request, err = service.validateAndFormatBulk(request)
	if err != nil {
		return err
	}

	var stockModel []*model.Stock
	err = structs.Automapper(request.Stock, &stockModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockRepository.StoreBulk(txCtx, stockModel)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *stockServiceImpl) validateAndFormat(req entity.CreateStock) (entity.CreateStock, error) {

	err := service.Validator.ValidateStructReturnError(req)
	if err != nil {
		log.Println("validateAndFormat, err:", err)
		return req, err
	}

	// parse time format YYYY-mm-dd to Rfc3339
	stockDate, err := str.DateStrToRfc3339String(req.StockDate)
	if err != nil {
		err := fmt.Errorf(errmsg.ERROR_DATE_FORMAT, "StockDate", req.StockDate)
		return req, err
	}
	req.StockDate = stockDate

	return req, err
}

func (service *stockServiceImpl) validateAndFormatBulk(req entity.CreateBulkStock) (entity.CreateBulkStock, error) {

	err := service.Validator.ValidateStructReturnError(req)
	if err != nil {
		log.Println("validateAndFormatBulk, err:", err)
		return req, err
	}

	for i, _ := range req.Stock {
		// parse time format YYYY-mm-dd to Rfc3339
		stockDate, err := str.DateStrToRfc3339String(req.Stock[i].StockDate)
		if err != nil {
			err := fmt.Errorf(errmsg.ERROR_DATE_FORMAT, "StockDate", req.Stock[i].StockDate)
			return req, err
		}
		req.Stock[i].StockDate = stockDate
	}

	return req, err
}

func (service *stockServiceImpl) Report(dataFilter entity.StockReportQueryFilter) (data []entity.StockReport, total int64, lastPage int, err error) {

	stocks, total, lastPage, err := service.StockRepository.Report(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range stocks {
		var vResp entity.StockReport

		qtyUnit := &conversion.Qty{
			Qty:       int(row.Qty),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyUnitConversion conversion.QtyConversionResult
		if row.Qty >= 0 {
			qtyUnitConversion = qtyUnit.ConvToQtyConversion()
		}

		qtyUnitOrder := &conversion.Qty{
			Qty:       int(row.QtyOrder),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}

		var qtyUnitOrderConversion conversion.QtyConversionResult
		if row.Qty >= 0 {
			qtyUnitOrderConversion = qtyUnitOrder.ConvToQtyConversion()
		}

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
		vResp.TotalQty = row.Qty
		vResp.Qty1 = float64(qtyUnitConversion.Qty1) + row.OrderQty1
		vResp.Qty2 = float64(qtyUnitConversion.Qty2) + row.OrderQty2
		vResp.Qty3 = float64(qtyUnitConversion.Qty3) + row.OrderQty3

		vResp.TotalQtyOrder = row.QtyOrder
		vResp.QtyOrder1 = float64(qtyUnitOrderConversion.Qty1)
		vResp.QtyOrder2 = float64(qtyUnitOrderConversion.Qty2)
		vResp.QtyOrder3 = float64(qtyUnitOrderConversion.Qty3)

		vResp.TotalQtyIncOnOrder = row.QtyOrder + row.QtyOrder
		vResp.QtyIncOnOrder1 = float64(qtyIncOrderConversion.Qty1)
		vResp.QtyIncOnOrder2 = float64(qtyIncOrderConversion.Qty2)
		vResp.QtyIncOnOrder3 = float64(qtyIncOrderConversion.Qty3)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *stockServiceImpl) OpnameLookup(dataFilter entity.StockOpnameLookupQueryFilter) (data []entity.StockOpnameLookup, total int64, lastPage int, err error) {

	stocks, total, lastPage, err := service.StockRepository.OpnameLookup(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range stocks {
		var vResp entity.StockOpnameLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
