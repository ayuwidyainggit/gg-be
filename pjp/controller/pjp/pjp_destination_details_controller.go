package pjp

import (
	"context"
	"net/http"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetDestinationDetails godoc
// @Summary     Get destination details.
// @Description Return outlet and distributor destination details from destination history.
// @Param       pjp_id path int true "PJP ID"
// @Param       page query int false "Page"
// @Param       limit query int false "Limit"
// @Param       sort_order query string false "Sort order by destination ID"
// @Param       date query string true "Date in YYYY-MM-DD format"
// @Produce     application/json
// @Tags        pjp
// @Success     200 {object} response.DestinationDetailsResponse{}
// @Router      /web/destination-details/{pjp_id} [get]
// @Security    Bearer
func (controller *pjpController) GetDestinationDetails(ctx *gin.Context) {
	customerID, exists := helper.GetCurrentCustomerId(ctx)
	requestID := uuid.New().String()
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message":    constant.ErrCustomerIDNotFound,
			"request_id": requestID,
		})
		return
	}
	pjpID, err := strconv.Atoi(ctx.Param("pjp_id"))
	if err != nil || pjpID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    constant.ErrInvalidRequestParams + "pjp_id must be a positive integer",
			"request_id": requestID,
		})
		return
	}

	var req request.DestinationDetailsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    constant.ErrInvalidRequestParams + err.Error(),
			"request_id": requestID,
		})
		return
	}
	req.PjpID = pjpID
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    constant.ErrInvalidRequestParams + "date must use YYYY-MM-DD",
			"request_id": requestID,
		})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	data, paging, err := controller.pjpService.GetDestinationDetails(reqCtx, req, customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":    err.Error(),
			"request_id": requestID,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.DestinationDetailsResponse{
		Message:   constant.MsgSuccess,
		Data:      data,
		Paging:    paging,
		RequestID: requestID,
	})
}
