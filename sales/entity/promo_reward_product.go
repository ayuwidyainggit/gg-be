package entity

type PromoRewardProduct struct {
	CustID        string `json:"cust_id,omitempty"`
	PromoID       string `json:"promo_id,omitempty"`
	PromoRewardID *int64 `json:"promo_reward_id"`
	ProID         int64  `json:"pro_id" validate:"required"`
	ProCode       string `json:"pro_code"`
	ProName       string `json:"pro_name"`
}
