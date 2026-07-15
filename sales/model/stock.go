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
