package controller

import (
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type MDivisionController struct {
	MDivisionService service.MDivisionService
	validator        *validation.Validate
}

func NewMDivisionController(divisionService service.MDivisionService, validator *validation.Validate) MDivisionController {
	return MDivisionController{
		MDivisionService: divisionService,
		validator:        validator,
	}
}

func (controller *MDivisionController) Route(app *fiber.App) {
	qParamId := ":division_id"
	MDivisionRouteV1 := app.Group("/v1/divisions", middleware.JWTProtected())
	MDivisionRouteV1.Get("/"+qParamId, controller.Detail)
	MDivisionRouteV1.Get("", controller.List)
	MDivisionRouteV1.Post("", controller.Create)
	MDivisionRouteV1.Patch("/"+qParamId, controller.Update)
	MDivisionRouteV1.Delete("/"+qParamId, controller.Delete)

}

func (controller *MDivisionController) List(c *fiber.Ctx) error {
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
		log.Println("MDivisionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.MDivisionService.LookupList(dataFilter)
		if err != nil {
			log.Println("EmployeeController, Lookup List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.MDivisionService.List(dataFilter)
		if err != nil {
			log.Println("MDivisionController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	if total == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg(constant.SUCCESS_NO_DATA)
		responsePayload.Setdata(data)
	}
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MDivisionController) Detail(c *fiber.Ctx) error {
	var params entity.DetailMDivisionParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("MDivisionController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("MDivisionController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	// log.Println("MDivisionController, Detail, CustId:", custId)

	data, err := controller.MDivisionService.Detail(params.MDivisionId, custId)
	if err != nil {
		log.Println("MDivisionController, Detail, FindOneByMDivisionId, err:", err.Error())
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

func (controller *MDivisionController) Create(c *fiber.Ctx) error {
	var request entity.CreateDivisionBody
	if err := c.BodyParser(&request); err != nil {
		log.Println("MDivisionController, Create, BodyParser:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("MDivisionController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("MDivisionController, Create, ValidateStruct, errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	_, err := controller.MDivisionService.Store(request)
	if err != nil {
		log.Println("MDivisionController, Create, Store, err:", err.Error())
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

func (controller *MDivisionController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateMDivisionParams
		request entity.UpdateDivisionBody
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("divisionController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("divisionController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("divisionController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("divisionController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("divisionController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	err := controller.MDivisionService.Update(params.MDivisionId, request)
	if err != nil {
		log.Println("divisionController, Update, Service.Update, err:", err.Error())
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

func (controller *MDivisionController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteMDivisionParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("divisionController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("divisionController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("SpPriceController, Delete, CustId:", custId)

	err := controller.MDivisionService.Delete(custId, params.MDivisionId, userId)
	if err != nil {
		log.Println("divisionController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
