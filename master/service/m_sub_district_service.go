package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SubDistrictService interface {
	Detail(string, string) (entity.SubDistrictResponse, error)
	Store(entity.CreateSubDistrictBody) (entity.SubDistrictResponse, error)
	Update(string, entity.UpdateSubDistrictRequest) error
	Delete(string, int64, string) error
	List(entity.SubDistrictQueryFilter, string) (data []entity.SubDistrictResponse, total int, lastPage int, err error)
	LookupList(entity.SubDistrictQueryFilter, string) (data []entity.SubDistrictLookupResponse, total int, lastPage int, err error)
}

func NewSubDistrictService(subDistrictRepository repository.SubDistrictRepository) *SubDistrictServiceImpl {
	return &SubDistrictServiceImpl{
		SubDistrictRepository: subDistrictRepository,
	}
}

type SubDistrictServiceImpl struct {
	SubDistrictRepository repository.SubDistrictRepository
}

func (service *SubDistrictServiceImpl) Detail(subDistrictId string, custId string) (response entity.SubDistrictResponse, err error) {
	subDistrict, err := service.SubDistrictRepository.FindOneBySubDistrictId(subDistrictId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(subDistrict, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *SubDistrictServiceImpl) Store(request entity.CreateSubDistrictBody) (response entity.SubDistrictResponse, err error) {

	subDistrict, err := service.SubDistrictRepository.FindOneBySubDistrictId(request.SubDistrictId, request.CustId)
	if err == nil {
		return response, errors.New("sub_district_id: " + subDistrict.SubDistrictId + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	subDistrictData := model.SubDistrict{
		CustId:        request.CustId,
		SubDistrictId: request.SubDistrictId,
		SubDistrict:   request.SubDistrict,
		ProvinceId:    request.ProvinceId,
		RegencyId:     request.RegencyId,
		IsActive:      true,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	subDistrictId, err := service.SubDistrictRepository.Store(subDistrictData)
	if err != nil {
		return response, err
	}

	response.SubDistrictId = subDistrictId

	return response, err
}

func (service *SubDistrictServiceImpl) Update(subDistrictId string, request entity.UpdateSubDistrictRequest) (err error) {

	subDistrict, err := service.SubDistrictRepository.FindOneBySubDistrictId(request.SubDistrictId, request.CustId)
	if err == nil && subDistrict.SubDistrictId != subDistrictId {
		return errors.New("sub_district_id: " + subDistrict.SubDistrictId + " is already exists")
	}

	err = service.SubDistrictRepository.Update(subDistrictId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *SubDistrictServiceImpl) Delete(subDistrictId string, userId int64, custId string) (err error) {

	err = service.SubDistrictRepository.Delete(subDistrictId, userId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *SubDistrictServiceImpl) List(dataFilter entity.SubDistrictQueryFilter, custId string) (data []entity.SubDistrictResponse, total int, lastPage int, err error) {
	regencys, total, lastPage, err := service.SubDistrictRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regencys {
		var vResp entity.SubDistrictResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SubDistrictServiceImpl) LookupList(dataFilter entity.SubDistrictQueryFilter, custId string) (data []entity.SubDistrictLookupResponse, total int, lastPage int, err error) {
	regencys, total, lastPage, err := service.SubDistrictRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regencys {
		var vResp entity.SubDistrictLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
