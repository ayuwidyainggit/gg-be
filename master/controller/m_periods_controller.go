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

type MPeriodsController struct {
	MPeriodsService service.MPeriodsService
	validator       *validation.Validate
}

func NewMPeriodsController(mPeriodsService service.MPeriodsService, validator *validation.Validate) *MPeriodsController {
	return &MPeriodsController{
		MPeriodsService: mPeriodsService,
		validator:       validator,
	}
}

func (controller *MPeriodsController) Route(app *fiber.App) {
	qParamId1 := ":per_year"
	qParamId2 := ":per_id"
	mPeriods := app.Group("/v1/periods", middleware.JWTProtected())
	mPeriods.Get("/"+qParamId1+"/"+qParamId2, controller.Detail)
	mPeriods.Get("", controller.List)
	mPeriods.Get("/years", controller.ListYear)
	mPeriods.Post("", controller.Create)
	mPeriods.Patch("/"+qParamId1+"/"+qParamId2, controller.Update)
	mPeriods.Delete("/"+qParamId1+"/"+qParamId2, controller.Delete)
}

func (controller *MPeriodsController) Detail(c *fiber.Ctx) error {
	var params entity.DetailCreateMPeriodsParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("MPeriodsController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("MPeriodsController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	custId := c.Locals("parent_cust_id").(string)

	data, err := controller.MPeriodsService.Detail(params.PerYear, params.PerId, custId)
	if err != nil {
		log.Println("MPeriodsController, Detail, err:", err.Error())
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

func (controller *MPeriodsController) List(c *fiber.Ctx) error {
	var dataFilter entity.MPeriodsQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("MPeriodsController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// custId := c.Locals("parent_cust_id").(string)
	data, total, lastPage, err := controller.MPeriodsService.List(dataFilter, custId)
	if err != nil {
		log.Println("MPeriodsController, List, data, err:", err.Error())
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

func (controller *MPeriodsController) Create(c *fiber.Ctx) error {
	log.Println("masuk controller Create")
	var request entity.CreateMPeriodsBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("MPeriodsController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	log.Println("MPeriodsController, Create, userId:", c.Locals("user_id"))
	userId := c.Locals("user_id").(int64)
	log.Println("MPeriodsController, Create, userId:", userId)
	// log.Println("MPeriodsController, Create, CustId:", custId)

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = int64(userId)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("MPeriodsController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	log.Println("pass controller validateStruct")

	_, err := controller.MPeriodsService.Store(request)
	if err != nil {
		log.Println("MPeriodsController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *MPeriodsController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateMPeriodsParams
		request entity.UpdateMPeriodsRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("MPeriodsController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("MPeriodsController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("MPeriodsController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// log.Println("MPeriodsController, Update, CustId:", custId)
	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("MPeriodsController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.ParentCustId = parentCustId
	err := controller.MPeriodsService.Update(params.PerYear, params.PerId, request)
	if err != nil {
		log.Println("MPeriodsController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPeriodsController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteMPeriodsParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("MPeriodsController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("MPeriodsController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	userName := c.Locals("user_name").(string)
	// log.Println("MPeriodsController, Delete, CustId:", custId)

	err := controller.MPeriodsService.Delete(custId, params.PerYear, params.PerId, userId, userName)
	if err != nil {
		log.Println("MPeriodsController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPeriodsController) ListYear(c *fiber.Ctx) error {
	var dataFilter entity.MPeriodsQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("MPeriodsController, List, query parser filter:", err.Error())
	}

	custId := c.Locals("parent_cust_id").(string)
	data, err := controller.MPeriodsService.ListYear(dataFilter, custId)
	if err != nil {
		log.Println("MPeriodsController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
