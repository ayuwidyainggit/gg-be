package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type OutletGroupService interface {
	Detail(int64, string) (entity.OutletGroupResponse, error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.OutletGroupLookupResponse, total int, lastPage int, err error)
	List(entity.GeneralQueryFilter, string) (data []entity.OutletGroupResponse, total int, lastPage int, err error)
	Store(entity.CreateOutletGroupBody) (entity.OutletGroupResponse, error)
	Update(int, entity.UpdateOutletGroupRequest) error
	Delete(string, int, int64) error
}

func NewOutletGroupService(outletGroupRepository repository.OutletGroupRepository) *outletGroupServiceImpl {
	return &outletGroupServiceImpl{
		OutletGroupRepository: outletGroupRepository,
	}
}

type outletGroupServiceImpl struct {
	OutletGroupRepository repository.OutletGroupRepository
}

func (service *outletGroupServiceImpl) Detail(outletGroupId int64, custId string) (response entity.OutletGroupResponse, err error) {
	outletGroup, err := service.OutletGroupRepository.FindOneByOutletGroupIdAndCustId(outletGroupId, custId)
	if err != nil {
		return response, err
	}

	response.OtGrpId = outletGroup.OutletGroupId
	response.OtGrpCode = outletGroup.OutletGroupCode
	response.OtGrpName = outletGroup.OutletGroupName
	response.IsActive = outletGroup.IsActive
	response.UpdatedBy = outletGroup.UpdatedBy
	response.UpdatedAt = outletGroup.UpdatedAt

	return response, err
}

func (service *outletGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletGroupResponse, total int, lastPage int, err error) {
	outletGroups, total, lastPage, err := service.OutletGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletGroups {
		var vResp entity.OutletGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletGroupLookupResponse, total int, lastPage int, err error) {
	outletGroups, total, lastPage, err := service.OutletGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletGroups {
		var vResp entity.OutletGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletGroupServiceImpl) Store(request entity.CreateOutletGroupBody) (response entity.OutletGroupResponse, err error) {
	outletGroup, err := service.OutletGroupRepository.FindOneByOutletGroupCodeAndCustId(request.OtGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("ot_grp_code: " + outletGroup.OutletGroupCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	outletGroupData := model.OutletGroup{
		CustId:          request.CustId,
		OutletGroupCode: request.OtGrpCode,
		OutletGroupName: request.OtGrpName,
		IsActive:        request.IsActive,
		CreatedAt:       &timeNow,
		CreatedBy:       &request.CreatedBy,
		UpdatedAt:       &timeNow,
		UpdatedBy:       &request.CreatedBy,
	}

	outletGroupId, err := service.OutletGroupRepository.Store(outletGroupData)
	if err != nil {
		return response, err
	}

	response.OtGrpId = outletGroupId

	return response, err
}

func (service *outletGroupServiceImpl) Update(outletGroupId int, request entity.UpdateOutletGroupRequest) (err error) {

	// outletGroup_code & cust id validation, if err == nil and params outletGroupId != outletGroup.Id, this means that code & cust id already exists
	outletGroup, err := service.OutletGroupRepository.FindOneByOutletGroupCodeAndCustId(request.OtGrpCode, request.CustId)
	if err == nil && outletGroup.OutletGroupId != outletGroupId {
		return errors.New("ot_grp_code: " + outletGroup.OutletGroupCode + " is already exists")
	}

	err = service.OutletGroupRepository.Update(outletGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *outletGroupServiceImpl) Delete(custId string, outletGroupId int, userId int64) (err error) {

	err = service.OutletGroupRepository.Delete(custId, outletGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
