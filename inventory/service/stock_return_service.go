package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/structs"
	"inventory/repository"
	"time"
)

type StockReturnService interface {
	List(dataFilter entity.StockReturnQueryFilter) (data []entity.StockReturnListResponse, total int64, lastPage int, err error)
	Detail(returnNo string, custID string, parentCustID string) (response entity.StockReturnResponse, err error)
	Update(returnNo string, custID string, parentCustID string, request entity.StockReturnUpdateBody) (err error)
	Updatebatch(custID string, parentCustID string, request entity.StockReturnUpdateBatchBody) (err error)
}

func NewStockReturnService(
	stockReturnRepository repository.StockReturnRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction) *StockReturnServiceImpl {
	return &StockReturnServiceImpl{
		StockReturnRepository: stockReturnRepository,
		Transaction:           transaction,
		StockRepository:       stockRepository,
	}
}

type StockReturnServiceImpl struct {
	StockReturnRepository repository.StockReturnRepository
	Transaction           repository.Dbtransaction
	StockRepository       repository.StockRepository
}

func (service *StockReturnServiceImpl) List(dataFilter entity.StockReturnQueryFilter) (data []entity.StockReturnListResponse, total int64, lastPage int, err error) {
	rtns, total, lastPage, err := service.StockReturnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rtns {
		var vResp entity.StockReturnListResponse
		structs.Automapper(row, &vResp)

		if row.InvoiceDate != nil {
			InvDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &InvDate

		}

		returnStatusName := vResp.GenerateReturnStatusName()
		vResp.DataStatusName = &returnStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *StockReturnServiceImpl) Detail(returnNo string, custID string, parentCustID string) (response entity.StockReturnResponse, err error) {
	rtn, err := service.StockReturnRepository.FindOneByReturnNo(returnNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(rtn, &response)
	if err != nil {
		return response, err
	}

	details, err := service.StockReturnRepository.FindReturnDetail(returnNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var DetailsData []entity.StockReturnDetailResponse

	for _, detail := range details {
		var detailData entity.StockReturnDetailResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		if response.InvoiceNo != nil {
			returnedProducts, err := service.StockReturnRepository.CountReturnedProductQty(*response.InvoiceNo, detail.ProductID, custID)
			if err != nil {
				return response, err
			}

			detailData.RemainingQty1 = detail.InvoiceQty1 - returnedProducts.RemainingQty1
			detailData.RemainingQty2 = detail.InvoiceQty2 - returnedProducts.RemainingQty2
			detailData.RemainingQty3 = detail.InvoiceQty3 - returnedProducts.RemainingQty3
		}

		itemConditionName := detailData.GenerateItemConditionName()
		detailData.ItemCndName = &itemConditionName

		DetailsData = append(DetailsData, detailData)
	}

	if rtn.ReturnDate != nil {
		returnDate := rtn.ReturnDate.Format("2006-01-02")
		response.ReturnDate = &returnDate
	}

	if rtn.InvoiceDate != nil {
		invoiceDate := rtn.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	returnStatusName := response.GenerateReturnStatusName()
	response.DataStatusName = &returnStatusName

	response.Details = DetailsData
	return response, nil
}

func (service *StockReturnServiceImpl) Update(returnNo string, custID string, parentCustID string, request entity.StockReturnUpdateBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		now := time.Now()
		err = service.StockReturnRepository.UpdateStatus(txCtx, custID, returnNo, request.DataStatus)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var returnDetailModel model.StockReturnDetail
			err = structs.Automapper(detail, &returnDetailModel)
			if err != nil {
				return err
			}

			returnDetailModel.CustID = ""
			err = service.StockReturnRepository.UpdateDetail(txCtx, &returnDetailModel)
			if err != nil {
				return err
			}
		}

		if request.DataStatus == entity.COMPLETED {
			rtn, err := service.StockReturnRepository.FindOneByReturnNo(returnNo, custID, parentCustID)
			if err != nil {
				return err
			}

			err = service.StockReturnRepository.UpdateClosedAt(txCtx, custID, returnNo, now)
			if err != nil {
				return err
			}

			details, err := service.StockReturnRepository.FindReturnDetail(returnNo, custID, parentCustID)
			if err != nil {
				return err
			}
			var stockUpdateEntities []*entity.StockUpdate

			for _, detail := range details {

				QtyShipUnit := &conversion.QtyUnit{
					Qty1:      int(detail.Qty1),
					Qty2:      int(detail.Qty2),
					Qty3:      int(detail.Qty3),
					ConvUnit2: int(*detail.ConvUnit2),
					ConvUnit3: int(*detail.ConvUnit3),
				}

				totalQtyShip, err := QtyShipUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    custID,
					WhID:      *detail.WhId,
					ProID:     detail.ProductID,
					StockDate: *rtn.ReturnDate,
					TrCode:    detail.ReturnNo[0:2],
					TrNo:      detail.ReturnNo,
					QtyIn:     float64(totalQtyShip),
					UnitPrice: detail.SellPrice1,
					RefDetId:  detail.ReturnDetailID,
				}

				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
			}

			if len(stockUpdateEntities) > 0 {
				err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *StockReturnServiceImpl) Updatebatch(custID string, parentCustID string, request entity.StockReturnUpdateBatchBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		stockReturns, err := service.StockReturnRepository.FindReturnDetailByListNo(request.ReturnsNo, custID, parentCustID)
		if err != nil {
			return err
		}
		var stockReturnsMap = model.MapStockReturn{}

		for _, stockReturn := range stockReturns {
			stockReturnsMap.Set(stockReturn.ReturnNo, stockReturn)
		}

		now := time.Now()
		err = service.StockReturnRepository.UpdatebatchStatus(txCtx, custID, request.ReturnsNo, request.DataStatus)
		if err != nil {
			return err
		}

		if request.DataStatus == entity.COMPLETED {
			err = service.StockReturnRepository.UpdatebatchClosedAt(txCtx, custID, request.ReturnsNo, now)
			if err != nil {
				return err
			}
			for _, returnNo := range request.ReturnsNo {
				rtn, err := stockReturnsMap.GetByID(returnNo)
				if err != nil {
					return err
				}

				details, err := service.StockReturnRepository.FindReturnDetail(returnNo, custID, parentCustID)
				if err != nil {
					return err
				}
				var stockUpdateEntities []*entity.StockUpdate

				for _, detail := range details {

					QtyShipUnit := &conversion.QtyUnit{
						Qty1:      int(detail.Qty1),
						Qty2:      int(detail.Qty2),
						Qty3:      int(detail.Qty3),
						ConvUnit2: int(*detail.ConvUnit2),
						ConvUnit3: int(*detail.ConvUnit3),
					}

					totalQtyShip, err := QtyShipUnit.ToTotalQuantity()
					if err != nil {
						return err
					}

					stockUpdateEntity := entity.StockUpdate{
						CustID:    custID,
						WhID:      *detail.WhId,
						ProID:     detail.ProductID,
						StockDate: *rtn.ReturnDate,
						TrCode:    detail.ReturnNo[0:2],
						TrNo:      detail.ReturnNo,
						QtyIn:     float64(totalQtyShip),
						UnitPrice: detail.SellPrice1,
						RefDetId:  detail.ReturnDetailID,
					}

					stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
				}

				if len(stockUpdateEntities) > 0 {
					err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
					if err != nil {
						return err
					}
				}
			}

		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
