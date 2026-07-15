package controller

import (
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type WhTrfController struct {
	whTrfService service.WhTrfService
	validator    *validation.Validate
}

func NewWhTrfController(whTrfService service.WhTrfService, validator *validation.Validate) *WhTrfController {
	return &WhTrfController{
		whTrfService: whTrfService,
		validator:    validator,
	}
}

func (controller *WhTrfController) Route(app *fiber.App) {
	qParamId := ":wh_trf_no"
	qParamWhMethod := ":method"
	grRouteV1 := app.Group("/v1/stock-transfers", middleware.JWTProtected())
	grRouteV1.Post("", controller.Create)
	grRouteV1.Get("/"+qParamId, controller.Detail)
	grRouteV1.Get("", controller.List)
	grRouteV1.Delete("/"+qParamId, controller.Delete)
	grRouteV1.Patch("/"+qParamId, controller.Update)
	grRouteV1.Get("/warehouse/list/"+qParamWhMethod, controller.ListWarehouse)
}
func (controller *WhTrfController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateWhTrfBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("WarehouseTrfController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("TprController, Create, CustId:", custId)

	request.CustID = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.whTrfService.Store(request)
	if err != nil {
		log.Error("WarehouseTrfController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
func (controller *WhTrfController) Detail(c *fiber.Ctx) error {
	var params entity.DetailWhTrfParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("WarehouseTrfController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	distributorID := c.Locals("distributor_id").(int64)
	tokenString := c.Get("Authorization")

	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.whTrfService.Detail(params.WhTrfNo, custId, parentCustId, distributorID, tokenString)
	if err != nil {
		log.Error("WarehouseTrfController, Detail, err:", err.Error())
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
func (controller *WhTrfController) List(c *fiber.Ctx) error {
	var dataFilter entity.WhQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("WarehouseTrfController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("BpprController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.whTrfService.List(dataFilter, custId)
	if err != nil {
		log.Error("WarehouseTrfController, List, data, err:", err.Error())
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
func (controller *WhTrfController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailWhTrfParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("WarehouseTrfController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.whTrfService.Delete(custId, params.WhTrfNo, userId)
	if err != nil {
		log.Error("WarehouseTrfController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
func (controller *WhTrfController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateWhTrfParams
		request entity.UpdateWhTrfBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("WarehouseTrfController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("WarehouseTrfController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustID = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.whTrfService.Update(params.WhTrfNo, request)
	if err != nil {
		log.Error("WarehouseTrfController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *WhTrfController) ListWarehouse(c *fiber.Ctx) error {
	var dataFilter entity.StockTranferWarehouseQueryFilter
	var params entity.StockTranferWarehouseQueryFilterParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("WarehouseTrfController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("WarehouseTrfController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("WarehouseTrfController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)
	dataFilter.Method = params.Method
	data, total, lastPage, err := controller.whTrfService.ListWarehouse(dataFilter, custId)
	if err != nil {
		log.Error("WarehouseTrfController, List, data, err:", err.Error())
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
