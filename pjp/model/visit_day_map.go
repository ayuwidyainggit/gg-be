package model

import "time"

type VisitDayMap struct {
	ID          int       `json:"id"`
	RouteCode   int       `json:"route_code"`
	RouteName   string    `json:"route_name"`
	TotalOutlet int       `json:"total_outlet"`
	Week        int       `json:"week"`
	Date        time.Time `json:"date"`
}
