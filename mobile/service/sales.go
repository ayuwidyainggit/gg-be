package service

import (
	"context"
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/repository"
	"time"
)

type SalesService interface {
	SalesSummary(request entity.SalesSummaryRequest, custID string, empID int64) (entity.SalesSummaryResponse, error)
}

type salesServiceImpl struct {
	Config          env.ConfigEnv
	Transaction     repository.Dbtransaction
	SalesRepository repository.SalesRepository
}

// NewSalesService creates a new instance of SalesService with required dependencies
func NewSalesService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
	salesRepository repository.SalesRepository,
) *salesServiceImpl {
	return &salesServiceImpl{
		Config:          config,
		Transaction:     transaction,
		SalesRepository: salesRepository,
	}
}

// SalesSummary calculates the current sales summary for a salesman.
// It calculates current_sales as total order minus total return within a dynamic date range:
// - start_date: first day of the current month
// - end_date: latest clock in date (type = 1) from mobile.attendances, or current date if no clock in found
// Returns SalesSummaryResponse with current_sales, daily_target, and monthly_sales_target
func (service *salesServiceImpl) SalesSummary(request entity.SalesSummaryRequest, custID string, empID int64) (response entity.SalesSummaryResponse, err error) {
	// Calculate start_date: first day of current month at 00:00:00
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get end_date: latest clock in date from attendance records
	ctx := context.Background()
	endDate, err := service.SalesRepository.GetLatestClockInDate(ctx, custID, empID)
	if err != nil {
		return response, err
	}

	// If no clock in found, use current date as end_date
	if endDate == nil {
		today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		endDate = &today
	} else {
		// Set time to end of day (23:59:59) to include all records for that day
		*endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
	}

	// Get total order amount within date range
	totalOrder, err := service.SalesRepository.GetTotalOrder(ctx, custID, empID, startDate, *endDate)
	if err != nil {
		return response, err
	}

	// Get total return amount within date range
	totalReturn, err := service.SalesRepository.GetTotalReturn(ctx, custID, empID, startDate, *endDate)
	if err != nil {
		return response, err
	}

	// Get monthly sales target from m_sales_target and m_sales_allocated
	monthlySalesTarget, err := service.SalesRepository.GetMonthlySalesTarget(ctx, custID, empID, int(now.Month()), now.Year())
	if err != nil {
		return response, err
	}

	// Calculate current sales = total order - total return
	response.CurrentSales = totalOrder - totalReturn
	// Default value, not yet implemented per enhancement requirements
	response.DailyTarget = 0
	// Set monthly sales target from database
	response.MonthlySalesTarget = monthlySalesTarget

	return response, err
}
