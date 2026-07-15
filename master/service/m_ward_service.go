package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type WardService interface {
	Detail(string, string) (entity.WardResponse, error)
	Store(entity.CreateWardBody) (entity.WardResponse, error)
	Update(string, entity.UpdateWardRequest) error
	Delete(string, int64, string) error
	List(entity.WardQueryFilter, string) (data []entity.WardResponse, total int, lastPage int, err error)
	LookupList(entity.WardQueryFilter, string) (data []entity.WardLookupResponse, total int, lastPage int, err error)
}

func NewWardService(wardRepository repository.WardRepository) *WardServiceImpl {
	return &WardServiceImpl{
		WardRepository: wardRepository,
	}
}

type WardServiceImpl struct {
	WardRepository repository.WardRepository
}

func (service *WardServiceImpl) Detail(ward string, custId string) (response entity.WardResponse, err error) {
	wards, err := service.WardRepository.FindOneByWardId(ward, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(wards, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *WardServiceImpl) Store(request entity.CreateWardBody) (response entity.WardResponse, err error) {

	// regency_code & cust id validation, if err == nil, this means that code & cust id already exists
	wards, err := service.WardRepository.FindOneByWardId(request.WardId, request.CustId)
	if err == nil {
		return response, errors.New("ward_id: " + wards.WardId + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	wardData := model.Ward{
		CustId:        request.CustId,
		WardId:        request.WardId,
		Ward:          request.Ward,
		ProvinceId:    request.ProvinceId,
		RegencyId:     request.RegencyId,
		SubDistrictId: request.SubDistrictId,
		IsActive:      true,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	wardId, err := service.WardRepository.Store(wardData)
	if err != nil {
		return response, err
	}

	response.WardId = wardId

	return response, err
}

func (service *WardServiceImpl) Update(wardId string, request entity.UpdateWardRequest) (err error) {

	// regency_code & cust id validation, if err == nil and params provinceId != province.Id, this means that code & cust id already exists
	wards, err := service.WardRepository.FindOneByWardId(request.WardId, request.CustId)
	if err == nil && wards.WardId != wardId {
		return errors.New("ward_id: " + wards.WardId + " is already exists")
	}

	err = service.WardRepository.Update(wardId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *WardServiceImpl) Delete(wardId string, userId int64, custId string) (err error) {

	err = service.WardRepository.Delete(wardId, userId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *WardServiceImpl) List(dataFilter entity.WardQueryFilter, custId string) (data []entity.WardResponse, total int, lastPage int, err error) {
	wards, total, lastPage, err := service.WardRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range wards {
		var vResp entity.WardResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *WardServiceImpl) LookupList(dataFilter entity.WardQueryFilter, custId string) (data []entity.WardLookupResponse, total int, lastPage int, err error) {
	wards, total, lastPage, err := service.WardRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range wards {
		var vResp entity.WardLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
