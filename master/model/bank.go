package model

import (
	"time"
)

type Bank struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	BankId        int        `db:"bank_id" json:"bank_id"`
	BankCode      string     `db:"bank_code" json:"bank_code"`
	BankName      string     `db:"bank_name" json:"bank_name"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	CreatedBy     *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName *string    `db:"updated_by_name" json:"updated_by_name"`
	DeletedBy     *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
}

type BankUpdate struct {
	BankCode  *string    `json:"bank_code,omitempty" sql:"bank_code"`
	BankName  *string    `json:"bank_name,omitempty" sql:"bank_name"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}

type BankLookup struct {
	BankId   int    `db:"bank_id" json:"bank_id"`
	BankCode string `db:"bank_code" json:"bank_code"`
	BankName string `db:"bank_name" json:"bank_name"`
}
