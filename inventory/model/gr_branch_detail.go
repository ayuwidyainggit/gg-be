package model

import (
	"time"
)

type GrBranchDetailCreate struct {
	GrBranchDetId int64    `gorm:"column:gr_branch_det_id;primaryKey" json:"gr_branch_det_id"`
	CustID        string   `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrBranchNo    string   `gorm:"column:gr_branch_no;primaryKey" json:"gr_branch_no"`
	SeqNo         int      `gorm:"column:seq_no" json:"seq_no"`
	ProID         int64    `gorm:"column:pro_id" json:"pro_id"`
	ItemType      int      `gorm:"column:item_type" json:"item_type"`
	Qty           int      `gorm:"column:qty" json:"qty"`
	QtyStr        *string  `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1    float64  `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2    float64  `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3    float64  `gorm:"column:unit_price3" json:"unit_price3"`
	UnitId1       *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2       *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3       *string  `gorm:"column:unit_id3" json:"unit_id3"`
	Vat           float64  `gorm:"column:vat" json:"vat"`
	VatValue      float64  `gorm:"column:vat_value" json:"vat_value"`
	VatBg         *float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch    *float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	Amount        float64  `gorm:"column:amount" json:"amount"`
	ConvUnit2     *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	// ConvUnit4     *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	// ConvUnit5     *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
	QtyShip      int `gorm:"column:qty_ship" json:"qty_ship"`
	QtyShip1     int `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2     int `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3     int `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyReceived  int `gorm:"column:qty_received" json:"qty_received"`
	QtyReceived1 int `gorm:"column:qty_received1" json:"qty_received1"`
	QtyReceived2 int `gorm:"column:qty_received2" json:"qty_received2"`
	QtyReceived3 int `gorm:"column:qty_received3" json:"qty_received3"`
}

func (GrBranchDetailCreate) TableName() string {
	return "inv.gr_branch_det"
}

// func (m *GrBranchDetailCreate) BeforeCreate(trx *gorm.DB) (err error) {
// 	return nil
// }

type GrBranchDet struct {
	GrBranchDetId int        `gorm:"column:gr_branch_det_id;primaryKey" json:"gr_branch_det_id"`
	CustID        string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrBranchNo    string     `gorm:"column:gr_branch_no;primaryKey" json:"gr_branch_no"`
	SeqNo         int        `gorm:"column:seq_no" json:"seq_no"`
	ProID         int        `gorm:"column:pro_id" json:"pro_id"`
	ItemType      int        `gorm:"column:item_type" json:"item_type"`
	Qty           *float64   `gorm:"column:qty" json:"qty"`
	QtyStr        *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1    *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2    *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3    *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4    *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5    *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	UnitId1       *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2       *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3       *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4       *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5       *string    `gorm:"column:unit_id5" json:"unit_id5"`
	EmbInc        *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc        *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo     *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo       *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate       *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Amount        *float64   `gorm:"column:amount" json:"amount"`
	Vat           *float64   `gorm:"column:vat" json:"vat"`
	VatValue      *float64   `gorm:"column:vat_value" json:"vat_value"`
	VatBg         *float64   `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch    *float64   `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate    *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax     *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	Qty1          *float64   `gorm:"column:qty1" json:"qty1"`
	Qty2          *float64   `gorm:"column:qty2" json:"qty2"`
	Qty3          *float64   `gorm:"column:qty3" json:"qty3"`
	QtyShip1      *float64   `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2      *float64   `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3      *float64   `gorm:"column:qty_ship3" json:"qty_ship3"`
}

func (GrBranchDet) TableName() string {
	return "inv.gr_branch_det"
}

type GrBranchDets struct {
	GrBranchDetId int64 `gorm:"column:gr_branch_det_id" json:"gr_branch_det_id"`
}

func (GrBranchDets) TableName() string {
	return "inv.gr_branch_det"
}

type GrBranchDetailList struct {
	GrBranchDetId int     `gorm:"column:gr_branch_det_id;primaryKey" json:"gr_branch_det_id"`
	CustID        string  `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrBranchNo    string  `gorm:"column:gr_branch_no;primaryKey" json:"gr_branch_no"`
	SeqNo         int     `gorm:"column:seq_no" json:"seq_no"`
	ProID         int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string  `gorm:"column:pro_code" json:"pro_code"`
	ProName       string  `gorm:"column:pro_name" json:"pro_name"`
	ConvUnit2     float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	// ConvUnit4    float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	// ConvUnit5    float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	ItemType int `gorm:"column:item_type" json:"item_type"`
	// Qty          float64    `gorm:"column:qty" json:"qty"`
	// QtyStr       *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitId1 *string `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2 *string `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3 *string `gorm:"column:unit_id3" json:"unit_id3"`
	// UnitId4      *string    `gorm:"column:unit_id4" json:"unit_id4"`
	// UnitId5      *string    `gorm:"column:unit_id5" json:"unit_id5"`
	UnitPrice1 float64 `gorm:"column:unit_price1" json:"purch_price1"`
	UnitPrice2 float64 `gorm:"column:unit_price2" json:"purch_price2"`
	UnitPrice3 float64 `gorm:"column:unit_price3" json:"purch_price3"`
	// EmbInc       *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	// EmbExc       *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	// InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	// BatchNo      *string    `gorm:"column:batch_no" json:"batch_no"`
	// ExpDate      *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat float64 `gorm:"column:vat" json:"vat"`
	// VatBg        float64    `gorm:"column:vat_bg" json:"vat_bg"`
	// VatLgPurch   float64    `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	// ExciseRate   *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	// ExciseTax    *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	// DiscP        *float64   `gorm:"column:disc_p" json:"disc_p"`
	// Qty1         float64    `gorm:"column:qty1" json:"qty1"`
	// Qty2         float64    `gorm:"column:qty2" json:"qty2"`
	// Qty3         float64    `gorm:"column:qty3" json:"qty3"`
	// QtyRemaining float64    `gorm:"column:qty_remaining" json:"qty_remaining"`
	QtyShip1     *float64 `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2     *float64 `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3     *float64 `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyShip      *float64 `gorm:"column:qty_ship" json:"qty_ship"`
	QtyReceived1 *float64 `gorm:"column:qty_received1" json:"qty_received1"`
	QtyReceived2 *float64 `gorm:"column:qty_received2" json:"qty_received2"`
	QtyReceived3 *float64 `gorm:"column:qty_received3" json:"qty_received3"`
	QtyReceived  *float64 `gorm:"column:qty_received" json:"qty_received"`
	VatValue     *float64 `gorm:"column:vat_value" json:"vat_value"`
	Amount       *float64 `gorm:"column:amount" json:"amount"`
	Qty1Alloc    *float64 `gorm:"column:qty1_alloc" json:"qty1_alloc"`
	Qty2Alloc    *float64 `gorm:"column:qty2_alloc" json:"qty2_alloc"`
	Qty3Alloc    *float64 `gorm:"column:qty3_alloc" json:"qty3_alloc"`
	// WhQty        float64    `gorm:"column:wh_qty" json:"wh_qty"`
	// Discount     *float64   `gorm:"column:discount" json:"discount"`
}

func (GrBranchDetailList) TableName() string {
	return "inv.gr_branch_det"
}

type GrBranchDetJoinGrBranchList struct {
	GrBranchDetId int64      `gorm:"column:gr_branch_det_id;primaryKey" json:"gr_branch_det_id"`
	CustID        string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrBranchNo    string     `gorm:"column:gr_branch_no;primaryKey" json:"gr_branch_no"`
	SeqNo         int        `gorm:"column:seq_no" json:"seq_no"`
	ProID         int64      `gorm:"column:pro_id" json:"pro_id"`
	WhID          int64      `gorm:"column:wh_id" json:"wh_id"`
	ItemType      int        `gorm:"column:item_type" json:"item_type"`
	Qty           *float64   `gorm:"column:qty" json:"qty"`
	QtyStr        *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1    *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2    *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3    *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4    *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5    *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	UnitId1       *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2       *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3       *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4       *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5       *string    `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2     float64    `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     float64    `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4     float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5     float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	EmbInc        *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc        *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo     *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo       *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate       *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Amount        *float64   `gorm:"column:amount" json:"amount"`
	Vat           *float64   `gorm:"column:vat" json:"vat"`
	VatValue      *float64   `gorm:"column:vat_value" json:"vat_value"`
	VatBg         *float64   `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch    *float64   `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate    *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax     *float64   `gorm:"column:excise_tax" json:"excise_tax"`
}

func (GrBranchDetJoinGrBranchList) TableName() string {
	return "inv.gr_branch_det"
}
