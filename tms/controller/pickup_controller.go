package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PickUpController struct {
	pickUpService service.PickUpService
}

func NewPickUpController(pickUpService service.PickUpService) *PickUpController {
	return &PickUpController{
		pickUpService: pickUpService,
	}
}

// Note             godoc
//
//	@Summary		pickup all
//	@Description	save pickup all to db.
//	@Param			data	body	entity.PickUpRequest	true	"pickup all"
//	@Produce		application/json
//	@Tags			pickup
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/pickup [post]
func (controller *PickUpController) PickUpAll(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.PickUpRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.pickUpService.PickUpAll(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "PickUp All Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		skip pickup
//	@Description	save skip pickup to db.
//	@Param			data	body	entity.SkipPickUpRequest	true	"skip pickup"
//	@Produce		application/json
//	@Tags			skip pickup
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/pickup/skip [post]
func (controller *PickUpController) SkipPickUp(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.SkipPickUpRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.pickUpService.SkipPickUp(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Skip PickUp Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
// @Summary		pickup partial
// @Description	save pickup partial to db.
// @Param		data	body	entity.PickupPartialRequest	true	"pickup partial"
// @Produce		application/json
// @Tags		pickup
// @Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
// @Router		/mobile/pickup/partial [post]
func (controller *PickUpController) PickUpPartial(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.PickupPartialRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.pickUpService.PickUpPartial(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "PickUp Partial Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
