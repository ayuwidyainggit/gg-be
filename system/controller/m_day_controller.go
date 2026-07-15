package controller

import (
	"system/entity"
	"system/pkg/constant"
	"system/pkg/middleware"
	"system/pkg/responsebuild"
	"system/pkg/validation"
	"system/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type MDayController struct {
	MDayService service.MDayService
	validator   *validation.Validate
}

func NewMDayController(mDayService service.MDayService, validator *validation.Validate) *MDayController {
	return &MDayController{
		MDayService: mDayService,
		validator:   validator,
	}
}

func (controller *MDayController) Route(app *fiber.App) {
	qParamId := ":day_id"
	RouteV1 := app.Group("/v1/m-days", middleware.JWTProtected())
	RouteV1.Get("", controller.List)
	RouteV1.Get("/"+qParamId, controller.Detail)
}

func (controller *MDayController) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("MDayController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	langId := headerAcceptLang
	if langId == "" {
		langId = c.Locals("user_lang").(string)
	}

	data, total, lastPage, err := controller.MDayService.List(dataFilter, langId)
	if err != nil {
		log.Error("MDayController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MDayController) Detail(c *fiber.Ctx) error {
	var params entity.DetailMDayBodyParam
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("UserController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("UserController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	langId := headerAcceptLang
	if langId == "" {
		langId = c.Locals("user_lang").(string)
	}

	data, err := controller.MDayService.Detail(params.DayId, langId)
	if err != nil {
		log.Error("UserController, Detail, FindOneByOutletId, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}
