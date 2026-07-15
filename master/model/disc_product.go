package model

import "time"

type DiscProduct struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	DiscId        int        `json:"disc_id" db:"disc_id"`
	ProId         int        `json:"pro_id,omitempty" db:"pro_id"`
	MinQty        *float64   `json:"min_qty" db:"min_qty"`
	MinQtyStr     *string    `json:"min_qty_str" db:"min_qty_str"`
	MaxQty        *float64   `json:"max_qty" db:"max_qty"`
	MaxQtyStr     *string    `json:"max_qty_str" db:"max_qty_str"`
	DiscPerc      *float64   `json:"disc_perc" db:"disc_perc"`
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

type DiscProductUpdate struct {
	DiscId    *int       `json:"disc_id" sql:"disc_id"`
	ProId     *int       `json:"pro_id" sql:"pro_id"`
	MinQty    *float64   `json:"min_qty" sql:"min_qty"`
	MinQtyStr *string    `json:"min_qty_str" sql:"min_qty_str"`
	MaxQty    *float64   `json:"max_qty" sql:"max_qty"`
	MaxQtyStr *string    `json:"max_qty_str" sql:"max_qty_str"`
	DiscPerc  *float64   `json:"disc_perc" sql:"disc_perc"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
}
type DiscProductRead struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	DiscId        int        `json:"disc_id" db:"disc_id"`
	DiscCode      *string    `json:"disc_code" db:"disc_code" `
	DiscName      *string    `json:"disc_name" db:"disc_name"`
	ProId         int        `json:"pro_id,omitempty" db:"pro_id"`
	ProCode       string     `json:"pro_code" db:"pro_code"`
	ProName       string     `json:"pro_name" db:"pro_name"`
	MinQty        *float64   `json:"min_qty" db:"min_qty"`
	MinQtyStr     *string    `json:"min_qty_str" db:"min_qty_str"`
	MaxQty        *float64   `json:"max_qty" db:"max_qty"`
	MaxQtyStr     *string    `json:"max_qty_str" db:"max_qty_str"`
	DiscPerc      *float64   `json:"disc_perc" db:"disc_perc"`
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
