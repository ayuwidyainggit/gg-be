package controller

import (
	"errors"
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/middleware"
	"finance/pkg/responsebuild"
	"finance/pkg/validation"
	"finance/service"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ExpenseController struct {
	ExpenseService service.ExpenseService
	validator      *validation.Validate
}

func NewExpenseController(expenseService service.ExpenseService, validator *validation.Validate) *ExpenseController {
	return &ExpenseController{
		ExpenseService: expenseService,
		validator:      validator,
	}
}

func (controller *ExpenseController) Route(app *fiber.App) {
	expenseRouteV1 := app.Group("/v1/expense-type", middleware.JWTProtected())
	expenseRouteV1.Get("", controller.List)
	expenseRouteV1.Post("", controller.Create)
	expenseRouteV1.Patch("/:expense_type_id", controller.Update)
	expenseRouteV1.Delete("/:expense_type_id", controller.Delete)
}

func (controller *ExpenseController) getAcceptLanguage(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}

func (controller *ExpenseController) List(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var dataFilter entity.ExpenseQueryFilter
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("ExpenseController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// Set default values
	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseController, List, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.ExpenseService.List(dataFilter)
	if err != nil {
		log.Error("ExpenseController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Handle empty data
	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setdata(data)
	}

	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseController) Create(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var request entity.CreateExpenseBody
	if err := c.BodyParser(&request); err != nil {
		log.Error("ExpenseController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	request.CustId = c.Locals("cust_id").(string)
	request.ParentCustId = c.Locals("parent_cust_id").(string)

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseController, Create, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ExpenseService.Store(request, userId)
	if err != nil {
		log.Error("ExpenseController, Create, Store, err:", err.Error())
		if errors.Is(err, service.ErrExpenseTypeExists) {
			responsePayload.Setmsg("Data already exists")
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data saved successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseController) Update(c *fiber.Ctx) error {
	var (
		params  entity.UpdateExpenseParams
		request entity.UpdateExpenseBody
	)
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("ExpenseController, Update, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseController, Update, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&request); err != nil {
		log.Error("ExpenseController, Update, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	errs = controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseController, Update, ValidateStruct(request), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.ExpenseService.Update(parentCustId, params.ExpenseTypeID, request, userId)
	if err != nil {
		log.Error("ExpenseController, Update, Service.Update, err:", err.Error())
		if errors.Is(err, service.ErrExpenseTypeExists) {
			responsePayload.Setmsg("Data already exists")
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		if errors.Is(err, service.ErrExpenseTypeNotFound) {
			responsePayload.Setmsg("Record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to update data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data updated successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseController) Delete(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.DeleteExpenseParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ExpenseController, Delete, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseController, Delete, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	err := controller.ExpenseService.Delete(parentCustId, params.ExpenseTypeID, userId)
	if err != nil {
		log.Error("ExpenseController, Delete, Service.Delete, err:", err.Error())
		if errors.Is(err, service.ErrExpenseTypeNotFound) {
			responsePayload.Setmsg("Record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to delete data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data deleted successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
