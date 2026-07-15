package controller

import (
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type GrController struct {
	GrService service.GrService
	validator *validation.Validate
}

func NewGrController(grService service.GrService, validator *validation.Validate) *GrController {
	return &GrController{
		GrService: grService,
		validator: validator,
	}
}

func (controller *GrController) Route(app *fiber.App) {
	qParamId := ":gr_no"
	qParamInvoice := ":invoice_no"
	grRouteV1 := app.Group("/v1/goods-receipts", middleware.JWTProtected())
	grRouteV1.Get("/lookup-ap", controller.LookupGrAP)
	grRouteV1.Get("/download", controller.Download)
	grRouteV1.Get("/"+qParamId, controller.Detail)
	grRouteV1.Get("/invoice/"+qParamInvoice, controller.DetailInvoice)
	grRouteV1.Get("/suppliers/list", controller.ListSupplier)
	grRouteV1.Get("/warehouses/list", controller.ListWarehouse)
	grRouteV1.Get("/distributors/list", controller.ListDistributor)
	grRouteV1.Post("", controller.Create)
	grRouteV1.Get("", controller.List)

	// grRouteV1.Patch("/"+qParamId, controller.Update)
	// grRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *GrController) Detail(c *fiber.Ctx) error {
	var params entity.DetailGrParams
	var queryParam entity.GrDetailQuery

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("GrController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.QueryParser(&queryParam); err != nil {
		log.Error("GrController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.GrService.Detail(params.GrNo, custId, parentCustId, queryParam.IsAp)
	if err != nil {
		log.Error("GrController, Detail, err:", err.Error())
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
	responsePayload.Setmsg("Success")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *GrController) DetailInvoice(c *fiber.Ctx) error {
	var params entity.DetailGrInvoiceParams
	var queryParam entity.GrDetailQuery

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("GrController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.QueryParser(&queryParam); err != nil {
		log.Error("GrController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.GrService.DetailByInvoice(params.InvoiceNo, custId, parentCustId, queryParam.IsAp)
	if err != nil {
		log.Error("GrController, Detail, err:", err.Error())
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

func (controller *GrController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateGrBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("GrController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)
	distributorID := c.Locals("distributor_id").(int64)
	// log.Println("GrController, Create, CustId:", custId)

	request.CustID = custId
	request.ParentCustID = parentCustId
	request.DistributorID = distributorID
	request.CreatedBy = userId
	request.UpdatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// details.normal.pro_id validation ( must be unique )
	var detailsNormalProIds entity.DetailsNormalProductId
	for _, prod := range request.Details.Normal {
		detailsNormalProIds.ProductIds = append(detailsNormalProIds.ProductIds, entity.ProductId{
			Product: entity.Product{
				Id: prod.ProID,
			},
		})
	}
	errs = controller.validator.ValidateStruct(detailsNormalProIds, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Create, Detail Normal Product ID ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// details.promo.pro_id validation ( must be unique )
	var detailsPromoProIds entity.DetailsPromoProductId
	for _, prod := range request.Details.Promo {
		detailsPromoProIds.ProductIds = append(detailsPromoProIds.ProductIds, entity.ProductId{
			Product: entity.Product{
				Id: prod.ProID,
			},
		})
	}
	errs = controller.validator.ValidateStruct(detailsPromoProIds, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Create, Detail Promo Product ID ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.GrService.Store(request)
	if err != nil {
		log.Error("GrController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *GrController) List(c *fiber.Ctx) error {
	var dataFilter entity.GrQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("GrController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("BpprController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.GrService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("GrController, List, data, err:", err.Error())
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

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
	} else {
		responsePayload.Setmsg("Success")
	}

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *GrController) ListSupplier(c *fiber.Ctx) error {
	var dataFilter entity.GrSupplierQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("GrController, ListSupplier, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, ListSupplier, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.GrService.ListSupplier(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("GrController, ListSupplier, data, err:", err.Error())
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

// func (controller *GrController) Update(c *fiber.Ctx) error {
// 	var (
// 		params  entity.UpdateGrParams
// 		request entity.UpdateGrRequest
// 	)
// 	var headerAcceptLang string
// 	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
// 		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
// 	}
// 	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

// 	if err := c.ParamsParser(&params); err != nil {
// 		log.Error("GrController, Update, ParamsParser(params):", err.Error())
// 		responsePayload.Setmsg(err.Error())
// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
// 	}

// 	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
// 	if errs != nil {
// 		log.Error("GrController, Update, ValidateStruct(params), errs:", errs)
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	if err := c.BodyParser(&request); err != nil {
// 		log.Error("GrController, Update, BodyParser(request), err:", err.Error())
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
// 	}

// 	custId := c.Locals("cust_id").(string)
// 	parentCustId := c.Locals("parent_cust_id").(string)
// 	userId := c.Locals("user_id").(int64)
// 	// log.Println("BankController, Update, CustId:", custId)
// 	request.CustID = custId
// 	request.ParentCustID = parentCustId
// 	request.UpdatedBy = userId

// 	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
// 	if errs != nil {
// 		log.Error("GrController, Update, ValidateStruct(request), errs:", errs)
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		responsePayload.Seterrors(errs)
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	// details.normal.pro_id validation ( must be unique )
// 	var detailsNormalProIds entity.DetailsNormalProductId
// 	for _, prod := range request.Details.Normal {
// 		detailsNormalProIds.ProductIds = append(detailsNormalProIds.ProductIds, entity.ProductId{
// 			Product: entity.Product{
// 				Id: prod.ProID,
// 			},
// 		})
// 	}
// 	errs = controller.validator.ValidateStruct(detailsNormalProIds, headerAcceptLang)
// 	if errs != nil {
// 		log.Error("GrController, Create, Detail Normal Product ID ValidateStruct, errs:", errs)
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		responsePayload.Seterrors(errs)
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	// details.promo.pro_id validation ( must be unique )
// 	var detailsPromoProIds entity.DetailsPromoProductId
// 	for _, prod := range request.Details.Promo {
// 		detailsPromoProIds.ProductIds = append(detailsPromoProIds.ProductIds, entity.ProductId{
// 			Product: entity.Product{
// 				Id: prod.ProID,
// 			},
// 		})
// 	}
// 	errs = controller.validator.ValidateStruct(detailsPromoProIds, headerAcceptLang)
// 	if errs != nil {
// 		log.Error("GrController, Create, Detail Promo Product ID ValidateStruct, errs:", errs)
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		responsePayload.Seterrors(errs)
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	err := controller.GrService.Update(params.GrNo, request)
// 	if err != nil {
// 		log.Error("GrController, Update, Service.Update, err:", err.Error())
// 		responsePayload.Setmsg(err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}
// 	responsePayload.Setmsg("Updated Successfully")
// 	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
// }

// func (controller *GrController) Delete(c *fiber.Ctx) error {
// 	var headerAcceptLang string
// 	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
// 		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
// 	}
// 	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
// 	var params entity.DeleteGrParams
// 	if err := c.ParamsParser(&params); err != nil {
// 		log.Error("GrController, Update, ParamsParser(params):", err.Error())
// 		responsePayload.Setmsg(err.Error())
// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
// 	}
// 	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
// 	if errs != nil {
// 		log.Error("GrController, Create, ValidateStruct, errs:", errs)
// 		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
// 		responsePayload.Seterrors(errs)
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	custId := c.Locals("cust_id").(string)
// 	userId := c.Locals("user_id").(int64)
// 	// log.Println("VehicleController, Delete, CustId:", custId)

// 	err := controller.GrService.Delete(custId, params.GrNo, userId)
// 	if err != nil {
// 		log.Error("GrController, Delete, err:", err.Error())
// 		responsePayload.Setmsg(err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
// 	}

// 	responsePayload.Setmsg("Deleted Successfully")
// 	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
// }

func (controller *GrController) ListWarehouse(c *fiber.Ctx) error {
	var dataFilter entity.GrWarehouseQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("GrController, ListWarehouse, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, ListWarehouse, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.GrService.ListWarehouse(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("GrController, ListWarehouse, data, err:", err.Error())
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

func (controller *GrController) ListDistributor(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.GrService.ListDistributor(custId, parentCustId)
	if err != nil {
		log.Error("GrController, ListDistributor, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *GrController) Download(c *fiber.Ctx) error {
	var dataFilter entity.GrDownloadQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("GrController, Download, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("GrController, Download, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, err := controller.GrService.Download(dataFilter.GrNo, custId, parentCustId)
	if err != nil {
		log.Error("GrController, Download, err:", err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = constant.DATA_NOT_FOUND
		}

		// Check if it's the "in progress" error message
		if errMsg == "Processing time may vary by file size. Please check Download History to access the file" {
			responsePayload.Setmsg(errMsg)
			responsePayload.Setdata(nil)
			return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.SUCCESS)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *GrController) LookupGrAP(c *fiber.Ctx) error {
	var dataFilter entity.GrLookupQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("GrController Lookup, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("GrController Lookup, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.GrService.ListLookupGrAp(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("GrController Lookup, List, data, err:", err.Error())
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
