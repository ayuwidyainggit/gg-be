package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type FlavorService interface {
	Detail(int, string) (entity.FlavorResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.FlavorListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.FlavorLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateFlavorBody) (entity.FlavorResponse, error)
	Update(int, entity.UpdateFlavorRequest) error
	Delete(string, int, int64) error
}

func NewFlavorService(flavorRepository repository.FlavorRepository, mProductRepository repository.ProductRepository) *flavorServiceImpl {
	return &flavorServiceImpl{
		FlavorRepository:   flavorRepository,
		MProductRepository: mProductRepository,
	}
}

type flavorServiceImpl struct {
	FlavorRepository   repository.FlavorRepository
	MProductRepository repository.ProductRepository
}

func (service *flavorServiceImpl) Detail(cProId int, custId string) (response entity.FlavorResponse, err error) {
	flavor, err := service.FlavorRepository.FindOneByFlavorIdAndCustId(cProId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(flavor, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *flavorServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.FlavorListResponse, total int, lastPage int, err error) {
	flavors, total, lastPage, err := service.FlavorRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range flavors {
		var vResp entity.FlavorListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *flavorServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.FlavorLookupResponse, total int, lastPage int, err error) {
	flavors, total, lastPage, err := service.FlavorRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range flavors {
		var vResp entity.FlavorLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *flavorServiceImpl) Store(request entity.CreateFlavorBody) (response entity.FlavorResponse, err error) {

	// flavor_code & cust id validation, if err == nil, this means that code & cust id already exists
	flavor, err := service.FlavorRepository.FindOneByFlavorCodeAndCustId(request.FlavorCode, request.CustId)
	if err == nil {
		return response, errors.New("flavor_code: " + flavor.FlavorCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.Flavor{
		CustId:     request.CustId,
		FlavorCode: request.FlavorCode,
		FlavorName: request.FlavorName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	cProId, err := service.FlavorRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.FlavorId = cProId

	return response, err
}

func (service *flavorServiceImpl) Update(cProId int, request entity.UpdateFlavorRequest) (err error) {

	// flavor_code & cust id validation, if err == nil and params cProId != flavor.Id, this means that code & cust id already exists
	flavor, err := service.FlavorRepository.FindOneByFlavorCodeAndCustId(request.FlavorCode, request.CustId)
	if err == nil && flavor.FlavorId != cProId {
		return errors.New("flavor_code: " + flavor.FlavorCode + " is already exists")
	}

	err = service.FlavorRepository.Update(cProId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *flavorServiceImpl) Delete(custId string, cProId int, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(cProId, custId, "pcat_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("pcat_id is still being used")
	}

	err = service.FlavorRepository.Delete(custId, cProId, userId)
	if err != nil {
		return err
	}

	return err
}
