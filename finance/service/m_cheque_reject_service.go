package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
)

type MChequeRejectService interface {
	Store(request entity.CreateMChequeRejectBody) (err error)
	Detail(ChqRejectId int, custID string) (response entity.MChequeRejectResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MChequeRejectListResponse, total int64, lastPage int, err error)
	Update(ChqRejectId int, request entity.UpdateMChequeRejectBody) (err error)
	Delete(custId string, ChqRejectId int, userId int64) (err error)
}

type MChequeRejectServiceImpl struct {
	MChequeRejectSoRepository repository.MChequeRejectRepository
	Transaction               repository.Dbtransaction
}

func NewMChequeRejectService(MChequeRejectSoRepository repository.MChequeRejectRepository, transaction repository.Dbtransaction) *MChequeRejectServiceImpl {
	return &MChequeRejectServiceImpl{
		MChequeRejectSoRepository: MChequeRejectSoRepository,
		Transaction:               transaction,
	}
}

func (service *MChequeRejectServiceImpl) Store(request entity.CreateMChequeRejectBody) (err error) {
	c := context.Background()

	var MChequeRejectModel model.MChequeReject
	err = structs.Automapper(request, &MChequeRejectModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.MChequeRejectSoRepository.Store(txCtx, &MChequeRejectModel)
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

func (service *MChequeRejectServiceImpl) Detail(ChqRejectId int, custID string) (response entity.MChequeRejectResponse, err error) {
	MChequeReject, err := service.MChequeRejectSoRepository.FindById(ChqRejectId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(MChequeReject, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *MChequeRejectServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MChequeRejectListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.MChequeRejectSoRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(whAdjs) > 0 {
		for _, row := range whAdjs {
			var vResp entity.MChequeRejectListResponse
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *MChequeRejectServiceImpl) Update(ChqRejectId int, request entity.UpdateMChequeRejectBody) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MChequeReject
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MChequeRejectSoRepository.Update(txCtx, ChqRejectId, request.CustId, Model)
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

func (service *MChequeRejectServiceImpl) Delete(custId string, ChqRejectId int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MChequeRejectSoRepository.Delete(txCtx, custId, ChqRejectId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
