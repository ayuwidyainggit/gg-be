package service

import (
	"context"
	"errors"
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/pkg/constant"
	"mobile/pkg/str"
	"mobile/repository"
	"time"
)

// PrintDailyReportService interface for print daily report service operations
type PrintDailyReportService interface {
	GetDailyReport(request entity.PrintDailyReportRequest, custID string, userID int64) (entity.PrintDailyReportResponse, error)
}

type printDailyReportServiceImpl struct {
	Config                     env.ConfigEnv
	Transaction                repository.Dbtransaction
	PrintDailyReportRepository repository.PrintDailyReportRepository
	UserRepository             repository.UserRepository
}

// NewPrintDailyReportService creates a new print daily report service
func NewPrintDailyReportService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
	printDailyReportRepo repository.PrintDailyReportRepository,
	userRepository repository.UserRepository,
) PrintDailyReportService {
	return &printDailyReportServiceImpl{
		Config:                     config,
		Transaction:                transaction,
		PrintDailyReportRepository: printDailyReportRepo,
		UserRepository:             userRepository,
	}
}

// GetDailyReport retrieves daily report data including payment and expense data
func (service *printDailyReportServiceImpl) GetDailyReport(request entity.PrintDailyReportRequest, custID string, userID int64) (entity.PrintDailyReportResponse, error) {
	var response entity.PrintDailyReportResponse
	ctx := context.Background()

	// Determine target date: use provided date or default to today
	targetDate := time.Now().UTC()
	if request.Date != nil {
		targetDate = str.UnixTimestampToUtcTime(*request.Date)
	}

	// Determine user ID: use provided user_id or default to logged-in user
	targetUserID := userID
	if request.UserID != nil {
		targetUserID = *request.UserID
	}

	// Get user information
	user, err := service.UserRepository.FindOneByUserID(targetUserID)
	if err != nil {
		if err.Error() == constant.STATUS_DB_NOT_FOUND || err.Error() == "record not found" {
			return response, errors.New("record not found")
		}
		return response, err
	}

	// Get sales data
	salesData, err := service.PrintDailyReportRepository.FindSalesDataByCustIdAndDate(ctx, request.EmpID, custID, targetDate)
	if err != nil {
		return response, err
	}

	// Get expense data
	expenseData, err := service.PrintDailyReportRepository.FindExpenseDataByCustIdAndDate(ctx, custID, targetDate, request.EmpID)
	if err != nil {
		return response, err
	}

	// Check attendance status for is_clock_out
	// Logic: 1 = clock in, 2 = clock out
	// Note: Field IsClockOut will be added to response entity in next task
	attendanceType, err := service.PrintDailyReportRepository.FindAttendanceByUserAndDate(ctx, targetUserID, custID, targetDate)
	var isClockOut int
	if err == nil {
		isClockOut = attendanceType
	}
	response.IsClockOut = isClockOut

	// Build response
	response.Date = targetDate.Format("02/01/2006")

	if user.UserId != nil {
		response.SalesmanID = *user.UserId
	}
	if user.Fullname != nil {
		response.SalesmanName = *user.Fullname
	}

	// Map sales data
	response.SalesData.Items = make([]entity.SalesDataDetail, 0)

	for _, data := range salesData {
		salesItem := entity.SalesDataDetail{
			SellingType: data.PaymentTypeCode,
			Amount:      data.PaymentAmount,
		}
		response.SalesData.Items = append(response.SalesData.Items, salesItem)
		response.SalesData.SalesTotal += data.PaymentAmount
	}

	response.PaymentData = []entity.PaymentData{}

	// Map expense data
	response.ExpenseData.Items = make([]entity.ExpenseDataDetail, 0)
	for _, expense := range expenseData {
		expenseItem := entity.ExpenseDataDetail{
			Amount: expense.Amount,
		}
		if expense.ExpenseTypeName != nil {
			expenseItem.ExpenseName = *expense.ExpenseTypeName
		}
		response.ExpenseData.Items = append(response.ExpenseData.Items, expenseItem)
		response.ExpenseData.ExpenseTotal += expense.Amount
	}

	return response, nil
}
