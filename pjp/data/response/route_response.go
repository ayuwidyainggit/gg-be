package response

type RouteResponse struct {
	ID          int               `json:"id"`
	RouteCode   int               `json:"route_code"`
	RouteName   string            `json:"route_name"`
	IsAssign    bool              `json:"is_assign"`
	IsAssignPjp bool              `json:"is_assign_pjp"`
	TotalOutlet int               `json:"total_outlet"`
	Outlets     *[]OutletResponse `json:"outlets"`
}

type RouteDailyResponse struct {
	ID          int                    `json:"id"`
	RouteCode   int                    `json:"route_code"`
	RouteName   string                 `json:"route_name"`
	IsAssign    bool                   `json:"is_assign"`
	IsAssignPjp bool                   `json:"is_assign_pjp"`
	TotalOutlet int                    `json:"total_outlet"`
	Outlets     *[]OutletDailyResponse `json:"outlets"`
}

type OutletDailyResponse struct {
	OutletID      int     `json:"outlet_id"`
	OutletCode    string  `json:"outlet_code"`
	OutletName    string  `json:"outlet_name"`
	Longitude     string  `json:"longitude"`
	Latitude      string  `json:"latitude"`
	OutletStatus  int     `json:"outlet_status"`
	OutletAddress string  `json:"outlet_address"`
	AvgSalesWeek  float64 `json:"avg_sales_week"`
	Status        string  `json:"status"`
	PjpCode       int     `json:"pjp_code"`
	PjpID         int     `json:"pjp_id"`
	RouteCode     int     `json:"route_code"`
}
