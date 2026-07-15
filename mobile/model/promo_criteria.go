package model

type PromoCriteria struct {
	CustID         string  `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID        string  `gorm:"column:promo_id" json:"promo_id"`
	SlabID         *int64  `gorm:"column:slab_id;primaryKey;autoIncrement:true" json:"slab_id"`
	SlabDesc       string  `gorm:"column:slab_desc" json:"slab_desc"`
	SlabRuleType   int64   `gorm:"column:slab_rule_type" json:"slab_rule_type"`
	SlabRuleFrom   float64 `gorm:"column:slab_rule_from" json:"slab_rule_from"`
	SlabRuleTo     float64 `gorm:"column:slab_rule_to" json:"slab_rule_to"`
	SlabRuleUom    int     `gorm:"column:slab_rule_uom" json:"slab_rule_uom"`
	SlabRewardType int     `gorm:"column:slab_reward_type" json:"slab_reward_type"`
	SlabReward     float64 `gorm:"column:slab_reward" json:"slab_reward"`
	SlabRewardUom  int     `gorm:"column:slab_reward_uom" json:"slab_reward_uom"`
}

func (PromoCriteria) TableName() string {
	return "sls.promo_criterias"
}

type ConsultPromoCriteria struct {
	CustID         string  `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID        string  `gorm:"column:promo_id" json:"promo_id"`
	PromoDesc      string  `gorm:"column:promo_desc" json:"promo_desc"`
	IsMultiplied   bool    `gorm:"column:is_multiplied" json:"is_multiplied"`
	SlabID         *int64  `gorm:"column:slab_id;primaryKey;autoIncrement:true" json:"slab_id"`
	SlabDesc       string  `gorm:"column:slab_desc" json:"slab_desc"`
	SlabRule       int64   `gorm:"column:slab_rule" json:"slab_rule"`
	SlabRuleType   int64   `gorm:"column:slab_rule_type" json:"slab_rule_type"`
	SlabRuleFrom   float64 `gorm:"column:slab_rule_from" json:"slab_rule_from"`
	SlabRuleTo     float64 `gorm:"column:slab_rule_to" json:"slab_rule_to"`
	SlabRuleUom    int     `gorm:"column:slab_rule_uom" json:"slab_rule_uom"`
	SlabRewardType int     `gorm:"column:slab_reward_type" json:"slab_reward_type"`
	SlabReward     float64 `gorm:"column:slab_reward" json:"slab_reward"`
	SlabRewardUom  int     `gorm:"column:slab_reward_uom" json:"slab_reward_uom"`
}

func (ConsultPromoCriteria) TableName() string {
	return "sls.promo_criterias"
}
