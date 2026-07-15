package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SupplierService interface {
	Detail(int, string) (entity.SupplierResponse, error)
	List(entity.SupplierQueryFilter, string) (data []entity.SupplierResponse, total int, lastPage int, err error)
	LookupList(entity.SupplierQueryFilter, entity.SupplierLookupScope) (data []entity.SupplierLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateSupplierBody) (entity.SupplierResponse, error)
	Update(int, entity.UpdateSupplierRequest) error
	Delete(string, int, int64) error
}

func NewSupplierService(supplierRepository repository.SupplierRepository) *supplierServiceImpl {
	return &supplierServiceImpl{
		SupplierRepository: supplierRepository,
	}
}

type supplierServiceImpl struct {
	SupplierRepository repository.SupplierRepository
}

func (service *supplierServiceImpl) Detail(supplierId int, custId string) (response entity.SupplierResponse, err error) {
	supplier, err := service.SupplierRepository.FindOneBySupplierIdAndCustId(supplierId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(supplier, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *supplierServiceImpl) List(dataFilter entity.SupplierQueryFilter, custId string) (data []entity.SupplierResponse, total int, lastPage int, err error) {
	suppliers, total, lastPage, err := service.SupplierRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range suppliers {
		var vResp entity.SupplierResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *supplierServiceImpl) LookupList(dataFilter entity.SupplierQueryFilter, scope entity.SupplierLookupScope) (data []entity.SupplierLookupResponse, total int, lastPage int, err error) {
	suppliers, total, lastPage, err := service.SupplierRepository.FindAllLookupByCustId(dataFilter, scope)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range suppliers {
		var vResp entity.SupplierLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *supplierServiceImpl) Store(request entity.CreateSupplierBody) (response entity.SupplierResponse, err error) {

	log.Println("request:", structs.StructToJson(request))

	// sup_code & cust id validation, if err == nil, this means that code & cust id already exists
	supplier, err := service.SupplierRepository.FindOneBySupplierCodeAndCustId(request.SupplierCode, request.CustId)
	if err == nil {
		return response, errors.New("sup_code: " + supplier.SupplierCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	supplierData := model.Supplier{}
	err = structs.Automapper(request, &supplierData)
	if err != nil {
		return response, err
	}
	supplierData.CustId = request.CustId
	supplierData.IsActive = request.IsActive
	supplierData.CreatedAt = &timeNow
	supplierData.CreatedBy = &request.CreatedBy
	supplierData.UpdatedAt = &timeNow
	supplierData.UpdatedBy = &request.CreatedBy

	log.Println("supplierData:", structs.StructToJson(supplierData))

	supplierId, err := service.SupplierRepository.Store(supplierData)
	if err != nil {
		return response, err
	}

	response.SupplierId = supplierId

	return response, err
}

func (service *supplierServiceImpl) Update(supplierId int, request entity.UpdateSupplierRequest) (err error) {

	// sup_code & cust id validation, if err == nil and params supplierId != supplier.Id, this means that code & cust id already exists
	supplier, err := service.SupplierRepository.FindOneBySupplierCodeAndCustId(request.SupplierCode, request.CustId)
	if err == nil && supplier.SupplierId != supplierId {
		return errors.New("sup_code: " + supplier.SupplierCode + " is already exists")
	}

	err = service.SupplierRepository.Update(supplierId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *supplierServiceImpl) Delete(custId string, supplierId int, userId int64) (err error) {

	err = service.SupplierRepository.Delete(custId, supplierId, userId)
	if err != nil {
		return err
	}

	return err
}
