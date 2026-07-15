package model

type ProductDist struct {
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ProID          *int64   `gorm:"column:pro_id" json:"pro_id"`
	Cogs           *float64 `gorm:"column:cogs" json:"cogs"`
	IsActive       bool     `gorm:"column:is_active" json:"is_active" `
	IsAlloc        bool     `gorm:"column:is_alloc" json:"is_alloc"`
	MinStock       float64  `gorm:"column:min_stock" json:"min_stock" `
	MinStockStr    string   `gorm:"column:min_stock_str" json:"min_stock_str"`
	SafetyStock    float64  `gorm:"column:safety_stock" json:"safety_stock"`
	SafetyStockStr string   `gorm:"column:safety_stock_str" json:"safety_stock_str"`
	PoFormula      int      `gorm:"column:po_formula" json:"po_formula"`
	IsNewPro       bool     `gorm:"column:is_new_pro" json:"is_new_pro"`
	Vat            float64  `gorm:"column:vat" json:"vat"`
	VatBg          float64  `gorm:"column:vat_bg" json:"vat_bg"`
	VatLgPurch     float64  `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	VatLgSell      float64  `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
}

func (ProductDist) TableName() string {
	return "mst.m_product_dist"
}

type ProductStockList struct {
	CustID        string   `gorm:"column:cust_id" json:"cust_id"`
	ProId         int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string   `gorm:"column:pro_code" json:"pro_code"`
	ProName       string   `gorm:"column:pro_name" json:"pro_name"`
	SupId         *int64   `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string  `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string  `gorm:"column:sup_name" json:"sup_name"`
	BrandId       *int     `gorm:"column:brand_id" json:"brand_id"`
	BrandCode     *string  `gorm:"column:brand_code" json:"brand_code"`
	BrandName     *string  `gorm:"column:brand_name" json:"brand_name"`
	SBrand1Id     *int     `gorm:"column:sbrand1_id" json:"sbrand1_id"`
	SBrand1Code   *string  `gorm:"column:sbrand1_code" json:"sbrand1_code"`
	SBrand1Name   *string  `gorm:"column:sbrand1_name" json:"sbrand1_name"`
	UnitId1       *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2       *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3       *string  `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4       *string  `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5       *string  `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2     *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3     *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4     *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5     *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
	Qty           *float64 `gorm:"column:qty" json:"qty"`
	QtyOnOrder    *float64 `gorm:"column:qty_on_order" json:"qty_on_order"`
	QtyOnShipping *float64 `gorm:"column:qty_on_shipping" json:"qty_on_shipping"`
	QtyBs         *float64 `gorm:"column:qty_bs" json:"qty_bs"`
	QtyExp        *float64 `gorm:"column:qty_exp" json:"qty_exp"`
}

func (ProductStockList) TableName() string {
	return "mst.m_product_dist"
}

type CheckProductStockList struct {
	CustID         string   `gorm:"column:cust_id" json:"cust_id"`
	ProId          int64    `gorm:"column:pro_id" json:"pro_id"`
	ProCode        string   `gorm:"column:pro_code" json:"pro_code"`
	ProName        string   `gorm:"column:pro_name" json:"pro_name"`
	SupId          *int64   `gorm:"column:sup_id" json:"sup_id"`
	SupCode        *string  `gorm:"column:sup_code" json:"sup_code"`
	SupName        *string  `gorm:"column:sup_name" json:"sup_name"`
	BrandId        *int     `gorm:"column:brand_id" json:"brand_id"`
	BrandCode      *string  `gorm:"column:brand_code" json:"brand_code"`
	BrandName      *string  `gorm:"column:brand_name" json:"brand_name"`
	SBrand1Id      *int     `gorm:"column:sbrand1_id" json:"sbrand1_id"`
	SBrand1Code    *string  `gorm:"column:sbrand1_code" json:"sbrand1_code"`
	SBrand1Name    *string  `gorm:"column:sbrand1_name" json:"sbrand1_name"`
	UnitId1        *string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2        *string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3        *string  `gorm:"column:unit_id3" json:"unit_id3"`
	UnitId4        *string  `gorm:"column:unit_id4" json:"unit_id4"`
	UnitId5        *string  `gorm:"column:unit_id5" json:"unit_id5"`
	ConvUnit2      *float64 `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3      *float64 `gorm:"column:conv_unit3" json:"conv_unit3"`
	ConvUnit4      *float64 `gorm:"column:conv_unit4" json:"conv_unit4"`
	ConvUnit5      *float64 `gorm:"column:conv_unit5" json:"conv_unit5"`
	Qty1           *float64 `gorm:"column:qty1" json:"qty1"`
	Qty2           *float64 `gorm:"column:qty2" json:"qty2"`
	Qty3           *float64 `gorm:"column:qty3" json:"qty3"`
	Qty4           *float64 `gorm:"column:qty4" json:"qty4"`
	Qty5           *float64 `gorm:"column:qty5" json:"qty5"`
	QtyOnOrder1    *float64 `gorm:"column:qty_on_order1" json:"qty_on_order1"`
	QtyOnOrder2    *float64 `gorm:"column:qty_on_order2" json:"qty_on_order2"`
	QtyOnOrder3    *float64 `gorm:"column:qty_on_order3" json:"qty_on_order3"`
	QtyOnOrder4    *float64 `gorm:"column:qty_on_order4" json:"qty_on_order4"`
	QtyOnOrder5    *float64 `gorm:"column:qty_on_order5" json:"qty_on_order5"`
	QtyOnShipping1 *float64 `gorm:"column:qty_on_shipping1" json:"qty_on_shipping1"`
	QtyOnShipping2 *float64 `gorm:"column:qty_on_shipping2" json:"qty_on_shipping2"`
	QtyOnShipping3 *float64 `gorm:"column:qty_on_shipping3" json:"qty_on_shipping3"`
	QtyOnShipping4 *float64 `gorm:"column:qty_on_shipping4" json:"qty_on_shipping4"`
	QtyOnShipping5 *float64 `gorm:"column:qty_on_shipping5" json:"qty_on_shipping5"`
	QtyBs1         *float64 `gorm:"column:qty_bs1" json:"qty_bs1"`
	QtyBs2         *float64 `gorm:"column:qty_bs2" json:"qty_bs2"`
	QtyBs3         *float64 `gorm:"column:qty_bs3" json:"qty_bs3"`
	QtyBs4         *float64 `gorm:"column:qty_bs4" json:"qty_bs4"`
	QtyBs5         *float64 `gorm:"column:qty_bs5" json:"qty_bs5"`
	QtyExp         *float64 `gorm:"column:qty_exp" json:"qty_exp"`
	QtyExp1        *float64 `gorm:"column:qty_exp1" json:"qty_exp1"`
	QtyExp2        *float64 `gorm:"column:qty_exp2" json:"qty_exp2"`
	QtyExp3        *float64 `gorm:"column:qty_exp3" json:"qty_exp3"`
	QtyExp4        *float64 `gorm:"column:qty_exp4" json:"qty_exp4"`
}

func (CheckProductStockList) TableName() string {
	return "mst.m_product_dist"
}
