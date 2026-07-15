package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GamificationList struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	GamificationNo   string     `gorm:"column:gamification_no;primaryKey" json:"gamification_no"`
	GamificationDate *time.Time `gorm:"column:gamification_date" json:"gamification_date"`
	Title            *string    `gorm:"column:title" json:"title"`
	// Description      *string    `gorm:"column:description" json:"description"`
	StartDate          *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate            *time.Time `gorm:"column:end_date" json:"end_date"`
	FinishedDate       *time.Time `gorm:"column:finished_date" json:"finished_date"`
	GamificationStatus *int64     `gorm:"column:gamification_status" json:"gamification_status"`
	// CreatedBy        *int64     `gorm:"column:created_by" json:"created_by"`
	// CreatedByName    *string    `gorm:"column:created_by_name" json:"created_by_name"`
	// CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GamificationList) TableName() string {
	return "sls.gamification"
}

type GamificationRead struct {
	CustID             string     `gorm:"column:cust_id" json:"cust_id"`
	GamificationNo     string     `gorm:"column:gamification_no;primaryKey" json:"gamification_no"`
	GamificationDate   time.Time  `gorm:"column:gamification_date" json:"gamification_date"`
	Title              string     `gorm:"column:title" json:"title"`
	Description        *string    `gorm:"column:description" json:"description"`
	EmpGrpID           int64      `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode         *string    `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGrpName         *string    `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	ProductFilter      string     `gorm:"column:product_filter" json:"product_filter"`
	OutletFilter       string     `gorm:"column:outlet_filter" json:"outlet_filter"`
	StartDate          time.Time  `gorm:"column:start_date" json:"start_date"`
	EndDate            time.Time  `gorm:"column:end_date" json:"end_date"`
	FinishedDate       *time.Time `gorm:"column:finished_date" json:"finished_date"`
	SourceID           int64      `gorm:"column:source_id" json:"source_id"`
	MeasurementID      int64      `gorm:"column:measurement_id" json:"measurement_id"`
	AnnouncementID     int64      `gorm:"column:announcement_id" json:"announcement_id"`
	AnnouncementTarget *float64   `gorm:"column:announcement_target" json:"announcement_target"`
	SubAnnouncementID  int64      `gorm:"column:subannouncement_id" json:"subannouncement_id"`
	GamificationStatus int64      `gorm:"column:gamification_status" json:"gamification_status"`
}

func (GamificationRead) TableName() string {
	return "sls.gamification"
}

type GamificationStatusesFilter struct {
	GamificationStatus int `gorm:"column:gamification_status" json:"gamification_status"`
}

func (GamificationStatusesFilter) TableName() string {
	return "sls.gamification"
}

type EmployeeGroups struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	EmpGroupId    int             `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGroupCode  string          `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGroupName  string          `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (EmployeeGroups) TableName() string {
	return "mst.m_emp_group"
}

type SalesTeams struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SalesTeamId   int             `gorm:"column:sales_team_id" json:"sales_team_id"`
	SalesTeamCode string          `gorm:"column:sales_team_code" json:"sales_team_code"`
	SalesTeamName string          `gorm:"column:sales_team_name" json:"sales_team_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SalesTeams) TableName() string {
	return "mst.m_sales_team"
}

type OperationTypes struct {
	OperationTypeCode string          `gorm:"column:operation_type_code" json:"operation_type_code"`
	OperationTypeName string          `gorm:"column:operation_type_name" json:"operation_type_name"`
	CreatedBy         *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt         *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy         *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt         *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName     *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy         *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt         *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OperationTypes) TableName() string {
	return "mst.m_operation_type"
}

type Warehouses struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	WhId          int             `gorm:"column:wh_id" json:"wh_id"`
	WhCode        string          `gorm:"column:wh_code" json:"wh_code"`
	WhName        string          `gorm:"column:wh_name" json:"wh_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Warehouses) TableName() string {
	return "mst.m_warehouse"
}

type Salesmans struct {
	SalesmanId   int     `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode *string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName *string `gorm:"column:salesman_name" json:"salesman_name"`
}

func (Salesmans) TableName() string {
	return "mst.m_salesman"
}

type ProductCategories struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	PCatId        int             `gorm:"column:pcat_id" json:"pcat_id"`
	PCatCode      string          `gorm:"column:pcat_code" json:"pcat_code"`
	PCatName      string          `gorm:"column:pcat_name" json:"pcat_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ProductCategories) TableName() string {
	return "mst.m_product_cat"
}

type ProductLines struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	PLId          int             `gorm:"column:pl_id" json:"pl_id"`
	PLCode        string          `gorm:"column:pl_code" json:"pl_code"`
	PLName        string          `gorm:"column:pl_name" json:"pl_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ProductLines) TableName() string {
	return "mst.m_product_line"
}

type Brands struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	BrandId       int             `gorm:"column:brand_id" json:"brand_id"`
	BrandCode     string          `gorm:"column:brand_code" json:"brand_code"`
	BrandName     string          `gorm:"column:brand_name" json:"brand_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Brands) TableName() string {
	return "mst.m_brand"
}

type SubBrands1 struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SBrand1Id     int             `gorm:"column:sbrand1_id" json:"sbrand1_id"`
	SBrand1Code   string          `gorm:"column:sbrand1_code" json:"sbrand1_code"`
	SBrand1Name   string          `gorm:"column:sbrand1_name" json:"sbrand1_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SubBrands1) TableName() string {
	return "mst.m_sub_brand1"
}

type SubBrands2 struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SBrand2Id     int             `gorm:"column:sbrand2_id" json:"sbrand2_id"`
	SBrand2Code   string          `gorm:"column:sbrand2_code" json:"sbrand2_code"`
	SBrand2Name   string          `gorm:"column:sbrand2_name" json:"sbrand2_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SubBrands2) TableName() string {
	return "mst.m_sub_brand2"
}

type Products struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	ProId         int             `gorm:"column:pro_id" json:"pro_id"`
	ProCode       string          `gorm:"column:pro_code" json:"pro_code"`
	ProName       string          `gorm:"column:pro_name" json:"pro_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Products) TableName() string {
	return "mst.m_product"
}

type OutletLocations struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OtLocId       int             `gorm:"column:ot_loc_id" json:"ot_loc_id"`
	OtLocCode     string          `gorm:"column:ot_loc_code" json:"ot_loc_code"`
	OtLocName     string          `gorm:"column:ot_loc_name" json:"ot_loc_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletLocations) TableName() string {
	return "mst.m_outlet_loc"
}

type OutletGroups struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OtGrpId       int             `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	OtGrpCode     string          `gorm:"column:ot_grp_code" json:"ot_grp_code"`
	OtGrpName     string          `gorm:"column:ot_grp_name" json:"ot_grp_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletGroups) TableName() string {
	return "mst.m_outlet_group"
}

type OutletClasses struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OtClassId     int             `gorm:"column:ot_class_id" json:"ot_class_id"`
	OtClassCode   string          `gorm:"column:ot_class_code" json:"ot_class_code"`
	OtClassName   string          `gorm:"column:ot_class_name" json:"ot_class_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletClasses) TableName() string {
	return "mst.m_outlet_class"
}

type Markets struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	MarketId      int             `gorm:"column:market_id" json:"market_id"`
	MarketCode    string          `gorm:"column:market_code" json:"market_code"`
	MarketName    string          `gorm:"column:market_name" json:"market_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Markets) TableName() string {
	return "mst.m_market"
}

type Districts struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	DistrictId    int             `gorm:"column:district_id" json:"district_id"`
	DistrictCode  string          `gorm:"column:district_code" json:"district_code"`
	DistrictName  string          `gorm:"column:district_name" json:"district_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Districts) TableName() string {
	return "mst.m_district"
}

type Industries struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	IndustryId    int             `gorm:"column:industry_id" json:"industry_id"`
	IndustryCode  string          `gorm:"column:industry_code" json:"industry_code"`
	IndustryName  string          `gorm:"column:industry_name" json:"industry_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Industries) TableName() string {
	return "mst.m_industry"
}

type Outlets struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OutletId      int             `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string          `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string          `gorm:"column:outlet_name" json:"outlet_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Outlets) TableName() string {
	return "mst.m_outlet"
}

type Gamification struct {
	CustID             string          `gorm:"column:cust_id" json:"cust_id"`
	GamificationNo     string          `gorm:"column:gamification_no;primaryKey" json:"gamification_no"`
	GamificationDate   time.Time       `gorm:"column:gamification_date" json:"gamification_date"`
	Title              string          `gorm:"column:title" json:"title"`
	Description        *string         `gorm:"column:description" json:"description"`
	EmpGrpID           int64           `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	ProductFilter      string          `gorm:"column:product_filter" json:"product_filter"`
	OutletFilter       string          `gorm:"column:outlet_filter" json:"outlet_filter"`
	SourceId           int64           `gorm:"column:source_id" json:"source_id"`
	MeasurementId      int64           `gorm:"column:measurement_id" json:"measurement_id"`
	StartDate          time.Time       `gorm:"column:start_date" json:"start_date"`
	EndDate            time.Time       `gorm:"column:end_date" json:"end_date"`
	FinishedDate       *time.Time      `gorm:"column:finished_date" json:"finished_date"`
	AnnouncementId     int64           `gorm:"column:announcement_id" json:"announcement_id"`
	SubAnnouncementId  int64           `gorm:"column:subannouncement_id" json:"subannouncement_id"`
	AnnouncementTarget float64         `gorm:"column:announcement_target" json:"announcement_target"`
	GamificationStatus int64           `gorm:"column:gamification_status" json:"gamification_status"`
	CreatedBy          int64           `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          int64           `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time       `gorm:"column:updated_at" json:"updated_at"`
	IsDel              bool            `gorm:"column:is_del" json:"is_del"`
	DeletedBy          *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt          *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Gamification) TableName() string {
	return "sls.gamification"
}

type GamificationNo struct {
	GamificationNo string `gorm:"column:get_no_fn"`
}

func (m *Gamification) BeforeCreate(trx *gorm.DB) (err error) {
	var gamificationNo GamificationNo
	trCode := "GM"
	gamificatioDateStr := m.GamificationDate.Format("2006-01-02")
	gamificatioDateSubtr := gamificatioDateStr[2:4] + gamificatioDateStr[5:7] + gamificatioDateStr[8:10]

	queryStr := fmt.Sprintf(`SELECT
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(gamification_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn
	FROM sls.gamification
	WHERE substr(gamification_no,3,6) = '%v' AND cust_id = '%v'`, gamificatioDateSubtr, strings.ToUpper(m.CustID))
	err = trx.Raw(queryStr).Scan(&gamificationNo).Error
	if err != nil {
		return err
	}

	m.GamificationNo = trCode + gamificatioDateSubtr + gamificationNo.GamificationNo
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

func (m *Gamification) BeforeUpdate(trx *gorm.DB) (err error) {
	// m.UpdatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now()

	return nil
}
