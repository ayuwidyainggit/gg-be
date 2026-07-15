package model

import "time"

type ValidateOrderStok struct {
	CustID     string `gorm:"cust_id" json:"cust_id"`
	RoNo       string `gorm:"ro_no" json:"ro_no"`
	SalesmanId *int64 `gorm:"salesman_id" json:"salesman_id"`
}

func (ValidateOrderStok) TableName() string {
	return "sls.order"
}

type StockReport struct {
	ProID       int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode     string  `gorm:"column:pro_code" json:"pro_code"`
	ProName     string  `gorm:"column:pro_name" json:"pro_name"`
	UnitId1     string  `gorm:"column:unit_id1" json:"unit_id1"`
	UnitId2     string  `gorm:"column:unit_id2" json:"unit_id2"`
	UnitId3     string  `gorm:"column:unit_id3" json:"unit_id3"`
	ConvUnit2   int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3   int     `gorm:"column:conv_unit3" json:"conv_unit3"`
	PurchPrice1 float64 `gorm:"column:purch_price1" json:"purch_price1"`
	PurchPrice2 float64 `gorm:"column:purch_price2" json:"purch_price2"`
	PurchPrice3 float64 `gorm:"column:purch_price3" json:"purch_price3"`
	SellPrice1  float64 `gorm:"column:sell_price1" json:"sell_price1"`
	SellPrice2  float64 `gorm:"column:sell_price2" json:"sell_price2"`
	SellPrice3  float64 `gorm:"column:sell_price3" json:"sell_price3"`
	Qty         float64 `gorm:"column:qty" json:"qty"`
	IsActive    bool    `gorm:"column:is_active" json:"is_active"`
	Vat         float64 `gorm:"column:vat" json:"vat"`
	VatLgPurch  float64 `gorm:"column:vat_lg_purch" json:"vat_lg_purch"`
	VatLgSell   float64 `gorm:"column:vat_lg_sell" json:"vat_lg_sell"`
}

type WarehouseStockValidation struct {
	ProID     int64   `gorm:"column:pro_id" json:"pro_id"`
	Qty       float64 `gorm:"column:qty" json:"qty"`
	ConvUnit2 int     `gorm:"column:conv_unit2" json:"conv_unit2"`
	ConvUnit3 int     `gorm:"column:conv_unit3" json:"conv_unit3"`
}

type InvoicePaidAmount struct {
	PaidAmount float64 `gorm:"column:paid_amount" json:"paid_amount"`
}

func (InvoicePaidAmount) TableName() string {
	return "acf.deposit_detail"
}

type ArList struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	SalesmanId    int64      `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode  *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName  *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID      int64      `gorm:"outlet_id" json:"outlet_id"`
	OutletCode    *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceAmount float64    `gorm:"invoice_amount" json:"invoice_amount"`
	InvoiceNo     string     `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time `gorm:"invoice_date" json:"invoice_date"`
	DueDate       *time.Time `gorm:"due_date" json:"due_date"`
}

func (ArList) TableName() string {
	return "sls.order"
}

type Outlet struct {
	CustId             string     `json:"cust_id" db:"cust_id"`
	OutletId           int        `json:"outlet_id" db:"outlet_id"`
	OutletCode         *string    `json:"outlet_code,omitempty" db:"outlet_code"`
	OutletName         *string    `json:"outlet_name,omitempty" db:"outlet_name"`
	Barcode            *string    `json:"barcode" db:"barcode"`
	OutletStatus       *int       `json:"outlet_status" db:"outlet_status"`
	Address1           *string    `json:"address1" db:"address1"`
	Address2           *string    `json:"address2" db:"address2"`
	City               *string    `json:"city" db:"city"`
	ZipCode            *string    `json:"zip_code" db:"zip_code"`
	PhoneNo            *string    `json:"phone_no" db:"phone_no"`
	WaNo               *string    `json:"wa_no" db:"wa_no"`
	FaxNo              *string    `json:"fax_no" db:"fax_no"`
	Email              *string    `json:"email" db:"email"`
	DiscGrpId          *int       `json:"disc_grp_id" db:"disc_grp_id"`
	OtLocId            *int       `json:"ot_loc_id" db:"ot_loc_id"`
	OtGrpId            *int       `json:"ot_grp_id" db:"ot_grp_id"`
	PriceGrpId         *int       `json:"price_grp_id" db:"price_grp_id"`
	DistrictId         *int       `json:"district_id" db:"district_id"`
	BeatId             *int       `json:"beat_id" db:"beat_id"`
	SbeatId            *int       `json:"sbeat_id" db:"sbeat_id"`
	OtClassId          *int       `json:"ot_class_id" db:"ot_class_id"`
	IndustryId         *int       `json:"industry_id" db:"industry_id"`
	MarketId           *int       `json:"market_id" db:"market_id"`
	Top                *int       `json:"top" db:"top"`
	PaymentType        *int       `json:"payment_type" db:"payment_type"`
	IsContraBon        *bool      `json:"is_contra_bon" db:"is_contra_bon"`
	PluGrpId           *int       `json:"plu_grp_id" db:"plu_grp_id"`
	ConvGrpId          *int       `json:"conv_grp_id" db:"conv_grp_id"`
	DiscInvId          *int       `json:"disc_inv_id" db:"disc_inv_id"`
	AgentFrom          *string    `json:"agent_from" db:"agent_from"`
	CreditLimitType    *int       `json:"credit_limit_type" db:"credit_limit_type"`
	CreditLimit        *float64   `json:"credit_limit" db:"credit_limit"`
	SalesInvLimitType  *int       `json:"sales_inv_limit_type" db:"sales_inv_limit_type"`
	SalesInvLimit      *int       `json:"sales_inv_limit" db:"sales_inv_limit"`
	AvgSalesWeek       *float64   `json:"avg_sales_week" db:"avg_sales_week"`
	AvgSalesMonth      *float64   `json:"avg_sales_month" db:"avg_sales_month"`
	FirstTransDate     *string    `json:"first_trans_date" db:"first_trans_date"`
	LastTransDate      *string    `json:"last_trans_date" db:"last_trans_date"`
	FirstWeekNo        *int       `json:"first_week_no" db:"first_week_no"`
	OtStartDate        *string    `json:"ot_start_date" db:"ot_start_date"`
	OtRegDate          *string    `json:"ot_reg_date" db:"ot_reg_date"`
	BuldingOwn         *int       `json:"building_own" db:"building_own"`
	Dob                *string    `json:"dob" db:"dob"`
	ArStatus           *int       `json:"ar_status" db:"ar_status"`
	ArTotal            *float64   `json:"ar_total" db:"ar_total"`
	CloseDate          *time.Time `json:"closed_date" db:"closed_date"`
	IsEmbBail          *bool      `json:"is_emb_bail" db:"is_emb_bail"`
	TaxName            *string    `json:"tax_name" db:"tax_name"`
	TaxAddr1           *string    `json:"tax_addr1" db:"tax_addr1"`
	TaxAddr2           *string    `json:"tax_addr2" db:"tax_addr2"`
	TaxCity            *string    `json:"tax_city" db:"tax_city"`
	TaxNo              *string    `json:"tax_no" db:"tax_no"`
	TaxInvoiceForm     *int       `json:"tax_invoice_form" db:"tax_invoice_form"`
	OwnerName          *string    `json:"owner_name" db:"owner_name"`
	OwnerAdd1          *string    `json:"owner_addr1" db:"owner_addr1"`
	OwnerAddr2         *string    `json:"owner_addr2" db:"owner_addr2"`
	OwnerCity          *string    `json:"owner_city" db:"owner_city"`
	OwnerPhoneNo       *string    `json:"owner_phone_no" db:"owner_phone_no"`
	OwnerIdNo          *string    `json:"owner_id_no" db:"owner_id_no"`
	DelvAdd1           *string    `json:"delv_addr1" db:"delv_addr1"`
	DelvAddr2          *string    `json:"delv_addr2" db:"delv_addr2"`
	DelvCity           *string    `json:"delv_city" db:"delv_city"`
	InvAddr1           *string    `json:"inv_addr1" db:"inv_addr1"`
	InvAddr2           *string    `json:"inv_addr2" db:"inv_addr2"`
	InvCity            *string    `json:"inv_city" db:"inv_city"`
	Latitude           *string    `json:"latitude" db:"latitude"`
	Longitude          *string    `json:"longitude" db:"longitude"`
	ImageUrl           *string    `json:"image_url" db:"image_url"`
	IsActive           bool       `db:"is_active" json:"is_active"`
	IsDel              bool       `db:"is_del" json:"is_del"`
	CreatedBy          *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt          *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy          *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName      *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt          *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy          *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt          *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	OtTypeId           *int       `json:"ot_type_id" db:"ot_type_id"`
	OtTypeName         *string    `json:"ot_type_name" db:"ot_type_name"`
	IsObs              *bool      `json:"is_obs" db:"is_obs"`
	Obs                *int       `json:"obs" db:"obs"`
	ObsType            *int       `json:"obs_type" db:"obs_type"`
	OutletWardId       *string    `db:"outlet_ward_id" json:"outlet_ward_id"`
	IsWaNo             *bool      `json:"is_wa_no" db:"is_wa_no"`
	DelvWardId         *string    `db:"delv_ward_id" json:"delv_ward_id"`
	DelvZipCode        *string    `json:"delv_zip_code" db:"delv_zip_code"`
	DelvIsSameAddress  *bool      `json:"delv_is_same_addr" db:"delv_is_same_addr"`
	InvWardId          *string    `db:"inv_ward_id" json:"inv_ward_id"`
	InvZipCode         *string    `json:"inv_zip_code" db:"inv_zip_code"`
	InvIsSameAddress   *bool      `json:"inv_is_same_addr" db:"inv_is_same_addr"`
	VerificationStatus *int       `json:"verification_status" db:"verification_status"`
}

func (Outlet) TableName() string {
	return "mst.m_outlet"
}
