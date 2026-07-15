package controller

import (
	"finance/entity"
	"finance/pkg/constant"
	"finance/pkg/middleware"
	"finance/pkg/responsebuild"
	"finance/pkg/str"
	"finance/service"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type PaymentDepositReportController struct {
	Service service.PaymentDepositReportService
}

func NewPaymentDepositReportController(service service.PaymentDepositReportService, _ interface{}) *PaymentDepositReportController {
	return &PaymentDepositReportController{
		Service: service,
	}
}

func (controller *PaymentDepositReportController) Route(app *fiber.App) {
	grRouteV1 := app.Group("/v1/reports/payment-deposit", middleware.JWTProtected())
	grRouteV1.Get("", controller.List)
	grRouteV1.Get("/download", controller.Download)
}

func (controller *PaymentDepositReportController) List(c *fiber.Ctx) error {
	var dataFilter entity.PaymentDepositReportQueryFilter
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PaymentDepositReportController, List, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	normalizePaymentDepositFilter(&dataFilter)

	if err := normalizeAndValidatePaymentDepositListFilter(&dataFilter); err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Locals
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	// Service Call
	data, err := controller.Service.ListReport(dataFilter)
	if err != nil {
		log.Error("PaymentDepositReportController, List, Service.ListReport, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Response
	responsePayload.Setmsg("Data berhasil ditampilkan")
	responsePayload.Setdata(data)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *PaymentDepositReportController) Download(c *fiber.Ctx) error {
	var dataFilter entity.PaymentDepositReportQueryFilter
	var headerAcceptLang string

	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}

	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("PaymentDepositReportController, Download, QueryParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	normalizePaymentDepositFilter(&dataFilter)

	if err := normalizeAndValidatePaymentDepositDownloadFilter(&dataFilter); err != nil {
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Locals
	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

	// Identify User (CreatedBy)
	createdBy := "System"
	if userName, ok := c.Locals("user_name").(string); ok && userName != "" {
		createdBy = userName
	} else if email, ok := c.Locals("email").(string); ok && email != "" {
		createdBy = email
	}

	// Service Call
	_, err := controller.Service.DownloadReport(dataFilter, createdBy)
	if err != nil {
		log.Error("PaymentDepositReportController, Download, Service.DownloadReport, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Response
	responsePayload.Setmsg(entity.PaymentDepositReportProcessingMessage)
	responsePayload.Setdata(nil)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func normalizeAndValidatePaymentDepositFilter(filter *entity.PaymentDepositReportQueryFilter) error {
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Page < 1 {
		return fmt.Errorf("invalid page: must be greater than 0")
	}

	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit < 1 {
		return fmt.Errorf("invalid limit: must be greater than 0")
	}
	if filter.Limit > 9999 {
		filter.Limit = 9999
	}

	filter.DepositType = normalizeCSVValues(filter.DepositType)
	filter.EmpID = normalizeCSVValues(filter.EmpID)
	filter.SalesmanID = normalizeCSVValues(filter.SalesmanID)
	filter.DepositNo = normalizeCSVValues(filter.DepositNo)

	if len(filter.EmpID) == 0 && len(filter.SalesmanID) > 0 {
		filter.EmpID = append(filter.EmpID, filter.SalesmanID...)
	}

	depositTypes, err := normalizeAndValidateDepositTypes(filter.DepositType, true)
	if err != nil {
		return err
	}
	filter.DepositType = depositTypes

	if strings.TrimSpace(filter.StartDate) == "" {
		return fmt.Errorf("start_date is required")
	}
	if strings.TrimSpace(filter.EndDate) == "" {
		return fmt.Errorf("end_date is required")
	}
	if len(filter.DepositType) == 0 {
		return fmt.Errorf("deposit_type is required")
	}

	if normalizedSort, err := validateAndNormalizeSort(filter.Sort); err != nil {
		return err
	} else {
		filter.Sort = normalizedSort
	}

	startDate, err := normalizeDateInput(filter.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := normalizeDateInput(filter.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date: %w", err)
	}
	filter.StartDate = startDate
	filter.EndDate = endDate

	if filter.StartDate != "" && filter.EndDate != "" {
		start, _ := time.Parse("2006-01-02", filter.StartDate)
		end, _ := time.Parse("2006-01-02", filter.EndDate)
		if end.Before(start) {
			return fmt.Errorf("invalid date range: end_date must be greater than or equal to start_date")
		}
	}

	if err := validateIntList(filter.EmpID, "emp_id"); err != nil {
		return err
	}
	if err := validateIntList(filter.SalesmanID, "salesman_id"); err != nil {
		return err
	}

	return nil
}

func normalizeAndValidatePaymentDepositListFilter(filter *entity.PaymentDepositReportQueryFilter) error {
	return normalizeAndValidatePaymentDepositFilter(filter)
}

func normalizeAndValidatePaymentDepositDownloadFilter(filter *entity.PaymentDepositReportQueryFilter) error {
	if len(filter.DepositType) == 0 {
		filter.DepositType = []string{"AR", "AP"}
	}
	if err := normalizeAndValidatePaymentDepositFilter(filter); err != nil {
		return err
	}
	return nil
}

func normalizePaymentDepositFilter(filter *entity.PaymentDepositReportQueryFilter) {
	filter.DepositType = normalizeCSVValues(filter.DepositType)
	filter.EmpID = normalizeCSVValues(filter.EmpID)
	filter.SalesmanID = normalizeCSVValues(filter.SalesmanID)
	filter.DepositNo = normalizeCSVValues(filter.DepositNo)
}

func validateAndNormalizeSort(sortValue string) (string, error) {
	defaultSort := "deposit_date:desc"
	if strings.TrimSpace(sortValue) == "" {
		return defaultSort, nil
	}

	allowedFields := map[string]bool{
		"created_date":   true,
		"deposit_date":   true,
		"deposit_no":     true,
		"deposit_type":   true,
		"collector_name": true,
		"total_payment":  true,
	}

	normalizedSort := make([]string, 0)
	for _, raw := range strings.Split(sortValue, ",") {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}

		sortPart := strings.SplitN(part, ":", 2)
		if len(sortPart) != 2 {
			return "", fmt.Errorf("invalid sort format: use field:direction")
		}

		field := strings.TrimSpace(strings.ToLower(sortPart[0]))
		direction := strings.TrimSpace(strings.ToLower(sortPart[1]))

		if !allowedFields[field] {
			return "", fmt.Errorf("invalid sort field: %s", field)
		}
		if direction != "asc" && direction != "desc" {
			return "", fmt.Errorf("invalid sort direction: %s", direction)
		}

		if field == "created_date" {
			field = "deposit_date"
		}

		normalizedSort = append(normalizedSort, fmt.Sprintf("%s:%s", field, direction))
	}

	if len(normalizedSort) == 0 {
		return defaultSort, nil
	}

	return strings.Join(normalizedSort, ","), nil
}

func normalizeAndValidateDepositTypes(values []string, requireValue bool) ([]string, error) {
	set := map[string]struct{}{}
	for _, item := range values {
		for _, token := range strings.Split(item, ",") {
			value := strings.ToUpper(strings.TrimSpace(token))
			if value == "" {
				continue
			}
			if value == "ALL" {
				set["AR"] = struct{}{}
				set["AP"] = struct{}{}
				continue
			}
			if value != "AR" && value != "AP" {
				return nil, fmt.Errorf("invalid deposit_type: %s", value)
			}
			set[value] = struct{}{}
		}
	}
	if len(set) == 0 {
		if !requireValue {
			return []string{}, nil
		}
		return nil, fmt.Errorf("deposit_type is required")
	}
	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Strings(result)
	return result, nil
}

func normalizeCSVValues(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, item := range values {
		for _, token := range strings.Split(item, ",") {
			value := strings.TrimSpace(token)
			if value == "" {
				continue
			}
			normalized = append(normalized, value)
		}
	}
	return normalized
}

func validateIntList(values []string, field string) error {
	for _, item := range values {
		if _, err := strconv.Atoi(item); err != nil {
			return fmt.Errorf("invalid %s: must be integer", field)
		}
	}
	return nil
}

func normalizeDateInput(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}

	epochVal, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		t := str.UnixTimestampToUtcTime(epochVal)
		return t.Format("2006-01-02"), nil
	}

	if _, errDate := time.Parse("2006-01-02", value); errDate == nil {
		return value, nil
	}

	return "", fmt.Errorf("must be epoch or YYYY-MM-DD")
}
