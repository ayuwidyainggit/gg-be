package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type IncentiveGroupService interface {
	Detail(int, string) (entity.IncentiveGroupResponse, error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.IncentiveGroupLookupResponse, total int, lastPage int, err error)
	List(entity.GeneralQueryFilter, string) (data []entity.IncentiveGroupResponse, total int, lastPage int, err error)
	Store(entity.CreateIncentiveGroupBody) (entity.IncentiveGroupResponse, error)
	Update(int, entity.UpdateIncentiveGroupRequest) error
	Delete(string, int, int64) error
}

func NewIncentiveGroupService(incentiveGroupRepository repository.IncentiveGroupRepository) *incentiveGroupServiceImpl {
	return &incentiveGroupServiceImpl{
		IncentiveGroupRepository: incentiveGroupRepository,
	}
}

type incentiveGroupServiceImpl struct {
	IncentiveGroupRepository repository.IncentiveGroupRepository
}

func (service *incentiveGroupServiceImpl) Detail(incentiveGroupId int, custId string) (response entity.IncentiveGroupResponse, err error) {
	incentiveGroup, err := service.IncentiveGroupRepository.FindOneByIncentiveGroupIdAndCustId(incentiveGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(incentiveGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *incentiveGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.IncentiveGroupLookupResponse, total int, lastPage int, err error) {
	incentiveGroups, total, lastPage, err := service.IncentiveGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range incentiveGroups {
		var vResp entity.IncentiveGroupLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *incentiveGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.IncentiveGroupResponse, total int, lastPage int, err error) {
	incentiveGroups, total, lastPage, err := service.IncentiveGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range incentiveGroups {
		var vResp entity.IncentiveGroupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *incentiveGroupServiceImpl) Store(request entity.CreateIncentiveGroupBody) (response entity.IncentiveGroupResponse, err error) {
	incentiveGroup, err := service.IncentiveGroupRepository.FindOneByIncentiveGroupCodeAndCustId(request.IncGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("incentiveGroup_code: " + incentiveGroup.IncGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	incentiveGroupData := model.IncentiveGroup{
		CustID:     request.CustId,
		IncGrpCode: request.IncGrpCode,
		IncGrpName: request.IncGrpName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	incentiveGroupId, err := service.IncentiveGroupRepository.Store(incentiveGroupData)
	if err != nil {
		return response, err
	}

	response.IncGrpID = incentiveGroupId

	return response, err
}

func (service *incentiveGroupServiceImpl) Update(incentiveGroupId int, request entity.UpdateIncentiveGroupRequest) (err error) {

	incentiveGroup, err := service.IncentiveGroupRepository.FindOneByIncentiveGroupCodeAndCustId(request.IncGrpCode, request.CustId)
	if err == nil && incentiveGroup.IncGrpID != incentiveGroupId {
		return errors.New("incentiveGroup_code: " + incentiveGroup.IncGrpCode + " is already exists")
	}

	err = service.IncentiveGroupRepository.Update(incentiveGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *incentiveGroupServiceImpl) Delete(custId string, incentiveGroupId int, userId int64) (err error) {

	err = service.IncentiveGroupRepository.Delete(custId, incentiveGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
