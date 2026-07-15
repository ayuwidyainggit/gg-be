package response

type OutletMap struct {
	OutletID      int     `json:"outlet_id"`
	OutletCode    string  `json:"outlet_code"`
	OutletName    string  `json:"outlet_name"`
	Longitude     string  `json:"longitude"`
	Latitude      string  `json:"latitude"`
	OutletStatus  string  `json:"outlet_status"`
	OutletAddress string  `json:"outlet_address"`
	AvgSalesWeek  float64 `json:"avg_sales_week"`
}

type RouteMap struct {
	RouteCode *int        `json:"route_code"`
	RouteName string      `json:"route_name"`
	Outlets   []OutletMap `json:"outlets"`
}

type RoutePopDailyResponse struct {
	ID           int        `json:"id"`
	Routes       []RouteMap `json:"routes"`
	PjpID        *int       `json:"pjp_id"`
	PjpCode      *int       `json:"pjp_code"`
	SalesmanName *string    `json:"salesman_name"`
	Status       string     `json:"status"`
	Week         int        `json:"week"`
}

type OutletMapTes struct {
	OutletCode    string `json:"outlet_code"`
	OutletName    string `json:"outlet_name"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	OutletStatus  string `json:"outlet_status"`
	OutletAddress string `json:"outlet_address"`
}

type RoutesMapTes struct {
	RouteCode int            `json:"route_code"`
	Outlets   []OutletMapTes `json:"outlets"`
}

type RoutePopDailyTesResponse struct {
	ID           int            `json:"id"`
	Routes       []RoutesMapTes `json:"routes"`
	PjpID        *int           `json:"pjp_id"`
	PjpCode      *int           `json:"pjp_code"`
	SalesmanName *string        `json:"salesman_name"`
	Status       string         `json:"status"`
	Week         int            `json:"week"`
}
