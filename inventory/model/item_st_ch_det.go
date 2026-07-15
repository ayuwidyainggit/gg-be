package model

import "time"

type ItemStChDet struct {
	IscDetId *int       `gorm:"column:isc_det_id;primaryKey" json:"isc_det_id"`
	CustID   string     `gorm:"column:cust_id" json:"cust_id"`
	IscNo    string     `gorm:"column:isc_no" json:"isc_no"`
	SeqNo    int        `gorm:"column:seq_no" json:"seq_no"`
	ProID    int        `gorm:"column:pro_id" json:"pro_id"`
	Qty      *float64   `gorm:"column:qty" json:"qty"`
	QtyStr   *string    `gorm:"column:qty_str" json:"qty_str"`
	BatchNo  *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate  *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ItemStChDet) TableName() string {
	return "inv.item_st_ch_det"
}

type ItemStChDetResponse struct {
	IscDetId *int       `gorm:"column:isc_det_id;primaryKey" json:"isc_det_id"`
	CustID   string     `gorm:"column:cust_id" json:"cust_id"`
	IscNo    string     `gorm:"column:isc_no" json:"isc_no"`
	SeqNo    int        `gorm:"column:seq_no" json:"seq_no"`
	ProID    int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode  string     `gorm:"column:pro_code" json:"pro_code"`
	ProName  string     `gorm:"column:pro_name" json:"pro_name"`
	Qty      *float64   `gorm:"column:qty" json:"qty"`
	QtyStr   *string    `gorm:"column:qty_str" json:"qty_str"`
	BatchNo  *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate  *time.Time `gorm:"column:exp_date" json:"exp_date"`
}

func (ItemStChDetResponse) TableName() string {
	return "inv.item_st_ch_det"
}
