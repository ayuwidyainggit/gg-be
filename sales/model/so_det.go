package model

import "time"

type SoDet struct {
	CustID      string     `gorm:"column:cust_id" json:"cust_id"`
	SoNo        string     `gorm:"column:so_no" json:"so_no"`
	SeqNo       int        `gorm:"column:seq_no" json:"seq_no"`
	SoDetID     *int       `gorm:"column:so_det_id;primaryKey" json:"so_det_id"`
	ProID       int        `gorm:"column:pro_id" json:"pro_id"`
	ItemType    int        `gorm:"column:item_type" json:"item_type"`
	Qty         *float64   `gorm:"column:qty" json:"qty"`
	QtyStr      *string    `gorm:"column:qty_str" json:"qty_str"`
	PurchPrice1 *float64   `gorm:"column:purch_price1" json:"purch_price_1"`
	PurchPrice2 *float64   `gorm:"column:purch_price2" json:"purch_price_2"`
	PurchPrice3 *float64   `gorm:"column:purch_price3" json:"purch_price_3"`
	SellPrice1  *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2  *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3  *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Amount      *float64   `gorm:"column:amount" json:"amount"`
	DiscValue   *float64   `gorm:"column:disc_value" json:"disc_value"`
	BatchNo     *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (SoDet) TableName() string {
	return "sls.so_det"
}

type SoDetRead struct {
	CustID      string     `gorm:"column:cust_id" json:"cust_id"`
	SoNo        string     `gorm:"column:so_no" json:"so_no"`
	SeqNo       int        `gorm:"column:seq_no" json:"seq_no"`
	SoDetID     *int       `gorm:"column:so_det_id;primaryKey" json:"so_det_id"`
	ProID       int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string     `gorm:"column:pro_code" json:"pro_code"`
	ProName     string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType    int        `gorm:"column:item_type" json:"item_type"`
	Qty         *float64   `gorm:"column:qty" json:"qty"`
	QtyStr      *string    `gorm:"column:qty_str" json:"qty_str"`
	PurchPrice1 *float64   `gorm:"column:purch_price1" json:"purch_price_1"`
	PurchPrice2 *float64   `gorm:"column:purch_price2" json:"purch_price_2"`
	PurchPrice3 *float64   `gorm:"column:purch_price3" json:"purch_price_3"`
	SellPrice1  *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2  *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3  *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Amount      *float64   `gorm:"column:amount" json:"amount"`
	DiscValue   *float64   `gorm:"column:disc_value" json:"disc_value"`
	BatchNo     *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (SoDetRead) TableName() string {
	return "sls.so_det"
}
