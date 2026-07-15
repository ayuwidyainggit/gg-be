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

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ProductAssignmentController struct {
	ProductAssignmentService service.ProductAssignmentService
	validator                *validation.Validate
}

func NewProductAssignmentController(productAssignmentService service.ProductAssignmentService, validator *validation.Validate) *ProductAssignmentController {
	productAssignmentValidator := validation.NewProductAssignmentValidator()
	return &ProductAssignmentController{
		ProductAssignmentService: productAssignmentService,
		validator:                productAssignmentValidator,
	}
}

func (controller *ProductAssignmentController) Route(app *fiber.App) {
	productAssignmentsRouteV1 := app.Group("/v1/product-assignments", middleware.JWTProtected())
	productAssignmentsRouteV1.Get("", controller.List)
	productAssignmentsRouteV1.Get("/template", controller.DownloadTemplate)
	productAssignmentsRouteV1.Get("/export", controller.Export)
	productAssignmentsRouteV1.Post("/import", controller.Import)
	productAssignmentsRouteV1.Post("/remove", controller.RemoveAssignment)
}

func (controller *ProductAssignmentController) List(c *fiber.Ctx) error {
	var dataFilter entity.ProductAssignmentQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductAssignmentController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.AssignmentType = parseStringSliceQuery(c.Context().QueryArgs(), "assignment_type", "assignment_type[]")

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	dataFilter.DistributorId = c.Locals("distributor_id").(int64)

	if dataFilter.DistributorId != 0 {
		responsePayload.Setmsg("Only principal user can access this resource")
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Info("ProductAssignmentController, List, ValidateStruct(dataFilter), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.ProductAssignmentService.List(dataFilter)
	if err != nil {
		log.Info("ProductAssignmentController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("List product assignment has been processed successfully")
	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductAssignmentController) Export(c *fiber.Ctx) error {
	var dataFilter entity.ProductAssignmentQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductAssignmentController, Export, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.AssignmentType = parseStringSliceQuery(c.Context().QueryArgs(), "assignment_type", "assignment_type[]")

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	dataFilter.DistributorId = c.Locals("distributor_id").(int64)

	if dataFilter.DistributorId != 0 {
		responsePayload.Setmsg("Only principal user can access this resource")
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	buffer, contentType, filename, err := controller.ProductAssignmentService.Export(dataFilter)
	if err != nil {
		log.Info("ProductAssignmentController, Export, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if _, err := c.Write(buffer.Bytes()); err != nil {
		log.Info("ProductAssignmentController, Export, Write response, err:", err.Error())
		return c.Status(http.StatusInternalServerError).SendString("Failed to export file")
	}

	return nil
}

func (controller *ProductAssignmentController) DownloadTemplate(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	buffer, contentType, filename, err := controller.ProductAssignmentService.DownloadTemplate()
	if err != nil {
		log.Info("ProductAssignmentController, DownloadTemplate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(buffer.Bytes())
}

func (controller *ProductAssignmentController) Import(c *fiber.Ctx) error {
	var req entity.ProductAssignmentImportRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&req); err != nil {
		log.Info("ProductAssignmentController, Import, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(req, headerAcceptLang)
	if errs != nil {
		log.Info("ProductAssignmentController, Import, ValidateStruct(req), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	createdBy, _ := c.Locals("user_id").(int64)
	distId, _ := c.Locals("distributor_id").(int64)

	if distId != 0 {
		responsePayload.Setmsg("Only principal user can access this resource")
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ProductAssignmentService.Import(req, custId, createdBy)
	if err != nil {
		log.Info("ProductAssignmentController, Import, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata([]entity.ProductAssignmentImportResponse{data})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata([]entity.ProductAssignmentImportResponse{data})
	responsePayload.Setmsg("Import product assignment has been processed successfully")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductAssignmentController) RemoveAssignment(c *fiber.Ctx) error {
	var req entity.ProductAssignmentImportRequest
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&req); err != nil {
		log.Info("ProductAssignmentController, RemoveAssignment, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(req, headerAcceptLang)
	if errs != nil {
		log.Info("ProductAssignmentController, RemoveAssignment, ValidateStruct(req), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userID, _ := c.Locals("user_id").(int64)
	distID, _ := c.Locals("distributor_id").(int64)

	if distID != 0 {
		responsePayload.Setmsg("Only principal user can access this resource")
		return c.Status(fiber.StatusForbidden).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ProductAssignmentService.RemoveAssignment(req, custId, userID)
	if err != nil {
		log.Info("ProductAssignmentController, RemoveAssignment, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata([]entity.ProductAssignmentImportResponse{data})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata([]entity.ProductAssignmentImportResponse{data})
	responsePayload.Setmsg("Remove assignment has been processed successfully")
	return c.JSON(responsePayload.GetRespPayload())
}
