package entity

import (
	"time"
)

type RejectReasonQueryFilter struct {
	CustId         string
	ParentCustId   string
	Page           int    `query:"page"`
	Limit          int    `query:"limit" validate:"required"`
	Query          string `query:"q"`
	Mode           string `query:"mode"`
	Sort           string `query:"sort"`
	IsActive       *int   `query:"is_active"`
	RejectReasonId int    `query:"reject_reason_id"`
}

type RejectReasonResponse struct {
	RejectReasonId   int        `json:"reject_reason_id"`
	RejectReasonCode string     `json:"reject_reason_code"`
	RejectReasonName string     `json:"reject_reason_name"`
	IsActive         bool       `json:"is_active"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedByName    string     `json:"updated_by_name"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
type RejectReasonListResponse struct {
	RejectReasonId   int        `json:"reject_reason_id"`
	RejectReasonCode string     `json:"reject_reason_code"`
	RejectReasonName string     `json:"reject_reason_name"`
	IsActive         bool       `json:"is_active"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedByName    string     `json:"updated_by_name"`
}

type RejectReasonLookupResponse struct {
	RejectReasonId   int    `json:"reject_reason_id"`
	RejectReasonName string `json:"reject_reason_name"`
}

type CreateRejectReasonBody struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	CreatedBy        int64  `json:"created_by" validate:"required"`
	RejectReasonCode string `json:"reject_reason_code" validate:"required,max=5"`
	RejectReasonName string `json:"reject_reason_name" validate:"required,max=50,alphanumericSpace"`
	IsActive         bool   `json:"is_active"`
}

type UpdateRejectReasonRequest struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy        int64  `json:"updated_by" validate:"required"`
	RejectReasonCode string `json:"reject_reason_code" validate:"required,max=5"`
	RejectReasonName string `json:"reject_reason_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	IsActive         *bool  `json:"is_active,omitempty"`
}

type DetailRejectReasonParams struct {
	RejectReasonId int `params:"reject_reason_id" validate:"required"`
}

type UpdateRejectReasonParams struct {
	RejectReasonId int `params:"reject_reason_id" validate:"required"`
}

type DeleteRejectReasonParams struct {
	RejectReasonId int `params:"reject_reason_id" validate:"required"`
}
