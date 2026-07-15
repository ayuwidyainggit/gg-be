package entity

import (
	"time"
)

type OutletClassResponse struct {
	OtClassId     int        `json:"ot_class_id"`
	OtClassCode   string     `json:"ot_class_code"`
	OtClassName   string     `json:"ot_class_name"`
	OtClassLimit  float64    `json:"ot_class_limit"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type OutletClassLookupResponse struct {
	OtClassId    int     `json:"ot_class_id"`
	OtClassCode  string  `json:"ot_class_code"`
	OtClassName  string  `json:"ot_class_name"`
	OtClassLimit float64 `json:"ot_class_limit"`
}

type CreateOutletClassBody struct {
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy    int64   `json:"created_by" validate:"required"`
	OtClassCode  string  `json:"ot_class_code" validate:"required,max=5,alphanumericSpace"`
	OtClassName  string  `json:"ot_class_name" validate:"required,max=40,alphanumericSpace"`
	OtClassLimit float64 `json:"ot_class_limit" validate:""`
	IsActive     bool    `json:"is_active"`
}

type DetailOutletClassParams struct {
	OtClassId int `params:"ot_class_id" validate:"required"`
}

type UpdateOutletClassParams struct {
	OtClassId int `params:"ot_class_id" validate:"required"`
}

type DeleteOutletClassParams struct {
	OtClassId int `params:"ot_class_id" validate:"required"`
}

type UpdateOutletClassRequest struct {
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	UpdatedBy    int64   `json:"updated_by" validate:"required"`
	OtClassCode  string  `json:"ot_class_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	OtClassName  string  `json:"ot_class_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	OtClassLimit float64 `json:"ot_class_limit,omitempty" validate:"omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}
