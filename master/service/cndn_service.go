package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type CndnService interface {
	Detail(int, string) (entity.CndnResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.CndnResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.CndnLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateCndnBody) (entity.CndnResponse, error)
	Update(int, entity.UpdateCndnRequest) error
	Delete(string, int, int64) error
}

func NewCndnService(cndnRepository repository.CndnRepository) *cndnServiceImpl {
	return &cndnServiceImpl{
		CndnRepository: cndnRepository,
	}
}

type cndnServiceImpl struct {
	CndnRepository repository.CndnRepository
}

func (service *cndnServiceImpl) Detail(cndnId int, custId string) (response entity.CndnResponse, err error) {
	cndn, err := service.CndnRepository.FindOneByCndnIdAndCustId(cndnId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(cndn, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *cndnServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.CndnResponse, total int, lastPage int, err error) {
	cndns, total, lastPage, err := service.CndnRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range cndns {
		var vResp entity.CndnResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *cndnServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.CndnLookupResponse, total int, lastPage int, err error) {
	cndns, total, lastPage, err := service.CndnRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range cndns {
		var vResp entity.CndnLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *cndnServiceImpl) Store(request entity.CreateCndnBody) (response entity.CndnResponse, err error) {

	// cndn_code & cust id validation, if err == nil, this means that code & cust id already exists
	cndn, err := service.CndnRepository.FindOneByCndnCodeAndCustId(request.CndnCode, request.CustId)
	if err == nil {
		return response, errors.New("cndn_code: " + cndn.CndnCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	cndnData := model.Cndn{
		CustId:    request.CustId,
		CndnCode:  request.CndnCode,
		CndnName:  request.CndnName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	cndnId, err := service.CndnRepository.Store(cndnData)
	if err != nil {
		return response, err
	}

	response.CndnId = cndnId

	return response, err
}

func (service *cndnServiceImpl) Update(cndnId int, request entity.UpdateCndnRequest) (err error) {

	// cndn_code & cust id validation, if err == nil and params cndnId != cndn.Id, this means that code & cust id already exists
	cndn, err := service.CndnRepository.FindOneByCndnCodeAndCustId(request.CndnCode, request.CustId)
	if err == nil && cndn.CndnId != cndnId {
		return errors.New("cndn_code: " + cndn.CndnCode + " is already exists")
	}

	err = service.CndnRepository.Update(cndnId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *cndnServiceImpl) Delete(custId string, cndnId int, userId int64) (err error) {

	err = service.CndnRepository.Delete(custId, cndnId, userId)
	if err != nil {
		return err
	}

	return err
}
