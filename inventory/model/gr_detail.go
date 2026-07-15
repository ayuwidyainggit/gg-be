package model

import (
	"time"

	"gorm.io/gorm"
)

type GrDetCreate struct {
	ID         int64      `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID     string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo       string     `gorm:"column:gr_no;primaryKey" json:"gr_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ProID      int64      `gorm:"column:pro_id" json:"pro_id"`
	ItemType   int        `gorm:"column:item_type" json:"item_type"`
	Qty        int        `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1 float64    `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2 float64    `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3 float64    `gorm:"column:unit_price3" json:"unit_price3"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat        *float64   `gorm:"column:vat" json:"vat"`
	VatBg      *float64   `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch *float64   `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax  *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	ConvUnit2  *float64   `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  *float64   `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  *float64   `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  *float64   `gorm:"column:conv_unit5" json:"conv_unit5"`
	Qty1       int        `gorm:"column:qty1" json:"qty1"`
	Qty2       int        `gorm:"column:qty2" json:"qty2"`
	Qty3       int        `gorm:"column:qty3" json:"qty3"`
	QtyShip1   int        `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2   int        `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3   int        `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyShip    int        `gorm:"column:qty_ship" json:"qty_ship"`
}

func (GrDetCreate) TableName() string {
	return "inv.gr_det"
}

func (m *GrDetCreate) BeforeCreate(trx *gorm.DB) (err error) {
	return nil
}

type GrDet struct {
	ID         int        `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID     string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo       string     `gorm:"column:gr_no;primaryKey" json:"gr_no"`
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
	UnitId1    *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2    *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3    *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4    *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5    *string    `gorm:"column:unit_id5" json:"unit_id5"`
	EmbInc     *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc     *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo  *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat        *float64   `gorm:"column:vat" json:"vat"`
	VatBg      *float64   `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch *float64   `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax  *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	Qty1       *float64   `gorm:"column:qty1" json:"qty1"`
	Qty2       *float64   `gorm:"column:qty2" json:"qty2"`
	Qty3       *float64   `gorm:"column:qty3" json:"qty3"`
	QtyShip1   *float64   `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2   *float64   `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3   *float64   `gorm:"column:qty_ship3" json:"qty_ship3"`
}

func (GrDet) TableName() string {
	return "inv.gr_det"
}

type GrDetList struct {
	ID           int        `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID       string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo         string     `gorm:"column:gr_no;primaryKey" json:"gr_no"`
	SeqNo        int        `gorm:"column:seq_no" json:"seq_no"`
	ProID        int64      `gorm:"column:pro_id" json:"pro_id"`
	ProCode      string     `gorm:"column:pro_code" json:"pro_code"`
	ProName      string     `gorm:"column:pro_name" json:"pro_name"`
	ConvUnit2    float64    `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3    float64    `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4    float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5    float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	ItemType     int        `gorm:"column:item_type" json:"item_type"`
	Qty          float64    `gorm:"column:qty" json:"qty"`
	QtyStr       *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitId1      *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2      *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3      *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4      *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5      *string    `gorm:"column:unit_id5" json:"unit_id5"`
	UnitPrice1   float64    `gorm:"column:unit_price1" json:"purch_price1"`
	UnitPrice2   float64    `gorm:"column:unit_price2" json:"purch_price2"`
	UnitPrice3   float64    `gorm:"column:unit_price3" json:"purch_price3"`
	EmbInc       *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc       *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo    *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo      *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate      *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat          float64    `gorm:"column:vat" json:"vat"`
	VatBg        float64    `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgPurch   float64    `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate   *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax    *float64   `gorm:"column:excise_tax" json:"excise_tax"`
	DiscP        *float64   `gorm:"column:disc_p" json:"disc_p"`
	Qty1         float64    `gorm:"column:qty1" json:"qty1"`
	Qty2         float64    `gorm:"column:qty2" json:"qty2"`
	Qty3         float64    `gorm:"column:qty3" json:"qty3"`
	QtyRemaining float64    `gorm:"column:qty_remaining" json:"qty_remaining"`
	QtyShip1     *float64   `gorm:"column:qty_ship1" json:"qty_ship1"`
	QtyShip2     *float64   `gorm:"column:qty_ship2" json:"qty_ship2"`
	QtyShip3     *float64   `gorm:"column:qty_ship3" json:"qty_ship3"`
	QtyShip      *float64   `gorm:"column:qty_ship" json:"qty_ship"`
	WhQty        float64    `gorm:"column:wh_qty" json:"wh_qty"`
	Discount     *float64   `gorm:"column:discount" json:"discount"`
}

func (GrDetList) TableName() string {
	return "inv.gr_det"
}

type GrDetJoinGrList struct {
	ID         int64      `gorm:"column:gr_det_id;primaryKey" json:"gr_det_id"`
	CustID     string     `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	GrNo       string     `gorm:"column:gr_no;primaryKey" json:"gr_no"`
	SeqNo      int        `gorm:"column:seq_no" json:"seq_no"`
	ProID      int64      `gorm:"column:pro_id" json:"pro_id"`
	WhID       int64      `gorm:"column:wh_id" json:"wh_id"`
	ItemType   int        `gorm:"column:item_type" json:"item_type"`
	Qty        *float64   `gorm:"column:qty" json:"qty"`
	QtyStr     *string    `gorm:"column:qty_str" json:"qty_str"`
	UnitPrice1 *float64   `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2 *float64   `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3 *float64   `gorm:"column:unit_price3" json:"unit_price3"`
	UnitPrice4 *float64   `gorm:"column:unit_price4" json:"unit_price4"`
	UnitPrice5 *float64   `gorm:"column:unit_price5" json:"unit_price5"`
	UnitId1    *string    `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2    *string    `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3    *string    `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4    *string    `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5    *string    `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2  float64    `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3  float64    `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4  float64    `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5  float64    `gorm:"column:conv_unit5" json:"conv_unit5"`
	EmbInc     *float64   `gorm:"column:emb_inc" json:"emb_inc"`
	EmbExc     *float64   `gorm:"column:emb_exc" json:"emb_exc"`
	InvoiceNo  *string    `gorm:"column:invoice_no" json:"invoice_no"`
	BatchNo    *string    `gorm:"column:batch_no" json:"batch_no"`
	ExpDate    *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat        *float64   `gorm:"column:vat" json:"vat"`
	VatBg      *float64   `gorm:"column:vat_bg" json:"vat_bg"`
	VatVgPurch *float64   `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	ExciseRate *float64   `gorm:"column:excise_rate" json:"excise_rate"`
	ExciseTax  *float64   `gorm:"column:excise_tax" json:"excise_tax"`
}

func (GrDetJoinGrList) TableName() string {
	return "inv.gr_det"
}

type ProductInvoiceBalances struct {
	CustID    string `gorm:"column:cust_id;primaryKey" json:"cust_id"`
	InvoiceNo string `gorm:"column:invoice_no;primaryKey" json:"invoice_no"`
	ProID     int64  `gorm:"column:pro_id;primaryKey" json:"pro_id"`
	Qty       int    `gorm:"column:qty" json:"qty"`
}

func (ProductInvoiceBalances) TableName() string {
	return "inv.product_invoice_balances"
}
