package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type OutletLocService interface {
	Detail(int, string) (entity.OutletLocResponse, error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.OutletLocLookupResponse, total int, lastPage int, err error)
	List(entity.GeneralQueryFilter, string) (data []entity.OutletLocResponse, total int, lastPage int, err error)
	Store(entity.CreateOutletLocBody) (entity.OutletLocResponse, error)
	Update(int, entity.UpdateOutletLocRequest) error
	Delete(string, int, int64) error
}

func NewOutletLocService(outletLocRepository repository.OutletLocRepository) *outletLocServiceImpl {
	return &outletLocServiceImpl{
		OutletLocRepository: outletLocRepository,
	}
}

type outletLocServiceImpl struct {
	OutletLocRepository repository.OutletLocRepository
}

func (service *outletLocServiceImpl) Detail(outletLocId int, custId string) (response entity.OutletLocResponse, err error) {
	outletLoc, err := service.OutletLocRepository.FindOneByOutletLocIdAndCustId(outletLocId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(outletLoc, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *outletLocServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletLocResponse, total int, lastPage int, err error) {
	outletLocs, total, lastPage, err := service.OutletLocRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletLocs {
		var vResp entity.OutletLocResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletLocServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.OutletLocLookupResponse, total int, lastPage int, err error) {
	outletLocs, total, lastPage, err := service.OutletLocRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range outletLocs {
		var vResp entity.OutletLocLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *outletLocServiceImpl) Store(request entity.CreateOutletLocBody) (response entity.OutletLocResponse, err error) {

	// outletLoc_code & cust id validation, if err == nil, this means that code & cust id already exists
	outletLoc, err := service.OutletLocRepository.FindOneByOutletLocCodeAndCustId(request.OtLocCode, request.CustId)
	if err == nil {
		return response, errors.New("ot_loc_code: " + outletLoc.OtLocCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	var otLocData model.OutletLoc
	err = structs.Automapper(request, &otLocData)
	if err != nil {
		return response, err
	}

	otLocData.CreatedAt = &timeNow
	otLocData.CreatedBy = &request.CreatedBy
	otLocData.UpdatedAt = &timeNow
	otLocData.UpdatedBy = &request.CreatedBy

	outletLocId, err := service.OutletLocRepository.Store(otLocData)
	if err != nil {
		return response, err
	}

	response.OtLocId = outletLocId

	return response, err
}

func (service *outletLocServiceImpl) Update(outletLocId int, request entity.UpdateOutletLocRequest) (err error) {

	// outletLoc_code & cust id validation, if err == nil and params outletLocId != outletLoc.Id, this means that code & cust id already exists
	outletLoc, err := service.OutletLocRepository.FindOneByOutletLocCodeAndCustId(request.OtLocCode, request.CustId)
	if err == nil && outletLoc.OtLocId != outletLocId {
		return errors.New("ot_loc_code: " + outletLoc.OtLocCode + " is already exists")
	}

	err = service.OutletLocRepository.Update(outletLocId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *outletLocServiceImpl) Delete(custId string, outletLocId int, userId int64) (err error) {

	err = service.OutletLocRepository.Delete(custId, outletLocId, userId)
	if err != nil {
		return err
	}

	return err
}
