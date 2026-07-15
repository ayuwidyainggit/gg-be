package route

import (
	"context"
	"errors"
	"net/http"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// FindAllTags 		godoc
// @Summary			Get all approval route.
// @Description		Return list of approval route.
// @Param		    limit query string false "Limit"
// @Param		    page query string false "Page"
// @Param		    pjp_code query string false "pjp_code"
// @Param		    route_code query string false "route_code"
// @Param		    status query string false "status"
// @Param		    salesman_name query string false "salesman_name"
// @Param		    salesman_code query string false "salesman_code"
// @Param           start_date query string false "start date"
// @Param           end_date query string false "end date"
// @Produce		    application/json
// @Tags			approval route
// @Success         200 {object} response.Response{}
// @Router			/web/approval-routes-enhance [get]
// @Security        Bearer
func (controller *routeController) FindAllApprovalEnhance(ctx *gin.Context) {
	log.Info().Msg("findAll approval list route")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	limit, _ := strconv.Atoi(ctx.Query("limit"))

	page, _ := strconv.Atoi(ctx.Query("page"))

	filters := make(map[string]interface{})
	filters["pjp_code"] = ctx.Query("pjp_code")
	filters["route_code"] = ctx.Query("route_code")
	filters["status"] = ctx.Query("status")
	filters["salesman_name"] = ctx.Query("salesman_name")
	filters["salesman_code"] = ctx.Query("salesman_code")

	startDate, err := helper.ParseDateFilter(ctx.Query("start_date"), "2006-01-02") //yyyy-mm-dd
	if err != nil {
		helper.ErrorPanic(err)
	}

	endDate, err := helper.ParseDateFilter(ctx.Query("end_date"), "2006-01-02") //yyyy-mm-dd
	if err != nil {
		helper.ErrorPanic(err)
	}

	if limit < 1 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	filters["start_date"] = startDate
	filters["end_date"] = endDate

	responses, pagination, err := controller.routeService.GetAllEnhance(c, page, limit, filters, currentCustomerId)
	helper.ErrorPanic(err)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
		Meta:   &pagination,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
