package response

import "scyllax-pjp/model"

type PjpResponse struct {
	ID            int    `json:"id"`
	PjpCode       string `json:"pjp_code"`
	OperationType string `json:"operation_type"`
	// EmployeeCode      string         `json:"emp_code"`
	TeamSalesMan      string         `json:"team_salesman"`
	SalesManID        int            `json:"salesman_id"`
	SalesmanName      string         `json:"salesman_name"`
	SalesmanCode      string         `json:"salesman_code"`
	WarehouseID       int            `json:"warehouse_id"`
	WarehouseName     string         `json:"warehouse_name"`
	PjpMode           string         `json:"pjp_mode"`
	Status            bool           `json:"status"`
	TotalRoute        int            `json:"total_route"`
	ApprovalStatus    string         `json:"approval_status"`
	TotalDestinations int            `json:"total_destinations"`
	TotalOutlet       int            `json:"total_outlet"`
	TotalDistributor  int            `json:"total_distributor"`
	Route             []RoutesEntity `json:"routes,omitempty"`
}

type RoutesEntity struct {
	RouteCode   int `json:"route_code"`
	TotalOutlet int `json:"total_outlet"`
}

type PjpAutoResponse struct {
	Day   []string                `json:"Day"`
	Route [][]DestinationResponse `json:"Route"`
	Sales [][]PjpResponse         `json:"Sales"`
}

type OperationTypeResponse struct {
	OperationType string `json:"operation_type"`
}

type SalesmanTeamResponse struct {
	TeamSalesMan string `json:"team_salesman"`
}

type PjpWithRouteRow struct {
	model.Pjp
	RouteCode   int `gorm:"column:route_code"`
	TotalOutlet int `gorm:"column:total_outlet"`
}
