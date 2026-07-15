package controller

import (
	"context"
	"errors"
	"net/http"
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

type DailyRouteMapController struct {
	routePopService service.RoutePopService
}

func NewDailyRouteMapController(service service.RoutePopService) *DailyRouteMapController {
	return &DailyRouteMapController{
		routePopService: service,
	}
}

// CreateTags		godoc
// @Summary			Copy all to daily.
// @Description     Copy all permanent to daily.
// @Param			approval body request.CopyAllRequest true "copy all permanent to daily"
// @Produce			application/json
// @Tags			daily route map
// @Success			200 {object} response.Response{}
// @Router			/web/daily-route-maps/all [post]
// @Security        Bearer
func (controller *DailyRouteMapController) CopyAllToDaily(ctx *gin.Context) {
	log.Info().Msg("copy all to daily")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.CopyAllRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.CopyAllPermanentToDaily(c, Request)
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

// CreateTags		godoc
// @Summary			Copy partial to daily.
// @Description     Copy partial permanent to daily.
// @Param			approval body request.CopyPartialRequest true "copy partial permanent to daily"
// @Produce			application/json
// @Tags			daily route map
// @Success			200 {object} response.Response{}
// @Router			/web/daily-route-maps/partial [post]
// @Security        Bearer
func (controller *DailyRouteMapController) CopyPartialToDaily(ctx *gin.Context) {
	log.Info().Msg("copy specific to daily")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.CopyPartialRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.CopyPartialPermanentToDaily(c, Request)
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

// CreateTags		godoc
// @Summary			Copy specific to daily.
// @Description     Copy specific permanent to daily.
// @Param			approval body request.RoutesMapping true "copy specific permanent to daily"
// @Produce			application/json
// @Tags			daily route map
// @Success			200 {object} response.Response{}
// @Router			/web/daily-route-maps/specific [post]
// @Security        Bearer
func (controller *DailyRouteMapController) CopySpecificToDaily(ctx *gin.Context) {
	log.Info().Msg("copy partial to daily")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.RoutesMapping{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.CopyToSpecificDaily(c, Request, currentCustomerId)
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

// CreateTags		godoc
// @Summary			Copy route daily to daily.
// @Description     Copy route daily to daily.
// @Param			approval body request.RoutesMapping true "copy route daily to daily"
// @Produce			application/json
// @Tags			daily route map
// @Success			200 {object} response.Response{}
// @Router			/web/daily-route-maps/to/daily [post]
// @Security        Bearer
func (controller *DailyRouteMapController) CopyRouteDailyToDaily(ctx *gin.Context) {
	log.Info().Msg("copy partial to daily")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.RoutesMapping{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.CopyRouteDailyToDaily(c, Request, currentCustomerId)
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
func (controller *DailyRouteMapController) FindAllPermanent(ctx *gin.Context) {
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

	responses := controller.routePopService.FindAllPermanent(c, filters, currentCustomerId)

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
// @Summary			Get all route pop daily.
// @Description		Return list of route pop daily.
// @Param		    pjp_code query string false "Pjp Code"
// @Param		    salesman_name query string false "Salesman Name"
// @Param		    route_code query string false "Route Code"
// @Produce		    application/json
// @Tags			daily route map
// @Success         200 {object} response.Response{}
// @Router			/web/daily-route-maps/daily [get]
// @Security        Bearer
func (controller *DailyRouteMapController) FindAllDaily(ctx *gin.Context) {
	log.Info().Msg("find all pop daily") // TODO Add params route_code
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	filters := make(map[string]interface{})
	filters["pjp_code"] = ctx.Query("pjp_code")
	filters["salesman_name"] = ctx.Query("salesman_name")
	filters["route_code"] = ctx.Query("route_code")

	responses := controller.routePopService.FindAllDaily(c, filters, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

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
func (controller *RouteMappingController) FindDailyRouteByPjpCode(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	pjpCode := ctx.Param("pjpCode")
	pjp_code, err := strconv.Atoi(pjpCode)
	helper.ErrorPanic(err)

	routeCode := ctx.Param("routeCode")
	route_code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	// Ambil nilai date dari context
	date := ctx.Query("date")

	log.Printf("Date ditemukan: %s", date)

	responses := controller.routeService.FindDailyRouteByPjpCode(c, pjp_code, route_code, date)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// CreateTags		godoc
// @Summary			add outlet to route for mobile.
// @Description     add outlet to route for mobile.
// @Param			save body request.AddOutletToRouteRequest true "save add outlet to route"
// @Produce			application/json
// @Tags			mobile
// @Success			201 {object} response.Response{}
// @Router			/mobile/add-outlets [post]
// @Security        Bearer
func (controller *DailyRouteMapController) MobileAddOutletToRoute(ctx *gin.Context) {
	log.Info().Msg("mobile add outlet to rute")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}
	customerCode, exists := ctx.Get("empCode")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "empCode not found",
		})
		return
	}

	empCodeStr, ok := customerCode.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "empCode is not a valid string",
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.AddOutletToRouteRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.SaveOutletToRoute(c, Request, currentCustomerId, empCodeStr)
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

// CreateTags		godoc
// @Summary			cancel add outlet to route for mobile.
// @Description     cancel add outlet to route for mobile.
// @Param			save body request.CancelAddOutletToRouteRequest true "cancel add outlet to route"
// @Produce			application/json
// @Tags			mobile
// @Success			200 {object} response.Response{}
// @Router			/mobile/cancel/add-outlets [post]
// @Security        Bearer
func (controller *DailyRouteMapController) MobileCancelAddOutletToRoute(ctx *gin.Context) {
	log.Info().Msg("mobile cancel add outlet to rute")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.CancelAddOutletToRouteRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routePopService.CancelOutletToRoute(c, Request)
	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Deleted successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
