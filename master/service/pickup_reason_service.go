package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PickupReasonService interface {
	Detail(int, string) (entity.PickupReasonResponse, error)
	// LookupList(entity.EmployeeQueryFilter, string) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupList(entity.PickupReasonQueryFilter, string) (data []entity.PickupReasonLookupResponse, total int, lastPage int, err error)
	List(entity.PickupReasonQueryFilter, string) (data []entity.PickupReasonResponse, total int, lastPage int, err error)
	Store(entity.CreatePickupReasonBody) (entity.PickupReasonResponse, error)
	Update(int, entity.UpdatePickupReasonRequest) error
	Delete(string, int, int64) error
}

func NewPickupReasonService(pickupReasonRepository repository.PickupReasonRepository) *pickupReasonServiceImpl {
	return &pickupReasonServiceImpl{
		PickupReasonRepository: pickupReasonRepository,
	}
}

type pickupReasonServiceImpl struct {
	PickupReasonRepository repository.PickupReasonRepository
}

func (service *pickupReasonServiceImpl) Detail(pickupReasonId int, custId string) (response entity.PickupReasonResponse, err error) {
	pickupReason, err := service.PickupReasonRepository.FindOneByPickupReasonIdAndCustId(pickupReasonId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(pickupReason, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *pickupReasonServiceImpl) LookupList(dataFilter entity.PickupReasonQueryFilter, custId string) (data []entity.PickupReasonLookupResponse, total int, lastPage int, err error) {
	var pickupReasons []model.PickupReason

	pickupReasons, total, lastPage, err = service.PickupReasonRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	// pickupReasonsDebug, _ := json.Marshal(pickupReasons)
	// fmt.Println("pickupReasonDebug:", string(pickupReasonsDebug))

	for _, row := range pickupReasons {
		var vResp entity.PickupReasonLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *pickupReasonServiceImpl) List(dataFilter entity.PickupReasonQueryFilter, custId string) (data []entity.PickupReasonResponse, total int, lastPage int, err error) {

	pickupReasons, total, lastPage, err := service.PickupReasonRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range pickupReasons {
		var vResp entity.PickupReasonResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *pickupReasonServiceImpl) Store(request entity.CreatePickupReasonBody) (response entity.PickupReasonResponse, err error) {

	// pickup_reason_code & cust id validation, if err == nil, this means that code & cust id already exists
	pickupReason, err := service.PickupReasonRepository.FindOneByPickupReasonCodeAndCustId(request.PickupReasonCode, request.CustId)
	if err == nil {
		return response, errors.New("pickup_reason_code: " + pickupReason.PickupReasonCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	pickupReasonData := model.PickupReason{
		CustId:           request.CustId,
		PickupReasonCode: request.PickupReasonCode,
		PickupReasonName: request.PickupReasonName,
		IsActive:         request.IsActive,
		CreatedAt:        &timeNow,
		CreatedBy:        &request.CreatedBy,
		UpdatedAt:        &timeNow,
		UpdatedBy:        &request.CreatedBy,
	}

	pickupReasonId, err := service.PickupReasonRepository.Store(pickupReasonData)
	if err != nil {
		return response, err
	}

	response.PickupReasonId = pickupReasonId

	return response, err
}

func (service *pickupReasonServiceImpl) Update(pickupReasonId int, request entity.UpdatePickupReasonRequest) (err error) {

	// pickup_reason_code & cust id validation, if err == nil and params pickupReasonId != pickupReason.Id, this means that code & cust id already exists
	pickupReason, err := service.PickupReasonRepository.FindOneByPickupReasonCodeAndCustId(request.PickupReasonCode, request.CustId)
	if err == nil && pickupReason.PickupReasonId != pickupReasonId {
		return errors.New("pickup_reason_code: " + pickupReason.PickupReasonCode + " is already exists")
	}

	err = service.PickupReasonRepository.Update(pickupReasonId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *pickupReasonServiceImpl) Delete(custId string, pickupReasonId int, userId int64) (err error) {

	err = service.PickupReasonRepository.Delete(custId, pickupReasonId, userId)
	if err != nil {
		return err
	}

	return err
}
