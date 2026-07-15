package controller

import (
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/str"
	"inventory/pkg/validation"
	"inventory/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
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
	stockRouteV1 := app.Group("/v1/stocks", middleware.JWTProtected())
	stockRouteV1.Get("", controller.List)
	stockRouteV1.Get("/opname-lookup", controller.OpnameLookup)
	stockRouteV1.Get("report", controller.Report)
	stockRouteV1.Post("", controller.Create)
	stockRouteV1.Post("/bulk", controller.CreateBulk)
}

func (controller *StockController) List(c *fiber.Ctx) error {
	var dataFilter entity.StockQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.StockService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("StockController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateStock
	if err := c.BodyParser(&request); err != nil {
		log.Error("StockController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)

	request.CustID = custID
	request.ParentCustID = parentCustID

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.StockService.Store(request)
	if err != nil {
		log.Error("StockController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *StockController) CreateBulk(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateBulkStock
	if err := c.BodyParser(&request); err != nil {
		log.Error("StockController, CreateBulk, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)

	for i, _ := range request.Stock {
		request.Stock[i].CustID = custID
		request.Stock[i].ParentCustID = parentCustID
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockController, CreateBulk, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.StockService.StoreBulk(request)
	if err != nil {
		log.Error("StockController, CreateBulk, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *StockController) Report(c *fiber.Ctx) error {
	var dataFilter entity.StockReportQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockController, Report, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockController, Report, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.StockService.Report(dataFilter)
	if err != nil {
		log.Error("StockController, Report, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.SetFilter(dataFilter)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockController) OpnameLookup(c *fiber.Ctx) error {
	var dataFilter entity.StockOpnameLookupQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockController, OpnameLookup, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockController, OpnameLookup, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)

	timeNowJkt := str.GetJakartaDate()
	dataFilter.Date = timeNowJkt.Format("2006-01-02")

	data, total, lastPage, err := controller.StockService.OpnameLookup(dataFilter)
	if err != nil {
		log.Error("StockController, OpnameLookup, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: 1,
		PageLimit:   0,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
