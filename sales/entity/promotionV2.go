package entity

import (
	"fmt"
	"strings"
	"unicode"
)

// --- Enum types (strings to match Postgres enums) ---

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

type PromotionV2Status string // promo.promotion_status
const (
	PromoStatusDraft    PromotionV2Status = "draft"
	PromoStatusSubmit   PromotionV2Status = "submit"
	PromoStatusApproved PromotionV2Status = "approved"
	PromoStatusRejected PromotionV2Status = "rejected"
	PromoStatusInactive PromotionV2Status = "inactive"
	PromoStatusActive   PromotionV2Status = "active"
	PromoStatusClosed   PromotionV2Status = "closed"
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

type RewardType string // promo.reward_type
const (
	RewardTypePercentage RewardType = "percentage"
	RewardTypeFixedValue RewardType = "fixed_value"
	RewardTypeProduct    RewardType = "product"
)

type UOMType string // promo.uom_type (smallest/middle/largest)
const (
	UOMSmallest UOMType = "smallest"
	UOMMiddle   UOMType = "middle"
	UOMLargest  UOMType = "largest"
)

// --- End of Enum types (strings to match Postgres enums) ---

// --- Request body aligned to DDL column names ---

type CreatePromotionV2Body struct {
	// Keys / ownership
	CustID            string `json:"cust_id" validate:"required,max=10,alphanum"`
	ParentCustID      string `json:"parent_cust_id" validate:"omitempty,max=10,alphanum"`
	DistributorCustID string `json:"distributor_cust_id" validate:"omitempty,max=10,alphanum"`

	// Identity
	PromoID   string `json:"promo_id" validate:"required,max=50"`
	PromoDesc string `json:"promo_desc" validate:"required,max=100"`

	// Types & lifecycle (string enums to match DB)
	PromoType         PromotionType     `json:"promo_type" validate:"required,oneof=slab strata"`
	PromoCreationType CreationType      `json:"promo_creation_type" validate:"required,oneof=new replacement"`
	ExistingPromoID   string            `json:"existing_promo_id" validate:"omitempty,max=50"`
	PromoStatus       PromotionV2Status `json:"promo_status" validate:"required,oneof=draft submit approved rejected inactive active closed"`

	// Multipliers / global flags
	SlabMultiplied   *bool `json:"slab_multiplied,omitempty"`   // applies to SLAB only (optional)
	StrataSequential *bool `json:"strata_sequential,omitempty"` // applies to STRATA only (optional)

	// Budget settings
	IsBudgetReference  bool          `json:"is_budget_reference"`
	BudgetRefType      BudgetRefType `json:"budget_ref_type" validate:"omitempty,oneof=unlimited limited"`
	BudgetReferenceID  int64         `json:"budget_reference_id"`
	BudgetControlLevel ControlLevel  `json:"budget_control_level" validate:"omitempty,oneof=region area distributor salesman"`
	BudgetAmount       float64       `json:"budget_amount" validate:"omitempty,gte=0"`
	ExecutionLevel     ControlLevel  `json:"execution_level" validate:"omitempty,oneof=region area distributor salesman"`

	// Coverage
	Coverage CoverageType `json:"coverage" validate:"omitempty,oneof=national by_distributor"`

	// Date range (keep as YYYY-MM-DD strings, DB column is DATE)
	EffectiveFrom string `json:"effective_from" validate:"required,datetime=2006-01-02"`
	EffectiveTo   string `json:"effective_to"   validate:"required,datetime=2006-01-02"`

	// Budget id & claim period (optional; claim dates accept DD/MM/YYYY or YYYY-MM-DD)
	BudgetID      string   `json:"budget_id" validate:"omitempty,max=50,budgetID"`
	ClaimDateFrom string   `json:"claim_date_from" validate:"omitempty"`
	ClaimDateTo   string   `json:"claim_date_to" validate:"omitempty"`
	VatRate       *float64 `json:"vat_rate,omitempty" validate:"omitempty,gte=0,lte=100"`
	WhtRate       *float64 `json:"wht_rate,omitempty" validate:"omitempty,gte=0,lte=100"`

	// Claim settings
	IsClaimable         bool      `json:"is_claimable"`
	ClaimType           ClaimType `json:"claim_type" validate:"omitempty,oneof=full partial"`
	ClaimStartAfterDays int       `json:"claim_start_after_days" validate:"omitempty,min=0,max=9999"`
	ClaimRealizationPct float64   `json:"claim_realization_pct" validate:"omitempty,gte=0,lte=100"`

	// Per-outlet limits / caps
	MaxTotalRewardType  RewardCapType `json:"max_total_reward_type" validate:"omitempty,oneof=amount qty"`
	MaxTotalRewardValue float64       `json:"max_total_reward_value" validate:"omitempty,gte=0"`           // was MaxDiscountOutlet
	MaxInvoicePerOutlet float64       `json:"max_invoice_per_outlet" validate:"omitempty,gte=0,lte=99999"` // was MaxInvoiceOutlet

	// Audit / misc (timestamps are DB-managed; keep who did it)
	CreatedBy string `json:"created_by" validate:"omitempty,max=150"`
	UpdatedBy string `json:"updated_by" validate:"omitempty,max=150"`
	Remarks   string `json:"remarks" validate:"omitempty,max=255"`

	// Related detail payloads (map to child tables; keep as-is)
	Slabs  []PromoSlabItem   `json:"slabs,omitempty" validate:"dive"`
	Strata []PromoStrataItem `json:"strata,omitempty" validate:"dive"`

	// product criteria
	MinimumSKU           int                                  `json:"minimum_sku"`
	ProductCriteria      []CreatePromotionProductCriteria     `json:"product_criteria" validate:"required,dive"` // optional embedded items
	RewardProducts       []CreatePromotionRewardProduct       `json:"reward_products,omitempty" validate:"dive"`
	CoverageDistributors []CreatePromotionCoverageDistributor `json:"coverage_distributors,omitempty" validate:"dive"`
	OutletCriteria       CreatePromotionOutletCriteria        `json:"outlet_criteria" validate:"required,dive"`
}

type CreatePromotionCoverageDistributor struct {
	DistributorID int64 `json:"distributor_id" validate:"required"`
}

type CreatePromotionOutletCriteria struct {
	SelectionType  string  `json:"selection_type" validate:"required,oneof=by_outlet by_attribute"`
	OutletIDs      []int64 `json:"outlet_ids,omitempty" validate:"omitempty,dive,min=1"` // for by_outlet selection
	OutletTypeIDs  []int64 `json:"outlet_type_ids,omitempty" validate:"omitempty"`
	OutletGroupIDs []int64 `json:"outlet_group_ids,omitempty" validate:"omitempty"`
	OutletClassIDs []int64 `json:"outlet_class_ids,omitempty" validate:"omitempty"`
	SalesTeamIDs   []int64 `json:"sales_team_ids,omitempty" validate:"omitempty"`
}

type CreatePromotionOutletAttributeType struct {
	OutletTypeID int64 `json:"outlet_type_id" validate:"required"`
}

type CreatePromotionOutletAttributeGroup struct {
	OutletGroupID int64 `json:"outlet_group_id" validate:"required"`
}

type CreatePromotionOutletAttributeClass struct {
	OutletClassID int64 `json:"outlet_class_id" validate:"required"`
}

type CreatePromotionOutletAttributeSalesTeam struct {
	SalesTeamID int64 `json:"sales_team_id" validate:"required"`
}

type PromoSlabItem struct {
	CustID      string `json:"cust_id,omitempty" validate:"required,max=10"`
	PromoID     string `json:"promo_id,omitempty" validate:"required,max=50"`
	ProCode     string `json:"pro_code,omitempty"`
	ProName     string `json:"pro_name,omitempty"`
	Ordinal     int    `json:"ordinal" validate:"required,min=1"`
	Description string `json:"description,omitempty" validate:"omitempty,max=50"`

	RuleType  RuleType `json:"rule_type" validate:"required,oneof=quantity value"`
	RuleUom   UOMType  `json:"rule_uom" validate:"omitempty,oneof=smallest middle largest"`
	RangeFrom *float64 `json:"range_from"` // nullable (wajib NULL jika slab_multiplied = true)
	RangeTo   float64  `json:"range_to" validate:"required"`

	RewardType  RewardType `json:"reward_type" validate:"required,oneof=percentage fixed_value product"`
	RewardValue float64    `json:"reward_value" validate:"required"` // nullable: wajib diisi utk percentage/fixed_value; NULL utk product
	RewardUom   UOMType    `json:"reward_uom" validate:"omitempty,oneof=smallest middle largest"`

	// Hanya valid utk reward_type = fixed_value
	PerScope string `json:"per_scope" validate:"omitempty,oneof=per_product per_order"`
}

type PromoStrataItem struct {
	CustID      string `json:"cust_id,omitempty" validate:"required,max=10"`
	PromoID     string `json:"promo_id,omitempty" validate:"required,max=50"`
	Ordinal     int    `json:"ordinal" validate:"required,min=1,max=5"`
	Description string `json:"description,omitempty" validate:"omitempty,max=50"`

	RuleType RuleType `json:"rule_type" validate:"required,oneof=quantity value"`
	RuleUOM  UOMType  `json:"rule_uom" validate:"omitempty,oneof=smallest middle largest"`

	RangeFrom float64 `json:"range_from" validate:"required"`
	RangeTo   float64 `json:"range_to" validate:"required"` // NOTE: DB enforce range_to > range_from

	RewardType  RewardType `json:"reward_type" validate:"required,oneof=percentage fixed_value product"`
	RewardValue float64    `json:"reward_value" validate:"required"` // nullable (NULL utk product; 0..100 utk percentage)
	RewardUom   UOMType    `json:"reward_uom" validate:"omitempty,oneof=smallest middle largest"`

	// Hanya valid utk reward_type = fixed_value
	PerScope string `json:"per_scope,omitempty" validate:"omitempty,oneof=per_product per_order"`

	Claimable           bool     `json:"claimable"`
	ClaimRealizationPct *float64 `json:"claim_realization_pct,omitempty" validate:"omitempty,gte=0,lte=100"`
}

// ==== End of STRATA payload (array) ====

type DetailPromotionV2Params struct {
	PromoID      string `params:"promo_id" validate:"required,max=30"`
	CustID       string
	ParentCustId string
}

type DeletePromotionV2Params struct {
	PromoID string `params:"promo_id" validate:"required,max=30"`
}

type UpdatePromotionV2Params struct {
	PromoID      string `params:"promo_id" validate:"required,max=30"`
	CustID       string
	ParentCustId string
}

type UpdatePromotionV2StatusParams struct {
	PromoID      string            `params:"promo_id" validate:"required,max=30"`
	PromoStatus  PromotionV2Status `params:"promo_status" validate:"required,oneof=draft submit approved rejected inactive active closed"`
	CustID       string
	ParentCustId string
}

type PromotionV2 struct {
	// Keys / ownership
	CustID            string `json:"cust_id"`
	DistributorCustID string `json:"distributor_cust_id,omitempty"`

	// Distributor ID
	DistributorIDs []int64 `json:"distributor_ids"`
	// DistributorID     *int64 `json:"distributor_id"`

	// Identity
	PromoID   string `json:"promo_id"`
	PromoDesc string `json:"promo_desc"`

	// Types & lifecycle (string enums to match DB)
	PromoType         PromotionType     `json:"promo_type"`          // slab | strata
	PromoCreationType CreationType      `json:"promo_creation_type"` // new | replacement
	ExistingPromoID   string            `json:"existing_promo_id"`
	PromoStatus       PromotionV2Status `json:"promo_status"` // draft|submit|approved|rejected|inactive|active|closed

	// Multipliers / global flags
	SlabMultiplied   *bool `json:"slab_multiplied,omitempty"`   // SLAB only
	StrataSequential *bool `json:"strata_sequential,omitempty"` // STRATA only

	// Budget settings
	IsBudgetReference  bool           `json:"is_budget_reference"`
	BudgetRefType      *BudgetRefType `json:"budget_ref_type,omitempty"` // unlimited | limited
	BudgetReferenceID  int64          `json:"budget_reference_id,omitempty"`
	BudgetControlLevel *ControlLevel  `json:"budget_control_level,omitempty"` // region|area|distributor|salesman
	BudgetAmount       float64        `json:"budget_amount,omitempty"`
	ExecutionLevel     ControlLevel   `json:"execution_level"` // region|area|distributor|salesman

	// Coverage
	Coverage CoverageType `json:"coverage,omitempty"` // national | by_distributor

	// Date range (YYYY-MM-DD strings)
	EffectiveFrom string `json:"effective_from"`
	EffectiveTo   string `json:"effective_to"`

	// Budget id & claim period (optional; claim dates as DD/MM/YYYY in responses)
	BudgetID      string   `json:"budget_id"`
	ClaimDateFrom string   `json:"claim_date_from"`
	ClaimDateTo   string   `json:"claim_date_to"`
	VatRate       *float64 `json:"vat_rate,omitempty"`
	WhtRate       *float64 `json:"wht_rate,omitempty"`

	// Claim settings
	IsClaimable         bool       `json:"is_claimable"`
	ClaimType           *ClaimType `json:"claim_type,omitempty"` // full | partial
	ClaimStartAfterDays *int       `json:"claim_start_after_days,omitempty"`
	ClaimRealizationPct *float64   `json:"claim_realization_pct,omitempty"` // 0..100

	// Per-outlet limits / caps
	MaxTotalRewardType  *RewardCapType `json:"max_total_reward_type,omitempty"` // amount | qty
	MaxTotalRewardValue *float64       `json:"max_total_reward_value,omitempty"`
	MaxInvoicePerOutlet *float64       `json:"max_invoice_per_outlet,omitempty"`

	// Details
	Slabs  []PromoSlabItem   `json:"slabs,omitempty"`
	Strata []PromoStrataItem `json:"strata,omitempty"`

	MinimumSKU      int                      `json:"minimum_sku"`
	ProductCriteria []PromoProductCriteria   `json:"product_criteria"`
	RewardProducts  []PromotionRewardProduct `json:"reward_products,omitempty"`

	CoverageDistributors []PromoCoverageDistributor `json:"coverage_distributors"`
	OutletCriteria       PromoOutletCriteria        `json:"outlet_criteria"`

	// Optional: DB-derived tracking (present in DDL, not in create body)
	BudgetRealization *float64 `json:"budget_realization,omitempty"`
	RemainingBudget   *float64 `json:"remaining_budget,omitempty"`

	// Audit / misc
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	Remarks   string `json:"remarks,omitempty"`
}

type PromoOutletCriteria struct {
	CustID                string                         `json:"cust_id,omitempty"`
	PromoID               string                         `json:"promo_id,omitempty"`
	SelectionType         string                         `json:"selection_type"`
	SelectedOutlets       []PromoOutletSelected          `json:"selected_outlets,omitempty"` // for by_outlet selection
	SelectedOutletTypes   []PromoOutletTypeSelected      `json:"selected_outlet_types,omitempty"`
	SelectedOutletGroups  []PromoOutletGroupSelected     `json:"selected_outlet_groups,omitempty"`
	SelectedOutletClasses []PromoOutletClassSelected     `json:"selected_outlet_classes,omitempty"`
	SelectedSalesTeams    []PromoOutletSalesTeamSelected `json:"selected_sales_teams,omitempty"`
}

type PromoOutletSelected struct {
	OutletID          int64  `json:"outlet_id"`
	OutletCode        string `json:"outlet_code"`
	OutletName        string `json:"outlet_name"`
	DistributorCode   string `json:"distributor_code"`
	DistributorName   string `json:"distributor_name"`
}

type PromoOutletTypeSelected struct {
	OutletTypeID   int64  `json:"outlet_type_id"`
	OutletTypeCode string `json:"outlet_type_code"`
	OutletTypeName string `json:"outlet_type_name"`
}

type PromoOutletGroupSelected struct {
	OutletGroupID   int64  `json:"outlet_group_id"`
	OutletGroupCode string `json:"outlet_group_code"`
	OutletGroupName string `json:"outlet_group_name"`
}

type PromoOutletClassSelected struct {
	OutletClassID   int64  `json:"outlet_class_id"`
	OutletClassCode string `json:"outlet_class_code"`
	OutletClassName string `json:"outlet_class_name"`
}

type PromoOutletSalesTeamSelected struct {
	SalesTeamID   int64  `json:"sales_team_id"`
	SalesTeamCode string `json:"sales_team_code"`
	SalesTeamName string `json:"sales_team_name"`
}

type PromoCoverageDistributor struct {
	CustID          string `json:"cust_id,omitempty"`
	PromoID         string `json:"promo_id,omitempty"`
	DistributorID   int64  `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
}

type PromotionV2QueryFilter struct {
	PromoStatus   []string `query:"promo_status"`
	CustID        string
	ParentCustID  string
	EffectiveFrom *int64 `query:"effective_from" validate:"required_with=EffectiveTo,omitempty,gte=1000000000"`
	EffectiveTo   *int64 `query:"effective_to" validate:"required_with=EffectiveFrom,omitempty,lte=9999999999,gtefield=EffectiveFrom"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	PromoID       string `query:"promo_id"`
	PromoDesc     string `query:"promo_desc"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
	DistributorID int64  `query:"distributor_id"`
	TokenDistID   int64
}

// For embedded items (no criteria_id yet — server will attach them to the new criteria)
type CreatePromotionProductCriteria struct {
	ProID       int64     `json:"pro_id" validate:"required"` // bigint
	Mandatory   bool      `json:"mandatory"`
	MinBuyType  *RuleType `json:"min_buy_type,omitempty" validate:"omitempty,oneof=quantity value"`
	MinBuyQty   *float64  `json:"min_buy_qty,omitempty" validate:"omitempty,gte=0"`
	MinBuyValue *float64  `json:"min_buy_value,omitempty" validate:"omitempty,gte=0"`
	MinBuyUom   *UOMType  `json:"min_buy_uom,omitempty" validate:"omitempty,oneof=smallest middle largest"`
}

type PromoProductCriteria struct {
	CustID      string    `json:"cust_id,omitempty"`
	PromoID     string    `json:"promo_id,omitempty"`
	ProID       int64     `json:"pro_id"`
	ProCode     string    `json:"pro_code"`
	ProName     string    `json:"pro_name"`
	Mandatory   bool      `json:"mandatory"`
	MinBuyType  *RuleType `json:"min_buy_type,omitempty"`
	MinBuyQty   *float64  `json:"min_buy_qty,omitempty"`
	MinBuyValue *float64  `json:"min_buy_value,omitempty"`
	MinBuyUom   *UOMType  `json:"min_buy_uom,omitempty"`
}

type CreatePromotionRewardProduct struct {
	ProID   int64 `json:"pro_id" validate:"required"`
	Ordinal int32 `json:"ordinal" validate:"required,min=1"`
}

type PromotionRewardProduct struct {
	CustID  string `json:"cust_id,omitempty"`
	PromoID string `json:"promo_id,omitempty"`
	ProID   int64  `json:"pro_id" validate:"required"`
	ProCode string `json:"pro_code"`
	ProName string `json:"pro_name"`
	Ordinal int32  `json:"ordinal" validate:"required,min=1"`
}

type UpdatePromotionV2Body struct {
	// Keys / ownership
	PromoID           string `json:"promo_id" validate:"required,max=50"`
	CustID            string `json:"cust_id" validate:"required,max=10,alphanum"`
	ParentCustID      string `json:"parent_cust_id" validate:"omitempty,max=10,alphanum"`
	DistributorCustID string `json:"distributor_cust_id" validate:"omitempty,max=10,alphanum"`

	// Identity
	PromoDesc string `json:"promo_desc" validate:"required,max=100"`

	// Types & lifecycle (string enums to match DB)
	PromoType         PromotionType     `json:"promo_type" validate:"required,oneof=slab strata"`
	PromoCreationType CreationType      `json:"promo_creation_type" validate:"required,oneof=new replacement"`
	ExistingPromoID   string            `json:"existing_promo_id" validate:"omitempty,max=50"`
	PromoStatus       PromotionV2Status `json:"promo_status" validate:"required,oneof=draft submit approved rejected inactive active closed"`

	// Multipliers / global flags
	SlabMultiplied   *bool `json:"slab_multiplied,omitempty"`   // applies to SLAB only (optional)
	StrataSequential *bool `json:"strata_sequential,omitempty"` // applies to STRATA only (optional)

	// Budget settings
	IsBudgetReference  bool          `json:"is_budget_reference"`
	BudgetRefType      BudgetRefType `json:"budget_ref_type" validate:"omitempty,oneof=unlimited limited"`
	BudgetReferenceID  int64         `json:"budget_reference_id"`
	BudgetControlLevel ControlLevel  `json:"budget_control_level" validate:"omitempty,oneof=region area distributor salesman"`
	BudgetAmount       float64       `json:"budget_amount" validate:"omitempty,gte=0"`
	ExecutionLevel     ControlLevel  `json:"execution_level" validate:"omitempty,oneof=region area distributor salesman"`

	// Coverage
	Coverage CoverageType `json:"coverage" validate:"omitempty,oneof=national by_distributor"`

	// Date range (keep as YYYY-MM-DD strings, DB column is DATE)
	EffectiveFrom string `json:"effective_from" validate:"required,datetime=2006-01-02"`
	EffectiveTo   string `json:"effective_to"   validate:"required,datetime=2006-01-02"`

	// Budget id & claim period (optional; claim dates accept DD/MM/YYYY or YYYY-MM-DD)
	BudgetID      string   `json:"budget_id" validate:"omitempty,max=50,budgetID"`
	ClaimDateFrom string   `json:"claim_date_from" validate:"omitempty"`
	ClaimDateTo   string   `json:"claim_date_to" validate:"omitempty"`
	VatRate       *float64 `json:"vat_rate,omitempty" validate:"omitempty,gte=0,lte=100"`
	WhtRate       *float64 `json:"wht_rate,omitempty" validate:"omitempty,gte=0,lte=100"`

	// Claim settings
	IsClaimable         bool      `json:"is_claimable"`
	ClaimType           ClaimType `json:"claim_type" validate:"omitempty,oneof=full partial"`
	ClaimStartAfterDays int       `json:"claim_start_after_days" validate:"omitempty,min=0,max=9999"`
	ClaimRealizationPct float64   `json:"claim_realization_pct" validate:"omitempty,gte=0,lte=100"`

	// Per-outlet limits / caps
	MaxTotalRewardType  RewardCapType `json:"max_total_reward_type" validate:"omitempty,oneof=amount qty"`
	MaxTotalRewardValue float64       `json:"max_total_reward_value" validate:"omitempty,gte=0"`           // was MaxDiscountOutlet
	MaxInvoicePerOutlet float64       `json:"max_invoice_per_outlet" validate:"omitempty,gte=0,lte=99999"` // was MaxInvoiceOutlet

	// Audit / misc (timestamps are DB-managed; keep who did it)
	UpdatedBy string `json:"updated_by" validate:"omitempty,max=150"`
	Remarks   string `json:"remarks" validate:"omitempty,max=255"`

	// Related detail payloads (map to child tables; keep as-is)
	Slabs  []PromoSlabItem   `json:"slabs,omitempty" validate:"dive"`
	Strata []PromoStrataItem `json:"strata,omitempty" validate:"dive"`

	// product criteria
	MinimumSKU           int                                  `json:"minimum_sku"`
	ProductCriteria      []CreatePromotionProductCriteria     `json:"product_criteria" validate:"required,dive"` // optional embedded items
	RewardProducts       []CreatePromotionRewardProduct       `json:"reward_products,omitempty" validate:"dive"`
	CoverageDistributors []CreatePromotionCoverageDistributor `json:"coverage_distributors,omitempty" validate:"dive"`
	OutletCriteria       CreatePromotionOutletCriteria        `json:"outlet_criteria" validate:"required,dive"`

	// Reward type for validation
	RewardType RewardType `json:"reward_type" validate:"required,oneof=percentage fixed_value product"`
}

type UpdateStatusPromotionV2Body struct {
	CustID       string `json:"cust_id" validate:""`
	ParentCustID string `json:"parent_cust_id" validate:""`
	Remarks      string `json:"remarks" validate:"max:255"`
	EffectiveTo  string `json:"effective_to" validate:"omitempty,datetime=2006-01-02"`
}

type ConsultPromoV2Req struct {
	CustID        string          `json:"cust_id" validate:"required"`
	ParentCustID  string          `json:"parent_cust_id" validate:"required"`
	OrderDate     string          `json:"order_date" validate:"required"`
	OutletID      int             `json:"outlet_id" validate:"required"`
	SalesmanID    int             `json:"salesman_id" validate:"required"`
	WhID          int             `json:"wh_id" validate:"required"`
	Details       []ConPromoV2Det `json:"details" validate:"required"`
	PromoIDs      []string        `json:"promo_id,omitempty"`
	DistributorID int64
}

type ConPromoV2Det struct {
	ProID      int     `json:"pro_id"`
	Qty1       float64 `json:"qty1"`
	Qty2       float64 `json:"qty2"`
	Qty3       float64 `json:"qty3"`
	Total      float64 `json:"total"`
	GrossValue int     `json:"gross_value"`
	SubTotal   int     `json:"sub_total"`
}

type ConsultPromoResp struct {
	PromoID        string     `json:"promo_id"`
	PromoDesc      string     `json:"promo_desc"`
	SlabID         string     `json:"slab_id,omitempty"`
	SlabDesc       string     `json:"slab_desc,omitempty"`
	SlabReward     float64    `json:"slab_reward,omitempty"`
	SlabRuleType   RuleType   `json:"slab_rule_type,omitempty"`
	SlabRuleUom    UOMType    `json:"slab_rule_uom,omitempty"`
	SlabRewardUom  UOMType    `json:"slab_reward_uom,omitempty"`
	SlabRewardType RewardType `json:"slab_reward_type,omitempty"`
	SlabPerScope   string     `json:"slab_per_scope,omitempty"`
	// Strata fields (when promo uses strata, especially sequential)
	StrataID                []string                `json:"strata_id,omitempty"`
	StrataDesc              []string                `json:"strata_desc,omitempty"`
	StrataReward            []float64               `json:"strata_reward,omitempty"`
	StrataRuleType          RuleType                `json:"strata_rule_type,omitempty"`
	StrataRuleUom           UOMType                 `json:"strata_rule_uom,omitempty"`
	StrataRewardUom         UOMType                 `json:"strata_reward_uom,omitempty"`
	StrataRewardType        RewardType              `json:"strata_reward_type,omitempty"`
	StrataPerScope          string                  `json:"strata_per_scope,omitempty"`
	TotalGrossValue         float64                 `json:"total_gross_value"`
	ProductsEligible        []int                   `json:"products_eligible"`
	RewardPrice             []PromoRewardPrice      `json:"reward_price,omitempty"`
	RewardValue             []PromoRewardValue      `json:"reward_value,omitempty"`
	RewardPercentage        []PromoRewardPercentage `json:"reward_percentage,omitempty"`
	RewardProduct           []PromoRewardProductDet `json:"reward_product,omitempty"`
	RewardUnfulfilledReason *string                 `json:"reward_unfulfilled_reason,omitempty"`
}

// RewardUnfulfilledReason values when slab is eligible but product reward is not given.
const RewardUnfulfilledInsufficientStock = "insufficient_stock"

func NormalizePromoIDList(raw []string) ([]string, error) {
	result := make([]string, 0)
	seen := make(map[string]struct{})

	for _, value := range raw {
		for _, part := range strings.Split(value, ",") {
			cleaned := strings.TrimSpace(part)
			if cleaned == "" {
				continue
			}
			for _, r := range cleaned {
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
					continue
				}
				return nil, fmt.Errorf("invalid promo_id value %q", cleaned)
			}
			if _, exists := seen[cleaned]; exists {
				continue
			}
			seen[cleaned] = struct{}{}
			result = append(result, cleaned)
		}
	}

	return result, nil
}

type PromoRewardPrice struct {
	ProID    int     `json:"pro_id"`
	SubTotal float64 `json:"sub_total"`
	Promo1   float64 `json:"promo1"`
	Promo2   float64 `json:"promo2"`
	Promo3   float64 `json:"promo3"`
	Promo4   float64 `json:"promo4"`
	Promo5   float64 `json:"promo5"`
	NetValue float64 `json:"net_value"`
}

type PromoRewardValue struct {
	ProID      int     `json:"pro_id"`
	GrossValue float64 `json:"gross_value"`
	Promo1     float64 `json:"promo1"`
	Promo2     float64 `json:"promo2"`
	Promo3     float64 `json:"promo3"`
	Promo4     float64 `json:"promo4"`
	Promo5     float64 `json:"promo5"`
	NetValue   float64 `json:"net_value"`
}

type PromoRewardPercentage struct {
	ProID      int     `json:"pro_id"`
	GrossValue float64 `json:"gross_value"`
	Promo1     float64 `json:"promo1"`
	Promo2     float64 `json:"promo2"`
	Promo3     float64 `json:"promo3"`
	Promo4     float64 `json:"promo4"`
	Promo5     float64 `json:"promo5"`
	NetValue   float64 `json:"net_value"`
}

type PromoRewardProductDet struct {
	ProID      int     `json:"pro_id"`
	Qty1       float64 `json:"qty1"`
	Qty2       float64 `json:"qty2"`
	Qty3       float64 `json:"qty3"`
	GrossValue float64 `json:"gross_value"`
	Promo1     float64 `json:"promo1"`
	Promo2     float64 `json:"promo2"`
	Promo3     float64 `json:"promo3"`
	Promo4     float64 `json:"promo4"`
	Promo5     float64 `json:"promo5"`
}

type PromoConversionReq struct {
	CustID    string `json:"cust_id" validate:"required,max=10"`
	ProductID int64  `json:"pro_id"`
	Qty1      int64  `json:"qty1" validate:"omitempty,numeric"`
	Qty2      int64  `json:"qty2" validate:"omitempty,numeric"`
	Qty3      int64  `json:"qty3" validate:"omitempty,numeric"`
}

type ConversionWithPriceResp struct {
	Qty1       int64   `json:"qty1"`
	Qty2       int64   `json:"qty2"`
	Qty3       int64   `json:"qty3"`
	TotalQty   int64   `json:"total_qty"`
	SellPrice1 float64 `json:"sell_price1"`
	SellPrice2 float64 `json:"sell_price2"`
	SellPrice3 float64 `json:"sell_price3"`
}

type PromoConverResp struct {
	Qty1     int64 `json:"qty1"`
	Qty2     int64 `json:"qty2"`
	Qty3     int64 `json:"qty3"`
	TotalQty int64 `json:"total_qty"`
}
