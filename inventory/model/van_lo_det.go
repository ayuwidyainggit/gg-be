package model

import "time"

type VanLoDet struct {
	VanLoDetID int64      `gorm:"column:van_lo_det_id;primaryKey" json:"van_lo_det_id"`
	CustID     string     `gorm:"column:cust_id" json:"cust_id"`
	VanLoNo    string     `gorm:"column:van_lo_no" json:"van_lo_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ProID      int        `gorm:"column:pro_id" json:"pro_id"`
	ItemType   int        `gorm:"column:item_type" json:"item_type"`
	Qty        *float64   `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1 *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2 *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3 *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4 *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5 *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	EmbInc     *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc     *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
	ConvUnit2  *float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  *float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  *float64   `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  *float64   `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (VanLoDet) TableName() string {
	return "inv.van_lo_det"
}

type VanLoDetRead struct {
	VanLoDetID int64      `gorm:"column:van_lo_det_id;primaryKey" json:"van_lo_det_id"`
	CustID     string     `gorm:"column:cust_id" json:"cust_id"`
	VanLoNo    string     `gorm:"column:van_lo_no" json:"van_lo_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ProID      int        `gorm:"column:pro_id" json:"pro_id"`
	ProCode    string     `gorm:"column:pro_code" json:"pro_code"`
	ProName    string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType   int        `gorm:"column:item_type" json:"item_type"`
	Qty        *float64   `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1 *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2 *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3 *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4 *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5 *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	EmbInc     *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc     *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
	ConvUnit2  *float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  *float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  *float64   `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  *float64   `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (VanLoDetRead) TableName() string {
	return "inv.van_lo_det"
}
