package controller

import (
	"fmt"
	"net/url"
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/validation"
	"sales/service"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

const (
	LogValidateStructParams = "PromotionController, ValidateStruct(params), errs:"
)

type PromotionController struct {
	PromotionService service.PromotionService
	validator        *validation.Validate
}

func NewPromotionController(roService service.PromotionService, validator *validation.Validate) *PromotionController {
	return &PromotionController{
		PromotionService: roService,
		validator:        validator,
	}
}

func (controller *PromotionController) Route(app *fiber.App) {
	qParamId := ":promo_id"
	qParamStatus := ":promo_status"
	roRouteV1 := app.Group("/v1/promotions", middleware.JWTProtected())
	roRouteV1.Post("", controller.Create)
	roRouteV1.Get("/statuses", controller.PromoStatus)
	roRouteV1.Get("/"+qParamId, controller.Detail)
	roRouteV1.Get("", controller.List)
	roRouteV1.Patch("/"+qParamId, controller.Update)
	roRouteV1.Post("/bulk-update-status", controller.BulkUpdateStatus)
	roRouteV1.Delete("/"+qParamId, controller.Delete)
	roRouteV1.Post("/consult", controller.Consult)

	roRouteV2 := app.Group("/v2/promotions", middleware.JWTProtected())
	roRouteV2.Post("", controller.CreateV2)
	roRouteV2.Patch("/update/"+qParamId+"/*", controller.UpdateV2)
	roRouteV2.Patch("/update-status/"+qParamStatus+"/"+qParamId+"/*", controller.UpdateV2Status)
	roRouteV2.Post("/duplicate/"+qParamId+"/*", controller.DuplicateV2)
	roRouteV2.Get("/statuses", controller.PromoStatusV2)
	roRouteV2.Get("", controller.ListV2)
	roRouteV2.Post("/consult", controller.ConsultV2)
	roRouteV2.Get("/"+qParamId+"/*", controller.DetailV2)
	roRouteV2.Post("/unit-conversion", controller.Conversion)
}

func (controller *PromotionController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreatePromotionBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("cust_id").(string)
	request.CreatedBy = c.Locals("user_fullname").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// reward_products.pro_id validation ( must be unique )
	var rewardProducts entity.UniqueRewardProductID
	for _, row := range request.RewardProduct {
		prp := entity.PromoRewardProduct{
			ProID: row.ProID,
		}
		rewardProducts.RewardProductID = append(rewardProducts.RewardProductID, prp)
	}
	errs = controller.validator.ValidateStruct(rewardProducts, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromotionService.Store(request)
	if err != nil {
		log.Error("PromotionController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Successfully added")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) Detail(c *fiber.Ctx) error {
	var params entity.DetailPromotionParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.PromotionService.Detail(params)
	if err != nil {
		log.Error("PromotionController, Detail, err:", err.Error())
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

func (controller *PromotionController) List(c *fiber.Ctx) error {
	var (
		dataFilter entity.PromotionQueryFilter
		data       []entity.Promotion
	)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromotionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error(LogValidateStructParams, errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.PromotionService.List(dataFilter)
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

func (controller *PromotionController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdatePromotionParams
		request entity.UpdatePromotionBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(LogValidateStructParams, errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// reward_products.pro_id validation ( must be unique )
	var rewardProducts entity.UniqueRewardProductID
	for _, row := range request.RewardProduct {
		prp := entity.PromoRewardProduct{
			ProID: row.ProID,
		}
		rewardProducts.RewardProductID = append(rewardProducts.RewardProductID, prp)
	}
	errs = controller.validator.ValidateStruct(rewardProducts, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Update, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromotionService.Update(params.PromoID, request)
	if err != nil {
		log.Error("PromotionController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailPromotionParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, Delete, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	deletedBy := c.Locals("user_fullname").(string)

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	err := controller.PromotionService.Delete(params, deletedBy)
	if err != nil {
		log.Error("PromotionController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) PromoStatus(c *fiber.Ctx) error {
	promoStatuses := make([]entity.PromotionStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.PromoStatusDesc {
		promoStatus := entity.PromotionStatus{
			PromoStatusID:   index,
			PromoStatusDesc: element,
		}
		promoStatuses = append(promoStatuses, promoStatus)
	}

	promoStatusesSorted := make(entity.PromoStatusDescSlice, 0)
	for _, row := range promoStatuses {
		promoStatusesSorted = append(promoStatusesSorted, row)
	}
	sort.Sort(promoStatusesSorted)
	responsePayload.Setdata(promoStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) BulkUpdateStatus(c *fiber.Ctx) error {
	var (
		request entity.BulkUpdateStatusPromotionBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, BulkUpdateStatus, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, BulkUpdateStatus, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromotionService.BulkUpdateStatus(request)
	if err != nil {
		log.Error("PromotionController, BulkUpdateStatus, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	promo := entity.Promotion{
		PromoStatusID: request.PromoStatusID,
	}
	promoStatusDesc := promo.GetPromoStatusDesc()
	responsePayload.Setmsg(promoStatusDesc + " Successfully")
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

func (controller *PromotionController) CreateV2(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreatePromotionV2Body
	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, CreateV2, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = fmt.Sprint(c.Locals("parent_cust_id"))
	request.ParentCustID = fmt.Sprint(c.Locals("parent_cust_id"))
	request.DistributorCustID = fmt.Sprint(c.Locals("cust_id"))
	request.CreatedBy = fmt.Sprint(c.Locals("user_fullname"))
	request.UpdatedBy = fmt.Sprint(c.Locals("user_fullname"))
	if request.PromoStatus == "" {
		request.PromoStatus = entity.PromoStatusDraft
	}
	if request.Coverage == "" {
		request.Coverage = entity.CoverageNational
	}

	return controller.processCreatePromotionV2(c, request, headerAcceptLang, "")
}

func (controller *PromotionController) processCreatePromotionV2(c *fiber.Ctx, request entity.CreatePromotionV2Body, headerAcceptLang string, successMsg string) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	rewardType := entity.RewardType("")
	for i, _ := range request.Slabs {
		request.Slabs[i].CustID = request.CustID
		request.Slabs[i].PromoID = request.PromoID
		rewardType = request.Slabs[i].RewardType
	}
	for i, _ := range request.Strata {
		request.Strata[i].CustID = request.CustID
		request.Strata[i].PromoID = request.PromoID
		rewardType = request.Strata[i].RewardType
		if request.IsClaimable {
			if request.ClaimType == entity.ClaimFull {
				// Full claim → strata pct harus kosong agar masuk validasi
				request.Strata[i].ClaimRealizationPct = nil
			} else if request.ClaimType == entity.ClaimPartial {
				if !request.Strata[i].Claimable {
					// Partial tapi strata tidak claimable  kosong
					request.Strata[i].ClaimRealizationPct = nil
				}
				// Jika claimable=true  biarkan nilai dari JSON (BodyParser)
			}
		} else {
			// isClaimable = false  semua harus kosong
			request.Strata[i].ClaimRealizationPct = nil
		}
	}

	// Tag-based validation
	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, CreateV2, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Header cross-field validation (dates, budget/claim, etc.)
	if problems := validateCreatePromotionV2(request, controller.PromotionService); len(problems) > 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(problems)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// SLAB / STRATA validation
	var childErrs []FieldError
	switch request.PromoType {
	case entity.PromotionTypeSlab:
		if len(request.Strata) > 0 {
			childErrs = append(childErrs, FieldError{Field: "strata", Message: "must be empty when promo_type='slab'"})
		}
		childErrs = append(childErrs, validateSlabs(request.Slabs, request)...)

	case entity.PromotionTypeStrata:
		if len(request.Slabs) > 0 {
			childErrs = append(childErrs, FieldError{Field: "slabs", Message: "must be empty when promo_type='strata'"})
		}

		// Validate strata_sequential field
		if len(request.Strata) > 0 {
			if request.StrataSequential == nil {
				childErrs = append(childErrs, FieldError{Field: "strata_sequential", Message: "required when promo_type='strata' and strata are provided"})
			} else {
				// Additional validation: if sequential is true, strata ranges should be chained
				if *request.StrataSequential && len(request.Strata) > 1 {
					// Sort strata by ordinal to check sequential chaining
					sortedStrata := make([]entity.PromoStrataItem, len(request.Strata))
					copy(sortedStrata, request.Strata)
					sort.Slice(sortedStrata, func(i, j int) bool {
						return sortedStrata[i].Ordinal < sortedStrata[j].Ordinal
					})

					// Check that ranges are properly chained (next.from >= prev.to)
					for i := 1; i < len(sortedStrata); i++ {
						if sortedStrata[i].RangeFrom < sortedStrata[i-1].RangeTo {
							childErrs = append(childErrs, FieldError{Field: "strata_sequential", Message: "when sequential=true, strata ranges must be chained (next range_from >= previous range_to)"})
							break
						}
					}
				}
			}
		}

		childErrs = append(childErrs, validateStrata(request.Strata, request)...)
	}

	if rewardType == entity.RewardTypeProduct {
		// validate reward products not empty
		if len(request.RewardProducts) == 0 {
			childErrs = append(childErrs, FieldError{Field: "reward_products", Message: "must be not empty when reward_type 'product'"})
		}
	} else {
		if len(request.RewardProducts) > 0 {
			childErrs = append(childErrs, FieldError{Field: "reward_products", Message: "must be empty when reward_type is not 'product'"})
		}
		// make reward products empty slice
		request.RewardProducts = make([]entity.CreatePromotionRewardProduct, 0)
	}

	// Budget Reference Validation
	if request.IsBudgetReference {
		if request.BudgetRefType == "" {
			childErrs = append(childErrs, FieldError{Field: "budget_ref_type", Message: "required when is_budget_reference=true"})
		} else {
			switch request.BudgetRefType {
			case entity.BudgetRefLimited:
				if request.BudgetAmount <= 0 {
					childErrs = append(childErrs, FieldError{Field: "budget_amount", Message: "required > 0 when budget_ref_type='limited'"})
				}
			case entity.BudgetRefUnlimited:
				if request.BudgetAmount > 0 {
					childErrs = append(childErrs, FieldError{Field: "budget_amount", Message: "must be 0 or empty when budget_ref_type='unlimited'"})
				}
			}
		}

		// Budget Control Level validation
		if request.BudgetControlLevel == "" {
			childErrs = append(childErrs, FieldError{Field: "budget_control_level", Message: "required when is_budget_reference=true"})
		}
	}

	// Product Criteria Validation
	childErrs = append(childErrs, validateProductCriteria(request)...)

	if len(childErrs) > 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(childErrs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.PromotionService.StoreV2(request)
	if err != nil {
		log.Error("PromotionController, processCreatePromotionV2, StoreV2, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if successMsg == "" {
		successMsg = constant.MsgOpenAPICreatePromotionSuccess
	}
	responsePayload.Setmsg(successMsg)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) DetailV2(c *fiber.Ctx) error {
	var params entity.DetailPromotionParams
	var headerAcceptLang string
	first := c.Params("promo_id")

	// Ambil ASD/ZZZ dari wildcard *
	rest := c.Params("*")
	if rest != "" {
		rest = strings.TrimPrefix(rest, "/")
	}

	// Rekonstruksi ID penuh
	fullPromoID := first
	if rest != "" {
		fullPromoID = first + "/" + rest
	}
	rawURL := c.OriginalURL()
	if strings.HasSuffix(rawURL, "/") && !strings.HasSuffix(fullPromoID, "/") {
		fullPromoID += "/"
	}
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, DetailV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	params.PromoID = fullPromoID

	// URL decode the promo_id parameter
	decodedPromoID, err := url.PathUnescape(params.PromoID)
	if err != nil {
		log.Error("PromotionController, DetailV2, PathUnescape:", err.Error())
		responsePayload.Setmsg("Invalid promo_id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	params.PromoID = decodedPromoID

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, DetailV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.PromotionService.DetailV2(params)
	if err != nil {
		log.Error("PromotionController, DetailV2, err:", err.Error())
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

func (controller *PromotionController) ListV2(c *fiber.Ctx) error {
	var (
		dataFilter entity.PromotionV2QueryFilter
		data       []entity.PromotionV2
	)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PromotionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error(LogValidateStructParams, errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustID = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)
	dataFilter.TokenDistID = c.Locals("distributor_id").(int64)

	data, total, lastPage, err := controller.PromotionService.ListV2(dataFilter)
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

func (controller *PromotionController) UpdateV2(c *fiber.Ctx) error {
	var (
		params       entity.UpdatePromotionV2Params
		paramsDetail entity.DetailPromotionParams
		request      entity.UpdatePromotionV2Body
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	first := c.Params("promo_id")

	// Ambil ASD/ZZZ dari wildcard *
	rest := c.Params("*")
	if rest != "" {
		rest = strings.TrimPrefix(rest, "/")
	}

	// Rekonstruksi ID penuh
	fullPromoID := first
	if rest != "" {
		fullPromoID = first + "/" + rest
	}

	rawURL := c.OriginalURL()
	if strings.HasSuffix(rawURL, "/") && !strings.HasSuffix(fullPromoID, "/") {
		fullPromoID += "/"
	}

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, UpdateV2, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(LogValidateStructParams, errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, UpdateV2, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	params.PromoID = fullPromoID
	log.Info("Full Promo ID:", params.PromoID)

	// URL decode the promo_id parameter
	decodedPromoID, err := url.PathUnescape(params.PromoID)
	if err != nil {
		log.Error("PromotionController, UpdateV2, PathUnescape:", err.Error())
		responsePayload.Setmsg("Invalid promo_id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	params.PromoID = decodedPromoID

	request.CustID = fmt.Sprint(c.Locals("cust_id"))
	request.ParentCustID = fmt.Sprint(c.Locals("parent_cust_id"))
	request.UpdatedBy = fmt.Sprint(c.Locals("user_fullname"))
	paramsDetail.PromoID = params.PromoID
	paramsDetail.CustID = fmt.Sprint(c.Locals("cust_id"))
	paramsDetail.ParentCustId = fmt.Sprint(c.Locals("parent_cust_id"))

	promoDetail, err := controller.PromotionService.DetailV2ForUpdate(paramsDetail)
	if err != nil {
		log.Error("PromotionController, DetailV2ForUpdate, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if promoDetail.CustID != request.CustID {
		responsePayload.Setmsg("Edit promotion can only be done by users with the same level")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if string(promoDetail.PromoStatus) != string(entity.PromoStatusDraft) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not draft")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.PromoStatus = promoDetail.PromoStatus
	request.PromoID = params.PromoID

	for i, _ := range request.Slabs {
		request.Slabs[i].CustID = request.CustID
		request.Slabs[i].PromoID = request.PromoID
		request.RewardType = request.Slabs[i].RewardType
	}
	for i, _ := range request.Strata {
		request.Strata[i].CustID = request.CustID
		request.Strata[i].PromoID = request.PromoID
		request.RewardType = request.Strata[i].RewardType
	}

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, UpdateV2, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Header cross-field validation (dates, budget/claim, etc.)
	if problems := validateUpdatePromotionV2(request, controller.PromotionService); len(problems) > 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(problems)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// SLAB / STRATA validation
	var childErrs []FieldError
	switch request.PromoType {
	case entity.PromotionTypeSlab:
		if len(request.Strata) > 0 {
			childErrs = append(childErrs, FieldError{Field: "strata", Message: "must be empty when promo_type='slab'"})
		}
		childErrs = append(childErrs, validateSlabsUpdate(request.Slabs, request)...)

	case entity.PromotionTypeStrata:
		if len(request.Slabs) > 0 {
			childErrs = append(childErrs, FieldError{Field: "slabs", Message: "must be empty when promo_type='strata'"})
		}

		// Validate strata_sequential field
		if len(request.Strata) > 0 {
			if request.StrataSequential == nil {
				childErrs = append(childErrs, FieldError{Field: "strata_sequential", Message: "required when promo_type='strata' and strata are provided"})
			} else {
				// Additional validation: if sequential is true, strata ranges should be chained
				if *request.StrataSequential && len(request.Strata) > 1 {
					// Sort strata by ordinal to check sequential chaining
					sortedStrata := make([]entity.PromoStrataItem, len(request.Strata))
					copy(sortedStrata, request.Strata)
					sort.Slice(sortedStrata, func(i, j int) bool {
						return sortedStrata[i].Ordinal < sortedStrata[j].Ordinal
					})

					// Check that ranges are properly chained (next.from >= prev.to)
					for i := 1; i < len(sortedStrata); i++ {
						if sortedStrata[i].RangeFrom < sortedStrata[i-1].RangeTo {
							childErrs = append(childErrs, FieldError{Field: "strata_sequential", Message: "when sequential=true, strata ranges must be chained (next range_from >= previous range_to)"})
							break
						}
					}
				}
			}
		}

		childErrs = append(childErrs, validateStrataUpdate(request.Strata, request)...)
	}

	if request.RewardType == entity.RewardTypeProduct {
		// validate reward products not empty
		if len(request.RewardProducts) == 0 {
			childErrs = append(childErrs, FieldError{Field: "reward_products", Message: "must be not empty when reward_type 'product'"})
		}
	} else {
		if len(request.RewardProducts) > 0 {
			childErrs = append(childErrs, FieldError{Field: "reward_products", Message: "must be empty when reward_type is not 'product'"})
		}
		// make reward products empty slice
		request.RewardProducts = make([]entity.CreatePromotionRewardProduct, 0)
	}

	// Budget Reference Validation
	if request.IsBudgetReference {
		if request.BudgetRefType == "" {
			childErrs = append(childErrs, FieldError{Field: "budget_ref_type", Message: "required when is_budget_reference=true"})
		} else {
			switch request.BudgetRefType {
			case entity.BudgetRefLimited:
				if request.BudgetAmount <= 0 {
					childErrs = append(childErrs, FieldError{Field: "budget_amount", Message: "required > 0 when budget_ref_type='limited'"})
				}
			case entity.BudgetRefUnlimited:
				if request.BudgetAmount > 0 {
					childErrs = append(childErrs, FieldError{Field: "budget_amount", Message: "must be 0 or empty when budget_ref_type='unlimited'"})
				}
			}
		}

		// Budget Control Level validation
		if request.BudgetControlLevel == "" {
			childErrs = append(childErrs, FieldError{Field: "budget_control_level", Message: "required when is_budget_reference=true"})
		}
	}

	// Product Criteria Validation
	childErrs = append(childErrs, validateProductCriteriaUpdate(request)...)

	if len(childErrs) > 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(childErrs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.PromotionService.UpdateV2(params.PromoID, request)
	if err != nil {
		log.Error("PromotionController, UpdateV2, Service.UpdateV2, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) UpdateV2Status(c *fiber.Ctx) error {
	var (
		params       entity.UpdatePromotionV2StatusParams
		paramsDetail entity.DetailPromotionParams
		request      entity.UpdateStatusPromotionV2Body
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	first := c.Params("promo_id")

	// Ambil ASD/ZZZ dari wildcard *
	rest := c.Params("*")
	if rest != "" {
		rest = strings.TrimPrefix(rest, "/")
	}

	// Rekonstruksi ID penuh
	fullPromoID := first
	if rest != "" {
		fullPromoID = first + "/" + rest
	}

	rawURL := c.OriginalURL()
	if strings.HasSuffix(rawURL, "/") && !strings.HasSuffix(fullPromoID, "/") {
		fullPromoID += "/"
	}

	if err := c.ParamsParser(&params); err != nil {
		log.Error("PromotionController, UpdateV2Status, ParamsParser(params), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("PromotionController, UpdateV2Status, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.PromoID = fullPromoID

	// URL decode the promo_id parameter
	decodedPromoID, err := url.PathUnescape(params.PromoID)
	if err != nil {
		log.Error("PromotionController, UpdateV2Status, PathUnescape:", err.Error())
		responsePayload.Setmsg("Invalid promo_id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	params.PromoID = decodedPromoID

	paramsDetail.PromoID = params.PromoID
	paramsDetail.CustID = fmt.Sprint(c.Locals("cust_id"))
	paramsDetail.ParentCustId = fmt.Sprint(c.Locals("parent_cust_id"))
	request.CustID = paramsDetail.CustID
	request.ParentCustID = paramsDetail.ParentCustId

	promoDetail, err := controller.PromotionService.DetailV2ForUpdate(paramsDetail)
	if err != nil {
		log.Error("PromotionController, UpdateV2Status, DetailV2ForUpdate, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if promoDetail.CustID != fmt.Sprint(c.Locals("cust_id")) {
		responsePayload.Setmsg("Edit promotion can only be done by users with the same level")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status submit
	if params.PromoStatus == entity.PromoStatusSubmit && string(promoDetail.PromoStatus) != string(entity.PromoStatusDraft) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not draft")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status approved
	if params.PromoStatus == entity.PromoStatusApproved && string(promoDetail.PromoStatus) != string(entity.PromoStatusSubmit) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not submit")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status rejected
	if params.PromoStatus == entity.PromoStatusRejected && string(promoDetail.PromoStatus) != string(entity.PromoStatusSubmit) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not submit")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status inactive
	if params.PromoStatus == entity.PromoStatusInactive && string(promoDetail.PromoStatus) != string(entity.PromoStatusActive) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not active")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status active
	if params.PromoStatus == entity.PromoStatusActive &&
		string(promoDetail.PromoStatus) != string(entity.PromoStatusInactive) &&
		string(promoDetail.PromoStatus) != string(entity.PromoStatusApproved) &&
		string(promoDetail.PromoStatus) != string(entity.PromoStatusActive) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not inactive or approved")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// validate promo status closed
	if params.PromoStatus == entity.PromoStatusClosed && string(promoDetail.PromoStatus) != string(entity.PromoStatusActive) && string(promoDetail.PromoStatus) != string(entity.PromoStatusInactive) {
		responsePayload.Setmsg("You can't edit this data. Promotion status is not active")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("PromotionController, UpdateV2Status, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	err = controller.PromotionService.UpdateV2Status(params.PromoID, params.PromoStatus, request)
	if err != nil {
		log.Error("PromotionController, UpdateV2Status, Service.UpdateV2, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Status Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) DuplicateV2(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// Get promotion ID from URL parameter
	promoID := c.Params("promo_id")
	rest := c.Params("*")
	if rest != "" {
		rest = strings.TrimPrefix(rest, "/")
	}
	if promoID == "" {
		responsePayload.Setmsg("Promotion ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	promoIDFull := promoID
	if rest != "" {
		promoIDFull = promoID + "/" + rest
	}

	rawURL := c.OriginalURL()
	if strings.HasSuffix(rawURL, "/") && !strings.HasSuffix(promoIDFull, "/") {
		promoIDFull += "/"
	}

	// Get customer information from context
	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)

	// Create parameters for service call
	params := entity.DetailPromotionParams{
		PromoID:      promoIDFull,
		CustID:       custID,
		ParentCustId: parentCustID,
		UserFullname: fmt.Sprint(c.Locals("user_fullname")),
	}

	promoDetail, err := controller.PromotionService.DetailV2ForUpdate(params)
	if err != nil {
		log.Error("PromotionController, DuplicateV2, DetailV2ForUpdate, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if promoDetail.CustID != fmt.Sprint(c.Locals("cust_id")) {
		responsePayload.Setmsg("Duplicate promotion can only be done by users with the same level")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Call service to duplicate promotion
	newPromoID, err := controller.PromotionService.DuplicateV2(params)
	if err != nil {
		log.Error("Error duplicating promotion:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	// Return success response with new promotion ID
	responseData := map[string]string{
		"new_promo_id": newPromoID,
		"message":      "Promotion successfully duplicated",
	}
	responsePayload.Setdata(responseData)
	responsePayload.Setmsg("Promotion successfully duplicated")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) PromoStatusV2(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	statuses := []entity.PromotionV2Status{
		entity.PromoStatusDraft,
		entity.PromoStatusSubmit,
		entity.PromoStatusApproved,
		entity.PromoStatusRejected,
		entity.PromoStatusInactive,
		entity.PromoStatusActive,
		entity.PromoStatusClosed,
	}
	result := make([]string, 0, len(statuses))
	for _, s := range statuses {
		result = append(result, string(s))
	}

	responsePayload.Setdata(result)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

type FieldError struct {
	Field   string `json:"key"`
	Message string `json:"message"`
}

func validateCreatePromotionV2(req entity.CreatePromotionV2Body, svc service.PromotionService) (errs []FieldError) {
	add := func(f, m string) { errs = append(errs, FieldError{Field: f, Message: m}) }

	// Dates: format + range
	const layout = "2006-01-02"
	from, err1 := time.Parse(layout, req.EffectiveFrom)
	to, err2 := time.Parse(layout, req.EffectiveTo)
	if err1 != nil {
		add("effective_from", "must be YYYY-MM-DD")
	}
	if err2 != nil {
		add("effective_to", "must be YYYY-MM-DD")
	}
	if err1 == nil && err2 == nil && to.Before(from) {
		add("effective_to", "must be >= effective_from")
	}

	// Creation type: replacement needs existing_promo_id and must belong to same cust
	if req.PromoCreationType == entity.CreationTypeReplacement {
		if strings.TrimSpace(req.ExistingPromoID) == "" {
			add("existing_promo_id", "required when promo_creation_type is 'replacement'")
		} else {
			ok, err := svc.ExistsPromo(custID(req.CustID), req.ExistingPromoID) // implement in service/repo
			if err != nil {
				add("existing_promo_id", "lookup failed")
			} else if !ok {
				add("existing_promo_id", "not found for this cust_id")
			}
		}
	}

	// Budget rules
	if req.IsBudgetReference {
		if req.BudgetRefType == "" {
			add("budget_ref_type", "required when is_budget_reference = true")
		}
		if req.BudgetRefType == entity.BudgetRefLimited && req.BudgetAmount < 0 {
			add("budget_amount", "must be >= 0 when budget_ref_type = 'limited'")
		}
	} else {
		// normalize: ignore any budget settings when not referenced
	}

	// Claim rules
	if req.IsClaimable {
		if req.ClaimType == "" {
			add("claim_type", "required when is_claimable = true")
		}
		if req.ClaimType == entity.ClaimPartial && (req.ClaimRealizationPct < 0 || req.ClaimRealizationPct > 100) {
			add("claim_realization_pct", "must be between 0 and 100 when claim_type = 'partial'")
		}
	} else {
		// normalize: ignore claim fields
	}

	// Caps (type+value must be consistent)
	// If value > 0 then type must be provided; if type provided then value must be >= 0
	if req.MaxTotalRewardType == "" && req.MaxTotalRewardValue > 0 {
		add("max_total_reward_type", "required when max_total_reward_value > 0")
	}
	if req.MaxTotalRewardType != "" && req.MaxTotalRewardValue < 0 {
		add("max_total_reward_value", "must be >= 0")
	}

	// Header-vs-type consistency (no hard error here; child insert will be validated)
	if req.PromoType == entity.PromotionTypeSlab && req.StrataSequential != nil && *req.StrataSequential {
		add("strata_sequential", "must be false for 'slab' promo_type")
	}
	if req.PromoType == entity.PromotionTypeStrata && req.SlabMultiplied != nil && *req.SlabMultiplied {
		add("slab_multiplied", "must be false for 'strata' promo_type")
	}

	if req.PromoType == entity.PromotionTypeSlab && len(req.Slabs) == 0 {
		add("slabs", "required for 'slab' promo_type")
	}
	if req.PromoType == entity.PromotionTypeStrata && len(req.Strata) == 0 {
		add("strata", "required for 'strata' promo_type")
	}

	// Principal vs Distributor coverage validation
	// Principal: CustID == ParentCustID (coverage must not be empty)
	// Distributor: CustID != ParentCustID (coverage can be empty)
	if req.CustID == req.ParentCustID {
		// This is a Principal - coverage must not be empty
		if req.Coverage == "" {
			add("coverage", "required for login user as Principal")
		}

		if len(req.CoverageDistributors) == 0 && req.Coverage == entity.CoverageByDistributor {
			add("coverage_distributors", "required for login user as Principal and coverage is by distributor")
		}

		if len(req.CoverageDistributors) > 0 && req.Coverage == entity.CoverageNational {
			add("coverage_distributors", "must be empty for login user as Principal and coverage is national")
		}
	} else {
		// This is a Distributor - coverage can be empty
		if len(req.CoverageDistributors) > 0 {
			add("coverage_distributors", "must be empty for login user as Distributor")
		}
	}

	return
}

func validateSlabs(slabs []entity.PromoSlabItem, hdr entity.CreatePromotionV2Body) (errs []FieldError) {
	if len(slabs) == 0 {
		return
	}

	// Validate ordinals
	errs = append(errs, validateSlabOrdinals(slabs)...)

	// Sort slabs by ordinal for consistency checks
	copySlabs := append([]entity.PromoSlabItem(nil), slabs...)
	sort.Slice(copySlabs, func(i, j int) bool { return copySlabs[i].Ordinal < copySlabs[j].Ordinal })

	// Validate consistency across slabs
	errs = append(errs, validateSlabConsistency(copySlabs)...)

	// Validate multiplied rules
	errs = append(errs, validateSlabMultipliedRules(copySlabs, hdr)...)

	// Validate non-multiplied rules
	errs = append(errs, validateSlabNonMultipliedRules(copySlabs, hdr)...)

	// Validate individual slab properties
	errs = append(errs, validateSlabProperties(copySlabs)...)

	return errs
}

func validateSlabOrdinals(slabs []entity.PromoSlabItem) []FieldError {
	var errs []FieldError
	seen := map[int]bool{}
	ordinals := make([]int, 0, len(slabs))

	for i, s := range slabs {
		if s.Ordinal <= 0 {
			errs = append(errs, FieldError{fmt.Sprintf("slab[%d].ordinal", i), "must be >= 1"})
		}
		if seen[s.Ordinal] {
			errs = append(errs, FieldError{fmt.Sprintf("slab[%d].ordinal", i), "duplicate ordinal"})
		}
		seen[s.Ordinal] = true
		ordinals = append(ordinals, s.Ordinal)
	}

	// Check contiguous ordinals starting from 1
	sort.Ints(ordinals)
	for i := 0; i < len(ordinals); i++ {
		if ordinals[i] != i+1 {
			errs = append(errs, FieldError{fmt.Sprintf("ordinals must be contiguous starting from 1 (expected %d, got %d)", i+1, ordinals[i]), ""})
			break
		}
	}

	return errs
}

func validateSlabConsistency(slabs []entity.PromoSlabItem) []FieldError {
	var errs []FieldError
	var ruleType *entity.RuleType
	var rewardType *entity.RewardType

	for i, s := range slabs {
		if ruleType == nil {
			rt := s.RuleType
			ruleType = &rt
		} else if *ruleType != s.RuleType {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].rule_type", i), "must be same across slabs"})
		}

		if rewardType == nil {
			r := s.RewardType
			rewardType = &r
		} else if *rewardType != s.RewardType {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_type", i), "must be same across slabs"})
		}
	}

	return errs
}

func validateSlabMultipliedRules(slabs []entity.PromoSlabItem, hdr entity.CreatePromotionV2Body) []FieldError {
	var errs []FieldError

	if hdr.SlabMultiplied == nil || !*hdr.SlabMultiplied {
		return errs
	}

	if len(slabs) > 1 {
		errs = append(errs, FieldError{"slabs", "must be 1 when slab_multiplied=true"})
	}

	for i, s := range slabs {
		if s.RangeFrom != nil {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].range_from", i), "must be null when slab_multiplied=true"})
		}
		if s.RewardType == entity.RewardTypePercentage {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_type", i), "percentage not allowed when slab_multiplied=true"})
		}
	}

	return errs
}

func validateSlabNonMultipliedRules(slabs []entity.PromoSlabItem, hdr entity.CreatePromotionV2Body) []FieldError {
	var errs []FieldError

	// Only validate when slab_multiplied is explicitly false
	if hdr.SlabMultiplied != nil && !*hdr.SlabMultiplied {
		for i, s := range slabs {
			if s.RangeFrom == nil {
				errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].range_from", i), "not allowed NULL when slab_multiplied=false"})
			}
		}
	}

	return errs
}

func validateSlabProperties(slabs []entity.PromoSlabItem) []FieldError {
	var errs []FieldError
	var prevReward float64

	for i, s := range slabs {
		// Validate range
		errs = append(errs, validateSlabRange(s, i)...)

		// Validate reward value based on type
		errs = append(errs, validateSlabRewardValue(s, i)...)

		// Validate rule_uom when rule_type = "quantity"
		errs = append(errs, validateSlabRuleUom(s, i)...)

		// Validate monotonic reward
		if s.RewardType != entity.RewardTypeProduct {
			if s.RewardValue <= prevReward {
				errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_value", i), "must be strictly increasing by ordinal"})
			}
			v := s.RewardValue
			prevReward = v
		}
	}

	return errs
}

func validateSlabRange(s entity.PromoSlabItem, index int) []FieldError {
	var errs []FieldError

	if s.RangeFrom != nil && !(*s.RangeFrom < s.RangeTo) {
		errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].range_to", index), "must be > range_from"})
	}
	if s.RangeFrom == nil && !(0 < s.RangeTo) {
		errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].range_to", index), "must be > 0"})
	}

	return errs
}

func validateSlabRewardValue(s entity.PromoSlabItem, index int) []FieldError {
	var errs []FieldError

	if s.RewardValue < 1 {
		errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_value", index), fmt.Sprintf("required > 0 for reward type = %s", s.RewardType)})
	}

	switch s.RewardType {
	case "percentage":
		if s.RewardValue > 100 {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_value", index), "required 1..100 for percentage"})
		}
	case "fixed_value":
		if s.RewardValue < 1 {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_value", index), "required > 0 for fixed_value"})
		}
		if s.PerScope != "per_product" && s.PerScope != "per_order" {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].per_scope", index), "required (per_product|per_order) for fixed_value"})
		}
	case "product":
		// if s.RewardUom == "" {
		// 	errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].reward_uom", index), "required for product reward"})
		// }
	}

	return errs
}

func validateSlabRuleUom(s entity.PromoSlabItem, index int) []FieldError {
	var errs []FieldError

	// rule_uom is required when rule_type = "quantity"
	if s.RuleType == entity.RuleTypeQuantity {
		if s.RuleUom == "" {
			errs = append(errs, FieldError{fmt.Sprintf("slabs[%d].rule_uom", index), "required when rule_type='quantity'"})
		}
	}

	return errs
}

// validation_strata
func validateStrata(items []entity.PromoStrataItem, hdr entity.CreatePromotionV2Body) (errs []FieldError) {
	if len(items) == 0 {
		return
	}
	if len(items) > 5 {
		errs = append(errs, FieldError{"strata", "maximum 5 strata allowed"})
	}

	// promo_id must match header; ordinals unique & 1..5; sort by ordinal
	seen := map[int]struct{}{}
	copyStrata := append([]entity.PromoStrataItem(nil), items...)
	sort.Slice(copyStrata, func(i, j int) bool { return copyStrata[i].Ordinal < copyStrata[j].Ordinal })

	for i, s := range copyStrata {
		if s.PromoID != hdr.PromoID {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].promo_id", i), "must match header promo_id"})
		}
		if s.Ordinal < 1 || s.Ordinal > 5 {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].ordinal", i), "must be between 1 and 5"})
		}
		if _, ok := seen[s.Ordinal]; ok {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].ordinal", i), "duplicate ordinal"})
		}
		seen[s.Ordinal] = struct{}{}

		if !(s.RangeTo > s.RangeFrom) {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].range_to", i), "must be > range_from"})
		}

		// reward value presence/limits
		switch s.RewardType {
		case entity.RewardTypePercentage:
			if s.RewardValue < 0 || s.RewardValue > 100 {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].reward_value", i), "required 0..100 for percentage"})
			}
		case entity.RewardTypeFixedValue:
			if s.RewardValue < 0 {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].reward_value", i), "required >=0 for fixed_value"})
			}
			if s.PerScope != "per_product" && s.PerScope != "per_order" {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].per_scope", i), "required (per_product|per_order) for fixed_value"})
			}
		case entity.RewardTypeProduct:
			if s.RewardUom == "" {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].reward_uom", i), "required for product reward"})
			}
		}

		// rule_uom is required when rule_type = "quantity"
		if s.RuleType == entity.RuleTypeQuantity {
			if s.RuleUOM == "" {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].rule_uom", i), "required when rule_type='quantity'"})
			}
		}
	}

	// consistency: single rule_type & reward_type across strata; chain ranges & increasing rewards
	var ruleType *entity.RuleType
	var rewardType *entity.RewardType
	// var prevTo *float64
	// var prevReward *float64

	for i, s := range copyStrata {
		if ruleType == nil {
			rt := s.RuleType
			ruleType = &rt
		} else if *ruleType != s.RuleType {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].rule_type", i), "must be same across strata"})
		}
		if rewardType == nil {
			r := s.RewardType
			rewardType = &r
		} else if *rewardType != s.RewardType {
			errs = append(errs, FieldError{fmt.Sprintf("strata[%d].reward_type", i), "must be same across strata"})
		}

		// chain ranges: next.from >= prev.to
		// if prevTo != nil && !(s.RangeFrom >= *prevTo) {
		// 	errs = append(errs, FieldError{fmt.Sprintf("strata[%d].range_from", i), "must be >= previous range_to"})
		// }
		// t := s.RangeTo
		// prevTo = &t

		// increasing rewards when numeric
		// if s.RewardType != entity.RewardTypeProduct {
		// 	if prevReward != nil && s.RewardValue <= *prevReward {
		// 		errs = append(errs, FieldError{fmt.Sprintf("strata[%d].reward_value", i), "must be strictly increasing by ordinal"})
		// 	}
		// 	v := s.RewardValue
		// 	prevReward = &v
		// }

		// claim validation vs header claim settings
		if !hdr.IsClaimable {
			if s.Claimable || s.ClaimRealizationPct != nil {
				errs = append(errs, FieldError{fmt.Sprintf("strata[%d].claimable", i), "claim fields must be empty when header.is_claimable=false"})
			}
		} else {
			switch hdr.ClaimType {
			case entity.ClaimFull:
				if s.ClaimRealizationPct != nil {
					errs = append(errs, FieldError{fmt.Sprintf("strata[%d].claim_realization_pct", i), "must be empty when header.claim_type='full'"})
				}
			case entity.ClaimPartial:
				if s.Claimable {
					if s.ClaimRealizationPct != nil {
						if *s.ClaimRealizationPct < 0 || *s.ClaimRealizationPct > 100 {
							errs = append(errs, FieldError{fmt.Sprintf("strata[%d].claim_realization_pct", i), "required 0..100 when claimable=true & header.claim_type='partial'"})
						}
					}
				}
				// else if s.ClaimRealizationPct != nil {
				// 	errs = append(errs, FieldError{fmt.Sprintf("strata[%d].claim_realization_pct", i), "must be empty when claimable=false"})
				// }
			}
		}
	}

	return errs
}

// validateProductCriteria validates product criteria based on DDL and business rules
func validateProductCriteria(req entity.CreatePromotionV2Body) (errs []FieldError) {
	add := func(f, m string) { errs = append(errs, FieldError{Field: f, Message: m}) }

	// Product criteria validation
	if len(req.ProductCriteria) > 0 {
		// Validate individual product criteria items
		totalMandatory := 0
		seenProIDs := make(map[int64]bool)
		for i, criteria := range req.ProductCriteria {
			prefix := fmt.Sprintf("product_criteria[%d]", i)

			if criteria.Mandatory {
				totalMandatory++
			}

			// Check for duplicate pro_id within the same promotion
			if seenProIDs[criteria.ProID] {
				add(fmt.Sprintf("%s.pro_id", prefix), "duplicate product ID")
			}
			seenProIDs[criteria.ProID] = true

			// Validate pro_id is positive
			if criteria.ProID <= 0 {
				add(fmt.Sprintf("%s.pro_id", prefix), "must be > 0")
			}

			// Validate min_buy_type consistency based on DDL constraints
			if criteria.MinBuyType != nil {
				switch *criteria.MinBuyType {
				case entity.RuleTypeQuantity:
					// For quantity type: min_buy_qty required, min_buy_value must be null, min_buy_uom required
					if criteria.MinBuyQty == nil || *criteria.MinBuyQty < 0 {
						add(fmt.Sprintf("%s.min_buy_qty", prefix), "required >= 0 when min_buy_type='quantity'")
					}
					if criteria.MinBuyValue != nil {
						add(fmt.Sprintf("%s.min_buy_value", prefix), "must be empty when min_buy_type='quantity'")
					}
					if criteria.MinBuyUom == nil {
						add(fmt.Sprintf("%s.min_buy_uom", prefix), "required when min_buy_type='quantity'")
					}
				case entity.RuleTypeValue:
					// For value type: min_buy_value required, min_buy_qty must be null, min_buy_uom must be null
					if criteria.MinBuyValue == nil || *criteria.MinBuyValue < 0 {
						add(fmt.Sprintf("%s.min_buy_value", prefix), "required >= 0 when min_buy_type='value'")
					}
					if criteria.MinBuyQty != nil {
						add(fmt.Sprintf("%s.min_buy_qty", prefix), "must be empty when min_buy_type='value'")
					}
					if criteria.MinBuyUom != nil {
						add(fmt.Sprintf("%s.min_buy_uom", prefix), "must be empty when min_buy_type='value'")
					}
				}
			} else {
				// When min_buy_type is null, all min_buy fields should be null (DDL constraint)
				if criteria.MinBuyQty != nil {
					add(fmt.Sprintf("%s.min_buy_qty", prefix), "must be empty when min_buy_type is not provided")
				}
				if criteria.MinBuyValue != nil {
					add(fmt.Sprintf("%s.min_buy_value", prefix), "must be empty when min_buy_type is not provided")
				}
				if criteria.MinBuyUom != nil {
					add(fmt.Sprintf("%s.min_buy_uom", prefix), "must be empty when min_buy_type is not provided")
				}
			}
		}

		// Business rule: If product criteria are provided, there should be at least one mandatory product
		// This ensures the promotion has meaningful criteria

		// disable this rule for now
		/*
			hasMandatory := false
			for _, criteria := range req.ProductCriteria {
				if criteria.Mandatory {
					hasMandatory = true
					break
				}
			}
			if !hasMandatory {
				add("product_criteria", "at least one product must be marked as mandatory")
			}

			if req.MinimumSKU < totalMandatory {
				add("product_criteria", "minimum_sku must be >= total mandatory product")
			}
		*/
	}

	return errs
}

// helper if you want a strong type
func custID(v string) string { return v }

// validateUpdatePromotionV2 performs header-level cross-field validation for UpdatePromotionV2Body
func validateUpdatePromotionV2(req entity.UpdatePromotionV2Body, service service.PromotionService) (problems []FieldError) {
	// Date range validation
	if req.EffectiveFrom != "" && req.EffectiveTo != "" {
		from, err := time.Parse("2006-01-02", req.EffectiveFrom)
		if err == nil {
			to, err := time.Parse("2006-01-02", req.EffectiveTo)
			if err == nil {
				if to.Before(from) || to.Equal(from) {
					problems = append(problems, FieldError{Field: "effective_to", Message: "must be after effective_from"})
				}
			}
		}
	}

	// Budget validation
	if req.IsBudgetReference {
		if req.BudgetRefType == "" {
			problems = append(problems, FieldError{Field: "budget_ref_type", Message: "required when is_budget_reference=true"})
		} else {
			switch req.BudgetRefType {
			case entity.BudgetRefLimited:
				if req.BudgetAmount <= 0 {
					problems = append(problems, FieldError{Field: "budget_amount", Message: "required > 0 when budget_ref_type='limited'"})
				}
			case entity.BudgetRefUnlimited:
				if req.BudgetAmount > 0 {
					problems = append(problems, FieldError{Field: "budget_amount", Message: "must be 0 or empty when budget_ref_type='unlimited'"})
				}
			}
		}

		if req.BudgetControlLevel == "" {
			problems = append(problems, FieldError{Field: "budget_control_level", Message: "required when is_budget_reference=true"})
		}
	}

	// Coverage validation
	if req.Coverage == entity.CoverageByDistributor {
		if len(req.CoverageDistributors) == 0 {
			problems = append(problems, FieldError{Field: "coverage_distributors", Message: "required when coverage='by_distributor'"})
		}
	} else if req.Coverage == entity.CoverageNational {
		if len(req.CoverageDistributors) > 0 {
			problems = append(problems, FieldError{Field: "coverage_distributors", Message: "must be empty when coverage='national'"})
		}
	}

	// Claim validation
	if req.IsClaimable {
		if req.ClaimType == "" {
			problems = append(problems, FieldError{Field: "claim_type", Message: "required when is_claimable=true"})
		}
		if req.ClaimStartAfterDays < 0 {
			problems = append(problems, FieldError{Field: "claim_start_after_days", Message: "must be >= 0 when is_claimable=true"})
		}
		if req.ClaimRealizationPct < 0 || req.ClaimRealizationPct > 100 {
			problems = append(problems, FieldError{Field: "claim_realization_pct", Message: "must be between 0 and 100 when is_claimable=true"})
		}
	} else {
		// When claimable is false, ensure claim fields are empty/default
		if req.ClaimType != "" {
			problems = append(problems, FieldError{Field: "claim_type", Message: "must be empty when is_claimable=false"})
		}
		if req.ClaimStartAfterDays > 0 {
			problems = append(problems, FieldError{Field: "claim_start_after_days", Message: "must be 0 or empty when is_claimable=false"})
		}
		if req.ClaimRealizationPct > 0 {
			problems = append(problems, FieldError{Field: "claim_realization_pct", Message: "must be 0 or empty when is_claimable=false"})
		}
	}

	// Per-outlet limits validation
	if req.MaxTotalRewardType != "" {
		if req.MaxTotalRewardValue <= 0 {
			problems = append(problems, FieldError{Field: "max_total_reward_value", Message: "required > 0 when max_total_reward_type is set"})
		}
	}

	// Outlet criteria validation
	if req.OutletCriteria.SelectionType == "" {
		problems = append(problems, FieldError{Field: "outlet_criteria.selection_type", Message: "required"})
	} else {
		switch req.OutletCriteria.SelectionType {
		case "by_attribute":
			if len(req.OutletCriteria.OutletTypeIDs) == 0 && len(req.OutletCriteria.OutletGroupIDs) == 0 && len(req.OutletCriteria.OutletClassIDs) == 0 && len(req.OutletCriteria.SalesTeamIDs) == 0 {
				problems = append(problems, FieldError{Field: "outlet_criteria", Message: "at least one attribute (type, group, or class, sales team) must be selected when selection_type='by_attribute'"})
			}
		case "by_outlet":
			if len(req.OutletCriteria.OutletIDs) == 0 && len(req.OutletCriteria.SalesTeamIDs) == 0 {
				problems = append(problems, FieldError{Field: "outlet_criteria", Message: "at least one outlet or sales team must be selected when selection_type='by_outlet'"})
			}
		}
	}

	return
}

// validateSlabsUpdate validates slabs for UpdatePromotionV2Body
func validateSlabsUpdate(slabs []entity.PromoSlabItem, hdr entity.UpdatePromotionV2Body) (problems []FieldError) {
	if len(slabs) == 0 {
		problems = append(problems, FieldError{Field: "slabs", Message: "required when promo_type='slab'"})
		return
	}

	// Check for duplicate ordinals
	ordinalMap := make(map[int]bool)
	for i, slab := range slabs {
		if ordinalMap[slab.Ordinal] {
			problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].ordinal", i), Message: "duplicate ordinal"})
		}
		ordinalMap[slab.Ordinal] = true

		// Validate range
		if slab.RangeFrom != nil && *slab.RangeFrom < 0 {
			problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].range_from", i), Message: "must be >= 0"})
		}
		if slab.RangeTo <= 0 {
			problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].range_to", i), Message: "must be > 0"})
		}
		if slab.RangeFrom != nil && *slab.RangeFrom >= slab.RangeTo {
			problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].range_to", i), Message: "must be > range_from"})
		}

		// Validate reward
		if slab.RewardType == "" {
			problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].reward_type", i), Message: "required"})
		} else {
			switch slab.RewardType {
			case "percentage":
				if slab.RewardValue <= 0 || slab.RewardValue > 100 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].reward_value", i), Message: "must be between 0 and 100 for percentage"})
				}
			case "fixed_value":
				if slab.RewardValue <= 0 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("slabs[%d].reward_value", i), Message: "must be > 0 for fixed_value"})
				}
			}
		}
	}

	return
}

// validateStrataUpdate validates strata for UpdatePromotionV2Body
func validateStrataUpdate(strata []entity.PromoStrataItem, hdr entity.UpdatePromotionV2Body) (problems []FieldError) {
	if len(strata) == 0 {
		problems = append(problems, FieldError{Field: "strata", Message: "required when promo_type='strata'"})
		return
	}

	// Check for duplicate ordinals
	ordinalMap := make(map[int]bool)
	for i, s := range strata {
		if ordinalMap[s.Ordinal] {
			problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].ordinal", i), Message: "duplicate ordinal"})
		}
		ordinalMap[s.Ordinal] = true

		// Validate range
		if s.RangeFrom < 0 {
			problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].range_from", i), Message: "must be >= 0"})
		}
		if s.RangeTo <= 0 {
			problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].range_to", i), Message: "must be > 0"})
		}
		if s.RangeFrom >= s.RangeTo {
			problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].range_to", i), Message: "must be > range_from"})
		}

		// Validate reward
		if s.RewardType == "" {
			problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].reward_type", i), Message: "required"})
		} else {
			switch s.RewardType {
			case "percentage":
				if s.RewardValue <= 0 || s.RewardValue > 100 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].reward_value", i), Message: "must be between 0 and 100 for percentage"})
				}
			case "fixed_value":
				if s.RewardValue <= 0 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("strata[%d].reward_value", i), Message: "must be > 0 for fixed_value"})
				}
			}
		}
	}

	return
}

// validateProductCriteriaUpdate validates product criteria for UpdatePromotionV2Body
func validateProductCriteriaUpdate(req entity.UpdatePromotionV2Body) (problems []FieldError) {
	if len(req.ProductCriteria) == 0 {
		problems = append(problems, FieldError{Field: "product_criteria", Message: "required"})
		return
	}

	// Check for duplicate product IDs
	productMap := make(map[int64]bool)
	for i, criteria := range req.ProductCriteria {
		if productMap[criteria.ProID] {
			problems = append(problems, FieldError{Field: fmt.Sprintf("product_criteria[%d].pro_id", i), Message: "duplicate product ID"})
		}
		productMap[criteria.ProID] = true

		// Validate minimum buy requirements
		if criteria.MinBuyType != nil {
			switch *criteria.MinBuyType {
			case "qty":
				if criteria.MinBuyQty != nil && *criteria.MinBuyQty <= 0 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("product_criteria[%d].min_buy_qty", i), Message: "must be > 0 when min_buy_type='qty'"})
				}
			case "value":
				if criteria.MinBuyValue != nil && *criteria.MinBuyValue <= 0 {
					problems = append(problems, FieldError{Field: fmt.Sprintf("product_criteria[%d].min_buy_value", i), Message: "must be > 0 when min_buy_type='value'"})
				}
			}
		}
	}

	return
}

func (controller *PromotionController) ConsultV2(c *fiber.Ctx) error {
	var (
		request entity.ConsultPromoV2Req
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	tokenDistributorID := c.Locals("distributor_id").(int64)
	request.DistributorID = tokenDistributorID

	if request.CustID == request.ParentCustID {
		responsePayload.Setmsg("This endpoint only for Distributor user")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for i, detail := range request.Details {
		if detail.Qty1+detail.Qty2+detail.Qty3 <= 0 {
			responsePayload.Setmsg("Qty1, Qty2, Qty3 must be greater than 0")
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if detail.GrossValue == 0 && detail.SubTotal > 0 {
			request.Details[i].GrossValue = detail.SubTotal
		}
	}

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responses, err := controller.PromotionService.ConsultV2(request)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata([]entity.ConsultPromoResp{})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(responses)
	responsePayload.Setmsg("Consulted V2 Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PromotionController) Conversion(c *fiber.Ctx) error {
	var (
		request entity.PromoConversionReq
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("parent_cust_id").(string)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responses, err := controller.PromotionService.PromoConversion(request, request.CustID)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(responses)
	responsePayload.Setmsg("Conversion Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
