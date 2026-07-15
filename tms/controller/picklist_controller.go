package controller

import (
	"context"
	"net/http"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/service"
	"scyllax-tms/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PicklistController struct {
	picklistService service.PicklistService
}

func NewPicklistController(picklistService service.PicklistService) *PicklistController {
	return &PicklistController{
		picklistService: picklistService,
	}
}

// CreatePicklist godoc
//
// @Summary      Create a new picklist
// @Description  Create a new picklist and save it to the database.
// @Param        data  body  entity.CreatePicklistRequest  true  "Create Picklist"
// @Produce      application/json
// @Tags         picklist
// @Success      200  {object}  entity.JsonSuccess{data=entity.PicklistResponse}  "Data"
// @Failure      400  {object}  entity.JsonBadRequest{}                           "Validation error"
// @Failure      500  {object}  entity.JsonInternalServerError{}                  "Internal server error"
// @Router       /picklists [post]
//
// @security 	Bearer
func (controller *PicklistController) CreatePicklist(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	customerId, err := helper.DecodeBearerToken(ctx)
	if err != nil {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "invalid or expired token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	request := entity.CreatePicklistRequest{}
	err = ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := controller.picklistService.Create(c, request, customerId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Picklist created successfully",
		Data:    response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// UpdatePicklist godoc
//
// @Summary      Update an existing picklist
// @Description  Update an existing picklist in the database.
// @Param        data  body  entity.UpdatePicklistRequest  true  "Update Picklist"
// @Produce      application/json
// @Tags         picklist
// @Success      200  {object}  entity.JsonSuccess{data=entity.PicklistResponse}  "Data"
// @Failure      400  {object}  entity.JsonBadRequest{}                           "Validation error"
// @Failure      500  {object}  entity.JsonInternalServerError{}                  "Internal server error"
// @Router       /picklists [put]
func (controller *PicklistController) UpdatePicklist(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.UpdatePicklistRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := controller.picklistService.Update(c, request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Picklist updated successfully",
		Data:    response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// DeletePicklist godoc
//
// @Summary      Delete a picklist
// @Description  Delete a picklist from the database.
// @Param        data  body  entity.DeletePicklistRequest  true  "Delete Picklist"
// @Produce      application/json
// @Tags         picklist
// @Success      200  {object}  entity.JsonSuccess{data=entity.PicklistResponse}  "Data"
// @Failure      400  {object}  entity.JsonBadRequest{}                           "Validation error"
// @Failure      500  {object}  entity.JsonInternalServerError{}                  "Internal server error"
// @Router       /picklists [delete]
func (controller *PicklistController) DeletePicklist(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.DeletePicklistRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	response, err := controller.picklistService.Delete(c, request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Picklist deleted successfully",
		Data:    response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// GetPicklist godoc
//
// @Summary      Get a picklist
// @Description  Get a picklist from the database.
// @Param        id  path  string  true  "Picklist ID"
// @Produce      application/json
// @Tags         picklist
// @Success      200  {object}  entity.JsonSuccess{data=entity.PicklistResponse}  "Data"
// @Failure      400  {object}  entity.JsonBadRequest{}                           "Validation error"
// @Failure      404  {object}  entity.JsonNotFound{}                             "Data not found"
// @Failure      500  {object}  entity.JsonInternalServerError{}                  "Internal server error"
// @Router       /picklists/{id} [get]
func (controller *PicklistController) GetPicklist(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	request := entity.GetPicklistRequest{
		PicklistNo: ctx.Params("id"),
	}

	response, err := controller.picklistService.GetPicklist(c, request)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	webResponse := entity.Response{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "Picklist retrieved successfully",
		Data:    response,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// GetAllPicklists godoc
//
// @Summary      Get all picklists
// @Description  Get all picklists from the database.
// @Produce      application/json
//
//	@Param		limit		query	string	false	"Limit"
//	@Param		page		query	string	false	"Page"
//	@Param		driver		query	string	false	"Driver"
//	@Param		vehicle		query	string	false	"Vehicle"
//	@Param		start_date		query	string	false	"Start Date"
//	@Param		end_date		query	string	false	"End Date"
//
// @Tags         picklist
// @Success      200  {object}  entity.Response{data=[]entity.PicklistResponse,meta=entity.Meta}  "Data"
// @Failure      400  {object}  entity.JsonBadRequest{}                             "Validation error"
// @Failure      500  {object}  entity.JsonInternalServerError{}                    "Internal server error"
// @Router       /picklists [get]
//
// @security 	Bearer
func (controller *PicklistController) GetAllPicklists(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	defer cancel()

	customerId, err := helper.DecodeBearerToken(ctx)
	if err != nil {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "invalid or expired token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}

	var dataFilter entity.GeneralQueryFilter
	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	dataFilter.Page, _ = strconv.Atoi(ctx.Query("page"))

	var picklistFilter entity.PicklistFilter
	picklistFilter.Driver = ctx.Query("driver")
	picklistFilter.Vehicle = ctx.Query("vehicle")
	picklistFilter.StartDate = ctx.Query("start_date")
	picklistFilter.EndDate = ctx.Query("end_date")

	response, paging, err := controller.picklistService.GetAll(c, dataFilter, picklistFilter, customerId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}

// Note             godoc
//
//	@Summary		Get list of returns.
//	@Description	Return list of returns.
//	@Produce		application/json
//	@Param			limit		query	string	false	"Limit"
//	@Param			outlet_id   query	string	false	"OutletId"
//	@Param			document_type   query	string	false	"DocumentType"
//	@Param			emp_id   query	string	false	"EmpId"
//	@Tags			picklist
//	@Success		200	{object}	entity.Response{data=[]entity.CustomShipmentInvoice}	"Data"
//	@Failure		400	{object}	entity.JsonBadRequest{}								"Validation error"
//	@Failure		401	{object}	entity.JsonUnauthorized{}								"Unauthorized"
//	@Failure		404	{object}	entity.JsonNotFound{}								"Data not found"
//	@Failure		500	{object}	entity.JsonInternalServerError{}					"Internal server error"
//	@Router			/picklists/invoices [get]
//
// @Security	Bearer
func (controller *PicklistController) GetPicklistInvoice(ctx *fiber.Ctx) error {
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

	customerId, err := helper.DecodeBearerToken(ctx)
	if err != nil {
		webResponse := entity.Response{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Message: "invalid or expired token",
		}
		utils.ResponseInterceptor(c, &webResponse)
		return ctx.Status(http.StatusUnauthorized).JSON(webResponse)
	}



	var dataFilter entity.GeneralQueryFilter

	dataFilter.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	outletId := ctx.Query("outlet_id")
	documentType := ctx.Query("document_type")
	empId := ctx.Query("emp_id")
	dataFilter.OutletId = outletId
	dataFilter.DocumentType = documentType
	dataFilter.EmpId = empId

	headers := make(map[string]string)
	headers["Authorization"] = ctx.Get("Authorization")
	headers["Accept"] = ctx.Get("application/json")

	response, paging, err := controller.picklistService.GetListInvoice(c, dataFilter, headers, customerId)
	helper.ErrorPanic(err)

	webResponse := entity.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   response,
		Meta:   &paging,
	}
	utils.ResponseInterceptor(c, &webResponse)
	return ctx.Status(http.StatusOK).JSON(webResponse)
}
