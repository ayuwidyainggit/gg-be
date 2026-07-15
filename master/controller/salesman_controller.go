package controller

import (
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SalesmanController struct {
	SalesmanService service.SalesmanService
	validator       *validation.Validate
}

func NewSalesmanController(salesmanService service.SalesmanService, validator *validation.Validate) *SalesmanController {
	return &SalesmanController{
		SalesmanService: salesmanService,
		validator:       validator,
	}
}

func (controller *SalesmanController) Route(app *fiber.App) {
	qParamId := ":emp_id"
	salesmanRouteV1Scheduler := app.Group("/v1/salesman/scheduler")
	salesmanRouteV1Scheduler.Post("/isactive", controller.UpdateIsActive)
	salesmanRouteV1Scheduler.Post("/deactive", controller.UpdateDeActive)
	salesmanRouteV1Scheduler.Post("/custom", controller.CustomScheduler)
	salesmansRouteV1 := app.Group("/v1/salesman", middleware.JWTProtected())
	salesmansRouteV1.Get("/job-type", controller.LookupJobType)
	salesmansRouteV1.Get("/tax-option", controller.LookupTaxOption)
	salesmansRouteV1.Get("/"+qParamId, controller.Detail)
	salesmansRouteV1.Get("", controller.List)
	salesmansRouteV1.Post("", controller.Create)
	salesmansRouteV1.Patch("/"+qParamId, controller.Update)
	salesmansRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SalesmanController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSalesmanParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SalesmanController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustId = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("SalesmanController, Detail, CustId:", custId)s

	data, err := controller.SalesmanService.Detail(params)
	if err != nil {
		log.Println("SalesmanController, Detail, FindOneByEmpId, err:", err.Error())
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

func parseSalesmanDistributorIDs(rawValue string) ([]int, error) {
	values, err := parseCSVIntValues(rawValue, "distributor_id")
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, len(values))
	for _, value := range values {
		if value < 0 {
			continue
		}
		result = append(result, value)
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func parseSalesmanDistributorIDQuery(args *fasthttp.Args) ([]int, error) {
	return parseIntSliceQueryAllowZero(args, "distributor_id", "distributor_id", "distributor_id[]")
}

func (controller *SalesmanController) List(c *fiber.Ctx) error {
	var (
		err            error
		dataFilter     entity.SalesmanQueryFilter
		data           interface{}
		total          int
		lastPage       int
		salesmanList   []entity.SalesmanListResponse
		salesmanLookup []entity.SalesmanLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("SalesmanController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	distributorIDs, err := parseSalesmanDistributorIDQuery(c.Context().QueryArgs())
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorID = distributorIDs
	}

	custIDs := parseStringSliceQuery(c.Context().QueryArgs(), "cust_id", "cust_id[]")
	if len(custIDs) > 0 {
		dataFilter.CustIds = custIDs
	}

	// custId := c.Locals("cust_id").(string)
	parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.SalesmanService.FindParentCustId(custId)
		if err != nil {
			log.Println("SalesmanController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		parentCustId = parentCust.ParentCustId
	}
	log.Println("custId:", custId)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.SalesmanService.List(dataFilter, custId, parentCustId)
		if err != nil {
			log.Println("SalesmanController, Lookup, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(salesmanLookup, &data)
		if err != nil {
			log.Println("SalesmanController, Lookup, Automapper data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.SalesmanService.List(dataFilter, custId, parentCustId)
		if err != nil {
			log.Println("SalesmanController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		err = structs.Automapper(salesmanList, &data)
		if err != nil {
			log.Println("SalesmanController, List, Automapper data, err:", err.Error())
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

func (controller *SalesmanController) Create(c *fiber.Ctx) error {
	var request entity.CreateSalesmanBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SalesmanController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.CreatedBy = c.Locals("user_id").(int64)
	request.ParentCustId = c.Locals("parent_cust_id").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SalesmanService.Store(request)
	if err != nil {
		log.Println("SalesmanController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSalesmanParams
		request entity.UpdateSalesmanRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SalesmanController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_id").(int64)

	// requestCanvas.CustId = c.Locals("cust_id").(string)
	// requestCanvas.ParentCustId = c.Locals("parent_cust_id").(string)
	// requestCanvas.UpdatedBy = c.Locals("user_id").(int64)

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SalesmanController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesmanService.Update(params.EmpId, request)
	if err != nil {
		log.Println("SalesmanController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSalesmanParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SalesmanController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SalesmanController, Delete, CustId:", custId)

	err := controller.SalesmanService.Delete(custId, params.EmpId, userId)
	if err != nil {
		log.Println("SalesmanController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) LookupJobType(c *fiber.Ctx) error {
	var (
		err        error
		data       interface{}
		dataLookup []entity.JobTypeLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// custId := c.Locals("cust_id").(string)
	// parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		// parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.SalesmanService.FindParentCustId(custId)
		if err != nil {
			log.Println("SalesmanController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		// parentCustId = parentCust.ParentCustId
	}
	log.Println("custId:", custId)

	data, err = controller.SalesmanService.LookupJobType()
	if err != nil {
		log.Println("SalesmanController, Lookup, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	err = structs.Automapper(dataLookup, &data)
	if err != nil {
		log.Println("SalesmanController, Lookup, Automapper data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) LookupTaxOption(c *fiber.Ctx) error {
	var (
		err        error
		data       interface{}
		dataLookup []entity.TaxOptionLookupResponse
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// custId := c.Locals("cust_id").(string)
	// parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		// parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.SalesmanService.FindParentCustId(custId)
		if err != nil {
			log.Println("SalesmanController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		// parentCustId = parentCust.ParentCustId
	}
	log.Println("custId:", custId)

	data, err = controller.SalesmanService.LookupTaxOption()
	if err != nil {
		log.Println("SalesmanController, Lookup, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	err = structs.Automapper(dataLookup, &data)
	if err != nil {
		log.Println("SalesmanController, Lookup, Automapper data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) UpdateIsActive(c *fiber.Ctx) error {
	var request entity.UpdateIsActiveRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SalesmanController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesmanService.UpdateIsActive(request.CustId, int64(request.EmpId), int64(request.UserId))
	if err != nil {
		log.Println("SalesmanController, Update Is Active, Service.UpdateIsActive, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Update Is Active Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) UpdateDeActive(c *fiber.Ctx) error {
	var request entity.UpdateIsActiveRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SalesmanController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesmanService.UpdateDeActive(request.CustId, int64(request.EmpId), int64(request.UserId))
	if err != nil {
		log.Println("SalesmanController, Update De Active, Service.UpdateIsActive, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Update De Active Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesmanController) CustomScheduler(c *fiber.Ctx) error {
	var request entity.CustomSchedulerRequest
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("SalesmanController, Scheduler custom, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SalesmanController, Scheduler custom, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	err := controller.SalesmanService.GoSchedulerCustom(request.StartDate, request.Url)
	if err != nil {
		log.Println("SalesmanController, Scheduler custom, Service.GoSchedulerCustom, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Scheduler custom Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
