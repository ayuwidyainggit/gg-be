package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"strconv"
	"time"
)

type OutletController struct {
	outletService service.OutletService
}

func NewOutletController(outletService service.OutletService) *OutletController {
	return &OutletController{
		outletService: outletService,
	}
}

// Note 		        godoc
//
//	@Summary		Get outlets.
//	@Param			shipment_no	query	string	false	"shipment_no"
//	@Param			cust_id		query	string	false	"cust_id"
//	@Param			driver_id	query	string	false	"driver_id"
//	@Param			product_id	query	string	false	"product_id"
//	@Param			sort		query	string	false	"sort"
//	@Description	Return the outlets.
//	@Produce		application/json
//	@Tags			outlet
//	@Success		200	{object}	entity.JsonSuccess{data=[]entity.OutletResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/outlets [get]
func (controller *OutletController) GetOutlet(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentInvoicesQueryFilter

	dataFilter.DriverID, _ = strconv.Atoi(ctx.Query("driver_id"))
	dataFilter.ProductID, _ = strconv.Atoi(ctx.Query("product_id"))
	dataFilter.ShipmentNo = ctx.Query("shipment_no")
	dataFilter.CustID = ctx.Query("cust_id")
	dataFilter.Sort = ctx.Query("sort")

	data := controller.outletService.GetOutlet(c, dataFilter)

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
//	@Summary		Get outlet by driver_id, outlet_id and shipment_no.
//	@Param			driverId	path	string	true	"Driver ID"
//	@Param			outletId	path	string	true	"Outlet ID"
//	@Param			shipmentNo	path	string	true	"ShipmentNo"
//	@Description	Return the outlet.
//	@Produce		application/json
//	@Tags			outlet
//	@Success		200	{object}	entity.JsonSuccess{data=entity.OutletResponse}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}							"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}							"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}				"Internal server error"
//	@Router			/mobile/outlet/{driverId}/{outletId}/{shipmentNo} [get]
func (controller *OutletController) GetOutletByParams(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var params entity.OutletParams

	if err := ctx.ParamsParser(&params); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.outletService.GetOutletByParams(c, params)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(fiber.StatusOK).JSON(webResponse)
}
