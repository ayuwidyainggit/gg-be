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

type AreaController struct {
	AreaService service.AreaService
	validator   *validation.Validate
}

func NewAreaController(areaService service.AreaService, validator *validation.Validate) *AreaController {
	return &AreaController{
		AreaService: areaService,
		validator:   validator,
	}
}

func (controller *AreaController) Route(app *fiber.App) {
	areaRouteV1 := app.Group("/v1/area", middleware.JWTProtected())
	areaRouteV1.Get("", controller.List)
}

func (controller *AreaController) List(c *fiber.Ctx) error {
	var dataFilter entity.AreaQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("AreaController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("AreaController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.AreaService.List(c.Context(), dataFilter.CustID)
	if err != nil {
		log.Error("AreaController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
