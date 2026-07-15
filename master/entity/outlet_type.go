package entity

import (
	"time"
)

type OutletTypeResponse struct {
	OtTypeId      int        `json:"ot_type_id"`
	OtTypeCode    string     `json:"ot_type_code"`
	OtTypeName    string     `json:"ot_type_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type OutletTypeLookupResponse struct {
	OtTypeId   int    `json:"ot_type_id"`
	OtTypeCode string `json:"ot_type_code"`
	OtTypeName string `json:"ot_type_name"`
}

type CreateOutletTypeBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	OtTypeCode string `json:"ot_type_code" validate:"required,max=3,numeric"`
	OtTypeName string `json:"ot_type_name" validate:"required,max=40,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailOutletTypeParams struct {
	OtTypeId int64 `params:"ot_type_id" validate:"required"`
}

type UpdateOutletTypeParams struct {
	OtTypeId int `params:"ot_type_id" validate:"required"`
}

type DeleteOutletTypeParams struct {
	OtTypeId int `params:"ot_type_id" validate:"required"`
}

type UpdateOutletTypeRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	OtTypeCode string `json:"ot_type_code,omitempty" validate:"required,max=3,omitempty,numeric"`
	OtTypeName string `json:"ot_type_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
