package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ReturnReasonService interface {
	Detail(int, string) (entity.ReturnReasonResponse, error)
	// LookupList(entity.EmployeeQueryFilter, string) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupList(entity.ReturnReasonQueryFilter, string) (data []entity.ReturnReasonLookupResponse, total int, lastPage int, err error)
	List(entity.ReturnReasonQueryFilter, string) (data []entity.ReturnReasonResponse, total int, lastPage int, err error)
	Store(entity.CreateReturnReasonBody) (entity.ReturnReasonResponse, error)
	Update(int, entity.UpdateReturnReasonRequest) error
	Delete(string, int, int64) error
}

func NewReturnReasonService(returnReasonRepository repository.ReturnReasonRepository) *returnReasonServiceImpl {
	return &returnReasonServiceImpl{
		ReturnReasonRepository: returnReasonRepository,
	}
}

type returnReasonServiceImpl struct {
	ReturnReasonRepository repository.ReturnReasonRepository
}

func (service *returnReasonServiceImpl) Detail(returnReasonId int, custId string) (response entity.ReturnReasonResponse, err error) {
	returnReason, err := service.ReturnReasonRepository.FindOneByReturnReasonIdAndCustId(returnReasonId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(returnReason, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *returnReasonServiceImpl) LookupList(dataFilter entity.ReturnReasonQueryFilter, custId string) (data []entity.ReturnReasonLookupResponse, total int, lastPage int, err error) {
	var returnReasons []model.ReturnReason

	returnReasons, total, lastPage, err = service.ReturnReasonRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	// returnReasonsDebug, _ := json.Marshal(returnReasons)
	// fmt.Println("returnReasonDebug:", string(returnReasonsDebug))

	for _, row := range returnReasons {
		var vResp entity.ReturnReasonLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnReasonServiceImpl) List(dataFilter entity.ReturnReasonQueryFilter, custId string) (data []entity.ReturnReasonResponse, total int, lastPage int, err error) {

	returnReasons, total, lastPage, err := service.ReturnReasonRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range returnReasons {
		var vResp entity.ReturnReasonResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	
	return data, total, lastPage, err
}

func (service *returnReasonServiceImpl) Store(request entity.CreateReturnReasonBody) (response entity.ReturnReasonResponse, err error) {

	// return_reason_code & cust id validation, if err == nil, this means that code & cust id already exists
	returnReason, err := service.ReturnReasonRepository.FindOneByReturnReasonCodeAndCustId(request.ReturnReasonCode, request.CustId)
	if err == nil {
		return response, errors.New("return_reason_code: " + returnReason.ReturnReasonCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	returnReasonData := model.ReturnReason{
		CustId:           request.CustId,
		ReturnReasonCode: request.ReturnReasonCode,
		ReturnReasonName: request.ReturnReasonName,
		ReturnReasonType: request.ReturnReasonType,
		IsActive:         request.IsActive,
		CreatedAt:        &timeNow,
		CreatedBy:        &request.CreatedBy,
		UpdatedAt:        &timeNow,
		UpdatedBy:        &request.CreatedBy,
	}

	returnReasonId, err := service.ReturnReasonRepository.Store(returnReasonData)
	if err != nil {
		return response, err
	}

	response.ReturnReasonId = returnReasonId

	return response, err
}

func (service *returnReasonServiceImpl) Update(returnReasonId int, request entity.UpdateReturnReasonRequest) (err error) {

	// return_reason_code & cust id validation, if err == nil and params returnReasonId != returnReason.Id, this means that code & cust id already exists
	returnReason, err := service.ReturnReasonRepository.FindOneByReturnReasonCodeAndCustId(request.ReturnReasonCode, request.CustId)
	if err == nil && returnReason.ReturnReasonId != returnReasonId {
		return errors.New("return_reason_code: " + returnReason.ReturnReasonCode + " is already exists")
	}

	err = service.ReturnReasonRepository.Update(returnReasonId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *returnReasonServiceImpl) Delete(custId string, returnReasonId int, userId int64) (err error) {

	err = service.ReturnReasonRepository.Delete(custId, returnReasonId, userId)
	if err != nil {
		return err
	}

	return err
}
