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

type InvoicesController struct {
	InvoicesService service.InvoicesService
	validator       *validation.Validate
}

func NewInvoicesController(
	InvoicesService service.InvoicesService,
	validator *validation.Validate,
) *InvoicesController {
	return &InvoicesController{
		InvoicesService: InvoicesService,
		validator:       validator,
	}
}

func (controller *InvoicesController) Route(app *fiber.App) {
	userRouteV1 := app.Group("/v1/invoices", middleware.JWTProtected())
	userRouteV1.Get("/payment/:invoice_no", controller.GetPayment)
	userRouteV1.Post("/payment", controller.CreatePayment)

}
func (controller *InvoicesController) CreatePayment(c *fiber.Ctx) error {
	var (
		ctx             = c.UserContext()
		responsePayload = responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
		request         entity.InvoicesPaymentCreate
	)

	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	var (
		custID, _   = c.Locals("cust_id").(string)
		userID, _   = c.Locals("user_id").(int64)
		EmpID, _    = c.Locals("emp_id").(int64)
		EmpGrpID, _ = c.Locals("emp_grp_id").(int64)
	)

	request.CustID = custID
	request.CreatedBy = userID
	request.EmpGrpID = EmpGrpID
	request.EmpID = EmpID
	request.SalesmanID = EmpID

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.InvoicesService.CreatePayment(ctx, request)
	if err != nil {
		log.Error("CollectionController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *InvoicesController) GetPayment(c *fiber.Ctx) error {
	// var (
	// 	headerAcceptLang string
	// 	request          entity.InvoicesListReq
	// )
	// if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
	// 	headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	// }
	// responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	// if err := c.QueryParser(&request); err != nil {
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	// }
	// errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	// if errs != nil {
	// 	// log.Println("errs, UserController-Login:", structs.StructToJson(errs))
	// 	responsePayload.Setmsg(fiber.ErrBadRequest.Message)
	// 	responsePayload.Seterrors(errs)
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }
	// datas, err := controller.InvoicesService.GetInvoices(request)
	// if err != nil {
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }
	// // log.Println("response:", response)
	// responsePayload.Setmsg(constant.STATUS_OK)
	// responsePayload.Setdata(datas)
	// return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

	var params entity.DetailInvoiceParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("CollectionController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.InvoicesService.DetailInvoice(params.InvoiceNo, custId)
	if err != nil {
		log.Error("CollectionController, Detail, FindOneByOutletId, err:", err.Error())
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
