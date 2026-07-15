package entity

import (
	"time"
)

type PickupReasonQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	IsActive *int   `query:"is_active"`
}

type PickupReasonResponse struct {
	PickupReasonId   int        `json:"pickup_reason_id"`
	PickupReasonCode string     `json:"pickup_reason_code"`
	PickupReasonName string     `json:"pickup_reason_name"`
	IsActive         bool       `json:"is_active"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedByName    string     `json:"updated_by_name"`
}

type PickupReasonLookupResponse struct {
	PickupReasonId   int    `json:"pickup_reason_id"`
	PickupReasonCode string `json:"pickup_reason_code"`
	PickupReasonName string `json:"pickup_reason_name"`
}

type CreatePickupReasonBody struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	CreatedBy        int64  `json:"created_by" validate:"required"`
	PickupReasonCode string `json:"pickup_reason_code" validate:"required,max=5,numeric"`
	PickupReasonName string `json:"pickup_reason_name" validate:"required,max=50,alphanumericSpace"`
	IsActive         bool   `json:"is_active"`
}

type DetailPickupReasonParams struct {
	PickupReasonId int `params:"pickup_reason_id" validate:"required"`
}

type UpdatePickupReasonParams struct {
	PickupReasonId int `params:"pickup_reason_id" validate:"required"`
}

type DeletePickupReasonParams struct {
	PickupReasonId int `params:"pickup_reason_id" validate:"required"`
}

type UpdatePickupReasonRequest struct {
	CustId           string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy        int64  `json:"updated_by" validate:"required"`
	PickupReasonCode string `json:"pickup_reason_code,omitempty" validate:"max=5,omitempty,numeric"`
	PickupReasonName string `json:"pickup_reason_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	IsActive         *bool  `json:"is_active,omitempty"`
}
