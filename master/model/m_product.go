package model

import (
	"time"
)

type MProduct struct {
	CustId    string     `db:"cust_id"`
	ProId     int        `db:"pro_id"`
	ProCode   string     `db:"pro_code"`
	ProName   string     `db:"pro_name"`
	IsActive  bool       `db:"is_active"`
	IsDel     bool       `db:"is_del"`
	CreatedBy *int64     `db:"created_by,omitempty"`
	CreatedAt *time.Time `db:"created_at,omitempty"`
	UpdatedBy *int64     `db:"updated_by,omitempty"`
	UpdatedAt *time.Time `db:"updated_at,omitempty"`
	DeletedBy *int64     `db:"deleted_by,omitempty"`
	DeletedAt *time.Time `db:"deleted_at,omitempty"`
}

type MProductUpdate struct {
	MProCode  *string    `json:"pro_code,omitempty" sql:"pro_code"`
	MProName  *string    `json:"pro_name,omitempty" sql:"pro_name"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
