package controller

import (
	"fmt"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ProductMappingController struct {
	ProductMappingService service.ProductMappingService
	validator             *validation.Validate
}

func NewProductMappingController(productMappingService service.ProductMappingService, validator *validation.Validate) *ProductMappingController {
	return &ProductMappingController{
		ProductMappingService: productMappingService,
		validator:             validation.NewProductMappingValidator(),
	}
}

func (controller *ProductMappingController) Route(app *fiber.App) {
	route := app.Group("/v1/product-mapping", middleware.JWTProtected())
	route.Get("", controller.List)
	route.Get("/export-template", controller.DownloadTemplate)
	route.Post("/import", controller.Import)
	route.Get("/:distributor_id", controller.Detail)
	route.Put("/:pro_id", controller.Update)
	route.Delete("/:pro_id", controller.Delete)
}

func (controller *ProductMappingController) ensurePrincipal(c *fiber.Ctx, responsePayload *responsebuild.DataRespReq) bool {
	distID, _ := c.Locals("distributor_id").(int64)
	if distID != 0 {
		responsePayload.Setmsg("Only principal user can access this resource")
		return false
	}
	return true
}

func (controller *ProductMappingController) List(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	var dataFilter entity.ProductMappingListQueryFilter
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductMappingController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.ProductMappingService.List(dataFilter)
	if err != nil {
		log.Info("ProductMappingController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}

	responsePayload.Setmsg("List product mapping has been processed successfully")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductMappingController) Detail(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	distributorID, err := strconv.ParseInt(c.Params("distributor_id"), 10, 64)
	if err != nil || distributorID <= 0 {
		responsePayload.Setmsg("invalid distributor_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	var dataFilter entity.ProductMappingDetailQueryFilter
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductMappingController, Detail, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	dataFilter.DistributorId = distributorID

	data, total, lastPage, err := controller.ProductMappingService.Detail(dataFilter)
	if err != nil {
		log.Info("ProductMappingController, Detail, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}

	responsePayload.Setmsg("Product mapping detail has been processed successfully")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductMappingController) Update(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	proID, err := strconv.ParseInt(c.Params("pro_id"), 10, 64)
	if err != nil || proID <= 0 {
		responsePayload.Setmsg("invalid pro_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	var req entity.ProductMappingUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Info("ProductMappingController, Update, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(req, headerAcceptLang)
	if errs != nil {
		log.Info("ProductMappingController, Update, ValidateStruct(req), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	principalCustID := c.Locals("cust_id").(string)
	updatedBy, _ := c.Locals("user_id").(int64)

	if err := controller.ProductMappingService.Update(proID, req, principalCustID, updatedBy); err != nil {
		log.Info("ProductMappingController, Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Product Mapping berhasil diperbarui")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductMappingController) Delete(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	proID, err := strconv.ParseInt(c.Params("pro_id"), 10, 64)
	if err != nil || proID <= 0 {
		responsePayload.Setmsg("invalid pro_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	principalCustID := c.Locals("cust_id").(string)
	deletedBy, _ := c.Locals("user_id").(int64)

	if err := controller.ProductMappingService.Delete(proID, principalCustID, deletedBy); err != nil {
		log.Info("ProductMappingController, Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Product Mapping berhasil dihapus")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductMappingController) DownloadTemplate(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	buffer, contentType, filename, err := controller.ProductMappingService.DownloadTemplate()
	if err != nil {
		log.Info("ProductMappingController, DownloadTemplate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *ProductMappingController) Import(c *fiber.Ctx) error {
	headerAcceptLang := getAcceptLang(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if !controller.ensurePrincipal(c, responsePayload) {
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	var req entity.ProductMappingImportRequest
	if err := c.BodyParser(&req); err != nil {
		log.Info("ProductMappingController, Import, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	fileURL := req.URL
	if fileURL == "" {
		fileURL = req.FileURL
	}
	if fileURL == "" {
		responsePayload.Setmsg("url is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	principalCustID := c.Locals("cust_id").(string)
	createdBy, _ := c.Locals("user_id").(int64)

	data, err := controller.ProductMappingService.Import(req, principalCustID, createdBy)
	if err != nil {
		log.Info("ProductMappingController, Import, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata(data)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg("Product mapping imported successfully")
	return c.JSON(responsePayload.GetRespPayload())
}

func getAcceptLang(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}
