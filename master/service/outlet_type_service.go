package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type OutletTypeService interface {
	Detail(int64, string) (entity.OutletTypeResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.OutletTypeResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.OutletTypeLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateOutletTypeBody) (entity.OutletTypeResponse, error)
	Update(int, entity.UpdateOutletTypeRequest) error
	Delete(string, int, int64) error
}

func NewOutletTypeService(outletTypeRepository repository.OutletTypeRepository) *outletTypeServiceImpl {
	return &outletTypeServiceImpl{
		OutletTypeRepository: outletTypeRepository,
	}
}

type outletTypeServiceImpl struct {
	OutletTypeRepository repository.OutletTypeRepository
}

func (service *outletTypeServiceImpl) Detail(outletTypeId int64, custId string) (response entity.OutletTypeResponse, err error) {
	outletType, err := service.OutletTypeRepository.FindOneByOutletTypeIdAndCustId(outletTypeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(outletType, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *outletTypeServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletTypeResponse, total int, lastPage int, err error) {
	outletTypes, total, lastPage, err := service.OutletTypeRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletTypes {
		var vResp entity.OutletTypeResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletTypeServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletTypeLookupResponse, total int, lastPage int, err error) {
	outletTypes, total, lastPage, err := service.OutletTypeRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletTypes {
		var vResp entity.OutletTypeLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletTypeServiceImpl) Store(request entity.CreateOutletTypeBody) (response entity.OutletTypeResponse, err error) {

	// outletType_code & cust id validation, if err == nil, this means that code & cust id already exists
	outletType, err := service.OutletTypeRepository.FindOneByOutletTypeCodeAndCustId(request.OtTypeCode, request.CustId)
	if err == nil {
		return response, errors.New("ot_type_code: " + outletType.OtTypeCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	outletTypeData := model.OutletType{
		CustId:     request.CustId,
		OtTypeCode: request.OtTypeCode,
		OtTypeName: request.OtTypeName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	outletTypeId, err := service.OutletTypeRepository.Store(outletTypeData)
	if err != nil {
		return response, err
	}

	response.OtTypeId = outletTypeId

	return response, err
}

func (service *outletTypeServiceImpl) Update(outletTypeId int, request entity.UpdateOutletTypeRequest) (err error) {

	// outletType_code & cust id validation, if err == nil and params outletTypeId != outletType.Id, this means that code & cust id already exists
	outletType, err := service.OutletTypeRepository.FindOneByOutletTypeCodeAndCustId(request.OtTypeCode, request.CustId)
	if err == nil && outletType.OtTypeId != outletTypeId {
		return errors.New("ot_type_code: " + outletType.OtTypeCode + " is already exists")
	}

	err = service.OutletTypeRepository.Update(outletTypeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *outletTypeServiceImpl) Delete(custId string, outletTypeId int, userId int64) (err error) {

	err = service.OutletTypeRepository.Delete(custId, outletTypeId, userId)
	if err != nil {
		return err
	}

	return err
}
