package entity

import (
	"time"
)

type ProductsQueryFilter struct {
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
	SupId        *int   `query:"sup_id"`
	WhId         *int64 `query:"wh_id"`
	EmpId        int64  `query:"emp_id"`
}

var StatusManageMinimumPrice = map[int]string{
	0: "Non Active",
	1: "Submit",
	2: "Active",
}

func ConvStatusManageMinimumPrice(param int) string {
	result, ok := StatusManageMinimumPrice[param] // langsung gunakan int sebagai key
	if !ok {
		result = "Unknown"
	}
	return result
}

type ProductsResp struct {
	WhId                         int64   `json:"wh_id"`
	ProId                        int64   `json:"pro_id"`
	ProCode                      string  `json:"pro_code"`
	ProName                      string  `json:"pro_name"`
	BarCode                      string  `json:"bar_code"`
	UnitId1                      string  `json:"unit_id1"`
	UnitId2                      string  `json:"unit_id2"`
	UnitId3                      string  `json:"unit_id3"`
	UnitId4                      string  `json:"unit_id4"`
	UnitId5                      string  `json:"unit_id5"`
	UnitName1                    string  `json:"unit_name1"`
	UnitName2                    string  `json:"unit_name2"`
	UnitName3                    string  `json:"unit_name3"`
	UnitName4                    string  `json:"unit_name4"`
	UnitName5                    string  `json:"unit_name5"`
	ConvUnit1                    float32 `json:"conv_unit1"`
	ConvUnit2                    float32 `json:"conv_unit2"`
	ConvUnit3                    float32 `json:"conv_unit3"`
	ConvUnit4                    float32 `json:"conv_unit4"`
	ConvUnit5                    float32 `json:"conv_unit5"`
	Stock1                       float64 `json:"stock1"`
	Stock2                       float64 `json:"stock2"`
	Stock3                       float64 `json:"stock3"`
	Stock4                       float64 `json:"stock4"`
	Stock5                       float64 `json:"stock5"`
	Qty1                         float64 `json:"qty1"`
	Qty2                         float64 `json:"qty2"`
	Qty3                         float64 `json:"qty3"`
	PoFormula                    int     `json:"po_formula"`
	Vat                          float64 `json:"vat"`
	VatBg                        float64 `json:"vat_bg"`
	VatLgPurch                   float64 `json:"vat_lg_purch"`
	VatLgSell                    float64 `json:"vat_lg_sell"`
	ExciseRate                   float64 `json:"excise_rate"`
	ExciseTax                    float64 `json:"excise_tax"`
	ImageUrl                     string  `json:"image_url"`
	Cogs                         float64 `json:"cogs"`
	Price1                       float64 `json:"price1"`
	Price2                       float64 `json:"price2"`
	Price3                       float64 `json:"price3"`
	Price4                       float64 `json:"price4"`
	Price5                       float64 `json:"price5"`
	PromoCode                    string  `json:"promo_code"`
	CtgId1                       string  `json:"ctg_id1"`
	CtgId2                       string  `json:"ctg_id2"`
	CtgId3                       string  `json:"ctg_id3"`
	Price1Minimum                float64 `json:"price1_minimum"`
	Price2Minimum                float64 `json:"price2_minimum"`
	Price3Minimum                float64 `json:"price3_minimum"`
	Price4Minimum                float64 `json:"price4_minimum"`
	Price5Minimum                float64 `json:"price5_minimum"`
	LimitAction                  int     `json:"limit_action"`
	LimitActionName              string  `json:"limit_action_name"`
	StatusManageMinimumPriceName string  `json:"status_manage_minimum_price_name"`
}

type DetailProductParams struct {
	ProductId    int64 `params:"pro_id" validate:"required"`
	CustID       string
	ParentCustID string
}
type ProductDetailResponse struct {
	ProductId         int64      `json:"pro_id"`
	ProductCode       string     `json:"pro_code"`
	BarCode           *string    `json:"bar_code"`
	ProductName       string     `json:"pro_name"`
	PCatId            *int       `json:"pcat_id"`
	PCatCode          string     `json:"pcat_code"`
	PCatName          string     `json:"pcat_name"`
	ProductLineID     int        `json:"pl_id"`
	ProductLineCode   string     `json:"pl_code"`
	ProductLineName   string     `json:"pl_name"`
	BrandId           int        `json:"brand_id"`
	BrandCode         string     `json:"brand_code"`
	BrandName         string     `json:"brand_name"`
	Sbrand1           int        `json:"sbrand1_id"`
	Sbrand1Code       string     `json:"sbrand1_code"`
	Sbrand1Name       string     `json:"sbrand1_name"`
	Sbrand2           int        `json:"sbrand2_id"`
	Sbrand2Code       string     `json:"sbrand2_code"`
	Sbrand2Name       string     `json:"sbrand2_name"`
	FlavorId          int        `json:"flavor_id"`
	FlavorCode        string     `json:"flavor_code"`
	FlavorName        string     `json:"flavor_name"`
	PTypeId           int        `json:"ptype_id"`
	PTypeCode         string     `json:"ptype_code"`
	PTypeName         string     `json:"ptype_name"`
	PSizeId           int        `json:"psize_id"`
	PSizeCode         string     `json:"psize_code"`
	PSizeName         string     `json:"psize_name"`
	SupId             int        `json:"sup_id"`
	SupCode           string     `json:"sup_code"`
	SupName           string     `json:"sup_name"`
	PrincipalId       int        `json:"principal_id"`
	PrincipalCode     string     `json:"principal_code"`
	PrincipalName     string     `json:"principal_name"`
	CProId            int        `json:"c_pro_id"`
	CProCode          string     `json:"c_pro_code"`
	CProName          string     `json:"c_pro_name"`
	IsMainPro         bool       `json:"is_main_pro"`
	SortNo            int        `json:"sort_no"`
	ItemNo            int        `json:"item_no"`
	UnitId1           string     `json:"unit_id1"`
	UnitId2           string     `json:"unit_id2"`
	UnitId3           string     `json:"unit_id3"`
	UnitId4           string     `json:"unit_id4"`
	UnitId5           string     `json:"unit_id5"`
	ConvUnit2         float32    `json:"conv_unit2"`
	ConvUnit3         float32    `json:"conv_unit3"`
	ConvUnit4         float32    `json:"conv_unit4"`
	ConvUnit5         float32    `json:"conv_unit5"`
	IsBatch           bool       `json:"is_batch"`
	IsExpDate         bool       `json:"is_exp_date"`
	Length            float64    `json:"length"`
	Width             float64    `json:"width"`
	Height            float64    `json:"height"`
	Weight            float64    `json:"weight"`
	Volume            float64    `json:"volume"`
	ParentProId       int        `json:"parent_pro_id"`
	PurchPrice1       float64    `json:"purch_price1"`
	PurchPrice2       float64    `json:"purch_price2"`
	PurchPrice3       float64    `json:"purch_price3"`
	PurchPrice4       float64    `json:"purch_price4"`
	PurchPrice5       float64    `json:"purch_price5"`
	SellPrice1        float64    `json:"sell_price1"`
	SellPrice2        float64    `json:"sell_price2"`
	SellPrice3        float64    `json:"sell_price3"`
	SellPrice4        float64    `json:"sell_price4"`
	SellPrice5        float64    `json:"sell_price5"`
	Length1           *float64   `json:"length1"`
	Length2           *float64   `json:"length2"`
	Length3           *float64   `json:"length3"`
	Length4           *float64   `json:"length4"`
	Length5           *float64   `json:"length5"`
	Width1            *float64   `json:"width1"`
	Width2            *float64   `json:"width2"`
	Width3            *float64   `json:"width3"`
	Width4            *float64   `json:"width4"`
	Width5            *float64   `json:"width5"`
	Height1           *float64   `json:"height1"`
	Height2           *float64   `json:"height2"`
	Height3           *float64   `json:"height3"`
	Height4           *float64   `json:"height4"`
	Height5           *float64   `json:"height5"`
	Weight1           *float64   `json:"weight1"`
	Weight2           *float64   `json:"weight2"`
	Weight3           *float64   `json:"weight3"`
	Weight4           *float64   `json:"weight4"`
	Weight5           *float64   `json:"weight5"`
	Volume1           *float64   `json:"volume1"`
	Volume2           *float64   `json:"volume2"`
	Volume3           *float64   `json:"volume3"`
	Volume4           *float64   `json:"volume4"`
	Volume5           *float64   `json:"volume5"`
	SafStockQty       float64    `json:"saf_stock_qty"`
	SafStockUnitId    *string    `json:"saf_stock_unit_id"`
	SafStockUnitName  *string    `json:"saf_stock_unit_name"`
	MinStockQty       float64    `json:"min_stock_qty"`
	MinStockUnitId    *string    `json:"min_stock_unit_id"`
	MinStockUnitName  *string    `json:"min_stock_unit_name"`
	ExciseRate        float64    `json:"excise_rate"`
	ExciseTax         float64    `json:"excise_tax"`
	IsActive          bool       `json:"is_active"`
	ImageUrl          *string    `json:"image_url"`
	UpdatedBy         *int64     `json:"updated_by"`
	UpdatedByName     string     `json:"updated_by_name"`
	UpdatedAt         *time.Time `json:"updated_at"`
	Vat               *float64   `json:"vat"`
	VatBg             *float64   `json:"vat_bg"`
	VatLgPurch        *float64   `json:"vat_lg_purch"`
	VatLgSell         *float64   `json:"vat_lg_sell"`
	Cogs              *float64   `json:"cogs"`
	ProStatus         *int       `json:"pro_status"`
	ParentProductCode *string    `json:"parent_pro_code"`
	ParentProductName *string    `json:"parent_pro_name"`
}

/*
type ProductStock struct {
	UnitId1   string  `json:"unit_id1"`
	UnitId2   string  `json:"unit_id2"`
	UnitId3   string  `json:"unit_id3"`
	UnitId4   string  `json:"unit_id4"`
	UnitId5   string  `json:"unit_id5"`
	UnitName1 string  `json:"unit_name1"`
	UnitName2 string  `json:"unit_name2"`
	UnitName3 string  `json:"unit_name3"`
	UnitName4 string  `json:"unit_name4"`
	UnitName5 string  `json:"unit_name5"`
	ConvUnit1 float32 `json:"conv_unit1"`
	ConvUnit2 float32 `json:"conv_unit2"`
	ConvUnit3 float32 `json:"conv_unit3"`
	ConvUnit4 float32 `json:"conv_unit4"`
	ConvUnit5 float32 `json:"conv_unit5"`
	Qty1      float64 `json:"qty1"`
	Qty2      float64 `json:"qty2"`
	Qty3      float64 `json:"qty3"`
	Qty4      float64 `json:"qty4"`
	Qty5      float64 `json:"qty5"`
}
*/
