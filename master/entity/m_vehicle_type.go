package entity

import (
	"time"
)

type VehicletypeQueryFilter struct {
	CustId        string
	ParentCustId  string
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
	IsActive      *int   `query:"is_active"`
	VehicleTypeId int    `query:"vehicle_type_id"`
}

type VehicletypeResponse struct {
	VehicleTypeId   int        `json:"vehicle_type_id"`
	VehicleTypeName string     `json:"vehicle_type_name"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedByName   string     `json:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
type VehicletypeListResponse struct {
	VehicleTypeId   int        `json:"vehicle_type_id"`
	VehicleTypeName string     `json:"vehicle_type_name"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
	UpdatedByName   string     `json:"updated_by_name"`
}

type VehicletypeLookupResponse struct {
	VehicleTypeId   int    `json:"vehicle_type_id"`
	VehicleTypeName string `json:"vehicle_type_name"`
}

type CreateVehicletypeBody struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	CreatedBy       int64  `json:"created_by" validate:"required"`
	VehicleTypeName string `json:"vehicle_type_name" validate:"required,max=25,alphanumericSpace"`
	IsActive        bool   `json:"is_active"`
}

type UpdateVehicletypeRequest struct {
	CustId          string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy       int64  `json:"updated_by" validate:"required"`
	VehicleTypeName string `json:"vehicle_type_name,omitempty" validate:"max=25,omitempty,alphanumericSpace"`
	IsActive        *bool  `json:"is_active,omitempty"`
}

type DetailVehicletypeParams struct {
	VehicleTypeId int `params:"vehicle_type_id" validate:"required"`
}

type UpdateVehicletypeParams struct {
	VehicleTypeId int `params:"vehicle_type_id" validate:"required"`
}

type DeleteVehicletypeParams struct {
	VehicleTypeId int `params:"vehicle_type_id" validate:"required"`
}
