package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type PromoTemplate struct {
	CustID                string    `gorm:"column:cust_id;primaryKey;autoIncrement:false" json:"cust_id"`
	PromoTemplateID       string    `gorm:"column:promo_template_id;primaryKey;autoIncrement:false" json:"promo_template_id"`
	PromoDesc             string    `gorm:"column:promo_desc" json:"promo_desc"`
	PromoTemplateStatusID int64     `gorm:"column:promo_template_status_id" json:"promo_template_status_id"`
	IsMultiplied          bool      `gorm:"column:is_multiplied" json:"is_multiplied"`
	IsBudgetReference     bool      `gorm:"column:is_budget_reference" json:"is_budget_reference"`
	BudgetReferenceType   int       `gorm:"column:budget_reference_type" json:"budget_reference_type"`
	BudgetReferenceID     int64     `gorm:"column:budget_reference_id" json:"budget_reference_id"`
	BudgetControlLevel    int64     `gorm:"column:budget_control_level" json:"budget_control_level"`
	BudgetAmount          float64   `gorm:"column:budget_amount" json:"budget_amount"`
	ExecutionLevel        int64     `gorm:"column:execution_level" json:"execution_level"`
	IsClaimable           bool      `gorm:"column:is_claimable" json:"is_claimable"`
	ClaimDays             int64     `gorm:"column:claim_days" json:"claim_days"`
	MaxDiscountType       int64     `gorm:"column:max_discount_type" json:"max_discount_type"`
	MaxDiscountOutlet     float64   `gorm:"column:max_discount_outlet" json:"max_discount_outlet"`
	MaxInvoiceOutlet      float64   `gorm:"column:max_invoice_outlet" json:"max_invoice_outlet"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy             string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy             string    `gorm:"column:updated_by" json:"updated_by"`
	MaxDiscountOutletUom  int       `gorm:"max_discount_outlet_uom" json:"max_discount_outlet_uom"`
}

func (PromoTemplate) TableName() string {
	return "sls.promo_templates"
}

func (m *PromoTemplate) BeforeCreate(trx *gorm.DB) (err error) {
	guid := xid.New()
	m.PromoTemplateID = guid.String()
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now().UTC()
	return nil
}
