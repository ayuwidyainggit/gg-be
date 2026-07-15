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
	"strconv"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type AutomaticReplenishmentProductController struct {
	AutomaticReplenishmentProductService service.AutomaticReplenishmentProductService
	validator                            *validation.Validate
}

func NewAutomaticReplenishmentProductController(automaticReplenishmentProductService service.AutomaticReplenishmentProductService, validator *validation.Validate) *AutomaticReplenishmentProductController {
	automaticReplenishmentProductValidator := validation.NewValidator()
	return &AutomaticReplenishmentProductController{
		AutomaticReplenishmentProductService: automaticReplenishmentProductService,
		validator:                            automaticReplenishmentProductValidator,
	}
}

func (controller *AutomaticReplenishmentProductController) Route(app *fiber.App) {
	qParamId := ":id"
	productReplenishmentsRouteV1 := app.Group("/v1/product-replenishments", middleware.JWTProtected())
	productReplenishmentsRouteV1.Get("", controller.List)
	productReplenishmentsRouteV1.Get("/template", controller.DownloadTemplate)
	productReplenishmentsRouteV1.Get("/export", controller.Export)
	productReplenishmentsRouteV1.Post("/import", controller.Import)
	productReplenishmentsRouteV1.Get("/"+qParamId, controller.Detail)
	productReplenishmentsRouteV1.Post("", controller.Create)
	productReplenishmentsRouteV1.Put("/"+qParamId, controller.Update)
	productReplenishmentsRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *AutomaticReplenishmentProductController) List(c *fiber.Ctx) error {
	var filter entity.AutomaticReplenishmentProductQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&filter); err != nil {
		log.Info("AutomaticReplenishmentProductController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	if filter.Sort == "" {
		filter.Sort = "created_at:desc"
	}

	errs := controller.validator.ValidateStruct(filter, headerAcceptLang)
	if errs != nil {
		log.Info("AutomaticReplenishmentProductController, List, ValidateStruct(filter), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	filter.CustId = c.Locals("cust_id").(string)
	filter.ParentCustId = c.Locals("parent_cust_id").(string)
	filter.JwtDistributorId = c.Locals("distributor_id").(int64)

	data, total, lastPage, err := controller.AutomaticReplenishmentProductService.List(filter, filter.CustId)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, List, service.List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	paging := map[string]interface{}{
		"total_record": total,
		"page_current": filter.Page,
		"page_limit":   filter.Limit,
		"page_total":   lastPage,
		"request_id":   c.Locals("requestid").(string),
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(paging)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *AutomaticReplenishmentProductController) Detail(c *fiber.Ctx) error {
	var params entity.DetailAutomaticReplenishmentProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("AutomaticReplenishmentProductController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("AutomaticReplenishmentProductController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustId = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.AutomaticReplenishmentProductService.Detail(params)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Detail, service.Detail, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "automatic replenishment product not found" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *AutomaticReplenishmentProductController) Export(c *fiber.Ctx) error {
	var filter entity.AutomaticReplenishmentProductQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&filter); err != nil {
		log.Info("AutomaticReplenishmentProductController, Export, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if filter.Sort == "" {
		filter.Sort = "created_at:desc"
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

	buffer, contentType, filename, err := controller.AutomaticReplenishmentProductService.Export(filter, filter.CustId)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Export, service.Export, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if _, err := c.Write(buffer.Bytes()); err != nil {
		log.Info("AutomaticReplenishmentProductController, Export, Write response, err:", err.Error())
		return c.Status(http.StatusInternalServerError).SendString("Failed to export file")
	}

	return nil
}

func (controller *AutomaticReplenishmentProductController) DownloadTemplate(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	format := c.Query("format", "xlsx")
	switch format {
	case "csv", "xls", "xlsx":
	default:
		responsePayload.Setmsg("Unsupported format. Use csv, xls, or xlsx")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	buffer, contentType, filename, err := controller.AutomaticReplenishmentProductService.DownloadTemplate(format)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, DownloadTemplate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *AutomaticReplenishmentProductController) Import(c *fiber.Ctx) error {
	var request entity.AutomaticReplenishmentProductImportRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("AutomaticReplenishmentProductController, Import, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("AutomaticReplenishmentProductController, Import, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	createdBy := c.Locals("user_id").(int64)

	data, err := controller.AutomaticReplenishmentProductService.Import(request, custId, createdBy)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Import, service.Import, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File imported successfully")
	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *AutomaticReplenishmentProductController) Create(c *fiber.Ctx) error {
	var request []*entity.CreateAutomaticReplenishmentProductRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("AutomaticReplenishmentProductController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	for i, item := range request {
		if item == nil {
			continue
		}

		if errs := controller.validator.ValidateStruct(item, headerAcceptLang); errs != nil {
			log.Info("Validation failed at index", i, ":", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	custId := c.Locals("cust_id").(string)
	createdBy := c.Locals("user_id").(int64)

	data, err := controller.AutomaticReplenishmentProductService.Create(request, custId, createdBy)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Create, service.Create, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *AutomaticReplenishmentProductController) Update(c *fiber.Ctx) error {
	var request entity.UpdateAutomaticReplenishmentProductRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Update, ParseInt:", err.Error())
		responsePayload.Setmsg("Invalid ID")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Info("AutomaticReplenishmentProductController, Update, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("AutomaticReplenishmentProductController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	updatedBy := c.Locals("user_id").(int64)

	err = controller.AutomaticReplenishmentProductService.Update(id, request, custId, updatedBy)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Update, service.Update, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "automatic replenishment product not found" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data updated successfully")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *AutomaticReplenishmentProductController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Delete, ParseInt:", err.Error())
		responsePayload.Setmsg("Invalid ID")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	deletedBy := c.Locals("user_id").(int64)

	err = controller.AutomaticReplenishmentProductService.Delete(custId, id, deletedBy)
	if err != nil {
		log.Info("AutomaticReplenishmentProductController, Delete, service.Delete, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "automatic replenishment product not found" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data deleted successfully")
	return c.JSON(responsePayload.GetRespPayload())
}
