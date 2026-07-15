package model

type SupplierReturnDet struct {
	ID               int64   `gorm:"column:supplier_return_det_id;primaryKey" json:"supplier_return_det_id"`
	CustID           string  `gorm:"column:cust_id" json:"cust_id"`
	SupplierReturnNo string  `gorm:"column:supplier_return_no" json:"supplier_return_no"`
	SeqNo            int     `gorm:"column:seq_no" json:"seq_no"`
	ProID            int     `gorm:"column:pro_id" json:"pro_id"`
	Qty              int     `gorm:"column:qty" json:"qty"`
	ItemCdn          *int64  `gorm:"column:item_cdn" json:"item_cdn"`
	ReturnReasonID   *int64  `gorm:"column:return_reason_id" json:"return_reason_id"`
	UnitPrice1       float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2       float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3       float64 `gorm:"column:unit_price3" json:"unit_price3"`
	SubTotal         float64 `gorm:"column:sub_total" json:"sub_total"`
	Discount         float64 `gorm:"column:discount" json:"discount"`
	DiscountValue    float64 `gorm:"column:discount_value" json:"discount_value"`
	Vat              float64 `gorm:"column:vat" json:"vat"`
	VatValue         float64 `gorm:"column:vat_value" json:"vavat_valuet"`
	VatLg            float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue       float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg            float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue       float64 `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	Total            float64 `gorm:"column:total" json:"total"`
}

func (SupplierReturnDet) TableName() string {
	return "inv.supplier_return_details"
}

type SupplierReturnDetGet struct {
	ID               int64   `gorm:"column:supplier_return_det_id;primaryKey" json:"supplier_return_det_id"`
	CustID           string  `gorm:"column:cust_id" json:"cust_id"`
	SupplierReturnNo string  `gorm:"column:supplier_return_no" json:"supplier_return_no"`
	SeqNo            int     `gorm:"column:seq_no" json:"seq_no"`
	ProID            int     `gorm:"column:pro_id" json:"pro_id"`
	ProCode          string  `gorm:"column:pro_code" json:"pro_code"`
	ProName          string  `gorm:"column:pro_name" json:"pro_name"`
	Qty              float64 `gorm:"column:qty" json:"qty"`
	UnitPrice1       float64 `gorm:"column:unit_price1" json:"unit_price1"`
	UnitPrice2       float64 `gorm:"column:unit_price2" json:"unit_price2"`
	UnitPrice3       float64 `gorm:"column:unit_price3" json:"unit_price3"`
	UnitId1          *string `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2          *string `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3          *string `gorm:"column:unit_id3" json:"unit_id3"`
	InvoiceQty       float64 `gorm:"column:invoice_qty" json:"invoice_qty"`
	RemainingQty     float64 `gorm:"column:remaining_qty" json:"remaining_qty"`
	ConvUnit2        float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3        float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ItemCdn          *int64  `gorm:"column:item_cdn" json:"item_cdn"`
	ReturnReasonID   int64   `gorm:"column:return_reason_id" json:"return_reason_id"`
	ReturnReasonName *string `gorm:"column:return_reason_name" json:"return_reason_name"`
	Discount         float64 `gorm:"column:discount" json:"discount"`
	DiscountValue    float64 `gorm:"column:discount_value" json:"discount_value"`
	Vat              float64 `gorm:"column:vat" json:"vat"`
	VatValue         float64 `gorm:"column:vat_value" json:"vat_value"`
	VatLg            float64 `gorm:"column:vat_lg" json:"vat_lg"`
	VatLgValue       float64 `gorm:"column:vat_lg_value" json:"vat_lg_value"`
	VatBg            float64 `gorm:"column:vat_bg" json:"vat_bg"`
	VatBgValue       float64 `gorm:"column:vat_bg_value" json:"vat_bg_value"`
	SubTotal         float64 `gorm:"column:sub_total" json:"sub_total"`
	Total            float64 `gorm:"column:total" json:"total"`
	WhQty            float64 `gorm:"column:wh_qty" json:"wh_qty"`
}

func (SupplierReturnDetGet) TableName() string {
	return "inv.supplier_return_details"
}
