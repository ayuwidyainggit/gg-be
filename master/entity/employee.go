package entity

import (
	"time"
)

type EmployeeQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	EmpTypeId    string `query:"emp_type_id"`
	EmpGrpId     int    `query:"emp_grp_id"`
	EmpGrpName   string `query:"emp_grp_name"`
	DivisionId   int    `query:"division_id"`
	NoSalesman   bool   `query:"no_salesman"`
	IsActive     *int   `query:"is_active"`
	Status       string `query:"status"`
	Format       string `query:"format"`
}

type EmployeeResponse struct {
	EmployeeId    int        `json:"emp_id"`
	EmployeeCode  string     `json:"emp_code"`
	EmployeeName  string     `json:"emp_name"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	LastEducation string     `json:"last_education"`
	PhoneNo       string     `json:"phone_no"`
	WaNo          string     `json:"wa_no"`
	Email         string     `json:"email"`
	EmpTypeId     string     `json:"emp_type_id"`
	EmpTypeName   string     `json:"emp_type_name"`
	EmpGrpId      *int       `json:"emp_grp_id"`
	EmpGrpCode    *string    `json:"emp_grp_code"`
	EmpGrpName    *string    `json:"emp_grp_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	Dob           *string    `json:"dob"`
	WorkDate      *string    `json:"work_date"`
	ImageUrl      *string    `json:"image_url"`
	UpdatedAt     *time.Time `json:"updated_at"`
	ProvinceId    string     `json:"province_id"`
	Province      string     `json:"province"`
	CityId        string     `json:"city_id"`
	SubDistrictId string     `json:"sub_district_id"`
	SubDistrict   string     `json:"sub_district"`
	WardId        string     `json:"ward_id"`
	Ward          string     `json:"ward"`
	IdentityNo    string     `json:"identity_no"`
	IsWaNo        bool       `json:"is_wa_no"`
	PostCode      string     `json:"post_code"`
	DivisionId    *int       `json:"division_id"`
	DivisionName  *string    `json:"division_name"`

	RegionScope      string                               `json:"region_scope,omitempty"`
	AreaScope        string                               `json:"area_scope,omitempty"`
	DistributorScope string                               `json:"distributor_scope,omitempty"`
	Regions          []RegionIdCodeNameResp               `json:"regions,omitempty"`
	Areas            []EmployeeAreaMappingResponse        `json:"areas,omitempty"`
	Distributors     []EmployeeDistributorMappingResponse `json:"distributors,omitempty"`
}

type EmployeeAreaMappingResponse struct {
	AreaId     int    `json:"area_id"`
	AreaCode   string `json:"area_code"`
	AreaName   string `json:"area_name"`
	RegionId   int    `json:"region_id"`
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
}

type EmployeeDistributorMappingResponse struct {
	DistributorId   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	AreaId          int    `json:"area_id"`
	AreaCode        string `json:"area_code"`
	AreaName        string `json:"area_name"`
	RegionId        int    `json:"region_id"`
	RegionCode      string `json:"region_code"`
	RegionName      string `json:"region_name"`
}

type EmployeeLookupResponse struct {
	EmployeeId    int     `json:"emp_id"`
	EmployeeCode  string  `json:"emp_code"`
	EmployeeName  string  `json:"emp_name"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	LastEducation string  `json:"last_education"`
	PhoneNo       string  `json:"phone_no"`
	WaNo          string  `json:"wa_no"`
	Email         string  `json:"email"`
	EmpTypeId     string  `json:"emp_type_id"`
	EmpTypeName   string  `json:"emp_type_name"`
	EmpGrpId      *int    `json:"emp_grp_id"`
	EmpGrpCode    *string `json:"emp_grp_code"`
	EmpGrpName    *string `json:"emp_grp_name"`
	Dob           *string `json:"dob"`
	WorkDate      *string `json:"work_date"`
	ImageUrl      string  `json:"image_url"`
	ProvinceId    *string `json:"province_id"`
	Province      string  `json:"province"`
	CityId        string  `json:"city_id"`
	SubDistrictId string  `json:"sub_district_id"`
	SubDistrict   string  `json:"sub_district"`
	WardId        string  `json:"ward_id"`
	Ward          string  `json:"ward"`
	IdentityNo    string  `json:"identity_no"`
	IsWaNo        bool    `json:"is_wa_no"`
	PostCode      string  `json:"post_code"`
	DivisionId    *int    `json:"division_id"`
	DivisionName  *string `json:"division_name"`
	EmpIdSalesman int     `json:"emp_id_salesman"`
	SalesName     string  `json:"sales_name"`
	DeleteSales   bool    `json:"delete_sales"`
}

type CreateMultipleEmployeeBody struct {
	CustId       string
	ParentCustId string
	Employees    []CreateEmployeeBody `json:"employees"`
}
type CreateEmployeeBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	ParentCustId  string `json:"parent_cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	EmployeeCode  string `json:"emp_code" validate:"required,max=10,alphanumericSpace"`
	EmployeeName  string `json:"emp_name" validate:"required,max=100"`
	Address       string `json:"address" validate:"required,max=150"`
	LastEducation string `json:"last_education"`
	PhoneNo       string `json:"phone_no" validate:"numeric"`
	WaNo          string `json:"wa_no" validate:"numeric"`
	Email         string `json:"email" validate:"email"`
	EmpTypeId     string `json:"emp_type_id" validate:"omitempty,max=1"`
	EmpGrpId      int    `json:"emp_grp_id"`
	IsActive      bool   `json:"is_active"`
	UpdatedBy     *int64 `json:"updated_by"`
	Dob           string `json:"dob"`
	WorkDate      string `json:"work_date"`
	ImageUrl      string `json:"image_url"`
	ProvinceId    string `json:"province_id" validate:"required"`
	CityId        string `json:"city_id" validate:"required,alphanumericSpace"`
	SubDistrictId string `json:"sub_district_id" validate:"required,alphanumericSpace"`
	WardId        string `json:"ward_id" validate:"required,alphanumericSpace"`
	IdentityNo    string `json:"identity_no" validate:"alphanumericSpace"`
	IsWaNo        bool   `json:"is_wa_no"`
	PostCode      string `json:"post_code" validate:"numeric"`
	DivisionId    int    `json:"division_id"`

	RegionScope      string `json:"region_scope"`
	AreaScope        string `json:"area_scope"`
	DistributorScope string `json:"distributor_scope"`
	RegionIds        []int  `json:"region_ids" validate:"omitempty,dive,gt=0"`
	AreaIds          []int  `json:"area_ids" validate:"omitempty,dive,gt=0"`
	DistributorIds   []int  `json:"distributor_ids" validate:"omitempty,dive,gt=0"`
}

type DetailEmployeeParams struct {
	CustId       string
	ParentCustId string
	EmployeeCode string
	EmployeeId   int `params:"emp_id" validate:"required"`
}

type UpdateEmployeeParams struct {
	EmployeeId int `params:"emp_id" validate:"required"`
}

type DeleteEmployeeParams struct {
	EmployeeId int `params:"emp_id" validate:"required"`
}

type UpdateEmployeeRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	ParentCustId  string `json:"parent_cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	EmployeeCode  string `json:"emp_code,omitempty" validate:"max=10,alphanumericSpace"`
	EmployeeName  string `json:"emp_name,omitempty" validate:"max=100"`
	IsActive      *bool  `json:"is_active,omitempty"`
	Address       string `json:"address,omitempty" validate:"max=150"`
	City          string `json:"city,omitempty"`
	LastEducation string `json:"last_education,omitempty"`
	PhoneNo       string `json:"phone_no,omitempty" validate:"numeric"`
	WaNo          string `json:"wa_no,omitempty" validate:"numeric"`
	Email         string `json:"email,omitempty" validate:"email"`
	EmpTypeId     string `json:"emp_type_id,omitempty" validate:"omitempty,max=1"`
	EmpGrpId      *int   `json:"emp_grp_id"`
	Dob           string `json:"dob,omitempty"`
	WorkDate      string `json:"work_date,omitempty"`
	ImageUrl      string `json:"image_url"`
	ProvinceId    string `json:"province_id,omitempty" validate:"required,alphanumericSpace"`
	Province      string `json:"province,omitempty" validate:"max=50,omitempty"`
	CityId        string `json:"city_id" validate:"required,alphanumericSpace"`
	SubDistrictId string `json:"sub_district_id,omitempty" validate:"required,alphanumericSpace"`
	SubDistrict   string `json:"sub_district,omitempty"`
	WardId        string `json:"ward_id,omitempty" validate:"required,alphanumericSpace"`
	Ward          string `json:"ward,omitempty"`
	IdentityNo    string `json:"identity_no,omitempty"`
	IsWaNo        bool   `json:"is_wa_no"`
	PostCode      string `json:"post_code,omitempty" validate:"numeric"`
	DivisionId    *int   `json:"division_id"`

	RegionScope      string `json:"region_scope,omitempty"`
	AreaScope        string `json:"area_scope,omitempty"`
	DistributorScope string `json:"distributor_scope,omitempty"`
	RegionIds        []int  `json:"region_ids" validate:"omitempty,dive,gt=0"`
	AreaIds          []int  `json:"area_ids" validate:"omitempty,dive,gt=0"`
	DistributorIds   []int  `json:"distributor_ids" validate:"omitempty,dive,gt=0"`

	TerritoryMappingProvided bool `json:"-"`
}

type ImportEmployeeTemp struct {
	HistoryID     int64  `db:"history_id"`
	CustID        string `db:"cust_id"`
	EmpCode       string `db:"emp_code"`
	EmpName       string `db:"emp_name"`
	Address       string `db:"address"`
	EmpTypeID     string `db:"emp_type_id"`
	EmpTypeCode   string `db:"emp_type_code"`
	EmpTypeName   string `db:"emp_type_name"`
	EmpGrpID      string `db:"emp_grp_id"`
	EmpGrpCode    string `db:"emp_grp_code"`
	EmpGrpName    string `db:"emp_grp_name"`
	DivisionID    string `db:"division_id"`
	DivisionCode  string `db:"division_code"`
	DivisionName  string `db:"division_name"`
	WorkDate      string `db:"work_date"`
	LastEducation string `db:"last_education"`
	Dob           string `db:"dob"`
	PhoneNo       string `db:"phone_no"`
	WaNo          string `db:"wa_no"`
	Email         string `db:"email"`
	IsActive      string `db:"is_active"`
	IsDel         string `db:"is_del"`
	DeviceID      string `db:"device_id"`
	MacAddress    string `db:"mac_address"`
	ImageURL      string `db:"image_url"`
	IdentityNo    string `db:"identity_no"`
	IsWaNo        string `db:"is_wa_no"`
	ProvinceID    string `db:"province_id"`
	Province      string `db:"province"`
	CityID        string `db:"city_id"`
	City          string `db:"city"`
	SubDistrictID string `db:"sub_district_id"`
	SubDistrict   string `db:"sub_district"`
	WardID        string `db:"ward_id"`
	Ward          string `db:"ward"`
	PostCode      string `db:"post_code"`
	StatusInsert  string `db:"status_insert"`
	ErrorMessage  string `db:"error_message"`
}

// EmployeePJPQueryFilter for employee-pjp endpoint
type EmployeePJPQueryFilter struct {
	CustId        string  // from JWT context
	ParentCustId  string  // from JWT context
	Page          int     `query:"page"`
	Limit         int     `query:"limit"`
	Query         string  `query:"q"`
	Sort          string  `query:"sort"`
	IsActive      []int   `query:"is_active"`      // Array: 0=inactive, 1=active
	FilterCustId  *string `query:"cust_id"`        // Optional: STRING type (e.g., 'C22001')
	DistributorId *int    `query:"distributor_id"` // Optional: Integer
}

// EmployeePJPResponse simplified response for employee-pjp endpoint
type EmployeePJPResponse struct {
	EmpId   int    `json:"emp_id"`
	EmpCode string `json:"emp_code"`
	EmpName string `json:"emp_name"`
}

type ImportEmployeeUpdateTemp struct {
	HistoryID     int64  `db:"history_id"`
	CustID        string `db:"cust_id"`
	EmpID         string `db:"emp_id"`
	EmpCode       string `db:"emp_code"`
	EmpName       string `db:"emp_name"`
	Address       string `db:"address"`
	EmpTypeID     string `db:"emp_type_id"`
	EmpTypeCode   string `db:"emp_type_code"`
	EmpTypeName   string `db:"emp_type_name"`
	EmpGrpID      string `db:"emp_grp_id"`
	EmpGrpCode    string `db:"emp_grp_code"`
	EmpGrpName    string `db:"emp_grp_name"`
	DivisionID    string `db:"division_id"`
	DivisionCode  string `db:"division_code"`
	DivisionName  string `db:"division_name"`
	WorkDate      string `db:"work_date"`
	LastEducation string `db:"last_education"`
	Dob           string `db:"dob"`
	PhoneNo       string `db:"phone_no"`
	WaNo          string `db:"wa_no"`
	Email         string `db:"email"`
	IsActive      string `db:"is_active"`
	IsDel         string `db:"is_del"`
	DeviceID      string `db:"device_id"`
	MacAddress    string `db:"mac_address"`
	ImageURL      string `db:"image_url"`
	IdentityNo    string `db:"identity_no"`
	IsWaNo        string `db:"is_wa_no"`
	ProvinceID    string `db:"province_id"`
	Province      string `db:"province"`
	CityID        string `db:"city_id"`
	City          string `db:"city"`
	SubDistrictID string `db:"sub_district_id"`
	SubDistrict   string `db:"sub_district"`
	WardID        string `db:"ward_id"`
	Ward          string `db:"ward"`
	PostCode      string `db:"post_code"`
	StatusInsert  string `db:"status_insert"`
	ErrorMessage  string `db:"error_message"`
}

// EmployeeLookupAPIQuery query params for GET /v1/employee-lookup.
type EmployeeLookupAPIQuery struct {
	Q     string `query:"q"`
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

// EmployeeLookupAPIFilter for GET /v1/employee-lookup.
type EmployeeLookupAPIFilter struct {
	CustId        string
	ParentCustId  string
	Query         string
	Page          int
	Limit         int
	Sort          string
	FilterCustIds []string
}

// EmployeeLookupMinimalItem response shape for employee-lookup (minimal fields).
type EmployeeLookupMinimalItem struct {
	EmpId   int    `json:"emp_id"`
	EmpCode string `json:"emp_code"`
	EmpName string `json:"emp_name"`
}
