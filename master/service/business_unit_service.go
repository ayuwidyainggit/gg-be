package service

import (
	"errors"
	"master/entity"
	"master/repository"
)

func intFromPtr(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

type BusinessUnitService interface {
	GetBusinessUnit(dataFilter entity.BusinessUnitQueryFilter) (interface{}, int, int, error)
}

type businessUnitServiceImpl struct {
	Repository        repository.BusinessUnitRepository
	EmployeeScopeRepo repository.EmployeeScopeRepository
}

func NewBusinessUnitService(repo repository.BusinessUnitRepository, employeeScopeRepo repository.EmployeeScopeRepository) *businessUnitServiceImpl {
	return &businessUnitServiceImpl{Repository: repo, EmployeeScopeRepo: employeeScopeRepo}
}

func (s *businessUnitServiceImpl) GetBusinessUnit(dataFilter entity.BusinessUnitQueryFilter) (interface{}, int, int, error) {
	// Get user info
	userInfo, err := s.Repository.FindUserByUsername(dataFilter.UserName)
	if err != nil {
		return nil, 0, 0, err
	}

	// Check if user is principal or distributor
	if dataFilter.DistributorId == nil || *dataFilter.DistributorId == 0 {
		if dataFilter.EmployeeId == 0 {
			return nil, 0, 0, errors.New("employee_id is required for principal")
		}
		employee, err := s.EmployeeScopeRepo.FindEmployeeDropdownScope(dataFilter.EmployeeId, dataFilter.CustId)
		if err != nil {
			return nil, 0, 0, err
		}
		dataFilter.Scope = NormalizeScopeSet(employee.RegionScope, employee.AreaScope, employee.DistributorScope)
		// User Principal - return array of distributors
		distributors, total, lastPage, err := s.Repository.FindDistributorsByCustId(dataFilter)
		if err != nil {
			return nil, 0, 0, err
		}
		custName, err := s.Repository.FindCustomerNameByCustId(dataFilter.CustId)
		if err != nil {
			return nil, 0, 0, err
		}

		// Map to response
		distributorData := []entity.BusinessUnitDistributorData{}
		for _, d := range distributors {
			areaCode := ""
			areaName := ""
			regionCode := ""
			regionName := ""
			if d.AreaCode != nil {
				areaCode = *d.AreaCode
			}
			if d.AreaName != nil {
				areaName = *d.AreaName
			}
			if d.RegionCode != nil {
				regionCode = *d.RegionCode
			}
			if d.RegionName != nil {
				regionName = *d.RegionName
			}

			distributorData = append(distributorData, entity.BusinessUnitDistributorData{
				CustId:          d.CustId,
				DistributorId:   d.DistributorId,
				DistributorCode: d.DistributorCode,
				DistributorName: d.DistributorName,
				AreaId:          intFromPtr(d.AreaId),
				AreaCode:        areaCode,
				AreaName:        areaName,
				RegionId:        intFromPtr(d.RegionId),
				RegionCode:      regionCode,
				RegionName:      regionName,
			})
		}

		response := entity.BusinessUnitPrincipalResponse{
			CustId:          dataFilter.CustId,
			UserId:          userInfo.UserId,
			UserFullname:    custName,
			CustName:        custName,
			DistributorId:   "",
			DistributorData: distributorData,
		}

		return response, total, lastPage, nil
	} else {
		// User Distributor - return single distributor
		distributor, err := s.Repository.FindDistributorByDistributorId(*dataFilter.DistributorId, dataFilter.CustId)
		if err != nil {
			return nil, 0, 0, err
		}

		areaCode := ""
		areaName := ""
		regionCode := ""
		regionName := ""
		if distributor.AreaCode != nil {
			areaCode = *distributor.AreaCode
		}
		if distributor.AreaName != nil {
			areaName = *distributor.AreaName
		}
		if distributor.RegionCode != nil {
			regionCode = *distributor.RegionCode
		}
		if distributor.RegionName != nil {
			regionName = *distributor.RegionName
		}

		response := entity.BusinessUnitDistributorResponse{
			CustId:          distributor.CustId,
			UserId:          userInfo.UserId,
			UserFullname:    userInfo.UserFullname,
			DistributorId:   distributor.DistributorId,
			DistributorCode: distributor.DistributorCode,
			DistributorName: distributor.DistributorName,
			AreaId:          intFromPtr(distributor.AreaId),
			AreaCode:        areaCode,
			AreaName:        areaName,
			RegionId:        intFromPtr(distributor.RegionId),
			RegionCode:      regionCode,
			RegionName:      regionName,
		}

		return response, 1, 1, nil
	}
}
