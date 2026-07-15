package controller

import (
	"fmt"
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	OrderService    service.OrderService
	DiscountService service.DiscountService
	validator       *validation.Validate
}

func NewOrderController(roService service.OrderService, discountService service.DiscountService, validator *validation.Validate) *OrderController {
	return &OrderController{
		OrderService:    roService,
		DiscountService: discountService,
		validator:       validator,
	}
}

func (controller *OrderController) Route(app *fiber.App) {
	qParamId := ":ro_no"
	roRouteV1 := app.Group("/v1/orders", middleware.JWTProtected())
	roRouteV1.Post("", controller.Create)
	roRouteV1.Post("/no-order", controller.CreateNoOrder)
	roRouteV1.Get("/no-order", controller.ListNoOrder)
	roRouteV1.Get("/sales-report", controller.SummaryBySalesman)
	roRouteV1.Get("/statistics", controller.StatisticReport)
	roRouteV1.Get("/"+qParamId, controller.Detail)
	roRouteV1.Get("", controller.List)
	roRouteV1.Patch("/final/"+qParamId, controller.UpdateFinal)
	roRouteV1.Patch("/status", controller.UpdateStatus)
	roRouteV1.Patch("/"+qParamId, controller.Update)
	roRouteV1.Delete("/"+qParamId, controller.Delete)
	roRouteV1.Post("/conversion", controller.Conversion)

	//routeoutlet
	roRouteOutletV1 := app.Group("/v1/outlets", middleware.JWTProtected())
	roRouteOutletV1.Get("", controller.LookupSalesman)

}

func (controller *OrderController) Create(c *fiber.Ctx) error {
	var (
		request entity.CreateOrderBody
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OrderService.Store(c.Context(), request)
	if err != nil {
		log.Error("OrderController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())

}

func (controller *OrderController) CreateNoOrder(c *fiber.Ctx) error {
	var (
		request entity.CreateNoOrderBody
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OrderService.StoreNoOrder(request)
	if err != nil {
		log.Error("OrderController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())

}

func (controller *OrderController) Detail(c *fiber.Ctx) error {
	var (
		params entity.DetailOrderParams
	)
	// var params entity.DetailOrderParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.OrderService.Detail(params.RoNo, custId)
	if err != nil {
		log.Error("OrderController, Detail, FindOneByOutletId, err:", err.Error())
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

func (controller *OrderController) List(c *fiber.Ctx) error {
	var dataFilter entity.OrderQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.OrderService.List(dataFilter)
	if err != nil {
		log.Error("OrderController, List, data, err:", err.Error())
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

func (controller *OrderController) ListNoOrder(c *fiber.Ctx) error {
	var dataFilter entity.NoOrderQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.OrderService.ListNoOrder(dataFilter)
	if err != nil {
		log.Error("OrderController, List, data, err:", err.Error())
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

func (controller *OrderController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateOrderParams
		request entity.UpdateOrderBody
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("VanOrderController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OrderService.Update(params.RoNo, request)
	if err != nil {
		log.Error("VanOrderController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) Delete(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	var params entity.DetailOrderParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.OrderService.Delete(custId, params.RoNo, userId)
	if err != nil {
		log.Error("OrderController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) Conversion(c *fiber.Ctx) error {
	// var params entity.DetailProductParams
	// responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG])

	var request entity.CreateConversionBody
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Conversion, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("ProductController, Create, CustId:", custId)

	request.CustId = custId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Conversion, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.OrderService.Conversion(request, custId, parentCustId)
	if err != nil {
		log.Error("OrderController, Conversion, FindOneProductByProductIdAndCustId, err:", err.Error())
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

func (controller *OrderController) UpdateFinal(c *fiber.Ctx) error {
	var (
		params  entity.UpdateOrderParams
		request entity.UpdateOrderDetailFinal
	)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, UpdateFinal, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, UpdateFinal, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, UpdateFinal, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, UpdateFinal, CustId:", custId)
	request.CustId = custId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, UpdateFinal, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OrderService.UpdateFinal(params.RoNo, request)
	if err != nil {
		log.Error("OrderController, UpdateFinal, Service.UpdateFinal, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) UpdateStatus(c *fiber.Ctx) error {
	var request entity.BulkUpdateStatusOrder
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, BulkUpdate, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// request.CustId = custId
	// request.UpdatedBy = userId

	if errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG); errs != nil {
		log.Error("OrderController, BulkUpdate, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for index := range request.Orders {
		// request.Orders[index].CustId = custId
		request.Orders[index].UpdatedBy = userId

		if errs := controller.validator.ValidateStruct(request.Orders[index], constant.HEADER_ACCEPT_LANG); errs != nil {
			log.Error("OrderController, BulkUpdate, ValidateStruct Order with RoNo "+fmt.Sprint(request.Orders[index].RoNo)+", errs:", errs)
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(errs)
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	}

	if err := controller.OrderService.BulkUpdateStatus(custId, request); err != nil {
		log.Error("OrderController, Update, Service.BulkUpdate, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// err := controller.InvoiceService.Update(params.RoNo, request)
	// if err != nil {
	// 	log.Println("OrderController, Update, Service.Update, err:", err.Error())
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }
	responsePayload.Setmsg("Updated Status Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) LookupSalesman(c *fiber.Ctx) error {
	var dataFilter entity.OrderQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.OrderService.LookupSalesman(dataFilter)
	if err != nil {
		log.Error("OrderController, List, data, err:", err.Error())
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

func (controller *OrderController) SummaryBySalesman(c *fiber.Ctx) error {
	var dataFilter entity.SummaryTotalFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.OrderService.SummaryBySalesman(dataFilter)
	if err != nil {
		log.Error("OrderController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) StatisticReport(c *fiber.Ctx) error {
	var dataFilter entity.StatisticReportFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, StatisticReport, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("OrderController, StatisticReport, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	isDistributor, _ := c.Locals("is_distributor").(bool)
	dataFilter.IsDistributor = isDistributor

	data, err := controller.OrderService.StatisticReport(dataFilter)
	if err != nil {
		log.Error("OrderController, StatisticReport, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("OK")
	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}
