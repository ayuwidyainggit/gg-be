package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"
)

type DiscProductService interface {
	Detail(int, int, string) (entity.DiscProductResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.DiscProductListResponse, total int, lastPage int, err error)
	Store(entity.CreateDiscProductBody) (entity.DiscProductResponse, error)
	Update(int, int, entity.UpdateDiscProductRequest) error
	Delete(string, int, int, int64) error
}

func NewDiscProductService(discProductRepository repository.DiscProductRepository) *discProductServiceImpl {
	return &discProductServiceImpl{
		DiscProductRepository: discProductRepository,
	}
}

type discProductServiceImpl struct {
	DiscProductRepository repository.DiscProductRepository
}

func (service *discProductServiceImpl) Detail(discId int, proId int, custId string) (response entity.DiscProductResponse, err error) {
	discProducts, err := service.DiscProductRepository.FindOneByDiscProductIdAndCustId(discId, proId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(discProducts, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *discProductServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DiscProductListResponse, total int, lastPage int, err error) {
	discProducts, total, lastPage, err := service.DiscProductRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range discProducts {
		var vResp entity.DiscProductListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *discProductServiceImpl) Store(request entity.CreateDiscProductBody) (response entity.DiscProductResponse, err error) {
	discProducts, err := service.DiscProductRepository.FindOneByDiscProductIdAndCustId(request.DiscId, request.ProId, request.CustId)
	if err == nil {
		return response, errors.New("pro_id: " + strconv.Itoa(discProducts.ProId) + " is already exists")
	}

	var discProductsData model.DiscProduct

	timeNow := time.Now().In(time.UTC)
	structs.Automapper(request, &discProductsData)
	discProductsData.CreatedAt = &timeNow
	discProductsData.CreatedBy = &request.CreatedBy
	discProductsData.UpdatedBy = &request.CreatedBy
	discProductsData.UpdatedAt = &timeNow
	ProId, err := service.DiscProductRepository.Store(discProductsData)
	if err != nil {
		return response, err
	}

	response.ProId = ProId

	return response, err
}

func (service *discProductServiceImpl) Update(discId int, proId int, request entity.UpdateDiscProductRequest) (err error) {
	// outlet_code & cust id validation, if err == nil and params outletId != outlet.Id, this means that code & cust id already exists
	discProducts, err := service.DiscProductRepository.FindOneByDiscProductIdAndCustId(request.DiscId, request.ProId, request.CustId)
	if err == nil && discProducts.ProId != proId {
		return errors.New("pro_id: " + strconv.Itoa(discProducts.ProId) + " is already exists")
	}
	err = service.DiscProductRepository.Update(discId, proId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *discProductServiceImpl) Delete(custId string, discId int, prodId int, userId int64) (err error) {

	err = service.DiscProductRepository.Delete(custId, discId, prodId, userId)
	if err != nil {
		return err
	}
	return err
}
