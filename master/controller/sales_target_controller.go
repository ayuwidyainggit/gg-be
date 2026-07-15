package controller

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type SalesTargetController struct {
	SalesTargetService service.SalesTargetService
	validator          *validation.Validate
}

func NewSalesTargetController(salesTargetService service.SalesTargetService, validator *validation.Validate) *SalesTargetController {
	return &SalesTargetController{
		SalesTargetService: salesTargetService,
		validator:          validator,
	}
}

func (controller *SalesTargetController) Route(app *fiber.App) {
	salesTargetRouteV1 := app.Group("/v1/sales-target", middleware.JWTProtected())
	salesTargetRouteV1.Get("", controller.List)
	salesTargetRouteV1.Get("/:sales_target_id", controller.Detail)
	salesTargetRouteV1.Post("", controller.Create)
	salesTargetRouteV1.Patch("/:sales_target_id", controller.Update)
}

func (controller *SalesTargetController) getAcceptLang(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	if val, ok := headers[constant.HEADER_ACCEPT_LANG]; ok {
		if len(val) > 0 && val[0] != "" {
			return val[0]
		}
	}
	return ""
}

// List - GET /master/v1/sales-target
func (controller *SalesTargetController) List(c *fiber.Ctx) error {
	var dataFilter entity.SalesTargetQueryFilter
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 20
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	data, total, lastPage, err := controller.SalesTargetService.List(dataFilter, custId)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Handle empty data
	if data == nil {
		data = []entity.SalesTargetListResponse{}
	}

	if total == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
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

// Detail - GET /master/v1/sales-target/:sales_target_id
func (controller *SalesTargetController) Detail(c *fiber.Ctx) error {
	var params entity.SalesTargetParams
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	data, err := controller.SalesTargetService.Detail(params.SalesTargetId, custId)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, constant.ErrSalesTargetNotFound) {
			statusCode = fiber.StatusNotFound
			errMsg = "record not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: len(data.Details),
		PageCurrent: 1,
		PageLimit:   len(data.Details),
		PageTotal:   1,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// Create - POST /master/v1/sales-target
func (controller *SalesTargetController) Create(c *fiber.Ctx) error {
	var request entity.CreateSalesTargetRequest
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.CreatedBy = c.Locals("user_id").(int64)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesTargetService.Store(request)
	if err != nil {
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

// Update - PATCH /master/v1/sales-target/:sales_target_id
func (controller *SalesTargetController) Update(c *fiber.Ctx) error {
	var params entity.SalesTargetParams
	var request entity.UpdateSalesTargetRequest
	headerAcceptLang := controller.getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_id").(int64)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesTargetService.Update(params.SalesTargetId, request)
	if err != nil {
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
