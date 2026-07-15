package model

import (
	"time"
)

type MissedPaymentReasons struct {
	CustId                   string     `db:"cust_id" json:"cust_id"`
	MissedPaymentReasonsId   int        `db:"missed_payment_reasons_id" json:"missed_payment_reasons_id"`
	MissedPaymentReasonsCode string     `db:"missed_payment_reasons_code" json:"missed_payment_reasons_code"`
	MissedPaymentReasonsName string     `db:"missed_payment_reasons_name" json:"missed_payment_reasons_name"`
	ImageUrl                 string     `db:"image_url" json:"image_url"`
	IsActive                 bool       `db:"is_active" json:"is_active"`
	CreatedBy                *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt                *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy                *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName            *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt                *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel                    bool       `db:"is_del" json:"is_del"`
	DeletedBy                *int8      `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt                *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type MissedPaymentReasonsUpdate struct {
	MissedPaymentReasonsCode *string    `json:"missed_payment_reasons_code" sql:"missed_payment_reasons_code"`
	MissedPaymentReasonsName *string    `json:"missed_payment_reasons_name" sql:"missed_payment_reasons_name"`
	ImageUrl                 *string    `db:"image_url" json:"image_url"`
	IsActive                 *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt                *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy                *int64     `json:"updated_by" sql:"updated_by"`
}
