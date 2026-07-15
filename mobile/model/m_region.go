package model

import "time"

type Region struct {
	RegionID   int        `db:"region_id" json:"region_id"`
	CustID     string     `db:"cust_id" json:"cust_id"`
	RegionCode string     `db:"region_code" json:"region_code"`
	RegionName string     `db:"region_name" json:"region_name"`
	IsActive   *bool      `db:"is_active" json:"is_active"`
	CreatedBy  *int64     `db:"created_by" json:"created_by"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy  *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
	IsDel      *bool      `db:"is_del" json:"is_del"`
	DeletedBy  *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}
