package service

import (
	"context"
	"errors"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type BpprService interface {
	Detail(bpprNo, custId, parentCustId string) (response entity.BpprResponse, err error)
	Store(request entity.CreateBpprBody) (response entity.BpprResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.BpprResponse, total int64, lastPage int, err error)
	Update(bpprNo string, request entity.UpdateBpprRequest) (err error)
	Delete(custId string, bpprNo string, userId int64) (err error)
}

func NewBpprService(
	bpprRepository repository.BpprRepository,
	whStockRepository repository.WhStockRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction) *bpprServiceImpl {
	return &bpprServiceImpl{
		BpprRepository:    bpprRepository,
		WhStockRepository: whStockRepository,
		StockRepository:   stockRepository,
		Transaction:       transaction,
	}
}

type bpprServiceImpl struct {
	BpprRepository    repository.BpprRepository
	WhStockRepository repository.WhStockRepository
	StockRepository   repository.StockRepository
	Transaction       repository.Dbtransaction
}

func (service *bpprServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.BpprResponse, total int64, lastPage int, err error) {
	bpprs, total, lastPage, err := service.BpprRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range bpprs {
		var vResp entity.BpprResponse
		structs.Automapper(row, &vResp)
		bpprDate := row.BpprDate.Format("2006-01-02")
		var returnDate string
		if row.ReturnDate != nil {
			returnDate = row.ReturnDate.Format("2006-01-02")
		}
		vResp.BpprDate = bpprDate
		vResp.ReturnDate = returnDate
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *bpprServiceImpl) Detail(bpprNo, custId, parentCustId string) (response entity.BpprResponse, err error) {
	bppr, err := service.BpprRepository.FindByNo(bpprNo, custId, parentCustId)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(bppr, &response)
	if err != nil {
		return response, err
	}

	bpprDate := bppr.BpprDate.Format("2006-01-02")
	var returnDate string
	if bppr.ReturnDate != nil {
		returnDate = bppr.ReturnDate.Format("2006-01-02")
	}
	response.BpprDate = bpprDate
	response.ReturnDate = returnDate

	bpprDetails, err := service.BpprRepository.FindBpprDetails(bpprNo, custId)
	if err != nil {
		return response, err
	}
	var bpprDetailsData []entity.BpprDet
	for _, bpprDetail := range bpprDetails {
		var bpprDetailData entity.BpprDet
		err = structs.Automapper(bpprDetail, &bpprDetailData)
		if err != nil {
			return response, err
		}
		if bpprDetail.ExpDate != nil {
			bpprDetailDate := bpprDetail.ExpDate.Format("2006-01-02")
			bpprDetailData.ExpDate = &bpprDetailDate
		}

		bpprDetailsData = append(bpprDetailsData, bpprDetailData)
	}
	response.Details = bpprDetailsData
	return response, nil
}

func (service *bpprServiceImpl) Store(request entity.CreateBpprBody) (response entity.BpprResponse, err error) {
	c := context.Background()

	if len(request.Details) == 0 {
		return response, errors.New("item details is required")
	}

	// parse time format YYYY-mm-dd to Rfc3339
	var returnDate, bpprDate string
	if request.ReturnDate != nil {
		if *request.ReturnDate != "" {
			returnDate, err = str.DateStrToRfc3339String(*request.ReturnDate)
			if err != nil {
				return response, err
			}
			request.ReturnDate = &returnDate
		} else {
			request.ReturnDate = nil
		}
	}

	if request.BpprDate != nil {
		if *request.BpprDate != "" {
			bpprDate, err = str.DateStrToRfc3339String(*request.BpprDate)
			if err != nil {
				return response, err
			}
			request.BpprDate = &bpprDate
		} else {
			request.BpprDate = nil
		}
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var bpprModel *model.Bppr
	err = structs.Automapper(request, &bpprModel)
	if err != nil {
		return response, err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.BpprRepository.Store(txCtx, bpprModel)
		if err != nil {
			return err
		}
		response.BpprNo = bpprModel.BpprNo

		for index, detail := range request.Details {
			var bpprDetailModel model.BpprDet
			seq := index + 1
			if detail.ExpDate != nil {
				// parse time format YYYY-mm-dd to Rfc3339
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &bpprDetailModel)
			if err != nil {
				return err
			}
			bpprDetailModel.SeqNo = seq
			bpprDetailModel.CustID = request.CustID
			bpprDetailModel.BpprNo = bpprModel.BpprNo
			_, err = service.BpprRepository.CreateBpprDetail(txCtx, &bpprDetailModel)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return response, err
}

func (service *bpprServiceImpl) Update(bpprNo string, request entity.UpdateBpprRequest) (err error) {
	c := context.Background()

	workDayActive, err := service.BpprRepository.FindActiveWorkDay(request.ParentCustID)
	if err != nil {
		return err
	}

	bppr, err := service.BpprRepository.FindByNo(bpprNo, request.CustID, request.ParentCustID)
	if err != nil {
		return err
	}

	if bppr.DataStatus == 2 {
		errorMsg := "your document has been returned"
		return errors.New(errorMsg)
	}

	// set value if return process
	if request.DataStatus == 2 {
		// request.ReturnNo = ""
		returnDate := workDayActive.WorkDate.Format("2006-01-02")
		request.ReturnDate = &returnDate

		if request.ReturnReasonID < 1 {
			return errors.New("return reason is required")
		}

		var bpprDate string
		if request.BpprDate != nil {
			if *request.BpprDate != "" {
				bpprDate, err = str.DateStrToRfc3339String(*request.BpprDate)
				if err != nil {
					return err
				}
				request.BpprDate = &bpprDate
			} else {
				return errors.New("bppr date is required")
			}
		}
	}

	if request.ReturnDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.ReturnDate != "" {
			returnDate, err := str.DateStrToRfc3339String(*request.ReturnDate)
			if err != nil {
				return err
			}
			request.ReturnDate = &returnDate
		}
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var bpprModel model.Bppr
	err = structs.Automapper(request, &bpprModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// returnNo, err := service.BpprRepository.Update(txCtx, bpprNo, bpprModel)
		_, err := service.BpprRepository.Update(txCtx, bpprNo, bpprModel)
		if err != nil {
			return err
		}

		err = service.BpprRepository.DeleteBpprDetailByBpprNo(txCtx, bpprNo)
		if err != nil {
			return err
		}

		var stringNil string
		bpprModel.CustID = &stringNil
		for index, detail := range request.Details {
			sequence := index + 1

			var bpprDetailModel model.BpprDet

			if detail.ExpDate != nil {
				if *detail.ExpDate != "" {
					expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
					if err != nil {
						return err
					}
					detail.ExpDate = &expDate

				}
			}

			err = structs.Automapper(detail, &bpprDetailModel)
			if err != nil {
				return err
			}
			bpprDetailModel.SeqNo = sequence
			bpprDetailModel.CustID = request.CustID
			bpprDetailModel.BpprNo = bpprNo

			// bpprDet, err := service.BpprRepository.CreateBpprDetail(txCtx, &bpprDetailModel)
			_, err := service.BpprRepository.CreateBpprDetail(txCtx, &bpprDetailModel)
			if err != nil {
				return err
			}

			if request.DataStatus == 2 {
				/*
					whStockQuery := entity.WhStockQuery{
						CustID: bppr.CustID,
						WhId:   *bppr.WhID,
						ProId:  int64(detail.ProID),
					}
					whStock, err := service.WhStockRepository.FindByWhIdAndProId(whStockQuery)
					if err != nil {
						return err
					}

					newWhStock := model.WhStock{}
					if *bppr.ItemCdn == 1 {
						oldQty := *whStock.Qty
						newQty := oldQty - detail.Qty
						newWhStock.Qty = &newQty
						if *newWhStock.Qty < 0 {
							return errors.New("quantity stock must be greater than 0")
						}
					} else if *bppr.ItemCdn == 2 {
						oldQty := *whStock.QtyBs
						newQty := oldQty - detail.Qty
						newWhStock.QtyBs = &newQty
						if *newWhStock.QtyBs < 0 {
							return errors.New("quantity bad stock must be greater than 0")
						}
					} else if *bppr.ItemCdn == 3 {
						oldQty := *whStock.QtyExp
						newQty := oldQty - detail.Qty
						newWhStock.QtyExp = &newQty
						if *newWhStock.QtyExp < 0 {
							return errors.New("quantity expired stock must be greater than 0")
						}
					}

					err = service.WhStockRepository.UpdateWhStockByWhIdAndProId(txCtx, whStock.CustID, *whStock.WhID, *whStock.ProID, newWhStock)
					if err != nil {
						return err
					}

					var stock model.Stock
					err = structs.Automapper(whStock, &stock)
					if err != nil {
						return err
					}
					tempCogs := float64(0)
					stock.Cogs = &tempCogs
					stock.ItemCnd = bppr.ItemCdn
					stock.Qty = &detail.Qty
					stock.TrCode = request.TrCode
					stock.TrNo = returnNo
					stock.StockDate = bpprModel.ReturnDate
					stock.WhIDFrom = bppr.WhID
					stock.UnitPrice = detail.UnitPrice1
					stock.RefDetId = &bpprDet.ID
					err = service.StockRepository.StoreStock(txCtx, &stock)
					if err != nil {
						return err
					}
				*/
			}
		}

		return nil
	})

	return err
}

func (service *bpprServiceImpl) Delete(custId string, bpprNo string, userId int64) (err error) {
	err = service.BpprRepository.Delete(custId, bpprNo, userId)
	if err != nil {
		return err
	}

	return err
}
