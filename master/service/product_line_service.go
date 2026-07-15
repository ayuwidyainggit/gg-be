package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ProductLineService interface {
	Detail(int, string) (entity.ProductLineResponse, error)
	List(entity.ProductLineQueryFilter, string) (data []entity.ProductLineResponse, total int, lastPage int, err error)
	LookupList(entity.ProductLineQueryFilter, string) (data []entity.ProductLineLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateProductLineBody) (entity.ProductLineResponse, error)
	Update(int, entity.UpdateProductLineRequest) error
	Delete(string, int, int64) error
}

func NewProductLineService(productLineRepository repository.ProductLineRepository, brandRepository repository.BrandRepository) *productLineServiceImpl {
	return &productLineServiceImpl{
		ProductLineRepository: productLineRepository,
		BrandRepository:       brandRepository,
	}
}

type productLineServiceImpl struct {
	ProductLineRepository repository.ProductLineRepository
	BrandRepository       repository.BrandRepository
}

func (service *productLineServiceImpl) Detail(pLId int, custId string) (response entity.ProductLineResponse, err error) {
	productLine, err := service.ProductLineRepository.FindOneByPLIdAndCustId(pLId, custId)
	if err != nil {
		return response, err
	}

	response.PLId = productLine.PlId
	response.PLCode = productLine.PlCode
	response.PLName = productLine.PlName
	response.EffCall = productLine.EffCall
	response.MinItem = productLine.MinItem
	response.IsActive = productLine.IsActive
	response.UpdatedBy = productLine.UpdatedBy
	response.UpdatedAt = productLine.UpdatedAt

	return response, err
}

func (service *productLineServiceImpl) List(dataFilter entity.ProductLineQueryFilter, custId string) (data []entity.ProductLineResponse, total int, lastPage int, err error) {
	productLines, total, lastPage, err := service.ProductLineRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range productLines {
		var vResp entity.ProductLineResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productLineServiceImpl) LookupList(dataFilter entity.ProductLineQueryFilter, custId string) (data []entity.ProductLineLookupResponse, total int, lastPage int, err error) {
	productLines, total, lastPage, err := service.ProductLineRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range productLines {
		var vResp entity.ProductLineLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productLineServiceImpl) Store(request entity.CreateProductLineBody) (response entity.ProductLineResponse, err error) {

	// pl_code & cust id validation, if err == nil, this means that code & cust id already exists
	productLine, err := service.ProductLineRepository.FindOneByPLCodeAndCustId(request.PlCode, request.CustId)
	if err == nil {
		return response, errors.New("pl_code: " + productLine.PlCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.ProductLine{
		CustId:    request.CustId,
		PlCode:    request.PlCode,
		PlName:    request.PlName,
		EffCall:   request.EffCall,
		MinItem:   request.MinItem,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	pLId, err := service.ProductLineRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.PLId = pLId

	return response, err
}

func (service *productLineServiceImpl) Update(pLId int, request entity.UpdateProductLineRequest) (err error) {

	// pl_code & cust id validation, if err == nil and params pLId != productLine.Id, this means that code & cust id already exists
	productLine, err := service.ProductLineRepository.FindOneByPLCodeAndCustId(request.PlCode, request.CustId)
	if err == nil && productLine.PlId != pLId {
		return errors.New("pl_code: " + productLine.PlCode + " is already exists")
	}

	err = service.ProductLineRepository.Update(pLId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *productLineServiceImpl) Delete(custId string, pLId int, userId int64) (err error) {

	isExists, err := service.BrandRepository.ProductLineIsUsed(pLId, custId)
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("product line is still being used")
	}

	err = service.ProductLineRepository.Delete(custId, pLId, userId)
	if err != nil {
		return err
	}

	return err
}
