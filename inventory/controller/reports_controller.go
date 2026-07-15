package controller

import (
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ReportsController struct {
	ReportsService service.ReportsService
	validator      *validation.Validate
}

func NewReportsController(reportsService service.ReportsService, validator *validation.Validate) *ReportsController {
	return &ReportsController{
		ReportsService: reportsService,
		validator:      validator,
	}
}

func (controller *ReportsController) Route(app *fiber.App) {
	reportsRouteV1 := app.Group("/v1/reports", middleware.JWTProtected())
	reportsRouteV1.Get("/stock_movememt", controller.StockMovement)
	reportsRouteV1.Get("/stock_movement/download", controller.DownloadStockMovement)
	reportsRouteV1.Get("/stock_movement/download/preview", controller.PreviewDownloadStockMovement)
}

func (controller *ReportsController) StockMovement(c *fiber.Ctx) error {
	var dataFilter entity.StockMovementReportQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReportsController, StockMovement, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReportsController, StockMovement, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)

	now := time.Now()
	if dataFilter.Month == 0 {
		dataFilter.Month = int(now.Month())
	}
	if dataFilter.Year == 0 {
		dataFilter.Year = now.Year()
	}

	data, err := controller.ReportsService.GetStockMovementReport(dataFilter)
	if err != nil {
		log.Error("ReportsController, StockMovement, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if data == nil || (len(data.WhTotalStock) == 0 && len(data.StockMovement) == 0 && len(data.TopProductIn) == 0 && len(data.TopProductOut) == 0) {
		responsePayload.Setmsg(constant.DATA_NOT_FOUND)
		responsePayload.Setdata(nil)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.SUCCESS)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportsController) PreviewDownloadStockMovement(c *fiber.Ctx) error {
	var dataFilter entity.PreviewDownloadStockMovementReportQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReportsController, StockMovement, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReportsController, StockMovement, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)
	dataFilter.UserID = c.Locals("user_id").(int64)
	if fullname, ok := c.Locals("user_fullname").(string); ok {
		dataFilter.UserFullName = fullname
	}

	data, total, err := controller.ReportsService.PreviewDownloadStockMovementReport(dataFilter)
	if err != nil {
		log.Error("ReportsController, StockMovement, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	page := dataFilter.Page
	if page <= 0 {
		page = 1
	}

	limit := dataFilter.Limit
	pageTotal := 1
	if limit > 0 {
		pageTotal = int((total + int64(limit) - 1) / int64(limit))
		if pageTotal == 0 {
			pageTotal = 1
		}

	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: page,
		PageLimit:   limit,
		PageTotal:   pageTotal,
	})
	responsePayload.Setmsg(constant.SUCCESS)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReportsController) DownloadStockMovement(c *fiber.Ctx) error {
	var dataFilter entity.DownloadStockMovementReportQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReportsController, StockMovement, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReportsController, StockMovement, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)
	dataFilter.UserID = c.Locals("user_id").(int64)
	if fullname, ok := c.Locals("user_fullname").(string); ok {
		dataFilter.UserFullName = fullname
	}

	data, err := controller.ReportsService.DownloadStockMovementReport(dataFilter)
	if err != nil {
		log.Error("ReportsController, StockMovement, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageCurrent: 0,
		PageLimit:   0,
		PageTotal:   1,
	})
	responsePayload.Setmsg(constant.SUCCESS)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
