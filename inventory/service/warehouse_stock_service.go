package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/structs"
	"inventory/pkg/validation"
	"inventory/repository"
	"log"
)

type WarehouseStockService interface {
	List(dataFilter entity.DistributorStockQueryFilter) (data []entity.DistributorStockList, total int64, lastPage int, err error)
	WarehouseList(dataFilter entity.WarehouseStockWhListQueryFilter) (data []entity.WarehouseStockWhList, total int64, lastPage int, err error)
	Upsert(request entity.UpsertWarehouseStock) (err error)
	UpsertBulk(request entity.UpsertBulkWarehouseStock) (err error)
	ProductList(dataFilter entity.ProductWarehouseListQueryFilter) (data []entity.ProductWarehouseList, total int64, lastPage int, err error)
}

func NewWarehouseStockService(stockRepository repository.WarehouseStockRepository,
	transaction repository.Dbtransaction,
	validator *validation.Validate) *warehouseStockImpl {
	return &warehouseStockImpl{
		WarehouseStockRepository: stockRepository,
		Transaction:              transaction,
		Validator:                validator,
	}
}

type warehouseStockImpl struct {
	WarehouseStockRepository repository.WarehouseStockRepository
	Transaction              repository.Dbtransaction
	Validator                *validation.Validate
}

func (service *warehouseStockImpl) List(dataFilter entity.DistributorStockQueryFilter) (data []entity.DistributorStockList, total int64, lastPage int, err error) {

	stocks, total, lastPage, err := service.WarehouseStockRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	// log.Println("stocks:", structs.StructToJson(stocks))
	for _, row := range stocks {
		var vResp entity.DistributorStockList
		qty := &conversion.Qty{
			Qty:       int(row.Qty),
			ConvUnit2: int(row.ConvUnit2),
			ConvUnit3: int(row.ConvUnit3),
		}
		qtyConversion := qty.ConvToQtyConversion()
		structs.Automapper(row, &vResp)
		vResp.TotalQty = float64(row.Qty)
		vResp.Qty1 = float64(qtyConversion.Qty1)
		vResp.Qty2 = float64(qtyConversion.Qty2)
		vResp.Qty3 = float64(qtyConversion.Qty3)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *warehouseStockImpl) WarehouseList(dataFilter entity.WarehouseStockWhListQueryFilter) (data []entity.WarehouseStockWhList, total int64, lastPage int, err error) {

	stocks, total, lastPage, err := service.WarehouseStockRepository.FindAllWarehouse(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	// log.Println("stocks:", structs.StructToJson(stocks))
	for _, row := range stocks {
		var vResp entity.WarehouseStockWhList
		structs.Automapper(row, &vResp)
		// stockDate := row.WarehouseStockDate.Format("2006-01-02")
		// vResp.WarehouseStockDate = stockDate
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *warehouseStockImpl) Upsert(request entity.UpsertWarehouseStock) (err error) {
	c := context.Background()

	// request.ItemCdn = 4
	//validate & format input data
	request, err = service.validateAndFormat(request)
	if err != nil {
		return err
	}

	var stockModel model.WarehouseStock
	err = structs.Automapper(request, &stockModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.WarehouseStockRepository.Upsert(txCtx, &stockModel)
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

func (service *warehouseStockImpl) UpsertBulk(request entity.UpsertBulkWarehouseStock) (err error) {
	c := context.Background()

	// request.ItemCdn = 4
	//validate & format input data
	request, err = service.validateAndFormatBulk(request)
	if err != nil {
		return err
	}

	var stockModel []*model.WarehouseStock
	err = structs.Automapper(request.WarehouseStock, &stockModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.WarehouseStockRepository.UpsertBulk(txCtx, stockModel)
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

func (service *warehouseStockImpl) validateAndFormat(req entity.UpsertWarehouseStock) (entity.UpsertWarehouseStock, error) {

	err := service.Validator.ValidateStructReturnError(req)
	if err != nil {
		log.Println("validateAndFormat, err:", err)
		return req, err
	}

	// parse time format YYYY-mm-dd to Rfc3339
	// stockDate, err := str.DateStrToRfc3339String(req.WarehouseStockDate)
	// if err != nil {
	// 	err := fmt.Errorf(errmsg.ERROR_DATE_FORMAT, "WarehouseStockDate", req.WarehouseStockDate)
	// 	return req, err
	// }
	// req.WarehouseStockDate = stockDate

	return req, err
}

func (service *warehouseStockImpl) validateAndFormatBulk(req entity.UpsertBulkWarehouseStock) (entity.UpsertBulkWarehouseStock, error) {

	err := service.Validator.ValidateStructReturnError(req)
	if err != nil {
		log.Println("validateAndFormatBulk, err:", err)
		return req, err
	}

	// for i, _ := range req.WarehouseStock {
	// 	// parse time format YYYY-mm-dd to Rfc3339
	// 	stockDate, err := str.DateStrToRfc3339String(req.WarehouseStock[i].WarehouseStockDate)
	// 	if err != nil {
	// 		err := fmt.Errorf(errmsg.ERROR_DATE_FORMAT, "WarehouseStockDate", req.WarehouseStock[i].WarehouseStockDate)
	// 		return req, err
	// 	}
	// 	req.WarehouseStock[i].WarehouseStockDate = stockDate
	// }

	return req, err
}

func (service *warehouseStockImpl) ProductList(dataFilter entity.ProductWarehouseListQueryFilter) (data []entity.ProductWarehouseList, total int64, lastPage int, err error) {
	productWarehouses, total, lastPage, err := service.WarehouseStockRepository.ProductList(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, productWarehouse := range productWarehouses {
		var vResp entity.ProductWarehouseList
		qty := &conversion.Qty{
			Qty:       int(productWarehouse.Qty),
			ConvUnit2: int(productWarehouse.ConvUnit2),
			ConvUnit3: int(productWarehouse.ConvUnit3),
		}
		qtyConversion := qty.ConvToQtyConversion()
		structs.Automapper(productWarehouse, &vResp)
		vResp.Qty1 = float64(qtyConversion.Qty1)
		vResp.Qty2 = float64(qtyConversion.Qty2)
		vResp.Qty3 = float64(qtyConversion.Qty3)

		data = append(data, vResp)
	}
	return
}
