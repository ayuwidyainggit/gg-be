package controller

import (
	"strings"

	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type PromotionController struct {
	PromotionService service.PromotionService
	validator        *validation.Validate
}

func NewPromotionController(
	promotionService service.PromotionService,
	validator *validation.Validate,
) *PromotionController {
	return &PromotionController{
		PromotionService: promotionService,
		validator:        validator,
	}
}

func (controller *PromotionController) Route(app *fiber.App) {
	// qParamId := ":promotion_id"
	promotionRouteV1 := app.Group("/v1/promotions", middleware.JWTProtected())
	promotionRouteV1.Get("", controller.List)
	promotionRouteV1.Post("/consult", controller.Consult)
	promotionRouteV1.Get("/list", controller.ListMobile)
	promotionRouteV1.Get("/:promo_id", controller.DetailMobile)
	promotionRouteV1.Get("/outlet/:ot_type_id", controller.OutletList)
}

func (controller *PromotionController) List(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.PromotionsQueryFilter
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromotionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	// log.Println("custId:", custId)
	// log.Println("parentCustId:", parentCustId)

	data, total, lastPage, err := controller.PromotionService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("PromotionController, List, data, err:", err.Error())
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

func (controller *PromotionController) Consult(c *fiber.Ctx) error {
	var (
		request entity.ConsultPromotionBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, Consult, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("PromotionController, Consult, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responses, err := controller.PromotionService.ConsultPromotion(request)
	if err != nil {
		log.Error("PromotionController, ConsultPromotion, Service.ConsultPromotion, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(responses)
	responsePayload.Setmsg("Promotion Consulted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// ListMobile returns active promotions for mobile app
func (controller *PromotionController) ListMobile(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.PromotionMobileListQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromotionController, ListMobile, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 20
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "promo_desc:asc"
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, ListMobile, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.PromotionService.ListMobile(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("PromotionController, ListMobile, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Handle empty data case
	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
	} else {
		responsePayload.Setmsg("Success")
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

func (controller *PromotionController) DetailMobile(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		params           entity.PromotionMobileDetailParams
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, DetailMobile, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, DetailMobile, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, err := controller.PromotionService.DetailMobile(params.PromoID, custId, parentCustId)
	if err != nil {
		log.Error("PromotionController, DetailMobile, data, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "record not found" || strings.Contains(err.Error(), "not found") {
			statusCode = fiber.StatusNotFound
			errMsg = "Promotion not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg("Success")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) OutletList(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		params           entity.PromotionOutletListParams
		dataFilter       entity.PromotionOutletListQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, OutletList, ParamsParser:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromotionController, OutletList, query parser filter:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}
	if dataFilter.Limit > 100 {
		dataFilter.Limit = 100
	}

	// Validate params
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, OutletList, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.PromotionService.OutletList(dataFilter, params.OtTypeID, custId, parentCustId)
	if err != nil {
		log.Error("PromotionController, OutletList, data, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := fiber.ErrBadRequest.Message
		if err.Error() == constant.STATUS_DB_NOT_FOUND || strings.Contains(err.Error(), constant.NOT_FOUND) || strings.Contains(err.Error(), constant.RECORD_NOT_FOUND) {
			statusCode = fiber.StatusNotFound
			errMsg = constant.STATUS_DB_NOT_FOUND
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	// Handle empty data case
	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
	} else {
		responsePayload.Setmsg(constant.STATUS_OK)
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
