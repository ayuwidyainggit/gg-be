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
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type SpPriceController struct {
	SpPriceService service.SpService
	validator      *validation.Validate
}

func NewSpPriceController(spPrice service.SpService, validator *validation.Validate) *SpPriceController {
	return &SpPriceController{
		SpPriceService: spPrice,
		validator:      validator,
	}
}

func (controller *SpPriceController) Route(app *fiber.App) {
	// go rabbitmq.Subscribe(constant.RMQ_OUTLET_PRICE_START_EVENT, controller.processPublishOutletPriceMessage)
	// go rabbitmq.Subscribe(constant.RMQ_OUTLET_PRICE_END_EVENT, controller.processInactiveOutletPriceMessage)
	qParamId := ":sp_price_id"
	spPriceRouteV1Scheduler := app.Group("/v1/outlet-prices/scheduler")
	spPriceRouteV1Scheduler.Post("/publish-unpublish", controller.PublishOrUnpublish)
	spPriceRouteV1 := app.Group("/v1/outlet-prices", middleware.JWTProtected())
	spPriceRouteV1.Get("/statuses", controller.GetStatuses)
	spPriceRouteV1.Get("/"+qParamId, controller.Detail)
	spPriceRouteV1.Post("preview", controller.Preview)
	spPriceRouteV1.Post("", controller.Create)
	spPriceRouteV1.Get("", controller.List)
	spPriceRouteV1.Patch("/"+qParamId+"/cancel", controller.Cancel)
	spPriceRouteV1.Patch("/"+qParamId, controller.Update)
	spPriceRouteV1.Delete("/"+qParamId, controller.Delete)
}

func (controller *SpPriceController) GetStatuses(c *fiber.Ctx) error {
	mSpPriceStatuses := make([]entity.MSpPriceStatus, 0)

	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	for index, element := range entity.MSpPriceStatusDesc {
		status := entity.MSpPriceStatus{
			StatusID:   index,
			StatusDesc: element,
		}
		mSpPriceStatuses = append(mSpPriceStatuses, status)
	}

	mpriceStatusesSorted := make(entity.MSpPriceStatusDescSlice, 0)
	for _, row := range mSpPriceStatuses {
		mpriceStatusesSorted = append(mpriceStatusesSorted, row)
	}
	sort.Sort(mpriceStatusesSorted)
	responsePayload.Setdata(mpriceStatusesSorted)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateMSpPriceBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("SpPriceController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	userFullname := c.Locals("user_fullname").(string)

	request.CustID = custId
	request.CreatedBy = userFullname
	request.ParentCustID = parentCustId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SpPriceService.Store(request)
	if err != nil {
		log.Error("SpPriceController, Create, Store, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) List(c *fiber.Ctx) error {
	var dataFilter entity.MSpPriceQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("SpPriceController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)
	dataFilter.CustID = custId
	dataFilter.ParentCustID = parentCustId

	data, total, lastPage, err := controller.SpPriceService.List(dataFilter, custId)
	if err != nil {
		log.Error("SpPriceController, List, data, err:", err.Error())
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

func (controller *SpPriceController) Detail(c *fiber.Ctx) error {
	var params entity.MSpPriceParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("SpPriceController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)

	data, err := controller.SpPriceService.Detail(params)
	if err != nil {
		log.Error("SpPriceController, Detail, err:", err.Error())
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

func (controller *SpPriceController) Delete(c *fiber.Ctx) error {
	var params entity.MSpPriceDeleteParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("SpPriceController, Delete, ParamsParser, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)

	err := controller.SpPriceService.Delete(custID, params.SpPriceID)
	if err != nil {
		log.Error("SpPriceController, Delete, Service.Delete, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) Update(c *fiber.Ctx) error {
	var (
		params  entity.MSpPriceParams
		request entity.UpdateMSpPriceBody
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("SpPriceController, Update, ParamsParser(params):", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("SpPriceController, Update, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userFullname := c.Locals("user_fullname").(string)

	request.CustID = custId
	request.UpdatedBy = userFullname
	request.ParentCustID = parentCustID

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SpPriceService.Update(params.SpPriceID, request)
	if err != nil {
		log.Error("SpPriceController, Update, Service.Update, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) Preview(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateMSpPriceBody
	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	parentCustID := c.Locals("parent_cust_id").(string)
	userFullname := c.Locals("user_fullname").(string)

	request.CustID = custID
	request.ParentCustID = parentCustID
	request.CreatedBy = userFullname

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	response, err := controller.SpPriceService.Preview(request)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Data = response
	responsePayload.Setmsg(constant.SUCCESSFULLY_PREVIEWED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) Cancel(c *fiber.Ctx) error {
	var (
		params entity.MSpPriceCancelParams
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error(err)
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	log.Info("params:", structs.StructToJson(params))

	params.CustID = c.Locals("cust_id").(string)
	params.ParentCustID = c.Locals("parent_cust_id").(string)
	params.UpdatedBy = c.Locals("user_fullname").(string)

	err := controller.SpPriceService.Cancel(params)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *SpPriceController) processPublishOutletPriceMessage(msg amqp.Delivery) {
	log.Infof("Processing message: %s", string(msg.Body))

	// Step 1: Unmarshal the message body into a JSON object
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(msg.Body, &jsonBody); err != nil {
		log.Errorf("Failed to unmarshal JSON body: %v", err)
		// Optionally reject the message if unmarshalling fails
		msg.Nack(false, false) // Requeue: false, multiple: false
	}

	// Step 2: Map the JSON body to the request struct
	var request entity.PublishUnpublishSPriceReq
	if err := structs.Automapper(jsonBody, &request); err != nil {
		log.Errorf("Failed to map JSON body to request struct: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Step 3: Call the service layer to process the message
	if err := controller.SpPriceService.PublishOrInactive(request); err != nil {
		log.Errorf("Failed to process message in service: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Acknowledge the message after successful processing
	// if err := msg.Ack(true); err != nil {
	// 	log.Errorf("Failed to acknowledge message: %v", err)
	// }

	log.Infof("Message processed successfully: %s", string(msg.Body))
}

func (controller *SpPriceController) processInactiveOutletPriceMessage(msg amqp.Delivery) {
	log.Infof("Processing message: %s", string(msg.Body))

	// Step 1: Unmarshal the message body into a JSON object
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(msg.Body, &jsonBody); err != nil {
		log.Errorf("Failed to unmarshal JSON body: %v", err)
		// Optionally reject the message if unmarshalling fails
		msg.Nack(false, false) // Requeue: false, multiple: false
	}

	// Step 2: Map the JSON body to the request struct
	var request entity.PublishUnpublishSPriceReq
	if err := structs.Automapper(jsonBody, &request); err != nil {
		log.Errorf("Failed to map JSON body to request struct: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Step 3: Call the service layer to process the message
	if err := controller.SpPriceService.PublishOrInactive(request); err != nil {
		log.Errorf("Failed to process message in service: %v", err)
		// Optionally reject the message
		msg.Nack(false, false)
	}

	// Acknowledge the message after successful processing
	// if err := msg.Ack(true); err != nil {
	// 	log.Errorf("Failed to acknowledge message: %v", err)
	// }

	log.Infof("Message processed successfully: %s", string(msg.Body))
}

func (controller *SpPriceController) PublishOrUnpublish(c *fiber.Ctx) error {
	var (
		request entity.PublishUnpublishSPriceReq
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error("SpPriceController, PublishOrUnpublish, BodyParser(request), err:", err.Error())
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("SpPriceController, PublishOrUnpublish, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := controller.SpPriceService.PublishOrInactive(request); err != nil {
		log.Errorf("Failed to PublishOrUnpublish in service: %v", err)
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
