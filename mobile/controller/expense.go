package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"mobile/entity"
	"mobile/pkg/constant"
	"mobile/pkg/middleware"
	"mobile/pkg/responsebuild"
	"mobile/pkg/validation"
	"mobile/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	expenseRouteV1 := app.Group("/v1/expense", middleware.JWTProtected())
	expenseRouteV1.Get("", controller.List)
	expenseRouteV1.Get("/:expense_id", controller.Detail)
	expenseRouteV1.Post("/", controller.Create)
	expenseRouteV1.Patch("/:expense_id", controller.Update)
	expenseRouteV1.Delete("/:expense_id", controller.Delete)

	// Lookup endpoints
	expenseTypeRouteV1 := app.Group("/v1/expense_type", middleware.JWTProtected())
	expenseTypeRouteV1.Get("", controller.ExpenseTypeLookup)

	outletRouteV1 := app.Group("/v1/outlet", middleware.JWTProtected())
	outletRouteV1.Get("", controller.OutletLookup)
}

func (controller *ExpenseController) List(c *fiber.Ctx) error {
	var dataFilter entity.ExpenseQueryFilter
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	empID := c.Locals("emp_id").(int64)
	dataFilter.CustId = custId
	dataFilter.EmpID = empID

	// Set default pagination
	dataFilter.Limit, dataFilter.Page = controller.setPaginationDefaults(dataFilter.Limit, dataFilter.Page, 20)

	// Get user ID for attendance check
	userId := c.Locals("user_id").(int64)

	data, total, lastPage, err := controller.ExpenseService.List(dataFilter, custId, userId)
	if err != nil {
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

func (controller *ExpenseController) Detail(c *fiber.Ctx) error {
	var params entity.DetailExpenseParams
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	// Parse expense_id from string to int64
	expenseId, err := controller.parseExpenseID(params.ExpenseID)
	if err != nil {
		responsePayload.Setmsg("invalid expense_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	userId := c.Locals("user_id").(int64)

	data, err := controller.ExpenseService.Detail(expenseId, custId, userId)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}
		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: 1,
		PageCurrent: 1,
		PageLimit:   1,
		PageTotal:   1,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseController) ExpenseTypeLookup(c *fiber.Ctx) error {
	var dataFilter entity.ExpenseTypeQueryFilter
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Set default pagination
	dataFilter.Limit, dataFilter.Page = controller.setPaginationDefaults(dataFilter.Limit, dataFilter.Page, 100)

	data, total, lastPage, err := controller.ExpenseService.ExpenseTypeLookup(dataFilter)
	if err != nil {
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

func (controller *ExpenseController) OutletLookup(c *fiber.Ctx) error {
	var dataFilter entity.OutletLookupQueryFilter
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(dataFilter, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)

	// Try to get emp_id from JWT, fallback to user_id if emp_id is 0 or not available
	var salesmanId int64
	if empId, ok := c.Locals("emp_id").(int64); ok && empId > 0 {
		salesmanId = empId
	} else {
		salesmanId = c.Locals("user_id").(int64) // Fallback to user_id
	}

	// Set default pagination
	dataFilter.Limit, dataFilter.Page = controller.setPaginationDefaults(dataFilter.Limit, dataFilter.Page, 100)

	data, total, lastPage, err := controller.ExpenseService.OutletLookupByPJP(dataFilter, salesmanId, custId)
	if err != nil {
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

func (controller *ExpenseController) Create(c *fiber.Ctx) error {
	var request entity.CreateExpenseBody
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Parse form values
	request, err = controller.parseExpenseFormValues(c, form)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Get files
	files := form.File["files"]
	if len(files) == 0 {
		responsePayload.Setmsg("files is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custID := c.Locals("cust_id").(string)
	userID := c.Locals("user_id").(int64)
	empID := c.Locals("emp_id").(int64)

	request.CustID = custID
	request.EmpID = empID

	errs := controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ExpenseService.Create(c.UserContext(), request, files, userID)
	if err != nil {
		message := "bad request"
		if errors.Is(err, service.ErrExpenseTypeNotFound) ||
			errors.Is(err, service.ErrOutletNotFound) ||
			errors.Is(err, service.ErrFileLimitExceeded) {
			message = err.Error()
		}
		responsePayload.Setmsg(message)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(map[string]interface{}{
		"expense_id": data.ExpenseID,
	})
	responsePayload.Setmsg(constant.SUCCESSFULLY_ADDED)
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
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Parse expense_id from string to int64
	expenseId, err := controller.parseExpenseID(params.ExpenseID)
	if err != nil {
		responsePayload.Setmsg("invalid expense_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Parse form values
	parsedRequest, err := controller.parseExpenseFormValues(c, form)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	// Convert CreateExpenseBody to UpdateExpenseBody
	request.ExpenseTypeID = parsedRequest.ExpenseTypeID
	request.OutletID = parsedRequest.OutletID
	request.Amount = parsedRequest.Amount
	request.Note = parsedRequest.Note
	request.Folder = parsedRequest.Folder

	// Parse delete_file_ids (can be array or comma-separated)
	deleteFileIds, err := controller.parseDeleteFileIds(c, form)
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	request.DeleteFileIDs = deleteFileIds

	// Get files (optional for update)
	var (
		files  = form.File["files"]
		custId = c.Locals("cust_id").(string)
		userId = c.Locals("user_id").(int64)
		empID  = c.Locals("emp_id").(int64)
	)

	request.CustID = custId
	request.EmpID = empID

	errs = controller.validator.ValidateStruct(request, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	data, err := controller.ExpenseService.Update(expenseId, request, files, userId, custId)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		message := "update failed"
		if errors.Is(err, service.ErrExpenseNotFound) {
			statusCode = fiber.StatusNotFound
			message = err.Error()
		} else if errors.Is(err, service.ErrExpenseTypeNotFound) ||
			errors.Is(err, service.ErrOutletNotFound) ||
			errors.Is(err, service.ErrFileLimitExceeded) ||
			errors.Is(err, service.ErrDuplicateDeleteFileIDs) ||
			errors.Is(err, service.ErrFileNotFoundOrMismatch) {
			message = err.Error()
		}
		responsePayload.Setmsg(message)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(map[string]interface{}{
		"expense_id": data.ExpenseID,
	})
	responsePayload.Setmsg(constant.SUCCESSFULLY_UPDATED)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *ExpenseController) Delete(c *fiber.Ctx) error {
	var params entity.DeleteExpenseParams
	headerAcceptLang := controller.getAcceptLanguage(c)
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, constant.HEADER_ACCEPT_LANG)
	if errs != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)

	// Parse expense_id from string to int64
	expenseId, err := controller.parseExpenseID(params.ExpenseID)
	if err != nil {
		responsePayload.Setmsg("invalid expense_id")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err = controller.ExpenseService.Delete(expenseId, custId, userId)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		message := "delete failed"
		if errors.Is(err, service.ErrExpenseNotFound) {
			statusCode = fiber.StatusNotFound
			message = err.Error()
		}
		responsePayload.Setmsg(message)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(nil)
	responsePayload.Setmsg(constant.STATUS_OK)
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// Helper methods

// getAcceptLanguage extracts Accept-Language header from request
func (controller *ExpenseController) getAcceptLanguage(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}

// parseExpenseID parses expense_id from string to int64
func (controller *ExpenseController) parseExpenseID(expenseIDStr string) (int64, error) {
	return strconv.ParseInt(expenseIDStr, 10, 64)
}

// setPaginationDefaults sets default values for pagination
func (controller *ExpenseController) setPaginationDefaults(limit, page int, defaultLimit int) (int, int) {
	if limit <= 0 || limit > 9999 {
		limit = defaultLimit
	}
	if page <= 0 {
		page = 1
	}
	return limit, page
}

// parseExpenseFormValues parses multipart form values into CreateExpenseBody or UpdateExpenseBody
func (controller *ExpenseController) parseExpenseFormValues(c *fiber.Ctx, form *multipart.Form) (entity.CreateExpenseBody, error) {
	var request entity.CreateExpenseBody

	// Parse expense_type_id
	if expenseTypeIdStr := c.FormValue("expense_type_id"); expenseTypeIdStr != "" {
		if err := controller.validator.Validator.Var(expenseTypeIdStr, "numeric"); err != nil {
			return request, fmt.Errorf("invalid expense_type_id")
		}
		expenseTypeId, _ := strconv.Atoi(expenseTypeIdStr)
		request.ExpenseTypeID = expenseTypeId
	}

	// Parse outlet_id[] array
	if outletIds := form.Value["outlet_id[]"]; len(outletIds) > 0 {
		request.OutletID = make([]int, 0)
		for _, outletIdStr := range outletIds {
			if err := controller.validator.Validator.Var(outletIdStr, "numeric"); err != nil {
				return request, fmt.Errorf("invalid outlet_id")
			}
			outletId, _ := strconv.Atoi(outletIdStr)
			request.OutletID = append(request.OutletID, outletId)
		}
	}

	// Parse amount
	if amountStr := c.FormValue("amount"); amountStr != "" {
		// Use validator to validate numeric format (supports float)
		if err := controller.validator.Validator.Var(amountStr, "numeric"); err != nil {
			return request, fmt.Errorf("invalid amount")
		}
		amount, _ := strconv.ParseFloat(amountStr, 64)
		request.Amount = amount
	}

	request.Note = c.FormValue("note")
	request.Folder = c.FormValue("folder")

	return request, nil
}

// parseDeleteFileIds parses delete_file_ids from form (supports array or comma-separated)
func (controller *ExpenseController) parseDeleteFileIds(c *fiber.Ctx, form *multipart.Form) ([]int64, error) {
	var deleteFileIds []int64

	// Try to get as array first (delete_file_ids[]=1&delete_file_ids[]=2)
	if deleteFileIdStrs := form.Value["delete_file_ids[]"]; len(deleteFileIdStrs) > 0 {
		deleteFileIds = make([]int64, 0, len(deleteFileIdStrs))
		for _, idStr := range deleteFileIdStrs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid delete_file_ids: %s", idStr)
			}
			deleteFileIds = append(deleteFileIds, id)
		}
		return deleteFileIds, nil
	}

	// Try to get as single value (delete_file_ids=1,2,3)
	if deleteFileIdsStr := c.FormValue("delete_file_ids"); deleteFileIdsStr != "" {
		// Split by comma
		idStrs := strings.Split(deleteFileIdsStr, ",")
		deleteFileIds = make([]int64, 0, len(idStrs))
		for _, idStr := range idStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			// Use validator to validate number format
			if err := controller.validator.Validator.Var(idStr, "numeric"); err != nil {
				return nil, fmt.Errorf("invalid delete_file_ids: %s", idStr)
			}
			id, _ := strconv.ParseInt(idStr, 10, 64)
			deleteFileIds = append(deleteFileIds, id)
		}
		return deleteFileIds, nil
	}

	// No delete_file_ids provided, return empty slice
	return []int64{}, nil
}
