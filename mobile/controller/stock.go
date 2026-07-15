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

type StockController struct {
	StockService service.StockService
	validator    *validation.Validate
}

func NewStockController(stockService service.StockService, validator *validation.Validate) *StockController {
	return &StockController{
		StockService: stockService,
		validator:    validator,
	}
}

func (controller *StockController) Route(app *fiber.App) {
	stocksRouteV1 := app.Group("/v1/stocks", middleware.JWTProtected())
	stocksRouteV1.Get("/gudang-utama", controller.ListGudangUtama)
	stocksRouteV1.Get("/gudang-canvas", controller.ListGudangCanvas)
}

func (controller *StockController) ListGudangUtama(c *fiber.Ctx) error {
	var dataFilter entity.StockQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("StockController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	CustId := c.Locals("cust_id").(string)
	EmpId := c.Locals("emp_id").(int64)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.StockService.ListGudangUtama(dataFilter, EmpId, CustId)
	if err != nil {
		log.Error("OrderController, Detail, FindOneByOutletId, err:", err.Error())
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
	return c.JSON(responsePayload.GetRespPayload())

}

func (controller *StockController) ListGudangCanvas(c *fiber.Ctx) error {
	var dataFilter entity.StockQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("StockController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	CustId := c.Locals("cust_id").(string)
	EmpId := c.Locals("emp_id").(int64)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.StockService.ListGudangCanvas(dataFilter, EmpId, CustId)
	if err != nil {
		log.Error("OrderController, Detail, FindOneByOutletId, err:", err.Error())
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
	return c.JSON(responsePayload.GetRespPayload())

}
