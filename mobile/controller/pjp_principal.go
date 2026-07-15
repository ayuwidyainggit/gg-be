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

type PjpPrincipalController struct {
	PjpPrincipalService service.PjpPrincipalService
	validator           *validation.Validate
}

func NewPjpPrincipalController(
	pjpPrincipalService service.PjpPrincipalService,
	validator *validation.Validate,
) *PjpPrincipalController {
	return &PjpPrincipalController{
		PjpPrincipalService: pjpPrincipalService,
		validator:           validator,
	}
}

func (controller *PjpPrincipalController) Route(app *fiber.App) {
	// Mobile route for submit PJP principal
	mobilePjpRouteV1 := app.Group("/v1", middleware.JWTProtected())
	mobilePjpRouteV1.Post("/pjp-principal", controller.SubmitPjpPrincipal)
	mobilePjpRouteV1.Put("/pjp-principal/:pjp_code", controller.UpdatePjpPrincipal)
}

func (controller *PjpPrincipalController) SubmitPjpPrincipal(c *fiber.Ctx) error {
	var (
		request          entity.SubmitPjpPrincipalRequest
		headerAcceptLang string
		ctx              = c.UserContext()
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	requestID, ok := c.Locals("requestid").(string)
	if !ok {
		requestID = ""
	}
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("PjpPrincipalController, SubmitPjpPrincipal, BodyParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId, ok := c.Locals("cust_id").(string)
	if !ok {
		responsePayload.Setmsg("Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	isDistributor := c.Locals("is_distributor").(bool)
	if isDistributor {
		responsePayload.Setmsg("Unauthorized, you are not a principal")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PjpPrincipalController, SubmitPjpPrincipal, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PjpPrincipalService.SubmitPjpPrincipal(ctx, request, custId)
	if err != nil {
		log.Error("PjpPrincipalController, SubmitPjpPrincipal, err:", err.Error())
		responsePayload.Setmsg("Failed to create PJP principal")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Created successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *PjpPrincipalController) UpdatePjpPrincipal(c *fiber.Ctx) error {
	var (
		payload          entity.UpdatePjpPrincipalRequest
		headerAcceptLang string
		custID           = c.Locals("cust_id").(string)
		ctx              = c.UserContext()
	)
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	requestID, ok := c.Locals("requestid").(string)
	if !ok {
		requestID = ""
	}
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.ParamsParser(&payload); err != nil {
		log.Error("PjpPrincipalController, UpdatePjpPrincipal, ParamsParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Error("PjpPrincipalController, UpdatePjpPrincipal, BodyParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	payload.CustomerID = custID
	errs := controller.validator.ValidateStruct(payload, headerAcceptLang)
	if errs != nil {
		log.Error("PjpPrincipalController, UpdatePjpPrincipal, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PjpPrincipalService.UpdatePjpPrincipal(ctx, payload)
	if err != nil {
		log.Error("PjpPrincipalController, UpdatePjpPrincipal, err:", err.Error())
		responsePayload.Setmsg("Failed to update PJP principal")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Updated successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
