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

type PjpDistributorController struct {
	PjpService service.PjpService
	validator  *validation.Validate
}

func NewPjpDistributorController(
	pjpService service.PjpService,
	validator *validation.Validate,
) *PjpDistributorController {
	return &PjpDistributorController{
		PjpService: pjpService,
		validator:  validator,
	}
}

func (controller *PjpDistributorController) Route(app *fiber.App) {
	// Mobile route for submit PJP distributor
	mobilePjpRouteV1 := app.Group("/v1", middleware.JWTProtected())
	mobilePjpRouteV1.Post("/pjp-distributor", controller.SubmitPjpDistributor)
	mobilePjpRouteV1.Put("/pjp-distributor/:pjp_code", controller.UpdatePjpDistributor)
}

func (controller *PjpDistributorController) SubmitPjpDistributor(c *fiber.Ctx) error {
	var (
		request          entity.SubmitPjpDistributorRequest
		headerAcceptLang string
		ctx              = c.UserContext()
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("PjpDistributorController, SubmitPjpDistributor, BodyParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PjpDistributorController, SubmitPjpDistributor, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PjpService.SubmitPjpDistributor(ctx, request, custId)
	if err != nil {
		log.Error("PjpDistributorController, SubmitPjpDistributor, err:", err.Error())
		responsePayload.Setmsg("Failed to create PJP distributor")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Created successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *PjpDistributorController) UpdatePjpDistributor(c *fiber.Ctx) error {
	var (
		payload          entity.UpdatePJPDistributorRequest
		headerAcceptLang string
		custID           = c.Locals("cust_id").(string)
		ctx              = c.UserContext()
	)
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&payload); err != nil {
		log.Error("PjpDistributorController, UpdatePJPDistributor, ParamsParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&payload); err != nil {
		log.Error("PjpDistributorController, UpdatePJPDistributor, BodyParser:", err.Error())
		responsePayload.Setmsg("Invalid request body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	payload.CustomerID = custID
	errs := controller.validator.ValidateStruct(payload, headerAcceptLang)
	if errs != nil {
		log.Error("PjpDistributorController, UpdatePJPDistributor, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PjpService.UpdatePJPDistributor(ctx, payload)
	if err != nil {
		log.Error("PjpDistributorController, UpdatePJPDistributor, err:", err.Error())
		responsePayload.Setmsg("Failed to update PJP distributor")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Updated successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
