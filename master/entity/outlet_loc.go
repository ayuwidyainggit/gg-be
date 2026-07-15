package entity

import (
	"time"
)

type OutletLocResponse struct {
	OtLocId       int        `json:"ot_loc_id"`
	OtLocCode     string     `json:"ot_loc_code"`
	OtLocName     string     `json:"ot_loc_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type OutletLocLookupResponse struct {
	OtLocId   int    `json:"ot_loc_id"`
	OtLocCode string `json:"ot_loc_code"`
	OtLocName string `json:"ot_loc_name"`
}

type CreateOutletLocBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	OtLocCode string `json:"ot_loc_code" validate:"required,max=5,alphanumericSpace"`
	OtLocName string `json:"ot_loc_name" validate:"required,max=40,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailOutletLocParams struct {
	OtLocId int `params:"ot_loc_id" validate:"required"`
}

type UpdateOutletLocParams struct {
	OtLocId int `params:"ot_loc_id" validate:"required"`
}

type DeleteOutletLocParams struct {
	OtLocId int `params:"ot_loc_id" validate:"required"`
}

type UpdateOutletLocRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	OtLocCode string `json:"ot_loc_code,omitempty" validate:"required,max=5,omitempty,alphanumericSpace"`
	OtLocName string `json:"ot_loc_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
