package entity

import "time"

type SubDistrictResponse struct {
	SubDistrictId string     `json:"sub_district_id"`
	SubDistrict   string     `json:"sub_district"`
	ProvinceId    string     `json:"province_id"`
	ProvinceName  *string    `json:"province"`
	RegencyId     string     `json:"regency_id"`
	RegencyName   string     `json:"regency"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type SubDistrictLookupResponse struct {
	SubDistrictId string  `json:"sub_district_id"`
	SubDistrict   string  `json:"sub_district"`
	ProvinceId    string  `json:"province_id"`
	ProvinceName  *string `json:"province"`
	RegencyId     string  `json:"regency_id"`
	RegencyName   string  `json:"regency"`
}

type CreateSubDistrictBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	SubDistrictId string `json:"sub_district_id" validate:"required,alphanumericSpace"`
	SubDistrict   string `json:"sub_district" validate:"required"`
	ProvinceId    string `json:"province_id" validate:"required,alphanumericSpace"`
	RegencyId     string `json:"regency_id" validate:"required,alphanumericSpace"`
	IsActive      bool   `json:"is_active"`
}

type DetailSubDistrictParams struct {
	SubDistrictId string `params:"sub_district_id" validate:"required"`
}

type UpdateSubDistrictParams struct {
	SubDistrictId string `params:"sub_district_id" validate:"required"`
}

type DeleteSubDistrictParams struct {
	SubDistrictId string `params:"sub_district_id" validate:"required"`
}

type UpdateSubDistrictRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	SubDistrictId string `json:"sub_district_id,omitempty" validate:"required,alphanumericSpace"`
	SubDistrict   string `json:"sub_district,omitempty" validate:"required"`
	ProvinceId    string `json:"province_id,omitempty" validate:"required,alphanumericSpace"`
	RegencyId     string `json:"regency_id" validate:"required,alphanumericSpace"`
	IsActive      *bool  `json:"is_active,omitempty"`
}

type SubDistrictQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	ProvinceId   string `query:"province_id"`
	RegencyId    string `query:"regency_id"`
	IsActive     *int   `query:"is_active"`
}
