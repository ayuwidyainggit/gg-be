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

type ReturnController struct {
	ReturnService service.ReturnService
	validator     *validation.Validate
}

func NewReturnController(Service service.ReturnService, validator *validation.Validate) *ReturnController {
	return &ReturnController{
		ReturnService: Service,
		validator:     validator,
	}
}
func (controller *ReturnController) Route(app *fiber.App) {
	returnRouteV1 := app.Group("/v1/returns", middleware.JWTProtected())
	returnRouteV1.Post("", controller.Create)
	returnRouteV1.Get("/master/return-reasons", controller.ReturnReasonList)
	returnRouteV1.Post("/status", controller.UpdateStatus)

	qParamId := ":return_no"
	returnRouteV1.Patch("/quantity/"+qParamId, controller.UpdateQuantity)

	returnRouteV2 := app.Group("/v1/return-reasons", middleware.JWTProtected())

	// returnRouteV1.Get("/", controller.ReturnList)
	returnRouteV2.Get("/", controller.Returns)
}

func (controller *ReturnController) Create(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateReturnBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("ReturnController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("ReturnController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ReturnService.Store(request)
	if err != nil {
		log.Error("ReturnController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ReturnController) UpdateStatus(c *fiber.Ctx) error {
	var (
		request entity.UpdateStatusReturnBody
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ReturnController, UpdateStatus, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	request.CustID = custId
	request.UpdatedBy = userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("ReturnController, UpdateStatus, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ReturnService.UpdateStatus(request)
	if err != nil {
		log.Error("ReturnController, UpdateStatus, Service.UpdateStatus, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReturnController) Returns(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		request          entity.ReturnReasonsRequest
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {

		// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ReturnService.ReturnReasons(request)
	if err != nil {
		log.Error("ReturnsSummaryDaily, Detail, FindOneReturnsSummary, err:", err.Error())
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
	// log.Println("response:", response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ReturnController) ReturnReasonList(c *fiber.Ctx) error {
	var (
		dataFilter entity.GeneralQueryFilter
		data       interface{}
		err        error
		total      int64
		lastPage   int
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("returnController, ReturnReasonListCreate, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("returnController, ReturnReasonListCreate, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err = controller.ReturnService.ReturnReasonLookupList(dataFilter)
	if err != nil {
		log.Error("returnController, ReturnReasonListCreate, data, err:", err.Error())
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

func (controller *ReturnController) UpdateQuantity(c *fiber.Ctx) error {
	var (
		params  entity.UpdateReturnParams
		request entity.UpdateQuantityReturnBody
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ReturnController, UpdateQuantity, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("ReturnController, UpdateQuantity, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("ReturnController, UpdateQuantity, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	request.CustID = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("ReturnController, UpdateQuantity, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ReturnService.UpdateQuantity(params.ReturnNo, request)
	if err != nil {
		log.Error("ReturnController, UpdateQuantity, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
