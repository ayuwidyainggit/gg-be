package model

import "sales/entity"

type OrderReward struct {
	CustId        string `gorm:"cust_id" json:"cust_id"`
	RoNo          string `gorm:"ro_no" json:"ro_no"`
	OrderRewardId int    `gorm:"column:order_reward_id;primaryKey" json:"order_reward_id"`
	ReffId        string `gorm:"reff_id" json:"reff_id"`
	SlabId        *int   `gorm:"slab_id" json:"slab_id"`
	SlabDesc      string `gorm:"slab_desc" json:"slab_desc"`
	RewardTypeId  int    `gorm:"reward_type_id" json:"reward_type_id"`
}

func (OrderReward) TableName() string {
	return "sls.order_reward"
}

type OrderRewardRead struct {
	CustId        string `gorm:"cust_id" json:"cust_id"`
	RoNo          string `gorm:"ro_no" json:"ro_no"`
	OrderRewardId int    `gorm:"column:order_reward_id;primaryKey" json:"order_reward_id"`
	ReffId        string `gorm:"reff_id" json:"reff_id"`
	ReffName      string `gorm:"reff_name" json:"reff_name"`
	SlabId        *int   `gorm:"slab_id" json:"slab_id"`
	SlabDesc      string `gorm:"slab_desc" json:"slab_desc"`
	RewardTypeId  int    `gorm:"reward_type_id" json:"reward_type_id"`
}

func (OrderRewardRead) TableName() string {
	return "sls.order_reward"
}

type FullPromoRewardRead struct {
	CustId        string `gorm:"cust_id" json:"cust_id"`
	RoNo          string `gorm:"ro_no" json:"ro_no"`
	OrderRewardId int    `gorm:"column:order_reward_id;primaryKey" json:"order_reward_id"`
	// ReffId         string  `gorm:"reff_id" json:"reff_id"`
	// ReffName       string  `gorm:"reff_name" json:"reff_name"`
	SlabId         *int                   `gorm:"slab_id" json:"slab_id"`
	SlabDesc       string                 `gorm:"slab_desc" json:"slab_desc"`
	RewardTypeId   int                    `gorm:"reward_type_id" json:"reward_type_id"`
	PromoID        string                 `gorm:"column:promo_id" json:"promo_id"`
	PromoDesc      string                 `gorm:"column:promo_desc" json:"promo_desc"`
	IsMultiplied   bool                   `gorm:"column:is_multiplied" json:"is_multiplied"`
	SlabRuleType   int64                  `gorm:"column:slab_rule_type" json:"slab_rule_type"`
	SlabRuleFrom   float64                `gorm:"column:slab_rule_from" json:"slab_rule_from"`
	SlabRuleTo     float64                `gorm:"column:slab_rule_to" json:"slab_rule_to"`
	SlabRuleUom    int                    `gorm:"column:slab_rule_uom" json:"slab_rule_uom"`
	SlabRewardType entity.PromoRewardType `gorm:"column:slab_reward_type" json:"slab_reward_type"`
	SlabReward     float64                `gorm:"column:slab_reward" json:"slab_reward"`
	SlabRewardUom  int                    `gorm:"column:slab_reward_uom" json:"slab_reward_uom"`
}

func (FullPromoRewardRead) TableName() string {
	return "sls.order_reward"
}
