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

type LeaderboardsController struct {
	LeaderboardsService service.LeaderboardsService
	validator           *validation.Validate
}

func NewLeaderboardsController(
	LeaderboardsService service.LeaderboardsService,
	validator *validation.Validate,
) *LeaderboardsController {
	return &LeaderboardsController{
		LeaderboardsService: LeaderboardsService,
		validator:           validator,
	}
}

func (controller *LeaderboardsController) Route(app *fiber.App) {
	LeaderboardsRouteV1 := app.Group("/v1/leaderboards", middleware.JWTProtected())
	LeaderboardsRouteV1.Get("/", controller.Leaderboards)
}

func (controller *LeaderboardsController) Leaderboards(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.LeaderboardsRequest
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

	data, err := controller.LeaderboardsService.Leaderboards(request)
	if err != nil {
		log.Error("LeaderboardsSummaryDaily, Detail, FindOneLeaderboardsSummary, err:", err.Error())
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
