package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"
)

type MOpexService interface {
	Store(request entity.CreateMOpexBody) (err error)
	Detail(OpexId int, custID string) (response entity.MOpexResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MOpexListResponse, total int64, lastPage int, err error)
	Update(OpexId int, request entity.UpdateMOpexBody) (err error)
	Delete(custId string, OpexId int, userId int64) (err error)
}

type mOpexServiceImpl struct {
	MOpexSoRepository repository.MOpexRepository
	Transaction       repository.Dbtransaction
}

func NewMOpexService(mOpexSoRepository repository.MOpexRepository, transaction repository.Dbtransaction) *mOpexServiceImpl {
	return &mOpexServiceImpl{
		MOpexSoRepository: mOpexSoRepository,
		Transaction:       transaction,
	}
}

func (service *mOpexServiceImpl) Store(request entity.CreateMOpexBody) (err error) {
	c := context.Background()

	var mOpexModel model.MOpex
	err = structs.Automapper(request, &mOpexModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if mOpexModel.OpexCode != "" {
			opex, err := service.MOpexSoRepository.FindByCode(mOpexModel.OpexCode, request.CustId)
			if err != nil {

			}
			if opex.OpexId != nil {
				return errors.New(fmt.Sprintf("opex_code already exist with id %d", *opex.OpexId))
			}
		}
		err := service.MOpexSoRepository.Store(txCtx, &mOpexModel)
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

func (service *mOpexServiceImpl) Detail(OpexId int, custID string) (response entity.MOpexResponse, err error) {
	mOpex, err := service.MOpexSoRepository.FindById(OpexId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mOpex, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mOpexServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MOpexListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.MOpexSoRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.MOpexListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mOpexServiceImpl) Update(OpexId int, request entity.UpdateMOpexBody) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MOpex
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if Model.OpexCode != "" {
			opex, err := service.MOpexSoRepository.FindByCode(Model.OpexCode, request.CustId)
			if err != nil {

			}
			if opex.OpexId != nil {
				return errors.New(fmt.Sprintf("opex_code already exist with id %d", *opex.OpexId))
			}
		}
		err = service.MOpexSoRepository.Update(txCtx, OpexId, Model)
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

func (service *mOpexServiceImpl) Delete(custId string, OpexId int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MOpexSoRepository.Delete(txCtx, custId, OpexId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
