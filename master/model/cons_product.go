package model

import (
	"time"
)

type ConsProduct struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	CProId        int        `db:"c_pro_id" json:"c_pro_id"`
	CProCode      string     `db:"c_pro_code" json:"c_pro_code"`
	CProName      string     `db:"c_pro_name" json:"c_pro_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name" `
}

type ConsProductUpdate struct {
	CProCode  *string    `json:"c_pro_code,omitempty" sql:"c_pro_code"`
	CProName  *string    `json:"c_pro_name,omitempty" sql:"c_pro_name"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
