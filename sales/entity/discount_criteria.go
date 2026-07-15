package entity

type DiscountCriteria struct {
	CustID             string             `json:"cust_id,omitempty"`
	DiscountID         string             `json:"discount_id,omitempty"`
	SlabDesc           string             `json:"slab_desc" validate:"max=100"`
	SlabRuleFrom       float64            `json:"slab_rule_from"`
	SlabRuleTo         float64            `json:"slab_rule_to"`
	SlabRewardType     DiscountRewardType `json:"slab_reward_type" validate:"oneof=1 2"`
	SlabRewardTypeName string             `json:"slab_reward_type_name"`
	SlabReward         float64            `json:"slab_reward"`
}
