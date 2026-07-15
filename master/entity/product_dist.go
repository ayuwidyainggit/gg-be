package entity

import (
	"time"
)

type ProductDistQueryFilter struct {
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
	PcatId           []int  `query:"pcat_id"`
	SupCode          string `query:"sup_code"`
	DistPriceGroupId int    `query:"dist_price_group_id"`
}

type ProductDistResponse struct {
	ProductId       int        `json:"pro_id"`
	ProductCode     string     `json:"pro_code"`
	BarCode         *string    `json:"bar_code"`
	ProductName     string     `json:"pro_name"`
	PCatId          *int       `json:"pcat_id"`
	PCatCode        string     `json:"pcat_code"`
	PCatName        string     `json:"pcat_name"`
	ProductLineID   int        `json:"pl_id"`
	ProductLineCode string     `json:"pl_code"`
	ProductLineName string     `json:"pl_name"`
	BrandId         int        `json:"brand_id"`
	BrandCode       string     `json:"brand_code"`
	BrandName       string     `json:"brand_name"`
	Sbrand1         int        `json:"sbrand1_id"`
	Sbrand1Code     string     `json:"sbrand1_code"`
	Sbrand1Name     string     `json:"sbrand1_name"`
	Sbrand2         int        `json:"sbrand2_id"`
	Sbrand2Code     string     `json:"sbrand2_code"`
	Sbrand2Name     string     `json:"sbrand2_name"`
	FlavorId        int        `json:"flavor_id"`
	FlavorCode      string     `json:"flavor_code"`
	FlavorName      string     `json:"flavor_name"`
	PTypeId         int        `json:"ptype_id"`
	PTypeCode       string     `json:"ptype_code"`
	PTypeName       string     `json:"ptype_name"`
	PSizeId         int        `json:"psize_id"`
	PSizeCode       string     `json:"psize_code"`
	PSizeName       string     `json:"psize_name"`
	SupId           int        `json:"sup_id"`
	SupCode         string     `json:"sup_code"`
	SupName         string     `json:"sup_name"`
	PrincipalId     int        `json:"principal_id"`
	PrincipalCode   string     `json:"principal_code"`
	PrincipalName   string     `json:"principal_name"`
	CProId          int        `json:"c_pro_id"`
	CProCode        string     `json:"c_pro_code"`
	CProName        string     `json:"c_pro_name"`
	IsMainPro       bool       `json:"is_main_pro"`
	IsAlloc         bool       `json:"is_alloc"`
	SMweek1         int        `json:"s_mweek1"`
	SMweek2         int        `json:"s_mweek2"`
	SortNo          int        `json:"sort_no"`
	ItemNo          int        `json:"item_no"`
	UnitId1         string     `json:"unit_id1"`
	UnitId2         string     `json:"unit_id2"`
	UnitId3         string     `json:"unit_id3"`
	UnitId4         string     `json:"unit_id4"`
	UnitId5         string     `json:"unit_id5"`
	UnitName1       string     `json:"unit_name1"`
	UnitName2       string     `json:"unit_name2"`
	UnitName3       string     `json:"unit_name3"`
	UnitName4       string     `json:"unit_name4"`
	UnitName5       string     `json:"unit_name5"`
	ConvUnit2       float32    `json:"conv_unit2"`
	ConvUnit3       float32    `json:"conv_unit3"`
	ConvUnit4       float32    `json:"conv_unit4"`
	ConvUnit5       float32    `json:"conv_unit5"`
	IsBatch         bool       `json:"is_batch"`
	IsExpDate       bool       `json:"is_exp_date"`
	Length          float64    `json:"length"`
	Width           float64    `json:"width"`
	Height          float64    `json:"height"`
	Weight          float64    `json:"weight"`
	Volume          float64    `json:"volume"`
	MinStock        float64    `json:"min_stock"`
	MinStockStr     string     `json:"min_stock_str"`
	SafetyStock     float64    `json:"safety_stock"`
	SafetyStockStr  string     `json:"safety_stock_str"`
	PoFormula       int        `json:"po_formula"`
	ParentProId     int        `json:"parent_pro_id"`
	IsNewPro        bool       `json:"is_new_pro"`
	Vat             float64    `json:"vat"`
	VatBg           float64    `json:"vat_bg"`
	VatLgPurch      float64    `json:"vat_lg_purch"`
	VatLgSell       float64    `json:"vat_lg_sell"`
	ExciseRate      float64    `json:"excise_rate"`
	ExciseTax       float64    `json:"excise_tax"`
	IsActive        bool       `json:"is_active"`
	ImageUrl        *string    `json:"image_url"`
	UpdatedBy       *int64     `json:"updated_by"`
	UpdatedByName   string     `json:"updated_by_name"`
	UpdatedAt       *time.Time `json:"updated_at"`
	Cogs            float64    `json:"cogs"`
}

type ProductDistLookupResponse struct {
	ProductId       int      `json:"pro_id"`
	ProductCode     string   `json:"pro_code"`
	BarCode         *string  `json:"bar_code"`
	ProductName     string   `json:"pro_name"`
	PCatId          *int     `json:"pcat_id"`
	PCatCode        string   `json:"pcat_code"`
	PCatName        string   `json:"pcat_name"`
	ProductLineID   int      `json:"pl_id"`
	ProductLineCode string   `json:"pl_code"`
	ProductLineName string   `json:"pl_name"`
	BrandId         int      `json:"brand_id"`
	BrandCode       string   `json:"brand_code"`
	Sbrand1         int      `json:"sbrand1_id"`
	Sbrand1Code     string   `json:"sbrand1_code"`
	Sbrand1Name     string   `json:"sbrand1_name"`
	Sbrand2         int      `json:"sbrand2_id"`
	Sbrand2Code     string   `json:"sbrand2_code"`
	Sbrand2Name     string   `json:"sbrand2_name"`
	FlavorId        int      `json:"flavor_id"`
	FlavorCode      string   `json:"flavor_code"`
	FlavorName      string   `json:"flavor_name"`
	PTypeId         int      `json:"ptype_id"`
	PTypeCode       string   `json:"ptype_code"`
	PTypeName       string   `json:"ptype_name"`
	PSizeId         int      `json:"psize_id"`
	PSizeCode       string   `json:"psize_code"`
	PSizeName       string   `json:"psize_name"`
	SupId           int      `json:"sup_id"`
	SupCode         string   `json:"sup_code"`
	SupName         string   `json:"sup_name"`
	PrincipalId     int      `json:"principal_id"`
	PrincipalCode   string   `json:"principal_code"`
	PrincipalName   string   `json:"principal_name"`
	CProId          int      `json:"c_pro_id"`
	CProCode        string   `json:"c_pro_code"`
	CProName        string   `json:"c_pro_name"`
	IsMainPro       bool     `json:"is_main_pro"`
	IsAlloc         bool     `json:"is_alloc"`
	SMweek1         int      `json:"s_mweek1"`
	SMweek2         int      `json:"s_mweek2"`
	SortNo          int      `json:"sort_no"`
	ItemNo          int      `json:"item_no"`
	UnitId1         string   `json:"unit_id1"`
	UnitId2         string   `json:"unit_id2"`
	UnitId3         string   `json:"unit_id3"`
	UnitId4         string   `json:"unit_id4"`
	UnitId5         string   `json:"unit_id5"`
	UnitName1       string   `json:"unit_name1"`
	UnitName2       string   `json:"unit_name2"`
	UnitName3       string   `json:"unit_name3"`
	UnitName4       string   `json:"unit_name4"`
	UnitName5       string   `json:"unit_name5"`
	ConvUnit2       float32  `json:"conv_unit2"`
	ConvUnit3       float32  `json:"conv_unit3"`
	ConvUnit4       float32  `json:"conv_unit4"`
	ConvUnit5       float32  `json:"conv_unit5"`
	Margin          *float64 `json:"margin"`
	IsBatch         bool     `json:"is_batch"`
	IsExpDate       bool     `json:"is_exp_date"`
	Length          float64  `json:"length"`
	Width           float64  `json:"width"`
	Height          float64  `json:"height"`
	Weight          float64  `json:"weight"`
	Volume          float64  `json:"volume"`
	MinStock        float64  `json:"min_stock"`
	MinStockStr     string   `json:"min_stock_str"`
	SafetyStock     float64  `json:"safety_stock"`
	SafetyStockStr  string   `json:"safety_stock_str"`
	PoFormula       int      `json:"po_formula"`
	ParentProId     int      `json:"parent_pro_id"`
	IsNewPro        bool     `json:"is_new_pro"`
	Vat             float64  `json:"vat"`
	VatBg           float64  `json:"vat_bg"`
	VatLgPurch      float64  `json:"vat_lg_purch"`
	VatLgSell       float64  `json:"vat_lg_sell"`
	ExciseRate      float64  `json:"excise_rate"`
	ExciseTax       float64  `json:"excise_tax"`
	IsActive        bool     `json:"is_active"`
	ImageUrl        *string  `json:"image_url"`
	Cogs            float64  `json:"cogs"`
}

type ProductDistSearchResponse struct {
	ProductId   int     `json:"pro_id"`
	ProductCode string  `json:"pro_code"`
	BarCode     *string `json:"bar_code"`
	ProductName string  `json:"pro_name"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	UnitId4     string  `json:"unit_id4"`
	UnitId5     string  `json:"unit_id5"`
	UnitName1   string  `json:"unit_name1"`
	UnitName2   string  `json:"unit_name2"`
	UnitName3   string  `json:"unit_name3"`
	UnitName4   string  `json:"unit_name4"`
	UnitName5   string  `json:"unit_name5"`
	ConvUnit2   float32 `json:"conv_unit2"`
	ConvUnit3   float32 `json:"conv_unit3"`
	ConvUnit4   float32 `json:"conv_unit4"`
	ConvUnit5   float32 `json:"conv_unit5"`
	Cogs        float64 `json:"cogs"`
}

type CreateDistProductBody struct {
	CustId         string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy      int64   `json:"created_by" validate:"required"`
	BarCode        string  `json:"bar_code" validate:"max=50,alphanumericSpace"`
	ProductCode    string  `json:"pro_code" validate:"required,max=30,alphanumericSpace"`
	ProductName    string  `json:"pro_name" validate:"required,max=150"`
	PCatId         int     `json:"pcat_id" validate:"numeric"`
	Sbrand1        int     `json:"sbrand1_id" validate:"numeric"`
	Sbrand2        int     `json:"sbrand2_id" validate:"numeric"`
	FlavorId       int     `json:"flavor_id" validate:"numeric"`
	PTypeId        int     `json:"ptype_id" validate:"numeric"`
	PSizeId        int     `json:"psize_id" validate:"numeric"`
	SupId          int     `json:"sup_id" validate:"numeric"`
	PrincipalId    int     `json:"principal_id" validate:"numeric"`
	CProId         int     `json:"c_pro_id" validate:"numeric"`
	IsMainPro      bool    `json:"is_main_pro"`
	IsAlloc        bool    `json:"is_alloc"`
	SMweek1        int     `json:"s_mweek1" validate:"numeric"`
	SMweek2        int     `json:"s_mweek2" validate:"numeric"`
	SortNo         int     `json:"sort_no" validate:"numeric"`
	ItemNo         int     `json:"item_no" validate:"numeric"`
	UnitId1        string  `json:"unit_id1" validate:"max=5"`
	UnitId2        string  `json:"unit_id2" validate:"max=5"`
	UnitId3        string  `json:"unit_id3" validate:"max=5"`
	UnitId4        string  `json:"unit_id4" validate:"max=5"`
	UnitId5        string  `json:"unit_id5" validate:"max=5"`
	ConvUnit2      float32 `json:"conv_unit2" validate:"numeric"`
	ConvUnit3      float32 `json:"conv_unit3" validate:"numeric"`
	ConvUnit4      float32 `json:"conv_unit4" validate:"numeric"`
	ConvUnit5      float32 `json:"conv_unit5" validate:"numeric"`
	PurchPrice     float64 `json:"purch_price" validate:"numeric"`
	SellPrice1     float64 `json:"sell_price1" validate:"numeric"`
	SellPrice2     float64 `json:"sell_price2" validate:"numeric"`
	SellPrice3     float64 `json:"sell_price3" validate:"numeric"`
	SellPrice4     float64 `json:"sell_price4" validate:"numeric"`
	SellPrice5     float64 `json:"sell_price5" validate:"numeric"`
	Margin         float64 `json:"margin" validate:"numeric"`
	IsBatch        bool    `json:"is_batch"`
	IsExpDate      bool    `json:"is_exp_date"`
	Length         float64 `json:"length" validate:"numeric"`
	Weight         float64 `json:"weight" validate:"numeric"`
	Height         float64 `json:"height" validate:"numeric"`
	Volume         float64 `json:"volume" validate:"numeric"`
	MinStock       float64 `json:"min_stock" validate:"numeric"`
	MinStockStr    string  `json:"min_stock_str"`
	SafetyStock    float64 `json:"safety_stock"`
	SafetyStockStr string  `json:"safety_stock_str"`
	PoFormula      int     `json:"po_formula" validate:"numeric"`
	ParentProId    int     `json:"parent_pro_id" validate:"numeric"`
	IsNewPro       bool    `json:"is_new_pro"`
	Vat            float64 `json:"vat" validate:"numeric"`
	VatBg          float64 `json:"vat_bg" validate:"numeric"`
	VatLgPurch     float64 `json:"vat_lg_purch" validate:"numeric"`
	VatLgSell      float64 `json:"vat_lg_sell" validate:"numeric"`
	ExciseRate     float64 `json:"excise_rate" validate:"numeric"`
	ExciseTax      float64 `json:"excise_tax" validate:"numeric"`
	IsActive       bool    `json:"is_active"`
	Cogs           float64 `json:"cogs"`
}

type DetailProductDistParams struct {
	ProductId int64 `params:"pro_id" validate:"required"`
}

type UpdateProductDistParams struct {
	ProductId int64 `params:"pro_id" validate:"required"`
}

type DeleteProductDistParams struct {
	ProductId int64 `params:"pro_id" validate:"required"`
}

type UpdateProductDistRequest struct {
	CustId         string  `json:"cust_id" validate:"required,max=10"`
	UpdatedBy      int64   `json:"updated_by" validate:"required"`
	IsAlloc        *bool   `json:"is_alloc,omitempty"`
	SMweek1        int     `json:"s_mweek1,omitempty"`
	SMweek2        int     `json:"s_mweek2,omitempty"`
	MinStock       float64 `json:"min_stock,omitempty"`
	MinStockStr    string  `json:"min_stock_str,omitempty"`
	SafetyStock    float64 `json:"safety_stock,omitempty"`
	SafetyStockStr string  `json:"safety_stock_str,omitempty"`
	PoFormula      int     `json:"po_formula,omitempty"`
	IsNewPro       *bool   `json:"is_new_pro,omitempty"`
	Vat            float64 `json:"vat,omitempty"`
	VatBg          float64 `json:"vat_bg,omitempty"`
	VatLgPurch     float64 `json:"vat_lg_purch,omitempty"`
	VatLgSell      float64 `json:"vat_lg_sell,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
	Cogs           float64 `json:"cogs,omitempty"`
}
