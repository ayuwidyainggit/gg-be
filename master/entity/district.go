package entity

import (
	"time"
)

type DistrictResponse struct {
	DistrictId    int        `json:"district_id"`
	DistrictCode  string     `json:"district_code"`
	DistrictName  string     `json:"district_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type DistrictLookupResponse struct {
	DistrictId   int    `json:"district_id"`
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

type CreateDistrictBody struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	CreatedBy    int64  `json:"created_by" validate:"required"`
	DistrictCode string `json:"district_code" validate:"required,max=5,alphanumericSpace"`
	DistrictName string `json:"district_name" validate:"required,max=40,alphanumericSpace"`
	IsActive     bool   `json:"is_active"`
}

type DetailDistrictParams struct {
	DistrictId int `params:"district_id" validate:"required"`
}

type UpdateDistrictParams struct {
	DistrictId int `params:"district_id" validate:"required"`
}

type DeleteDistrictParams struct {
	DistrictId int `params:"district_id" validate:"required"`
}

type UpdateDistrictRequest struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	DistrictCode string `json:"district_code,omitempty" validate:"required,max=5,alphanumericSpace"`
	DistrictName string `json:"district_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive     *bool  `json:"is_active,omitempty"`
}
