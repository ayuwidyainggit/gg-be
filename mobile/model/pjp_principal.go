package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// PermanentJourneyPlanPrincipal represents pjp_principles.permanent_journey_plans table
type PermanentJourneyPlanPrincipal struct {
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

func (PermanentJourneyPlanPrincipal) TableName() string {
	return "pjp_principles.permanent_journey_plans"
}

// RoutePrincipal represents pjp_principles.routes table
type RoutePrincipal struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode int64      `gorm:"column:route_code" json:"route_code"` // int8 in database
	RouteName string     `gorm:"column:route_name" json:"route_name"`
	Sequence  int64      `gorm:"column:sequence" json:"sequence"`
	CustId    string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	PjpId     int64      `gorm:"column:pjp_id" json:"pjp_id"`
}

func (RoutePrincipal) TableName() string {
	return "pjp_principles.routes"
}

// DestinationPrincipal represents pjp_principles.destinations table
type DestinationPrincipal struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode          string     `gorm:"column:route_code" json:"route_code"`
	RouteName          string     `gorm:"column:route_name" json:"route_name"`
	Status             string     `gorm:"column:status" json:"status"`
	VerifiedDate       *time.Time `gorm:"column:verified_date" json:"verified_date"`
	DestinationId      int64      `gorm:"column:destination_id" json:"destination_id"`
	DestinationCode    string     `gorm:"column:destination_code" json:"destination_code"`
	DestinationStatus  string     `gorm:"column:destination_status" json:"destination_status"`
	DestinationName    string     `gorm:"column:destination_name" json:"destination_name"`
	DestinationAddress string     `gorm:"column:destination_address" json:"destination_address"`
	DestinationType    string     `gorm:"column:destination_type" json:"destination_type"` // "outlet" or "distributor"
	Longitude          string     `gorm:"column:longitude" json:"longitude"`
	Latitude           string     `gorm:"column:latitude" json:"latitude"`
	PjpId              int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode            int64      `gorm:"column:pjp_code" json:"pjp_code"`
	OldPjpId           *int64     `gorm:"column:old_pjp_id" json:"old_pjp_id"`
	OldPjpCode         *int64     `gorm:"column:old_pjp_code" json:"old_pjp_code"`
	OldRouteCode       *string    `gorm:"column:old_route_code" json:"old_route_code"`
	OldRouteName       *string    `gorm:"column:old_route_name" json:"old_route_name"`
	Photo              *string    `gorm:"column:photo" json:"photo"`
	Signature          *string    `gorm:"column:signature" json:"signature"`
	AvgSalesWeek       *float64   `gorm:"column:avg_sales_week" json:"avg_sales_week"`
	CustId             string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DestinationPrincipal) TableName() string {
	return "pjp_principles.destinations"
}

// DestinationHistoryPrincipal represents pjp_principles.destinations_history table
type DestinationHistoryPrincipal struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RouteCode          string     `gorm:"column:route_code" json:"route_code"`
	RouteName          string     `gorm:"column:route_name" json:"route_name"`
	VerifiedDate       *time.Time `gorm:"column:verified_date" json:"verified_date"`
	Date               time.Time  `gorm:"column:date" json:"date"`
	Week               int        `gorm:"column:week" json:"week"`
	Year               int        `gorm:"column:year" json:"year"`
	IndexDay           int        `gorm:"column:index_day" json:"index_day"`
	StartWeek          time.Time  `gorm:"column:start_week" json:"start_week"`
	IsInCurrentYear    bool       `gorm:"column:is_in_current_year" json:"is_in_current_year"`
	IsAdditional       *bool      `gorm:"column:is_additional" json:"is_additional"`
	DestinationId      int64      `gorm:"column:destination_id" json:"destination_id"`
	DestinationCode    string     `gorm:"column:destination_code" json:"destination_code"`
	DestinationStatus  string     `gorm:"column:destination_status" json:"destination_status"`
	DestinationName    string     `gorm:"column:destination_name" json:"destination_name"`
	DestinationAddress string     `gorm:"column:destination_address" json:"destination_address"`
	DestinationType    string     `gorm:"column:destination_type" json:"destination_type"` // "outlet" or "distributor"
	Longitude          string     `gorm:"column:longitude" json:"longitude"`
	Latitude           string     `gorm:"column:latitude" json:"latitude"`
	PjpId              int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode            int64      `gorm:"column:pjp_code" json:"pjp_code"`
	OldPjpId           *int64     `gorm:"column:old_pjp_id" json:"old_pjp_id"`
	OldPjpCode         *int64     `gorm:"column:old_pjp_code" json:"old_pjp_code"`
	OldRouteCode       *string    `gorm:"column:old_route_code" json:"old_route_code"`
	OldRouteName       *string    `gorm:"column:old_route_name" json:"old_route_name"`
	Photo              *string    `gorm:"column:photo" json:"photo"`
	Signature          *string    `gorm:"column:signature" json:"signature"`
	AvgSalesWeek       *float64   `gorm:"column:avg_sales_week" json:"avg_sales_week"`
	CustId             string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DestinationHistoryPrincipal) TableName() string {
	return "pjp_principles.destinations_history"
}

// RoutePopPermanentPrincipal represents pjp_principles.route_pop_permanent table
type RoutePopPermanentPrincipal struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Year      int        `gorm:"column:year" json:"year"`
	Week      int        `gorm:"column:week" json:"week"`
	Date      time.Time  `gorm:"column:date" json:"date"`
	Day       string     `gorm:"column:day" json:"day"`
	RouteCode int64      `gorm:"column:route_code" json:"route_code"`
	PjpId     int64      `gorm:"column:pjp_id" json:"pjp_id"`
	PjpCode   int64      `gorm:"column:pjp_code" json:"pjp_code"`
	CustId    string     `gorm:"column:cust_id" json:"cust_id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (RoutePopPermanentPrincipal) TableName() string {
	return "pjp_principles.route_pop_permanent"
}

func (r *RoutePrincipal) BeforeCreate(tx *gorm.DB) (err error) {
	if r.RouteCode != 0 {
		return nil
	}

	var lastCode int64
	err = tx.Table("pjp_principles.routes AS r").
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
