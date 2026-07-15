package model

import "time"

type Outlet struct {
	CustId                  string     `json:"cust_id" db:"cust_id"`
	OutletId                int64      `json:"outlet_id" db:"outlet_id"`
	OutletCode              *string    `json:"outlet_code,omitempty" db:"outlet_code"`
	OutletName              *string    `json:"outlet_name,omitempty" db:"outlet_name"`
	Barcode                 *string    `json:"barcode" db:"barcode"`
	OutletStatus            *int       `json:"outlet_status" db:"outlet_status"`
	Address1                *string    `json:"address1" db:"address1"`
	Address2                *string    `json:"address2" db:"address2"`
	City                    *string    `json:"city" db:"city"`
	ZipCode                 *string    `json:"zip_code" db:"zip_code"`
	PhoneNo                 *string    `json:"phone_no" db:"phone_no"`
	WaNo                    *string    `json:"wa_no" db:"wa_no"`
	FaxNo                   *string    `json:"fax_no" db:"fax_no"`
	Email                   *string    `json:"email" db:"email"`
	DiscGrpId               *int       `json:"disc_grp_id" db:"disc_grp_id"`
	OtLocId                 *int       `json:"ot_loc_id" db:"ot_loc_id"`
	OtGrpId                 *int       `json:"ot_grp_id" db:"ot_grp_id"`
	OtGrpName               *string    `json:"ot_grp_name" db:"ot_grp_name"`
	PriceGrpId              *int       `json:"price_grp_id" db:"price_grp_id"`
	DistrictId              *int       `json:"district_id" db:"district_id"`
	BeatId                  *int       `json:"beat_id" db:"beat_id"`
	SbeatId                 *int       `json:"sbeat_id" db:"sbeat_id"`
	OtClassId               *int       `json:"ot_class_id" db:"ot_class_id"`
	IndustryId              *int       `json:"industry_id" db:"industry_id"`
	MarketId                *int       `json:"market_id" db:"market_id"`
	Top                     *int       `json:"top" db:"top"`
	PaymentType             *int       `json:"payment_type" db:"payment_type"`
	IsContraBon             *bool      `json:"is_contra_bon" db:"is_contra_bon"`
	PluGrpId                *int       `json:"plu_grp_id" db:"plu_grp_id"`
	ConvGrpId               *int       `json:"conv_grp_id" db:"conv_grp_id"`
	DiscInvId               *int       `json:"disc_inv_id" db:"disc_inv_id"`
	AgentFrom               *string    `json:"agent_from" db:"agent_from"`
	CreditLimitType         *int       `json:"credit_limit_type" db:"credit_limit_type"`
	CreditLimitTypeName     *string    `json:"credit_limit_type_name" db:"credit_limit_type_name"`
	CreditLimit             *float64   `json:"credit_limit" db:"credit_limit"`
	CreditLimitAction       *int       `json:"credit_limit_action" db:"credit_limit_action"`
	CreditLimitActionName   *string    `json:"credit_limit_action_name" db:"credit_limit_action_name"`
	SalesInvLimitType       *int       `json:"sales_inv_limit_type" db:"sales_inv_limit_type"`
	SalesInvLimitTypeName   *string    `json:"sales_inv_limit_type_name" db:"sales_inv_limit_type_name"`
	SalesInvLimit           *int       `json:"sales_inv_limit" db:"sales_inv_limit"`
	SalesInvLimitAction     *int       `json:"sales_inv_limit_action" db:"sales_inv_limit_action"`
	SalesInvLimitActionName *string    `json:"sales_inv_limit_action_name" db:"sales_inv_limit_action_name"`
	AvgSalesWeek            *float64   `json:"avg_sales_week" db:"avg_sales_week"`
	AvgSalesMonth           *float64   `json:"avg_sales_month" db:"avg_sales_month"`
	FirstTransDate          *string    `json:"first_trans_date" db:"first_trans_date"`
	LastTransDate           *string    `json:"last_trans_date" db:"last_trans_date"`
	FirstWeekNo             *int       `json:"first_week_no" db:"first_week_no"`
	OtStartDate             *string    `json:"ot_start_date" db:"ot_start_date"`
	OtRegDate               *string    `json:"ot_reg_date" db:"ot_reg_date"`
	BuldingOwn              *int       `json:"building_own" db:"building_own"`
	Dob                     *string    `json:"dob" db:"dob"`
	ArStatus                *int       `json:"ar_status" db:"ar_status"`
	ArTotal                 *float64   `json:"ar_total" db:"ar_total"`
	CloseDate               *time.Time `json:"closed_date" db:"closed_date"`
	IsEmbBail               *bool      `json:"is_emb_bail" db:"is_emb_bail"`
	IsPkpOutlet             *bool      `json:"is_pkp_outlet" db:"is_pkp_outlet"`
	TaxName                 *string    `json:"tax_name" db:"tax_name"`
	TaxAddr1                *string    `json:"tax_addr1" db:"tax_addr1"`
	TaxAddr2                *string    `json:"tax_addr2" db:"tax_addr2"`
	TaxCity                 *string    `json:"tax_city" db:"tax_city"`
	TaxNo                   *string    `json:"tax_no" db:"tax_no"`
	TaxInvoiceForm          *int       `json:"tax_invoice_form" db:"tax_invoice_form"`
	OwnerName               *string    `json:"owner_name" db:"owner_name"`
	OwnerAdd1               *string    `json:"owner_addr1" db:"owner_addr1"`
	OwnerAddr2              *string    `json:"owner_addr2" db:"owner_addr2"`
	OwnerCity               *string    `json:"owner_city" db:"owner_city"`
	OwnerPhoneNo            *string    `json:"owner_phone_no" db:"owner_phone_no"`
	OwnerIdNo               *string    `json:"owner_id_no" db:"owner_id_no"`
	DelvAdd1                *string    `json:"delv_addr1" db:"delv_addr1"`
	DelvCity                *string    `json:"delv_city" db:"delv_city"`
	DelvLatitude            *string    `json:"delv_latitude" db:"delv_latitude"`
	DelvLongitude           *string    `json:"delv_longitude" db:"delv_longitude"`
	DelvAddr2               *string    `json:"delv_addr2" db:"delv_addr2"`
	DelvCity2               *string    `json:"delv_city2" db:"delv_city2"`
	DelvLatitude2           *string    `json:"delv_latitude2" db:"delv_latitude2"`
	DelvLongitude2          *string    `json:"delv_longitude2" db:"delv_longitude2"`
	InvAddr1                *string    `json:"inv_addr1" db:"inv_addr1"`
	InvAddr2                *string    `json:"inv_addr2" db:"inv_addr2"`
	InvCity                 *string    `json:"inv_city" db:"inv_city"`
	Latitude                *string    `json:"latitude" db:"latitude"`
	Longitude               *string    `json:"longitude" db:"longitude"`
	ImageUrl                *string    `json:"image_url" db:"image_url"`
	FileUrl                 *string    `json:"file_url" db:"file_url"`
	IsActive                bool       `db:"is_active" json:"is_active"`
	IsDel                   bool       `db:"is_del" json:"is_del"`
	CreatedBy               *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt               *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy               *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName           *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt               *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy               *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt               *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	OtTypeId                *int       `json:"ot_type_id" db:"ot_type_id"`
	OtTypeName              *string    `json:"ot_type_name" db:"ot_type_name"`
	IsObs                   *bool      `json:"is_obs" db:"is_obs"`
	ObsType                 *int       `json:"obs_type" db:"obs_type"`
	Obs                     *int       `json:"obs" db:"obs"`
	ObsLimitAction          *int       `json:"obs_limit_action" db:"obs_limit_action"`
	OutletProvinceId        *string    `db:"outlet_province_id" json:"outlet_province_id"`
	OutletRegencyId         *string    `db:"outlet_regency_id" json:"outlet_regency_id"`
	OutletSubDistrictId     *string    `db:"outlet_sub_district_id" json:"outlet_sub_district_id"`
	OutletWardId            *string    `db:"outlet_ward_id" json:"outlet_ward_id"`
	IsWaNo                  *bool      `json:"is_wa_no" db:"is_wa_no"`
	DelvWardId              *string    `db:"delv_ward_id" json:"delv_ward_id"`
	DelvZipCode             *string    `json:"delv_zip_code" db:"delv_zip_code"`
	DelvIsSameAddress       *bool      `json:"delv_is_same_addr" db:"delv_is_same_addr"`
	DelvWardId2             *string    `db:"delv_ward_id2" json:"delv_ward_id2"`
	DelvZipCode2            *string    `json:"delv_zip_code2" db:"delv_zip_code2"`
	InvWardId               *string    `db:"inv_ward_id" json:"inv_ward_id"`
	InvZipCode              *string    `json:"inv_zip_code" db:"inv_zip_code"`
	InvIsSameAddress        *bool      `json:"inv_is_same_addr" db:"inv_is_same_addr"`
	VerificationStatus      *int       `json:"verification_status" db:"verification_status"`
	OutletEstablishmentDate *string    `db:"outlet_establishment_date" json:"outlet_establishment_date"`
	OutletPrincipalCode     *string    `json:"outlet_principal_code,omitempty" db:"outlet_principal_code"`
}

type MOutlet struct {
	CustId     string     `db:"cust_id"`
	OutletId   int        `db:"outlet_id"`
	OutletCode string     `db:"outlet_code"`
	OutletName string     `db:"outlet_name"`
	IsActive   bool       `db:"is_active"`
	IsDel      bool       `db:"is_del"`
	CreatedBy  *int64     `db:"created_by,omitempty"`
	CreatedAt  *time.Time `db:"created_at,omitempty"`
	UpdatedBy  *int64     `db:"updated_by,omitempty"`
	UpdatedAt  *time.Time `db:"updated_at,omitempty"`
	DeletedBy  *int64     `db:"deleted_by,omitempty"`
	DeletedAt  *time.Time `db:"deleted_at,omitempty"`
}

type MOutletUpdate struct {
	MOutletCode *string    `json:"outlet_code,omitempty" sql:"outlet_code"`
	MOutletName *string    `json:"outlet_name,omitempty" sql:"outlet_name"`
	IsActive    *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy   *int64     `json:"updated_by" sql:"updated_by"`
}

type OutletTemp struct {
	HistoryId    string `db:"history_id"`
	CustId       string `db:"cust_id"`
	OutletCode   string `db:"outlet_code"`
	OutletName   string `db:"outlet_name"`
	Address1     string `db:"address1"`
	Address2     string `db:"address2"`
	City         string `db:"city"`
	PhoneNo      string `db:"phone_no"`
	WaNo         string `db:"wa_no"`
	Email        string `db:"email"`
	DiscGrpCode  string `db:"disc_grp_code"`
	OtClassCode  string `db:"ot_class_code"`
	OtGrpCode    string `db:"ot_grp_code"`
	OtTypeCode   string `db:"ot_type_code"`
	DistrictCode string `db:"district_code"`
	MarketCode   string `db:"market_code"`
	IndustryCode string `db:"industry_code"`
	StatusInsert string `db:"status_insert"`
}

type OutletUpdate struct {
	OutletCode              *string    `json:"outlet_code,omitempty" sql:"outlet_code"`
	OutletName              *string    `json:"outlet_name,omitempty" sql:"outlet_name"`
	Barcode                 *string    `json:"barcode" sql:"barcode"`
	OutletStatus            *int       `json:"outlet_status" sql:"outlet_status"`
	Address1                *string    `json:"address1" sql:"address1"`
	Address2                *string    `json:"address2" sql:"address2"`
	OutletProvinceId        *string    `json:"outlet_province_id" sql:"outlet_province_id"`
	OutletRegencyId         *string    `json:"outlet_regency_id" sql:"outlet_regency_id"`
	OutletSubDistrictId     *string    `json:"outlet_sub_district_id" sql:"outlet_sub_district_id"`
	City                    *string    `json:"city" sql:"city"`
	ZipCode                 *string    `json:"zip_code" sql:"zip_code"`
	PhoneNo                 *string    `json:"phone_no" sql:"phone_no"`
	WaNo                    *string    `json:"wa_no" sql:"wa_no"`
	FaxNo                   *string    `json:"fax_no" sql:"fax_no"`
	Email                   *string    `json:"email" sql:"email"`
	IsPkpOutlet             *bool      `json:"is_pkp_outlet" sql:"is_pkp_outlet"`
	IdentityType            *string    `json:"identity_type" sql:"identity_type"`
	DiscGrpId               *int       `json:"disc_grp_id" sql:"disc_grp_id"`
	OtLocId                 *int       `json:"ot_loc_id" sql:"ot_loc_id"`
	OtGrpId                 *int       `json:"ot_grp_id" sql:"ot_grp_id"`
	PriceGrpId              *int       `json:"price_grp_id" sql:"price_grp_id"`
	DistrictId              *int       `json:"district_id" sql:"district_id"`
	BeatId                  *int       `json:"beat_id" sql:"beat_id"`
	SbeatId                 *int       `json:"sbeat_id" sql:"sbeat_id"`
	OtClassId               *int       `json:"ot_class_id" sql:"ot_class_id"`
	IndustryId              *int       `json:"industry_id" sql:"industry_id"`
	MarketId                *int       `json:"market_id" sql:"market_id"`
	Top                     *int       `json:"top" sql:"top"`
	PaymentType             *int       `json:"payment_type" sql:"payment_type"`
	IsContraBon             *bool      `json:"is_contra_bon" sql:"is_contra_bon"`
	PluGrpId                *int       `json:"plu_grp_id" sql:"plu_grp_id"`
	ConvGrpId               *int       `json:"conv_grp_id" sql:"conv_grp_id"`
	DiscInvId               *int       `json:"disc_inv_id" sql:"disc_inv_id"`
	AgentFrom               *string    `json:"agent_from" sql:"agent_from"`
	CreditLimitType         *int       `json:"credit_limit_type" sql:"credit_limit_type"`
	CreditLimitTypeName     *string    `json:"credit_limit_type_name" sql:"credit_limit_type_name"`
	CreditLimit             *float64   `json:"credit_limit" sql:"credit_limit"`
	CreditLimitAction       *int       `json:"credit_limit_action" sql:"credit_limit_action"`
	CreditLimitActionName   *string    `json:"credit_limit_action_name" sql:"credit_limit_action_name"`
	SalesInvLimitType       *int       `json:"sales_inv_limit_type" sql:"sales_inv_limit_type"`
	SalesInvLimitTypeName   *string    `json:"sales_inv_limit_type_name" sql:"sales_inv_limit_type_name"`
	SalesInvLimit           *int       `json:"sales_inv_limit" sql:"sales_inv_limit"`
	SalesInvLimitAction     *int       `json:"sales_inv_limit_action" sql:"sales_inv_limit_action"`
	SalesInvLimitActionName *string    `json:"sales_inv_limit_action_name" sql:"sales_inv_limit_action_name"`
	AvgSalesWeek            *float64   `json:"avg_sales_week" sql:"avg_sales_week"`
	AvgSalesMonth           *float64   `json:"avg_sales_month" sql:"avg_sales_month"`
	FirstTransDate          *string    `json:"first_trans_date" sql:"first_trans_date"`
	LastTransDate           *string    `json:"last_trans_date" sql:"last_trans_date"`
	FirstWeekNo             *int       `json:"first_week_no" sql:"first_week_no"`
	OtStartDate             *string    `json:"ot_start_date" sql:"ot_start_date"`
	OtRegDate               *string    `json:"ot_reg_date" sql:"ot_reg_date"`
	BuldingOwn              *int       `json:"building_own" sql:"building_own"`
	Dob                     *string    `json:"dob" sql:"dob"`
	ArStatus                *int       `json:"ar_status" sql:"ar_status"`
	ArTotal                 *float64   `json:"ar_total" sql:"ar_total"`
	CloseDate               *string    `json:"closed_date" sql:"closed_date"`
	IsEmbBail               *bool      `json:"is_emb_bail" sql:"is_emb_bail"`
	TaxName                 *string    `json:"tax_name" sql:"tax_name"`
	TaxAddr1                *string    `json:"tax_addr1" sql:"tax_addr1"`
	TaxAddr2                *string    `json:"tax_addr2" sql:"tax_addr2"`
	TaxCity                 *string    `json:"tax_city" sql:"tax_city"`
	TaxNo                   *string    `json:"tax_no" sql:"tax_no"`
	Nitku                   *string    `json:"nitku" sql:"nitku"`
	TaxType                 *string    `json:"tax_type" sql:"tax_type"`
	TaxInvoiceForm          *int       `json:"tax_invoice_form" sql:"tax_invoice_form"`
	OwnerName               *string    `json:"owner_name" sql:"owner_name"`
	OwnerAdd1               *string    `json:"owner_addr1" sql:"owner_addr1"`
	OwnerAddr2              *string    `json:"owner_addr2" sql:"owner_addr2"`
	OwnerCity               *string    `json:"owner_city" sql:"owner_city"`
	OwnerPhoneNo            *string    `json:"owner_phone_no" sql:"owner_phone_no"`
	OwnerIdNo               *string    `json:"owner_id_no" sql:"owner_id_no"`
	DelvAdd1                *string    `json:"delv_addr1" sql:"delv_addr1"`
	DelvCity                *string    `json:"delv_city" sql:"delv_city"`
	DelvLatitude            *string    `json:"delv_latitude" sql:"delv_latitude"`
	DelvLongitude           *string    `json:"delv_longitude" sql:"delv_longitude"`
	DelvAddr2               *string    `json:"delv_addr2" sql:"delv_addr2"`
	DelvCity2               *string    `json:"delv_city2" sql:"delv_city2"`
	DelvLatitude2           *string    `json:"delv_latitude2" sql:"delv_latitude2"`
	DelvLongitude2          *string    `json:"delv_longitude2" sql:"delv_longitude2"`
	InvAddr1                *string    `json:"inv_addr1" sql:"inv_addr1"`
	InvAddr2                *string    `json:"inv_addr2" sql:"inv_addr2"`
	InvCity                 *string    `json:"inv_city" sql:"inv_city"`
	Latitude                *string    `json:"latitude" sql:"latitude"`
	Longitude               *string    `json:"longitude" sql:"longitude"`
	ImageUrl                *string    `json:"image_url" db:"image_url"`
	IsActive                *bool      `json:"is_active" sql:"is_active"`
	UpdatedAt               *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy               *int64     `json:"updated_by" sql:"updated_by"`
	OtTypeId                *int       `json:"ot_type_id" db:"ot_type_id"`
	IsObs                   *bool      `json:"is_obs" db:"is_obs"`
	ObsType                 *int       `json:"obs_type" db:"obs_type"`
	Obs                     *int       `json:"obs" db:"obs"`
	ObsLimitAction          *int       `json:"obs_limit_action" sql:"obs_limit_action"`
	OutletWardId            *string    `db:"outlet_ward_id" json:"outlet_ward_id"`
	IsWaNo                  *bool      `json:"is_wa_no" db:"is_wa_no"`
	DelvWardId              *string    `db:"delv_ward_id" json:"delv_ward_id"`
	DelvZipCode             *string    `json:"delv_zip_code" db:"delv_zip_code"`
	DelvIsSameAddress       *bool      `json:"delv_is_same_addr" db:"delv_is_same_addr" sql:"delv_is_same_addr"`
	DelvProvinceId1         *string    `json:"delv_province_id1" sql:"-"`  // derived from delv_ward_id, not stored in m_outlet
	DelvRegencyId1          *string    `json:"delv_regency_id1" sql:"-"`   // derived from delv_ward_id
	DelvSubDistrictId1      *string    `json:"delv_sub_district_id1" sql:"-"` // derived from delv_ward_id
	DelvWardId2             *string    `db:"delv_ward_id2" json:"delv_ward_id2"`
	DelvZipCode2            *string    `json:"delv_zip_code2" db:"delv_zip_code2"`
	DelvProvinceId2         *string    `json:"delv_province_id2" sql:"-"`  // derived from delv_ward_id2
	DelvRegencyId2          *string    `json:"delv_regency_id2" sql:"-"`   // derived from delv_ward_id2
	DelvSubDistrictId2      *string    `json:"delv_sub_district_id2" sql:"-"` // derived from delv_ward_id2
	InvWardId               *string    `db:"inv_ward_id" json:"inv_ward_id"`
	InvZipCode              *string    `json:"inv_zip_code" db:"inv_zip_code"`
	InvIsSameAddress        *bool      `json:"inv_is_same_addr" db:"inv_is_same_addr" sql:"inv_is_same_addr"`
	InvProvinceId           *string    `json:"inv_province_id" sql:"-"`    // derived from inv_ward_id, not stored in m_outlet
	InvRegencyId            *string    `json:"inv_regency_id" sql:"-"`     // derived from inv_ward_id
	InvSubDistrictId        *string    `json:"inv_sub_district_id" sql:"-"` // derived from inv_ward_id
	VerificationStatus      *int       `json:"verification_status" db:"verification_status"`
	OutletEstablishmentDate *string    `db:"outlet_establishment_date" json:"outlet_establishment_date"`
}

type OutletRead struct {
	VerificationStatusName string     `json:"-" db:"-"`
	OutletPrincipalCode    *string    `json:"outlet_principal_code,omitempty" db:"outlet_principal_code"`
	IsPkpOutlet            *bool      `json:"is_pkp_outlet,omitempty" db:"is_pkp_outlet"`
	CustId                 string     `json:"cust_id" db:"cust_id"`
	OutletId               int        `json:"outlet_id" db:"outlet_id"`
	Source                 *int       `json:"source,omitempty" db:"source"`
	OutletCode             *string    `json:"outlet_code,omitempty" db:"outlet_code"`
	OutletName             *string    `json:"outlet_name,omitempty" db:"outlet_name"`
	Barcode                *string    `json:"barcode" db:"barcode"`
	OutletStatus           *int       `json:"outlet_status" db:"outlet_status"`
	OutletStatusCode       *string    `json:"-" db:"outlet_status_code"`
	OutletStatusDesc       *string    `json:"-" db:"outlet_status_desc"`
	Address1               *string    `json:"address1" db:"address1"`
	Address2               *string    `json:"address2" db:"address2"`
	City                   *string    `json:"city" db:"city"`
	ZipCode                *string    `json:"zip_code" db:"zip_code"`
	PhoneNo                *string    `json:"phone_no" db:"phone_no"`
	WaNo                   *string    `json:"wa_no" db:"wa_no"`
	FaxNo                  *string    `json:"fax_no" db:"fax_no"`
	Email                  *string    `json:"email" db:"email"`
	DiscGrpId              *int       `json:"disc_grp_id" db:"disc_grp_id"`
	DiscGrpCode            *string    `json:"disc_grp_code" db:"disc_grp_code"`
	DiscGrpName            *string    `json:"disc_grp_name" db:"disc_grp_name"`
	OtLocId                *int       `json:"ot_loc_id" db:"ot_loc_id"`
	OtLocCode              *string    `json:"ot_loc_code" db:"ot_loc_code"`
	OtLocName              *string    `json:"ot_loc_name" db:"ot_loc_name"`
	OtGrpId                *int       `json:"ot_grp_id" db:"ot_grp_id"`
	OtGrpCode              *string    `json:"ot_grp_code" db:"ot_grp_code"`
	OtGrpName              *string    `json:"ot_grp_name" db:"ot_grp_name"`
	PriceGrpId             *int       `json:"price_grp_id" db:"price_grp_id"`
	PriceGrpCode           *string    `json:"price_grp_code" db:"price_grp_code"`
	PriceGrpName           *string    `json:"price_grp_name" db:"price_grp_name"`
	DistrictId             *int       `json:"district_id" db:"district_id"`
	DistrictCode           *string    `json:"district_code" db:"district_code"`
	DistrictName           *string    `json:"district_name" db:"district_name"`
	BeatId                 *int       `json:"beat_id" db:"beat_id"`
	BeatCode               *string    `json:"beat_code" db:"beat_code"`
	BeatName               *string    `json:"beat_name" db:"beat_name"`
	SbeatId                *int       `json:"sbeat_id" db:"sbeat_id"`
	SBeatCode              *string    `json:"sbeat_code" db:"sbeat_code"`
	SBeatName              *string    `json:"sbeat_name" db:"sbeat_name"`
	OtClassId              *int       `json:"ot_class_id" db:"ot_class_id"`
	OtClassCode            *string    `json:"ot_class_code" db:"ot_class_code"`
	OtClassName            *string    `json:"ot_class_name" db:"ot_class_name"`
	IndustryId             *int       `json:"industry_id" db:"industry_id"`
	IndustryCode           *string    `json:"industry_code" db:"industry_code"`
	Industryname           *string    `json:"industry_name" db:"industry_name"`
	MarketId               *int       `json:"market_id" db:"market_id"`
	MarketCode             *string    `json:"market_code" db:"market_code"`
	MarketName             *string    `json:"market_name" db:"market_name"`
	Top                    *int       `json:"top" db:"top"`
	PaymentType            *int       `json:"payment_type" db:"payment_type"`
	PaymentTypeName        *string    `json:"payment_type_name" db:"payment_type_name"`
	IsContraBon            *bool      `json:"is_contra_bon" db:"is_contra_bon"`
	PluGrpId               *int       `json:"plu_grp_id" db:"plu_grp_id"`
	PluGrpCode             *string    `json:"plu_grp_code" db:"plu_grp_code"`
	PluGrpName             *string    `json:"plu_grp_name" db:"plu_grp_name"`
	ConvGrpId              *int       `json:"conv_grp_id" db:"conv_grp_id"`
	ConvGrpCode            *string    `json:"conv_grp_code" db:"conv_grp_code"`
	ConvGrpName            *string    `json:"conv_grp_name" db:"conv_grp_name"`
	DiscInvId              *int       `json:"disc_inv_id" db:"disc_inv_id"`
	DiscInvCode            *string    `json:"disc_inv_code" db:"disc_inv_code"`
	DiscInvName            *string    `json:"disc_inv_name" db:"disc_inv_name"`
	AgentFrom              *string    `json:"agent_from" db:"agent_from"`
	CreditLimitType        *int       `json:"credit_limit_type" db:"credit_limit_type"`
	CreditLimitTypeName    *string    `json:"credit_limit_type_name" db:"credit_limit_type_name"`
	CreditLimit            *float64   `json:"credit_limit" db:"credit_limit"`
	SalesInvLimitType      *int       `json:"sales_inv_limit_type" db:"sales_inv_limit_type"`
	SalesInvLimitTypeName  *string    `json:"sales_inv_limit_type_name" db:"sales_inv_limit_type_name"`
	SalesInvLimit          *int       `json:"sales_inv_limit" db:"sales_inv_limit"`
	AvgSalesWeek           *float64   `json:"avg_sales_week" db:"avg_sales_week"`
	AvgSalesMonth          *float64   `json:"avg_sales_month" db:"avg_sales_month"`
	FirstTransDate         *string    `json:"first_trans_date" db:"first_trans_date"`
	LastTransDate          *string    `json:"last_trans_date" db:"last_trans_date"`
	PrevTransDate          *string    `json:"prev_trans_date,omitempty" db:"prev_trans_date"`
	FirstWeekNo            *int       `json:"first_week_no" db:"first_week_no"`
	OtStartDate            *string    `json:"ot_start_date" db:"ot_start_date"`
	OtRegDate              *string    `json:"ot_reg_date" db:"ot_reg_date"`
	BuldingOwn             *int       `json:"building_own" db:"building_own"`
	Dob                    *string    `json:"dob" db:"dob"`
	ArStatus               *int       `json:"ar_status" db:"ar_status"`
	ArTotal                *float64   `json:"ar_total" db:"ar_total"`
	CloseDate              *time.Time `json:"closed_date" db:"closed_date"`
	IsEmbBail              *bool      `json:"is_emb_bail" db:"is_emb_bail"`
	TaxName                *string    `json:"tax_name" db:"tax_name"`
	TaxAddr1               *string    `json:"tax_addr1" db:"tax_addr1"`
	TaxAddr2               *string    `json:"tax_addr2" db:"tax_addr2"`
	TaxCity                *string    `json:"tax_city" db:"tax_city"`
	TaxNo                  *string    `json:"tax_no" db:"tax_no"`
	TaxInvoiceForm         *int       `json:"tax_invoice_form" db:"tax_invoice_form"`
	OwnerName              *string    `json:"owner_name" db:"owner_name"`
	OwnerAdd1              *string    `json:"owner_addr1" db:"owner_addr1"`
	OwnerAddr2             *string    `json:"owner_addr2" db:"owner_addr2"`
	OwnerCity              *string    `json:"owner_city" db:"owner_city"`
	OwnerPhoneNo           *string    `json:"owner_phone_no" db:"owner_phone_no"`
	OwnerIdNo              *string    `json:"owner_id_no" db:"owner_id_no"`
	DelvAdd1               *string    `json:"delv_addr1" db:"delv_addr1"`
	DelvCity               *string    `json:"delv_city" db:"delv_city"`
	DelvLatitude           *string    `json:"delv_latitude" db:"delv_latitude"`
	DelvLongitude          *string    `json:"delv_longitude" db:"delv_longitude"`
	DelvAddr2              *string    `json:"delv_addr2" db:"delv_addr2"`
	DelvCity2              *string    `json:"delv_city2" db:"delv_city2"`
	DelvLatitude2          *string    `json:"delv_latitude2" db:"delv_latitude2"`
	DelvLongitude2         *string    `json:"delv_longitude2" db:"delv_longitude2"`
	InvAddr1               *string    `json:"inv_addr1" db:"inv_addr1"`
	InvAddr2               *string    `json:"inv_addr2" db:"inv_addr2"`
	InvCity                *string    `json:"inv_city" db:"inv_city"`
	Latitude               *string    `json:"latitude" db:"latitude"`
	Longitude              *string    `json:"longitude" db:"longitude"`
	ImageUrl               *string    `json:"image_url" db:"image_url"`
	FileUrl                *string    `json:"file_url" db:"file_url"`
	IsActive               bool       `db:"is_active" json:"is_active"`
	IsDel                  bool       `db:"is_del" json:"is_del"`
	CreatedBy              *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedByName          *string    `db:"created_by_name" json:"created_by_name"`
	CreatedAt              *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy              *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedByName          *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt              *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	DeletedBy              *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt              *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	OtTypeId               *int       `db:"ot_type_id" json:"ot_type_id"`
	OtTypeCode             *string    `db:"ot_type_code" json:"ot_type_code"`
	OtTypeName             *string    `db:"ot_type_name" json:"ot_type_name"`
	IsObs                  *bool      `db:"is_obs" json:"is_obs"`
	Obs                    *int       `db:"obs" json:"obs"`

	OutletWardId            *string    `db:"outlet_ward_id" json:"outlet_ward_id"`
	OutletWard              *string    `db:"outlet_ward" json:"outlet_ward"`
	OutletSubDistrictId     *string    `db:"outlet_sub_district_id" json:"outlet_sub_district_id"`
	OutletSubDistrict       *string    `db:"outlet_sub_district" json:"outlet_sub_district"`
	OutletRegencyId         *string    `db:"outlet_regency_id" json:"outlet_regency_id"`
	OutletRegency           *string    `db:"outlet_regency" json:"outlet_regency"`
	OutletProvinceId        *string    `db:"outlet_province_id" json:"outlet_province_id"`
	OutletProvince          *string    `db:"outlet_province" json:"outlet_province"`
	IsWaNo                  *bool      `db:"is_wa_no" json:"is_wa_no"`
	DelvWardId              *string    `db:"delv_ward_id" json:"delv_ward_id"`
	DelvWard                *string    `db:"delv_ward" json:"delv_ward"`
	DelvSubDistrictId       *string    `db:"delv_sub_district_id" json:"delv_sub_district_id"`
	DelvSubDistrict         *string    `db:"delv_sub_district" json:"delv_sub_district"`
	DelvRegencyId           *string    `db:"delv_regency_id" json:"delv_regency_id"`
	DelvRegency             *string    `db:"delv_regency" json:"delv_regency"`
	DelvProvinceId          *string    `db:"delv_province_id" json:"delv_province_id"`
	DelvProvince            *string    `db:"delv_province" json:"delv_province"`
	DelvZipCode             *string    `db:"delv_zip_code" json:"delv_zip_code"`
	DelvIsSameAddress       *bool      `db:"delv_is_same_addr" json:"delv_is_same_addr"`
	DelvWardId2             *string    `db:"delv_ward_id2" json:"delv_ward_id2"`
	DelvWard2               *string    `db:"delv_ward2" json:"delv_ward2"`
	DelvSubDistrictId2      *string    `db:"delv_sub_district_id2" json:"delv_sub_district_id2"`
	DelvSubDistrict2        *string    `db:"delv_sub_district2" json:"delv_sub_district2"`
	DelvRegencyId2          *string    `db:"delv_regency_id2" json:"delv_regency_id2"`
	DelvRegency2            *string    `db:"delv_regency2" json:"delv_regency2"`
	DelvProvinceId2         *string    `db:"delv_province_id2" json:"delv_province_id2"`
	DelvProvince2           *string    `db:"delv_province2" json:"delv_province2"`
	DelvZipCode2            *string    `db:"delv_zip_code2" json:"delv_zip_code2"`
	InvWardId               *string    `db:"inv_ward_id" json:"inv_ward_id"`
	InvWard                 *string    `db:"inv_ward" json:"inv_ward"`
	InvSubDistrictId        *string    `db:"inv_sub_district_id" json:"inv_sub_district_id"`
	InvSubDistrict          *string    `db:"inv_sub_district" json:"inv_sub_district"`
	InvRegencyId            *string    `db:"inv_regency_id" json:"inv_regency_id"`
	InvRegency              *string    `db:"inv_regency" json:"inv_regency"`
	InvProvinceId           *string    `db:"inv_province_id" json:"inv_province_id"`
	InvProvince             *string    `db:"inv_province" json:"inv_province"`
	InvZipCode              *string    `db:"inv_zip_code" json:"inv_zip_code"`
	InvIsSameAddress        *bool      `db:"inv_is_same_addr" json:"inv_is_same_addr"`
	VerificationStatus      *int       `db:"verification_status" json:"verification_status"`
	VerifiedBy              *int64     `db:"verified_by,omitempty" json:"verified_by"`
	VerifiedByName          *string    `db:"verified_by_name" json:"verified_by_name"`
	VerifiedAt              *time.Time `db:"verified_at,omitempty" json:"verified_at"`
	ObsType                 *int       `db:"obs_type" json:"obs_type"`
	ObsTypeName             *string    `db:"obs_type_name" json:"obs_type_name"`
	CreditLimitAction       *int       `db:"credit_limit_action" json:"credit_limit_action"`
	CreditLimitActionName   *string    `db:"credit_limit_action_name" json:"credit_limit_action_name"`
	SalesInvLimitAction     *int       `db:"sales_inv_limit_action" json:"sales_inv_limit_action"`
	SalesInvLimitActionName *string    `db:"sales_inv_limit_action_name" json:"sales_inv_limit_action_name"`
	ObsLimitAction          *int       `db:"obs_limit_action" json:"obs_limit_action"`
	ObsLimitActionName      *string    `db:"obs_limit_action_name" json:"obs_limit_action_name"`
	OutletEstablishmentDate *string    `db:"outlet_establishment_date" json:"outlet_establishment_date"`
	IdentityType            *string    `db:"identity_type" json:"identity_type"`
	IdentityNo              *string    `db:"identity_no" json:"identity_no"`
	TaxInvoiceId            *int64     `db:"tax_invoice_id" json:"tax_invoice_id"`
	TaxType                 *string    `db:"tax_type" json:"tax_type"`
	Nitku                   *string    `db:"nitku" json:"nitku"`
	AddressTax              *string    `db:"address_tax" json:"address_tax"`
	TaxIdentifierType       *string    `db:"tax_identifier_type" json:"tax_identifier_type"`
	TaxIdentifierNo         *string    `db:"tax_identifier_no" json:"tax_identifier_no"`
	BankId                  *string    `db:"bank_id" json:"bank_id"`
	BankCode                *string    `db:"bank_code" json:"bank_code"`
	BankName                *string    `db:"bank_name" json:"bank_name"`
	AccountNo               *string    `db:"account_no" json:"account_no"`
	AccountName             *string    `db:"account_name" json:"account_name"`
	ContactName             *string    `db:"contact_name" json:"contact_name"`
	JobTitle                *string    `db:"job_title" json:"job_title"`
	ContactPhoneNo          *string    `db:"contact_phone_no" json:"contact_phone_no"`
	ContactWaNo             *string    `db:"contact_wa_no" json:"contact_wa_no"`
	ContactEmail            *string    `db:"contact_email" json:"contact_email"`
	ContactIsWaNo           *bool      `db:"contact_is_wa_no" json:"contact_is_wa_no"`
	PreDormantStatus        *int       `db:"pre_dormant_status" json:"pre_dormant_status,omitempty"`
}

// Alias used by export-only repository method so we can reuse the same struct
// without duplicating fields. This keeps service export creators working.
// NOTE: This is a type alias (not a new type), so []OutletExport == []OutletRead.
type OutletExport = OutletRead

type OutletApprove struct {
	VerifiedAt         *time.Time `json:"verified_at" sql:"verified_at"`
	VerifiedBy         *int64     `json:"verified_by" sql:"verified_by"`
	VerificationStatus *int       `json:"verification_status" db:"verification_status"`
}

type OutletReject struct {
	VerifiedAt         *time.Time `json:"verified_at" sql:"verified_at"`
	VerifiedBy         *int64     `json:"verified_by" sql:"verified_by"`
	VerificationStatus *int       `json:"verification_status" db:"verification_status"`
}
