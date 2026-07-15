package controller

import (
	"fmt"
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type OutletController struct {
	OutletService service.OutletService
	validator     *validation.Validate
}

func NewOutletController(outletService service.OutletService, validator *validation.Validate) *OutletController {
	return &OutletController{
		OutletService: outletService,
		validator:     validator,
	}
}

func (controller *OutletController) Route(app *fiber.App) {

	qParamId := ":outlet_id"
	outletRouteV1 := app.Group("/v1/outlets", middleware.JWTProtected())
	outletRouteV1.Get("/verification-status", controller.VerificationStatusList)
	outletRouteV1.Get("/export", controller.Export)
	outletRouteV1.Get("/export-template", controller.ExportTemplate)
	outletRouteV1.Get("/export-template-new", controller.ExportTemplateNew)
	outletRouteV1.Get("/import-secondary-check", controller.ImportSecondaryCheck)
	outletRouteV1.Get("/export-template-update", controller.ExportTemplateUpdate)
	outletRouteV1.Post("/import", controller.Import)
	outletRouteV1.Post("/import-new", controller.ImportNew)
	outletRouteV1.Post("/import-update", controller.ImportUpdate)
	outletRouteV1.Get("/list-by-distributor", controller.ListByDistributor)
	outletRouteV1.Get("/"+qParamId, controller.Detail)
	outletRouteV1.Get("", controller.List)

	outletRouteV1.Post("", controller.Create)
	outletRouteV1.Post("/approve", controller.Approve)
	outletRouteV1.Post("/reject", controller.Reject)
	outletRouteV1.Post("/update-status", controller.UpdateMasterStatus)
	outletRouteV1.Patch("/update-status/"+qParamId, controller.UpdateStatus)
	outletRouteV1.Patch("/"+qParamId, controller.Update)
	outletRouteV1.Delete("/"+qParamId, controller.Delete)

	outletRouteV2 := app.Group("/v1/dropdown-outlet-type", middleware.JWTProtected())
	outletRouteV2.Get("/", controller.OutletTypes)
	outletRouteV3 := app.Group("/v1/dropdown-outlet-group", middleware.JWTProtected())
	outletRouteV3.Get("/", controller.OutletGroups)

	outletListRouteV1 := app.Group("/v1/outlet-list", middleware.JWTProtected())
	outletListRouteV1.Get("", controller.OutletListApproval)
	outletListRouteV1.Patch("/approval", controller.ApproveOutletList)
}

func (controller *OutletController) Detail(c *fiber.Ctx) error {
	var params entity.DetailOutletParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("OutletController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)
	lang := c.Locals("user_lang").(string)
	data, err := controller.OutletService.Detail(params.OutletId, custId, parentCustId, lang)
	if err != nil {
		log.Println("OutletController, Detail, err:", err.Error())
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
	responsePayload.Setmsg(constant.SUCCESS_NO_DATA)
	return c.JSON(responsePayload.GetRespPayload())
}

func parseOutletClassIDs(rawValue string) ([]int, error) {
	return parseCSVIntValues(rawValue, "ot_class_id")
}

func (controller *OutletController) List(c *fiber.Ctx) error {
	var dataFilter entity.OutletQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, List, query parser filter:", err.Error())
	}

	queryArgs := c.Context().QueryArgs()

	verificationStatus, err := parseIntSliceQuery(queryArgs, "verification_status", "verification_status", "verification_status[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(verificationStatus) > 0 {
		dataFilter.VerificationStatus = verificationStatus
	}

	outletIDs, err := parseIntSliceQuery(queryArgs, "outlet_id", "outlet_id", "outlet_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(outletIDs) > 0 {
		dataFilter.OutletID = outletIDs
	}

	otClassIDs, err := parseIntSliceQuery(queryArgs, "ot_class_id", "ot_class_id", "ot_class_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(otClassIDs) > 0 {
		dataFilter.OtClassID = otClassIDs
	}

	otTypeIDs, err := parseIntSliceQuery(queryArgs, "ot_type_id", "ot_type_id", "ot_type_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(otTypeIDs) > 0 {
		dataFilter.OtTypeID = otTypeIDs
	}

	otGrpIDs, err := parseIntSliceQuery(queryArgs, "ot_grp_id", "ot_grp_id", "ot_grp_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(otGrpIDs) > 0 {
		dataFilter.OtGrpID = otGrpIDs
	}

	distributorIDs, err := parseIntSliceQueryAllowZero(queryArgs, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorID = distributorIDs
	}

	outletStatusIDs, err := parseIntSliceQuery(queryArgs, "outlet_status", "outlet_status", "outlet_status[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(outletStatusIDs) > 0 {
		if !containsInt(outletStatusIDs, 0) {
			dataFilter.OutletStatusIDs = outletStatusIDs
		}
	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit < 1 {
		dataFilter.Limit = 10
	}

	parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.OutletService.FindParentCustId(custId)
		if err != nil {
			log.Println("OutletController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		parentCustId = parentCust.ParentCustId
	}

	data, total, lastPage, err := controller.OutletService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Println("OutletController, List, data, err:", err.Error())
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

func (controller *OutletController) OutletTypes(c *fiber.Ctx) error {

	var dataFilter entity.OutletQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, List, query parser filter:", err.Error())
	}

	parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.OutletService.FindParentCustId(custId)
		if err != nil {
			log.Println("OutletController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		parentCustId = parentCust.ParentCustId
	}
	log.Println("parentCustId:", parentCustId)
	log.Println("custId:", custId)

	// buat function di service untuk dapatkan data Outlet Types List yg ada di table mst.m_outlet

	data, total, lastPage, err := controller.OutletService.OutletTypeList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Println("OutletController, List, data, err:", err.Error())
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

func (controller *OutletController) OutletGroups(c *fiber.Ctx) error {

	var dataFilter entity.OutletQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, List, query parser filter:", err.Error())
	}

	parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.OutletService.FindParentCustId(custId)
		if err != nil {
			log.Println("OutletController, List, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		parentCustId = parentCust.ParentCustId
	}
	log.Println("parentCustId:", parentCustId)
	log.Println("custId:", custId)

	// buat function di service untuk dapatkan data Outlet Group List yg ada di table mst.m_outlet

	data, total, lastPage, err := controller.OutletService.OutletGroupList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Println("OutletController, List, data, err:", err.Error())
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

func (controller *OutletController) Create(c *fiber.Ctx) error {
	var request entity.CreateOutletBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("OutletController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)
	userName := ""
	if v := c.Locals("user_name"); v != nil {
		if s, ok := v.(string); ok {
			userName = s
		}
	}
	// log.Println("OutletController, Create, CustId:", custId)

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.CreatedBy = userId
	request.UpdatedBy = userId
	request.CreatedByName = userName

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.OutletService.Store(request)
	if err != nil {
		log.Println("OutletController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateOutletParams
		request entity.UpdateOutletRequest
	)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("OutletController, Update, ParamsParser(params):", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	var headerAcceptLang string
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Update, ValidateStruct(params), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("OutletController, Update, BodyParser(request), err:", err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Update, ValidateStruct(request), errs:", errs)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   fiber.ErrBadRequest.Message,
			Errors:    errs,
		})
	}

	request.ParentCustId = parentCustId
	err := controller.OutletService.Update(params.OutletId, request)
	if err != nil {
		log.Println("OutletController, Update, Service.Update, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.ApiResponse{
		RequestId: c.Locals("requestid").(string),
		Message:   constant.SUCCESSFULLY_UPDATED,
	})
}

func (controller *OutletController) UpdateStatus(c *fiber.Ctx) error {
	var params entity.UpdateOutletParams
	var request entity.UpdateOutletStatusRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("OutletController, UpdateStatus, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.BodyParser(&request); err != nil {
		log.Println("OutletController, UpdateStatus, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	err := controller.OutletService.UpdateStatus(int64(params.OutletId), custId, parentCustId, request, userId)
	if err != nil {
		log.Println("OutletController, UpdateStatus, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		if err.Error() == "outlet not found" {
			statusCode = fiber.StatusNotFound
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) UpdateMasterStatus(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	rowsAffected, err := controller.OutletService.UpdateStatuses()
	if err != nil {
		log.Println("OutletController, UpdateMasterStatus, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	responsePayload.Setdata(fiber.Map{
		"rows_affected": rowsAffected,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteOutletParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Println("OutletController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("OutletController, Delete, CustId:", custId)

	err := controller.OutletService.Delete(custId, params.OutletId, userId)
	if err != nil {
		log.Println("OutletController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) Approve(c *fiber.Ctx) error {
	var request entity.ApproveOutletBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("OutletController, Approve, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.VerifiedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Approve, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OutletService.Approve(request)
	if err != nil {
		log.Println("OutletController, Approve, Approve, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Successfully Approved")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) Reject(c *fiber.Ctx) error {
	var request entity.RejectOutletBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Println("OutletController, Reject, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.VerifiedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Println("OutletController, Reject, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := controller.OutletService.Reject(request); err != nil {
		log.Println("OutletController, Reject, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Successfully Rejected")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) VerificationStatusList(c *fiber.Ctx) error {

	var dataFilter entity.OutletQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, VerificationStatusList, query parser filter:", err.Error())
	}

	parentCustId := "" // c.Locals("parent_cust_id").(string)
	headerCustId := ""
	if len(c.GetReqHeaders()[constant.CUST_ID]) > 0 {
		headerCustId = c.GetReqHeaders()[constant.CUST_ID][0]
	}
	custId := headerCustId
	if custId == "" {
		custId = c.Locals("cust_id").(string)
		parentCustId = c.Locals("parent_cust_id").(string)
	} else {
		parentCust, err := controller.OutletService.FindParentCustId(custId)
		if err != nil {
			log.Println("OutletController, VerificationStatusList, query get parent cust id:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		custId = parentCust.CustId
		parentCustId = parentCust.ParentCustId
	}

	// buat function di service untuk dapatkan data Outlet Types List yg ada di table mst.m_outlet

	data, total, lastPage, err := controller.OutletService.VerificationStatusList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Println("OutletController, VerificationStatusList, data, err:", err.Error())
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

func (controller *OutletController) Export(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.OutletQueryFilter
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, Export, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	queryArgs := c.Context().QueryArgs()
	distributorIDs, err := parseIntSliceQuery(queryArgs, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorID = distributorIDs
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	// Panggil service yang sekarang mengembalikan buffer, content-type, dan nama file
	buffer, contentType, filename, err := controller.OutletService.Export(dataFilter)
	if err != nil {
		log.Println("OutletController, Export, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Set header HTTP secara dinamis berdasarkan hasil dari service
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if _, err := c.Write(buffer.Bytes()); err != nil {
		log.Println("OutletController, Export, Write response, err:", err.Error())
		return c.Status(http.StatusInternalServerError).SendString("Failed to export file")
	}

	return nil
}

func (controller *OutletController) ExportTemplate(c *fiber.Ctx) error {
	// Get format from query parameter (default ke xlsx)
	format := c.Query("format", "xlsx")

	// Validasi format yang didukung
	switch format {
	case "csv", "xls", "xlsx":
		// Format valid
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tidak didukung. Gunakan csv, xls, atau xlsx",
		})
	}

	additionalRaw := c.Query("additional", c.Query("aditional", ""))
	var additional []string
	if additionalRaw != "" {
		for _, s := range strings.Split(additionalRaw, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				additional = append(additional, s)
			}
		}
	}

	fieldsParam := c.Query("fields", "")
	var fields []string
	if fieldsParam != "" {
		for _, s := range strings.Split(fieldsParam, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				fields = append(fields, s)
			}
		}
	}

	buffer, contentType, filename, err := controller.OutletService.ExportTemplate(format, additional, fields)
	if err != nil {
		log.Println("OutletController, ExportTemplate, error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set header untuk download file
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(buffer.Bytes())
}

func (controller *OutletController) ExportTemplateNew(c *fiber.Ctx) error {
	format := c.Query("format", "xlsx")

	switch format {
	case "csv", "xls", "xlsx":
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tidak didukung. Gunakan csv, xls, atau xlsx",
		})
	}

	custId := c.Locals("cust_id").(string)
	if err := controller.OutletService.ValidateUploadSecondarySalesPermission(jwtDistributorID(c), custId); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	buffer, contentType, filename, err := controller.OutletService.ExportTemplateNew(format)
	if err != nil {
		log.Println("OutletController, ExportTemplateNew, error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buffer.Bytes())
}

func jwtDistributorID(c *fiber.Ctx) int64 {
	if v := c.Locals("distributor_id"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func (controller *OutletController) ImportSecondaryCheck(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	custId := c.Locals("cust_id").(string)
	data, err := controller.OutletService.ImportSecondaryCheck(jwtDistributorID(c), custId)
	if err != nil {
		log.Println("OutletController, ImportSecondaryCheck, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) ExportTemplateUpdate(c *fiber.Ctx) error {
	// Ambil format (default xlsx)
	format := c.Query("format", "xlsx")
	custId := c.Locals("cust_id").(string)

	// Ambil fields yang dipilih user
	fieldsParam := c.Query("fields", "")
	if fieldsParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fields harus diisi, contoh: ?fields=outlet_code,outlet_name,address",
		})
	}
	fields := strings.Split(fieldsParam, ",")

	// Panggil service
	buffer, contentType, filename, err := controller.OutletService.ExportTemplateUpdate(custId, format, fields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set header untuk download
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buffer.Bytes())
}

func (controller *OutletController) Import(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := strings.ToLower(strings.TrimSpace(c.Query("format")))
	if format == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx/xls) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		log.Println("Error getting file from form:", err)
		responsePayload.Setmsg("File upload is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if fileHeader.Size == 0 {
		responsePayload.Setmsg("File upload empty. Use form field 'file_upload' with Excel/CSV file containing data (do not send body without file). For curl, use: -F \"file_upload=@/path/to/file.xlsx\"")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Error opening uploaded file:", err)
		responsePayload.Setmsg("Failed to process uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)
	userName := ""
	if v := c.Locals("user_name"); v != nil {
		if s, ok := v.(string); ok {
			userName = s
		}
	}

	importRequest := entity.ImportRequest{
		File:          file,
		CustId:        custId,
		UserId:        userId,
		Filename:      fileHeader.Filename,
		Format:        format,
		ParentCustId:  parentCustId,
		CreatedByName: userName,
	}

	var importErr error
	switch format {
	case "csv":
		log.Println("Processing CSV file import...")
		importErr = controller.OutletService.ImportOutletsFromCSV(importRequest)
	case "xlsx":
		log.Println("Processing XLSX file import...")
		importErr = controller.OutletService.ImportOutletsFromXLSX(importRequest)
	case "xls":
		log.Println("Processing XLS file import...")
		importErr = controller.OutletService.ImportOutletsFromXLSX(importRequest)
	default:
		responsePayload.Setmsg(fmt.Sprintf("Unsupported file format: '%s'. Only 'csv', 'xlsx', and 'xls' are supported.", format))
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(responsePayload.GetRespPayload())
	}

	if importErr != nil {
		log.Println("Error during import process:", importErr)
		responsePayload.Setmsg(importErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) ImportNew(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := strings.ToLower(strings.TrimSpace(c.Query("format")))
	if format == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx/xls) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	if err := controller.OutletService.ValidateUploadSecondarySalesPermission(jwtDistributorID(c), custId); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		log.Println("OutletController, ImportNew, get file:", err.Error())
		responsePayload.Setmsg("File upload is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if fileHeader.Size == 0 {
		responsePayload.Setmsg("File upload empty. Use form field 'file_upload' with Excel/CSV file containing data.")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println("OutletController, ImportNew, open file:", err.Error())
		responsePayload.Setmsg("Failed to process uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)
	userName := ""
	if v := c.Locals("user_name"); v != nil {
		if s, ok := v.(string); ok {
			userName = s
		}
	}

	importRequest := entity.ImportRequest{
		File:          file,
		CustId:        custId,
		UserId:        userId,
		Filename:      fileHeader.Filename,
		Format:        format,
		ParentCustId:  parentCustId,
		CreatedByName: userName,
		IsImportNew:   true,
	}

	var importErr error
	switch format {
	case "csv":
		importErr = controller.OutletService.ImportOutletsFromCSV(importRequest)
	case "xlsx", "xls":
		importErr = controller.OutletService.ImportOutletsFromXLSX(importRequest)
	default:
		responsePayload.Setmsg(fmt.Sprintf("Unsupported file format: '%s'. Only 'csv', 'xlsx', and 'xls' are supported.", format))
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(responsePayload.GetRespPayload())
	}

	if importErr != nil {
		log.Println("OutletController, ImportNew, err:", importErr.Error())
		responsePayload.Setmsg(importErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OutletController) ImportUpdate(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := c.Query("format")
	if format == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		log.Println("Error getting file from form:", err)
		responsePayload.Setmsg("File upload is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if fileHeader.Size == 0 {
		responsePayload.Setmsg("File upload empty. Use form field 'file_upload' with file containing data. For curl, use: -F \"file_upload=@/path/to/file.xlsx\"")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Error opening uploaded file:", err)
		responsePayload.Setmsg("Failed to process uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	importRequest := entity.ImportRequest{
		File:         file,
		CustId:       custId,
		UserId:       userId,
		ParentCustId: parentCustId,
		Filename:     fileHeader.Filename,
		Format:       format,
	}

	var importErr error
	switch format {
	case "csv":
		log.Println("Processing CSV file update import...")
		importErr = controller.OutletService.ImportUpdateCSV(importRequest)
	case "xlsx":
		log.Println("Processing XLSX file update import...")
		importErr = controller.OutletService.ImportUpdateXLSX(importRequest)
	case "xls":
		log.Println("Processing XLS file update import...")
		importErr = controller.OutletService.ImportUpdateXLS(importRequest)
	default:
		responsePayload.Setmsg(fmt.Sprintf("Unsupported file format: '%s'. Only 'csv', 'xlsx', and 'xls' are supported.", format))
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(responsePayload.GetRespPayload())
	}

	if importErr != nil {
		log.Println("Error during import update process:", importErr)
		responsePayload.Setmsg(importErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// func (controller *OutletController) ExportTemplateUpdate(c *fiber.Ctx) error {
// 	// Ambil format (default xlsx)
// 	format := c.Query("format", "xlsx")
// 	custId := c.Locals("cust_id").(string)

// 	// Ambil fields yang dipilih user
// 	fieldsParam := c.Query("fields", "")
// 	if fieldsParam == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "fields harus diisi, contoh: ?fields=outlet_code,outlet_name,address",
// 		})
// 	}
// 	fields := strings.Split(fieldsParam, ",")

// 	// Panggil service
// 	buffer, contentType, filename, err := controller.OutletService.ExportTemplateUpdate(custId, format, fields)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	// Set header untuk download
// 	c.Set("Content-Type", contentType)
// 	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
// 	return c.Send(buffer.Bytes())
// }

func (controller *OutletController) ListByDistributor(c *fiber.Ctx) error {
	var dataFilter entity.OutletQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Println("OutletController, ListByDistributor, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	queryArgs := c.Context().QueryArgs()

	otClassIDs, err := parseIntSliceQuery(queryArgs, "ot_class_id", "ot_class_id", "ot_class_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(otClassIDs) > 0 {
		dataFilter.OtClassID = otClassIDs
	}

	distributorIDs, err := parseIntSliceQuery(queryArgs, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorID = distributorIDs
	}

	custID := fmt.Sprint(c.Locals("cust_id"))
	parentCustID := fmt.Sprint(c.Locals("parent_cust_id"))

	dataFilter.CustId = custID
	dataFilter.ParentCustId = parentCustID

	if dataFilter.CustId != dataFilter.ParentCustId && len(dataFilter.DistributorID) == 0 {
		if distID, ok := c.Locals("distributor_id").(int64); ok {
			dataFilter.DistributorID = []int{int(distID)}
		}
	}

	data, total, lastPage, err := controller.OutletService.ListByDistributor(dataFilter)
	if err != nil {
		log.Println("OutletController, ListByDistributor, data, err:", err.Error())
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

// OutletListApproval handles GET request to retrieve outlet change request list with pagination.
// Parses query parameters, validates input, extracts cust_id from JWT context, and calls service layer.
// Returns paginated list of outlet change requests filtered by status, or "No Data" message if empty.
func (controller *OutletController) OutletListApproval(c *fiber.Ctx) error {
	var dataFilter entity.OutletListApprovalQueryFilter
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	data, total, lastPage, err := controller.OutletService.OutletListApproval(dataFilter, custId)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setmsg("")
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

// ApproveOutletList handles PATCH request to approve or reject outlet change requests.
// Parses request body, validates input, extracts cust_id and user_id from JWT context.
// Calls service layer to process approval/rejection and returns success or error response.
func (controller *OutletController) ApproveOutletList(c *fiber.Ctx) error {
	var request entity.OutletListApprovalRequest
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OutletService.ApproveOutletList(request, custId, userId)
	if err != nil {
		responsePayload.Setmsg("Approval failed")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Approval successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
