package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PrincipalService interface {
	Detail(custId string, principalId int) (entity.PrincipalResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.PrincipalListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.PrincipalLookupResponse, total int, lastPage int, err error)
	Store(entity.CreatePrincipalBody) (entity.PrincipalResponse, error)
	Update(int, entity.UpdatePrincipalRequest) error
	Delete(string, int, int64) error
}

func NewPrincipalService(principalRepository repository.PrincipalRepository, mProductRepository repository.ProductRepository) *principalServiceImpl {
	return &principalServiceImpl{
		PrincipalRepository: principalRepository,
		MProductRepository:  mProductRepository,
	}
}

type principalServiceImpl struct {
	PrincipalRepository repository.PrincipalRepository
	MProductRepository  repository.ProductRepository
}

func (service *principalServiceImpl) Detail(custId string, principalId int) (response entity.PrincipalResponse, err error) {
	principal, err := service.PrincipalRepository.FindOneByPrincipalIdAndCustId(custId, principalId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(principal, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *principalServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PrincipalLookupResponse, total int, lastPage int, err error) {
	principals, total, lastPage, err := service.PrincipalRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range principals {
		var vResp entity.PrincipalLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *principalServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PrincipalListResponse, total int, lastPage int, err error) {
	principals, total, lastPage, err := service.PrincipalRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range principals {
		var vResp entity.PrincipalListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *principalServiceImpl) Store(request entity.CreatePrincipalBody) (response entity.PrincipalResponse, err error) {

	princ, err := service.PrincipalRepository.FindOneByPrincipalCodeAndCustId(request.PrincipalCode, request.CustId)
	if err == nil {
		return response, errors.New("principal_code: " + princ.PrincipalCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	principal := model.Principal{
		CustId:        request.CustId,
		PrincipalCode: request.PrincipalCode,
		PrincipalName: request.PrincipalName,
		IsActive:      request.IsActive,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	principalId, err := service.PrincipalRepository.Store(principal)
	if err != nil {
		return response, err
	}

	response.PrincipalId = principalId

	return response, err
}

func (service *principalServiceImpl) Update(principalId int, request entity.UpdatePrincipalRequest) (err error) {

	princ, err := service.PrincipalRepository.FindOneByPrincipalCodeAndCustId(request.PrincipalCode, request.CustId)
	if err == nil && princ.PrincipalId != principalId {
		return errors.New("principal_code: " + princ.PrincipalCode + " is already exists")
	}

	err = service.PrincipalRepository.Update(principalId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *principalServiceImpl) Delete(custId string, principalId int, deletedBy int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(principalId, custId, "principal_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("principal_id is still being used")
	}

	err = service.PrincipalRepository.Delete(custId, principalId, deletedBy)
	if err != nil {
		return err
	}

	return err
}
