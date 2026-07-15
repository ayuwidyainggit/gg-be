package request

type Week struct {
	Date string `json:"date"`
	Day  string `json:"day"`
	Week int    `json:"week" validate:"required"`
	Year int    `json:"year" validate:"required"`
}

type RouteData struct {
	PjpCode   string `validate:"required" json:"pjp_code"`
	PjpID     int    `validate:"required" json:"pjp_id"`
	RouteCode int    `validate:"required" json:"route_code"`
	Weeks     []Week `validate:"required,dive" json:"weeks"`
}

type RouteDailyData struct {
	PjpCode   string    `validate:"required" json:"pjp_code"`
	PjpID     int       `validate:"required" json:"pjp_id"`
	RouteCode int       `validate:"required" json:"route_code"`
	Weeks     []Week    `validate:"required,dive" json:"weeks"`
	Outlets   []Outlets `validate:"required,dive" json:"outlets"`
}

type AddOutletToRouteData struct {
	PjpCode   string `validate:"required" json:"pjp_code" example:"2"`
	PjpID     int    `validate:"required" json:"pjp_id" example:"53"`
	RouteCode int    `validate:"required" json:"route_code" example:"7170"`
	// Date      string `validate:"required,date"    json:"date"`
	Outlets []Outlets `validate:"required,dive" json:"outlets"`
}

type SaveWeeklyRequest struct {
	Data []RouteData `json:"data"  validate:"required,dive"`
}

type SaveDailyRouteMap struct {
	Data []RouteDailyData `json:"data"  validate:"required,dive"`
}

type AddOutletToRouteRequest struct {
	Data []AddOutletToRouteData `json:"data"  validate:"required,dive"`
}

type SaveDelegateRequest struct {
	RouteCode int    `validate:"required"         json:"route_code"`
	PjpCode   string `validate:"required"         json:"pjp_code"`
	PjpID     int    `validate:"required"         json:"pjp_id"`
	Week      int    `validate:"required"         json:"week"`
	Day       string `validate:"required"         json:"day"`
	Date      string `validate:"required,date"    json:"date"`
	Year      int    `validate:"required"         json:"year"`
}

type CopyAllRequest struct {
	PjpCode []int `validate:"required,notEmptyIntSlice" json:"pjp_code"`
}

type RoutesData struct {
	PjpCode int      `validate:"required" json:"pjp_code"`
	Routes  []Routes `json:"routes"`
}

type CopyPartialRequest struct {
	Data []RoutesData `json:"data"`
}

type Routes struct {
	RouteCode int `validate:"required" json:"route_code"`
}

type RoutesMapping struct {
	PjpCode   int      `validate:"required" json:"pjp_code"`
	PjpID     int      `validate:"required" json:"pjp_id"`
	RouteCode int      `validate:"required" json:"route_code"`
	Day       string   `validate:"required" json:"day"`
	Week      int      `validate:"required" json:"week"`
	Date      string   `validate:"required" json:"date"`
	Year      int      `validate:"required" json:"year"`
	Routes    []Routes `json:"routes"`
}

type CancelAddOutletToRouteData struct {
	PjpID      int    `validate:"required" json:"pjp_id" example:"53"`
	RouteCode  int    `validate:"required" json:"route_code" example:"7170"`
	OutletCode string `validate:"required" json:"outlet_code" example:"20241126"`
}

type CancelAddOutletToRouteRequest struct {
	Data []CancelAddOutletToRouteData `json:"data"  validate:"required,dive"`
}
