package controller

import (
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type DistPriceController struct {
	DistPriceService service.DistPriceService
	validator        *validation.Validate
}

func NewDistPriceController(distPriceService service.DistPriceService, validator *validation.Validate) DistPriceController {
	return DistPriceController{
		DistPriceService: distPriceService,
		validator:        validator,
	}
}

func (controller *DistPriceController) Route(app *fiber.App) {
	qParamId := ":dist_price_id"
	distPriceRouteV1Scheduler := app.Group("/v1/distributor-prices/scheduler")
	distPriceRouteV1Scheduler.Post("/publish-unpublish", controller.PublishOrUnpublish)
	distPriceRouteV1 := app.Group("/v1/distributor-prices", middleware.JWTProtected())
	distPriceRouteV1.Get("/"+qParamId, controller.Detail)
	distPriceRouteV1.Get("", controller.List)
	distPriceRouteV1.Post("", controller.Create)
	distPriceRouteV1.Patch("/"+qParamId, controller.Update)
	distPriceRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *DistPriceController) Detail(c *fiber.Ctx) error {
	var params entity.DetailDistPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DistPriceController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	data, err := controller.DistPriceService.Detail(params.DistPriceId, custId)
	if err != nil {
		log.Error("DistPriceController, Detail, FindOneByDiscId, err:", err.Error())
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

func (controller *DistPriceController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.DistPriceQueryFilter
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
		log.Info("DistPriceController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	parentCustId := c.Locals("parent_cust_id").(string)
	custId := c.Locals("cust_id").(string)
	distPriceGrpId := c.Locals("dist_price_grp_id").(int)
	if custId != parentCustId {
		dataFilter.DistPriceGroupId = &distPriceGrpId
	}

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.DistPriceService.LookupList(dataFilter, parentCustId, custId)
		if err != nil {
			log.Info("DistPriceController, LookupList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
	case "lookup_product":
		data, total, lastPage, err = controller.DistPriceService.LookupProductList(dataFilter, parentCustId)
		if err != nil {
			log.Info("DistPriceController, LookupProductList, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.DistPriceService.List(dataFilter, parentCustId)
		if err != nil {
			log.Info("DistPriceController, List, data, err:", err.Error())
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

func (controller *DistPriceController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateDistPriceBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("DistPriceController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("DistPriceController, Create, CustId:", claims.CustId)

	request.CustId = custID
	request.CreatedBy = userId

	if custID != parentCustID {
		log.Info(c.Locals("requestid").(string) + ", this endpoint only for principal")
		responsePayload.Setmsg("this endpoint only for principal")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.DistPriceService.Store(request)
	if err != nil {
		log.Error("DistPriceController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *DistPriceController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateDistPriceParams
		request entity.UpdateDistPriceRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DistPriceController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("DistPriceController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	jwtClaims := entity.TokenMetadata{}
	err := structs.Automapper(c.Locals("claims"), &jwtClaims)
	if err != nil {
		log.Error("DistPriceController, Delete, Automapper, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.CustId = jwtClaims.CustId
	request.UpdatedBy = jwtClaims.UserId

	err = controller.DistPriceService.Update(jwtClaims, params.DistPriceId, request)
	if err != nil {
		log.Error("DistPriceController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DistPriceController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteDistPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("DistPriceController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	jwtClaims := entity.TokenMetadata{}
	err := structs.Automapper(c.Locals("claims"), &jwtClaims)
	if err != nil {
		log.Error("DistPriceController, Delete, Automapper, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.DistPriceService.Delete(jwtClaims, params.DistPriceId, jwtClaims.UserId)
	if err != nil {
		log.Error("DistPriceController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DistPriceController) PublishOrUnpublish(c *fiber.Ctx) error {
	var (
		request entity.PublishUnpublishDistPriceReq
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("DistPriceController, PublishOrUnpublish, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("DistPriceController, PublishOrUnpublish, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := controller.DistPriceService.PublishOrInactive(request); err != nil {
		log.Errorf("DistPriceController, Failed to PublishOrUnpublish in service: %v", err)
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
