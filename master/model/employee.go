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
	EmpIdSalesman int        `db:"emp_id_salesman" json:"emp_id_salesman"`
	SalesName     string     `db:"sales_name" json:"sales_name"`
	DeleteSales   bool       `db:"delete_sales" json:"delete_sales"`

	RegionScope      string `db:"region_scope" json:"region_scope"`
	AreaScope        string `db:"area_scope" json:"area_scope"`
	DistributorScope string `db:"distributor_scope" json:"distributor_scope"`
	RegionIds        []int  `db:"-" json:"region_ids"`
	AreaIds          []int  `db:"-" json:"area_ids"`
	DistributorIds   []int  `db:"-" json:"distributor_ids"`
}

type EmployeeTerritoryMapping struct {
	RegionIds      []int
	AreaIds        []int
	DistributorIds []int
}

type EmployeeTerritoryDetail struct {
	Regions      []EmployeeRegionMappingDetail
	Areas        []EmployeeAreaMappingDetail
	Distributors []EmployeeDistributorMappingDetail
}

type EmployeeRegionMappingDetail struct {
	RegionId   int    `db:"region_id" json:"region_id"`
	RegionCode string `db:"region_code" json:"region_code"`
	RegionName string `db:"region_name" json:"region_name"`
}

type EmployeeAreaMappingDetail struct {
	AreaId     int    `db:"area_id" json:"area_id"`
	AreaCode   string `db:"area_code" json:"area_code"`
	AreaName   string `db:"area_name" json:"area_name"`
	RegionId   int    `db:"region_id" json:"region_id"`
	RegionCode string `db:"region_code" json:"region_code"`
	RegionName string `db:"region_name" json:"region_name"`
}

type EmployeeDistributorMappingDetail struct {
	DistributorId   int    `db:"distributor_id" json:"distributor_id"`
	DistributorCode string `db:"distributor_code" json:"distributor_code"`
	DistributorName string `db:"distributor_name" json:"distributor_name"`
	AreaId          int    `db:"area_id" json:"area_id"`
	AreaCode        string `db:"area_code" json:"area_code"`
	AreaName        string `db:"area_name" json:"area_name"`
	RegionId        int    `db:"region_id" json:"region_id"`
	RegionCode      string `db:"region_code" json:"region_code"`
	RegionName      string `db:"region_name" json:"region_name"`
}

type EmployeeExport struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	EmployeeId    int64      `db:"emp_id" json:"emp_id"`
	EmployeeCode  *string    `db:"emp_code" json:"emp_code"`
	EmployeeName  *string    `db:"emp_name" json:"emp_name"`
	Address       *string    `db:"address" json:"address"`
	EmpTypeId     *string    `db:"emp_type_id" json:"emp_type_id"`
	EmpTypeName   *string    `db:"emp_type_name" json:"emp_type_name"`
	EmpGrpId      *int64     `db:"emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode    *string    `db:"emp_grp_code" json:"emp_grp_code"`
	EmpGrpName    *string    `db:"emp_grp_name" json:"emp_grp_name"`
	WorkDate      *time.Time `db:"work_date" json:"work_date"`
	LastEducation *string    `db:"last_education" json:"last_education"`
	Dob           *time.Time `db:"dob" json:"dob"`
	PhoneNo       *string    `db:"phone_no" json:"phone_no"`
	WaNo          *string    `db:"wa_no" json:"wa_no"`
	Email         *string    `db:"email" json:"email"`
	IsActive      *bool      `db:"is_active" json:"is_active"`
	IsDel         *bool      `db:"is_del" json:"is_del"`
	DeviceID      *string    `db:"device_id" json:"device_id"`
	MacAddress    *string    `db:"mac_address" json:"mac_address"`
	ImageURL      *string    `db:"image_url" json:"image_url"`
	IdentityNo    *string    `db:"identity_no" json:"identity_no"`
	IsWaNo        *bool      `db:"is_wa_no" json:"is_wa_no"`
	ProvinceID    *string    `db:"province_id" json:"province_id"`
	Province      *string    `db:"province" json:"province"`
	CityID        *string    `db:"city_id" json:"city_id"`
	City          *string    `db:"city" json:"city"`
	SubDistrictID *string    `db:"sub_district_id" json:"sub_district_id"`
	SubDistrict   *string    `db:"sub_district" json:"sub_district"`
	WardID        *string    `db:"ward_id" json:"ward_id"`
	Ward          *string    `db:"ward" json:"ward"`
	PostCode      *string    `db:"post_code" json:"post_code"`
	DivisionID    *int64     `db:"division_id" json:"division_id"`
	DivisionCode  *string    `db:"division_code" json:"division_code"`
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

	RegionScope      *string `json:"region_scope,omitempty" sql:"region_scope"`
	AreaScope        *string `json:"area_scope,omitempty" sql:"area_scope"`
	DistributorScope *string `json:"distributor_scope,omitempty" sql:"distributor_scope"`
}

// EmployeePJP simplified model for PJP list
type EmployeePJP struct {
	EmpId   int    `db:"emp_id"`
	EmpCode string `db:"emp_code"`
	EmpName string `db:"emp_name"`
}

// EmployeeLookupMinimal row for GET /v1/employee-lookup.
type EmployeeLookupMinimal struct {
	EmployeeId   int    `db:"emp_id"`
	EmployeeCode string `db:"emp_code"`
	EmployeeName string `db:"emp_name"`
}
