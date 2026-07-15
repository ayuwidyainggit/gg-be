package entity

import (
	"mime/multipart"
	"time"
)

type ProductQueryFilter struct {
	CustId           string
	ParentCustId     string
	JwtDistributorId int64   // from JWT claim, to be used for authorization check, not for query filter
	DistributorID    int64   `query:"distributor_id"`
	DistPriceGroupId int     `query:"dist_price_group_id"`
	Page             int     `query:"page"`
	Limit            int     `query:"limit" validate:"required"`
	Query            string  `query:"q"`
	Mode             string  `query:"mode"`
	Sort             string  `query:"sort"`
	IsActive         *int    `query:"is_active"`
	IncludeDeleted   bool    `query:"include_deleted"`
	PcatId           []int   `query:"pcat_id"`
	ProID            []int   `query:"pro_id"`
	SupCode          string  `query:"sup_code"`
	SupID            []int   `query:"sup_id"`
	PrincipalID      []int   `query:"principal_id"`
	BrandID          []int   `query:"brand_id"`
	ProductLineId    []int   `query:"pl_id"`
	SubBrand1Id      []int   `query:"sbrand1_id"`
	OutletId         int     `query:"outlet_id"`
	OrderDate        string  `query:"order_date"`
	Format           string  `query:"format"`
	Status           string  `query:"status"`
	DistributorIds   []int64 `query:"distributor_ids"`
	AllDistributor   bool    `query:"all_distributor"`
}

type ProductReportQueryFilter struct {
	CustIDs   []string `json:"cust_id[]"`
	Query     string   `query:"q"`
	Page      int      `query:"page"`
	Limit     int      `query:"limit"`
	SortBy    string   `query:"sort_by"`
	SortOrder string   `query:"sort_order"`
}

type ProductReportResponse struct {
	CustID              string  `json:"cust_id" db:"cust_id"`
	ProductID           int64   `json:"pro_id" db:"pro_id"`
	ProductCode         string  `json:"pro_code" db:"pro_code"`
	ProductName         string  `json:"pro_name" db:"pro_name"`
	OriginalCustID      *string `json:"original_cust_id" db:"original_cust_id"`
	OriginalProductID   *int64  `json:"original_pro_id" db:"original_pro_id"`
	OriginalProductCode *string `json:"original_pro_code" db:"original_pro_code"`
	OriginalParentID    *int64  `json:"original_parent_pro_id" db:"original_parent_pro_id"`
	Type                string  `json:"type" db:"type"`
}

type ProductPrincipalQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Query        string `query:"q"`
}

type ProductCategoryQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Query        string `query:"q"`
}

type ProductBrandQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Query        string `query:"q"`
}

type ImportProductRequest struct {
	CustId        string
	ParentCustId  string
	DistributorId int64
	File          multipart.File
	FileName      string
	CreatedBy     int64
}

type ImportProductRow struct {
	HistoryId      int64  `db:"history_id"`
	CustId         string `db:"cust_id"`
	ProCode        string `db:"pro_code"`
	BarCode        string `db:"bar_code"`
	ProName        string `db:"pro_name"`
	PcatId         string `db:"pcat_id"`
	PcatCode       string `db:"pcat_code"`
	PcatName       string `db:"pcat_name"`
	PlId           string `db:"pl_id"`
	PlCode         string `db:"pl_code"`
	PlName         string `db:"pl_name"`
	BrandId        string `db:"brand_id"`
	BrandCode      string `db:"brand_code"`
	BrandName      string `db:"brand_name"`
	Sbrand1Id      string `db:"sbrand1_id"`
	Sbrand1Code    string `db:"sbrand1_code"`
	Sbrand1Name    string `db:"sbrand1_name"`
	Sbrand2Id      string `db:"sbrand2_id"`
	Sbrand2Code    string `db:"sbrand2_code"`
	Sbrand2Name    string `db:"sbrand2_name"`
	FlavorId       string `db:"flavor_id"`
	FlavorCode     string `db:"flavor_code"`
	FlavorName     string `db:"flavor_name"`
	PTypeId        string `db:"ptype_id"`
	PTypeCode      string `db:"ptype_code"`
	PTypeName      string `db:"ptype_name"`
	PSizeId        string `db:"psize_id"`
	PSizeCode      string `db:"psize_code"`
	PSizeName      string `db:"psize_name"`
	SupId          string `db:"sup_id"`
	SupCode        string `db:"sup_code"`
	SupName        string `db:"sup_name"`
	PrincipalId    string `db:"principal_id"`
	PrincipalCode  string `db:"principal_code"`
	PrincipalName  string `db:"principal_name"`
	CProId         string `db:"c_pro_id"`
	CProCode       string `db:"c_pro_code"`
	CProName       string `db:"c_pro_name"`
	IsMainPro      string `db:"is_main_pro"`
	SortNo         string `db:"sort_no"`
	ItemNo         string `db:"item_no"`
	UnitId1        string `db:"unit_id1"`
	UnitId2        string `db:"unit_id2"`
	UnitId3        string `db:"unit_id3"`
	UnitId4        string `db:"unit_id4"`
	UnitId5        string `db:"unit_id5"`
	UnitName1      string `db:"unit_name1"`
	UnitName2      string `db:"unit_name2"`
	UnitName3      string `db:"unit_name3"`
	ConvUnit2      string `db:"conv_unit2"`
	ConvUnit3      string `db:"conv_unit3"`
	ConvUnit4      string `db:"conv_unit4"`
	ConvUnit5      string `db:"conv_unit5"`
	IsBatch        string `db:"is_batch"`
	IsExpDate      string `db:"is_exp_date"`
	Length         string `db:"length"`
	Width          string `db:"width"`
	Height         string `db:"height"`
	Weight         string `db:"weight"`
	Volume         string `db:"volume"`
	ParentProId    string `db:"parent_pro_id"`
	IsNewPro       string `db:"is_new_pro"`
	PurchPrice1    string `db:"purch_price1"`
	PurchPrice2    string `db:"purch_price2"`
	PurchPrice3    string `db:"purch_price3"`
	PurchPrice4    string `db:"purch_price4"`
	PurchPrice5    string `db:"purch_price5"`
	SellPrice1     string `db:"sell_price1"`
	SellPrice2     string `db:"sell_price2"`
	SellPrice3     string `db:"sell_price3"`
	SellPrice4     string `db:"sell_price4"`
	SellPrice5     string `db:"sell_price5"`
	Length1        string `db:"length1"`
	Length2        string `db:"length2"`
	Length3        string `db:"length3"`
	Length4        string `db:"length4"`
	Length5        string `db:"length5"`
	Width1         string `db:"width1"`
	Width2         string `db:"width2"`
	Width3         string `db:"width3"`
	Width4         string `db:"width4"`
	Width5         string `db:"width5"`
	Height1        string `db:"height1"`
	Height2        string `db:"height2"`
	Height3        string `db:"height3"`
	Height4        string `db:"height4"`
	Height5        string `db:"height5"`
	Weight1        string `db:"weight1"`
	Weight2        string `db:"weight2"`
	Weight3        string `db:"weight3"`
	Weight4        string `db:"weight4"`
	Weight5        string `db:"weight5"`
	Volume1        string `db:"volume1"`
	Volume2        string `db:"volume2"`
	Volume3        string `db:"volume3"`
	Volume4        string `db:"volume4"`
	Volume5        string `db:"volume5"`
	SafStockQty    string `db:"saf_stock_qty"`
	SafStockUnitId string `db:"saf_stock_unit_id"`
	MinStockQty    string `db:"min_stock_qty"`
	MinStockUnitId string `db:"min_stock_unit_id"`
	ExciseRate     string `db:"excise_rate"`
	ExciseTax      string `db:"excise_tax"`
	IsActive       string `db:"is_active"`
	IsDel          string `db:"is_del"`
	ImageUrl       string `db:"image_url"`
	Vat            string `db:"vat"`
	VatBg          string `db:"vat_bg"`
	VatLgPurch     string `db:"vat_lg_purch"`
	VatLgSell      string `db:"vat_lg_sell"`
	Cogs           string `db:"cogs"`
	ProStatus      string `db:"pro_status"`
	ProCodeCoretax string `db:"pro_code_coretax"`
	ProNameCoretax string `db:"pro_name_coretax"`
	DistributorId  int64  `db:"distributor_id"`
	Level          int    `db:"level"`
	StatusInsert   string `db:"status_insert"`
	ErrorMessage   string `db:"error_message,omitempty"`
}

type ImportProductUpdateTemp struct {
	HistoryId     int64  `db:"history_id"`
	CustId        string `db:"cust_id"`
	BrandId       string `db:"brand_id"`
	BrandCode     string `db:"brand_code"`
	BrandName     string `db:"brand_name"`
	PcatId        string `db:"pcat_id"`
	PcatCode      string `db:"pcat_code"`
	PcatName      string `db:"pcat_name"`
	PlId          string `db:"pl_id"`
	PlCode        string `db:"pl_code"`
	PlName        string `db:"pl_name"`
	Sbrand1Id     string `db:"sbrand1_id"`
	Sbrand1Code   string `db:"sbrand1_code"`
	Sbrand1Name   string `db:"sbrand1_name"`
	Sbrand2Id     string `db:"sbrand2_id"`
	Sbrand2Code   string `db:"sbrand2_code"`
	Sbrand2Name   string `db:"sbrand2_name"`
	FlavorId      string `db:"flavor_id"`
	FlavorCode    string `db:"flavor_code"`
	FlavorName    string `db:"flavor_name"`
	UnitId        string `db:"unit_id"`
	UnitName      string `db:"unit_name"`
	PTypeId       string `db:"ptype_id"`
	PTypeCode     string `db:"ptype_code"`
	PTypeName     string `db:"ptype_name"`
	PSizeId       string `db:"psize_id"`
	PSizeCode     string `db:"psize_code"`
	PSizeName     string `db:"psize_name"`
	SupId         string `db:"sup_id"`
	SupCode       string `db:"sup_code"`
	SupName       string `db:"sup_name"`
	PrincipalId   string `db:"principal_id"`
	PrincipalCode string `db:"principal_code"`
	PrincipalName string `db:"principal_name"`
	CProId        string `db:"c_pro_id"`
	CProCode      string `db:"c_pro_code"`
	CProName      string `db:"c_pro_name"`
	// Product core fields
	ProId       string `db:"pro_id"`
	ProCode     string `db:"pro_code"`
	ProName     string `db:"pro_name"`
	Barcode     string `db:"bar_code"`
	Cogs        string `db:"cogs"`
	ProStatus   string `db:"pro_status"`
	IsActive    string `db:"is_active"`
	IsMainPro   string `db:"is_main_pro"`
	SortNo      string `db:"sort_no"`
	ItemNo      string `db:"item_no"`
	ParentProId string `db:"parent_pro_id"`

	// Units
	UnitId1   string `db:"unit_id1"`
	UnitName1 string `db:"unit_name1"`
	UnitId2   string `db:"unit_id2"`
	UnitName2 string `db:"unit_name2"`
	UnitId3   string `db:"unit_id3"`
	UnitName3 string `db:"unit_name3"`
	UnitId4   string `db:"unit_id4"`
	UnitId5   string `db:"unit_id5"`

	ConvUnit2 string `db:"conv_unit2"`
	ConvUnit3 string `db:"conv_unit3"`
	ConvUnit4 string `db:"conv_unit4"`
	ConvUnit5 string `db:"conv_unit5"`

	Length string `db:"length"`
	Width  string `db:"width"`
	Height string `db:"height"`
	Weight string `db:"weight"`
	Volume string `db:"volume"`

	Length1 string `db:"length1"`
	Length2 string `db:"length2"`
	Length3 string `db:"length3"`
	Length4 string `db:"length4"`
	Length5 string `db:"length5"`
	Width1  string `db:"width1"`
	Width2  string `db:"width2"`
	Width3  string `db:"width3"`
	Width4  string `db:"width4"`
	Width5  string `db:"width5"`
	Height1 string `db:"height1"`
	Height2 string `db:"height2"`
	Height3 string `db:"height3"`
	Height4 string `db:"height4"`
	Height5 string `db:"height5"`
	Weight1 string `db:"weight1"`
	Weight2 string `db:"weight2"`
	Weight3 string `db:"weight3"`
	Weight4 string `db:"weight4"`
	Weight5 string `db:"weight5"`
	Volume1 string `db:"volume1"`
	Volume2 string `db:"volume2"`
	Volume3 string `db:"volume3"`
	Volume4 string `db:"volume4"`
	Volume5 string `db:"volume5"`

	PurchPrice1 string `db:"purch_price1"`
	PurchPrice2 string `db:"purch_price2"`
	PurchPrice3 string `db:"purch_price3"`
	PurchPrice4 string `db:"purch_price4"`
	PurchPrice5 string `db:"purch_price5"`
	SellPrice1  string `db:"sell_price1"`
	SellPrice2  string `db:"sell_price2"`
	SellPrice3  string `db:"sell_price3"`
	SellPrice4  string `db:"sell_price4"`
	SellPrice5  string `db:"sell_price5"`
	// dst...

	SafStockQty      string `db:"saf_stock_qty"`
	SafStockUnitId   string `db:"saf_stock_unit_id"`
	SafStockUnitName string `db:"saf_stock_unit_name"`

	MinStockQty      string `db:"min_stock_qty"`
	MinStockUnitId   string `db:"min_stock_unit_id"`
	MinStockUnitName string `db:"min_stock_unit_name"`

	ExciseRate string `db:"excise_rate"`
	ExciseTax  string `db:"excise_tax"`
	Vat        string `db:"vat"`
	VatBg      string `db:"vat_bg"`
	VatLgPurch string `db:"vat_lg_purch"`
	VatLgSell  string `db:"vat_lg_sell"`

	IsBatch   string `db:"is_batch"`
	IsExpDate string `db:"is_exp_date"`
	IsNewPro  string `db:"is_new_pro"`
	IsDel     string `db:"is_del"`
	ImageURL  string `db:"image_url"`

	ProCodeCoretax string `db:"pro_code_coretax"`
	ProNameCoretax string `db:"pro_name_coretax"`

	StatusInsert string     `db:"status_insert"`
	ErrorMessage string     `db:"error_message"`
	CreatedAt    *time.Time `db:"created_at"`
}

type ProcessedProductRow struct {
	CustId         string     `db:"cust_id"`
	ProId          int        `db:"pro_id"`
	ProCode        string     `db:"pro_code"`
	BarCode        string     `db:"bar_code"`
	ProName        string     `db:"pro_name"`
	PcatId         int64      `db:"pcat_id"`
	BrandId        int64      `db:"brand_id"`
	Sbrand1Id      int64      `db:"sbrand1_id"`
	Sbrand2Id      int64      `db:"sbrand2_id"`
	FlavorId       int64      `db:"flavor_id"`
	PTypeId        int64      `db:"ptype_id"`
	PSizeId        int64      `db:"psize_id"`
	SupId          int64      `db:"sup_id"`
	PrincipalId    int64      `db:"principal_id"`
	CProId         int64      `db:"c_pro_id"`
	IsMainPro      bool       `db:"is_main_pro"`
	SortNo         int        `db:"sort_no"`
	ItemNo         int        `db:"item_no"`
	UnitId1        string     `db:"unit_id1"`
	UnitId2        string     `db:"unit_id2"`
	UnitId3        string     `db:"unit_id3"`
	ConvUnit2      float64    `db:"conv_unit2"`
	ConvUnit3      float64    `db:"conv_unit3"`
	Weight         float64    `db:"weight"`
	IsBatch        bool       `db:"is_batch"`
	IsExpDate      bool       `db:"is_exp_date"`
	Length         float64    `db:"length"`
	Width          float64    `db:"width"`
	Height         float64    `db:"height"`
	Volume         float64    `db:"volume"`
	ParentProId    int64      `db:"parent_pro_id"`
	IsNewPro       bool       `db:"is_new_pro"`
	Vat            float64    `db:"vat"`
	VatBg          float64    `db:"vat_bg"`
	VatLgPurch     float64    `db:"vat_lg_purch"`
	VatLgSell      float64    `db:"vat_lg_sell"`
	ExciseRate     float64    `db:"excise_rate"`
	ExciseTax      float64    `db:"excise_tax"`
	IsActive       bool       `db:"is_active"`
	CreatedBy      *int64     `db:"created_by"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedBy      *int64     `db:"updated_by"`
	UpdatedAt      *time.Time `db:"updated_at"`
	IsDel          bool       `db:"is_del"`
	DeletedBy      *int64     `db:"deleted_by"`
	DeletedAt      *time.Time `db:"deleted_at"`
	ImageUrl       string     `db:"image_url"`
	UnitId4        string     `db:"unit_id4"`
	UnitId5        string     `db:"unit_id5"`
	ConvUnit4      float64    `db:"conv_unit4"`
	ConvUnit5      float64    `db:"conv_unit5"`
	ProStatus      int        `db:"pro_status"`
	PurchPrice1    float64    `db:"purch_price1"`
	PurchPrice2    float64    `db:"purch_price2"`
	PurchPrice3    float64    `db:"purch_price3"`
	PurchPrice4    float64    `db:"purch_price4"`
	PurchPrice5    float64    `db:"purch_price5"`
	SellPrice1     float64    `db:"sell_price1"`
	SellPrice2     float64    `db:"sell_price2"`
	SellPrice3     float64    `db:"sell_price3"`
	SellPrice4     float64    `db:"sell_price4"`
	SellPrice5     float64    `db:"sell_price5"`
	Cogs           float64    `db:"cogs"`
	Weight1        float64    `db:"weight1"`
	Weight2        float64    `db:"weight2"`
	Weight3        float64    `db:"weight3"`
	Weight4        float64    `db:"weight4"`
	Weight5        float64    `db:"weight5"`
	Length1        float64    `db:"length1"`
	Length2        float64    `db:"length2"`
	Length3        float64    `db:"length3"`
	Length4        float64    `db:"length4"`
	Length5        float64    `db:"length5"`
	Width1         float64    `db:"width1"`
	Width2         float64    `db:"width2"`
	Width3         float64    `db:"width3"`
	Width4         float64    `db:"width4"`
	Width5         float64    `db:"width5"`
	Height1        float64    `db:"height1"`
	Height2        float64    `db:"height2"`
	Height3        float64    `db:"height3"`
	Height4        float64    `db:"height4"`
	Height5        float64    `db:"height5"`
	Volume1        float64    `db:"volume1"`
	Volume2        float64    `db:"volume2"`
	Volume3        float64    `db:"volume3"`
	Volume4        float64    `db:"volume4"`
	Volume5        float64    `db:"volume5"`
	SafStockUnitId string     `db:"saf_stock_unit_id"`
	SafStockQty    float64    `db:"saf_stock_qty"`
	MinStockUnitId string     `db:"min_stock_unit_id"`
	MinStockQty    float64    `db:"min_stock_qty"`
	ProCodeCoretax string     `db:"pro_code_coretax"`
	DistributorId  *int64     `db:"distributor_id"`
	Level          int        `db:"level"`
	Origin         string     `db:"origin"`
	AssignerUserID *int64     `db:"assigner_user_id"`
}

type ProductResponse struct {
	ProductId        int64      `json:"pro_id"`
	ProductCode      string     `json:"pro_code"`
	BarCode          *string    `json:"bar_code"`
	ProductName      string     `json:"pro_name"`
	PCatId           *int       `json:"pcat_id"`
	PCatCode         string     `json:"pcat_code"`
	PCatName         string     `json:"pcat_name"`
	ProductLineID    int        `json:"pl_id"`
	ProductLineCode  string     `json:"pl_code"`
	ProductLineName  string     `json:"pl_name"`
	BrandId          int        `json:"brand_id"`
	BrandCode        string     `json:"brand_code"`
	BrandName        string     `json:"brand_name"`
	Sbrand1          int        `json:"sbrand1_id"`
	Sbrand1Code      string     `json:"sbrand1_code"`
	Sbrand1Name      string     `json:"sbrand1_name"`
	Sbrand2          int        `json:"sbrand2_id"`
	Sbrand2Code      string     `json:"sbrand2_code"`
	Sbrand2Name      string     `json:"sbrand2_name"`
	FlavorId         int        `json:"flavor_id"`
	FlavorCode       string     `json:"flavor_code"`
	FlavorName       string     `json:"flavor_name"`
	PTypeId          int        `json:"ptype_id"`
	PTypeCode        string     `json:"ptype_code"`
	PTypeName        string     `json:"ptype_name"`
	PSizeId          int        `json:"psize_id"`
	PSizeCode        string     `json:"psize_code"`
	PSizeName        string     `json:"psize_name"`
	SupId            int        `json:"sup_id"`
	SupCode          string     `json:"sup_code"`
	SupName          string     `json:"sup_name"`
	PrincipalId      int        `json:"principal_id"`
	PrincipalCode    string     `json:"principal_code"`
	PrincipalName    string     `json:"principal_name"`
	CProId           int        `json:"c_pro_id"`
	CProCode         string     `json:"c_pro_code"`
	CProName         string     `json:"c_pro_name"`
	IsMainPro        bool       `json:"is_main_pro"`
	SortNo           int        `json:"sort_no"`
	ItemNo           int        `json:"item_no"`
	UnitId1          string     `json:"unit_id1"`
	UnitId2          string     `json:"unit_id2"`
	UnitId3          string     `json:"unit_id3"`
	UnitId4          string     `json:"unit_id4"`
	UnitId5          string     `json:"unit_id5"`
	ConvUnit2        float32    `json:"conv_unit2"`
	ConvUnit3        float32    `json:"conv_unit3"`
	ConvUnit4        float32    `json:"conv_unit4"`
	ConvUnit5        float32    `json:"conv_unit5"`
	IsBatch          bool       `json:"is_batch"`
	IsExpDate        bool       `json:"is_exp_date"`
	Length           float64    `json:"length"`
	Width            float64    `json:"width"`
	Height           float64    `json:"height"`
	Weight           float64    `json:"weight"`
	Volume           float64    `json:"volume"`
	PurchPrice1      float64    `json:"purch_price1"`
	PurchPrice2      float64    `json:"purch_price2"`
	PurchPrice3      float64    `json:"purch_price3"`
	PurchPrice4      float64    `json:"purch_price4"`
	PurchPrice5      float64    `json:"purch_price5"`
	SellPrice1       float64    `json:"sell_price1"`
	SellPrice2       float64    `json:"sell_price2"`
	SellPrice3       float64    `json:"sell_price3"`
	SellPrice4       float64    `json:"sell_price4"`
	SellPrice5       float64    `json:"sell_price5"`
	Length1          *float64   `json:"length1"`
	Length2          *float64   `json:"length2"`
	Length3          *float64   `json:"length3"`
	Length4          *float64   `json:"length4"`
	Length5          *float64   `json:"length5"`
	Width1           *float64   `json:"width1"`
	Width2           *float64   `json:"width2"`
	Width3           *float64   `json:"width3"`
	Width4           *float64   `json:"width4"`
	Width5           *float64   `json:"width5"`
	Height1          *float64   `json:"height1"`
	Height2          *float64   `json:"height2"`
	Height3          *float64   `json:"height3"`
	Height4          *float64   `json:"height4"`
	Height5          *float64   `json:"height5"`
	Weight1          *float64   `json:"weight1"`
	Weight2          *float64   `json:"weight2"`
	Weight3          *float64   `json:"weight3"`
	Weight4          *float64   `json:"weight4"`
	Weight5          *float64   `json:"weight5"`
	Volume1          *float64   `json:"volume1"`
	Volume2          *float64   `json:"volume2"`
	Volume3          *float64   `json:"volume3"`
	Volume4          *float64   `json:"volume4"`
	Volume5          *float64   `json:"volume5"`
	SafStockQty      float64    `json:"saf_stock_qty"`
	SafStockUnitId   *string    `json:"saf_stock_unit_id"`
	SafStockUnitName *string    `json:"saf_stock_unit_name"`
	MinStockQty      float64    `json:"min_stock_qty"`
	MinStockUnitId   *string    `json:"min_stock_unit_id"`
	MinStockUnitName *string    `json:"min_stock_unit_name"`
	ParentProId      int        `json:"parent_pro_id"`
	ExciseRate       float64    `json:"excise_rate"`
	ExciseTax        float64    `json:"excise_tax"`
	IsActive         bool       `json:"is_active"`
	ImageUrl         *string    `json:"image_url"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedByName    string     `json:"updated_by_name"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Vat              *float64   `json:"vat"`
	VatBg            *float64   `json:"vat_bg"`
	VatLgPurch       *float64   `json:"vat_lg_purch"`
	VatLgSell        *float64   `json:"vat_lg_sell"`
	Cogs             *float64   `json:"cogs"`
	ProStatus        *int       `json:"pro_status"`
	ProCodeCoreTax   *string    `json:"pro_code_coretax"`
	ProNameCoreTax   *string    `json:"pro_name_coretax"`
	DistributorId    *int64     `json:"distributor_id"`
	DistrbutorName   *string    `json:"distributor_name"`
	CreatedBy        *int64     `json:"created_by"`
	CreatedByName    *string    `json:"created_by_name"`
	CreatedAt        *time.Time `json:"created_at"`
}

type ProductExportResponse struct {
	ProductId        int64      `json:"pro_id"`
	ProductCode      string     `json:"pro_code"`
	BarCode          *string    `json:"bar_code"`
	ProductName      string     `json:"pro_name"`
	PCatId           *int       `json:"pcat_id"`
	PCatCode         string     `json:"pcat_code"`
	PCatName         string     `json:"pcat_name"`
	ProductLineID    int        `json:"pl_id"`
	ProductLineCode  string     `json:"pl_code"`
	ProductLineName  string     `json:"pl_name"`
	BrandId          int        `json:"brand_id"`
	BrandCode        string     `json:"brand_code"`
	BrandName        string     `json:"brand_name"`
	Sbrand1          int        `json:"sbrand1_id"`
	Sbrand1Code      string     `json:"sbrand1_code"`
	Sbrand1Name      string     `json:"sbrand1_name"`
	Sbrand2          int        `json:"sbrand2_id"`
	Sbrand2Code      string     `json:"sbrand2_code"`
	Sbrand2Name      string     `json:"sbrand2_name"`
	FlavorId         int        `json:"flavor_id"`
	FlavorCode       string     `json:"flavor_code"`
	FlavorName       string     `json:"flavor_name"`
	PTypeId          int        `json:"ptype_id"`
	PTypeCode        string     `json:"ptype_code"`
	PTypeName        string     `json:"ptype_name"`
	PSizeId          int        `json:"psize_id"`
	PSizeCode        string     `json:"psize_code"`
	PSizeName        string     `json:"psize_name"`
	SupId            int        `json:"sup_id"`
	SupCode          string     `json:"sup_code"`
	SupName          string     `json:"sup_name"`
	PrincipalId      int        `json:"principal_id"`
	PrincipalCode    string     `json:"principal_code"`
	PrincipalName    string     `json:"principal_name"`
	CProId           int        `json:"c_pro_id"`
	CProCode         string     `json:"c_pro_code"`
	CProName         string     `json:"c_pro_name"`
	IsMainPro        bool       `json:"is_main_pro"`
	SortNo           int        `json:"sort_no"`
	ItemNo           int        `json:"item_no"`
	UnitId1          string     `json:"unit_id1"`
	UnitName1        string     `json:"unit_name1"`
	UnitId2          string     `json:"unit_id2"`
	UnitName2        string     `json:"unit_name2"`
	UnitId3          string     `json:"unit_id3"`
	UnitName3        string     `json:"unit_name3"`
	UnitId4          string     `json:"unit_id4"`
	UnitId5          string     `json:"unit_id5"`
	ConvUnit2        float32    `json:"conv_unit2"`
	ConvUnit3        float32    `json:"conv_unit3"`
	ConvUnit4        float32    `json:"conv_unit4"`
	ConvUnit5        float32    `json:"conv_unit5"`
	IsBatch          bool       `json:"is_batch"`
	IsExpDate        bool       `json:"is_exp_date"`
	Length           float64    `json:"length"`
	Width            float64    `json:"width"`
	Height           float64    `json:"height"`
	Weight           float64    `json:"weight"`
	Volume           float64    `json:"volume"`
	PurchPrice1      float64    `json:"purch_price1"`
	PurchPrice2      float64    `json:"purch_price2"`
	PurchPrice3      float64    `json:"purch_price3"`
	PurchPrice4      float64    `json:"purch_price4"`
	PurchPrice5      float64    `json:"purch_price5"`
	SellPrice1       float64    `json:"sell_price1"`
	SellPrice2       float64    `json:"sell_price2"`
	SellPrice3       float64    `json:"sell_price3"`
	SellPrice4       float64    `json:"sell_price4"`
	SellPrice5       float64    `json:"sell_price5"`
	Length1          *float64   `json:"length1"`
	Length2          *float64   `json:"length2"`
	Length3          *float64   `json:"length3"`
	Length4          *float64   `json:"length4"`
	Length5          *float64   `json:"length5"`
	Width1           *float64   `json:"width1"`
	Width2           *float64   `json:"width2"`
	Width3           *float64   `json:"width3"`
	Width4           *float64   `json:"width4"`
	Width5           *float64   `json:"width5"`
	Height1          *float64   `json:"height1"`
	Height2          *float64   `json:"height2"`
	Height3          *float64   `json:"height3"`
	Height4          *float64   `json:"height4"`
	Height5          *float64   `json:"height5"`
	Weight1          *float64   `json:"weight1"`
	Weight2          *float64   `json:"weight2"`
	Weight3          *float64   `json:"weight3"`
	Weight4          *float64   `json:"weight4"`
	Weight5          *float64   `json:"weight5"`
	Volume1          *float64   `json:"volume1"`
	Volume2          *float64   `json:"volume2"`
	Volume3          *float64   `json:"volume3"`
	Volume4          *float64   `json:"volume4"`
	Volume5          *float64   `json:"volume5"`
	SafStockQty      float64    `json:"saf_stock_qty"`
	SafStockUnitId   *string    `json:"saf_stock_unit_id"`
	SafStockUnitName *string    `json:"saf_stock_unit_name"`
	MinStockQty      float64    `json:"min_stock_qty"`
	MinStockUnitId   *string    `json:"min_stock_unit_id"`
	MinStockUnitName *string    `json:"min_stock_unit_name"`
	ParentProId      int        `json:"parent_pro_id"`
	ExciseRate       float64    `json:"excise_rate"`
	ExciseTax        float64    `json:"excise_tax"`
	IsActive         bool       `json:"is_active"`
	IsDel            bool       `json:"is_del"`
	ImageUrl         *string    `json:"image_url"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedByName    string     `json:"updated_by_name"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Vat              *float64   `json:"vat"`
	VatBg            *float64   `json:"vat_bg"`
	VatLgPurch       *float64   `json:"vat_lg_purch"`
	VatLgSell        *float64   `json:"vat_lg_sell"`
	Cogs             *float64   `json:"cogs"`
	ProStatus        *int       `json:"pro_status"`
	ProCodeCoreTax   *string    `json:"pro_code_coretax"`
	ProNameCoreTax   *string    `json:"pro_name_coretax"`
	DistributorId    *int64     `json:"distributor_id"`
	DistrbutorName   *string    `json:"distributor_name"`
	CreatedBy        *int64     `json:"created_by"`
	CreatedByName    *string    `json:"created_by_name"`
	CreatedAt        *time.Time `json:"created_at"`
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
	UnitName1         string     `json:"unit_name1"`
	UnitName2         string     `json:"unit_name2"`
	UnitName3         string     `json:"unit_name3"`
	UnitIdCoreTax1    *string    `json:"unit_id_coretax1"`
	UnitIdCoreTax2    *string    `json:"unit_id_coretax2"`
	UnitIdCoreTax3    *string    `json:"unit_id_coretax3"`
	UnitNameCoreTax1  *string    `json:"unit_name_coretax1"`
	UnitNameCoreTax2  *string    `json:"unit_name_coretax2"`
	UnitNameCoreTax3  *string    `json:"unit_name_coretax3"`
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
	ProCodeCoreTax    *string    `json:"pro_code_coretax"`
	ProNameCoreTax    *string    `json:"pro_name_coretax"`
	DistributorId     *int64     `json:"distributor_id"`
	DistrbutorName    *string    `json:"distributor_name"`
	CreatedBy         *int64     `json:"created_by"`
	CreatedByName     *string    `json:"created_by_name"`
	CreatedAt         *time.Time `json:"created_at"`
}

type ProductLookupResponse struct {
	ProductId        int      `json:"pro_id"`
	ProductCode      string   `json:"pro_code"`
	BarCode          *string  `json:"bar_code"`
	ProductName      string   `json:"pro_name"`
	PCatId           *int     `json:"pcat_id"`
	PCatCode         string   `json:"pcat_code"`
	PCatName         string   `json:"pcat_name"`
	ProductLineID    int      `json:"pl_id"`
	ProductLineCode  string   `json:"pl_code"`
	ProductLineName  string   `json:"pl_name"`
	BrandId          int      `json:"brand_id"`
	BrandCode        string   `json:"brand_code"`
	Sbrand1          int      `json:"sbrand1_id"`
	Sbrand1Code      string   `json:"sbrand1_code"`
	Sbrand1Name      string   `json:"sbrand1_name"`
	Sbrand2          int      `json:"sbrand2_id"`
	Sbrand2Code      string   `json:"sbrand2_code"`
	Sbrand2Name      string   `json:"sbrand2_name"`
	FlavorId         int      `json:"flavor_id"`
	FlavorCode       string   `json:"flavor_code"`
	FlavorName       string   `json:"flavor_name"`
	PTypeId          int      `json:"ptype_id"`
	PTypeCode        string   `json:"ptype_code"`
	PTypeName        string   `json:"ptype_name"`
	PSizeId          int      `json:"psize_id"`
	PSizeCode        string   `json:"psize_code"`
	PSizeName        string   `json:"psize_name"`
	SupId            int      `json:"sup_id"`
	SupCode          string   `json:"sup_code"`
	SupName          string   `json:"sup_name"`
	PrincipalId      int      `json:"principal_id"`
	PrincipalCode    string   `json:"principal_code"`
	PrincipalName    string   `json:"principal_name"`
	CProId           int      `json:"c_pro_id"`
	CProCode         string   `json:"c_pro_code"`
	CProName         string   `json:"c_pro_name"`
	IsMainPro        bool     `json:"is_main_pro"`
	SortNo           int      `json:"sort_no"`
	ItemNo           int      `json:"item_no"`
	UnitId1          string   `json:"unit_id1"`
	UnitId2          string   `json:"unit_id2"`
	UnitId3          string   `json:"unit_id3"`
	UnitId4          string   `json:"unit_id4"`
	UnitId5          string   `json:"unit_id5"`
	ConvUnit2        float32  `json:"conv_unit2"`
	ConvUnit3        float32  `json:"conv_unit3"`
	ConvUnit4        float32  `json:"conv_unit4"`
	ConvUnit5        float32  `json:"conv_unit5"`
	IsBatch          bool     `json:"is_batch"`
	IsExpDate        bool     `json:"is_exp_date"`
	Length           float64  `json:"length"`
	Width            float64  `json:"width"`
	Height           float64  `json:"height"`
	Weight           float64  `json:"weight"`
	Volume           float64  `json:"volume"`
	PurchPrice1      float64  `json:"purch_price1"`
	PurchPrice2      float64  `json:"purch_price2"`
	PurchPrice3      float64  `json:"purch_price3"`
	PurchPrice4      float64  `json:"purch_price4"`
	PurchPrice5      float64  `json:"purch_price5"`
	SellPrice1       float64  `json:"sell_price1"`
	SellPrice2       float64  `json:"sell_price2"`
	SellPrice3       float64  `json:"sell_price3"`
	SellPrice4       float64  `json:"sell_price4"`
	SellPrice5       float64  `json:"sell_price5"`
	Length1          *float64 `json:"length1"`
	Length2          *float64 `json:"length2"`
	Length3          *float64 `json:"length3"`
	Length4          *float64 `json:"length4"`
	Length5          *float64 `json:"length5"`
	Width1           *float64 `json:"width1"`
	Width2           *float64 `json:"width2"`
	Width3           *float64 `json:"width3"`
	Width4           *float64 `json:"width4"`
	Width5           *float64 `json:"width5"`
	Height1          *float64 `json:"height1"`
	Height2          *float64 `json:"height2"`
	Height3          *float64 `json:"height3"`
	Height4          *float64 `json:"height4"`
	Height5          *float64 `json:"height5"`
	Weight1          *float64 `json:"weight1"`
	Weight2          *float64 `json:"weight2"`
	Weight3          *float64 `json:"weight3"`
	Weight4          *float64 `json:"weight4"`
	Weight5          *float64 `json:"weight5"`
	Volume1          *float64 `json:"volume1"`
	Volume2          *float64 `json:"volume2"`
	Volume3          *float64 `json:"volume3"`
	Volume4          *float64 `json:"volume4"`
	Volume5          *float64 `json:"volume5"`
	SafStockQty      float64  `json:"saf_stock_qty"`
	SafStockUnitId   *string  `json:"saf_stock_unit_id"`
	SafStockUnitName *string  `json:"saf_stock_unit_name"`
	MinStockQty      float64  `json:"min_stock_qty"`
	MinStockUnitId   *string  `json:"min_stock_unit_id"`
	MinStockUnitName *string  `json:"min_stock_unit_name"`
	ParentProId      int      `json:"parent_pro_id"`
	ExciseRate       float64  `json:"excise_rate"`
	ExciseTax        float64  `json:"excise_tax"`
	IsActive         bool     `json:"is_active"`
	ImageUrl         *string  `json:"image_url"`
	Vat              *float64 `json:"vat"`
	VatBg            *float64 `json:"vat_bg"`
	VatLgPurch       *float64 `json:"vat_lg_purch"`
	VatLgSell        *float64 `json:"vat_lg_sell"`
	Cogs             *float64 `json:"cogs"`
	ProStatus        *int     `json:"pro_status"`
}

type ProductSearchResponse struct {
	ProductId   int     `json:"pro_id"`
	ProductCode string  `json:"pro_code"`
	BarCode     *string `json:"bar_code"`
	ProductName string  `json:"pro_name"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	UnitId4     string  `json:"unit_id4"`
	UnitId5     string  `json:"unit_id5"`
	ConvUnit2   float32 `json:"conv_unit2"`
	ConvUnit3   float32 `json:"conv_unit3"`
	ConvUnit4   float32 `json:"conv_unit4"`
	ConvUnit5   float32 `json:"conv_unit5"`
}

type CreateProductBody struct {
	DistributorId  int64    `json:"distributor_id"` // From JWT
	CustId         string   `json:"cust_id" validate:"required,max=10"`
	CreatedBy      int64    `json:"created_by" validate:"required"`
	BarCode        string   `json:"bar_code" validate:"max=50"`
	ProductCode    string   `json:"pro_code" validate:"required,max=30,alphanumericSpaceDash"`
	ProductName    string   `json:"pro_name" validate:"required,max=150"`
	PCatId         int      `json:"pcat_id" validate:"numeric"`
	Sbrand1        int      `json:"sbrand1_id" validate:"numeric"`
	Sbrand2        int      `json:"sbrand2_id" validate:"numeric"`
	FlavorId       int      `json:"flavor_id" validate:"numeric"`
	PTypeId        int      `json:"ptype_id" validate:"numeric"`
	PSizeId        int      `json:"psize_id" validate:"numeric"`
	SupId          int      `json:"sup_id" validate:"numeric"`
	PrincipalId    int      `json:"principal_id" validate:"numeric"`
	CProId         int      `json:"c_pro_id" validate:"numeric"`
	IsMainPro      bool     `json:"is_main_pro"`
	SortNo         int      `json:"sort_no" validate:"numeric"`
	ItemNo         int      `json:"item_no" validate:"numeric"`
	UnitId1        string   `json:"unit_id1" validate:"max=5"`
	UnitId2        string   `json:"unit_id2" validate:"max=5"`
	UnitId3        string   `json:"unit_id3" validate:"max=5"`
	UnitId4        string   `json:"unit_id4" validate:"max=5"`
	UnitId5        string   `json:"unit_id5" validate:"max=5"`
	ConvUnit2      float32  `json:"conv_unit2" validate:"numeric"`
	ConvUnit3      float32  `json:"conv_unit3" validate:"numeric"`
	ConvUnit4      float32  `json:"conv_unit4" validate:"numeric"`
	ConvUnit5      float32  `json:"conv_unit5" validate:"numeric"`
	IsBatch        bool     `json:"is_batch"`
	IsExpDate      bool     `json:"is_exp_date"`
	Length         float64  `json:"length" validate:"numeric"`
	Width          float64  `json:"width" validate:"numeric"`
	Height         float64  `json:"height" validate:"numeric"`
	Weight         float64  `json:"weight" validate:"numeric"`
	Volume         float64  `json:"volume" validate:"numeric"`
	PurchPrice1    float64  `json:"purch_price1" validate:"numeric"`
	PurchPrice2    float64  `json:"purch_price2" validate:"numeric"`
	PurchPrice3    float64  `json:"purch_price3" validate:"numeric"`
	PurchPrice4    float64  `json:"purch_price4" validate:"numeric"`
	PurchPrice5    float64  `json:"purch_price5" validate:"numeric"`
	SellPrice1     float64  `json:"sell_price1" validate:"numeric"`
	SellPrice2     float64  `json:"sell_price2" validate:"numeric"`
	SellPrice3     float64  `json:"sell_price3" validate:"numeric"`
	SellPrice4     float64  `json:"sell_price4" validate:"numeric"`
	SellPrice5     float64  `json:"sell_price5" validate:"numeric"`
	Length1        float64  `json:"length1" validate:"numeric"`
	Length2        float64  `json:"length2" validate:"numeric"`
	Length3        float64  `json:"length3" validate:"numeric"`
	Length4        float64  `json:"length4" validate:"numeric"`
	Length5        float64  `json:"length5" validate:"numeric"`
	Width1         float64  `json:"width1" validate:"numeric"`
	Width2         float64  `json:"width2" validate:"numeric"`
	Width3         float64  `json:"width3" validate:"numeric"`
	Width4         float64  `json:"width4" validate:"numeric"`
	Width5         float64  `json:"width5" validate:"numeric"`
	Height1        float64  `json:"height1" validate:"numeric"`
	Height2        float64  `json:"height2" validate:"numeric"`
	Height3        float64  `json:"height3" validate:"numeric"`
	Height4        float64  `json:"height4" validate:"numeric"`
	Height5        float64  `json:"height5" validate:"numeric"`
	Weight1        float64  `json:"weight1" validate:"numeric"`
	Weight2        float64  `json:"weight2" validate:"numeric"`
	Weight3        float64  `json:"weight3" validate:"numeric"`
	Weight4        float64  `json:"weight4" validate:"numeric"`
	Weight5        float64  `json:"weight5" validate:"numeric"`
	Volume1        float64  `json:"volume1" validate:"numeric"`
	Volume2        float64  `json:"volume2" validate:"numeric"`
	Volume3        float64  `json:"volume3" validate:"numeric"`
	Volume4        float64  `json:"volume4" validate:"numeric"`
	Volume5        float64  `json:"volume5" validate:"numeric"`
	SafStockQty    float64  `json:"saf_stock_qty" validate:"numeric"`
	SafStockUnitId string   `json:"saf_stock_unit_id"`
	MinStockQty    float64  `json:"min_stock_qty" validate:"numeric"`
	MinStockUnitId string   `json:"min_stock_unit_id"`
	ParentProId    int      `json:"parent_pro_id" validate:"numeric"`
	IsNewPro       bool     `json:"is_new_pro"`
	ExciseRate     float64  `json:"excise_rate" validate:"numeric,max=1000000000"`
	ExciseTax      float64  `json:"excise_tax" validate:"numeric,max=1000000"`
	IsActive       bool     `json:"is_active"`
	ImageUrl       string   `json:"image_url,omitempty"`
	Vat            *float64 `json:"vat"`
	VatBg          *float64 `json:"vat_bg"`
	VatLgPurch     *float64 `json:"vat_lg_purch"`
	VatLgSell      *float64 `json:"vat_lg_sell"`
	Cogs           *float64 `json:"cogs"`
	ProStatus      *int     `json:"pro_status"`
	ProCodeCoreTax *string  `json:"pro_code_coretax"`
}

type BulkProductBody struct {
	Products []CreateProductBody `json:"products" validate:"min=1"`
}

type BulkProductResponse struct {
	Products []ProductResponse `json:"products"`
}

type DetailProductParams struct {
	ProductId     int64 `params:"pro_id" validate:"required"`
	CustID        string
	DistributorID int64 `params:"distributor_id"`
	ParentCustID  string
}

type UpdateProductParams struct {
	ProductId int64 `params:"pro_id" validate:"required"`
}

type DeleteProductParams struct {
	ProductId int64 `params:"pro_id" validate:"required"`
}

type DeleteMultipleProductBody struct {
	ProductId []int64 `json:"pro_id" validate:"min=1"`
}

type UpdateProductRequest struct {
	CustId         string   `json:"cust_id" validate:"required,max=10"`
	UpdatedBy      int64    `json:"updated_by" validate:"required"`
	BarCode        string   `json:"bar_code,omitempty" validate:"omitempty,max=50"`
	ProductCode    string   `json:"pro_code,omitempty" validate:"omitempty,max=30,alphanumericSpaceDash"`
	ProductName    string   `json:"pro_name,omitempty" validate:"omitempty,max=150"`
	PCatId         int      `json:"pcat_id,omitempty"`
	Sbrand1        int      `json:"sbrand1_id,omitempty"`
	Sbrand2        int      `json:"sbrand2_id,omitempty"`
	FlavorId       int      `json:"flavor_id,omitempty"`
	PTypeId        int      `json:"ptype_id,omitempty"`
	PSizeId        int      `json:"psize_id,omitempty"`
	SupId          int      `json:"sup_id,omitempty"`
	PrincipalId    int      `json:"principal_id,omitempty"`
	CProId         int      `json:"c_pro_id,omitempty"`
	IsMainPro      *bool    `json:"is_main_pro,omitempty"`
	SortNo         int      `json:"sort_no,omitempty"`
	ItemNo         int      `json:"item_no,omitempty"`
	UnitId1        string   `json:"unit_id1,omitempty"`
	UnitId2        string   `json:"unit_id2,omitempty"`
	UnitId3        string   `json:"unit_id3,omitempty"`
	UnitId4        string   `json:"unit_id4,omitempty"`
	UnitId5        string   `json:"unit_id5,omitempty"`
	ConvUnit2      float32  `json:"conv_unit2,omitempty"`
	ConvUnit3      float32  `json:"conv_unit3,omitempty"`
	ConvUnit4      float32  `json:"conv_unit4,omitempty"`
	ConvUnit5      float32  `json:"conv_unit5,omitempty"`
	IsBatch        *bool    `json:"is_batch,omitempty"`
	IsExpDate      *bool    `json:"is_exp_date,omitempty"`
	Length         float64  `json:"length,omitempty"`
	Width          float64  `json:"width,omitempty"`
	Height         float64  `json:"height,omitempty"`
	Weight         float64  `json:"weight,omitempty"`
	Volume         float64  `json:"volume,omitempty"`
	PurchPrice1    float64  `json:"purch_price1" validate:"numeric"`
	PurchPrice2    float64  `json:"purch_price2" validate:"numeric"`
	PurchPrice3    float64  `json:"purch_price3" validate:"numeric"`
	PurchPrice4    float64  `json:"purch_price4" validate:"numeric"`
	PurchPrice5    float64  `json:"purch_price5" validate:"numeric"`
	SellPrice1     float64  `json:"sell_price1" validate:"numeric"`
	SellPrice2     float64  `json:"sell_price2" validate:"numeric"`
	SellPrice3     float64  `json:"sell_price3" validate:"numeric"`
	SellPrice4     float64  `json:"sell_price4" validate:"numeric"`
	SellPrice5     float64  `json:"sell_price5" validate:"numeric"`
	Length1        float64  `json:"length1" validate:"numeric"`
	Length2        float64  `json:"length2" validate:"numeric"`
	Length3        float64  `json:"length3" validate:"numeric"`
	Length4        float64  `json:"length4" validate:"numeric"`
	Length5        float64  `json:"length5" validate:"numeric"`
	Width1         float64  `json:"width1" validate:"numeric"`
	Width2         float64  `json:"width2" validate:"numeric"`
	Width3         float64  `json:"width3" validate:"numeric"`
	Width4         float64  `json:"width4" validate:"numeric"`
	Width5         float64  `json:"width5" validate:"numeric"`
	Height1        float64  `json:"height1" validate:"numeric"`
	Height2        float64  `json:"height2" validate:"numeric"`
	Height3        float64  `json:"height3" validate:"numeric"`
	Height4        float64  `json:"height4" validate:"numeric"`
	Height5        float64  `json:"height5" validate:"numeric"`
	Weight1        float64  `json:"weight1" validate:"numeric"`
	Weight2        float64  `json:"weight2" validate:"numeric"`
	Weight3        float64  `json:"weight3" validate:"numeric"`
	Weight4        float64  `json:"weight4" validate:"numeric"`
	Weight5        float64  `json:"weight5" validate:"numeric"`
	Volume1        float64  `json:"volume1" validate:"numeric"`
	Volume2        float64  `json:"volume2" validate:"numeric"`
	Volume3        float64  `json:"volume3" validate:"numeric"`
	Volume4        float64  `json:"volume4" validate:"numeric"`
	Volume5        float64  `json:"volume5" validate:"numeric"`
	SafStockQty    float64  `json:"saf_stock_qty" validate:"numeric"`
	SafStockUnitId string   `json:"saf_stock_unit_id"`
	MinStockQty    float64  `json:"min_stock_qty" validate:"numeric"`
	MinStockUnitId string   `json:"min_stock_unit_id"`
	ParentProId    int      `json:"parent_pro_id"`
	ExciseRate     float64  `json:"excise_rate,omitempty" validate:"omitempty,max=1000000000"`
	ExciseTax      float64  `json:"excise_tax,omitempty" validate:"omitempty,max=1000000"`
	IsActive       *bool    `json:"is_active,omitempty"`
	ImageUrl       string   `json:"image_url"`
	Vat            *float64 `json:"vat,omitempty"`
	VatBg          *float64 `json:"vat_bg,omitempty"`
	VatLgPurch     *float64 `json:"vat_lg_purch,omitempty"`
	VatLgSell      *float64 `json:"vat_lg_sell,omitempty"`
	Cogs           *float64 `json:"cogs"`
	ProStatus      *int     `json:"pro_status,omitempty"`
	ProCodeCoretax *string  `json:"pro_code_coretax"`
}

type ProductLookupDistPrice struct {
	ProductId   int     `json:"pro_id"`
	ProductCode string  `json:"pro_code"`
	ProductName string  `json:"pro_name"`
	UnitId1     string  `json:"unit_id1"`
	UnitId2     string  `json:"unit_id2"`
	UnitId3     string  `json:"unit_id3"`
	UnitId4     string  `json:"unit_id4"`
	UnitId5     string  `json:"unit_id5"`
	ConvUnit2   float32 `json:"conv_unit2"`
	ConvUnit3   float32 `json:"conv_unit3"`
	ConvUnit4   float32 `json:"conv_unit4"`
	ConvUnit5   float32 `json:"conv_unit5"`
	// DistPriceId int     `json:"dist_price_id"`
	PurchPrice1     float64  `json:"purch_price1"`
	PurchPrice2     float64  `json:"purch_price2"`
	PurchPrice3     float64  `json:"purch_price3"`
	PurchPrice4     float64  `json:"purch_price4"`
	PurchPrice5     float64  `json:"purch_price5"`
	SellPrice1      float64  `json:"sell_price1"`
	SellPrice2      float64  `json:"sell_price2"`
	SellPrice3      float64  `json:"sell_price3"`
	SellPrice4      float64  `json:"sell_price4"`
	SellPrice5      float64  `json:"sell_price5"`
	Vat             *float64 `json:"vat"`
	ProductLineID   int      `json:"pl_id"`
	ProductLineCode string   `json:"pl_code"`
	ProductLineName string   `json:"pl_name"`
	BrandId         int      `json:"brand_id"`
	BrandCode       string   `json:"brand_code"`
	Sbrand1         int      `json:"sbrand1_id"`
	Sbrand1Code     string   `json:"sbrand1_code"`
	Sbrand1Name     string   `json:"sbrand1_name"`
}

type ProductCategoryList struct {
	PCatId   int    `json:"pcat_id"`
	PCatCode string `json:"pcat_code"`
	PCatName string `json:"pcat_name"`
}

type ProductBrandList struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

// Automatic Replenishment Product structures
type AutomaticReplenishmentProductQueryFilter struct {
	CustId           string `query:"cust_id"`
	ParentCustId     string
	JwtDistributorId int64
	DistributorID    []int64 `query:"distributor_id"`
	ProID            []int64 `query:"pro_id"`
	Page             int     `query:"page" validate:"required,min=1"`
	Limit            int     `query:"limit" validate:"required,min=1,max=100"`
	Query            string  `query:"q"`
	Sort             string  `query:"sort"`
	Format           string  `query:"format"`
}

type DetailAutomaticReplenishmentProductParams struct {
	CustId       string
	ParentCustId string
	Id           int64 `params:"id" validate:"required,min=1"`
}

type CreateAutomaticReplenishmentProductRequest struct {
	DistributorId   int64  `json:"distributor_id" validate:"required,min=1"`
	ProId           int64  `json:"pro_id" validate:"required,min=1"`
	LimitAction     string `json:"limit_action" validate:"required,oneof=RESTRICTED WARNING UNRESTRICTED"`
	MaxOrderQty     int    `json:"max_order_qty" validate:"min=0"`
	MaxOrderType    string `json:"max_order_type" validate:"omitempty,oneof=S M L"`
	MinStockQty     int    `json:"min_stock" validate:"min=0"`
	MinStockType    string `json:"min_stock_type" validate:"omitempty,oneof=S M L"`
	SafetyStockQty  int    `json:"saf_stock" validate:"min=0"`
	SafetyStockType string `json:"saf_stock_type" validate:"omitempty,oneof=S M L"`
	MinOrderQty     int    `json:"min_order_qty" validate:"min=0"`
	MinOrderType    string `json:"min_order_type" validate:"omitempty,oneof=S M L"`
}

type UpdateAutomaticReplenishmentProductRequest struct {
	DistributorId   int64  `json:"distributor_id" validate:"required,min=1"`
	ProId           int64  `json:"pro_id" validate:"required,min=1"`
	LimitAction     string `json:"limit_action" validate:"required,oneof=RESTRICTED WARNING UNRESTRICTED"`
	MaxOrderQty     int    `json:"max_order_qty" validate:"min=0"`
	MaxOrderType    string `json:"max_order_type" validate:"omitempty,oneof=S M L"`
	MinStockQty     int    `json:"min_stock" validate:"min=0"`
	MinStockType    string `json:"min_stock_type" validate:"omitempty,oneof=S M L"`
	SafetyStockQty  int    `json:"saf_stock" validate:"min=0"`
	SafetyStockType string `json:"saf_stock_type" validate:"omitempty,oneof=S M L"`
	MinOrderQty     int    `json:"min_order_qty" validate:"min=0"`
	MinOrderType    string `json:"min_order_type" validate:"omitempty,oneof=S M L"`
}

type AutomaticReplenishmentProductResponse struct {
	CustId          string  `json:"cust_id"`
	Id              int64   `json:"id"`
	ProId           int64   `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	DistributorId   int64   `json:"distributor_id"`
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	LimitAction     string  `json:"limit_action"`
	MaxOrderQty     int     `json:"max_order_qty"`
	MaxOrderType    string  `json:"max_order_type"`
	MinStockQty     int     `json:"min_stock_qty"`
	MinStockType    string  `json:"min_stock_type"`
	SafetyStockQty  int     `json:"safety_stock_qty"`
	SafetyStockType string  `json:"safety_stock_type"`
	MinOrderQty     int     `json:"min_order_qty"`
	MinOrderType    string  `json:"min_order_type"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedBy       *int64  `json:"updated_by"`
	UpdatedByName   *string `json:"updated_by_name"`
	UpdatedAt       *string `json:"updated_at"`
}

type AutomaticReplenishmentProductImportRequest struct {
	FileURL string `json:"file_url" validate:"required"`
}

type AutomaticReplenishmentProductImportResponse struct {
	FileURL       string   `json:"file_url"`
	FileName      string   `json:"file_name"`
	TotalRow      int      `json:"total_row"`
	SuccessRow    int      `json:"success_row"`
	FailedRow     int      `json:"failed_row"`
	FailedReasons []string `json:"failed_reasons"`
	ProcessedAt   string   `json:"processed_at"`
}

type AutomaticReplenishmentProductDetailResponse struct {
	CustId          string  `json:"cust_id"`
	Id              int64   `json:"id"`
	ProId           int64   `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	DistributorId   int64   `json:"distributor_id"`
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	LimitAction     string  `json:"limit_action"`
	MaxOrderQty     int     `json:"max_order_qty"`
	MaxOrderType    string  `json:"max_order_type"`
	MinStockQty     int     `json:"min_stock_qty"`
	MinStockType    string  `json:"min_stock_type"`
	SafetyStockQty  int     `json:"safety_stock_qty"`
	SafetyStockType string  `json:"safety_stock_type"`
	MinOrderQty     int     `json:"min_order_qty"`
	MinOrderType    string  `json:"min_order_type"`
	IsActive        *bool   `json:"is_active"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedBy       *int64  `json:"updated_by"`
	UpdatedByName   *string `json:"updated_by_name"`
	UpdatedAt       *string `json:"updated_at"`
	DeletedBy       *int64  `json:"deleted_by"`
	DeletedByName   *string `json:"deleted_by_name"`
	DeletedAt       *string `json:"deleted_at"`
	IsDel           *bool   `json:"is_del"`
}
