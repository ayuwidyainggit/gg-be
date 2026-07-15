package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type UnitCoreTaxService interface {
	Detail(string, string) (entity.UnitCoreTaxResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.UnitCoreTaxListResponse, total int, lastPage int, err error)
	Store(entity.CreateUnitCoreTaxBody) (entity.UnitCoreTaxResponse, error)
	Update(string, entity.UpdateUnitCoreTaxRequest) error
	Delete(string, string, int64) error
}

func NewUnitCoreTaxService(coreTaxRepository repository.UnitCoreTaxRepository) *coreTaxServiceImpl {
	return &coreTaxServiceImpl{
		UnitCoreTaxRepository: coreTaxRepository,
	}
}

type coreTaxServiceImpl struct {
	UnitCoreTaxRepository repository.UnitCoreTaxRepository
}

func (service *coreTaxServiceImpl) Detail(coreTaxId string, custId string) (response entity.UnitCoreTaxResponse, err error) {
	coreTax, err := service.UnitCoreTaxRepository.FindOneByUnitIdCoreTaxAndCustId(coreTaxId, custId)
	if err != nil {
		return response, err
	}

	response.UnitIdCoreTax = coreTax.UnitIdCoreTax
	response.UnitNameCoreTax = coreTax.UnitNameCoreTax
	response.IsActive = coreTax.IsActive
	response.UpdatedBy = coreTax.UpdatedBy
	response.UpdatedAt = coreTax.UpdatedAt

	return response, err
}

func (service *coreTaxServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.UnitCoreTaxListResponse, total int, lastPage int, err error) {
	coreTaxs, total, lastPage, err := service.UnitCoreTaxRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(coreTaxs) > 0 {
		for _, row := range coreTaxs {
			coreTaxmodel := entity.UnitCoreTaxListResponse{
				UnitIdCoreTax:   row.UnitIdCoreTax,
				UnitNameCoreTax: row.UnitNameCoreTax,
				IsActive:        row.IsActive,
				UpdatedBy:       row.UpdatedBy,
				UpdatedAt:       row.UpdatedAt,
			}
			if row.UpdatedByName != nil {
				coreTaxmodel.UpdatedByName = *row.UpdatedByName
			}
			data = append(data, coreTaxmodel)
		}
	}
	return data, total, lastPage, err
}

func (service *coreTaxServiceImpl) Store(request entity.CreateUnitCoreTaxBody) (response entity.UnitCoreTaxResponse, err error) {

	// coreTax_code & cust id validation, if err == nil, this means that code & cust id already exists
	coreTax, err := service.UnitCoreTaxRepository.FindOneByUnitIdCoreTaxAndCustId(request.UnitIdCoreTax, request.CustId)
	if err == nil {
		return response, errors.New("coreTax_id: " + coreTax.UnitIdCoreTax + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	coreTaxData := model.UnitCoreTax{
		CustId:          request.CustId,
		UnitIdCoreTax:   request.UnitIdCoreTax,
		UnitNameCoreTax: request.UnitNameCoreTax,
		IsActive:        request.IsActive,
		CreatedAt:       &timeNow,
		CreatedBy:       &request.CreatedBy,
		UpdatedAt:       &timeNow,
		UpdatedBy:       &request.CreatedBy,
	}

	coreTaxId, err := service.UnitCoreTaxRepository.Store(coreTaxData)
	if err != nil {
		return response, err
	}

	response.UnitIdCoreTax = coreTaxId

	return response, err
}

func (service *coreTaxServiceImpl) Update(coreTaxId string, request entity.UpdateUnitCoreTaxRequest) (err error) {

	// coreTax_code & cust id validation, if err == nil and params coreTaxId != coreTax.Id, this means that code & cust id already exists
	coreTax, err := service.UnitCoreTaxRepository.FindOneByUnitIdCoreTaxAndCustId(request.UnitIdCoreTax, request.CustId)
	if err == nil && coreTax.UnitIdCoreTax != coreTaxId {
		return errors.New("coreTax_id: " + coreTax.UnitIdCoreTax + " is already exists")
	}

	err = service.UnitCoreTaxRepository.Update(coreTaxId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *coreTaxServiceImpl) Delete(custId string, coreTaxId string, userId int64) (err error) {
	err = service.UnitCoreTaxRepository.Delete(custId, coreTaxId, userId)
	if err != nil {
		return err
	}

	return err
}
