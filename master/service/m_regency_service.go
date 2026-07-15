package service

import (
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type RegencyService interface {
	Detail(string, string) (entity.RegencyResponse, error)
	Store(entity.CreateRegencyBody) (entity.RegencyResponse, error)
	Update(string, entity.UpdateRegencyRequest) error
	Delete(string, int64, string) error
	List(entity.RegencyQueryFilter, string) (data []entity.RegencyResponse, total int, lastPage int, err error)
	LookupList(entity.RegencyQueryFilter, string) (data []entity.RegencyLookupResponse, total int, lastPage int, err error)
}

func NewRegencyService(regencyRepository repository.RegencyRepository) *regencyServiceImpl {
	return &regencyServiceImpl{
		RegencyRepository: regencyRepository,
	}
}

type regencyServiceImpl struct {
	RegencyRepository repository.RegencyRepository
}

func (service *regencyServiceImpl) Detail(regencyId string, custId string) (response entity.RegencyResponse, err error) {
	regency, err := service.RegencyRepository.FindOneByRegencyId(regencyId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(regency, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *regencyServiceImpl) Store(request entity.CreateRegencyBody) (response entity.RegencyResponse, err error) {

	regency, err := service.RegencyRepository.FindOneByRegencyId(request.RegencyId, request.CustId)
	if err == nil {
		return response, errors.New("regency_id: " + regency.RegencyId + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	provinceData := model.Regency{
		CustId:     request.CustId,
		RegencyId:  request.RegencyId,
		Regency:    request.Regency,
		ProvinceId: request.ProvinceId,
		IsActive:   true,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	regencyId, err := service.RegencyRepository.Store(provinceData)
	if err != nil {
		return response, err
	}

	response.RegencyId = regencyId

	return response, err
}

func (service *regencyServiceImpl) Update(regencyId string, request entity.UpdateRegencyRequest) (err error) {

	regency, err := service.RegencyRepository.FindOneByRegencyId(request.RegencyId, request.CustId)
	if err == nil && regency.RegencyId != regencyId {
		return errors.New("regency_id: " + regency.RegencyId + " is already exists")
	}

	err = service.RegencyRepository.Update(regencyId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *regencyServiceImpl) Delete(regencyId string, userId int64, custId string) (err error) {

	err = service.RegencyRepository.Delete(regencyId, userId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *regencyServiceImpl) List(dataFilter entity.RegencyQueryFilter, custId string) (data []entity.RegencyResponse, total int, lastPage int, err error) {
	regencys, total, lastPage, err := service.RegencyRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regencys {
		var vResp entity.RegencyResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *regencyServiceImpl) LookupList(dataFilter entity.RegencyQueryFilter, custId string) (data []entity.RegencyLookupResponse, total int, lastPage int, err error) {
	fmt.Println("FILTER-FILTER ====>", dataFilter)
	regencys, total, lastPage, err := service.RegencyRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regencys {
		var vResp entity.RegencyLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
