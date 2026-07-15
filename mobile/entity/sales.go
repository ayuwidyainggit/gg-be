package entity

import "errors"

type SalesSummaryRequest struct {
}

type SalesSummaryResponse struct {
	CurrentSales       int64 `json:"current_sales"`
	DailyTarget        int64 `json:"daily_target"`
	MonthlySalesTarget int64 `json:"monthly_sales_target"`
}

var (
	ErrUserNotSalesman            = errors.New("user should be salesman canvas or taking order")
	ErrUserNotSalesmanCanvas      = errors.New("user should be salesman canvas")
	ErrUserNotSalesmanTakingOrder = errors.New("user should be salesman canvas taking order")
)
