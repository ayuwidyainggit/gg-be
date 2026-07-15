package controller

import (
	"database/sql"
	"errors"
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/sql_helper"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MOutletController struct {
	MOutletService service.MOutletService
	validator      *validation.Validate
}

func NewMOutletController(mOutletService service.MOutletService, validator *validation.Validate) *MOutletController {
	return &MOutletController{
		MOutletService: mOutletService,
		validator:      validator,
	}
}

func (controller *MOutletController) Route(app *fiber.App) {

	qParamId := ":outlet_id"
	mOutletRouteV1 := app.Group("/v1/m-outlets", middleware.JWTProtected())
	mOutletRouteV1.Post("", controller.Create)
	mOutletRouteV1.Post("/from-list", controller.CreateFromList)
	mOutletRouteV1.Get("", controller.List)
	mOutletRouteV1.Get("/additionals", controller.ListOutletAdditionals)
	mOutletRouteV1.Delete("/"+qParamId, controller.Delete)
	mOutletRouteV1.Get("/m-setup-outlet-check", controller.MSetupOutletCheck)
	mOutletRouteV1.Get("/"+qParamId, controller.Detail)

	// Mobile route for outlet list - sesuai docs: /v1/outlet-list
	mobileOutletRouteV1 := app.Group("/v1/outlet-list", middleware.JWTProtected())
	mobileOutletRouteV1.Get("", controller.MobileOutletList)
	mobileOutletRouteV1.Get("/"+qParamId, controller.MobileOutletDetail)
	mobileOutletRouteV1.Delete("/"+qParamId, controller.MobileOutletDelete)

	outletPJPRouteV1 := app.Group("/v1/outlet-pjp", middleware.JWTProtected())
	outletPJPRouteV1.Get("", controller.OutletPJPList)

	// Route for region info by distributor ID
	mOutletRouteV1.Get("/distributor/:distributor_id/region", controller.GetRegionByDistributorID)
}

func (controller *MOutletController) Create(c *fiber.Ctx) error {
	var request entity.CreateMOutletBody
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// userId := c.Locals("user_id").(int64)
	empId := c.Locals("emp_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Create, CustId:", custId)

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.CreatedBy = empId
	request.UpdatedBy = empId

	// Custom validation rules
	var customErrors []map[string]interface{}

	if request.OutletCode != "" && len(request.OutletCode) > 10 {
		customErrors = append(customErrors, map[string]interface{}{
			"key":     "outlet_code",
			"message": "Maksimal Outlet Code 10 Karakter",
		})
	}

	if len(request.OutletName) > 75 {
		customErrors = append(customErrors, map[string]interface{}{
			"key":     "outlet_name",
			"message": "Maksimal Outlet Name 75 Karakter",
		})
	}

	// Validate contact details
	for i, contact := range request.Details.OutletContact {
		if contact.ContactName != nil && len(*contact.ContactName) > 50 {
			customErrors = append(customErrors, map[string]interface{}{
				"key":     "contact_name",
				"message": "Maksimal Contact Name 50 Karakter",
			})
			_ = i // suppress unused warning
		}
		if len(contact.JobTitle) > 20 {
			customErrors = append(customErrors, map[string]interface{}{
				"key":     "job_title",
				"message": "Maksimal position 20 Karakter",
			})
		}
	}

	if len(customErrors) > 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(customErrors)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if len(errs) != 0 {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Set source to 1 (Mobile)
	request.Source = 1
	request.CreditLimitAction = 1
	request.CreditLimitActionName = "Warning"

	ctx := c.UserContext()
	custID := c.Locals("cust_id").(string)
	distributorCode := c.Locals("distributor_code").(string)
	distributorID := c.Locals("distributor_id").(int64)
	_, err := controller.MOutletService.Store(ctx, custID, distributorCode, distributorID, request)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())

}

func (controller *MOutletController) CreateFromList(c *fiber.Ctx) error {
	var request entity.ExtraCallOutlet
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	var (
		ctx           = c.UserContext()
		custID        = c.Locals("cust_id").(string)
		empID         = c.Locals("emp_id").(int64)
		parentCustID  = c.Locals("parent_cust_id").(string)
		isDistributor = c.Locals("is_distributor").(bool)
	)

	request.CustID = custID
	request.ParentCustID = parentCustID
	request.EmpID = empID
	request.IsDistributor = isDistributor

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.MOutletService.StoreFromList(ctx, request)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())

}

func (controller *MOutletController) List(c *fiber.Ctx) error {
	var dataFilter entity.MOutletQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OutletController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.MOutletService.List(dataFilter, custId, parentCustId, empId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
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

func (controller *MOutletController) ListOutletAdditionals(c *fiber.Ctx) error {
	var dataFilter entity.MOutletQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OutletController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.MOutletService.ListOutletAdditionals(dataFilter, custId, empId, parentCustId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
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

func (controller *MOutletController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteOutletParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("OutletController, Delete, CustId:", custId)

	err := controller.MOutletService.Delete(custId, params.OutletId, userId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) Detail(c *fiber.Ctx) error {
	var (
		params       entity.DetailOutletParams
		custID       = c.Locals("cust_id").(string)
		parentCustID = c.Locals("parent_cust_id").(string)
		empID        = c.Locals("emp_id").(int64)
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.MOutletService.Detail(params.OutletId, empID, custID, parentCustID)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		errMsg := fiber.ErrBadRequest.Message
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) MobileOutletList(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.MobileOutletListQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}

	if dataFilter.Limit <= 0 || dataFilter.Limit > 999 {
		dataFilter.Limit = 5
	}

	if dataFilter.Sort == "" {
		dataFilter.Sort = "outlet_code:asc"
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	empId := c.Locals("emp_id").(int64)

	data, total, lastPage, err := controller.MOutletService.MobileOutletList(dataFilter, custId, empId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrInternalServerError.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.STATUS_DB_NOT_FOUND)
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) MobileOutletDelete(c *fiber.Ctx) error {
	var params entity.DeleteOutletParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	err := controller.MOutletService.Delete(custId, params.OutletId, userId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) MSetupOutletCheck(c *fiber.Ctx) error {
	var query struct {
		Year   int      `query:"year" binding:"required"`
		Status []string `query:"status" binding:"required"`
	}

	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	custId := c.Locals("cust_id").(string)

	outlet, err := controller.MOutletService.MOutletCheck(c.UserContext(), query.Year, custId, query.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": outlet})
}

func (controller *MOutletController) MobileOutletDetail(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		params           entity.DetailOutletParams
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	data, err := controller.MOutletService.MobileOutletDetail(params.OutletId, custId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responsePayload.Setmsg("record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(fiber.ErrInternalServerError.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) OutletPJPList(c *fiber.Ctx) error {
	var params entity.OutletPJPListQuery
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.QueryParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Retrieve User ID
	empID := c.Locals("emp_id").(int64)
	if empID == 0 {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Re-assign
	if params.Sort != "" {
		params.Sort, params.SortOrder = sql_helper.ParseSort(params.Sort)
	}

	params.EmpID = empID
	data, paging, err := controller.MOutletService.OutletPJPList(params)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrInternalServerError.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(paging)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MOutletController) GetRegionByDistributorID(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	distributorID, err := c.ParamsInt("distributor_id")
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	ctx := c.UserContext()
	region, err := controller.MOutletService.GetRegionByDistributorID(ctx, distributorID)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrInternalServerError.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.STATUS_OK)
	responsePayload.Setdata(region)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
