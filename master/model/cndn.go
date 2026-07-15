package model

import (
	"time"
)

type Cndn struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	CndnId        int        `db:"cndn_id" json:"cndn_id"`
	CndnCode      string     `db:"cndn_code" json:"cndn_code"`
	CndnName      string     `db:"cndn_name" json:"cndn_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type CndnUpdate struct {
	CndnCode  *string    `json:"cndn_code,omitempty" sql:"cndn_code"`
	CndnName  *string    `json:"cndn_name,omitempty" sql:"cndn_name"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
