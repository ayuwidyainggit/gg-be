package controller

import (
	"fmt"
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

type HistoryController struct {
	svc       service.HistoryService
	validator *validation.Validate
}

func NewHistoryController(svc service.HistoryService, validator *validation.Validate) *HistoryController {
	return &HistoryController{svc: svc, validator: validator}
}

func (h *HistoryController) Route(app *fiber.App) {
	grp := app.Group("/v1/history", middleware.JWTProtected())
	grp.Get("", h.List)
	// Support both path-param and query-param based downloads
	grp.Get("/download", h.Download)
	grp.Get("/download/:id", h.Download)
	grp.Get("/:id", h.Detail)
	grp.Post("/reupload/file", h.ReuploadFile)
}

func (h *HistoryController) List(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	resp := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	uploadType := c.Query("upload_type")
	search := c.Query("search")
	custId := c.Locals("cust_id").(string)
	// Allow searching by explicit file_name param as an alias for search
	if strings.TrimSpace(search) == "" {
		if fn := strings.TrimSpace(c.Query("file_name")); fn != "" {
			search = fn
		}
	}
	var ut *string
	if strings.TrimSpace(uploadType) != "" {
		ut = &uploadType
	}
	var s *string
	if strings.TrimSpace(search) != "" {
		s = &search
	}
	rows, total, last, err := h.svc.List(ut, s, page, limit, custId)
	if err != nil {
		resp.Setmsg(err.Error())
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setdata(rows)
	resp.Setpaging(entity.Pagination{TotalRecord: total, PageCurrent: page, PageLimit: limit, PageTotal: last})
	return c.JSON(resp.GetRespPayload())
}

func (h *HistoryController) Download(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	resp := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// Allow both path param /download/:id and query ?id= / ?history_id=
	idParam := strings.TrimSpace(c.Params("id"))
	if idParam == "" {
		idParam = strings.TrimSpace(c.Query("id"))
		if idParam == "" {
			idParam = strings.TrimSpace(c.Query("history_id"))
		}
	}
	if idParam == "" {
		resp.Setmsg("history_id is required")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	var historyId int64
	if _, err := fmt.Sscan(idParam, &historyId); err != nil || historyId <= 0 {
		resp.Setmsg("invalid history_id")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	format := c.Query("format", "xlsx")
	custId := c.Locals("cust_id").(string)
	buf, ct, fn, err := h.svc.Download(custId, historyId, format)
	if err != nil {
		resp.Setmsg(err.Error())
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	c.Set("Content-Type", ct)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fn))
	return c.Send(buf.Bytes())
}

func (h *HistoryController) Detail(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	resp := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	idParam := strings.TrimSpace(c.Params("id"))
	if idParam == "" {
		resp.Setmsg("history_id is required")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	var historyId int64
	if _, err := fmt.Sscan(idParam, &historyId); err != nil || historyId <= 0 {
		resp.Setmsg("invalid history_id")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	row, err := h.svc.GetDetail(historyId)
	if err != nil {
		resp.Setmsg(err.Error())
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setdata(row)
	return c.JSON(resp.GetRespPayload())
}

func (h *HistoryController) ReuploadFile(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	resp := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := c.Query("format", "xlsx")
	idParam := c.Query("id")
	createdBy := c.Locals("user_id").(int64)
	var historyId int64
	if _, err := fmt.Sscan(idParam, &historyId); err != nil || historyId == 0 {
		resp.Setmsg("id (history_id) is required")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		resp.Setmsg("file_upload is required")
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	f, err := fileHeader.Open()
	if err != nil {
		resp.Setmsg("failed to open file")
		return c.Status(http.StatusInternalServerError).JSON(resp.GetRespPayload())
	}
	defer f.Close()
	custId := c.Locals("cust_id").(string)
	req := entity.ImportRequest{File: f, CustId: custId, Filename: fileHeader.Filename, Format: format, UserId: createdBy}
	if err := h.svc.ReuploadFile(custId, historyId, format, req); err != nil {
		resp.Setmsg("reupload failed")
		resp.Setdata(map[string]interface{}{"error": err.Error()})
		return c.Status(http.StatusBadRequest).JSON(resp.GetRespPayload())
	}
	resp.Setmsg("reupload successfull. Processing in background, check import history for success or failure details.")
	return c.JSON(resp.GetRespPayload())
}
