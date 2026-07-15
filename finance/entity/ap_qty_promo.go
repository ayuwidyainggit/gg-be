package entity

type CreateApQtyPromoBody struct {
	ProID           int      `json:"pro_id"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	Total           *float64 `json:"total"`
	SeqNo           *int64   `json:"seq_no"`
	ConvUnit2       float64  `json:"conv_unit2"`
	ConvUnit3       float64  `json:"conv_unit3"`
	ConvUnit4       float64  `json:"conv_unit4"`
	ConvUnit5       float64  `json:"conv_unit5"`
}

type ApQtyPromoResponse struct {
	ApQtyPromoID    int64    `json:"ap_qty_promo_id"`
	ProID           int      `json:"pro_id"`
	ProCode         string   `json:"pro_code"`
	ProName         string   `json:"pro_name"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	Total           *float64 `json:"total"`
	SeqNo           *int64   `json:"seq_no"`
	ConvUnit2       float64  `json:"conv_unit2"`
	ConvUnit3       float64  `json:"conv_unit3"`
	ConvUnit4       float64  `json:"conv_unit4"`
	ConvUnit5       float64  `json:"conv_unit5"`
}
type UpdateApQtyPromoBody struct {
	ApQtyPromoID    *int64   `json:"ap_qty_promo_id"`
	ProID           int      `json:"pro_id"`
	PurchPrice      *float64 `json:"purch_price"`
	PurchPriceLevel *int64   `json:"purch_price_level"`
	Qty             *float64 `json:"qty"`
	QtyStr          *string  `json:"qty_str"`
	Total           *float64 `json:"total"`
	SeqNo           *int64   `json:"seq_no"`
}
