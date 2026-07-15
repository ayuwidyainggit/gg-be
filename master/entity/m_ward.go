package entity

import "time"

type WardResponse struct {
	WardId        string     `json:"ward_id"`
	Ward          string     `json:"ward"`
	ProvinceId    string     `json:"province_id"`
	Province      string     `json:"province"`
	RegencyId     string     `json:"regency_id"`
	Regency       string     `json:"regency"`
	SubDistrictId string     `json:"sub_district_id"`
	SubDistrict   string     `json:"sub_district"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type WardLookupResponse struct {
	WardId        string `json:"ward_id"`
	Ward          string `json:"ward"`
	ProvinceId    string `json:"province_id"`
	Province      string `json:"province"`
	RegencyId     string `json:"regency_id"`
	Regency       string `json:"regency"`
	SubDistrictId string `json:"sub_district_id"`
	SubDistrict   string `json:"sub_district"`
}

type CreateWardBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	WardId        string `json:"ward_id" validate:"required,alphanumericSpace"`
	Ward          string `json:"ward" validate:"required"`
	ProvinceId    string `json:"province_id" validate:"required,alphanumericSpace"`
	RegencyId     string `json:"regency_id" validate:"required,alphanumericSpace"`
	SubDistrictId string `json:"sub_district_id" validate:"required,alphanumericSpace"`
	IsActive      bool   `json:"is_active"`
}

type DetailWardParams struct {
	WardId string `params:"ward_id" validate:"required"`
}

type UpdateWardParams struct {
	WardId string `params:"ward_id" validate:"required"`
}

type DeleteWardParams struct {
	WardId string `params:"ward_id" validate:"required"`
}

type UpdateWardRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	WardId        string `json:"ward_id,omitempty" validate:"required,alphanumericSpace"`
	Ward          string `json:"ward,omitempty" validate:"required"`
	ProvinceId    string `json:"province_id,omitempty" validate:"required,alphanumericSpace"`
	RegencyId     string `json:"regency_id" validate:"required,alphanumericSpace"`
	SubDistrictId string `json:"sub_district_id" validate:"required,alphanumericSpace"`
	IsActive      *bool  `json:"is_active,omitempty"`
}

type WardQueryFilter struct {
	CustId        string
	ParentId      string
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
	ProvinceId    string `query:"province_id"`
	RegencyId     string `query:"regency_id"`
	SubDistrictId string `query:"sub_district_id"`
	IsActive      *int   `query:"is_active"`
}
