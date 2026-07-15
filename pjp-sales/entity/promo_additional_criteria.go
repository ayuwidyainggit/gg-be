package entity

/* API Spec
"promo_additional_criterias":[
  {
    "attribute": "PRO",
    // PRO=product, OCL=Outlet Class,
    // OCH=Outlet Channel, OTY=Outlet Type
    // OTG=Outlet Group, STY=Sales Type
    // STE=Sales Team
    "condition": "I", // I=include, E=Exclude
    "condition_name": "Include",
    "reference_id": 1, // product id / outlet class id dst
    "reference_code": "458391123",
    "reference_name": "SoKlin Detergen Cair Pink 10x6x60 ml",
    "is_mandatory": false,
    "min_buy_type": 1, // 1=qty, 2=amount
    "min_buy_type_name": "Quantity",
    "min_buy_value": 1,
    "min_buy_uom": 1, // 1=smallest, 2=middle, 3=largest
    "min_buy_uom_name": "Smallest",
  }
]
*/

type PromoAdditionalCriteria struct {
	CustID             string  `json:"cust_id,omitempty"`
	PromoID            string  `json:"promo_id,omitempty"`
	PromoAddCriteriaID *int64  `json:"promo_add_criteria_id"`
	Attribute          string  `json:"attribute" validate:"oneof=PRO OCL OCH OTY OTG STY STE"` // PRO=product, OCL=Outlet Class,
	AttributeName      string  `json:"attribute_name"`
	Condition          string  `json:"condition" validate:"oneof=I E"` // I=include, E=Exclude
	ConditionName      string  `json:"condition_name"`
	ReferenceID        int64   `json:"reference_id"`   // product id / outlet class id dst
	ReferenceCode      string  `json:"reference_code"` // "458391123",
	ReferenceName      string  `json:"reference_name"` // "SoKlin Detergen Cair Pink 10x6x60 ml",
	IsMandatory        bool    `json:"is_mandatory"`
	MinBuyType         int     `json:"min_buy_type" validate:"oneof=0 1 2"` // 1=qty, 2=amount
	MinBuyTypeName     string  `json:"min_buy_type_name"`
	MinBuyValue        float64 `json:"min_buy_value"`
	MinBuyUom          int     `json:"min_buy_uom" validate:"oneof=0 1 2 3"` // 1=smallest, 2=middle, 3=largest
	MinBuyUomName      string  `json:"min_buy_uom_name"`
}
