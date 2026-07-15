package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SubDistributorGroupService interface {
	Store(request entity.CreateSubDistributorGroupBody) (response entity.CreateSubDistributorGroupBody, err error)
	List(dataFilter entity.SubDistributorGroupQueryFilter, custId string) (data []entity.SubDistributorGroupListResponse, total int, lastPage int, err error)
	LookupList(dataFilter entity.SubDistributorGroupQueryFilter, custId string) (data []entity.SubDistributorGroupLookupResponse, total int, lastPage int, err error)
	Detail(channelID int, custId string) (response entity.SubDistributorGroupResponse, err error)
	Update(channelID int, request entity.SubDistributorGroupUpdateRequest) (err error)
	Delete(custId string, channelId int, userId int64) (err error)
}
type SubDistributorGroupServiceImpl struct {
	SubDistributorGroupRepository repository.SubDistributorGroupRepository
}

func NewSubDistributorGroupService(mSubDistributorGroupRepository repository.SubDistributorGroupRepository) *SubDistributorGroupServiceImpl {
	return &SubDistributorGroupServiceImpl{
		SubDistributorGroupRepository: mSubDistributorGroupRepository,
	}
}

func (service *SubDistributorGroupServiceImpl) Store(request entity.CreateSubDistributorGroupBody) (response entity.CreateSubDistributorGroupBody, err error) {
	_, err = service.SubDistributorGroupRepository.FindOneBySubDistributorGroupCodeAndCustId(request.SubDistributorGroupCode, request.CustID)
	if err == nil {
		return response, errors.New("sub_distributor_group_code: " + request.SubDistributorGroupCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	channel := model.SubDistributorGroup{
		CustID:                  request.CustID,
		SubDistributorGroupCode: request.SubDistributorGroupCode,
		SubDistributorGroupName: request.SubDistributorGroupName,
		IsActive:                request.IsActive,
		CreatedBy:               &request.CreatedBy,
		CreatedAt:               &timeNow,
		UpdatedBy:               &request.CreatedBy,
		UpdatedAt:               &timeNow,
	}

	id, err := service.SubDistributorGroupRepository.Store(channel)
	if err != nil {
		return response, err
	}

	response.SubDistributorGroupID = id
	return response, nil
}

func (service *SubDistributorGroupServiceImpl) List(dataFilter entity.SubDistributorGroupQueryFilter, custId string) (data []entity.SubDistributorGroupListResponse, total int, lastPage int, err error) {
	channels, total, lastPage, err := service.SubDistributorGroupRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range channels {
		var vResp entity.SubDistributorGroupListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SubDistributorGroupServiceImpl) LookupList(dataFilter entity.SubDistributorGroupQueryFilter, custId string) (data []entity.SubDistributorGroupLookupResponse, total int, lastPage int, err error) {
	channels, total, lastPage, err := service.SubDistributorGroupRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range channels {
		var vResp entity.SubDistributorGroupLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SubDistributorGroupServiceImpl) Detail(channelID int, custId string) (response entity.SubDistributorGroupResponse, err error) {
	brand, err := service.SubDistributorGroupRepository.FindOneBySubDistributorGroupIdAndCustId(channelID, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(brand, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *SubDistributorGroupServiceImpl) Update(channelID int, request entity.SubDistributorGroupUpdateRequest) (err error) {

	channel, err := service.SubDistributorGroupRepository.FindOneBySubDistributorGroupCodeAndCustId(request.SubDistributorGroupCode, request.CustID)
	if err == nil && channel.SubDistributorGroupID != int64(channelID) {
		return errors.New("sub_distributor_group_code: " + channel.SubDistributorGroupCode + " is already exists")
	}

	err = service.SubDistributorGroupRepository.Update(channelID, request)
	if err != nil {
		return err
	}

	return err
}

func (service *SubDistributorGroupServiceImpl) Delete(custId string, channelId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(brandId, custId, "brand_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("brand_id is still being used")
	// }

	err = service.SubDistributorGroupRepository.Delete(custId, channelId, userId)
	if err != nil {
		return err
	}

	return err
}
