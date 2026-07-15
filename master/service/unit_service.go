package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type UnitService interface {
	Detail(string, string) (entity.UnitResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.UnitListResponse, total int, lastPage int, err error)
	Store(entity.CreateUnitBody) (entity.UnitResponse, error)
	Update(string, entity.UpdateUnitRequest) error
	Delete(string, string, int64) error
}

func NewUnitService(unitRepository repository.UnitRepository, mProductRepository repository.ProductRepository) *unitServiceImpl {
	return &unitServiceImpl{
		UnitRepository:     unitRepository,
		MProductRepository: mProductRepository,
	}
}

type unitServiceImpl struct {
	UnitRepository     repository.UnitRepository
	MProductRepository repository.ProductRepository
}

func (service *unitServiceImpl) Detail(unitId string, custId string) (response entity.UnitResponse, err error) {
	unit, err := service.UnitRepository.FindOneByUnitIdAndCustId(unitId, custId)
	if err != nil {
		return response, err
	}

	response.UnitId = unit.UnitId
	response.UnitName = unit.UnitName
	response.UnitIdCoreTax = unit.UnitIdCoreTax
	response.UnitNameCoreTax = unit.UnitNameCoreTax
	response.IsActive = unit.IsActive
	response.UpdatedBy = unit.UpdatedBy
	response.UpdatedAt = unit.UpdatedAt

	return response, err
}

func (service *unitServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.UnitListResponse, total int, lastPage int, err error) {
	units, total, lastPage, err := service.UnitRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(units) > 0 {
		for _, row := range units {
			unitmodel := entity.UnitListResponse{
				UnitId:          row.UnitId,
				UnitName:        row.UnitName,
				UnitIdCoreTax:   row.UnitIdCoreTax,
				UnitNameCoreTax: row.UnitNameCoreTax,
				IsActive:        row.IsActive,
				UpdatedBy:       row.UpdatedBy,
				UpdatedAt:       row.UpdatedAt,
			}
			if row.UpdatedByName != nil {
				unitmodel.UpdatedByName = *row.UpdatedByName
			}
			data = append(data, unitmodel)
		}
	}
	return data, total, lastPage, err
}

func (service *unitServiceImpl) Store(request entity.CreateUnitBody) (response entity.UnitResponse, err error) {

	// unit_code & cust id validation, if err == nil, this means that code & cust id already exists
	unit, err := service.UnitRepository.FindOneByUnitIdAndCustId(request.UnitId, request.CustId)
	if err == nil {
		return response, errors.New("unit_id: " + unit.UnitId + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	unitData := model.Unit{
		CustId:        request.CustId,
		UnitId:        request.UnitId,
		UnitName:      request.UnitName,
		UnitIdCoreTax: &request.UnitIdCoreTax,
		IsActive:      request.IsActive,
		CreatedAt:     &timeNow,
		CreatedBy:     &request.CreatedBy,
		UpdatedAt:     &timeNow,
		UpdatedBy:     &request.CreatedBy,
	}

	unitId, err := service.UnitRepository.Store(unitData)
	if err != nil {
		return response, err
	}

	response.UnitId = unitId

	return response, err
}

func (service *unitServiceImpl) Update(unitId string, request entity.UpdateUnitRequest) (err error) {

	// unit_code & cust id validation, if err == nil and params unitId != unit.Id, this means that code & cust id already exists
	unit, err := service.UnitRepository.FindOneByUnitIdAndCustId(request.UnitId, request.CustId)
	if err == nil && unit.UnitId != unitId {
		return errors.New("unit_id: " + unit.UnitId + " is already exists")
	}

	err = service.UnitRepository.Update(unitId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *unitServiceImpl) Delete(custId string, unitId string, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsKeyExists(unitId, custId, "unit_id1")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("unit_id is still being used")
	}

	err = service.UnitRepository.Delete(custId, unitId, userId)
	if err != nil {
		return err
	}

	return err
}
