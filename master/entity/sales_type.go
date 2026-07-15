package entity

import (
	"time"
)

type SalesTypeResponse struct {
	SalesTypeId   int        `json:"sales_type_id"`
	SalesTypeCode string     `json:"sales_type_code"`
	SalesTypeName string     `json:"sales_type_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type SalesTypeLookupResponse struct {
	SalesTypeId   int    `json:"sales_type_id"`
	SalesTypeCode string `json:"sales_type_code"`
	SalesTypeName string `json:"sales_type_name"`
}

type CreateSalesTypeBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	SalesTypeCode string `json:"sales_type_code" validate:"required,max=1,alphanum"`
	SalesTypeName string `json:"sales_type_name" validate:"required,max=25,alphanumericSpace"`
	IsActive      bool   `json:"is_active"`
}

type DetailSalesTypeParams struct {
	SalesTypeId int `params:"sales_type_id" validate:"required"`
}

type UpdateSalesTypeParams struct {
	SalesTypeId int `params:"sales_type_id" validate:"required"`
}

type DeleteSalesTypeParams struct {
	SalesTypeId int `params:"sales_type_id" validate:"required"`
}

type UpdateSalesTypeRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	SalesTypeCode string `json:"sales_type_code,omitempty" validate:"required,max=3,numeric"`
	SalesTypeName string `json:"sales_type_name,omitempty" validate:"max=25,alphanumericSpace,omitempty"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
