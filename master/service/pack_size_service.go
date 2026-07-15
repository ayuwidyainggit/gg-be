package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PackSizeService interface {
	Detail(int, string) (entity.PackSizeResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.PackSizeListResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.PackSizeLookupResponse, total int, lastPage int, err error)
	Store(entity.CreatePackSizeBody) (entity.PackSizeResponse, error)
	Update(int, entity.UpdatePackSizeRequest) error
	Delete(string, int, int64) error
}

func NewPackSizeService(packSizeRepository repository.PackSizeRepository, mProductRepository repository.ProductRepository) *packSizeServiceImpl {
	return &packSizeServiceImpl{
		PackSizeRepository: packSizeRepository,
		MProductRepository: mProductRepository,
	}
}

type packSizeServiceImpl struct {
	PackSizeRepository repository.PackSizeRepository
	MProductRepository repository.ProductRepository
}

func (service *packSizeServiceImpl) Detail(pSizeId int, custId string) (response entity.PackSizeResponse, err error) {
	packSize, err := service.PackSizeRepository.FindOneByPSizeIdAndCustId(pSizeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(packSize, &response)
	if err != nil {
		return response, err
	}
	return response, err
}

func (service *packSizeServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PackSizeListResponse, total int, lastPage int, err error) {
	packSizes, total, lastPage, err := service.PackSizeRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range packSizes {
		var vResp entity.PackSizeListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *packSizeServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PackSizeLookupResponse, total int, lastPage int, err error) {
	packSizes, total, lastPage, err := service.PackSizeRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range packSizes {
		var vResp entity.PackSizeLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *packSizeServiceImpl) Store(request entity.CreatePackSizeBody) (response entity.PackSizeResponse, err error) {

	// psize_code & cust id validation, if err == nil, this means that code & cust id already exists
	packSize, err := service.PackSizeRepository.FindOneByPSizeCodeAndCustId(request.PsizeCode, request.CustId)
	if err == nil {
		return response, errors.New("psize_code: " + packSize.PsizeCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.PackSize{
		CustId:    request.CustId,
		PsizeCode: request.PsizeCode,
		PsizeName: request.PsizeName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	pSizeId, err := service.PackSizeRepository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.PSizeId = pSizeId

	return response, err
}

func (service *packSizeServiceImpl) Update(pSizeId int, request entity.UpdatePackSizeRequest) (err error) {

	// psize_code & cust id validation, if err == nil and params pSizeId != packSize.Id, this means that code & cust id already exists
	packSize, err := service.PackSizeRepository.FindOneByPSizeCodeAndCustId(request.PsizeCode, request.CustId)
	if err == nil && packSize.PsizeId != pSizeId {
		return errors.New("psize_code: " + packSize.PsizeCode + " is already exists")
	}

	err = service.PackSizeRepository.Update(pSizeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *packSizeServiceImpl) Delete(custId string, pSizeId int, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(pSizeId, custId, "psize_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("psize_id is still being used")
	}

	err = service.PackSizeRepository.Delete(custId, pSizeId, userId)
	if err != nil {
		return err
	}

	return err
}
