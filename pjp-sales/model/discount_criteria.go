package model

import "sales/entity"

type DiscountCriteria struct {
	CustID         string                    `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	DiscountID     string                    `gorm:"column:discount_id;primaryKey;autoIncrement:false" json:"discount_id"`
	SlabDesc       string                    `gorm:"column:slab_desc;primaryKey;autoIncrement:false" json:"slab_desc"`
	SlabRuleFrom   float64                   `gorm:"column:slab_rule_from" json:"slab_rule_from"`
	SlabRuleTo     float64                   `gorm:"column:slab_rule_to" json:"slab_rule_to"`
	SlabRewardType entity.DiscountRewardType `gorm:"column:slab_reward_type" json:"slab_reward_type"`
	SlabReward     float64                   `gorm:"column:slab_reward" json:"slab_reward"`
}

func (DiscountCriteria) TableName() string {
	return "sls.discount_criterias"
}
