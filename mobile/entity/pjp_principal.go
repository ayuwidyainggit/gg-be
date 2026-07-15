package entity

type SubmitPjpPrincipalRequest struct {
	ApprovalStatus string                 `json:"approval_status" validate:"required"`
	SalesmanId     int64                  `json:"salesman_id" validate:"required"`
	SalesmanCode   int64                  `json:"salesman_code" validate:"required"`
	SalesmanName   string                 `json:"salesman_name" validate:"required"`
	SalesmanTeam   string                 `json:"salesman_team" validate:"required"`
	WarehouseId    int64                  `json:"warehouse_id" validate:"required"`
	OperationType  string                 `json:"operation_type" validate:"required"`
	PjpCode        int64                  `json:"pjp_code" validate:"required"`
	IsActive       bool                   `json:"is_active" validate:"required"`
	Status         bool                   `json:"status" validate:"required"`
	VisitDay       []PjpPrincipalVisitDay `json:"visit_day" validate:"required,max=7,dive"`
}

type PjpPrincipalVisitDay struct {
	Day             string            `json:"day" validate:"required"`
	IndexDay        int               `json:"index_day" validate:"required"`
	Date            FlexibleTime      `json:"date" validate:"required"`
	Week            int               `json:"week" validate:"required"`
	Year            int               `json:"year" validate:"required"`
	StartWeek       FlexibleTime      `json:"start_week" validate:"required"`
	IsInCurrentYear bool              `json:"is_in_current_year" validate:"required"`
	Visit           PjpPrincipalVisit `json:"visit" validate:"required"`
}

type PjpPrincipalVisit struct {
	RouteName               string                      `json:"route_name" validate:"required"`
	OutletDestinations      []PjpOutletDestination      `json:"outlet_destinations" validate:"required_without=DistributorDestinations,dive"`
	DistributorDestinations []PjpDistributorDestination `json:"distributor_destinations" validate:"required_without=OutletDestinations,dive"`
}

type PjpDistributorDestination struct {
	DistributorId   int64  `json:"distributor_id" validate:"required"`
	DistributorName string `json:"distributor_name" validate:"required"`
	DistributorCode string `json:"distributor_code" validate:"required"`
	Address         string `json:"address" validate:"required"`
	Latitude        string `json:"latitude" validate:"latitude"`
	Longitude       string `json:"longitude" validate:"longitude"`
}

type UpdatePjpPrincipalRequest struct {
	PJPCode        int64                        `params:"pjp_code" validate:"required"`
	CustomerID     string                       `json:"customer_id" validate:"required"`
	ApprovalStatus string                       `json:"approval_status" validate:"required"`
	SalesmanID     int64                        `json:"salesman_id" validate:"required"`
	SalesmanCode   int64                        `json:"salesman_code" validate:"required"`
	SalesmanName   string                       `json:"salesman_name" validate:"required"`
	SalesmanTeam   string                       `json:"salesman_team" validate:"required"`
	WarehouseID    int64                        `json:"warehouse_id" validate:"required"`
	OperationType  string                       `json:"operation_type" validate:"required"`
	IsActive       bool                         `json:"is_active"`
	Status         bool                         `json:"status"`
	VisitDay       []UpdatePjpPrincipalVisitDay `json:"visit_day" validate:"dive"`
}

type UpdatePjpPrincipalVisitDay struct {
	Day             string                        `json:"day" validate:"required,oneof=Sun Mon Tue Wed Thu Fri Sat"`
	IndexDay        int                           `json:"index_day" validate:"required,gte=1,lte=7"`
	Date            FlexibleTime                  `json:"date" validate:"required"`
	Week            int                           `json:"week" validate:"required"`
	Year            int                           `json:"year" validate:"required"`
	StartWeek       FlexibleTime                  `json:"start_week" validate:"required"`
	IsInCurrentYear bool                          `json:"is_in_current_year" validate:"required"`
	Visit           UpdatePjpPrincipalVisitDetail `json:"visit"`
}

type UpdatePjpPrincipalVisitDetail struct {
	RouteCode               int64                                      `json:"route_code" validate:"required"`
	RouteName               string                                     `json:"route_name" validate:"required"`
	OutletDestinations      []UpdatePjpPrincipalOutletDestination      `json:"outlet_destinations" validate:"required_without=DistributorDestinations,dive"`
	DistributorDestinations []UpdatePjpPrincipalDistributorDestination `json:"distributor_destinations" validate:"required_without=OutletDestinations,dive"`
}

type UpdatePjpPrincipalOutletDestination struct {
	OutletID     int64    `json:"outlet_id" validate:"required"`
	OutletName   string   `json:"outlet_name" validate:"required"`
	OutletCode   string   `json:"outlet_code"`
	Address1     string   `json:"address1"`
	OutletStatus string   `json:"outlet_status"`
	Latitude     string   `json:"latitude"`
	Longitude    string   `json:"longitude"`
	AvgSalesWeek *float64 `json:"avg_sales_week"`
}

type UpdatePjpPrincipalDistributorDestination struct {
	DistributorID   int64  `json:"distributor_id" validate:"required"`
	DistributorName string `json:"distributor_name" validate:"required"`
	DistributorCode string `json:"distributor_code"`
	Address         string `json:"address"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
}

type RoutePopDailyPrincipalResult struct {
	RouteCode string `gorm:"column:route_code"`
	PJPID     int64  `gorm:"column:pjp_id"`
	PJPCode   string `gorm:"column:pjp_code"`
}
