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

// GetDistributorMonitoring godoc
// @Summary     Get location monitoring data for Distributor
// @Description Return list of employee location monitoring data for distributor users
// @Param       region_id query int true "Region ID"
// @Param       area_id query int true "Area ID"
// @Param       distributor_id query int true "Distributor ID"
// @Param       date query int true "Date in epoch format"
// @Param       emp_id[] query []int false "Employee IDs filter"
// @Param       status[] query []string true "Approval status filter (e.g., Approved)"
// @Param       page query int false "Page number"
// @Param       limit query int false "Page limit"
// @Produce     application/json
// @Tags        live-monitoring
// @Success     200 {object} response.LiveMonitoringResponse{}
// @Router      /v1/live-monitoring-distributor [get]
// @Security    Bearer
func (c *liveMonitoringController) GetDistributorMonitoring(ctx *gin.Context) {
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

	// Generate request ID
	requestID := uuid.New().String()

	// Set request context with timeout
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	// Bind query parameters
	var req request.LiveMonitoringRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    constant.ErrInvalidRequestParams + err.Error(),
			"request_id": requestID,
		})
		return
	}

	// Call service
	data, paging, err := c.service.GetDistributorMonitoring(reqCtx, req, customerID)
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
	if len(data) == 0 {
		message = constant.MsgNoData
		responseData = nil
	} else {
		message = constant.MsgSuccess
		responseData = data
	}

	webResponse := gin.H{
		"message":    message,
		"data":       responseData,
		"paging":     paging,
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
