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

type ApPayService interface {
	Store(request entity.CreateApPayBody) (err error)
	Detail(apPayNo string, custID, parentCustId string) (response entity.ApPayRespone, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ApPayListResponse, total int64, lastPage int, err error)
	Delete(custId string, apPayNo string, userId int64) (err error)
	Update(apPayNo string, request entity.UpdateApPayBody) (err error)
}

type ApPayServiceImpl struct {
	Repository  repository.ApPayRepository
	Transaction repository.Dbtransaction
}

func NewApPayService(repository repository.ApPayRepository, transaction repository.Dbtransaction) *ApPayServiceImpl {
	return &ApPayServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}
func (service *ApPayServiceImpl) Store(request entity.CreateApPayBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ApPayDate != nil {
		ApPayDate, err := str.DateStrToRfc3339String(*request.ApPayDate)
		if err != nil {
			return err
		}
		request.ApPayDate = &ApPayDate
	}

	var Apmodel model.ApPay
	err = structs.Automapper(request, &Apmodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &Apmodel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var detModel model.ApPayDet
			err = structs.Automapper(Detail, &detModel)
			if err != nil {
				return err
			}
			detModel.CustID = request.CustID
			detModel.ApPayNo = Apmodel.ApPayNo
			err = service.Repository.StoreDetail(txCtx, &detModel)
			if err != nil {
				return err
			}
		}

		for _, MoneyPromoDetail := range request.ApPayMethodDetails {
			var apPayMethodModel model.ApPayMethod
			err = structs.Automapper(MoneyPromoDetail, &apPayMethodModel)
			if err != nil {
				return err
			}
			apPayMethodModel.CustID = request.CustID
			apPayMethodModel.ApPayNo = Apmodel.ApPayNo
			err = service.Repository.StoreApPaymethod(txCtx, &apPayMethodModel)
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

func (service *ApPayServiceImpl) Detail(apPayNo string, custID, parentCustId string) (response entity.ApPayRespone, err error) {
	ap, err := service.Repository.FindByNo(apPayNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ap, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(apPayNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ApPayDetResponse
	for _, detail := range Details {
		var detailData entity.ApPayDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		DetailsData = append(DetailsData, detailData)
	}

	ApPatMethodDetails, err := service.Repository.FindApPayMethod(apPayNo, custID)
	if err != nil {
		return response, err
	}
	var ApPayMethodData []entity.ApPayMethodRespone
	for _, ApPayMethodDetail := range ApPatMethodDetails {
		var detailData entity.ApPayMethodRespone
		err = structs.Automapper(ApPayMethodDetail, &detailData)
		if err != nil {
			return response, err
		}

		ApPayMethodData = append(ApPayMethodData, detailData)
	}

	ApPayDate := ap.ApPayDate.Format("2006-01-02")
	response.ApPayDate = &ApPayDate

	response.Details = DetailsData
	response.ApPayMethodDetails = ApPayMethodData
	return response, nil
}

func (service *ApPayServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ApPayListResponse, total int64, lastPage int, err error) {
	aps, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range aps {
		var vResp entity.ApPayListResponse
		structs.Automapper(row, &vResp)
		if row.ApPayDate != nil {
			ApPayDate := row.ApPayDate.Format("2006-01-02")
			vResp.ApPayDate = &ApPayDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ApPayServiceImpl) Delete(custId string, apPayNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, apPayNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ApPayServiceImpl) Update(apPayNo string, request entity.UpdateApPayBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ApPayDate != nil {
		ApPayDate, err := str.DateStrToRfc3339String(*request.ApPayDate)
		if err != nil {
			return err
		}
		request.ApPayDate = &ApPayDate
	}
	// End parse time format YYYY-mm-dd to Rfc339

	var Model model.ApPay
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
		err = service.Repository.Update(txCtx, apPayNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}
		ApPayMethodDetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.ApPayDetId != nil {
				DetailIds = append(DetailIds, *detail.ApPayDetId)
			}
		}
		for _, detail := range request.ApPayMethodDetails {
			if detail.ApPayMethodId != nil {
				ApPayMethodDetailIds = append(ApPayMethodDetailIds, *detail.ApPayMethodId)
			}
		}

		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, apPayNo, DetailIds)
			if err != nil {
				return err
			}
		}
		if len(ApPayMethodDetailIds) > 0 {
			err := service.Repository.DeleteApPayMethodDetailNotInIDs(txCtx, apPayNo, ApPayMethodDetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			var apDetModel model.ApPayDet

			err = structs.Automapper(detail, &apDetModel)
			if err != nil {
				return err
			}
			apDetModel.CustID = request.CustID
			apDetModel.ApPayNo = apPayNo
			if detail.ApPayDetId == nil || *detail.ApPayDetId == 0 {
				apDetModel.ApPayDetId = 0
				err = service.Repository.StoreDetail(txCtx, &apDetModel)
				if err != nil {
					return err
				}
			} else {
				apDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &apDetModel)
				if err != nil {
					return err
				}

			}
		}

		for _, detail := range request.ApPayMethodDetails {
			var apPayMethodModel model.ApPayMethod

			err = structs.Automapper(detail, &apPayMethodModel)
			if err != nil {
				return err
			}
			apPayMethodModel.CustID = request.CustID
			apPayMethodModel.ApPayNo = apPayNo
			if detail.ApPayMethodId == nil || *detail.ApPayMethodId == 0 {
				apPayMethodModel.ApPayMethodId = nil
				err = service.Repository.StoreApPaymethod(txCtx, &apPayMethodModel)
				if err != nil {
					return err
				}
			} else {
				apPayMethodModel.CustID = ""
				err = service.Repository.UpdateApPayMethodDetail(txCtx, &apPayMethodModel)
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
