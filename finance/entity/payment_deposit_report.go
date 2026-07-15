package entity

import "time"

const (
	PaymentDepositReportDownloadPrefix       = "DownloadDepositPayment"
	PaymentDepositReportProcessingMessage    = "Processing time may vary by file size. Please check Download History to access the file"
	PaymentDepositReportStatusProcessing     = 0
	PaymentDepositReportStatusReady          = 1
	PaymentDepositReportStatusNameProcessing = "Processing"
	PaymentDepositReportStatusNameReady      = "Ready"
)

// PaymentDepositReportQueryFilter - query params untuk endpoint
type PaymentDepositReportQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int      `query:"page"`
	Limit        int      `query:"limit"`
	Sort         string   `query:"sort"`
	StartDate    string   `query:"start_date"`
	EndDate      string   `query:"end_date"`
	Q            string   `query:"q"`
	DepositType  []string `query:"deposit_type"`
	EmpID        []string `query:"emp_id"`
	SalesmanID   []string `query:"salesman_id"`
	DepositNo    []string `query:"deposit_no"`
}

// PaymentDepositReportItem - single row response
type PaymentDepositReportItem struct {
	DepositDate       string  `json:"deposit_date"`
	DepositType       string  `json:"deposit_type"`
	DepositNo         string  `json:"deposit_no"`
	CollectorID       *int    `json:"collector_id"`
	CollectorCode     *string `json:"collector_code"`
	CollectorName     *string `json:"collector_name"`
	CashAmount        float64 `json:"cash_amount"`
	ChequeAmount      float64 `json:"cheque_amount"`
	TransferAmount    float64 `json:"transfer_amount"`
	ReturnAmount      float64 `json:"return_amount"`
	CreditDebitAmount float64 `json:"credit_debit_amount"`
	ExpenseAmount     float64 `json:"expense_amount"`
	TotalPayment      float64 `json:"total_payment"`
}

// PaymentDepositReportDownloadRow - export row for XLSX download.
type PaymentDepositReportDownloadRow struct {
	DepositDate    string  `json:"deposit_date"`
	DepositType    string  `json:"deposit_type"`
	DepositNo      string  `json:"deposit_no"`
	Collector      string  `json:"collector"`
	DocumentDate   string  `json:"document_date"`
	Code           string  `json:"code"`
	BusinessName   string  `json:"business_name"`
	DocumentNo     string  `json:"document_no"`
	Cash           float64 `json:"cash"`
	ChequeGiro     float64 `json:"cheque_giro"`
	Transfer       float64 `json:"transfer"`
	ReturnAmount   float64 `json:"return_amount"`
	CreditDebit    float64 `json:"credit_debit"`
	Discount       float64 `json:"discount"`
	PaymentBalance float64 `json:"payment_balance"`
	Expense        float64 `json:"expense"`
	ExpenseName    string  `json:"expense_name"`
}

// PaymentDepositReportSummary - summary totals
type PaymentDepositReportSummary struct {
	TotalAmount      float64 `json:"total_amount"`
	TotalCash        float64 `json:"total_cash"`
	TotalTransfer    float64 `json:"total_transfer"`
	TotalCheque      float64 `json:"total_cheque"`
	TotalReturn      float64 `json:"total_return"`
	TotalCreditDebit float64 `json:"total_credit_debit"`
	TotalExpense     float64 `json:"total_expense"`
	GrandTotal       float64 `json:"grand_total"`
}

type PaymentDepositReportSummaryByDepositTypeItem struct {
	DepositTypeLabel    string  `json:"deposit_type_label"`
	SummaryCash         float64 `json:"summary_cash"`
	ChequeGiro          float64 `json:"cheque_giro"`
	Transfer            float64 `json:"transfer"`
	ReturnAmount        float64 `json:"return_amount"`
	CreditDebit         float64 `json:"credit_debit"`
	Discount            float64 `json:"discount"`
	TotalPaymentBalance float64 `json:"total_payment_balance"`
	TotalExpense        float64 `json:"total_expense"`
}

// PaymentDepositReportPagination
type PaymentDepositReportPagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

// PaymentDepositReportResponse - full response wrapper
type PaymentDepositReportResponse struct {
	Items                 []PaymentDepositReportItem                      `json:"items"`
	Summary               PaymentDepositReportSummary                     `json:"summary"`
	SummaryByDepositType  []PaymentDepositReportSummaryByDepositTypeItem `json:"summary_by_deposit_type"`
	Pagination            PaymentDepositReportPagination                  `json:"pagination"`
}

// ReportListCreate - untuk insert ke report.list
type ReportListCreate struct {
	CustID     string
	ReportID   string
	ReportName string
	StartDate  time.Time
	EndDate    time.Time
	FileStatus int
	FileURL    string
	FileBase64 string
	CreatedBy  string
	CreatedAt  time.Time
}

// ReportListResponse - response ketika download masih processing atau sukses
type ReportListResponse struct {
	ReportID       string `json:"report_id"`
	ReportName     string `json:"report_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	FileStatus     int    `json:"file_status"`
	FileStatusName string `json:"file_status_name"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      string `json:"created_at"`
}
