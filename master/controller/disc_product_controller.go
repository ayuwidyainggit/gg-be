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

type DiscProductController struct {
	DiscProductService service.DiscProductService
	validator          *validation.Validate
}

func NewDiscProductController(discProductService service.DiscProductService, validator *validation.Validate) *DiscProductController {
	return &DiscProductController{
		DiscProductService: discProductService,
		validator:          validator,
	}
}

func (controller *DiscProductController) Route(app *fiber.App) {

	qParamId1 := ":disc_id"
	qParamId2 := ":pro_id"
	discProductRouteV1 := app.Group("/v1/disc-product", middleware.JWTProtected())
	discProductRouteV1.Get("/"+qParamId1+"/"+qParamId2, controller.Detail)
	discProductRouteV1.Get("", controller.List)
	discProductRouteV1.Post("", controller.Create)
	discProductRouteV1.Patch("/"+qParamId1+"/"+qParamId2, controller.Update)
	discProductRouteV1.Delete("/"+qParamId1+"/"+qParamId2, controller.Delete)
}

func (controller *DiscProductController) Detail(c *fiber.Ctx) error {
	var params entity.DetailDiscProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("DiscProductController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("DiscProductController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("DiscProductController, Detail, CustId:", custId)

	data, err := controller.DiscProductService.Detail(params.DiscId, params.ProId, custId)
	if err != nil {
		log.Println("DiscProductController, Detail, err:", err.Error())
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

func (controller *DiscProductController) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("DiscProductController, List, query parser filter:", err.Error())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("DiscProductController, List, CustId:", custId)

	data, total, lastPage, err := controller.DiscProductService.List(dataFilter, custId)
	if err != nil {
		log.Println("DiscProductController, List, data, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	// dataPrint, _ := json.Marshal(data)
	// log.Println("### DiscProductController, List, dataPrint ###")
	// log.Println(string(dataPrint))
	// log.Println("### End Of dataPrint ###")

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscProductController) Create(c *fiber.Ctx) error {
	var request entity.CreateDiscProductBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("DiscProductController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// log.Println("DiscProductController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId
	request.UpdatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("DiscProductController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.DiscProductService.Store(request)
	if err != nil {
		log.Println("DiscProductController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *DiscProductController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateDiscProductParams
		request entity.UpdateDiscProductRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("DiscProductController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("DiscProductController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("DiscProductController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// log.Println("DiscProductController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("DiscProductController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	// fmt.Println("REQ UPDATE >>>>>", params.ProId)
	err := controller.DiscProductService.Update(params.DiscId, params.ProId, request)
	if err != nil {
		log.Println("DiscProductController, Update, Service.Update, err:", err.Error())
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

func (controller *DiscProductController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteDiscProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("DiscProductController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("DiscProductController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	log.Println("DiscProductController, Delete, CustId:", custId)

	err := controller.DiscProductService.Delete(custId, params.DiscId, params.ProId, userId)
	if err != nil {
		log.Println("DiscProductController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
