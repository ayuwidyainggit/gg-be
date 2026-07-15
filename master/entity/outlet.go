package entity

import (
	"encoding/json"
	"mime/multipart"
	"strings"
	"time"
)

type FlexString string

func (s *FlexString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*s = ""
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = FlexString(str)
		return nil
	}
	// number (e.g. 0)
	*s = FlexString(strings.TrimSpace(string(data)))
	return nil
}

// COMMENT ON COLUMN "mst"."m_outlet"."ar_status" IS '1: Normal; 2: 1x Giro tolakan; 3: 2x Giro tolakan; 4: 3x Giro tolakan';
var (
	AR_TYPE_NORMAL          string = "Normal"
	AR_TYPE_1X_GIRO_TOLAKAN string = "1x Giro tolakan"
	AR_TYPE_2X_GIRO_TOLAKAN string = "2x Giro tolakan"
	AR_TYPE_3X_GIRO_TOLAKAN string = "3x Giro tolakan"

	AR_TYPE_NORMAL_ID          int = 1
	AR_TYPE_1X_GIRO_TOLAKAN_ID int = 2
	AR_TYPE_2X_GIRO_TOLAKAN_ID int = 3
	AR_TYPE_3X_GIRO_TOLAKAN_ID int = 4

	OUTLET_STATUS_ACTIVE   int = 1
	OUTLET_STATUS_DEACTIVE int = 2
)

// master_outlet

type ImportRow struct {
	CustId                  string `json:"cust_id"`
	OutletCode              string `json:"outlet_code"`
	OutletName              string `json:"outlet_name"`
	Address1                string `json:"address1"`
	Address2                string `json:"address2"`
	City                    string `json:"city"`
	PhoneNo                 string `json:"phone_no"`
	WaNo                    string `json:"wa_no"`
	Email                   string `json:"email"`
	ZipCode                 string `json:"zip_code"`
	FaxNo                   string `json:"fax_no"`
	DiscGrpId               string `json:"disc_grp_id"`
	DiscGrpCode             string `json:"disc_grp_code"`
	DiscGrpName             string `json:"disc_grp_name"`
	OtLocId                 string `json:"ot_loc_id"`
	OtLocCode               string `json:"ot_loc_code"`
	OtLocName               string `json:"ot_loc_name"`
	OtGrpId                 int64  `json:"ot_grp_id"`
	OtGrpCode               string `json:"ot_grp_code"`
	OtGrpName               string `json:"ot_grp_name"`
	PriceGrpId              string `json:"price_grp_id"`
	PriceGrpCode            string `json:"price_grp_code"`
	PriceGrpName            string `json:"price_grp_name"`
	DistrictId              int64  `json:"district_id"`
	DistrictCode            string `json:"district_code"`
	DistrictName            string `json:"district_name"`
	BeatId                  string `json:"beat_id"`
	SbeatId                 string `json:"sbeat_id"`
	OtClassId               int64  `json:"ot_class_id"`
	OtClassCode             string `json:"ot_class_code"`
	OtClassName             string `json:"ot_class_name"`
	IndustryId              int64  `json:"industry_id"`
	IndustryCode            string `json:"industry_code"`
	IndustryName            string `json:"industry_name"`
	MarketId                int64  `json:"market_id"`
	MarketCode              string `json:"market_code"`
	MarketName              string `json:"market_name"`
	Top                     string `json:"top"`
	PaymentType             string `json:"payment_type"`
	PaymentTypeName         string `json:"payment_type_name"`
	IsContraBon             string `json:"is_contra_bon"`
	PluGrpId                string `json:"plu_grp_id"`
	ConvGrpId               string `json:"conv_grp_id"`
	DiscInvId               string `json:"disc_inv_id"`
	AgentFrom               string `json:"agent_from"`
	CreditLimitType         string `json:"credit_limit_type"`
	CreditLimitTypeName     string `json:"credit_limit_type_name"`
	CreditLimit             string `json:"credit_limit"`
	SalesInvLimitType       string `json:"sales_inv_limit_type"`
	SalesInvLimitTypeName   string `json:"sales_inv_limit_type_name"`
	SalesInvLimit           string `json:"sales_inv_limit"`
	AvgSalesWeek            string `json:"avg_sales_week"`
	AvgSalesMonth           string `json:"avg_sales_month"`
	FirstTransDate          string `json:"first_trans_date"`
	LastTransDate           string `json:"last_trans_date"`
	FirstWeekNo             string `json:"first_week_no"`
	OtStartDate             string `json:"ot_start_date"`
	OtRegDate               string `json:"ot_reg_date"`
	BuildingOwn             string `json:"building_own"`
	Dob                     string `json:"dob"`
	ArStatus                string `json:"ar_status"`
	ArStatusName            string `json:"ar_status_name"`
	ArTotal                 string `json:"ar_total"`
	ClosedDate              string `json:"closed_date"`
	IsEmbBail               string `json:"is_emb_bail"`
	TaxName                 string `json:"tax_name"`
	TaxAddr1                string `json:"tax_addr1"`
	TaxAddr2                string `json:"tax_addr2"`
	TaxCity                 string `json:"tax_city"`
	TaxNo                   string `json:"tax_no"`
	OwnerName               string `json:"owner_name"`
	OwnerAddr1              string `json:"owner_addr1"`
	OwnerAddr2              string `json:"owner_addr2"`
	OwnerCity               string `json:"owner_city"`
	OwnerPhoneNo            string `json:"owner_phone_no"`
	OwnerIdNo               string `json:"owner_id_no"`
	DelvAddr1               string `json:"delv_addr1"`
	DelvAddr2               string `json:"delv_addr2"`
	DelvCity                string `json:"delv_city"`
	InvAddr1                string `json:"inv_addr1"`
	InvAddr2                string `json:"inv_addr2"`
	InvCity                 string `json:"inv_city"`
	IsActive                string `json:"is_active"`
	IsDel                   string `json:"is_del"`
	Latitude                string `json:"latitude"`
	Longitude               string `json:"longitude"`
	ImageUrl                string `json:"image_url"`
	OtTypeId                int64  `json:"ot_type_id"`
	OtTypeCode              string `json:"ot_type_code"`
	OtTypeName              string `json:"ot_type_name"`
	IsObs                   string `json:"is_obs"`
	Obs                     string `json:"obs"`
	OutletWardId            string `json:"outlet_ward_id"`
	OutletWard              string `json:"outlet_ward"`
	OutletSubDistrictId     string `json:"outlet_sub_district_id"`
	OutletSubDistrict       string `json:"outlet_sub_district"`
	OutletRegencyId         string `json:"outlet_regency_id"`
	OutletRegency           string `json:"outlet_regency"`
	OutletProvinceId        string `json:"outlet_province_id"`
	OutletProvince          string `json:"outlet_province"`
	IsWaNo                  string `json:"is_wa_no"`
	DelvWardId              string `json:"delv_ward_id"`
	DelvWard                string `json:"delv_ward"`
	DelvSubDistrictId       string `json:"delv_sub_district_id"`
	DelvSubDistrict         string `json:"delv_sub_district"`
	DelvRegencyId           string `json:"delv_regency_id"`
	DelvRegency             string `json:"delv_regency"`
	DelvProvinceId          string `json:"delv_province_id"`
	DelvProvince            string `json:"delv_province"`
	DelvZipCode             string `json:"delv_zip_code"`
	DelvIsSameAddr          string `json:"delv_is_same_addr"`
	InvWardId               string `json:"inv_ward_id"`
	InvWard                 string `json:"inv_ward"`
	InvSubDistrictId        string `json:"inv_sub_district_id"`
	InvSubDistrict          string `json:"inv_sub_district"`
	InvRegencyId            string `json:"inv_regency_id"`
	InvRegency              string `json:"inv_regency"`
	InvProvinceId           string `json:"inv_province_id"`
	InvProvince             string `json:"inv_province"`
	InvZipCode              string `json:"inv_zip_code"`
	InvIsSameAddr           string `json:"inv_is_same_addr"`
	VerificationStatus      string `json:"verification_status"`
	TaxInvoiceForm          string `json:"tax_invoice_form"`
	TaxInvoiceFormName      string `json:"tax_invoice_form_name"`
	ObsType                 string `json:"obs_type"`
	ObyTypeName             string `json:"oby_type_name"`
	CreditLimitAction       string `json:"credit_limit_action"`
	CreditLimitActionName   string `json:"credit_limit_action_name"`
	SalesInvLimitAction     string `json:"sales_inv_limit_action"`
	SalesInvLimitActionName string `json:"sales_inv_limit_action_name"`
	ObsLimitAction          string `json:"obs_limit_action"`
	ObsLimitActionName      string `json:"obs_limit_action_name"`
	OutletEstablishmentDate string `json:"outlet_establishment_date"`
	DelvCity2               string `json:"delv_city2"`
	DelvLatitude            string
	DelvLongitude           string
	DelvLatitude2           string
	DelvLongitude2          string
	DelvWardId2             string
	DelvZipCode2            string
	BankId                  string
	BankCode                string
	BankName                string
	AccountNo               string
	AccountName             string
	ContactName             string
	JobTitle                string
	ContactPhoneNo          string
	ContactWaNo             string
	ContactEmail            string
	IdentityNo              string
	ContactIsWaNo           string
	IdentityType            string
	FaxNumber               string
	TaxInvoiceId            string
	IsEmbBail2              string
	TaxNo2                  string
	TaxName2                string
	TaxCity2                string
	TaxAddr1_2              string
	TaxAddr2_2              string
	TaxType                 string
	Nitku                   string
	AddressTax              string
	TaxIdentifierType       string
	TaxIdentifierNo         string
	StatusInsert            string
	UpdatedAt               *time.Time
	DeletedAt               *time.Time
}

// ProcessedOutlet represents processed outlet data ready for insertion
type ProcessedOutlet struct {
	CustId                  string     `db:"cust_id"`
	OutletId                int64      `db:"outlet_id"`
	OutletCode              string     `db:"outlet_code"`
	Barcode                 string     `db:"barcode"`
	OutletName              string     `db:"outlet_name"`
	OutletStatus            int16      `db:"outlet_status"`
	Address1                string     `db:"address1"`
	Address2                string     `db:"address2"`
	City                    string     `db:"city"`
	ZipCode                 string     `db:"zip_code"`
	PhoneNo                 string     `db:"phone_no"`
	WaNo                    string     `db:"wa_no"`
	FaxNo                   string     `db:"fax_no"`
	Email                   string     `db:"email"`
	DiscGrpId               int64      `db:"disc_grp_id"`
	OtLocId                 int64      `db:"ot_loc_id"`
	OtGrpId                 int64      `db:"ot_grp_id"`
	PriceGrpId              int64      `db:"price_grp_id"`
	DistrictId              int64      `db:"district_id"`
	BeatId                  int64      `db:"beat_id"`
	SbeatId                 int64      `db:"sbeat_id"`
	OtClassId               int64      `db:"ot_class_id"`
	IndustryId              int64      `db:"industry_id"`
	MarketId                int64      `db:"market_id"`
	Top                     int32      `db:"top"`
	PaymentType             int16      `db:"payment_type"`
	IsContraBon             bool       `db:"is_contra_bon"`
	PluGrpId                int64      `db:"plu_grp_id"`
	ConvGrpId               int64      `db:"conv_grp_id"`
	DiscInvId               int64      `db:"disc_inv_id"`
	AgentFrom               string     `db:"agent_from"`
	CreditLimitType         int16      `db:"credit_limit_type"`
	CreditLimitTypeName     string     `db:"credit_limit_type_name"`
	CreditLimit             float64    `db:"credit_limit"`
	SalesInvLimitType       int16      `db:"sales_inv_limit_type"`
	SalesInvLimitTypeName   string     `db:"sales_inv_limit_type_name"`
	SalesInvLimit           int16      `db:"sales_inv_limit"`
	AvgSalesWeek            float64    `db:"avg_sales_week"`
	AvgSalesMonth           float64    `db:"avg_sales_month"`
	FirstTransDate          *time.Time `db:"first_trans_date"`
	LastTransDate           *time.Time `db:"last_trans_date"`
	FirstWeekNo             int16      `db:"first_week_no"`
	OtStartDate             *time.Time `db:"ot_start_date"`
	OtRegDate               *time.Time `db:"ot_reg_date"`
	BuildingOwn             int16      `db:"building_own"`
	Dob                     *time.Time `db:"dob"`
	ArStatus                int16      `db:"ar_status"`
	ArTotal                 float64    `db:"ar_total"`
	ClosedDate              *time.Time `db:"closed_date"`
	IsEmbBail               bool       `db:"is_emb_bail"`
	IsPkpOutlet             bool       `db:"is_pkp_outlet"`
	TaxName                 string     `db:"tax_name"`
	TaxAddr1                string     `db:"tax_addr1"`
	TaxAddr2                string     `db:"tax_addr2"`
	TaxCity                 string     `db:"tax_city"`
	TaxNo                   string     `db:"tax_no"`
	OwnerName               string     `db:"owner_name"`
	OwnerAddr1              string     `db:"owner_addr1"`
	OwnerAddr2              string     `db:"owner_addr2"`
	OwnerCity               string     `db:"owner_city"`
	OwnerPhoneNo            string     `db:"owner_phone_no"`
	OwnerIdNo               string     `db:"owner_id_no"`
	DelvAddr1               string     `db:"delv_addr1"`
	DelvAddr2               string     `db:"delv_addr2"`
	DelvCity                string     `db:"delv_city"`
	InvAddr1                string     `db:"inv_addr1"`
	InvAddr2                string     `db:"inv_addr2"`
	InvCity                 string     `db:"inv_city"`
	IsActive                bool       `db:"is_active"`
	CreatedBy               *int64     `db:"created_by"`
	CreatedAt               *time.Time `db:"created_at"`
	UpdatedBy               *int64     `db:"updated_by"`
	UpdatedAt               *time.Time `db:"updated_at"`
	IsDel                   bool       `db:"is_del"`
	DeletedBy               *int64     `db:"deleted_by"`
	DeletedAt               *time.Time `db:"deleted_at"`
	Latitude                string     `db:"latitude"`
	Longitude               string     `db:"longitude"`
	ImageUrl                string     `db:"image_url"`
	OtTypeId                int64      `db:"ot_type_id"`
	IsObs                   bool       `db:"is_obs"`
	Obs                     int64      `db:"obs"`
	OutletWardId            string     `db:"outlet_ward_id"`
	IsWaNo                  bool       `db:"is_wa_no"`
	DelvWardId              string     `db:"delv_ward_id"`
	DelvZipCode             string     `db:"delv_zip_code"`
	DelvIsSameAddr          bool       `db:"delv_is_same_addr"`
	InvWardId               string     `db:"inv_ward_id"`
	InvZipCode              string     `db:"inv_zip_code"`
	InvIsSameAddr           bool       `db:"inv_is_same_addr"`
	VerificationStatus      int16      `db:"verification_status"`
	VerifiedAt              *time.Time `db:"verified_at"`
	VerifiedBy              *int64     `db:"verified_by"`
	TaxInvoiceForm          int16      `db:"tax_invoice_form"`
	ObsType                 int64      `db:"obs_type"`
	CreditLimitAction       int64      `db:"credit_limit_action"`
	CreditLimitActionName   string     `db:"credit_limit_action_name"`
	SalesInvLimitAction     int64      `db:"sales_inv_limit_action"`
	SalesInvLimitActionName string     `db:"sales_inv_limit_action_name"`
	ObsLimitAction          int64      `db:"obs_limit_action"`
	OutletProvinceId        string     `db:"outlet_province_id"`
	OutletRegencyId         string     `db:"outlet_regency_id"`
	OutletSubDistrictId     string     `db:"outlet_sub_district_id"`
	OutletEstablishmentDate *time.Time `db:"outlet_establishment_date"`
	DelvCity2               string     `db:"delv_city2"`
	DelvLatitude            string     `db:"delv_latitude"`
	DelvLongitude           string     `db:"delv_longitude"`
	DelvLatitude2           string     `db:"delv_latitude2"`
	DelvLongitude2          string     `db:"delv_longitude2"`
	DelvWardId2             string     `db:"delv_ward_id2"`
	DelvZipCode2            string     `db:"delv_zip_code2"`
	OutletPrincipalCode     string     `db:"outlet_principal_code"`
}

// Temp_outlet

type ProcessedRow struct {
	HistoryId                string
	OutletCode               string
	OutletName               string
	Barcode                  string
	OutletStatus             string
	Address1                 string
	Address2                 string
	City                     string
	ZipCode                  string
	PhoneNo                  string
	WaNo                     string
	FaxNo                    string
	Email                    string
	DiscGrpId                int64
	DiscGrpCode              string
	DiscGrpName              string
	OtLocId                  int64
	OtLocCode                string
	OtLocName                string
	OtGrpId                  int64
	OtGrpCode                string
	OtGrpName                string
	PriceGrpId               int64
	PriceGrpCode             string
	PriceGrpName             string
	DistrictId               int64
	DistrictCode             string
	DistrictName             string
	BeatId                   int64
	SbeatId                  int64
	OtClassId                int64
	OtClassCode              string
	OtClassName              string
	IndustryId               int64
	IndustryCode             string
	IndustryName             string
	MarketId                 int64
	MarketCode               string
	MarketName               string
	Top                      string
	PaymentType              string
	PaymentTypeName          string
	IsContraBon              bool
	PluGrpId                 int64
	ConvGrpId                int64
	DiscInvId                int64
	AgentFrom                string
	CreditLimitType          string
	CreditLimitTypeName      string
	CreditLimit              float64
	SalesInvLimitType        string
	SalesInvLimitTypeName    string
	SalesInvLimit            float64
	AvgSalesWeek             float64
	AvgSalesMonth            float64
	FirstTransDate           *time.Time
	LastTransDate            *time.Time
	FirstWeekNo              int
	OtStartDate              *time.Time
	OtRegDate                *time.Time
	BuildingOwn              string
	Dob                      *time.Time
	ArStatus                 string
	ArStatusName             string
	ArTotal                  float64
	ClosedDate               *time.Time
	IsEmbBail                bool
	TaxName                  string
	TaxAddr1                 string
	TaxAddr2                 string
	TaxCity                  string
	TaxNo                    string
	OwnerName                string
	OwnerAddr1               string
	OwnerAddr2               string
	OwnerCity                string
	OwnerPhoneNo             string
	OwnerIdNo                string
	DelvAddr1                string
	DelvAddr2                string
	DelvCity                 string
	InvAddr1                 string
	InvAddr2                 string
	InvCity                  string
	IsActive                 bool
	IsDel                    bool
	Latitude                 float64
	Longitude                float64
	ImageUrl                 string
	OtTypeId                 int64
	OtTypeCode               string
	OtTypeName               string
	IsObs                    bool
	Obs                      string
	OutletWardId             int64
	OutletWard               string
	OutletSubDistrictId      int64
	OutletSubDistrict        string
	OutletRegencyId          int64
	OutletRegency            string
	OutletProvinceId         int64
	OutletProvince           string
	IsWaNo                   bool
	DelvWardId               int64
	DelvWard                 string
	DelvSubDistrictId        int64
	DelvSubDistrict          string
	DelvRegencyId            int64
	DelvRegency              string
	DelvProvinceId           int64
	DelvProvince             string
	DelvZipCode              string
	DelvIsSameAddr           bool
	InvWardId                int64
	InvWard                  string
	InvSubDistrictId         int64
	InvSubDistrict           string
	InvRegencyId             int64
	InvRegency               string
	InvProvinceId            int64
	InvProvince              string
	InvZipCode               string
	InvIsSameAddr            bool
	VerificationStatus       string
	TaxInvoiceForm           string
	TaxInvoiceFormName       string
	ObsType                  string
	ObyTypeName              string
	CreditLimitAction        string
	CreditLimitActionName    string
	SalesInvLimitAction      string
	SalesInvLimitActionName  string
	ObsLimitAction           string
	ObsLimitActionName       string
	OutletEstablishmentDate  *time.Time
	MOutletEstablishmentDate *time.Time
	DelvCity2                string
	DelvLatitude             float64
	DelvLongitude            float64
	DelvLatitude2            float64
	DelvLongitude2           float64
	DelvWardId2              int64
	DelvZipCode2             string
	BankId                   int64
	BankCode                 string
	BankName                 string
	AccountNo                string
	AccountName              string
	ContactName              string
	JobTitle                 string
	ContactPhoneNo           string
	ContactWaNo              string
	ContactEmail             string
	IdentityNo               string
	ContactIsWaNo            bool
	IdentityType             string
	FaxNumber                string
	TaxInvoiceId             int64
	IsEmbBail2               bool
	TaxNo2                   string
	TaxName2                 string
	TaxCity2                 string
	TaxAddr1_2               string
	TaxAddr2_2               string
	TaxType                  string
	Nitku                    string
	AddressTax               string
	TaxIdentifierType        string
	TaxIdentifierNo          string
	StatusInsert             string
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
}

type OutletTemp struct {
	CustId                   string     `db:"cust_id"`
	HistoryId                int64      `db:"history_id"`
	OutletCode               string     `db:"outlet_code"`
	OutletName               string     `db:"outlet_name"`
	Barcode                  string     `db:"barcode"`
	OutletStatus             string     `db:"outlet_status"`
	Address1                 string     `db:"address1"`
	Address2                 string     `db:"address2"`
	City                     string     `db:"city"`
	ZipCode                  string     `db:"zip_code"`
	PhoneNo                  string     `db:"phone_no"`
	WaNo                     string     `db:"wa_no"`
	FaxNo                    string     `db:"fax_no"`
	Email                    string     `db:"email"`
	DiscGrpId                string     `db:"disc_grp_id"`
	DiscGrpCode              string     `db:"disc_grp_code"`
	DiscGrpName              string     `db:"disc_grp_name"`
	OtLocId                  string     `db:"ot_loc_id"`
	OtLocCode                string     `db:"ot_loc_code"`
	OtLocName                string     `db:"ot_loc_name"`
	OtGrpId                  string     `db:"ot_grp_id"`
	OtGrpCode                string     `db:"ot_grp_code"`
	OtGrpName                string     `db:"ot_grp_name"`
	PriceGrpId               string     `db:"price_grp_id"`
	PriceGrpCode             string     `db:"price_grp_code"`
	PriceGrpName             string     `db:"price_grp_name"`
	DistrictId               string     `db:"district_id"`
	DistrictCode             string     `db:"district_code"`
	DistrictName             string     `db:"district_name"`
	BeatId                   string     `db:"beat_id"`
	SbeatId                  string     `db:"sbeat_id"`
	OtClassId                string     `db:"ot_class_id"`
	OtClassCode              string     `db:"ot_class_code"`
	OtClassName              string     `db:"ot_class_name"`
	IndustryId               string     `db:"industry_id"`
	IndustryCode             string     `db:"industry_code"`
	IndustryName             string     `db:"industry_name"`
	MarketId                 string     `db:"market_id"`
	MarketCode               string     `db:"market_code"`
	MarketName               string     `db:"market_name"`
	Top                      string     `db:"top"`
	PaymentType              string     `db:"payment_type"`
	PaymentTypeName          string     `db:"payment_type_name"`
	IsContraBon              string     `db:"is_contra_bon"`
	PluGrpId                 string     `db:"plu_grp_id"`
	ConvGrpId                string     `db:"conv_grp_id"`
	DiscInvId                string     `db:"disc_inv_id"`
	AgentFrom                string     `db:"agent_from"`
	CreditLimitType          string     `db:"credit_limit_type"`
	CreditLimitTypeName      string     `db:"credit_limit_type_name"`
	CreditLimit              string     `db:"credit_limit"`
	SalesInvLimitType        string     `db:"sales_inv_limit_type"`
	SalesInvLimitTypeName    string     `db:"sales_inv_limit_type_name"`
	SalesInvLimit            string     `db:"sales_inv_limit"`
	AvgSalesWeek             string     `db:"avg_sales_week"`
	AvgSalesMonth            string     `db:"avg_sales_month"`
	FirstTransDate           string     `db:"first_trans_date"`
	LastTransDate            string     `db:"last_trans_date"`
	FirstWeekNo              string     `db:"first_week_no"`
	OtStartDate              string     `db:"ot_start_date"`
	OtRegDate                string     `db:"ot_reg_date"`
	BuildingOwn              string     `db:"building_own"`
	Dob                      string     `db:"dob"`
	ArStatus                 string     `db:"ar_status"`
	ArStatusName             string     `db:"ar_status_name"`
	ArTotal                  string     `db:"ar_total"`
	ClosedDate               string     `db:"closed_date"`
	IsEmbBail                string     `db:"is_emb_bail"`
	TaxName                  string     `db:"tax_name"`
	TaxAddr1                 string     `db:"tax_addr1"`
	TaxAddr2                 string     `db:"tax_addr2"`
	TaxCity                  string     `db:"tax_city"`
	TaxNo                    string     `db:"tax_no"`
	OwnerName                string     `db:"owner_name"`
	OwnerAddr1               string     `db:"owner_addr1"`
	OwnerAddr2               string     `db:"owner_addr2"`
	OwnerCity                string     `db:"owner_city"`
	OwnerPhoneNo             string     `db:"owner_phone_no"`
	OwnerIdNo                string     `db:"owner_id_no"`
	DelvAddr1                string     `db:"delv_addr1"`
	DelvAddr2                string     `db:"delv_addr2"`
	DelvCity                 string     `db:"delv_city"`
	InvAddr1                 string     `db:"inv_addr1"`
	InvAddr2                 string     `db:"inv_addr2"`
	InvCity                  string     `db:"inv_city"`
	IsActive                 string     `db:"is_active"`
	IsDel                    string     `db:"is_del"`
	Latitude                 string     `db:"latitude"`
	Longitude                string     `db:"longitude"`
	ImageUrl                 string     `db:"image_url"`
	OtTypeId                 string     `db:"ot_type_id"`
	OtTypeCode               string     `db:"ot_type_code"`
	OtTypeName               string     `db:"ot_type_name"`
	IsObs                    string     `db:"is_obs"`
	Obs                      string     `db:"obs"`
	ObsTypeName              string     `db:"obs_type_name"`
	OutletWardId             string     `db:"outlet_ward_id"`
	OutletWard               string     `db:"outlet_ward"`
	OutletSubDistrictId      string     `db:"outlet_sub_district_id"`
	OutletSubDistrict        string     `db:"outlet_sub_district"`
	OutletRegencyId          string     `db:"outlet_regency_id"`
	OutletRegency            string     `db:"outlet_regency"`
	OutletProvinceId         string     `db:"outlet_province_id"`
	OutletProvince           string     `db:"outlet_province"`
	IsWaNo                   string     `db:"is_wa_no"`
	DelvWardId               string     `db:"delv_ward_id"`
	DelvWard                 string     `db:"delv_ward"`
	DelvSubDistrictId        string     `db:"delv_sub_district_id"`
	DelvSubDistrict          string     `db:"delv_sub_district"`
	DelvRegencyId            string     `db:"delv_regency_id"`
	DelvRegency              string     `db:"delv_regency"`
	DelvProvinceId           string     `db:"delv_province_id"`
	DelvProvince             string     `db:"delv_province"`
	DelvZipCode              string     `db:"delv_zip_code"`
	DelvIsSameAddr           string     `db:"delv_is_same_addr"`
	InvWardId                string     `db:"inv_ward_id"`
	InvWard                  string     `db:"inv_ward"`
	InvSubDistrictId         string     `db:"inv_sub_district_id"`
	InvSubDistrict           string     `db:"inv_sub_district"`
	InvRegencyId             string     `db:"inv_regency_id"`
	InvRegency               string     `db:"inv_regency"`
	InvProvinceId            string     `db:"inv_province_id"`
	InvProvince              string     `db:"inv_province"`
	InvZipCode               string     `db:"inv_zip_code"`
	InvIsSameAddr            string     `db:"inv_is_same_addr"`
	VerificationStatus       string     `db:"verification_status"`
	TaxInvoiceForm           string     `db:"tax_invoice_form"`
	TaxInvoiceFormName       string     `db:"tax_invoice_form_name"`
	ObsType                  string     `db:"obs_type"`
	ObyTypeName              string     `db:"oby_type_name"`
	CreditLimitAction        string     `db:"credit_limit_action"`
	CreditLimitActionName    string     `db:"credit_limit_action_name"`
	SalesInvLimitAction      string     `db:"sales_inv_limit_action"`
	SalesInvLimitActionName  string     `db:"sales_inv_limit_action_name"`
	ObsLimitAction           string     `db:"obs_limit_action"`
	ObsLimitActionName       string     `db:"obs_limit_action_name"`
	MOutletEstablishmentDate string     `db:"outlet_establishment_date"`
	OutletEstablishmentDate  string     `db:"outlet_establishment_date"`
	DelvCity2                string     `db:"delv_city2"`
	DelvLatitude             string     `db:"delv_latitude"`
	DelvLongitude            string     `db:"delv_longitude"`
	DelvLatitude2            string     `db:"delv_latitude2"`
	DelvLongitude2           string     `db:"delv_longitude2"`
	DelvWardId2              string     `db:"delv_ward_id2"`
	DelvZipCode2             string     `db:"delv_zip_code2"`
	BankId                   string     `db:"bank_id"`
	BankCode                 string     `db:"bank_code"`
	BankName                 string     `db:"bank_name"`
	AccountNo                string     `db:"account_no"`
	AccountName              string     `db:"account_name"`
	ContactName              string     `db:"contact_name"`
	JobTitle                 string     `db:"job_title"`
	ContactPhoneNo           string     `db:"contact_phone_no"`
	ContactWaNo              string     `db:"contact_wa_no"`
	ContactEmail             string     `db:"contact_email"`
	IdentityNo               string     `db:"identity_no"`
	ContactIsWaNo            string     `db:"contact_is_wa_no"`
	IdentityType             string     `db:"identity_type"`
	FaxNumber                string     `db:"fax_number"`
	TaxInvoiceId             string     `db:"tax_invoice_id"`
	IsEmbBail2               string     `db:"is_emb_bail2"`
	TaxNo2                   string     `db:"tax_no2"`
	TaxName2                 string     `db:"tax_name2"`
	TaxCity2                 string     `db:"tax_city2"`
	TaxAddr1_2               string     `db:"tax_addr1_2"`
	TaxAddr2_2               string     `db:"tax_addr2_2"`
	TaxType                  string     `db:"tax_type"`
	Nitku                    string     `db:"nitku"`
	AddressTax               string     `db:"address_tax"`
	TaxIdentifierType        string     `db:"tax_identifier_type"`
	TaxIdentifierNo          string     `db:"tax_identifier_no"`
	StatusInsert             string     `db:"status_insert"`
	ErrorMessage             string     `db:"error_message"`
	UpdatedAt                *time.Time `db:"updated_at"`
	DeletedAt                *time.Time `db:"deleted_at"`
}

// ImportOutletUpdateTemp represents a failed update row captured in import.outlet_update_temp
type ImportOutletUpdateTemp struct {
	HistoryID               int64  `db:"history_id"`
	CustID                  string `db:"cust_id"`
	OutletLocCode           string `db:"outlet_loc_code"`
	OutletLocName           string `db:"outlet_loc_name"`
	OutletTypeCode          string `db:"outlet_type_code"`
	OutletTypeName          string `db:"outlet_type_name"`
	OutletGrpCode           string `db:"outlet_grp_code"`
	OutletGrpName           string `db:"outlet_grp_name"`
	DistrictCode            string `db:"district_code"`
	DistrictName            string `db:"district_name"`
	OtClassCode             string `db:"ot_class_code"`
	OtClassName             string `db:"ot_class_name"`
	DiscGrpCode             string `db:"disc_grp_code"`
	DiscGrpName             string `db:"disc_grp_name"`
	MarketCode              string `db:"market_code"`
	MarketName              string `db:"market_name"`
	IndustryCode            string `db:"industry_code"`
	IndustryName            string `db:"industry_name"`
	PriceGrpCode            string `db:"price_grp_code"`
	PriceGrpName            string `db:"price_grp_name"`
	BankCode                string `db:"bank_code"`
	BankName                string `db:"bank_name"`
	ProvinceID              string `db:"province_id"`
	ProvinceName            string `db:"province_name"`
	RegencyID               string `db:"regency_id"`
	RegencyName             string `db:"regency_name"`
	SubDistrictID           string `db:"sub_district_id"`
	SubDistrictName         string `db:"sub_district_name"`
	WardID                  string `db:"ward_id"`
	WardName                string `db:"ward_name"`
	StatusInsert            string `db:"status_insert"`
	ErrorMessage            string `db:"error_message"`
	OutletCode              string `db:"outlet_code"`
	OutletName              string `db:"outlet_name"`
	Address1                string `db:"address1"`
	OutletProvince          string `db:"outlet_province"`
	OutletRegency           string `db:"outlet_regency"`
	OutletSubDistrict       string `db:"outlet_sub_district"`
	OutletWard              string `db:"outlet_ward"`
	ZipCode                 string `db:"zip_code"`
	Longitude               string `db:"longitude"`
	Latitude                string `db:"latitude"`
	BuildingOwn             string `db:"building_own"`
	OutletEstablishmentDate string `db:"outlet_establishment_date"`

	PhoneNo         string `db:"phone_no"`
	FaxNo           string `db:"fax_no"`
	Barcode         string `db:"barcode"`
	ContactName     string `db:"contact_name"`
	JobTitle        string `db:"job_title"`
	IdentityType    string `db:"identity_type"`
	IdentityNo      string `db:"identity_no"`
	ContactPhoneNo  string `db:"contact_phone_no"`
	ContactIsWaNo   string `db:"contact_is_wa_no"`
	ContactWaNo     string `db:"contact_wa_no"`
	ContactEmail    string `db:"contact_email"`
	OutletContactID string `db:"outlet_contact_id"`

	TaxInvoiceFormName string `db:"tax_invoice_form_name"`
	TaxIdentifierType  string `db:"tax_identifier_type"`
	TaxIdentifierNo    string `db:"tax_identifier_no"`
	Nitku              string `db:"nitku"`
	TaxName            string `db:"tax_name"`
	AddressTax         string `db:"address_tax"`
	OutletTaxID        string `db:"outlet_tax_id"`

	IsContraBon string `db:"is_contra_bon"`
	AgentFrom   string `db:"agent_from"`

	DelvAddr1       string `db:"delv_addr1"`
	DelvProvince    string `db:"delv_province"`
	DelvRegency     string `db:"delv_regency"`
	DelvSubDistrict string `db:"delv_sub_district"`
	DelvWard        string `db:"delv_ward"`
	DelvLongitude   string `db:"delv_longitude"`
	DelvLatitude    string `db:"delv_latitude"`
	DelvZipCode     string `db:"delv_zip_code"`
	DelvIsSameAddr  string `db:"delv_is_same_addr"`

	InvAddr1       string `db:"inv_addr1"`
	InvProvince    string `db:"inv_province"`
	InvRegency     string `db:"inv_regency"`
	InvSubDistrict string `db:"inv_sub_district"`
	InvWard        string `db:"inv_ward"`
	InvZipCode     string `db:"inv_zip_code"`
	InvIsSameAddr  string `db:"inv_is_same_addr"`

	PaymentTypeName         string `db:"payment_type_name"`
	ArStatusName            string `db:"ar_status_name"`
	AccountNo               string `db:"account_no"`
	AccountName             string `db:"account_name"`
	OutletBankID            string `db:"outlet_bank_id"`
	Top                     string `db:"top"`
	CreditLimitTypeName     string `db:"credit_limit_type_name"`
	CreditLimit             string `db:"credit_limit"`
	CreditLimitActionName   string `db:"credit_limit_action_name"`
	SalesInvLimitTypeName   string `db:"sales_inv_limit_type_name"`
	SalesInvLimit           string `db:"sales_inv_limit"`
	SalesInvLimitActionName string `db:"sales_inv_limit_action_name"`
	ObsTypeName             string `db:"obs_type_name"`
	Obs                     string `db:"obs"`
	ObsLimitActionName      string `db:"obs_limit_action_name"`
	OutletID                string `db:"outlet_id"`
}

type ImportInstruction struct {
	InstructionID   int64   `db:"instruction_id" json:"instruction_id"`
	InstructionType string  `db:"instruction_type" json:"instruction_type"`
	Kolom           string  `db:"kolom" json:"kolom"`
	Mandatory       bool    `db:"mandatory" json:"mandatory"`
	Keterangan      string  `db:"keterangan" json:"keterangan"`
	Step            *string `db:"step" json:"step"`
}

type OutletQueryFilter struct {
	CustId                        string
	ParentCustId                  string
	Page                          int      `query:"page"`
	Limit                         int      `query:"limit" validate:"required"`
	VerificationStatus            []int    `query:"verification_status"`
	IdentityType                  []string `query:"identity_type"`
	IdentityNo                    []string `query:"identity_no"`
	OutletID                      []int    `query:"outlet_id"`
	OtClassID                     []int    `query:"ot_class_id"`
	OtTypeID                      []int    `query:"ot_type_id"`
	OtGrpID                       []int    `query:"ot_grp_id"`
	Query                         string   `query:"q"`
	Mode                          string   `query:"mode"`
	Sort                          string   `query:"sort"`
	IsActive                      *int     `query:"is_active"`
	IncludeInactive               *int     `query:"include_inactive"`
	Status                        string   `query:"status"`
	OutletStatus                  *int     `query:"outlet_status"` // 0=all, 1=New Open, 2=Covered, 3=Non Active, 4=Closed, 5=Dormant, 6=Registered, 7=Active
	OutletStatusIDs               []int    `query:"-"`
	Format                        string   `query:"format"`
	DistributorID                 []int
	ResolvedCustIdsForDistributor []string
}

type OutletRespone struct {
	OutletId                int64            `json:"outlet_id"`
	OutletCode              string           `json:"outlet_code"`
	OutletName              string           `json:"outlet_name"`
	OutletPrincipalCode     string           `json:"outlet_principal_code"`
	IsPkpOutlet             bool             `json:"is_pkp_outlet"`
	Barcode                 string           `json:"barcode"`
	OutletStatus            int              `json:"status"` // dari mst.m_outlet.outlet_status (1=New Open, 2=Covered, 3=Non Active, 4=Close, 5=Dormant, 6=Registered, 7=Active)
	Address1                string           `json:"address1"`
	Address2                string           `json:"address2"`
	City                    string           `json:"city"`
	ZipCode                 string           `json:"zip_code"`
	PhoneNo                 string           `json:"phone_no"`
	WaNo                    string           `json:"wa_no"`
	FaxNo                   string           `json:"fax_no"`
	Email                   string           `json:"email"`
	DiscGrpId               int              `json:"disc_grp_id"`
	DiscGrpCode             string           `json:"disc_grp_code"`
	DiscGrpName             string           `json:"disc_grp_name"`
	OtLocId                 int              `json:"ot_loc_id"`
	OtGrpId                 int              `json:"ot_grp_id"`
	PriceGrpId              int              `json:"price_grp_id"`
	DistrictId              int              `json:"district_id"`
	BeatId                  int              `json:"beat_id"`
	SbeatId                 int              `json:"sbeat_id"`
	OtClassId               int              `json:"ot_class_id"`
	IndustryId              int              `json:"industry_id"`
	MarketId                int              `json:"market_id"`
	Top                     int              `json:"top"`
	Duedate                 string           `json:"due_date"`
	PaymentType             int              `json:"payment_type"`
	PaymentTypeName         *string          `json:"payment_type_name"`
	IsContraBon             bool             `json:"is_contra_bon"`
	PluGrpId                int              `json:"plu_grp_id"`
	ConvGrpId               int              `json:"conv_grp_id"`
	DiscInvId               int              `json:"disc_inv_id"`
	AgentFrom               string           `json:"agent_from"`
	CreditLimitType         *int             `json:"credit_limit_type"`
	CreditLimitTypeName     *string          `json:"credit_limit_type_name"`
	CreditLimit             float64          `json:"credit_limit"`
	CreditLimitAction       *int             `json:"credit_limit_action"`
	CreditLimitActionName   *string          `json:"credit_limit_action_name"`
	SalesInvLimitType       *int             `json:"sales_inv_limit_type"`
	SalesInvLimitTypeName   *string          `json:"sales_inv_limit_type_name"`
	SalesInvLimit           int              `json:"sales_inv_limit"`
	SalesInvLimitAction     *int             `json:"sales_inv_limit_action"`
	SalesInvLimitActionName *string          `json:"sales_inv_limit_action_name"`
	AvgSalesWeek            float64          `json:"avg_sales_week"`
	AvgSalesMonth           float64          `json:"avg_sales_month"`
	FirstTransDate          string           `json:"first_trans_date"`
	LastTransDate           string           `json:"last_trans_date"`
	FirstWeekNo             int              `json:"first_week_no"`
	OtStartDate             string           `json:"ot_start_date"`
	OtRegDate               string           `json:"ot_reg_date"`
	BuldingOwn              int              `json:"building_own"`
	Dob                     string           `json:"dob"`
	ArStatus                int              `json:"ar_status"`
	ArStatusName            string           `json:"ar_status_name"`
	ArTotal                 float64          `json:"ar_total"`
	CloseDate               string           `json:"closed_date"`
	IsEmbBail               bool             `json:"is_emb_bail"`
	TaxName                 string           `json:"tax_name"`
	TaxAddr1                string           `json:"tax_addr1"`
	TaxAddr2                string           `json:"tax_addr2"`
	TaxCity                 string           `json:"tax_city"`
	TaxNo                   string           `json:"tax_no"`
	TaxInvoiceForm          int              `json:"tax_invoice_form"`
	TaxInvoiceFormName      string           `json:"tax_invoice_form_name"`
	OwnerName               string           `json:"owner_name"`
	OwnerAdd1               string           `json:"owner_addr1"`
	OwnerAddr2              string           `json:"owner_addr2"`
	OwnerCity               string           `json:"owner_city"`
	OwnerPhoneNo            string           `json:"owner_phone_no"`
	OwnerIdNo               string           `json:"owner_id_no"`
	DelvAdd1                string           `json:"delv_addr1"`
	DelvCity                string           `json:"delv_city"`
	DelvLatitude            string           `json:"delv_latitude"`
	DelvLongitude           string           `json:"delv_longitude"`
	DelvAddr2               string           `json:"delv_addr2"`
	DelvCity2               string           `json:"delv_city2"`
	DelvLatitude2           string           `json:"delv_latitude2"`
	DelvLongitude2          string           `json:"delv_longitude2"`
	InvAddr1                string           `json:"inv_addr1"`
	InvAddr2                string           `json:"inv_addr2"`
	InvCity                 string           `json:"inv_city"`
	IsActive                bool             `json:"is_active"`
	UpdatedBy               *int64           `json:"updated_by"`
	UpdatedAt               *time.Time       `json:"updated_at"`
	UpdatedByName           string           `json:"updated_by_name" db:"updated_by_name"`
	Details                 OutletDetailRead `json:"details" validate:"dive"`
	Latitude                string           `json:"latitude"`
	Longitude               string           `json:"longitude"`
	ImageUrl                string           `json:"image_url,omitempty"`
	FileUrl                 *string          `json:"file_url"`
	OtTypeId                int              `json:"ot_type_id"`
	OtTypeName              string           `json:"ot_type_name"`
	IsObs                   bool             `json:"is_obs"`
	Obs                     int              `json:"obs"`
	ObsType                 *int             `json:"obs_type"`
	ObsTypeName             *string          `json:"obs_type_name"`
	ObsLimitAction          *int             `json:"obs_limit_action"`
	ObsLimitActionName      *string          `json:"obs_limit_action_name"`
	OutletLocationName      string           `json:"ot_loc_name"`
	MarketName              string           `json:"market_name"`
	OutletGroupName         string           `json:"ot_grp_name"`
	PriceGroupName          string           `json:"price_grp_name"`
	DistrictName            string           `json:"district_name"`
	IndustryName            string           `json:"industry_name"`
	OutletClassName         string           `json:"ot_class_name"`
	OutletWardId            *string          `json:"outlet_ward_id"`
	OutletWard              *string          `json:"outlet_ward"`
	OutletSubDistrictId     *string          `json:"outlet_sub_district_id"`
	OutletSubDistrict       *string          `json:"outlet_sub_district"`
	OutletRegencyId         *string          `json:"outlet_regency_id"`
	OutletRegency           *string          `json:"outlet_regency"`
	OutletProvinceId        *string          `json:"outlet_province_id"`
	OutletProvince          *string          `json:"outlet_province"`
	IsWaNo                  *bool            `json:"is_wa_no"`
	DelvWardId              *string          `json:"delv_ward_id"`
	DelvWard                *string          `json:"delv_ward"`
	DelvSubDistrictId       *string          `json:"delv_sub_district_id"`
	DelvSubDistrict         *string          `json:"delv_sub_district"`
	DelvRegencyId           *string          `json:"delv_regency_id"`
	DelvRegency             *string          `json:"delv_regency"`
	DelvProvinceId          *string          `json:"delv_province_id"`
	DelvProvince            *string          `json:"delv_province"`
	DelvZipCode             *string          `json:"delv_zip_code"`
	DelvIsSameAddress       *bool            `json:"delv_is_same_addr"`
	DelvWardId2             string           `json:"delv_ward_id2"`
	DelvWard2               *string          `json:"delv_ward2"`
	DelvSubDistrictId2      *string          `json:"delv_sub_district_id2"`
	DelvSubDistrict2        *string          `json:"delv_sub_district2"`
	DelvRegencyId2          *string          `json:"delv_regency_id2"`
	DelvRegency2            *string          `json:"delv_regency2"`
	DelvProvinceId2         *string          `json:"delv_province_id2"`
	DelvProvince2           *string          `json:"delv_province2"`
	DelvZipCode2            *string          `json:"delv_zip_code2"`
	InvWardId               *string          `json:"inv_ward_id"`
	InvWard                 *string          `json:"inv_ward"`
	InvSubDistrictId        *string          `json:"inv_sub_district_id"`
	InvSubDistrict          *string          `json:"inv_sub_district"`
	InvRegencyId            *string          `json:"inv_regency_id"`
	InvRegency              *string          `json:"inv_regency"`
	InvProvinceId           *string          `json:"inv_province_id"`
	InvProvince             *string          `json:"inv_province"`
	InvZipCode              *string          `json:"inv_zip_code"`
	InvIsSameAddress        *bool            `json:"inv_is_same_addr"`
	VerificationStatus      *int             `json:"verification_status"`
	VerificationStatusName  *string          `json:"verification_status_name"`
	VerifiedBy              *int64           `json:"verified_by"`
	VerifiedAt              *time.Time       `json:"verified_at"`
	VerifiedByName          string           `json:"verified_by_name"`
	OutletEstablishmentDate *string          `json:"outlet_establishment_date"`
}

func (o OutletRespone) GenerateTaxInvoiceFormName() string {
	return dataOutletTaxInvoiceFormStatusName[o.TaxInvoiceForm]
}

var dataOutletTaxInvoiceFormStatusName = map[int]string{
	1: "Standart",
	2: "Gabungan",
}

type OutletListRespone struct {
	OutletId                int        `json:"outlet_id"`
	OutletCode              string     `json:"outlet_code"`
	OutletName              string     `json:"outlet_name"`
	Barcode                 string     `json:"barcode"`
	OutletStatus            string     `json:"outlet_status"`
	OutletStatusName        string     `json:"outlet_status_name"`
	PreDormantStatus        *int       `json:"pre_dormant_status"`
	Address1                string     `json:"address1"`
	Address2                string     `json:"address2"`
	City                    string     `json:"city"`
	ZipCode                 string     `json:"zip_code"`
	ContactName             string     `json:"contact_name"`
	ContactPhoneNo          string     `json:"contact_phone_no"`
	PhoneNo                 string     `json:"phone_no"`
	WaNo                    string     `json:"wa_no"`
	FaxNo                   string     `json:"fax_no"`
	Email                   string     `json:"email"`
	DiscGrpId               int        `json:"disc_grp_id"`
	DiscGrpCode             string     `json:"disc_grp_code"`
	DiscGrpName             string     `json:"disc_grp_name"`
	OtLocId                 int        `json:"ot_loc_id"`
	OtLocCode               string     `json:"ot_loc_code"`
	OtLocName               string     `json:"ot_loc_name"`
	OtGrpId                 int        `json:"ot_grp_id"`
	OtGrpCode               string     `json:"ot_grp_code"`
	OtGrpName               string     `json:"ot_grp_name"`
	PriceGrpId              int        `json:"price_grp_id"`
	PriceGrpCode            string     `json:"price_grp_code"`
	PriceGrpName            string     `json:"price_grp_name"`
	DistrictId              int        `json:"district_id"`
	DistrictCode            string     `json:"district_code"`
	DistrictName            string     `json:"district_name"`
	BeatId                  int        `json:"beat_id"`
	BeatCode                string     `json:"beat_code"`
	BeatName                string     `json:"beat_name"`
	SbeatId                 int        `json:"sbeat_id"`
	SBeatCode               string     `json:"sbeat_code"`
	SBeatName               string     `json:"sbeat_name"`
	OtClassId               int        `json:"ot_class_id"`
	OtClassCode             string     `json:"ot_class_code"`
	OtClassName             string     `json:"ot_class_name"`
	IndustryId              int        `json:"industry_id"`
	IndustryCode            string     `json:"industry_code"`
	Industryname            string     `json:"industry_name"`
	MarketId                int        `json:"market_id"`
	MarketCode              string     `json:"market_code"`
	MarketName              string     `json:"market_name"`
	Top                     int        `json:"top"`
	Duedate                 string     `json:"due_date"`
	PaymentType             int        `json:"payment_type"`
	IsContraBon             bool       `json:"is_contra_bon"`
	PluGrpId                int        `json:"plu_grp_id"`
	PluGrpCode              string     `json:"plu_grp_code"`
	PluGrpName              string     `json:"plu_grp_name"`
	ConvGrpId               int        `json:"conv_grp_id"`
	ConvGrpCode             string     `json:"conv_grp_code"`
	ConvGrpName             string     `json:"conv_grp_name"`
	DiscInvId               int        `json:"disc_inv_id"`
	DiscInvCode             string     `json:"disc_inv_code"`
	DiscInvName             string     `json:"disc_inv_name"`
	AgentFrom               string     `json:"agent_from"`
	CreditLimitType         *int       `json:"credit_limit_type"`
	CreditLimit             float64    `json:"credit_limit"`
	SalesInvLimitType       *int       `json:"sales_inv_limit_type"`
	SalesInvLimit           int        `json:"sales_inv_limit"`
	AvgSalesWeek            float64    `json:"avg_sales_week"`
	AvgSalesMonth           float64    `json:"avg_sales_month"`
	FirstTransDate          string     `json:"first_trans_date"`
	LastTransDate           string     `json:"last_trans_date"`
	FirstWeekNo             int        `json:"first_week_no"`
	OtStartDate             string     `json:"ot_start_date"`
	OtRegDate               string     `json:"ot_reg_date"`
	BuldingOwn              int        `json:"building_own"`
	Dob                     string     `json:"dob"`
	ArStatus                int        `json:"ar_status"`
	ArTotal                 float64    `json:"ar_total"`
	CloseDate               string     `json:"closed_date"`
	IsEmbBail               bool       `json:"is_emb_bail"`
	TaxName                 string     `json:"tax_name"`
	TaxAddr1                string     `json:"tax_addr1"`
	TaxAddr2                string     `json:"tax_addr2"`
	TaxCity                 string     `json:"tax_city"`
	TaxNo                   string     `json:"tax_no"`
	TaxInvoiceForm          int        `json:"tax_invoice_form"`
	TaxInvoiceFormName      string     `json:"tax_invoice_form_name"`
	OwnerName               string     `json:"owner_name"`
	OwnerAdd1               string     `json:"owner_addr1"`
	OwnerAddr2              string     `json:"owner_addr2"`
	OwnerCity               string     `json:"owner_city"`
	OwnerPhoneNo            string     `json:"owner_phone_no"`
	OwnerIdNo               string     `json:"owner_id_no"`
	DelvAdd1                string     `json:"delv_addr1"`
	DelvCity                string     `json:"delv_city"`
	DelvLatitude            string     `json:"delv_latitude"`
	DelvLongitude           string     `json:"delv_longitude"`
	DelvAddr2               string     `json:"delv_addr2"`
	DelvCity2               string     `json:"delv_city2"`
	DelvLatitude2           string     `json:"delv_latitude2"`
	DelvLongitude2          string     `json:"delv_longitude2"`
	InvAddr1                string     `json:"inv_addr1"`
	InvAddr2                string     `json:"inv_addr2"`
	InvCity                 string     `json:"inv_city"`
	IsActive                bool       `json:"is_active"`
	CreatedBy               *int64     `json:"created_by"`
	CreatedByName           string     `json:"created_by_name"`
	CreatedAt               *time.Time `json:"created_at"`
	UpdatedBy               *int64     `json:"updated_by"`
	UpdatedAt               *time.Time `json:"updated_at"`
	UpdatedByName           string     `json:"updated_by_name" db:"updated_by_name"`
	Latitude                string     `json:"latitude"`
	Longitude               string     `json:"longitude"`
	ImageUrl                string     `json:"image_url,omitempty"`
	FileUrl                 *string    `json:"file_url"`
	OtTypeId                int        `json:"ot_type_id"`
	OtTypeCode              *string    `json:"ot_type_code"`
	OtTypeName              *string    `json:"ot_type_name"`
	IsObs                   bool       `json:"is_obs"`
	Obs                     int        `json:"obs"`
	OutletWardId            *string    `json:"outlet_ward_id"`
	OutletWard              *string    `json:"outlet_ward"`
	OutletSubDistrictId     *string    `json:"outlet_sub_district_id"`
	OutletSubDistrict       *string    `json:"outlet_sub_district"`
	OutletRegencyId         *string    `json:"outlet_regency_id"`
	OutletRegency           *string    `json:"outlet_regency"`
	OutletProvinceId        *string    `json:"outlet_province_id"`
	OutletProvince          *string    `json:"outlet_province"`
	IsWaNo                  *bool      `json:"is_wa_no"`
	DelvWardId              *string    `json:"delv_ward_id"`
	DelvWard                *string    `json:"delv_ward"`
	DelvSubDistrictId       *string    `json:"delv_sub_district_id"`
	DelvSubDistrict         *string    `json:"delv_sub_district"`
	DelvRegencyId           *string    `json:"delv_regency_id"`
	DelvRegency             *string    `json:"delv_regency"`
	DelvProvinceId          *string    `json:"delv_province_id"`
	DelvProvince            *string    `json:"delv_province"`
	DelvZipCode             *string    `json:"delv_zip_code"`
	DelvIsSameAddress       *bool      `json:"delv_is_same_addr"`
	DelvWardId2             *string    `json:"delv_ward_id2"`
	DelvWard2               *string    `json:"delv_ward2"`
	DelvSubDistrictId2      *string    `json:"delv_sub_district_id2"`
	DelvSubDistrict2        *string    `json:"delv_sub_district2"`
	DelvRegencyId2          *string    `json:"delv_regency_id2"`
	DelvRegency2            *string    `json:"delv_regency2"`
	DelvProvinceId2         *string    `json:"delv_province_id2"`
	DelvProvince2           *string    `json:"delv_province2"`
	DelvZipCode2            *string    `json:"delv_zip_code2"`
	InvWardId               *string    `json:"inv_ward_id"`
	InvWard                 *string    `json:"inv_ward"`
	InvSubDistrictId        *string    `json:"inv_sub_district_id"`
	InvSubDistrict          *string    `json:"inv_sub_district"`
	InvRegencyId            *string    `json:"inv_regency_id"`
	InvRegency              *string    `json:"inv_regency"`
	InvProvinceId           *string    `json:"inv_province_id"`
	InvProvince             *string    `json:"inv_province"`
	InvZipCode              *string    `json:"inv_zip_code"`
	InvIsSameAddress        *bool      `json:"inv_is_same_addr"`
	VerificationStatus      *int       `json:"verification_status"`
	VerificationStatusName  *string    `json:"verification_status_name"`
	ObsType                 *string    `json:"obs_type"`
	CreditLimitAction       *int       `json:"credit_limit_action"`
	SalesInvLimitAction     *int       `json:"sales_inv_limit_action"`
	ObsLimitAction          *int       `json:"obs_limit_action"`
	OutletEstablishmentDate *string    `json:"outlet_establishment_date,omitempty"`
	IdentityType            *string    `json:"identity_type"`
	IdentityNo              *string    `json:"identity_no"`
}

func (o OutletListRespone) GenerateTaxInvoiceFormName() string {
	return dataOutletTaxInvoiceFormStatusName[o.TaxInvoiceForm]
}

type OutletTypeListRespone struct {
	OtTypeId   int     `json:"ot_type_id"`
	OtTypeCode *string `json:"ot_type_code"`
	OtTypeName *string `json:"ot_type_name"`
}

type OutletGroupListRespone struct {
	OtGrpId   int     `json:"ot_grp_id"`
	OtGrpCode *string `json:"ot_grp_code"`
	OtGrpName *string `json:"ot_grp_name"`
}

type ProcessedUpdateOutlet struct {
	// Kunci Utama (Wajib diisi untuk WHERE clause)
	CustId   string `json:"cust_id"`
	OutletId int64  `json:"outlet_id"`

	// Informasi Dasar Outlet
	OutletCode   string `json:"outlet_code,omitempty"`
	Barcode      string `json:"barcode,omitempty"`
	OutletName   string `json:"outlet_name,omitempty"`
	OutletStatus int32  `json:"outlet_status,omitempty"` // int2 di postgres

	// Alamat
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city,omitempty"`
	ZipCode  string `json:"zip_code,omitempty"`

	// Kontak
	PhoneNo string `json:"phone_no,omitempty"`
	WaNo    string `json:"wa_no,omitempty"`
	FaxNo   string `json:"fax_no,omitempty"`
	Email   string `json:"email,omitempty"`

	// Grup & Klasifikasi
	DiscGrpId    int64  `json:"disc_grp_id,omitempty"`
	OtLocId      int64  `json:"ot_loc_id,omitempty"`
	OtGrpId      int64  `json:"ot_grp_id,omitempty"`
	PriceGrpId   int64  `json:"price_grp_id,omitempty"`
	DistrictId   int64  `json:"district_id,omitempty"`
	BeatId       int64  `json:"beat_id,omitempty"`
	SbeatId      int64  `json:"sbeat_id,omitempty"`
	OtClassId    int64  `json:"ot_class_id,omitempty"`
	IndustryId   int64  `json:"industry_id,omitempty"`
	MarketId     int64  `json:"market_id,omitempty"`
	OtTypeId     int64  `json:"ot_type_id,omitempty"`
	OutletWardId string `json:"outlet_ward_id,omitempty"`

	// Informasi Penjualan & Pembayaran
	Top                 int32   `json:"top,omitempty"`          // int4 di postgres
	PaymentType         int32   `json:"payment_type,omitempty"` // int2
	IsContraBon         bool    `json:"is_contra_bon,omitempty"`
	CreditLimitType     int32   `json:"credit_limit_type,omitempty"`    // int2
	CreditLimit         float64 `json:"credit_limit,omitempty"`         // numeric
	SalesInvLimitType   int32   `json:"sales_inv_limit_type,omitempty"` // int2
	SalesInvLimit       int32   `json:"sales_inv_limit,omitempty"`      // int2
	AvgSalesWeek        float64 `json:"avg_sales_week,omitempty"`       // numeric
	AvgSalesMonth       float64 `json:"avg_sales_month,omitempty"`      // numeric
	CreditLimitAction   int64   `json:"credit_limit_action,omitempty"`
	SalesInvLimitAction int64   `json:"sales_inv_limit_action,omitempty"`

	// PLU, Konversi, & Diskon Invoice
	PluGrpId  int64 `json:"plu_grp_id,omitempty"`
	ConvGrpId int64 `json:"conv_grp_id,omitempty"`
	DiscInvId int64 `json:"disc_inv_id,omitempty"`

	// Informasi Agen & Transaksi
	AgentFrom      string    `json:"agent_from,omitempty"`
	FirstTransDate time.Time `json:"first_trans_date,omitempty"`
	LastTransDate  time.Time `json:"last_trans_date,omitempty"`
	FirstWeekNo    int32     `json:"first_week_no,omitempty"` // int2

	// Informasi Registrasi & Kepemilikan
	OtStartDate             time.Time `json:"ot_start_date,omitempty"`
	OtRegDate               time.Time `json:"ot_reg_date,omitempty"`
	BuildingOwn             int32     `json:"building_own,omitempty"` // int2
	OutletEstablishmentDate time.Time `json:"outlet_establishment_date,omitempty"`

	// Informasi Piutang (AR)
	ArStatus int32   `json:"ar_status,omitempty"` // int2
	ArTotal  float64 `json:"ar_total,omitempty"`  // numeric

	// Status & Tanggal Lainnya
	ClosedDate  time.Time `json:"closed_date,omitempty"`
	IsEmbBail   bool      `json:"is_emb_bail,omitempty"`
	IsPkpOutlet *bool     `json:"is_pkp_outlet,omitempty"`
	IsActive    *bool     `json:"is_active,omitempty"`

	// Informasi Pajak
	TaxName        string `json:"tax_name,omitempty"`
	TaxAddr1       string `json:"tax_addr1,omitempty"`
	TaxAddr2       string `json:"tax_addr2,omitempty"`
	TaxCity        string `json:"tax_city,omitempty"`
	TaxNo          string `json:"tax_no,omitempty"`
	TaxInvoiceForm int32  `json:"tax_invoice_form,omitempty"` // int2

	// Informasi Pemilik
	OwnerName    string    `json:"owner_name,omitempty"`
	OwnerAddr1   string    `json:"owner_addr1,omitempty"`
	OwnerAddr2   string    `json:"owner_addr2,omitempty"`
	OwnerCity    string    `json:"owner_city,omitempty"`
	OwnerPhoneNo string    `json:"owner_phone_no,omitempty"`
	OwnerIdNo    string    `json:"owner_id_no,omitempty"`
	Dob          time.Time `json:"dob,omitempty"` // Date of Birth

	// Alamat Pengiriman (Delivery)
	DelvAddr1      string `json:"delv_addr1,omitempty"`
	DelvAddr2      string `json:"delv_addr2,omitempty"`
	DelvCity       string `json:"delv_city,omitempty"`
	DelvCity2      string `json:"delv_city2,omitempty"`
	DelvWardId     string `json:"delv_ward_id,omitempty"`
	DelvWardId2    string `json:"delv_ward_id2,omitempty"`
	DelvZipCode    string `json:"delv_zip_code,omitempty"`
	DelvZipCode2   string `json:"delv_zip_code2,omitempty"`
	DelvIsSameAddr bool   `json:"delv_is_same_addr,omitempty"`
	DelvLatitude   string `json:"delv_latitude,omitempty"`
	DelvLongitude  string `json:"delv_longitude,omitempty"`
	DelvLatitude2  string `json:"delv_latitude2,omitempty"`
	DelvLongitude2 string `json:"delv_longitude2,omitempty"`

	// Alamat Penagihan (Invoice)
	InvAddr1      string `json:"inv_addr1,omitempty"`
	InvAddr2      string `json:"inv_addr2,omitempty"`
	InvCity       string `json:"inv_city,omitempty"`
	InvWardId     string `json:"inv_ward_id,omitempty"`
	InvZipCode    string `json:"inv_zip_code,omitempty"`
	InvIsSameAddr bool   `json:"inv_is_same_addr,omitempty"`

	// Geolocation & Gambar
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`

	// Observasi (OBS)
	IsObs          bool  `json:"is_obs,omitempty"`
	Obs            int64 `json:"obs,omitempty"`
	ObsType        int64 `json:"obs_type,omitempty"`
	ObsLimitAction int64 `json:"obs_limit_action,omitempty"`

	// Lain-lain
	IsWaNo             bool      `json:"is_wa_no,omitempty"`
	VerificationStatus int32     `json:"verification_status,omitempty"` // int2
	VerifiedAt         time.Time `json:"verified_at,omitempty"`
	VerifiedBy         int64     `json:"verified_by,omitempty"`

	// Meta (Wajib diisi)
	UpdatedBy int64 `json:"updated_by"`
}

type CreateOutletBody struct {
	ParentCustId            string
	CustId                  string       `json:"cust_id" validate:"required,max=10"`
	OutletCode              string       `json:"outlet_code" validate:"omitempty,max=10,alphanumericSpace"`
	OutletName              string       `json:"outlet_name" validate:"required,max=75"`
	Barcode                 string       `json:"barcode" validate:"max=50,alphanumericSpace"`
	OutletStatus            int          `json:"outlet_status" validate:"max=5"`
	Status                  int          `json:"status"`
	Address1                string       `json:"address1" validate:"max=150"`
	Address2                string       `json:"address2" validate:"max=150"`
	City                    string       `json:"city" validate:"max=100"`
	ZipCode                 string       `json:"zip_code" validate:"omitempty,max=6,number"`
	PhoneNo                 string       `json:"phone_no" validate:"omitempty,max=20,number"`
	WaNo                    string       `json:"wa_no" validate:"max=20"`
	FaxNo                   string       `json:"fax_no" validate:"max=20"`
	Email                   string       `json:"email" validate:"omitempty,email,max=20"`
	DiscGrpId               int          `json:"disc_grp_id"`
	OtLocId                 int          `json:"ot_loc_id"`
	OtGrpId                 int          `json:"ot_grp_id"`
	PriceGrpId              int          `json:"price_grp_id"`
	DistrictId              int          `json:"district_id"`
	BeatId                  int          `json:"beat_id"`
	SbeatId                 int          `json:"sbeat_id"`
	OtClassId               int          `json:"ot_class_id"`
	IndustryId              int          `json:"industry_id"`
	MarketId                int          `json:"market_id"`
	Top                     int          `json:"top"`
	PaymentType             int          `json:"payment_type" validate:"max=5"`
	IsContraBon             bool         `json:"is_contra_bon"`
	PluGrpId                int          `json:"plu_grp_id"`
	ConvGrpId               int          `json:"conv_grp_id"`
	DiscInvId               int          `json:"disc_inv_id"`
	AgentFrom               string       `json:"agent_from" validate:"max=50"`
	CreditLimitType         *int         `json:"credit_limit_type"`
	CreditLimitTypeName     string       `json:"credit_limit_type_name"`
	CreditLimit             float64      `json:"credit_limit"`
	CreditLimitAction       *int         `json:"credit_limit_action"`
	CreditLimitActionName   string       `json:"credit_limit_action_name"`
	SalesInvLimitType       *int         `json:"sales_inv_limit_type"`
	SalesInvLimitTypeName   string       `json:"sales_inv_limit_type_name"`
	SalesInvLimit           int          `json:"sales_inv_limit"`
	SalesInvLimitAction     *int         `json:"sales_inv_limit_action"`
	SalesInvLimitActionName string       `json:"sales_inv_limit_action_name"`
	AvgSalesWeek            float64      `json:"avg_sales_week"`
	AvgSalesMonth           float64      `json:"avg_sales_month"`
	FirstTransDate          *string      `json:"first_trans_date,omitempty"`
	LastTransDate           *string      `json:"last_trans_date,omitempty"`
	FirstWeekNo             int          `json:"first_week_no"`
	OtStartDate             *string      `json:"ot_start_date,omitempty"`
	OtRegDate               *string      `json:"ot_reg_date,omitempty"`
	BuldingOwn              int          `json:"building_own"`
	Dob                     *string      `json:"dob,omitempty"`
	ArStatus                int          `json:"ar_status"`
	ArTotal                 float64      `json:"ar_total"`
	CloseDate               *string      `json:"closed_date,omitempty"`
	IsEmbBail               bool         `json:"is_emb_bail"`
	TaxName                 string       `json:"tax_name"`
	TaxAddr1                string       `json:"tax_addr1"`
	TaxAddr2                string       `json:"tax_addr2"`
	TaxCity                 string       `json:"tax_city"`
	TaxNo                   string       `json:"tax_no"`
	TaxInvoiceForm          int          `json:"tax_invoice_form" validate:"omitempty,oneof=1 2"`
	OwnerName               string       `json:"owner_name"`
	OwnerAdd1               string       `json:"owner_addr1"`
	OwnerAddr2              string       `json:"owner_addr2"`
	OwnerCity               string       `json:"owner_city"`
	OwnerPhoneNo            string       `json:"owner_phone_no"`
	OwnerIdNo               string       `json:"owner_id_no"`
	DelvAdd1                string       `json:"delv_addr1"`
	DelvCity                string       `json:"delv_city"`
	DelvLatitude            string       `json:"delv_latitude"`
	DelvLongitude           string       `json:"delv_longitude"`
	DelvAddr2               string       `json:"delv_addr2"`
	DelvCity2               string       `json:"delv_city2"`
	DelvLatitude2           string       `json:"delv_latitude2"`
	DelvLongitude2          string       `json:"delv_longitude2"`
	InvAddr1                string       `json:"inv_addr1"`
	InvAddr2                string       `json:"inv_addr2"`
	InvCity                 FlexString   `json:"inv_city"`
	IsActive                bool         `json:"is_active"`
	CreatedBy               int64        `json:"created_by" validate:"required"`
	UpdatedBy               int64        `json:"updated_by"`
	Details                 OutletDetail `json:"details"`
	Latitude                string       `json:"latitude"`
	Longitude               string       `json:"longitude"`
	ImageUrl                string       `json:"image_url,omitempty"`
	OtTypeId                int          `json:"ot_type_id"`
	IsObs                   bool         `json:"is_obs"`
	ObsType                 *int         `json:"obs_type"`
	Obs                     int          `json:"obs"`
	ObsLimitAction          *int         `json:"obs_limit_action"`
	OutletWardId            *FlexString  `json:"outlet_ward_id"`
	OutletProvinceId        *FlexString  `json:"outlet_province_id,omitempty"`
	OutletProvince          string       `json:"outlet_province,omitempty"`
	OutletRegencyId         *FlexString  `json:"outlet_regency_id,omitempty"`
	OutletRegency           string       `json:"outlet_regency,omitempty"`
	OutletSubDistrictId     *FlexString  `json:"outlet_sub_district_id,omitempty"`
	OutletSubDistrict       string       `json:"outlet_sub_district,omitempty"`
	OutletWard              string       `json:"outlet_ward,omitempty"`
	IsWaNo                  *bool        `json:"is_wa_no"`
	DelvWardId              *FlexString  `json:"delv_ward_id"`
	DelvZipCode             *FlexString  `json:"delv_zip_code"`
	DelvIsSameAddress       *bool        `json:"delv_is_same_addr"`
	DelvWardId2             *FlexString  `json:"delv_ward_id2"`
	DelvZipCode2            *FlexString  `json:"delv_zip_code2"`
	InvWardId               *FlexString  `json:"inv_ward_id"`
	InvZipCode              *FlexString  `json:"inv_zip_code"`
	InvIsSameAddress        *bool        `json:"inv_is_same_addr"`
	VerificationStatus      *int         `json:"verification_status"`
	OutletEstablishmentDate *string      `json:"outlet_establishment_date,omitempty"`
	IsPkpOutlet             bool         `json:"is_pkp_outlet"`
	CreatedByName           string       `json:"-"`
}

type CreateOutletBodyV2 struct {
	ParentCustId string
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	OutletName   string  `json:"outlet_name" validate:"required,max=150"`
	OutletCode   string  `json:"outlet_code" validate:"required,max=30,alphanum"`
	Barcode      string  `json:"barcode" validate:"max=50,alphanum"`
	OutletStatus int     `json:"outlet_status" validate:"max=5"`
	Address1     string  `json:"address1" validate:"max=150"`
	Address2     string  `json:"address2" validate:"max=150"`
	City         string  `json:"city" validate:"max=100"`
	ZipCode      string  `json:"zip_code" validate:"max=6,number"`
	OtLocId      int     `json:"ot_loc_id"`
	Latitude     string  `json:"latitude"`
	Longitude    string  `json:"longitude"`
	Email        string  `json:"email" validate:"email,max=20"`
	PhoneNo      string  `json:"phone_no" validate:"max=20,number"`
	WaNo         string  `json:"wa_no" validate:"max=20"`
	FaxNo        string  `json:"fax_no" validate:"max=20"`
	MarketId     int     `json:"market_id"`
	BuldingOwn   int     `json:"building_own"`
	OtRegDate    *string `json:"ot_reg_date,omitempty"`
	OtStartDate  *string `json:"ot_start_date,omitempty"`
	Dob          *string `json:"dob,omitempty"`
	CloseDate    *string `json:"closed_date,omitempty"`

	// TAX
	// TaxInvId int  `json:"tax_invoice_id"`
	IsEmbBail bool   `json:"is_emb_bail"`
	TaxNo     string `json:"tax_no"`
	TaxName   string `json:"tax_name"`
	TaxCity   string `json:"tax_city"`
	TaxAddr1  string `json:"tax_addr1"`
	TaxAddr2  string `json:"tax_addr2"`

	DiscGrpId   int  `json:"disc_grp_id"`
	OtGrpId     int  `json:"ot_grp_id"`
	IsContraBon bool `json:"is_contra_bon"`
	PriceGrpId  int  `json:"price_grp_id"`
	DistrictId  int  `json:"district_id"`
	IndustryId  int  `json:"industry_id"`
	OtClassId   int  `json:"ot_class_id"`
	// OtTypeId int  `json:"ot_type_id"`
	AgentFrom string `json:"agent_from" validate:"max=50"`

	DelvAdd1  string `json:"delv_addr1"`
	DelvAddr2 string `json:"delv_addr2"`
	DelvCity  string `json:"delv_city"`
	InvAddr1  string `json:"inv_addr1"`
	InvAddr2  string `json:"inv_addr2"`
	InvCity   string `json:"inv_city"`

	PaymentType       int     `json:"payment_type" validate:"max=5"`
	CreditLimitType   int     `json:"credit_limit_type" validate:"max=5"`
	CreditLimit       float64 `json:"credit_limit"`
	SalesInvLimitType int     `json:"sales_inv_limit_type" validate:"max=5"`
	SalesInvLimit     int     `json:"sales_inv_limit"`
	AvgSalesWeek      float64 `json:"avg_sales_week"`
	AvgSalesMonth     float64 `json:"avg_sales_month"`
	ArStatus          int     `json:"ar_status"`
	ArTotal           float64 `json:"ar_total"`
	Top               int     `json:"top"`
	// IsObs             bool    `json:"is_obs"`
	// Obs               int  `json:"obs"`

	BeatId         int     `json:"beat_id"`
	SbeatId        int     `json:"sbeat_id"`
	PluGrpId       int     `json:"plu_grp_id"`
	ConvGrpId      int     `json:"conv_grp_id"`
	DiscInvId      int     `json:"disc_inv_id"`
	FirstTransDate *string `json:"first_trans_date,omitempty"`
	LastTransDate  *string `json:"last_trans_date,omitempty"`
	FirstWeekNo    int     `json:"first_week_no"`
	TaxInvoiceForm int     `json:"tax_invoice_form"  validate:"omitempty,oneof=1 2"`
	OwnerName      string  `json:"owner_name"`
	OwnerAdd1      string  `json:"owner_addr1"`
	OwnerAddr2     string  `json:"owner_addr2"`
	OwnerCity      string  `json:"owner_city"`
	OwnerPhoneNo   string  `json:"owner_phone_no"`
	OwnerIdNo      string  `json:"owner_id_no"`

	IsActive  bool         `json:"is_active"`
	CreatedBy int64        `json:"created_by" validate:"required"`
	UpdatedBy int64        `json:"updated_by"`
	Details   OutletDetail `json:"details"`
	ImageUrl  string       `json:"image_url,omitempty"`
}

type DetailOutletParams struct {
	OutletId int64 `params:"outlet_id" validate:"required"`
}

type UpdateOutletParams struct {
	OutletId int `params:"outlet_id" validate:"required"`
}

// UpdateOutletStatusRequest body for PATCH /v1/outlets/update-status/:outlet_id (e.g. Close Outlet: status=4)
type UpdateOutletStatusRequest struct {
	Status *int `json:"status"` // optional; 4 = closed
}

type DeleteOutletParams struct {
	OutletId int `params:"outlet_id" validate:"required"`
}

type UpdateOutletRequest struct {
	ParentCustId            string
	CustId                  string       `json:"cust_id" validate:"required,max=10"`
	OutletCode              string       `json:"outlet_code,omitempty" validate:"required,max=30,omitempty"`
	OutletName              string       `json:"outlet_name,omitempty" validate:"max=50,omitempty"`
	Barcode                 string       `json:"barcode"`
	OutletStatus            int          `json:"outlet_status"`
	Status                  *int         `json:"status,omitempty"`
	Address1                string       `json:"address1"`
	Address2                string       `json:"address2"`
	OutletProvinceId        *string      `json:"outlet_province_id,omitempty"`
	OutletRegencyId         *string      `json:"outlet_regency_id,omitempty"`
	OutletSubDistrictId     *string      `json:"outlet_sub_district_id,omitempty"`
	OutletWardId            *string      `json:"outlet_ward_id,omitempty"`
	OutletWard              *string      `json:"outlet_ward,omitempty"`
	City                    string       `json:"city"`
	ZipCode                 string       `json:"zip_code"`
	PhoneNo                 string       `json:"phone_no"`
	WaNo                    string       `json:"wa_no"`
	FaxNo                   string       `json:"fax_no"`
	Email                   string       `json:"email"`
	ContactName             *string      `json:"contact_name,omitempty"`
	Positions               *string      `json:"positions,omitempty"`
	IsPkpOutlet             *bool        `json:"is_pkp_outlet,omitempty"`
	IdentityType            *string      `json:"identity_type,omitempty"`
	DiscGrpId               int          `json:"disc_grp_id"`
	OtLocId                 int          `json:"ot_loc_id"`
	OtGrpId                 int          `json:"ot_grp_id"`
	PriceGrpId              int          `json:"price_grp_id"`
	DistrictId              int          `json:"district_id"`
	BeatId                  int          `json:"beat_id"`
	SbeatId                 int          `json:"sbeat_id"`
	OtClassId               int          `json:"ot_class_id"`
	IndustryId              int          `json:"industry_id"`
	MarketId                int          `json:"market_id"`
	Top                     int          `json:"top"`
	PaymentType             int          `json:"payment_type"`
	IsContraBon             bool         `json:"is_contra_bon"`
	PluGrpId                int          `json:"plu_grp_id"`
	ConvGrpId               int          `json:"conv_grp_id"`
	DiscInvId               int          `json:"disc_inv_id"`
	AgentFrom               string       `json:"agent_from"`
	CreditLimitType         *int         `json:"credit_limit_type"`
	CreditLimit             float64      `json:"credit_limit"`
	CreditLimitAction       *int         `json:"credit_limit_action"`
	SalesInvLimitType       *int         `json:"sales_inv_limit_type"`
	SalesInvLimit           int          `json:"sales_inv_limit"`
	SalesInvLimitAction     *int         `json:"sales_inv_limit_action"`
	AvgSalesWeek            float64      `json:"avg_sales_week"`
	AvgSalesMonth           float64      `json:"avg_sales_month"`
	FirstTransDate          *string      `json:"first_trans_date"`
	LastTransDate           *string      `json:"last_trans_date"`
	FirstWeekNo             int          `json:"first_week_no"`
	OtStartDate             *string      `json:"ot_start_date"`
	OtRegDate               *string      `json:"ot_reg_date"`
	BuldingOwn              int          `json:"building_own"`
	Dob                     *string      `json:"dob"`
	ArStatus                int          `json:"ar_status"`
	ArTotal                 float64      `json:"ar_total"`
	CloseDate               *string      `json:"closed_date"`
	IsEmbBail               bool         `json:"is_emb_bail"`
	TaxName                 string       `json:"tax_name"`
	TaxAddr1                string       `json:"tax_addr1"`
	TaxAddr2                string       `json:"tax_addr2"`
	TaxCity                 string       `json:"tax_city"`
	TaxNo                   string       `json:"tax_no"`
	Nitku                   *string      `json:"nitku,omitempty"`
	TaxType                 *string      `json:"tax_type,omitempty"`
	TaxInvoiceForm          int          `json:"tax_invoice_form"  validate:"omitempty,oneof=1 2"`
	OwnerName               string       `json:"owner_name"`
	OwnerAdd1               string       `json:"owner_addr1"`
	OwnerAddr2              string       `json:"owner_addr2"`
	OwnerCity               string       `json:"owner_city"`
	OwnerPhoneNo            string       `json:"owner_phone_no"`
	OwnerIdNo               string       `json:"owner_id_no"`
	DelvAdd1                string       `json:"delv_addr1"`
	DelvCity                string       `json:"delv_city"`
	DelvLatitude            string       `json:"delv_latitude"`
	DelvLongitude           string       `json:"delv_longitude"`
	DelvAddr2               string       `json:"delv_addr2"`
	DelvCity2               string       `json:"delv_city2"`
	DelvLatitude2           string       `json:"delv_latitude2"`
	DelvLongitude2          string       `json:"delv_longitude2"`
	InvAddr1                string       `json:"inv_addr1"`
	InvAddr2                string       `json:"inv_addr2"`
	InvCity                 string       `json:"inv_city"`
	IsActive                *bool        `json:"is_active,omitempty"`
	UpdatedBy               int64        `json:"updated_by" validate:"required"`
	Details                 OutletDetail `json:"details" validate:""`
	Latitude                string       `json:"latitude"`
	Longitude               string       `json:"longitude"`
	ImageUrl                string       `json:"image_url,omitempty"`
	OtTypeId                int          `json:"ot_type_id"`
	IsObs                   bool         `json:"is_obs"`
	ObsType                 *int         `json:"obs_type"`
	Obs                     int          `json:"obs"`
	ObsLimitAction          *int         `json:"obs_limit_action"`
	IsWaNo                  *bool        `json:"is_wa_no"`
	DelvWardId              *string      `json:"delv_ward_id"`
	DelvZipCode             *string      `json:"delv_zip_code"`
	DelvIsSameAddress       *bool        `json:"delv_is_same_addr"`
	DelvProvinceId1         *string      `json:"delv_province_id1,omitempty"`
	DelvRegencyId1          *string      `json:"delv_regency_id1,omitempty"`
	DelvSubDistrictId1      *string      `json:"delv_sub_district_id1,omitempty"`
	DelvWardId2             *string      `json:"delv_ward_id2"`
	DelvZipCode2            *string      `json:"delv_zip_code2"`
	DelvProvinceId2         *string      `json:"delv_province_id2,omitempty"`
	DelvRegencyId2          *string      `json:"delv_regency_id2,omitempty"`
	DelvSubDistrictId2      *string      `json:"delv_sub_district_id2,omitempty"`
	InvWardId               *string      `json:"inv_ward_id"`
	InvZipCode              *string      `json:"inv_zip_code"`
	InvIsSameAddress        *bool        `json:"inv_is_same_addr"`
	InvProvinceId           *string      `json:"inv_province_id,omitempty"`
	InvRegencyId            *string      `json:"inv_regency_id,omitempty"`
	InvSubDistrictId        *string      `json:"inv_sub_district_id,omitempty"`
	VerificationStatus      *int         `json:"verification_status"`
	OutletEstablishmentDate string       `json:"outlet_establishment_date"`
}

type OutletSalesman struct {
	OutletSalesId *int64 `json:"outlet_sales_id"`
	SalesID       *int   `json:"sales_id" validate:"gt=0"`
	W1            *bool  `json:"w1"`
	W2            *bool  `json:"w2"`
	W3            *bool  `json:"w3"`
	W4            *bool  `json:"w4"`
	RouteID       *int   `json:"route_id" validate:"gt=0"`
	DayID         *int   `json:"day_id" validate:"gt=0"`
}

type OutletBank struct {
	OutletBankId *int64 `json:"outlet_bank_id"`
	BankId       *int64 `json:"bank_id"`
	AccountNo    string `json:"account_no" validate:"max=50"`
	AccountName  string `json:"account_name" validate:"max=150"`
}

type OutletContact struct {
	OutletContactId *int64  `json:"outlet_contact_id"`
	ContactName     *string `json:"contact_name" validate:"max=50"`
	JobTitle        string  `json:"job_title" validate:"max=20"`
	IdentityNo      string  `json:"identity_no" validate:"max=100"`
	IsWaNo          bool    `json:"is_wa_no" `
	PhoneNo         string  `json:"phone_no" validate:"max=20"`
	WaNo            string  `json:"wa_no" validate:"max=20"`
	Email           string  `json:"email" validate:""`
	IdentityType    string  `json:"identity_type"`
	FaxNumber       string  `json:"fax_number"`
	// OutletEstablishmentDate *string `json:"outlet_establishment_date"`
}

type OutletTax struct {
	OutletTaxId       *int64 `json:"outlet_tax_id"`
	TaxInvId          int    `json:"tax_invoice_id"`
	IsEmbBail         bool   `json:"is_emb_bail"`
	TaxNo             string `json:"tax_no"`
	TaxName           string `json:"tax_name"`
	TaxCity           string `json:"tax_city"`
	TaxAddr1          string `json:"tax_addr1"`
	TaxAddr2          string `json:"tax_addr2"`
	TaxType           string `json:"tax_type"`
	Nitku             string `json:"nitku"`
	AddressTax        string `json:"address_tax"`
	TaxIdentifierType string `json:"tax_identifier_type"`
	TaxIdentifierNo   string `json:"tax_identifier_no"`
}

type OutletDetail struct {
	OutletSalesman []OutletSalesman `json:"salesman" validate:"dive"`
	OutletBank     []OutletBank     `json:"bank" validate:"dive"`
	OutletContact  []OutletContact  `json:"contact" validate:"dive"`
	OutletTax      []OutletTax      `json:"tax" validate:"dive"`
}

type OutletDetailRead struct {
	OutletSalesman []OutletSalesmanRead `json:"salesman"  validate:"dive"`
	OutletBank     []OutletBankRead     `json:"bank" validate:"dive"`
	OutletContact  []OutletContactRead  `json:"contact"  validate:"dive"`
	OutletTax      []OutletTaxRead      `json:"tax"  validate:"dive"`
}

type OutletBankRead struct {
	OutletBankId *int64  `json:"outlet_bank_id"`
	OutletId     *int64  `json:"outlet_id"`
	BankId       *int64  `json:"bank_id"`
	BankCode     *string `json:"bank_code"`
	BankName     *string `json:"bank_name"`
	AccountNo    string  `json:"account_no" validate:"required,max=50"`
	AccountName  string  `json:"account_name" validate:"required,max=150"`
}

type OutletSalesmanRead struct {
	OutletSalesId *int64  `json:"outlet_sales_id"`
	OutletId      *int64  `json:"outlet_id"`
	SalesID       *int    `json:"sales_id" validate:"gt=0"`
	SalesName     *string `json:"sales_name"`
	W1            *bool   `json:"w1"`
	W2            *bool   `json:"w2"`
	W3            *bool   `json:"w3"`
	W4            *bool   `json:"w4"`
	RouteID       *int    `json:"route_id" validate:"gt=0"`
	DayID         *int    `json:"day_id" validate:"gt=0"`
	DayName       *string `json:"day_name" validate:"gt=0"`
}
type OutletContactRead struct {
	OutletContactId         *int64  `json:"outlet_contact_id"`
	OutletId                *int64  `json:"outlet_id"`
	ContactName             *string `json:"contact_name" validate:"required,max=150"`
	JobTitle                string  `json:"job_title" validate:"required,max=100"`
	PhoneNo                 string  `json:"phone_no" validate:"max=20"`
	WaNo                    string  `json:"wa_no" validate:"max=20"`
	Email                   string  `json:"email" validate:"email,max=100"`
	IdentityNo              string  `json:"identity_no" validate:"max=100"`
	IsWaNo                  *bool   `json:"is_wa_no" validate:"max=100"`
	IdentityType            string  `json:"identity_type"`
	FaxNumber               string  `json:"fax_number"`
	OutletEstablishmentDate string  `json:"outlet_establishment_date"`
}
type OutletTaxRead struct {
	OutletTaxId       *int64 `json:"outlet_tax_id"`
	OutletId          *int64 `json:"outlet_id"`
	TaxInvId          int    `json:"tax_invoice_id"`
	IsEmbBail         bool   `json:"is_emb_bail"`
	TaxNo             string `json:"tax_no"`
	TaxName           string `json:"tax_name"`
	TaxCity           string `json:"tax_city"`
	TaxAddr1          string `json:"tax_addr1"`
	TaxAddr2          string `json:"tax_addr2"`
	TaxType           string `json:"tax_type"`
	Nitku             string `json:"nitku"`
	AdddressTax       string `json:"address_tax"`
	TaxIdentifierType string `json:"tax_identifier_type"`
	TaxIdentifierNo   string `json:"tax_identifier_no"`
}

var dataOutletVerificationStatusName = map[int]string{
	1: "Approved",
	2: "In Review",
	3: "Rejected",
}

func (outlet OutletListRespone) GenerateOutletVerificationStatusName() string {
	if outlet.VerificationStatus != nil {
		return dataOutletVerificationStatusName[*outlet.VerificationStatus]
	}
	return ""
}

func (outlet OutletRespone) GenerateOutletVerificationStatusName() string {
	if outlet.VerificationStatus != nil {
		return dataOutletVerificationStatusName[*outlet.VerificationStatus]
	}
	return ""
}

type ApproveOutletBody struct {
	CustId  string  `json:"cust_id" validate:"required,max=10"`
	Outlets []int64 `json:"outlets"`
	// UpdatedBy          int64   `json:"updated_by" validate:"required"`
	VerifiedBy         int64     `json:"verified_by"`
	VerifiedAt         time.Time `json:"updated_at"`
	VerificationStatus int       `json:"verification_status"`
}

type RejectOutletBody struct {
	CustId  string  `json:"cust_id" validate:"required,max=10"`
	Outlets []int64 `json:"outlets"`
	// UpdatedBy          int64   `json:"updated_by" validate:"required"`
	VerifiedBy         int64     `json:"verified_by"`
	VerifiedAt         time.Time `json:"updated_at"`
	VerificationStatus int       `json:"verification_status"`
}

type VerificationStatusListRespone struct {
	VerificationStatus     *int    `json:"verification_status"`
	VerificationStatusName *string `json:"verification_status_name"`
}

func (outlet VerificationStatusListRespone) GenerateOutletVerificationStatusName() string {
	if outlet.VerificationStatus != nil {
		return dataOutletVerificationStatusName[*outlet.VerificationStatus]
	}
	return ""
}

// ImportRequest represents a request to import data from a file
type ImportRequest struct {
	File          multipart.File
	CustId        string
	UserId        int64
	Filename      string
	Format        string
	ParentCustId  string
	CreatedByName string
	IsImportNew   bool
}

// ImportOutletRow represents a row from an imported outlet file
type ImportOutletRow struct {
	OutletCode     string
	OutletName     string
	Address        string
	Phone          string
	ContactPerson  string
	OutletCatCode  string
	OutletCatName  string
	AreaCode       string
	AreaName       string
	RegionCode     string
	RegionName     string
	RayonCode      string
	RayonName      string
	PriceGroupCode string
	PriceGroupName string
	SalesmanCode   string
	SalesmanName   string
	IsActive       string
	Latitude       string
	Longitude      string
	CreditLimit    string
}

type OutletListByDistributorRespone struct {
	OutletId     int    `json:"outlet_id"`
	OutletCode   string `json:"outlet_code"`
	OutletName   string `json:"outlet_name"`
	OutletStatus int    `json:"outlet_status"`
	Address1     string `json:"address1"`
	City         string `json:"city"`
	ZipCode      string `json:"zip_code"`
	PhoneNo      string `json:"phone_no"`
	WaNo         string `json:"wa_no"`
	FaxNo        string `json:"fax_no"`
	Email        string `json:"email"`
}

// OutletListApprovalQueryFilter represents query parameters for filtering outlet approval list
type OutletListApprovalQueryFilter struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Status string `query:"status"`
	Q      string `query:"q"`
}

// OutletListApprovalResponse represents the response structure for outlet approval list
type OutletListApprovalResponse struct {
	OutletCrId  int64   `json:"outlet_cr_id"`
	OutletId    int64   `json:"outlet_id"`
	OutletCode  string  `json:"outlet_code"`
	OutletName  string  `json:"outlet_name"`
	CurrentLong *string `json:"current_long"`
	CurrentLat  *string `json:"current_lat"`
	NewLong     *string `json:"new_long"`
	NewLat      *string `json:"new_lat"`
	Source      int     `json:"source"`
	Status      int     `json:"status"`
	StatusDesc  string  `json:"status_desc"`
	RequestBy   *string `json:"request_by"`
	RequestDate string  `json:"request_date"`
}

// OutletListApprovalRequest represents the request body for approving/rejecting outlet changes
type OutletListApprovalRequest struct {
	OutletCrId []int `json:"outlet_cr_id" validate:"required,min=1"`
	Status     int   `json:"status" validate:"required,oneof=2 3"`
}

type OutletImportSecondaryCheckData struct {
	DistributorId             int64 `json:"distributor_id"`
	AllowUploadSecondarySales bool  `json:"allow_upload_secondary_sales"`
}
