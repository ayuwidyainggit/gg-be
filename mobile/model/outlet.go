package model

import "time"

type Outlet struct {
	CustId            string    `gorm:"cust_id" json:"cust_id"`
	OutletId          int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode        string    `gorm:"outlet_code" json:"outlet_code"`
	Barcode           string    `gorm:"barcode" json:"barcode"`
	OutletName        string    `gorm:"outlet_name" json:"outlet_name"`
	OutletStatus      int64     `gorm:"outlet_status" json:"outlet_status"`
	Address1          string    `gorm:"address1" json:"address1"`
	Address2          string    `gorm:"address2" json:"address2"`
	City              string    `gorm:"city" json:"city"`
	ZipCode           string    `gorm:"zip_code" json:"zip_code"`
	PhoneNo           string    `gorm:"phone_no" json:"phone_no"`
	WaNo              string    `gorm:"wa_no" json:"wa_no"`
	FaxNo             string    `gorm:"fax_no" json:"fax_no"`
	Email             string    `gorm:"email" json:"email"`
	DiscGrpId         int64     `gorm:"disc_grp_id" json:"disc_grp_id"`
	OtLocId           int64     `gorm:"ot_loc_id" json:"ot_loc_id"`
	OtGrpId           int64     `gorm:"ot_grp_id" json:"ot_grp_id"`
	PriceGrpId        int64     `gorm:"price_grp_id" json:"price_grp_id"`
	DistrictId        int64     `gorm:"district_id" json:"district_id"`
	BeatId            int64     `gorm:"beat_id" json:"beat_id"`
	SbeatId           int64     `gorm:"sbeat_id" json:"sbeat_id"`
	OtClassId         int64     `gorm:"ot_class_id" json:"ot_class_id"`
	IndustryId        int64     `gorm:"industry_id" json:"industry_id"`
	MarketId          int64     `gorm:"market_id" json:"market_id"`
	Top               int64     `gorm:"top" json:"top"`
	PaymentType       int64     `gorm:"payment_type" json:"payment_type"`
	IsContraBon       bool      `gorm:"is_contra_bon" json:"is_contra_bon"`
	PluGrpId          int64     `gorm:"plu_grp_id" json:"plu_grp_id"`
	ConvGrpId         int64     `gorm:"conv_grp_id" json:"conv_grp_id"`
	DiscInvId         int64     `gorm:"disc_inv_id" json:"disc_inv_id"`
	AgentFrom         string    `gorm:"agent_from" json:"agent_from"`
	CreditLimitType   int64     `gorm:"credit_limit_type" json:"credit_limit_type"`
	CreditLimit       float64   `gorm:"credit_limit" json:"credit_limit"`
	SalesInvLimitType int64     `gorm:"sales_inv_limit_type" json:"sales_inv_limit_type"`
	SalesInvLimit     int64     `gorm:"sales_inv_limit" json:"sales_inv_limit"`
	AvgSalesWeek      float64   `gorm:"avg_sales_week" json:"avg_sales_week"`
	AvgSalesMonth     float64   `gorm:"avg_sales_month" json:"avg_sales_month"`
	FirstTransDate    time.Time `gorm:"first_trans_date" json:"first_trans_date"`
	LastTransDate     time.Time `gorm:"last_trans_date" json:"last_trans_date"`
	FirstWeekNo       int64     `gorm:"first_week_no" json:"first_week_no"`
	OtStartDate       time.Time `gorm:"ot_start_date" json:"ot_start_date"`
	OtRegDate         time.Time `gorm:"ot_reg_date" json:"ot_reg_date"`
	BuildingOwn       int64     `gorm:"building_own" json:"building_own"`
	Dob               time.Time `gorm:"dob" json:"dob"`
	ArStatus          int64     `gorm:"ar_status" json:"ar_status"`
	ArTotal           float64   `gorm:"ar_total" json:"ar_total"`
	ClosedDate        time.Time `gorm:"closed_date" json:"closed_date"`
	IsEmbBail         bool      `gorm:"is_emb_bail" json:"is_emb_bail"`
	TaxName           string    `gorm:"tax_name" json:"tax_name"`
	TaxAddr1          string    `gorm:"tax_addr1" json:"tax_addr_1"`
	TaxAddr2          string    `gorm:"tax_addr2" json:"tax_addr_2"`
	TaxCity           string    `gorm:"tax_city" json:"tax_city"`
	TaxNo             string    `gorm:"tax_no" json:"tax_no"`
	TaxInvoiceForm    string    `gorm:"tax_invoice_form" json:"tax_invoice_form"`
	OwnerName         string    `gorm:"owner_name" json:"owner_name"`
	OwnerAddr1        string    `gorm:"owner_addr1" json:"owner_addr_1"`
	OwnerAddr2        string    `gorm:"owner_addr2" json:"owner_addr_2"`
	OwnerCity         string    `gorm:"owner_city" json:"owner_city"`
	OwnerPhoneNo      string    `gorm:"owner_phone_no" json:"owner_phone_no"`
	OwnerIdNo         string    `gorm:"owner_id_no" json:"owner_id_no"`
	DelvAddr1         string    `gorm:"delv_addr1" json:"delv_addr_1"`
	DelvAddr2         string    `gorm:"delv_addr2" json:"delv_addr_2"`
	DelvCity          string    `gorm:"delv_city" json:"delv_city"`
	InvAddr1          string    `gorm:"inv_addr1" json:"inv_addr_1"`
	InvAddr2          string    `gorm:"inv_addr2" json:"inv_addr_2"`
	InvCity           string    `gorm:"inv_city" json:"inv_city"`
	IsActive          bool      `gorm:"is_active" json:"is_active"`
	CreatedBy         int64     `gorm:"created_by" json:"created_by"`
	CreatedAt         time.Time `gorm:"created_at" json:"created_at"`
	UpdatedBy         int64     `gorm:"updated_by" json:"updated_by"`
	UpdatedAt         time.Time `gorm:"updated_at" json:"updated_at"`
	IsDel             bool      `gorm:"is_del" json:"is_del"`
	DeletedBy         int64     `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt         time.Time `gorm:"deleted_at" json:"deleted_at"`
	Latitude          string    `gorm:"latitude" json:"latitude"`
	Longitude         string    `gorm:"longitude" json:"longitude"`
	ImageURL          string    `gorm:"image_url" json:"image_url"`
	OtTypeId          int64     `gorm:"ot_type_id" json:"ot_type_id"`
	IsObs             bool      `gorm:"is_obs" json:"is_obs"`
	Obs               int       `gorm:"obs" json:"obs"`
	NoOrderId         int       `gorm:"no_order_id" json:"no_order_id"`
}

func (Outlet) TableName() string {
	return "mst.m_outlet"
}
