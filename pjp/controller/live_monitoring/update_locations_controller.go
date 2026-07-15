package live_monitoring

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c *liveMonitoringController) GetUpdateLocations(ctx *gin.Context) {
	customerID, exists := helper.GetCurrentCustomerId(ctx)
	requestID := uuid.New().String()
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": constant.MsgUnauthorized, "request_id": requestID})
		return
	}
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()
	var req request.UpdateLocationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": constant.ErrInvalidRequestParams + err.Error(), "request_id": requestID})
		return
	}
	data, err := c.service.GetUpdateLocations(reqCtx, req, customerID)
	if err != nil {
		status := http.StatusInternalServerError
		message := err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
			message = http.StatusText(http.StatusNotFound)
		}
		ctx.JSON(status, gin.H{"message": message, "request_id": requestID})
		return
	}
	message := constant.MsgSuccess
	if len(data.Timeline) == 0 {
		message = constant.MsgNoData
	}
	ctx.JSON(http.StatusOK, gin.H{"message": message, "data": data, "request_id": requestID})
}
