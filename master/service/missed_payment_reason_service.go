package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type MissedPaymentReasonsService interface {
	Detail(int, string) (entity.MissedPaymentReasonsResponse, error)
	LookupList(entity.MissedPaymentReasonsQueryFilter, string) (data []entity.MissedPaymentReasonsLookupResponse, total int, lastPage int, err error)
	List(entity.MissedPaymentReasonsQueryFilter, string) (data []entity.MissedPaymentReasonsResponse, total int, lastPage int, err error)
	Store(entity.CreateMissedPaymentReasonsBody) (entity.MissedPaymentReasonsResponse, error)
	Update(int, entity.UpdateMissedPaymentReasonsRequest) error
	Delete(string, int, int64) error
}

func NewMissedPaymentReasonsService(MissedPaymentReasonsRepository repository.MissedPaymentReasonsRepository) *MissedPaymentReasonsServiceImpl {
	return &MissedPaymentReasonsServiceImpl{
		MissedPaymentReasonsRepository: MissedPaymentReasonsRepository,
	}
}

type MissedPaymentReasonsServiceImpl struct {
	MissedPaymentReasonsRepository repository.MissedPaymentReasonsRepository
}

func (service *MissedPaymentReasonsServiceImpl) Detail(MissedPaymentReasonsId int, custId string) (response entity.MissedPaymentReasonsResponse, err error) {
	MissedPaymentReasons, err := service.MissedPaymentReasonsRepository.FindOneByMissedPaymentReasonsIdAndCustId(MissedPaymentReasonsId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(MissedPaymentReasons, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *MissedPaymentReasonsServiceImpl) List(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) (data []entity.MissedPaymentReasonsResponse, total int, lastPage int, err error) {
	MissedPaymentReasonss, total, lastPage, err := service.MissedPaymentReasonsRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MissedPaymentReasonss {
		var vResp entity.MissedPaymentReasonsResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MissedPaymentReasonsServiceImpl) LookupList(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) (data []entity.MissedPaymentReasonsLookupResponse, total int, lastPage int, err error) {
	MissedPaymentReasonss, total, lastPage, err := service.MissedPaymentReasonsRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MissedPaymentReasonss {
		var vResp entity.MissedPaymentReasonsLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MissedPaymentReasonsServiceImpl) Store(request entity.CreateMissedPaymentReasonsBody) (response entity.MissedPaymentReasonsResponse, err error) {
	// validate reason (taking_order_name) exist or not
	reason, err := service.MissedPaymentReasonsRepository.FindOneByReasonAndCustId(request.MissedPaymentReasonsCode, request.CustId)
	if err == nil {
		return response, errors.New("reason: " + reason.MissedPaymentReasonsName + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	MissedPaymentReasonsData := model.MissedPaymentReasons{
		CustId:                   request.CustId,
		MissedPaymentReasonsCode: request.MissedPaymentReasonsCode,
		MissedPaymentReasonsName: request.MissedPaymentReasonsName,
		ImageUrl:                 request.ImageUrl,
		IsActive:                 request.IsActive,
		CreatedAt:                &timeNow,
		CreatedBy:                &request.CreatedBy,
		UpdatedAt:                &timeNow,
		UpdatedBy:                &request.CreatedBy,
	}

	MissedPaymentReasonsId, err := service.MissedPaymentReasonsRepository.Store(MissedPaymentReasonsData)
	if err != nil {
		return response, err
	}

	response.MissedPaymentReasonsId = MissedPaymentReasonsId

	return response, err
}

func (service *MissedPaymentReasonsServiceImpl) Update(MissedPaymentReasonsId int, request entity.UpdateMissedPaymentReasonsRequest) (err error) {

	// validate reason (taking_order_name) exist or not
	reason, err := service.MissedPaymentReasonsRepository.FindOneByReasonAndCustId(request.MissedPaymentReasonsName, request.CustId)
	if err == nil && reason.MissedPaymentReasonsId != MissedPaymentReasonsId {
		return errors.New("reason: " + reason.MissedPaymentReasonsName + " is already exists")
	}

	err = service.MissedPaymentReasonsRepository.Update(MissedPaymentReasonsId, request)

	if err != nil {
		return err
	}

	return err
}

func (service *MissedPaymentReasonsServiceImpl) Delete(custId string, MissedPaymentReasonsId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(MissedPaymentReasonsId, custId, "taking_order_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("taking_order_id is still being used")
	// }

	err = service.MissedPaymentReasonsRepository.Delete(custId, MissedPaymentReasonsId, userId)
	if err != nil {
		return err
	}

	return err
}
