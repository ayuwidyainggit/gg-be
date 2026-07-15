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
	IsActive     *int   `query:"is_active"`
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
	EmpTypeId     string `json:"emp_type_id"`
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
	EmpTypeId     string `json:"emp_type_id,omitempty"`
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
}
