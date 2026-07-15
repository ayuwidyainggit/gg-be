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

type RouteMappingController struct {
	routeService service.RouteService
}

func NewRouteMappingController(service service.RouteService) *RouteMappingController {
	return &RouteMappingController{
		routeService: service,
	}
}

// CreateTags		godoc
// @Summary			Create route.
// @Description		Save route data in Db.
// @Param			route body request.CreateRouteRequest true "Create routes"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings [post]
// @Security        Bearer
func (controller *RouteMappingController) Create(ctx *gin.Context) {
	log.Info().Msg("create route")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.CreateRouteRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.Create(c, Request, currentCustomerId)
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
// @Summary			Assign outlet to routeCode.
// @Description		Save assign outlet to routeCode data in Db.
// @Param			outlet body request.SaveOutletRequest true "save assign outlet to routeCode data in Db"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/outlets [post]
// @Security        Bearer
func (controller *RouteMappingController) SaveOutlet(ctx *gin.Context) {
	log.Info().Msg("add outlet in route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.SaveOutletRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.SaveOutlet(c, Request)
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
// @Summary			Assign outlet to pjpCode.
// @Description		Save assign outlet to pjpCode data in Db.
// @Param			pjp body request.SavePjpRequest true "save assign outlet to pjpCode data in Db"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/pjp [post]
// @Security        Bearer
func (controller *RouteMappingController) SavePjp(ctx *gin.Context) {
	log.Info().Msg("add outlet in pjp")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.SavePjpRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.SavePjp(c, Request)
	webResponse := response.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created successfully",
		Data:    nil,
	}
	// utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, webResponse)
}

// CreateTags		godoc
// @Summary			Delete outlet in assign route.
// @Description		Delete outlet in assign route.
// @Param			outlet body request.DeleteOutletRequest true "delete outlet in assign route"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/outlets/update [patch]
// @Security        Bearer
func (controller *RouteMappingController) DeleteOutlet(ctx *gin.Context) {
	log.Info().Msg("delete outlet in route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.DeleteOutletRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.DeleteOutlet(c, Request)
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

// CreateTags		godoc
// @Summary			Delete outlet additional in assign route.
// @Description		Delete outlet additional in assign route.
// @Param			outlet body request.DeleteOutletAdditionalRequest true "delete outlet additional in assign route"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/outlets-additional/update [patch]
// @Security        Bearer
func (controller *RouteMappingController) DeleteOutletAdditional(ctx *gin.Context) {
	log.Info().Msg("delete outlet additional in route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.DeleteOutletAdditionalRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.DeleteOutletAdditional(c, Request)
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

// CreateTags		godoc
// @Summary			Update pjp in route outlet.
// @Description		Update pjpId in route route outlet.
// @Param			pjp body request.UpdatePjpInRouteRequest true "update pjp in route outlet"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/pjp [patch]
// @Security        Bearer
func (controller *RouteMappingController) UpdatePjp(ctx *gin.Context) {
	log.Info().Msg("update pjp in route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.UpdatePjpInRouteRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.UpdatePjp(c, Request)
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

// DeleteTags		godoc
// @Summary			Delete route.
// @Description		Remove route data by id.
// @Param		    routeId path string true "remove route data by id"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/{routeId} [delete]
// @Security        Bearer
func (controller *RouteMappingController) Delete(ctx *gin.Context) {
	log.Info().Msg("delete route")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	routeId := ctx.Param("routeId")

	id, err := strconv.Atoi(routeId)
	helper.ErrorPanic(err)

	controller.routeService.Delete(c, id)

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

// FindAllTags 		godoc
// @Summary			Get all route.
// @Description		Return list of route.
// @Produce		    application/json
// @Param		    route_code query string false "route_code"
// @Param		    is_assign query string false "is_assign"
// @Param		    route_name query string false "route_name"
// @Tags			route mapping
// @Success         200 {object} response.Response{}
// @Router			/web/route-mappings [get]
// @Security        Bearer
func (controller *RouteMappingController) FindAll(ctx *gin.Context) {
	log.Info().Msg("findAll list route")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	filters := make(map[string]interface{})
	filters["route_code"] = ctx.Query("route_code")
	filters["is_assign"] = ctx.Query("is_assign")
	filters["route_name"] = ctx.Query("route_name")

	responses := controller.routeService.FindAllRoute(c, filters, currentCustomerId)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAllTags 		    godoc
// @Summary				Get route outlet by routeCode and pjpCode.
// @Param				routeCode path string true "route_code"
// @Param				pjpCode path string true "pjp_code"
// @Description			Return list of route outlet by route_code.
// @Produce				application/json
// @Tags				route mapping
// @Success				200 {object} response.Response{}
// @Router				/web/route-mappings/{routeCode}/{pjpCode} [get]
// @Security            Bearer
func (controller *RouteMappingController) FindByRouteCode(ctx *gin.Context) {
	log.Info().Msg("find by route codes")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	routeCode := ctx.Param("routeCode")
	route_code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	pjpCode := ctx.Param("pjpCode")
	pjp_code, err := strconv.Atoi(pjpCode)
	helper.ErrorPanic(err)

	responses := controller.routeService.FindByRouteOutlet(ctx, route_code, pjp_code)

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
// @Summary			Get all route by pjp_code and route_code.
// @Description		Return list of route by pjp_code and route_code.
// @Param		    pjpCode path string true "pjp_code"
// @Param		    routeCode path string true "route_code"
// @Produce		    application/json
// @Tags			route mapping
// @Success         200 {object} response.Response{}
// @Router			/web/route-mappings/pjp/{pjpCode}/{routeCode} [get]
// @Security        Bearer
func (controller *RouteMappingController) FindRouteByPjpCode(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	pjpCode := ctx.Param("pjpCode")
	pjp_code, err := strconv.Atoi(pjpCode)
	helper.ErrorPanic(err)

	routeCode := ctx.Param("routeCode")
	route_code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	responses := controller.routeService.FindRouteByPjpCode(c, pjp_code, route_code)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateTags		godoc
// @Summary			Update route.
// @Description		Update route data.
// @Param			routeId path string true "update route by id"
// @Param			route body request.UpdateRoutesRequest true  "Update route"
// @Tags			route mapping
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/{routeId} [patch]
// @Security        Bearer
func (controller *RouteMappingController) Update(ctx *gin.Context) {
	log.Info().Msg("update route")
	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.UpdateRoutesRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	ParamID := ctx.Param("routeId")
	id, err := strconv.Atoi(ParamID)
	helper.ErrorPanic(err)
	Request.ID = id

	controller.routeService.UpdateRoute(c, Request, currentCustomerId)

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

// CreateTags		godoc
// @Summary			Save route confirmation.
// @Description		Save route confirmation in Db.
// @Param			raw body request.SaveRouteConfirmationRequest true "save route confirmation"
// @Produce			application/json
// @Tags			route mapping
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/save/route [post]
// @Security        Bearer
func (controller *RouteMappingController) SaveRouteConfirmation(ctx *gin.Context) {
	log.Info().Msg("save route confirmation")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.SaveRouteConfirmationRequest{}

	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.SaveRouteConfirmation(c, Request)
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

// UpdateTags		godoc
// @Summary			Delete route in pjp.
// @Description		Delete route in pjp.
// @Param			route body request.DeletePjpRequest true  "delete pjp code & pjp id in route"
// @Tags			route mapping
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/remove/route [patch]
// @Security        Bearer
func (controller *RouteMappingController) RemoveRouteInPjp(ctx *gin.Context) {
	log.Info().Msg("delete route in pjp")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.DeletePjpRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.DeletePjp(c, Request)

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

// UpdateTags		godoc
// @Summary			New route propose.
// @Description		New route propose.
// @Param			raw body request.NewRouteRequest true "new route with propose"
// @Tags			route mapping
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/new-route [patch]
// @Security        Bearer
func (controller *RouteMappingController) NewRoutePropose(ctx *gin.Context) {
	log.Info().Msg("new route propose")

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.NewRouteRequest{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	controller.routeService.UpdateNewRoute(c, Request)

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

// FindAllTags 		    godoc
// @Summary				Get route outlet by routeCode.
// @Param				routeCode path string true "route_code"
// @Description			Return list of route outlet by route_code.
// @Produce				application/json
// @Tags				route mapping
// @Success				200 {object} response.Response{}
// @Router				/web/route-mappings/{routeCode} [get]
// @Security            Bearer
func (controller *RouteMappingController) FindByRouteOutlet(ctx *gin.Context) {
	log.Info().Msg("find by route codes")
	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	routeCode := ctx.Param("routeCode")
	route_code, err := strconv.Atoi(routeCode)
	helper.ErrorPanic(err)

	responses := controller.routeService.FindByRouteCode(ctx, route_code)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   responses,
	}
	utils.ResponseInterceptor(c, &webResponse)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateTags		godoc
// @Summary			Duplicate route.
// @Description		Duplicate route.
// @Param			raw body request.DuplicateRoute true "duplicate route"
// @Tags			route mapping
// @Produce			application/json
// @Success			200 {object} response.Response{}
// @Router			/web/route-mappings/duplicate [post]
// @Security        Bearer
func (controller *RouteMappingController) RouteDuplicate(ctx *gin.Context) {
	log.Info().Msg("duplicate route")

	currentCustomerId, exists := helper.GetCurrentCustomerId(ctx)
	if !exists {
		helper.ErrorPanic(errors.New("customer_id not found"))
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	Request := request.DuplicateRoute{}
	err := ctx.ShouldBindJSON(&Request)
	helper.ErrorPanic(err)

	err = controller.routeService.DuplicateRoute(c, Request, currentCustomerId)
	if err != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Failed",
			Message: err.Error(),
			Data:    nil,
		}

		utils.ResponseInterceptor(c, &webResponse)
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	webResponse := response.Response{
		Code:    http.StatusOK,
		Status:  "Ok",
		Message: "Duplicate successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
