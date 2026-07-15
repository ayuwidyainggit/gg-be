package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type BrandService interface {
	Detail(int, string) (entity.BrandResponse, error)
	LookupList(entity.BrandQueryFilter, string) (data []entity.BrandLookupResponse, total int, lastPage int, err error)
	List(entity.BrandQueryFilter, string) (data []entity.BrandListResponse, total int, lastPage int, err error)
	Store(entity.CreateBrandBody) (entity.BrandResponse, error)
	Update(int, entity.UpdateBrandRequest) error
	Delete(string, int, int64) error
}

func NewBrandService(brandRepository repository.BrandRepository) *brandServiceImpl {
	return &brandServiceImpl{
		BrandRepository: brandRepository,
	}
}

type brandServiceImpl struct {
	BrandRepository repository.BrandRepository
}

func (service *brandServiceImpl) Detail(brandId int, custId string) (response entity.BrandResponse, err error) {
	brand, err := service.BrandRepository.FindOneByBrandIdAndCustId(brandId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(brand, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *brandServiceImpl) List(dataFilter entity.BrandQueryFilter, custId string) (data []entity.BrandListResponse, total int, lastPage int, err error) {
	brands, total, lastPage, err := service.BrandRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range brands {
		var vResp entity.BrandListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *brandServiceImpl) LookupList(dataFilter entity.BrandQueryFilter, custId string) (data []entity.BrandLookupResponse, total int, lastPage int, err error) {
	brands, total, lastPage, err := service.BrandRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range brands {
		var vResp entity.BrandLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *brandServiceImpl) Store(request entity.CreateBrandBody) (response entity.BrandResponse, err error) {

	// brand_code & cust id validation, if err == nil, this means that code & cust id already exists
	brand, err := service.BrandRepository.FindOneByBrandCodeAndCustId(request.BrandCode, request.CustId)
	if err == nil {
		return response, errors.New("brand_code: " + brand.BrandCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	brandData := model.Brand{
		CustId:    request.CustId,
		BrandCode: request.BrandCode,
		BrandName: request.BrandName,
		PlId:      request.PlId,
		EffCall:   request.EffCall,
		MinItem:   request.MinItem,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	brandId, err := service.BrandRepository.Store(brandData)
	if err != nil {
		return response, err
	}

	response.BrandId = brandId

	return response, err
}

func (service *brandServiceImpl) Update(brandId int, request entity.UpdateBrandRequest) (err error) {

	brand, err := service.BrandRepository.FindOneByBrandCodeAndCustId(request.BrandCode, request.CustId)
	if err == nil && brand.BrandId != brandId {
		return errors.New("brand_code: " + brand.BrandCode + " is already exists")
	}

	err = service.BrandRepository.Update(brandId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *brandServiceImpl) Delete(custId string, brandId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(brandId, custId, "brand_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("brand_id is still being used")
	// }

	err = service.BrandRepository.Delete(custId, brandId, userId)
	if err != nil {
		return err
	}

	return err
}
