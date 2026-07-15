package service

import (
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/structs"
	"mobile/repository"
	"time"
)

type EmployeeService interface {
	List(entity.EmployeeQueryFilter) (data []entity.EmployeeResponse, total int64, lastPage int, err error)
	LookupList(entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error)
	Detail(entity.DetailEmployeeParams) (entity.EmployeeResponse, error)
	Store(entity.CreateEmployeeBody) (entity.EmployeeResponse, error)
	Update(int, entity.UpdateEmployeeRequest) error
	Delete(string, int, int64) error
	StoreMultiple(entity.CreateMultipleEmployeeBody, string, int64) (entity.EmployeeResponse, error)
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository) *employeeServiceImpl {
	return &employeeServiceImpl{
		EmployeeRepository: employeeRepository,
	}
}

type employeeServiceImpl struct {
	EmployeeRepository repository.EmployeeRepository
}

func (service *employeeServiceImpl) List(dataFilter entity.EmployeeQueryFilter) (data []entity.EmployeeResponse, total int64, lastPage int, err error) {
	var employees []model.Employee

	employees, total, lastPage, err = service.EmployeeRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range employees {
		var vResp entity.EmployeeResponse
		structs.Automapper(row, &vResp)
		if row.WorkDate != nil {
			workDate := row.WorkDate.Format("2006-01-02")
			vResp.WorkDate = &workDate
		}
		if row.Dob != nil {
			dob := row.Dob.Format("2006-01-02")
			vResp.Dob = &dob
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *employeeServiceImpl) LookupList(dataFilter entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error) {
	var employees []model.Employee

	employees, total, lastPage, err = service.EmployeeRepository.FindAllByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range employees {
		var vResp entity.EmployeeLookupResponse
		structs.Automapper(row, &vResp)
		if row.WorkDate != nil {
			workDate := row.WorkDate.Format("2006-01-02")
			vResp.WorkDate = &workDate
		}
		if row.Dob != nil {
			dob := row.Dob.Format("2006-01-02")
			vResp.Dob = &dob
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *employeeServiceImpl) Detail(params entity.DetailEmployeeParams) (response entity.EmployeeResponse, err error) {
	employee, err := service.EmployeeRepository.FindOneByEmployeeIdAndCustId(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(employee, &response)
	if err != nil {
		return response, err
	}

	if employee.WorkDate != nil {
		workDate := employee.WorkDate.Format("2006-01-02")
		response.WorkDate = &workDate
	}
	if employee.Dob != nil {
		dob := employee.Dob.Format("2006-01-02")
		response.Dob = &dob
	}

	return response, err
}

func (service *employeeServiceImpl) Store(request entity.CreateEmployeeBody) (response entity.EmployeeResponse, err error) {
	detailEmployee := entity.DetailEmployeeParams{
		CustId:       request.CustId,
		ParentCustId: request.ParentCustId,
		EmployeeCode: request.EmployeeCode,
	}
	employee, err := service.EmployeeRepository.FindOneByEmployeeCodeAndCustId(detailEmployee)
	if err == nil {
		return response, errors.New("emp_code: " + employee.EmployeeCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	var employeeData model.Employee
	structs.Automapper(request, &employeeData)
	workDate, err := time.Parse("2006-01-02", request.WorkDate)
	if err != nil {
		return response, err
	}
	dob, err := time.Parse("2006-01-02", request.Dob)
	if err != nil {
		return response, err
	}
	employeeData.WorkDate = &workDate
	employeeData.Dob = &dob
	employeeData.CreatedAt = &timeNow
	employeeData.CreatedBy = &request.CreatedBy
	employeeData.UpdatedAt = &timeNow
	employeeData.UpdatedBy = &request.CreatedBy
	employeeData.ImageUrl = &request.ImageUrl
	employeeData.ProvinceId = &request.ProvinceId
	employeeData.CityId = &request.CityId
	employeeData.SubDistrictId = &request.SubDistrictId
	employeeData.WardId = &request.WardId
	employeeData.IdentityNo = &request.IdentityNo
	employeeData.IsWaNo = request.IsWaNo
	employeeData.PostCode = &request.PostCode
	employeeData.DivisionId = &request.DivisionId

	employeeId, err := service.EmployeeRepository.Store(employeeData)
	if err != nil {
		return response, err
	}

	response.EmployeeId = employeeId

	return response, err
}

func (service *employeeServiceImpl) Update(employeeId int, request entity.UpdateEmployeeRequest) (err error) {

	detailEmployee := entity.DetailEmployeeParams{
		CustId:       request.CustId,
		ParentCustId: request.ParentCustId,
		EmployeeCode: request.EmployeeCode,
	}

	employee, err := service.EmployeeRepository.FindOneByEmployeeCodeAndCustId(detailEmployee)
	if err == nil && employee.EmployeeId != employeeId {
		return errors.New("emp_code: " + employee.EmployeeCode + " is already exists")
	}

	err = service.EmployeeRepository.Update(employeeId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *employeeServiceImpl) Delete(custId string, employeeId int, userId int64) (err error) {

	err = service.EmployeeRepository.Delete(custId, employeeId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *employeeServiceImpl) StoreMultiple(request entity.CreateMultipleEmployeeBody, tempCustId string, tempUserId int64) (response entity.EmployeeResponse, err error) {

	detailEmployee := entity.DetailEmployeeParams{
		CustId:       request.CustId,
		ParentCustId: request.ParentCustId,
	}

	for _, row := range request.Employees {
		employee, err := service.EmployeeRepository.FindOneByEmployeeCodeAndCustId(detailEmployee)
		if err == nil {
			return response, errors.New("emp_code: " + employee.EmployeeCode + " is already exists")
		}

		timeNow := time.Now().In(time.UTC)
		var employeeData model.Employee
		structs.Automapper(row, &employeeData)
		workDate, err := time.Parse("2006-01-02", row.WorkDate)
		if err != nil {
			return response, err
		}
		dob, err := time.Parse("2006-01-02", row.Dob)
		if err != nil {
			return response, err
		}
		employeeData.CustId = tempCustId
		employeeData.WorkDate = &workDate
		employeeData.Dob = &dob
		employeeData.CreatedAt = &timeNow
		employeeData.CreatedBy = &tempUserId
		employeeData.UpdatedAt = &timeNow
		employeeData.UpdatedBy = &tempUserId
		employeeData.ProvinceId = &row.ProvinceId
		// employeeData.Province = &row.Province
		employeeData.CityId = &row.CityId
		employeeData.SubDistrictId = &row.SubDistrictId
		// employeeData.SubDistrict = &row.SubDistrict
		employeeData.WardId = &row.WardId
		// employeeData.Ward = &row.Ward
		employeeData.IdentityNo = &row.IdentityNo
		employeeData.IsWaNo = row.IsWaNo
		employeeData.PostCode = &row.PostCode

		employeeId, err := service.EmployeeRepository.Store(employeeData)
		if err != nil {
			return response, err
		}

		response.EmployeeId = employeeId

	}

	return response, err
}
