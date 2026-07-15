package entity

type CreateApMoneyPromoBody struct {
	ProID      int      `json:"pro_id"`
	MoneyPromo *float64 `json:"money_promo"`
	SeqNo      int      `json:"seq_no"`
}
type ApMoneyPromoResponse struct {
	ApMoneyPromoID int64    `json:"ap_money_promo_id"`
	ProID          int      `json:"pro_id"`
	ProCode        string   `json:"pro_code"`
	ProName        string   `json:"pro_name"`
	MoneyPromo     *float64 `json:"money_promo"`
	SeqNo          int      `json:"seq_no"`
	ConvUnit2      float64  `json:"conv_unit2"`
	ConvUnit3      float64  `json:"conv_unit3"`
	ConvUnit4      float64  `json:"conv_unit4"`
	ConvUnit5      float64  `json:"conv_unit5"`
}

type UpdateApMoneyPromoBody struct {
	ApMoneyPromoID *int64   `json:"ap_money_promo_id"`
	ProID          int      `json:"pro_id"`
	MoneyPromo     *float64 `json:"money_promo"`
	SeqNo          int      `json:"seq_no"`
}
