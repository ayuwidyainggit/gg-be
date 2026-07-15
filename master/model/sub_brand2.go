package model

import (
	"time"
)

type SubBrand2 struct {
	CustId        string     `db:"cust_id"`
	SBrand2Id     int        `db:"sbrand2_id"`
	SBrand2Code   string     `db:"sbrand2_code"`
	SBrand2Name   string     `db:"sbrand2_name"`
	IsActive      bool       `db:"is_active"`
	IsDel         bool       `db:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty"`
	CreatedAt     *time.Time `db:"created_at,omitempty"`
	UpdatedBy     *int64     `db:"updated_by,omitempty"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty"`
	DeletedBy     *int64     `db:"deleted_by,omitempty"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
}

type SubBrand2Update struct {
	SBrand2Code *string    `json:"sbrand2_code,omitempty" sql:"sbrand2_code"`
	SBrand2Name *string    `json:"sbrand2_name,omitempty" sql:"sbrand2_name"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
