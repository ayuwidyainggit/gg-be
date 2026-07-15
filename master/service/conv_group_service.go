package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ConvGroupService interface {
	Detail(int, string) (entity.ConvGroupResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.ConvGroupResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.ConvGroupLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateConvGroupBody) (entity.ConvGroupResponse, error)
	Update(int, entity.UpdateConvGroupRequest) error
	Delete(string, int, int64) error
}

func NewConvGroupService(convGroupRepository repository.ConvGroupRepository) *convGroupServiceImpl {
	return &convGroupServiceImpl{
		ConvGroupRepository: convGroupRepository,
	}
}

type convGroupServiceImpl struct {
	ConvGroupRepository repository.ConvGroupRepository
}

func (service *convGroupServiceImpl) Detail(convGroupId int, custId string) (response entity.ConvGroupResponse, err error) {
	convGroup, err := service.ConvGroupRepository.FindOneByConvGroupIdAndCustId(convGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(convGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *convGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ConvGroupResponse, total int, lastPage int, err error) {
	convGroups, total, lastPage, err := service.ConvGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range convGroups {
		var vResp entity.ConvGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *convGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ConvGroupLookupResponse, total int, lastPage int, err error) {
	convGroups, total, lastPage, err := service.ConvGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range convGroups {
		var vResp entity.ConvGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *convGroupServiceImpl) Store(request entity.CreateConvGroupBody) (response entity.ConvGroupResponse, err error) {

	// conv_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	convGroup, err := service.ConvGroupRepository.FindOneByConvGroupCodeAndCustId(request.ConvGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("conv_grp_code: " + convGroup.ConvGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	convGroupData := model.ConvGroup{
		CustId:      request.CustId,
		ConvGrpCode: request.ConvGrpCode,
		ConvGrpName: request.ConvGrpName,
		IsActive:    request.IsActive,
		CreatedAt:   &timeNow,
		CreatedBy:   &request.CreatedBy,
		UpdatedAt:   &timeNow,
		UpdatedBy:   &request.CreatedBy,
	}

	convGroupId, err := service.ConvGroupRepository.Store(convGroupData)
	if err != nil {
		return response, err
	}

	response.ConvGrpId = convGroupId

	return response, err
}

func (service *convGroupServiceImpl) Update(convGroupId int, request entity.UpdateConvGroupRequest) (err error) {

	// conv_grp_code & cust id validation, if err == nil and params convGroupId != convGroup.Id, this means that code & cust id already exists
	convGroup, err := service.ConvGroupRepository.FindOneByConvGroupCodeAndCustId(request.ConvGrpCode, request.CustId)
	if err == nil && convGroup.ConvGrpId != convGroupId {
		return errors.New("conv_grp_code: " + convGroup.ConvGrpCode + " is already exists")
	}

	err = service.ConvGroupRepository.Update(convGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *convGroupServiceImpl) Delete(custId string, convGroupId int, userId int64) (err error) {

	err = service.ConvGroupRepository.Delete(custId, convGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
