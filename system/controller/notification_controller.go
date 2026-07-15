package controller

import (
	"system/entity"
	"system/pkg/constant"
	"system/pkg/responsebuild"
	"system/pkg/validation"
	"system/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	NotificationService service.NotificationService
	validator           *validation.Validate
}

func NewNotificationController(userService service.NotificationService, validator *validation.Validate) *NotificationController {
	return &NotificationController{
		NotificationService: userService,
		validator:           validator,
	}
}

func (controller *NotificationController) Route(app *fiber.App) {
	app.Post("v1/notifications/whatsapp-cicd", controller.WhatsappCicd)
}

func (controller *NotificationController) WhatsappCicd(c *fiber.Ctx) error {
	var request entity.NotifyCicdWaReq

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	err := c.BodyParser(&request)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("errs, NotificationController, WhatsappCicd:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.NotificationService.WhatsappCicd(request)
	if err != nil {
		log.Error("NotificationController, WhatsappCicd, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}
