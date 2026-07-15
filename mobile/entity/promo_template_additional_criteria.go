package entity

type PromoTemplateAdditionalCriteria struct {
	CustID                     string  `json:"cust_id,omitempty"`
	PromoTemplateID            string  `json:"promo_template_id,omitempty"`
	PromoTemplateAddCriteriaID *int64  `json:"promo_template_add_criteria_id"`
	Attribute                  string  `json:"attribute" validate:"oneof=PRO OCL OCH OTY OTG STY STE"` // PRO=product, OCL=Outlet Class,
	AttributeName              string  `json:"attribute_name"`
	Condition                  string  `json:"condition" validate:"oneof=I E"` // I=include, E=Exclude
	ConditionName              string  `json:"condition_name"`
	ReferenceID                int64   `json:"reference_id"`   // product id / outlet class id dst
	ReferenceCode              string  `json:"reference_code"` // "458391123",
	ReferenceName              string  `json:"reference_name"` // "SoKlin Detergen Cair Pink 10x6x60 ml",
	IsMandatory                bool    `json:"is_mandatory"`
	MinBuyType                 int     `json:"min_buy_type" validate:"oneof=0 1 2"` // 1=qty, 2=amount
	MinBuyTypeName             string  `json:"min_buy_type_name"`
	MinBuyValue                float64 `json:"min_buy_value"`
	MinBuyUom                  int     `json:"min_buy_uom" validate:"oneof=0 1 2 3"` // 1=smallest, 2=middle, 3=largest
	MinBuyUomName              string  `json:"min_buy_uom_name"`
}
