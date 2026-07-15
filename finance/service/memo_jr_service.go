package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"time"
)

type MemoJrService interface {
	Store(request entity.CreateMemoJrBody) (err error)
	Detail(mjNo string, custID string) (response entity.MemoJrResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MemoJrListResponse, total int64, lastPage int, err error)
	Delete(custId string, mjNo string, userId int64) (err error)
	Update(mjNo string, request entity.UpdateMemoJrBody) (err error)
}

type memoJrServiceImpl struct {
	Repository  repository.MemoJrRepository
	Transaction repository.Dbtransaction
}

func NewMemoJrService(repository repository.MemoJrRepository, transaction repository.Dbtransaction) *memoJrServiceImpl {
	return &memoJrServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}
func (service *memoJrServiceImpl) Store(request entity.CreateMemoJrBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.MjDate != nil {
		MjDate, err := str.DateStrToRfc3339String(*request.MjDate)
		if err != nil {
			return err
		}
		request.MjDate = &MjDate
	}

	var memoModel model.MemoJr
	err = structs.Automapper(request, &memoModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &memoModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var detModel model.MemoJrDet

			err = structs.Automapper(Detail, &detModel)
			if err != nil {
				return err
			}
			detModel.CustID = request.CustID
			detModel.MjNo = memoModel.MjNo
			err = service.Repository.StoreDetail(txCtx, &detModel)
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
func (service *memoJrServiceImpl) Detail(mjNo string, custID string) (response entity.MemoJrResponse, err error) {
	arPay, err := service.Repository.FindByNo(mjNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(arPay, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(mjNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.MemoJrDetResponse
	for _, detail := range Details {
		var detailData entity.MemoJrDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		DetailsData = append(DetailsData, detailData)
	}

	MjDate := arPay.MjDate.Format("2006-01-02")
	response.MjDate = &MjDate

	response.Details = DetailsData
	return response, nil
}
func (service *memoJrServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MemoJrListResponse, total int64, lastPage int, err error) {
	arpays, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range arpays {
		var vResp entity.MemoJrListResponse
		structs.Automapper(row, &vResp)
		if row.MjDate != nil {
			MjDate := row.MjDate.Format("2006-01-02")
			vResp.MjDate = &MjDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *memoJrServiceImpl) Delete(custId string, mjNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, mjNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *memoJrServiceImpl) Update(mjNo string, request entity.UpdateMemoJrBody) (err error) {
	c := context.Background()

	if request.MjDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.MjDate != "" {
			MjDate, err := str.DateStrToRfc3339String(*request.MjDate)
			if err != nil {
				return err
			}
			request.MjDate = &MjDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MemoJr
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	if Model.IsPosted != nil {
		if *Model.IsPosted {
			now := time.Now()
			Model.PostedAt = &now
		}
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, mjNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.MemoJrDetID != nil {
				DetailIds = append(DetailIds, *detail.MemoJrDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, mjNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var memoJrDetModel model.MemoJrDet

			err = structs.Automapper(detail, &memoJrDetModel)
			if err != nil {
				return err
			}
			memoJrDetModel.CustID = request.CustID
			memoJrDetModel.MjNo = mjNo
			if detail.MemoJrDetID == nil || *detail.MemoJrDetID == 0 {
				memoJrDetModel.MemoJrDetID = nil
				err = service.Repository.StoreDetail(txCtx, &memoJrDetModel)
				if err != nil {
					return err
				}
			} else {
				memoJrDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &memoJrDetModel)
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
