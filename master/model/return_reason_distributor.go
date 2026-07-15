package model

import (
	"time"
)

type ReturnReasonDistributor struct {
	CustId                      string     `db:"cust_id" json:"cust_id"`
	ReturnReasonDistributorId   int        `db:"return_reason_id" json:"return_reason_id"`
	ReturnReasonDistributorCode string     `db:"return_reason_code" json:"return_reason_code"`
	ReturnReasonDistributorName string     `db:"return_reason_name" json:"return_reason_name"`
	ReturnReasonDistributorType *string    `db:"return_reason_type" json:"return_reason_type"`
	IsActive                    bool       `db:"is_active" json:"is_active"`
	IsDel                       bool       `db:"is_del" json:"is_del"`
	CreatedBy                   *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt                   *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy                   *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName               *string    `db:"updated_by_name" json:"updated_by_name" `
	UpdatedAt                   *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy                   *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt                   *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type ReturnReasonDistributorUpdate struct {
	ReturnReasonDistributorCode *string    `json:"return_reason_code,omitempty" sql:"return_reason_code"`
	ReturnReasonDistributorName *string    `json:"return_reason_name,omitempty" sql:"return_reason_name"`
	ReturnReasonDistributorType *string    `json:"return_reason_type,omitempty" sql:"return_reason_type"`
	IsActive                    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt                   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy                   *int64     `json:"updated_by" sql:"updated_by"`
}
