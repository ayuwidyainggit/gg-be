package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PackTypeService interface {
	Detail(int, string) (entity.PackTypeResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.PackTypeListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.PackTypeLookupResponse, total int, lastPage int, err error)
	Store(entity.CreatePackTypeBody) (entity.PackTypeResponse, error)
	Update(int, entity.UpdatePackTypeRequest) error
	Delete(string, int, int64) error
}

func NewPackTypeService(packTypeRepository repository.PackTypeRepository, mProductRepository repository.ProductRepository) *packTypeServiceImpl {
	return &packTypeServiceImpl{
		PackTypeRepository: packTypeRepository,
		MProductRepository: mProductRepository,
	}
}

type packTypeServiceImpl struct {
	PackTypeRepository repository.PackTypeRepository
	MProductRepository repository.ProductRepository
}

func (service *packTypeServiceImpl) Detail(pTypeId int, custId string) (response entity.PackTypeResponse, err error) {
	packType, err := service.PackTypeRepository.FindOneByPTypeIdAndCustId(pTypeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(packType, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *packTypeServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PackTypeListResponse, total int, lastPage int, err error) {
	packTypes, total, lastPage, err := service.PackTypeRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range packTypes {
		var vResp entity.PackTypeListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *packTypeServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PackTypeLookupResponse, total int, lastPage int, err error) {
	packTypes, total, lastPage, err := service.PackTypeRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range packTypes {
		var vResp entity.PackTypeLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *packTypeServiceImpl) Store(request entity.CreatePackTypeBody) (response entity.PackTypeResponse, err error) {

	// ptype_code & cust id validation, if err == nil, this means that code & cust id already exists
	packType, err := service.PackTypeRepository.FindOneByPTypeCodeAndCustId(request.PtypeCode, request.CustId)
	if err == nil {
		return response, errors.New("ptype_code: " + packType.PtypeCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.PackType{
		CustId:    request.CustId,
		PtypeCode: request.PtypeCode,
		PtypeName: request.PtypeName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	pTypeId, err := service.PackTypeRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.PTypeId = pTypeId

	return response, err
}

func (service *packTypeServiceImpl) Update(pTypeId int, request entity.UpdatePackTypeRequest) (err error) {

	// ptype_code & cust id validation, if err == nil and params pTypeId != packType.Id, this means that code & cust id already exists
	packType, err := service.PackTypeRepository.FindOneByPTypeCodeAndCustId(request.PtypeCode, request.CustId)
	if err == nil && packType.PtypeId != pTypeId {
		return errors.New("ptype_code: " + packType.PtypeCode + " is already exists")
	}

	err = service.PackTypeRepository.Update(pTypeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *packTypeServiceImpl) Delete(custId string, pTypeId int, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(pTypeId, custId, "ptype_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("ptype_id is still being used")
	}

	err = service.PackTypeRepository.Delete(custId, pTypeId, userId)
	if err != nil {
		return err
	}

	return err
}
