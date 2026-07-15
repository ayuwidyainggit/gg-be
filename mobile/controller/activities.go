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

type ActivitiesController struct {
	ActivitiesService service.ActivitiesService
	validator         *validation.Validate
}

func NewActivitiesController(
	ActivitiesService service.ActivitiesService,
	validator *validation.Validate,
) *ActivitiesController {
	return &ActivitiesController{
		ActivitiesService: ActivitiesService,
		validator:         validator,
	}
}

func (controller *ActivitiesController) Route(app *fiber.App) {
	ActivitiesRouteV1 := app.Group("/v1/activities", middleware.JWTProtected())
	ActivitiesRouteV1.Get("/summary/daily", controller.Summary)
}

func (controller *ActivitiesController) Summary(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.SummaryDailyRequest
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

	EmpId := c.Locals("emp_id").(int64)
	request.CustId = c.Locals("cust_id").(string)
	request.EmployeeId = EmpId

	data, err := controller.ActivitiesService.ActivitiesSummaryDaily(request)
	if err != nil {
		log.Error("ActivitiesSummaryDaily, Detail, FindOneActivitiesSummary, err:", err.Error())
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
