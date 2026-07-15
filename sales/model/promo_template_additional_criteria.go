package model

type PromoTemplateAdditionalCriteria struct {
	CustID                      string  `gorm:"column:cust_id" json:"cust_id"`
	PromoTemplateID             string  `gorm:"column:promo_template_id" json:"promo_template_id"`
	PromoTeamplateAddCriteriaID *int64  `gorm:"column:promo_template_add_criteria_id;primaryKey;autoIncrement:true" json:"promo_template_add_criteria_id"`
	Attribute                   string  `gorm:"column:attribute" json:"attribute"`
	Condition                   string  `gorm:"column:condition" json:"condition"`
	ReferenceID                 int64   `gorm:"column:reference_id" json:"reference_id"`
	IsMandatory                 bool    `gorm:"column:is_mandatory" json:"is_mandatory"`
	MinBuyType                  int64   `gorm:"column:min_buy_type" json:"min_buy_type"`
	MinBuyValue                 float64 `gorm:"column:min_buy_value" json:"min_buy_value"`
	MinBuyUom                   int64   `gorm:"column:min_buy_uom" json:"min_buy_uom"`
}

func (PromoTemplateAdditionalCriteria) TableName() string {
	return "sls.promo_template_additional_criterias"
}

type TempProductAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempProductAdditionalCriteria) TableName() string {
	return "mst.m_product"
}

type TempOutletClassAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempOutletClassAdditionalCriteria) TableName() string {
	return "mst.m_outlet_class"
}

type TempOutletTypeAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempOutletTypeAdditionalCriteria) TableName() string {
	return "mst.m_outlet_type"
}

type TempOutletGroupAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempOutletGroupAdditionalCriteria) TableName() string {
	return "mst.m_outlet_group"
}

type TempSalesTypeAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempSalesTypeAdditionalCriteria) TableName() string {
	return "mst.m_sales_type"
}

type TempSalesTeamAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (TempSalesTeamAdditionalCriteria) TableName() string {
	return "mst.m_sales_team"
}
