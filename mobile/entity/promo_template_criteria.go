package entity

type PromoTemplateCriteria struct {
	CustID              string  `json:"cust_id,omitempty"`
	PromoTemplateID     string  `json:"promo_template_id,omitempty"`
	PromoTemplateSlabID *int64  `json:"promo_template_slab_id"`
	SlabDesc            string  `json:"slab_desc" validate:"max=100"`
	SlabRuleType        int     `json:"slab_rule_type" validate:"oneof=1 2"` // 1=quantity, 2=value
	SlabRuleTypeName    string  `json:"slab_rule_type_name"`
	SlabRuleFrom        float64 `json:"slab_rule_from"`
	SlabRuleTo          float64 `json:"slab_rule_to"`
	SlabRuleUom         int     `json:"slab_rule_uom" validate:"oneof=1 2 3"` // 1=smallest, 2=middle, 3=largest
	SlabRuleUomName     string  `json:"slab_rule_uom_name"`
	SlabRewardType      int     `json:"slab_reward_type" validate:"oneof=1 2 3"` // 1=Quantity, 2=Fixed Value, 3=Percentage
	SlabRewardTypeName  string  `json:"slab_reward_type_name"`
	SlabReward          float64 `json:"slab_reward"`
	SlabRewardUom       int     `json:"slab_reward_uom" validate:"oneof=1 2 3"` // 1=smallest, 2=middle, 3=largest
	SlabRewardUomName   string  `json:"slab_reward_uom_name"`
}
