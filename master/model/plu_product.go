package model

import (
	"time"
)

type PluProduct struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	PluProId      int        `json:"plu_pro_id" db:"plu_pro_id"`
	PluGrpId      int        `json:"plu_grp_id" db:"plu_grp_id"`
	ProId         int        `json:"pro_id" db:"pro_id"`
	ProCode       string     `json:"pro_code" db:"pro_code"`
	ProName       string     `json:"pro_name" db:"pro_name"`
	PluNo         string     `json:"plu_no" db:"plu_no"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	CreatedBy     *int64     `json:"created_by" db:"created_by,omitempty"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at,omitempty"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at,omitempty"`
}

type PluProductUpdate struct {
	PluGrpId  *int       `json:"plu_grp_id,omitempty" sql:"plu_grp_id"`
	ProId     *int       `json:"pro_id,omitempty" sql:"pro_id"`
	PluNo     *string    `json:"plu_no,omitempty" sql:"plu_no"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
