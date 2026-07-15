package model

import (
	"time"
)

type SubBrand1 struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	BrandId       int        `db:"brand_id" json:"brand_id"`
	SBrand1Id     int        `db:"sbrand1_id" json:"sbrand1_id"`
	SBrand1Code   string     `db:"sbrand1_code" json:"sbrand1_code"`
	SBrand1Name   string     `db:"sbrand1_name" json:"sbrand1_name"`
	EffCall       int        `db:"eff_call" json:"eff_call"`
	MinItem       int        `db:"min_item" json:"min_item"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by,omitempty"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by,omitempty"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	BrandCode     *string    `db:"brand_code" json:"brand_code"`
	BrandName     *string    `db:"brand_name" json:"brand_name"`
	PlId          *int       `db:"pl_id" json:"pl_id"`
	PlCode        *string    `db:"pl_code" json:"pl_code"`
	PlName        *string    `db:"pl_name" json:"pl_name"`
	MatGroupCode  *string    `db:"mat_group_code" json:"mat_group_code"`
	MatGroupName  *string    `db:"mat_group_name" json:"mat_group_name"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
}

type SubBrand1Update struct {
	BrandId     *int       `json:"brand_id,omitempty" sql:"brand_id"`
	SBrand1Code *string    `json:"sbrand1_code,omitempty" sql:"sbrand1_code"`
	SBrand1Name *string    `json:"sbrand1_name,omitempty" sql:"sbrand1_name"`
	EffCall     *int       `json:"eff_call" sql:"eff_call"`
	MinItem     *int       `json:"min_item" sql:"min_item"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}
