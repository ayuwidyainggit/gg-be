package model

import "time"

// PaymentDepositReportRow - GORM model untuk result query GROUP BY
type PaymentDepositReportRow struct {
	DepositDate       time.Time `gorm:"column:deposit_date"`
	DepositType       string    `gorm:"column:deposit_type"`
	DepositNo         string    `gorm:"column:deposit_no"`
	CollectorID       *int      `gorm:"column:collector_id"`
	CollectorCode     *string   `gorm:"column:collector_code"`
	CollectorName     *string   `gorm:"column:collector_name"`
	TotalPayment      float64   `gorm:"column:total_payment"`
	CashAmount        float64   `gorm:"column:cash_amount"`
	ChequeAmount      float64   `gorm:"column:cheque_amount"`
	TransferAmount    float64   `gorm:"column:transfer_amount"`
	ReturnAmount      float64   `gorm:"column:return_amount"`
	CreditDebitAmount float64   `gorm:"column:credit_debit_amount"`
	ExpenseAmount     float64   `gorm:"column:expense_amount"`
}

// PaymentDepositReportDownloadRow represents detail rows used for XLSX export.
type PaymentDepositReportDownloadRow struct {
	DepositDate    time.Time  `gorm:"column:deposit_date"`
	DepositType    string     `gorm:"column:deposit_type"`
	DepositNo      string     `gorm:"column:deposit_no"`
	Collector      *string    `gorm:"column:collector"`
	DocumentDate   *time.Time `gorm:"column:document_date"`
	Code           *string    `gorm:"column:code"`
	BusinessName   *string    `gorm:"column:business_name"`
	DocumentNo     *string    `gorm:"column:document_no"`
	Cash           float64    `gorm:"column:cash"`
	ChequeGiro     float64    `gorm:"column:cheque_giro"`
	Transfer       float64    `gorm:"column:transfer"`
	ReturnAmount   float64    `gorm:"column:return_amount"`
	CreditDebit    float64    `gorm:"column:credit_debit"`
	Discount       float64    `gorm:"column:discount"`
	PaymentBalance float64    `gorm:"column:payment_balance"`
	Expense        float64    `gorm:"column:expense"`
	ExpenseName    *string    `gorm:"column:expense_name"`
}

// PaymentDepositReportSummaryRow - GORM model untuk summary aggregate
type PaymentDepositReportSummaryRow struct {
	TotalCash        float64 `gorm:"column:total_cash"`
	TotalCheque      float64 `gorm:"column:total_cheque"`
	TotalTransfer    float64 `gorm:"column:total_transfer"`
	TotalReturn      float64 `gorm:"column:total_return"`
	TotalCreditDebit float64 `gorm:"column:total_credit_debit"`
	TotalExpense     float64 `gorm:"column:total_expense"`
}

type PaymentDepositReportRecapRow struct {
	DepositType    string  `gorm:"column:deposit_type"`
	Cash           float64 `gorm:"column:cash"`
	ChequeGiro     float64 `gorm:"column:cheque_giro"`
	Transfer       float64 `gorm:"column:transfer"`
	ReturnAmount   float64 `gorm:"column:return_amount"`
	CreditDebit    float64 `gorm:"column:credit_debit"`
	Discount       float64 `gorm:"column:discount"`
	PaymentBalance float64 `gorm:"column:payment_balance"`
	Expense        float64 `gorm:"column:expense"`
}
