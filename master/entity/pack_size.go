package entity

import (
	"time"
)

type PackSizeResponse struct {
	PSizeId   int        `json:"psize_id"`
	PSizeCode string     `json:"psize_code"`
	PSizeName string     `json:"psize_name"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type PackSizeListResponse struct {
	PsizeId       int        `json:"psize_id"`
	PsizeCode     string     `json:"psize_code"`
	PsizeName     string     `json:"psize_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type PackSizeLookupResponse struct {
	PsizeId   int    `json:"psize_id"`
	PsizeCode string `json:"psize_code"`
	PsizeName string `json:"psize_name"`
}

type CreatePackSizeBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	PsizeCode string `json:"psize_code" validate:"required,max=5,alphanumericSpace"`
	PsizeName string `json:"psize_name" validate:"required,max=30"`
	IsActive  bool   `json:"is_active"`
}

type DetailPackSizeParams struct {
	PSizeId int `params:"psize_id" validate:"required"`
}

type UpdatePackSizeParams struct {
	PSizeId int `params:"psize_id" validate:"required"`
}

type DeletePackSizeParams struct {
	PSizeId int `params:"psize_id" validate:"required"`
}

type UpdatePackSizeRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	PsizeCode string `json:"psize_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	PsizeName string `json:"psize_name,omitempty" validate:"max=30,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
