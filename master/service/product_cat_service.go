package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type ProductCatService interface {
	Detail(int, string) (entity.ProductCatResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.ProductCatListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.ProductCatLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateProductCatBody) (entity.ProductCatResponse, error)
	Update(int, entity.UpdateProductCatRequest) error
	Delete(string, int, int64) error
}

func NewProductCatService(productCatRepository repository.ProductCatRepository, mProductRepository repository.ProductRepository) *productCatServiceImpl {
	return &productCatServiceImpl{
		ProductCatRepository: productCatRepository,
		MProductRepository:   mProductRepository,
	}
}

type productCatServiceImpl struct {
	ProductCatRepository repository.ProductCatRepository
	MProductRepository   repository.ProductRepository
}

func (service *productCatServiceImpl) Detail(cProId int, custId string) (response entity.ProductCatResponse, err error) {
	productCat, err := service.ProductCatRepository.FindOneByPCatIdAndCustId(cProId, custId)
	if err != nil {
		return response, err
	}

	response.PCatId = productCat.PCatId
	response.PCatCode = productCat.PCatCode
	response.PCatName = productCat.PCatName
	response.IsActive = productCat.IsActive
	response.UpdatedBy = productCat.UpdatedBy
	response.UpdatedAt = productCat.UpdatedAt

	return response, err
}

func (service *productCatServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ProductCatListResponse, total int, lastPage int, err error) {
	productCats, total, lastPage, err := service.ProductCatRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(productCats) > 0 {
		for _, row := range productCats {
			pcModel := entity.ProductCatListResponse{
				PCatId:    row.PCatId,
				PCatCode:  row.PCatCode,
				PCatName:  row.PCatName,
				IsActive:  row.IsActive,
				UpdatedBy: row.UpdatedBy,
				UpdatedAt: row.UpdatedAt,
			}
			if row.UpdatedByName != nil {
				pcModel.UpdatedByName = *row.UpdatedByName
			}
			data = append(data, pcModel)
		}
	}
	return data, total, lastPage, err
}

func (service *productCatServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ProductCatLookupResponse, total int, lastPage int, err error) {
	productCats, total, lastPage, err := service.ProductCatRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(productCats) > 0 {
		for _, row := range productCats {
			data = append(data, entity.ProductCatLookupResponse{
				PCatId:   row.PCatId,
				PCatCode: row.PCatCode,
				PCatName: row.PCatName,
				IsActive: row.IsActive,
			})
		}
	}
	return data, total, lastPage, err
}

func (service *productCatServiceImpl) Store(request entity.CreateProductCatBody) (response entity.ProductCatResponse, err error) {

	// pcat_code & cust id validation, if err == nil, this means that code & cust id already exists
	productCat, err := service.ProductCatRepository.FindOneByPCatCodeAndCustId(request.PCatCode, request.CustId)
	if err == nil {
		return response, errors.New("pcat_code: " + productCat.PCatCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.ProductCat{
		CustId:    request.CustId,
		PCatCode:  request.PCatCode,
		PCatName:  request.PCatName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	cProId, err := service.ProductCatRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.PCatId = cProId

	return response, err
}

func (service *productCatServiceImpl) Update(cProId int, request entity.UpdateProductCatRequest) (err error) {

	// pcat_code & cust id validation, if err == nil and params cProId != productCat.Id, this means that code & cust id already exists
	productCat, err := service.ProductCatRepository.FindOneByPCatCodeAndCustId(request.PCatCode, request.CustId)
	if err == nil && productCat.PCatId != cProId {
		return errors.New("pcat_code: " + productCat.PCatCode + " is already exists")
	}

	err = service.ProductCatRepository.Update(cProId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *productCatServiceImpl) Delete(custId string, cProId int, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(cProId, custId, "pcat_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("pcat_id is still being used")
	}

	err = service.ProductCatRepository.Delete(custId, cProId, userId)
	if err != nil {
		return err
	}

	return err
}
