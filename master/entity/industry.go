package entity

import (
	"time"
)

type IndustryResponse struct {
	IndustryId    int        `json:"industry_id"`
	IndustryCode  string     `json:"industry_code"`
	IndustryName  string     `json:"industry_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type IndustryLookupResponse struct {
	IndustryId   int    `json:"industry_id"`
	IndustryCode string `json:"industry_code"`
	IndustryName string `json:"industry_name"`
}

type CreateIndustryBody struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	CreatedBy    int64  `json:"created_by" validate:"required"`
	IndustryCode string `json:"industry_code" validate:"required,max=5,alphanumericSpace"`
	IndustryName string `json:"industry_name" validate:"required,max=40,alphanumericSpace"`
	IsActive     bool   `json:"is_active"`
}

type DetailIndustryParams struct {
	IndustryId int `params:"industry_id" validate:"required"`
}

type UpdateIndustryParams struct {
	IndustryId int `params:"industry_id" validate:"required"`
}

type DeleteIndustryParams struct {
	IndustryId int `params:"industry_id" validate:"required"`
}

type UpdateIndustryRequest struct {
	CustId       string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64  `json:"updated_by" validate:"required"`
	IndustryCode string `json:"industry_code,omitempty" validate:"required,max=5,alphanumericSpace"`
	IndustryName string `json:"industry_name,omitempty" validate:"max=40,alphanumericSpace"`
	IsActive     *bool  `json:"is_active,omitempty"`
}
