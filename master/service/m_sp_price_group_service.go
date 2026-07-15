package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SpecialPriceGroupService interface {
	Detail(int, string) (entity.SpecialPriceGroupResponse, error)
	LookupList(entity.SpecialPriceGroupQueryFilter, string) (data []entity.SpecialPriceGroupLookupResponse, total int, lastPage int, err error)
	List(entity.SpecialPriceGroupQueryFilter, string) (data []entity.SpecialPriceGroupResponse, total int, lastPage int, err error)
	Store(entity.CreateSpecialPriceGroupBody) (entity.SpecialPriceGroupResponse, error)
	Update(int, entity.UpdateSpecialPriceGroupRequest) error
	Delete(string, int, int64) error
}

func NewSpecialPriceGroupService(specialPriceGroupRepository repository.SpecialPriceGroupRepository) *SpecialPriceGroupServiceImpl {
	return &SpecialPriceGroupServiceImpl{
		SpecialPriceGroupRepository: specialPriceGroupRepository,
	}
}

type SpecialPriceGroupServiceImpl struct {
	SpecialPriceGroupRepository repository.SpecialPriceGroupRepository
}

func (service *SpecialPriceGroupServiceImpl) Detail(specialPriceGroupId int, custId string) (response entity.SpecialPriceGroupResponse, err error) {
	specialPriceGroup, err := service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupIdAndCustId(specialPriceGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(specialPriceGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *SpecialPriceGroupServiceImpl) List(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) (data []entity.SpecialPriceGroupResponse, total int, lastPage int, err error) {
	specialPriceGroups, total, lastPage, err := service.SpecialPriceGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range specialPriceGroups {
		var vResp entity.SpecialPriceGroupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SpecialPriceGroupServiceImpl) LookupList(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) (data []entity.SpecialPriceGroupLookupResponse, total int, lastPage int, err error) {
	specialPriceGroups, total, lastPage, err := service.SpecialPriceGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range specialPriceGroups {
		var vResp entity.SpecialPriceGroupLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SpecialPriceGroupServiceImpl) Store(request entity.CreateSpecialPriceGroupBody) (response entity.SpecialPriceGroupResponse, err error) {

	_, err = service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupCodeAndCustId(request.SpecialPriceGroupCode, request.CustId)
	if err == nil {
		return response, errors.New("Code : " + request.SpecialPriceGroupCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	specialPriceGroupData := model.SpecialPriceGroup{
		CustId:                request.CustId,
		SpecialPriceGroupCode: request.SpecialPriceGroupCode,
		SpecialPriceGroupName: request.SpecialPriceGroupName,
		IsActive:              request.IsActive,
		CreatedAt:             &timeNow,
		CreatedBy:             &request.CreatedBy,
		UpdatedAt:             &timeNow,
		UpdatedBy:             &request.CreatedBy,
	}

	specialPriceGroupId, err := service.SpecialPriceGroupRepository.Store(specialPriceGroupData)
	if err != nil {
		return response, err
	}

	response.SpecialPriceGroupId = specialPriceGroupId

	return response, err
}

func (service *SpecialPriceGroupServiceImpl) Update(specialPriceGroupId int, request entity.UpdateSpecialPriceGroupRequest) (err error) {
	specialPriceGroup, err := service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupCodeAndCustId(request.SpecialPriceGroupCode, request.CustId)
	if err == nil && specialPriceGroup.SpecialPriceGroupCode != request.SpecialPriceGroupCode {
		return errors.New("Code : " + request.SpecialPriceGroupCode + " is already exists")
	}

	err = service.SpecialPriceGroupRepository.Update(specialPriceGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *SpecialPriceGroupServiceImpl) Delete(custId string, specialPriceGroupId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(specialPriceGroupId, custId, "sp_price_grp_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("sp_price_grp_id is still being used")
	// }

	err = service.SpecialPriceGroupRepository.Delete(custId, specialPriceGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
