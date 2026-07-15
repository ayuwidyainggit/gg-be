package model

import "time"

type MobileDistributorList struct {
	DistributorId   int    `db:"distributor_id" json:"distributor_id"`
	DistributorCode string `db:"distributor_code" json:"distributor_code"`
	DistributorName string `db:"distributor_name" json:"distributor_name"`
	Address         string `db:"address" json:"address"`
	Latitude        string `db:"latitude" json:"latitude"`
	Longitude       string `db:"longitude" json:"longitude"`
	RegionID        int    `db:"region_id" json:"region_id"`
	AreaID          int    `db:"area_id" json:"area_id"`
}

type Distributor struct {
	CustID                    string     `gorm:"column:cust_id"`
	DistributorID             int64      `gorm:"column:distributor_id;primaryKey"`
	DistributorCode           *string    `gorm:"column:distributor_code"`
	DistributorName           *string    `gorm:"column:distributor_name"`
	Barcode                   *string    `gorm:"column:barcode"`
	RegionID                  *int64     `gorm:"column:region_id"`
	AreaID                    *int64     `gorm:"column:area_id"`
	ChannelID                 *int64     `gorm:"column:channel_id"`
	SubDistributorGroupID     *int64     `gorm:"column:sub_distributor_group_id"`
	DistPriceGrpID            *int64     `gorm:"column:dist_price_grp_id"`
	Address                   *string    `gorm:"column:address"`
	ProvinceID                *string    `gorm:"column:province_id"`
	RegencyID                 *string    `gorm:"column:regency_id"`
	SubDistrictID             *string    `gorm:"column:sub_district_id"`
	WardID                    *string    `gorm:"column:ward_id"`
	ZipCode                   *string    `gorm:"column:zip_code"`
	OtLocID                   *int64     `gorm:"column:ot_loc_id"`
	Latitude                  *string    `gorm:"column:latitude"`
	Longitude                 *string    `gorm:"column:longitude"`
	IsActive                  *bool      `gorm:"column:is_active"`
	CreatedBy                 *int64     `gorm:"column:created_by"`
	CreatedAt                 *time.Time `gorm:"column:created_at"`
	UpdatedBy                 *int64     `gorm:"column:updated_by"`
	UpdatedAt                 *time.Time `gorm:"column:updated_at"`
	IsDel                     *bool      `gorm:"column:is_del"`
	DeletedBy                 *int64     `gorm:"column:deleted_by"`
	DeletedAt                 *time.Time `gorm:"column:deleted_at"`
	Phone                     *string    `gorm:"column:phone"`
	FaxNumber                 *string    `gorm:"column:fax_number"`
	AllowAddProduct           *bool      `gorm:"column:allow_add_product"`
	AllowEditProduct          *bool      `gorm:"column:allow_edit_product"`
	AllowManagePricing        *bool      `gorm:"column:allow_manage_pricing"`
	AllowUploadSecondarySales *bool      `gorm:"column:allow_upload_secondary_sales"`
	ParentCustID              *string    `gorm:"column:parent_cust_id"`
}

func (Distributor) TableName() string {
	return "mst.m_distributor"
}

type PrincipalInfo struct {
	CustID        string `gorm:"column:cust_id"`
	CustName      string `gorm:"column:cust_name"`
	DistributorID int64  `gorm:"column:distributor_id"`
}
