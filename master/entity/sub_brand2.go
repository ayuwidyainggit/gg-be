package entity

import (
	"time"
)

type SubBrand2Response struct {
	SBrand2Id   int        `json:"sbrand2_id"`
	SBrand2Code string     `json:"sbrand2_code"`
	SBrand2Name string     `json:"sbrand2_name"`
	IsActive    bool       `json:"is_active"`
	UpdatedBy   *int64     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type SubBrand2ListResponse struct {
	SBrand2Id     int        `json:"sbrand2_id"`
	SBrand2Code   string     `json:"sbrand2_code"`
	SBrand2Name   string     `json:"sbrand2_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type CreateSubBrand2Body struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64  `json:"created_by" validate:"required"`
	SBrand2Code string `json:"sbrand2_code" validate:"required,max=5,alphanumericSpace"`
	SBrand2Name string `json:"sbrand2_name" validate:"required,max=40,alphanumericSpace"`
	IsActive    bool   `json:"is_active"`
}

type DetailSubBrand2Params struct {
	SBrand2Id int `params:"sbrand2_id" validate:"required"`
}

type UpdateSubBrand2Params struct {
	SBrand2Id int `params:"sbrand2_id" validate:"required"`
}

type DeleteSubBrand2Params struct {
	SBrand2Id int `params:"sbrand2_id" validate:"required"`
}

type UpdateSubBrand2Request struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64  `json:"updated_by" validate:"required"`
	SBrand2Code string `json:"sbrand2_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	SBrand2Name string `json:"sbrand2_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
