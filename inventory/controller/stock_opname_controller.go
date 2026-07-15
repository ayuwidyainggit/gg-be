package controller

import (
	"encoding/json"
	"errors"
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type StockOpnameController struct {
	StockOpnameService service.StockOpnameService
	validator          *validation.Validate
}

func NewStockOpnameController(StockOpnameService service.StockOpnameService, validator *validation.Validate) *StockOpnameController {
	return &StockOpnameController{
		StockOpnameService: StockOpnameService,
		validator:          validator,
	}
}

func (controller *StockOpnameController) Route(app *fiber.App) {
	qParamId := ":doc_no"
	stockOpnameRouteV1 := app.Group("/v1/stock-opname", middleware.JWTProtected())
	stockOpnameRouteV1.Get("", controller.List)
	stockOpnameRouteV1.Get("product-hierarchy", controller.GetProductHierarchy)
	stockOpnameRouteV1.Get("statuses", controller.GetStatuses)
	stockOpnameRouteV1.Get("/"+qParamId, controller.Report)
	stockOpnameRouteV1.Post("", controller.Create)
	stockOpnameRouteV1.Patch("/"+qParamId+"/cancel", controller.Cancel)

	// V2 routes Stock Opname (underscore URL)
	stockOpnameRouteV2 := app.Group("/v1/stock_opname", middleware.JWTProtected())
	stockOpnameRouteV2.Get("", controller.ListV2)
	stockOpnameRouteV2.Post("", controller.CreateV2)
	stockOpnameRouteV2.Get("/product-list", controller.ProductList)
	stockOpnameRouteV2.Get("/template", controller.DownloadTemplate) // Must be before /:doc_no route
	stockOpnameRouteV2.Get("/download", controller.DownloadReport)   // Download final stock opname report
	stockOpnameRouteV2.Post("/bulk_upload/"+qParamId, controller.BulkUpload)
	stockOpnameRouteV2.Get("/"+qParamId, controller.DetailV2)
	stockOpnameRouteV2.Patch("/"+qParamId, controller.UpdateStatusV2)
	stockOpnameRouteV2.Patch("/revised/"+qParamId, controller.RevisedV2)
	stockOpnameRouteV2.Patch("/start/"+qParamId, controller.StartV2)
	stockOpnameRouteV2.Patch("/submit/"+qParamId, controller.SubmitV2)
}

func (controller *StockOpnameController) Report(c *fiber.Ctx) error {
	var params entity.ReportStockOpanmeParams

	// var dataFilter entity.StockOpnameQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, Report, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, Report, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)

	data, err := controller.StockOpnameService.Reports(params)
	if err != nil {
		log.Error("StockOpnameController, Report, err:", err.Error())
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

func (controller *StockOpnameController) List(c *fiber.Ctx) error {
	var dataFilter entity.StockOpnameQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockOpnameController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("StockOpnameController, List, custId, parentCustId:", custId, parentCustId)

	dataFilter.CustID = custId
	dataFilter.ParentCustID = parentCustId

	data, total, lastPage, err := controller.StockOpnameService.List(dataFilter)
	if err != nil {
		log.Error("StockOpnameController, List, data, err:", err.Error())
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

func (controller *StockOpnameController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateStockOpname
	if err := c.BodyParser(&request); err != nil {
		log.Error("StockOpnameController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("StockOpnameController, Create, CustId:", custId)

	request.CustID = custId
	request.ParentCustID = parentCustId
	request.CreatedBy = c.Locals("user_id").(int64)
	request.UpdatedBy = request.CreatedBy

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.StockOpnameService.Store(request)
	if err != nil {
		log.Error("StockOpnameController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.DATA_CREATED_SUCCESSFULLY)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) CreateV2(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateStockOpnameV2
	if err := c.BodyParser(&request); err != nil {
		log.Error("StockOpnameController, CreateV2, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.ParentCustID = parentCustId
	request.CreatedBy = userId

	if len(request.PrincipalID) == 0 && len(request.PLLane) == 0 && len(request.BrandID) == 0 && len(request.SBrand1ID) == 0 {
		responsePayload.Setmsg(constant.STOCK_OPNAME_FILTER_REQUIRED)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, CreateV2, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.StockOpnameService.StoreV2(request)
	if err != nil {
		log.Error("StockOpnameController, CreateV2, StoreV2, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_SAVED_SUCCESSFULLY)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) GetStatuses(c *fiber.Ctx) error {
	stockOpnameStatuses := make([]entity.StockOpnameStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.StockOpnameStatusDesc {
		stockOpnameStatus := entity.StockOpnameStatus{
			StatusID:   index,
			StatusDesc: element,
		}
		stockOpnameStatuses = append(stockOpnameStatuses, stockOpnameStatus)
	}

	statusesSorted := make(entity.StockOpnameStatusDescSlice, 0)
	for _, row := range stockOpnameStatuses {
		statusesSorted = append(statusesSorted, row)
	}
	sort.Sort(statusesSorted)
	responsePayload.Setdata(statusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) GetProductHierarchy(c *fiber.Ctx) error {
	productHierarchies := make([]entity.ProductHierarchy, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.ProductHierarchyDesc {
		productHierarchy := entity.ProductHierarchy{
			ProductHierarchyID:   index,
			ProductHierarchyDesc: element,
		}
		productHierarchies = append(productHierarchies, productHierarchy)
	}

	productHierarchySorted := make(entity.ProductHierarchyDescSlice, 0)
	for _, row := range productHierarchies {
		productHierarchySorted = append(productHierarchySorted, row)
	}
	sort.Sort(productHierarchySorted)
	responsePayload.Setdata(productHierarchySorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) Cancel(c *fiber.Ctx) error {
	var params entity.CancelStockOpanmeParams

	// var dataFilter entity.StockOpnameQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, Cancel, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, Cancel, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.CancelBy = c.Locals("user_id").(int64)

	err := controller.StockOpnameService.Cancel(params)
	if err != nil {
		log.Error("StockOpnameController, Cancel, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) ListV2(c *fiber.Ctx) error {
	var dataFilter entity.StockOpnameListV2QueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockOpnameController, ListV2, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set defaults
	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	isAdmin := c.Locals("is_admin").(bool)

	dataFilter.CustID = custId
	dataFilter.ParentCustID = parentCustId
	dataFilter.UserID = userId
	dataFilter.IsAdmin = isAdmin

	data, total, lastPage, err := controller.StockOpnameService.ListV2(dataFilter)
	if err != nil {
		log.Error("StockOpnameController, ListV2, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_DISPLAYED_SUCCESSFULLY)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) DetailV2(c *fiber.Ctx) error {
	var params entity.StockOpnameDetailV2Params
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, DetailV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, DetailV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)

	data, err := controller.StockOpnameService.DetailV2(params)
	if err != nil {
		log.Error("StockOpnameController, DetailV2, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			statusCode = fiber.StatusNotFound
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_DISPLAYED_SUCCESSFULLY)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) UpdateStatusV2(c *fiber.Ctx) error {
	var params entity.UpdateStockOpnameStatusV2Params
	var request entity.UpdateStockOpnameStatusV2Request
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, UpdateStatusV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("StockOpnameController, UpdateStatusV2, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, UpdateStatusV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, UpdateStatusV2, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	err := controller.StockOpnameService.UpdateStatusV2(params, request)
	if err != nil {
		log.Error("StockOpnameController, UpdateStatusV2, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()

		if errMsg == "record not found" {
			statusCode = fiber.StatusNotFound
		} else if errMsg == "you are not authorized to process or assign this stock opname. Only the creator can perform this action" {
			statusCode = fiber.StatusForbidden
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_UPDATED_SUCCESSFULLY)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) ProductList(c *fiber.Ctx) error {
	var dataFilter entity.StockOpnameProductListQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockOpnameController, ProductList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	dataFilter.CustID = custId
	dataFilter.ParentCustID = parentCustId

	if len(dataFilter.PrincipalID) == 0 && len(dataFilter.PLID) == 0 && len(dataFilter.BrandID) == 0 && len(dataFilter.SBrand1ID) == 0 {
		responsePayload.Setmsg(constant.STOCK_OPNAME_FILTER_REQUIRED_V2)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.StockOpnameService.ProductList(dataFilter)
	if err != nil {
		log.Error("StockOpnameController, ProductList, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_DISPLAYED_SUCCESSFULLY)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) RevisedV2(c *fiber.Ctx) error {
	var params entity.RevisedStockOpnameV2Params
	var request entity.RevisedStockOpnameV2Request
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, RevisedV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("StockOpnameController, RevisedV2, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, RevisedV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, RevisedV2, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	err := controller.StockOpnameService.RevisedV2(params, request)
	if err != nil {
		log.Error("StockOpnameController, RevisedV2, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()

		if errMsg == "record not found" || errMsg == "stock opname detail not found" {
			statusCode = fiber.StatusNotFound
		} else if errMsg == "you are not authorized to revised this stock opname. Only the creator can perform this action" {
			statusCode = fiber.StatusForbidden
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_UPDATED_SUCCESSFULLY)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) DownloadTemplate(c *fiber.Ctx) error {
	var params entity.StockOpnameTemplateDownloadParams
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&params); err != nil {
		log.Error("StockOpnameController, DownloadTemplate, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if params.DocNo == "" {
		params.DocNo = c.Query("doc_no")
	}
	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)

	data, err := controller.StockOpnameService.DownloadTemplate(params)
	if err != nil {
		log.Error("StockOpnameController, DownloadTemplate, err:", err.Error())
		statusCode := fiber.StatusBadRequest

		if errors.Is(err, constant.ErrRecordNotFound) {
			statusCode = fiber.StatusNotFound
		}

		responsePayload.Setmsg(constant.STOCK_OPNAME_TEMPLATE_DOWNLOAD_FAILED)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STOCK_OPNAME_TEMPLATE_DOWNLOAD_SUCCESS)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) DownloadReport(c *fiber.Ctx) error {
	var params entity.StockOpnameDownloadParams
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&params); err != nil {
		log.Error("StockOpnameController, DownloadReport, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if params.DocNo == "" {
		params.DocNo = c.Query("doc_no")
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	if params.DocNo == "" {
		responsePayload.Setmsg("doc_no is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.StockOpnameService.DownloadReport(params)
	if err != nil {
		log.Error("StockOpnameController, DownloadReport, err:", err.Error())
		statusCode := fiber.StatusBadRequest

		if errors.Is(err, constant.ErrRecordNotFound) {
			statusCode = fiber.StatusNotFound
		}

		responsePayload.Setmsg("Failed to download stock opname report")
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Stock opname report downloaded successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) BulkUpload(c *fiber.Ctx) error {
	var params entity.BulkUploadStockOpnameV2Params
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, BulkUpload, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	file, err := c.FormFile("file")
	if err != nil || file == nil {
		log.Error("StockOpnameController, BulkUpload, form file:", err)
		responsePayload.Setmsg("file is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	fileReader, err := file.Open()
	if err != nil {
		log.Error("StockOpnameController, BulkUpload, file open:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	defer fileReader.Close()

	fileBytes := make([]byte, file.Size)
	_, err = fileReader.Read(fileBytes)
	if err != nil {
		log.Error("StockOpnameController, BulkUpload, file read:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.StockOpnameService.BulkUpload(params, fileBytes, file.Filename)
	if err != nil {
		log.Error("StockOpnameController, BulkUpload, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		// Return specific status for known errors
		if strings.Contains(err.Error(), "invalid file format") {
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if strings.Contains(err.Error(), "exceeds maximum size") {
			responsePayload.Setdata(fiber.Map{"max_size_mb": 100})
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if strings.Contains(err.Error(), "cannot be updated in current status") {
			responsePayload.Setdata(fiber.Map{"status": "Completed"})
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Bulk stock opname upload processed successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) StartV2(c *fiber.Ctx) error {
	var params entity.StartStockOpnameV2Params
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, StartV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, StartV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	data, err := controller.StockOpnameService.StartV2(params)
	if err != nil {
		log.Error("StockOpnameController, StartV2, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()

		if errors.Is(err, constant.ErrRecordNotFound) {
			statusCode = fiber.StatusNotFound
		} else if errors.Is(err, constant.ErrStockOpnameCannotBeStarted) {
			statusCode = fiber.StatusBadRequest
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.DATA_STARTED_SUCCESSFULLY)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockOpnameController) SubmitV2(c *fiber.Ctx) error {
	var params entity.SubmitStockOpnameV2Params
	var request entity.SubmitStockOpnameV2Request
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockOpnameController, SubmitV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("StockOpnameController, SubmitV2, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Get file from form-data if available
	fileHeader, _ := c.FormFile("file")

	if len(request.Details) == 0 {
		itemsStr := c.FormValue("items")
		if itemsStr != "" {
			if err := json.Unmarshal([]byte(itemsStr), &request.Details); err != nil {
				log.Error("StockOpnameController, SubmitV2, json.Unmarshal(items):", err.Error())
				responsePayload.Setmsg("Invalid items JSON: " + err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
			}
		}
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, SubmitV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("StockOpnameController, SubmitV2, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UserID = c.Locals("user_id").(int64)

	data, err := controller.StockOpnameService.SubmitV2(params, request, fileHeader)
	if err != nil {
		log.Error("StockOpnameController, SubmitV2, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()

		if errors.Is(err, constant.ErrRecordNotFound) {
			statusCode = fiber.StatusNotFound
		} else if strings.Contains(errMsg, "not found in document") {
			statusCode = fiber.StatusBadRequest
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Stock opname details updated successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
