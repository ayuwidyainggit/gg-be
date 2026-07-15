package controller

import (
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/validation"
	"sales/service"
	"sort"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type PromoTemplateController struct {
	PromoTemplateService service.PromoTemplateService
	validator            *validation.Validate
}

func NewPromoTemplateController(roService service.PromoTemplateService, validator *validation.Validate) *PromoTemplateController {
	return &PromoTemplateController{
		PromoTemplateService: roService,
		validator:            validator,
	}
}

func (controller *PromoTemplateController) Route(app *fiber.App) {
	qParamId := ":promo_template_id"
	roRouteV1 := app.Group("/v1/promo-templates", middleware.JWTProtected())
	roRouteV1.Post("", controller.Create)
	roRouteV1.Get("/statuses", controller.PromoStatus)
	roRouteV1.Get("/"+qParamId, controller.Detail)
	roRouteV1.Get("", controller.List)
	roRouteV1.Patch("/"+qParamId, controller.Update)
	roRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *PromoTemplateController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreatePromoTemplateBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("PromoTemplateController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("cust_id").(string)
	request.CreatedBy = c.Locals("user_fullname").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// reward_products.pro_id validation ( must be unique )
	var rewardProducts entity.PromoTemplateUniqueRewardProductID
	for _, row := range request.PromoTemplateRewardProduct {
		prp := entity.PromoTemplateRewardProduct{
			ProID: row.ProID,
		}
		rewardProducts.RewardProductID = append(rewardProducts.RewardProductID, prp)
	}
	errs = controller.validator.ValidateStruct(rewardProducts, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromoTemplateService.Store(request)
	if err != nil {
		log.Error("PromoTemplateController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Successfully added")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *PromoTemplateController) Detail(c *fiber.Ctx) error {
	var params entity.DetailPromoTemplateParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromoTemplateController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.PromoTemplateService.Detail(params)
	if err != nil {
		log.Error("PromoTemplateController, Detail, err:", err.Error())
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

func (controller *PromoTemplateController) List(c *fiber.Ctx) error {
	var (
		dataFilter entity.PromoTemplateQueryFilter
		data       []entity.PromoTemplate
	)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromoTemplateController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.PromoTemplateService.List(dataFilter)
	if err != nil {
		log.Error("PromoTemplateController, List, data, err:", err.Error())
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

func (controller *PromoTemplateController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdatePromoTemplateParams
		request entity.UpdatePromoTemplateBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromoTemplateController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromoTemplateController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// reward_products.pro_id validation ( must be unique )
	var rewardProducts entity.PromoTemplateUniqueRewardProductID
	for _, row := range request.PromoTemplateRewardProduct {
		prp := entity.PromoTemplateRewardProduct{
			ProID: row.ProID,
		}
		rewardProducts.RewardProductID = append(rewardProducts.RewardProductID, prp)
	}
	errs = controller.validator.ValidateStruct(rewardProducts, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Update, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromoTemplateService.Update(params.PromoTemplateID, request)
	if err != nil {
		log.Error("PromoTemplateController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromoTemplateController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailPromoTemplateParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromoTemplateController, Delete, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromoTemplateController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	deletedBy := c.Locals("user_fullname").(string)

	err := controller.PromoTemplateService.Delete(custId, params.PromoTemplateID, deletedBy)
	if err != nil {
		log.Error("PromoTemplateController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromoTemplateController) PromoStatus(c *fiber.Ctx) error {
	promoStatuses := make([]entity.PromoTemplateStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.PromoTemplateStatusDesc {
		promoTemplateStatus := entity.PromoTemplateStatus{
			PromoTemplateStatusID:   index,
			PromoTemplateStatusDesc: element,
		}
		promoStatuses = append(promoStatuses, promoTemplateStatus)
	}

	promoStatusesSorted := make(entity.PromoTemplateStatusDescSlice, 0)
	for _, row := range promoStatuses {
		promoStatusesSorted = append(promoStatusesSorted, row)
	}
	sort.Sort(promoStatusesSorted)
	responsePayload.Setdata(promoStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
