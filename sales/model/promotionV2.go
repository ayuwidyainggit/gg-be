package model

import (
	"time"

	"gorm.io/gorm"
)

// ---------- Enum types that mirror Postgres enums ----------

type PromotionType string // promo.promotion_type
const (
	PromotionTypeSlab   PromotionType = "slab"
	PromotionTypeStrata PromotionType = "strata"
)

type CreationType string // promo.creation_type
const (
	CreationTypeNew         CreationType = "new"
	CreationTypeReplacement CreationType = "replacement"
)

type PromotionStatus string // promo.promotion_status
const (
	PromoStatusDraft    PromotionStatus = "draft"
	PromoStatusSubmit   PromotionStatus = "submit"
	PromoStatusApproved PromotionStatus = "approved"
	PromoStatusRejected PromotionStatus = "rejected"
	PromoStatusInactive PromotionStatus = "inactive"
	PromoStatusActive   PromotionStatus = "active"
	PromoStatusClosed   PromotionStatus = "closed"
)

type BudgetRefType string // promo.budget_ref_type
const (
	BudgetRefUnlimited BudgetRefType = "unlimited"
	BudgetRefLimited   BudgetRefType = "limited"
)

type ControlLevel string // promo.control_level
const (
	CtrlRegion      ControlLevel = "region"
	CtrlArea        ControlLevel = "area"
	CtrlDistributor ControlLevel = "distributor"
	CtrlSalesman    ControlLevel = "salesman"
)

type ClaimType string // promo.claim_type
const (
	ClaimFull    ClaimType = "full"
	ClaimPartial ClaimType = "partial"
)

type RewardCapType string // promo.reward_cap_type
const (
	CapAmount RewardCapType = "amount"
	CapQty    RewardCapType = "qty"
)

type CoverageType string // promo.coverage_type
const (
	CoverageNational      CoverageType = "national"
	CoverageByDistributor CoverageType = "by_distributor"
)

type RuleType string // promo.rule_type
const (
	RuleTypeQuantity RuleType = "quantity"
	RuleTypeValue    RuleType = "value"
)

type UomType string // promo.uom_type
const (
	UomTypeSmallest UomType = "smallest"
	UomTypeMiddle   UomType = "middle"
	UomTypeLargest  UomType = "largest"
)

type RewardType string // promo.reward_type
const (
	RewardTypePercentage RewardType = "percentage"
	RewardTypeFixedValue RewardType = "fixed_value"
	RewardTypeProduct    RewardType = "product"
)

type PerScope string // promo.per_scope
const (
	PerScopeProduct PerScope = "per_product"
	PerScopeOrder   PerScope = "per_order"
)

type OutletSelType string // promo.outlet_sel_type
const (
	OutletSelByAttribute OutletSelType = "by_attribute"
	OutletSelBySelection OutletSelType = "by_selection"
)

// ---------- Main model aligned with DDL ----------

type PromotionV2 struct {
	// Composite PK
	CustID            string `gorm:"column:cust_id;primaryKey;size:10;autoIncrement:false" json:"cust_id"`
	DistributorCustID string `gorm:"column:distributor_cust_id;size:10" json:"distributor_cust_id"`
	PromoID           string `gorm:"column:promo_id;primaryKey;size:20;autoIncrement:false" json:"promo_id"`

	// Distributor IDs
	DistributorIDs []int64 `gorm:"-" json:"distributor_ids"`
	// DistributorID     *int64 `gorm:"column:distributor_id" json:"distributor_id"`

	// Identity & lifecycle
	PromoDesc         string          `gorm:"column:promo_desc;size:100" json:"promo_desc"`
	PromoType         PromotionType   `gorm:"column:promo_type;type:promo.promotion_type" json:"promo_type"`
	PromoCreationType CreationType    `gorm:"column:promo_creation_type;type:promo.creation_type" json:"promo_creation_type"`
	ExistingPromoID   *string         `gorm:"column:existing_promo_id;size:20" json:"existing_promo_id,omitempty"` // FK pairs with (cust_id, existing_promo_id)
	PromoStatus       PromotionStatus `gorm:"column:promo_status;type:promo.promotion_status;default:draft" json:"promo_status"`

	// Flags
	SlabMultiplied   *bool `gorm:"column:slab_multiplied" json:"slab_multiplied,omitempty"`     // SLAB only
	StrataSequential *bool `gorm:"column:strata_sequential" json:"strata_sequential,omitempty"` // STRATA only

	// Budget settings
	IsBudgetReference  bool           `gorm:"column:is_budget_reference" json:"is_budget_reference"`
	BudgetRefType      *BudgetRefType `gorm:"column:budget_ref_type;type:promo.budget_ref_type" json:"budget_ref_type,omitempty"`
	BudgetReferenceID  *int32         `gorm:"column:budget_reference_id" json:"budget_reference_id,omitempty"`
	BudgetControlLevel *ControlLevel  `gorm:"column:budget_control_level;type:promo.control_level" json:"budget_control_level,omitempty"`
	BudgetAmount       float64        `gorm:"column:budget_amount;type:numeric(20,4)" json:"budget_amount"`
	ExecutionLevel     *ControlLevel  `gorm:"column:execution_level;type:promo.control_level" json:"execution_level"`

	// Coverage
	Coverage CoverageType `gorm:"column:coverage;type:promo.coverage_type;default:national" json:"coverage"`

	// Dates (DDL uses DATE; gorm will persist date-only with type:date)
	EffectiveFrom time.Time `gorm:"column:effective_from;type:date" json:"effective_from"`
	EffectiveTo   time.Time `gorm:"column:effective_to;type:date" json:"effective_to"`

	BudgetID      *string    `gorm:"column:budget_id;size:50" json:"budget_id,omitempty"`
	ClaimDateFrom *time.Time `gorm:"column:claim_date_from;type:date" json:"-"`
	ClaimDateTo   *time.Time `gorm:"column:claim_date_to;type:date" json:"-"`
	VatRate       *float64   `gorm:"column:vat_rate;type:numeric(5,2)" json:"vat_rate,omitempty"`
	WhtRate       *float64   `gorm:"column:wht_rate;type:numeric(5,2)" json:"wht_rate,omitempty"`

	// Claim settings
	IsClaimable         bool       `gorm:"column:is_claimable" json:"is_claimable"`
	ClaimType           *ClaimType `gorm:"column:claim_type;type:promo.claim_type" json:"claim_type,omitempty"`
	ClaimStartAfterDays *int32     `gorm:"column:claim_start_after_days" json:"claim_start_after_days,omitempty"`
	ClaimRealizationPct *float64   `gorm:"column:claim_realization_pct;type:numeric(5,2)" json:"claim_realization_pct,omitempty"`

	// Per-outlet caps / limits
	MaxTotalRewardType  *RewardCapType `gorm:"column:max_total_reward_type;type:promo.reward_cap_type" json:"max_total_reward_type,omitempty"`
	MaxTotalRewardValue *float64       `gorm:"column:max_total_reward_value;type:numeric(20,4)" json:"max_total_reward_value,omitempty"`
	MaxInvoicePerOutlet *float64       `gorm:"column:max_invoice_per_outlet;type:numeric(10,2)" json:"max_invoice_per_outlet,omitempty"`

	MinimumSKU int `gorm:"column:minimum_sku;not null;default:1" json:"minimum_sku"`

	// Audit
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CreatedBy string    `gorm:"column:created_by;size:150" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by;size:150" json:"updated_by"`
	Remarks   string    `gorm:"column:remarks;size:255" json:"remarks,omitempty"`

	// Budget tracking
	BudgetRealization float64  `gorm:"column:budget_realization;type:numeric(20,4)" json:"budget_realization"`
	RemainingBudget   *float64 `gorm:"column:remaining_budget;type:numeric(20,4);->" json:"remaining_budget,omitempty"` // generated always as ..., read-only
}

// Tell GORM the exact table (with schema).
func (PromotionV2) TableName() string { return "promo.promotions" }

// BeforeSave handles data conversion before saving to database
func (p *PromotionV2) BeforeSave(tx *gorm.DB) error {
	// Convert empty string to nil for ExistingPromoID
	if p.ExistingPromoID != nil && *p.ExistingPromoID == "" {
		p.ExistingPromoID = nil
	}

	if p.BudgetRefType != nil && *p.BudgetRefType == "" {
		p.BudgetRefType = nil
	}
	if p.BudgetControlLevel != nil && *p.BudgetControlLevel == "" {
		p.BudgetControlLevel = nil
	}
	if p.ExecutionLevel != nil && *p.ExecutionLevel == "" {
		p.ExecutionLevel = nil
	}
	if p.ClaimType != nil && *p.ClaimType == "" {
		p.ClaimType = nil
	}
	if p.MaxTotalRewardType != nil && *p.MaxTotalRewardType == "" {
		p.MaxTotalRewardType = nil
	}

	return nil
}

type PromotionV2Slabs struct {
	CustID string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	// Primary key
	ID string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`

	// Foreign key to promotions table
	PromoID string `gorm:"column:promo_id;size:50;not null" json:"promo_id"`

	// Slab ordering and description
	Ordinal     int     `gorm:"column:ordinal;not null" json:"ordinal"` // 1..N
	Description *string `gorm:"column:description;size:50" json:"description,omitempty"`

	// Rule configuration
	RuleType RuleType `gorm:"column:rule_type;type:promo.rule_type;not null" json:"rule_type"` // quantity/value
	RuleUom  *UomType `gorm:"column:rule_uom;type:promo.uom_type;" json:"rule_uom"`            // smallest/middle/largest

	// Range configuration
	RangeFrom *float64 `gorm:"column:range_from;type:numeric(20,4)" json:"range_from,omitempty"` // may be NULL when multiplied
	RangeTo   float64  `gorm:"column:range_to;type:numeric(20,4);not null" json:"range_to"`      // CHECK (range_to > COALESCE(range_from, 0))

	// Reward configuration
	RewardType  RewardType `gorm:"column:reward_type;type:promo.reward_type;not null" json:"reward_type"` // percentage/fixed_value/product
	RewardValue *float64   `gorm:"column:reward_value;type:numeric(20,4)" json:"reward_value,omitempty"`  // % or fixed amount; NULL if product
	RewardUom   *UomType   `gorm:"column:reward_uom;type:promo.uom_type;" json:"reward_uom"`              // smallest/middle/largest

	// Scope configuration
	PerScope *string `gorm:"column:per_scope;size:16" json:"per_scope"` // 'per_product'|'per_outlet' if fixed_value chosen

	// Business constraints (enforced at application level):
	// - NOT (reward_type = 'percentage') OR (reward_value IS NOT NULL AND reward_value BETWEEN 0 AND 100)
	// - NOT (reward_type = 'fixed_value') OR (reward_value IS NOT NULL AND reward_value >= 0)
	// - NOT (reward_type = 'product') OR reward_value IS NULL
	// - per_scope IS NULL OR per_scope IN ('per_product','per_outlet')
	// - UNIQUE(promo_id, ordinal)

	IsMultiplied bool `gorm:"column:is_multiplied;->" json:"is_multiplied"` // read-only
}

func (PromotionV2Slabs) TableName() string { return "promo.promotion_slabs" }

// BeforeSave handles data conversion before saving to database
func (p *PromotionV2Slabs) BeforeSave(tx *gorm.DB) error {
	// Convert empty string to nil
	if p.RuleUom != nil && *p.RuleUom == "" {
		p.RuleUom = nil
	}
	if p.RewardUom != nil && *p.RewardUom == "" {
		p.RewardUom = nil
	}
	if p.PerScope != nil && *p.PerScope == "" {
		p.PerScope = nil
	}

	return nil
}

type PromotionV2Strata struct {
	CustID string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	// Primary key
	ID string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`

	// Foreign key to promotions table
	PromoID string `gorm:"column:promo_id;size:50;not null" json:"promo_id"`

	// Strata ordering and description
	Ordinal     int     `gorm:"column:ordinal;not null" json:"ordinal"` // 1..5 (CHECK ordinal BETWEEN 1 AND 5)
	Description *string `gorm:"column:description;size:50" json:"description,omitempty"`

	// Rule configuration
	RuleType RuleType `gorm:"column:rule_type;type:promo.rule_type;not null" json:"rule_type"` // quantity/value
	RuleUom  *UomType `gorm:"column:rule_uom;type:promo.uom_type" json:"rule_uom"`             // smallest/middle/largest

	// Range configuration
	RangeFrom float64 `gorm:"column:range_from;type:numeric(20,4);not null" json:"range_from"` // CHECK (range_to > range_from)
	RangeTo   float64 `gorm:"column:range_to;type:numeric(20,4);not null" json:"range_to"`

	// Reward configuration
	RewardType  RewardType `gorm:"column:reward_type;type:promo.reward_type;not null" json:"reward_type"` // percentage/fixed_value/product
	RewardValue *float64   `gorm:"column:reward_value;type:numeric(20,4)" json:"reward_value,omitempty"`  // % or fixed amount; NULL if product
	RewardUom   *UomType   `gorm:"column:reward_uom;type:promo.uom_type" json:"reward_uom"`               // smallest/middle/largest
	PerScope    *string    `gorm:"column:per_scope;size:16" json:"per_scope"`                             // 'per_product'|'per_outlet' if fixed_value chosen

	// Claim configuration
	Claimable           bool     `gorm:"column:claimable;not null;default:false" json:"claimable"`                              // claimable per strata
	ClaimRealizationPct *float64 `gorm:"column:claim_realization_pct;type:numeric(5,2)" json:"claim_realization_pct,omitempty"` // only if claimable && header claim_type=partial

	// Business constraints (enforced at application level):
	// - NOT (reward_type = 'percentage') OR (reward_value IS NOT NULL AND reward_value BETWEEN 0 AND 100)
	// - NOT (reward_type = 'fixed_value') OR (reward_value IS NOT NULL AND reward_value >= 0)
	// - NOT (reward_type = 'product') OR reward_value IS NULL
	// - claim_realization_pct IS NULL OR (claim_realization_pct >= 0 AND claim_realization_pct <= 100)
	// - UNIQUE(promo_id, ordinal)
	// - ordinal BETWEEN 1 AND 5
}

type RewardContext struct {
	PromoID   string
	RewardUom *UomType
}

func (PromotionV2Strata) TableName() string { return "promo.promotion_strata" }

// BeforeSave handles data conversion before saving to database
func (p *PromotionV2Strata) BeforeSave(tx *gorm.DB) error {
	// log.Info("p:", structs.StructToJson(p))
	// Convert empty string to nil
	if p.RuleUom != nil && *p.RuleUom == "" {
		p.RuleUom = nil
	}
	if p.RewardUom != nil && *p.RewardUom == "" {
		p.RewardUom = nil
	}
	if p.PerScope != nil && *p.PerScope == "" {
		p.PerScope = nil
	}

	return nil
}

// runtime error: invalid memory address or nil pointer dereference

// ================================
// promo.promotion_product_criteria
// ================================
type PromotionProductCriteria struct {
	CustID    string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ProCode   string `gorm:"column:pro_code;size:50;->" json:"pro_code"`
	ProName   string `gorm:"column:pro_name;size:100;->" json:"pro_name"`
	ID        string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	PromoID   string `gorm:"column:promo_id;size:50;not null" json:"promo_id"`                   // FK to promo.promotions(promo_id)
	ProID     int64  `gorm:"column:pro_id;not null;uniqueIndex:uniq_promo_pro_id" json:"pro_id"` // part of UNIQUE(promo_id, pro_id)
	Mandatory bool   `gorm:"column:mandatory;not null;default:false" json:"mandatory"`

	MinBuyType  *RuleType `gorm:"column:min_buy_type;type:promo.rule_type" json:"min_buy_type,omitempty"` // nullable
	MinBuyQty   *float64  `gorm:"column:min_buy_qty;type:numeric(20,4)" json:"min_buy_qty,omitempty"`     // nullable
	MinBuyValue *float64  `gorm:"column:min_buy_value;type:numeric(20,4)" json:"min_buy_value,omitempty"` // nullable
	MinBuyUom   *UomType  `gorm:"column:min_buy_uom;type:promo.uom_type" json:"min_buy_uom,omitempty"`    // nullable

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PromotionProductCriteria) TableName() string { return "promo.promotion_product_criteria" }

// ================================
// promo.promotion_reward_products
// ================================
type PromotionRewardProduct struct {
	CustID    string  `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ProCode   string  `gorm:"column:pro_code;size:50;->" json:"pro_code"`  // read-only
	ProName   string  `gorm:"column:pro_name;size:100;->" json:"pro_name"` // read-only
	ID        string  `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	PromoID   string  `gorm:"column:promo_id;size:50;not null;index:idx_promo_reward_products_promo" json:"promo_id"` // FK to promo.promotions(promo_id)
	ProID     int64   `gorm:"column:pro_id;not null;uniqueIndex:uniq_promo_pro_id" json:"pro_id"`                     // part of UNIQUE(promo_id, pro_id)
	Ordinal   int32   `gorm:"column:ordinal;not null;uniqueIndex:uniq_promo_ordinal" json:"ordinal"`                  // part of UNIQUE(promo_id, ordinal)
	UnitId    string  `gorm:"column:unit_id;->" json:"unit_id"`                                                       // read-only
	QtyStock  int64   `gorm:"column:qty_stock;->" json:"qty_stock"`                                                   // read-only
	ConvUnit2 float64 `gorm:"column:conv_unit2;->" json:"conv_unit2"`                                                 // read-only
	ConvUnit3 float64 `gorm:"column:conv_unit3;->" json:"conv_unit3"`                                                 // read-only
}

func (PromotionRewardProduct) TableName() string { return "promo.promotion_reward_products" }

// ===================================
// promo.promotion_coverage_distributors
// ===================================
type PromotionCoverageDistributors struct {
	CustID          string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID              string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	PromoID         string `gorm:"column:promo_id;size:50;not null;index:idx_promo_cov_dist_promo" json:"promo_id"`         // FK to promo.promotions(promo_id)
	DistributorID   int64  `gorm:"column:distributor_id;not null;uniqueIndex:uniq_promo_distributor" json:"distributor_id"` // part of UNIQUE(promo_id, distributor_id)
	DistributorCode string `gorm:"column:distributor_code;size:50;->" json:"distributor_code"`
	DistributorName string `gorm:"column:distributor_name;size:100;->" json:"distributor_name"`
}

func (PromotionCoverageDistributors) TableName() string {
	return "promo.promotion_coverage_distributors"
}

// ================================
// promo.promotion_outlet_criteria
// ================================
type PromotionOutletCriteria struct {
	CustID        string        `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID            string        `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	PromoID       string        `gorm:"column:promo_id;size:50;not null" json:"promo_id"` // FK to promo.promotions(promo_id)
	SelectionType OutletSelType `gorm:"column:selection_type;type:promo.outlet_sel_type;default:by_attribute" json:"selection_type"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Optional: preload related data (not part of table schema)
	SelectedOutlets     []PromotionOutletsSelected          `gorm:"foreignKey:CriteriaID;references:ID" json:"selected_outlets,omitempty"`
	AttributeTypes      []PromotionOutletAttributeType      `gorm:"foreignKey:CriteriaID;references:ID" json:"attribute_types,omitempty"`
	AttributeSalesTeams []PromotionOutletAttributeSalesTeam `gorm:"foreignKey:CriteriaID;references:ID" json:"attribute_sales_teams,omitempty"`
	AttributeGroups     []PromotionOutletAttributeGroup     `gorm:"foreignKey:CriteriaID;references:ID" json:"attribute_groups,omitempty"`
	AttributeClasses    []PromotionOutletAttributeClass     `gorm:"foreignKey:CriteriaID;references:ID" json:"attribute_classes,omitempty"`
}

func (PromotionOutletCriteria) TableName() string { return "promo.promotion_outlet_criteria" }

// ================================
// promo.promotion_outlets_selected
// ================================
type PromotionOutletsSelected struct {
	CustID     string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID         string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	CriteriaID string `gorm:"column:criteria_id;size:30;not null;index:idx_promo_outlet_sel_crit;uniqueIndex:uniq_criteria_outlet" json:"criteria_id"` // FK to promo.promotion_outlet_criteria(id)
	OutletID          int64  `gorm:"column:outlet_id;not null;uniqueIndex:uniq_criteria_outlet" json:"outlet_id"`
	OutletCode        string `gorm:"column:outlet_code;size:50;->" json:"outlet_code"`
	OutletName        string `gorm:"column:outlet_name;size:100;->" json:"outlet_name"`
	DistributorCode   string `gorm:"column:distributor_code;size:50;->" json:"distributor_code"`
	DistributorName   string `gorm:"column:distributor_name;size:100;->" json:"distributor_name"`
	// part of UNIQUE(criteria_id, outlet_id)
}

func (PromotionOutletsSelected) TableName() string { return "promo.promotion_outlets_selected" }

// ===================================
// promo.promotion_outlet_attribute_type
// ===================================
type PromotionOutletAttributeType struct {
	CustID         string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID             string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	CriteriaID     string `gorm:"column:criteria_id;size:30;not null;uniqueIndex:uniq_criteria_outlet_type" json:"criteria_id"` // FK to promo.promotion_outlet_criteria(id)
	OutletTypeID   int64  `gorm:"column:outlet_type_id;not null;uniqueIndex:uniq_criteria_outlet_type" json:"outlet_type_id"`   // part of UNIQUE(criteria_id, outlet_type_id)
	OutletTypeCode string `gorm:"column:outlet_type_code;size:50;->" json:"outlet_type_code"`
	OutletTypeName string `gorm:"column:outlet_type_name;size:100;->" json:"outlet_type_name"`
}

func (PromotionOutletAttributeType) TableName() string {
	return "promo.promotion_outlet_attribute_type"
}

// ==========================================
// promo.promotion_outlet_attribute_sales_team
// ==========================================
type PromotionOutletAttributeSalesTeam struct {
	CustID        string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID            string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	CriteriaID    string `gorm:"column:criteria_id;size:30;not null;uniqueIndex:uniq_criteria_sales_team" json:"criteria_id"` // FK to promo.promotion_outlet_criteria(id)
	SalesTeamID   int64  `gorm:"column:sales_team_id;not null;uniqueIndex:uniq_criteria_sales_team" json:"sales_team_id"`     // part of UNIQUE(criteria_id, sales_team_id)
	SalesTeamCode string `gorm:"column:sales_team_code;size:50;->" json:"sales_team_code"`
	SalesTeamName string `gorm:"column:sales_team_name;size:100;->" json:"sales_team_name"`
}

func (PromotionOutletAttributeSalesTeam) TableName() string {
	return "promo.promotion_outlet_attribute_sales_team"
}

// ====================================
// promo.promotion_outlet_attribute_group
// ====================================
type PromotionOutletAttributeGroup struct {
	CustID          string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID              string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	CriteriaID      string `gorm:"column:criteria_id;size:30;not null;uniqueIndex:uniq_criteria_outlet_group" json:"criteria_id"` // FK to promo.promotion_outlet_criteria(id)
	OutletGroupID   int64  `gorm:"column:outlet_group_id;not null;uniqueIndex:uniq_criteria_outlet_group" json:"outlet_group_id"` // part of UNIQUE(criteria_id, outlet_group_id)
	OutletGroupCode string `gorm:"column:outlet_group_code;size:50;->" json:"outlet_group_code"`
	OutletGroupName string `gorm:"column:outlet_group_name;size:100;->" json:"outlet_group_name"`
}

func (PromotionOutletAttributeGroup) TableName() string {
	return "promo.promotion_outlet_attribute_group"
}

// ====================================
// promo.promotion_outlet_attribute_class
// ====================================
type PromotionOutletAttributeClass struct {
	CustID          string `gorm:"column:cust_id;size:10;not null" json:"cust_id"`
	ID              string `gorm:"column:id;primaryKey;size:30;autoIncrement:true" json:"id"`
	CriteriaID      string `gorm:"column:criteria_id;size:30;not null;uniqueIndex:uniq_criteria_outlet_class" json:"criteria_id"` // FK to promo.promotion_outlet_criteria(id)
	OutletClassID   int64  `gorm:"column:outlet_class_id;not null;uniqueIndex:uniq_criteria_outlet_class" json:"outlet_class_id"` // part of UNIQUE(criteria_id, outlet_class_id)
	OutletClassCode string `gorm:"column:outlet_class_code;size:50;->" json:"outlet_class_code"`
	OutletClassName string `gorm:"column:outlet_class_name;size:100;->" json:"outlet_class_name"`
}

func (PromotionOutletAttributeClass) TableName() string {
	return "promo.promotion_outlet_attribute_class"
}

type OutletPromo struct {
	OutletID   int    `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName string `gorm:"column:outlet_name" json:"outlet_name"`
	OtGrpID    int    `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	OtClassID  int    `gorm:"column:ot_class_id" json:"ot_class_id"`
	OtTypeID   int    `gorm:"column:ot_type_id" json:"ot_type_id"`
}

func (OutletPromo) TableName() string {
	return "mst.m_outlet"
}

type SalesmanPromo struct {
	SalesmanId   int    `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName string `gorm:"column:salesman_name" json:"salesman_name"`
	WhId         int    `gorm:"column:wh_id" json:"wh_id"`
	SalesTeamId  int    `gorm:"column:sales_team_id" json:"sales_team_id"`
}

func (SalesmanPromo) TableName() string {
	return "mst.m_salesman"
}

type WarehousePromo struct {
	WhID   int    `gorm:"column:wh_id" json:"wh_id"`
	WhCode string `gorm:"column:wh_code" json:"wh_code"`
	WhName string `gorm:"column:wh_name" json:"wh_name"`
}

func (WarehousePromo) TableName() string {
	return "mst.m_warehouse"
}
