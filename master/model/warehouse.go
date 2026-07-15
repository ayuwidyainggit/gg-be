package model

import (
	"time"
)

type Warehouse struct {
	CustId        string     `db:"cust_id" json:"cust_id"`
	WarehouseId   int        `db:"wh_id" json:"wh_id"`
	WarehouseCode string     `db:"wh_code" json:"wh_code"`
	WarehouseName string     `db:"wh_name" json:"wh_name"`
	StockType     *string    `db:"stock_type" json:"stock_type"`
	Latitude      *string    `db:"latitude" json:"latitude"`
	Longitude     *string    `db:"longitude" json:"longitude"`
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

type WarehouseUpdate struct {
	WarehouseCode *string    `json:"wh_code,omitempty" sql:"wh_code"`
	WarehouseName *string    `json:"wh_name,omitempty" sql:"wh_name"`
	StockType     *string    `json:"stock_type,omitempty" sql:"stock_type"`
	Latitude      *string    `json:"latitude" sql:"latitude"`
	Longitude     *string    `json:"longitude" sql:"longitude"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
}

type TotalWarehouse struct {
	TotalWarehouse int `db:"total" json:"total"`
}
