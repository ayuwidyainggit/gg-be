package entity

import (
	"time"
)

type FlavorResponse struct {
	FlavorId   int        `json:"flavor_id"`
	FlavorCode string     `json:"flavor_code"`
	FlavorName string     `json:"flavor_name"`
	IsActive   bool       `json:"is_active"`
	UpdatedBy  *int64     `json:"updated_by"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
type FlavorListResponse struct {
	FlavorId      int        `json:"flavor_id"`
	FlavorCode    string     `json:"flavor_code"`
	FlavorName    string     `json:"flavor_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type FlavorLookupResponse struct {
	FlavorId   int    `json:"flavor_id"`
	FlavorCode string `json:"flavor_code"`
	FlavorName string `json:"flavor_name"`
}

type CreateFlavorBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	FlavorCode string `json:"flavor_code" validate:"required,max=5,alphanumericSpace"`
	FlavorName string `json:"flavor_name" validate:"required,max=40,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailFlavorParams struct {
	FlavorId int `params:"flavor_id" validate:"required"`
}

type UpdateFlavorParams struct {
	FlavorId int `params:"flavor_id" validate:"required"`
}

type DeleteFlavorParams struct {
	FlavorId int `params:"flavor_id" validate:"required"`
}

type UpdateFlavorRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	FlavorCode string `json:"flavor_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	FlavorName string `json:"flavor_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}
