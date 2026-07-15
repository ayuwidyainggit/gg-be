package controller

import (
	"fmt"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SalesTeamController struct {
	SalesTeamService service.SalesTeamService
	validator        *validation.Validate
}

func NewSalesTeamController(salesTeamService service.SalesTeamService, validator *validation.Validate) *SalesTeamController {
	return &SalesTeamController{
		SalesTeamService: salesTeamService,
		validator:        validator,
	}
}

func (controller *SalesTeamController) Route(app *fiber.App) {
	qParamId := ":sales_team_id"
	outletTypesRouteV1 := app.Group("/v1/sales-teams", middleware.JWTProtected())
	outletTypesRouteV1.Get("/list-by-distributor", controller.ListByDistributor)
	outletTypesRouteV1.Get("/"+qParamId, controller.Detail)
	outletTypesRouteV1.Get("", controller.List)
	outletTypesRouteV1.Post("", controller.Create)
	outletTypesRouteV1.Patch("/"+qParamId, controller.Update)
	outletTypesRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SalesTeamController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSalesTeamParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SalesTeamController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SalesTeamController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Info("SalesTeamController, Detail, CustId:", custId)

	data, err := controller.SalesTeamService.Detail(params.SalesTeamId, custId)
	if err != nil {
		log.Info("SalesTeamController, Detail, FindOneBySalesTeamId, err:", err.Error())
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

func (controller *SalesTeamController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.GeneralQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("SalesTeamController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	distributorIDs, err := parseIntSliceQueryAllowZero(c.Context().QueryArgs(), "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorIDs = distributorIDs
		if len(distributorIDs) == 1 && distributorIDs[0] > 0 {
			dataFilter.DistributorID = int64(distributorIDs[0])
		}
	}

	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	parentCustId := ""
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		customer, err := controller.SalesTeamService.FindParentCustId(custId)
		if err != nil {
			log.Info("SalesTeamController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = customer.CustId
		parentCustId = customer.ParentCustId
	}
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId
	log.Info("custId:", custId)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.SalesTeamService.LookupList(dataFilter, custId)
		if err != nil {
			log.Info("SalesTeamController, LookupList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.SalesTeamService.List(dataFilter, custId)
		if err != nil {
			log.Info("SalesTeamController, List, data, err:", err.Error())
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

func (controller *SalesTeamController) Create(c *fiber.Ctx) error {
	var request entity.CreateSalesTeamBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("SalesTeamController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("SalesTeamController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("SalesTeamController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SalesTeamService.Store(request)
	if err != nil {
		log.Info("SalesTeamController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesTeamController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSalesTeamParams
		request entity.UpdateSalesTeamRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SalesTeamController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SalesTeamController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Info("SalesTeamController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("SalesTeamController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("SalesTeamController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesTeamService.Update(params.SalesTeamId, request)
	if err != nil {
		log.Info("SalesTeamController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SalesTeamController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSalesTeamParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SalesTeamController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SalesTeamController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	err := controller.SalesTeamService.Delete(custId, params.SalesTeamId, userId)
	if err != nil {
		log.Info("SalesTeamController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *SalesTeamController) ListByDistributor(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.GeneralQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("SalesTeamController, ListByDistributor, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := fmt.Sprint(c.Locals("cust_id"))
	parentCustID := fmt.Sprint(c.Locals("parent_cust_id"))

	dataFilter.CustId = custID
	dataFilter.ParentCustId = parentCustID

	if dataFilter.CustId != dataFilter.ParentCustId {
		dataFilter.DistributorID = c.Locals("distributor_id").(int64)
	}

	data, total, lastPage, err = controller.SalesTeamService.ListByDistributor(dataFilter)
	if err != nil {
		log.Info("SalesTeamController, ListByDistributor, data, err:", err.Error())
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
