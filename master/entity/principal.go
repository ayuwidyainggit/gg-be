package entity

import (
	"time"
)

type PrincipalResponse struct {
	PrincipalId   int        `json:"principal_id"`
	PrincipalCode string     `json:"principal_code"`
	PrincipalName string     `json:"principal_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type PrincipalListResponse struct {
	PrincipalId   int        `json:"principal_id"`
	PrincipalCode string     `json:"principal_code"`
	PrincipalName string     `json:"principal_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type PrincipalLookupResponse struct {
	PrincipalId   int    `json:"principal_id"`
	PrincipalCode string `json:"principal_code"`
	PrincipalName string `json:"principal_name"`
}

type CreatePrincipalBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	PrincipalCode string `json:"principal_code" validate:"required,max=5,alphanumericSpace"`
	PrincipalName string `json:"principal_name" validate:"required,max=50,alphanumericSpace"`
	IsActive      bool   `json:"is_active"`
}

type DetailPrincipalParams struct {
	PrincipalId int `params:"principal_id" validate:"required"`
}

type UpdatePrincipalParams struct {
	PrincipalId int `params:"principal_id" validate:"required"`
}

type DeletePrincipalParams struct {
	PrincipalId int `params:"principal_id" validate:"required"`
}

type UpdatePrincipalRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	PrincipalCode string `json:"principal_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	PrincipalName string `json:"principal_name,omitempty" validate:"max=50,omitempty,alphanumericSpace"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
