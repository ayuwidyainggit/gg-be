package entity

import (
	"time"
)

type PackTypeResponse struct {
	PTypeId   int        `json:"ptype_id"`
	PTypeCode string     `json:"ptype_code"`
	PTypeName string     `json:"ptype_name"`
	IsActive  bool       `json:"is_active"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type PackTypeListResponse struct {
	PtypeId       int        `json:"ptype_id"`
	PtypeCode     string     `json:"ptype_code"`
	PtypeName     string     `json:"ptype_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type PackTypeLookupResponse struct {
	PtypeId   int    `json:"ptype_id"`
	PtypeCode string `json:"ptype_code"`
	PtypeName string `json:"ptype_name"`
}

type CreatePackTypeBody struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	CreatedBy int64  `json:"created_by" validate:"required"`
	PtypeCode string `json:"ptype_code" validate:"required,max=5,alphanumericSpace"`
	PtypeName string `json:"ptype_name" validate:"required,max=40,alphanumericSpace"`
	IsActive  bool   `json:"is_active"`
}

type DetailPackTypeParams struct {
	PtypeId int `params:"ptype_id" validate:"required"`
}

type UpdatePackTypeParams struct {
	PtypeId int `params:"ptype_id" validate:"required"`
}

type DeletePackTypeParams struct {
	PtypeId int `params:"ptype_id" validate:"required"`
}

type UpdatePackTypeRequest struct {
	CustId    string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy int64  `json:"updated_by" validate:"required"`
	PtypeCode string `json:"ptype_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	PtypeName string `json:"ptype_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	IsActive  *bool  `json:"is_active,omitempty"`
}
