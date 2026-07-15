package daily

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get all route pop permanent.
// @Description		Return list of route pop permanent.
// @Param		    pjp_code query string false "Pjp Code"
// @Param		    salesman_name query string false "Salesman Name"
// @Param		    week query string false "Week"
// @Produce		    application/json
// @Tags			daily route map
// @Success         200 {object} response.Response{}
// @Router			/web/daily-route-maps [get]
// @Security        Bearer
func (controller *dailyRouteMapController) GetAll(ctx *gin.Context) {
	log.Info().Msg("find all permanent pjp")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	filters := make(map[string]interface{})
	filters["pjp_code"] = ctx.Query("pjp_code")
	filters["salesman_name"] = ctx.Query("salesman_name")
	filters["week"] = ctx.Query("week")

	responses := controller.routePopService.GetAll(c, filters, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}

	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
