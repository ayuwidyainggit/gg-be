package model

import (
	"time"
)

type Principal struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	PrincipalId   int        `db:"principal_id" json:"principal_id"`
	PrincipalCode string     `db:"principal_code" json:"principal_code"`
	PrincipalName string     `db:"principal_name" json:"principal_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
}

type PrincipalUpdate struct {
	PrincipalCode *string    `json:"principal_code,omitempty" sql:"principal_code"`
	PrincipalName *string    `json:"principal_name,omitempty" sql:"principal_name"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
