package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type MDivisionService interface {
	Detail(MDivisionId int64, custId string) (response entity.MDivisionDetailsResponse, err error)
	Store(request entity.CreateDivisionBody) (response entity.MDivisionDetailsResponse, err error)
	Update(divisionID int64, request entity.UpdateDivisionBody) (err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.MDivisionDetailsResponse, total int, lastPage int, err error)
	LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MDivisionLookupResponse, total int, lastPage int, err error)
	Delete(custId string, MDivisionId int64, userId int64) (err error)
}
type MDivisionServiceImpl struct {
	MDivisionRepository repository.MDivisionRepository
}

func NewMDivisionService(mdivisionRepository repository.MDivisionRepository) *MDivisionServiceImpl {
	return &MDivisionServiceImpl{MDivisionRepository: mdivisionRepository}
}
func (service *MDivisionServiceImpl) Detail(MDivisionId int64, custId string) (response entity.MDivisionDetailsResponse, err error) {
	MDivision, err := service.MDivisionRepository.FindOneByMDivisionIdAndCustId(MDivisionId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(MDivision, &response)
	if err != nil {
		return response, err
	}
	return
}
func (service *MDivisionServiceImpl) Store(request entity.CreateDivisionBody) (response entity.MDivisionDetailsResponse, err error) {

	// emp_code & cust id validation, if err == nil, this means that code & cust id already exists
	division, err := service.MDivisionRepository.FindOneByMDivisionCodeAndCustId(request.DivisionCode, request.CustId)
	if err == nil {
		return response, errors.New("division_code: " + division.DivisionCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	var divisionData model.MDivision
	structs.Automapper(request, &divisionData)

	divisionData.CreatedAt = &timeNow
	divisionData.CreatedBy = &request.CreatedBy
	divisionData.UpdatedAt = &timeNow
	divisionData.UpdatedBy = &request.CreatedBy

	divisionId, err := service.MDivisionRepository.Store(divisionData)
	if err != nil {
		return response, err
	}

	response.DivisionID = divisionId

	return response, err
}

func (service *MDivisionServiceImpl) Update(divisionID int64, request entity.UpdateDivisionBody) (err error) {

	// emp_code & cust id validation, if err == nil and params employeeId != employee.Id, this means that code & cust id already exists
	division, err := service.MDivisionRepository.FindOneByMDivisionCodeAndCustId(request.DivisionCode, request.CustId)
	if err == nil && division.DivisionID != divisionID {
		return errors.New("division_code: " + division.DivisionCode + " is already exists")
	}

	err = service.MDivisionRepository.Update(divisionID, request)
	if err != nil {
		return err
	}

	return err
}
func (service *MDivisionServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.MDivisionDetailsResponse, total int, lastPage int, err error) {

	divisions, total, lastPage, err := service.MDivisionRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range divisions {
		var vResp entity.MDivisionDetailsResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MDivisionServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MDivisionLookupResponse, total int, lastPage int, err error) {
	var divisions []model.MDivision

	divisions, total, lastPage, err = service.MDivisionRepository.FindAllByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range divisions {
		var vResp entity.MDivisionLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MDivisionServiceImpl) Delete(custId string, MDivisionId int64, userId int64) (err error) {
	err = service.MDivisionRepository.Delete(custId, MDivisionId, userId)
	if err != nil {
		return err
	}

	return err
}
