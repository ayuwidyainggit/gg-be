package live_monitoring

import (
	"context"
	"net/http"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetMonitoringDetail godoc
// @Summary     Get location monitoring detail for an employee
// @Description Return detailed monitoring information including sales, returns, expenses, and shipments
// @Param       emp_id query int true "Employee ID"
// @Param       distributor_id query int false "Distributor ID (if NULL = Principal, if NOT NULL = Distributor)"
// @Param       date query string true "Date in YYYY-MM-DD format"
// @Produce     application/json
// @Tags        live-monitoring
// @Success     200 {object} response.LiveMonitoringDetailResponse{}
// @Router      /v1/monitoring_locations/details [get]
// @Security    Bearer
func (c *liveMonitoringController) GetMonitoringDetail(ctx *gin.Context) {
	// Get customer ID from JWT context
	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  constant.MsgUnauthorized,
			Message: constant.ErrCustomerIDNotFound,
		})
		return
	}

	// Get user ID from JWT context for expense filtering
	userID, userExists := helper.GetCurrentUserId(ctx)
	if !userExists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  constant.MsgUnauthorized,
			Message: "user_id not found",
		})
		return
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Set request context with timeout
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Bind query parameters
	var req request.LiveMonitoringDetailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    constant.ErrInvalidRequestParams + err.Error(),
			"request_id": requestID,
		})
		return
	}

	// Call service with userID for expense filtering
	data, err := c.service.GetMonitoringDetail(reqCtx, req, customerID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":    err.Error(),
			"request_id": requestID,
		})
		return
	}

	// Build response
	var message string
	var responseData interface{}
	if data == nil {
		message = constant.MsgNoData
		responseData = nil
	} else {
		message = constant.MsgSuccess
		responseData = []interface{}{data}
	}

	webResponse := gin.H{
		"message":    message,
		"data":       responseData,
		"request_id": requestID,
	}

	utils.ResponseInterceptor(reqCtx, &response.Response{
		Code:    http.StatusOK,
		Status:  constant.MsgOK,
		Data:    responseData,
		Message: message,
	})

	ctx.JSON(http.StatusOK, webResponse)
}
