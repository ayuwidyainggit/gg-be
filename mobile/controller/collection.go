package controller

import (
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type CollectionController struct {
	CollectionService service.CollectionService
	validator         *validation.Validate
}

func NewCollectionController(collectionService service.CollectionService, validator *validation.Validate) *CollectionController {
	return &CollectionController{
		CollectionService: collectionService,
		validator:         validator,
	}
}

func (controller *CollectionController) Route(app *fiber.App) {
	collectionRouteV1 := app.Group("/v1/collections", middleware.JWTProtected())
	collectionRouteV1.Post("", controller.CreateCollection)
	collectionRouteV1.Get("/missed-payment-reasons", controller.ListMissedPayment) // Static route first
	collectionRouteV1.Get("", controller.List)

	// Use unique path segments for deposit_no and invoice_no
	collectionRouteV1.Get("/deposit/:deposit_no", controller.Detail)        // Deposit-specific route
	collectionRouteV1.Get("/invoice/:invoice_no", controller.DetailInvoice) // Invoice-specific route

	collectionRouteV1.Post("/no-payment", controller.StoreCollectionNoPayment)
	collectionListRouteV1 := app.Group("/v1/collection-list", middleware.JWTProtected())
	collectionListRouteV1.Get("", controller.CollectionList)
	collectionListRouteV1.Post("", controller.CreateCollectionList)
}

func (controller *CollectionController) CreateCollection(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateDepositBodyByCollection
	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.CollectionService.StoreCollection(request)
	if err != nil {
		log.Error("CollectionController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
func (controller *CollectionController) StoreCollectionNoPayment(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateCollectionNoPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustID = custId
	request.CreatedBy = &userId

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.CollectionService.StoreCollectionNoPayment(request)
	if err != nil {
		log.Error("CollectionController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionController) List(c *fiber.Ctx) error {
	var dataFilter entity.CollectionQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("CollectionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, total, lastPage, err := controller.CollectionService.List(dataFilter)
	if err != nil {
		log.Error("CollectionController, List, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.PaginationWithTotalAmount{
		TotalRecord:  total.Total,
		TotalInvoice: total.TotalInvoice,
		PageCurrent:  dataFilter.Page,
		PageLimit:    dataFilter.Limit,
		PageTotal:    lastPage,
	})

	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionController) Detail(c *fiber.Ctx) error {
	var params entity.DetailDepositParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("CollectionController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.CollectionService.Detail(params.DepositNo, custId)
	if err != nil {
		log.Error("CollectionController, Detail, FindOneByOutletId, err:", err.Error())
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

func (controller *CollectionController) DetailInvoice(c *fiber.Ctx) error {
	var params entity.DetailInvoiceParams
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("CollectionController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	// log.Println("OutletController, Detail, CustId:", custId)

	data, err := controller.CollectionService.DetailInvoice(params.InvoiceNo, custId)
	if err != nil {
		log.Error("CollectionController, Detail, FindOneByOutletId, err:", err.Error())
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

func (controller *CollectionController) ListMissedPayment(c *fiber.Ctx) error {
	var (
		headerAcceptLang string
		dataFilter       entity.GeneralQueryFilter
	)

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, err := controller.CollectionService.ListMissedPayment(dataFilter)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setdata(data)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionController) CollectionList(c *fiber.Ctx) error {
	var dataFilter entity.CollectionQueryFilter
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("CollectionController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, CollectionList, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.CollectionService.CollectionList(dataFilter)
	if err != nil {
		log.Error("CollectionController, CollectionList, data, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.PaginationWithTotalAmount{
		TotalRecord:  total.Total,
		TotalInvoice: total.TotalInvoice,
		PageCurrent:  dataFilter.Page,
		PageLimit:    dataFilter.Limit,
		PageTotal:    lastPage,
	})

	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *CollectionController) CreateCollectionList(c *fiber.Ctx) error {
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), constant.HEADER_ACCEPT_LANG)

	var request entity.CreateCollectionListBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("CollectionController, CreateCollectionList, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		log.Error("CollectionController, CreateCollectionList, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.CollectionService.StoreCollectionList(request)
	if err != nil {
		log.Error("CollectionController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}
