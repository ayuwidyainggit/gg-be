package controller

import (
	"encoding/json"
	"fmt"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/rabbitmq"
	"master/pkg/responsebuild"
	"master/pkg/structs"
	"master/pkg/validation"
	"master/service"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
	"github.com/valyala/fasthttp"
)

type MPriceController struct {
	MPriceService service.MPriceService
	validator     *validation.Validate
}

func NewMPriceController(priceService service.MPriceService, validator *validation.Validate) MPriceController {
	return MPriceController{
		MPriceService: priceService,
		validator:     validator,
	}
}

func (controller *MPriceController) Route(app *fiber.App) {
	go rabbitmq.Subscribe(constant.RMQ_MANAGE_PRICE_CREATE_EVENT, controller.processPublishManagePriceMessage)
	qParamId := ":price_id"
	priceRouteV1 := app.Group("/v1/prices", middleware.JWTProtected())
	priceRouteV1.Get("/statuses", controller.GetStatuses)
	priceRouteV1.Get("/template", controller.Template)
	priceRouteV1.Get("/export", controller.Export)
	priceRouteV1.Get("/"+qParamId, controller.Detail)
	priceRouteV1.Get("", controller.List)
	priceRouteV1.Post("", controller.Create)
	priceRouteV1.Post("/import", controller.Import)
	priceRouteV1.Patch("/publish/"+qParamId, controller.Publish)
	priceRouteV1.Patch("/"+qParamId+"/publish", controller.Publish)
	priceRouteV1.Patch("/cancel/"+qParamId, controller.Cancel)
	priceRouteV1.Patch("/"+qParamId+"/cancel", controller.Cancel)
	priceRouteV1.Patch("/"+qParamId, controller.Update)
	priceRouteV1.Put("/"+qParamId, controller.Update)
	priceRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *MPriceController) Detail(c *fiber.Ctx) error {
	var params entity.DetailMPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("MPriceController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)

	// log.Println("MPriceController, Detail, CustId:", custId)

	data, err := controller.MPriceService.Detail(params)
	if err != nil {
		log.Error("MPriceController, Detail, FindOneByMPriceId, err:", err.Error())
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

func (controller *MPriceController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.MPriceQueryFilter
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
		log.Error("MPriceController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	statuses, err := parseMPriceStatusQuery(c.Context().QueryArgs())
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(statuses) > 0 {
		dataFilter.Status = statuses
	}

	distributorIDs, err := parseMPriceDistributorIDQuery(c.Context().QueryArgs())
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorIDs = distributorIDs
	}

	custId := c.Locals("cust_id").(string)

	data, total, lastPage, err = controller.MPriceService.List(dataFilter, custId)
	if err != nil {
		log.Error("MPriceController, List, data, err:", err.Error())
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

func (controller *MPriceController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var request entity.CreateMPriceBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("MPriceController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.CustID = c.Locals("cust_id").(string)
	request.CreatedBy = c.Locals("user_fullname").(string)
	userID := c.Locals("user_id").(int64)
	request.CreatedByID = &userID
	request.DistributorID = c.Locals("distributor_id").(int64)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.MPriceService.Store(request)
	if err != nil {
		log.Error("MPriceController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateMPriceParams
		request entity.UpdateMPriceRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("MPriceController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if err := c.BodyParser(&request); err != nil {
		log.Error("MPriceController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	request.CustID = c.Locals("cust_id").(string)
	request.ParentCustID = c.Locals("parent_cust_id").(string)
	request.UpdatedBy = c.Locals("user_fullname").(string)
	userID := c.Locals("user_id").(int64)
	request.UpdatedByID = &userID
	request.DistributorID = c.Locals("distributor_id").(int64)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.MPriceService.Update(params.PriceID, request)
	if err != nil {
		log.Error("MPriceController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteMPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("MPriceController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	// log.Println("MPriceController, Delete, CustId:", custId)

	err := controller.MPriceService.Delete(custId, params.PriceId, userId)
	if err != nil {
		log.Error("MPriceController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) GetStatuses(c *fiber.Ctx) error {
	mpriceStatuses := make([]entity.MPriceStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.MPriceStatusDesc {
		status := entity.MPriceStatus{
			StatusID:   index,
			StatusDesc: element,
		}
		mpriceStatuses = append(mpriceStatuses, status)
	}

	mpriceStatusesSorted := make(entity.MPriceStatusDescSlice, 0)
	for _, row := range mpriceStatuses {
		mpriceStatusesSorted = append(mpriceStatusesSorted, row)
	}
	sort.Sort(mpriceStatusesSorted)
	responsePayload.Setdata(mpriceStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) Cancel(c *fiber.Ctx) error {
	var (
		params entity.CancelMPriceParams
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("MPriceController, Cancel, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UpdatedBy = c.Locals("user_fullname").(string)
	userID := c.Locals("user_id").(int64)
	params.UpdatedByID = &userID

	if err := controller.MPriceService.Cancel(params); err != nil {
		log.Error("MPriceController, Cancel, Service.Cancel, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_CANCELLED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) Publish(c *fiber.Ctx) error {
	var params entity.PublishMPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("MPriceController, Publish, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("MPriceController, Publish, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.DistributorID = c.Locals("distributor_id").(int64)
	params.UpdatedBy = c.Locals("user_fullname").(string)
	userID := c.Locals("user_id").(int64)
	params.UpdatedByID = &userID

	if err := controller.MPriceService.Publish(params); err != nil {
		log.Error("MPriceController, Publish, Service.Publish, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Published Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) Template(c *fiber.Ctx) error {
	buffer, contentType, filename, err := controller.MPriceService.Template(
		c.Query("format"),
		c.Locals("cust_id").(string),
		c.Locals("parent_cust_id").(string),
		c.Locals("distributor_id").(int64),
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *MPriceController) Export(c *fiber.Ctx) error {
	var dataFilter entity.MPriceQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	statuses, err := parseMPriceStatusQuery(c.Context().QueryArgs())
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(statuses) > 0 {
		dataFilter.Status = statuses
	}

	distributorIDs, err := parseMPriceDistributorIDQuery(c.Context().QueryArgs())
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	if len(distributorIDs) > 0 {
		dataFilter.DistributorIDs = distributorIDs
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	buffer, contentType, filename, err := controller.MPriceService.Export(dataFilter, custID, parentCustID)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(buffer.Bytes())
}

func (controller *MPriceController) Import(c *fiber.Ctx) error {
	var (
		request entity.MPriceImportRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userID := c.Locals("user_id").(int64)
	distributorID := c.Locals("distributor_id").(int64)
	userFullName := c.Locals("user_fullname").(string)

	data, err := controller.MPriceService.Import(request, custID, parentCustID, userID, distributorID, userFullName)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		responsePayload.Setdata([]entity.MPriceImportResponse{data})
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata([]entity.MPriceImportResponse{data})
	responsePayload.Setmsg("Price data upload processed successfully")
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *MPriceController) processPublishManagePriceMessage(msg amqp.Delivery) {
	log.Infof("Processing message: %s", string(msg.Body))

	// Step 1: Unmarshal the message body into a JSON object
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(msg.Body, &jsonBody); err != nil {
		log.Errorf("Failed to unmarshal JSON body: %v", err)
		// Optionally reject the message if unmarshalling fails
		msg.Nack(false, false) // Requeue: false, multiple: false
	}

	// Step 2: Map the JSON body to the request struct
	var request entity.PublishByRmqMPriceReq
	if err := structs.Automapper(jsonBody, &request); err != nil {
		log.Errorf("Failed to map JSON body to request struct: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Step 3: Call the service layer to process the message
	if err := controller.MPriceService.PublishByRMQ(request); err != nil {
		log.Errorf("Failed to process message in service: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Acknowledge the message after successful processing
	if err := msg.Ack(true); err != nil {
		log.Errorf("Failed to acknowledge message: %v", err)
	}

	log.Infof("Message processed successfully: %s", string(msg.Body))
}

func parseMPriceStatusQuery(args *fasthttp.Args) ([]int, error) {
	return parseIntSliceQuery(args, "status", "status", "status[]")
}

func parseMPriceDistributorIDQuery(args *fasthttp.Args) ([]int64, error) {
	values, err := parseIntSliceQuery(args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}

	result := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		result = append(result, int64(value))
	}
	return result, nil
}
