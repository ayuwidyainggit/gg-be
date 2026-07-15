package entity

type DiscountStatusDescSlice []DiscountStatus

func (p DiscountStatusDescSlice) Len() int {
	return len(p)
}

func (p DiscountStatusDescSlice) Less(i, j int) bool {
	return p[i].DiscountStatusID < p[j].DiscountStatusID
}

func (p DiscountStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PublishStatusDescSlice []PublishStatus

func (p PublishStatusDescSlice) Len() int {
	return len(p)
}

func (p PublishStatusDescSlice) Less(i, j int) bool {
	return p[i].PublishStatusID < p[j].PublishStatusID
}

func (p PublishStatusDescSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type DiscountQueryFilter struct {
	DiscountStatusID []int `query:"discount_status_id"`
	CustId           string
	ParentCustId     string
	EffectiveFrom    *int64 `query:"effective_from" validate:"required_with=EffectiveTo,omitempty,gte=1000000000"`
	EffectiveTo      *int64 `query:"effective_to" validate:"required_with=EffectiveFrom,omitempty,lte=9999999999,gtefield=EffectiveFrom"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	DiscountID       string `query:"discount_id"`
	DiscountDesc     string `query:"discount_desc"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
}

type DetailDiscountParams struct {
	DiscountID   string `params:"discount_id" validate:"required,max=20,alphanumericSpaceSlashPercent"`
	CustID       string
	ParentCustId string
}

type DeleteDiscountParams struct {
	DiscountID string `params:"discount_id" validate:"required,max=20,alphanum"`
}

type UpdateDiscountParams struct {
	DiscountID string `params:"discount_id" validate:"required,max=20,alphanum"`
}

var DiscountStatusDesc = map[int]string{
	1: "Inactive", 2: "Active",
}

func (discount Discount) GetDiscountStatusDesc() string {
	return DiscountStatusDesc[discount.DiscountStatusID]
}

var PublishStatusDesc = map[int]string{
	1: "New", 2: "Published",
}

func (discount Discount) GetPublishStatusDesc() string {
	return PublishStatusDesc[discount.PublishStatusID]
}

type Discount struct {
	CustID             string              `json:"cust_id,omitempty"`
	ParentCustID       string              `json:"parent_cust_id,omitempty"`
	DiscountID         string              `json:"discount_id"`
	DiscountDesc       string              `json:"discount_desc"`
	DiscountStatusID   int                 `json:"discount_status_id"`
	DiscountStatusDesc string              `json:"discount_status_desc"`
	PublishStatusID    int                 `json:"publish_status_id"`
	PublishStatusDesc  string              `json:"publish_status_desc"`
	EffectiveFrom      string              `json:"effective_from"`
	EffectiveTo        string              `json:"effective_to"`
	CreatedAt          string              `json:"created_at"`
	CreatedBy          string              `json:"created_by"`
	DiscountPrincipals []DiscountPrincipal `json:"discount_principals,omitempty"`
	DiscountGroups     []DiscountGroup     `json:"discount_groups,omitempty"`
	DiscountCriterias  []DiscountCriteria  `json:"discount_criterias,omitempty"`
}

type CreateDiscountBody struct {
	CustID             string              `json:"cust_id"`
	ParentCustID       string              `json:"parent_cust_id"`
	DiscountID         string              `json:"discount_id" validate:"required,max=10,alphanum"`
	DiscountDesc       string              `json:"discount_desc" validate:"required,max=100"`
	DiscountStatusID   int                 `json:"discount_status_id" validate:"required,oneof=1 2"`
	EffectiveFrom      string              `json:"effective_from" validate:"required"`
	EffectiveTo        string              `json:"effective_to" validate:"required"`
	CreatedBy          string              `json:"created_by"`
	DiscountPrincipals []DiscountPrincipal `json:"discount_principals" validate:"dive"`
	DiscountGroups     []DiscountGroup     `json:"discount_groups" validate:"dive"`
	DiscountCriterias  []DiscountCriteria  `json:"discount_criterias" validate:"dive"`
}

type UpdateDiscountBody struct {
	CustID             string              `json:"cust_id"`
	ParentCustID       string              `json:"parent_cust_id"`
	DiscountDesc       string              `json:"discount_desc" validate:"required,max=100"`
	DiscountStatusID   int                 `json:"discount_status_id" validate:"oneof=1 2"`
	EffectiveFrom      string              `json:"effective_from" validate:"required"`
	EffectiveTo        string              `json:"effective_to" validate:"required"`
	UpdatedBy          string              `json:"updated_by"`
	DiscountPrincipals []DiscountPrincipal `json:"discount_principals" validate:"dive"`
	DiscountGroups     []DiscountGroup     `json:"discount_groups" validate:"dive"`
	DiscountCriterias  []DiscountCriteria  `json:"discount_criterias" validate:"dive"`
}

type DiscountStatus struct {
	DiscountStatusID   int    `json:"discount_status_id"`
	DiscountStatusDesc string `json:"discount_status_desc"`
}

type PublishStatus struct {
	PublishStatusID   int    `json:"publish_status_id"`
	PublishStatusDesc string `json:"publish_status_desc"`
}

type PublishDiscountBody struct {
	CustID       string   `json:"cust_id" validate:"required"`
	ParentCustID string   `json:"parent_cust_id" validate:"required"`
	DiscountID   []string `json:"discount_id" validate:"required"`
	UpdatedBy    string   `json:"updated_by" validate:"required"`
}

type ConsultDiscountBody struct {
	CustID       string                   `json:"cust_id" validate:"required"`
	ParentCustID string                   `json:"parent_cust_id" validate:"required"`
	OrderDate    string                   `json:"order_date" validate:"required"`
	OutletId     int                      `json:"outlet_id" validate:"required"`
	Details      []ConsultDiscountSubBody `json:"details" validate:"required"`
}

type ConsultDiscountSubBody struct {
	ProID    int `json:"pro_id"`
	SubTotal int `json:"sub_total"`
}

type ConsultDiscountResponse struct {
	ProID             int    `json:"pro_id"`
	SubTotal          int    `json:"sub_total"`
	SubTotalPrincipal int    `json:"sub_total_principal"`
	DiscountID        string `json:"discount_id"`
	PrincipalID       []int  `json:"principal_id"`
	SlabDesc          string `json:"slab_desc"`
	SlabReward        int    `json:"slab_reward"`
	RewardProduct     int    `json:"reward_product"`
}

type ConsultDiscountResponseNew struct {
	DiscountID   string                          `json:"discount_id"`
	DiscountDesc string                          `json:"discount_desc"`
	SlabId       int64                           `json:"slab_id"`
	SlabDesc     string                          `json:"slab_desc"`
	SlabReward   float64                         `json:"slab_reward"`
	Products     []int                           `json:"products"`
	Rewards      []ConsultDiscountRewardResponse `json:"rewards"`
}

type ConsultDiscountRewardResponse struct {
	ProID    int     `json:"pro_id"`
	SubTotal float64 `json:"sub_total"`
	Reward   float64 `json:"reward"`
	Total    float64 `json:"total"`
}

// DiscountStatusID int      `json:"discount_status_id" validate:"required"`
// Remarks          string   `json:"remarks"`
