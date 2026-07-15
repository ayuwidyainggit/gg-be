package model

import (
	"time"
)

type SkipReason struct {
	CustId         string `db:"cust_id" json:"cust_id"`
	SkipReasonId   int    `db:"skip_reason_id" json:"skip_reason_id"`
	SkipReasonCode string `db:"skip_reason_code" json:"skip_reason_code"`
	SkipReasonName string `db:"skip_reason_name" json:"skip_reason_name"`

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

type SkipReasonUpdate struct {
	// CustId string `db:"cust_id" json:"cust_id"`
	// SkipReasonId int8   `db:"skip_reason_id" json:"skip_reason_id"`
	SkipReasonCode *string `db:"skip_reason_code" json:"skip_reason_code"`
	SkipReasonName *string `json:"skip_reason_name" sql:"skip_reason_name"`

	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
