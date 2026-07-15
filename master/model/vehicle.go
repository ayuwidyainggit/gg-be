package model

import (
	"time"
)

type Vehicle struct {
	CustId          string     `db:"cust_id" json:"cust_id"`
	VehicleId       int64      `db:"vehicle_id" json:"vehicle_id"`
	VehicleNo       string     `db:"vehicle_no" json:"vehicle_no"`
	VehicleDesc     string     `db:"vehicle_desc" json:"vehicle_desc"`
	VehicleType     int        `db:"vehicle_type" json:"vehicle_type"`
	VehicleTypeName string     `db:"vehicle_type_name" json:"vehicle_type_name"`
	Length          float64    `db:"length" json:"length"`
	Width           float64    `db:"width" json:"width"`
	Height          float64    `db:"height" json:"height"`
	Weight          float64    `db:"weight" json:"weight"`
	Volume          float64    `db:"volume" json:"volume"`
	DriverId        int64      `db:"driver_id" json:"driver_id"`
	HelperId        int64      `db:"helper_id" json:"helper_id"`
	DriverName      string     `db:"driver_name" json:"driver_name"`
	HelperName      string     `db:"helper_name" json:"helper_name"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	IsDel           bool       `db:"is_del" json:"is_del"`
	CreatedBy       *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt       *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy       *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName   *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt       *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy       *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt       *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type VehicleUpdate struct {
	VehicleId   *string    `json:"vehicle_id,omitempty" sql:"vehicle_id"`
	VehicleNo   *string    `json:"vehicle_no,omitempty" sql:"vehicle_no"`
	VehicleDesc *string    `json:"vehicle_desc" sql:"vehicle_desc"`
	VehicleType *int       `json:"vehicle_type" sql:"vehicle_type"`
	Length      *float64   `json:"length" sql:"length"`
	Width       *float64   `json:"width" sql:"width"`
	Height      *float64   `json:"height" sql:"height"`
	Weight      *float64   `json:"weight" sql:"weight"`
	Volume      *float64   `json:"volume" sql:"volume"`
	DriverId    *int64     `json:"driver_id" sql:"driver_id"`
	HelperId    *int64     `json:"helper_id" sql:"helper_id"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
