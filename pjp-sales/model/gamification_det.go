package model

import (
	"time"
)

type GamificationParticipant struct {
	CustID                    string    `gorm:"column:cust_id" json:"cust_id"`
	GamificationParticipantID int64     `gorm:"column:gamification_participant_id;primaryKey" json:"gamification_participant_id"`
	GamificationNo            string    `gorm:"column:gamification_no" json:"gamification_no"`
	EmpID                     int64     `gorm:"column:emp_id" json:"emp_id"`
	CreatedBy                 int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt                 time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GamificationParticipant) TableName() string {
	return "sls.gamification_participants"
}

type GamificationProduct struct {
	CustID                string    `gorm:"column:cust_id" json:"cust_id"`
	GamificationProductID int64     `gorm:"column:gamification_product_id;primaryKey" json:"gamification_product_id"`
	GamificationNo        string    `gorm:"column:gamification_no" json:"gamification_no"`
	ProID                 int64     `gorm:"column:pro_id" json:"pro_id"`
	CreatedBy             int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GamificationProduct) TableName() string {
	return "sls.gamification_products"
}

type GamificationOutlet struct {
	CustID               string    `gorm:"column:cust_id" json:"cust_id"`
	GamificationOutletID int64     `gorm:"column:gamification_outlet_id;primaryKey" json:"gamification_outlet_id"`
	GamificationNo       string    `gorm:"column:gamification_no" json:"gamification_no"`
	OutletID             int64     `gorm:"column:outlet_id" json:"outlet_id"`
	CreatedBy            int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GamificationOutlet) TableName() string {
	return "sls.gamification_outlets"
}

type GamificationAnnouncement struct {
	CustID                     string    `gorm:"column:cust_id" json:"cust_id"`
	GamificationAnnouncementID int64     `gorm:"column:gamification_announcement_id;primaryKey" json:"gamification_announcement_id"`
	GamificationNo             string    `gorm:"column:gamification_no" json:"gamification_no"`
	AnnouncementDate           time.Time `gorm:"column:announcement_date" json:"announcement_date"`
	CreatedBy                  int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt                  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GamificationAnnouncement) TableName() string {
	return "sls.gamification_announcements"
}

type GamificationParticipantRead struct {
	GamificationParticipantID int64   `gorm:"column:gamification_participant_id;primaryKey" json:"gamification_participant_id"`
	SalesmanID                int64   `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode              *string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName              *string `gorm:"column:salesman_name" json:"salesman_name"`
}

func (GamificationParticipantRead) TableName() string {
	return "sls.gamification_participants"
}

type GamificationProductRead struct {
	GamificationProductID int64   `gorm:"column:gamification_product_id;primaryKey" json:"gamification_product_id"`
	ProID                 int64   `gorm:"column:pro_id" json:"pro_id"`
	ProCode               *string `gorm:"column:pro_code" json:"pro_code"`
	ProName               *string `gorm:"column:pro_name" json:"pro_name"`
}

func (GamificationProductRead) TableName() string {
	return "sls.gamification_products"
}

type GamificationOutletRead struct {
	GamificationOutletID int64   `gorm:"column:gamification_outlet_id;primaryKey" json:"gamification_outlet_id"`
	OutletID             int64   `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode           *string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName           *string `gorm:"column:outlet_name" json:"outlet_name"`
}

func (GamificationOutletRead) TableName() string {
	return "sls.gamification_outlets"
}

type GamificationAnnouncementRead struct {
	GamificationAnnouncementID int64     `gorm:"column:gamification_announcement_id;primaryKey" json:"gamification_announcement_id"`
	AnnouncementDate           time.Time `gorm:"column:announcement_date" json:"announcement_date"`
}

func (GamificationAnnouncementRead) TableName() string {
	return "sls.gamification_announcements"
}

type GamificationRankingRead struct {
	SalesmanID   int64   `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode *string `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName *string `gorm:"column:salesman_name" json:"salesman_name"`
	Total        float64 `gorm:"column:total" json:"total"`
}

func (GamificationRankingRead) TableName() string {
	return "sls.gamification_participants"
}

type GamificationFilter struct {
	GamificationFilterID int64  `gorm:"column:gamification_filter_id;primaryKey" json:"gamification_filter_id"`
	GamificationNo       string `gorm:"column:gamification_no" json:"gamification_no"`
	SalesTeam            string `gorm:"column:sales_team" json:"sales_team"`
	OperationType        string `gorm:"column:operation_type" json:"operation_type"`
	Warehouse            string `gorm:"column:warehouse" json:"warehouse"`
	ProductCategory      string `gorm:"column:product_category" json:"product_category"`
	ProductLine          string `gorm:"column:product_line" json:"product_line"`
	Brand                string `gorm:"column:brand" json:"brand"`
	SubBrands1           string `gorm:"column:sub_brand1" json:"sub_brand1"`
	SubBrands2           string `gorm:"column:sub_brand2" json:"sub_brand2"`
	OutletLocation       string `gorm:"column:outlet_location" json:"outlet_location"`
	OutletGroup          string `gorm:"column:outlet_group" json:"outlet_group"`
	OutletClass          string `gorm:"column:outlet_class" json:"outlet_class"`
	Market               string `gorm:"column:market" json:"market"`
	District             string `gorm:"column:district" json:"district"`
	Industry             string `gorm:"column:industry" json:"industry"`
}

func (GamificationFilter) TableName() string {
	return "sls.gamification_filters"
}

type GamificationFilterRead struct {
	SalesTeam       string `gorm:"column:sales_team" json:"sales_team"`
	OperationType   string `gorm:"column:operation_type" json:"operation_type"`
	Warehouse       string `gorm:"column:warehouse" json:"warehouse"`
	ProductCategory string `gorm:"column:product_category" json:"product_category"`
	ProductLine     string `gorm:"column:product_line" json:"product_line"`
	Brand           string `gorm:"column:brand" json:"brand"`
	SubBrand1       string `gorm:"column:sub_brand1" json:"sub_brand1"`
	SubBrand2       string `gorm:"column:sub_brand2" json:"sub_brand2"`
	OutletLocation  string `gorm:"column:outlet_location" json:"outlet_location"`
	OutletGroup     string `gorm:"column:outlet_group" json:"outlet_group"`
	OutletClass     string `gorm:"column:outlet_class" json:"outlet_class"`
	Market          string `gorm:"column:market" json:"market"`
	District        string `gorm:"column:district" json:"district"`
	Industry        string `gorm:"column:industry" json:"industry"`
}

func (GamificationFilterRead) TableName() string {
	return "sls.gamification_filters"
}
