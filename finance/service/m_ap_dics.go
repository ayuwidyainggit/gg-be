package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type MApDiscService interface {
	Store(request entity.CreateMApDiscBody) (err error)
	Detail(apDiscID int64, custID string) (response entity.MApDiscResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MApDiscListResponse, total int64, lastPage int, err error)
	Delete(custId string, apDiscID int64, userId int64) (err error)
	Update(apDiscID int64, request entity.UpdateMApDiscBody) (err error)
}

type mApDiscServiceImpl struct {
	Repository  repository.MApDiscRepository
	Transaction repository.Dbtransaction
}

func NewMApDiscService(Repository repository.MApDiscRepository, transaction repository.Dbtransaction) *mApDiscServiceImpl {
	return &mApDiscServiceImpl{
		Repository:  Repository,
		Transaction: transaction,
	}
}
func (service *mApDiscServiceImpl) Store(request entity.CreateMApDiscBody) (err error) {
	c := context.Background()

	var apDiscModel model.MApDisc
	err = structs.Automapper(request, &apDiscModel)
	if err != nil {
		return err
	}

	mApDisc, err := service.Repository.FindByProId(*request.ProID, request.CustID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("mApDiscServiceImpl, Store, ErrRecordNotFound")
		} else {
			return err
		}
	}

	if mApDisc.ApDiscID > 0 {
		proIdStr := strconv.Itoa(int(*request.ProID))
		return errors.New("pro_id: " + proIdStr + " ( pro_code: " + *mApDisc.ProCode + " ) is already exists")
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &apDiscModel)
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
func (service *mApDiscServiceImpl) Detail(apDiscID int64, custID string) (response entity.MApDiscResponse, err error) {
	mApDisc, err := service.Repository.FindByNo(apDiscID, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mApDisc, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
func (service *mApDiscServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MApDiscListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.MApDiscListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *mApDiscServiceImpl) Delete(custId string, apDiscID int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, apDiscID, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *mApDiscServiceImpl) Update(apDiscID int64, request entity.UpdateMApDiscBody) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MApDisc
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, apDiscID, Model)
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
