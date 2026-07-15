package entity

type PromoTemplateStatusDescSlice []PromoTemplateStatus

func (p PromoTemplateStatusDescSlice) Len() int {
	return len(p)
}

func (p PromoTemplateStatusDescSlice) Less(i, j int) bool {
	return p[i].PromoTemplateStatusID < p[j].PromoTemplateStatusID
}

func (p PromoTemplateStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PromoTemplateQueryFilter struct {
	PromoTemplateStatusID []int `query:"promo_template_status_id"`
	CustId                string
	ParentCustId          string
	From                  *int64 `query:"from" validate:"required_with=From,omitempty,gte=1000000000"`
	To                    *int64 `query:"to" validate:"required_with=To,omitempty,lte=9999999999,gtefield=From"`
	Page                  int    `query:"page"`
	Limit                 int    `query:"limit" validate:"required"`
	Query                 string `query:"q"`
	PromoTemplateID       string `query:"promo_template_id"`
	PromoDesc             string `query:"promo_desc"`
	Mode                  string `query:"mode"`
	Sort                  string `query:"sort"`
}

type CreatePromoTemplateBody struct {
	CustID                          string                            `json:"cust_id"`
	ParentCustID                    string                            `json:"parent_cust_id"`
	PromoTemplateID                 string                            `json:"promo_template_id" validate:""`
	PromoDesc                       string                            `json:"promo_desc" validate:"required,max=100"`
	PromoTemplateStatusID           int                               `json:"promo_template_status_id" validate:"required,oneof=1 2"`
	IsMultiplied                    bool                              `json:"is_multiplied"`
	IsBudgetReference               bool                              `json:"is_budget_reference"`
	BudgetReferenceType             int                               `json:"budget_reference_type" validate:"required,oneof=1 2"` // 1=unlimited, 2=manual input\
	BudgetReferenceID               int64                             `json:"budget_reference_id"`
	BudgetAmount                    float64                           `json:"budget_amount"`
	IsClaimable                     bool                              `json:"is_claimable"`
	ClaimDays                       int                               `json:"claim_days" validate:"min=0,max=9999"`
	MaxDiscountType                 int64                             `json:"max_discount_type" validate:"required,oneof=1 2"` // 1=quantity, 2=amount
	MaxDiscountOutlet               float64                           `json:"max_discount_outlet" validate:"max=999999999"`
	MaxInvoiceOutlet                float64                           `json:"max_invoice_outlet" validate:"max=99999"`
	CreatedBy                       string                            `json:"created_by"`
	PromoTemplateCriteria           []PromoTemplateCriteria           `json:"promo_criterias" validate:"dive"`
	PromoTemplateRewardProduct      []PromoTemplateRewardProduct      `json:"reward_products" validate:"dive"`
	PromoTemplateAdditionalCriteria []PromoTemplateAdditionalCriteria `json:"promo_additional_criterias" validate:"dive"`
	MaxDiscountOutletUom            int                               `json:"max_discount_outlet_uom" validate:"oneof=0 1 2 3"`
}

type DetailPromoTemplateParams struct {
	PromoTemplateID string `params:"promo_template_id" validate:"required,max=20"`
	CustID          string
	ParentCustId    string
}

type DeletePromoTemplateParams struct {
	PromoTemplateID string `params:"promo_template_id" validate:"required,max=20"`
}

type UpdatePromoTemplateParams struct {
	PromoTemplateID string `params:"promo_template_id" validate:"required,max=20"`
}

var PromoTemplateStatusDesc = map[int]string{
	1: "Showed", 2: "Hidden",
}

func (promo PromoTemplate) GetPromoBudgetReferenceTypeName() string {
	return PromoBudgetReferenceTypeName[promo.BudgetReferenceType]
}

func (promo PromoTemplate) GetPromoStatusDesc() string {
	return PromoTemplateStatusDesc[promo.PromoTemplateStatusID]
}

type PromoTemplate struct {
	PromoTemplateID          string                            `json:"promo_template_id"`
	PromoDesc                string                            `json:"promo_desc"`
	PromoTemplateStatusID    int                               `json:"promo_template_status_id"`
	PromoStatusDesc          string                            `json:"promo_template_status_desc"`
	IsMultiplied             bool                              `json:"is_multiplied"`
	IsBudgetReference        bool                              `json:"is_budget_reference"`
	BudgetReferenceType      int                               `json:"budget_reference_type"`
	BudgetReferenceTypeName  string                            `json:"budget_reference_type_name"`
	BudgetReferenceID        int64                             `json:"budget_reference_id"`
	BudgetAmount             float64                           `json:"budget_amount"`
	IsClaimable              bool                              `json:"is_claimable"`
	ClaimDays                int64                             `json:"claim_days"`
	MaxDiscountType          int                               `json:"max_discount_type"`
	MaxDiscountTypeName      string                            `json:"max_discount_type_name"`
	MaxDiscountOutlet        float64                           `json:"max_discount_outlet"`
	MaxInvoiceOutlet         float64                           `json:"max_invoice_outlet"`
	PromoCriterias           []PromoTemplateCriteria           `json:"promo_criterias"`
	RewardProduct            []PromoTemplateRewardProduct      `json:"reward_products"`
	PromoAdditionalCriterias []PromoTemplateAdditionalCriteria `json:"promo_additional_criterias"`
	CreatedAt                string                            `json:"created_at"`
	UpdatedAt                string                            `json:"updated_at,omitempty"`
	CreatedBy                string                            `json:"created_by"`
	UpdatedBy                string                            `json:"updated_by,omitempty"`
	MaxDiscountOutletUom     int                               `json:"max_discount_outlet_uom"`
	MaxDiscountOutletUomName string                            `json:"max_discount_outlet_uom_name"`
}

type UpdatePromoTemplateBody struct {
	CustID                          string                            `json:"cust_id"`
	ParentCustID                    string                            `json:"parent_cust_id"`
	PromoDesc                       string                            `json:"promo_desc" validate:"max=100"`
	PromoTemplateStatusID           int64                             `json:"promo_template_status_id" validate:"oneof=1 2"`
	IsMultiplied                    bool                              `json:"is_multiplied"`
	IsBudgetReference               bool                              `json:"is_budget_reference"`
	BudgetReferenceType             int                               `json:"budget_reference_type" validate:"required,oneof=1 2"` // 1=unlimited, 2=manual input\
	BudgetReferenceID               int64                             `json:"budget_reference_id"`
	BudgetAmount                    float64                           `json:"budget_amount" validate:"min=0"`
	IsClaimable                     bool                              `json:"is_claimable"`
	ClaimDays                       int64                             `json:"claim_days" validate:"min=0,max=9999"`
	MaxDiscountType                 int64                             `json:"max_discount_type" validate:"oneof=1 2"`
	MaxDiscountOutlet               float64                           `json:"max_discount_outlet" validate:"min=0,max=999999999"`
	MaxInvoiceOutlet                float64                           `json:"max_invoice_outlet" validate:"min=0,max=99999"`
	UpdatedBy                       string                            `json:"updated_by"`
	PromoTemplateCriteria           []PromoTemplateCriteria           `json:"promo_criterias" validate:"dive"`
	PromoTemplateRewardProduct      []PromoTemplateRewardProduct      `json:"reward_products" validate:"dive"`
	PromoTemplateAdditionalCriteria []PromoTemplateAdditionalCriteria `json:"promo_additional_criterias" validate:"dive"`
	MaxDiscountOutletUom            int                               `json:"max_discount_outlet_uom" validate:"oneof=0 1 2 3"`
}

type PromoTemplateStatus struct {
	PromoTemplateStatusID   int    `json:"promo_template_status_id"`
	PromoTemplateStatusDesc string `json:"promo_template_status_desc"`
}

type PromoTemplateUniqueRewardProductID struct {
	RewardProductID []PromoTemplateRewardProduct `json:"reward_products[].pro_id" validate:"unique,dive"`
}
