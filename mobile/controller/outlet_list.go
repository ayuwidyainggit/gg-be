package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type OutletListController struct {
	Service   service.OutletListService
	Validator *validation.Validate
}

func NewOutletListController(svc service.OutletListService, validator *validation.Validate) *OutletListController {
	return &OutletListController{
		Service:   svc,
		Validator: validator,
	}
}

func (ctrl *OutletListController) Route(app *fiber.App) {
	// Routes for /v1/outlet-list
	outletListRoute := app.Group("/v1/outlet-list", middleware.JWTProtected())
	outletListRoute.Get("", ctrl.List)
	outletListRoute.Delete("/:outlet_id", ctrl.Delete)

	// Route for PATCH /v1/m-outlets/:outlet_id
	mOutletsRoute := app.Group("/v1/m-outlets", middleware.JWTProtected())
	mOutletsRoute.Patch("/:outlet_id", ctrl.Update)
}

// List - GET /v1/outlet-list
func (ctrl *OutletListController) List(c *fiber.Ctx) error {
	var filter entity.OutletListQueryFilter

	var requestId string
	if c.Locals("requestid") != nil {
		requestId = c.Locals("requestid").(string)
	}

	resp := responsebuild.BuildResponse(requestId, constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&filter); err != nil {
		log.Error("OutletListController, List, QueryParser:", err.Error())
	}

	var custId string
	if c.Locals("cust_id") != nil {
		custId = c.Locals("cust_id").(string)
	}

	data, total, lastPage, err := ctrl.Service.List(filter, custId)
	if err != nil {
		log.Error("OutletListController, List, Service.List:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	if len(data) == 0 {
		resp.Setmsg("No Data")
		resp.Setdata(nil)
	} else {
		resp.Setmsg("success")
		resp.Setdata(data)
	}

	resp.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

// Delete - DELETE /v1/outlet-list/:outlet_id
func (ctrl *OutletListController) Delete(c *fiber.Ctx) error {
	var params entity.OutletListParams
	requestId, _ := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestId, constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OutletListController, Delete, ParamsParser:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(resp.GetRespPayload())
	}

	errs := ctrl.Validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OutletListController, Delete, ValidateStruct:", errs)
		resp.Setmsg(fiber.ErrBadRequest.Message)
		resp.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	custId, _ := c.Locals("cust_id").(string)
	var userId int64
	if val, ok := c.Locals("user_id").(int64); ok {
		userId = val
	}

	err := ctrl.Service.Delete(custId, params.OutletId, userId)
	if err != nil {
		log.Error("OutletListController, Delete, Service.Delete:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			statusCode = fiber.StatusNotFound
		}
		resp.Setmsg(err.Error())
		return c.Status(statusCode).JSON(resp.GetRespPayload())
	}

	resp.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

// Update - PATCH /v1/m-outlets/:outlet_id
func (ctrl *OutletListController) Update(c *fiber.Ctx) error {
	var params entity.OutletListParams
	var body entity.UpdateOutletBody

	// SAFE ASSERTION
	var requestId string
	if c.Locals("requestid") != nil {
		requestId = c.Locals("requestid").(string)
	}

	resp := responsebuild.BuildResponse(requestId, constant.HEADER_ACCEPT_LANG)

	// Parse path params
	if err := c.ParamsParser(&params); err != nil {
		log.Error("OutletListController, Update, ParamsParser:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(resp.GetRespPayload())
	}

	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		log.Error("OutletListController, Update, BodyParser:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(resp.GetRespPayload())
	}

	// Validate
	errs := ctrl.Validator.ValidateStruct(body, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OutletListController, Update, ValidateStruct:", errs)
		resp.Setmsg(fiber.ErrBadRequest.Message)
		resp.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	// Set context values
	body.OutletId = params.OutletId

	// SAFE EXTRACTION
	if c.Locals("cust_id") != nil {
		body.CustId = c.Locals("cust_id").(string)
	}
	if c.Locals("parent_cust_id") != nil {
		body.ParentCustId = c.Locals("parent_cust_id").(string)
	}
	if val, ok := c.Locals("user_id").(int64); ok {
		body.UpdatedBy = val
	}

	err := ctrl.Service.Update(body)
	if err != nil {
		log.Error("OutletListController, Update, Service.Update:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			statusCode = fiber.StatusNotFound
		}
		resp.Setmsg(err.Error())
		return c.Status(statusCode).JSON(resp.GetRespPayload())
	}

	resp.Setmsg("Data successfully updated")
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}
