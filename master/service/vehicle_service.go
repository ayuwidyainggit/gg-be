package service

import (
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"
)

type VehicleService interface {
	FindParentCustId(string) (entity.MCustomerResp, error)
	Detail(int64, string) (entity.VehicleResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.VehicleResponse, total int, lastPage int, err error)
	Store(entity.CreateVehicleBody) (entity.VehicleResponse, error)
	Update(int64, entity.UpdateVehicleRequest) error
	Delete(string, int64, int64) error
}

func NewVehicleService(vehicleRepository repository.VehicleRepository) *vehicleServiceImpl {
	return &vehicleServiceImpl{
		VehicleRepository: vehicleRepository,
	}
}

type vehicleServiceImpl struct {
	VehicleRepository repository.VehicleRepository
	// MProductRepository repository.MProductRepository
}

func (service *vehicleServiceImpl) FindParentCustId(custId string) (response entity.MCustomerResp, err error) {
	mCustomer, err := service.VehicleRepository.FindOneParentCustId(custId)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(mCustomer, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *vehicleServiceImpl) Detail(vehicleId int64, custId string) (response entity.VehicleResponse, err error) {
	vehicle, err := service.VehicleRepository.FindOneByVehicleIdAndCustId(vehicleId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(vehicle, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *vehicleServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.VehicleResponse, total int, lastPage int, err error) {
	vehicles, total, lastPage, err := service.VehicleRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range vehicles {
		var vResp entity.VehicleResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *vehicleServiceImpl) Store(request entity.CreateVehicleBody) (response entity.VehicleResponse, err error) {
	vehicleByDriver, err := service.VehicleRepository.FindOneByDriverAndCustId(request.DriverId, request.CustId)
	if err == nil {
		return response, errors.New(fmt.Sprintf("driver is already exists on vehicle id %v", vehicleByDriver.VehicleId))
	}

	vehicleByHelper, err := service.VehicleRepository.FindOneByHelperAndCustId(request.HelperId, request.CustId)
	if err == nil {
		return response, errors.New(fmt.Sprintf("helper is already exists on vehicle id %v", vehicleByHelper.VehicleId))
	}
	// vehicle_code & cust id validation, if err == nil, this means that code & cust id already exists
	vehicle, err := service.VehicleRepository.FindOneByVehicleNoAndCustId(request.VehicleNo, request.CustId)
	if err == nil {
		return response, errors.New("vehicle_no: " + vehicle.VehicleNo + " is already exists")
	}
	vehicleType, err := service.VehicleRepository.FindOneByVehicleTypeAndCustId(request.VehicleType)
	if err != nil {
		return response, errors.New("vehicle_type: " + vehicleType.VehicleTypeName + strconv.Itoa(request.VehicleType) + " not exists")
	}

	timeNow := time.Now().In(time.UTC)
	var vehicleData model.Vehicle
	structs.Automapper(request, &vehicleData)
	vehicleData.CreatedAt = &timeNow
	vehicleData.CreatedBy = &request.CreatedBy
	vehicleData.UpdatedAt = &timeNow

	vehicleId, err := service.VehicleRepository.Store(vehicleData)
	if err != nil {
		return response, err
	}

	response.VehicleId = vehicleId

	return response, err
}

func (service *vehicleServiceImpl) Update(vehicleId int64, request entity.UpdateVehicleRequest) (err error) {
	vehicleByDriver, err := service.VehicleRepository.FindOneByDriverAndCustId(request.DriverId, request.CustId)
	if err == nil && vehicleByDriver.VehicleId != vehicleId {
		return errors.New(fmt.Sprintf("driver is already exists on vehicle id %v", vehicleByDriver.VehicleId))
	}

	vehicleByHelper, err := service.VehicleRepository.FindOneByHelperAndCustId(request.HelperId, request.CustId)
	if err == nil && vehicleByHelper.VehicleId != vehicleId {
		return errors.New(fmt.Sprintf("helper is already exists on vehicle id %v", vehicleByHelper.VehicleId))
	}
	// vehicle_code & cust id validation, if err == nil and params vehicleId != vehicle.Id, this means that code & cust id already exists
	vehicle, err := service.VehicleRepository.FindOneByVehicleNoAndCustId(request.VehicleNo, request.CustId)
	if err == nil && vehicle.VehicleId != vehicleId {
		return errors.New("vehicle_no: " + vehicle.VehicleNo + " is already exists")
	}

	vehicleType, err := service.VehicleRepository.FindOneByVehicleTypeAndCustId(request.VehicleType)
	if err != nil {
		return errors.New("vehicle_type: " + vehicleType.VehicleTypeName + strconv.Itoa(request.VehicleType) + " not exists")
	}

	err = service.VehicleRepository.Update(vehicleId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *vehicleServiceImpl) Delete(custId string, vehicleId int64, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsKeyExists(vehicleId, custId, "vehicle_id1")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("vehicle_id is still being used")
	// }

	err = service.VehicleRepository.Delete(custId, vehicleId, userId)
	if err != nil {
		return err
	}

	return err
}
