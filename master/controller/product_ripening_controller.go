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
	"github.com/gofiber/fiber/v2/log"
)

type ProductRipeningController struct {
	service   service.ProductRipeningService
	validator *validation.Validate
}

func NewProductRipeningController(service service.ProductRipeningService, validator *validation.Validate) *ProductRipeningController {
	return &ProductRipeningController{service: service, validator: validator}
}

func (controller *ProductRipeningController) Route(app *fiber.App) {
	g := app.Group("/v1/product-ripenings", middleware.JWTProtected())
	g.Get("", controller.List)
	g.Get("/template", controller.DownloadTemplate)
	g.Get("/export", controller.Export)
	g.Post("/import", controller.Import)
	g.Get("/:distributor_id/:per_year/:week_id", controller.Detail)
	g.Put("/:distributor_id/:per_year/:week_id", controller.Update)
	// g.Get("/history", controller.ListHistory)
	// g.Get("/history/:id", controller.HistoryDetail)
}

func (controller *ProductRipeningController) List(c *fiber.Ctx) error {
	var filter entity.ProductRipeningQueryFilter
	responsePayload, headerAcceptLang := ripeningResponsePayload(c)
	if err := c.QueryParser(&filter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Sort == "" {
		filter.Sort = "per_year:desc,per_id:desc,week_id:desc"
	}
	filter.CustId = c.Locals("cust_id").(string)
	filter.ParentCustId = c.Locals("parent_cust_id").(string)
	filter.JwtDistributorId = c.Locals("distributor_id").(int64)

	if errs := controller.validator.ValidateStruct(filter, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.service.List(filter, c.Locals("user_id").(int64))
	if err != nil {
		log.Info("ProductRipeningController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(map[string]interface{}{
		"total_record": total,
		"page_current": filter.Page,
		"page_limit":   filter.Limit,
		"page_total":   lastPage,
		"request_id":   c.Locals("requestid").(string),
	})
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductRipeningController) Detail(c *fiber.Ctx) error {
	var params entity.ProductRipeningDetailParams
	responsePayload, headerAcceptLang := ripeningResponsePayload(c)
	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.service.Detail(params, c.Locals("cust_id").(string), c.Locals("parent_cust_id").(string), c.Locals("user_id").(int64))
	if err != nil {
		log.Info("ProductRipeningController, Detail, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductRipeningController) Update(c *fiber.Ctx) error {
	var params entity.ProductRipeningDetailParams
	var req entity.ProductRipeningUpdateRequest
	responsePayload, headerAcceptLang := ripeningResponsePayload(c)
	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if err := c.BodyParser(&req); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(req, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.service.Update(params, req, c.Locals("cust_id").(string), c.Locals("parent_cust_id").(string), c.Locals("user_id").(int64))
	if err != nil {
		log.Info("ProductRipeningController, Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Product ripening updated successfully")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductRipeningController) DownloadTemplate(c *fiber.Ctx) error {
	responsePayload, _ := ripeningResponsePayload(c)

	format := c.Query("format", "xlsx")
	switch format {
	case "csv", "xls", "xlsx":
	default:
		responsePayload.Setmsg("Unsupported format. Use csv, xls, or xlsx")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	buffer, contentType, filename, err := controller.service.DownloadTemplate(format, c.Locals("cust_id").(string), c.Locals("parent_cust_id").(string), c.Locals("user_id").(int64))
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *ProductRipeningController) Export(c *fiber.Ctx) error {
	var filter entity.ProductRipeningQueryFilter
	responsePayload, _ := ripeningResponsePayload(c)
	if err := c.QueryParser(&filter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if filter.Sort == "" {
		filter.Sort = "per_year:desc,per_id:desc,week_id:desc"
	}
	if filter.Format == "" {
		filter.Format = "xlsx"
	}
	switch filter.Format {
	case "csv", "xls", "xlsx":
	default:
		responsePayload.Setmsg("Unsupported format. Use csv, xls, or xlsx")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	filter.CustId = c.Locals("cust_id").(string)
	filter.ParentCustId = c.Locals("parent_cust_id").(string)
	filter.JwtDistributorId = c.Locals("distributor_id").(int64)

	buffer, contentType, filename, err := controller.service.Export(filter, c.Locals("user_id").(int64))
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	if _, err := c.Write(buffer.Bytes()); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to export file")
	}
	return nil
}

func (controller *ProductRipeningController) Import(c *fiber.Ctx) error {
	var req entity.ProductRipeningImportRequest
	responsePayload, headerAcceptLang := ripeningResponsePayload(c)
	if err := c.BodyParser(&req); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if errs := controller.validator.ValidateStruct(req, headerAcceptLang); errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.service.Import(req, c.Locals("cust_id").(string), c.Locals("parent_cust_id").(string), c.Locals("user_id").(int64))
	if err != nil {
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata(data)
		return c.Status(ripeningStatusCode(err)).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("File imported successfully")
	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func ripeningResponsePayload(c *fiber.Ctx) (*responsebuild.DataRespReq, string) {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang), headerAcceptLang
}

func ripeningStatusCode(err error) int {
	if err == nil {
		return fiber.StatusOK
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "only available for users assigned as distributor pic"):
		return fiber.StatusForbidden
	case strings.Contains(message, "not found"):
		return fiber.StatusNotFound
	default:
		return fiber.StatusBadRequest
	}
}
