package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type IndustryService interface {
	Detail(int, string) (entity.IndustryResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.IndustryResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.IndustryLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateIndustryBody) (entity.IndustryResponse, error)
	Update(int, entity.UpdateIndustryRequest) error
	Delete(string, int, int64) error
}

func NewIndustryService(industryRepository repository.IndustryRepository) *industryServiceImpl {
	return &industryServiceImpl{
		IndustryRepository: industryRepository,
	}
}

type industryServiceImpl struct {
	IndustryRepository repository.IndustryRepository
}

func (service *industryServiceImpl) Detail(industryId int, custId string) (response entity.IndustryResponse, err error) {
	industry, err := service.IndustryRepository.FindOneByIndustryIdAndCustId(industryId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(industry, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *industryServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.IndustryResponse, total int, lastPage int, err error) {
	industries, total, lastPage, err := service.IndustryRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range industries {
		var vResp entity.IndustryResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *industryServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.IndustryLookupResponse, total int, lastPage int, err error) {
	industries, total, lastPage, err := service.IndustryRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range industries {
		var vResp entity.IndustryLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *industryServiceImpl) Store(request entity.CreateIndustryBody) (response entity.IndustryResponse, err error) {

	// industry_code & cust id validation, if err == nil, this means that code & cust id already exists
	industry, err := service.IndustryRepository.FindOneByIndustryCodeAndCustId(request.IndustryCode, request.CustId)
	if err == nil {
		return response, errors.New("industry_code: " + industry.IndustryCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	industryData := model.Industry{
		CustId:       request.CustId,
		IndustryCode: request.IndustryCode,
		IndustryName: request.IndustryName,
		IsActive:     request.IsActive,
		CreatedAt:    &timeNow,
		CreatedBy:    &request.CreatedBy,
		UpdatedAt:    &timeNow,
		UpdatedBy:    &request.CreatedBy,
	}

	industryId, err := service.IndustryRepository.Store(industryData)
	if err != nil {
		return response, err
	}

	response.IndustryId = industryId

	return response, err
}

func (service *industryServiceImpl) Update(industryId int, request entity.UpdateIndustryRequest) (err error) {

	// industry_code & cust id validation, if err == nil and params industryId != industry.Id, this means that code & cust id already exists
	industry, err := service.IndustryRepository.FindOneByIndustryCodeAndCustId(request.IndustryCode, request.CustId)
	if err == nil && industry.IndustryId != industryId {
		return errors.New("industry_code: " + industry.IndustryCode + " is already exists")
	}

	err = service.IndustryRepository.Update(industryId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *industryServiceImpl) Delete(custId string, industryId int, userId int64) (err error) {

	err = service.IndustryRepository.Delete(custId, industryId, userId)
	if err != nil {
		return err
	}

	return err
}
