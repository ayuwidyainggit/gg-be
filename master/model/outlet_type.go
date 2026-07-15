package model

import (
	"time"
)

type OutletType struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	OtTypeId      int        `db:"ot_type_id" json:"ot_type_id"`
	OtTypeCode    string     `db:"ot_type_code" json:"ot_type_code"`
	OtTypeName    string     `db:"ot_type_name" json:"ot_type_name"`
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

type OutletTypeUpdate struct {
	OtTypeCode *string    `json:"ot_type_code,omitempty" sql:"ot_type_code"`
	OtTypeName *string    `json:"ot_type_name,omitempty" sql:"ot_type_name"`
	IsActive   *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt  *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by" sql:"updated_by"`
}
