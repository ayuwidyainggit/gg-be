package model

import (
	"time"
)

type VehicleType struct {
	CustId        string `db:"cust_id" json:"cust_id"`
	VehicleTypeId int    `db:"vehicle_type_id" json:"vehicle_type_id"`
	// VehicleTypeCode string `db:"vehicle_type_code" json:"vehicle_type_code"`
	VehicleTypeName string `db:"vehicle_type_name" json:"vehicle_type_name"`

	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int8      `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type VehicleTypeUpdate struct {
	// CustId string `db:"cust_id" json:"cust_id"`
	// VehicleTypeId int8   `db:"vehicle_type_id" json:"vehicle_type_id"`
	// VehicleTypeCode string `db:"vehicle_type_code" json:"vehicle_type_code"`
	VehicleTypeName *string `json:"vehicle_type_name" sql:"vehicle_type_name"`

	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
