package entity

type CreateApDetBody struct {
	GrNo            string   `json:"gr_no"`
	SeqNo           int      `json:"seq_no"`
	ProID           int      `json:"pro_id"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	SubTotal        *float64 `json:"sub_total"`
	Disc            *float64 `json:"disc"`
	DiscValue       *float64 `json:"disc_value"`
	SubTotalBtax    *float64 `json:"sub_total_btax"`
	Vat             *float64 `json:"vat"`
	VatValue        *float64 `json:"vat_value"`
	VatLg           *float64 `json:"vat_lg"`
	VatLgValue      *float64 `json:"vat_lg_value"`
	Total           *float64 `json:"total"`
	VatBg           *float64 `json:"vat_bg"`
	VatBgValue      *float64 `json:"vat_bg_value"`
	ConvUnit2       float64  `json:"conv_unit2"`
	ConvUnit3       float64  `json:"conv_unit3"`
	ConvUnit4       float64  `json:"conv_unit4"`
	ConvUnit5       float64  `json:"conv_unit5"`
}
type UpdateApDetBody struct {
	ApDetID         *int64   `json:"ap_det_id"`
	GrNo            string   `json:"gr_no"`
	SeqNo           int      `json:"seq_no"`
	ProID           int      `json:"pro_id"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	SubTotal        *float64 `json:"sub_total"`
	Disc            *float64 `json:"disc"`
	DiscValue       *float64 `json:"disc_value"`
	SubTotalBtax    *float64 `json:"sub_total_btax"`
	Vat             *float64 `json:"vat"`
	VatValue        *float64 `json:"vat_value"`
	VatLg           *float64 `json:"vat_lg"`
	VatLgValue      *float64 `json:"vat_lg_value"`
	Total           *float64 `json:"total"`
	VatBg           *float64 `json:"vat_bg"`
	VatBgValue      *float64 `json:"vat_bg_value"`
}

type ApDetResponse struct {
	ApDetID         int64    `json:"ap_det_id"`
	GrNo            string   `json:"gr_no"`
	SeqNo           int      `json:"seq_no"`
	ProID           int      `json:"pro_id"`
	ProCode         string   `json:"pro_code"`
	ProName         string   `json:"pro_name"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	SubTotal        *float64 `json:"sub_total"`
	Disc            *float64 `json:"disc"`
	DiscValue       *float64 `json:"disc_value"`
	SubTotalBtax    *float64 `json:"sub_total_btax"`
	Vat             *float64 `json:"vat"`
	VatValue        *float64 `json:"vat_value"`
	VatLg           *float64 `json:"vat_lg"`
	VatLgValue      *float64 `json:"vat_lg_value"`
	Total           *float64 `json:"total"`
	VatBg           *float64 `json:"vat_bg"`
	VatBgValue      *float64 `json:"vat_bg_value"`
	ConvUnit2       float64  `json:"conv_unit2"`
	ConvUnit3       float64  `json:"conv_unit3"`
	ConvUnit4       float64  `json:"conv_unit4"`
	ConvUnit5       float64  `json:"conv_unit5"`
}
