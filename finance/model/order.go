package model

import (
	"errors"
	"fmt"
	"time"
)

type Order struct {
	InvoiceNo      string `json:"invoice_no" gorm:"column:invoice_no"`
	OutletID       int64  `json:"outlet_id" gorm:"column:outlet_id"`
	TaxInvoiceForm int    `json:"tax_invoice_form" gorm:"tax_invoice_form"`
}

func (Order) TableName() string {
	return "sls.order"
}

type OrderList struct {
	TaxesId int        `gorm:"column:taxes_id" json:"taxes_id"`
	MTaxId  int        `gorm:"column:m_tax_id" json:"m_tax_id"`
	TaxNo   string     `gorm:"column:tax_no" json:"tax_no"`
	TaxDate *time.Time `gorm:"invoice_date" json:"invoice_date"`
	Type    string     `gorm:"type" json:"type"`

	RoDate          *time.Time `gorm:"ro_date" json:"ro_date"`
	ValDate         *time.Time `gorm:"val_date" json:"val_date"`
	DueDate         *time.Time `gorm:"due_date" json:"due_date"`
	SalesmanId      *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode    *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesName       *string    `gorm:"sales_name" json:"sales_name"`
	WhId            *int64     `gorm:"wh_id" json:"wh_id"`
	WhCode          *string    `gorm:"wh_code" json:"wh_code"`
	WhName          *string    `gorm:"wh_name" json:"wh_name"`
	WhLatitude      *string    `gorm:"column:wh_latitude" json:"wh_latitude"`
	WhLongitude     *string    `gorm:"column:wh_longitude" json:"wh_longitude"`
	OutletID        *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode      *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      *string    `gorm:"column:outlet_name" json:"outlet_name"`
	OutletLatitude  *string    `gorm:"column:outlet_latitude" json:"outlet_latitude"`
	OutletLongitude *string    `gorm:"column:outlet_longitude" json:"outlet_longitude"`
	OutletAddress   *string    `gorm:"column:outlet_address" json:"outlet_address"`
	DeliveryDate    *time.Time `gorm:"delivery_date" json:"delivery_date"`
	OrderNo         string     `gorm:"order_no" json:"order_no"`
	PoNo            *string    `gorm:"po_no" json:"po_no"`
	VehicleNo       *string    `gorm:"vehicle_no" json:"vehicle_no"`
	PayType         *int64     `gorm:"pay_type" json:"pay_type"`
	ReffNo          *string    `gorm:"reff_no" json:"reff_no"`
	MobileID        *int64     `gorm:"mobile_id" json:"mobile_id"`
	SubTotal        *float64   `gorm:"sub_total" json:"sub_total"`
	Disc            *float64   `gorm:"disc" json:"disc"`
	DiscValue       *float64   `gorm:"disc_value" json:"disc_value"`
	PromoValue      *float64   `gorm:"promo_value" json:"promo_value"`
	CashDiscValue   *float64   `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1        *float64   `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2        *float64   `gorm:"tot_disc2" json:"tot_disc2"`
	Vat             *float64   `gorm:"vat" json:"vat"`
	VatValue        *float64   `gorm:"vat_value" json:"vat_value"`
	Total           *float64   `gorm:"total" json:"total"`
	DataStatus      *int64     `gorm:"data_status" json:"data_status"`
	InvoiceNo       *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time `gorm:"invoice_date" json:"invoice_date"`
	Status          *int       `gorm:"column:status" json:"status"`
}

// // TableName sets the insert table name for this struct type
// func (Taxes) TableName() string {
// 	return "acf.taxes"
// }

// TableName sets the insert table name for this struct type
func (OrderList) TableName() string {
	return "sls.order"
}

type CoretaxInvoiceOrderList struct {
	RoNo              string     `gorm:"column:ro_no" json:"ro_no"`
	RoDate            *time.Time `gorm:"column:ro_date" json:"ro_date"`
	InvoiceNo         *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate       *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	TaxIdentifierType string     `gorm:"column:tax_identifier_type" json:"tax_identifier_type"`
	TaxIdentifierNo   string     `gorm:"column:tax_identifier_no" json:"tax_identifier_no"`
	IdentityNo        string     `gorm:"column:identity_no" json:"identity_no"`
	TaxName           string     `gorm:"column:tax_name" json:"tax_name"`
	AddressTax        string     `gorm:"column:address_tax" json:"address_tax"`
	NITKU             string     `gorm:"column:nitku" json:"nitku"`
	SalesId           *int64     `gorm:"column:sales_id" json:"sales_id"`
	SalesCode         string     `gorm:"column:sales_code" json:"sales_code"`
	SalesName         string     `gorm:"column:sales_name" json:"sales_name"`
	OutletID          *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode        string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName        string     `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress1    string     `gorm:"column:outlet_address1" json:"outlet_address1"`
	OutletAddress2    string     `gorm:"column:outlet_address2" json:"outlet_address2"`
	OutletTaxAddress1 string     `gorm:"column:outlet_tax_address1" json:"outlet_tax_address1"`
	OutletTaxAddress2 string     `gorm:"column:outlet_tax_address2" json:"outlet_tax_address2"`
	TaxExtractDate    *time.Time `gorm:"column:tax_extract_date" json:"tax_extract_date"`
	Vat               float64    `gorm:"column:vat" json:"vat"`
	VatValue          float64    `gorm:"column:vat_value" json:"vat_value"`
	VatValueFinal     float64    `gorm:"column:vat_value_final" json:"vat_value_final"`
	SubTotalFinal     float64    `gorm:"column:sub_total_final" json:"sub_total_final"`
	DiscValueFinal    float64    `gorm:"column:disc_value_final" json:"disc_value_final"`
	PromoValueFinal   float64    `gorm:"column:promo_value_final" json:"promo_value_final"`
	Total             float64    `gorm:"column:total" json:"total"`
	TotalFinal        float64    `gorm:"column:total_final" json:"total_final"`
	DPP               float64    `gorm:"column:dpp" json:"dpp"`
}

func (CoretaxInvoiceOrderList) TableName() string {
	return "sls.order"
}

type CoretaxInvoiceOrderDetailRead struct {
	CustId          string     `gorm:"column:cust_id" json:"cust_id"`
	RoNo            string     `gorm:"column:ro_no" json:"ro_no"`
	SeqNo           int        `gorm:"column:seq_no" json:"seq_no"`
	OrderDetailID   *int       `gorm:"column:order_detail_id;primaryKey" json:"order_detail_id"`
	ProId           int64      `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string     `gorm:"column:pro_code" json:"pro_code"`
	ProCodeCoretax  string     `gorm:"column:pro_code_coretax" json:"pro_code_coretax"`
	ProName         string     `gorm:"column:pro_name" json:"pro_name"`
	ItemType        int        `gorm:"column:item_type" json:"item_type"`
	QtyFinal        float64    `gorm:"column:qty_final" json:"qty_final"`
	QtyPo           float64    `gorm:"column:qty_po" json:"qty_po"`
	Qty             float64    `gorm:"column:qty" json:"qty"`
	Qty1            float64    `gorm:"column:qty1" json:"qty1"`
	Qty2            float64    `gorm:"column:qty2" json:"qty2"`
	Qty3            float64    `gorm:"column:qty3" json:"qty3"`
	Qty4            float64    `gorm:"column:qty4" json:"qty4"`
	Qty5            float64    `gorm:"column:qty5" json:"qty5"`
	Qty1Final       float64    `gorm:"column:qty1_final" json:"qty1_final"`
	Qty2Final       float64    `gorm:"column:qty2_final" json:"qty2_final"`
	Qty3Final       float64    `gorm:"column:qty3_final" json:"qty3_final"`
	Qty4Final       float64    `gorm:"column:qty4_final" json:"qty4_final"`
	Qty5Final       float64    `gorm:"column:qty5_final" json:"qty5_final"`
	Qty1Stok        float64    `gorm:"column:qty1_stok" json:"qty1_stok"`
	Qty2Stok        float64    `gorm:"column:qty2_stok" json:"qty2_stok"`
	Qty3Stok        float64    `gorm:"column:qty3_stok" json:"qty3_stok"`
	PurchPrice1     float64    `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2     float64    `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3     float64    `gorm:"column:purch_price3" json:"purch_price3"`
	PurchPrice4     float64    `gorm:"column:purch_price4" json:"purch_price4"`
	PurchPrice5     float64    `gorm:"column:purch_price5" json:"purch_price5"`
	SellPrice1      float64    `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2      float64    `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3      float64    `gorm:"column:sell_price3" json:"sell_price3"`
	SellPrice4      float64    `gorm:"column:sell_price4" json:"sell_price4"`
	SellPrice5      float64    `gorm:"column:sell_price5" json:"sell_price5"`
	Amount          float64    `gorm:"column:amount" json:"amount"`
	AmountFinal     float64    `gorm:"column:amount" json:"amount_final"`
	DiscValue       float64    `gorm:"column:disc_value" json:"disc_value"`
	DiscValueFinal  float64    `gorm:"column:disc_value_final" json:"disc_value_final"`
	PromoValue      float64    `gorm:"column:promo_value" json:"promo_value"`
	PromoValueFinal float64    `gorm:"column:promo_value_final" json:"promo_value_final"`
	BatchNo         string     `gorm:"column:batch_no" json:"batch_no"`
	ExpDate         *time.Time `gorm:"column:exp_date" json:"exp_date"`
	Vat             float64    `gorm:"column:vat" json:"vat"`
	VatBg           float64    `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgSell       float64    `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
	VatValue        float64    `gorm:"column:vat_value" json:"vat_value"`
	VatValueFinal   float64    `gorm:"column:vat_value_final" json:"vat_value_final"`
	VatBgValue      float64    `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	VatLgSellValue  float64    `gorm:"column:vat_lg_sell_value" json:"vat_lg_sell_value"`
	UnitId1         string     `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2         string     `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3         string     `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4         string     `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5         string     `gorm:"column:unit_id5" json:"unit_id5"`
	UnitIdCoreTax1  string     `gorm:"column:unit_id_coretax1" json:"unit_id_coretax1"`
	UnitIdCoreTax2  string     `gorm:"column:unit_id_coretax2" json:"unit_id_coretax2"`
	UnitIdCoreTax3  string     `gorm:"column:unit_id_coretax3" json:"unit_id_coretax3"`
	ConvUnit2       int        `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       int        `gorm:"column:conv_unit3" json:"conv_unit3"`
	MpConvUnit2     int        `gorm:"column:mconv_unit2" json:"mconv_unit2"`
	MpConvUnit3     int        `gorm:"column:mconv_unit3" json:"mconv_unit3"`
	ConvUnit4       int        `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5       int        `gorm:"column:conv_unit5" json:"conv_unit5"`
	Notes           string     `gorm:"column:notes" json:"notes"`
	DiscountID      string     `gorm:"column:discount_id" json:"discount_id"`
	Extracted       bool
}

func (CoretaxInvoiceOrderDetailRead) TableName() string {
	return "sls.order_detail"
}

func (c *CoretaxInvoiceOrderDetailRead) SetExtracted() {
	c.Extracted = true
}

type CoretaxInvoiceOrderDetailReadMap map[int64]*CoretaxInvoiceOrderDetailRead

func (m CoretaxInvoiceOrderDetailReadMap) SetTempEmployeeValidationMap(id int64, value *CoretaxInvoiceOrderDetailRead) {
	value.Extracted = false
	m[id] = value
}
func (m CoretaxInvoiceOrderDetailReadMap) GetByID(id int64) (value *CoretaxInvoiceOrderDetailRead, err error) {
	val, ok := m[id]

	// If the key exists
	if !ok {
		return value, errors.New(fmt.Sprintf("%v Not Found", id))
	}

	return val, nil
}
