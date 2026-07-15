package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SalesTeamService interface {
	FindParentCustId(string) (entity.MCustomerResp, error)
	Detail(int, string) (entity.SalesTeamResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.SalesTeamResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.SalesTeamLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateSalesTeamBody) (entity.SalesTeamResponse, error)
	Update(int, entity.UpdateSalesTeamRequest) error
	Delete(string, int, int64) error
	ListByDistributor(entity.GeneralQueryFilter) (data []entity.SalesTeamResponse, total int, lastPage int, err error)
}

func NewSalesTeamService(salesTeamRepository repository.SalesTeamRepository) *salesTeamServiceImpl {
	return &salesTeamServiceImpl{
		SalesTeamRepository: salesTeamRepository,
	}
}

type salesTeamServiceImpl struct {
	SalesTeamRepository repository.SalesTeamRepository
}

func (service *salesTeamServiceImpl) FindParentCustId(custId string) (response entity.MCustomerResp, err error) {
	mCustomer, err := service.SalesTeamRepository.FindOneParentCustId(custId)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(mCustomer, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *salesTeamServiceImpl) Detail(salesTeamId int, custId string) (response entity.SalesTeamResponse, err error) {
	salesTeam, err := service.SalesTeamRepository.FindOneBySalesTeamIdAndCustId(salesTeamId, custId)
	if err != nil {
		return response, err
	}

	response.SalesTeamId = salesTeam.SalesTeamId
	response.SalesTeamCode = salesTeam.SalesTeamCode
	response.SalesTeamName = salesTeam.SalesTeamName
	response.IsActive = salesTeam.IsActive
	response.UpdatedBy = salesTeam.UpdatedBy
	response.UpdatedAt = salesTeam.UpdatedAt

	return response, err
}

func (service *salesTeamServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SalesTeamResponse, total int, lastPage int, err error) {
	salesTeams, total, lastPage, err := service.SalesTeamRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesTeams {
		var vResp entity.SalesTeamResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesTeamServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SalesTeamLookupResponse, total int, lastPage int, err error) {
	salesTeams, total, lastPage, err := service.SalesTeamRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesTeams {
		var vResp entity.SalesTeamLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesTeamServiceImpl) Store(request entity.CreateSalesTeamBody) (response entity.SalesTeamResponse, err error) {

	// salesTeam_code & cust id validation, if err == nil, this means that code & cust id already exists
	salesTeam, err := service.SalesTeamRepository.FindOneBySalesTeamCodeAndCustId(request.SalesTeamCode, request.CustId)
	if err == nil {
		return response, errors.New("sales_team_code: " + salesTeam.SalesTeamCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	salesTeamData := model.SalesTeam{
		CustId:        request.CustId,
		SalesTeamCode: request.SalesTeamCode,
		SalesTeamName: request.SalesTeamName,
		IsActive:      request.IsActive,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	salesTeamId, err := service.SalesTeamRepository.Store(salesTeamData)
	if err != nil {
		return response, err
	}

	response.SalesTeamId = salesTeamId

	return response, err
}

func (service *salesTeamServiceImpl) Update(salesTeamId int, request entity.UpdateSalesTeamRequest) (err error) {

	// salesTeam_code & cust id validation, if err == nil and params salesTeamId != salesTeam.Id, this means that code & cust id already exists
	salesTeam, err := service.SalesTeamRepository.FindOneBySalesTeamCodeAndCustId(request.SalesTeamCode, request.CustId)
	if err == nil && salesTeam.SalesTeamId != salesTeamId {
		return errors.New("sales_team_code: " + salesTeam.SalesTeamCode + " is already exists")
	}

	err = service.SalesTeamRepository.Update(salesTeamId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *salesTeamServiceImpl) Delete(custId string, salesTeamId int, userId int64) (err error) {

	err = service.SalesTeamRepository.Delete(custId, salesTeamId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *salesTeamServiceImpl) ListByDistributor(dataFilter entity.GeneralQueryFilter) (data []entity.SalesTeamResponse, total int, lastPage int, err error) {

	mCustomer, err := service.SalesTeamRepository.FindOneCustomerByDistributorID(dataFilter.DistributorID)
	if err != nil {
		return data, total, lastPage, err
	}
	dataFilter.CustId = mCustomer.CustId
	dataFilter.ParentCustId = mCustomer.ParentCustId

	salesTeams, total, lastPage, err := service.SalesTeamRepository.FindAllByCustId(dataFilter, mCustomer.CustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(salesTeams) > 0 {
		for _, row := range salesTeams {
			var vResp entity.SalesTeamResponse
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}

	return data, total, lastPage, err
}	