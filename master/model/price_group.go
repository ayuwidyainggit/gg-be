package model

import (
	"time"
)

type PriceGroup struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	PriceGrpId    int        `db:"price_grp_id" json:"price_grp_id"`
	PriceGrpCode  string     `db:"price_grp_code" json:"price_grp_code"`
	PriceGrpName  string     `db:"price_grp_name" json:"price_grp_name"`
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

type PriceGroupUpdate struct {
	PriceGroupCode *string    `json:"price_grp_code,omitempty" sql:"price_grp_code"`
	PriceGroupName *string    `json:"price_grp_name,omitempty" sql:"price_grp_name"`
	IsActive       *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt      *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy      *int64     `json:"updated_by" sql:"updated_by"`
}
