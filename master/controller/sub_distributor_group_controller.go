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

type SubDistributorGroupController struct {
	SubDistributorGroupService service.SubDistributorGroupService
	validator                  *validation.Validate
}

func NewSubDistributorGroupController(mSubDistributorGroupService service.SubDistributorGroupService, validator *validation.Validate) *SubDistributorGroupController {
	return &SubDistributorGroupController{
		SubDistributorGroupService: mSubDistributorGroupService,
		validator:                  validator,
	}
}

func (controller *SubDistributorGroupController) Route(app *fiber.App) {
	qParamId := ":sub_distributor_group_id"
	subDistributorGroupRouteV1 := app.Group("/v1/sub-distributor-groups", middleware.JWTProtected())
	subDistributorGroupRouteV1.Post("", controller.Create)
	subDistributorGroupRouteV1.Get("", controller.List)
	subDistributorGroupRouteV1.Get("/"+qParamId, controller.Detail)
	subDistributorGroupRouteV1.Patch("/"+qParamId, controller.Update)
	subDistributorGroupRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SubDistributorGroupController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateSubDistributorGroupBody
	if err := c.BodyParser(&request); err != nil {
		log.Println("SubDistributorGroupController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BrandController, Create, CustId:", custId)

	request.CustID = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubDistributorGroupController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SubDistributorGroupService.Store(request)
	if err != nil {
		log.Println("SubDistributorGroupController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SubDistributorGroupController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.SubDistributorGroupQueryFilter
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
		log.Println("SubDistributorGroupController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("BrandController, List, CustId:", custId)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.SubDistributorGroupService.LookupList(dataFilter, custId)
		if err != nil {
			log.Println("SubDistributorGroupController, LookupList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.SubDistributorGroupService.List(dataFilter, custId)
		if err != nil {
			log.Println("SubDistributorGroupController, List, data, err:", err.Error())
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

func (controller *SubDistributorGroupController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSubDistributorGroupParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubDistributorGroupController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubDistributorGroupController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("BrandController, Detail, CustId:", custId)

	data, err := controller.SubDistributorGroupService.Detail(params.SubDistributorGroupId, custId)
	if err != nil {
		log.Println("SubDistributorGroupController, Detail, FindOneByBrandId, err:", err.Error())
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

func (controller *SubDistributorGroupController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSubDistributorGroupParams
		request entity.SubDistributorGroupUpdateRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubDistributorGroupController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubDistributorGroupController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("SubDistributorGroupController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("divisionController, Update, CustId:", custId)

	request.CustID = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("SubDistributorGroupController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	err := controller.SubDistributorGroupService.Update(params.SubDistributorGroupId, request)
	if err != nil {
		log.Println("SubDistributorGroupController, Update, Service.Update, err:", err.Error())
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

func (controller *SubDistributorGroupController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteSubDistributorGroupParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("SubDistributorGroupController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("SubDistributorGroupController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BrandController, Delete, CustId:", custId)

	err := controller.SubDistributorGroupService.Delete(custId, params.SubDistributorGroupId, userId)
	if err != nil {
		log.Println("BrandController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
