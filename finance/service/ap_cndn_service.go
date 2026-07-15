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

type ApCndnService interface {
	Store(request entity.CreateApCndnBody) (err error)
	Detail(ApCndnNo string, custID string) (response entity.ApCndnResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ApCndnListResponse, total int64, lastPage int, err error)
	Delete(custId string, ApCndnNo string, userId int64) (err error)
	Update(ApCndnNo string, request entity.UpdateApCndnBody) (err error)
}

type ApCndnServiceImpl struct {
	Repository  repository.ApCndnRepository
	Transaction repository.Dbtransaction
}

func NewApCndnService(repository repository.ApCndnRepository, transaction repository.Dbtransaction) *ApCndnServiceImpl {
	return &ApCndnServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

func (service *ApCndnServiceImpl) Store(request entity.CreateApCndnBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ApCndnDate != nil {
		ApCndnDate, err := str.DateStrToRfc3339String(*request.ApCndnDate)
		if err != nil {
			return err
		}
		request.ApCndnDate = &ApCndnDate
	}

	var apCndnmodel model.ApCndn
	err = structs.Automapper(request, &apCndnmodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &apCndnmodel)
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

func (service *ApCndnServiceImpl) Detail(ApCndnNo string, custID string) (response entity.ApCndnResponse, err error) {
	apCndn, err := service.Repository.FindByNo(ApCndnNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(apCndn, &response)
	if err != nil {
		return response, err
	}

	apCndnDate := apCndn.ApCndnDate.Format("2006-01-02")
	response.ApCndnDate = &apCndnDate

	return response, nil
}
func (service *ApCndnServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ApCndnListResponse, total int64, lastPage int, err error) {
	apCndns, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range apCndns {
		var vResp entity.ApCndnListResponse
		structs.Automapper(row, &vResp)
		if row.ApCndnDate != nil {
			apCndnDate := row.ApCndnDate.Format("2006-01-02")
			vResp.ApCndnDate = &apCndnDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *ApCndnServiceImpl) Delete(custId string, ApCndnNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, ApCndnNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ApCndnServiceImpl) Update(ApCndnNo string, request entity.UpdateApCndnBody) (err error) {
	c := context.Background()

	if request.ApCndnDate != nil {
		apcndnDate, err := str.DateStrToRfc3339String(*request.ApCndnDate)
		if err != nil {
			return err
		}
		request.ApCndnDate = &apcndnDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ApCndn
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
		err = service.Repository.Update(txCtx, ApCndnNo, Model)
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
