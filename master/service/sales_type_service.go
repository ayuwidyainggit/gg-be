package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SalesTypeService interface {
	Detail(int, string) (entity.SalesTypeResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.SalesTypeResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.SalesTypeLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateSalesTypeBody) (entity.SalesTypeResponse, error)
	Update(int, entity.UpdateSalesTypeRequest) error
	Delete(string, int, int64) error
}

func NewSalesTypeService(salesTypeRepository repository.SalesTypeRepository) *salesTypeServiceImpl {
	return &salesTypeServiceImpl{
		SalesTypeRepository: salesTypeRepository,
	}
}

type salesTypeServiceImpl struct {
	SalesTypeRepository repository.SalesTypeRepository
}

func (service *salesTypeServiceImpl) Detail(salesTypeId int, custId string) (response entity.SalesTypeResponse, err error) {
	salesType, err := service.SalesTypeRepository.FindOneBySalesTypeIdAndCustId(salesTypeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(salesType, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *salesTypeServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SalesTypeResponse, total int, lastPage int, err error) {
	salesTypes, total, lastPage, err := service.SalesTypeRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesTypes {
		var vResp entity.SalesTypeResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesTypeServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SalesTypeLookupResponse, total int, lastPage int, err error) {
	salesTypes, total, lastPage, err := service.SalesTypeRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesTypes {
		var vResp entity.SalesTypeLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesTypeServiceImpl) Store(request entity.CreateSalesTypeBody) (response entity.SalesTypeResponse, err error) {

	// salesType_code & cust id validation, if err == nil, this means that code & cust id already exists
	salesType, err := service.SalesTypeRepository.FindOneBySalesTypeCodeAndCustId(request.SalesTypeCode, request.CustId)
	if err == nil {
		return response, errors.New("salesType_code: " + salesType.SalesTypeCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	salesTypeData := model.SalesType{
		CustId:        request.CustId,
		SalesTypeCode: request.SalesTypeCode,
		SalesTypeName: request.SalesTypeName,
		IsActive:      request.IsActive,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	salesTypeId, err := service.SalesTypeRepository.Store(salesTypeData)
	if err != nil {
		return response, err
	}

	response.SalesTypeId = salesTypeId

	return response, err
}

func (service *salesTypeServiceImpl) Update(salesTypeId int, request entity.UpdateSalesTypeRequest) (err error) {

	// salesType_code & cust id validation, if err == nil and params salesTypeId != salesType.Id, this means that code & cust id already exists
	salesType, err := service.SalesTypeRepository.FindOneBySalesTypeCodeAndCustId(request.SalesTypeCode, request.CustId)
	if err == nil && salesType.SalesTypeId != salesTypeId {
		return errors.New("salesType_code: " + salesType.SalesTypeCode + " is already exists")
	}

	err = service.SalesTypeRepository.Update(salesTypeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *salesTypeServiceImpl) Delete(custId string, salesTypeId int, userId int64) (err error) {

	err = service.SalesTypeRepository.Delete(custId, salesTypeId, userId)
	if err != nil {
		return err
	}

	return err
}
