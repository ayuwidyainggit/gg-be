package geotaging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"scyllax-pjp/config"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// DefaultRadius is the fallback radius in meters when API call fails or salesman has no radius set.
const DefaultRadius = 100

// GeotagingService defines the interface for geotaging validation operations.
type GeotagingService interface {
	ValidateGeotaging(ctx context.Context, req request.GeotagingValidationRequest, headers map[string]string, custID string) response.GeotagingValidationResponse
}

type geotagingServiceImpl struct {
	validate *validator.Validate
	db       *gorm.DB
	cfg      config.Config
}

// NewGeotagingService creates a new instance of GeotagingService.
func NewGeotagingService(validate *validator.Validate, db *gorm.DB, cfg config.Config) GeotagingService {
	return &geotagingServiceImpl{
		validate: validate,
		db:       db,
		cfg:      cfg,
	}
}

// ValidateGeotaging validates the user's GPS location against the outlet location.
// It calculates the distance and determines if the user is within the allowed radius.
func (s *geotagingServiceImpl) ValidateGeotaging(ctx context.Context, req request.GeotagingValidationRequest, headers map[string]string, custID string) response.GeotagingValidationResponse {
	requestID := generateRequestID()

	// Check if GPS is disabled (arrival coordinates are empty)
	if req.ArrivalLongitude == "" || req.ArrivalLatitude == "" {
		return s.buildGPSDisabledResponse(requestID)
	}

	// Get allowed radius from salesman data
	allowedRadius := s.fetchSalesmanRadius(ctx, req.EmpId, headers, custID)

	// Parse and validate coordinates
	outletLat, outletLng := s.parseCoordinates(req.OutletLatitude, req.OutletLongitude, "outlet")
	arrivalLat, arrivalLng := s.parseCoordinates(req.ArrivalLatitude, req.ArrivalLongitude, "arrival")

	// Calculate distance using Haversine formula
	distanceMeter := helper.CalculateHaversineDistance(outletLat, outletLng, arrivalLat, arrivalLng)

	// Determine location status based on distance
	locationStatus := s.determineLocationStatus(distanceMeter, allowedRadius)

	return s.buildSuccessResponse(requestID, locationStatus, distanceMeter, allowedRadius)
}

// generateRequestID creates a unique request ID for tracking.
func generateRequestID() string {
	return uuid.New().String()[:24]
}

// buildGPSDisabledResponse creates an error response for GPS disabled case.
func (s *geotagingServiceImpl) buildGPSDisabledResponse(requestID string) response.GeotagingValidationResponse {
	errorCode := response.ErrorCodeGPSDisabled
	return response.GeotagingValidationResponse{
		Message:   "GPS is required to confirm arrival",
		Data:      nil,
		ErrorCode: &errorCode,
		RequestId: requestID,
	}
}

// parseCoordinates parses latitude and longitude strings to float64.
// Returns 0 for invalid coordinates and logs the error.
func (s *geotagingServiceImpl) parseCoordinates(latStr, lngStr, label string) (lat, lng float64) {
	var err error

	lat, err = strconv.ParseFloat(latStr, 64)
	if err != nil {
		log.Warn().Err(err).Str("label", label).Str("latitude", latStr).Msg("failed to parse latitude")
		lat = 0
	}

	lng, err = strconv.ParseFloat(lngStr, 64)
	if err != nil {
		log.Warn().Err(err).Str("label", label).Str("longitude", lngStr).Msg("failed to parse longitude")
		lng = 0
	}

	return lat, lng
}

// determineLocationStatus determines the location status based on distance and allowed radius.
func (s *geotagingServiceImpl) determineLocationStatus(distanceMeter, allowedRadius int) string {
	if distanceMeter <= allowedRadius {
		return response.LocationStatusInRadius
	}
	return response.LocationStatusForcedOutOfRadius
}

// buildSuccessResponse creates a successful validation response.
func (s *geotagingServiceImpl) buildSuccessResponse(requestID, locationStatus string, distanceMeter, allowedRadius int) response.GeotagingValidationResponse {
	respData := &response.GeotagingValidationData{
		LocationStatus: locationStatus,
		DistanceMeter:  distanceMeter,
		AllowedRadius:  allowedRadius,
	}

	// Add confirmed_at timestamp only if user is within radius
	if locationStatus == response.LocationStatusInRadius {
		confirmedAt := time.Now().Format(time.RFC3339)
		respData.ConfirmedAt = &confirmedAt
	}

	return response.GeotagingValidationResponse{
		Message:   "Arrival confirmed successfully",
		Data:      respData,
		RequestId: requestID,
	}
}

// fetchSalesmanRadius retrieves the allowed radius from salesman data via master service.
// Returns DefaultRadius if the API call fails or salesman has no radius configured.
func (s *geotagingServiceImpl) fetchSalesmanRadius(ctx context.Context, empID int64, headers map[string]string, custID string) int {
	endpointURL := fmt.Sprintf("%s/v1/salesman/%d", s.cfg.KongUrl, empID)
	log.Debug().Str("url", endpointURL).Int64("emp_id", empID).Msg("fetching salesman radius")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create HTTP request for salesman")
		return DefaultRadius
	}

	// Set required headers
	req.Header.Set("cust_id", custID)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch salesman data")
		return DefaultRadius
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status_code", resp.StatusCode).Msg("salesman API returned non-OK status")
		return DefaultRadius
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read salesman response body")
		return DefaultRadius
	}

	var result struct {
		Message   string            `json:"message"`
		Data      model.NewSalesman `json:"data"`
		RequestID string            `json:"request_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal salesman response")
		return DefaultRadius
	}

	if result.Data.SmRadius > 0 {
		return result.Data.SmRadius
	}

	return DefaultRadius
}
