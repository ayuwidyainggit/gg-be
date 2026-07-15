package entity

import (
	"time"
)

type ReturnReasonQueryFilter struct {
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	ReturnReasonType string `query:"return_reason_type"`
	IsActive         *int   `query:"is_active"`
}

type ReturnReasonResponse struct {
	ReturnReasonId   int        `json:"return_reason_id"`
	ReturnReasonCode string     `json:"return_reason_code"`
	ReturnReasonName string     `json:"return_reason_name"`
	ReturnReasonType *string    `json:"return_reason_type"`
	IsActive         bool       `json:"is_active"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedByName    string     `json:"updated_by_name"`
}

type ReturnReasonLookupResponse struct {
	ReturnReasonId   int    `json:"return_reason_id"`
	ReturnReasonCode string `json:"return_reason_code"`
	ReturnReasonName string `json:"return_reason_name"`
	ReturnReasonType string `json:"return_reason_type"`
}

type CreateReturnReasonBody struct {
	CustId           string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy        int64   `json:"created_by" validate:"required"`
	ReturnReasonCode string  `json:"return_reason_code" validate:"required,max=10,alphanumericSpace"`
	ReturnReasonName string  `json:"return_reason_name" validate:"required,max=150"`
	ReturnReasonType *string `json:"return_reason_type" validate:"required,oneof='A' 'O' 'P' 'S'"`
	IsActive         bool    `json:"is_active"`
}

type DetailReturnReasonParams struct {
	ReturnReasonId int `params:"return_reason_id" validate:"required"`
}

type UpdateReturnReasonParams struct {
	ReturnReasonId int `params:"return_reason_id" validate:"required"`
}

type DeleteReturnReasonParams struct {
	ReturnReasonId int `params:"return_reason_id" validate:"required"`
}

type UpdateReturnReasonRequest struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy        int64  `json:"updated_by" validate:"required"`
	ReturnReasonCode string `json:"return_reason_code,omitempty" validate:"max=10,omitempty,alphanumericSpace"`
	ReturnReasonName string `json:"return_reason_name,omitempty" validate:"max=150,omitempty"`
	ReturnReasonType string `json:"return_reason_type" validate:"required,oneof='A' 'O' 'P' 'S'"`
	IsActive         *bool  `json:"is_active,omitempty"`
}
