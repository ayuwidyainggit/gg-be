package model

import "time"

type ConsignmentDet struct {
	CustID     string     `gorm:"column:cust_id" json:"cust_id"`
	ConsNo     string     `gorm:"column:cons_no" json:"cons_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ConsDetID  *int       `gorm:"column:cons_det_id;primaryKey" json:"cons_det_id"`
	ProID      int        `gorm:"column:pro_id" json:"pro_id"`
	SellPrice1 *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2 *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3 *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Qty        *float64   `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	TotAmount  *float64   `gorm:"column:tot_amount" json:"tot_amount"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ConsignmentDet) TableName() string {
	return "sls.consign_det"
}

type ConsignmentDetRead struct {
	CustID     string     `gorm:"column:cust_id" json:"cust_id"`
	ConsNo     string     `gorm:"column:cons_no" json:"cons_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ConsDetID  *int       `gorm:"column:cons_det_id;primaryKey" json:"cons_det_id"`
	ProID      int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode    string     `gorm:"column:pro_code" json:"pro_code"`
	ProName    string     `gorm:"column:pro_name" json:"pro_name"`
	SellPrice1 *float64   `gorm:"column:sell_price1" json:"sell_price_1"`
	SellPrice2 *float64   `gorm:"column:sell_price2" json:"sell_price_2"`
	SellPrice3 *float64   `gorm:"column:sell_price3" json:"sell_price_3"`
	Qty        *float64   `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	TotAmount  *float64   `gorm:"column:tot_amount" json:"tot_amount"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ConsignmentDetRead) TableName() string {
	return "sls.consign_det"
}
