package controller

import (
	"fmt"
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type EmployeeController struct {
	EmployeeService service.EmployeeService
	validator       *validation.Validate
}

func NewEmployeeController(empGroupService service.EmployeeService, validator *validation.Validate) EmployeeController {
	return EmployeeController{
		EmployeeService: empGroupService,
		validator:       validator,
	}
}

func (controller *EmployeeController) Route(app *fiber.App) {
	qParamId := ":emp_id"
	// employeeRouteV2 := app.Group("/create-multiple", middleware.JWTProtected())
	// employeeRouteV2.Post("", controller.CreateMultiple)

	employeeRouteV1 := app.Group("/v1/employees", middleware.JWTProtected())
	employeeRouteV1.Get("/export", controller.Export)
	employeeRouteV1.Get("/export-template", controller.ExportTemplate)
	employeeRouteV1.Get("/export-template-update", controller.ExportTemplateUpdate)
	employeeRouteV1.Post("/import", controller.Import)
	employeeRouteV1.Post("/import-update", controller.ImportUpdate)
	employeeRouteV1.Get("/"+qParamId, controller.Detail)
	employeeRouteV1.Get("", controller.List)
	employeeRouteV1.Post("", controller.Create)
	employeeRouteV1.Patch("/"+qParamId, controller.Update)
	employeeRouteV1.Delete("/"+qParamId, controller.Delete)
	employeeRouteV1.Post("/create-multiple", controller.CreateMultiple)
	employeeRouteV1.Get("/list/without-salesman", controller.ListWithoutSalesman)

	qParamId = ":emp_id"
	// NEW ROUTE - must be registered before param routes
	employeePJPRoute := app.Group("/v1/employee-pjp", middleware.JWTProtected())
	employeePJPRoute.Get("", controller.ListPJP)

	employeeLookupRoute := app.Group("/v1/employee-lookup", middleware.JWTProtected())
	employeeLookupRoute.Get("", controller.EmployeeLookup)
}

// EmployeeLookup GET /v1/employee-lookup — employees with a row in sys.m_user; paging; optional cust_id filter.
func (controller *EmployeeController) EmployeeLookup(c *fiber.Ctx) error {
	var q entity.EmployeeLookupAPIQuery
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&q); err != nil {
		log.Println("EmployeeController, EmployeeLookup, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	filter := entity.EmployeeLookupAPIFilter{
		CustId:        c.Locals("cust_id").(string),
		ParentCustId:  c.Locals("parent_cust_id").(string),
		Query:         q.Q,
		Page:          q.Page,
		Limit:         q.Limit,
		Sort:          q.Sort,
		FilterCustIds: queryStringSliceFromRequest(c, "cust_id"),
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 5
	}
	if filter.Sort == "" {
		filter.Sort = "created_date:desc"
	}

	data, total, lastPage, err := controller.EmployeeService.LookupAPI(filter)
	if err != nil {
		log.Println("EmployeeController, EmployeeLookup, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if data == nil {
		data = []entity.EmployeeLookupMinimalItem{}
	}

	responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func queryStringSliceFromRequest(c *fiber.Ctx, key string) []string {
	var out []string
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		if string(k) != key {
			return
		}
		s := strings.TrimSpace(string(v))
		if s == "" {
			return
		}
		if strings.Contains(s, ",") {
			for _, p := range strings.Split(s, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					out = append(out, p)
				}
			}
			return
		}
		out = append(out, s)
	})
	return out
}

func (controller *EmployeeController) ListPJP(c *fiber.Ctx) error {
	var dataFilter entity.EmployeePJPQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.EmployeeService.ListPJP(dataFilter)
	if err != nil {
		if total == 0 {
			responsePayload.Setmsg("No Data")
			responsePayload.Setpaging(entity.Pagination{
				TotalRecord: 0,
				PageCurrent: dataFilter.Page,
				PageLimit:   dataFilter.Limit,
				PageTotal:   0,
			})
			return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
		}
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

func (controller *EmployeeController) List(c *fiber.Ctx) error {
	var (
		err             error
		dataFilter      entity.EmployeeQueryFilter
		data            interface{}
		total           int
		lastPage        int
		employees       []entity.EmployeeResponse
		employeesLookup []entity.EmployeeLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("EmployeeController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.EmployeeService.LookupList(dataFilter)
		if err != nil {
			log.Println("EmployeeController, Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employeesLookup, &data)
		if err != nil {
			log.Println("EmployeeController, Automapper Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.EmployeeService.List(dataFilter)
		if err != nil {
			log.Println("EmployeeController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employees, &data)
		if err != nil {
			log.Println("EmployeeController, Automapper List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
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

func (controller *EmployeeController) ListWithoutSalesman(c *fiber.Ctx) error {

	var (
		err             error
		dataFilter      entity.EmployeeQueryFilter
		data            interface{}
		total           int
		lastPage        int
		employees       []entity.EmployeeResponse
		employeesLookup []entity.EmployeeLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("EmployeeController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.EmployeeService.LookupListWithoutSalesman(dataFilter)
		if err != nil {
			log.Println("EmployeeController, Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employeesLookup, &data)
		if err != nil {
			log.Println("EmployeeController, Automapper Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.EmployeeService.List(dataFilter)
		if err != nil {
			log.Println("EmployeeController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		err = structs.Automapper(employees, &data)
		if err != nil {
			log.Println("EmployeeController, Automapper List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	responsePayload.Setdata(data)
	if data == nil {
		data = []entity.EmployeeLookupResponse{}
	}
	if total == 0 {
		responsePayload.Setmsg("No Data")
		lastPage = 0
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

func (controller *EmployeeController) Export(c *fiber.Ctx) error {
	var dataFilter entity.EmployeeQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	dataFilter.Format = c.Query("format")

	buffer, contentType, filename, err := controller.EmployeeService.Export(dataFilter)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *EmployeeController) ExportTemplate(c *fiber.Ctx) error {
	format := c.Query("format", "xlsx")
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	buffer, contentType, filename, err := controller.EmployeeService.ExportTemplate(format)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buffer.Bytes())
}

func (controller *EmployeeController) ExportTemplateUpdate(c *fiber.Ctx) error {
	format := c.Query("format", "xlsx")
	fieldsParam := c.Query("fields")
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if strings.TrimSpace(fieldsParam) == "" {
		responsePayload.Setmsg("fields query parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	parts := strings.Split(fieldsParam, ",")
	fields := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			fields = append(fields, trimmed)
		}
	}
	if len(fields) == 0 {
		responsePayload.Setmsg("fields query parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	buffer, contentType, filename, err := controller.EmployeeService.ExportTemplateUpdate(custId, format, fields)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buffer.Bytes())
}

func (controller *EmployeeController) Import(c *fiber.Ctx) error {
	format := c.Query("format")
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if strings.TrimSpace(format) == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		responsePayload.Setmsg("file_upload is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := fileHeader.Open()
	if err != nil {
		responsePayload.Setmsg("failed to open uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	req := entity.ImportRequest{
		File:         file,
		CustId:       c.Locals("cust_id").(string),
		ParentCustId: c.Locals("parent_cust_id").(string),
		UserId:       c.Locals("user_id").(int64),
		Filename:     fileHeader.Filename,
		Format:       format,
	}

	if err := controller.EmployeeService.ImportEmployees(req); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *EmployeeController) ImportUpdate(c *fiber.Ctx) error {
	format := c.Query("format")
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if strings.TrimSpace(format) == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		responsePayload.Setmsg("file_upload is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := fileHeader.Open()
	if err != nil {
		responsePayload.Setmsg("failed to open uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	req := entity.ImportRequest{
		File:         file,
		CustId:       c.Locals("cust_id").(string),
		ParentCustId: c.Locals("parent_cust_id").(string),
		UserId:       c.Locals("user_id").(int64),
		Filename:     fileHeader.Filename,
		Format:       format,
	}

	if err := controller.EmployeeService.ImportEmployeesUpdate(req); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
func (controller *EmployeeController) Detail(c *fiber.Ctx) error {
	var params entity.DetailEmployeeParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("EmployeeController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("EmployeeController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustId = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("EmployeeController, Detail, CustId:", custId)

	data, err := controller.EmployeeService.Detail(params)
	if err != nil {
		log.Println("EmployeeController, Detail, FindOneByEmployeeId, err:", err.Error())
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

func (controller *EmployeeController) Create(c *fiber.Ctx) error {
	var request entity.CreateEmployeeBody
	if err := c.BodyParser(&request); err != nil {
		log.Println("EmployeeController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.CreatedBy = c.Locals("user_id").(int64)
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("EmployeeController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.EmployeeService.Store(request)
	if err != nil {
		log.Println("EmployeeController, Create, Store, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   constant.SUCCESSFULLY_ADDED,
	})
}

func (controller *EmployeeController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateEmployeeParams
		request entity.UpdateEmployeeRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("EmployeeController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_id").(int64)

	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("EmployeeController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("EmployeeController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_id").(int64)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		// log.Println("EmployeeController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	err := controller.EmployeeService.Update(params.EmployeeId, request)
	if err != nil {
		log.Println("EmployeeController, Update, Service.Update, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   constant.SUCCESSFULLY_UPDATED,
	})
}

func (controller *EmployeeController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteEmployeeParams
	if err := c.ParamsParser(&params); err != nil {
		log.Println("EmployeeController, Delete, ParamsParser, err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("EmployeeController, Delete, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmployeeController, Delete, CustId:", custId)

	err := controller.EmployeeService.Delete(custId, params.EmployeeId, userId)
	if err != nil {
		log.Println("EmployeeController, Delete, Service.Delete, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   "Deleted Successfully",
	})
}

func (controller *EmployeeController) CreateMultiple(c *fiber.Ctx) error {
	var request entity.CreateMultipleEmployeeBody
	if err := c.BodyParser(&request); err != nil {
		log.Println("EmployeeController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmployeeController, Create, CustId:", custId)

	// request.CustId = custId
	// request.CreatedBy = userId
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("EmployeeController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.EmployeeService.StoreMultiple(request, custId, userId)
	if err != nil {
		log.Println("EmployeeController, Create, Store, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   "Created Successfully",
	})
}
