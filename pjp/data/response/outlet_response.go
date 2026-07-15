package response

type OutletResponse struct {
	OutletID      int     `json:"outlet_id"`
	OutletCode    string  `json:"outlet_code"`
	OutletName    string  `json:"outlet_name"`
	Longitude     string  `json:"longitude"`
	Latitude      string  `json:"latitude"`
	OutletStatus  string  `json:"outlet_status"`
	OutletAddress string  `json:"outlet_address"`
	AvgSalesWeek  float64 `json:"avg_sales_week"`
	Status        string  `json:"status"`
	PjpCode       int     `json:"pjp_code"`
	PjpID         int     `json:"pjp_id"`
	RouteCode     int     `json:"route_code"`
}

type DailyOutletResponse struct {
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
