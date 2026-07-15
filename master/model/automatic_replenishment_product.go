package model

import (
	"time"
)

type AutomaticReplenishmentProduct struct {
	CustId          string     `json:"cust_id" db:"cust_id"`
	Id              int64      `json:"id" db:"id"`
	ProId           int64      `json:"pro_id" db:"pro_id"`
	ProCode         string     `json:"pro_code" db:"pro_code"`
	ProName         string     `json:"pro_name" db:"pro_name"`
	DistributorId   int64      `json:"distributor_id" db:"distributor_id"`
	DistributorCode string     `json:"distributor_code" db:"distributor_code"`
	DistributorName string     `json:"distributor_name" db:"distributor_name"`
	LimitAction     string     `json:"limit_action" db:"limit_action"`
	MaxOrderQty     int        `json:"max_order_qty" db:"max_order_qty"`
	MaxOrderType    string     `json:"max_order_type" db:"max_order_type"`
	MinStockQty     int        `json:"min_stock_qty" db:"min_stock_qty"`
	MinStockType    string     `json:"min_stock_type" db:"min_stock_type"`
	SafetyStockQty  int        `json:"safety_stock_qty" db:"safety_stock_qty"`
	SafetyStockType string     `json:"safety_stock_type" db:"safety_stock_type"`
	MinOrderQty     int        `json:"min_order_qty" db:"min_order_qty"`
	MinOrderType    string     `json:"min_order_type" db:"min_order_type"`
	IsActive        *bool      `json:"is_active" db:"is_active"`
	CreatedBy       int64      `json:"created_by" db:"created_by"`
	CreatedByName   string     `json:"created_by_name" db:"created_by_name"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedBy       *int64     `json:"updated_by" db:"updated_by"`
	UpdatedByName   *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
	DeletedBy       *int64     `json:"deleted_by" db:"deleted_by"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at"`
	IsDel           *bool      `json:"is_del" db:"is_del"`
}
