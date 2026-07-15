package response

import "time"

type RoutesMap struct {
	RouteCode        int       `json:"route_code"`
	RouteName        string    `json:"route_name"`
	Week             int       `json:"week"`
	Date             time.Time `json:"date"`
	TotalOutlet      int       `json:"total_outlet"`
	TotalDistributor int       `json:"total_distributor"`
}

type RoutePopPermanentResponse struct {
	ID           int         `json:"id"`
	Routes       []RoutesMap `json:"routes"`
	PjpID        *int        `json:"pjp_id"`
	PjpCode      *int        `json:"pjp_code"`
	SalesmanName *string     `json:"salesman_name"`
}
