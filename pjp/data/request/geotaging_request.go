package request

// GeotagingValidationRequest represents the request body for geotaging validation
type GeotagingValidationRequest struct {
	EmpId            int64  `json:"emp_id" validate:"required" example:"1"`
	UserId           int64  `json:"user_id" validate:"required" example:"43"`
	CustId           string `json:"cust_id" validate:"required" example:"C220010001"`
	Activity         string `json:"activity" validate:"required,oneof=arrive leave" example:"arrive"`
	CurrentTime      int64  `json:"current_time" validate:"required" example:"1704556800"`
	OutletId         int64  `json:"outlet_id" validate:"required" example:"190"`
	OutletLongitude  string `json:"outlet_longitude" validate:"required" example:"106.845599"`
	OutletLatitude   string `json:"outlet_latitude" validate:"required" example:"-6.208763"`
	ArrivalLongitude string `json:"arrival_longitude" example:"106.845600"`
	ArrivalLatitude  string `json:"arrival_latitude" example:"-6.208760"`
}
