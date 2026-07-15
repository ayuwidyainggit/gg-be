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

type AnnouncementsController struct {
	AnnouncementsService service.AnnouncementsService
	validator            *validation.Validate
}

func NewAnnouncementsController(
	AnnouncementsService service.AnnouncementsService,
	validator *validation.Validate,
) *AnnouncementsController {
	return &AnnouncementsController{
		AnnouncementsService: AnnouncementsService,
		validator:            validator,
	}
}

func (controller *AnnouncementsController) Route(app *fiber.App) {
	AnnouncementsRouteV1 := app.Group("/v1/announcements", middleware.JWTProtected())
	AnnouncementsRouteV1.Get("/", controller.Announcements)
}

func (controller *AnnouncementsController) Announcements(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.AnnouncementsRequest
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

	data, err := controller.AnnouncementsService.Announcements(request)
	if err != nil {
		log.Error("AnnouncementsSummaryDaily, Detail, FindOneAnnouncementsSummary, err:", err.Error())
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
