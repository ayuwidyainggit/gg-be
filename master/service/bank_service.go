package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type BankService interface {
	Detail(int, string) (entity.BankResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.BankResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.BankLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateBankBody) (entity.BankResponse, error)
	Update(int, entity.UpdateBankRequest) error
	Delete(string, int, int64) error

	LookupListOutletBank(dataFilter entity.QueryFilterOutletBank, custId string) (data []entity.BankLookupResponse, total int, lastPage int, err error)
	LookupListOutletBankByBankID(dataFilter entity.QueryFilterOutletBank, custId string) (data []entity.OutletBankList, total int, lastPage int, err error)
}

func NewBankService(bankRepository repository.BankRepository) *bankServiceImpl {
	return &bankServiceImpl{
		BankRepository: bankRepository,
	}
}

type bankServiceImpl struct {
	BankRepository repository.BankRepository
}

func (service *bankServiceImpl) Detail(bankId int, custId string) (response entity.BankResponse, err error) {
	bank, err := service.BankRepository.FindOneByBankIdAndCustId(bankId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(bank, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *bankServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.BankResponse, total int, lastPage int, err error) {
	banks, total, lastPage, err := service.BankRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range banks {
		var vResp entity.BankResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *bankServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.BankLookupResponse, total int, lastPage int, err error) {
	banks, total, lastPage, err := service.BankRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range banks {
		var vResp entity.BankLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *bankServiceImpl) Store(request entity.CreateBankBody) (response entity.BankResponse, err error) {

	// bank_code & cust id validation, if err == nil, this means that code & cust id already exists
	bank, err := service.BankRepository.FindOneByBankCodeAndCustId(request.BankCode, request.CustId)
	if err == nil {
		return response, errors.New("bank_code: " + bank.BankCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	bankData := model.Bank{
		CustId:    request.CustId,
		BankCode:  request.BankCode,
		BankName:  request.BankName,
		IsActive:  request.IsActive,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
		UpdatedBy: &request.CreatedBy,
	}

	bankId, err := service.BankRepository.Store(bankData)
	if err != nil {
		return response, err
	}

	response.BankId = bankId

	return response, err
}

func (service *bankServiceImpl) Update(bankId int, request entity.UpdateBankRequest) (err error) {

	// bank_code & cust id validation, if err == nil and params bankId != bank.Id, this means that code & cust id already exists
	bank, err := service.BankRepository.FindOneByBankCodeAndCustId(request.BankCode, request.CustId)
	if err == nil && bank.BankId != bankId {
		return errors.New("bank_code: " + bank.BankCode + " is already exists")
	}

	err = service.BankRepository.Update(bankId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *bankServiceImpl) Delete(custId string, bankId int, userId int64) (err error) {

	err = service.BankRepository.Delete(custId, bankId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *bankServiceImpl) LookupListOutletBank(dataFilter entity.QueryFilterOutletBank, custId string) (data []entity.BankLookupResponse, total int, lastPage int, err error) {
	banks, total, lastPage, err := service.BankRepository.FindDistictOutletBankByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range banks {
		var vResp entity.BankLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *bankServiceImpl) LookupListOutletBankByBankID(dataFilter entity.QueryFilterOutletBank, custId string) (data []entity.OutletBankList, total int, lastPage int, err error) {
	banks, total, lastPage, err := service.BankRepository.FindBankOutletByBankIdAndCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range banks {
		var vResp entity.OutletBankList
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
