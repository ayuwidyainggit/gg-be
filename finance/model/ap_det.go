package model

type ApDet struct {
	ApDetID         int64    `gorm:"column:ap_det_id;primaryKey" json:"ap_det_id"`
	CustID          string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo            string   `gorm:"column:ap_no" json:"ap_no"`
	GrNo            string   `gorm:"column:gr_no" json:"gr_no"`
	SeqNo           int      `gorm:"column:seq_no" json:"seq_no"`
	ProID           int      `gorm:"column:pro_id" json:"pro_id"`
	PurchPrice      *float64 `gorm:"column:purch_price" json:"purch_price"`
	PurchPriceLevel *int64   `gorm:"column:purch_price_level" json:"purch_price_level"`
	Qty             *float64 `gorm:"column:qty" json:"qty"`
	QtyStr          *string  `gorm:"column:qty_str" json:"qty_str"`
	SubTotal        *float64 `gorm:"column:sub_total" json:"sub_total"`
	Disc            *float64 `gorm:"column:disc" json:"disc"`
	DiscValue       *float64 `gorm:"column:disc_value" json:"disc_value"`
	SubTotalBtax    *float64 `gorm:"column:sub_total_btax" json:"sub_total_btax"`
	Vat             *float64 `gorm:"column:vat" json:"vat"`
	VatValue        *float64 `gorm:"column:vat_value" json:"vat_value"`
	VatLg           *float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue      *float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total           *float64 `gorm:"column:total" json:"total"`
	VatBg           *float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue      *float64 `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	ConvUnit2       *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4       *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5       *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (ApDet) TableName() string {
	return "acf.ap_det"
}

type ApDetRead struct {
	ApDetID         int64    `gorm:"column:ap_det_id;primaryKey" json:"ap_det_id"`
	CustID          string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo            string   `gorm:"column:ap_no" json:"ap_no"`
	GrNo            string   `gorm:"column:gr_no" json:"gr_no"`
	SeqNo           int      `gorm:"column:seq_no" json:"seq_no"`
	ProID           int      `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string   `gorm:"column:pro_code" json:"pro_code"`
	ProName         string   `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice      *float64 `gorm:"column:purch_price" json:"purch_price"`
	PurchPriceLevel *int64   `gorm:"column:purch_price_level" json:"purch_price_level"`
	Qty             *float64 `gorm:"column:qty" json:"qty"`
	QtyStr          *string  `gorm:"column:qty_str" json:"qty_str"`
	SubTotal        *float64 `gorm:"column:sub_total" json:"sub_total"`
	Disc            *float64 `gorm:"column:disc" json:"disc"`
	DiscValue       *float64 `gorm:"column:disc_value" json:"disc_value"`
	SubTotalBtax    *float64 `gorm:"column:sub_total_btax" json:"sub_total_btax"`
	Vat             *float64 `gorm:"column:vat" json:"vat"`
	VatValue        *float64 `gorm:"column:vat_value" json:"vat_value"`
	VatLg           *float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue      *float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	Total           *float64 `gorm:"column:total" json:"total"`
	VatBg           *float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue      *float64 `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	ConvUnit2       float64  `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64  `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4       float64  `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5       float64  `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (ApDetRead) TableName() string {
	return "acf.ap_det"
}
