package service

import (
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/structs"
	"mobile/repository"
)

type PickupReasonService interface {
	// LookupList(entity.EmployeeQueryFilter, string) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupList(entity.PickupReasonQueryFilter, string) (data []entity.PickupReasonLookupResponse, total int64, lastPage int, err error)
	List(entity.PickupReasonQueryFilter, string) (data []entity.PickupReasonResponse, total int64, lastPage int, err error)
}

func NewPickupReasonService(pickupReasonRepository repository.PickupReasonRepository) *pickupReasonServiceImpl {
	return &pickupReasonServiceImpl{
		PickupReasonRepository: pickupReasonRepository,
	}
}

type pickupReasonServiceImpl struct {
	PickupReasonRepository repository.PickupReasonRepository
}

func (service *pickupReasonServiceImpl) LookupList(dataFilter entity.PickupReasonQueryFilter, custId string) (data []entity.PickupReasonLookupResponse, total int64, lastPage int, err error) {
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

func (service *pickupReasonServiceImpl) List(dataFilter entity.PickupReasonQueryFilter, custId string) (data []entity.PickupReasonResponse, total int64, lastPage int, err error) {

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
