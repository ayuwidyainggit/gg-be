package daily

import (
	"context"
	"log"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FindAllTags 		godoc
// @Summary			Get all route by pjp_code and route_code.
// @Description		Return list of route by pjp_code and route_code.
// @Param		    pjpCode path string true "pjp_code"
// @Param		    routeCode path string true "route_code"
// @Param		    date query string true "date"
// @Produce		    application/json
// @Tags			daily route map
// @Success         200 {object} response.Response{}
// @Router			/web/daily-route-maps/pjp/{pjpCode}/{routeCode} [get]
// @Security        Bearer
func (controller *dailyRouteMapController) GetByPjpAndRouteCode(ctx *gin.Context) {
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

	pjpCode := ctx.Param("pjpCode")
	pjp_code, err := strconv.Atoi(pjpCode)
	helper.ErrorPanic(err)

	routeCode := ctx.Param("routeCode")
	route_code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	// Ambil nilai date dari context
	date := ctx.Query("date")

	log.Printf("Date ditemukan: %s", date)

	responses := controller.routePopService.GetByPjpAndRouteCode(c, pjp_code, route_code, date, customerID)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
