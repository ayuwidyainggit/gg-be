package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type RejectReasonService interface {
	FindParentCustId(string) (entity.MCustomerResp, error)
	Detail(int, string) (entity.RejectReasonResponse, error)
	LookupList(entity.RejectReasonQueryFilter, string) (data []entity.RejectReasonLookupResponse, total int, lastPage int, err error)
	List(entity.RejectReasonQueryFilter, string) (data []entity.RejectReasonResponse, total int, lastPage int, err error)
	Store(entity.CreateRejectReasonBody) (entity.RejectReasonResponse, error)
	Update(int, entity.UpdateRejectReasonRequest) error
	Delete(string, int, int64) error
}

func NewRejectReasonService(rejectReasonRepository repository.RejectReasonRepository) *RejectReasonServiceImpl {
	return &RejectReasonServiceImpl{
		RejectReasonRepository: rejectReasonRepository,
	}
}

type RejectReasonServiceImpl struct {
	RejectReasonRepository repository.RejectReasonRepository
}

func (service *RejectReasonServiceImpl) FindParentCustId(custId string) (response entity.MCustomerResp, err error) {
	mCustomer, err := service.RejectReasonRepository.FindOneParentCustId(custId)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(mCustomer, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *RejectReasonServiceImpl) Detail(rejectReasonId int, custId string) (response entity.RejectReasonResponse, err error) {
	rejectReason, err := service.RejectReasonRepository.FindOneByRejectReasonIdAndCustId(rejectReasonId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(rejectReason, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *RejectReasonServiceImpl) List(dataFilter entity.RejectReasonQueryFilter, custId string) (data []entity.RejectReasonResponse, total int, lastPage int, err error) {
	rejectReasons, total, lastPage, err := service.RejectReasonRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rejectReasons {
		var vResp entity.RejectReasonResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *RejectReasonServiceImpl) LookupList(dataFilter entity.RejectReasonQueryFilter, custId string) (data []entity.RejectReasonLookupResponse, total int, lastPage int, err error) {
	rejectReasons, total, lastPage, err := service.RejectReasonRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rejectReasons {
		var vResp entity.RejectReasonLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *RejectReasonServiceImpl) Store(request entity.CreateRejectReasonBody) (response entity.RejectReasonResponse, err error) {
	rejectReason, err := service.RejectReasonRepository.FindOneByRejectReasonCodeAndCustId(request.RejectReasonCode, request.CustId)
	if err == nil {
		return response, errors.New("reject_reason_code: " + rejectReason.RejectReasonCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	rejectReasonData := model.RejectReason{
		CustId:           request.CustId,
		RejectReasonCode: request.RejectReasonCode,
		RejectReasonName: request.RejectReasonName,
		IsActive:         request.IsActive,
		CreatedAt:        &timeNow,
		CreatedBy:        &request.CreatedBy,
		UpdatedAt:        &timeNow,
		UpdatedBy:        &request.CreatedBy,
	}

	rejectReasonId, err := service.RejectReasonRepository.Store(rejectReasonData)
	if err != nil {
		return response, err
	}

	response.RejectReasonId = rejectReasonId

	return response, err
}

func (service *RejectReasonServiceImpl) Update(rejectReasonId int, request entity.UpdateRejectReasonRequest) (err error) {

	rejectReason, err := service.RejectReasonRepository.FindOneByRejectReasonCodeAndCustId(request.RejectReasonCode, request.CustId)
	if err == nil && rejectReason.RejectReasonId != rejectReasonId {
		return errors.New("reject_reason_code: " + rejectReason.RejectReasonCode + " is already exists")
	}

	err = service.RejectReasonRepository.Update(rejectReasonId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *RejectReasonServiceImpl) Delete(custId string, rejectReasonId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(rejectReasonId, custId, "reject_reason_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("reject_reason_id is still being used")
	// }

	err = service.RejectReasonRepository.Delete(custId, rejectReasonId, userId)
	if err != nil {
		return err
	}

	return err
}
