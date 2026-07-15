package controller

import (
	"encoding/json"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type DistributorController struct {
	DistributorService service.DistributorService
	validator          *validation.Validate
}

func getDistributorAcceptLanguage(c *fiber.Ctx) string {
	return c.Get(constant.HEADER_ACCEPT_LANG)
}

func NewDistributorController(distributorService service.DistributorService, validator *validation.Validate) *DistributorController {
	return &DistributorController{
		DistributorService: distributorService,
		validator:          validator,
	}
}

func (controller *DistributorController) Route(app *fiber.App) {
	qParamId := ":distributor_id"
	distributorRouteV1 := app.Group("/v1/distributors", middleware.JWTProtected())
	distributorRouteV1.Post("", controller.Create)
	distributorRouteV1.Get("/customers", controller.ListWithCustomer)
	distributorRouteV1.Get("/"+qParamId, controller.Detail)
	distributorRouteV1.Get("", controller.List)
	distributorRouteV1.Delete("/"+qParamId, controller.Delete)
	distributorRouteV1.Patch("/"+qParamId, controller.Update)
}

func (controller *DistributorController) Create(c *fiber.Ctx) error {
	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateDistributorBody
	if err := c.BodyParser(&request); err != nil {
		log.Info("distributorController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)
	request.CreatedBy = &userId
	// log.Info("distributorController, Create, CustId:", custId)

	if request.CustId != request.ParentCustId {
		responsePayload.Setmsg("Distributor cannot be created outside of the parent customer")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("distributorController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.DistributorService.Store(request)
	if err != nil {
		log.Info("distributorController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *DistributorController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.DistributorQueryFilter
		data       interface{}
		total      int
		lastPage   int
		sup        []entity.DistributorResponse
		supLookup  []entity.DistributorLookupResponse
	)

	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("DistributorController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	tokenDistributorID := c.Locals("distributor_id").(int64)
	if tokenDistributorID > 0 {
		dataFilter.JwtDistributorId = tokenDistributorID
	}
	// log.Info("BankController, List, CustId:", custId)

	switch dataFilter.Mode {
	case "lookup":
		data, total, lastPage, err = controller.DistributorService.LookupList(dataFilter, dataFilter.CustId)
		if err != nil {
			log.Info("DistributorController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
		err = structs.Automapper(supLookup, &data)
		if err != nil {
			log.Info("DistributorController, List, data, err:", err.Error())
			responsePayload.Setmsg(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
	default:
		data, total, lastPage, err = controller.DistributorService.List(dataFilter, dataFilter.CustId)
		if err != nil {
			log.Info("DistributorController, List, data, err:", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
				RequestId: c.Locals("requestid").(string),
				Message:   err.Error(),
			})
		}
		err = structs.Automapper(sup, &data)
		if err != nil {
			log.Info("DistributorController, List, data, err:", err.Error())
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

func (controller *DistributorController) Detail(c *fiber.Ctx) error {
	var params entity.DetailDistributorParams
	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SpPriceController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SpPriceController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustId = c.Locals("cust_id").(string)
	params.ParentCustId = c.Locals("parent_cust_id").(string)
	tokenDistributorID := c.Locals("distributor_id").(int64)
	if tokenDistributorID > 0 {
		params.JwtDistributorId = tokenDistributorID
	}
	// log.Info("OutletController, Detail, CustId:", custId)

	data, err := controller.DistributorService.Detail(params)
	if err != nil {
		log.Info("SpPriceController, Detail, err:", err.Error())
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

func normalizeUpdateDistributorRequest(request *entity.UpdateDistributorRequest) {
	if request.Barcode != nil && strings.TrimSpace(*request.Barcode) == "" {
		request.Barcode = nil
	}

	if request.ZipCode != nil && strings.TrimSpace(*request.ZipCode) == "" {
		request.ZipCode = nil
	}

	for i := range request.Contacts {
		if request.Contacts[i].Email == nil {
			emptyEmail := ""
			request.Contacts[i].Email = &emptyEmail
		}
	}
}

func preserveNullableDistributorStringPatch(rawRequest map[string]json.RawMessage, key string, target **string) {
	rawValue, ok := rawRequest[key]
	if !ok {
		return
	}

	trimmedValue := strings.TrimSpace(string(rawValue))
	if trimmedValue == "null" {
		emptyValue := ""
		*target = &emptyValue
		return
	}

	var parsedValue string
	if err := json.Unmarshal(rawValue, &parsedValue); err != nil {
		return
	}

	if strings.TrimSpace(parsedValue) == "" {
		emptyValue := ""
		*target = &emptyValue
	}
}

func preserveNullableDistributorIntPatch(rawRequest map[string]json.RawMessage, key string, target **int) {
	rawValue, ok := rawRequest[key]
	if !ok {
		return
	}

	trimmedValue := strings.TrimSpace(string(rawValue))
	if trimmedValue == "null" {
		zeroValue := 0
		*target = &zeroValue
		return
	}

	var parsedValue int
	if err := json.Unmarshal(rawValue, &parsedValue); err != nil {
		return
	}

	*target = &parsedValue
}

func markNullableDistributorStringPresence(rawRequest map[string]json.RawMessage, key string) bool {
	_, ok := rawRequest[key]
	return ok
}

func (controller *DistributorController) Update(c *fiber.Ctx) error {
	var (
		params     entity.UpdateDistributorParams
		request    entity.UpdateDistributorRequest
		rawRequest map[string]json.RawMessage
	)
	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SpPriceController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SpPriceController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Info("SpPriceController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := json.Unmarshal(c.Body(), &rawRequest); err != nil {
		log.Info("SpPriceController, Update, Unmarshal(rawRequest), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("BankController, Update, CustId:", custId)

	request.CustId = custId
	request.UpdatedBy = &userId
	normalizeUpdateDistributorRequest(&request)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Info("SpPriceController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	preserveNullableDistributorStringPatch(rawRequest, "barcode", &request.Barcode)
	preserveNullableDistributorStringPatch(rawRequest, "province_id", &request.ProvinceId)
	preserveNullableDistributorStringPatch(rawRequest, "regency_id", &request.RegencyId)
	preserveNullableDistributorStringPatch(rawRequest, "sub_district_id", &request.SubDistrictId)
	preserveNullableDistributorStringPatch(rawRequest, "ward_id", &request.WardId)
	preserveNullableDistributorStringPatch(rawRequest, "zip_code", &request.ZipCode)
	preserveNullableDistributorIntPatch(rawRequest, "ot_loc_id", &request.OtLocId)
	preserveNullableDistributorStringPatch(rawRequest, "phone", &request.Phone)
	preserveNullableDistributorStringPatch(rawRequest, "fax_number", &request.FaxNumber)
	request.BarcodeProvided = markNullableDistributorStringPresence(rawRequest, "barcode")
	request.ProvinceIdProvided = markNullableDistributorStringPresence(rawRequest, "province_id")
	request.RegencyIdProvided = markNullableDistributorStringPresence(rawRequest, "regency_id")
	request.SubDistrictIdProvided = markNullableDistributorStringPresence(rawRequest, "sub_district_id")
	request.WardIdProvided = markNullableDistributorStringPresence(rawRequest, "ward_id")
	request.ZipCodeProvided = markNullableDistributorStringPresence(rawRequest, "zip_code")
	request.OtLocIdProvided = markNullableDistributorStringPresence(rawRequest, "ot_loc_id")
	request.PhoneProvided = markNullableDistributorStringPresence(rawRequest, "phone")
	request.FaxNumberProvided = markNullableDistributorStringPresence(rawRequest, "fax_number")

	err := controller.DistributorService.Update(params.DistributorId, request)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if service.IsDistributorNotFoundError(err) {
			statusCode = fiber.StatusNotFound
		}

		log.Info("SpPriceController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DistributorController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteDistributorParams
	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Info("SpPriceController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Info("SpPriceController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Info("SpPriceController, Delete, CustId:", custId)

	err := controller.DistributorService.Delete(custId, int64(params.DistributorId), userId)
	if err != nil {
		log.Info("SpPriceController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *DistributorController) ListWithCustomer(c *fiber.Ctx) error {
	var (
		err          error
		dataFilter   entity.DistributorQueryFilter
		data         interface{}
		total        int
		lastPage     int
		distributors []entity.DistributorCustomerResp
	)

	log.Infof("ListWithCustomer")

	headerAcceptLang := getDistributorAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Info("DistributorController, ListWithCustomer, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	data, total, lastPage, err = controller.DistributorService.ListWithCustomer(dataFilter, custId)
	if err != nil {
		log.Info("DistributorController, ListWithCustomer, data, err:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			RequestId: c.Locals("requestid").(string),
			Message:   err.Error(),
		})
	}
	err = structs.Automapper(distributors, &data)
	if err != nil {
		log.Info("DistributorController, ListWithCustomer, data, err:", err.Error())
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
