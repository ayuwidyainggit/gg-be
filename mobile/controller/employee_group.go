package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type EmpGroupController struct {
	EmpGroupService service.EmpGroupService
	validator       *validation.Validate
}

func NewEmpGroupController(empGroupService service.EmpGroupService, validator *validation.Validate) EmpGroupController {
	return EmpGroupController{
		EmpGroupService: empGroupService,
		validator:       validator,
	}
}

func (controller *EmpGroupController) Route(app *fiber.App) {
	qParamId := ":emp_grp_id"
	empGroupsRouteV1 := app.Group("/v1/emp-groups", middleware.JWTProtected())
	empGroupsRouteV1.Get("/"+qParamId, controller.Detail)
	empGroupsRouteV1.Get("", controller.List)
	empGroupsRouteV1.Post("", controller.Create)
	empGroupsRouteV1.Patch("/"+qParamId, controller.Update)
	empGroupsRouteV1.Delete("/"+qParamId, controller.Delete)

	qParamTypeId := ":emp_type_id"
	empTypessRouteV1 := app.Group("/v1/emp-types", middleware.JWTProtected())
	empTypessRouteV1.Get("", controller.EmpTypesList)
	empTypessRouteV1.Get("/"+qParamTypeId, controller.TypeDetail)
}

func (controller *EmpGroupController) EmpTypesList(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("EmpGroupController, EmpTypesList, query parser filter:", err.Error())
	}

	custId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.EmpGroupService.EmpTypeList(dataFilter, custId)
	if err != nil {
		log.Error("EmpGroupController, EmpTypesList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
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

func (controller *EmpGroupController) TypeDetail(c *fiber.Ctx) error {
	var params entity.DetailEmpTypeParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmpGroupController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	// log.Println("EmpGroupController, Detail, CustId:", custId)

	data, err := controller.EmpGroupService.TypeDetail(params.EmpTypeId, custId)
	if err != nil {
		log.Error("EmpGroupController, Detail, FindOneByEmpGroupId, err:", err.Error())
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

func (controller *EmpGroupController) Detail(c *fiber.Ctx) error {
	var params entity.DetailEmpGroupParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmpGroupController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	// log.Println("EmpGroupController, Detail, CustId:", custId)

	data, err := controller.EmpGroupService.Detail(params.EmpGroupId, custId)
	if err != nil {
		log.Error("EmpGroupController, Detail, FindOneByEmpGroupId, err:", err.Error())
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

func (controller *EmpGroupController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.GeneralQueryFilter
		data       interface{}
		total      int64
		lastPage   int
	)

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("EmpGroupController, List, query parser filter:", err.Error())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("EmpGroupController, List, dataFilter:", structs.StructToJson(dataFilter))
	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.EmpGroupService.LookupList(dataFilter)
		if err != nil {
			log.Error("EmpGroupController, LookupList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.EmpGroupService.List(dataFilter)
		if err != nil {
			log.Error("EmpGroupController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
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

func (controller *EmpGroupController) Create(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateEmpGroupBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("EmpGroupController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmpGroupController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.EmpGroupService.Store(request)
	if err != nil {
		log.Error("EmpGroupController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *EmpGroupController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateEmpGroupParams
		request entity.UpdateEmpGroupRequest
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmpGroupController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("EmpGroupController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmpGroupController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.EmpGroupService.Update(params.EmpGroupId, request)
	if err != nil {
		log.Error("EmpGroupController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *EmpGroupController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteEmpGroupParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("EmpGroupController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("EmpGroupController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("EmpGroupController, Delete, CustId:", custId)

	err := controller.EmpGroupService.Delete(custId, params.EmpGroupId, userId)
	if err != nil {
		log.Error("EmpGroupController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
