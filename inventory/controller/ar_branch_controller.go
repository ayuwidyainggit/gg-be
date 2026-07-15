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

type ArBranchController struct {
	ArBranchService service.ArBranchService
	validator       *validation.Validate
}

func NewArBranchController(arBranchService service.ArBranchService, validator *validation.Validate) *ArBranchController {
	return &ArBranchController{
		ArBranchService: arBranchService,
		validator:       validator,
	}
}

func (controller *ArBranchController) Route(app *fiber.App) {
	qParamId := ":gr_branch_no"
	// qParamOrderBookingId := ":order_booking_id"
	// qParamInvoice := ":invoice_no"
	arBranchRouteV1 := app.Group("/v1/ar-branch", middleware.JWTProtected())
	arBranchRouteV1.Get("/"+qParamId, controller.Detail)
	// arBranchRouteV1.Get("/invoice/"+qParamInvoice, controller.DetailInvoice)
	// arBranchRouteV1.Get("/suppliers/list", controller.ListSupplier)
	// arBranchRouteV1.Get("/warehouses/list", controller.ListWarehouse)
	// arBranchRouteV1.Get("/order-bookings/list", controller.OrderBookingList)
	// arBranchRouteV1.Get("/order-bookings/"+qParamOrderBookingId, controller.OrderBookingDetail)
	arBranchRouteV1.Get("", controller.List)
	arBranchRouteV1.Get("/filter/distributors", controller.DistributorsFilter)
	arBranchRouteV1.Get("/filter/suppliers", controller.SuppliersFilter)
	arBranchRouteV1.Post("payments/"+qParamId, controller.CreateArBranchPayment)
	// arBranchRouteV1.Patch("/status", controller.UpdateStatus)
	// arBranchRouteV1.Patch("/print", controller.Print)

	// arBranchRouteV1.Patch("/"+qParamId, controller.Update)
	// arBranchRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *ArBranchController) Detail(c *fiber.Ctx) error {
	var params entity.DetailArBranchParams
	var queryParam entity.ArBranchDetailQuery

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ArBranchController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.QueryParser(&queryParam); err != nil {
		log.Error("ArBranchController, Detail, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	// data, err := controller.ArBranchService.Detail(params.ArBranchNo, custId, parentCustId, queryParam.IsAp)
	data, err := controller.ArBranchService.Detail(params.GrBranchNo, custId, parentCustId, queryParam)
	if err != nil {
		log.Error("ArBranchController, Detail, err:", err.Error())
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

/*
	func (controller *ArBranchController) DetailInvoice(c *fiber.Ctx) error {
		var params entity.DetailArBranchInvoiceParams
		var queryParam entity.ArBranchDetailInvoiceQuery

		var headerAcceptLang string
		if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
			headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
		}
		responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

		if err := c.ParamsParser(&params); err != nil {
			log.Error("ArBranchController, Detail, ParamsParser:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}

		if err := c.QueryParser(&queryParam); err != nil {
			log.Error("ArBranchController, List, query parser filter:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}

		errs := controller.validator.ValidateStruct(params, headerAcceptLang)
		if errs != nil {
			log.Error("ArBranchController, Detail, ValidateStruct(params), errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		custId := c.Locals("cust_id").(string)
		parentCustId := c.Locals("parent_cust_id").(string)
		// log.Println("OutletController, Detail, CustId:", custId)

		data, err := controller.ArBranchService.DetailByInvoice(params.InvoiceNo, custId, parentCustId, queryParam.IsAp)
		if err != nil {
			log.Error("ArBranchController, Detail, err:", err.Error())
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

	func (controller *ArBranchController) Create(c *fiber.Ctx) error {
		var headerAcceptLang string
		if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
			headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
		}
		responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

		var request entity.CreateArBranchBody
		if err := c.BodyParser(&request); err != nil {
			log.Error("ArBranchController, Create, BodyParser:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}

		custId := c.Locals("cust_id").(string)
		userId := c.Locals("user_id").(int64)
		parentCustId := c.Locals("parent_cust_id").(string)

		// log.Println("ArBranchController, Create, CustId:", custId)

		request.CustID = custId
		request.ParentCustID = parentCustId
		request.CreatedBy = userId
		request.UpdatedBy = userId

		errs := controller.validator.ValidateStruct(request, headerAcceptLang)
		if errs != nil {
			log.Error("ArBranchController, Create, ValidateStruct, errs:", errs)
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
			log.Error("ArBranchController, Create, Detail Normal Product ID ValidateStruct, errs:", errs)
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
			log.Error("ArBranchController, Create, Detail Promo Product ID ValidateStruct, errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		_, err := controller.ArBranchService.Store(request)
		if err != nil {
			log.Error("ArBranchController, Create, Store, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Created Successfully")
		return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
	}
*/
func (controller *ArBranchController) List(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, List, query parser filter:", err.Error())
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

	data, total, lastPage, err := controller.ArBranchService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, List, data, err:", err.Error())
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

/*
func (controller *ArBranchController) ListSupplier(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchSupplierQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, ListSupplier, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, ListSupplier, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.ArBranchService.ListSupplier(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, ListSupplier, data, err:", err.Error())
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

func (controller *ArBranchController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateArBranchParams
		request entity.UpdateArBranchRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ArBranchController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("ArBranchController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustID = custId
	request.ParentCustID = parentCustId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, Update, ValidateStruct(request), errs:", errs)
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
		log.Error("ArBranchController, Create, Detail Normal Product ID ValidateStruct, errs:", errs)
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
		log.Error("ArBranchController, Create, Detail Promo Product ID ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ArBranchService.Update(params.ArBranchNo, request)
	if err != nil {
		log.Error("ArBranchController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ArBranchController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DeleteArBranchParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ArBranchController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.ArBranchService.Delete(custId, params.ArBranchNo, userId)
	if err != nil {
		log.Error("ArBranchController, Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ArBranchController) ListWarehouse(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchWarehouseQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, ListWarehouse, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, ListWarehouse, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.ArBranchService.ListWarehouse(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, ListWarehouse, data, err:", err.Error())
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

func (controller *ArBranchController) OrderBookingList(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchOrderBookingListQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, OrderBookingList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, OrderBookingList, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.ArBranchService.OrderBookingList(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, OrderBookingList, data, err:", err.Error())
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

func (controller *ArBranchController) OrderBookingDetail(c *fiber.Ctx) error {
	var params entity.ArBranchOrderBookingDetailParams

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ArBranchController, OrderBookingDetail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, OrderBookingDetail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.ArBranchService.OrderBookingDetail(params.OrderBookingId, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, OrderBookingDetail, err:", err.Error())
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

func (controller *ArBranchController) UpdateStatus(c *fiber.Ctx) error {
	var request entity.ArBranchBulkUpdateDataStatus
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ArBranchController, BulkUpdate, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("ArBranchController, BulkUpdate, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for index := range request.ArBranches {
		request.ArBranches[index].UpdatedBy = userId
		// request.ArBranches[index].CustId = custId

		if errs := controller.validator.ValidateStruct(request.ArBranches[index], headerAcceptLang); errs != nil {
			log.Error("ArBranchController, BulkUpdate, ValidateStruct Order with ArBranchNo "+fmt.Sprint(request.ArBranches[index].ArBranchNo)+", errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	if err := controller.ArBranchService.BulkUpdateStatus(request, custId, parentCustId); err != nil {
		log.Error("ArBranchController, Update, Service.BulkUpdate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Updated Status Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ArBranchController) Print(c *fiber.Ctx) error {
	var request entity.ArBranchBulkPrint
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("ArBranchController, BulkPrint, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("ArBranchController, BulkPrint, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for index := range request.ArBranches {
		if errs := controller.validator.ValidateStruct(request.ArBranches[index], headerAcceptLang); errs != nil {
			log.Error("ArBranchController, BulkPrint, ValidateStruct Order with ArBranchNo "+fmt.Sprint(request.ArBranches[index].ArBranchNo)+", errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}
	log.Error("ArBranchController Print")

	if err := controller.ArBranchService.BulkPrint(request, custId, userId); err != nil {
		log.Error("ArBranchController, BulkPrint, Service.BulkUpdate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Printed Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
*/

func (controller *ArBranchController) DistributorsFilter(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchDistributorsFilterQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, DistributorsFilter, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, DistributorsFilter, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.ArBranchService.DistributorsFilter(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, DistributorsFilter, data, err:", err.Error())
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

func (controller *ArBranchController) SuppliersFilter(c *fiber.Ctx) error {
	var dataFilter entity.ArBranchSuppliersFilterQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ArBranchController, SuppliersFilter, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, SuppliersFilter, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.ArBranchService.SuppliersFilter(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("ArBranchController, SuppliersFilter, data, err:", err.Error())
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

func (controller *ArBranchController) CreateArBranchPayment(c *fiber.Ctx) error {

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.CreateArBranchPaymentParams

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ArBranchController, CreateArBranchPayment, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, CreateArBranchPayment, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	var request entity.CreateArBranchPaymentBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("ArBranchController, CreateArBranchPayment, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	// log.Println("ArBranchController, Create, CustId:", custId)

	// request.CustID = custId
	request.ParentCustID = &parentCustId
	request.CreatedBy = &userId
	request.UpdatedBy = &userId
	request.GrBranchNo = &params.GrBranchNo

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ArBranchController, CreateArBranchPayment, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.ArBranchService.StoreArBranchPayment(request)
	if err != nil {
		log.Error("ArBranchController, CreateArBranchPayment, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
