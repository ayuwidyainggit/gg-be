package model

type PromoTemplateCriteria struct {
	CustID              string  `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoTemplateID     string  `gorm:"column:promo_template_id" json:"promo_template_id"`
	PromoTemplateSlabID *int64  `gorm:"column:promo_template_slab_id;primaryKey;autoIncrement:true" json:"promo_template_slab_id"`
	SlabDesc            string  `gorm:"column:slab_desc" json:"slab_desc"`
	SlabRuleType        int64   `gorm:"column:slab_rule_type" json:"slab_rule_type"`
	SlabRuleFrom        float64 `gorm:"column:slab_rule_from" json:"slab_rule_from"`
	SlabRuleTo          float64 `gorm:"column:slab_rule_to" json:"slab_rule_to"`
	SlabRuleUom         int     `gorm:"column:slab_rule_uom" json:"slab_rule_uom"`
	SlabRewardType      int     `gorm:"column:slab_reward_type" json:"slab_reward_type"`
	SlabReward          float64 `gorm:"column:slab_reward" json:"slab_reward"`
	SlabRewardUom       int     `gorm:"column:slab_reward_uom" json:"slab_reward_uom"`
}

func (PromoTemplateCriteria) TableName() string {
	return "sls.promo_template_criterias"
}
