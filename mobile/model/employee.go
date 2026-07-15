package model

import (
	"time"
)

type Employee struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	EmployeeId    int        `db:"emp_id" json:"emp_id"`
	EmployeeCode  string     `db:"emp_code" json:"emp_code"`
	EmployeeName  string     `db:"emp_name" json:"emp_name"`
	Address       *string    `db:"address" json:"address"`
	City          *string    `db:"city" json:"city"`
	LastEducation *string    `db:"last_education" json:"last_education"`
	PhoneNo       *string    `db:"phone_no" json:"phone_no"`
	WaNo          *string    `db:"wa_no" json:"wa_no"`
	Email         *string    `db:"email" json:"email"`
	EmpTypeId     *string    `db:"emp_type_id" json:"emp_type_id"`
	EmpTypeName   *string    `db:"emp_type_name" json:"emp_type_name"`
	EmpGrpId      *int       `db:"emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode    *string    `db:"emp_grp_code" json:"emp_grp_code"`
	EmpGrpName    *string    `db:"emp_grp_name" json:"emp_grp_name"`
	Dob           *time.Time `db:"dob" json:"dob"`
	WorkDate      *time.Time `db:"work_date" json:"work_date"`
	ImageUrl      *string    `db:"image_url" json:"image_url"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by,omitempty"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	ProvinceId    *string    `db:"province_id" json:"province_id"`
	Province      *string    `db:"province" json:"province"`
	CityId        *string    `db:"city_id" json:"city_id"`
	SubDistrictId *string    `db:"sub_district_id" json:"sub_district_id"`
	SubDistrict   *string    `db:"sub_district" json:"sub_district"`
	WardId        *string    `db:"ward_id" json:"ward_id"`
	Ward          *string    `db:"ward" json:"ward"`
	PostCode      *string    `db:"post_code" json:"post_code"`
	IdentityNo    *string    `db:"identity_no" json:"identity_no"`
	IsWaNo        bool       `db:"is_wa_no" json:"is_wa_no"`
	DivisionId    *int       `db:"division_id" json:"division_id"`
	DivisionName  *string    `db:"division_name" json:"division_name"`
}

type EmployeeUpdate struct {
	EmployeeCode  *string    `json:"emp_code,omitempty" sql:"emp_code"`
	EmployeeName  *string    `json:"emp_name,omitempty" sql:"emp_name"`
	Address       *string    `json:"address" sql:"address"`
	City          *string    `json:"city" sql:"city"`
	LastEducation *string    `json:"last_education" sql:"last_education"`
	PhoneNo       *string    `json:"phone_no" sql:"phone_no"`
	WaNo          *string    `json:"wa_no" sql:"wa_no"`
	Email         *string    `json:"email" sql:"email"`
	EmpTypeId     *string    `json:"emp_type_id" sql:"emp_type_id"`
	EmpGrpId      *int       `json:"emp_grp_id" sql:"emp_grp_id"`
	Dob           *string    `json:"dob" sql:"dob"`
	WorkDate      *string    `json:"work_date" sql:"work_date"`
	ImageUrl      *string    `json:"image_url" sql:"image_url"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
	ProvinceId    *string    `json:"province_id,omitempty" sql:"province_id"`
	Province      *string    `json:"province,omitempty" sql:"province"`
	CityId        *string    `json:"city_id,omitempty" sql:"city_id"`
	SubDistrictId *string    `json:"sub_district_id,omitempty" sql:"sub_district_id"`
	SubDistrict   *string    `json:"sub_district,omitempty" sql:"sub_district"`
	WardId        *string    `json:"ward_id,omitempty" sql:"ward_id"`
	Ward          *string    `json:"ward,omitempty" sql:"ward"`
	PostCode      *string    `json:"post_code" sql:"post_code" `
	IdentityNo    *string    `json:"identity_no" sql:"identity_no"`
	IsWaNo        *bool      `json:"is_wa_no" sql:"is_wa_no"`
	DivisionId    *int       `json:"division_id" sql:"division_id"`
}
