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

type RegionController struct {
	RegionService service.RegionService
	validator     *validation.Validate
}

func NewRegionController(regionService service.RegionService, validator *validation.Validate) *RegionController {
	return &RegionController{
		RegionService: regionService,
		validator:     validator,
	}
}

func (controller *RegionController) Route(app *fiber.App) {
	regionRouteV1 := app.Group("/v1/region", middleware.JWTProtected())
	regionRouteV1.Get("", controller.List)
}

func (controller *RegionController) List(c *fiber.Ctx) error {
	var dataFilter entity.RegionQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("RegionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("RegionController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.RegionService.List(c.Context(), dataFilter.CustID)
	if err != nil {
		log.Error("RegionController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
