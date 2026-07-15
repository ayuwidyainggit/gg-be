package response

// GeotagingValidationData represents the data returned from geotaging validation
type GeotagingValidationData struct {
	LocationStatus string  `json:"location_status" example:"IN_RADIUS"`
	DistanceMeter  int     `json:"distance_meter" example:"135"`
	AllowedRadius  int     `json:"allowed_radius" example:"100"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty" example:"2026-01-06T19:45:03+07:00"`
}

// GeotagingValidationResponse represents the response for geotaging validation API
type GeotagingValidationResponse struct {
	Message   string                   `json:"message" example:"Arrival confirmed successfully"`
	Data      *GeotagingValidationData `json:"data"`
	ErrorCode *string                  `json:"error_code,omitempty" example:"GPS_DISABLED"`
	RequestId string                   `json:"request_id" example:"6915a6dd2395083c685e8e16"`
}

// Location status constants
const (
	LocationStatusInRadius          = "IN_RADIUS"
	LocationStatusForcedOutOfRadius = "FORCED_OUT_OF_RADIUS"
	LocationStatusOutOfRadius       = "OUT_OF_RADIUS"
	ErrorCodeGPSDisabled            = "GPS_DISABLED"
)
