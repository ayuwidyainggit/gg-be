package controller

import (
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/validation"
	"sales/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ValidateOrderController struct {
	ValidateOrderService service.ValidateOrderService
	validator            *validation.Validate
}

func NewValidateOrderController(invoiceService service.ValidateOrderService, validator *validation.Validate) *ValidateOrderController {
	return &ValidateOrderController{
		ValidateOrderService: invoiceService,
		validator:            validator,
	}
}
func (controller *ValidateOrderController) Route(app *fiber.App) {
	invoiceRouteV1 := app.Group("/v1/validate-order", middleware.JWTProtected())
	invoiceRouteV1.Post("/detail", controller.ValidateOrderDetail)
	invoiceRouteV1.Post("/", controller.ValidateOrder)
}

func (controller *ValidateOrderController) ValidateOrder(c *fiber.Ctx) error {
	var request entity.ValidateOrderBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ValidateOrderController, BulkValidateOrder, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	ParentCustID := c.Locals("parent_cust_id").(string)

	request.CustID = custId
	request.ParentCustID = ParentCustID

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("ValidateOrderController, BulkValidateOrder, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, _, _, err := controller.ValidateOrderService.ValidateOrder(request)
	if err != nil {
		log.Error("InvoiceController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ValidateOrderController) ValidateOrderDetail(c *fiber.Ctx) error {
	var request entity.ValidateOrderDetailBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ValidateOrderController, BulkValidateOrder, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	ParentCustID := c.Locals("parent_cust_id").(string)

	request.CustID = custId
	request.ParentCustID = ParentCustID

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("ValidateOrderController, BulkValidateOrder, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, _, _, err := controller.ValidateOrderService.ValidateOrderDetail(request)
	if err != nil {
		log.Error("InvoiceController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
