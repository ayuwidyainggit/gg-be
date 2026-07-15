package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type DistPriceGroupService interface {
	Detail(int, string) (entity.DistPriceGroupResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.DistPriceGroupResponse, total int, lastPage int, err error)
	Store(entity.CreateDistPriceGroupBody) (entity.DistPriceGroupResponse, error)
	Update(int, entity.UpdateDistPriceGroupRequest) error
	Delete(string, int, int64) error
}

func NewDistPriceGroupService(DistPriceGroupRepository repository.DistPriceGroupRepository) *DistPriceGroupServiceImpl {
	return &DistPriceGroupServiceImpl{
		DistPriceGroupRepository: DistPriceGroupRepository,
	}
}

type DistPriceGroupServiceImpl struct {
	DistPriceGroupRepository repository.DistPriceGroupRepository
}

func (service *DistPriceGroupServiceImpl) Detail(DistPriceGroupId int, custId string) (response entity.DistPriceGroupResponse, err error) {
	distPriceGroup, err := service.DistPriceGroupRepository.FindOneByDistPriceGroupIdAndCustId(DistPriceGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(distPriceGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *DistPriceGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DistPriceGroupResponse, total int, lastPage int, err error) {
	distPriceGroups, total, lastPage, err := service.DistPriceGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distPriceGroups {
		var vResp entity.DistPriceGroupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DistPriceGroupServiceImpl) Store(request entity.CreateDistPriceGroupBody) (response entity.DistPriceGroupResponse, err error) {

	// disc_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	DistPriceGroup, err := service.DistPriceGroupRepository.FindOneByDistPriceGroupCodeAndCustId(request.DistPriceGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("disc_grp_code: " + DistPriceGroup.DistPriceGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	DistPriceGroupData := model.DistPriceGroup{
		CustId:           request.CustId,
		DistPriceGrpCode: request.DistPriceGrpCode,
		DistPriceGrpName: request.DistPriceGrpName,
		IsActive:         request.IsActive,
		CreatedAt:        &timeNow,
		CreatedBy:        &request.CreatedBy,
		UpdatedAt:        &timeNow,
		UpdatedBy:        &request.CreatedBy,
	}

	DistPriceGroupId, err := service.DistPriceGroupRepository.Store(DistPriceGroupData)
	if err != nil {
		return response, err
	}

	response.DistPriceGrpId = DistPriceGroupId

	return response, err
}

func (service *DistPriceGroupServiceImpl) Update(distPriceGroupId int, request entity.UpdateDistPriceGroupRequest) (err error) {

	// disc_grp_code & cust id validation, if err == nil and params DistPriceGroupId != DistPriceGroup.Id, this means that code & cust id already exists
	distPriceGroup, err := service.DistPriceGroupRepository.FindOneByDistPriceGroupCodeAndCustId(request.DistPriceGrpCode, request.CustId)
	if err == nil && distPriceGroup.DistPriceGrpId != distPriceGroupId {
		return errors.New("disc_grp_code: " + distPriceGroup.DistPriceGrpCode + " is already exists")
	}

	err = service.DistPriceGroupRepository.Update(distPriceGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *DistPriceGroupServiceImpl) Delete(custId string, DistPriceGroupId int, userId int64) (err error) {

	err = service.DistPriceGroupRepository.Delete(custId, DistPriceGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
