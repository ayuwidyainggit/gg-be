package entity

import (
	"time"
)

type VehicleResponse struct {
	VehicleId       int64      `json:"vehicle_id"`
	VehicleNo       string     `json:"vehicle_no"`
	VehicleDesc     string     `json:"vehicle_desc"`
	VehicleType     int        `json:"vehicle_type"`
	VehicleTypeName string     `json:"vehicle_type_name"`
	Length          float64    `json:"length"`
	Width           float64    `json:"width"`
	Height          float64    `json:"height"`
	Weight          float64    `json:"weight"`
	Volume          float64    `json:"volume"`
	DriverId        int64      `json:"driver_id"`
	HelperId        int64      `json:"helper_id"`
	DriverName      string     `json:"driver_name"`
	HelperName      string     `json:"helper_name"`
	IsActive        bool       `json:"is_active"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedByName   *string    `json:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type CreateVehicleBody struct {
	CustId      string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64   `json:"created_by" validate:"required"`
	VehicleNo   string  `json:"vehicle_no" validate:"required,max=12"`
	VehicleDesc string  `json:"vehicle_desc" validate:"required,max=150"`
	VehicleType int     `json:"vehicle_type" validate:"required"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Weight      float64 `json:"weight"`
	Volume      float64 `json:"volume"`
	DriverId    int64   `json:"driver_id"`
	HelperId    int64   `json:"helper_id"`
	IsActive    bool    `json:"is_active"`
}

type DetailVehicleParams struct {
	VehicleId int64 `params:"vehicle_id" validate:"required"`
}

type UpdateVehicleParams struct {
	VehicleId int64 `params:"vehicle_id" validate:"required"`
}

type DeleteVehicleParams struct {
	VehicleId int64 `params:"vehicle_id" validate:"required"`
}

type UpdateVehicleRequest struct {
	CustId      string  `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64   `json:"updated_by" validate:"required"`
	VehicleNo   string  `json:"vehicle_no,omitempty" validate:"max=12,omitempty"`
	VehicleDesc string  `json:"vehicle_desc,omitempty" validate:"max=150,omitempty"`
	VehicleType int     `json:"vehicle_type,omitempty" validate:"omitempty"`
	Length      float64 `json:"length"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Weight      float64 `json:"weight"`
	Volume      float64 `json:"volume"`
	DriverId    int64   `json:"driver_id"`
	HelperId    int64   `json:"helper_id"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
