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

type OutletConfigController struct {
	OutletConfigService       service.OutletConfigService
	OutletConfigStatusService service.OutletConfigStatusService
	validator                 *validation.Validate
}

func NewOutletConfigController(svc service.OutletConfigService, statusSvc service.OutletConfigStatusService, validator *validation.Validate) *OutletConfigController {
	return &OutletConfigController{
		OutletConfigService:       svc,
		OutletConfigStatusService: statusSvc,
		validator:                 validator,
	}
}

func (c *OutletConfigController) Route(app *fiber.App) {
	qParamId := ":outlet_config_id"
	g := app.Group("/v1/outlet_config", middleware.JWTProtected())
	g.Get("", c.List)
	g.Post("", c.Create)
	g.Get("/"+qParamId, c.Detail)
	g.Put("/"+qParamId, c.Update)
	g.Delete("/"+qParamId, c.Delete)

	gStatus := app.Group("/v1/outlet_config_status", middleware.JWTProtected())
	gStatus.Get("", c.StatusList)
}

func (c *OutletConfigController) List(ctx *fiber.Ctx) error {
	var filter entity.OutletConfigListFilter
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)

	if err := ctx.QueryParser(&filter); err != nil {
		log.Println("OutletConfigController List QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 5
	}
	if filter.Sort == "" {
		filter.Sort = "created_at:desc"
	}

	custID := ctx.Locals("cust_id").(string)
	parentCustID := ""
	if v := ctx.Locals("parent_cust_id"); v != nil {
		parentCustID = v.(string)
	}
	data, total, lastPage, err := c.OutletConfigService.List(filter, custID, parentCustID)
	if err != nil {
		log.Println("OutletConfigController List:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("Success")
		responsePayload.Setdata(data)
	}
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})
	return ctx.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (c *OutletConfigController) Detail(ctx *fiber.Ctx) error {
	var params entity.OutletConfigDetailParams
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)

	if err := ctx.ParamsParser(&params); err != nil {
		log.Println("OutletConfigController Detail ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := c.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("OutletConfigController Detail ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := ctx.Locals("cust_id").(string)
	parentCustID := ""
	if v := ctx.Locals("parent_cust_id"); v != nil {
		parentCustID = v.(string)
	}
	data, err := c.OutletConfigService.Detail(params.OutletConfigId, custID, parentCustID)
	if err != nil {
		log.Println("OutletConfigController Detail:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return ctx.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	return ctx.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (c *OutletConfigController) Create(ctx *fiber.Ctx) error {
	var body entity.CreateOutletConfigBody
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)

	if err := ctx.BodyParser(&body); err != nil {
		log.Println("OutletConfigController Create BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := c.validator.ValidateStruct(body, headerAcceptLang)
	if errs != nil {
		log.Println("OutletConfigController Create ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := ctx.Locals("cust_id").(string)
	userID := ctx.Locals("user_id").(int64)
	if err := c.OutletConfigService.Create(body, custID, userID); err != nil {
		log.Println("OutletConfigController Create:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Outlet Status Configuration successfully added")
	return ctx.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (c *OutletConfigController) Update(ctx *fiber.Ctx) error {
	var params entity.OutletConfigDetailParams
	var body entity.CreateOutletConfigBody
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)

	if err := ctx.ParamsParser(&params); err != nil {
		log.Println("OutletConfigController Update ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := ctx.BodyParser(&body); err != nil {
		log.Println("OutletConfigController Update BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := c.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		log.Println("OutletConfigController Update ValidateStruct params:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := c.validator.ValidateStruct(body, headerAcceptLang); errs != nil {
		log.Println("OutletConfigController Update ValidateStruct body:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := ctx.Locals("cust_id").(string)
	parentCustID := ""
	if v := ctx.Locals("parent_cust_id"); v != nil {
		parentCustID = v.(string)
	}
	userID := ctx.Locals("user_id").(int64)
	if err := c.OutletConfigService.Update(params.OutletConfigId, custID, parentCustID, body, userID); err != nil {
		log.Println("OutletConfigController Update:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if errMsg == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		if errMsg == "not allowed to edit this outlet config" {
			statusCode = fiber.StatusForbidden
		}
		responsePayload.Setmsg(errMsg)
		return ctx.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Outlet Status Configuration successfully updated")
	return ctx.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (c *OutletConfigController) Delete(ctx *fiber.Ctx) error {
	var params entity.OutletConfigDetailParams
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)
	if err := ctx.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := c.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custID := ctx.Locals("cust_id").(string)
	parentCustID := ""
	if v := ctx.Locals("parent_cust_id"); v != nil {
		parentCustID = v.(string)
	}
	userID := ctx.Locals("user_id").(int64)
	if err := c.OutletConfigService.Delete(params.OutletConfigId, custID, parentCustID, userID); err != nil {
		log.Println("OutletConfigController Delete:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if errMsg == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		if errMsg == "not allowed to delete this outlet config" {
			statusCode = fiber.StatusForbidden
		}
		responsePayload.Setmsg(errMsg)
		return ctx.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Outlet Status Configuration successfully deleted")
	return ctx.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (c *OutletConfigController) StatusList(ctx *fiber.Ctx) error {
	var filter entity.OutletConfigStatusListFilter
	headerAcceptLang := ""
	if len(ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = ctx.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(ctx.Locals("requestid").(string), headerAcceptLang)

	if err := ctx.QueryParser(&filter); err != nil {
		log.Println("OutletConfigController StatusList QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Sort == "" {
		filter.Sort = "sort_order:asc"
	}

	data, total, lastPage, err := c.OutletConfigStatusService.List(filter)
	if err != nil {
		log.Println("OutletConfigController StatusList:", err.Error())
		responsePayload.Setmsg(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("Success")
		responsePayload.Setdata(data)
	}
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})
	return ctx.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
