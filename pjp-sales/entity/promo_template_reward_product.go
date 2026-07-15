package entity

type PromoTemplateRewardProduct struct {
	CustID                string `json:"cust_id,omitempty"`
	PromoTemplateID       string `json:"promo_template_id,omitempty"`
	PromoTemplateRewardID *int64 `json:"promo_template_reward_id"`
	ProID                 int64  `json:"pro_id" validate:"required"`
	ProCode               string `json:"pro_code"`
	ProName               string `json:"pro_name"`
}
