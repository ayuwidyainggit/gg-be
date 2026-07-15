package entity

import (
	"time"
)

type TakingOrderQueryFilter struct {
	CustId        string
	ParentCustId  string
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
	IsActive      *int   `query:"is_active"`
	TakingOrderId int    `query:"taking_order_id"`
}

type TakingOrderResponse struct {
	TakingOrderId   int        `json:"taking_order_id"`
	TakingOrderName string     `json:"taking_order_name"`
	ImageUrl        string     `json:"image_url"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedByName   string     `json:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
type TakingOrderListResponse struct {
	TakingOrderId   int        `json:"taking_order_id"`
	TakingOrderName string     `json:"taking_order_name"`
	ImageUrl        string     `json:"image_url"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
	UpdatedByName   string     `json:"updated_by_name"`
}

type TakingOrderLookupResponse struct {
	TakingOrderId   int    `json:"taking_order_id"`
	TakingOrderName string `json:"taking_order_name"`
}

type CreateTakingOrderBody struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	CreatedBy       int64  `json:"created_by" validate:"required"`
	TakingOrderName string `json:"taking_order_name" validate:"required,max=50,alphanumericSpace"`
	ImageUrl        string `json:"image_url"`
	IsActive        bool   `json:"is_active"`
}

type UpdateTakingOrderRequest struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy       int64  `json:"updated_by" validate:"required"`
	TakingOrderName string `json:"taking_order_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	ImageUrl        string `json:"image_url"`
	IsActive        *bool  `json:"is_active,omitempty"`
}

type DetailTakingOrderParams struct {
	TakingOrderId int `params:"taking_order_id" validate:"required"`
}

type UpdateTakingOrderParams struct {
	TakingOrderId int `params:"taking_order_id" validate:"required"`
}

type DeleteTakingOrderParams struct {
	TakingOrderId int `params:"taking_order_id" validate:"required"`
}
