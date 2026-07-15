package entity

import "time"

var (
	AR_TYPE_NORMAL          string = "Normal"
	AR_TYPE_1X_GIRO_TOLAKAN string = "1x Giro tolakan"
	AR_TYPE_2X_GIRO_TOLAKAN string = "2x Giro tolakan"
	AR_TYPE_3X_GIRO_TOLAKAN string = "3x Giro tolakan"

	AR_TYPE_NORMAL_ID          int = 1
	AR_TYPE_1X_GIRO_TOLAKAN_ID int = 2
	AR_TYPE_2X_GIRO_TOLAKAN_ID int = 3
	AR_TYPE_3X_GIRO_TOLAKAN_ID int = 4
)

type CreateMOutletBody struct {
	ParentCustId          string
	CustId                string        `json:"cust_id" validate:"required,max=10"`
	OutletCode            string        `json:"outlet_code"`
	SalesmanId            int           `json:"salesman_id"`
	OutletId              []int         `json:"outlet_id"`
	PjpId                 []int         `json:"pjp_id"`
	PjpCode               []int         `json:"pjp_code"`
	RouteCode             []int         `json:"route_code"`
	RouteName             []string      `json:"route_name"`
	OutletName            string        `json:"outlet_name" validate:"max=150"`
	Address               string        `json:"address" validate:"max=150"`
	PhoneNo               string        `json:"phone_no" validate:"max=20"`
	WaNo                  string        `json:"wa_no" validate:"max=20"`
	FaxNo                 string        `json:"fax_no" validate:"max=20"`
	BuldingOwn            int           `json:"building_own"`
	IsActive              bool          `json:"is_active"`
	CreatedBy             int64         `json:"created_by" validate:"required"`
	UpdatedBy             int64         `json:"updated_by"`
	Details               MOutletDetail `json:"details"`
	Latitude              string        `json:"latitude"`
	Longitude             string        `json:"longitude"`
	FileUrl               string        `json:"file_url"`
	VerificationStatus    *int          `json:"verification_status"`
	IsAddiitional         *bool         `json:"is_additional"`
	Source                int           `json:"source"`
	CreditLimitAction     int           `json:"credit_limit_action"`
	CreditLimitActionName string        `json:"credit_limit_action_name"`
}

type MOutletDetail struct {
	OutletContact []MOutletContact `json:"contact" validate:"dive"`
}

type MOutletContact struct {
	OutletContactId *int64  `json:"outlet_contact_id"`
	ContactName     *string `json:"contact_name" validate:"max=150"`
	JobTitle        string  `json:"job_title" validate:"max=100"`
	IdentityNo      string  `json:"identity_no" validate:"max=100"`
	PhoneNo         string  `json:"phone_no" validate:"max=20"`
	WaNo            string  `json:"wa_no" validate:"max=20"`
	Email           string  `json:"email" validate:""`
}

type MOutletRespone struct {
	OutletId               int              `json:"outlet_id"`
	SalesmanId             int              `json:"salesman_id"`
	OutletName             string           `json:"outlet_name"`
	Address                string           `json:"address"`
	PhoneNo                string           `json:"phone_no"`
	WaNo                   string           `json:"wa_no"`
	FaxNo                  string           `json:"fax_no"`
	Email                  string           `json:"email"`
	BuldingOwn             int              `json:"building_own"`
	IsActive               bool             `json:"is_active"`
	UpdatedBy              *int64           `json:"updated_by"`
	UpdatedAt              *time.Time       `json:"updated_at"`
	UpdatedByName          string           `json:"updated_by_name" db:"updated_by_name"`
	Details                OutletDetailRead `json:"details" validate:"dive"`
	Latitude               string           `json:"latitude"`
	Longitude              string           `json:"longitude"`
	VerificationStatus     *int             `json:"verification_status"`
	VerificationStatusName *string          `json:"verification_status_name"`
	VerifiedBy             *int64           `json:"verified_by"`
	VerifiedAt             *time.Time       `json:"verified_at"`
	VerifiedByName         string           `json:"verified_by_name" db:"verified_by_name"`
}

type MOutletListRespone struct {
	OutletId               int        `json:"outlet_id"`
	SalesmanId             int        `json:"salesman_id"`
	OutletName             string     `json:"outlet_name"`
	OutletCode             string     `json:"outlet_code"`
	OutletStatus           int        `json:"outlet_status"`
	OutletStatusName       string     `json:"outlet_status_name"`
	Address                string     `json:"address"`
	PhoneNo                string     `json:"phone_no"`
	WaNo                   string     `json:"wa_no"`
	FaxNo                  string     `json:"fax_no"`
	Email                  string     `json:"email"`
	BuldingOwn             int        `json:"building_own"`
	IsActive               bool       `json:"is_active"`
	UpdatedBy              *int64     `json:"updated_by"`
	UpdatedAt              *time.Time `json:"updated_at"`
	UpdatedByName          string     `json:"updated_by_name" db:"updated_by_name"`
	Latitude               string     `json:"latitude"`
	Longitude              string     `json:"longitude"`
	VerificationStatus     *int       `json:"verification_status"`
	VerificationStatusName *string    `json:"verification_status_name"`
	VerifiedBy             *int64     `json:"verified_by"`
	VerifiedAt             *time.Time `json:"verified_at"`
	VerifiedByName         string     `json:"verified_by_name" db:"verified_by_name"`
}

type OutletDetailRead struct {
	OutletContact []MOutletContactRead `json:"contact"  validate:"dive"`
	OutletPayment []MOutletPayment     `json:"payment_and_limit"  validate:"dive"`
}

type MOutletContactRead struct {
	OutletContactId *int64  `json:"outlet_contact_id"`
	OutletId        *int64  `json:"outlet_id"`
	ContactName     *string `json:"contact_name" validate:"required,max=150"`
	JobTitle        string  `json:"job_title" validate:"required,max=100"`
	PhoneNo         string  `json:"phone_no" validate:"max=20"`
	WaNo            string  `json:"wa_no" validate:"max=20"`
	Email           string  `json:"email" validate:"email,max=100"`
	IdentityNo      string  `json:"identity_no" validate:"max=100"`
}

type MOutletPayment struct {
	PaymentType       string  `json:"payment_type"`
	CreditLimitType   string  `json:"credit_limit_type"`
	CreditLimit       float64 `json:"credit_limit"`
	SalesInvLimitType string  `json:"sales_inv_limit_type"`
	SalesInvLimit     float64 `json:"sales_inv_limit"`
}

type MOutletQueryFilter struct {
	CustId             string
	ParentCustId       string
	Page               int    `query:"page"`
	Limit              int    `query:"limit" validate:"required"`
	VerificationStatus []int  `query:"verification_status"`
	Query              string `query:"q"`
	Mode               string `query:"mode"`
	Sort               string `query:"sort"`
	IsActive           *int   `query:"is_active"`
	SalesID            int    `query:"sales_id"`
	CreatedBy          int    `query:"created_by"`
	IsAdditional       *bool  `query:"is_additional"`
}

var dataOutletVerificationStatusName = map[int]string{
	1: "Approved",
	2: "In Review",
	3: "Rejected",
}

func (outlet MOutletListRespone) GenerateOutletVerificationStatusName() string {
	if outlet.VerificationStatus != nil {
		return dataOutletVerificationStatusName[*outlet.VerificationStatus]
	}
	return ""
}

type DetailOutletParams struct {
	OutletId int64 `params:"outlet_id" validate:"required"`
}

type DeleteOutletParams struct {
	OutletId int `params:"outlet_id" validate:"required"`
}

type OutletContactReads struct {
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
type OutletRespone struct {
	OutletId                   int64            `json:"outlet_id"`
	OutletCode                 string           `json:"outlet_code"`
	OutletName                 string           `json:"outlet_name"`
	Barcode                    string           `json:"barcode"`
	OutletStatus               int              `json:"outlet_status"`
	Address1                   string           `json:"address1"`
	Address2                   string           `json:"address2"`
	City                       string           `json:"city"`
	ZipCode                    string           `json:"zip_code"`
	PhoneNo                    string           `json:"phone_no"`
	WaNo                       string           `json:"wa_no"`
	FaxNo                      string           `json:"fax_no"`
	Email                      string           `json:"email"`
	DiscGrpId                  int              `json:"disc_grp_id"`
	DiscGrpCode                string           `json:"disc_grp_code"`
	DiscGrpName                string           `json:"disc_grp_name"`
	OtLocId                    int              `json:"ot_loc_id"`
	OtGrpId                    int              `json:"ot_grp_id"`
	PriceGrpId                 int              `json:"price_grp_id"`
	DistrictId                 int              `json:"district_id"`
	BeatId                     int              `json:"beat_id"`
	SbeatId                    int              `json:"sbeat_id"`
	OtClassId                  int              `json:"ot_class_id"`
	IndustryId                 int              `json:"industry_id"`
	MarketId                   int              `json:"market_id"`
	Top                        int              `json:"top"`
	Duedate                    string           `json:"due_date"`
	PaymentType                int              `json:"payment_type"`
	PaymentTypeName            *string          `json:"payment_type_name"`
	IsContraBon                bool             `json:"is_contra_bon"`
	PluGrpId                   int              `json:"plu_grp_id"`
	ConvGrpId                  int              `json:"conv_grp_id"`
	DiscInvId                  int              `json:"disc_inv_id"`
	AgentFrom                  string           `json:"agent_from"`
	CreditLimitType            *int             `json:"credit_limit_type"`
	CreditLimitTypeName        *string          `json:"credit_limit_type_name"`
	CreditLimit                float64          `json:"credit_limit"`
	CreditLimitAction          *int             `json:"credit_limit_action"`
	CreditLimitActionName      *string          `json:"credit_limit_action_name"`
	SalesInvLimitType          *int             `json:"sales_inv_limit_type"`
	SalesInvLimitTypeName      *string          `json:"sales_inv_limit_type_name"`
	SalesInvLimit              int              `json:"sales_inv_limit"`
	SalesInvLimitAction        *int             `json:"sales_inv_limit_action"`
	SalesInvLimitActionName    *string          `json:"sales_inv_limit_action_name"`
	AvgSalesWeek               float64          `json:"avg_sales_week"`
	AvgSalesMonth              float64          `json:"avg_sales_month"`
	FirstTransDate             string           `json:"first_trans_date"`
	LastTransDate              string           `json:"last_trans_date"`
	FirstWeekNo                int              `json:"first_week_no"`
	OtStartDate                string           `json:"ot_start_date"`
	OtRegDate                  string           `json:"ot_reg_date"`
	BuldingOwn                 int              `json:"building_own"`
	Dob                        string           `json:"dob"`
	ArStatus                   int              `json:"ar_status"`
	ArStatusName               string           `json:"ar_status_name"`
	ArTotal                    float64          `json:"ar_total"`
	CloseDate                  string           `json:"closed_date"`
	IsEmbBail                  bool             `json:"is_emb_bail"`
	TaxName                    string           `json:"tax_name"`
	TaxAddr1                   string           `json:"tax_addr1"`
	TaxAddr2                   string           `json:"tax_addr2"`
	TaxCity                    string           `json:"tax_city"`
	TaxNo                      string           `json:"tax_no"`
	TaxInvoiceForm             int              `json:"tax_invoice_form"`
	TaxInvoiceFormName         string           `json:"tax_invoice_form_name"`
	OwnerName                  string           `json:"owner_name"`
	OwnerAdd1                  string           `json:"owner_addr1"`
	OwnerAddr2                 string           `json:"owner_addr2"`
	OwnerCity                  string           `json:"owner_city"`
	OwnerPhoneNo               string           `json:"owner_phone_no"`
	OwnerIdNo                  string           `json:"owner_id_no"`
	DelvAdd1                   string           `json:"delv_addr1"`
	DelvCity                   string           `json:"delv_city"`
	DelvLatitude               string           `json:"delv_latitude"`
	DelvLongitude              string           `json:"delv_longitude"`
	DelvAddr2                  string           `json:"delv_addr2"`
	DelvCity2                  string           `json:"delv_city2"`
	DelvLatitude2              string           `json:"delv_latitude2"`
	DelvLongitude2             string           `json:"delv_longitude2"`
	InvAddr1                   string           `json:"inv_addr1"`
	InvAddr2                   string           `json:"inv_addr2"`
	InvCity                    string           `json:"inv_city"`
	IsActive                   bool             `json:"is_active"`
	UpdatedBy                  *int64           `json:"updated_by"`
	UpdatedAt                  *time.Time       `json:"updated_at"`
	UpdatedByName              string           `json:"updated_by_name" db:"updated_by_name"`
	Details                    OutletDetailRead `json:"details" validate:"dive"`
	Latitude                   string           `json:"latitude"`
	Longitude                  string           `json:"longitude"`
	ImageUrl                   string           `json:"image_url,omitempty"`
	OtTypeId                   int              `json:"ot_type_id"`
	OtTypeName                 string           `json:"ot_type_name"`
	IsObs                      bool             `json:"is_obs"`
	Obs                        int              `json:"obs"`
	ObsType                    *int             `json:"obs_type"`
	ObsTypeName                *string          `json:"obs_type_name"`
	ObsLimitAction             *int             `json:"obs_limit_action"`
	ObsLimitActionName         *string          `json:"obs_limit_action_name"`
	OutletLocationName         string           `json:"ot_loc_name"`
	MarketName                 string           `json:"market_name"`
	OutletGroupName            string           `json:"ot_grp_name"`
	PriceGroupName             string           `json:"price_grp_name"`
	DistrictName               string           `json:"district_name"`
	IndustryName               string           `json:"industry_name"`
	OutletClassName            string           `json:"ot_class_name"`
	OutletWardId               *string          `json:"outlet_ward_id"`
	OutletWard                 *string          `json:"outlet_ward"`
	OutletSubDistrictId        *string          `json:"outlet_sub_district_id"`
	OutletSubDistrict          *string          `json:"outlet_sub_district"`
	OutletRegencyId            *string          `json:"outlet_regency_id"`
	OutletRegency              *string          `json:"outlet_regency"`
	OutletProvinceId           *string          `json:"outlet_province_id"`
	OutletProvince             *string          `json:"outlet_province"`
	IsWaNo                     *bool            `json:"is_wa_no"`
	DelvWardId                 *string          `json:"delv_ward_id"`
	DelvWard                   *string          `json:"delv_ward"`
	DelvSubDistrictId          *string          `json:"delv_sub_district_id"`
	DelvSubDistrict            *string          `json:"delv_sub_district"`
	DelvRegencyId              *string          `json:"delv_regency_id"`
	DelvRegency                *string          `json:"delv_regency"`
	DelvProvinceId             *string          `json:"delv_province_id"`
	DelvProvince               *string          `json:"delv_province"`
	DelvZipCode                *string          `json:"delv_zip_code"`
	DelvIsSameAddress          *bool            `json:"delv_is_same_addr"`
	DelvWardId2                string           `json:"delv_ward_id2"`
	DelvWard2                  *string          `json:"delv_ward2"`
	DelvSubDistrictId2         *string          `json:"delv_sub_district_id2"`
	DelvSubDistrict2           *string          `json:"delv_sub_district2"`
	DelvRegencyId2             *string          `json:"delv_regency_id2"`
	DelvRegency2               *string          `json:"delv_regency2"`
	DelvProvinceId2            *string          `json:"delv_province_id2"`
	DelvProvince2              *string          `json:"delv_province2"`
	DelvZipCode2               *string          `json:"delv_zip_code2"`
	InvWardId                  *string          `json:"inv_ward_id"`
	InvWard                    *string          `json:"inv_ward"`
	InvSubDistrictId           *string          `json:"inv_sub_district_id"`
	InvSubDistrict             *string          `json:"inv_sub_district"`
	InvRegencyId               *string          `json:"inv_regency_id"`
	InvRegency                 *string          `json:"inv_regency"`
	InvProvinceId              *string          `json:"inv_province_id"`
	InvProvince                *string          `json:"inv_province"`
	InvZipCode                 *string          `json:"inv_zip_code"`
	InvIsSameAddress           *bool            `json:"inv_is_same_addr"`
	VerificationStatus         *int             `json:"verification_status"`
	VerificationStatusName     *string          `json:"verification_status_name"`
	VerifiedBy                 *int64           `json:"verified_by"`
	VerifiedAt                 *time.Time       `json:"verified_at"`
	VerifiedByName             string           `json:"verified_by_name"`
	OutletEstablishmentDate    *string          `json:"outlet_establishment_date"`
	RemainingOutstandingAmount float64          `json:"remaining_amount"`
	FileURL                    *string          `json:"file_url"`
}

type MobileOutletListQueryFilter struct {
	Query        string `query:"q"`
	Page         int    `query:"page" validate:"required"`
	Limit        int    `query:"limit" validate:"required"`
	Sort         string `query:"sort" validate:"required"`
	OutletStatus []int  `query:"outlet_status"`
}

type MobileOutletListResponse struct {
	OutletId     int     `json:"outlet_id"`
	OutletCode   string  `json:"outlet_code"`
	OutletName   string  `json:"outlet_name"`
	OutletStatus int     `json:"outlet_status"`
	Address1     string  `json:"address1"`
	Latitude     string  `json:"latitude"`
	Longitude    string  `json:"longitude"`
	AvgSalesWeek float64 `json:"avg_sales_week"`
}

// MobileOutletDetailResponse represents outlet detail for mobile
type MobileOutletDetailResponse struct {
	OutletId     int64                      `json:"outlet_id"`
	OutletCode   string                     `json:"outlet_code"`
	OutletName   string                     `json:"outlet_name"`
	Address1     string                     `json:"address1"`
	PhoneNo      string                     `json:"phone_no"`
	BuildingOwn  int                        `json:"building_own"`
	FileUrl      string                     `json:"file_url"`
	Longitude    string                     `json:"longitude"`
	Latitude     string                     `json:"latitude"`
	OtherContact *MobileOutletContactDetail `json:"other_contact"`
}

// MobileOutletContactDetail represents contact detail for outlet
type MobileOutletContactDetail struct {
	ContactName string `json:"contact_name"`
	JobTitle    string `json:"job_title"`
	PhoneNo     string `json:"phone_no"`
	WaNo        string `json:"wa_no"`
	Email       string `json:"email"`
}

type OutletPJPListQuery struct {
	EmpID         int64  `query:"-"`
	Date          string `query:"date" validate:"required,len=10,datetime=2006-01-02"`
	DistributorID int    `query:"distributor_id"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	Sort          string `query:"sort"`
	SortOrder     string `query:"-"`
}

type OutletPJPListResponse struct {
	RouteCode    int                      `json:"route_code"`
	RouteName    string                   `json:"route_name"`
	Week         int                      `json:"week"`
	Year         int                      `json:"year"`
	Date         time.Time                `json:"date"`
	Outlets      []OutletPJPResponse      `json:"outlets"`
	Distributors []DistributorPJPResponse `json:"distributors"`
}

type OutletPJPResponse struct {
	OutletID      int    `json:"outlet_id"`
	OutletCode    string `json:"outlet_code"`
	OutletName    string `json:"outlet_name"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	OutletStatus  string `json:"outlet_status"`
	OutletAddress string `json:"outlet_address"`
}

type DistributorPJPResponse struct {
	DistributorID      int    `json:"distributor_id"`
	DistributorCode    string `json:"distributor_code"`
	DistributorName    string `json:"distributor_name"`
	DistributorStatus  string `json:"distributor_status"`
	DistributorAddress string `json:"distributor_address"`
	Longitude          string `json:"longitude"`
	Latitude           string `json:"latitude"`
}
