package model

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	CustID      string    `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	StockID     int64     `gorm:"column:stock_id;primaryKey;autoIncrement" json:"stock_id"`
	StockDate   time.Time `gorm:"column:stock_date" json:"stock_date"`
	TrCode      string    `gorm:"column:tr_code" json:"tr_code"`
	TrNo        string    `gorm:"column:tr_no" json:"tr_no"`
	WhID        int64     `gorm:"column:wh_id" json:"wh_id"`
	ProID       int64     `gorm:"column:pro_id" json:"pro_id"`
	ItemCdn     int64     `gorm:"column:item_cdn" json:"item_cdn"`
	QtyIn       float64   `gorm:"column:qty_in" json:"qty_in"`
	QtyOut      float64   `gorm:"column:qty_out" json:"qty_out"`
	QtyInOrder  float64   `gorm:"column:qty_in_order" json:"qty_in_order"`
	QtyOutOrder float64   `gorm:"column:qty_out_order" json:"qty_out_order"`
	UnitPrice   float64   `gorm:"column:unit_price" json:"unit_price"`
	Cogs        float64   `gorm:"column:cogs" json:"cogs"`
	RefDetId    int64     `gorm:"ref_det_id" json:"ref_det_id"`
	CreatedAt   int64     `gorm:"created_at" json:"created_at"`
}

func (Stock) TableName() string {
	return "inv.stock"
}

func (m *Stock) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UTC().Unix()
	return nil
}

func (m *Stock) BeforeUpdate(trx *gorm.DB) (err error) {
	return nil
}

type StockReport struct {
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	UnitId1     string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2   int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	PurchPrice1 float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2 float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3 float64 `gorm:"column:purch_price3" json:"purch_price3"`
	SellPrice1  float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2  float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3  float64 `gorm:"column:sell_price3" json:"sell_price3"`
	Qty         float64 `gorm:"column:qty" json:"qty"`
	QtyOrder    float64 `gorm:"column:qty_order" json:"qty_order"`
	IsActive    bool    `gorm:"column:is_active" json:"is_active"`
	Vat         float64 `gorm:"column:vat" json:"vat"`
	VatLgPurch  float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	VatLgSell   float64 `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
	// Order detail qty from sls.order_detail for the specified outlet
	OrderQty1 float64 `gorm:"column:order_qty1" json:"-"`
	OrderQty2 float64 `gorm:"column:order_qty2" json:"-"`
	OrderQty3 float64 `gorm:"column:order_qty3" json:"-"`
}

type StockOpnameLookup struct {
	ID   int64  `gorm:"column:id" json:"id"`
	Code string `gorm:"column:code" json:"code"`
	Name string `gorm:"column:name" json:"name"`
}
