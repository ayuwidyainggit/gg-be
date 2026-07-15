package entity

// PrintDailyReportRequest represents query parameters for print daily report
type PrintDailyReportRequest struct {
	EmpID  int64  `query:"-"`
	Date   *int64 `query:"date"`
	UserID *int64 `query:"user_id"`
}

// PrintDailyReportResponse represents the response for print daily report
type PrintDailyReportResponse struct {
	Date         string        `json:"date"`
	SalesmanID   int64         `json:"salesman_id"`
	SalesmanName string        `json:"salesman_name"`
	IsClockOut   int           `json:"is_clock_out"`
	SalesData    SalesData     `json:"sales_data"`
	PaymentData  []PaymentData `json:"payment_data"`
	ExpenseData  ExpenseData   `json:"expense_data"`
}

// SalesData represents sales data in daily report
type SalesData struct {
	SalesTotal float64           `json:"sales_total"`
	Items      []SalesDataDetail `json:"items"`
}

type SalesDataDetail struct {
	SellingType string  `json:"selling_type"`
	Amount      float64 `json:"amount"`
}

// PaymentData represents payment data in daily report
type PaymentData struct {
	PaymentType string  `json:"payment_type"`
	Amount      float64 `json:"amount"`
}

// ExpenseData represents expense data in daily report
type ExpenseData struct {
	ExpenseTotal float64             `json:"expense_total"`
	Items        []ExpenseDataDetail `json:"items"`
}

type ExpenseDataDetail struct {
	ExpenseName string  `json:"expense_name"`
	Amount      float64 `json:"amount"`
}
