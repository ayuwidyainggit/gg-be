package response

type RouteResponse struct {
	ID           int                    `json:"id"`
	RouteCode    int                    `json:"route_code"`
	RouteName    string                 `json:"route_name"`
	TotalOutlet  int                    `json:"total_outlet"`
	Destinations *[]DestinationResponse `json:"destinations"`
}

type RouteDailyResponse struct {
	ID           int                    `json:"id"`
	RouteCode    int                    `json:"route_code"`
	RouteName    string                 `json:"route_name"`
	IsAssign     bool                   `json:"is_assign"`
	IsAssignPjp  bool                   `json:"is_assign_pjp"`
	TotalOutlet  int                    `json:"total_outlet"`
	Destinations *[]OutletDailyResponse `json:"destinations"`
}

type OutletDailyResponse struct {
	DestinationID      int     `json:"outlet_id"`
	DestinationCode    string  `json:"outlet_code"`
	DestinationName    string  `json:"outlet_name"`
	Longitude          string  `json:"longitude"`
	Latitude           string  `json:"latitude"`
	DestinationStatus  int     `json:"outlet_status"`
	DestinationAddress string  `json:"outlet_address"`
	AvgSalesWeek       float64 `json:"avg_sales_week"`
	Status             string  `json:"status"`
	PjpCode            int     `json:"pjp_code"`
	PjpID              int     `json:"pjp_id"`
	RouteCode          int     `json:"route_code"`
}
