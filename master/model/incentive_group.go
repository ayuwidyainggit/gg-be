package model

import (
	"time"
)

type IncentiveGroup struct {
	CustID     string     `db:"cust_id" json:"cust_id"`
	IncGrpID   int        `db:"inc_grp_id" json:"inc_grp_id"`
	IncGrpCode string     `db:"inc_grp_code" json:"inc_grp_code"`
	IncGrpName string     `db:"inc_grp_name" json:"inc_grp_name"`
	IsActive   bool       `db:"is_active" json:"is_active"`
	CreatedBy  *int64     `db:"created_by" json:"created_by"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy  *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`
	IsDel      bool       `db:"is_del" json:"is_del"`
	DeletedBy  *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`
}

type IncentiveGroupList struct {
	CustID        string     `db:"cust_id" json:"cust_id"`
	IncGrpID      int        `db:"inc_grp_id" json:"inc_grp_id"`
	IncGrpCode    string     `db:"inc_grp_code" json:"inc_grp_code"`
	IncGrpName    string     `db:"inc_grp_name" json:"inc_grp_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
}

type IncentiveGroupUpdate struct {
	IncGrpCode *string    `json:"inc_grp_code,omitempty" sql:"inc_grp_code"`
	IncGrpName *string    `json:"inc_grp_name,omitempty" sql:"inc_grp_name"`
	IsActive   *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt  *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy  *int64     `json:"updated_by" sql:"updated_by"`
}
