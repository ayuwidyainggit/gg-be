package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type RegionService interface {
	Detail(int, string) (entity.RegionResponse, error)
	LookupList(entity.RegionQueryFilter) (data []entity.RegionLookupResponse, total int, lastPage int, err error)
	List(entity.RegionQueryFilter) (data []entity.RegionResponse, total int, lastPage int, err error)
	Store(entity.CreateRegionBody) (entity.RegionResponse, error)
	Update(int, entity.UpdateRegionRequest) error
	Delete(string, int, int64) error
}

func NewRegionService(regionRepository repository.RegionRepository, employeeScopeRepo repository.EmployeeScopeRepository) *regionServiceImpl {
	return &regionServiceImpl{
		RegionRepository:   regionRepository,
		EmployeeScopeRepo: employeeScopeRepo,
	}
}

type regionServiceImpl struct {
	RegionRepository   repository.RegionRepository
	EmployeeScopeRepo repository.EmployeeScopeRepository
}

func (service *regionServiceImpl) Detail(regionId int, custId string) (response entity.RegionResponse, err error) {
	region, err := service.RegionRepository.FindOneByRegionIdAndCustId(regionId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(region, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *regionServiceImpl) Store(request entity.CreateRegionBody) (response entity.RegionResponse, err error) {
	// bank_code & cust id validation, if err == nil, this means that code & cust id already exists
	region, err := service.RegionRepository.FindOneByRegionCodeAndCustId(request.RegionCode, request.CustId)
	if err == nil {
		return response, errors.New("region_code: " + region.RegionCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	regionData := model.Region{
		CustId:     request.CustId,
		RegionCode: request.RegionCode,
		RegionName: request.RegionName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	regionId, err := service.RegionRepository.Store(regionData)
	if err != nil {
		return response, err
	}

	response.RegionId = regionId

	return response, err
}

func (service *regionServiceImpl) applyPrincipalScope(dataFilter entity.RegionQueryFilter) (entity.RegionQueryFilter, error) {
	if !IsPrincipalDistributor(dataFilter.DistributorId) {
		return dataFilter, nil
	}
	if dataFilter.EmployeeId == 0 {
		return dataFilter, errors.New("employee_id is required for principal")
	}

	employee, err := service.EmployeeScopeRepo.FindEmployeeDropdownScope(dataFilter.EmployeeId, dataFilter.CustId)
	if err != nil {
		return dataFilter, err
	}
	dataFilter.Scope = NormalizeScopeSet(employee.RegionScope, employee.AreaScope, employee.DistributorScope)
	return dataFilter, nil
}

func (service *regionServiceImpl) LookupList(dataFilter entity.RegionQueryFilter) (data []entity.RegionLookupResponse, total int, lastPage int, err error) {
	dataFilter, err = service.applyPrincipalScope(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	regions, total, lastPage, err := service.RegionRepository.FindAllByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regions {
		var vResp entity.RegionLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *regionServiceImpl) List(dataFilter entity.RegionQueryFilter) (data []entity.RegionResponse, total int, lastPage int, err error) {
	dataFilter, err = service.applyPrincipalScope(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	regions, total, lastPage, err := service.RegionRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regions {
		var vResp entity.RegionResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *regionServiceImpl) Update(regionId int, request entity.UpdateRegionRequest) (err error) {

	// bank_code & cust id validation, if err == nil and params bankId != bank.Id, this means that code & cust id already exists
	region, err := service.RegionRepository.FindOneByRegionCodeAndCustId(request.RegionCode, request.CustId)
	if err == nil && region.RegionId != regionId {
		return errors.New("region_code: " + region.RegionCode + " is already exists")
	}

	err = service.RegionRepository.Update(regionId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *regionServiceImpl) Delete(custId string, regionId int, userId int64) (err error) {

	err = service.RegionRepository.Delete(custId, regionId, userId)
	if err != nil {
		return err
	}

	return err
}
