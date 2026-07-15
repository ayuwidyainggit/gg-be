package entity

/* API Spec
"promo_criterias":[
  {
    "slab_desc": "",
    "slab_rule_type": 1, // 1=quantity, 2=value
    "slab_rule_type_name": "Quantity",
    "slab_rule_from": 1,
    "slab_rule_to": 5,
    "slab_rule_uom": 1, // 1=smallest, 2=middle, 3=largest
    "slab_rule_uom_name": "Smallest",
    "slab_reward_type": 1, // 1=Quantity, 2=Fixed Value, 3=Percentage
    "slab_reward_type_name": "Quantity",
    "slab_reward": 5
    "slab_reward_uom": 1, // 1=smallest, 2=middle, 3=largest
    "slab_rule_uom_name": "Smallest",
  }
],
*/

type PromoCriteria struct {
	CustID             string          `json:"cust_id,omitempty"`
	PromoID            string          `json:"promo_id,omitempty"`
	SlabID             *int64          `json:"slab_id"`
	SlabDesc           string          `json:"slab_desc" validate:"max=100"`
	SlabRuleType       int             `json:"slab_rule_type" validate:"oneof=1 2"` // 1=quantity, 2=value
	SlabRuleTypeName   string          `json:"slab_rule_type_name"`
	SlabRuleFrom       float64         `json:"slab_rule_from"`
	SlabRuleTo         float64         `json:"slab_rule_to"`
	SlabRuleUom        int             `json:"slab_rule_uom" validate:"oneof=1 2 3"` // 1=smallest, 2=middle, 3=largest
	SlabRuleUomName    string          `json:"slab_rule_uom_name"`
	SlabRewardType     PromoRewardType `json:"slab_reward_type" validate:"oneof=1 2 3"`
	SlabRewardTypeName string          `json:"slab_reward_type_name"`
	SlabReward         float64         `json:"slab_reward"`
	SlabRewardUom      int             `json:"slab_reward_uom" validate:"oneof=0 1 2 3"` // 1=smallest, 2=middle, 3=largest
	SlabRewardUomName  string          `json:"slab_reward_uom_name"`
}
