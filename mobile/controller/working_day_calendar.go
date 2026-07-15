package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WorkingDayCalendarController struct {
	WorkingDayCalendarService service.WorkingDayCalendarService
	validator                 *validation.Validate
}

func NewWorkingDayCalendarController(
	workingDayCalendarService service.WorkingDayCalendarService,
	validator *validation.Validate,
) *WorkingDayCalendarController {
	return &WorkingDayCalendarController{
		WorkingDayCalendarService: workingDayCalendarService,
		validator:                 validator,
	}
}

func (controller *WorkingDayCalendarController) Route(app *fiber.App) {
	workingDayCalendarRouteV1 := app.Group("/v1/working-days", middleware.JWTProtected())
	workingDayCalendarRouteV1.Get("", controller.List)
	workingDayCalendarRouteV1.Get("/month", controller.ListMonths)
}

func (controller *WorkingDayCalendarController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.WorkingDayCalendarQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("WorkingDayCalendarController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("WorkingDayCalendarController, List, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.WorkingDayCalendarService.List(dataFilter)
	if err != nil {
		log.Error("WorkingDayCalendarController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
		responsePayload.Setdata(nil)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WorkingDayCalendarController) ListMonths(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.WorkingDayCalendarMonthQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("WorkingDayCalendarController, ListMonths, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("WorkingDayCalendarController, ListMonths, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.WorkingDayCalendarService.ListMonths(dataFilter)
	if err != nil {
		log.Error("WorkingDayCalendarController, ListMonths, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
		responsePayload.Setdata(nil)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
