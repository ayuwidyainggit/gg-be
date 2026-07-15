package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type ProductCoreTaxService interface {
	Detail(string, string) (entity.ProductCoreTaxResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.ProductCoreTaxListResponse, total int, lastPage int, err error)
	Store(entity.CreateProductCoreTaxBody) (entity.ProductCoreTaxResponse, error)
	Update(string, entity.UpdateProductCoreTaxRequest) error
	Delete(string, string, int64) error
}

func NewProductCoreTaxService(coreTaxRepository repository.ProductCoreTaxRepository) *productCoreTaxServiceImpl {
	return &productCoreTaxServiceImpl{
		ProductCoreTaxRepository: coreTaxRepository,
	}
}

type productCoreTaxServiceImpl struct {
	ProductCoreTaxRepository repository.ProductCoreTaxRepository
}

func (service *productCoreTaxServiceImpl) Detail(coreTaxId string, custId string) (response entity.ProductCoreTaxResponse, err error) {
	coreTax, err := service.ProductCoreTaxRepository.FindOneByProCodeCoreTaxAndCustId(coreTaxId, custId)
	if err != nil {
		return response, err
	}

	response.ProCodeCoreTax = coreTax.ProCodeCoreTax
	response.CatCoreTax = coreTax.CatCoreTax
	response.ProNameCoreTax = coreTax.ProNameCoreTax
	response.IsActive = coreTax.IsActive
	response.UpdatedBy = coreTax.UpdatedBy
	response.UpdatedAt = coreTax.UpdatedAt

	return response, err
}

func (service *productCoreTaxServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ProductCoreTaxListResponse, total int, lastPage int, err error) {
	coreTaxs, total, lastPage, err := service.ProductCoreTaxRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(coreTaxs) > 0 {
		for _, row := range coreTaxs {
			coreTaxmodel := entity.ProductCoreTaxListResponse{
				ProCodeCoreTax: row.ProCodeCoreTax,
				CatCoreTax:     row.CatCoreTax,
				ProNameCoreTax: row.ProNameCoreTax,
				IsActive:       row.IsActive,
				UpdatedBy:      row.UpdatedBy,
				UpdatedAt:      row.UpdatedAt,
			}
			if row.UpdatedByName != nil {
				coreTaxmodel.UpdatedByName = *row.UpdatedByName
			}
			data = append(data, coreTaxmodel)
		}
	}
	return data, total, lastPage, err
}

func (service *productCoreTaxServiceImpl) Store(request entity.CreateProductCoreTaxBody) (response entity.ProductCoreTaxResponse, err error) {

	// coreTax_code & cust id validation, if err == nil, this means that code & cust id already exists
	coreTax, err := service.ProductCoreTaxRepository.FindOneByProCodeCoreTaxAndCustId(request.ProCodeCoreTax, request.CustId)
	if err == nil {
		return response, errors.New("product code coreTax: " + coreTax.ProCodeCoreTax + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	coreTaxData := model.ProductCoreTax{
		CustId:         request.CustId,
		ProCodeCoreTax: request.ProCodeCoreTax,
		ProNameCoreTax: request.ProNameCoreTax,
		IsActive:       request.IsActive,
		CreatedAt:      &timeNow,
		CreatedBy:      &request.CreatedBy,
		UpdatedAt:      &timeNow,
		UpdatedBy:      &request.CreatedBy,
	}

	coreTaxId, err := service.ProductCoreTaxRepository.Store(coreTaxData)
	if err != nil {
		return response, err
	}

	response.ProCodeCoreTax = coreTaxId

	return response, err
}

func (service *productCoreTaxServiceImpl) Update(coreTaxId string, request entity.UpdateProductCoreTaxRequest) (err error) {

	// coreTax_code & cust id validation, if err == nil and params coreTaxId != coreTax.Id, this means that code & cust id already exists
	coreTax, err := service.ProductCoreTaxRepository.FindOneByProCodeCoreTaxAndCustId(request.ProCodeCoreTax, request.CustId)
	if err == nil && coreTax.ProCodeCoreTax != coreTaxId {
		return errors.New("product code coreTax: " + coreTax.ProCodeCoreTax + " is already exists")
	}

	err = service.ProductCoreTaxRepository.Update(coreTaxId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *productCoreTaxServiceImpl) Delete(custId string, coreTaxId string, userId int64) (err error) {
	err = service.ProductCoreTaxRepository.Delete(custId, coreTaxId, userId)
	if err != nil {
		return err
	}

	return err
}
