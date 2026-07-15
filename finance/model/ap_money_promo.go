package model

type ApMoneyPromo struct {
	ApMoneyPromoID int64    `gorm:"column:ap_money_promo_id;primaryKey" json:"ap_money_promo_id"`
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo           string   `gorm:"column:ap_no" json:"ap_no"`
	ProID          int      `gorm:"column:pro_id" json:"pro_id"`
	MoneyPromo     *float64 `gorm:"column:money_promo" json:"money_promo"`
	SeqNo          int      `gorm:"column:seq_no" json:"seq_no"`
}

func (ApMoneyPromo) TableName() string {
	return "acf.ap_money_promo"
}

type ApMoneyPromoRead struct {
	ApMoneyPromoID int64    `gorm:"column:ap_money_promo_id;primaryKey" json:"ap_money_promo_id"`
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ApNo           string   `gorm:"column:ap_no" json:"ap_no"`
	ProID          int      `gorm:"column:pro_id" json:"pro_id"`
	ProCode        string   `gorm:"column:pro_code" json:"pro_code"`
	ProName        string   `gorm:"column:pro_name" json:"pro_name"`
	MoneyPromo     *float64 `gorm:"column:money_promo" json:"money_promo"`
	SeqNo          int      `gorm:"column:seq_no" json:"seq_no"`
	ConvUnit2      float64  `json:"conv_unit2"`
	ConvUnit3      float64  `json:"conv_unit3"`
	ConvUnit4      float64  `json:"conv_unit4"`
	ConvUnit5      float64  `json:"conv_unit5"`
}

func (ApMoneyPromoRead) TableName() string {
	return "acf.ap_money_promo"
}
