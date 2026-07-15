package controller

import (
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/middleware"
	"finance/pkg/responsebuild"
	"finance/pkg/validation"
	"finance/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CoreTaxVatExtractController struct {
	Service   service.CoreTaxVatExtractService
	validator *validation.Validate
}

func NewCoreTaxVatExtractController(service service.CoreTaxVatExtractService, validator *validation.Validate) *CoreTaxVatExtractController {
	return &CoreTaxVatExtractController{
		Service:   service,
		validator: validator,
	}
}

func (controller *CoreTaxVatExtractController) Route(app *fiber.App) {
	qParamId := ":coretax_vat_extract_id"
	CoreTaxVatExtractRouteV1 := app.Group("/v1/coretax-vat-extract", middleware.JWTProtected())
	CoreTaxVatExtractRouteV1.Post("", controller.Extract)
	CoreTaxVatExtractRouteV1.Get("", controller.List)
	CoreTaxVatExtractRouteV1.Get("/download-result/"+qParamId, controller.ExtractDownloadDetail)

}

func (controller *CoreTaxVatExtractController) Extract(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CoreTaxExtractReq
	if err := c.BodyParser(&request); err != nil {
		log.Error("CoretaxVatExtractController, Generate, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("CoretaxVatExtractController, Generate, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.Service.Extract(request)
	if err != nil {
		log.Error("CoretaxVatExtractController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)

	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())

}

func (controller *CoreTaxVatExtractController) List(c *fiber.Ctx) error {
	var dataFilter entity.CoreTaxVatExtractQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("CoretaxTaxesController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("CoretaxTaxesController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	datas, total, lastPage, err := controller.Service.List(dataFilter)
	if err != nil {
		log.Error("CoretaxTaxesController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(datas)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())

}

func (controller *CoreTaxVatExtractController) ExtractDownloadDetail(c *fiber.Ctx) error {
	var params entity.CoretaxVatExtractParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("CoretaxTaxesController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("CoretaxTaxesController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.Service.ExtractDownloadResult(params.CoretaxVatExtractID, custId, parentCustId)
	if err != nil {
		log.Error("CndnController, Detail, FindOneByOutletId, err:", err.Error())
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
