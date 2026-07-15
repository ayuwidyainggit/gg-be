package model

import (
	"time"
)

type Brand struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	BrandId       int        `db:"brand_id" json:"brand_id"`
	BrandCode     string     `db:"brand_code" json:"brand_code"`
	BrandName     string     `db:"brand_name" json:"brand_name"`
	PlId          int        `db:"pl_id" json:"pl_id"`
	PlCode        string     `db:"pl_code" json:"pl_code"`
	PlName        string     `db:"pl_name" json:"pl_name"`
	EffCall       float32    `db:"eff_call" json:"eff_call"`
	MinItem       float32    `db:"min_item" json:"min_item"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type BrandUpdate struct {
	PlId      *int       `json:"pl_id,omitempty" sql:"pl_id"`
	BrandCode *string    `json:"brand_code,omitempty" sql:"brand_code"`
	BrandName *string    `json:"brand_name,omitempty" sql:"brand_name"`
	EffCall   *float32   `json:"eff_call,omitempty" sql:"eff_call"`
	MinItem   *float32   `json:"min_item,omitempty" sql:"min_item"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
