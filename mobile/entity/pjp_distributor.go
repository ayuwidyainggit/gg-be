package entity

type SubmitPjpDistributorRequest struct {
	ApprovalStatus string        `json:"approval_status" validate:"required"`
	SalesmanId     int64         `json:"salesman_id" validate:"required"`
	SalesmanCode   int64         `json:"salesman_code" validate:"required"`
	SalesmanName   string        `json:"salesman_name" validate:"required"`
	SalesmanTeam   string        `json:"salesman_team" validate:"required"`
	WarehouseId    int64         `json:"warehouse_id" validate:"required"`
	OperationType  string        `json:"operation_type" validate:"required"`
	PjpCode        int64         `json:"pjp_code" validate:"required"`
	IsActive       bool          `json:"is_active" validate:"required"`
	Status         bool          `json:"status" validate:"required"`
	VisitDay       []PjpVisitDay `json:"visit_day" validate:"required,max=7,dive"`
}

type PjpVisitDay struct {
	Day             string       `json:"day" validate:"required"`
	IndexDay        int          `json:"index_day" validate:"required"`
	Date            FlexibleTime `json:"date" validate:"required"`
	Week            int          `json:"week" validate:"required"`
	Year            int          `json:"year" validate:"required"`
	StartWeek       FlexibleTime `json:"start_week" validate:"required"`
	IsInCurrentYear bool         `json:"is_in_current_year" validate:"required"`
	Visit           PjpVisit     `json:"visit" validate:"required"`
}

type PjpVisit struct {
	// RouteId            FlexibleRouteId        `json:"route_id" validate:"required"`
	RouteName          string                 `json:"route_name" validate:"required"`
	OutletDestinations []PjpOutletDestination `json:"outlet_destinations" validate:"required,dive"`
}

type PjpOutletDestination struct {
	OutletId     int64    `json:"outlet_id" validate:"required"`
	OutletName   string   `json:"outlet_name" validate:"required"`
	OutletCode   string   `json:"outlet_code"`
	Address1     string   `json:"address1"`
	OutletStatus string   `json:"outlet_status"`
	Latitude     string   `json:"latitude" validate:"latitude"`
	Longitude    string   `json:"longitude" validate:"longitude"`
	AvgSalesWeek *float64 `json:"avg_sales_week"`
}

type UpdatePJPDistributorRequest struct {
	PJPCode        int64                          `params:"pjp_code" validate:"required"`
	CustomerID     string                         `json:"customer_id" validate:"required"`
	ApprovalStatus string                         `json:"approval_status" validate:"required"`
	SalesmanID     int64                          `json:"salesman_id" validate:"required"`
	SalesmanCode   int64                          `json:"salesman_code" validate:"required"`
	SalesmanName   string                         `json:"salesman_name" validate:"required"`
	SalesmanTeam   string                         `json:"salesman_team" validate:"required"`
	WarehouseID    int64                          `json:"warehouse_id" validate:"required"`
	OperationType  string                         `json:"operation_type" validate:"required"`
	IsActive       bool                           `json:"is_active"`
	Status         bool                           `json:"status"`
	VisitDay       []UpdatePJPDistributorVisitDay `json:"visit_day" validate:"dive"`
}

type UpdatePJPDistributorVisitDay struct {
	Day             string                          `json:"day" validate:"required,oneof=Sun Mon Tue Wed Thu Fri Sat"` // e.g., "Sun", "Mon"
	IndexDay        int                             `json:"index_day" validate:"required,gte=1,lte=7"`                 // 1-7
	Date            FlexibleTime                    `json:"date" validate:"required"`                                  // Using string for "YYYY-MM-DD"
	Week            int                             `json:"week" validate:"required"`
	Year            int                             `json:"year" validate:"required"`
	StartWeek       FlexibleTime                    `json:"start_week" validate:"required"`
	IsInCurrentYear bool                            `json:"is_in_current_year" validate:"required"`
	Visit           UpdatePJPDistributorVisitDetail `json:"visit"`
}

type UpdatePJPDistributorVisitDetail struct {
	RouteCode          int64                                   `json:"route_code" validate:"required"`
	RouteName          string                                  `json:"route_name" validate:"required"`
	OutletDestinations []UpdatePJPDistributorOutletDestination `json:"outlet_destinations" validate:"required"`
}

type UpdatePJPDistributorOutletDestination struct {
	OutletID     int64    `json:"outlet_id" validate:"required"`
	OutletName   string   `json:"outlet_name" validate:"required"`
	OutletCode   string   `json:"outlet_code"`
	Address1     string   `json:"address1"`
	OutletStatus string   `json:"outlet_status"`
	Latitude     string   `json:"latitude"`
	Longitude    string   `json:"longitude"`
	AvgSalesWeek *float64 `json:"avg_sales_week"` // mapped from numeric(10,2)
}

type SalesmanAndTeamResult struct {
	EmpID         int64  `gorm:"column:emp_id"`
	EmpCode       string `gorm:"column:emp_code"`
	SalesName     string `gorm:"column:sales_name"`
	SalesTeamName string `gorm:"column:sales_team_name"`
}

type RoutePopDailyResult struct {
	RouteCode string `gorm:"column:route_code"`
	PJPID     int64  `gorm:"column:pjp_id"`
	PJPCode   string `gorm:"column:pjp_code"`
}
