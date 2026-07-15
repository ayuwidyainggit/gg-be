package request

type CreateRouteRequest struct {
	RouteName string `validate:"required,min=2,max=125" json:"route_name"`
}

type SaveOutletRequest struct {
	RouteCode int       `validate:"required" json:"route_code"`
	Outlets   []Outlets `validate:"required,dive" json:"outlets"`
}

type SavePjpRequest struct {
	RouteCode []int  `validate:"required,notEmptyIntSlice" json:"route_code"`
	PjpCode   string `validate:"required" json:"pjp_code"`
	PjpID     int    `validate:"required" json:"pjp_id"`
}

type DeletePjpRequest struct {
	RouteCode int    `validate:"required" json:"route_code"`
	PjpCode   string `validate:"required" json:"pjp_code"`
	PjpID     int    `validate:"required" json:"pjp_id"`
}

type NewRouteRequest struct {
	NewRoutePropose []NewRouteBody `validate:"required,dive" json:"new_route_propose"`
}

type NewRouteBody struct {
	DestinationID   int    `validate:"required" json:"outlet_id"`
	DestinationCode string `validate:"required" json:"outlet_code"`
	PjpCode         string `validate:"required" json:"pjp_code"`
	PjpID           int    `validate:"required" json:"pjp_id"`
	RouteCode       int    `validate:"required" json:"route_code"`
	RouteName       string `validate:"required" json:"route_name"`
	OldPjpCode      int    `validate:"required" json:"old_pjp_code"`
	OldPjpID        int    `validate:"required" json:"old_pjp_id"`
	OldRouteCode    int    `validate:"required" json:"old_route_code"`
	OldRouteName    string `validate:"required" json:"old_route_name"`
}

type UpdateRoutesRequest struct {
	ID        int    `validate:"required"`
	RouteName string `validate:"required" json:"route_name"`
}

type DeleteOutletRequest struct {
	RouteCode       int      `validate:"required" json:"route_code"`
	DestinationCode []string `validate:"required,notEmptyStringSlice" json:"outlet_code"`
	Week            int      `validate:"omitempty" json:"week"`
	Date            string   `validate:"omitempty" json:"date"`
}

type DeleteOutletAdditionalRequest struct {
	RouteCode       int      `validate:"required" json:"route_code"`
	DestinationCode []string `validate:"required,notEmptyStringSlice" json:"outlet_code"`
	Week            int      `validate:"required" json:"week"`
	Date            string   `validate:"required" json:"date"`
}

type UpdatePjpInRouteRequest struct {
	RouteCode []int `validate:"required,notEmptyIntSlice" json:"route_code"`
	PjpCode   int   `validate:"required" json:"pjp_code"`
}

type UpdateStatusRequest struct {
	ID     []int  `validate:"required,notEmptyIntSlice" json:"id"`
	Status string `validate:"required" json:"status"`
}
type UpdateStatusEnhanceRequest struct {
	PjpCode []string `validate:"required,notEmptyStringSlice" json:"pjp_code"`
	Status  string   `validate:"required" json:"status"`
}

type AddRouteVisitDayRequest struct {
	RouteCode int    `validate:"required" json:"route_code"`
	Day       string `validate:"required" json:"day"`
}

type Outlets struct {
	DestinationID      int     `validate:"required" json:"destination_id" example:"5"`
	DestinationCode    string  `validate:"required" json:"destination_code" example:"OUT 5"`
	DestinationName    string  `validate:"required" json:"destination_name" example:"Destination okok"`
	Longitude          string  `validate:"required" json:"longitude" example:"106.816666"`
	Latitude           string  `validate:"required" json:"latitude" example:"-6.200000"`
	DestinationStatus  int     `validate:"required" json:"destination_status" example:"1"`
	DestinationAddress string  `validate:"required" json:"destination_address" example:"mampang"`
	DestinationType    string  `validate:"required" json:"destination_type" example:"outlet"`
	AvgSalesWeek       float64 `json:"avg_sales_week" example:"100"`
}

type OutletProcess struct {
	DestinationID int    `validate:"required" json:"outlet_id"`
	Longitude     string `validate:"required" json:"longitude"`
	Latitude      string `validate:"required" json:"latitude"`
	AvgSalesWeek  string `json:"avg_sales_week"`
}

type ParamsID struct {
	RouteID []int `validate:"required,notEmptyIntSlice" json:"id"`
}

type Route struct {
	RouteCode          int    `validate:"required" json:"route_code"`
	DestinationID      int    `validate:"required" json:"outlet_id"`
	DestinationCode    string `validate:"required" json:"outlet_code"`
	DestinationName    string `validate:"required" json:"outlet_name"`
	Longitude          string `validate:"required" json:"longitude"`
	Latitude           string `validate:"required" json:"latitude"`
	DestinationStatus  int    `validate:"required" json:"outlet_status"`
	DestinationAddress string `validate:"required" json:"outlet_address"`
}

type SaveRouteConfirmationRequest struct {
	Routes []Route `json:"new_routes"`
}

type DuplicateRoute struct {
	RouteCode int `validate:"required" json:"route_code"`
}
