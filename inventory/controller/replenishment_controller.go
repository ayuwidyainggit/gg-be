package controller

import (
	"fmt"
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ReplenishmentController struct {
	ReplenishmentService service.ReplenishmentService
	validator            *validation.Validate
}

func NewReplenishmentController(replenishmentService service.ReplenishmentService, validator *validation.Validate) *ReplenishmentController {
	return &ReplenishmentController{
		ReplenishmentService: replenishmentService,
		validator:            validator,
	}
}

func (controller *ReplenishmentController) Route(app *fiber.App) {
	replenishmentRouteV1 := app.Group("/v1/replenishment-order", middleware.JWTProtected())
	replenishmentRouteV1.Get("", controller.List)
	replenishmentRouteV1.Post("", controller.Create)
	replenishmentRouteV1.Get("/products", controller.ProductList)
	replenishmentRouteV1.Get("/:replenishment_id", controller.Detail)

	// Product GR List endpoint
	productGrListRouteV1 := app.Group("/v1/product_gr_list", middleware.JWTProtected())
	productGrListRouteV1.Get("", controller.ProductGrList)

	// PO List endpoint
	poListRouteV1 := app.Group("/v1/po_list", middleware.JWTProtected())
	poListRouteV1.Get("", controller.PoList)

	// Replenishment Approval routes
	replenishmentApprovalRouteV1 := app.Group("/v1/replenishment-approval", middleware.JWTProtected())
	replenishmentApprovalRouteV1.Get("", controller.ApprovalList)
	replenishmentApprovalRouteV1.Get("/products", controller.ApprovalProductList)
	replenishmentApprovalRouteV1.Patch("/batch", controller.BatchUpdateApproval)
	replenishmentApprovalRouteV1.Patch("/:replenishment_id", controller.UpdateApproval)

	// Replenishment summarize route
	summarizeReplanishmentRouteV1 := app.Group("/v1/summarize-replanishment", middleware.JWTProtected())
	summarizeReplanishmentRouteV1.Get("", controller.SummarizeReplanishment)
}

// buildResponse builds response payload with accept language header
func (controller *ReplenishmentController) buildResponse(c *fiber.Ctx) *responsebuild.DataRespReq {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
}

// handleServiceError handles service errors and maps them to appropriate HTTP status codes
func (controller *ReplenishmentController) handleServiceError(c *fiber.Ctx, err error, responsePayload *responsebuild.DataRespReq) error {
	statusCode := fiber.StatusBadRequest
	errMsg := err.Error()
	if err.Error() == "sql: no rows in result set" {
		statusCode = fiber.StatusNotFound
		errMsg = "Not found"
	}
	responsePayload.Setmsg(errMsg)
	return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
}

// setPagination sets pagination information in response payload
func (controller *ReplenishmentController) setPagination(responsePayload *responsebuild.DataRespReq, total int64, page, limit, lastPage int) {
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: page,
		PageLimit:   limit,
		PageTotal:   lastPage,
	})
}

// parseCommaSeparatedInt64 parses comma-separated string to []int64
func (controller *ReplenishmentController) parseCommaSeparatedInt64(str string) []int64 {
	if str == "" {
		return []int64{}
	}
	parts := strings.Split(str, ",")
	var result []int64
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if val, err := strconv.ParseInt(part, 10, 64); err == nil {
				result = append(result, val)
			}
		}
	}
	return result
}

func (controller *ReplenishmentController) parseReplanishmentIDs(rawQuery string) []int64 {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return []int64{}
	}

	collector := make([]int64, 0)
	seen := make(map[int64]struct{})
	appendUnique := func(id int64) {
		if id <= 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		collector = append(collector, id)
	}

	for _, key := range []string{"replanishment_id", "replanishment_id[]"} {
		rawVals := values[key]
		for _, rawVal := range rawVals {
			for _, parsed := range controller.parseCommaSeparatedInt64(rawVal) {
				appendUnique(parsed)
			}
		}
	}

	return collector
}

// parseCommaSeparatedInt parses comma-separated string to []int
func (controller *ReplenishmentController) parseCommaSeparatedInt(str string) []int {
	if str == "" {
		return []int{}
	}
	parts := strings.Split(str, ",")
	var result []int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if val, err := strconv.Atoi(part); err == nil {
				result = append(result, val)
			}
		}
	}
	return result
}

func (controller *ReplenishmentController) List(c *fiber.Ctx) error {
	var dataFilter entity.ReplenishmentQueryFilter
	responsePayload := controller.buildResponse(c)

	// Parse query parameters
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit < 1 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	// Parse comma-separated arrays
	dataFilter.SupIDParsed = controller.parseCommaSeparatedInt64(dataFilter.SupID)
	dataFilter.StatusParsed = controller.parseCommaSeparatedInt(dataFilter.Status)

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, List, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Extract cust_id and parent_cust_id from JWT
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	employeeID := int64(0)
	if v := c.Locals("employee_id"); v != nil {
		switch val := v.(type) {
		case int64:
			employeeID = val
		case int:
			employeeID = int64(val)
		case float64:
			employeeID = int64(val)
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				employeeID = parsed
			}
		}
	}

	var distributorIDFromToken *int64
	if v := c.Locals("distributor_id"); v != nil {
		switch val := v.(type) {
		case int64:
			distributorIDFromToken = &val
		case int:
			tmp := int64(val)
			distributorIDFromToken = &tmp
		case float64:
			tmp := int64(val)
			distributorIDFromToken = &tmp
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				distributorIDFromToken = &parsed
			}
		}
	}

	// Check if user is principal
	isPrincipal, err := controller.ReplenishmentService.CheckIsPrincipal(custId)
	if err != nil {
		log.Error("ReplenishmentController, List, CheckIsPrincipal, err:", err.Error())
		responsePayload.Setmsg("failed to check principal status")
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	// Set in filter
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId
	dataFilter.IsPrincipal = isPrincipal
	dataFilter.UserID = userId
	dataFilter.EmpID = employeeID
	dataFilter.DistributorIDFromToken = distributorIDFromToken

	// Call service
	data, total, lastPage, err := controller.ReplenishmentService.List(dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		log.Error("ReplenishmentController, List, Service.List, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("success")
		responsePayload.Setdata(data)
	}
	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) ProductGrList(c *fiber.Ctx) error {
	var dataFilter entity.ProductGrListQueryFilter
	responsePayload := controller.buildResponse(c)

	// Parse query params
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, ProductGrList, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Inject cust IDs from middleware
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, ProductGrList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Service call
	data, total, lastPage, err := controller.ReplenishmentService.ProductGrList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ReplenishmentController, ProductGrList, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if data == nil || len(data.Details) == 0 {
		responsePayload.Setmsg("")
		responsePayload.Data = nil
	} else {
		responsePayload.Setdata(data)
		responsePayload.Setmsg("success")
	}

	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) Create(c *fiber.Ctx) error {
	var request entity.CreateReplenishmentOrderBody
	responsePayload := controller.buildResponse(c)

	// Parse request body
	if err := c.BodyParser(&request); err != nil {
		log.Error("ReplenishmentController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Extract cust_id, parent_cust_id, and user_id from JWT
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	employeeID := int64(0)
	if v := c.Locals("employee_id"); v != nil {
		switch val := v.(type) {
		case int64:
			employeeID = val
		case int:
			employeeID = int64(val)
		case float64:
			employeeID = int64(val)
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				employeeID = parsed
			}
		}
	}

	// Set values from JWT
	request.CustID = custId
	request.ParentCustID = parentCustId
	request.CreatedBy = userId
	request.CreatedEmpID = employeeID

	if request.DistributorID == nil || *request.DistributorID <= 0 {
		if v := c.Locals("distributor_id"); v != nil {
			switch val := v.(type) {
			case int64:
				if val > 0 {
					tmp := val
					request.DistributorID = &tmp
				}
			case int:
				if val > 0 {
					tmp := int64(val)
					request.DistributorID = &tmp
				}
			case float64:
				if val > 0 {
					tmp := int64(val)
					request.DistributorID = &tmp
				}
			case string:
				if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil && parsed > 0 {
					tmp := parsed
					request.DistributorID = &tmp
				}
			}
		}
	}

	// Validate struct
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Validate unique pro_id in data array
	var detailsProIds entity.ReplenishmentDetailsProductId
	for _, prod := range request.Data {
		detailsProIds.ProductIds = append(detailsProIds.ProductIds, entity.ProductId{
			Product: entity.Product{
				Id: prod.ProID,
			},
		})
	}
	errs = controller.validator.ValidateStruct(detailsProIds, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, Create, Detail Product ID ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Call service
	err := controller.ReplenishmentService.Store(request)
	if err != nil {
		log.Error("ReplenishmentController, Create, Service.Store, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	responsePayload.Setmsg("Product succsessfully aded")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) Detail(c *fiber.Ctx) error {
	var params entity.DetailReplenishmentParams
	responsePayload := controller.buildResponse(c)

	// Parse path and query parameters
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ReplenishmentController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.QueryParser(&params); err != nil {
		log.Error("ReplenishmentController, Detail, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, Detail, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Extract cust_id and parent_cust_id from JWT
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	// Check if user is principal
	isPrincipal, err := controller.ReplenishmentService.CheckIsPrincipal(custId)
	if err != nil {
		log.Error("ReplenishmentController, Detail, CheckIsPrincipal, err:", err.Error())
		responsePayload.Setmsg("failed to check principal status")
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	// Call service
	data, err := controller.ReplenishmentService.Detail(params.ReplenishmentID, params.Type, params.Status, custId, parentCustId, isPrincipal)
	if err != nil {
		log.Error("ReplenishmentController, Detail, Service.Detail, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if data == nil {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("success")
		responsePayload.Setdata(data)
	}
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) ProductList(c *fiber.Ctx) error {
	var dataFilter entity.ReplenishmentProductQueryFilter
	responsePayload := controller.buildResponse(c)

	// Parse query parameters
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, ProductList, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit < 1 {
		dataFilter.Limit = constant.DEFAULT_PAGE_LIMIT
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, ProductList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Extract cust_id and parent_cust_id from JWT
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	// Set in filter
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId

	// Call service
	data, total, lastPage, err := controller.ReplenishmentService.ProductList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ReplenishmentController, ProductList, Service.ProductList, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("success")
		responsePayload.Setdata(data)
	}
	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) PoList(c *fiber.Ctx) error {
	var dataFilter entity.PoListQueryFilter
	responsePayload := controller.buildResponse(c)

	// Parse query params (page and limit)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, PoList, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Inject cust IDs from middleware
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, PoList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Service call
	data, total, lastPage, err := controller.ReplenishmentService.PoList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ReplenishmentController, PoList, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if len(data) == 0 {
		responsePayload.Setmsg("")
		responsePayload.Data = nil
	} else {
		responsePayload.Setdata(data)
		responsePayload.Setmsg("success")
	}

	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) ApprovalList(c *fiber.Ctx) error {
	var dataFilter entity.ReplenishmentApprovalListQueryFilter
	responsePayload := controller.buildResponse(c)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, ApprovalList, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	dataFilter.StatusParsed = controller.parseCommaSeparatedInt(dataFilter.Status)

	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, ApprovalList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	dataFilter.UserID = userId

	var empID int64
	if v := c.Locals("employee_id"); v != nil {
		switch e := v.(type) {
		case int64:
			empID = e
		case int:
			empID = int64(e)
		case float64:
			empID = int64(e)
		}
	}
	dataFilter.EmpID = empID

	var distributorIDFromToken *int64
	if v := c.Locals("distributor_id"); v != nil {
		switch val := v.(type) {
		case int64:
			distributorIDFromToken = &val
		case int:
			tmp := int64(val)
			distributorIDFromToken = &tmp
		case float64:
			tmp := int64(val)
			distributorIDFromToken = &tmp
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				distributorIDFromToken = &parsed
			}
		}
	}

	isPrincipal, err := controller.ReplenishmentService.CheckIsPrincipal(custId)
	if err != nil {
		log.Error("ReplenishmentController, ApprovalList, CheckIsPrincipal, err:", err.Error())
		responsePayload.Setmsg("failed to check principal status")
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId
	dataFilter.IsPrincipal = isPrincipal
	dataFilter.DistributorIDFromToken = distributorIDFromToken

	data, total, lastPage, err := controller.ReplenishmentService.ApprovalList(dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		log.Error("ReplenishmentController, ApprovalList, Service.ApprovalList, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	if len(data) == 0 {
		responsePayload.Setmsg("")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("success")
		responsePayload.Setdata(data)
	}
	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) ApprovalProductList(c *fiber.Ctx) error {
	var dataFilter entity.ReplenishmentApprovalProductQueryFilter
	responsePayload := controller.buildResponse(c)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ReplenishmentController, ApprovalProductList, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId
	dataFilter.UserID = userID

	errs := controller.validator.ValidateStruct(dataFilter, "")
	if errs != nil {
		log.Error("ReplenishmentController, ApprovalProductList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Check if user is principal
	isPrincipal, err := controller.ReplenishmentService.CheckIsPrincipal(custId)
	if err != nil {
		log.Error("ReplenishmentController, ApprovalProductList, CheckIsPrincipal, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.ReplenishmentService.ApprovalProductList(dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		log.Error("ReplenishmentController, ApprovalProductList, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	responsePayload.Setdata(data)
	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)
	responsePayload.Setmsg("success")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) UpdateApproval(c *fiber.Ctx) error {
	var request entity.UpdateReplenishmentApprovalRequest
	responsePayload := controller.buildResponse(c)

	replenishmentIDStr := c.Params("replenishment_id")
	if strings.TrimSpace(replenishmentIDStr) == "" {
		responsePayload.Setmsg("replenishment_id is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	replenishmentID, errParse := strconv.ParseInt(strings.TrimSpace(replenishmentIDStr), 10, 64)
	if errParse != nil || replenishmentID <= 0 {
		responsePayload.Setmsg("invalid replenishment_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Parse request body
	if err := c.BodyParser(&request); err != nil {
		log.Error("ReplenishmentController, UpdateApproval, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	employeeID := int64(0)
	if v := c.Locals("employee_id"); v != nil {
		switch val := v.(type) {
		case int64:
			employeeID = val
		case int:
			employeeID = int64(val)
		case float64:
			employeeID = int64(val)
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				employeeID = parsed
			}
		}
	}

	// Validate
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, UpdateApproval, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Service call
	err := controller.ReplenishmentService.UpdateApproval(replenishmentID, request, custId, userId, employeeID)
	if err != nil {
		log.Error("ReplenishmentController, UpdateApproval, err:", err.Error())
		if strings.Contains(strings.ToLower(err.Error()), "only setup replenishment pic is allowed to approve replenishment data") {
			responsePayload.Setmsg("Bad Request")
			responsePayload.Seterrors([]map[string]string{
				{
					"key":     "user_role",
					"message": "Only setup replenishment PIC is allowed to approve replenishment data",
				},
			})
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		return controller.handleServiceError(c, err, responsePayload)
	}

	// Set success message based on approval status
	if request.Approval != nil && *request.Approval {
		responsePayload.Setmsg("Replenishment data has been successfully approved.")
	} else {
		responsePayload.Setmsg("Replenishment data has been rejected.")
	}
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) BatchUpdateApproval(c *fiber.Ctx) error {
	var request entity.BatchReplenishmentApprovalRequest
	responsePayload := controller.buildResponse(c)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ReplenishmentController, BatchUpdateApproval, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	userID := c.Locals("user_id").(int64)
	employeeID := int64(0)
	if v := c.Locals("employee_id"); v != nil {
		switch val := v.(type) {
		case int64:
			employeeID = val
		case int:
			employeeID = int64(val)
		case float64:
			employeeID = int64(val)
		case string:
			if parsed, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				employeeID = parsed
			}
		}
	}

	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, BatchUpdateApproval, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	isPIC, err := controller.ReplenishmentService.IsAnyApprovalPICUser(userID, employeeID)
	if err != nil {
		log.Error("ReplenishmentController, BatchUpdateApproval, IsAnyApprovalPICUser, err:", err.Error())
		responsePayload.Setmsg("failed to verify PIC role")
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	if !isPIC {
		responsePayload.Setmsg("Forbidden")
		responsePayload.Seterrors([]map[string]string{
			{
				"key":     "user_role",
				"message": "Only PIC (Principal) is allowed to approve replenishment data",
			},
		})
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ReplenishmentService.BatchUpdateApproval(request, custID, userID, employeeID)
	if err != nil {
		log.Error("ReplenishmentController, BatchUpdateApproval, err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors([]map[string]interface{}{
			{"key": "batch", "message": err.Error()},
		})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Batch replenishment processed successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReplenishmentController) SummarizeReplanishment(c *fiber.Ctx) error {
	responsePayload := controller.buildResponse(c)

	replanishmentIDs := controller.parseReplanishmentIDs(string(c.Context().URI().QueryString()))
	if len(replanishmentIDs) == 0 {
		responsePayload.Setmsg("Bad Request")
		responsePayload.Seterrors([]map[string]string{
			{
				"key":     "replanishment_id",
				"message": "replanishment_id is required",
			},
		})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	queryPayload := entity.SummarizeReplanishmentQuery{ReplanishmentID: replanishmentIDs}
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(queryPayload, headerAcceptLang)
	if errs != nil {
		log.Error("ReplenishmentController, SummarizeReplanishment, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	requestID, ok := c.Locals("requestid").(string)
	if !ok || strings.TrimSpace(requestID) == "" {
		requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	}

	isPrincipal, err := controller.ReplenishmentService.CheckIsPrincipal(custID)
	if err != nil {
		log.Error("ReplenishmentController, SummarizeReplanishment, CheckIsPrincipal, err:", err.Error())
		responsePayload.Setmsg("failed to check principal status")
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ReplenishmentService.SummarizeReplanishment(replanishmentIDs, custID, parentCustID, isPrincipal, requestID, userID)
	if err != nil {
		log.Error("ReplenishmentController, SummarizeReplanishment, service err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	responsePayload.Setmsg("Success get summarize replanishment")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
