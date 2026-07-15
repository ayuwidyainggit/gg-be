package controller

import (
	"fmt"
	"inventory/entity"
	"inventory/pkg/config"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterSAPReplenishmentRoutes(app *fiber.App, sapCfg config.SAPReplenishmentStatusConfig, c *ReplenishmentController) {
	grp := app.Group("/v1/replenishment-order/sap", middleware.SAPReplenishmentStatusProtected(sapCfg))
	grp.Get("", c.SAPGetReplenishmentExport)
	grp.Post("/status", c.SAPPostReplenishmentStatus)
}

func sapFallbackRequestID() string {
	return fmt.Sprintf("REQ-%s-ERR", time.Now().UTC().Format("20060102"))
}

func (controller *ReplenishmentController) SAPGetReplenishmentExport(c *fiber.Ctx) error {
	var query entity.SAPReplExportQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": constant.INVALID_JSON_BODY,
		})
	}

	lang := "en"
	if hl := strings.TrimSpace(c.Get(constant.HEADER_ACCEPT_LANG)); hl != "" {
		lang = hl
	}
	if errs := controller.validator.ValidateStruct(query, lang); len(errs) > 0 {
		var fe []entity.SAPReplFieldError
		for _, e := range errs {
			key := ""
			msg := ""
			if k, ok := e["key"].(string); ok {
				key = k
			}
			if m, ok := e["message"].(string); ok {
				msg = m
			}
			fe = append(fe, entity.SAPReplFieldError{Field: key, Message: msg})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": constant.ERR_VALIDATION,
			"errors":  fe,
		})
	}

	data, err := controller.ReplenishmentService.SAPGetReplenishmentExport(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	if data == nil {
		data = []entity.SAPReplExportItem{}
	}
	return c.Status(fiber.StatusOK).JSON(data)
}

func (controller *ReplenishmentController) SAPPostReplenishmentStatus(c *fiber.Ctx) error {
	var req entity.SAPReplStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.SAPReplStatusResponse{
			RequestID: sapFallbackRequestID(),
			Status:    "error",
			Message:   constant.INVALID_JSON_BODY,
			Errors: []entity.SAPReplStatusReplErrWrap{
				{Errors: []entity.SAPReplFieldError{{Field: "body", Message: err.Error()}}},
			},
		})
	}

	lang := "en"
	if hl := strings.TrimSpace(c.Get(constant.HEADER_ACCEPT_LANG)); hl != "" {
		lang = hl
	}
	if errs := controller.validator.ValidateStruct(req, lang); len(errs) > 0 {
		var fe []entity.SAPReplFieldError
		for _, e := range errs {
			key := ""
			msg := ""
			if k, ok := e["key"].(string); ok {
				key = k
			}
			if m, ok := e["message"].(string); ok {
				msg = m
			}
			fe = append(fe, entity.SAPReplFieldError{Field: key, Message: msg})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entity.SAPReplStatusResponse{
			RequestID: sapFallbackRequestID(),
			Status:    "error",
			Message:   constant.ERR_VALIDATION,
			Errors: []entity.SAPReplStatusReplErrWrap{
				{ReplenishmentNo: "", Errors: fe},
			},
		})
	}

	resp := controller.ReplenishmentService.SAPUpdateReplenishmentStatus(
		req,
		middleware.SAPReplenishmentStatusLocals(c),
	)

	code := fiber.StatusOK
	if resp.Status == "error" {
		code = fiber.StatusUnprocessableEntity
	}
	return c.Status(code).JSON(resp)
}
