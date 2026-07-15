package model

import (
	"time"
)

type Product struct {
	CustId            string     `json:"cust_id" db:"cust_id"`
	ProductId         int64      `json:"pro_id" db:"pro_id"`
	ProductCode       string     `json:"pro_code" db:"pro_code"`
	ProductName       string     `json:"pro_name" db:"pro_name"`
	BarCode           *string    `json:"bar_code" db:"bar_code"`
	PCatId            *int       `json:"pcat_id" db:"pcat_id"`
	PCatCode          *string    `json:"pcat_code" db:"pcat_code"`
	PCatName          *string    `json:"pcat_name" db:"pcat_name"`
	PlId              *int       `json:"pl_id" db:"pl_id"`
	PlCode            *string    `json:"pl_code" db:"pl_code"`
	PlName            *string    `json:"pl_name" db:"pl_name"`
	BrandId           *int       `json:"brand_id" db:"brand_id"`
	BrandCode         *string    `json:"brand_code" db:"brand_code"`
	BrandName         *string    `json:"brand_name" db:"brand_name"`
	Sbrand1           *int       `json:"sbrand1_id" db:"sbrand1_id"`
	Sbrand1Code       *string    `json:"sbrand1_code" db:"sbrand1_code"`
	Sbrand1Name       *string    `json:"sbrand1_name" db:"sbrand1_name"`
	Sbrand2           *int       `json:"sbrand2_id" db:"sbrand2_id"`
	Sbrand2Code       *string    `json:"sbrand2_code" db:"sbrand2_code"`
	Sbrand2Name       *string    `json:"sbrand2_name" db:"sbrand2_name"`
	FlavorId          *int       `json:"flavor_id" db:"flavor_id"`
	FlavorCode        *string    `json:"flavor_code" db:"flavor_code"`
	FlavorName        *string    `json:"flavor_name" db:"flavor_name"`
	PTypeId           *int       `json:"ptype_id" db:"ptype_id"`
	PTypeCode         *string    `json:"ptype_code" db:"ptype_code"`
	PTypeName         *string    `json:"ptype_name" db:"ptype_name"`
	PSizeId           *int       `json:"psize_id" db:"psize_id"`
	PSizeCode         *string    `json:"psize_code" db:"psize_code"`
	PsizeName         *string    `json:"psize_name" db:"psize_name"`
	SupId             *int       `json:"sup_id" db:"sup_id"`
	SupCode           *string    `json:"sup_code" db:"sup_code"`
	SupName           *string    `json:"sup_name" db:"sup_name"`
	PrincipalId       *int       `json:"principal_id" db:"principal_id"`
	PrincipalCode     *string    `json:"principal_code" db:"principal_code"`
	PrincipalName     *string    `json:"principal_name" db:"principal_name"`
	CProId            *int       `json:"c_pro_id" db:"c_pro_id"`
	CProCode          *string    `json:"c_pro_code" db:"c_pro_code"`
	CProName          *string    `json:"c_pro_name" db:"c_pro_name"`
	IsMainPro         bool       `json:"is_main_pro" db:"is_main_pro"`
	SortNo            int        `json:"sort_no" db:"sort_no"`
	ItemNo            int        `json:"item_no" db:"item_no"`
	UnitId1           string     `json:"unit_id1" db:"unit_id1"`
	UnitId2           string     `json:"unit_id2" db:"unit_id2"`
	UnitId3           string     `json:"unit_id3" db:"unit_id3"`
	UnitId4           *string    `json:"unit_id4" db:"unit_id4"`
	UnitId5           *string    `json:"unit_id5" db:"unit_id5"`
	UnitName1         *string    `json:"unit_name1" db:"unit_name1"`
	UnitName2         *string    `json:"unit_name2" db:"unit_name2"`
	UnitName3         *string    `json:"unit_name3" db:"unit_name3"`
	UnitIdCoreTax1    *string    `json:"unit_id_coretax1" db:"unit_id_coretax1"`
	UnitIdCoreTax2    *string    `json:"unit_id_coretax2" db:"unit_id_coretax2"`
	UnitIdCoreTax3    *string    `json:"unit_id_coretax3" db:"unit_id_coretax3"`
	UnitNameCoreTax1  *string    `json:"unit_name_coretax1" db:"unit_name_coretax1"`
	UnitNameCoreTax2  *string    `json:"unit_name_coretax2" db:"unit_name_coretax2"`
	UnitNameCoreTax3  *string    `json:"unit_name_coretax3" db:"unit_name_coretax3"`
	ConvUnit2         float32    `json:"conv_unit2" db:"conv_unit2"`
	ConvUnit3         float32    `json:"conv_unit3" db:"conv_unit3"`
	ConvUnit4         float32    `json:"conv_unit4" db:"conv_unit4"`
	ConvUnit5         float32    `json:"conv_unit5" db:"conv_unit5"`
	Weight            *float64   `json:"weight" db:"weight"`
	IsBatch           bool       `json:"is_batch" db:"is_batch"`
	IsExpDate         bool       `json:"is_exp_date" db:"is_exp_date"`
	Length            float64    `json:"length" db:"length"`
	Width             float64    `json:"width"  db:"width" `
	Height            float64    `json:"height" db:"height"`
	Volume            float64    `json:"volume" db:"volume"`
	PurchPrice1       float64    `json:"purch_price1" db:"purch_price1"`
	PurchPrice2       float64    `json:"purch_price2" db:"purch_price2"`
	PurchPrice3       float64    `json:"purch_price3" db:"purch_price3"`
	PurchPrice4       float64    `json:"purch_price4" db:"purch_price4"`
	PurchPrice5       float64    `json:"purch_price5" db:"purch_price5"`
	SellPrice1        float64    `json:"sell_price1" db:"sell_price1"`
	SellPrice2        float64    `json:"sell_price2" db:"sell_price2"`
	SellPrice3        float64    `json:"sell_price3" db:"sell_price3"`
	SellPrice4        float64    `json:"sell_price4" db:"sell_price4"`
	SellPrice5        float64    `json:"sell_price5" db:"sell_price5"`
	Length1           *float64   `json:"length1" db:"length1"`
	Length2           *float64   `json:"length2" db:"length2"`
	Length3           *float64   `json:"length3" db:"length3"`
	Length4           *float64   `json:"length4" db:"length4"`
	Length5           *float64   `json:"length5" db:"length5"`
	Width1            *float64   `json:"width1" db:"width1"`
	Width2            *float64   `json:"width2" db:"width2"`
	Width3            *float64   `json:"width3" db:"width3"`
	Width4            *float64   `json:"width4" db:"width4"`
	Width5            *float64   `json:"width5" db:"width5"`
	Height1           *float64   `json:"height1" db:"height1"`
	Height2           *float64   `json:"height2" db:"height2"`
	Height3           *float64   `json:"height3" db:"height3"`
	Height4           *float64   `json:"height4" db:"height4"`
	Height5           *float64   `json:"height5" db:"height5"`
	Weight1           *float64   `json:"weight1" db:"weight1"`
	Weight2           *float64   `json:"weight2" db:"weight2"`
	Weight3           *float64   `json:"weight3" db:"weight3"`
	Weight4           *float64   `json:"weight4" db:"weight4"`
	Weight5           *float64   `json:"weight5" db:"weight5"`
	Volume1           *float64   `json:"volume1" db:"volume1"`
	Volume2           *float64   `json:"volume2" db:"volume2"`
	Volume3           *float64   `json:"volume3" db:"volume3"`
	Volume4           *float64   `json:"volume4" db:"volume4"`
	Volume5           *float64   `json:"volume5" db:"volume5"`
	SafStockQty       float64    `json:"saf_stock_qty" db:"saf_stock_qty"`
	SafStockUnitId    *string    `json:"saf_stock_unit_id" db:"saf_stock_unit_id"`
	SafStockUnitName  *string    `json:"saf_stock_unit_name" db:"saf_stock_unit_name"`
	MinStockQty       float64    `json:"min_stock_qty" db:"min_stock_qty"`
	MinStockUnitId    *string    `json:"min_stock_unit_id" db:"min_stock_unit_id"`
	MinStockUnitName  *string    `json:"min_stock_unit_name" db:"min_stock_unit_name"`
	ParentProId       int        `json:"parent_pro_id" db:"parent_pro_id"`
	ExciseRate        float64    `json:"excise_rate" db:"excise_rate"`
	ExciseTax         float64    `json:"excise_tax" db:"excise_tax"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	IsDel             bool       `json:"is_del" db:"is_del"`
	CreatedBy         *int64     `json:"created_by,omitempty" db:"created_by,omitempty"`
	CreatedByName     *string    `json:"created_by_name" db:"created_by_name,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedBy         *int64     `json:"updated_by,omitempty" db:"updated_by,omitempty"`
	UpdatedByName     *string    `json:"updated_by_name" db:"updated_by_name,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	DeletedBy         *int64     `json:"deleted_by,omitempty" db:"deleted_by,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at,omitempty"`
	ImageUrl          *string    `json:"image_url" db:"image_url"`
	Vat               *float64   `json:"vat" db:"vat"`
	VatBg             *float64   `json:"vat_bg" db:"vat_bg"`
	VatLgPurch        *float64   `json:"vat_lg_purch" db:"vat_lg_purch"`
	VatLgSell         *float64   `json:"vat_lg_sell" db:"vat_lg_sell"`
	Cogs              *float64   `json:"cogs" db:"cogs"`
	ProStatus         *int       `json:"pro_status" db:"pro_status"`
	ParentProductCode *string    `json:"parent_pro_code" db:"parent_pro_code"`
	ParentProductName *string    `json:"parent_pro_name" db:"parent_pro_name"`
	ProCodeCoreTax    *string    `json:"pro_code_coretax" db:"pro_code_coretax"`
	ProNameCoreTax    *string    `json:"pro_name_coretax" db:"pro_name_coretax"`
	DistributorID     *int64     `json:"distributor_id" db:"distributor_id"`
	DistributorName   *string    `json:"distributor_name" db:"distributor_name,omitempty"`
	Level             int        `json:"level" db:"level"`
	Origin            string     `json:"origin" db:"origin"`
	AssignerUserID    *int64     `json:"assigner_user_id,omitempty" db:"assigner_user_id,omitempty"`
	IsProductMapping  bool       `json:"is_product_mapping" db:"is_product_mapping"`
}

type ProductTemp struct {
	HistoryId      string `db:"history_id"`
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
	StatusInsert   string `db:"status_insert"`
}

type ImportProductLine struct {
	PlId   int64  `db:"pl_id"`
	CustId string `db:"cust_id"`
	PlCode string `db:"pl_code"`
	PlName string `db:"pl_name"`
}

// Brand adalah representasi dari tabel m_brand.
type ImportBrand struct {
	BrandId   int64  `db:"brand_id"`
	CustId    string `db:"cust_id"`
	PlId      int64  `db:"pl_id"`
	PlCode    string `db:"pl_code"`
	PlName    string `db:"pl_name"`
	BrandCode string `db:"brand_code"`
	BrandName string `db:"brand_name"`
	// ... field lainnya ...
}

// SubBrand1 adalah representasi dari tabel m_sub_brand1.
type ImportSubBrand1 struct {
	Sbrand1Id   int64  `db:"sbrand1_id"`
	CustId      string `db:"cust_id"`
	BrandId     int64  `db:"brand_id"`
	Sbrand1Code string `db:"sbrand1_code"`
	Sbrand1Name string `db:"sbrand1_name"`
	// ... field lainnya ...
}

type ImportSubBrand2 struct {
	Sbrand2Id   int64  `db:"sbrand2_id"`
	CustId      string `db:"cust_id"`
	Sbrand1Code string `db:"sbrand2_code"`
	Sbrand1Name string `db:"sbrand2_name"`
	// ... field lainnya ...
}

type ImportInstruction struct {
	InstructionID   int64   `db:"instruction_id"`
	InstructionType string  `db:"instruction_type"`
	Kolom           string  `db:"kolom"`
	Mandatory       string  `db:"mandatory"`
	Step            string  `db:"step"`
	Color           string  `db:"color"`
	Keterangan      *string `db:"keterangan"` // nullable
}

type ImportFlavor struct {
	FlavorId   int64  `db:"flavor_id"`
	custId     string `db:"cust_id"`
	FlavorCode string `db:"flavor_code"`
	FlavorName string `db:"flavor_name"`
	// ... field lainnya ...
}

type ImportProductType struct {
	PtypeId   int64  `db:"ptype_id"`
	CustId    string `db:"cust_id"`
	PtypeCode string `db:"ptype_code"`
	PtypeName string `db:"ptype_name"`
	// ... field lainnya ...
}

type ImportProductSize struct {
	PsizeId   int64  `db:"psize_id"`
	CustId    string `db:"cust_id"`
	PsizeCode string `db:"psize_code"`
	PsizeName string `db:"psize_name"`
	// ... field lainnya ...
}

type ImportSupplier struct {
	SupId         int64  `db:"sup_id"`
	CustId        string `db:"cust_id"`
	SupCode       string `db:"sup_code"`
	SupName       string `db:"sup_name"`
	Phone         string `db:"phone"`
	Fax           string `db:"fax"`
	SupType       string `db:"sup_type"`
	ContactName   string `db:"contact_name"`
	TaxName       string `db:"tax_name"`
	Email         string `db:"email"`
	TaxIdentityNo string `db:"tax_identity_no"`
	TaxAddress    string `db:"tax_address"`
	// ... field lainnya ...
}

type ImportPrincipal struct {
	PrincipalId   int64  `db:"principal_id"`
	CustId        string `db:"cust_id"`
	PrincipalCode string `db:"principal_code"`
	PrincipalName string `db:"principal_name"`
}

type ImportConsPro struct {
	CProId   int64  `db:"c_pro_id"`
	CustId   string `db:"cust_id"`
	CProCode string `db:"c_pro_code"`
	CProName string `db:"c_pro_name"`
}

type ImportUnit struct {
	UnitId   int64  `db:"unit_id"`
	CustId   string `db:"cust_id"`
	UnitCode string `db:"unit_code"`
	UnitName string `db:"unit_name"`
	// ... field lainnya ...
}

type ImportUnitCoretax struct {
	UnitIdCoreTax   int64  `db:"unit_id_coretax"`
	CustId          string `db:"cust_id"`
	UnitNameCoreTax string `db:"unit_name_coretax"`
	IsActive        bool   `db:"is_active"`
	IsDel           bool   `db:"is_del"`
}

type ImportProductCoretax struct {
	CustId         string `db:"cust_id"`
	CatCoretax     string `db:"cat_coretax"`
	ProCodeCoretax string `db:"pro_code_coretax"`
	ProNameCoretax string `db:"pro_name_coretax"`
	IsActive       bool   `db:"is_active"`
	IsDel          bool   `db:"is_del"`
}

type ProductUpdate struct {
	ProductCode    *string    `json:"pro_code,omitempty" sql:"pro_code"`
	ProductName    *string    `json:"pro_name" sql:"pro_name"`
	BarCode        *string    `json:"bar_code" sql:"bar_code"`
	PCatId         *int       `json:"pcat_id" sql:"pcat_id"`
	Sbrand1        *int       `json:"sbrand1_id" sql:"sbrand1_id"`
	Sbrand2        *int       `json:"sbrand2_id" sql:"sbrand2_id"`
	FlavorId       *int       `json:"flavor_id" sql:"flavor_id"`
	PTypeId        *int       `json:"ptype_id" sql:"ptype_id"`
	PSizeId        *int       `json:"psize_id" sql:"psize_id"`
	SupId          *int       `json:"sup_id" sql:"sup_id"`
	PrincipalId    *int       `json:"principal_id" sql:"principal_id"`
	CProId         *int       `json:"c_pro_id" sql:"c_pro_id"`
	IsMainPro      *bool      `json:"is_main_pro" sql:"is_main_pro"`
	SortNo         *int       `json:"sort_no" sql:"sort_no"`
	ItemNo         *int       `json:"item_no" sql:"item_no"`
	UnitId1        *string    `json:"unit_id1" sql:"unit_id1"`
	UnitId2        *string    `json:"unit_id2" sql:"unit_id2"`
	UnitId3        *string    `json:"unit_id3" sql:"unit_id3"`
	UnitId4        *string    `json:"unit_id4" sql:"unit_id4"`
	UnitId5        *string    `json:"unit_id5" sql:"unit_id5"`
	ConvUnit2      *float32   `json:"conv_unit2" sql:"conv_unit2"`
	ConvUnit3      *float32   `json:"conv_unit3" sql:"conv_unit3"`
	ConvUnit4      *float32   `json:"conv_unit4" sql:"conv_unit4"`
	ConvUnit5      *float32   `json:"conv_unit5" sql:"conv_unit5"`
	Weight         *float64   `json:"weight" sql:"weight"`
	IsBatch        *bool      `json:"is_batch" sql:"is_batch"`
	IsExpDate      *bool      `json:"is_exp_date" sql:"is_exp_date"`
	Length         *float64   `json:"length" sql:"length"`
	Width          *float64   `json:"width"  sql:"width" `
	Height         *float64   `json:"height" sql:"height"`
	Volume         *float64   `json:"volume" sql:"volume"`
	PurchPrice1    *float64   `json:"purch_price1" sql:"purch_price1"`
	PurchPrice2    *float64   `json:"purch_price2" sql:"purch_price2"`
	PurchPrice3    *float64   `json:"purch_price3" sql:"purch_price3"`
	PurchPrice4    *float64   `json:"purch_price4" sql:"purch_price4"`
	PurchPrice5    *float64   `json:"purch_price5" sql:"purch_price5"`
	SellPrice1     *float64   `json:"sell_price1" sql:"sell_price1"`
	SellPrice2     *float64   `json:"sell_price2" sql:"sell_price2"`
	SellPrice3     *float64   `json:"sell_price3" sql:"sell_price3"`
	SellPrice4     *float64   `json:"sell_price4" sql:"sell_price4"`
	SellPrice5     *float64   `json:"sell_price5" sql:"sell_price5"`
	Length1        *float64   `json:"length1" sql:"length1"`
	Length2        *float64   `json:"length2" sql:"length2"`
	Length3        *float64   `json:"length3" sql:"length3"`
	Length4        *float64   `json:"length4" sql:"length4"`
	Length5        *float64   `json:"length5" sql:"length5"`
	Width1         *float64   `json:"width1" sql:"width1"`
	Width2         *float64   `json:"width2" sql:"width2"`
	Width3         *float64   `json:"width3" sql:"width3"`
	Width4         *float64   `json:"width4" sql:"width4"`
	Width5         *float64   `json:"width5" sql:"width5"`
	Height1        *float64   `json:"height1" sql:"height1"`
	Height2        *float64   `json:"height2" sql:"height2"`
	Height3        *float64   `json:"height3" sql:"height3"`
	Height4        *float64   `json:"height4" sql:"height4"`
	Height5        *float64   `json:"height5" sql:"height5"`
	Weight1        *float64   `json:"weight1" sql:"weight1"`
	Weight2        *float64   `json:"weight2" sql:"weight2"`
	Weight3        *float64   `json:"weight3" sql:"weight3"`
	Weight4        *float64   `json:"weight4" sql:"weight4"`
	Weight5        *float64   `json:"weight5" sql:"weight5"`
	Volume1        *float64   `json:"volume1" sql:"volume1"`
	Volume2        *float64   `json:"volume2" sql:"volume2"`
	Volume3        *float64   `json:"volume3" sql:"volume3"`
	Volume4        *float64   `json:"volume4" sql:"volume4"`
	Volume5        *float64   `json:"volume5" sql:"volume5"`
	SafStockQty    *float64   `json:"saf_stock_qty" sql:"saf_stock_qty"`
	SafStockUnitId *string    `json:"saf_stock_unit_id" sql:"saf_stock_unit_id"`
	MinStockQty    *float64   `json:"min_stock_qty" sql:"min_stock_qty"`
	MinStockUnitId *string    `json:"min_stock_unit_id" sql:"min_stock_unit_id"`
	ParentProId    *int       `json:"parent_pro_id" sql:"parent_pro_id"`
	ExciseRate     *float64   `json:"excise_rate" sql:"excise_rate"`
	ExciseTax      *float64   `json:"excise_tax" sql:"excise_tax"`
	IsActive       *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt      *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy      *int64     `json:"updated_by" sql:"updated_by"`
	ImageUrl       *string    `json:"image_url" sql:"image_url"`
	Vat            *float64   `json:"vat" sql:"vat"`
	VatBg          *float64   `json:"vat_bg" sql:"vat_bg"`
	VatLgPurch     *float64   `json:"vat_lg_purch" sql:"vat_lg_purch"`
	VatLgSell      *float64   `json:"vat_lg_sell" sql:"vat_lg_sell"`
	Cogs           *float64   `json:"cogs" sql:"cogs"`
	ProStatus      *int       `json:"pro_status" sql:"pro_status"`
	ProCodeCoretax *string    `json:"pro_code_coretax" db:"pro_code_coretax"`
}

type ProductDistPrice struct {
	ProductId   int64    `json:"pro_id" db:"pro_id"`
	ProductCode string   `json:"pro_code" db:"pro_code"`
	ProductName string   `json:"pro_name" db:"pro_name"`
	UnitId1     string   `json:"unit_id1" db:"unit_id1"`
	UnitId2     string   `json:"unit_id2" db:"unit_id2"`
	UnitId3     string   `json:"unit_id3" db:"unit_id3"`
	UnitId4     *string  `json:"unit_id4" db:"unit_id4"`
	UnitId5     *string  `json:"unit_id5" db:"unit_id5"`
	DistPriceId *int     `json:"dist_price_id" db:"dist_price_id"`
	ConvUnit2   *float32 `json:"conv_unit2" db:"conv_unit2"`
	ConvUnit3   *float32 `json:"conv_unit3" db:"conv_unit3"`
	ConvUnit4   *float32 `json:"conv_unit4" db:"conv_unit4"`
	ConvUnit5   *float32 `json:"conv_unit5" db:"conv_unit5"`
	PurchPrice1 *float64 `json:"purch_price1" db:"purch_price1"`
	PurchPrice2 *float64 `json:"purch_price2" db:"purch_price2"`
	PurchPrice3 *float64 `json:"purch_price3" db:"purch_price3"`
	PurchPrice4 *float64 `json:"purch_price4" db:"purch_price4"`
	PurchPrice5 *float64 `json:"purch_price5" db:"purch_price5"`
	SellPrice1  *float64 `json:"sell_price1" db:"sell_price1"`
	SellPrice2  *float64 `json:"sell_price2" db:"sell_price2"`
	SellPrice3  *float64 `json:"sell_price3" db:"sell_price3"`
	SellPrice4  *float64 `json:"sell_price4" db:"sell_price4"`
	SellPrice5  *float64 `json:"sell_price5" db:"sell_price5"`
	Vat         *float64 `json:"vat" db:"vat"`
	PlId        *int     `json:"pl_id" db:"pl_id"`
	PlCode      *string  `json:"pl_code" db:"pl_code"`
	PlName      *string  `json:"pl_name" db:"pl_name"`
	BrandId     *int     `json:"brand_id" db:"brand_id"`
	BrandCode   *string  `json:"brand_code" db:"brand_code"`
	BrandName   *string  `json:"brand_name" db:"brand_name"`
	Sbrand1     *int     `json:"sbrand1_id" db:"sbrand1_id"`
	Sbrand1Code *string  `json:"sbrand1_code" db:"sbrand1_code"`
	Sbrand1Name *string  `json:"sbrand1_name" db:"sbrand1_name"`
}
