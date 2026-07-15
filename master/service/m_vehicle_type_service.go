package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type VehicleTypeService interface {
	Detail(int, string) (entity.VehicletypeResponse, error)
	LookupList(entity.VehicletypeQueryFilter, string) (data []entity.VehicletypeLookupResponse, total int, lastPage int, err error)
	List(entity.VehicletypeQueryFilter, string) (data []entity.VehicletypeResponse, total int, lastPage int, err error)
	Store(entity.CreateVehicletypeBody) (entity.VehicletypeResponse, error)
	Update(int, entity.UpdateVehicletypeRequest) error
	Delete(string, int, int64) error
}

func NewVehicleTypeService(vehicleTypeRepository repository.VehicleTypeRepository) *VehicleTypeServiceImpl {
	return &VehicleTypeServiceImpl{
		VehicleTypeRepository: vehicleTypeRepository,
	}
}

type VehicleTypeServiceImpl struct {
	VehicleTypeRepository repository.VehicleTypeRepository
}

func (service *VehicleTypeServiceImpl) Detail(vehicleTypeId int, custId string) (response entity.VehicletypeResponse, err error) {
	vehicleType, err := service.VehicleTypeRepository.FindOneByVehicleTypeIdAndCustId(vehicleTypeId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(vehicleType, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *VehicleTypeServiceImpl) List(dataFilter entity.VehicletypeQueryFilter, custId string) (data []entity.VehicletypeResponse, total int, lastPage int, err error) {
	vehicleTypes, total, lastPage, err := service.VehicleTypeRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range vehicleTypes {
		var vResp entity.VehicletypeResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *VehicleTypeServiceImpl) LookupList(dataFilter entity.VehicletypeQueryFilter, custId string) (data []entity.VehicletypeLookupResponse, total int, lastPage int, err error) {
	vehicleTypes, total, lastPage, err := service.VehicleTypeRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range vehicleTypes {
		var vResp entity.VehicletypeLookupResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *VehicleTypeServiceImpl) Store(request entity.CreateVehicletypeBody) (response entity.VehicletypeResponse, err error) {
	_, err = service.VehicleTypeRepository.FindOneByVehicleTypeNameAndCustId(request.VehicleTypeName, request.CustId)
	if err == nil {
		return response, errors.New("Vehicle Type Name: " + request.VehicleTypeName + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	vehicleTypeData := model.VehicleType{
		CustId:          request.CustId,
		VehicleTypeName: request.VehicleTypeName,
		IsActive:        request.IsActive,
		CreatedAt:       &timeNow,
		CreatedBy:       &request.CreatedBy,
		UpdatedAt:       &timeNow,
		UpdatedBy:       &request.CreatedBy,
	}

	vehicleTypeId, err := service.VehicleTypeRepository.Store(vehicleTypeData)
	if err != nil {
		return response, err
	}

	response.VehicleTypeId = vehicleTypeId

	return response, err
}

func (service *VehicleTypeServiceImpl) Update(vehicleTypeId int, request entity.UpdateVehicletypeRequest) (err error) {

	vehicleType, err := service.VehicleTypeRepository.FindOneByVehicleTypeNameAndCustId(request.VehicleTypeName, request.CustId)
	if err == nil && vehicleType.VehicleTypeId != vehicleTypeId {
		return errors.New("Vehicle Type Name: " + request.VehicleTypeName + " is already exists")
	}

	err = service.VehicleTypeRepository.Update(vehicleTypeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *VehicleTypeServiceImpl) Delete(custId string, vehicleTypeId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsExists(vehicleTypeId, custId, "vehicle_type_id")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("vehicle_type_id is still being used")
	// }

	err = service.VehicleTypeRepository.Delete(custId, vehicleTypeId, userId)
	if err != nil {
		return err
	}

	return err
}
