package controller

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"

	"github.com/gofiber/fiber/v2"
)

type SalesTargetDistributorController struct {
	SalesTargetDistributorService service.SalesTargetDistributorService
	SalesTargetService            service.SalesTargetService
	validator                     *validation.Validate
}

// NewSalesTargetDistributorController creates a new instance of SalesTargetDistributorController
func NewSalesTargetDistributorController(salesTargetDistributorService service.SalesTargetDistributorService, salesTargetService service.SalesTargetService, validator *validation.Validate) SalesTargetDistributorController {
	return SalesTargetDistributorController{
		SalesTargetDistributorService: salesTargetDistributorService,
		SalesTargetService:            salesTargetService,
		validator:                     validator,
	}
}

// Route registers the routes for Sales Target Distributor feature
func (controller *SalesTargetDistributorController) Route(app *fiber.App) {
	route := app.Group("/v1/sales-target-distributor", middleware.JWTProtected())
	route.Get("", controller.List)
	route.Post("", controller.Create)
	route.Get("/:sales_target_distributor_yearly_id", controller.Detail)
	route.Patch("/:sales_target_distributor_yearly_id", controller.Update)
}

// List handles the request to retrieve a list of yearly sales targets for a distributor
func (controller *SalesTargetDistributorController) List(c *fiber.Ctx) error {
	var (
		err        error
		dataFilter entity.SalesTargetDistributorQueryFilter
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
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 20
	}
	if dataFilter.Limit > 9999 {
		dataFilter.Limit = 9999
	}
	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}

	custId := c.Locals("parent_cust_id").(string)
	data, total, lastPage, err = controller.SalesTargetDistributorService.List(dataFilter, custId)
	if err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
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

// Detail handles the request to retrieve details of a specific yearly sales target, including monthly targets
func (controller *SalesTargetDistributorController) Detail(c *fiber.Ctx) error {
	var params entity.DetailSalesTargetDistributorParams
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	data, err := controller.SalesTargetDistributorService.Detail(params.SalesTargetDistributorYearlyId, custId)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		errMsg := fiber.ErrBadRequest.Message
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = fiber.StatusNotFound
			errMsg = "record not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

// Create handles the request to add a new yearly sales target with its monthly details
func (controller *SalesTargetDistributorController) Create(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateSalesTargetDistributorBody
	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	_, err := controller.SalesTargetDistributorService.Store(request)
	if err != nil {
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

// Update handles the request to modify an existing yearly sales target and its monthly details
func (controller *SalesTargetDistributorController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateSalesTargetDistributorParams
		request entity.UpdateSalesTargetDistributorRequest
	)
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(fiber.ErrUnprocessableEntity.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("parent_cust_id").(string)
	userId := c.Locals("user_id").(int64)

	request.CustId = custId
	request.UpdatedBy = userId

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.SalesTargetDistributorService.Update(params.SalesTargetDistributorYearlyId, request)
	if err != nil {
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
