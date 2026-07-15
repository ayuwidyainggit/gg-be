package model

import "time"

type Area struct {
	AreaID     int        `db:"area_id" json:"area_id"`
	CustID     string     `db:"cust_id" json:"cust_id"`
	AreaCode   string     `db:"area_code" json:"area_code"`
	AreaName   string     `db:"area_name" json:"area_name"`
	RegionID   int        `db:"region_id" json:"region_id"`
	OfficialID *int       `db:"official_id" json:"official_id"`
	IsActive   *bool      `db:"is_active" json:"is_active"`
	CreatedBy  *int64     `db:"created_by" json:"created_by"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy  *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
	IsDel      *bool      `db:"is_del" json:"is_del"`
	DeletedBy  *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}
