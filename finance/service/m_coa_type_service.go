package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
)

type MCoaTypeService interface {
	Store(request entity.CreateMCoaTypeBody) (err error)
	Detail(MCoaTypeId int64, custID string) (response entity.MCoaTypeResponse, err error)
	LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaTypeLookupListResponse, total int64, lastPage int, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaTypeListResponse, total int64, lastPage int, err error)
	Delete(MCoaTypeId int64, userId int64) (err error)
	Update(MCoaTypeId int64, request entity.UpdateMCoaTypeBody) (err error)
}

type MCoaTypeServiceImpl struct {
	MCoaTypeRepository repository.MCoaTypeRepository
	Transaction        repository.Dbtransaction
}

func NewMCoaTypeService(repository repository.MCoaTypeRepository, transaction repository.Dbtransaction) *MCoaTypeServiceImpl {
	return &MCoaTypeServiceImpl{
		MCoaTypeRepository: repository,
		Transaction:        transaction,
	}
}
func (service *MCoaTypeServiceImpl) Store(request entity.CreateMCoaTypeBody) (err error) {
	c := context.Background()

	var MCoaTypemodel model.MCoaType
	err = structs.Automapper(request, &MCoaTypemodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.MCoaTypeRepository.Store(txCtx, &MCoaTypemodel)
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

func (service *MCoaTypeServiceImpl) Detail(MCoaTypeId int64, custID string) (response entity.MCoaTypeResponse, err error) {
	MCoaType, err := service.MCoaTypeRepository.FindByID(MCoaTypeId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(MCoaType, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
func (service *MCoaTypeServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaTypeListResponse, total int64, lastPage int, err error) {
	MCoaTypes, total, lastPage, err := service.MCoaTypeRepository.FindAll(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MCoaTypes {
		var vResp entity.MCoaTypeListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MCoaTypeServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaTypeLookupListResponse, total int64, lastPage int, err error) {
	MCoaTypes, total, lastPage, err := service.MCoaTypeRepository.FindAllLookup(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MCoaTypes {
		var vResp entity.MCoaTypeLookupListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MCoaTypeServiceImpl) Delete(MCoaTypeId int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MCoaTypeRepository.Delete(txCtx, MCoaTypeId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *MCoaTypeServiceImpl) Update(MCoaTypeId int64, request entity.UpdateMCoaTypeBody) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MCoaType
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MCoaTypeRepository.Update(txCtx, MCoaTypeId, Model)
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
