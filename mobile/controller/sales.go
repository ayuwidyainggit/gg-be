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

type SalesController struct {
	SalesService service.SalesService
	validator    *validation.Validate
}

func NewSalesController(
	salesService service.SalesService,
	validator *validation.Validate,
) *SalesController {
	return &SalesController{
		SalesService: salesService,
		validator:    validator,
	}
}

// Route registers sales routes with JWT protection middleware
func (controller *SalesController) Route(app *fiber.App) {
	salesRouteV1 := app.Group("/v1/sales", middleware.JWTProtected())
	salesRouteV1.Get("/summary", controller.Summary)
}

// Summary handles GET /v1/sales/summary endpoint.
// Extracts cust_id and emp_id from JWT context and retrieves sales summary data.
// Returns current_sales (total order - total return) and daily_target.
func (controller *SalesController) Summary(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.SalesSummaryRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Extract cust_id and emp_id from JWT context
	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)

	data, err := controller.SalesService.SalesSummary(request, custId, empId)
	if err != nil {
		log.Error("SalesSummaryController, Detail, FindOneSalesSummary, err:", err.Error())
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
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
