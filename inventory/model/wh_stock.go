package model

import (
	"gorm.io/gorm"
)

type WhStock struct {
	CustID        string   `gorm:"column:cust_id" json:"cust_id"`
	WhID          *int64   `gorm:"column:wh_id" json:"wh_id"`
	ProID         *int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty           *float64 `gorm:"column:qty" json:"qty"`
	QtyOnOrder    *float64 `gorm:"column:qty_on_order" json:"qty_on_order"`
	QtyOnShipping *float64 `gorm:"column:qty_on_shipping" json:"qty_on_shipping"`
	QtyBs         *float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyExp        *float64 `gorm:"column:qty_exp" json:"qty_exp"`
	StockId       *int64   `gorm:"column:stock_id" json:"stock_id"`
}

func (WhStock) TableName() string {
	return "inv.wh_stock"
}

func (m *WhStock) BeforeCreate(trx *gorm.DB) (err error) {

	return nil
}

func (m *WhStock) BeforeUpdate(trx *gorm.DB) (err error) {
	return nil
}

type WhStockList struct {
	CustID        string   `gorm:"column:cust_id" json:"cust_id"`
	WhID          *int64   `gorm:"column:wh_id" json:"wh_id"`
	ProID         *int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty           *float64 `gorm:"column:qty" json:"qty"`
	QtyOnOrder    *float64 `gorm:"column:qty_on_order" json:"qty_on_order"`
	QtyOnShipping *float64 `gorm:"column:qty_on_shipping" json:"qty_on_shipping"`
	QtyBs         *float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyExp        *float64 `gorm:"column:qty_exp" json:"qty_exp"`
	StockId       *int64   `gorm:"column:stock_id" json:"stock_id"`
}

func (WhStockList) TableName() string {
	return "inv.wh_stock"
}
