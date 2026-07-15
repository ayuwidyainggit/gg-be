package model

type ApQtyPromo struct {
	ApQtyPromoID    *int64   `gorm:"column:ap_qty_promo_id;primaryKey" json:"ap_qty_promo_id"`
	CustID          string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo            string   `gorm:"column:ap_no" json:"ap_no"`
	ProID           int      `gorm:"column:pro_id" json:"pro_id"`
	PurchPrice      *float64 `gorm:"column:purch_price" json:"purch_price"`
	PurchPriceLevel *int64   `gorm:"column:purch_price_level" json:"purch_price_level"`
	Qty             *float64 `gorm:"column:qty" json:"qty"`
	QtyStr          *string  `gorm:"column:qty_str" json:"qty_str"`
	Total           *float64 `gorm:"column:total" json:"total"`
	SeqNo           *int64   `gorm:"column:seq_no" json:"seq_no"`
	ConvUnit2       float64  `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64  `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4       float64  `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5       float64  `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (ApQtyPromo) TableName() string {
	return "acf.ap_qty_promo"
}

type ApQtyPromoRead struct {
	ApQtyPromoID    *int64   `gorm:"column:ap_qty_promo_id;primaryKey" json:"ap_qty_promo_id"`
	CustID          string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo            string   `gorm:"column:ap_no" json:"ap_no"`
	ProID           int      `gorm:"column:pro_id" json:"pro_id"`
	ProCode         string   `gorm:"column:pro_code" json:"pro_code"`
	ProName         string   `gorm:"column:pro_name" json:"pro_name"`
	PurchPrice      *float64 `gorm:"column:purch_price" json:"purch_price"`
	PurchPriceLevel *int64   `gorm:"column:purch_price_level" json:"purch_price_level"`
	Qty             *float64 `gorm:"column:qty" json:"qty"`
	QtyStr          *string  `gorm:"column:qty_str" json:"qty_str"`
	Total           *float64 `gorm:"column:total" json:"total"`
	SeqNo           *int64   `gorm:"column:seq_no" json:"seq_no"`
	ConvUnit2       float64  `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3       float64  `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4       float64  `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5       float64  `gorm:"column:conv_unit5" json:"conv_unit5"`
}

func (ApQtyPromoRead) TableName() string {
	return "acf.ap_qty_promo"
}
