package model

import "time"

type Distributor struct {
	CustId                    string     `db:"cust_id" json:"cust_id"`
	ParentCustId              string     `db:"parent_cust_id" json:"parent_cust_id"`
	DistributorId             int64      `db:"distributor_id" json:"distributor_id"`
	DistributorCode           string     `db:"distributor_code" json:"distributor_code"`
	DistributorName           string     `db:"distributor_name" json:"distributor_name"`
	Barcode                   string     `db:"barcode" json:"barcode"`
	RegionId                  int        `db:"region_id" json:"region_id"`
	AreaId                    int        `db:"area_id" json:"area_id"`
	ChannelId                 int        `db:"channel_id" json:"channel_id"`
	SubDistributorGroupId     int        `db:"sub_distributor_group_id" json:"sub_distributor_group_id"`
	DistPriceGrpId            int        `db:"dist_price_grp_id" json:"dist_price_grp_id"`
	Address                   string     `db:"address" json:"address"`
	ProvinceId                string     `db:"province_id" json:"province_id"`
	RegencyId                 string     `db:"regency_id" json:"regency_id"`
	SubDistrictId             string     `db:"sub_district_id" json:"sub_district_id"`
	WardId                    string     `db:"ward_id" json:"ward_id"`
	ZipCode                   string     `db:"zip_code" json:"zip_code"`
	OtLocId                   int        `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude                  string     `db:"latitude" json:"latitude"`
	Longitude                 string     `db:"longitude" json:"longitude"`
	Phone                     string     `db:"phone" json:"phone"`
	FaxNumber                 string     `db:"fax_number" json:"fax_number"`
	IsActive                  bool       `db:"is_active" json:"is_active"`
	CreatedBy                 *int64     `db:"created_by" json:"created_by"`
	CreatedAt                 *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy                 *int64     `db:"updated_by" json:"updated_by"`
	UpdatedByName             *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt                 *time.Time `db:"updated_at" json:"updated_at"`
	IsDel                     bool       `db:"is_del" json:"is_del"`
	DeletedBy                 *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt                 *time.Time `db:"deleted_at" json:"deleted_at"`
	AllowAddProduct           bool       `db:"allow_add_product" json:"allow_add_product"`
	AllowEditProduct          bool       `db:"allow_edit_product" json:"allow_edit_product"`
	AllowManagePricing        bool       `db:"allow_manage_pricing" json:"allow_manage_pricing"`
	AllowUploadSecondarySales bool       `db:"allow_upload_secondary_sales" json:"allow_upload_secondary_sales"`
}

type DistributorList struct {
	CustId                    string     `db:"cust_id" json:"cust_id"`
	ParentCustId              string     `db:"parent_cust_id" json:"parent_cust_id"`
	DistributorId             int64      `db:"distributor_id" json:"distributor_id"`
	DistributorCode           string     `db:"distributor_code" json:"distributor_code"`
	DistributorName           string     `db:"distributor_name" json:"distributor_name"`
	Barcode                   *string    `db:"barcode" json:"barcode"`
	RegionId                  int        `db:"region_id" json:"region_id"`
	AreaId                    int        `db:"area_id" json:"area_id"`
	ChannelId                 int        `db:"channel_id" json:"channel_id"`
	SubDistributorGroupId     int        `db:"sub_distributor_group_id" json:"sub_distributor_group_id"`
	SubDistributorGroupCode   string     `db:"sub_distributor_group_code" json:"sub_distributor_group_code"`
	SubDistributorGroupName   string     `db:"sub_distributor_group_name" json:"sub_distributor_group_name"`
	DistPriceGrpId            int        `db:"dist_price_grp_id" json:"dist_price_grp_id"`
	Address                   string     `db:"address" json:"address"`
	ProvinceId                string     `db:"province_id" json:"province_id"`
	RegencyId                 string     `db:"regency_id" json:"regency_id"`
	SubDistrictId             string     `db:"sub_district_id" json:"sub_district_id"`
	WardId                    string     `db:"ward_id" json:"ward_id"`
	ZipCode                   string     `db:"zip_code" json:"zip_code"`
	OtLocId                   int        `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude                  string     `db:"latitude" json:"latitude"`
	Longitude                 string     `db:"longitude" json:"longitude"`
	Phone                     *string    `db:"phone" json:"phone"`
	FaxNumber                 *string    `db:"fax_number" json:"fax_number"`
	ContactName               string     `db:"contact_name" json:"contact_name"`
	DistributorContactId      int        `db:"distributor_contact_id" json:"distributor_contact_id"`
	JobTitle                  string     `db:"job_title" json:"job_title"`
	PhoneNo                   string     `db:"phone_no" json:"phone_no"`
	IsWaNo                    bool       `db:"is_wa_no" json:"is_wa_no"`
	WaNo                      *string    `db:"wa_no" json:"wa_no"`
	Email                     string     `db:"email" json:"email"`
	DistPriceGrpCode          *string    `db:"dist_price_grp_code" json:"dist_price_grp_code"`
	DistPriceGrpName          *string    `db:"dist_price_grp_name" json:"dist_price_grp_name"`
	RegionCode                *string    `db:"region_code" json:"region_code"`
	RegionName                *string    `db:"region_name" json:"region_name"`
	ChannelCode               *string    `db:"channel_code" json:"channel_code"`
	ChannelName               *string    `db:"channel_name" json:"channel_name"`
	AreaCode                  *string    `db:"area_code" json:"area_code"`
	AreaName                  *string    `db:"area_name" json:"area_name"`
	ProvinceCode              *string    `db:"province_code" json:"province_code"`
	ProvinceName              *string    `db:"province_name" json:"province_name"`
	RegencyCode               *string    `db:"regency_code" json:"regency_code"`
	RegencyName               *string    `db:"regency_name" json:"regency_name"`
	SubDistrictCode           *string    `db:"sub_district_code" json:"sub_district_code"`
	SubDistrictName           *string    `db:"sub_district_name" json:"sub_district_name"`
	WardCode                  *string    `db:"ward_code" json:"ward_code"`
	WardName                  *string    `db:"ward_name" json:"ward_name"`
	IsActive                  bool       `db:"is_active" json:"is_active"`
	IsDel                     bool       `db:"is_del" json:"is_del"`
	CreatedBy                 *int64     `db:"created_by,omitempty" json:"created_by"`
	CreatedAt                 *time.Time `db:"created_at,omitempty" json:"created_at"`
	UpdatedBy                 *int64     `db:"updated_by,omitempty" json:"updated_by"`
	UpdatedAt                 *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedByName             *string    `db:"updated_by_name" json:"updated_by_name"`
	DeletedBy                 *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	DeletedAt                 *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	CustomerID                string     `db:"customer_id" json:"customer_id"`
	AllowAddProduct           bool       `db:"allow_add_product" json:"allow_add_product"`
	AllowEditProduct          bool       `db:"allow_edit_product" json:"allow_edit_product"`
	AllowManagePricing        bool       `db:"allow_manage_pricing" json:"allow_manage_pricing"`
	AllowUploadSecondarySales bool       `db:"allow_upload_secondary_sales" json:"allow_upload_secondary_sales"`
}

type DistributorUpdate struct {
	BarcodeProvided           bool       `db:"-" sql:"-" json:"-"`
	ProvinceIdProvided        bool       `db:"-" sql:"-" json:"-"`
	RegencyIdProvided         bool       `db:"-" sql:"-" json:"-"`
	SubDistrictIdProvided     bool       `db:"-" sql:"-" json:"-"`
	WardIdProvided            bool       `db:"-" sql:"-" json:"-"`
	DistributorCode           *string    `db:"distributor_code" json:"distributor_code"`
	DistributorName           *string    `db:"distributor_name" json:"distributor_name"`
	Barcode                   *string    `db:"barcode" json:"barcode"`
	RegionId                  *int       `db:"region_id" json:"region_id"`
	AreaId                    *int       `db:"area_id" json:"area_id"`
	ChannelId                 *int       `db:"channel_id" json:"channel_id"`
	SubDistributorGroupId     *int       `db:"sub_distributor_group_id" json:"sub_distributor_group_id"`
	DistPriceGrpId            *int       `db:"dist_price_grp_id" json:"dist_price_grp_id"`
	Address                   *string    `db:"address" json:"address"`
	ProvinceId                *string    `db:"province_id" json:"province_id"`
	RegencyId                 *string    `db:"regency_id" json:"regency_id"`
	SubDistrictId             *string    `db:"sub_district_id" json:"sub_district_id"`
	WardId                    *string    `db:"ward_id" json:"ward_id"`
	ZipCodeProvided           bool       `db:"-" sql:"-" json:"-"`
	OtLocIdProvided           bool       `db:"-" sql:"-" json:"-"`
	ZipCode                   *string    `db:"zip_code" json:"zip_code"`
	OtLocId                   *int       `db:"ot_loc_id" json:"ot_loc_id"`
	Latitude                  *string    `db:"latitude" json:"latitude"`
	Longitude                 *string    `db:"longitude" json:"longitude"`
	PhoneProvided             bool       `db:"-" sql:"-" json:"-"`
	Phone                     *string    `db:"phone" json:"phone"`
	FaxNumberProvided         bool       `db:"-" sql:"-" json:"-"`
	FaxNumber                 *string    `db:"fax_number" json:"fax_number"`
	IsActive                  *bool      `db:"is_active" json:"is_active"`
	CreatedBy                 *int64     `db:"created_by" json:"created_by"`
	CreatedAt                 *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy                 *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt                 *time.Time `db:"updated_at" json:"updated_at"`
	IsDel                     *bool      `db:"is_del" json:"is_del"`
	DeletedBy                 *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt                 *time.Time `db:"deleted_at" json:"deleted_at"`
	AllowAddProduct           *bool      `db:"allow_add_product" json:"allow_add_product"`
	AllowEditProduct          *bool      `db:"allow_edit_product" json:"allow_edit_product"`
	AllowManagePricing        *bool      `db:"allow_manage_pricing" json:"allow_manage_pricing"`
	AllowUploadSecondarySales *bool      `db:"allow_upload_secondary_sales" json:"allow_upload_secondary_sales"`
}

type DistributorAreaRegionDetail struct {
	DistributorID   int64  `db:"distributor_id" json:"distributor_id"`
	DistributorCode string `db:"distributor_code" json:"distributor_code"`
	DistributorName string `db:"distributor_name" json:"distributor_name"`
	RegionID        int    `db:"region_id" json:"region_id"`
	RegionName      string `db:"region_name" json:"region_name"`
	AreaID          int    `db:"area_id" json:"area_id"`
	AreaName        string `db:"area_name" json:"area_name"`
}

type DistributorWithCustomer struct {
	DistributorCustID string `db:"dist_cust_id" json:"dist_cust_id"`
	DistributorList
}
