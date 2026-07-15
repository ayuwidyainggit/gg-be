package entity

import "time"

type RegencyResponse struct {
	RegencyId     string     `json:"regency_id"`
	Regency       string     `json:"regency"`
	ProvinceId    string     `json:"province_id"`
	Province      string     `json:"province"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type RegencyLookupResponse struct {
	RegencyId  string `json:"regency_id"`
	Regency    string `json:"regency"`
	ProvinceId string `json:"province_id"`
	Province   string `json:"province"`
}

type CreateRegencyBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	RegencyId  string `json:"regency_id" validate:"required,alphanumericSpace,max=4"`
	Regency    string `json:"regency" validate:"required,max=50,alphanumericSpace"`
	ProvinceId string `json:"province_id" validate:"required,alphanumericSpace,max=4"`
	IsActive   bool   `json:"is_active"`
}

type DetailRegencyParams struct {
	RegencyId string `params:"regency_id" validate:"required"`
}

type UpdateRegencyParams struct {
	RegencyId string `params:"regency_id" validate:"required"`
}

type DeleteRegencyParams struct {
	RegencyId string `params:"regency_id" validate:"required"`
}

type UpdateRegencyRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	RegencyId  string `json:"regency_id" validate:"required,alphanumericSpace,max=4"`
	Regency    string `json:"regency" validate:"required,max=50,alphanumericSpace"`
	ProvinceId string `json:"province_id,omitempty" validate:"required,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type RegencyQueryFilter struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit" validate:"required"`
	Query      string `query:"q"`
	Mode       string `query:"mode"`
	Sort       string `query:"sort"`
	IsActive   *int   `query:"is_active"`
	ProvinceId string `query:"province_id"`
}
