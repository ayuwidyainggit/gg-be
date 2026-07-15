package request

type StartVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  int64  `json:"current_time" validate:"required"`
	EmpID        int
	// RouteCode    string `validate:"required" json:"route_code"`
}

type OutletVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	EmpID        int64
}

type FinishVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  int64  `json:"current_time" validate:"required"`
	// Id           int64  `validate:"required" json:"id"`
}

type SkipVisitRequest struct {
	Date             string   `validate:"required,date" json:"date"`
	SalesmanCode     string   `validate:"required" json:"salesman_code"`
	CustID           string   `validate:"required" json:"cust_id"`
	CurrentTime      int64    `json:"current_time" validate:"required"`
	Id               int64    `validate:"required" json:"id"`
	SkipReason       string   `validate:"required" json:"skip_reason"`
	FileUrl          string   `json:"file_url"`
	Latitude         *float64 `json:"latitude" form:"latitude"`
	Longitude        *float64 `json:"longitude" form:"longitude"`
	IsUpdateLocation *bool    `json:"is_update_location" form:"is_update_location"`
	InOutlet         bool     `json:"in_outlet" form:"in_outlet"`
}

type ResumeVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}

type ArriveVisitRequest struct {
	Date             string   `validate:"required,date" json:"date" form:"date"`
	SalesmanCode     string   `validate:"required" json:"salesman_code" form:"salesman_code"`
	CustID           string   `validate:"required" json:"cust_id" form:"cust_id"`
	CurrentTime      *int64   `json:"current_time" validate:"required" form:"current_time"`
	Id               int64    `validate:"required" json:"id" form:"id"`
	OutletID         *int64   `json:"outlet_id" form:"outlet_id"`
	Latitude         *float64 `json:"latitude" form:"latitude"`
	Longitude        *float64 `json:"longitude" form:"longitude"`
	IsUpdateLocation *bool    `validate:"required" json:"is_update_location" form:"is_update_location"`
	FileUrl          string   `validate:"required" json:"file_url" form:"file_url"`
	// New fields for arrival report
	UserId          *int64  `json:"user_id" form:"user_id"`
	Activity        *string `json:"activity" form:"activity"`
	OutletLongitude *string `json:"outlet_longitude" form:"outlet_longitude"`
	OutletLatitude  *string `json:"outlet_latitude" form:"outlet_latitude"`
	LocationStatus  *string `json:"location_status" form:"location_status"`
	DistanceMeter   *int    `json:"distance_meter" form:"distance_meter"`
	AllowedRadius   *int    `json:"allowed_radius" form:"allowed_radius"`
}

type LeaveVisitRequest struct {
	Date           string  `validate:"required,date" json:"date"`
	SalesmanCode   string  `validate:"required" json:"salesman_code"`
	CustID         string  `validate:"required" json:"cust_id"`
	CurrentTime    *int64  `json:"current_time" validate:"required"`
	Id             int64   `validate:"required" json:"id"`
	LeaveLatitude  *string `json:"leave_latitude"`
	LeaveLongitude *string `json:"leave_longitude"`
}

type OnholdVisitRequest struct {
	Date         string `validate:"required,date" json:"date"`
	SalesmanCode string `validate:"required" json:"salesman_code"`
	CustID       string `validate:"required" json:"cust_id"`
	CurrentTime  *int64 `json:"current_time" validate:"required"`
	Id           int64  `validate:"required" json:"id"`
}
