package response

import "time"

type PjpEnhanceResponse struct {
	Pjp      Pjp
	Routes   []Routes
	VisitDay []VisitDay `json:"visit_day"`
}

type VisitDay struct {
	ID                   int           `json:"id"`
	Day                  string        `json:"day" example:"Senin"`
	IndexDay             int           `json:"indexDay" example:"1"`
	Week                 int           `json:"week" example:"1"`
	WorkingDayCalendarID *int64        `json:"working_day_calendar_id"`
	StartWeek            string        `json:"startWeek" example:"2022-01-01"`
	Year                 int           `json:"year" example:"2022"`
	Date                 string        `json:"date" example:"2022-01-01"`
	IsInCurrentYear      bool          `json:"isInCurrentYear" example:"true"`
	Visit                RoutesHistory `json:"visit"`
}

type OutletHistory struct {
	ID              int        `gorm:"type:int;primary_key" json:"id"`
	RouteCode       int        `gorm:"column:route_code;type:int;not null" json:"route_code"`
	RouteName       string     `gorm:"column:route_name;type:varchar(125);not null" json:"route_name"`
	OutletID        int        `gorm:"column:outlet_id;type:int" json:"outlet_id"`
	OutletCode      string     `gorm:"column:outlet_code;type:varchar(125)" json:"outlet_code"`
	OutletName      string     `gorm:"column:outlet_name;type:varchar(125)" json:"outlet_name"`
	Longitude       string     `gorm:"column:longitude;type:varchar(125)" json:"longitude"`
	Latitude        string     `gorm:"column:latitude;type:varchar(125)" json:"latitude"`
	OutletStatus    string     `gorm:"column:outlet_status;type:varchar(125)" json:"outlet_status"`
	OutletAddress   string     `gorm:"column:outlet_address;type:varchar(125);null" json:"outlet_address"`
	AvgSalesWeek    float64    `gorm:"column:avg_sales_week;type:numeric(10,2)" json:"avg_sales_week"`
	PjpID           *int       `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode         *int       `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	Status          string     `gorm:"column:status;type:varchar(125);default:pending" json:"status"`
	CustID          string     `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	Year            int        `gorm:"column:year;null" json:"year"`
	Week            int        `gorm:"column:week;null" json:"week"`
	IndexDay        int        `gorm:"column:index_day;null" json:"indexDay" example:"1"`
	StartWeek       *time.Time `gorm:"column:start_week" json:"startWeek" example:"2022-01-01"`
	IsInCurrentYear bool       `gorm:"column:is_in_current_year;" json:"isInCurrentYear" example:"true"`
	IsAdditional    bool       `gorm:"column:is_additional;" json:"is_additional" example:"true"`
}

type Outlets struct {
	ID            int     `gorm:"type:int;primary_key" json:"id"`
	RouteCode     int     `gorm:"column:route_code;type:int;not null" json:"route_code"`
	RouteName     string  `gorm:"column:route_name;type:varchar(125);not null" json:"route_name"`
	OutletID      int     `gorm:"column:outlet_id;type:int" json:"outlet_id"`
	OutletCode    string  `gorm:"column:outlet_code;type:varchar(125)" json:"outlet_code"`
	OutletName    string  `gorm:"column:outlet_name;type:varchar(125)" json:"outlet_name"`
	Longitude     string  `gorm:"column:longitude;type:varchar(125)" json:"longitude"`
	Latitude      string  `gorm:"column:latitude;type:varchar(125)" json:"latitude"`
	OutletStatus  string  `gorm:"column:outlet_status;type:varchar(125)" json:"outlet_status"`
	OutletAddress string  `gorm:"column:outlet_address;type:varchar(125);null" json:"outlet_address"`
	AvgSalesWeek  float64 `gorm:"column:avg_sales_week;type:numeric(10,2)" json:"avg_sales_week"`
	PjpID         *int    `gorm:"column:pjp_id;type:int;null" json:"pjp_id"`
	PjpCode       *int    `gorm:"column:pjp_code;type:int;null" json:"pjp_code"`
	Status        string  `gorm:"column:status;type:varchar(125);default:pending" json:"status"`
	CustID        string  `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	Type          string  `json:"type"`
}

type Routes struct {
	ID        int    `gorm:"column:id;type:int;primary_key" json:"id"`
	RouteCode int    `gorm:"column:route_code;type:int;uniqueIndex;not null" json:"route_code"`
	RouteName string `gorm:"column:route_name;type:varchar(125);unique;not null" json:"route_name"`
	// IsAssign       bool      `gorm:"column:is_assign;type:bool;default:false" json:"is_assign"`
	// IsAssignPjp    bool      `gorm:"->" json:"is_assign_pjp"`
	// RoutePopStatus string    `gorm:"->" json:"route_pop_status"`
	CustID  string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	Outlets []Outlets `json:"destinations"`
}
type RoutesHistory struct {
	RouteCode int             `gorm:"column:route_code;type:int;uniqueIndex;not null" json:"route_code"`
	RouteName string          `gorm:"column:route_name;type:varchar(125);unique;not null" json:"route_name"`
	CustID    string          `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	Outlets   []OutletHistory `json:"destinations"`
}

type Pjp struct {
	ID             int    `gorm:"type:int;primary_key" json:"id"`
	PjpCode        string `gorm:"column:pjp_code;type:int;uniqueIndex;not null" json:"pjp_code"`
	OperationType  string `gorm:"column:operation_type;type:varchar(125);not null" json:"operation_type"`
	TeamSalesMan   string `gorm:"column:team_salesman;type:varchar(125)" json:"team_salesman"`
	SalesManID     int    `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName   string `gorm:"column:salesman_name;type:varchar(125)" json:"salesman_name"`
	WarehouseID    int    `gorm:"column:warehouse_id" json:"warehouse_id"`
	WarehouseName  string `gorm:"column:warehouse_name;type:varchar(125)" json:"warehouse_name"`
	SalesmanCode   string `gorm:"column:salesman_code" json:"salesman_code"`
	PjpMode        string `gorm:"column:pjp_mode;type:varchar(125);null" json:"pjp_mode"`
	EmpCode        string `gorm:"->" json:"emp_code"`
	Status         string `gorm:"column:status;type:varchar(125)" json:"status"`
	ApprovalStatus string `gorm:"type:varchar(32);not null;default:'Draft'" validate:"required,oneof='In Review' Draft Approved 'Approved With Propose' Reject" json:"approval_status"`
	CustID         string `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	RouteCode      int    `gorm:"->" json:"route_code"`
}
