package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ReturnReasonDistributorService interface {
	Detail(int, string) (entity.ReturnReasonDistributorResponse, error)
	// LookupList(entity.EmployeeQueryFilter, string) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupList(entity.ReturnReasonDistributorQueryFilter, string) (data []entity.ReturnReasonDistributorLookupResponse, total int, lastPage int, err error)
	List(entity.ReturnReasonDistributorQueryFilter, string) (data []entity.ReturnReasonDistributorResponse, total int, lastPage int, err error)
	Store(entity.CreateReturnReasonDistributorBody) (entity.ReturnReasonDistributorResponse, error)
	Update(int, entity.UpdateReturnReasonDistributorRequest) error
	Delete(string, int, int64) error
}

func NewReturnReasonDistributorService(returnReasonDistributorRepository repository.ReturnReasonDistributorRepository) *returnReasonDistributorServiceImpl {
	return &returnReasonDistributorServiceImpl{
		ReturnReasonDistributorRepository: returnReasonDistributorRepository,
	}
}

type returnReasonDistributorServiceImpl struct {
	ReturnReasonDistributorRepository repository.ReturnReasonDistributorRepository
}

func (service *returnReasonDistributorServiceImpl) Detail(returnReasonDistributorId int, custId string) (response entity.ReturnReasonDistributorResponse, err error) {
	returnReasonDistributor, err := service.ReturnReasonDistributorRepository.FindOneByReturnReasonDistributorIdAndCustId(returnReasonDistributorId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(returnReasonDistributor, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *returnReasonDistributorServiceImpl) LookupList(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) (data []entity.ReturnReasonDistributorLookupResponse, total int, lastPage int, err error) {
	var returnReasonDistributors []model.ReturnReasonDistributor

	returnReasonDistributors, total, lastPage, err = service.ReturnReasonDistributorRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	// returnReasonDistributorsDebug, _ := json.Marshal(returnReasonDistributors)
	// fmt.Println("returnReasonDistributorDebug:", string(returnReasonDistributorsDebug))

	for _, row := range returnReasonDistributors {
		var vResp entity.ReturnReasonDistributorLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnReasonDistributorServiceImpl) List(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) (data []entity.ReturnReasonDistributorResponse, total int, lastPage int, err error) {

	returnReasonDistributors, total, lastPage, err := service.ReturnReasonDistributorRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range returnReasonDistributors {
		var vResp entity.ReturnReasonDistributorResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnReasonDistributorServiceImpl) Store(request entity.CreateReturnReasonDistributorBody) (response entity.ReturnReasonDistributorResponse, err error) {

	// return_reason_code & cust id validation, if err == nil, this means that code & cust id already exists
	returnReasonDistributor, err := service.ReturnReasonDistributorRepository.FindOneByReturnReasonDistributorCodeAndCustId(request.ReturnReasonDistributorCode, request.CustId)
	if err == nil {
		return response, errors.New("return_reason_code: " + returnReasonDistributor.ReturnReasonDistributorCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	returnReasonDistributorData := model.ReturnReasonDistributor{
		CustId:                      request.CustId,
		ReturnReasonDistributorCode: request.ReturnReasonDistributorCode,
		ReturnReasonDistributorName: request.ReturnReasonDistributorName,
		ReturnReasonDistributorType: request.ReturnReasonDistributorType,
		IsActive:                    request.IsActive,
		CreatedAt:                   &timeNow,
		CreatedBy:                   &request.CreatedBy,
		UpdatedAt:                   &timeNow,
		UpdatedBy:                   &request.CreatedBy,
	}

	returnReasonDistributorId, err := service.ReturnReasonDistributorRepository.Store(returnReasonDistributorData)
	if err != nil {
		return response, err
	}

	response.ReturnReasonDistributorId = returnReasonDistributorId

	return response, err
}

func (service *returnReasonDistributorServiceImpl) Update(returnReasonDistributorId int, request entity.UpdateReturnReasonDistributorRequest) (err error) {

	// return_reason_code & cust id validation, if err == nil and params returnReasonDistributorId != returnReasonDistributor.Id, this means that code & cust id already exists
	returnReasonDistributor, err := service.ReturnReasonDistributorRepository.FindOneByReturnReasonDistributorCodeAndCustId(request.ReturnReasonDistributorCode, request.CustId)
	if err == nil && returnReasonDistributor.ReturnReasonDistributorId != returnReasonDistributorId {
		return errors.New("return_reason_code: " + returnReasonDistributor.ReturnReasonDistributorCode + " is already exists")
	}

	err = service.ReturnReasonDistributorRepository.Update(returnReasonDistributorId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *returnReasonDistributorServiceImpl) Delete(custId string, returnReasonDistributorId int, userId int64) (err error) {

	err = service.ReturnReasonDistributorRepository.Delete(custId, returnReasonDistributorId, userId)
	if err != nil {
		return err
	}

	return err
}
