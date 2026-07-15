package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type CashService interface {
	Store(request entity.CreateCashTrBody) (err error)
	Detail(CashTrNo string, custID string) (response entity.CashTrResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.CashTrListResponse, total int64, lastPage int, err error)
	Delete(custId string, CashTrNo string, userId int64) (err error)
	Update(CashTrNo string, request entity.UpdateCashTrBody) (err error)
}

type CashServiceImpl struct {
	Repository  repository.CashRepository
	Transaction repository.Dbtransaction
}

func NewCashService(CashSoRepository repository.CashRepository, transaction repository.Dbtransaction) *CashServiceImpl {
	return &CashServiceImpl{
		Repository:  CashSoRepository,
		Transaction: transaction,
	}
}

func (service *CashServiceImpl) Store(request entity.CreateCashTrBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.CashTrDate != nil {
		CashDate, err := str.DateStrToRfc3339String(*request.CashTrDate)
		if err != nil {
			return err
		}
		request.CashTrDate = &CashDate
	}

	var CashModel model.CashTr
	err = structs.Automapper(request, &CashModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &CashModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var CashDetModel model.CashTrDet
			err = structs.Automapper(Detail, &CashDetModel)
			if err != nil {
				return err
			}
			CashDetModel.CashTrDetId = nil
			CashDetModel.CustID = request.CustId
			CashDetModel.CashTrNo = CashModel.CashTrNo
			err = service.Repository.StoreDetail(txCtx, &CashDetModel)
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

func (service *CashServiceImpl) Detail(CashTrNo string, custID string) (response entity.CashTrResponse, err error) {
	Cash, err := service.Repository.FindByNo(CashTrNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Cash, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(CashTrNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.CashTrDetResponse
	for _, detail := range Details {
		var detailData entity.CashTrDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		DetailsData = append(DetailsData, detailData)
	}
	if Cash.CashTrDate != nil {
		CashDate := Cash.CashTrDate.Format("2006-01-02")
		response.CashTrDate = &CashDate
	}

	response.Details = DetailsData
	return response, nil
}

func (service *CashServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.CashTrListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.CashTrListResponse
		structs.Automapper(row, &vResp)
		if row.CashTrDate != nil {
			CashDate := row.CashTrDate.Format("2006-01-02")
			vResp.CashTrDate = &CashDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *CashServiceImpl) Delete(custId string, CashTrNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, CashTrNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *CashServiceImpl) Update(CashTrNo string, request entity.UpdateCashTrBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.CashTrDate != nil {
		CashDate, err := str.DateStrToRfc3339String(*request.CashTrDate)
		if err != nil {
			return err
		}
		request.CashTrDate = &CashDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.CashTr
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, CashTrNo, Model)
		if err != nil {
			return err
		}
		detailIds := []int64{}

		for _, detail := range request.Details {
			if detail.CashTrDetId != nil {
				detailIds = append(detailIds, *detail.CashTrDetId)
			}
		}

		if len(detailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, CashTrNo, detailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339
			var CashDetModel model.CashTrDet
			err = structs.Automapper(detail, &CashDetModel)
			if err != nil {
				return err
			}
			CashDetModel.CustID = request.CustId
			CashDetModel.CashTrNo = CashTrNo
			if detail.CashTrDetId == nil || *detail.CashTrDetId == 0 {
				CashDetModel.CashTrDetId = nil
				err = service.Repository.StoreDetail(txCtx, &CashDetModel)
				if err != nil {
					return err
				}
			} else {
				CashDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &CashDetModel)
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
