package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type DiscGroupService interface {
	Detail(int, string) (entity.DiscGroupResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.DiscGroupResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.DiscGroupLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateDiscGroupBody) (entity.DiscGroupResponse, error)
	Update(int, entity.UpdateDiscGroupRequest) error
	Delete(string, int, int64) error
}

func NewDiscGroupService(discGroupRepository repository.DiscGroupRepository) *discGroupServiceImpl {
	return &discGroupServiceImpl{
		DiscGroupRepository: discGroupRepository,
	}
}

type discGroupServiceImpl struct {
	DiscGroupRepository repository.DiscGroupRepository
}

func (service *discGroupServiceImpl) Detail(discGroupId int, custId string) (response entity.DiscGroupResponse, err error) {
	discGroup, err := service.DiscGroupRepository.FindOneByDiscGroupIdAndCustId(discGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(discGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *discGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DiscGroupResponse, total int, lastPage int, err error) {
	discGroups, total, lastPage, err := service.DiscGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range discGroups {
		var vResp entity.DiscGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *discGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DiscGroupLookupResponse, total int, lastPage int, err error) {
	discGroups, total, lastPage, err := service.DiscGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range discGroups {
		var vResp entity.DiscGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *discGroupServiceImpl) Store(request entity.CreateDiscGroupBody) (response entity.DiscGroupResponse, err error) {

	// disc_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	discGroup, err := service.DiscGroupRepository.FindOneByDiscGroupCodeAndCustId(request.DiscGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("disc_grp_code: " + discGroup.DiscGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	discGroupData := model.DiscGroup{
		CustId:        request.CustId,
		DiscGrpCode: request.DiscGrpCode,
		DiscGrpName: request.DiscGrpName,
		IsActive:      request.IsActive,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	discGroupId, err := service.DiscGroupRepository.Store(discGroupData)
	if err != nil {
		return response, err
	}

	response.DiscGrpId = discGroupId

	return response, err
}

func (service *discGroupServiceImpl) Update(discGroupId int, request entity.UpdateDiscGroupRequest) (err error) {

	// disc_grp_code & cust id validation, if err == nil and params discGroupId != discGroup.Id, this means that code & cust id already exists
	discGroup, err := service.DiscGroupRepository.FindOneByDiscGroupCodeAndCustId(request.DiscGrpCode, request.CustId)
	if err == nil && discGroup.DiscGrpId != discGroupId {
		return errors.New("disc_grp_code: " + discGroup.DiscGrpCode + " is already exists")
	}

	err = service.DiscGroupRepository.Update(discGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *discGroupServiceImpl) Delete(custId string, discGroupId int, userId int64) (err error) {

	err = service.DiscGroupRepository.Delete(custId, discGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
