package model

import (
	"time"
)

type Disc struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	DiscId        int64      `db:"disc_id" json:"disc_id"`
	DiscCode      string     `db:"disc_code" json:"disc_code"`
	DiscName      string     `db:"disc_name" json:"disc_name"`
	StartDate     *time.Time `db:"start_date" json:"start_date"`
	EndDate       *time.Time `db:"end_date" json:"end_date"`
	RangeType     *int       `db:"range_type" json:"range_type"`
	IsMultiple    bool       `db:"is_multiple" json:"is_multiple"`
	PurchaseLimit float64    `db:"purchase_limit" json:"purchase_limit"`
	DiscType      int        `db:"disc_type" json:"disc_type"`
	DiscPerc      float64    `db:"disc_perc" json:"disc_perc"`
	DiscValue     float64    `db:"disc_value" json:"disc_value"`
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

type DiscUpdate struct {
	DiscCode      *string    `json:"disc_code,omitempty" sql:"disc_code"`
	DiscName      *string    `json:"disc_name,omitempty" sql:"disc_name"`
	StartDate     *string    `json:"start_date,omitempty" sql:"start_date"`
	EndDate       *string    `json:"end_date,omitempty" sql:"end_date"`
	RangeType     *int       `json:"range_type,omitempty" sql:"range_type"`
	IsMultiple    *bool      `json:"is_multiple,omitempty" sql:"is_multiple"`
	PurchaseLimit *float64   `json:"purchase_limit,omitempty" sql:"purchase_limit"`
	DiscType      *int       `json:"disc_type,omitempty" sql:"disc_type"`
	DiscPerc      *float64   `json:"disc_perc,omitempty" sql:"disc_perc"`
	DiscValue     *float64   `json:"disc_value,omitempty" sql:"disc_value"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}
