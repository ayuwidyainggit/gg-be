package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PermanentJourneyPlan represents pjp.permanent_journey_plans table
type PermanentJourneyPlan struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PjpCode        int64      `gorm:"column:pjp_code" json:"pjp_code"`
	OperationType  string     `gorm:"column:operation_type" json:"operation_type"`
	TeamSalesman   string     `gorm:"column:team_salesman" json:"team_salesman"`
	SalesmanId     int64      `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName   string     `gorm:"column:salesman_name" json:"salesman_name"`
	WarehouseId    int64      `gorm:"column:warehouse_id" json:"warehouse_id"`
	WarehouseName  *string    `gorm:"column:warehouse_name" json:"warehouse_name"`
	PjpMode        string     `gorm:"column:pjp_mode" json:"pjp_mode"`
	Status         string     `gorm:"column:status" json:"status"`
	CustId         string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	SalesmanCode   string     `gorm:"column:salesman_code" json:"salesman_code"`
	ApprovalStatus *string    `gorm:"column:approval_status" json:"approval_status"`
}

func (PermanentJourneyPlan) TableName() string {
	return "pjp.permanent_journey_plans"
}

// Route represents pjp.routes table
type Route struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode int64      `gorm:"column:route_code" json:"route_code"` // int8 in database
	RouteName string     `gorm:"column:route_name" json:"route_name"`
	Sequence  int64      `gorm:"column:sequence" json:"sequence"`
	IsAssign  *bool      `gorm:"column:is_assign" json:"is_assign"` // Fixed: is_assign not is_assigne
	CustId    string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	PjpId     int64      `gorm:"column:pjp_id" json:"pjp_id"`
	IsPjpOld  *bool      `gorm:"column:is_pjp_old" json:"is_pjp_old"`
}

func (Route) TableName() string {
	return "pjp.routes" // Fixed: plural form
}

// RouteOutlet represents pjp.route_outlet table
type RouteOutlet struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode     string     `gorm:"column:route_code" json:"route_code"`
	RouteName     string     `gorm:"column:route_name" json:"route_name"`
	OutletId      int64      `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string     `gorm:"column:outlet_name" json:"outlet_name"`
	Longitude     string     `gorm:"column:longitude" json:"longitude"`
	Latitude      string     `gorm:"column:latitude" json:"latitude"`
	OutletStatus  string     `gorm:"column:outlet_status" json:"outlet_status"` // Changed to string (varchar in DB)
	OutletAddress string     `gorm:"column:outlet_address" json:"outlet_address"`
	PjpId         int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode       int64      `gorm:"column:pjp_code" json:"pjp_code"`
	CustId        string     `gorm:"column:cust_id" json:"cust_id"`
	Status        string     `gorm:"column:status" json:"status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
	VerifiedDate  *time.Time `gorm:"column:verified_date" json:"verified_date"`
	OldPjpId      *int64     `gorm:"column:old_pjp_id" json:"old_pjp_id"`
	OldPjpCode    *int64     `gorm:"column:old_pjp_code" json:"old_pjp_code"`
	OldRouteCode  *string    `gorm:"column:old_route_code" json:"old_route_code"`
	OldRouteName  *string    `gorm:"column:old_route_name" json:"old_route_name"`
	Photo         *string    `gorm:"column:photo" json:"photo"`
	Signature     *string    `gorm:"column:signature" json:"signature"`
	AvgSalesWeek  *float64   `gorm:"column:avg_sales_week" json:"avg_sales_week"`
}

type RouteOutletInUse struct {
	Total     int    `gorm:"column:total" json:"total"`
	RouteCode string `gorm:"column:route_code" json:"route_code"`
}

func (RouteOutlet) TableName() string {
	return "pjp.route_outlet"
}

// RouteOutletHistory represents pjp.route_outlet_history table
type RouteOutletHistory struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode       string     `gorm:"column:route_code" json:"route_code"`
	RouteName       string     `gorm:"column:route_name" json:"route_name"`
	OutletId        int64      `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode      string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      string     `gorm:"column:outlet_name" json:"outlet_name"`
	Longitude       string     `gorm:"column:longitude" json:"longitude"`
	Latitude        string     `gorm:"column:latitude" json:"latitude"`
	OutletStatus    string     `gorm:"column:outlet_status" json:"outlet_status"`
	OutletAddress   string     `gorm:"column:outlet_address" json:"outlet_address"`
	PjpId           int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode         int64      `gorm:"column:pjp_code" json:"pjp_code"`
	CustId          string     `gorm:"column:cust_id" json:"cust_id"`
	Status          string     `gorm:"column:status" json:"status"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
	VerifiedDate    *time.Time `gorm:"column:verified_date" json:"verified_date"`
	OldPjpId        *int64     `gorm:"column:old_pjp_id" json:"old_pjp_id"`
	OldPjpCode      *int64     `gorm:"column:old_pjp_code" json:"old_pjp_code"`
	OldRouteCode    *string    `gorm:"column:old_route_code" json:"old_route_code"`
	OldRouteName    *string    `gorm:"column:old_route_name" json:"old_route_name"`
	Photo           *string    `gorm:"column:photo" json:"photo"`
	Signature       *string    `gorm:"column:signature" json:"signature"`
	AvgSalesWeek    *float64   `gorm:"column:avg_sales_week" json:"avg_sales_week"`
	IndexDay        int        `gorm:"column:index_day" json:"index_day"`
	StartWeek       time.Time  `gorm:"column:start_week" json:"start_week"`
	IsInCurrentYear bool       `gorm:"column:is_in_current_year" json:"is_in_current_year"`
	Week            int        `gorm:"column:week" json:"week"`
	Year            int        `gorm:"column:year" json:"year"`
	Date            time.Time  `gorm:"column:date" json:"date"`
	IsAdditional    *bool      `gorm:"column:is_additional" json:"is_additional"`
}

func (RouteOutletHistory) TableName() string {
	return "pjp.route_outlet_history"
}

// RoutePopPermanent represents pjp.route_pop_permanent table
type RoutePopPermanent struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Year      int        `gorm:"column:year" json:"year"`
	Week      int        `gorm:"column:week" json:"week"`
	Date      time.Time  `gorm:"column:date" json:"date"`
	Day       string     `gorm:"column:day" json:"day"`
	RouteCode string     `gorm:"column:route_code" json:"route_code"`
	PjpId     int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode   int64      `gorm:"column:pjp_code" json:"pjp_code"`
	CustId    string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (RoutePopPermanent) TableName() string {
	return "pjp.route_pop_permanent"
}

func (r *Route) BeforeCreate(tx *gorm.DB) (err error) {
	if r.RouteCode != 0 {
		return nil
	}

	var lastCode int64
	// Using existing repository logic to ensure consistency
	err = tx.Table("pjp.routes AS r").
		Select("r.route_code").
		Order("r.route_code DESC").
		Limit(1).
		Scan(&lastCode).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	r.RouteCode = lastCode + 1
	return nil
}
