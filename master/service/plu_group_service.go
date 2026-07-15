package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PluGroupService interface {
	Detail(int, string) (entity.PluGroupResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.PluGroupResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.PluGroupLookupResponse, total int, lastPage int, err error)
	Store(entity.CreatePluGroupBody) (entity.PluGroupResponse, error)
	Update(int, entity.UpdatePluGroupRequest) error
	Delete(string, int, int64) error
}

func NewPluGroupService(pluGroupRepository repository.PluGroupRepository) *pluGroupServiceImpl {
	return &pluGroupServiceImpl{
		PluGroupRepository: pluGroupRepository,
	}
}

type pluGroupServiceImpl struct {
	PluGroupRepository repository.PluGroupRepository
}

func (service *pluGroupServiceImpl) Detail(pluGroupId int, custId string) (response entity.PluGroupResponse, err error) {
	pluGroup, err := service.PluGroupRepository.FindOneByPluGroupIdAndCustId(pluGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(pluGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *pluGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PluGroupResponse, total int, lastPage int, err error) {
	pluGroups, total, lastPage, err := service.PluGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range pluGroups {
		var vResp entity.PluGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *pluGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PluGroupLookupResponse, total int, lastPage int, err error) {
	pluGroups, total, lastPage, err := service.PluGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range pluGroups {
		var vResp entity.PluGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *pluGroupServiceImpl) Store(request entity.CreatePluGroupBody) (response entity.PluGroupResponse, err error) {

	// plu_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	pluGroup, err := service.PluGroupRepository.FindOneByPluGroupCodeAndCustId(request.PluGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("plu_grp_code: " + pluGroup.PluGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	pluGroupData := model.PluGroup{
		CustId:     request.CustId,
		PluGrpCode: request.PluGrpCode,
		PluGrpName: request.PluGrpName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	pluGroupId, err := service.PluGroupRepository.Store(pluGroupData)
	if err != nil {
		return response, err
	}

	response.PluGrpId = pluGroupId

	return response, err
}

func (service *pluGroupServiceImpl) Update(pluGroupId int, request entity.UpdatePluGroupRequest) (err error) {

	// plu_grp_code & cust id validation, if err == nil and params pluGroupId != pluGroup.Id, this means that code & cust id already exists
	pluGroup, err := service.PluGroupRepository.FindOneByPluGroupCodeAndCustId(request.PluGrpCode, request.CustId)
	if err == nil && pluGroup.PluGrpId != pluGroupId {
		return errors.New("plu_grp_code: " + pluGroup.PluGrpCode + " is already exists")
	}

	err = service.PluGroupRepository.Update(pluGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *pluGroupServiceImpl) Delete(custId string, pluGroupId int, userId int64) (err error) {

	err = service.PluGroupRepository.Delete(custId, pluGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
