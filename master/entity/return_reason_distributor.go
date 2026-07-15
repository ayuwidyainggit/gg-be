package entity

import (
	"time"
)

type ReturnReasonDistributorQueryFilter struct {
	Page                        int    `query:"page"`
	Limit                       int    `query:"limit" validate:"required"`
	Query                       string `query:"q"`
	Mode                        string `query:"mode"`
	Sort                        string `query:"sort"`
	ReturnReasonDistributorType string `query:"return_reason_type"`
	IsActive                    *int   `query:"is_active"`
}

type ReturnReasonDistributorResponse struct {
	ReturnReasonDistributorId   int        `json:"return_reason_id"`
	ReturnReasonDistributorCode string     `json:"return_reason_code"`
	ReturnReasonDistributorName string     `json:"return_reason_name"`
	ReturnReasonDistributorType *string    `json:"return_reason_type"`
	IsActive                    bool       `json:"is_active"`
	UpdatedBy                   *int64     `json:"updated_by"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	UpdatedByName               string     `json:"updated_by_name"`
}

type ReturnReasonDistributorLookupResponse struct {
	ReturnReasonDistributorId   int    `json:"return_reason_id"`
	ReturnReasonDistributorCode string `json:"return_reason_code"`
	ReturnReasonDistributorName string `json:"return_reason_name"`
	ReturnReasonDistributorType string `json:"return_reason_type"`
}

type CreateReturnReasonDistributorBody struct {
	CustId                      string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy                   int64   `json:"created_by" validate:"required"`
	ReturnReasonDistributorCode string  `json:"return_reason_code" validate:"required,max=10,alphanumericSpace"`
	ReturnReasonDistributorName string  `json:"return_reason_name" validate:"required,max=150"`
	ReturnReasonDistributorType *string `json:"return_reason_type" validate:"required,oneof='A' 'D' 'P' 'S'"`
	IsActive                    bool    `json:"is_active"`
}

type DetailReturnReasonDistributorParams struct {
	ReturnReasonDistributorId int `params:"return_reason_id" validate:"required"`
}

type UpdateReturnReasonDistributorParams struct {
	ReturnReasonDistributorId int `params:"return_reason_id" validate:"required"`
}

type DeleteReturnReasonDistributorParams struct {
	ReturnReasonDistributorId int `params:"return_reason_id" validate:"required"`
}

type UpdateReturnReasonDistributorRequest struct {
	CustId                      string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy                   int64  `json:"updated_by" validate:"required"`
	ReturnReasonDistributorCode string `json:"return_reason_code,omitempty" validate:"max=10,omitempty,alphanumericSpace"`
	ReturnReasonDistributorName string `json:"return_reason_name,omitempty" validate:"max=150,omitempty"`
	ReturnReasonDistributorType string `json:"return_reason_type" validate:"required,oneof='A' 'D' 'P' 'S'"`
	IsActive                    *bool  `json:"is_active,omitempty"`
}
