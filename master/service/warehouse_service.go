package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/errmsg"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type WarehouseService interface {
	Detail(int, string) (entity.WarehouseResponse, error)
	List(entity.WarehouseQueryFilter, string) (data []entity.WarehouseResponse, total int, lastPage int, err error)
	LookupList(entity.WarehouseQueryFilter, string) (data []entity.WarehouseLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateWarehouseBody) (entity.WarehouseResponse, error)
	Update(int, entity.UpdateWarehouseRequest) error
	Delete(string, int, int64) error
}

func NewWarehouseService(warehouseRepository repository.WarehouseRepository) *warehouseServiceImpl {
	return &warehouseServiceImpl{
		WarehouseRepository: warehouseRepository,
	}
}

type warehouseServiceImpl struct {
	WarehouseRepository repository.WarehouseRepository
}

func (service *warehouseServiceImpl) Detail(warehouseId int, custId string) (response entity.WarehouseResponse, err error) {
	warehouse, err := service.WarehouseRepository.FindOneByWarehouseIdAndCustId(warehouseId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(warehouse, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *warehouseServiceImpl) List(dataFilter entity.WarehouseQueryFilter, custId string) (data []entity.WarehouseResponse, total int, lastPage int, err error) {
	warehouses, total, lastPage, err := service.WarehouseRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range warehouses {
		var vResp entity.WarehouseResponse
		err := structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *warehouseServiceImpl) LookupList(dataFilter entity.WarehouseQueryFilter, custId string) (data []entity.WarehouseLookupResponse, total int, lastPage int, err error) {
	var warehouses []model.Warehouse
	if len(dataFilter.DistributorIDs) > 0 {
		warehouses, total, lastPage, err = service.WarehouseRepository.FindLookupByDistributorUnion(dataFilter, dataFilter.DistributorIDs)
	} else {
		warehouses, total, lastPage, err = service.WarehouseRepository.FindAllLookupByCustId(dataFilter, custId)
	}
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range warehouses {
		var vResp entity.WarehouseLookupResponse
		err := structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *warehouseServiceImpl) Store(request entity.CreateWarehouseBody) (response entity.WarehouseResponse, err error) {

	// wh_code & cust id validation, if err == nil, this means that code & cust id already exists
	warehouse, err := service.WarehouseRepository.FindOneByWarehouseCodeAndCustId(request.WarehouseCode, request.CustId)
	if err == nil {
		return response, errors.New("wh_code: " + warehouse.WarehouseCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	warehouseData := model.Warehouse{}

	err = structs.Automapper(request, &warehouseData)
	if err != nil {
		return response, err
	}

	warehouseData.CreatedAt = &timeNow
	warehouseData.CreatedBy = &request.CreatedBy
	warehouseData.UpdatedAt = &timeNow
	warehouseData.UpdatedBy = &request.CreatedBy

	warehouseId, err := service.WarehouseRepository.Store(warehouseData)
	if err != nil {
		return response, err
	}

	response.WarehouseId = warehouseId

	return response, err
}

func (service *warehouseServiceImpl) Update(warehouseId int, request entity.UpdateWarehouseRequest) (err error) {

	// wh_code & cust id validation, if err == nil and params warehouseId != warehouse.Id, this means that code & cust id already exists
	warehouse, err := service.WarehouseRepository.FindOneByWarehouseCodeAndCustId(request.WarehouseCode, request.CustId)
	if err == nil && warehouse.WarehouseId != warehouseId {
		return errors.New("wh_code: " + warehouse.WarehouseCode + " is already exists")
	} else {
		log.Println("warehouse:", structs.StructToJson(warehouse))
	}

	err = service.WarehouseRepository.Update(warehouseId, request)
	if err != nil {
		log.Println("2 err.Error() 2:", err.Error())
		return err
	}

	return err
}

func (service *warehouseServiceImpl) Delete(custId string, warehouseId int, userId int64) (err error) {

	// validasi warehouse yg ada stock tidak bisa di delete
	warehouse, err := service.WarehouseRepository.CountWarehouseInStockByCustId(warehouseId, custId)
	if err == nil && warehouse.TotalWarehouse > 0 {
		return errors.New(errmsg.ERROR_DEL_DATA_NOT_ALLOWED)
	}

	err = service.WarehouseRepository.Delete(custId, warehouseId, userId)
	if err != nil {
		return err
	}

	return err
}
