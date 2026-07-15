package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"inventory/pkg/conversion"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type WhAdjService interface {
	Store(request entity.CreateWhAdjBody) (err error)
	Detail(adjNo, custID, langId string) (response entity.WhAdjResponse, err error)
	List(dataFilter entity.WhAdjQueryFilter, custId, langId string) (data []entity.WhAdjListResponse, total int64, lastPage int, err error)
	Delete(custId string, adcNo string, userId int64) (err error)
	UpdateStatus(adjNo string, request entity.UpdateWhAdjStatusBody, langId string) (err error)
	ListWarehouse(dataFilter entity.WhAdjWarehouseQueryFilter, custId string) (data []entity.WarehouseAdjustment, total int64, lastPage int, err error)
}

type whAdjServiceImpl struct {
	WhAdjRepository repository.WhAdjRepository
	Transaction     repository.Dbtransaction
	StockRepository repository.StockRepository
}

func NewWhAdjService(WhAdjRepository repository.WhAdjRepository, transaction repository.Dbtransaction, stockRepository repository.StockRepository) *whAdjServiceImpl {
	return &whAdjServiceImpl{
		WhAdjRepository: WhAdjRepository,
		Transaction:     transaction,
		StockRepository: stockRepository,
	}
}
func (service *whAdjServiceImpl) Store(request entity.CreateWhAdjBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	var whAdjModel model.WhAdj

	err = structs.Automapper(request, &whAdjModel)
	if err != nil {
		return err
	}

	adjDate := str.GetJakartaDate()

	whAdjModel.AdjDate = adjDate
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.WhAdjRepository.Store(txCtx, &whAdjModel)
		if err != nil {
			return err
		}

		var productIDs []int64
		for _, detail := range request.Details {
			productIDs = append(productIDs, detail.ProID)
		}

		productsModel, err := service.WhAdjRepository.FindProductByListID(productIDs)
		if err != nil {
			return err
		}

		var productMap = model.MapProduct{}

		for _, productModel := range productsModel {
			productMap.SetProduct(productModel.ProductId, productModel)
		}

		for _, Detail := range request.Details {

			productModel, err := productMap.GetByID(Detail.ProID)
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(Detail.Qty1),
				Qty2:      int(Detail.Qty2),
				Qty3:      int(Detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			// parse time format YYYY-mm-dd to Rfc3339
			Detail.StockAdjNo = whAdjModel.AdjNo
			var whAdjDetModel model.WhAdjDet
			err = structs.Automapper(Detail, &whAdjDetModel)
			if err != nil {
				return err
			}

			whAdjDetModel.Qty = totalQty
			whAdjDetModel.CustID = request.CustID
			err = service.WhAdjRepository.StoreDetail(txCtx, &whAdjDetModel)
			if err != nil {
				return err
			}

		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *whAdjServiceImpl) Detail(adjNo, custID, langId string) (response entity.WhAdjResponse, err error) {
	whAdj, err := service.WhAdjRepository.FindByNo(adjNo, custID, langId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(whAdj, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.WhAdjRepository.FindDetail(adjNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.WhAdjDetresponse
	for _, detail := range Details {

		qty := &conversion.Qty{
			Qty:       int(detail.Qty),
			ConvUnit2: int(detail.ConvUnit2),
			ConvUnit3: int(detail.ConvUnit3),
		}

		qtyConversion := qty.ConvToQtyConversion()

		var detailData entity.WhAdjDetresponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		detailData.Qty1 = qtyConversion.Qty1
		detailData.Qty2 = qtyConversion.Qty2
		detailData.Qty3 = qtyConversion.Qty3
		DetailsData = append(DetailsData, detailData)
	}

	whAdjDate := whAdj.AdjDate.Format("2006-01-02")
	response.StockAdjDate = whAdjDate

	response.Details = DetailsData
	return response, nil
}
func (service *whAdjServiceImpl) List(dataFilter entity.WhAdjQueryFilter, custId, langId string) (data []entity.WhAdjListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.WhAdjRepository.FindAllByCustId(dataFilter, custId, langId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.WhAdjListResponse
		structs.Automapper(row, &vResp)
		if row.AdjDate != nil {
			whTrfDate := row.AdjDate.Format("2006-01-02")
			vResp.StockAdjDate = whTrfDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *whAdjServiceImpl) Delete(custId string, adjNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhAdjRepository.Delete(txCtx, custId, adjNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *whAdjServiceImpl) UpdateStatus(adjNo string, request entity.UpdateWhAdjStatusBody, langId string) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhAdjRepository.UpdateStatus(txCtx, adjNo, request.CustID, request.DataStatus)
		if err != nil {
			return err
		}

		if request.DataStatus == constant.STOCK_ADJUTMENT_STATUS_REJECT { // if reject, just update status on main data
			return nil
		} else if request.DataStatus == constant.STOCK_ADJUTMENT_STATUS_APPROVED {
			whAdj, err := service.WhAdjRepository.FindByNo(adjNo, request.CustID, langId)
			if err != nil {
				return err
			}

			details, err := service.WhAdjRepository.FindDetail(adjNo, request.CustID)
			if err != nil {
				return err
			}

			var stockUpdateEntities []*entity.StockUpdate
			for _, detail := range details {
				var qtyin, qtyout float64

				if detail.WhAdjDetType == constant.STOCK_ADJUSTMENT_DET_TYPE_ADD_STOCK {
					qtyin = detail.Qty
				} else {
					qtyout = detail.Qty
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    request.CustID,
					WhID:      whAdj.WhID,
					ProID:     int64(detail.ProID),
					StockDate: *whAdj.AdjDate,
					TrCode:    whAdj.AdjNo[0:2],
					TrNo:      whAdj.AdjNo,
					QtyIn:     qtyin,
					QtyOut:    qtyout,
					UnitPrice: 0,
					RefDetId:  int64(*detail.WhAdjDetId),
				}
				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
			}

			err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func (service *whAdjServiceImpl) ListWarehouse(dataFilter entity.WhAdjWarehouseQueryFilter, custId string) (data []entity.WarehouseAdjustment, total int64, lastPage int, err error) {
	warehouses, total, lastPage, err := service.WhAdjRepository.FindWarehouseStockAdjusment(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range warehouses {
		var vResp entity.WarehouseAdjustment
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}
