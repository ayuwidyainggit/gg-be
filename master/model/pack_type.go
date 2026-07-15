package model

import (
	"time"
)

type PackType struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	PtypeId       int        `db:"ptype_id" json:"ptype_id"`
	PtypeCode     string     `db:"ptype_code" json:"ptype_code"`
	PtypeName     string     `db:"ptype_name" json:"ptype_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
}

type PackTypeUpdate struct {
	PtypeCode *string    `json:"ptype_code,omitempty" sql:"ptype_code"`
	PtypeName *string    `json:"ptype_name,omitempty" sql:"ptype_name"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
