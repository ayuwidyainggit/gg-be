package model

import (
	"time"
)

type OutletLoc struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	OtLocId       int        `db:"ot_loc_id" json:"ot_loc_id"`
	OtLocCode     string     `db:"ot_loc_code" json:"ot_loc_code"`
	OtLocName     string     `db:"ot_loc_name" json:"ot_loc_name"`
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

type OutletLocUpdate struct {
	OutletLocCode *string    `json:"ot_loc_code,omitempty" sql:"ot_loc_code"`
	OutletLocName *string    `json:"ot_loc_name,omitempty" sql:"ot_loc_name"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
