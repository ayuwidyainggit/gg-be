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

	// "path/filepath"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductService service.ProductService
	validator      *validation.Validate
}

func NewProductController(productService service.ProductService, validator *validation.Validate) *ProductController {
	productValidator := validation.NewProductValidator()
	return &ProductController{
		ProductService: productService,
		validator:      productValidator,
	}
}

func (controller *ProductController) Route(app *fiber.App) {
	qParamId := ":pro_id"
	productsFileRouteV1 := app.Group("/v1/products-file", middleware.JWTProtected())
	productsFileRouteV1.Get("/export", controller.Export)
	productsFileRouteV1.Get("/export-instructions", controller.ExportImportInstructions)
	productsFileRouteV1.Get("/export-template", controller.ExportTemplate)
	productsFileRouteV1.Get("/export-template-update", controller.ExportTemplateUpdate)
	productsFileRouteV1.Post("/import", controller.Import)
	productsFileRouteV1.Post("/import-update", controller.ImportUpdate)
	productsRouteV1 := app.Group("/v1/products", middleware.JWTProtected())
	productsRouteV1.Get("/principals", controller.PrincipalList)
	productsRouteV1.Get("/categories", controller.CategoryList)
	productsRouteV1.Get("/brands", controller.BrandList)
	productsRouteV1.Post("/report", controller.Report)
	productsRouteV1.Get("/report", controller.Report)
	productsRouteV1.Get("/"+qParamId, controller.Detail)
	productsRouteV1.Get("", controller.List)
	productsRouteV1.Post("", controller.Create)
	productsRouteV1.Post("/bulk", controller.Bulk)
	productsRouteV1.Patch("/"+qParamId, controller.Update)
	productsRouteV1.Delete("/"+qParamId, controller.Delete)
	productsRouteV1.Delete("/", controller.DeleteMultiple)
}

func (controller *ProductController) Detail(c *fiber.Ctx) error {
	var params entity.DetailProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("ProductController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.DistributorID = c.Locals("distributor_id").(int64)
	// log.Debug(params.DistributorID)

	data, err := controller.ProductService.Detail(params)
	if err != nil {
		log.Info("ProductController, Detail, FindOneByProductId, err:", err.Error())
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
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	tokenDistributorID := c.Locals("distributor_id").(int64)
	if tokenDistributorID > 0 {
		dataFilter.JwtDistributorId = tokenDistributorID
		dataFilter.DistributorID = tokenDistributorID
	}

	switch dataFilter.Mode {
	case "search":
		data, total, lastPage, err = controller.ProductService.SearchList(dataFilter, dataFilter.CustId)
		if err != nil {
			log.Info("ProductController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	case "lookup":
		data, total, lastPage, err = controller.ProductService.LookupList(dataFilter, dataFilter.CustId)
		if err != nil {
			log.Info("ProductController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	case "lookup_dist_price":
		data, total, lastPage, err = controller.ProductService.LookupDistPrice(dataFilter)
		if err != nil {
			log.Info("ProductController, Lookup Dist Price, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.ProductService.List(dataFilter, dataFilter.CustId)
		if err != nil {
			log.Info("ProductController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
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

func (controller *ProductController) Create(c *fiber.Ctx) error {
	var request entity.CreateProductBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("ProductController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	distributorId := c.Locals("distributor_id").(int64)
	// log.Info("ProductController, Create, CustId:", custId)

	request.CustId = custId
	request.CreatedBy = userId
	request.DistributorId = distributorId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.ProductService.Store(request)
	if err != nil {
		log.Info("ProductController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) Bulk(c *fiber.Ctx) error {
	var request entity.BulkProductBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("ProductController, Bulk, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Info("ProductController, Bulk, ValidateRequest, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for index := range request.Products {
		request.Products[index].CustId = custId
		request.Products[index].CreatedBy = userId

		if errs := controller.validator.ValidateStruct(request.Products[index], headerAcceptLang); errs != nil {
			log.Info("ProductController, Bulk, ValidateStruct line "+fmt.Sprint(index+2)+", errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	if _, err := controller.ProductService.BulkStore(request); err != nil {
		log.Info("ProductController, Bulk, Bulk Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateProductParams
		request entity.UpdateProductRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("ProductController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Info("ProductController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("ProductController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ProductService.Update(params.ProductId, request)
	if err != nil {
		log.Info("ProductController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("ProductController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("ProductController, Delete, CustId:", custId)

	err := controller.ProductService.Delete(custId, params.ProductId, userId)
	if err != nil {
		log.Info("ProductController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) DeleteMultiple(c *fiber.Ctx) error {
	var request entity.DeleteMultipleProductBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Info("ProductController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("ProductController, DeleteMultiple, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	err := controller.ProductService.DeleteMultiple(custId, request.ProductId, userId)
	if err != nil {
		log.Info("ProductController, DeleteMultiple, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())

	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) PrincipalList(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductPrincipalQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductController, PrincipalList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Info("ProductController, List, CustId:", custId)

	data, total, lastPage, err = controller.ProductService.PrincipalList(dataFilter, dataFilter.CustId)
	if err != nil {
		log.Info("ProductController, PrincipalList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   0,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) CategoryList(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductCategoryQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductController, CategoryList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Info("ProductController, List, CustId:", custId)

	data, total, lastPage, err = controller.ProductService.CategoryList(dataFilter, dataFilter.CustId)
	if err != nil {
		log.Info("ProductController, CategoryList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   0,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) BrandList(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductBrandQueryFilter
		data       interface{}
		total      int
		lastPage   int
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductController, BrandList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Info("ProductController, List, CustId:", custId)

	data, total, lastPage, err = controller.ProductService.BrandList(dataFilter, dataFilter.CustId)
	if err != nil {
		log.Info("ProductController, BrandList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   0,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) Report(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductReportQueryFilter
		data       interface{}
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// Parse cust_id[] from query
	custIDsRaw := c.Request().URI().QueryArgs().PeekMulti("cust_id[]")
	if len(custIDsRaw) == 0 {
		responsePayload.Setmsg("cust_id[] is required")
		responsePayload.Seterrors(map[string]string{"cust_id": "missing or empty"})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Trim and validate each cust_id
	dataFilter.CustIDs = make([]string, 0, len(custIDsRaw))
	for _, id := range custIDsRaw {
		trimmed := strings.TrimSpace(string(id))
		if trimmed == "" {
			responsePayload.Setmsg("cust_id[] contains blank entry")
			responsePayload.Seterrors(map[string]string{"cust_id": "blank entry not allowed"})
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		dataFilter.CustIDs = append(dataFilter.CustIDs, trimmed)
	}

	// Parse optional query parameters
	dataFilter.Query = c.Query("q", "")
	dataFilter.Page = c.QueryInt("page", 1)
	dataFilter.Limit = c.QueryInt("limit", 20)
	dataFilter.SortBy = c.Query("sort_by", "pro_name")
	dataFilter.SortOrder = c.Query("sort_order", "asc")

	// Normalize page
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}

	// Validate sort_by
	validSortCols := map[string]bool{"pro_name": true, "pro_code": true, "type": true, "pro_id": true}
	if !validSortCols[dataFilter.SortBy] {
		responsePayload.Setmsg("sort_by must be one of: pro_name, pro_code, type, pro_id")
		responsePayload.Seterrors(map[string]string{"sort_by": "invalid value"})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Normalize sort_order
	dataFilter.SortOrder = strings.ToLower(dataFilter.SortOrder)
	if dataFilter.SortOrder != "asc" && dataFilter.SortOrder != "desc" {
		responsePayload.Setmsg("sort_order must be asc or desc")
		responsePayload.Seterrors(map[string]string{"sort_order": "invalid value"})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Call service
	var total, lastPage int
	data, total, lastPage, err = controller.ProductService.ReportList(dataFilter)
	if err != nil {
		log.Info("ProductController, Report, ReportList, err:", err.Error())
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

func (controller *ProductController) ExportImportInstructions(c *fiber.Ctx) error {

	// Ambil query param format, default ke xlsx
	format := c.Query("format", "xlsx")

	// Validasi format
	switch format {
	case "csv", "xls", "xlsx":
		// valid
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tidak didukung. Gunakan csv, xls, atau xlsx",
		})
	}

	// Panggil service
	buffer, contentType, filename, err := controller.ProductService.ExportImportInstructions(format)
	if err != nil {
		log.Info("ProductController, ExportImportInstructions, error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set header untuk download
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(buffer.Bytes())
}

func (controller *ProductController) Export(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.ProductQueryFilter
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("ProductController, Export, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	tokenDistributorID := c.Locals("distributor_id").(int64)
	if tokenDistributorID > 0 {
		dataFilter.JwtDistributorId = tokenDistributorID
		dataFilter.DistributorID = tokenDistributorID
	}

	// Panggil service yang sekarang mengembalikan buffer, content-type, dan nama file
	buffer, contentType, filename, err := controller.ProductService.Export(dataFilter)
	if err != nil {
		log.Info("ProductController, Export, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Set header HTTP secara dinamis berdasarkan hasil dari service
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if _, err := c.Write(buffer.Bytes()); err != nil {
		log.Info("ProductController, Export, Write response, err:", err.Error())
		return c.Status(http.StatusInternalServerError).SendString("Failed to export file")
	}

	return nil
}

func (controller *ProductController) ExportTemplate(c *fiber.Ctx) error {
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

	// Panggil service untuk export template
	buffer, contentType, filename, err := controller.ProductService.ExportTemplate(format)
	if err != nil {
		log.Info("ProductController, ExportTemplate, error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set header untuk download file
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(buffer.Bytes())
}

func (controller *ProductController) ExportTemplateUpdate(c *fiber.Ctx) error {
	// Ambil format (default xlsx)
	format := c.Query("format", "xlsx")
	CustId := c.Locals("cust_id").(string)

	// Ambil fields yang dipilih user
	fieldsParam := c.Query("fields", "")
	if fieldsParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fields harus diisi, contoh: ?fields=brand,category",
		})
	}
	fields := strings.Split(fieldsParam, ",")

	// Panggil service
	buffer, contentType, filename, err := controller.ProductService.ExportTemplateUpdate(CustId, format, fields)
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

func (controller *ProductController) Import(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// 1. Dapatkan format dari query parameter
	format := c.Query("format")
	if format == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// 2. Dapatkan file dari request
	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		responsePayload.Setmsg("File upload with key 'file_upload' is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// 3. Buka file
	file, err := fileHeader.Open()
	if err != nil {
		responsePayload.Setmsg("Failed to process uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	// 4. Siapkan request untuk service
	importReq := entity.ImportProductRequest{
		CustId:        c.Locals("cust_id").(string),
		ParentCustId:  c.Locals("parent_cust_id").(string),
		DistributorId: c.Locals("distributor_id").(int64),
		CreatedBy:     c.Locals("user_id").(int64),
		File:          file,
		FileName:      fileHeader.Filename,
	}

	// 5. Panggil service yang sesuai
	switch format {
	case "csv":
		err = controller.ProductService.ImportProductCSV4(importReq)
	case "xlsx":
		err = controller.ProductService.ImportProductXLSX4(importReq)
	case "xls":
		err = controller.ProductService.ImportProductXLSX4(importReq)
	default:
		err = fmt.Errorf("unsupported format '%s'", format)
	}

	// 6. Tangani hasil
	if err != nil {
		log.Info("Import failed: %v", err)  // Log error untuk debugging
		responsePayload.Setmsg(err.Error()) // Kirim pesan error ke klien
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ProductController) ImportUpdate(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// 1. Dapatkan format dari query parameter
	format := c.Query("format")
	if format == "" {
		responsePayload.Setmsg("Query parameter 'format' (csv/xlsx) is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// 2. Dapatkan file dari request
	fileHeader, err := c.FormFile("file_upload")
	if err != nil {
		responsePayload.Setmsg("File upload with key 'file_upload' is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// 3. Buka file
	file, err := fileHeader.Open()
	if err != nil {
		responsePayload.Setmsg("Failed to process uploaded file")
		return c.Status(http.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()

	// 4. Siapkan request untuk service
	importReq := entity.ImportProductRequest{
		CustId:       c.Locals("cust_id").(string),
		ParentCustId: c.Locals("parent_cust_id").(string),
		CreatedBy:    c.Locals("user_id").(int64),
		File:         file,
		FileName:     fileHeader.Filename,
	}

	// 5. Panggil service yang sesuai
	switch format {
	case "csv":
		err = controller.ProductService.ImportUpdateCSV(importReq)
	case "xlsx":
		err = controller.ProductService.ImportUpdateXLSX(importReq)
	case "xls":
		err = controller.ProductService.ImportUpdateXLSX(importReq)
	default:
		err = fmt.Errorf("unsupported format '%s'", format)
	}

	// 6. Tangani hasil
	if err != nil {
		log.Info("Import failed: %v", err)  // Log error untuk debugging
		responsePayload.Setmsg(err.Error()) // Kirim pesan error ke klien
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("File import accepted. Processing in background, check import history for success or failure details.")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
