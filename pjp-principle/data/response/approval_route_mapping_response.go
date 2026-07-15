package response

type ApprovalRouteMappingResponse struct {
	ID              int                  `json:"id"`
	RouteCode       int                  `json:"route_code"`
	RouteName       string               `json:"route_name"`
	Outlets         *DestinationResponse `json:"outlets"`
	PjpID           *int                 `json:"pjp_id"`
	PjpCode         *int                 `json:"pjp_code"`
	Status          string               `json:"status"`
	SalesmanName    *string              `json:"salesman_name"`
	SalesmanCode    *string              `json:"salesman_code"`
	Date            string               `json:"date"`
	VerifiedDate    string               `json:"verified_date"`
	NewSalesmanName *string              `json:"new_salesman_name"`
	NewSalesmanCode *string              `json:"new_salesman_code"`
	NewRouteCode    int                  `json:"new_route_code"`
	NewRouteName    string               `json:"new_route_name"`
}

type ApprovalRouteMappingEnhanceResponse struct {
	PjpID        int    `json:"pjp_id"`
	PjpCode      string `json:"pjp_code"`
	Status       string `json:"status"`
	SalesmanName string `json:"salesman_name"`
	SalesmanCode string `json:"salesman_code"`
	Date         string `json:"date"`
	VerifiedDate string `json:"verified_date"`
}
