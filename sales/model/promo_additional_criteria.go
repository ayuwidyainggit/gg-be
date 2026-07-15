package model

type PromoAdditionalCriteria struct {
	CustID             string  `gorm:"column:cust_id" json:"cust_id"`
	PromoID            string  `gorm:"column:promo_id" json:"promo_id"`
	PromoAddCriteriaID *int64  `gorm:"column:promo_add_criteria_id;primaryKey;autoIncrement:true" json:"promo_add_criteria_id"`
	Attribute          string  `gorm:"column:attribute" json:"attribute"`
	Condition          string  `gorm:"column:condition" json:"condition"`
	ReferenceID        int64   `gorm:"column:reference_id" json:"reference_id"`
	IsMandatory        bool    `gorm:"column:is_mandatory" json:"is_mandatory"`
	MinBuyType         int64   `gorm:"column:min_buy_type" json:"min_buy_type"`
	MinBuyValue        float64 `gorm:"column:min_buy_value" json:"min_buy_value"`
	MinBuyUom          int64   `gorm:"column:min_buy_uom" json:"min_buy_uom"`
}

func (PromoAdditionalCriteria) TableName() string {
	return "sls.promo_additional_criterias"
}

type ProductAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (ProductAdditionalCriteria) TableName() string {
	return "mst.m_product"
}

type OutletClassAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (OutletClassAdditionalCriteria) TableName() string {
	return "mst.m_outlet_class"
}

type OutletTypeAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (OutletTypeAdditionalCriteria) TableName() string {
	return "mst.m_outlet_type"
}

type OutletGroupAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (OutletGroupAdditionalCriteria) TableName() string {
	return "mst.m_outlet_group"
}

type SalesTypeAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (SalesTypeAdditionalCriteria) TableName() string {
	return "mst.m_sales_type"
}

type SalesTeamAdditionalCriteria struct {
	ReferenceID   int64  `gorm:"column:reference_id" json:"reference_id"`
	ReferenceCode string `gorm:"column:reference_code" json:"reference_code"`
	ReferenceName string `gorm:"column:reference_name" json:"reference_name"`
}

func (SalesTeamAdditionalCriteria) TableName() string {
	return "mst.m_sales_team"
}

type PromoAdditionalCriteriaByActivePromo struct {
	CustID             string  `gorm:"column:cust_id" json:"cust_id"`
	PromoID            string  `gorm:"column:promo_id" json:"promo_id"`
	PromoDesc          string  `json:"promo_desc"`
	IsMultiplied       bool    `json:"is_multiplied"`
	PromoAddCriteriaID *int64  `gorm:"column:promo_add_criteria_id;primaryKey;autoIncrement:true" json:"promo_add_criteria_id"`
	Attribute          string  `gorm:"column:attribute" json:"attribute"`
	Condition          string  `gorm:"column:condition" json:"condition"`
	ReferenceID        int64   `gorm:"column:reference_id" json:"reference_id"`
	IsMandatory        bool    `gorm:"column:is_mandatory" json:"is_mandatory"`
	MinBuyType         int64   `gorm:"column:min_buy_type" json:"min_buy_type"`
	MinBuyValue        float64 `gorm:"column:min_buy_value" json:"min_buy_value"`
	MinBuyUom          int64   `gorm:"column:min_buy_uom" json:"min_buy_uom"`
}

func (PromoAdditionalCriteriaByActivePromo) TableName() string {
	return "sls.promo_additional_criterias"
}
