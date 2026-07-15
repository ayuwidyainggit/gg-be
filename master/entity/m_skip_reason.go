package entity

import (
	"time"
)

type SkipReasonQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	SkipReasonId int    `query:"skip_reason_id"`
}

type SkipReasonResponse struct {
	SkipReasonId   int        `json:"skip_reason_id"`
	SkipReasonCode string     `json:"skip_reason_code"`
	SkipReasonName string     `json:"skip_reason_name"`
	IsActive       bool       `json:"is_active"`
	UpdatedBy      *int64     `json:"updated_by"`
	UpdatedByName  string     `json:"updated_by_name"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
type SkipReasonListResponse struct {
	SkipReasonId   int        `json:"skip_reason_id"`
	SkipReasonCode string     `json:"skip_reason_code"`
	SkipReasonName string     `json:"skip_reason_name"`
	IsActive       bool       `json:"is_active"`
	UpdatedBy      *int64     `json:"updated_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedByName  string     `json:"updated_by_name"`
}

type SkipReasonLookupResponse struct {
	SkipReasonId   int    `json:"skip_reason_id"`
	SkipReasonName string `json:"skip_reason_name"`
}

type CreateSkipReasonBody struct {
	CustId         string `json:"cust_id" validate:"required,max=10"`
	CreatedBy      int64  `json:"created_by" validate:"required"`
	SkipReasonCode string `json:"skip_reason_code" validate:"required,max=5"`
	SkipReasonName string `json:"skip_reason_name" validate:"required,max=50,alphanumericSpace"`
	IsActive       bool   `json:"is_active"`
}

type UpdateSkipReasonRequest struct {
	CustId         string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy      int64  `json:"updated_by" validate:"required"`
	SkipReasonCode string `json:"skip_reason_code" validate:"required,max=5"`
	SkipReasonName string `json:"skip_reason_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	IsActive       *bool  `json:"is_active,omitempty"`
}

type DetailSkipReasonParams struct {
	SkipReasonId int `params:"skip_reason_id" validate:"required"`
}

type UpdateSkipReasonParams struct {
	SkipReasonId int `params:"skip_reason_id" validate:"required"`
}

type DeleteSkipReasonParams struct {
	SkipReasonId int `params:"skip_reason_id" validate:"required"`
}
