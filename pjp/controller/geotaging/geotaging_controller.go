package geotaging

import (
	"net/http"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/service/geotaging"

	"github.com/gin-gonic/gin"
)

// GeotagingController defines the interface for geotaging controller
type GeotagingController interface {
	ValidateGeotaging(ctx *gin.Context)
}

type geotagingController struct {
	geotagingService geotaging.GeotagingService
}

// NewGeotagingController creates a new instance of GeotagingController
func NewGeotagingController(geotagingService geotaging.GeotagingService) GeotagingController {
	return &geotagingController{
		geotagingService: geotagingService,
	}
}

// ValidateGeotaging godoc
// @Summary Validate geotaging for arrival
// @Description Validates user location against outlet location and returns distance information
// @Tags Mobile - Visits
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body request.GeotagingValidationRequest true "Geotaging validation request"
// @Success 200 {object} response.GeotagingValidationResponse "Successful validation"
// @Failure 400 {object} response.GeotagingValidationResponse "GPS disabled or invalid request"
// @Failure 500 {object} response.GeotagingValidationResponse "Internal server error"
// @Router /mobile/visits/geotaging [post]
func (c *geotagingController) ValidateGeotaging(ctx *gin.Context) {
	var req request.GeotagingValidationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.GeotagingValidationResponse{
			Message:   "Invalid request body: " + err.Error(),
			Data:      nil,
			RequestId: "",
		})
		return
	}

	// Get headers for API call
	headers := make(map[string]string)
	headers["Authorization"] = ctx.GetHeader("Authorization")
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	// Get cust_id from header or context
	custId := ctx.GetHeader("cust_id")
	if custId == "" {
		custId = req.CustId
	}

	result := c.geotagingService.ValidateGeotaging(ctx.Request.Context(), req, headers, custId)

	// Determine HTTP status based on error code
	statusCode := http.StatusOK
	if result.ErrorCode != nil && *result.ErrorCode == response.ErrorCodeGPSDisabled {
		statusCode = http.StatusBadRequest
	}

	ctx.JSON(statusCode, result)
}
