package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
)

type TakingOrderController struct {
	TakingOrderService service.TakingOrderService
	validator          *validation.Validate
}

func NewTakingOrderController(
	TakingOrderService service.TakingOrderService,
	validator *validation.Validate,
) *TakingOrderController {
	return &TakingOrderController{
		TakingOrderService: TakingOrderService,
		validator:          validator,
	}
}

func (controller *TakingOrderController) Route(app *fiber.App) {
	RouteV1 := app.Group("/v1/no-order-reasons", middleware.JWTProtected())
	RouteV1.Get("", controller.OrderReasonseGet)
}
func (controller *TakingOrderController) OrderReasonseGet(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.GeneralQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.TakingOrderService.List(dataFilter)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
