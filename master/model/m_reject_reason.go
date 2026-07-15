package model

import (
	"time"
)

type RejectReason struct {
	CustId           string `db:"cust_id" json:"cust_id"`
	RejectReasonId   int    `db:"reject_reason_id" json:"reject_reason_id"`
	RejectReasonCode string `db:"reject_reason_code" json:"reject_reason_code"`
	RejectReasonName string `db:"reject_reason_name" json:"reject_reason_name"`

	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int       `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type RejectReasonUpdate struct {
	// CustId string `db:"cust_id" json:"cust_id"`
	// RejectReasonId int   `db:"reject_reason_id" json:"reject_reason_id"`
	RejectReasonCode *string `db:"reject_reason_code" json:"reject_reason_code"`
	RejectReasonName *string `json:"reject_reason_name" sql:"reject_reason_name"`

	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
