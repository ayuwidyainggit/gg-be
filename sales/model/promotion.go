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
