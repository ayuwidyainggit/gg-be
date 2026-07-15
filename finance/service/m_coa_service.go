package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"log"
)

type MCoaService interface {
	Store(request entity.CreateMCoaBody) (err error)
	Detail(MCoaId int64, custID string) (response entity.MCoaResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaListResponse, total int64, lastPage int, err error)
	LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaLookupListResponse, total int64, lastPage int, err error)
	Delete(custId string, MCoaId int64, userId int64) (err error)
	Update(mCoaId int64, request entity.UpdateMCoaBody, custId string) (err error)
}

type MCoaServiceImpl struct {
	MCoaRepository repository.MCoaRepository
	Transaction    repository.Dbtransaction
}

func NewMCoaService(repository repository.MCoaRepository, transaction repository.Dbtransaction) *MCoaServiceImpl {
	return &MCoaServiceImpl{
		MCoaRepository: repository,
		Transaction:    transaction,
	}
}
func (service *MCoaServiceImpl) Store(request entity.CreateMCoaBody) (err error) {
	c := context.Background()

	var mCoaModel model.MCoa
	err = structs.Automapper(request, &mCoaModel)
	if err != nil {
		return err
	}
	log.Println("mCoaModel:", structs.StructToJson(mCoaModel))
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.MCoaRepository.Store(txCtx, &mCoaModel)
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

func (service *MCoaServiceImpl) Detail(MCoaId int64, custID string) (response entity.MCoaResponse, err error) {
	MCoa, err := service.MCoaRepository.FindByID(MCoaId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(MCoa, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
func (service *MCoaServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaListResponse, total int64, lastPage int, err error) {
	MCoas, total, lastPage, err := service.MCoaRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MCoas {
		var vResp entity.MCoaListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MCoaServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MCoaLookupListResponse, total int64, lastPage int, err error) {
	mCoas, total, lastPage, err := service.MCoaRepository.FindAllByCustIdLookup(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range mCoas {
		var vResp entity.MCoaLookupListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MCoaServiceImpl) Delete(custId string, MCoaId int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MCoaRepository.Delete(txCtx, custId, MCoaId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *MCoaServiceImpl) Update(mCoaId int64, request entity.UpdateMCoaBody, custId string) (err error) {
	c := context.Background()

	// log.Println("request:", structs.StructToJson(request))
	// End parse time format YYYY-mm-dd to Rfc339
	var mCoa model.MCoa
	err = structs.Automapper(request, &mCoa)
	if err != nil {
		return err
	}
	mCoa.CustID = ""

	// log.Println("model.MCoa:", structs.StructToJson(mCoa))
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MCoaRepository.Update(txCtx, mCoaId, mCoa, custId)
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
