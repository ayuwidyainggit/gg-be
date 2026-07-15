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

type ArPayService interface {
	Store(request entity.CreateArPayBody) (err error)
	Detail(arPayNo string, custID string) (response entity.ArPayResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ArPayListResponse, total int64, lastPage int, err error)
	Delete(custId string, arPayNo string, userId int64) (err error)
	Update(arPayNo string, request entity.UpdateArPayBody) (err error)
}

type arPayServiceImpl struct {
	Repository  repository.ArPayRepository
	Transaction repository.Dbtransaction
}

func NewArPayService(repository repository.ArPayRepository, transaction repository.Dbtransaction) *arPayServiceImpl {
	return &arPayServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}
func (service *arPayServiceImpl) Store(request entity.CreateArPayBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ArPayDate != nil {
		arPayDate, err := str.DateStrToRfc3339String(*request.ArPayDate)
		if err != nil {
			return err
		}
		request.ArPayDate = &arPayDate
	}

	var ArPaymodel model.ArPay
	err = structs.Automapper(request, &ArPaymodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &ArPaymodel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var detModel model.ArPayDet

			err = structs.Automapper(Detail, &detModel)
			if err != nil {
				return err
			}
			detModel.CustID = request.CustID
			detModel.ArPayNo = ArPaymodel.ArPayNo
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
func (service *arPayServiceImpl) Detail(arPayNo string, custID string) (response entity.ArPayResponse, err error) {
	arPay, err := service.Repository.FindByNo(arPayNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(arPay, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(arPayNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ArPayDetResponse
	for _, detail := range Details {
		var detailData entity.ArPayDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		DetailsData = append(DetailsData, detailData)
	}

	arPayDate := arPay.ArPayDate.Format("2006-01-02")
	response.ArPayDate = &arPayDate

	response.Details = DetailsData
	return response, nil
}
func (service *arPayServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ArPayListResponse, total int64, lastPage int, err error) {
	arpays, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range arpays {
		var vResp entity.ArPayListResponse
		structs.Automapper(row, &vResp)
		if row.ArPayDate != nil {
			arPayDate := row.ArPayDate.Format("2006-01-02")
			vResp.ArPayDate = &arPayDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *arPayServiceImpl) Delete(custId string, arPayNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, arPayNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *arPayServiceImpl) Update(arPayNo string, request entity.UpdateArPayBody) (err error) {
	c := context.Background()

	if request.ArPayDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.ArPayDate != "" {
			arPayDate, err := str.DateStrToRfc3339String(*request.ArPayDate)
			if err != nil {
				return err
			}
			request.ArPayDate = &arPayDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ArPay
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
		err = service.Repository.Update(txCtx, arPayNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.ArPayDetID != nil {
				DetailIds = append(DetailIds, *detail.ArPayDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, arPayNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var arPayDetModel model.ArPayDet

			err = structs.Automapper(detail, &arPayDetModel)
			if err != nil {
				return err
			}
			arPayDetModel.CustID = request.CustID
			arPayDetModel.ArPayNo = arPayNo
			if detail.ArPayDetID == nil || *detail.ArPayDetID == 0 {
				arPayDetModel.ArPayDetID = nil
				err = service.Repository.StoreDetail(txCtx, &arPayDetModel)
				if err != nil {
					return err
				}
			} else {
				arPayDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &arPayDetModel)
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
