package model

type PromoTemplateRewardProduct struct {
	CustID                string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoTemplateID       string `gorm:"column:promo_template_id" json:"promo_template_id"`
	PromoTemplateRewardID *int64 `gorm:"column:promo_template_reward_id;primaryKey;autoIncrement:true" json:"promo_template_reward_id"`
	ProID                 int64  `gorm:"column:pro_id" json:"pro_id"`
}

func (PromoTemplateRewardProduct) TableName() string {
	return "sls.promo_template_reward_products"
}

type PromoTemplateRewardProductDetail struct {
	CustID  string `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	ProID   int64  `gorm:"column:pro_id" json:"pro_id"`
	ProCode string `gorm:"column:pro_code" json:"pro_code"`
	ProName string `gorm:"column:pro_name" json:"pro_name"`
}

func (PromoTemplateRewardProductDetail) TableName() string {
	return "mst.m_product"
}
