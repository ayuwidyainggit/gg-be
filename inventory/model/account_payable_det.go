package model

type RemainingQty struct {
	RemainingQty float64 `gorm:"column:remaining_qty" json:"remaining_qty"`
}

func (RemainingQty) TableName() string {
	return "acf.account_payable_detail"
}

type RemainingQtyProduct struct {
	GrNo         string `gorm:"column:gr_no;primaryKey" json:"gr_no"`
	ProID        int64  `gorm:"column:pro_id;primaryKey" json:"pro_id"`
	QtyRemaining int    `gorm:"column:qty_remaining" json:"qty_remaining"`
}

func (RemainingQtyProduct) TableName() string {
	return "acf.account_payable_detail"
}

type AccountPayableProductList struct {
	InvoiceNo         *string `gorm:"column:invoice_no" json:"invoice_no"`
	ProId             *int64  `gorm:"column:pro_id" json:"pro_id"`
	ProCode           *string `gorm:"column:pro_code" json:"pro_code"`
	ProName           *string `gorm:"column:pro_name" json:"pro_name"`
	UnitPrice1        float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2        float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3        float64 `gorm:"column:unit_price3" json:"unit_price3"`
	ConvUnit2         float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3         float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	SubTotal          float64 `gorm:"column:sub_total" json:"sub_total"`
	Disc              float64 `gorm:"column:disc" json:"disc"`
	DiscValue         float64 `gorm:"column:disc_value" json:"disc_value"`
	SubTotalBeforePpn float64 `gorm:"column:sub_total_before_ppn" json:"sub_total_before_ppn"`
	Vat               float64 `gorm:"column:vat" json:"vat"`
	VatValue          float64 `gorm:"column:vat_value" json:"vat_value"`
	Total             float64 `gorm:"column:total" json:"total"`
	VatLg             float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue        float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg             float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue        float64 `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	Qty               float64 `gorm:"column:qty" json:"qty"`
	QtyRemaining      float64 `gorm:"column:qty_remaining" json:"qty_remaining"`
}

func (AccountPayableProductList) TableName() string {
	return "acf.account_payable_detail"
}
