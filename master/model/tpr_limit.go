package model

import (
	"time"
)

type MTprLimit struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	TprLimitId    int        `db:"tpr_limit_id" json:"tpr_limit_id"`
	ProId         int        `db:"pro_id" json:"pro_id"`
	TprType       *int64     `db:"tpr_type" json:"tpr_type"`
	DateStart     *string    `db:"date_start" json:"date_start"`
	DateEnd       *string    `db:"date_end" json:"date_end"`
	ValueLimit    *float64   `db:"value_limit" json:"value_limit"`
	ValueUsed     *float64   `db:"value_used" json:"value_used"`
	ValueUsedStr  *string    `db:"value_used_str" json:"value_used_str"`
	VatType       *int64     `db:"vat_type" json:"vat_type"`
	CreatedBy     *int64     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

type MTprLimitUpdate struct {
	ProId        *int       `sql:"pro_id" json:"pro_id"`
	TprType      *int64     `sql:"tpr_type" json:"tpr_type"`
	DateStart    *string    `sql:"date_start" json:"date_start"`
	DateEnd      *string    `sql:"date_end" json:"date_end"`
	ValueLimit   *float64   `sql:"value_limit" json:"value_limit"`
	ValueUsed    *float64   `sql:"value_used" json:"value_used"`
	ValueUsedStr *string    `sql:"value_used_str" json:"value_used_str"`
	VatType      *int64     `sql:"vat_type" json:"vat_type"`
	UpdatedBy    *int64     `sql:"updated_by" json:"updated_by"`
	UpdatedAt    *time.Time `sql:"updated_at" json:"updated_at"`
}

type MTprLimitList struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	TprLimitId    int        `db:"tpr_limit_id" json:"tpr_limit_id"`
	ProId         int        `db:"pro_id" json:"pro_id"`
	ProCode       *string    `db:"pro_code" json:"pro_code"`
	ProName       *string    `db:"pro_name" json:"pro_name"`
	TprType       *int64     `db:"tpr_type" json:"tpr_type"`
	DateStart     *time.Time `db:"date_start" json:"date_start"`
	DateEnd       *time.Time `db:"date_end" json:"date_end"`
	ValueLimit    *float64   `db:"value_limit" json:"value_limit"`
	ValueUsed     *float64   `db:"value_used" json:"value_used"`
	ValueUsedStr  *string    `db:"value_used_str" json:"value_used_str"`
	VatType       *int64     `db:"vat_type" json:"vat_type"`
	CreatedBy     *int64     `db:"created_by" json:"created_by"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy     *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	IsDel         bool       `db:"is_del" json:"is_del"`
	DeletedBy     *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}
