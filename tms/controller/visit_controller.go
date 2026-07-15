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

type VisitController struct {
	visitService service.VisitService
}

func NewVisitController(service service.VisitService) *VisitController {
	return &VisitController{
		visitService: service,
	}
}

// Note                godoc
//
//	@Summary		Get summary by driver_id, cust_id.
//	@Param			driverId	path	int		true	"Driver ID"
//	@Param			custId		path	string	true	"Customer ID"
//	@Description	Return the summary.
//	@Produce		application/json
//	@Tags			activities
//	@Success		200	{object}	entity.JsonSuccess{data=entity.SummaryResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/visits/summary/{driverId}/{custId} [get]
func (controller *VisitController) GetSummaryByDriverIDAndCustID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.SummaryParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.visitService.GetSummary(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		Get summary daily by shipment_no.
//	@Param			shipmentNo	path	string	true	"shipment_no"
//	@Param			custId  	path	string	true	"cust_id"
//	@Description	Return the summary daily by shipment_no.
//	@Produce		application/json
//	@Tags			activities
//	@Success		200	{object}	entity.JsonSuccess{data=entity.SummaryDailyResponse{}}	    "Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}										"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}										"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}							"Internal server error"
//	@Router			/mobile/visits/daily/{shipmentNo}/{custId} [get]
func (controller *VisitController) GetSummaryDailyByParams(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.SummaryDailyParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.visitService.GetSummaryDailyByParams(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}

	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		Get daily activity by driver_id.
//	@Param			driverId	path	string	true	"Driver ID"
//	@Description	Return the daily activity is current_date.
//	@Produce		application/json
//	@Tags			activities
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.DailyActivityResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}										"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}										"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}							"Internal server error"
//	@Router			/mobile/visits/daily-activity/{driverId} [get]
func (controller *VisitController) GetDailyActivityByDriverID(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.DailyActivityParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.visitService.GetDailyActivity(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		start visit.
//	@Description	start visit in Db.
//	@Param			visit	body	entity.VisitRequest	true	"start visit and fill in current_time, you can use it from https://currentmillis.com/"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/visits/start [post]
func (controller *VisitController) Start(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.VisitRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.visitService.Start(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Start Visit Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		end visit.
//	@Description	end visit in Db.
//	@Param			visit	body	entity.VisitRequest	true	"end visit and fill in current_time, you can use it from https://currentmillis.com/"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/visits/end [post]
func (controller *VisitController) End(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.VisitRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.visitService.End(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "End Visit Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		leave
//	@Description	update leave_at to db.
//	@Param			data	body	entity.LeaveRequest	true	"leave"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/leave [post]
func (controller *VisitController) Leave(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.LeaveRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.visitService.Leave(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Leave Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		arrive
//	@Description	update arrive_at to db.
//	@Param			data	body	entity.ArriveRequest	true	"arrive"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/arrive [post]
func (controller *VisitController) Arrive(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.ArriveRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.visitService.Arrive(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Arrive Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		skip
//	@Description	update skip_at to db.
//	@Param			data	body	entity.SkipRequest	true	"skip"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}																	"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}																	"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}														"Internal server error"
//	@Router			/mobile/skip [post]
func (controller *VisitController) Skip(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.SkipRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.visitService.Skip(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Skip Successfully",
		Data:    nil,
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
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/unload [post]
//
// @Security	Bearer
func (controller *VisitController) Unload(ctx *fiber.Ctx) error {
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

	controller.visitService.Unload(c, headers, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Unload Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		resume
//	@Description	update resume_at to db.
//	@Param			data	body	entity.UnloadRequest	true	"resume_at"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/resume [post]
//
// @Security	Bearer
func (controller *UnloadController) Resume(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	// authHeader := ctx.Get("Authorization")
	// if authHeader == "" {
	// 	webResponse := entity.Response{
	// 		Code:    http.StatusUnauthorized,
	// 		Status:  "UNAUTHORIZED",
	// 		Message: "empty token",
	// 	}
	// 	utils.ResponseInterceptor(c, &webResponse)
	// 	return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	// }

	request := entity.UnloadRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	// headers := make(map[string]string)
	// headers["Authorization"] = ctx.Get("Authorization")
	// headers["Accept"] = ctx.Get("application/json")

	controller.unloadService.Resume(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Resume Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		onhold
//	@Description	update on_hold to db.
//	@Param			data	body	entity.UnloadRequest	true	"on_hold"
//	@Produce		application/json
//	@Tags			visit
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/mobile/onhold [post]
//
// @Security	Bearer
func (controller *UnloadController) OnHold(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	// authHeader := ctx.Get("Authorization")
	// if authHeader == "" {
	// 	webResponse := entity.Response{
	// 		Code:    http.StatusUnauthorized,
	// 		Status:  "UNAUTHORIZED",
	// 		Message: "empty token",
	// 	}
	// 	utils.ResponseInterceptor(c, &webResponse)
	// 	return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	// }

	request := entity.UnloadRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	// headers := make(map[string]string)
	// headers["Authorization"] = ctx.Get("Authorization")
	// headers["Accept"] = ctx.Get("application/json")

	controller.unloadService.Onhold(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Onhold Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
