package model

import (
	"time"

	"gorm.io/gorm"
)

type Promotion struct {
	CustID               string     `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID              string     `gorm:"column:promo_id;primaryKey;autoIncrement:false" json:"promo_id"`
	PromoDesc            string     `gorm:"column:promo_desc" json:"promo_desc"`
	PromoType            int64      `gorm:"column:promo_type" json:"promo_type"`
	ExistingPromoID      string     `gorm:"column:existing_promo_id" json:"existing_promo_id"`
	PromoStatusID        int        `gorm:"column:promo_status_id" json:"promo_status_id,omitempty"`
	IsMultiplied         bool       `gorm:"column:is_multiplied" json:"is_multiplied"`
	IsBudgetReference    bool       `gorm:"column:is_budget_reference" json:"is_budget_reference"`
	BudgetReferenceType  int        `gorm:"column:budget_reference_type" json:"budget_reference_type"`
	BudgetReferenceID    int64      `gorm:"column:budget_reference_id" json:"budget_reference_id"`
	BudgetControlLevel   int64      `gorm:"column:budget_control_level" json:"budget_control_level"`
	BudgetAmount         float64    `gorm:"column:budget_amount" json:"budget_amount"`
	ExecutionLevel       int64      `gorm:"column:execution_level" json:"execution_level"`
	EffectiveFrom        *time.Time `gorm:"column:effective_from" json:"effective_from"`
	EffectiveTo          *time.Time `gorm:"column:effective_to" json:"effective_to"`
	IsClaimable          bool       `gorm:"column:is_claimable" json:"is_claimable"`
	ClaimDays            int64      `gorm:"column:claim_days" json:"claim_days"`
	MaxDiscountType      int64      `gorm:"column:max_discount_type" json:"max_discount_type"`
	MaxDiscountOutlet    float64    `gorm:"column:max_discount_outlet" json:"max_discount_outlet"`
	MaxInvoiceOutlet     float64    `gorm:"column:max_invoice_outlet" json:"max_invoice_outlet"`
	CreatedAt            time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy            string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy            string     `gorm:"column:updated_by" json:"updated_by"`
	MaxDiscountOutletUom int        `gorm:"max_discount_outlet_uom" json:"max_discount_outlet_uom"`
	Remarks              string     `gorm:"column:remarks" json:"remarks,omitempty"`
}

func (Promotion) TableName() string {
	return "sls.promotions"
}

func (m *Promotion) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now().UTC()
	m.UpdatedAt = now

	return nil
}

type PromotionRead struct {
	CustID               string     `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoID              string     `gorm:"column:promo_id;primaryKey;autoIncrement:false" json:"promo_id"`
	PromoDesc            string     `gorm:"column:promo_desc" json:"promo_desc"`
	PromoType            int64      `gorm:"column:promo_type" json:"promo_type"`
	ExistingPromoID      string     `gorm:"column:existing_promo_id" json:"existing_promo_id"`
	PromoStatusID        int        `gorm:"column:promo_status_id" json:"promo_status_id,omitempty"`
	IsMultiplied         bool       `gorm:"column:is_multiplied" json:"is_multiplied"`
	IsBudgetReference    bool       `gorm:"column:is_budget_reference" json:"is_budget_reference"`
	BudgetReferenceType  int        `gorm:"column:budget_reference_type" json:"budget_reference_type"`
	BudgetReferenceID    int64      `gorm:"column:budget_reference_id" json:"budget_reference_id"`
	BudgetControlLevel   int64      `gorm:"column:budget_control_level" json:"budget_control_level"`
	BudgetAmount         float64    `gorm:"column:budget_amount" json:"budget_amount"`
	ExecutionLevel       int64      `gorm:"column:execution_level" json:"execution_level"`
	EffectiveFrom        *time.Time `gorm:"column:effective_from" json:"effective_from"`
	EffectiveTo          *time.Time `gorm:"column:effective_to" json:"effective_to"`
	IsClaimable          bool       `gorm:"column:is_claimable" json:"is_claimable"`
	ClaimDays            int64      `gorm:"column:claim_days" json:"claim_days"`
	MaxDiscountType      int64      `gorm:"column:max_discount_type" json:"max_discount_type"`
	MaxDiscountOutlet    float64    `gorm:"column:max_discount_outlet" json:"max_discount_outlet"`
	MaxInvoiceOutlet     float64    `gorm:"column:max_invoice_outlet" json:"max_invoice_outlet"`
	CreatedAt            time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy            string     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy            string     `gorm:"column:updated_by" json:"updated_by"`
	MaxDiscountOutletUom int        `gorm:"max_discount_outlet_uom" json:"max_discount_outlet_uom"`
	Remarks              string     `gorm:"column:remarks" json:"remarks,omitempty"`
}

func (PromotionRead) TableName() string {
	return "sls.promotions"
}

type PromotionMobileList struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	PromoID          string     `gorm:"column:promo_id" json:"promo_id"`
	PromoDesc        string     `gorm:"column:promo_desc" json:"promo_desc"`
	EffectiveFrom    *time.Time `gorm:"column:effective_from" json:"effective_from"`
	EffectiveTo      *time.Time `gorm:"column:effective_to" json:"effective_to"`
	IsMultiplied     bool       `gorm:"column:slab_multiplied" json:"is_multiplied"`
	MaxInvoiceOutlet float64    `gorm:"column:max_invoice_per_outlet" json:"max_invoice_outlet"`
}

func (PromotionMobileList) TableName() string {
	return "promo.promotions"
}

type PromotedProductRead struct {
	PromoID     string  `gorm:"column:promo_id" json:"promo_id"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	IsMandatory bool    `gorm:"column:is_mandatory" json:"is_mandatory"`
	MinBuyValue float64 `gorm:"column:min_buy_value" json:"min_buy_value"`
}

func (PromotedProductRead) TableName() string {
	return "sls.promo_additional_criterias"
}

type RewardProductRead struct {
	PromoID string `gorm:"column:promo_id" json:"promo_id"`
	ProID   int64  `gorm:"column:pro_id" json:"pro_id"`
	ProCode string `gorm:"column:pro_code" json:"pro_code"`
	ProName string `gorm:"column:pro_name" json:"pro_name"`
}

func (RewardProductRead) TableName() string {
	return "promo.promo_reward_products"
}

type AdditionalCriteriaRead struct {
	PromoID            string `gorm:"column:promo_id" json:"promo_id"`
	PromoAddCriteriaID int64  `gorm:"column:promo_add_criteria_id" json:"promo_add_criteria_id"`
	Attribute          string `gorm:"column:attribute" json:"attribute"`
	AttributeName      string `gorm:"column:attribute_name" json:"attribute_name"`
}

func (AdditionalCriteriaRead) TableName() string {
	return "promo.promo_additional_criterias"
}

type PromoCriteriaMinMax struct {
	PromoID     string  `gorm:"column:promo_id" json:"promo_id"`
	MinPurchase float64 `gorm:"column:min_purchase" json:"min_purchase"`
	MaxPurchase float64 `gorm:"column:max_purchase" json:"max_purchase"`
}

func (PromoCriteriaMinMax) TableName() string {
	return "promo.promo_criterias"
}

// Models for Promotion Mobile Detail
type PromotionMobileDetail struct {
	PromoID              string     `gorm:"column:promo_id" json:"promo_id"`
	PromoDesc            string     `gorm:"column:promo_desc" json:"promo_desc"`
	EffectiveFrom        *time.Time `gorm:"column:effective_from" json:"effective_from"`
	EffectiveTo          *time.Time `gorm:"column:effective_to" json:"effective_to"`
	IsMultiplied         bool       `gorm:"column:slab_multiplied" json:"is_multiplied"`
	MaxInvoiceOutlet     float64    `gorm:"column:max_invoice_per_outlet" json:"max_invoice_outlet"`
	MaxPromoUsage        float64    `gorm:"column:max_invoice_per_outlet" json:"max_promo_usage"`
	PromoType            string     `gorm:"column:promo_type" json:"promo_type"`
	MaxTotalRewardType   string     `gorm:"column:max_total_reward_type" json:"max_total_reward_type"`
	MaxTotalRewardValue  float64    `gorm:"column:max_total_reward_value" json:"max_total_reward_value"`
	RuleUom              string     `gorm:"column:rule_uom" json:"rule_uom"`
}

func (PromotionMobileDetail) TableName() string {
	return "promo.promotions"
}

type PromotedProductDetailRead struct {
	PromoID     string  `gorm:"column:promo_id" json:"promo_id"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	Mandatory   bool    `gorm:"column:mandatory" json:"mandatory"`
	MinBuyType  string  `gorm:"column:min_buy_type" json:"min_buy_type"`
	MinBuyQty   float64 `gorm:"column:min_buy_qty" json:"min_buy_qty"`
	MinBuyValue float64 `gorm:"column:min_buy_value" json:"min_buy_value"`
	MinBuyUom   string  `gorm:"column:min_buy_uom" json:"min_buy_uom"`
}

func (PromotedProductDetailRead) TableName() string {
	return "promo.promotion_product_criteria"
}

type PromotionCriteriaDetailRead struct {
	PromoID     string  `gorm:"column:promo_id" json:"promo_id"`
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	CountPromo  int     `gorm:"column:count_promo" json:"count_promo"`
	MinPurchase float64 `gorm:"column:min_purchase" json:"min_purchase"`
	MaxPurchase float64 `gorm:"column:max_purchase" json:"max_purchase"`
	Uom         string  `gorm:"column:uom" json:"uom"`
}

type PromotionRewardDetailRead struct {
	PromoID     string `gorm:"column:promo_id" json:"promo_id"`
	ProID       int64  `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string `gorm:"column:pro_code" json:"pro_code"`
	ProName     string `gorm:"column:pro_name" json:"pro_name"`
	RewardProID int64  `gorm:"column:id" json:"reward_pro_id"`
}

func (PromotionRewardDetailRead) TableName() string {
	return "promo.promotion_reward_products"
}

type PromotionSlabRead struct {
	PromoID     string  `gorm:"column:promo_id" json:"promo_id"`
	SlabID      string  `gorm:"column:id" json:"slab_id"`
	SlabName    string  `gorm:"column:description" json:"slab_name"`
	Ordinal     int     `gorm:"column:ordinal" json:"ordinal"`
	RuleType    string  `gorm:"column:rule_type" json:"rule_type"`
	RangeFrom   float64 `gorm:"column:range_from" json:"range_from"`
	RangeTo     float64 `gorm:"column:range_to" json:"range_to"`
	RuleUom     string  `gorm:"column:rule_uom" json:"rule_uom"`
	RewardType  string  `gorm:"column:reward_type" json:"reward_type"`
	RewardUom   string  `gorm:"column:reward_uom" json:"reward_uom"`
	RewardValue float64 `gorm:"column:reward_value" json:"reward_value"`
}

func (PromotionSlabRead) TableName() string {
	return "promo.promotion_slabs"
}

type PromotionStrataRead struct {
	PromoID      string  `gorm:"column:promo_id" json:"promo_id"`
	StrataID     string  `gorm:"column:id" json:"strata_id"`
	StrataName   string  `gorm:"column:description" json:"strata_name"`
	Ordinal      int     `gorm:"column:ordinal" json:"ordinal"`
	RuleType     string  `gorm:"column:rule_type" json:"rule_type"`
	RangeFrom    float64 `gorm:"column:range_from" json:"range_from"`
	RangeTo      float64 `gorm:"column:range_to" json:"range_to"`
	RuleUom      string  `gorm:"column:rule_uom" json:"rule_uom"`
	RewardType   string  `gorm:"column:reward_type" json:"reward_type"`
	RewardUom    string  `gorm:"column:reward_uom" json:"reward_uom"`
	RewardValue  float64 `gorm:"column:reward_value" json:"reward_value"`
}

func (PromotionStrataRead) TableName() string {
	return "promo.promotion_strata"
}

type OutletTypeDetailRead struct {
	PromoID        string `gorm:"column:promo_id" json:"promo_id"`
	OutletTypeID   int64  `gorm:"column:ot_type_id" json:"outlet_type_id"`
	OutletTypeCode string `gorm:"column:ot_type_code" json:"outlet_type_code"`
	OutletTypeName string `gorm:"column:ot_type_name" json:"outlet_type_name"`
}

func (OutletTypeDetailRead) TableName() string {
	return "promo.promo_additional_criterias"
}

type OutletGroupDetailRead struct {
	PromoID         string `gorm:"column:promo_id" json:"promo_id"`
	OutletGroupID   int64  `gorm:"column:ot_grp_id" json:"outlet_group_id"`
	OutletGroupCode string `gorm:"column:ot_grp_code" json:"outlet_group_code"`
	OutletGroupName string `gorm:"column:ot_grp_name" json:"outlet_group_name"`
}

func (OutletGroupDetailRead) TableName() string {
	return "promo.promo_additional_criterias"
}

type OutletClassDetailRead struct {
	PromoID         string `gorm:"column:promo_id" json:"promo_id"`
	OutletClassID   int64  `gorm:"column:ot_class_id" json:"outlet_class_id"`
	OutletClassCode string `gorm:"column:ot_class_code" json:"outlet_class_code"`
	OutletClassName string `gorm:"column:ot_class_name" json:"outlet_class_name"`
}

func (OutletClassDetailRead) TableName() string {
	return "promo.promo_additional_criterias"
}

type PromotionOutletList struct {
	OtTypeID   int64  `gorm:"column:ot_type_id" json:"ot_type_id"`
	OutletCode string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName string `gorm:"column:outlet_name" json:"outlet_name"`
	Address1   string `gorm:"column:address1" json:"address1"`
	TodayVisit *bool  `gorm:"column:today_visit" json:"today_visit"`
}
