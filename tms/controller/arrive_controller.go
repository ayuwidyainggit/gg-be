package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"time"
)

type ArriveController struct {
	arriveService service.ArriveService
}

func NewArriveController(arriveService service.ArriveService) *ArriveController {
	return &ArriveController{
		arriveService: arriveService,
	}
}

// Note             godoc
//
//	@Summary		arrive
//	@Description	update arrive_at to db.
//	@Param			data	body	entity.ArriveRequest	true	"arrive"
//	@Produce		application/json
//	@Tags			arrive
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/arrive [post]
func (controller *ArriveController) Arrive(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.ArriveRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.arriveService.Arrive(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Arrive Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
