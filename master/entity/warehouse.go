package entity

import (
	"time"
)

type WarehouseQueryFilter struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit" validate:"required"`
	Query     string `query:"q"`
	Mode      string `query:"mode"`
	Gudang    *int   `query:"gudang"`
	Sort      string `query:"sort"`
	IsActive  *int   `query:"is_active"`
	IsReplenishment *bool `query:"is_replenishment"`
	StockType string `query:"stock_type"`
	DistributorIDs []int `query:"-"`
}

type WarehouseResponse struct {
	WarehouseId   int        `json:"wh_id"`
	WarehouseCode string     `json:"wh_code"`
	WarehouseName string     `json:"wh_name"`
	StockType     string     `json:"stock_type"`
	Latitude      string     `json:"latitude"`
	Longitude     string     `json:"longitude"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type WarehouseLookupResponse struct {
	WarehouseId   int    `json:"wh_id"`
	WarehouseCode string `json:"wh_code"`
	WarehouseName string `json:"wh_name"`
	StockType     string `json:"stock_type"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
}

type CreateWarehouseBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	WarehouseCode string `json:"wh_code" validate:"required,max=3,numeric"`
	WarehouseName string `json:"wh_name" validate:"required,max=20,alphanumericSpace"`
	StockType     string `json:"stock_type"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	IsActive      bool   `json:"is_active"`
}

type DetailWarehouseParams struct {
	WarehouseId int `params:"wh_id" validate:"required"`
}

type UpdateWarehouseParams struct {
	WarehouseId int `params:"wh_id" validate:"required"`
}

type DeleteWarehouseParams struct {
	WarehouseId int `params:"wh_id" validate:"required"`
}

type UpdateWarehouseRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	WarehouseCode string `json:"wh_code,omitempty" validate:"required,numeric,max=3,omitempty"`
	WarehouseName string `json:"wh_name,omitempty" validate:"max=20,alphanumericSpace,omitempty"`
	StockType     string `json:"stock_type,omitempty" validate:"max=3,alphanumericSpace,omitempty"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
