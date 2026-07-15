package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sales/entity"
	"sales/pkg/constant"
	"sales/pkg/middleware"
	"sales/pkg/responsebuild"
	"sales/pkg/structs"
	"sales/pkg/validation"
	"sales/service"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	OrderService         service.OrderService
	ValidateOrderService service.ValidateOrderService
	validator            *validation.Validate
}

func NewOrderController(orderService service.OrderService, validateOrderService service.ValidateOrderService, validator *validation.Validate) *OrderController {
	return &OrderController{
		OrderService:         orderService,
		ValidateOrderService: validateOrderService,
		validator:            validator,
	}
}

func (controller *OrderController) Route(app *fiber.App) {
	qParamId := ":ro_no"
	qParamProductId := ":pro_id"

	roRouteV1 := app.Group("/v1/orders", middleware.JWTProtected())
	roRouteV1.Post("", controller.Create)
	roRouteV1.Get("/discount", controller.DetailDiscount)
	roRouteV1.Get("/export-template", controller.ExportTemplate)
	roRouteV1.Get("/minimum-price/"+qParamProductId+"/product", controller.GetMinimumPriceProduct)
	roRouteV1.Get("", controller.List)
	roRouteV1.Get("/"+qParamId, controller.Detail)
	roRouteV1.Patch("/final/"+qParamId, controller.UpdateFinal)
	roRouteV1.Patch("/status", controller.UpdateStatus)
	roRouteV1.Patch("/"+qParamId, controller.Update)
	roRouteV1.Patch("/enhance/:ro_no", controller.UpdateEnhance)
	roRouteV1.Delete("/"+qParamId, controller.Delete)
	roRouteV1.Post("/import", controller.Import)
	roRouteV1.Post("/export-template/import", controller.ImportFromUrl)
	roRouteV1.Post("/conversion", controller.Conversion)

	roRouteOutletV1 := app.Group("/v1/outlets", middleware.JWTProtected())
	roRouteOutletV1.Get("", controller.LookupSalesman)

	roRouteV2 := app.Group("/v2/orders", middleware.JWTProtected())
	roRouteV2.Get("/"+qParamId, controller.DetailV2)

	proformaRouteV1 := app.Group("/v1/proforma_invoice", middleware.JWTProtected())
	proformaRouteV1.Get("", controller.ProformaInvoiceList)

	printProformaRouteV1 := app.Group("/v1/print_proforma_invoice", middleware.JWTProtected())
	printProformaRouteV1.Post("", controller.PrintProformaInvoice)

}

func (controller *OrderController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateOrderBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.CreatedBy = &userId
	if request.OrderType != nil && strings.TrimSpace(*request.OrderType) == "" {
		request.OrderType = nil
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// VALIDATE ORDER
	var validateOrderRequest entity.ValidateOrderBody
	if err := structs.Automapper(request, &validateOrderRequest); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Main Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := structs.Automapper(request.Details.Normal, &validateOrderRequest.ProStok); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Detail Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// err := controller.OrderService.SetValidateOrderRequest("", &validateOrderRequest)
	// if err != nil {
	// 	log.Error("OrderController, Create, ValidateOrder, err:", err.Error())
	// 	responsePayload.Setmsg(err.Error())
	// 	return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	// }

	validationData := service.BuildCreateOrderValidationBypassResponse(request.OrderType)
	var validateErr error
	if service.IsTakingOrder(request.OrderType) {
		validationData, _, _, validateErr = controller.ValidateOrderService.ValidateOrderWithoutStock(validateOrderRequest)
	} else if service.ShouldValidateStockOnCreate(request.OrderType) {
		validationData, _, _, validateErr = controller.ValidateOrderService.ValidateOrder(validateOrderRequest)
	}
	if validateErr != nil {
		log.Error("OrderController, Create, ValidateOrder, err:", validateErr.Error())
		responsePayload.Setmsg(validateErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// STORE ORDER
	data, err := controller.OrderService.Store(request, validationData)
	if err != nil {
		log.Error("OrderController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Created Successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) Detail(c *fiber.Ctx) error {
	var params entity.DetailOrderParams
	var headerAcceptLang string
	var dataFilter entity.DetailOrderQueryParams

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if !dataFilter.NoCustID {
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
	} else {
		data, err := controller.OrderService.DetailNoCustID(params.RoNo, dataFilter.CustIDOrigin, dataFilter.EmpID)
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
}

func (controller *OrderController) DetailV2(c *fiber.Ctx) error {
	var params entity.DetailOrderParams
	var headerAcceptLang string
	var dataFilter entity.DetailOrderQueryParams

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, DetailV2, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, DetailV2, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, DetailV2, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if !dataFilter.NoCustID {
		custId := c.Locals("cust_id").(string)
		parentCustId := c.Locals("parent_cust_id").(string)

		data, err := controller.OrderService.DetailV2(params.RoNo, custId, parentCustId)
		if err != nil {
			log.Error("OrderController, DetailV2, DetailV2, err:", err.Error())
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
	} else {
		data, err := controller.OrderService.DetailNoCustID(params.RoNo, dataFilter.CustIDOrigin, dataFilter.EmpID)
		if err != nil {
			log.Error("OrderController, DetailV2, DetailNoCustID, err:", err.Error())
			statusCode := fiber.StatusBadRequest
			errMsg := err.Error()
			if err.Error() == "sql: no rows in result set" {
				statusCode = fiber.StatusNotFound
				errMsg = "Not found"
			}

			responsePayload.Setmsg(errMsg)
			return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
		}
		// Map data_source to source for DetailNoCustID response
		if data.DataSource != nil {
			data.Source = service.MapDataSourceToSource(data.DataSource)
		}
		// Set is_proforma_inv (default false if null)
		if data.IsProformaInv == nil {
			falseVal := false
			data.IsProformaInv = &falseVal
		}
		// Copy Details to PurchaseDetails
		data.PurchaseDetails.Normal = make([]entity.OrderDetResponse, len(data.Details.Normal))
		copy(data.PurchaseDetails.Normal, data.Details.Normal)
		data.PurchaseDetails.Promo = make([]entity.OrderDetResponse, len(data.Details.Promo))
		copy(data.PurchaseDetails.Promo, data.Details.Promo)
		responsePayload.Setdata(data)
		return c.JSON(responsePayload.GetRespPayload())
	}
}

func (controller *OrderController) List(c *fiber.Ctx) error {
	var dataFilter entity.OrderQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
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

func (controller *OrderController) ProformaInvoiceList(c *fiber.Ctx) error {
	var dataFilter entity.ProformaInvoiceQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, ProformaInvoiceList, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Manual parse array parameters (Fiber QueryParser might not handle array correctly)
	// Get all values for salesman_id (supports both salesman_id=12&salesman_id=15 and salesman_id[]=12&salesman_id[]=15)
	allSalesmanIds := c.Context().QueryArgs().PeekMulti("salesman_id")
	allSalesmanIdsBracket := c.Context().QueryArgs().PeekMulti("salesman_id[]")

	// Combine both formats
	dataFilter.SalesmanId = []int{}
	if len(allSalesmanIds) > 0 {
		for _, idBytes := range allSalesmanIds {
			if id, err := strconv.Atoi(string(idBytes)); err == nil {
				dataFilter.SalesmanId = append(dataFilter.SalesmanId, id)
			}
		}
	}
	if len(allSalesmanIdsBracket) > 0 {
		for _, idBytes := range allSalesmanIdsBracket {
			if id, err := strconv.Atoi(string(idBytes)); err == nil {
				dataFilter.SalesmanId = append(dataFilter.SalesmanId, id)
			}
		}
	}
	// Also try without [] bracket (for single value or comma-separated)
	if len(dataFilter.SalesmanId) == 0 {
		salesmanIds := c.Query("salesman_id")
		if salesmanIds != "" {
			ids := strings.Split(salesmanIds, ",")
			for _, idStr := range ids {
				if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					dataFilter.SalesmanId = append(dataFilter.SalesmanId, id)
				}
			}
		}
	}

	// Get all values for outlet_id
	allOutletIds := c.Context().QueryArgs().PeekMulti("outlet_id")
	allOutletIdsBracket := c.Context().QueryArgs().PeekMulti("outlet_id[]")

	// Combine both formats
	dataFilter.OutletID = []int{}
	if len(allOutletIds) > 0 {
		for _, idBytes := range allOutletIds {
			if id, err := strconv.Atoi(string(idBytes)); err == nil {
				dataFilter.OutletID = append(dataFilter.OutletID, id)
			}
		}
	}
	if len(allOutletIdsBracket) > 0 {
		for _, idBytes := range allOutletIdsBracket {
			if id, err := strconv.Atoi(string(idBytes)); err == nil {
				dataFilter.OutletID = append(dataFilter.OutletID, id)
			}
		}
	}
	// Also try without [] bracket
	if len(dataFilter.OutletID) == 0 {
		outletIds := c.Query("outlet_id")
		if outletIds != "" {
			ids := strings.Split(outletIds, ",")
			for _, idStr := range ids {
				if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					dataFilter.OutletID = append(dataFilter.OutletID, id)
				}
			}
		}
	}

	// Validate after manual parsing
	// Check if salesman_id is required and provided
	if len(dataFilter.SalesmanId) == 0 {
		errs := map[string]string{
			"salesman_id": "salesman_id is required",
		}
		log.Error("OrderController, ProformaInvoiceList, salesman_id is required but empty")
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, ProformaInvoiceList, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.OrderService.ProformaInvoiceList(dataFilter)
	if err != nil {
		log.Error("OrderController, ProformaInvoiceList, data, err:", err.Error())
		// Handle empty state
		if err.Error() == "record not found" || total == 0 {
			responsePayload.Setmsg("No Data")
			responsePayload.Setdata(nil)
			responsePayload.Setpaging(entity.Pagination{
				TotalRecord: 0,
				PageCurrent: dataFilter.Page,
				PageLimit:   dataFilter.Limit,
				PageTotal:   0,
			})
			return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Handle empty state
	if total == 0 || len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
		responsePayload.Setpaging(entity.Pagination{
			TotalRecord: 0,
			PageCurrent: dataFilter.Page,
			PageLimit:   dataFilter.Limit,
			PageTotal:   0,
		})
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
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

func (controller *OrderController) PrintProformaInvoice(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.PrintProformaInvoiceRequest
	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, PrintProformaInvoice, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	ctx := c.Context()
	response, err := controller.OrderService.PrintProformaInvoice(ctx, request, custId, userId)
	if err != nil {
		// Handle specific error cases
		if strings.Contains(err.Error(), "not found") {
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(response)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) DetailDiscount(c *fiber.Ctx) error {
	var dataFilter entity.OrderDiscountQuery
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, Detail discount, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, Detail discount, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.OrderService.DetailDiscount(dataFilter)

	if err != nil {
		log.Error("OrderController, Detail discount err:", err.Error())
		statusCode := fiber.StatusNotFound
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	if reflect.DeepEqual(data, entity.DiscountCriteria{}) {
		data.SlabDesc = "Discount Not Found"
	}

	responsePayload.Setmsg(data.SlabDesc)
	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) Update(c *fiber.Ctx) error {
	var params entity.UpdateOrderParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	body := c.Body()
	if isEnhancePatchRequest(body) {
		request := entity.EditOrderEnhanceBody{}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &request); err != nil {
				log.Error("OrderController, Update, Unmarshal enhance request, err:", err.Error())
				responsePayload.Setmsg(fiber.ErrBadRequest.Message)
				return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
			}
		}

		request.CustId = custId
		request.ParentCustId = parentCustId
		request.UpdatedBy = userId

		if isEmptyEnhanceRequest(request) {
			err := controller.OrderService.ProcessEnhanceWithoutProductEdit(c.UserContext(), params.RoNo, custId, userId)
			if err != nil {
				log.Error("OrderController, Update, Service.ProcessEnhanceWithoutProductEdit, err:", err.Error())
				responsePayload.Setmsg(err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
			}

			responsePayload.Setmsg("Updated Successfully")
			return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
		}

		err := controller.OrderService.UpdateEnhance(c.UserContext(), params.RoNo, request)
		if err != nil {
			log.Error("OrderController, Update, Service.UpdateEnhance, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}

		responsePayload.Setmsg("Updated Successfully")
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	var request entity.UpdateOrderBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = custId
	request.ParentCustId = parentCustId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("VanOrderController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	var validateOrderRequest entity.ValidateOrderBody
	if err := structs.Automapper(request, &validateOrderRequest); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Main Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := structs.Automapper(request.Details.Normal, &validateOrderRequest.ProStok); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Detail Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	validationData, _, _, err := controller.ValidateOrderService.ValidateOrder(validateOrderRequest)
	if err != nil {
		log.Error("OrderController, Create, ValidateOrder, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.OrderService.Update(params.RoNo, request, validationData)
	if err != nil {
		log.Error("OrderController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) UpdateFinal(c *fiber.Ctx) error {
	var (
		params  entity.UpdateOrderParams
		request entity.UpdateOrderDetailFinal
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
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
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("BankController, Update, CustId:", custId)
	request.CustId = custId
	request.ParentCustId = parentCustId
	request.UpdatedBy = userId

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("VanOrderController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// VALIDATE ORDER
	var validateOrderRequest entity.ValidateOrderBody
	if err := structs.Automapper(request, &validateOrderRequest); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Main Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := structs.Automapper(request.Details.Normal, &validateOrderRequest.ProStok); err != nil {
		log.Error("OrderController, Create, ValidateOrder Mapping Detail Request, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.OrderService.SetValidateOrderRequest(params.RoNo, &validateOrderRequest)
	if err != nil {
		log.Error("OrderController, Create, ValidateOrder, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	validationData, _, _, err := controller.ValidateOrderService.ValidateOrder(validateOrderRequest)
	if err != nil {
		log.Error("OrderController, Create, ValidateOrder, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	err = controller.OrderService.UpdateFinal(params.RoNo, request, validationData)
	if err != nil {
		log.Error("VanOrderController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Updated Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *OrderController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DetailOrderParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
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

func (controller *OrderController) ExportTemplate(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()["Accept-Language"]) > 0 {
		headerAcceptLang = c.GetReqHeaders()["Accept-Language"][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	format := strings.ToLower(strings.TrimSpace(c.Query("format", "xlsx")))
	if format == "xls" {
		format = "xlsx"
	}
	if format != "" && format != "xlsx" {
		responsePayload.Setmsg("unsupported format")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	requestedFormat := strings.ToLower(strings.TrimSpace(c.Query("format", "xlsx")))
	buf, contentType, filename, err := controller.OrderService.ExportTemplate(format)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if requestedFormat == "xls" {
		contentType = "application/vnd.ms-excel"
		filename = "order_import_template.xls"
	}
	c.Set(fiber.HeaderContentType, contentType)
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+filename+`"`)
	return c.Send(buf.Bytes())
}

func (controller *OrderController) Import(c *fiber.Ctx) error {
	headerAcceptLang := ""
	if len(c.GetReqHeaders()["Accept-Language"]) > 0 {
		headerAcceptLang = c.GetReqHeaders()["Accept-Language"][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		responsePayload.Setmsg("missing file field")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	file, err := fileHeader.Open()
	if err != nil {
		responsePayload.Setmsg("cannot open file")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	defer file.Close()
	result, errs, svcErr := controller.OrderService.ImportOrders(custId, parentCustId, userId, file, fileHeader.Filename)
	if svcErr != nil {
		var failedErr *entity.ImportFailedError
		if errors.As(svcErr, &failedErr) {
			failedSummary := entity.OrderImportSummary{
				StartDate:       result.StartDate,
				EndDate:         result.EndDate,
				NumberOfInvoice: result.NumberOfInvoice,
				NumberOfOutlet:  result.NumberOfOutlet,
				Amount:          result.Amount,
				FailedReasons:   failedErr.FailedReasons,
			}
			responsePayload.Setmsg("Secondary sales file validate failed")
			responsePayload.Setdata(failedSummary)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(svcErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	_ = errs
	responsePayload.Setmsg("Secondary sales file imported successfully.")
	responsePayload.Setdata(result)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// ImportFromUrl handles FE-uploaded-file-then-validate flow for the
// /v1/orders/export-template/import alias route. The FE uploads the
// file via master /v1/files/uploads, gets a URL, and POSTs
// {url, validate} here. When validate is "False" (or empty/falsey)
// only the validation summary is returned, otherwise a real import
// is performed via OrderService.ImportOrders.
func (controller *OrderController) ImportFromUrl(c *fiber.Ctx) error {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(c.Get(fiber.HeaderContentType))), "multipart/form-data") {
		return controller.Import(c)
	}
	headerAcceptLang := ""
	if len(c.GetReqHeaders()["Accept-Language"]) > 0 {
		headerAcceptLang = c.GetReqHeaders()["Accept-Language"][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	var request entity.OrderImportFromURLRequest
	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg("invalid json body")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	if strings.TrimSpace(request.URL) == "" {
		responsePayload.Setmsg("url is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Get(request.URL)
	if err != nil {
		responsePayload.Setmsg("cannot download file: " + err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(responsePayload.GetRespPayload())
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responsePayload.Setmsg("cannot download file: unexpected status")
		return c.Status(fiber.StatusBadGateway).JSON(responsePayload.GetRespPayload())
	}

	filename := filenameFromURL(request.URL)

	shouldImport := isTruthy(request.Validate)
	if !shouldImport {
		summary, svcErr := controller.OrderService.ValidateImport(custId, parentCustId, userId, resp.Body, filename)
		if svcErr != nil {
			responsePayload.Setmsg(svcErr.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if len(summary.FailedReasons) > 0 {
			responsePayload.Setmsg("Secondary sales file validate failed")
			responsePayload.Setdata(summary)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Secondary sales file validate successfully.")
		responsePayload.Setdata(summary)
		return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
	}

	result, errs, svcErr := controller.OrderService.ImportOrders(custId, parentCustId, userId, resp.Body, filename)
	if svcErr != nil {
		var failedErr *entity.ImportFailedError
		if errors.As(svcErr, &failedErr) {
			failedSummary := entity.OrderImportSummary{
				StartDate:       result.StartDate,
				EndDate:         result.EndDate,
				NumberOfInvoice: result.NumberOfInvoice,
				NumberOfOutlet:  result.NumberOfOutlet,
				Amount:          result.Amount,
				FailedReasons:   failedErr.FailedReasons,
			}
			responsePayload.Setmsg("Secondary sales file validate failed")
			responsePayload.Setdata(failedSummary)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(svcErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	_ = errs
	responsePayload.Setmsg("Secondary sales file imported successfully.")
	responsePayload.Setdata(result)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "", "false", "0", "no":
		return false
	}
	return true
}

func filenameFromURL(raw string) string {
	clean := raw
	if idx := strings.Index(clean, "?"); idx >= 0 {
		clean = clean[:idx]
	}
	if idx := strings.LastIndex(clean, "/"); idx >= 0 {
		clean = clean[idx+1:]
	}
	if clean == "" {
		return "imported.xlsx"
	}
	return clean
}

func (controller *OrderController) Conversion(c *fiber.Ctx) error {
	// var params entity.DetailProductParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateConversionBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, Conversion, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	// log.Println("ProductController, Create, CustId:", custId)

	request.CustId = custId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
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

func (controller *OrderController) LookupSalesman(c *fiber.Ctx) error {
	var dataFilter entity.OrderQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
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

func (controller *OrderController) UpdateStatus(c *fiber.Ctx) error {
	var request entity.BulkUpdateStatusOrder
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, BulkUpdate, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// request.CustId = custId
	// request.UpdatedBy = userId

	if errs := controller.validator.ValidateStruct(request, headerAcceptLang); errs != nil {
		log.Error("OrderController, BulkUpdate, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	for index := range request.Orders {
		// request.Orders[index].CustId = custId
		request.Orders[index].UpdatedBy = userId

		if errs := controller.validator.ValidateStruct(request.Orders[index], headerAcceptLang); errs != nil {
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

func (controller *OrderController) GetMinimumPriceProduct(c *fiber.Ctx) error {
	var request entity.OrderMinimumPriceFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&request); err != nil {
		log.Error("OrderController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.ParamsParser(&request); err != nil {
		log.Error("OrderController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, GetMinimumPriceProduct, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	// log.Println("BankController, List, CustId:", custId)

	data, err := controller.OrderService.GetMinimumPriceProduct(request)
	if err != nil {
		log.Error("OrderController, GetMinimumPriceProduct, err:", err.Error())
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

// UpdateEnhance handles enhanced edit order for specific tabs (Purchase Order, Sales Order, Final Order)
func (controller *OrderController) UpdateEnhance(c *fiber.Ctx) error {
	var (
		params  entity.UpdateOrderParams
		request entity.EditOrderEnhanceBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("OrderController, UpdateEnhance, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("OrderController, UpdateEnhance, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("OrderController, UpdateEnhance, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	request.CustId = custId
	request.ParentCustId = parentCustId
	request.UpdatedBy = userId

	err := controller.OrderService.UpdateEnhance(c.UserContext(), params.RoNo, request)
	if err != nil {
		log.Error("OrderController, UpdateEnhance, Service.UpdateEnhance, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Berhasil Diperbarui")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func isEnhancePatchRequest(body []byte) bool {
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" || trimmed == "{}" {
		return true
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}

	enhanceKeys := []string{
		"purchase_order",
		"purchase_details",
		"add_purchase_order",
		"add_purchase_details",
		"sales_order",
		"sales_order_details",
		"add_sales_order",
		"final_order",
		"final_order_details",
		"add_final_order",
	}

	for _, key := range enhanceKeys {
		if _, ok := payload[key]; ok {
			return true
		}
	}

	return false
}

func isEmptyEnhanceRequest(request entity.EditOrderEnhanceBody) bool {
	return len(request.PurchaseOrder) == 0 &&
		len(request.PurchaseDetails) == 0 &&
		len(request.AddPurchaseOrder) == 0 &&
		len(request.AddPurchaseDetails) == 0 &&
		len(request.SalesOrder) == 0 &&
		len(request.SalesOrderDetails) == 0 &&
		len(request.AddSalesOrder) == 0 &&
		len(request.FinalOrder) == 0 &&
		len(request.FinalOrderDetails) == 0 &&
		len(request.AddFinalOrder) == 0
}
