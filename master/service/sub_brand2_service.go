package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type SubBrand2Service interface {
	Detail(int, string) (entity.SubBrand2Response, error)
	List(entity.GeneralQueryFilter, string) (data []entity.SubBrand2ListResponse, total int, lastPage int, err error)
	Store(entity.CreateSubBrand2Body) (entity.SubBrand2Response, error)
	Update(int, entity.UpdateSubBrand2Request) error
	Delete(string, int, int64) error
}

func NewSubBrand2Service(subBrand2Repository repository.SubBrand2Repository, mProductRepository repository.ProductRepository) *subBrand2ServiceImpl {
	return &subBrand2ServiceImpl{
		SubBrand2Repository: subBrand2Repository,
		MProductRepository:  mProductRepository,
	}
}

type subBrand2ServiceImpl struct {
	SubBrand2Repository repository.SubBrand2Repository
	MProductRepository  repository.ProductRepository
}

func (service *subBrand2ServiceImpl) Detail(sbrand2Id int, custId string) (response entity.SubBrand2Response, err error) {
	subBrand2, err := service.SubBrand2Repository.FindOneBySBrand2IdAndCustId(sbrand2Id, custId)
	if err != nil {
		return response, err
	}

	response.SBrand2Id = subBrand2.SBrand2Id
	response.SBrand2Code = subBrand2.SBrand2Code
	response.SBrand2Name = subBrand2.SBrand2Name
	response.IsActive = subBrand2.IsActive
	response.UpdatedBy = subBrand2.UpdatedBy
	response.UpdatedAt = subBrand2.UpdatedAt

	return response, err
}

func (service *subBrand2ServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SubBrand2ListResponse, total int, lastPage int, err error) {
	var subBrand2s []model.SubBrand2

	switch dataFilter.Mode {
	case "lookup":
		subBrand2s, total, lastPage, err = service.SubBrand2Repository.FindAllByCustIdLookupMode(dataFilter, custId)
		if err != nil {
			return data, total, lastPage, err
		}
		if len(subBrand2s) > 0 {
			for _, row := range subBrand2s {
				sb2Model := entity.SubBrand2ListResponse{
					SBrand2Id:   row.SBrand2Id,
					SBrand2Code: row.SBrand2Code,
					SBrand2Name: row.SBrand2Name,
				}
				if row.UpdatedByName != nil {
					sb2Model.UpdatedByName = *row.UpdatedByName
				}
				data = append(data, sb2Model)
			}
		}
	default:
		subBrand2s, total, lastPage, err = service.SubBrand2Repository.FindAllByCustId(dataFilter, custId)
		if err != nil {
			return data, total, lastPage, err
		}
		if len(subBrand2s) > 0 {
			for _, row := range subBrand2s {
				sb2Model := entity.SubBrand2ListResponse{
					SBrand2Id:   row.SBrand2Id,
					SBrand2Code: row.SBrand2Code,
					SBrand2Name: row.SBrand2Name,
					IsActive:    row.IsActive,
					UpdatedBy:   row.UpdatedBy,
					UpdatedAt:   row.UpdatedAt,
				}
				if row.UpdatedByName != nil {
					sb2Model.UpdatedByName = *row.UpdatedByName
				}
				data = append(data, sb2Model)
			}
		}
	}
	return data, total, lastPage, err
}

func (service *subBrand2ServiceImpl) Store(request entity.CreateSubBrand2Body) (response entity.SubBrand2Response, err error) {

	// sbrand2_code & cust id validation, if err == nil, this means that code & cust id already exists
	subBrand2, err := service.SubBrand2Repository.FindOneBySBrand2CodeAndCustId(request.SBrand2Code, request.CustId)
	if err == nil {
		return response, errors.New("sbrand2_code: " + subBrand2.SBrand2Code + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	consPro := model.SubBrand2{
		CustId:      request.CustId,
		SBrand2Code: request.SBrand2Code,
		SBrand2Name: request.SBrand2Name,
		IsActive:    request.IsActive,
		CreatedAt:   &timeNow,
		CreatedBy:   &request.CreatedBy,
		UpdatedAt:   &timeNow,
		UpdatedBy:   &request.CreatedBy,
	}

	sbrand2Id, err := service.SubBrand2Repository.Store(consPro)
	if err != nil {
		return response, err
	}

	response.SBrand2Id = sbrand2Id

	return response, err
}

func (service *subBrand2ServiceImpl) Update(sbrand2Id int, request entity.UpdateSubBrand2Request) (err error) {

	// sbrand2_code & cust id validation, if err == nil and params sbrand2Id != subBrand2.Id, this means that code & cust id already exists
	subBrand2, err := service.SubBrand2Repository.FindOneBySBrand2CodeAndCustId(request.SBrand2Code, request.CustId)
	if err == nil && subBrand2.SBrand2Id != sbrand2Id {
		return errors.New("sbrand2_code: " + subBrand2.SBrand2Code + " is already exists")
	}

	err = service.SubBrand2Repository.Update(sbrand2Id, request)
	if err != nil {
		return err
	}

	return err
}

func (service *subBrand2ServiceImpl) Delete(custId string, sbrand2Id int, userId int64) (err error) {

	isExists, err := service.MProductRepository.IsExists(sbrand2Id, custId, "pcat_id")
	if err != nil {
		return err
	}

	if isExists {
		return errors.New("Sub Brand2 is still being used")
	}

	err = service.SubBrand2Repository.Delete(custId, sbrand2Id, userId)
	if err != nil {
		return err
	}

	return err
}
