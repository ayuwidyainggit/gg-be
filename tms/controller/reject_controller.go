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

type RejectController struct {
	rejectService service.RejectService
}

func NewRejectController(rejectService service.RejectService) *RejectController {
	return &RejectController{
		rejectService: rejectService,
	}
}

// Note             godoc
//
//	@Summary		Get all reject.
//	@Description	Return list of reject.
//	@Param			reason_id		query	string	false	"reason_id"
//	@Param			driver_id		query	string	false	"driver_id"
//	@Param			outlet_id		query	string	false	"outlet_id"
//	@Param			product_name	query	string	false	"product_name"
//	@Param			sort			query	string	false	"Sort"
//	@Param			shipment_no			query	string	false	"shipment_no"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.RejectResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/rejects [get]
func (controller *RejectController) GetReject(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	// TODO

	var dataFilter entity.RejectQueryFilter

	dataFilter.ReasonID, _ = strconv.Atoi(ctx.Query("reason_id"))
	dataFilter.DriverID, _ = strconv.Atoi(ctx.Query("driver_id"))
	dataFilter.OutletID, _ = strconv.Atoi(ctx.Query("outlet_id"))
	dataFilter.ProductName = ctx.Query("product_name")
	dataFilter.Sort = ctx.Query("sort")
	dataFilter.ShipmentNo = ctx.Query("shipment_no")

	data := controller.rejectService.GetReject(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get all reject partial.
//	@Description	Return list of reject partial.
//	@Param			reason_id		query	string	false	"reason_id"
//	@Param			driver_id		query	string	false	"driver_id"
//	@Param			outlet_id		query	string	false	"outlet_id"
//	@Param			product_name	query	string	false	"product_name"
//	@Param			sort			query	string	false	"Sort"
//	@Param			shipment_no			query	string	false	"shipment_no"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.RejectPartialResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/rejects/partial [get]
func (controller *RejectController) GetRejectPartial(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.RejectQueryFilter

	dataFilter.ReasonID, _ = strconv.Atoi(ctx.Query("reason_id"))
	dataFilter.DriverID, _ = strconv.Atoi(ctx.Query("driver_id"))
	dataFilter.OutletID, _ = strconv.Atoi(ctx.Query("outlet_id"))
	dataFilter.ProductName = ctx.Query("product_name")
	dataFilter.Sort = ctx.Query("sort")
	dataFilter.ShipmentNo = ctx.Query("shipment_no")

	data := controller.rejectService.GetRejectPartial(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		reject all
//	@Description	save reject all to db.
//	@Param			data	body	entity.RejectRequest	true	"reject all"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/rejects [post]
//
// @Security	Bearer
func (controller *RejectController) RejectAll(ctx *fiber.Ctx) error {
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

	request := entity.RejectRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.rejectService.RejectAll(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Reject All Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		reject partial
//	@Description	save reject partial to db.
//	@Param			data	body	entity.RejectPartialRequest	true	"reject partial"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/rejects/partial [post]
//
// @Security	Bearer
func (controller *RejectController) RejectPartial(ctx *fiber.Ctx) error {
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

	request := entity.RejectPartialRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.rejectService.RejectPartial(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Reject Partial Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		reject cancel
//	@Description	reject cancel to db.
//	@Param			data	body	entity.RejectCancelRequest	true	"reject cancel"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/rejects/cancel [post]
func (controller *RejectController) RejectCancel(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.RejectCancelRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.rejectService.RejectCancel(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Reject Cancel Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
