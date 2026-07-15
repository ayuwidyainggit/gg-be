package model

import (
	"time"
)

type ProductLine struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	PlId          int        `db:"pl_id" json:"pl_id"`
	PlCode        string     `db:"pl_code" json:"pl_code"`
	PlName        string     `db:"pl_name" json:"pl_name"`
	EffCall       int        `db:"eff_call" json:"eff_call"`
	MinItem       int        `db:"min_item" json:"min_item"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type ProductLineUpdate struct {
	PLCode    *string    `json:"pl_code,omitempty" sql:"pl_code"`
	PLName    *string    `json:"pl_name,omitempty" sql:"pl_name"`
	EffCall   *float32   `json:"eff_call" db:"eff_call"`
	MinItem   *float32   `json:"min_item" db:"min_item"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
