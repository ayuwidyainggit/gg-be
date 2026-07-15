package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

// Note             godoc
//
//	@Summary		get driver report.
//	@Description	Return driver report.
//	@Param			driver_id		query	string	true	"driver_id"
//	@Param			period		    query	string	true	"period (today/month)"
//	@Produce		application/json
//	@Tags			report
//	@Success		200	{object}	entity.JsonSuccess{data=entity.DriverReportResponse{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/mobile/driver/reports [get]
func (controller *ReportController) GetDriverReport(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.DriverReportQueryFilter

	if err := ctx.QueryParser(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.reportService.GetDriverReport(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// TODO Shipment Report

// Note             godoc
//
//	@Summary		get shipment report summary.
//	@Description	Return shipment report summary.
//	@Param			start_date		query	string	false	"start_date"
//	@Param			end_date		query	string	false	"end_date"
//	@Param			shipment_no		query	string	false	"shipment_no"
//	@Param			driver_name		query	string	false	"driver_name"
//	@Produce		application/json
//	@Tags			web(shipment_report)
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentReportSummary{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/summary [get]
func (controller *ReportController) GetShipmentReportSummary(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentReportQueryFilter

	if err := ctx.QueryParser(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	dataFilter.StartDate = ctx.Query("start_date")
	dataFilter.EndDate = ctx.Query("end_date")
	dataFilter.ShipmentNo = ctx.Query("shipment_no")
	dataFilter.DriverName = ctx.Query("driver_name")

	data := controller.reportService.GetShipmentReportSummary(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		get shipment report detail.
//	@Description	Return shipment report detail.
//	@Param			start_date		query	string	false	"start_date"
//	@Param			end_date		query	string	false	"end_date"
//	@Param			shipment_no		query	string	false	"shipment_no"
//	@Param			outlet		query	string	false	"outlet"
//	@Param			driver		query	string	false	"driver"
//	@Param			visited_status		query	string	false	"visited_status"
//	@Param			received_status		query	string	false	"received_status"
//	@Produce		application/json
//	@Tags			web(shipment_report)
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentReportDetail{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/detail [get]
func (controller *ReportController) GetShipmentReportDetail(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentReportDetailQueryFilter

	if err := ctx.QueryParser(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.reportService.GetShipmentReportDetail(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		get shipment report reject.
//	@Description	Return shipment report reject.
//	@Param			start_date		query	string	false	"start_date"
//	@Param			end_date		query	string	false	"end_date"
//	@Param			shipment_no		query	string	false	"shipment_no"
//	@Param			outlet		query	string	false	"outlet"
//	@Param			driver		query	string	false	"driver"
//	@Param			pcode		query	string	false	"pcode"
//	@Param			reason		query	string	false	"reason"
//	@Produce		application/json
//	@Tags			web(shipment_report)
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentReportReject{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/reject [get]
func (controller *ReportController) GetShipmentReportReject(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	var dataFilter entity.ShipmentReportRejectlQueryFilter

	if err := ctx.QueryParser(&dataFilter); err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	data := controller.reportService.GetShipmentReportReject(c, dataFilter)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// TODO Dropdown [DONE]

// Note             godoc
//
//	@Summary		get shipment no drop down.
//	@Description	Return shipment_no dropdown.
//	@Produce		application/json
//	@Tags			web(shipment_report)-dropdown
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ShipmentNoDropdown{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/shipment [get]
func (controller *ReportController) GetShipmentNumberDropdown(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data := controller.reportService.GetListShipmentNo(c)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// @Summary		get product code dropdown.
// @Description	Return product code dropdown.
// @Produce		application/json
// @Tags			web(shipment_report)-dropdown
// @Success		200	{object}	entity.JsonSuccess{data=entity.ProductCodeDropdown{}}	"Data"
// @Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
// @Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
// @Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
// @Router			/web/shipment-report/product [get]
func (controller *ReportController) GetProductCodeDropdown(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data := controller.reportService.GetListProductCode(c)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		get driver name dropdown.
//	@Description	Return driver_name dropdown.
//	@Produce		application/json
//	@Tags			web(shipment_report)-dropdown
//	@Success		200	{object}	entity.JsonSuccess{data=entity.DriverNameDropdown{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/driver [get]
func (controller *ReportController) GetDriverDropdown(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data := controller.reportService.GetListDriver(c)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		get outlet dropdown.
//	@Description	Return outlet (outlet_name & outlet_code) dropdown.
//	@Produce		application/json
//	@Tags			web(shipment_report)-dropdown
//	@Success		200	{object}	entity.JsonSuccess{data=entity.OutletDropdown{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/outlet [get]
func (controller *ReportController) GetOutletDropdown(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data := controller.reportService.GetListOutlet(c)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		get reason dropdown.
//	@Description	Return reason for dropdown.
//	@Produce		application/json
//	@Tags			web(shipment_report)-dropdown
//	@Success		200	{object}	entity.JsonSuccess{data=entity.ReasonDropdown{}}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/web/shipment-report/reason [get]
func (controller *ReportController) GetReasonDropdown(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	data := controller.reportService.GetListReasons(c)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
