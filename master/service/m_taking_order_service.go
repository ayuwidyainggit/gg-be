package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type TakingOrderService interface {
	Detail(int, string) (entity.TakingOrderResponse, error)
	LookupList(entity.TakingOrderQueryFilter, string) (data []entity.TakingOrderLookupResponse, total int, lastPage int, err error)
	List(entity.TakingOrderQueryFilter, string) (data []entity.TakingOrderResponse, total int, lastPage int, err error)
	Store(entity.CreateTakingOrderBody) (entity.TakingOrderResponse, error)
	Update(int, entity.UpdateTakingOrderRequest) error
	Delete(string, int, int64) error
}

func NewTakingOrderService(takingOrderRepository repository.TakingOrderRepository) *TakingOrderServiceImpl {
	return &TakingOrderServiceImpl{
		TakingOrderRepository: takingOrderRepository,
	}
}

type TakingOrderServiceImpl struct {
	TakingOrderRepository repository.TakingOrderRepository
}

func (service *TakingOrderServiceImpl) Detail(takingOrderId int, custId string) (response entity.TakingOrderResponse, err error) {
	takingOrder, err := service.TakingOrderRepository.FindOneByTakingOrderIdAndCustId(takingOrderId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(takingOrder, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *TakingOrderServiceImpl) List(dataFilter entity.TakingOrderQueryFilter, custId string) (data []entity.TakingOrderResponse, total int, lastPage int, err error) {
	takingOrders, total, lastPage, err := service.TakingOrderRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range takingOrders {
		var vResp entity.TakingOrderResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *TakingOrderServiceImpl) LookupList(dataFilter entity.TakingOrderQueryFilter, custId string) (data []entity.TakingOrderLookupResponse, total int, lastPage int, err error) {
	takingOrders, total, lastPage, err := service.TakingOrderRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range takingOrders {
		var vResp entity.TakingOrderLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *TakingOrderServiceImpl) Store(request entity.CreateTakingOrderBody) (response entity.TakingOrderResponse, err error) {
	// validate reason (taking_order_name) exist or not
	reason, err := service.TakingOrderRepository.FindOneByReasonAndCustId(request.TakingOrderName, request.CustId)
	if err == nil {
		return response, errors.New("reason: " + reason.TakingOrderName + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	takingOrderData := model.TakingOrder{
		CustId:          request.CustId,
		TakingOrderName: request.TakingOrderName,
		ImageUrl:        request.ImageUrl,
		IsActive:        request.IsActive,
		CreatedAt:       &timeNow,
		CreatedBy:       &request.CreatedBy,
		UpdatedAt:       &timeNow,
		UpdatedBy:       &request.CreatedBy,
	}

	takingOrderId, err := service.TakingOrderRepository.Store(takingOrderData)
	if err != nil {
		return response, err
	}

	response.TakingOrderId = takingOrderId

	return response, err
}

func (service *TakingOrderServiceImpl) Update(takingOrderId int, request entity.UpdateTakingOrderRequest) (err error) {

	// validate reason (taking_order_name) exist or not
	reason, err := service.TakingOrderRepository.FindOneByReasonAndCustId(request.TakingOrderName, request.CustId)
	if err == nil && reason.TakingOrderId != takingOrderId {
		return errors.New("reason: " + reason.TakingOrderName + " is already exists")
	}

	err = service.TakingOrderRepository.Update(takingOrderId, request)

	if err != nil {
		return err
	}

	return err
}

func (service *TakingOrderServiceImpl) Delete(custId string, takingOrderId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(takingOrderId, custId, "taking_order_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("taking_order_id is still being used")
	// }

	err = service.TakingOrderRepository.Delete(custId, takingOrderId, userId)
	if err != nil {
		return err
	}

	return err
}
