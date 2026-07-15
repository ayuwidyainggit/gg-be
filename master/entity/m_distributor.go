package entity

import "time"

// DistributorSetup represents the distributor permission/setup configuration
type DistributorSetup struct {
	AllowAddProduct           bool `json:"allow_add_product"`
	AllowEditProduct          bool `json:"allow_edit_product"`
	AllowManagePricing        bool `json:"allow_manage_pricing"`
	AllowUploadSecondarySales bool `json:"allow_upload_secondary_sales"`
}

type DistributorQueryFilter struct {
	CustId           string
	ParentCustId     string
	JwtDistributorId int64  // from JWT claim, to be used for authorization check, not for query filter
	Page             int    `query:"page"`
	Limit            int    `query:"limit" validate:"required"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
	AreaID           []int  `query:"area_id"`
}

type DistributorResponse struct {
	ParentCustId            string               `json:"parent_cust_id"`
	DistributorId           int                  `json:"distributor_id"`
	DistributorCode         string               `json:"distributor_code"`
	DistributorName         string               `json:"distributor_name"`
	Barcode                 string               `json:"barcode"`
	RegionId                int64                `json:"region_id"`
	AreaId                  int64                `json:"area_id"`
	ChannelId               int64                `json:"channel_id"`
	SubDistributorGroupId   int64                `json:"sub_distributor_group_id"`
	SubDistributorGroupCode string               `json:"sub_distributor_group_code"`
	SubDistributorGroupName string               `json:"sub_distributor_group_name"`
	DistPriceGrpId          int64                `json:"dist_price_grp_id"`
	Address                 string               `json:"address"`
	ProvinceId              string               `json:"province_id"`
	RegencyId               string               `json:"regency_id"`
	SubDistrictId           string               `json:"sub_district_id"`
	WardId                  string               `json:"ward_id"`
	ZipCode                 string               `json:"zip_code"`
	OtLocId                 int                  `json:"ot_loc_id"`
	Latitude                string               `json:"latitude"`
	Longitude               string               `json:"longitude"`
	Phone                   string               `json:"phone"`
	FaxNumber               string               `json:"fax_number"`
	DistPriceGrpCode        string               `json:"dist_price_grp_code"`
	DistPriceGrpName        string               `json:"dist_price_grp_name"`
	RegionCode              string               `json:"region_code"`
	RegionName              string               `json:"region_name"`
	ChannelCode             string               `json:"channel_code"`
	ChannelName             string               `json:"channel_name"`
	AreaCode                string               `json:"area_code"`
	AreaName                string               `json:"area_name"`
	ProvinceCode            string               `json:"province_code"`
	ProvinceName            string               `json:"province_name"`
	RegencyCode             string               `json:"regency_code"`
	RegencyName             string               `json:"regency_name"`
	SubDistrictCode         string               `json:"sub_district_code"`
	SubDistrictName         string               `json:"sub_district_name"`
	WardCode                string               `json:"ward_code"`
	WardName                string               `json:"ward_name"`
	IsActive                bool                 `json:"is_active"`
	UpdatedBy               int64                `json:"updated_by"`
	UpdatedByName           string               `json:"updated_by_name"`
	UpdatedAt               *time.Time           `json:"updated_at"`
	DistributorContact      []DistributorContact `json:"contacts"`
	DistributorTax          []DistributorTax     `json:"tax"`
	DistributorSetup        *DistributorSetup    `json:"distributor_setup"`
}

type DistributorListRespone struct {
	CustId                  string     `json:"cust_id"`
	ParentCustId            string     `json:"parent_cust_id"`
	DistributorId           int64      `json:"distributor_id"`
	DistributorCode         string     `json:"distributor_code"`
	DistributorName         string     `json:"distributor_name"`
	Barcode                 *string    `json:"barcode"`
	RegionId                int        `json:"region_id"`
	AreaId                  int        `json:"area_id"`
	ChannelId               int        `json:"channel_id"`
	SubDistributorGroupId   int        `json:"sub_distributor_group_id"`
	SubDistributorGroupCode string     `json:"sub_distributor_group_code"`
	SubDistributorGroupName string     `json:"sub_distributor_group_name"`
	DistPriceGrpId          int        `json:"dist_price_grp_id"`
	Address                 string     `json:"address"`
	ProvinceId              string     `json:"province_id"`
	RegencyId               string     `json:"regency_id"`
	SubDistrictId           string     `json:"sub_district_id"`
	WardId                  string     `json:"ward_id"`
	ZipCode                 string     `json:"zip_code"`
	OtLocId                 int        `json:"ot_loc_id"`
	Latitude                string     `json:"latitude"`
	Longitude               string     `json:"longitude"`
	Phone                   string     `json:"phone"`
	FaxNumber               string     `json:"fax_number"`
	ContactName             string     `json:"contact_name"`
	JobTitle                string     `json:"job_title"`
	PhoneNo                 string     `json:"phone_no"`
	WaNo                    string     `json:"wa_no"`
	Email                   string     `json:"email"`
	DistPriceGrpCode        string     `json:"dist_price_grp_code"`
	DistPriceGrpName        string     `json:"dist_price_grp_name"`
	RegionCode              string     `json:"region_code"`
	RegionName              string     `json:"region_name"`
	ChannelCode             string     `json:"channel_code"`
	ChannelName             string     `json:"channel_name"`
	AreaCode                string     `json:"area_code"`
	AreaName                string     `json:"area_name"`
	ProvinceCode            string     `json:"province_code"`
	ProvinceName            string     `json:"province_name"`
	RegencyCode             string     `json:"regency_code"`
	RegencyName             string     `json:"regency_name"`
	SubDistrictCode         string     `json:"sub_district_code"`
	SubDistrictName         string     `json:"sub_district_name"`
	WardCode                string     `json:"ward_code"`
	WardName                string     `json:"ward_name"`
	IsActive                bool       `json:"is_active"`
	IsDel                   bool       `json:"is_del"`
	CreatedBy               *int64     `json:"createdby"`
	CreatedAt               *time.Time `json:"created_at"`
	UpdatedBy               *int64     `json:"updatedby"`
	UpdatedAt               *time.Time `json:"updated_at"`
	UpdatedByName           string     `json:"updated_by_name"`
	DeletedBy               *int64     `json:"deleted_by"`
	DeletedAt               *time.Time `json:"deleted_at"`
	CustomerID              string     `json:"customer_id"`
}

type DistributorLookupResponse struct {
	DistributorId         int    `json:"distributor_id"`
	DistributorCode       string `json:"distributor_code"`
	DistributorName       string `json:"distributor_name"`
	Barcode               string `json:"barcode"`
	RegionId              int    `json:"region_id"`
	AreaId                int    `json:"area_id"`
	ChannelId             int    `json:"channel_id"`
	SubDistributorGroupId int    `json:"sub_distributor_group_id"`
	DistPriceGrpId        int    `json:"dist_price_grp_id"`
	Address               string `json:"address"`
	ProvinceId            string `json:"province_id"`
	RegencyId             string `json:"regency_id"`
	SubDistrictId         string `json:"sub_district_id"`
	WardId                string `json:"ward_id"`
	ZipCode               string `json:"zip_code"`
	OtLocId               int    `json:"ot_loc_id"`
	Latitude              string `json:"latitude"`
	Longitude             string `json:"longitude"`
	Phone                 int    `json:"phone"`
	FaxNumber             int    `json:"fax_number"`
	IsActive              bool   `json:"is_active"`
	CustomerID            string `json:"customer_id"`
}

type CreateDistributorBody struct {
	CustId                string               `json:"cust_id" validate:"required,max=10"`
	ParentCustId          string               `json:"parent_cust_id" validate:"required,max=10"`
	CreatedBy             *int64               `json:"created_by" validate:"required"`
	DistributorId         int                  `json:"distributor_id" validate:""`
	DistributorCode       string               `json:"distributor_code" validate:"required,alphanumDashUnderscore,max=20"`
	DistributorName       string               `json:"distributor_name" validate:"required,alphanumericSpace,max=150"`
	Barcode               string               `json:"barcode" validate:"omitempty,alphanumDashUnderscore,max=20"`
	RegionId              int                  `json:"region_id" validate:"required"`
	AreaId                int                  `json:"area_id" validate:"required"`
	ChannelId             int                  `json:"channel_id" validate:"required"`
	SubDistributorGroupId int                  `json:"sub_distributor_group_id" validate:"required"`
	DistPriceGrpId        int                  `json:"dist_price_grp_id" validate:"required"`
	Address               string               `json:"address" validate:"required,max=255"`
	ProvinceId            string               `json:"province_id" validate:""`
	RegencyId             string               `json:"regency_id" validate:""`
	SubDistrictId         string               `json:"sub_district_id" validate:""`
	WardId                string               `json:"ward_id" validate:""`
	ZipCode               string               `json:"zip_code" validate:"omitempty,alphanum,max=6"`
	OtLocId               *int                 `json:"ot_loc_id"`
	Latitude              string               `json:"latitude" validate:"required,max=50"`
	Longitude             string               `json:"longitude" validate:"required,max=50"`
	Phone                 string               `json:"phone" validate:"omitempty,max=25"`
	FaxNumber             string               `json:"fax_number" validate:"omitempty,max=25"`
	IsActive              *bool                `json:"is_active" validate:"required"`
	Contacts              []DistributorContact `json:"contacts"`
	Tax                   []DistributorTax     `json:"tax"`
	DistributorSetup      *DistributorSetup    `json:"distributor_setup"`
}

type DetailDistributorParams struct {
	CustId           string
	ParentCustId     string
	JwtDistributorId int64
	DistributorId    int `params:"distributor_id" validate:"required"`
}

type UpdateDistributorParams struct {
	DistributorId int64 `params:"distributor_id" validate:"required"`
}

type DeleteDistributorParams struct {
	DistributorId int `params:"distributor_id" validate:"required"`
}

type UpdateDistributorRequest struct {
	CustId                    string                     `json:"cust_id" validate:"required,max=10"`
	UpdatedBy                 *int64                     `json:"updated_by" validate:"required"`
	BarcodeProvided           bool                       `json:"-"`
	ProvinceIdProvided        bool                       `json:"-"`
	RegencyIdProvided         bool                       `json:"-"`
	SubDistrictIdProvided     bool                       `json:"-"`
	WardIdProvided            bool                       `json:"-"`
	DistributorId             int                        `json:"distributor_id"`
	DistributorCode           *string                    `json:"distributor_code" validate:"omitempty,alphanumDashUnderscore,max=20"`
	DistributorName           *string                    `json:"distributor_name" validate:"omitempty,alphanumericSpace,max=150"`
	Barcode                   *string                    `json:"barcode" validate:"omitempty,alphanumDashUnderscore,max=20"`
	RegionId                  *int                       `json:"region_id" validate:"omitempty,gt=0"`
	AreaId                    *int                       `json:"area_id" validate:"omitempty,gt=0"`
	ChannelId                 *int                       `json:"channel_id" validate:"omitempty,gt=0"`
	SubDistributorGroupId     *int                       `json:"sub_distributor_group_id" validate:"omitempty,gt=0"`
	DistPriceGrpId            *int                       `json:"dist_price_grp_id" validate:"omitempty,gt=0"`
	Address                   *string                    `json:"address" validate:"omitempty,max=255"`
	ProvinceId                *string                    `json:"province_id" validate:"omitempty,numeric,max=20"`
	RegencyId                 *string                    `json:"regency_id" validate:"omitempty,numeric,max=20"`
	SubDistrictId             *string                    `json:"sub_district_id" validate:"omitempty,numeric,max=20"`
	WardId                    *string                    `json:"ward_id" validate:"omitempty,numeric,max=20"`
	ZipCodeProvided           bool                       `json:"-"`
	OtLocIdProvided           bool                       `json:"-"`
	ZipCode                   *string                    `json:"zip_code" validate:"omitempty,alphanum,max=6"`
	OtLocId                   *int                       `json:"ot_loc_id"`
	Latitude                  *string                    `json:"latitude" validate:"omitempty,max=50"`
	Longitude                 *string                    `json:"longitude" validate:"omitempty,max=50"`
	PhoneProvided             bool                       `json:"-"`
	Phone                     *string                    `json:"phone" validate:"omitempty,max=25"`
	FaxNumberProvided         bool                       `json:"-"`
	FaxNumber                 *string                    `json:"fax_number" validate:"omitempty,max=25"`
	IsActive                  *bool                      `json:"is_active,omitempty" validate:"omitempty"`
	Contacts                  []DistributorContactUpdate `json:"contacts"`
	Tax                       []DistributorTaxUpdate     `json:"tax"`
	DistributorSetup          *DistributorSetup          `json:"distributor_setup"`
	AllowAddProduct           *bool                      `json:"allow_add_product,omitempty"`
	AllowEditProduct          *bool                      `json:"allow_edit_product,omitempty"`
	AllowManagePricing        *bool                      `json:"allow_manage_pricing,omitempty"`
	AllowUploadSecondarySales *bool                      `json:"allow_upload_secondary_sales,omitempty"`
}

type DistributorAreaRegionData struct {
	DistributorID   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	RegionID        int    `json:"region_id"`
	RegionName      string `json:"region_name"`
	AreaID          int    `json:"area_id"`
	AreaName        string `json:"area_name"`
	IsUpdated       bool   `json:"is_updated"`
	UpdateStatus    string `json:"update_status"`
}

type DistributorCustomerResp struct {
	ParentCustId            string               `json:"parent_cust_id"`
	DistributorId           int                  `json:"distributor_id"`
	DistributorCode         string               `json:"distributor_code"`
	DistributorName         string               `json:"distributor_name"`
	Barcode                 string               `json:"barcode"`
	RegionId                int64                `json:"region_id"`
	AreaId                  int64                `json:"area_id"`
	ChannelId               int64                `json:"channel_id"`
	SubDistributorGroupId   int64                `json:"sub_distributor_group_id"`
	SubDistributorGroupCode string               `json:"sub_distributor_group_code"`
	SubDistributorGroupName string               `json:"sub_distributor_group_name"`
	DistPriceGrpId          int64                `json:"dist_price_grp_id"`
	Address                 string               `json:"address"`
	ProvinceId              string               `json:"province_id"`
	RegencyId               string               `json:"regency_id"`
	SubDistrictId           string               `json:"sub_district_id"`
	WardId                  string               `json:"ward_id"`
	ZipCode                 string               `json:"zip_code"`
	OtLocId                 int                  `json:"ot_loc_id"`
	Latitude                string               `json:"latitude"`
	Longitude               string               `json:"longitude"`
	Phone                   string               `json:"phone"`
	FaxNumber               string               `json:"fax_number"`
	DistPriceGrpCode        string               `json:"dist_price_grp_code"`
	DistPriceGrpName        string               `json:"dist_price_grp_name"`
	RegionCode              string               `json:"region_code"`
	RegionName              string               `json:"region_name"`
	ChannelCode             string               `json:"channel_code"`
	ChannelName             string               `json:"channel_name"`
	AreaCode                string               `json:"area_code"`
	AreaName                string               `json:"area_name"`
	ProvinceCode            string               `json:"province_code"`
	ProvinceName            string               `json:"province_name"`
	RegencyCode             string               `json:"regency_code"`
	RegencyName             string               `json:"regency_name"`
	SubDistrictCode         string               `json:"sub_district_code"`
	SubDistrictName         string               `json:"sub_district_name"`
	WardCode                string               `json:"ward_code"`
	WardName                string               `json:"ward_name"`
	IsActive                bool                 `json:"is_active"`
	UpdatedBy               int64                `json:"updated_by"`
	UpdatedByName           string               `json:"updated_by_name"`
	UpdatedAt               *time.Time           `json:"updated_at"`
	DistributorContact      []DistributorContact `json:"contacts"`
	DistributorTax          []DistributorTax     `json:"tax"`
	CustomerID              string               `json:"customer_id"`
}
