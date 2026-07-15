package response

type DestinationResponse struct {
	DestinationID      int     `json:"destination_id"`
	DestinationCode    string  `json:"destination_code"`
	DestinationName    string  `json:"destination_name"`
	Longitude          string  `json:"longitude"`
	Latitude           string  `json:"latitude"`
	DestinationStatus  string  `json:"destination_status"`
	DestinationAddress string  `json:"destination_address"`
	DestinationType    *string `json:"destination_type"`
	AvgSalesWeek       float64 `json:"avg_sales_week"`
	Status             string  `json:"status"`
	PjpCode            int     `json:"pjp_code"`
	PjpID              int     `json:"pjp_id"`
	RouteCode          int     `json:"route_code"`
}

type DailyDestinationResponse struct {
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
