package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SkipReasonService interface {
	Detail(int, string) (entity.SkipReasonResponse, error)
	LookupList(entity.SkipReasonQueryFilter, string) (data []entity.SkipReasonLookupResponse, total int, lastPage int, err error)
	List(entity.SkipReasonQueryFilter, string) (data []entity.SkipReasonResponse, total int, lastPage int, err error)
	Store(entity.CreateSkipReasonBody) (entity.SkipReasonResponse, error)
	Update(int, entity.UpdateSkipReasonRequest) error
	Delete(string, int, int64) error
}

func NewSkipReasonService(skipReasonRepository repository.SkipReasonRepository) *SkipReasonServiceImpl {
	return &SkipReasonServiceImpl{
		SkipReasonRepository: skipReasonRepository,
	}
}

type SkipReasonServiceImpl struct {
	SkipReasonRepository repository.SkipReasonRepository
}

func (service *SkipReasonServiceImpl) Detail(skipReasonId int, custId string) (response entity.SkipReasonResponse, err error) {
	skipReason, err := service.SkipReasonRepository.FindOneBySkipReasonIdAndCustId(skipReasonId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(skipReason, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *SkipReasonServiceImpl) List(dataFilter entity.SkipReasonQueryFilter, custId string) (data []entity.SkipReasonResponse, total int, lastPage int, err error) {
	skipReasons, total, lastPage, err := service.SkipReasonRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range skipReasons {
		var vResp entity.SkipReasonResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SkipReasonServiceImpl) LookupList(dataFilter entity.SkipReasonQueryFilter, custId string) (data []entity.SkipReasonLookupResponse, total int, lastPage int, err error) {
	skipReasons, total, lastPage, err := service.SkipReasonRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range skipReasons {
		var vResp entity.SkipReasonLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SkipReasonServiceImpl) Store(request entity.CreateSkipReasonBody) (response entity.SkipReasonResponse, err error) {

	skipReason, err := service.SkipReasonRepository.FindOneBySkipReasonCodeAndCustId(request.SkipReasonCode, request.CustId)
	if err == nil {
		return response, errors.New("skip_reason_code: " + skipReason.SkipReasonCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	skipReasonData := model.SkipReason{
		CustId:         request.CustId,
		SkipReasonCode: request.SkipReasonCode,
		SkipReasonName: request.SkipReasonName,
		IsActive:       request.IsActive,
		CreatedAt:      &timeNow,
		CreatedBy:      &request.CreatedBy,
		UpdatedAt:      &timeNow,
		UpdatedBy:      &request.CreatedBy,
	}

	skipReasonId, err := service.SkipReasonRepository.Store(skipReasonData)
	if err != nil {
		return response, err
	}

	response.SkipReasonId = skipReasonId

	return response, err
}

func (service *SkipReasonServiceImpl) Update(skipReasonId int, request entity.UpdateSkipReasonRequest) (err error) {

	skipReason, err := service.SkipReasonRepository.FindOneBySkipReasonCodeAndCustId(request.SkipReasonCode, request.CustId)
	if err == nil && skipReason.SkipReasonId != skipReasonId {
		return errors.New("skip_reason_code: " + skipReason.SkipReasonCode + " is already exists")
	}

	err = service.SkipReasonRepository.Update(skipReasonId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *SkipReasonServiceImpl) Delete(custId string, skipReasonId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(skipReasonId, custId, "skip_reason_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("skip_reason_id is still being used")
	// }

	err = service.SkipReasonRepository.Delete(custId, skipReasonId, userId)
	if err != nil {
		return err
	}

	return err
}
