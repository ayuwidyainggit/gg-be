package model

import "time"

type RoDet struct {
	CustId      string     `gorm:"cust_id" json:"cust_id"`
	RoNo        string     `gorm:"ro_no" json:"ro_no"`
	SeqNo       int        `gorm:"seq_no" json:"seq_no"`
	RoDetID     *int       `gorm:"column:ro_det_id;primaryKey" json:"ro_det_id"`
	ProId       int        `gorm:"pro_id" json:"pro_id"`
	ItemType    int        `gorm:"item_type" json:"item_type"`
	Qty         *float64   `gorm:"qty" json:"qty"`
	QtyStr      *string    `gorm:"qty_str" json:"qty_str"`
	PurchPrice1 *float64   `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2 *float64   `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3 *float64   `gorm:"purch_price3" json:"purch_price3"`
	SellPrice1  *float64   `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2  *float64   `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3  *float64   `gorm:"sell_price3" json:"sell_price3"`
	Amount      *float64   `gorm:"amount" json:"amount"`
	DiscValue   *float64   `gorm:"disc_value" json:"disc_value"`
	BatchNo     *string    `gorm:"batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"exp_date" json:"exp_date"`
}

func (RoDet) TableName() string {
	return "sls.ro_det"
}

type RoDetRead struct {
	CustId      string     `gorm:"cust_id" json:"cust_id"`
	RoNo        string     `gorm:"ro_no" json:"ro_no"`
	SeqNo       int        `gorm:"seq_no" json:"seq_no"`
	RoDetID     *int       `gorm:"column:ro_det_id;primaryKey" json:"ro_det_id"`
	ProId       int        `gorm:"pro_id" json:"pro_id"`
	ProCode     string     `gorm:"column:pro_code" json:"pro_code"`
	ProName     string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType    int        `gorm:"item_type" json:"item_type"`
	Qty         *float64   `gorm:"qty" json:"qty"`
	QtyStr      *string    `gorm:"qty_str" json:"qty_str"`
	PurchPrice1 *float64   `gorm:"purch_price1" json:"purch_price1"`
	PurchPrice2 *float64   `gorm:"purch_price2" json:"purch_price2"`
	PurchPrice3 *float64   `gorm:"purch_price3" json:"purch_price3"`
	SellPrice1  *float64   `gorm:"sell_price1" json:"sell_price1"`
	SellPrice2  *float64   `gorm:"sell_price2" json:"sell_price2"`
	SellPrice3  *float64   `gorm:"sell_price3" json:"sell_price3"`
	Amount      *float64   `gorm:"amount" json:"amount"`
	DiscValue   *float64   `gorm:"disc_value" json:"disc_value"`
	BatchNo     *string    `gorm:"batch_no" json:"batch_no"`
	ExpDate     *time.Time `gorm:"exp_date" json:"exp_date"`
}

func (RoDetRead) TableName() string {
	return "sls.ro_det"
}
