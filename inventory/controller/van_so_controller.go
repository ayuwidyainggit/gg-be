package controller

import (
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type VanSoController struct {
	vanSoService service.VanSoService
	validator    *validation.Validate
}

func NewVanSoController(vanSoService service.VanSoService, validator *validation.Validate) *VanSoController {
	return &VanSoController{
		vanSoService: vanSoService,
		validator:    validator,
	}
}
func (controller *VanSoController) Route(app *fiber.App) {
	qParamId := ":van_so_no"
	grRouteV1 := app.Group("/v1/van-so", middleware.JWTProtected())
	grRouteV1.Post("", controller.Create)
	grRouteV1.Get("/"+qParamId, controller.Detail)
	grRouteV1.Get("", controller.List)
	grRouteV1.Delete("/"+qParamId, controller.Delete)
	grRouteV1.Patch("/"+qParamId, controller.Update)
}

func (controller *VanSoController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateVanSoBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("VanSoController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("VanSoController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.vanSoService.Store(request)
	if err != nil {
		log.Error("VanSoController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
func (controller *VanSoController) Detail(c *fiber.Ctx) error {
	var params entity.DetailVanSoParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("VanSoController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("VanSoController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.vanSoService.Detail(params.VanSoNo, custId)
	if err != nil {
		log.Error("VanSoController, Detail, err:", err.Error())
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
func (controller *VanSoController) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("VanSoController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("BpprController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.vanSoService.List(dataFilter)
	if err != nil {
		log.Error("VanSoController, List, data, err:", err.Error())
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

func (controller *VanSoController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailVanSoParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("VanSoController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("VanSoController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.vanSoService.Delete(custId, params.VanSoNo, userId)
	if err != nil {
		log.Error("VanSoController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
func (controller *VanSoController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateVanSoParams
		request entity.UpdateVanSoBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("VanSoController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("VanSoController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("VanSoController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustID = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("VanSoController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.vanSoService.Update(params.VanSoNo, request)
	if err != nil {
		log.Error("VanSoController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
