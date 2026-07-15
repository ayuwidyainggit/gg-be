package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DmsController struct {
	dmsService     service.DmsService
	vehicleService service.VehicleService
}

func NewDmsController(dmsService service.DmsService, vehicleService service.VehicleService) *DmsController {
	return &DmsController{
		dmsService:     dmsService,
		vehicleService: vehicleService,
	}
}

// Note             godoc
//
//	@Summary		Get All vehicles via mapping dms.
//	@Description	Return list of vehicles via mapping dms.
//	@Produce		application/json
//	@Param			limit		query	string	false	"Limit"
//	@Param			page		query	string	false	"Page"
//	@Param			is_active	query	string	false	"IsActive"
//	@Tags			vehicle
//	@Success		200	{object}	entity.Response{data=[]entity.VehicleResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}							"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}							"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}				"Internal server error"
//	@Router			/vehicles [get]
func (controller *DmsController) GetVehicleByDms(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.GeneralQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	dataFilter.Page, _ = strconv.Atoi(ctx.Query("page"))
	dataFilter.IsActive, _ = strconv.Atoi(ctx.Query("is_active"))

	response, paging, err := controller.dmsService.GetVehicleByDms(c, dataFilter)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get All vehicles.
//	@Description	Return list of vehicles.
//	@Produce		application/json
//	@Param			limit			query	string	false	"limit"
//	@Param			page			query	string	false	"page"
//	@Param			delivery_date	query	string	false	"delivery_date"
//	@Param			sort			query	string	false	"sort"
//	@Tags			vehicle
//	@Success		200	{object}	entity.Response{data=[]entity.VehicleResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}							"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}							"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}				"Internal server error"
//	@Router			/vehicles/dev [get]
func (controller *DmsController) GetVehicle(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.VehicleQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	dataFilter.Page, _ = strconv.Atoi(ctx.Query("page"))
	dataFilter.DeliveryDate = ctx.Query("delivery_date")
	dataFilter.Sort = ctx.Query("sort")

	response, paging, err := controller.vehicleService.GetVehicle(c, dataFilter)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get All reject reason.
//	@Description	Return list of reject reason.
//	@Produce		application/json
//	@Param			limit		query	string	false	"Limit"
//	@Param			page		query	string	false	"Page"
//	@Param			is_active	query	string	false	"IsActive"
//	@Tags			reason
//	@Success		200	{object}	entity.Response{data=[]entity.ReasonRejectResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/reject-reason [get]
func (controller *DmsController) GetRejectReason(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.GeneralQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	dataFilter.Page, _ = strconv.Atoi(ctx.Query("page"))
	dataFilter.IsActive, _ = strconv.Atoi(ctx.Query("is_active"))

	response, paging, err := controller.dmsService.GetRejectReason(c, dataFilter)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get list of returns.
//	@Description	Return list of returns.
//	@Produce		application/json
//	@Param			limit		query	string	false	"Limit"
//	@Param			outlet_id   query	string	false	"OutletId"
//	@Tags			invoice
//	@Success		200	{object}	entity.Response{data=[]entity.CustomShipmentInvoice}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		401	{object}	entity.JsonUnauthorized{}								"Unauthorized"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/invoices [get]
//
// @Security	Bearer
func (controller *DmsController) GetListInvoice(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	var dataFilter entity.GeneralQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	outletId := ctx.Query("outlet_id")
	dataFilter.OutletId = outletId

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	response, paging, err := controller.dmsService.GetListInvoice(c, dataFilter, headers)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get list of returns.
//	@Description	Return list of returns.
//	@Produce		application/json
//	@Param			limit		query	string	false	"Limit"
//	@Param			outlet_id   query	string	false	"OutletId"
//	@Tags			return
//	@Success		200	{object}	entity.Response{data=[]entity.ResponseReturn}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		401	{object}	entity.JsonUnauthorized{}								"Unauthorized"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/returns [get]
//
// @Security	Bearer
func (controller *DmsController) GetListReturn(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	var dataFilter entity.GeneralQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	// dataFilter.OutletId, _ = strconv.Atoi(ctx.Query("outlet_id"))
	outletId := ctx.Query("outlet_id")
	dataFilter.OutletId = outletId

	// if outletId != "" {
	// 	dataFilter.OutletId = strings.Split(outletId, ",")
	// }

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	response, paging, err := controller.dmsService.GetReturns(c, dataFilter, headers)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
