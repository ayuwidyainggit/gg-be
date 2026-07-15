package controller

import (
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/middleware"
	"finance/pkg/responsebuild"
	"finance/pkg/validation"
	"finance/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type DepositLookupController struct {
	DepositLookupService service.DepositLookupService
	validator            *validation.Validate
}

func NewDepositLookupController(DepositLookupService service.DepositLookupService, validator *validation.Validate) *DepositLookupController {
	return &DepositLookupController{
		DepositLookupService: DepositLookupService,
		validator:            validator,
	}
}

func (controller *DepositLookupController) Route(app *fiber.App) {

	grfRouteV1 := app.Group("/v1/deposit-lookup-filter", middleware.JWTProtected())
	grfRouteV1.Get("/", controller.LookupIndex)

	grfRouteV2 := app.Group("/v1/invoice-list-collection", middleware.JWTProtected())
	grfRouteV2.Get("/", controller.InvoiceCollectionList)

	grfRouteDepositBalance := app.Group("/v1/deposit-payment", middleware.JWTProtected())
	grfRouteDepositBalance.Get("/balance", controller.ListBalancePaymentDepositByCustId)

}

func (controller *DepositLookupController) parseAndValidateRequest(c *fiber.Ctx) (entity.GeneralQueryFilter, *responsebuild.DataRespReq, error) {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("DepositLookupController, parseAndValidateRequest, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return dataFilter, responsePayload, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("DepositLookupController, parseAndValidateRequest, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return dataFilter, responsePayload, fiber.NewError(fiber.StatusBadRequest, fiber.ErrBadRequest.Message)
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// // log.Println("BankController, List, CustId:", custId)

	return dataFilter, responsePayload, nil
}

func (controller *DepositLookupController) LookupIndex(c *fiber.Ctx) error {
	dataFilter, responsePayload, err := controller.parseAndValidateRequest(c)
	if err != nil {
		if err == fiber.ErrBadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	var (
		// err      error
		data     interface{}
		total    int64
		lastPage int
	)

	switch dataFilter.Mode {
	case "collection":
		data, total, lastPage, err = controller.DepositLookupService.LookupCollectionNo(dataFilter)
		if err != nil {
			log.Error("DepositLookupController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	case "deposit":
		data, total, lastPage, err = controller.DepositLookupService.LookupDepositNo(dataFilter)
		if err != nil {
			log.Error("DepositLookupController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.DepositLookupService.LookupDepositStatus(dataFilter)
		if err != nil {
			log.Error("DepositLookupController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
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

func (controller *DepositLookupController) InvoiceCollectionList(c *fiber.Ctx) error {
	dataFilter, responsePayload, err := controller.parseAndValidateRequest(c)
	if err != nil {
		if err == fiber.ErrBadRequest {
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.DepositLookupService.ListInvoiceByCollection(dataFilter)
	if err != nil {
		log.Error("DepositLookupController, List, data, err:", err.Error())
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

func (controller *DepositLookupController) ListBalancePaymentDepositByCustId(c *fiber.Ctx) error {
	var dataFilter entity.DepositLookupQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("DepositController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("DepositController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.DepositLookupService.ListBalancePaymentDepositByCustId(dataFilter)
	if err != nil {
		log.Error("DepositController, List, data, err:", err.Error())
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
