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

type EventsController struct {
	EventsService service.EventsService
	validator     *validation.Validate
}

func NewEventsController(
	EventsService service.EventsService,
	validator *validation.Validate,
) *EventsController {
	return &EventsController{
		EventsService: EventsService,
		validator:     validator,
	}
}

func (controller *EventsController) Route(app *fiber.App) {
	EventsRouteV1 := app.Group("/v1/events", middleware.JWTProtected())
	EventsRouteV1.Get("/", controller.Events)
}

func (controller *EventsController) Events(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.EventsRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {

		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.EventsService.Events(request)
	if err != nil {
		log.Error("EventsSummaryDaily, Detail, FindOneEventsSummary, err:", err.Error())
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
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
