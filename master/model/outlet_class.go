package model

import (
	"time"
)

type OutletClass struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	OtClassId     int        `db:"ot_class_id" json:"ot_class_id"`
	OtClassCode   string     `db:"ot_class_code" json:"ot_class_code"`
	OtClassName   string     `db:"ot_class_name" json:"ot_class_name"`
	OtClassLimit  float64    `db:"ot_class_limit" json:"ot_class_limit"`
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

type OutletClassUpdate struct {
	OtClassCode  *string    `json:"ot_class_code,omitempty" sql:"ot_class_code"`
	OtClassName  *string    `json:"ot_class_name,omitempty" sql:"ot_class_name"`
	OtClassLimit *float64   `json:"ot_class_limit,omitempty" sql:"ot_class_limit"`
	IsActive     *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt    *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy    *int64     `json:"updated_by" sql:"updated_by"`
}
