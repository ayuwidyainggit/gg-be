package service

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

type EmployeeService interface {
	List(entity.EmployeeQueryFilter) (data []entity.EmployeeResponse, total int, lastPage int, err error)
	LookupList(entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupListWithoutSalesman(entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error)
	LookupAPI(entity.EmployeeLookupAPIFilter) (data []entity.EmployeeLookupMinimalItem, total int, lastPage int, err error)
	Detail(entity.DetailEmployeeParams) (entity.EmployeeResponse, error)
	Store(entity.CreateEmployeeBody) (entity.EmployeeResponse, error)
	Update(int, entity.UpdateEmployeeRequest) error
	Delete(string, int, int64) error
	StoreMultiple(entity.CreateMultipleEmployeeBody, string, int64) (entity.EmployeeResponse, error)
	Export(filter entity.EmployeeQueryFilter) (*bytes.Buffer, string, string, error)
	ExportTemplate(format string) (*bytes.Buffer, string, string, error)
	ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error)
	ImportEmployees(req entity.ImportRequest) error
	ImportEmployeesUpdate(req entity.ImportRequest) error
	ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error
	ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error
	ListPJP(dataFilter entity.EmployeePJPQueryFilter) (data []entity.EmployeePJPResponse, total int, lastPage int, err error)
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository) *employeeServiceImpl {
	return &employeeServiceImpl{
		EmployeeRepository: employeeRepository,
	}
}

type employeeServiceImpl struct {
	EmployeeRepository repository.EmployeeRepository
}

var (
	employeeExportHeaders = []string{
		"emp_name", "emp_code", "identity_no", "email", "phone_no", "wa_no",
		"province", "city", "sub_district", "ward", "post_code",
		"address", "dob", "last_education", "work_date", "emp_grp_name", "division_name",
		"image_url", "is_active",
	}

	employeeHeaderDisplay = map[string]string{
		"cust_id":         "Customer ID",
		"emp_id":          "Employee ID",
		"emp_code":        "ID Employee",
		"emp_name":        "Employee Name",
		"address":         "Address",
		"emp_type_id":     "Employee Type ID",
		"emp_type_name":   "Employee Type",
		"emp_grp_id":      "Employee Group ID",
		"emp_grp_code":    "Job Title Code",
		"emp_grp_name":    "Job Title",
		"work_date":       "Join Date",
		"last_education":  "Education",
		"dob":             "Date Of Birth",
		"phone_no":        "Phone",
		"wa_no":           "WhatsApp",
		"email":           "Email",
		"is_active":       "Status",
		"is_del":          "Is Deleted",
		"device_id":       "Device ID",
		"mac_address":     "MAC Address",
		"image_url":       "Photo",
		"identity_no":     "Identity No",
		"is_wa_no":        "Set as Whatsapp",
		"province_id":     "Province",
		"province":        "Province",
		"city_id":         "City",
		"city":            "City",
		"sub_district_id": "Sub District",
		"sub_district":    "Sub District",
		"ward_id":         "Village",
		"ward":            "Village",
		"post_code":       "Postal Code",
		"division_id":     "Division ID",
		"division_code":   "Division Code",
		"division_name":   "Division",
	}

	employeeTemplateUpdateHeaders = map[string][]string{
		"employee": {
			"emp_code", "emp_name", "identity_no", "email", "phone_no",
			"province", "city", "sub_district", "ward", "post_code",
			"address", "emp_grp_name", "division_name", "image_url", "dob",
			"last_education", "work_date", "wa_no", "is_active",
		},
		"emp_group":    {"emp_grp_id", "emp_grp_code", "emp_grp_name"},
		"emp_type":     {"emp_type_id", "emp_type_name"},
		"division":     {"division_id", "division_code", "division_name"},
		"province":     {"province_id", "province"},
		"city":         {"city_id", "city", "province_id"},
		"sub_district": {"sub_district_id", "sub_district", "city_id"},
		"ward":         {"ward_id", "ward", "sub_district_id"},
	}

	employeeTemplateDatasetNames = map[string]string{
		"employee":     "Employee",
		"emp_group":    "Employee Group",
		"emp_type":     "Employee Type",
		"division":     "Division",
		"province":     "Province",
		"city":         "City",
		"sub_district": "Sub District",
		"ward":         "Ward",
	}
)

const employeeDateLayout = "2006-01-02"

const (
	employeeTerritoryScopeAll      = "ALL"
	employeeTerritoryScopeSelected = "SELECTED"
)

func normalizeEmployeeTerritoryRequest(request *entity.CreateEmployeeBody) (model.EmployeeTerritoryMapping, error) {
	request.RegionIds = uniquePositiveInts(request.RegionIds)
	request.AreaIds = uniquePositiveInts(request.AreaIds)
	request.DistributorIds = uniquePositiveInts(request.DistributorIds)

	request.RegionScope = normalizeEmployeeTerritoryScope(request.RegionScope, len(request.RegionIds) > 0)
	request.AreaScope = normalizeEmployeeTerritoryScope(request.AreaScope, len(request.AreaIds) > 0)
	request.DistributorScope = normalizeEmployeeTerritoryScope(request.DistributorScope, len(request.DistributorIds) > 0)

	if err := validateEmployeeTerritoryScope("region_scope", request.RegionScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}
	if err := validateEmployeeTerritoryScope("area_scope", request.AreaScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}
	if err := validateEmployeeTerritoryScope("distributor_scope", request.DistributorScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}

	if request.RegionScope == employeeTerritoryScopeAll {
		request.RegionIds = nil
	}
	if request.AreaScope == employeeTerritoryScopeAll {
		request.AreaIds = nil
	}
	if request.DistributorScope == employeeTerritoryScopeAll {
		request.DistributorIds = nil
	}

	if request.RegionScope == employeeTerritoryScopeSelected && len(request.RegionIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("region_ids is required when region_scope is SELECTED")
	}
	if request.AreaScope == employeeTerritoryScopeSelected && len(request.AreaIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("area_ids is required when area_scope is SELECTED")
	}
	if request.DistributorScope == employeeTerritoryScopeSelected && len(request.DistributorIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("distributor_ids is required when distributor_scope is SELECTED")
	}
	if len(request.AreaIds) > 0 && len(request.RegionIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("region_ids is required when area_ids is provided")
	}
	if len(request.DistributorIds) > 0 && len(request.AreaIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("area_ids is required when distributor_ids is provided")
	}

	return model.EmployeeTerritoryMapping{
		RegionIds:      request.RegionIds,
		AreaIds:        request.AreaIds,
		DistributorIds: request.DistributorIds,
	}, nil
}

func normalizeEmployeeTerritoryUpdateRequest(request *entity.UpdateEmployeeRequest) (model.EmployeeTerritoryMapping, error) {
	request.RegionIds = uniquePositiveInts(request.RegionIds)
	request.AreaIds = uniquePositiveInts(request.AreaIds)
	request.DistributorIds = uniquePositiveInts(request.DistributorIds)

	request.RegionScope = normalizeEmployeeTerritoryScope(request.RegionScope, len(request.RegionIds) > 0)
	request.AreaScope = normalizeEmployeeTerritoryScope(request.AreaScope, len(request.AreaIds) > 0)
	request.DistributorScope = normalizeEmployeeTerritoryScope(request.DistributorScope, len(request.DistributorIds) > 0)

	if err := validateEmployeeTerritoryScope("region_scope", request.RegionScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}
	if err := validateEmployeeTerritoryScope("area_scope", request.AreaScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}
	if err := validateEmployeeTerritoryScope("distributor_scope", request.DistributorScope); err != nil {
		return model.EmployeeTerritoryMapping{}, err
	}

	if request.RegionScope == employeeTerritoryScopeAll {
		request.RegionIds = nil
	}
	if request.AreaScope == employeeTerritoryScopeAll {
		request.AreaIds = nil
	}
	if request.DistributorScope == employeeTerritoryScopeAll {
		request.DistributorIds = nil
	}

	if request.RegionScope == employeeTerritoryScopeSelected && len(request.RegionIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("region_ids is required when region_scope is SELECTED")
	}
	if request.AreaScope == employeeTerritoryScopeSelected && len(request.AreaIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("area_ids is required when area_scope is SELECTED")
	}
	if request.DistributorScope == employeeTerritoryScopeSelected && len(request.DistributorIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("distributor_ids is required when distributor_scope is SELECTED")
	}
	if len(request.AreaIds) > 0 && len(request.RegionIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("region_ids is required when area_ids is provided")
	}
	if len(request.DistributorIds) > 0 && len(request.AreaIds) == 0 {
		return model.EmployeeTerritoryMapping{}, errors.New("area_ids is required when distributor_ids is provided")
	}

	return model.EmployeeTerritoryMapping{
		RegionIds:      request.RegionIds,
		AreaIds:        request.AreaIds,
		DistributorIds: request.DistributorIds,
	}, nil
}

func normalizeEmployeeTerritoryScope(scope string, hasIDs bool) string {
	scope = strings.ToUpper(strings.TrimSpace(scope))
	if scope != "" {
		return scope
	}
	if hasIDs {
		return employeeTerritoryScopeSelected
	}
	return employeeTerritoryScopeAll
}

func validateEmployeeTerritoryScope(field string, scope string) error {
	if scope == employeeTerritoryScopeAll || scope == employeeTerritoryScopeSelected {
		return nil
	}
	return fmt.Errorf("%s must be ALL or SELECTED", field)
}

func employeeTerritoryUpdateProvided(request entity.UpdateEmployeeRequest) bool {
	return strings.TrimSpace(request.RegionScope) != "" ||
		strings.TrimSpace(request.AreaScope) != "" ||
		strings.TrimSpace(request.DistributorScope) != "" ||
		request.RegionIds != nil ||
		request.AreaIds != nil ||
		request.DistributorIds != nil
}

func uniquePositiveInts(values []int) []int {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[int]struct{}, len(values))
	unique := make([]int, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}
	return unique
}

func (service *employeeServiceImpl) List(dataFilter entity.EmployeeQueryFilter) (data []entity.EmployeeResponse, total int, lastPage int, err error) {
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

func (service *employeeServiceImpl) LookupList(dataFilter entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error) {
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

func (service *employeeServiceImpl) LookupListWithoutSalesman(dataFilter entity.EmployeeQueryFilter) (data []entity.EmployeeLookupResponse, total int, lastPage int, err error) {
	var employees []model.Employee

	employees, total, lastPage, err = service.EmployeeRepository.FindAllByCustIdLookupModeWithoutSalesman(dataFilter)
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

		// if row.EmpIdSalesman == 0 && row.SalesName == "" || row.DeleteSales {
		// 	data = append(data, vResp)
		// }

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *employeeServiceImpl) LookupAPI(dataFilter entity.EmployeeLookupAPIFilter) (data []entity.EmployeeLookupMinimalItem, total int, lastPage int, err error) {
	rows, total, lastPage, err := service.EmployeeRepository.FindLookupMinimal(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, row := range rows {
		data = append(data, entity.EmployeeLookupMinimalItem{
			EmpId:   row.EmployeeId,
			EmpCode: row.EmployeeCode,
			EmpName: row.EmployeeName,
		})
	}
	return data, total, lastPage, nil
}

func (service *employeeServiceImpl) Detail(params entity.DetailEmployeeParams) (response entity.EmployeeResponse, err error) {
	employee, err := service.EmployeeRepository.FindOneByEmployeeIdAndCustId(params)
	if err != nil {
		return response, err
	}

	territory, err := service.EmployeeRepository.FindEmployeeTerritoryDetail(params)
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

	response.Regions = employeeRegionMappingResponses(territory.Regions)
	response.Areas = employeeAreaMappingResponses(territory.Areas)
	response.Distributors = employeeDistributorMappingResponses(territory.Distributors)

	return response, err
}

func employeeRegionMappingResponses(rows []model.EmployeeRegionMappingDetail) []entity.RegionIdCodeNameResp {
	regions := make([]entity.RegionIdCodeNameResp, 0, len(rows))
	for _, row := range rows {
		regions = append(regions, entity.RegionIdCodeNameResp{
			RegionId:   row.RegionId,
			RegionCode: row.RegionCode,
			RegionName: row.RegionName,
		})
	}
	return regions
}

func employeeAreaMappingResponses(rows []model.EmployeeAreaMappingDetail) []entity.EmployeeAreaMappingResponse {
	areas := make([]entity.EmployeeAreaMappingResponse, 0, len(rows))
	for _, row := range rows {
		areas = append(areas, entity.EmployeeAreaMappingResponse{
			AreaId:     row.AreaId,
			AreaCode:   row.AreaCode,
			AreaName:   row.AreaName,
			RegionId:   row.RegionId,
			RegionCode: row.RegionCode,
			RegionName: row.RegionName,
		})
	}
	return areas
}

func employeeDistributorMappingResponses(rows []model.EmployeeDistributorMappingDetail) []entity.EmployeeDistributorMappingResponse {
	distributors := make([]entity.EmployeeDistributorMappingResponse, 0, len(rows))
	for _, row := range rows {
		distributors = append(distributors, entity.EmployeeDistributorMappingResponse{
			DistributorId:   row.DistributorId,
			DistributorCode: row.DistributorCode,
			DistributorName: row.DistributorName,
			AreaId:          row.AreaId,
			AreaCode:        row.AreaCode,
			AreaName:        row.AreaName,
			RegionId:        row.RegionId,
			RegionCode:      row.RegionCode,
			RegionName:      row.RegionName,
		})
	}
	return distributors
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

	territory, err := normalizeEmployeeTerritoryRequest(&request)
	if err != nil {
		return response, err
	}
	if err = service.EmployeeRepository.ValidateEmployeeTerritoryMapping(request.ParentCustId, territory); err != nil {
		return response, err
	}

	timeNow := time.Now().In(time.UTC)
	var employeeData model.Employee
	structs.Automapper(request, &employeeData)
	workDate, err := parseDate(request.WorkDate)
	if err != nil {
		return response, err
	}
	dob, err := parseDate(request.Dob)
	if err != nil {
		return response, err
	}
	employeeData.WorkDate = workDate
	employeeData.Dob = dob
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
	employeeData.RegionScope = request.RegionScope
	employeeData.AreaScope = request.AreaScope
	employeeData.DistributorScope = request.DistributorScope
	employeeData.RegionIds = territory.RegionIds
	employeeData.AreaIds = territory.AreaIds
	employeeData.DistributorIds = territory.DistributorIds

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

	if employeeTerritoryUpdateProvided(request) {
		request.TerritoryMappingProvided = true
		territory, err := normalizeEmployeeTerritoryUpdateRequest(&request)
		if err != nil {
			return err
		}
		if err = service.EmployeeRepository.ValidateEmployeeTerritoryMapping(request.ParentCustId, territory); err != nil {
			return err
		}
	}

	err = service.EmployeeRepository.Update(employeeId, request)
	if err != nil {
		return err
	}

	if request.IsActive != nil && *request.IsActive == false {
		err = service.EmployeeRepository.UpdateSalesman(employeeId, *request.IsActive)
		if err != nil {
			return err
		}
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
		workDate, err := parseDate(row.WorkDate)
		if err != nil {
			return response, err
		}
		dob, err := parseDate(row.Dob)
		if err != nil {
			return response, err
		}
		employeeData.CustId = tempCustId
		employeeData.WorkDate = workDate
		employeeData.Dob = dob
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

func (service *employeeServiceImpl) Export(filter entity.EmployeeQueryFilter) (*bytes.Buffer, string, string, error) {
	if filter.IsActive == nil && strings.TrimSpace(filter.Status) != "" {
		status := strings.ToLower(strings.TrimSpace(filter.Status))
		switch status {
		case "active":
			val := 1
			filter.IsActive = &val
		case "deactive", "inactive":
			val := 2
			filter.IsActive = &val
		case "all":
			// keep IsActive nil to fetch all records
		default:
			return nil, "", "", fmt.Errorf("status must be Active, Deactive, or All")
		}
	}

	employees, _, err := service.EmployeeRepository.FindAllExport(filter)
	if err != nil {
		return nil, "", "", err
	}

	format := strings.ToLower(strings.TrimSpace(filter.Format))
	switch format {
	case "csv":
		buf, err := createEmployeeCSV(employees)
		return buf, "text/csv", "employees.csv", err
	case "xls":
		buf, err := createEmployeeXLS(employees)
		return buf, "application/vnd.ms-excel", "employees.xls", err
	default:
		buf, err := createEmployeeXLSX(employees)
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "employees.xlsx", err
	}
}

func reorderEmployeeSelectedKeys(selectedKeys []string) []string {
	// Buat map untuk akses cepat
	selectedMap := make(map[string]bool, len(selectedKeys))
	for _, key := range selectedKeys {
		selectedMap[key] = true
	}

	// Susun ulang sesuai urutan utama
	ordered := []string{}
	for _, key := range employeeExportHeaders {
		if selectedMap[key] {
			ordered = append(ordered, key)
		}
	}

	// Tambahkan sisa kolom yang tidak ada di urutan utama (misal kolom baru)
	for _, key := range selectedKeys {
		if !containsEmployee(ordered, key) {
			ordered = append(ordered, key)
		}
	}

	return ordered
}

// helper sederhana
func containsEmployee(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func (service *employeeServiceImpl) ExportTemplate(format string) (*bytes.Buffer, string, string, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	instructions, _ := service.EmployeeRepository.GetImportInstructions("employee")
	switch format {
	case "csv":
		templateBuf, err := createEmployeeTemplateCSV()
		if err != nil {
			return nil, "", "", err
		}
		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)
		cf, err := zw.Create("employee_template.csv")
		if err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if _, err := cf.Write(templateBuf.Bytes()); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if instructions, err := service.EmployeeRepository.GetImportInstructions("employee"); err == nil && len(instructions) > 0 {
			insFile, err := zw.Create("instructions.csv")
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(insFile)
			cw.Comma = ';'
			_ = cw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
			for _, ins := range instructions {
				mandatory := "Kolom Tidak Wajib Diisi"
				if ins.Mandatory {
					mandatory = "Kolom Wajib Diisi"
				}
				stepVal := ""
				if ins.Step != nil {
					stepVal = *ins.Step
				}
				stepLabel, colorLabel := instructionStepAndColor(ins.Kolom, stepVal)
				_ = cw.Write([]string{instructionColumnDisplay(ins.Kolom), mandatory, ins.Keterangan, stepLabel, colorLabel})
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
		}
		if err := zw.Close(); err != nil {
			return nil, "", "", err
		}
		return zipBuf, "application/zip", "employee_template.zip", nil
	case "xls":
		buf, err := createEmployeeTemplateXLS(instructions)
		return buf, "application/vnd.ms-excel", "employee_template.xls", err
	default:
		buf, err := createEmployeeTemplateXLSX(instructions)
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "employee_template.xlsx", err
	}
}

func (service *employeeServiceImpl) ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error) {
	if len(fields) == 0 {
		return nil, "", "", errors.New("fields parameter is required")
	}

	sourceHeaders := employeeTemplateUpdateHeaders["employee"]
	if len(sourceHeaders) == 0 {
		return nil, "", "", errors.New("employee headers not configured")
	}

	aliasToKey := map[string]string{}
	for _, header := range sourceHeaders {
		aliasToKey[canonicalHeader(header)] = header
		display := employeeDisplayHeader(header)
		aliasToKey[canonicalHeader(display)] = header
	}
	aliasToKey["city_regency"] = "city"
	aliasToKey["job_title"] = "emp_grp_name"
	aliasToKey["division"] = "division_name"
	aliasToKey["province_name"] = "province"
	aliasToKey["city_name"] = "city"
	aliasToKey["sub_district_name"] = "sub_district"
	aliasToKey["village_name"] = "ward"
	aliasToKey["emp_id"] = "emp_code"
	aliasToKey["employee_id"] = "emp_code"
	aliasToKey["id_employee"] = "emp_code"

	selectedKeys := make([]string, 0, len(fields))
	seen := map[string]struct{}{}
	full := false
	for _, raw := range fields {
		canon := canonicalHeader(raw)
		if canon == "" {
			continue
		}
		if canon == "employee" {
			full = true
			break
		}
		if key, ok := aliasToKey[canon]; ok {
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			selectedKeys = append(selectedKeys, key)
		}
	}

	if full || len(selectedKeys) == 0 {
		selectedKeys = append([]string{}, sourceHeaders...)
	}

	selectedKeys = ensureLeadingField(selectedKeys, "emp_code")
	selectedKeys = reorderEmployeeSelectedKeys(selectedKeys)

	rows := [][]string{}

	instructions, _ := service.EmployeeRepository.GetImportInstructions("employee")

	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "csv":
		buf, err := createEmployeeTemplateUpdateCSV(selectedKeys, sourceHeaders, rows, instructions)
		return buf, "application/zip", "employee_template_update.zip", err
	case "xls":
		buf, err := createEmployeeTemplateUpdateXLS(selectedKeys, sourceHeaders, rows, instructions)
		return buf, "application/vnd.ms-excel", "employee_template_update.xls", err
	default:
		buf, err := createEmployeeTemplateUpdateXLSX(selectedKeys, sourceHeaders, rows, instructions)
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "employee_template_update.xlsx", err
	}
}

func (service *employeeServiceImpl) ImportEmployees(req entity.ImportRequest) error {
	rows, err := readImportRows(req)
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return errors.New("File harus memiliki header dan minimal satu baris data.")
	}

	header := rows[0]
	dataRows := rows[1:]
	dataRows = filterEmptyRows(dataRows)
	total := len(dataRows)
	if total == 0 {
		return errors.New("Tidak ditemukan data pada file (semua baris kosong).")
	}

	headerMap := buildHeaderIndex(header)

	historyId, err := service.EmployeeRepository.CreateImportHistory("employee", req.Filename, req.CustId, req.UserId, total)
	if err != nil {
		return errors.New("Gagal membuat riwayat impor data karyawan.")
	}

	// === Validasi kolom mandatory dulu (sebelum bikin history)
	var (
		mu     sync.Mutex
		errs   []string
		maxErr = 50
		wg     sync.WaitGroup
	)

	for i, row := range dataRows {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()

			rowNum := i + 2 // baris ke berapa di file (header = 1)
			rowMap := buildRowMap(r, headerMap)
			localErrs := []string{}
			prefix := employeeValidationPrefix(rowMap, rowNum)

			// ✅ Cek mandatory field
			for field := range mandatoryEmployeeColumns {
				if val, ok := rowMap[field]; !ok || strings.TrimSpace(val) == "" {
					localErrs = append(localErrs, fmt.Sprintf("%s: kolom '%s' wajib diisi.", prefix, employeeDisplayHeader(field)))
				}
			}

			// ✅ Cek panjang karakter
			for field, maxLen := range headerEmployeeMaxLength {
				if val, ok := rowMap[field]; ok && len(val) > maxLen {
					localErrs = append(localErrs, fmt.Sprintf("%s: kolom '%s' melebihi batas maksimal %d karakter (panjang saat ini %d).",
						prefix, employeeDisplayHeader(field), maxLen, len(val)))
				}
			}

			// Gabungkan ke hasil global
			if len(localErrs) > 0 {
				errMsg := strings.Join(localErrs, "; ")
				temp := buildEmployeeTemp(historyId, req.CustId, rowMap, errors.New(errMsg), rowNum)
				if e := service.EmployeeRepository.InsertEmployeeTemp(temp); e != nil {
					log.Printf("Gagal menyimpan data sementara karyawan pada baris %d: %v", rowNum, e)
				}

				mu.Lock()
				if len(errs) < maxErr {
					errs = append(errs, fmt.Sprintf("Baris %d: %s", rowNum, errMsg)) // ✅ 1 baris = 1 error
				}
				mu.Unlock()
			}
		}(i, row)
	}

	wg.Wait()

	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		// Update history (status masih processing, tapi update progress awal)
		_ = service.EmployeeRepository.UpdateImportHistory(historyId, successCount, failedCount, false)
		return fmt.Errorf("Validasi data gagal (%d baris):\n%s", failedCount, strings.Join(errs, "\n"))
	}

	go func() {
		success, failed := service.processEmployeeInsertRows(header, dataRows, req, historyId)
		if err := service.EmployeeRepository.UpdateImportHistory(historyId, success, failed, failed > 0); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (service *employeeServiceImpl) ImportEmployeesUpdate(req entity.ImportRequest) error {
	rows, err := readImportRows(req)
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return errors.New("File harus memiliki header dan minimal satu baris data")
	}

	header := rows[0]
	dataRows := filterEmptyRows(rows[1:])
	total := len(dataRows)
	if total == 0 {
		return errors.New("Tidak ada data yang ditemukan di dalam file")
	}

	headerMap := buildHeaderIndex(header)

	historyId, err := service.EmployeeRepository.CreateImportHistory("employee-update", req.Filename, req.CustId, req.UserId, total)
	if err != nil {
		return err
	}

	// === 1️⃣ Validasi panjang karakter dulu
	var (
		mu     sync.Mutex
		errs   []string
		maxErr = 50
		wg     sync.WaitGroup
	)

	for i, row := range dataRows {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()

			rowNum := i + 2
			rowMap := buildRowMap(r, headerMap)
			localErrs := []string{}
			prefix := employeeValidationPrefix(rowMap, rowNum)

			for field, maxLen := range headerEmployeeMaxLength {
				if val, ok := rowMap[field]; ok && len(val) > maxLen {
					localErrs = append(localErrs,
						fmt.Sprintf("%s: kolom '%s' melebihi batas maksimal %d karakter (saat ini %d karakter)",
							prefix, employeeDisplayHeader(field), maxLen, len(val)))
				}
			}

			if len(localErrs) > 0 {
				errMsg := strings.Join(localErrs, "; ")
				temp := buildEmployeeUpdateTemp(historyId, req.CustId, rowMap, errors.New(errMsg), rowNum)
				if e := service.EmployeeRepository.InsertEmployeeUpdateTemp(temp); e != nil {
					log.Printf("Gagal menyimpan data sementara karyawan: %v", e)
				}

				mu.Lock()
				if len(errs) < maxErr {
					errs = append(errs, fmt.Sprintf("Baris %d: %s", rowNum, errMsg)) // ✅ 1 baris = 1 error
				}
				mu.Unlock()
			}
		}(i, row)
	}

	wg.Wait()

	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		// Update history (status masih processing, tapi update progress awal)
		_ = service.EmployeeRepository.UpdateImportHistory(historyId, successCount, failedCount, false)
		return fmt.Errorf("Validasi data gagal:\n%s", strings.Join(errs, "\n"))
	}

	go func() {
		success, failed := service.processEmployeeUpdateRows(header, dataRows, req, historyId)
		if err := service.EmployeeRepository.UpdateImportHistory(historyId, success, failed, failed > 0); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (service *employeeServiceImpl) ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error {
	req.CustId = custId
	if strings.TrimSpace(req.Format) == "" {
		req.Format = "xlsx"
	}
	return service.ImportEmployees(req)
}

func (service *employeeServiceImpl) ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error {
	req.CustId = custId
	if strings.TrimSpace(req.Format) == "" {
		req.Format = "xlsx"
	}
	return service.ImportEmployeesUpdate(req)
}

func (service *employeeServiceImpl) ListPJP(dataFilter entity.EmployeePJPQueryFilter) (data []entity.EmployeePJPResponse, total int, lastPage int, err error) {
	employees, total, lastPage, err := service.EmployeeRepository.FindAllForPJP(dataFilter)
	if err != nil {
		return nil, 0, 0, err
	}

	for _, emp := range employees {
		data = append(data, entity.EmployeePJPResponse{
			EmpId:   emp.EmpId,
			EmpCode: emp.EmpCode,
			EmpName: emp.EmpName,
		})
	}

	return data, total, lastPage, nil
}

var headerEmployeeMaxLength = map[string]int{
	"emp_name":       100,
	"emp_code":       10,
	"identity_no":    16,
	"email":          100,
	"phone_no":       15,
	"wa_no":          15,
	"post_code":      10,
	"address":        150,
	"last_education": 100,
}

var allowedEmployeeEducations = map[string]string{
	"sd":            "SD",
	"smp":           "SMP",
	"sma":           "SMA",
	"sma sederajat": "SMA Sederajat",
	"d3":            "D3",
	"s1":            "S1",
	"s2":            "S2",
	"s3":            "S3",
}

var mandatoryEmployeeColumns = map[string]bool{
	"emp_name":      true,
	"emp_code":      true,
	"identity_no":   true,
	"email":         true,
	"phone_no":      true,
	"city":          true,
	"post_code":     true,
	"address":       true,
	"province":      true,
	"sub_district":  true,
	"ward":          true,
	"emp_grp_name":  true,
	"division_name": true,
}

func createEmployeeCSV(employees []model.EmployeeExport) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = ';'
	headers := make([]string, len(employeeExportHeaders))
	log.Println(headers)
	for i, key := range employeeExportHeaders {
		headers[i] = employeeDisplayHeader(key)
	}
	if err := w.Write(headers); err != nil {
		return nil, err
	}
	for _, emp := range employees {
		if err := w.Write(employeeExportValues(emp)); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf, nil
}

func createEmployeeXLSX(employees []model.EmployeeExport) (*bytes.Buffer, error) {
	return createEmployeeSpreadsheet(employees)
}

func createEmployeeXLS(employees []model.EmployeeExport) (*bytes.Buffer, error) {
	return createEmployeeSpreadsheet(employees)
}

func createEmployeeSpreadsheet(employees []model.EmployeeExport) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Employees"
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for i, key := range employeeExportHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, employeeDisplayHeader(key))
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	for r, emp := range employees {
		values := employeeExportValues(emp)
		for c, val := range values {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	autoFitEmployeeColumns(f, sheetName, employeeExportHeaders)

	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createEmployeeTemplateCSV() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = ';'
	headers := make([]string, len(employeeExportHeaders))
	for i, key := range employeeExportHeaders {
		display := employeeDisplayHeader(key)

		// Tambahkan * untuk kolom mandatory
		if mandatoryEmployeeColumns[key] {
			display += "*"
		}

		// Tambahkan (maksimal xx karakter)
		if max, ok := headerEmployeeMaxLength[key]; ok {
			display += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		headers[i] = display
	}
	if err := w.Write(headers); err != nil {
		return nil, err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf, nil
}

func createEmployeeTemplateXLSX(instructions []entity.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Employee Template"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"C6EFCE"}, Pattern: 1},
	})
	for i, key := range employeeExportHeaders {
		display := employeeDisplayHeader(key)

		// Tambahkan * untuk mandatory
		if mandatoryEmployeeColumns[key] {
			display += "*"
		}

		// Tambahkan (maksimal xx karakter)
		if max, ok := headerEmployeeMaxLength[key]; ok {
			display += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, display)
		f.SetCellStyle(sheet, cell, cell, style)
	}

	autoFitEmployeeColumns(f, sheet, employeeExportHeaders)

	if len(instructions) > 0 {
		insSheet := "Instructions"
		f.NewSheet(insSheet)
		headers := []string{"kolom", "mandatory", "keterangan", "step", "color"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(insSheet, cell, h)
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{Bold: true},
				Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
			})
			f.SetCellStyle(insSheet, cell, cell, style)
		}
		for r, ins := range instructions {
			rowIdx := r + 2
			mand := "Kolom Tidak Wajib Diisi"
			if ins.Mandatory {
				mand = "Kolom Wajib Diisi"
			}
			stepVal := ""
			if ins.Step != nil {
				stepVal = *ins.Step
			}
			stepLabel, colorLabel := instructionStepAndColor(ins.Kolom, stepVal)
			f.SetCellValue(insSheet, fmt.Sprintf("A%d", rowIdx), instructionColumnDisplay(ins.Kolom))
			f.SetCellValue(insSheet, fmt.Sprintf("B%d", rowIdx), mand)
			f.SetCellValue(insSheet, fmt.Sprintf("C%d", rowIdx), ins.Keterangan)
			f.SetCellValue(insSheet, fmt.Sprintf("D%d", rowIdx), stepLabel)
			f.SetCellValue(insSheet, fmt.Sprintf("E%d", rowIdx), colorLabel)
		}
		autoFitEmployeeColumns(f, insSheet, headers)
	}

	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createEmployeeTemplateXLS(instructions []entity.ImportInstruction) (*bytes.Buffer, error) {
	return createEmployeeTemplateXLSX(instructions)
}

func createEmployeeTemplateUpdateXLSX(headers []string, sourceHeaders []string, rows [][]string, instructions []entity.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Employee Update"
	idx, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	sourceIndex := make(map[string]int, len(sourceHeaders))
	for i, header := range sourceHeaders {
		sourceIndex[header] = i
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"C6EFCE"}, Pattern: 1},
	})
	for c, header := range headers {
		display := employeeDisplayHeader(header)

		// Tambahkan info maksimal karakter jika ada
		if max, ok := headerEmployeeMaxLength[header]; ok {
			display += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue(sheetName, cell, display)
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	autoFitEmployeeColumns(f, sheetName, headers)

	for r, row := range rows {
		for c, header := range headers {
			idxSrc, ok := sourceIndex[header]
			if !ok || idxSrc >= len(row) {
				continue
			}
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheetName, cell, row[idxSrc])
		}
	}

	if len(instructions) > 0 {
		insSheet := "Instructions"
		f.NewSheet(insSheet)
		headersIns := []string{"kolom", "mandatory", "keterangan", "step", "color"}
		for i, h := range headersIns {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(insSheet, cell, h)
			style, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{Bold: true},
				Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
			})
			f.SetCellStyle(insSheet, cell, cell, style)
		}
		for r, ins := range instructions {
			rowIdx := r + 2
			mand := "Kolom Tidak Wajib Diisi"
			if ins.Mandatory {
				mand = "Kolom Wajib Diisi"
			}
			stepVal := ""
			if ins.Step != nil {
				stepVal = *ins.Step
			}
			stepLabel, colorLabel := instructionStepAndColor(ins.Kolom, stepVal)
			f.SetCellValue(insSheet, fmt.Sprintf("A%d", rowIdx), instructionColumnDisplay(ins.Kolom))
			f.SetCellValue(insSheet, fmt.Sprintf("B%d", rowIdx), mand)
			f.SetCellValue(insSheet, fmt.Sprintf("C%d", rowIdx), ins.Keterangan)
			f.SetCellValue(insSheet, fmt.Sprintf("D%d", rowIdx), stepLabel)
			f.SetCellValue(insSheet, fmt.Sprintf("E%d", rowIdx), colorLabel)
		}
		autoFitEmployeeColumns(f, insSheet, headersIns)
	}

	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createEmployeeTemplateUpdateXLS(headers []string, sourceHeaders []string, rows [][]string, instructions []entity.ImportInstruction) (*bytes.Buffer, error) {
	return createEmployeeTemplateUpdateXLSX(headers, sourceHeaders, rows, instructions)
}

func createEmployeeTemplateUpdateCSV(headers []string, sourceHeaders []string, rows [][]string, instructions []entity.ImportInstruction) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	sourceIndex := make(map[string]int, len(sourceHeaders))
	for i, header := range sourceHeaders {
		sourceIndex[header] = i
	}

	mainFile, err := zw.Create("employee_template_update.csv")
	if err != nil {
		_ = zw.Close()
		return nil, err
	}
	writer := csv.NewWriter(mainFile)
	writer.Comma = ';'
	display := make([]string, len(headers))
	for i, h := range headers {
		name := employeeDisplayHeader(h)

		// Tambahkan (maksimal xx karakter) tanpa *
		if max, ok := headerEmployeeMaxLength[h]; ok {
			name += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		display[i] = name
	}
	if err := writer.Write(display); err != nil {
		_ = zw.Close()
		return nil, err
	}
	for _, row := range rows {
		values := make([]string, len(headers))
		for i, header := range headers {
			if idxSrc, ok := sourceIndex[header]; ok && idxSrc < len(row) {
				values[i] = row[idxSrc]
			}
		}
		if err := writer.Write(values); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		_ = zw.Close()
		return nil, err
	}

	if len(instructions) > 0 {
		insFile, err := zw.Create("instructions.csv")
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		insWriter := csv.NewWriter(insFile)
		insWriter.Comma = ';'
		_ = insWriter.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
		for _, ins := range instructions {
			mand := "Kolom Tidak Wajib Diisi"
			if ins.Mandatory {
				mand = "Kolom Wajib Diisi"
			}
			stepVal := ""
			if ins.Step != nil {
				stepVal = *ins.Step
			}
			stepLabel, colorLabel := instructionStepAndColor(ins.Kolom, stepVal)
			_ = insWriter.Write([]string{instructionColumnDisplay(ins.Kolom), mand, ins.Keterangan, stepLabel, colorLabel})
		}
		insWriter.Flush()
		if err := insWriter.Error(); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func employeeExportValues(emp model.EmployeeExport) []string {
	values := make([]string, len(employeeExportHeaders))
	for i, key := range employeeExportHeaders {
		switch key {
		case "cust_id":
			values[i] = emp.CustId
		case "emp_id":
			values[i] = strconv.FormatInt(emp.EmployeeId, 10)
		case "emp_code":
			values[i] = stringPtrValue(emp.EmployeeCode)
		case "emp_name":
			values[i] = stringPtrValue(emp.EmployeeName)
		case "address":
			values[i] = stringPtrValue(emp.Address)
		case "emp_type_id":
			values[i] = stringPtrValue(emp.EmpTypeId)
		case "emp_grp_id":
			values[i] = intPtrValue(emp.EmpGrpId)
		case "emp_grp_code":
			values[i] = stringPtrValue(emp.EmpGrpCode)
		case "emp_grp_name":
			values[i] = stringPtrValue(emp.EmpGrpName)
		case "work_date":
			values[i] = timePtrValue(emp.WorkDate)
		case "last_education":
			values[i] = stringPtrValue(emp.LastEducation)
		case "dob":
			values[i] = timePtrValue(emp.Dob)
		case "phone_no":
			values[i] = stringPtrValue(emp.PhoneNo)
		case "wa_no":
			values[i] = stringPtrValue(emp.WaNo)
		case "email":
			values[i] = stringPtrValue(emp.Email)
		case "is_active":
			values[i] = boolPtrToActiveDeactive(emp.IsActive)
		case "is_del":
			values[i] = boolPtrValue(emp.IsDel)
		case "device_id":
			values[i] = stringPtrValue(emp.DeviceID)
		case "mac_address":
			values[i] = stringPtrValue(emp.MacAddress)
		case "image_url":
			values[i] = stringPtrValue(emp.ImageURL)
		case "identity_no":
			values[i] = stringPtrValue(emp.IdentityNo)
		case "is_wa_no":
			values[i] = boolPtrValue(emp.IsWaNo)
		case "province_id":
			values[i] = stringPtrValue(emp.ProvinceID)
		case "province":
			values[i] = stringPtrValue(emp.Province)
		case "city_id":
			values[i] = stringPtrValue(emp.CityID)
		case "city":
			values[i] = stringPtrValue(emp.City)
		case "sub_district_id":
			values[i] = stringPtrValue(emp.SubDistrictID)
		case "sub_district":
			values[i] = stringPtrValue(emp.SubDistrict)
		case "ward_id":
			values[i] = stringPtrValue(emp.WardID)
		case "ward":
			values[i] = stringPtrValue(emp.Ward)
		case "post_code":
			values[i] = stringPtrValue(emp.PostCode)
		case "division_id":
			values[i] = intPtrValue(emp.DivisionID)
		case "division_code":
			values[i] = stringPtrValue(emp.DivisionCode)
		case "division_name":
			values[i] = stringPtrValue(emp.DivisionName)
		default:
			values[i] = ""
		}
	}
	return values
}

func autoFitEmployeeColumns(f *excelize.File, sheet string, headers []string) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return
	}
	colCount := len(headers)
	for _, row := range rows {
		if len(row) > colCount {
			colCount = len(row)
		}
	}
	if colCount == 0 {
		return
	}
	for col := 0; col < colCount; col++ {
		maxLen := 0
		for _, row := range rows {
			if col >= len(row) {
				continue
			}
			length := utf8.RuneCountInString(strings.TrimSpace(row[col]))
			if length > maxLen {
				maxLen = length
			}
		}
		if maxLen < 8 {
			maxLen = 8
		} else {
			maxLen += 2
		}
		colName, err := excelize.ColumnNumberToName(col + 1)
		if err != nil {
			continue
		}
		_ = f.SetColWidth(sheet, colName, colName, float64(maxLen))
	}
}

func employeeDisplayHeader(key string) string {
	if label, ok := employeeHeaderDisplay[key]; ok {
		return label
	}
	key = strings.ReplaceAll(key, "_", " ")
	key = strings.Title(key)
	return key
}

func stringPtrValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func boolPtrValue(v *bool) string {
	if v == nil {
		return ""
	}
	return strconv.FormatBool(*v)
}

func boolPtrToActiveDeactive(v *bool) string {
	if v == nil {
		return ""
	}
	if *v {
		return "Active"
	}
	return "Deactive"
}

func timePtrValue(v *time.Time) string {
	if v == nil {
		return ""
	}
	return v.Format(employeeDateLayout)
}

func intPtrValue(v interface{}) string {
	switch val := v.(type) {
	case *int:
		if val == nil {
			return ""
		}
		return strconv.Itoa(*val)
	case *int64:
		if val == nil {
			return ""
		}
		return strconv.FormatInt(*val, 10)
	default:
		return ""
	}
}

func readImportRows(req entity.ImportRequest) ([][]string, error) {
	format := strings.ToLower(strings.TrimSpace(req.Format))
	switch format {
	case "csv":
		reader := csv.NewReader(req.File)
		reader.Comma = ';'
		reader.TrimLeadingSpace = true
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		return normalizeRows(records), nil
	case "xls", "xlsx":
		f, err := excelize.OpenReader(req.File)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close() }()
		sheetName := f.GetSheetName(0)
		if sheetName == "" {
			sheetName = "Sheet1"
		}
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}
		return normalizeRows(rows), nil
	default:
		return nil, fmt.Errorf("Gagal membaca file. Pastikan format file sesuai template yang benar. Format file terkirim: %s", req.Format)
	}
}

func normalizeRows(rows [][]string) [][]string {
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, normalizeRow(row))
	}
	return out
}

func normalizeRow(row []string) []string {
	out := make([]string, len(row))
	for i, col := range row {
		val := strings.TrimSpace(col)
		if i == 0 {
			val = strings.TrimPrefix(val, "\ufeff")
		}
		out[i] = val
	}
	return out
}

func filterEmptyRows(rows [][]string) [][]string {
	filtered := make([][]string, 0, len(rows))
	for _, row := range rows {
		empty := true
		for _, col := range row {
			if strings.TrimSpace(col) != "" {
				empty = false
				break
			}
		}
		if !empty {
			filtered = append(filtered, normalizeRow(row))
		}
	}
	return filtered
}

func (service *employeeServiceImpl) processEmployeeInsertRows(header []string, rows [][]string, req entity.ImportRequest, historyId int64) (int, int) {
	headerMap := buildHeaderIndex(header)
	success := 0
	failed := 0
	for idx, row := range rows {
		rowNumber := idx + 2
		rowMap := buildRowMap(row, headerMap)
		empCode := strings.TrimSpace(rowMap["emp_code"])

		exist, _ := service.EmployeeRepository.CheckEmployeeExists(req.CustId, empCode)
		if exist {
			errMsg := fmt.Sprintf("ID Karyawan %s sudah ada. Silakan gunakan ID lain.", empCode)
			temp := buildEmployeeTemp(historyId, req.CustId, rowMap, errors.New(errMsg), rowNumber)
			temp.StatusInsert = "failed-duplicate"
			if e := service.EmployeeRepository.InsertEmployeeTemp(temp); e != nil {
				log.Printf("Gagal menyimpan data sementara karyawan duplikat pada baris %d: %v", rowNumber, e)
			}
			failed++
			continue
		}

		if err := service.insertEmployeeRow(rowMap, req, rowNumber); err != nil {
			failed++
			temp := buildEmployeeTemp(historyId, req.CustId, rowMap, err, rowNumber)
			if e := service.EmployeeRepository.InsertEmployeeTemp(temp); e != nil {
				log.Printf("failed to record employee import error: %v", e)
			}
			continue
		}
		success++
	}
	return success, failed
}

func (service *employeeServiceImpl) processEmployeeUpdateRows(header []string, rows [][]string, req entity.ImportRequest, historyId int64) (int, int) {
	headerMap := buildHeaderIndex(header)
	success := 0
	failed := 0
	for idx, row := range rows {
		rowNumber := idx + 2
		rowMap := buildRowMap(row, headerMap)
		if err := service.updateEmployeeRow(rowMap, req, rowNumber); err != nil {
			failed++
			temp := buildEmployeeUpdateTemp(historyId, req.CustId, rowMap, err, rowNumber)
			if e := service.EmployeeRepository.InsertEmployeeUpdateTemp(temp); e != nil {
				log.Printf("failed to record employee update import error: %v", e)
			}
			continue
		}
		success++
	}
	return success, failed
}

func buildHeaderIndex(header []string) map[string]int {
	index := make(map[string]int, len(header))
	for i, h := range header {
		h = cleanHeaderEmployee(h)
		key := canonicalHeader(h)
		if key == "status" {
			key = "is_active"
		}
		if key == "" {
			continue
		}
		if _, exists := index[key]; !exists {
			index[key] = i
		}
	}
	return index
}

func cleanHeaderEmployee(h string) string {
	h = strings.TrimSpace(h)

	// Hapus bagian seperti "(maksimal 100 karakter)"
	re := regexp.MustCompile(`\(maksimal\s*\d+\s*karakter\)`)
	h = re.ReplaceAllString(h, "")

	// Hapus tanda * di depan / belakang
	h = strings.TrimPrefix(h, "*")
	h = strings.TrimSuffix(h, "*")

	return strings.TrimSpace(h)
}

func canonicalHeader(h string) string {
	key := strings.TrimSpace(strings.ToLower(h))
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ReplaceAll(key, "__", "_")
	return key
}

func buildRowMap(row []string, header map[string]int) map[string]string {
	rowMap := make(map[string]string, len(header))
	for key, idx := range header {
		if idx < len(row) {
			rowMap[key] = strings.TrimSpace(row[idx])
		} else {
			rowMap[key] = ""
		}
	}
	applyRowAliases(rowMap)
	return rowMap
}

func requireField(row map[string]string, label string, keys ...string) (string, error) {
	for _, key := range keys {
		if val, ok := row[key]; ok {
			trimmed := strings.TrimSpace(val)
			if trimmed != "" {
				row[key] = trimmed
				return trimmed, nil
			}
		}
	}
	return "", fmt.Errorf("%s is required", label)
}

func (service *employeeServiceImpl) lookupProvinceId(req entity.ImportRequest, name string) (string, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindProvinceIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}
	return "", fmt.Errorf("province '%s' tidak ditemukan", name)
}

func (service *employeeServiceImpl) lookupCityId(req entity.ImportRequest, name string) (string, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindCityIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}
	return "", fmt.Errorf("city '%s' tidak ditemukan", name)
}

func (service *employeeServiceImpl) lookupSubDistrictId(req entity.ImportRequest, name string) (string, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindSubDistrictIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}
	return "", fmt.Errorf("sub district '%s' tidak ditemukan", name)
}

func (service *employeeServiceImpl) lookupWardId(req entity.ImportRequest, name string) (string, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindWardIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}
	return "", fmt.Errorf("village '%s' tidak ditemukan", name)
}

func (service *employeeServiceImpl) lookupEmpGroupId(req entity.ImportRequest, name string) (int, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindEmpGroupIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}
	return 0, fmt.Errorf("job title '%s' tidak ditemukan", name)
}

func (service *employeeServiceImpl) lookupDivisionId(req entity.ImportRequest, name string) (int, error) {
	for _, cust := range uniqueCustIds(req.ParentCustId, req.CustId) {
		id, err := service.EmployeeRepository.FindDivisionIdByName(cust, name)
		if err == nil {
			return id, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}
	}
	return 0, fmt.Errorf("division '%s' tidak ditemukan", name)
}

func applyRowAliases(row map[string]string) {
	copyIfEmpty := func(from, to string) {
		if val, ok := row[from]; ok && strings.TrimSpace(val) != "" {
			if existing, ok2 := row[to]; !ok2 || strings.TrimSpace(existing) == "" {
				row[to] = val
			}
		}
	}

	copyIfEmpty("id_employee", "emp_code")
	copyIfEmpty("emp_code", "id_employee")
	copyIfEmpty("emp_id", "emp_code")
	copyIfEmpty("emp_code", "emp_id")
	copyIfEmpty("employee_id", "emp_code")
	copyIfEmpty("city_regency", "city")
	copyIfEmpty("job_title", "emp_grp_name")
	copyIfEmpty("division", "division_name")
	copyIfEmpty("village", "ward")
	copyIfEmpty("employee_name", "emp_name")
	copyIfEmpty("phone", "phone_no")
	copyIfEmpty("postal_code", "post_code")
	copyIfEmpty("photo", "image_url")
	copyIfEmpty("date_of_birth", "dob")
	copyIfEmpty("education", "last_education")
	copyIfEmpty("join_date", "work_date")
	copyIfEmpty("whatsapp", "wa_no")
	copyIfEmpty("status", "is_active")
}

func employeeValidationPrefix(row map[string]string, rowNumber int) string {
	var identifier string
	if row != nil {
		identifier = strings.TrimSpace(firstNonEmpty(
			row["emp_code"],
			row["id_employee"],
			row["emp_id"],
			row["employee_id"],
		))
	}
	if identifier != "" {
		return fmt.Sprintf("employee ID %s", identifier)
	}
	if rowNumber > 0 {
		return fmt.Sprintf("employee ID tidak diketahui (baris ke-%d)", rowNumber)
	}
	return "employee ID tidak diketahui"
}

func wrapEmployeeError(row map[string]string, rowNumber int, err error) error {
	if err == nil {
		return nil
	}
	msg := strings.TrimSpace(err.Error())
	prefix := employeeValidationPrefix(row, rowNumber)
	if prefix == "" {
		if msg == "" {
			return err
		}
		return errors.New(msg)
	}
	lowerPrefix := strings.ToLower(prefix)
	lowerMsg := strings.ToLower(msg)
	if lowerMsg == "" {
		return errors.New(prefix)
	}
	if strings.HasPrefix(lowerMsg, lowerPrefix) {
		return err
	}
	return fmt.Errorf("%s: %s", prefix, msg)
}

func isMasterNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	return strings.Contains(msg, "tidak ditemukan") || strings.Contains(msg, "not found")
}

func instructionColumnDisplay(column string) string {
	aliasMap := map[string]string{
		"emp_id":            "emp_code",
		"employee_id":       "emp_code",
		"id_employee":       "emp_code",
		"city_regency":      "city",
		"job_title":         "emp_grp_name",
		"division":          "division_name",
		"sub_district":      "sub_district",
		"village":           "ward",
		"province":          "province",
		"province_name":     "province",
		"city_name":         "city",
		"sub_district_name": "sub_district",
		"village_name":      "ward",
	}
	key := canonicalHeader(column)
	if mapped, ok := aliasMap[key]; ok {
		return employeeDisplayHeader(mapped)
	}
	if _, ok := employeeHeaderDisplay[column]; ok {
		return employeeDisplayHeader(column)
	}
	if _, ok := employeeHeaderDisplay[key]; ok {
		return employeeDisplayHeader(key)
	}
	return column
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func uniqueCustIds(ids ...string) []string {
	seen := map[string]struct{}{}
	ordered := make([]string, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		ordered = append(ordered, trimmed)
	}
	return ordered
}

func ensureLeadingField(list []string, field string) []string {
	field = strings.TrimSpace(field)
	if field == "" {
		return list
	}
	pos := -1
	for i, v := range list {
		if v == field {
			pos = i
			break
		}
	}
	if pos == 0 {
		return list
	}
	if pos > 0 {
		list = append(append([]string{}, list[:pos]...), list[pos+1:]...)
	}
	return append([]string{field}, list...)
}

func normalizeEducation(value string) string {
	val := strings.TrimSpace(strings.ToLower(value))
	if val == "" {
		return ""
	}
	val = strings.Join(strings.Fields(val), " ")
	if canonical, ok := allowedEmployeeEducations[val]; ok {
		return canonical
	}
	return ""
}

func (service *employeeServiceImpl) insertEmployeeRow(row map[string]string, req entity.ImportRequest, rowNumber int) error {
	empCode, err := requireField(row, "ID Employee", "emp_code", "id_employee", "employee_id")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["emp_code"] = empCode

	empName, err := requireField(row, "Employee Name", "emp_name")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["emp_name"] = empName

	identityNo, err := requireField(row, "Identity No", "identity_no")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["identity_no"] = identityNo

	email, err := requireField(row, "Email", "email")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["email"] = email

	phoneNo, err := requireField(row, "Phone", "phone_no")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["phone_no"] = phoneNo

	address, err := requireField(row, "Address", "address")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["address"] = address

	postCode, err := requireField(row, "Postal Code", "post_code")
	if err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	row["post_code"] = postCode

	missingMasters := make([]string, 0, 6)
	seenMissing := make(map[string]struct{})
	addMissing := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		msg := fmt.Sprintf("master %s '%s' tidak ditemukan", label, value)
		if _, exists := seenMissing[msg]; exists {
			return
		}
		seenMissing[msg] = struct{}{}
		missingMasters = append(missingMasters, msg)
	}

	provinceName := strings.TrimSpace(firstNonEmpty(row["province"], row["province_name"]))
	provinceId := strings.TrimSpace(row["province_id"])
	if provinceId == "" {
		if provinceName == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("Province is required"))
		}
		provinceId, err = service.lookupProvinceId(req, provinceName)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("province", provinceName)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		}
	}
	if provinceName == "" {
		provinceName = provinceId
	}
	row["province"] = provinceName
	row["province_id"] = provinceId

	cityName := strings.TrimSpace(firstNonEmpty(row["city"], row["city_regency"], row["city_name"]))
	cityId := strings.TrimSpace(row["city_id"])
	if cityId == "" {
		if cityName == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("City is required"))
		}
		cityId, err = service.lookupCityId(req, cityName)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("city", cityName)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		}
	}
	if cityName == "" {
		cityName = cityId
	}
	row["city"] = cityName
	row["city_id"] = cityId

	subDistrictName := strings.TrimSpace(firstNonEmpty(row["sub_district"], row["sub_district_name"]))
	subDistrictId := strings.TrimSpace(row["sub_district_id"])
	if subDistrictId == "" {
		if subDistrictName == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("Sub District is required"))
		}
		subDistrictId, err = service.lookupSubDistrictId(req, subDistrictName)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("sub district", subDistrictName)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		}
	}
	if subDistrictName == "" {
		subDistrictName = subDistrictId
	}
	row["sub_district"] = subDistrictName
	row["sub_district_id"] = subDistrictId

	wardName := strings.TrimSpace(firstNonEmpty(row["ward"], row["village"], row["village_name"]))
	wardId := strings.TrimSpace(row["ward_id"])
	if wardId == "" {
		if wardName == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("Village is required"))
		}
		wardId, err = service.lookupWardId(req, wardName)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("village", wardName)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		}
	}
	if wardName == "" {
		wardName = wardId
	}
	row["ward"] = wardName
	row["ward_id"] = wardId

	empGrpPtr, err := parseOptionalInt(row["emp_grp_id"])
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid emp_grp_id: %w", err))
	}
	jobTitle := strings.TrimSpace(firstNonEmpty(row["emp_grp_name"], row["job_title"]))
	if empGrpPtr == nil {
		if jobTitle == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("Job Title is required"))
		}
		id, errLookup := service.lookupEmpGroupId(req, jobTitle)
		if errLookup != nil {
			if isMasterNotFoundError(errLookup) {
				addMissing("job title", jobTitle)
			} else {
				return wrapEmployeeError(row, rowNumber, errLookup)
			}
		} else {
			idCopy := id
			empGrpPtr = &idCopy
		}
	}
	row["emp_grp_name"] = jobTitle
	if empGrpPtr != nil {
		row["emp_grp_id"] = strconv.Itoa(*empGrpPtr)
	}

	divisionPtr, err := parseOptionalInt(row["division_id"])
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid division_id: %w", err))
	}
	divisionName := strings.TrimSpace(firstNonEmpty(row["division_name"], row["division"]))
	if divisionPtr == nil {
		if divisionName == "" {
			return wrapEmployeeError(row, rowNumber, errors.New("Division is required"))
		}
		id, errLookup := service.lookupDivisionId(req, divisionName)
		if errLookup != nil {
			if isMasterNotFoundError(errLookup) {
				addMissing("division", divisionName)
			} else {
				return wrapEmployeeError(row, rowNumber, errLookup)
			}
		} else {
			idCopy := id
			divisionPtr = &idCopy
		}
	}
	row["division_name"] = divisionName
	if divisionPtr != nil {
		row["division_id"] = strconv.Itoa(*divisionPtr)
	}

	if len(missingMasters) > 0 {
		return wrapEmployeeError(row, rowNumber, errors.New(strings.Join(missingMasters, ", ")))
	}

	rawEducation := strings.TrimSpace(row["last_education"])
	normalizedEducation := normalizeEducation(rawEducation)
	if normalizedEducation == "" && rawEducation != "" {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("last education '%s' tidak termasuk dalam daftar yang diperbolehkan", rawEducation))
	}
	row["last_education"] = normalizedEducation

	timeNow := time.Now().In(time.UTC)
	createdBy := req.UserId

	provinceIDCopy := provinceId
	cityIDCopy := cityId
	subDistrictIDCopy := subDistrictId
	wardIDCopy := wardId

	emp := model.Employee{
		CustId:        req.CustId,
		EmployeeCode:  empCode,
		EmployeeName:  empName,
		Address:       stringPtr(address),
		EmpTypeId:     stringPtr(row["emp_type_id"]),
		EmpGrpId:      empGrpPtr,
		LastEducation: stringPtr(normalizedEducation),
		PhoneNo:       stringPtr(phoneNo),
		WaNo:          stringPtr(row["wa_no"]),
		Email:         stringPtr(email),
		ImageUrl:      stringPtr(row["image_url"]),
		IdentityNo:    stringPtr(identityNo),
		ProvinceId:    &provinceIDCopy,
		CityId:        &cityIDCopy,
		SubDistrictId: &subDistrictIDCopy,
		WardId:        &wardIDCopy,
		PostCode:      stringPtr(postCode),
		DivisionId:    divisionPtr,
		CreatedBy:     &createdBy,
		CreatedAt:     &timeNow,
		UpdatedBy:     &createdBy,
		UpdatedAt:     &timeNow,
		IsDel:         false,
	}

	workDate, err := parseDate(row["work_date"])
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid work_date: %w", err))
	}
	emp.WorkDate = workDate

	dob, err := parseDate(row["dob"])
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid dob: %w", err))
	}
	emp.Dob = dob

	isActive, err := parseBoolDefault(row["is_active"], true)
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid is_active: %w", err))
	}
	emp.IsActive = isActive

	isWaNo, err := parseBoolDefault(row["is_wa_no"], true)
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid is_wa_no: %w", err))
	}
	emp.IsWaNo = isWaNo

	if _, err := service.EmployeeRepository.Store(emp); err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	return nil
}

func (service *employeeServiceImpl) updateEmployeeRow(row map[string]string, req entity.ImportRequest, rowNumber int) error {
	empCode := strings.TrimSpace(row["emp_code"])
	if empCode == "" {
		empCode = strings.TrimSpace(row["id_employee"])
	}
	if empCode == "" {
		return wrapEmployeeError(row, rowNumber, errors.New("id employee (emp_code) is required"))
	}

	detail := entity.DetailEmployeeParams{
		CustId:       req.CustId,
		ParentCustId: req.ParentCustId,
		EmployeeCode: empCode,
	}

	current, err := service.EmployeeRepository.FindOneByEmployeeCodeAndCustId(detail)
	if err != nil {
		return wrapEmployeeError(row, rowNumber, fmt.Errorf("emp_code %s tidak ditemukan", empCode))
	}

	updateReq := entity.UpdateEmployeeRequest{
		CustId:       req.CustId,
		ParentCustId: req.ParentCustId,
		UpdatedBy:    req.UserId,
		EmployeeCode: empCode,
	}

	missingMasters := make([]string, 0, 4)
	seenMissing := make(map[string]struct{})
	addMissing := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		msg := fmt.Sprintf("master %s '%s' tidak ditemukan", label, value)
		if _, exists := seenMissing[msg]; exists {
			return
		}
		seenMissing[msg] = struct{}{}
		missingMasters = append(missingMasters, msg)
	}

	if v := strings.TrimSpace(row["emp_name"]); v != "" {
		updateReq.EmployeeName = v
	}
	if v := strings.TrimSpace(row["address"]); v != "" {
		updateReq.Address = v
	}
	if v := strings.TrimSpace(row["last_education"]); v != "" {
		norm := normalizeEducation(v)
		if norm == "" {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("last education '%s' tidak termasuk dalam daftar yang diperbolehkan", v))
		}
		updateReq.LastEducation = norm
		row["last_education"] = norm
	}
	if v := strings.TrimSpace(row["phone_no"]); v != "" {
		updateReq.PhoneNo = v
	}
	if v := strings.TrimSpace(row["wa_no"]); v != "" {
		updateReq.WaNo = v
	}
	if v := strings.TrimSpace(row["email"]); v != "" {
		updateReq.Email = v
	}
	if v := strings.TrimSpace(row["emp_type_id"]); v != "" {
		updateReq.EmpTypeId = v
	}
	if v := strings.TrimSpace(row["identity_no"]); v != "" {
		updateReq.IdentityNo = v
	}
	if v := strings.TrimSpace(row["post_code"]); v != "" {
		updateReq.PostCode = v
	}
	if v := strings.TrimSpace(row["image_url"]); v != "" {
		updateReq.ImageUrl = v
	}

	provinceName := strings.TrimSpace(firstNonEmpty(row["province"], row["province_name"]))
	provinceId := strings.TrimSpace(row["province_id"])
	if provinceName != "" {
		if provinceId == "" {
			provinceId, err = service.lookupProvinceId(req, provinceName)
			if err != nil {
				if isMasterNotFoundError(err) {
					addMissing("province", provinceName)
				} else {
					return wrapEmployeeError(row, rowNumber, err)
				}
			}
		}
		updateReq.ProvinceId = provinceId
		row["province_id"] = provinceId
	}
	if provinceId != "" && provinceName == "" {
		updateReq.ProvinceId = provinceId
	}

	cityName := strings.TrimSpace(firstNonEmpty(row["city"], row["city_regency"], row["city_name"]))
	cityId := strings.TrimSpace(row["city_id"])
	if cityName != "" {
		if cityId == "" {
			cityId, err = service.lookupCityId(req, cityName)
			if err != nil {
				if isMasterNotFoundError(err) {
					addMissing("city", cityName)
				} else {
					return wrapEmployeeError(row, rowNumber, err)
				}
			}
		}
		updateReq.CityId = cityId
		row["city_id"] = cityId
	}
	if cityId != "" && cityName == "" {
		updateReq.CityId = cityId
	}

	subDistrictName := strings.TrimSpace(firstNonEmpty(row["sub_district"], row["sub_district_name"]))
	subDistrictId := strings.TrimSpace(row["sub_district_id"])
	if subDistrictName != "" {
		if subDistrictId == "" {
			subDistrictId, err = service.lookupSubDistrictId(req, subDistrictName)
			if err != nil {
				if isMasterNotFoundError(err) {
					addMissing("sub district", subDistrictName)
				} else {
					return wrapEmployeeError(row, rowNumber, err)
				}
			}
		}
		updateReq.SubDistrictId = subDistrictId
		row["sub_district_id"] = subDistrictId
	}
	if subDistrictId != "" && subDistrictName == "" {
		updateReq.SubDistrictId = subDistrictId
	}

	wardName := strings.TrimSpace(firstNonEmpty(row["ward"], row["village"], row["village_name"]))
	wardId := strings.TrimSpace(row["ward_id"])
	if wardName != "" {
		if wardId == "" {
			wardId, err = service.lookupWardId(req, wardName)
			if err != nil {
				if isMasterNotFoundError(err) {
					addMissing("village", wardName)
				} else {
					return wrapEmployeeError(row, rowNumber, err)
				}
			}
		}
		updateReq.WardId = wardId
		row["ward_id"] = wardId
	}
	if wardId != "" && wardName == "" {
		updateReq.WardId = wardId
	}

	if v := strings.TrimSpace(row["work_date"]); v != "" {
		if _, err := parseDate(v); err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid work_date: %w", err))
		}
		updateReq.WorkDate = v
	}
	if v := strings.TrimSpace(row["dob"]); v != "" {
		if _, err := parseDate(v); err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid dob: %w", err))
		}
		updateReq.Dob = v
	}

	if v := strings.TrimSpace(row["emp_grp_id"]); v != "" {
		ptr, err := parseOptionalInt(v)
		if err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid emp_grp_id: %w", err))
		}
		updateReq.EmpGrpId = ptr
	} else if name := strings.TrimSpace(firstNonEmpty(row["emp_grp_name"], row["job_title"])); name != "" {
		id, err := service.lookupEmpGroupId(req, name)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("job title", name)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		} else {
			idCopy := id
			updateReq.EmpGrpId = &idCopy
		}
	}

	if v := strings.TrimSpace(row["division_id"]); v != "" {
		ptr, err := parseOptionalInt(v)
		if err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid division_id: %w", err))
		}
		updateReq.DivisionId = ptr
	} else if name := strings.TrimSpace(firstNonEmpty(row["division_name"], row["division"])); name != "" {
		id, err := service.lookupDivisionId(req, name)
		if err != nil {
			if isMasterNotFoundError(err) {
				addMissing("division", name)
			} else {
				return wrapEmployeeError(row, rowNumber, err)
			}
		} else {
			idCopy := id
			updateReq.DivisionId = &idCopy
		}
	}

	if v := strings.TrimSpace(firstNonEmpty(row["is_active"], row["status"])); v != "" {
		val, err := parseBoolValue(v)
		if err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid is_active: %w", err))
		}
		updateReq.IsActive = &val
	}
	if v := strings.TrimSpace(row["is_wa_no"]); v != "" {
		val, err := parseBoolValue(v)
		if err != nil {
			return wrapEmployeeError(row, rowNumber, fmt.Errorf("invalid is_wa_no: %w", err))
		}
		updateReq.IsWaNo = val
	}

	if len(missingMasters) > 0 {
		return wrapEmployeeError(row, rowNumber, errors.New(strings.Join(missingMasters, ", ")))
	}

	if err := service.EmployeeRepository.ImportUpdate(current.EmployeeId, updateReq); err != nil {
		return wrapEmployeeError(row, rowNumber, err)
	}
	return nil
}

func parseOptionalInt(input string) (*int, error) {
	val := strings.TrimSpace(input)
	if val == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func stringPtr(value string) *string {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}
	return &v
}

func parseDate(value string) (*time.Time, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(employeeDateLayout, v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseBoolDefault(value string, defaultVal bool) (bool, error) {
	v := strings.TrimSpace(strings.ToLower(value))
	if v == "" {
		return defaultVal, nil
	}
	switch v {
	case "1", "true", "yes", "y", "active", "aktif":
		return true, nil
	case "0", "false", "no", "n", "deactive", "inactive", "nonaktif", "non-active", "non active":
		return false, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal, err
	}
	return b, nil
}

func parseBoolValue(value string) (bool, error) {
	if strings.TrimSpace(value) == "" {
		return false, errors.New("value required")
	}
	return parseBoolDefault(value, false)
}

func buildEmployeeTemp(historyId int64, custId string, row map[string]string, err error, rowNumber int) entity.ImportEmployeeTemp {
	status := "failed"
	message := strings.TrimSpace(err.Error())
	prefix := employeeValidationPrefix(row, rowNumber)
	if prefix != "" {
		lowerPrefix := strings.ToLower(prefix)
		lowerMessage := strings.ToLower(message)
		if !strings.HasPrefix(lowerMessage, lowerPrefix) {
			if message != "" {
				message = fmt.Sprintf("%s: %s", prefix, message)
			} else {
				message = prefix
			}
		}
	}
	return entity.ImportEmployeeTemp{
		HistoryID:     historyId,
		CustID:        custId,
		EmpCode:       firstNonEmpty(row["emp_code"], row["id_employee"]),
		EmpName:       row["emp_name"],
		Address:       row["address"],
		EmpTypeID:     row["emp_type_id"],
		EmpTypeCode:   row["emp_type_code"],
		EmpTypeName:   row["emp_type_name"],
		EmpGrpID:      row["emp_grp_id"],
		EmpGrpCode:    row["emp_grp_code"],
		EmpGrpName:    row["emp_grp_name"],
		DivisionID:    row["division_id"],
		DivisionCode:  row["division_code"],
		DivisionName:  row["division_name"],
		WorkDate:      row["work_date"],
		LastEducation: row["last_education"],
		Dob:           row["dob"],
		PhoneNo:       row["phone_no"],
		WaNo:          row["wa_no"],
		Email:         row["email"],
		IsActive:      row["is_active"],
		IsDel:         row["is_del"],
		DeviceID:      row["device_id"],
		MacAddress:    row["mac_address"],
		ImageURL:      row["image_url"],
		IdentityNo:    row["identity_no"],
		IsWaNo:        row["is_wa_no"],
		ProvinceID:    row["province_id"],
		Province:      row["province"],
		CityID:        row["city_id"],
		City:          row["city"],
		SubDistrictID: row["sub_district_id"],
		SubDistrict:   row["sub_district"],
		WardID:        row["ward_id"],
		Ward:          row["ward"],
		PostCode:      row["post_code"],
		StatusInsert:  status,
		ErrorMessage:  message,
	}
}

func buildEmployeeUpdateTemp(historyId int64, custId string, row map[string]string, err error, rowNumber int) entity.ImportEmployeeUpdateTemp {
	temp := buildEmployeeTemp(historyId, custId, row, err, rowNumber)
	idValue := row["id_employee"]
	if strings.TrimSpace(idValue) == "" {
		idValue = row["emp_code"]
	}
	return entity.ImportEmployeeUpdateTemp{
		HistoryID:     temp.HistoryID,
		CustID:        temp.CustID,
		EmpID:         idValue,
		EmpCode:       temp.EmpCode,
		EmpName:       temp.EmpName,
		Address:       temp.Address,
		EmpTypeID:     temp.EmpTypeID,
		EmpTypeCode:   temp.EmpTypeCode,
		EmpTypeName:   temp.EmpTypeName,
		EmpGrpID:      temp.EmpGrpID,
		EmpGrpCode:    temp.EmpGrpCode,
		EmpGrpName:    temp.EmpGrpName,
		DivisionID:    temp.DivisionID,
		DivisionCode:  temp.DivisionCode,
		DivisionName:  temp.DivisionName,
		WorkDate:      temp.WorkDate,
		LastEducation: temp.LastEducation,
		Dob:           temp.Dob,
		PhoneNo:       temp.PhoneNo,
		WaNo:          temp.WaNo,
		Email:         temp.Email,
		IsActive:      temp.IsActive,
		IsDel:         temp.IsDel,
		DeviceID:      temp.DeviceID,
		MacAddress:    temp.MacAddress,
		ImageURL:      temp.ImageURL,
		IdentityNo:    temp.IdentityNo,
		IsWaNo:        temp.IsWaNo,
		ProvinceID:    temp.ProvinceID,
		Province:      temp.Province,
		CityID:        temp.CityID,
		City:          temp.City,
		SubDistrictID: temp.SubDistrictID,
		SubDistrict:   temp.SubDistrict,
		WardID:        temp.WardID,
		Ward:          temp.Ward,
		PostCode:      temp.PostCode,
		StatusInsert:  temp.StatusInsert,
		ErrorMessage:  temp.ErrorMessage,
	}
}
