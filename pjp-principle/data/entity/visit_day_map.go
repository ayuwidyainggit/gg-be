package entity

type VisitDayMapResponse struct {
	ID          int    `json:"id"`
	RouteCode   int    `json:"route_code"`
	RouteName   string `json:"route_name"`
	TotalOutlet int    `json:"total_outlet"`
	WeekNumber  int    `json:"week_number"`
	Date        string `json:"date"`
}

type VisitDayMapQueryFilter struct {
	PjpCode int    `query:"pjp_code"`
	Sort    string `query:"sort"`
}
