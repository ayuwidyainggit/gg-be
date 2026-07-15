package model

import (
	"time"
)

type Unit struct {
	CustId          string     `db:"cust_id"`
	UnitId          string     `db:"unit_id"`
	UnitName        string     `db:"unit_name"`
	UnitIdCoreTax   *string    `db:"unit_id_coretax"`
	UnitNameCoreTax *string    `db:"unit_name_coretax"`
	IsActive        bool       `db:"is_active"`
	IsDel           bool       `db:"is_del"`
	CreatedBy       *int64     `db:"created_by,omitempty"`
	CreatedAt       *time.Time `db:"created_at,omitempty"`
	UpdatedBy       *int64     `db:"updated_by,omitempty"`
	UpdatedAt       *time.Time `db:"updated_at,omitempty"`
	UpdatedByName   *string    `json:"updated_by_name" db:"updated_by_name"`
	DeletedBy       *int64     `db:"deleted_by,omitempty"`
	DeletedAt       *time.Time `db:"deleted_at,omitempty"`
}

type UnitUpdate struct {
	UnitId        *string    `json:"unit_id,omitempty" sql:"unit_id"`
	UnitName      *string    `json:"unit_name,omitempty" sql:"unit_name"`
	UnitIdCoreTax *string    `json:"unit_id_coretax,omitempty" sql:"unit_id_coretax"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
