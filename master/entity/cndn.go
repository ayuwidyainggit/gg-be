package entity

import (
	"time"
)

type CndnResponse struct {
	CndnId        int        `json:"cndn_id"`
	CndnCode      string     `json:"cndn_code"`
	CndnName      string     `json:"cndn_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type CndnLookupResponse struct {
	CndnId   int    `json:"cndn_id"`
	CndnCode string `json:"cndn_code"`
	CndnName string `json:"cndn_name"`
}

type CreateCndnBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	CndnCode  string `json:"cndn_code" validate:"required,max=10,alphanumericSpace"`
	CndnName  string `json:"cndn_name" validate:"required,max=150"`
	IsActive  bool   `json:"is_active"`
}

type DetailCndnParams struct {
	CndnId int `params:"cndn_id" validate:"required"`
}

type UpdateCndnParams struct {
	CndnId int `params:"cndn_id" validate:"required"`
}

type DeleteCndnParams struct {
	CndnId int `params:"cndn_id" validate:"required"`
}

type UpdateCndnRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	CndnCode  string `json:"cndn_code,omitempty" validate:"required,max=10,alphanumericSpace"`
	CndnName  string `json:"cndn_name,omitempty" validate:"max=150,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
