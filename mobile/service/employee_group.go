package service

import (
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/structs"
	"mobile/repository"
	"time"
)

type EmpGroupService interface {
	Detail(int, string) (entity.EmpGroupResponse, error)
	TypeDetail(string, string) (entity.EmpTypeResponse, error)
	List(entity.GeneralQueryFilter) (data []entity.EmpGroupResponse, total int64, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter) (data []entity.EmpGroupLookupResponse, total int64, lastPage int, err error)
	EmpTypeList(entity.GeneralQueryFilter, string) (data []entity.EmpTypeResponse, total int64, lastPage int, err error)
	Store(entity.CreateEmpGroupBody) (entity.EmpGroupResponse, error)
	Update(int, entity.UpdateEmpGroupRequest) error
	Delete(string, int, int64) error
}

func NewEmpGroupService(empGroupRepository repository.EmpGroupRepository) *empGroupServiceImpl {
	return &empGroupServiceImpl{
		EmpGroupRepository: empGroupRepository,
	}
}

type empGroupServiceImpl struct {
	EmpGroupRepository repository.EmpGroupRepository
}

func (service *empGroupServiceImpl) EmpTypeList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.EmpTypeResponse, total int64, lastPage int, err error) {
	empTypes, total, lastPage, err := service.EmpGroupRepository.EmpTypeFindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range empTypes {
		var vResp entity.EmpTypeResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *empGroupServiceImpl) Detail(empGroupId int, custId string) (response entity.EmpGroupResponse, err error) {
	empGroup, err := service.EmpGroupRepository.FindOneByEmpGroupIdAndCustId(empGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(empGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *empGroupServiceImpl) TypeDetail(empTypeId string, custId string) (response entity.EmpTypeResponse, err error) {
	empType, err := service.EmpGroupRepository.FindOneByEmpTypeIdAndCustId(empTypeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(empType, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *empGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.EmpGroupResponse, total int64, lastPage int, err error) {
	var empGroups []model.EmpGroup

	empGroups, total, lastPage, err = service.EmpGroupRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range empGroups {
		var vResp entity.EmpGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err

}

func (service *empGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmpGroupLookupResponse, total int64, lastPage int, err error) {
	var empGroups []model.EmpGroup

	empGroups, total, lastPage, err = service.EmpGroupRepository.FindAllByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range empGroups {
		var vResp entity.EmpGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err

}

func (service *empGroupServiceImpl) Store(request entity.CreateEmpGroupBody) (response entity.EmpGroupResponse, err error) {

	// emp_group_code & cust id validation, if err == nil, this means that code & cust id already exists
	empGroup, err := service.EmpGroupRepository.FindOneByEmpGroupCodeAndCustId(request.EmpGroupCode, request.CustId)
	if err == nil {
		return response, errors.New("emp_group_code: " + empGroup.EmpGroupCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	empGroupData := model.EmpGroup{
		CustId:       request.CustId,
		EmpGroupCode: request.EmpGroupCode,
		EmpGroupName: request.EmpGroupName,
		IsActive:     request.IsActive,
		CreatedAt:    &timeNow,
		CreatedBy:    &request.CreatedBy,
		UpdatedAt:    &timeNow,
		UpdatedBy:    &request.CreatedBy,
	}

	empGroupId, err := service.EmpGroupRepository.Store(empGroupData)
	if err != nil {
		return response, err
	}

	response.EmpGroupId = empGroupId

	return response, err
}

func (service *empGroupServiceImpl) Update(empGroupId int, request entity.UpdateEmpGroupRequest) (err error) {

	// emp_group_code & cust id validation, if err == nil and params empGroupId != empGroup.Id, this means that code & cust id already exists
	empGroup, err := service.EmpGroupRepository.FindOneByEmpGroupCodeAndCustId(request.EmpGroupCode, request.CustId)
	if err == nil && empGroup.EmpGroupId != empGroupId {
		return errors.New("emp_group_code: " + empGroup.EmpGroupCode + " is already exists")
	}

	err = service.EmpGroupRepository.Update(empGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *empGroupServiceImpl) Delete(custId string, empGroupId int, userId int64) (err error) {

	err = service.EmpGroupRepository.Delete(custId, empGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
