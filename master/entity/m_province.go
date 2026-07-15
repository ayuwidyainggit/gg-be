package entity

import "time"

type ProvinceResponse struct {
	ProvinceId    string     `json:"province_id"`
	Province      string     `json:"province"`
	IsActive      bool       `json:"is_active"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type ProvinceLookupResponse struct {
	ProvinceId string `json:"province_id"`
	Province   string `json:"province"`
}

type CreateProvinceBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	ProvinceId string `json:"province_id" validate:"required,alphanumericSpace,max=4"`
	Province   string `json:"province" validate:"required,max=30,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailProvinceParams struct {
	ProvinceId string `params:"province_id" validate:"required"`
}

type UpdateProvinceParams struct {
	ProvinceId string `params:"province_id" validate:"required"`
}

type DeleteProvinceParams struct {
	ProvinceId string `params:"province_id" validate:"required"`
}

type UpdateProvinceRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	ProvinceId string `json:"province_id,omitempty" validate:"required,alphanumericSpace,max=4"`
	Province   string `json:"province,omitempty" validate:"max=30,omitempty,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
