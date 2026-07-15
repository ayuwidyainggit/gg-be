package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type PjpController struct {
	PjpService service.PjpService
	validator  *validation.Validate
}

func NewPjpController(
	pjpService service.PjpService,
	validator *validation.Validate,
) *PjpController {
	return &PjpController{
		PjpService: pjpService,
		validator:  validator,
	}
}

func (controller *PjpController) Route(app *fiber.App) {
	pjpRouteV1 := app.Group("/v1/pjp", middleware.JWTProtected())
	pjpRouteV1.Get("/salesman", controller.SalesmanDetail)
}

func (controller *PjpController) SalesmanDetail(c *fiber.Ctx) error {
	var queryParams entity.PjpSalesmanQueryParams
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&queryParams); err != nil {
		log.Error("PjpController, SalesmanDetail, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(queryParams, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("PjpController, SalesmanDetail, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	isDistributor := c.Locals("is_distributor").(bool)
	data, err := controller.PjpService.GetSalesmanDetail(queryParams.EmpId, isDistributor)
	if err != nil {
		log.Error("PjpController, SalesmanDetail, err:", err.Error())

		if err.Error() == "No Data" {
			responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
			responsePayload.Setdata(nil)
			return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
		}

		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PjpController) getAcceptLanguage(c *fiber.Ctx) string {
	acceptLang := c.Get("Accept-Language", "id")
	if acceptLang == "" {
		acceptLang = "id"
	}
	return acceptLang
}
