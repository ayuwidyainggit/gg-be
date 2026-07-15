package controller

import (
	"errors"
	"log"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type OutletCodeController struct {
	OutletCodeService service.OutletCodeService
}

func NewOutletCodeController(svc service.OutletCodeService) *OutletCodeController {
	return &OutletCodeController{OutletCodeService: svc}
}

func (ctrl *OutletCodeController) Route(app *fiber.App) {
	g := app.Group("/v1/outlet-code", middleware.JWTProtected())
	g.Get("", ctrl.List)
	g.Post("", middleware.JWTProtected(), ctrl.Create)
	g.Put("/:id", middleware.JWTProtected(), ctrl.Update)
	g.Patch("/:id", middleware.JWTProtected(), ctrl.UpdateStatus)

	gSetup := app.Group("/v1/setup-outlet-check", middleware.JWTProtected())
	gSetup.Get("", ctrl.SetupOutletCheck)
}

func (ctrl *OutletCodeController) List(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestID, headerAcceptLang)

	var filter entity.OutletCodeListFilter
	if err := c.QueryParser(&filter); err != nil {
		log.Println("OutletCodeController List QueryParser:", err.Error())
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Sort == "" {
		filter.Sort = "created_at:desc"
	}

	custId := ""
	if v := c.Locals("cust_id"); v != nil {
		custId = v.(string)
	}
	if filter.CustId != "" {
		custId = filter.CustId
	}

	data, total, lastPage, err := ctrl.OutletCodeService.List(filter, custId)
	if err != nil {
		log.Println("OutletCodeController List service err:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	resp.Setdata(data)
	resp.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})
	if len(data) == 0 {
		resp.Setmsg("No Data")
	}
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

// Create POST /v1/outlet_code - Create Outlet Code
func (ctrl *OutletCodeController) Create(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestID, headerAcceptLang)

	var body entity.CreateOutletCodeBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("OutletCodeController Create BodyParser:", err.Error())
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if strings.TrimSpace(body.SerialCode) == "" || strings.TrimSpace(body.LastSequenceNo) == "" {
		resp.Setmsg("serial_code, year_code and last_sequence_no are required")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	custId := ""
	if v := c.Locals("cust_id"); v != nil {
		custId = v.(string)
	}
	if custId == "" {
		if v := c.Locals("parent_cust_id"); v != nil {
			custId = v.(string)
		}
	}
	createdBy := ""
	if v := c.Locals("user_name"); v != nil {
		createdBy = v.(string)
	}
	if createdBy == "" && c.Locals("user_id") != nil {
		switch uid := c.Locals("user_id").(type) {
		case int:
			createdBy = strconv.Itoa(uid)
		case int64:
			createdBy = strconv.FormatInt(uid, 10)
		case float64:
			createdBy = strconv.FormatInt(int64(uid), 10)
		}
	}

	err := ctrl.OutletCodeService.Create(body, custId, createdBy)
	if err != nil {
		if errors.Is(err, entity.ErrOutletCodeDuplicate) {
			resp.Setmsg("Data already exists")
			return c.Status(fiber.StatusConflict).JSON(resp.GetRespPayload())
		}
		log.Println("OutletCodeController Create service err:", err.Error())
		resp.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

func (ctrl *OutletCodeController) Update(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestID, headerAcceptLang)

	id := c.Params("id")
	if strings.TrimSpace(id) == "" {
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	var body entity.UpdateOutletCodeBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("OutletCodeController Update BodyParser:", err.Error())
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if strings.TrimSpace(body.SerialCode) == "" {
		resp.Setmsg("serial_code is required")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	custId := ""
	if v := c.Locals("cust_id"); v != nil {
		custId = v.(string)
	}
	if custId == "" {
		if v := c.Locals("parent_cust_id"); v != nil {
			custId = v.(string)
		}
	}
	updatedBy := ""
	if v := c.Locals("user_name"); v != nil {
		updatedBy = v.(string)
	}
	if updatedBy == "" && c.Locals("user_id") != nil {
		switch uid := c.Locals("user_id").(type) {
		case int:
			updatedBy = strconv.Itoa(uid)
		case int64:
			updatedBy = strconv.FormatInt(uid, 10)
		case float64:
			updatedBy = strconv.FormatInt(int64(uid), 10)
		}
	}

	err := ctrl.OutletCodeService.Update(id, body, custId, updatedBy)
	if err != nil {
		if errors.Is(err, entity.ErrOutletCodeDuplicate) {
			resp.Setmsg("Data already exists")
			return c.Status(fiber.StatusConflict).JSON(resp.GetRespPayload())
		}
		log.Println("OutletCodeController Update service err:", err.Error())
		resp.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

func (ctrl *OutletCodeController) UpdateStatus(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestID, headerAcceptLang)

	pathID := c.Params("id")
	if strings.TrimSpace(pathID) == "" {
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	var body entity.UpdateOutletCodeStatusBody
	if err := c.BodyParser(&body); err != nil {
		log.Println("OutletCodeController UpdateStatus BodyParser:", err.Error())
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if strings.TrimSpace(body.Status) == "" || strings.TrimSpace(body.Id) == "" {
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if strings.TrimSpace(body.Id) != strings.TrimSpace(pathID) {
		resp.Setmsg("Bad request")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	custId := ""
	if v := c.Locals("cust_id"); v != nil {
		custId = v.(string)
	}
	if custId == "" {
		if v := c.Locals("parent_cust_id"); v != nil {
			custId = v.(string)
		}
	}
	updatedBy := ""
	if v := c.Locals("user_name"); v != nil {
		updatedBy = v.(string)
	}
	if updatedBy == "" && c.Locals("user_id") != nil {
		switch uid := c.Locals("user_id").(type) {
		case int:
			updatedBy = strconv.Itoa(uid)
		case int64:
			updatedBy = strconv.FormatInt(uid, 10)
		case float64:
			updatedBy = strconv.FormatInt(int64(uid), 10)
		}
	}

	if err := ctrl.OutletCodeService.UpdateStatus(pathID, body, custId, updatedBy); err != nil {
		log.Println("OutletCodeController UpdateStatus service err:", err.Error())
		resp.Setmsg("Failed to update data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setmsg("Data update successfully")
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}

// SetupOutletCheck GET /v1/setup-outlet-check - cek apakah distributor punya setup outlet code (status Active); jika ada maka FE hide field outlet code
func (ctrl *OutletCodeController) SetupOutletCheck(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	requestID := c.Locals("requestid").(string)
	resp := responsebuild.BuildResponse(requestID, headerAcceptLang)

	year := c.QueryInt("year", 0)
	if year == 0 {
		resp.Setmsg("year is required")
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}

	var status []string
	for _, v := range c.Context().QueryArgs().PeekMulti("status") {
		if s := strings.TrimSpace(string(v)); s != "" {
			status = append(status, s)
		}
	}
	for _, v := range c.Context().QueryArgs().PeekMulti("status[]") {
		if s := strings.TrimSpace(string(v)); s != "" {
			status = append(status, s)
		}
	}
	if single := strings.TrimSpace(c.Query("status")); single != "" && len(status) == 0 {
		status = []string{single}
	}
	if len(status) == 0 {
		status = []string{"Active"}
	}

	custId := ""
	if v := c.Locals("cust_id"); v != nil {
		custId = v.(string)
	}
	if custId == "" {
		resp.Setmsg("Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(resp.GetRespPayload())
	}
	parentCustId := ""
	if v := c.Locals("parent_cust_id"); v != nil {
		parentCustId = v.(string)
	}

	createdBy := ""
	if v := c.Locals("user_name"); v != nil {
		createdBy = v.(string)
	}
	data, err := ctrl.OutletCodeService.SetupOutletCheck(custId, parentCustId, year, status, createdBy)
	if err != nil {
		log.Println("OutletCodeController SetupOutletCheck err:", err.Error())
		resp.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	if data != nil {
		resp.Setdata(data)
	} else {
		resp.Setdata(map[string]interface{}{})
	}
	return c.Status(fiber.StatusOK).JSON(resp.GetRespPayload())
}
