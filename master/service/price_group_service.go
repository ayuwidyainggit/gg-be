package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PriceGroupService interface {
	Detail(int, string) (entity.PriceGroupResponse, error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.PriceGroupLookupResponse, total int, lastPage int, err error)
	List(entity.GeneralQueryFilter, string) (data []entity.PriceGroupResponse, total int, lastPage int, err error)
	Store(entity.CreatePriceGroupBody) (entity.PriceGroupResponse, error)
	Update(int, entity.UpdatePriceGroupRequest) error
	Delete(string, int, int64) error
}

func NewPriceGroupService(priceGroupRepository repository.PriceGroupRepository) *priceGroupServiceImpl {
	return &priceGroupServiceImpl{
		PriceGroupRepository: priceGroupRepository,
	}
}

type priceGroupServiceImpl struct {
	PriceGroupRepository repository.PriceGroupRepository
}

func (service *priceGroupServiceImpl) Detail(priceGroupId int, custId string) (response entity.PriceGroupResponse, err error) {
	priceGroup, err := service.PriceGroupRepository.FindOneByPriceGroupIdAndCustId(priceGroupId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(priceGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *priceGroupServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PriceGroupResponse, total int, lastPage int, err error) {
	priceGroups, total, lastPage, err := service.PriceGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range priceGroups {
		var vResp entity.PriceGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *priceGroupServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PriceGroupLookupResponse, total int, lastPage int, err error) {
	priceGroups, total, lastPage, err := service.PriceGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range priceGroups {
		var vResp entity.PriceGroupLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *priceGroupServiceImpl) Store(request entity.CreatePriceGroupBody) (response entity.PriceGroupResponse, err error) {

	// price_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	priceGroup, err := service.PriceGroupRepository.FindOneByPriceGroupCodeAndCustId(request.PriceGrpCode, request.CustId)
	if err == nil {
		return response, errors.New("price_grp_code: " + priceGroup.PriceGrpCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	priceGroupData := model.PriceGroup{
		CustId:       request.CustId,
		PriceGrpCode: request.PriceGrpCode,
		PriceGrpName: request.PriceGrpName,
		IsActive:     request.IsActive,
		CreatedAt:    &timeNow,
		CreatedBy:    &request.CreatedBy,
		UpdatedAt:    &timeNow,
		UpdatedBy:    &request.CreatedBy,
	}

	priceGroupId, err := service.PriceGroupRepository.Store(priceGroupData)
	if err != nil {
		return response, err
	}

	response.PriceGrpId = priceGroupId

	return response, err
}

func (service *priceGroupServiceImpl) Update(priceGroupId int, request entity.UpdatePriceGroupRequest) (err error) {

	// price_grp_code & cust id validation, if err == nil and params priceGroupId != priceGroup.Id, this means that code & cust id already exists
	priceGroup, err := service.PriceGroupRepository.FindOneByPriceGroupCodeAndCustId(request.PriceGrpCode, request.CustId)
	if err == nil && priceGroup.PriceGrpId != priceGroupId {
		return errors.New("price_grp_code: " + priceGroup.PriceGrpCode + " is already exists")
	}

	err = service.PriceGroupRepository.Update(priceGroupId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *priceGroupServiceImpl) Delete(custId string, priceGroupId int, userId int64) (err error) {

	err = service.PriceGroupRepository.Delete(custId, priceGroupId, userId)
	if err != nil {
		return err
	}

	return err
}
