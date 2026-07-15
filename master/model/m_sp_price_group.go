package model

import (
	"time"
)

type SpecialPriceGroup struct {
	CustId                string `db:"cust_id" json:"cust_id"`
	SpecialPriceGroupId   int8   `db:"sp_price_grp_id" json:"sp_price_grp_id"`
	SpecialPriceGroupCode string `db:"sp_price_grp_code" json:"sp_price_grp_code"`
	SpecialPriceGroupName string `db:"sp_price_grp_name" json:"sp_price_grp_name"`

	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int8      `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type SpecialPriceGroupUpdate struct {
	// CustId string `db:"cust_id" json:"cust_id"`
	// SpecialPriceGroupId int8   `db:"sp_price_grp_id" json:"sp_price_grp_id"`
	SpecialPriceGroupCode *string `sql:"sp_price_grp_code" json:"sp_price_grp_code"`
	SpecialPriceGroupName *string `sql:"sp_price_grp_name" json:"sp_price_grp_name"`

	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
