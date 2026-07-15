package daily

import (
	"context"
	"errors"
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
// @Summary			Save daily route map additional.
// @Description     Save daily route map additional example date(yyyy-mm-dd).
// @Param			save body request.SaveDailyRouteMap true "save daily route map and example date (2006-01-02 || yyyy-mm-dd)"
// @Produce			application/json
// @Tags			daily route map
// @Success			200 {object} response.Response{}
// @Router			/web/daily-route-maps [post]
// @Security        Bearer
func (controller *dailyRouteMapController) SaveDailyRouteMap(ctx *gin.Context) {
	log.Info().Msg("save daily route map")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.SaveDailyRouteMap{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.SaveDailyRouteMap(c, Request, currentCustomerId)
	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}
