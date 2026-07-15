package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SubBrand1Service interface {
	Detail(int, string) (entity.SubBrand1Response, error)
	List(entity.SubBrand1QueryFilter, string) (data []entity.SubBrand1ListResponse, total int, lastPage int, err error)
	LookupList(entity.SubBrand1QueryFilter, string) (data []entity.SubBrand1LookupResponse, total int, lastPage int, err error)
	MatGroupList(entity.SubBrand1QueryFilter, string) (data []entity.SubBrand1MatGroupResponse, total int, lastPage int, err error)
	SubBrandList(entity.SubBrandQueryFilter, string) (data []entity.SubBrandResponse, total int, lastPage int, err error)
	Store(entity.CreateSubBrand1Body) (entity.SubBrand1Response, error)
	Update(int, entity.UpdateSubBrand1Request) error
	Delete(string, int, int64) error
}

func NewSubBrand1Service(subBrand1Repository repository.SubBrand1Repository) *subBrand1ServiceImpl {
	return &subBrand1ServiceImpl{
		SubBrand1Repository: subBrand1Repository,
	}
}

type subBrand1ServiceImpl struct {
	SubBrand1Repository repository.SubBrand1Repository
	// MProductRepository repository.MProductRepository
}

func (service *subBrand1ServiceImpl) Detail(subBrand1Id int, custId string) (response entity.SubBrand1Response, err error) {
	subBrand1, err := service.SubBrand1Repository.FindOneBySubBrand1IdAndCustId(subBrand1Id, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(subBrand1, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *subBrand1ServiceImpl) List(dataFilter entity.SubBrand1QueryFilter, custId string) (data []entity.SubBrand1ListResponse, total int, lastPage int, err error) {
	var subBrand1s []model.SubBrand1

	subBrand1s, total, lastPage, err = service.SubBrand1Repository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBrand1s {
		var vResp entity.SubBrand1ListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	// List
	return data, total, lastPage, err
}

func (service *subBrand1ServiceImpl) LookupList(dataFilter entity.SubBrand1QueryFilter, custId string) (data []entity.SubBrand1LookupResponse, total int, lastPage int, err error) {
	var subBrand1s []model.SubBrand1

	subBrand1s, total, lastPage, err = service.SubBrand1Repository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBrand1s {
		var vResp entity.SubBrand1LookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	// List
	return data, total, lastPage, err
}

func (service *subBrand1ServiceImpl) MatGroupList(dataFilter entity.SubBrand1QueryFilter, custId string) (data []entity.SubBrand1MatGroupResponse, total int, lastPage int, err error) {
	var subBrand1s []model.SubBrand1

	subBrand1s, total, lastPage, err = service.SubBrand1Repository.FindAllByCustIdMatGroupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBrand1s {
		var vResp entity.SubBrand1MatGroupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	// List
	return data, total, lastPage, err
}

func (service *subBrand1ServiceImpl) Store(request entity.CreateSubBrand1Body) (response entity.SubBrand1Response, err error) {

	// subBrand1_code & cust id validation, if err == nil, this means that code & cust id already exists
	subBrand1, err := service.SubBrand1Repository.FindOneBySubBrand1CodeAndCustId(request.Sbrand1Code, request.CustId)
	if err == nil {
		return response, errors.New("sbrand1_code: " + subBrand1.SBrand1Code + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	var subBrand1Data model.SubBrand1
	structs.Automapper(request, &subBrand1Data)
	subBrand1Data.CreatedAt = &timeNow
	subBrand1Data.CreatedBy = &request.CreatedBy
	subBrand1Data.UpdatedAt = &timeNow

	subBrand1Id, err := service.SubBrand1Repository.Store(subBrand1Data)
	if err != nil {
		return response, err
	}

	response.Sbrand1Id = subBrand1Id

	return response, err
}

func (service *subBrand1ServiceImpl) Update(subBrand1Id int, request entity.UpdateSubBrand1Request) (err error) {
	// subBrand1_code & cust id validation, if err == nil and params subBrand1Id != subBrand1.Id, this means that code & cust id already exists
	subBrand1, err := service.SubBrand1Repository.FindOneBySubBrand1CodeAndCustId(request.Sbrand1Code, request.CustId)
	if err == nil && subBrand1.SBrand1Id != subBrand1Id {
		return errors.New("sbrand1_code: " + subBrand1.SBrand1Code + " is already exists")
	}

	err = service.SubBrand1Repository.Update(subBrand1Id, request)
	if err != nil {
		return err
	}

	return err
}

func (service *subBrand1ServiceImpl) Delete(custId string, subBrand1Id int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsKeyExists(subBrand1Id, custId, "subBrand1_id1")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("subBrand1_id is still being used")
	// }

	err = service.SubBrand1Repository.Delete(custId, subBrand1Id, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *subBrand1ServiceImpl) SubBrandList(dataFilter entity.SubBrandQueryFilter, custId string) (data []entity.SubBrandResponse, total int, lastPage int, err error) {
	var subBrand1s []model.SubBrand1

	subBrand1s, total, lastPage, err = service.SubBrand1Repository.FindAllByCustIdSubBrand(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBrand1s {
		var vResp entity.SubBrandResponse
		vResp.BrandId = row.BrandId
		vResp.Sbrand1Id = row.SBrand1Id
		vResp.Sbrand1Code = row.SBrand1Code
		vResp.Sbrand1Name = row.SBrand1Name
		vResp.EffCall = row.EffCall
		vResp.MinItem = row.MinItem
		vResp.IsActive = row.IsActive
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
