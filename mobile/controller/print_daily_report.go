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

type PrintDailyReportController struct {
	PrintDailyReportService service.PrintDailyReportService
	validator               *validation.Validate
}

func NewPrintDailyReportController(
	printDailyReportService service.PrintDailyReportService,
	validator *validation.Validate,
) *PrintDailyReportController {
	return &PrintDailyReportController{
		PrintDailyReportService: printDailyReportService,
		validator:               validator,
	}
}

func (controller *PrintDailyReportController) Route(app *fiber.App) {
	printDailyReportRouteV1 := app.Group("/v1/print_daily_report", middleware.JWTProtected())
	printDailyReportRouteV1.Get("/", controller.GetDailyReport)
}

func (controller *PrintDailyReportController) GetDailyReport(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.PrintDailyReportRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	requestID, ok := c.Locals("requestid").(string)
	if !ok {
		requestID = ""
	}
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	// Parse query parameters
	if err := c.QueryParser(&request); err != nil {
		log.Error("PrintDailyReportController, GetDailyReport, QueryParser, err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Validate request
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PrintDailyReportController, GetDailyReport, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Extract JWT context
	custID, ok := c.Locals("cust_id").(string)
	if !ok {
		responsePayload.Setmsg("cust_id not found in token")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		responsePayload.Setmsg("user_id not found in token")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}

	empID, ok := c.Locals("emp_id").(int64)
	if !ok {
		responsePayload.Setmsg("emp_id not found in token")
		return c.Status(fiber.StatusUnauthorized).JSON(responsePayload.GetRespPayload())
	}
	request.EmpID = empID

	// Call service
	data, err := controller.PrintDailyReportService.GetDailyReport(request, custID, userID)
	if err != nil {
		log.Error("PrintDailyReportController, GetDailyReport, GetDailyReport, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()

		// Handle not found error
		if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "record not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
