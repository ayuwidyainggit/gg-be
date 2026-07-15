package controller

import (
	"errors"
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/middleware"
	"finance/pkg/responsebuild"
	"finance/pkg/validation"
	"finance/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type ExpenseEntryController struct {
	ExpenseService service.ExpenseEntryService
	validator      *validation.Validate
}

func NewExpenseEntryController(expenseService service.ExpenseEntryService, validator *validation.Validate) *ExpenseEntryController {
	return &ExpenseEntryController{ExpenseService: expenseService, validator: validator}
}

func (controller *ExpenseEntryController) Route(app *fiber.App) {
	expenseRouteV1 := app.Group("/v1/expense", middleware.JWTProtected())
	expenseRouteV1.Get("", controller.List)
	expenseRouteV1.Post("", controller.Create)
	expenseRouteV1.Get("/:expense_id", controller.Detail)
	expenseRouteV1.Patch("/:expense_id", controller.Update)
	expenseRouteV1.Delete("/:expense_id", controller.Delete)
}

func (controller *ExpenseEntryController) getAcceptLanguage(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}

func (controller *ExpenseEntryController) List(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var filter entity.ExpenseEntryQueryFilter
	if err := c.QueryParser(&filter); err != nil {
		log.Error("ExpenseEntryController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	collectorIDsRaw := c.Query("collector_id")
	if collectorIDsRaw != "" {
		collectorIDs := make([]int64, 0)
		for _, collectorID := range strings.Split(collectorIDsRaw, ",") {
			collectorID = strings.TrimSpace(collectorID)
			if collectorID == "" {
				continue
			}

			collectorIDInt, err := strconv.ParseInt(collectorID, 10, 64)
			if err != nil {
				responsePayload.Setmsg(fiber.ErrBadRequest.Message)
				responsePayload.Seterrors("collector_id must be integer list")
				return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
			}

			collectorIDs = append(collectorIDs, collectorIDInt)
		}
		filter.CollectorIDs = collectorIDs
	}

	filter.CustID = c.Locals("cust_id").(string)
	filter.ParentCustID = c.Locals("parent_cust_id").(string)
	filter.UserID = c.Locals("user_id").(int64)

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Sort == "" {
		filter.Sort = "created_date:desc"
	}

	errs := controller.validator.ValidateStruct(filter, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseEntryController, List, ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, total, lastPage, err := controller.ExpenseService.List(filter)
	if err != nil {
		log.Error("ExpenseEntryController, List, Service.List:", err.Error())
		if errors.Is(err, service.ErrInvalidExpenseDateFilter) {
			responsePayload.Setmsg(fiber.ErrBadRequest.Message)
			responsePayload.Seterrors(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("No Data")
		responsePayload.Setdata(nil)
	} else {
		responsePayload.Setdata(data)
	}

	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: filter.Page,
		PageLimit:   filter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseEntryController) Create(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var req entity.CreateExpenseEntryBody
	if err := c.BodyParser(&req); err != nil {
		log.Error("ExpenseEntryController, Create, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	custId := c.Locals("cust_id").(string)

	errs := controller.validator.ValidateStruct(req, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseEntryController, Create, ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ExpenseService.Store(custId, req, userId)
	if err != nil {
		log.Error("ExpenseEntryController, Create, Service.Store:", err.Error())
		responsePayload.Setmsg("Failed to save data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Expense created successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseEntryController) Detail(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.ExpenseEntryParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ExpenseEntryController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error("ExpenseEntryController, Detail, ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	data, err := controller.ExpenseService.Detail(custId, params.ExpenseID)
	if err != nil {
		log.Error("ExpenseEntryController, Detail, Service.Detail:", err.Error())
		if errors.Is(err, service.ErrExpenseNotFound) {
			responsePayload.Setmsg("Record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to fetch data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseEntryController) Update(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.ExpenseEntryParams
	var req entity.UpdateExpenseEntryBody
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ExpenseEntryController, Update, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		log.Error("ExpenseEntryController, Update, ValidateStruct(params):", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if err := c.BodyParser(&req); err != nil {
		log.Error("ExpenseEntryController, Update, BodyParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if errs := controller.validator.ValidateStruct(req, headerAcceptLang); errs != nil {
		log.Error("ExpenseEntryController, Update, ValidateStruct(req):", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	custId := c.Locals("cust_id").(string)

	data, err := controller.ExpenseService.Update(custId, params.ExpenseID, req, userId)
	if err != nil {
		log.Error("ExpenseEntryController, Update, Service.Update:", err.Error())
		if errors.Is(err, service.ErrExpenseNotFound) {
			responsePayload.Setmsg("Record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to update data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Expense updated successfully")
	responsePayload.Setdata(data)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseEntryController) Delete(c *fiber.Ctx) error {
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	var params entity.ExpenseEntryParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error("ExpenseEntryController, Delete, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	if errs := controller.validator.ValidateStruct(params, headerAcceptLang); errs != nil {
		log.Error("ExpenseEntryController, Delete, ValidateStruct:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)
	custId := c.Locals("cust_id").(string)

	if err := controller.ExpenseService.Delete(custId, params.ExpenseID, userId); err != nil {
		log.Error("ExpenseEntryController, Delete, Service.Delete:", err.Error())
		if errors.Is(err, service.ErrExpenseNotFound) {
			responsePayload.Setmsg("Record not found")
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg("Failed to delete data, please try again")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Data deleted successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
