package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type AreaService interface {
	Detail(int, string) (entity.AreaListResponse, error)
	LookupList(entity.AreaQueryFilter) (data []entity.AreaListResponse, total int, lastPage int, err error)
	List(entity.AreaQueryFilter) (data []entity.AreaListResponse, total int, lastPage int, err error)
	Store(entity.CreateAreaBody) (entity.AreaResponse, error)
	Update(int, entity.UpdateAreaRequest) error
	Delete(string, int, int64) error
}

func NewAreaService(areaRepository repository.AreaRepository, employeeScopeRepo repository.EmployeeScopeRepository) *areaServiceImpl {
	return &areaServiceImpl{
		AreaRepository:      areaRepository,
		EmployeeScopeRepo: employeeScopeRepo,
	}
}

type areaServiceImpl struct {
	AreaRepository      repository.AreaRepository
	EmployeeScopeRepo repository.EmployeeScopeRepository
}

func (service *areaServiceImpl) Detail(areaId int, custId string) (response entity.AreaListResponse, err error) {
	area, err := service.AreaRepository.FindOneByAreaIdAndCustId(areaId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(area, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *areaServiceImpl) Store(request entity.CreateAreaBody) (response entity.AreaResponse, err error) {
	// bank_code & cust id validation, if err == nil, this means that code & cust id already exists
	area, err := service.AreaRepository.FindOneByAreaCodeAndCustId(request.AreaCode, request.CustId)
	if err == nil {
		return response, errors.New("area_code: " + area.AreaCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	areaData := model.Area{
		CustId:     request.CustId,
		AreaCode:   request.AreaCode,
		AreaName:   request.AreaName,
		RegionId:   request.RegionId,
		OfficialId: request.OfficialId,
		IsActive:   true,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	areaId, err := service.AreaRepository.Store(areaData)
	if err != nil {
		return response, err
	}

	response.AreaId = areaId

	return response, err
}

func (service *areaServiceImpl) applyPrincipalScope(dataFilter entity.AreaQueryFilter) (entity.AreaQueryFilter, error) {
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

func (service *areaServiceImpl) LookupList(dataFilter entity.AreaQueryFilter) (data []entity.AreaListResponse, total int, lastPage int, err error) {
	dataFilter, err = service.applyPrincipalScope(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	regions, total, lastPage, err := service.AreaRepository.FindAllByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range regions {
		var vResp entity.AreaListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *areaServiceImpl) List(dataFilter entity.AreaQueryFilter) (data []entity.AreaListResponse, total int, lastPage int, err error) {
	dataFilter, err = service.applyPrincipalScope(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	areas, total, lastPage, err := service.AreaRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range areas {
		var vResp entity.AreaListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *areaServiceImpl) Update(areaId int, request entity.UpdateAreaRequest) (err error) {

	area, err := service.AreaRepository.FindOneByAreaCodeAndCustId(request.AreaCode, request.CustId)
	if err == nil && area.AreaId != areaId {
		return errors.New("area_code: " + area.AreaCode + " is already exists")
	}

	err = service.AreaRepository.Update(areaId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *areaServiceImpl) Delete(custId string, areaId int, userId int64) (err error) {

	err = service.AreaRepository.Delete(custId, areaId, userId)
	if err != nil {
		return err
	}

	return err
}
