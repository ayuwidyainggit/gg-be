package route

import (
	"context"
	"net/http"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CreateTags		godoc
// @Summary			Update status in route to approval with propose.
// @Description		Update status in route to approval with propose.
// @Param			approval body request.UpdateStatusEnhanceRequest true "update status approval with propose"
// @Produce			application/json
// @Tags			approval route
// @Success			200 {object} response.Response{}
// @Router			/web/approval-routes-enhance/status [patch]
// @Security        Bearer
func (controller *routeController) UpdateStatusEnhance(ctx *gin.Context) {
	log.Info().Msg("update status in route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	customerID, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Response{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "customer_id not found",
		})
		return
	}

	Request := request.UpdateStatusEnhanceRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.UpdateStatusEnhance(c, Request, customerID)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Updated successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
