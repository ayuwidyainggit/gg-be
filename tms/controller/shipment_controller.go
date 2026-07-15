package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ShipmentController struct {
	shipmentService service.ShipmentService
}

func NewShipmentController(service service.ShipmentService) *ShipmentController {
	return &ShipmentController{
		shipmentService: service,
	}
}

// TODO
// Note             godoc
//
//	@Summary		Create manual shipment
//	@Description	Save shipment data to db.
//	@Param			shipment	body	entity.CreateShipmentRequest	true	"Create manual shipment and example value delivery_date = 2006-01-02 (yyyy-mm-dd)"
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		201	{object}	entity.JsonCreated{data=string}"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/web/shipments [post]
//
// @Security	Bearer
func (controller *ShipmentController) CreateManual(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 60*time.Second)
	defer cancel()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	request := entity.CreateShipmentRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	data, error := controller.shipmentService.CreateManual(c, headers, request)
	helper.ErrorPanic(error)

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Manual Successfully",
		Data:    data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}

// TODO
// Note             godoc
//
//	@Summary		Create auto shipment
//	@Description	catatan untuk attr longitude dan latitude itu harus numeric!!!.
//	@Param			shipment	body	entity.CreateShipmentAutoRequest	true	"Create auto shipment and example value delivery_date = 2006-01-02 (yyyy-mm-dd)"
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		201	{object}	entity.JsonCreated{data=[]string}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/web/shipments/auto [post]
//
// @Security	Bearer
func (controller *ShipmentController) CreateAuto(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 60*time.Second)
	defer cancel()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	request := entity.CreateShipmentAutoRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	data, error := controller.shipmentService.CreateAuto(c, headers, request)
	if error != nil {
		return exception.NewInternalServerError(error.Error())
	}

	webResponse := entity.Response{
		Code:    http.StatusCreated,
		Status:  "Created",
		Message: "Created Auto Successfully",
		Data:    data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Submit shipment preview
//	@Description	Update route_id in Db.
//	@Param			shipment	body	entity.SubmitShipmentRequest	true	"Submit shipment preview"
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		200	{object}	entity.JsonSuccess{data=nil}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}																"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}																"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}													"Internal server error"
//	@Router			/web/shipments/submit [patch]
func (controller *ShipmentController) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.SubmitShipmentRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.shipmentService.SubmitShipment(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Submit Shipment Preview Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		Get shipment by shipment_no.
//	@Param			shipmentNo	path	string	true	"ShipmentNo"
//	@Description	Return the shipments.
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentPreviewResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipments/{shipmentNo} [get]
func (controller *ShipmentController) FindByShipmentNo(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.ShipmentParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.shipmentService.FindByShipmentNo(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note 		        godoc
//
//	@Summary		Get shipment by order_no.
//	@Param			shipmentNo	path	string	true	"ShipmentNo"
//	@Description	Return the shipments.
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentPickList{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipments/invoices/{shipmentNo} [get]
func (controller *ShipmentController) FindByOrderNo(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.ShipmentParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.shipmentService.FindShipmentInvoiceByShipmentNo(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get all shipments.
//	@Description	Return list of shipments.
//	@Produce		application/json
//	@Param			shipment_no		query	string	false	"shipment_no"
//	@Param			cust_id			query	string	false	"cust_id"
//	@Param			delivery_date	query	string	false	"delivery_date"
//	@Param			driver_id		query	string	false	"driver_id"
//	@Param			vehicle_id		query	string	false	"vehicle_id"
//	@Param			driver_name		query	string	false	"driver_name"
//	@Param			outlet_name		query	string	false	"outlet_name"
//	@Param			start_date		query	string	false	"start_date"
//	@Param			end_date		query	string	false	"end_date"
//	@Param			sort			query	string	false	"sort"
//	@Tags			web(shipment)
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.ShipmentResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}									"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}									"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}						"Internal server error"
//	@Router			/web/shipments [get]
func (controller *ShipmentController) FindAll(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentQueryFilter

	dataFilter.DriverName = ctx.Query("driver_name")
	dataFilter.DriverID, _ = strconv.Atoi(ctx.Query("driver_id"))
	dataFilter.VehicleID, _ = strconv.Atoi(ctx.Query("vehicle_id"))
	dataFilter.DeliveryDate = ctx.Query("delivery_date")
	dataFilter.OutletName = ctx.Query("outlet_name")
	dataFilter.StartDate = ctx.Query("start_date")
	dataFilter.ShipmentNo = ctx.Query("shipment_no")
	dataFilter.CustID = ctx.Query("cust_id")
	dataFilter.EndDate = ctx.Query("end_date")
	dataFilter.Sort = ctx.Query("sort")

	response := controller.shipmentService.FindAll(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note		        godoc
//
//	@Summary		Delete shipment.
//	@Param			shipmentNo	path	string	true	"ShipmentNo"
//	@Description	Remove shipment data by shipment_no.
//	@Produce		application/json
//	@Tags			web(shipment)
//	@Success		200	{object}	entity.JsonSuccess{data=nil}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/web/shipments/{shipmentNo} [delete]
// @Security	Bearer
func (controller *ShipmentController) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "empty token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	var params entity.ShipmentParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	controller.shipmentService.Delete(c, headers, params)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Shipment Deleted Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
// @Summary		    Delete bulk shipment
// @Description	    delete bulk shipment.
// @Param			shipment	body	entity.DeleteShipmentRequest	true	"delete bulk shipment"
// @Produce		    application/json
// @Tags			web(shipment)
// @Success		    200	{object}	entity.JsonSuccess{data=nil}	"Data"
// @Failure		    400	{object}	entity.JsonBadRequest{}					"Validation error"
// @Failure		    404	{object}	entity.JsonNotFound{}					"Data not found"
// @Failure		    500	{object}	entity.JsonInternalServerError{}		"Internal server error"
// @Router			/web/shipments/bulk [post]
func (controller *ShipmentController) DeleteBulk(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.DeleteShipmentRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	controller.shipmentService.DeleteBulk(c, request)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Delete Bulk Successfully",
		Data:    nil,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		login sendpick get token jwt
//	@Description	login sendpick get token jwt.
//	@Produce		application/json
//	@Tags			sendpick(thirdparty)
//	@Success		200	{object}	entity.JsonSuccess{data=string}		"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}				"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}				"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}	"Internal server error"
//	@Router			/login/send-pick [post]
func (controller *ShipmentController) LoginSendPick(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data, error := controller.shipmentService.LoginSendPick(c)
	helper.ErrorPanic(error)

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Login SendPick Successfully",
		Data:    data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}

// Note             godoc
//
// @Summary		    send to sendpick with result shipment_no and order_id
// @Description	    catatan untuk attr outlet_latitude, outlet_longitude, warehouse_latitude, warehouse_longitude wajib diisi dgn numeric/integer, contoh : -6.1365484 dan gk boleh kosong!!!!!!!!.
// @Param		    shipment	body	entity.CreateShipmentAutoRequest	true	"send pick auto shipment"
// @Produce		    application/json
// @Tags		    sendpick(thirdparty)
// @Success		    200	{object}	entity.JsonSuccess{data=[]entity.SendPickResponse}	"Data"
// @Failure		    400	{object}	entity.JsonBadRequest{}								"Validation error"
// @Failure	    	404	{object}	entity.JsonNotFound{}								"Data not found"
// @Failure		    500	{object}	entity.JsonInternalServerError{}					"Internal server error"
// @Router			/generate/send-pick [post]
func (controller *ShipmentController) GenerateSendPick(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.CreateShipmentAutoRequest{}
	err := ctx.BodyParser(&request)
	helper.ErrorPanic(err)

	data, error := controller.shipmentService.GenerateSendPick(c, request)
	if error != nil {
		return exception.NewInternalServerError(error.Error())
	}

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusCreated).JSON(webResponse)
}
