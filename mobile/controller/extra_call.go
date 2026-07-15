package controller

import (
	"github.com/gofiber/fiber/v2/log"
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
)

type ExtraCallController struct {
	ExtraCallService service.ExtraCallService
	validator        *validation.Validate
}

func NewExtraCallController(
	ExtraCallService service.ExtraCallService,
	validator *validation.Validate,
) *ExtraCallController {
	return &ExtraCallController{
		ExtraCallService: ExtraCallService,
		validator:        validator,
	}
}

func (controller *ExtraCallController) Route(app *fiber.App) {
	EventsRouteV1 := app.Group("/v1/extra-call", middleware.JWTProtected())
	EventsRouteV1.Post("/", controller.Create)
}

func (controller *ExtraCallController) Create(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.CreateExtraCallRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ExtraCall, Validate params, err:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.CustID = c.Locals("cust_id").(string)
	empID := c.Locals("emp_id").(int64)
	request.EmpID = int(empID)
	request.IsDistributor = c.Locals("is_distributor").(bool)

	err := controller.ExtraCallService.Create(request)
	if err != nil {
		log.Error("ExtraCall, Create (service), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
