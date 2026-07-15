package controller

import (
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/str"
	"sales/pkg/validation"
	"sales/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type SoController struct {
	SoService service.SoService
	validator *validation.Validate
}

func parseDownloadSalesmanIDs(c *fiber.Ctx) []int64 {
	allSalesmanIds := c.Context().QueryArgs().PeekMulti("salesman_id")
	allSalesmanIdsBracket := c.Context().QueryArgs().PeekMulti("salesman_id[]")

	salesmanIDs := make([]int64, 0)
	for _, idBytes := range allSalesmanIds {
		salesmanIDs = appendSalesmanIDValues(salesmanIDs, string(idBytes))
	}
	for _, idBytes := range allSalesmanIdsBracket {
		salesmanIDs = appendSalesmanIDValues(salesmanIDs, string(idBytes))
	}

	if len(salesmanIDs) > 0 {
		return salesmanIDs
	}

	fallbackSalesmanIDs := c.Query("salesman_id")
	if fallbackSalesmanIDs == "" {
		return []int64{}
	}

	return appendSalesmanIDValues(salesmanIDs, fallbackSalesmanIDs)
}

func appendSalesmanIDValues(existing []int64, rawValue string) []int64 {
	if rawValue == "" {
		return existing
	}

	parts := strings.Split(rawValue, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		parsedID, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			continue
		}

		existing = append(existing, parsedID)
	}

	return existing
}

func NewSoController(roService service.SoService, validator *validation.Validate) *SoController {
	return &SoController{
		SoService: roService,
		validator: validator,
	}
}
func (controller *SoController) Route(app *fiber.App) {
	qParamId := ":so_no"
	grRouteV1 := app.Group("/v1/so", middleware.JWTProtected())
	grRouteV1.Post("", controller.Create)
	grRouteV1.Get("/"+qParamId, controller.Detail)
	grRouteV1.Get("", controller.List)
	grRouteV1.Delete("/"+qParamId, controller.Delete)
	grRouteV1.Patch("/"+qParamId, controller.Update)

	// Download Excel endpoint - sesuai dokumentasi di /sales/v1/download
	grRouteDownload := app.Group("/v1", middleware.JWTProtected())
	grRouteDownload.Get("/download", controller.Download)
}
func (controller *SoController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateSoBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("SoController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SoService.Store(request)
	if err != nil {
		log.Error("SoController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
func (controller *SoController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSoParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("SoController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.SoService.Detail(params.SoNo, custId)
	if err != nil {
		log.Error("SoController, Detail, FindOneByOutletId, err:", err.Error())
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

func (controller *SoController) List(c *fiber.Ctx) error {
	var dataFilter entity.SoQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("SoController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.SoService.List(dataFilter)
	if err != nil {
		log.Error("SoController, List, data, err:", err.Error())
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
func (controller *SoController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailSoParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("SoController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.SoService.Delete(custId, params.SoNo, userId)
	if err != nil {
		log.Error("SoController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
func (controller *SoController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSoParams
		request entity.UpdateSoBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("SoController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("SoController, Update, BodyParser(request), err:", err.Error())
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
		log.Error("VanSoController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SoService.Update(params.SoNo, request)
	if err != nil {
		log.Error("VanSoController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SoController) Download(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var dataFilter entity.SoDownloadQueryFilter
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("SoController, Download, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.SalesmanId = parseDownloadSalesmanIDs(c)

	// Validate required fields
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("SoController, Download, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Validate date range (maximum 31 days)
	// Convert epoch time to time.Time for validation
	startDate := str.UnixTimestampToUtcDate(dataFilter.StartDate)
	endDate := str.UnixTimestampToUtcDate(dataFilter.EndDate)

	// Check if end_date is after start_date
	if endDate.Before(startDate) {
		log.Error("SoController, Download, end_date is before start_date")
		responsePayload.Setmsg("end_date must be after or equal to start_date")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Calculate date difference (inclusive: start_date and end_date both count)
	dateDiff := int(endDate.Sub(startDate).Hours() / 24)
	if dateDiff > 31 {
		log.Error("SoController, Download, date range exceeds 31 days, diff:", dateDiff)
		responsePayload.Setmsg("Date range cannot exceed 31 days")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userFullname := c.Locals("user_fullname").(string)

	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId
	dataFilter.ExportBy = userFullname

	data, err := controller.SoService.Download(dataFilter)
	if err != nil {
		log.Error("SoController, Download, Download, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data.FileStatusName = data.GetFileStatusName()
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageTotal:   1,
		PageCurrent: 0,
		PageLimit:   0,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
