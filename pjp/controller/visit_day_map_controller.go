package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/service"
	"scyllax-pjp/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type VisitDayMapController struct {
	routePopService service.RoutePopService
}

func NewVisitDayMapController(service service.RoutePopService) *VisitDayMapController {
	return &VisitDayMapController{
		routePopService: service,
	}
}

// CreateTags		godoc
// @Summary			Save weeklys.
// @Description     Recurring for all week,biweekly,next weekly.
// @Param			approval body request.SaveWeeklyRequest true "save recurring week and example date (2006-01-02 || yyyy-mm-dd)"
// @Produce			application/json
// @Tags			visit day map
// @Success			200 {object} response.Response{}
// @Router			/web/visit-day-maps [post]
// @Security        Bearer
func (controller *VisitDayMapController) SaveWeekly(ctx *gin.Context) {
	log.Info().Msg("save weekly")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), 120*time.Second)
	defer cancel()

	Request := request.SaveWeeklyRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.SaveWeekly(c, Request, currentCustomerId)
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

// FindAllTags 		    godoc
// @Summary				Get route outlet by additional routeCode.
// @Param				routeCode path string true "get route outlet by additional route_code"
// @Description			Return list of route outlet by additional route_code.
// @Produce				application/json
// @Tags				route mapping
// @Success				200 {object} response.Response{}
// @Router				/web/route-mappings/additional/{routeCode} [get]
// @Security            Bearer
func (controller *VisitDayMapController) FindByParentRoute(ctx *gin.Context) {
	log.Info().Msg("find by parent route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	routeCode := ctx.Param("routeCode")
	code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	responses := controller.routePopService.FindByRouteOutletAdditional(ctx, code, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAllTags 		godoc
// @Summary			Get all visit day map.
// @Description		Return list of visit day map.
// @Produce		    application/json
// @Param		    pjp_code query string false "pjp_code"
// @Param		    sort query string false "sort"
// @Tags			visit day map
// @Success         200 {object} response.Response{}
// @Router			/web/visit-day-maps [get]
// @Security        Bearer
func (controller *VisitDayMapController) GetAllVisitDayMap(ctx *gin.Context) {
	log.Info().Msg("get all visit day map")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.VisitDayMapQueryFilter
	dataFilter.PjpCode, _ = strconv.Atoi(ctx.Query("pjp_code"))
	dataFilter.Sort = ctx.Query("sort")

	fmt.Println("filter", dataFilter)
	responses := controller.routePopService.GetAllVisitDayMap(c, dataFilter, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
