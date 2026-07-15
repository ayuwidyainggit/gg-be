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

type ArCndnService interface {
	Store(request entity.CreateArCndnBody) (err error)
	Detail(ArCndnId int, custID, parentCustId string) (response entity.ArCndnResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ArCndnListResponse, total int64, lastPage int, err error)
	Delete(custId string, ArCndnId int, userId int64) (err error)
	Update(ArCndnId int, request entity.UpdateArCndnBody) (err error)
}

type ArCndnServiceImpl struct {
	ArCndnRepository repository.ArCndnRepository
	Transaction      repository.Dbtransaction
}

func NewArCndnService(repository repository.ArCndnRepository, transaction repository.Dbtransaction) *ArCndnServiceImpl {
	return &ArCndnServiceImpl{
		ArCndnRepository: repository,
		Transaction:      transaction,
	}
}

func (service *ArCndnServiceImpl) Store(request entity.CreateArCndnBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ArCndnDate != nil {
		ArCndnDate, err := str.DateStrToRfc3339String(*request.ArCndnDate)
		if err != nil {
			return err
		}
		request.ArCndnDate = &ArCndnDate
	}

	var ArCndnmodel model.ArCndn
	err = structs.Automapper(request, &ArCndnmodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.ArCndnRepository.Store(txCtx, &ArCndnmodel)
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

func (service *ArCndnServiceImpl) Detail(ArCndnId int, custID, parentCustId string) (response entity.ArCndnResponse, err error) {
	ArCndn, err := service.ArCndnRepository.FindById(ArCndnId, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ArCndn, &response)
	if err != nil {
		return response, err
	}

	ArCndnDate := ArCndn.ArCndnDate.Format("2006-01-02")
	response.ArCndnDate = &ArCndnDate

	return response, nil
}

func (service *ArCndnServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ArCndnListResponse, total int64, lastPage int, err error) {
	ArCndns, total, lastPage, err := service.ArCndnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ArCndns {
		var vResp entity.ArCndnListResponse
		structs.Automapper(row, &vResp)
		if row.ArCndnDate != nil {
			ArCndnDate := row.ArCndnDate.Format("2006-01-02")
			vResp.ArCndnDate = &ArCndnDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ArCndnServiceImpl) Delete(custId string, ArCndnId int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ArCndnRepository.Delete(txCtx, custId, ArCndnId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ArCndnServiceImpl) Update(ArCndnId int, request entity.UpdateArCndnBody) (err error) {
	c := context.Background()

	if request.ArCndnDate != nil {
		ArCndnDate, err := str.DateStrToRfc3339String(*request.ArCndnDate)
		if err != nil {
			return err
		}
		request.ArCndnDate = &ArCndnDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ArCndn
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
	Model.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ArCndnRepository.Update(txCtx, ArCndnId, Model)
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
