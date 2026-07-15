package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type MarketService interface {
	Detail(int, string) (entity.MarketResponse, error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.MarketLookupResponse, total int, lastPage int, err error)
	List(entity.GeneralQueryFilter, string) (data []entity.MarketResponse, total int, lastPage int, err error)
	Store(entity.CreateMarketBody) (entity.MarketResponse, error)
	Update(int, entity.UpdateMarketRequest) error
	Delete(string, int, int64) error
}

func NewMarketService(marketRepository repository.MarketRepository) *marketServiceImpl {
	return &marketServiceImpl{
		MarketRepository: marketRepository,
	}
}

type marketServiceImpl struct {
	MarketRepository repository.MarketRepository
}

func (service *marketServiceImpl) Detail(marketId int, custId string) (response entity.MarketResponse, err error) {
	market, err := service.MarketRepository.FindOneByMarketIdAndCustId(marketId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(market, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *marketServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MarketResponse, total int, lastPage int, err error) {
	markets, total, lastPage, err := service.MarketRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range markets {
		var vResp entity.MarketResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *marketServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MarketLookupResponse, total int, lastPage int, err error) {
	markets, total, lastPage, err := service.MarketRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range markets {
		var vResp entity.MarketLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *marketServiceImpl) Store(request entity.CreateMarketBody) (response entity.MarketResponse, err error) {

	// market_code & cust id validation, if err == nil, this means that code & cust id already exists
	market, err := service.MarketRepository.FindOneByMarketCodeAndCustId(request.MarketCode, request.CustId)
	if err == nil {
		return response, errors.New("market_code: " + market.MarketCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	marketData := model.Market{
		CustId:     request.CustId,
		MarketCode: request.MarketCode,
		MarketName: request.MarketName,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	marketId, err := service.MarketRepository.Store(marketData)
	if err != nil {
		return response, err
	}

	response.MarketId = marketId

	return response, err
}

func (service *marketServiceImpl) Update(marketId int, request entity.UpdateMarketRequest) (err error) {

	// market_code & cust id validation, if err == nil and params marketId != market.Id, this means that code & cust id already exists
	market, err := service.MarketRepository.FindOneByMarketCodeAndCustId(request.MarketCode, request.CustId)
	if err == nil && market.MarketId != marketId {
		return errors.New("market_code: " + market.MarketCode + " is already exists")
	}

	err = service.MarketRepository.Update(marketId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *marketServiceImpl) Delete(custId string, marketId int, userId int64) (err error) {

	err = service.MarketRepository.Delete(custId, marketId, userId)
	if err != nil {
		return err
	}

	return err
}
