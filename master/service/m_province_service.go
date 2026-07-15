package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ProvinceService interface {
	Detail(string, string) (entity.ProvinceResponse, error)
	Store(entity.CreateProvinceBody) (entity.ProvinceResponse, error)
	Update(string, entity.UpdateProvinceRequest) error
	Delete(string, int64, string) error
	List(entity.GeneralQueryFilter, string) (data []entity.ProvinceResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.ProvinceLookupResponse, total int, lastPage int, err error)
}

func NewProvinceService(provinceRepository repository.ProvinceRepository) *provinceServiceImpl {
	return &provinceServiceImpl{
		ProvinceRepository: provinceRepository,
	}
}

type provinceServiceImpl struct {
	ProvinceRepository repository.ProvinceRepository
}

func (service *provinceServiceImpl) Detail(provinceId string, custId string) (response entity.ProvinceResponse, err error) {
	province, err := service.ProvinceRepository.FindOneByProvinceId(provinceId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(province, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *provinceServiceImpl) Store(request entity.CreateProvinceBody) (response entity.ProvinceResponse, err error) {

	// // province_code & cust id validation, if err == nil, this means that code & cust id already exists
	province, err := service.ProvinceRepository.FindOneByProvinceId(request.ProvinceId, request.CustId)
	if err == nil {
		return response, errors.New("province_id: " + province.ProvinceId + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	provinceData := model.Province{
		CustId:     request.CustId,
		ProvinceId: request.ProvinceId,
		Province:   request.Province,
		IsActive:   true,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	provinceId, err := service.ProvinceRepository.Store(provinceData)
	if err != nil {
		return response, err
	}

	response.ProvinceId = provinceId

	return response, err
}

func (service *provinceServiceImpl) Update(provinceId string, request entity.UpdateProvinceRequest) (err error) {

	province, err := service.ProvinceRepository.FindOneByProvinceId(request.ProvinceId, request.CustId)
	if err == nil && province.ProvinceId != provinceId {
		return errors.New("province_id: " + province.ProvinceId + " is already exists")
	}

	err = service.ProvinceRepository.Update(provinceId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *provinceServiceImpl) Delete(provinceId string, userId int64, custId string) (err error) {

	err = service.ProvinceRepository.Delete(provinceId, userId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *provinceServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ProvinceResponse, total int, lastPage int, err error) {
	provinces, total, lastPage, err := service.ProvinceRepository.FindAll(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range provinces {
		var vResp entity.ProvinceResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *provinceServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ProvinceLookupResponse, total int, lastPage int, err error) {
	provinces, total, lastPage, err := service.ProvinceRepository.FindAllLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range provinces {
		var vResp entity.ProvinceLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
