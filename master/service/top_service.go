package service

import (
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type TopService interface {
	Detail(int, string) (entity.TopResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.TopResponse, total int, lastPage int, err error)
	Store(entity.CreateTopBody) (entity.TopResponse, error)
	Update(int, entity.UpdateTopRequest) error
	Delete(string, int, int64) error
}

func NewTopService(topRepository repository.TopRepository) *topServiceImpl {
	return &topServiceImpl{
		TopRepository: topRepository,
	}
}

type topServiceImpl struct {
	TopRepository repository.TopRepository
	// MProductRepository repository.MProductRepository
}

func (service *topServiceImpl) Detail(topId int, custId string) (response entity.TopResponse, err error) {
	top, err := service.TopRepository.FindOneByTopAndCustId(topId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(top, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *topServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.TopResponse, total int, lastPage int, err error) {
	tops, total, lastPage, err := service.TopRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(tops) > 0 {
		for _, row := range tops {
			var vResp entity.TopResponse
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *topServiceImpl) Store(request entity.CreateTopBody) (response entity.TopResponse, err error) {

	// top_code & cust id validation, if err == nil, this means that code & cust id already exists
	top, err := service.TopRepository.FindOneByTopAndCustId(request.Top, request.CustId)
	if err == nil {
		return response, fmt.Errorf("top: %d is already exists", top.Top)
	}

	timeNow := time.Now().In(time.UTC)
	var topData model.Top
	structs.Automapper(request, &topData)
	topData.CreatedAt = &timeNow
	topData.CreatedBy = &request.CreatedBy
	topData.UpdatedAt = &timeNow
	topData.UpdatedBy = &request.CreatedBy

	topId, err := service.TopRepository.Store(topData)
	if err != nil {
		return response, err
	}

	response.Top = topId

	return response, err
}

func (service *topServiceImpl) Update(topId int, request entity.UpdateTopRequest) (err error) {

	// top_code & cust id validation, if err == nil and params topId != top.Id, this means that code & cust id already exists
	top, err := service.TopRepository.FindOneByTopAndCustId(request.Top, request.CustId)
	if err == nil && top.Top != topId {
		return fmt.Errorf("top: %d is already exists", top.Top)
	}

	err = service.TopRepository.Update(topId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *topServiceImpl) Delete(custId string, topId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsKeyExists(topId, custId, "top_id1")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("top_id is still being used")
	// }

	err = service.TopRepository.Delete(custId, topId, userId)
	if err != nil {
		return err
	}

	return err
}
