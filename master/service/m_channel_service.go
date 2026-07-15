package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type MChannelService interface {
	Store(request entity.CreateChannelBody) (response entity.CreateChannelBody, err error)
	List(dataFilter entity.ChannelQueryFilter, custId string) (data []entity.ChannelListResponse, total int, lastPage int, err error)
	LookupList(dataFilter entity.ChannelQueryFilter, custId string) (data []entity.ChannelLookupResponse, total int, lastPage int, err error)
	Detail(channelID int, custId string) (response entity.ChannelResponse, err error)
	Update(channelID int, request entity.ChannelUpdateRequest) (err error)
	Delete(custId string, channelId int, userId int64) (err error)
}
type MChannelServiceImpl struct {
	MChannelRepository repository.MChannelRepository
}

func NewMChannelService(mChannelRepository repository.MChannelRepository) *MChannelServiceImpl {
	return &MChannelServiceImpl{
		MChannelRepository: mChannelRepository,
	}
}

func (service *MChannelServiceImpl) Store(request entity.CreateChannelBody) (response entity.CreateChannelBody, err error) {
	_, err = service.MChannelRepository.FindOneByChannelCodeAndCustId(request.ChannelCode, request.CustID)
	if err == nil {
		return response, errors.New("channel_code: " + request.ChannelCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	channel := model.MChannel{
		CustID:      request.CustID,
		ChannelCode: request.ChannelCode,
		ChannelName: request.ChannelName,
		IsActive:    true,
		CreatedBy:   &request.CreatedBy,
		CreatedAt:   &timeNow,
		UpdatedBy:   &request.CreatedBy,
		UpdatedAt:   &timeNow,
	}

	id, err := service.MChannelRepository.Store(channel)
	if err != nil {
		return response, err
	}

	response.ChannelID = id
	return response, nil
}

func (service *MChannelServiceImpl) List(dataFilter entity.ChannelQueryFilter, custId string) (data []entity.ChannelListResponse, total int, lastPage int, err error) {
	channels, total, lastPage, err := service.MChannelRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range channels {
		var vResp entity.ChannelListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MChannelServiceImpl) LookupList(dataFilter entity.ChannelQueryFilter, custId string) (data []entity.ChannelLookupResponse, total int, lastPage int, err error) {
	channels, total, lastPage, err := service.MChannelRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range channels {
		var vResp entity.ChannelLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MChannelServiceImpl) Detail(channelID int, custId string) (response entity.ChannelResponse, err error) {
	brand, err := service.MChannelRepository.FindOneByChannelIdAndCustId(channelID, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(brand, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *MChannelServiceImpl) Update(channelID int, request entity.ChannelUpdateRequest) (err error) {

	channel, err := service.MChannelRepository.FindOneByChannelCodeAndCustId(request.ChannelCode, request.CustID)
	if err == nil && channel.ChannelID != int64(channelID) {
		return errors.New("channel_code: " + channel.ChannelCode + " is already exists")
	}

	err = service.MChannelRepository.Update(channelID, request)
	if err != nil {
		return err
	}

	return err
}

func (service *MChannelServiceImpl) Delete(custId string, channelId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(brandId, custId, "brand_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("brand_id is still being used")
	// }

	err = service.MChannelRepository.Delete(custId, channelId, userId)
	if err != nil {
		return err
	}

	return err
}
