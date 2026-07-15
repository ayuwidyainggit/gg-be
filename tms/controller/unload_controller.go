package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UnloadController struct {
	unloadService service.UnloadService
}

func NewUnloadController(unloadService service.UnloadService) *UnloadController {
	return &UnloadController{
		unloadService: unloadService,
	}
}

// Note             godoc
//
//	@Summary		todo list
//	@Description	todo list arrive, unload, leave.
//	@Param			outletId	path	string	true	"OutletID"
//	@Param			shipmentNo	path	string	true	"ShipmentNo"
//	@Produce		application/json
//	@Tags			reject
//	@Success		200	{object}	entity.JsonSuccess{data=entity.TravelListResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}									"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}									"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}						"Internal server error"
//	@Router			/mobile/todo/list/{outletId}/{shipmentNo} [get]
func (controller *UnloadController) TravelList(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.TravelListParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.unloadService.TravelList(c, params)

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
//	@Summary		unload
//	@Description	update unload_at to db.
//	@Param			data	body	entity.UnloadRequest	true	"unload"
//	@Produce		application/json
//	@Tags			unload
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/unload [post]
//
// @Security	Bearer
func (controller *UnloadController) Unload(ctx *fiber.Ctx) error {
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

	request := entity.UnloadRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.unloadService.Unload(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Unload Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
