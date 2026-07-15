package entity

type PromotionsQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type PromotionsResp struct {
	Code  string   `json:"code"`
	Title string   `json:"title"`
	Tnc   []string `json:"tnc"`
}

type PromoStatusDescSlice []PromotionStatus

func (p PromoStatusDescSlice) Len() int {
	return len(p)
}

func (p PromoStatusDescSlice) Less(i, j int) bool {
	return p[i].PromoStatusID < p[j].PromoStatusID
}

func (p PromoStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PromotionQueryFilter struct {
	PromoStatusID []int `query:"promo_status_id"`
	CustId        string
	ParentCustId  string
	EffectiveFrom *int64 `query:"effective_from" validate:"required_with=EffectiveTo,omitempty,gte=1000000000"`
	EffectiveTo   *int64 `query:"effective_to" validate:"required_with=EffectiveFrom,omitempty,lte=9999999999,gtefield=EffectiveFrom"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required"`
	Query         string `query:"q"`
	PromoID       string `query:"promo_id"`
	PromoDesc     string `query:"promo_desc"`
	Mode          string `query:"mode"`
	Sort          string `query:"sort"`
}

/* API Spec:
	"promo_id": "PR0910234",
    "promo_desc": "Promo Pembelian ... ",
    "promo_type": 1, // 1 = NEW, 2 = Replacement,
    "existing_promo_id": "PR0910231",
    "promo_status_id": 1,
    "is_budget_reference": false, // false = No, true = YES
    "budget_reference_type": 1, // 1=unlimited, 2=manual input\
    "budget_reference_id": 1,
    "budget_amount": 0,
    "budget_control_level": 1, // 1=distributor, 2=salesman, 3=outlet, 4=area
    "execution_level": 1, // 1=distributor, 2=salesman, 3=outlet, 4=area
    "effective_from": "2024-06-25",
    "effective_to": "2024-06-29",
    "is_claimable": false,
    "claim_days": 4,
    "max_invoice_outlet": 99999,
    "max_discount_type": 1, // 1=quantity, 2=amount
    "max_discount_outlet": 50,
    "is_multiplied": false,
	"promo_criterias": [],
	"promo_additional_criterias": []
*/

type CreatePromotionBody struct {
	CustID                  string                    `json:"cust_id"`
	ParentCustID            string                    `json:"parent_cust_id"`
	PromoID                 string                    `json:"promo_id" validate:"required,max=20,alphanum"`
	PromoDesc               string                    `json:"promo_desc" validate:"required,max=100"`
	PromoType               int                       `json:"promo_type" validate:"required,oneof=1 2"` // 1 = NEW, 2 = Replacement,
	ExistingPromoID         string                    `json:"existing_promo_id" validate:"max=20"`
	PromoStatusID           int64                     `json:"promo_status_id" validate:"required,oneof=1 2 3 4 5 6 7"`
	IsMultiplied            bool                      `json:"is_multiplied"`
	IsBudgetReference       bool                      `json:"is_budget_reference"`
	BudgetReferenceType     int                       `json:"budget_reference_type" validate:"required,oneof=1 2"` // 1=unlimited, 2=manual input\
	BudgetReferenceID       int64                     `json:"budget_reference_id"`
	BudgetControlLevel      int                       `json:"budget_control_level" validate:"required,oneof=1 2 3 4"` // 1=distributor, 2=salesman, 3=outlet, 4=area
	BudgetAmount            float64                   `json:"budget_amount"`
	ExecutionLevel          int                       `json:"execution_level" validate:"required,oneof=1 2 3 4"` // 1=distributor, 2=salesman, 3=outlet, 4=area
	EffectiveFrom           string                    `json:"effective_from" validate:"required"`
	EffectiveTo             string                    `json:"effective_to" validate:"required"`
	IsClaimable             bool                      `json:"is_claimable"`
	ClaimDays               int64                     `json:"claim_days" validate:"min=0,max=9999"`
	MaxDiscountType         int64                     `json:"max_discount_type" validate:"required,oneof=1 2"` // 1=quantity, 2=amount
	MaxDiscountOutlet       float64                   `json:"max_discount_outlet" validate:"max=999999999"`
	MaxInvoiceOutlet        float64                   `json:"max_invoice_outlet" validate:"max=99999"`
	CreatedBy               string                    `json:"created_by"`
	PromoCriteria           []PromoCriteria           `json:"promo_criterias" validate:"dive"`
	RewardProduct           []PromoRewardProduct      `json:"reward_products" validate:"dive"`
	PromoAdditionalCriteria []PromoAdditionalCriteria `json:"promo_additional_criterias" validate:"dive"`
	MaxDiscountOutletUom    int                       `json:"max_discount_outlet_uom" validate:"oneof=0 1 2 3"`
}

type DetailPromotionParams struct {
	PromoID      string `params:"promo_id" validate:"required,max=20,alphanum"`
	CustID       string
	ParentCustId string
}

type DeletePromotionParams struct {
	PromoID string `params:"promo_id" validate:"required,max=20,alphanum"`
}

type UpdatePromotionParams struct {
	PromoID string `params:"promo_id" validate:"required,max=20,alphanum"`
}

var PromoTypeName = map[int]string{
	1: "New", 2: "Replacement",
}

var PromoBudgetReferenceTypeName = map[int]string{
	1: "Unlimited", 2: "Manual Input",
}

var PromoStatusDesc = map[int]string{
	1: "Draft", 2: "Submitted", 3: "Approved", 4: "Rejected", 5: "Expired", 6: "Active", 7: "Inactive",
}

func (promo Promotion) GetPromoTypeName() string {
	return PromoTypeName[promo.PromoType]
}

func (promo Promotion) GetPromoBudgetReferenceTypeName() string {
	return PromoBudgetReferenceTypeName[promo.BudgetReferenceType]
}

func (promo Promotion) GetPromoStatusDesc() string {
	return PromoStatusDesc[promo.PromoStatusID]
}

/* API Spec:
"promo_id": "PR0910231",
"promo_desc": "Promo Pembelian ... ",
"promo_type": 1, // 1 = NEW, 2 = Replacement,
"promo_type_name": "New",
"existing_promo_id": "PR0910231",
"promo_status_id": 1,
"promo_status_desc": "Draft",
"is_budget_reference": false, // false = No, true = YES
"budget_reference_type": 1, // 1=Unlimited, 2=Manual Input
"budget_reference_type_name": "Unlimited",
"budget_amount": 0,
"budget_control_level": 1, // 1=distributor, 2=salesman, 3=outlet, 4=area
"budget_control_level_name": "Distributor",
"execution_level": 1, // 1=distributor, 2=salesman, 3=outlet, 4=area
"execution_level_name": "Distributor",
"effective_from": "2024-06-25",
"effective_to": "2024-06-29",
"is_claimable": false,
"claim_days": 4,
"max_invoice_outlet": 99999,
"max_discount_type": 1, // 1=Quantity, 2=Amount
"max_discount_type_name": "Quantity",
"max_discount_outlet": 50,
"is_multiplied": false,
"created_at": "2024-05-14T10:59:50.819233Z"
*/

type Promotion struct {
	PromoID                  string                    `json:"promo_id"`
	PromoDesc                string                    `json:"promo_desc"`
	PromoType                int                       `json:"promo_type"`
	PromoTypeName            string                    `json:"promo_type_name"`
	ExistingPromoID          string                    `json:"existing_promo_id"`
	PromoStatusID            int                       `json:"promo_status_id"`
	PromoStatusDesc          string                    `json:"promo_status_desc"`
	IsMultiplied             bool                      `json:"is_multiplied"`
	IsBudgetReference        bool                      `json:"is_budget_reference"`
	BudgetReferenceType      int                       `json:"budget_reference_type"`
	BudgetReferenceTypeName  string                    `json:"budget_reference_type_name"`
	BudgetReferenceID        int64                     `json:"budget_reference_id"`
	BudgetControlLevel       int                       `json:"budget_control_level"`
	BudgetControlLevelName   string                    `json:"budget_control_level_name"`
	BudgetAmount             float64                   `json:"budget_amount"`
	ExecutionLevel           int                       `json:"execution_level"`
	ExecutionLevelName       string                    `json:"execution_level_name"`
	EffectiveFrom            string                    `json:"effective_from"`
	EffectiveTo              string                    `json:"effective_to"`
	IsClaimable              bool                      `json:"is_claimable"`
	ClaimDays                int64                     `json:"claim_days"`
	MaxDiscountType          int                       `json:"max_discount_type"`
	MaxDiscountTypeName      string                    `json:"max_discount_type_name"`
	MaxDiscountOutlet        float64                   `json:"max_discount_outlet"`
	MaxInvoiceOutlet         float64                   `json:"max_invoice_outlet"`
	PromoCriterias           []PromoCriteria           `json:"promo_criterias"`
	RewardProduct            []PromoRewardProduct      `json:"reward_products"`
	PromoAdditionalCriterias []PromoAdditionalCriteria `json:"promo_additional_criterias"`
	CreatedAt                string                    `json:"created_at"`
	UpdatedAt                string                    `json:"updated_at,omitempty"`
	CreatedBy                string                    `json:"created_by"`
	UpdatedBy                string                    `json:"updated_by,omitempty"`
	MaxDiscountOutletUom     int                       `json:"max_discount_outlet_uom"`
	MaxDiscountOutletUomName string                    `json:"max_discount_outlet_uom_name"`
	Remarks                  string                    `json:"remarks"`
}

type UpdatePromotionBody struct {
	CustID                  string                    `json:"cust_id"`
	ParentCustID            string                    `json:"parent_cust_id"`
	PromoDesc               string                    `json:"promo_desc" validate:"max=100"`
	PromoType               int                       `json:"promo_type" validate:"oneof=1 2"` // 1 = NEW, 2 = Replacement,
	ExistingPromoID         string                    `json:"existing_promo_id" validate:"max=20"`
	PromoStatusID           int                       `json:"promo_status_id" validate:"oneof=1 2 3 4 5 6 7"`
	IsMultiplied            bool                      `json:"is_multiplied"`
	IsBudgetReference       bool                      `json:"is_budget_reference"`
	BudgetReferenceType     int                       `json:"budget_reference_type" validate:"required,oneof=1 2"` // 1=unlimited, 2=manual input\
	BudgetReferenceID       int64                     `json:"budget_reference_id"`
	BudgetControlLevel      int64                     `json:"budget_control_level"` // 1=distributor, 2=salesman, 3=outlet, 4=area
	BudgetAmount            float64                   `json:"budget_amount" validate:"min=0"`
	ExecutionLevel          int64                     `json:"execution_level"` // 1=distributor, 2=salesman, 3=outlet, 4=area
	EffectiveFrom           string                    `json:"effective_from" validate:"required"`
	EffectiveTo             string                    `json:"effective_to" validate:"required"`
	IsClaimable             bool                      `json:"is_claimable"`
	ClaimDays               int64                     `json:"claim_days" validate:"min=0,max=9999"`
	MaxDiscountType         int64                     `json:"max_discount_type" validate:"oneof=1 2"`
	MaxDiscountOutlet       float64                   `json:"max_discount_outlet" validate:"min=0,max=999999999"`
	MaxInvoiceOutlet        float64                   `json:"max_invoice_outlet" validate:"min=0,max=99999"`
	UpdatedBy               string                    `json:"updated_by"`
	PromoCriteria           []PromoCriteria           `json:"promo_criterias" validate:"dive"`
	RewardProduct           []PromoRewardProduct      `json:"reward_products" validate:"dive"`
	PromoAdditionalCriteria []PromoAdditionalCriteria `json:"promo_additional_criterias" validate:"dive"`
	MaxDiscountOutletUom    int                       `json:"max_discount_outlet_uom" validate:"oneof=0 1 2 3"`
}
type PromotionStatus struct {
	PromoStatusID   int    `json:"promo_status_id"`
	PromoStatusDesc string `json:"promo_status_desc"`
}

type UniqueRewardProductID struct {
	RewardProductID []PromoRewardProduct `json:"reward_products[].pro_id" validate:"unique,dive"`
}

type BulkUpdateStatusPromotionBody struct {
	CustID        string   `json:"cust_id" validate:"required"`
	ParentCustID  string   `json:"parent_cust_id" validate:"required"`
	PromoID       []string `json:"promo_id" validate:"required"`
	PromoStatusID int      `json:"promo_status_id" validate:"required"`
	UpdatedBy     string   `json:"updated_by" validate:"required"`
	Remarks       string   `json:"remarks"`
}

type ConsultPromotionBody struct {
	CustID       string                    `json:"cust_id" validate:"required"`
	ParentCustID string                    `json:"parent_cust_id" validate:"required"`
	OrderDate    string                    `json:"order_date" validate:"required"`
	OutletId     int                       `json:"outlet_id" validate:"required"`
	SalesmanId   int                       `json:"salesman_id" validate:"required"`
	WhId         int                       `json:"wh_id"`
	Details      []ConsultPromotionSubBody `json:"details" validate:"required"`
}

type ConsultPromotionSubBody struct {
	ProID     int     `json:"pro_id"`
	Qty1      float64 `json:"qty1"`
	Qty2      float64 `json:"qty2"`
	Qty3      float64 `json:"qty3"`
	ConvUnit2 float64 `json:"conv_unit2"`
	ConvUnit3 float64 `json:"conv_unit3"`
	SubTotal  int     `json:"sub_total"`
}

type ConsultPromotionResponse struct {
	PromotionID   string                                  `json:"promo_id"`
	PromotionDesc string                                  `json:"promo_desc"`
	SlabId        int64                                   `json:"slab_id"`
	SlabDesc      string                                  `json:"slab_desc"`
	SlabReward    float64                                 `json:"slab_reward"`
	Products      []int                                   `json:"products"`
	RewardPrice   []ConsultPromotionRewardPriceResponse   `json:"reward_price"`
	RewardProduct []ConsultPromotionRewardProductResponse `json:"reward_product"`
}

type ConsultPromotionRewardProductResponse struct {
	ProID int     `json:"pro_id"`
	Qty1  float64 `json:"qty1"`
	Qty2  float64 `json:"qty2"`
	Qty3  float64 `json:"qty3"`
	// UnitId string  `json:"unit_id"`
	// Uom    int     `json:"uom"`
}

type ConsultPromotionRewardPriceResponse struct {
	ProID    int     `json:"pro_id"`
	SubTotal float64 `json:"sub_total"`
	Reward   float64 `json:"reward"`
	Total    float64 `json:"total"`
}

type PromoAdditionalCriteriaReferences struct {
	Condition   string `json:"condition"`
	ReferenceId []int  `json:"reference_id"`
	IsValid     bool   `json:"is_valid"`
}

type PromotionMobileListQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Sort         string `query:"sort"`
	PromoDesc    string `query:"promo_desc"`
}

type PromotionMobileListResponse struct {
	PromoID          string  `json:"promo_id"`
	PromoDesc        string  `json:"promo_desc"`
	EffectiveFrom    string  `json:"effective_from"`
	EffectiveTo      string  `json:"effective_to"`
	IsMultiplied     bool    `json:"is_multiplied"`
	MaxInvoiceOutlet float64 `json:"max_invoice_outlet"`
}

type PromotionMobilePromotedProduct struct {
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	IsMandatory bool    `json:"is_mandatory"`
	MinPurchase float64 `json:"min_purchase"`
}

type PromotionMobileRewardProduct struct {
	ProID   int64  `json:"pro_id"`
	ProCode string `json:"pro_code"`
	ProName string `json:"pro_name"`
}

type PromotionMobileAdditionalCriteria struct {
	PromoAddCriteriaID int64  `json:"promo_add_criteria_id"`
	Attribute          string `json:"attribute"`
}

// PromotionMobileDetailResponse for GET /mobile/v1/promotion/{promo_id}
type PromotionMobileDetailResponse struct {
	PromoID              string                                  `json:"promo_id"`
	PromoDesc            string                                  `json:"promo_desc"`
	PromoType            string                                  `json:"promo_type"`
	EffectiveFrom        string                                  `json:"effective_from"`
	EffectiveTo          string                                  `json:"effective_to"`
	IsMultiplied         bool                                    `json:"is_multiplied"`
	MaxInvoiceOutlet     float64                                 `json:"max_invoice_outlet"`
	MaxPromoUsage        float64                                 `json:"max_promo_usage"`
	MaxTotalRewardType   string                                  `json:"max_total_reward_type"`
	MaxTotalRewardValue  float64                                 `json:"max_total_reward_value"`
	PromotedProduct      []PromotionMobileDetailPromotedProduct  `json:"promoted_product"`
	PromotionReward      []PromotionMobileDetailRewardItem       `json:"promotion_reward"`
	Criteria             []PromotionMobileDetailCriteria         `json:"criteria"`
	AdditionalCriteria   PromotionMobileDetailAdditionalCriteria `json:"additional_criteria"`
}

type PromotionMobileDetailPromotedProduct struct {
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	Mandatory   bool    `json:"mandatory"`
	MinBuyType  string  `json:"min_buy_type"`
	MinBuyQty   float64 `json:"min_buy_qty"`
	MinBuyValue float64 `json:"min_buy_value"`
	MinBuyUom   string  `json:"min_buy_uom"`
}

type PromotionMobileDetailCriteria struct {
	ProID       int64   `json:"pro_id"`
	ProCode     string  `json:"pro_code"`
	ProName     string  `json:"pro_name"`
	CountPromo  int     `json:"count_promo"`
	MinPurchase float64 `json:"min_purchase"`
	MaxPurchase float64 `json:"max_purchase"`
	Uom         string  `json:"uom"`
}

type PromotionMobileDetailRewardItem struct {
	SlabID        *string                          `json:"slab_id,omitempty"`
	StrataID      *string                          `json:"strata_id,omitempty"`
	SlabName      *string                          `json:"slab_name,omitempty"`
	StrataName    *string                          `json:"strata_name,omitempty"`
	Ordinal       int                              `json:"ordinal"`
	RuleType      string                           `json:"rule_type"`
	RangeFrom     float64                          `json:"range_from"`
	RangeTo       float64                          `json:"range_to"`
	RuleUom       string                           `json:"rule_uom"`
	RewardType    string                           `json:"reward_type"`
	RewardUom     string                           `json:"reward_uom"`
	RewardValue   float64                          `json:"reward_value"`
	ProductReward []PromotionMobileDetailProductReward `json:"product_reward"`
}

type PromotionMobileDetailProductReward struct {
	RewardProID int64  `json:"reward_pro_id"`
	ProID       int64  `json:"pro_id"`
	ProCode     string `json:"pro_code"`
	ProName     string `json:"pro_name"`
}

type PromotionMobileDetailAdditionalCriteria struct {
	SelectedOutletTypes   []PromotionMobileOutletType  `json:"selected_outlet_types"`
	SelectedOutletGroups  []PromotionMobileOutletGroup `json:"selected_outlet_groups"`
	SelectedOutletClasses []PromotionMobileOutletClass `json:"selected_outlet_classes"`
}

type PromotionMobileOutletType struct {
	OutletTypeID   int64  `json:"outlet_type_id"`
	OutletTypeCode string `json:"outlet_type_code"`
	OutletTypeName string `json:"outlet_type_name"`
}

type PromotionMobileOutletGroup struct {
	OutletGroupID   int64  `json:"outlet_group_id"`
	OutletGroupCode string `json:"outlet_group_code"`
	OutletGroupName string `json:"outlet_group_name"`
}

type PromotionMobileOutletClass struct {
	OutletClassID   int64  `json:"outlet_class_id"`
	OutletClassCode string `json:"outlet_class_code"`
	OutletClassName string `json:"outlet_class_name"`
}

type PromotionMobileDetailParams struct {
	PromoID string `params:"promo_id" validate:"required"`
}

type PromotionOutletListQueryFilter struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	OtGrpID   *int64 `query:"ot_grp_id"`
	OtClassID *int64 `query:"ot_class_id"`
}

type PromotionOutletListParams struct {
	OtTypeID int64 `params:"ot_type_id" validate:"required"`
}

type PromotionOutletListResponse struct {
	OtTypeID   int64  `json:"ot_type_id"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
	Address1   string `json:"address1"`
	TodayVisit bool   `json:"today_visit"`
}
