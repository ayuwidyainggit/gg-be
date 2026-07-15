package response

type RouteDestinationResponse struct {
	ID        int    `json:"id"`
	RouteCode int    `json:"route_code"`
	RouteName string `json:"route_name"`
	PjpID     *int   `json:"pjp_id"`
	PjpCode   *int   `json:"pjp_code"`
	Status    string `json:"status"`
}

type RouteDetailResponse struct {
	ID        int                   `json:"id"`
	RouteCode int                   `json:"route_code"`
	RouteName string                `json:"route_name"`
	Outlets   []DestinationResponse `json:"outlets"`
}

type RouteOutletsResponse struct {
	RouteCode int                   `json:"route_code"`
	RouteName string                `json:"route_name"`
	Outlets   []DestinationResponse `json:"outlets"`
}
