package controller

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type DistributorReplenishmentSetupController struct {
	svc       service.DistributorReplenishmentSetupService
	validator *validation.Validate
}

func NewDistributorReplenishmentSetupController(svc service.DistributorReplenishmentSetupService, validator *validation.Validate) *DistributorReplenishmentSetupController {
	return &DistributorReplenishmentSetupController{svc: svc, validator: validator}
}

func (ctrl *DistributorReplenishmentSetupController) Route(app *fiber.App) {
	g := app.Group("/v1/distributor-replenishment-setup", middleware.JWTProtected())
	g.Get("", ctrl.List)
	g.Get("/distributor", ctrl.ListDistributorsForPIC)
	g.Get("/supplier", ctrl.ListSuppliersForPIC)
	g.Get("/pic/:user_id", ctrl.ListByPicUser)
	g.Get("/:id", ctrl.Detail)
	g.Post("", ctrl.Create)
	g.Put("/:id", ctrl.Update)
	g.Delete("/:id", ctrl.Delete)
}

func (ctrl *DistributorReplenishmentSetupController) ListDistributorsForPIC(c *fiber.Ctx) error {
	var f entity.DistributorReplenishmentDistributorQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.QueryParser(&f); err != nil {
		log.Info("DistributorReplenishmentSetupController, ListDistributorsForPIC, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 5
	}
	if f.Sort == "" {
		f.Sort = "created_date:desc"
	}
	f.CustId = c.Locals("cust_id").(string)
	f.ParentCustID = c.Locals("parent_cust_id").(string)

	if errs := ctrl.validator.ValidateStruct(f, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := ctrl.svc.ListDistributorsForPic(f)
	if err != nil {
		log.Info("DistributorReplenishmentSetupController, ListDistributorsForPIC, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if data == nil {
		data = []entity.DistributorReplenishmentDistributorItem{}
	}

	responsePayload.Setmsg(constant.SUCCESS_GET_DISTRIBUTOR_REPLENISHMENT_SETUP)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.DistributorReplenishmentSetupPaging{
		TotalRecord: total,
		PageCurrent: f.Page,
		PageLimit:   f.Limit,
		PageTotal:   lastPage,
		RequestID:   requestID,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// List GET /v1/distributor-replenishment-setup
func (ctrl *DistributorReplenishmentSetupController) List(c *fiber.Ctx) error {
	var dataFilter entity.DistributorReplenishmentSetupQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("DistributorReplenishmentSetupController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if ids := queryIntSliceFromRequest(c, "distributor_id"); len(ids) > 0 {
		dataFilter.DistributorIDs = ids
	}
	if ids := queryIntSliceFromRequest(c, "supplier_id"); len(ids) > 0 {
		dataFilter.SupplierIDs = ids
	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := ctrl.svc.List(dataFilter)
	if err != nil {
		log.Info("DistributorReplenishmentSetupController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if data == nil {
		data = []entity.DistributorReplenishmentSetupListItem{}
	}

	responsePayload.Setmsg(constant.SUCCESS_GET_DISTRIBUTOR_REPLENISHMENT_SETUP)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.DistributorReplenishmentSetupPaging{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
		RequestID:   requestID,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (ctrl *DistributorReplenishmentSetupController) ListSuppliersForPIC(c *fiber.Ctx) error {
	var f entity.DistributorReplenishmentSupplierQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.QueryParser(&f); err != nil {
		log.Info("DistributorReplenishmentSetupController, ListSuppliersForPIC, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if raw := strings.TrimSpace(c.Query("distributor_id")); raw != "" {
		if did, err := strconv.Atoi(raw); err == nil {
			f.DistributorID = &did
		}
	}

	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 5
	}
	if f.Sort == "" {
		f.Sort = "created_date:desc"
	}

	f.CustId = c.Locals("cust_id").(string)
	f.ParentCustID = c.Locals("parent_cust_id").(string)

	if errs := ctrl.validator.ValidateStruct(f, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := ctrl.svc.ListSuppliersForPic(f)
	if err != nil {
		log.Info("DistributorReplenishmentSetupController, ListSuppliersForPIC, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if data == nil {
		data = []entity.DistributorReplenishmentSupplierItem{}
	}

	responsePayload.Setmsg(constant.SUCCESS_GET_DISTRIBUTOR_REPLENISHMENT_SETUP)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.DistributorReplenishmentSetupPaging{
		TotalRecord: total,
		PageCurrent: f.Page,
		PageLimit:   f.Limit,
		PageTotal:   lastPage,
		RequestID:   requestID,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (ctrl *DistributorReplenishmentSetupController) ListByPicUser(c *fiber.Ctx) error {
	var params entity.DistributorReplenishmentSetupPicParams
	var dataFilter entity.DistributorReplenishmentSetupQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("DistributorReplenishmentSetupController, ListByPicUser, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("DistributorReplenishmentSetupController, ListByPicUser, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if ids := queryIntSliceFromRequest(c, "distributor_id"); len(ids) > 0 {
		dataFilter.DistributorIDs = ids
	}
	if ids := queryIntSliceFromRequest(c, "supplier_id"); len(ids) > 0 {
		dataFilter.SupplierIDs = ids
	}

	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustID = c.Locals("parent_cust_id").(string)

	errs := ctrl.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := ctrl.svc.ListByPicUserID(dataFilter, params.UserID)
	if err != nil {
		log.Info("DistributorReplenishmentSetupController, ListByPicUser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if data == nil {
		data = []entity.DistributorReplenishmentSetupListItem{}
	}

	responsePayload.Setmsg(constant.SUCCESS_GET_DISTRIBUTOR_REPLENISHMENT_SETUP)
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.DistributorReplenishmentSetupPaging{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
		RequestID:   requestID,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// Detail GET /v1/distributor-replenishment-setup/:id
func (ctrl *DistributorReplenishmentSetupController) Detail(c *fiber.Ctx) error {
	var params entity.DistributorReplenishmentSetupDetailParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("DistributorReplenishmentSetupController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := ctrl.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := ctrl.svc.Detail(params.ID, c.Locals("cust_id").(string), c.Locals("parent_cust_id").(string))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responsePayload.Setmsg(constant.RECORD_NOT_FOUND)
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		log.Info("DistributorReplenishmentSetupController, Detail, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if userName, ok := c.Locals("user_name").(string); ok && userName != "" {
		for i := range data.ApprovalData {
			if data.ApprovalData[i].BusinessUnit == 0 {
				data.ApprovalData[i].BusinessUnitName = userName
			}
		}
	}

	responsePayload.Setmsg(constant.SUCCESS_GET_DETAIL_DISTRIBUTOR_REPLENISHMENT_SETUP)
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (ctrl *DistributorReplenishmentSetupController) Create(c *fiber.Ctx) error {
	var body entity.DistributorReplenishmentSetupCreatePayload
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.BodyParser(&body); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	normalizeDistributorReplenishmentOptionalFields(&body)
	if errs := ctrl.validator.ValidateStruct(body, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	data, err := ctrl.svc.Create(body, custID, parentCustID, userID)
	if err != nil {
		log.Info("DistributorReplenishmentSetupController, Create, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (ctrl *DistributorReplenishmentSetupController) Update(c *fiber.Ctx) error {
	var params entity.DistributorReplenishmentSetupDetailParams
	var body entity.DistributorReplenishmentSetupCreatePayload
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.BodyParser(&body); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	normalizeDistributorReplenishmentOptionalFields(&body)
	if errs := ctrl.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := ctrl.validator.ValidateStruct(body, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	data, err := ctrl.svc.Update(params.ID, body, custID, parentCustID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "sql: no rows in result set" {
			responsePayload.Setmsg(constant.RECORD_NOT_FOUND)
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		log.Info("DistributorReplenishmentSetupController, Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (ctrl *DistributorReplenishmentSetupController) Delete(c *fiber.Ctx) error {
	var params entity.DistributorReplenishmentSetupDetailParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	responsePayload := responsebuild.BuildResponse(requestID, headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := ctrl.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	err := ctrl.svc.Delete(params.ID, custID, parentCustID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "sql: no rows in result set" {
			responsePayload.Setmsg(constant.RECORD_NOT_FOUND)
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		log.Info("DistributorReplenishmentSetupController, Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Success")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// queryIntSliceFromRequest supports distributor_id=1,2,3 and repeated distributor_id=1&distributor_id=2
func queryIntSliceFromRequest(c *fiber.Ctx, key string) []int {
	var out []int
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		if string(k) != key {
			return
		}
		s := strings.TrimSpace(string(v))
		if s == "" {
			return
		}
		if strings.Contains(s, ",") {
			out = append(out, parseCommaSeparatedInts(s)...)
			return
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		out = append(out, n)
	})
	return out
}

func normalizeDistributorReplenishmentOptionalFields(p *entity.DistributorReplenishmentSetupCreatePayload) {
	if p.WhLimitAction != nil && strings.TrimSpace(*p.WhLimitAction) == "" {
		p.WhLimitAction = nil
	}
}

func parseCommaSeparatedInts(raw string) []int {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		out = append(out, v)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
