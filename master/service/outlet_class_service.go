package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type OutletClassService interface {
	Detail(int, string) (entity.OutletClassResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.OutletClassResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.OutletClassLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateOutletClassBody) (entity.OutletClassResponse, error)
	Update(int, entity.UpdateOutletClassRequest) error
	Delete(string, int, int64) error
}

func NewOutletClassService(outletClassRepository repository.OutletClassRepository) *outletClassServiceImpl {
	return &outletClassServiceImpl{
		OutletClassRepository: outletClassRepository,
	}
}

type outletClassServiceImpl struct {
	OutletClassRepository repository.OutletClassRepository
}

func (service *outletClassServiceImpl) Detail(outletClassId int, custId string) (response entity.OutletClassResponse, err error) {
	outletClass, err := service.OutletClassRepository.FindOneByOutletClassIdAndCustId(outletClassId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(outletClass, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *outletClassServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletClassResponse, total int, lastPage int, err error) {
	outletClasss, total, lastPage, err := service.OutletClassRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletClasss {
		var vResp entity.OutletClassResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletClassServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletClassLookupResponse, total int, lastPage int, err error) {
	outletClasss, total, lastPage, err := service.OutletClassRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletClasss {
		var vResp entity.OutletClassLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletClassServiceImpl) Store(request entity.CreateOutletClassBody) (response entity.OutletClassResponse, err error) {

	// outletClass_code & cust id validation, if err == nil, this means that code & cust id already exists
	outletClass, err := service.OutletClassRepository.FindOneByOutletClassCodeAndCustId(request.OtClassCode, request.CustId)
	if err == nil {
		return response, errors.New("ot_class_code: " + outletClass.OtClassCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	outletClassData := model.OutletClass{
		CustId:       request.CustId,
		OtClassCode:  request.OtClassCode,
		OtClassName:  request.OtClassName,
		OtClassLimit: request.OtClassLimit,
		IsActive:     request.IsActive,
		CreatedAt:    &timeNow,
		CreatedBy:    &request.CreatedBy,
		UpdatedAt:    &timeNow,
		UpdatedBy:    &request.CreatedBy,
	}

	outletClassId, err := service.OutletClassRepository.Store(outletClassData)
	if err != nil {
		return response, err
	}

	response.OtClassId = outletClassId

	return response, err
}

func (service *outletClassServiceImpl) Update(outletClassId int, request entity.UpdateOutletClassRequest) (err error) {

	// outletClass_code & cust id validation, if err == nil and params outletClassId != outletClass.Id, this means that code & cust id already exists
	outletClass, err := service.OutletClassRepository.FindOneByOutletClassCodeAndCustId(request.OtClassCode, request.CustId)
	if err == nil && outletClass.OtClassId != outletClassId {
		return errors.New("ot_class_code: " + outletClass.OtClassCode + " is already exists")
	}

	err = service.OutletClassRepository.Update(outletClassId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *outletClassServiceImpl) Delete(custId string, outletClassId int, userId int64) (err error) {

	err = service.OutletClassRepository.Delete(custId, outletClassId, userId)
	if err != nil {
		return err
	}

	return err
}
