package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type DistrictService interface {
	Detail(int, string) (entity.DistrictResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.DistrictResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.DistrictLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateDistrictBody) (entity.DistrictResponse, error)
	Update(int, entity.UpdateDistrictRequest) error
	Delete(string, int, int64) error
}

func NewDistrictService(districtRepository repository.DistrictRepository) *districtServiceImpl {
	return &districtServiceImpl{
		DistrictRepository: districtRepository,
	}
}

type districtServiceImpl struct {
	DistrictRepository repository.DistrictRepository
}

func (service *districtServiceImpl) Detail(districtId int, custId string) (response entity.DistrictResponse, err error) {
	district, err := service.DistrictRepository.FindOneByDistrictIdAndCustId(districtId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(district, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *districtServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DistrictResponse, total int, lastPage int, err error) {
	districts, total, lastPage, err := service.DistrictRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range districts {
		var vResp entity.DistrictResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *districtServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DistrictLookupResponse, total int, lastPage int, err error) {
	districts, total, lastPage, err := service.DistrictRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range districts {
		var vResp entity.DistrictLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *districtServiceImpl) Store(request entity.CreateDistrictBody) (response entity.DistrictResponse, err error) {

	// district_code & cust id validation, if err == nil, this means that code & cust id already exists
	district, err := service.DistrictRepository.FindOneByDistrictCodeAndCustId(request.DistrictCode, request.CustId)
	if err == nil {
		return response, errors.New("district_code: " + district.DistrictCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	districtData := model.District{
		CustId:       request.CustId,
		DistrictCode: request.DistrictCode,
		DistrictName: request.DistrictName,
		IsActive:     true,
		CreatedAt:    &timeNow,
		CreatedBy:    &request.CreatedBy,
		UpdatedAt:    &timeNow,
		UpdatedBy:    &request.CreatedBy,
	}

	districtId, err := service.DistrictRepository.Store(districtData)
	if err != nil {
		return response, err
	}

	response.DistrictId = districtId

	return response, err
}

func (service *districtServiceImpl) Update(districtId int, request entity.UpdateDistrictRequest) (err error) {

	// district_code & cust id validation, if err == nil and params districtId != district.Id, this means that code & cust id already exists
	district, err := service.DistrictRepository.FindOneByDistrictCodeAndCustId(request.DistrictCode, request.CustId)
	if err == nil && district.DistrictId != districtId {
		return errors.New("district_code: " + district.DistrictCode + " is already exists")
	}

	err = service.DistrictRepository.Update(districtId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *districtServiceImpl) Delete(custId string, districtId int, userId int64) (err error) {

	err = service.DistrictRepository.Delete(custId, districtId, userId)
	if err != nil {
		return err
	}

	return err
}
