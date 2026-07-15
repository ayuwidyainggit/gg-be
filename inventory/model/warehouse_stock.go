package model

import (
	"time"

	"gorm.io/gorm"
)

type WarehouseStock struct {
	CustID        string  `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	WhID          int64   `gorm:"column:wh_id;primaryKey" json:"wh_id"`
	ProID         int64   `gorm:"column:pro_id;primaryKey" json:"pro_id"`
	Qty           float64 `gorm:"column:qty" json:"qty"`
	QtyOnOrder    float64 `gorm:"column:qty_on_order" json:"qty_on_order"`
	QtyOnShipping float64 `gorm:"column:qty_on_shipping" json:"qty_on_shipping"`
	QtyBs         float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyExp        float64 `gorm:"column:qty_exp" json:"qty_exp"`
	UpdatedAt     int64   `gorm:"column:updated_at" json:"updated_at"`
}

func (WarehouseStock) TableName() string {
	return "inv.warehouse_stock"
}

func (m *WarehouseStock) BeforeCreate(trx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().UTC().Unix()
	return nil
}

func (m *WarehouseStock) BeforeUpdate(trx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().UTC().Unix()
	return nil
}

type DistributorStockList struct {
	ProID         int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string  `gorm:"column:pro_code" json:"pro_code"`
	ProName       string  `gorm:"column:pro_name" json:"pro_name"`
	UnitId1       string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2       string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3       string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2     float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	PurchPrice1   float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2   float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3   float64 `gorm:"column:purch_price3" json:"purch_price3"`
	SellPrice1    float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2    float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3    float64 `gorm:"column:sell_price3" json:"sell_price3"`
	SupID         int64   `gorm:"column:sup_id" json:"sup_id"`
	SupCode       string  `gorm:"column:sup_code" json:"sup_code"`
	SupName       string  `gorm:"column:sup_name" json:"sup_name"`
	Qty           float64 `gorm:"column:qty" json:"qty"`
	QtyOnOrder    float64 `gorm:"column:qty_on_order" json:"qty_on_order"`
	QtyOnShipping float64 `gorm:"column:qty_on_shipping" json:"qty_on_shipping"`
	QtyBs         float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyExp        float64 `gorm:"column:qty_exp" json:"qty_exp"`
	UpdatedAt     int64   `gorm:"column:updated_at" json:"updated_at"`
}

func (DistributorStockList) TableName() string {
	return `"mst"."m_product"`
}

type WarehouseStockWhList struct {
	CustID    string `gorm:"column:cust_id" json:"cust_id"`
	WhID      int64  `gorm:"column:wh_id" json:"wh_id"`
	WhCode    string `gorm:"column:wh_code" json:"wh_code"`
	WhName    string `gorm:"column:wh_name" json:"wh_name"`
	StockType string `gorm:"stock_type" json:"stock_type"`
}

func (WarehouseStockWhList) TableName() string {
	return "inv.warehouse_stock"
}

type ProductWarehouseList struct {
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	Qty         float64 `gorm:"column:qty" json:"qty"`
	ConvUnit2   float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	SellPrice1  float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2  float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3  float64 `gorm:"column:sell_price3" json:"sell_price3"`
	UnitId1     string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     string  `gorm:"column:unit_id3" json:"unit_id3"`
	Vat         float64 `gorm:"column:vat" json:"vat"`
	VatLgPurch  float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	VatLgSell   float64 `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
	VatBg       float64 `gorm:"column:vat_bg" json:"vat_bg"`
	PurchPrice1 float64 `json:"purch_price1" gorm:"column:purch_price1"`
	PurchPrice2 float64 `json:"purch_price2" gorm:"column:purch_price2"`
	PurchPrice3 float64 `json:"purch_price3" gorm:"column:purch_price3"`
}

func (ProductWarehouseList) TableName() string {
	return "inv.warehouse_stock"
}
