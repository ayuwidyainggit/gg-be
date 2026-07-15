package entity

type TotalSummaryDeposit struct {
	TotalPayment        float64 `json:"total_payment"`
	TotalExpense        float64 `json:"total_expense"`
	TotalCollection     float64 `json:"total_collection"`
	TotalReceived       float64 `json:"total_received"`
	TotalMaterai        float64 `json:"total_materai"`
	TotalDiscount       float64 `json:"total_discount"`
	TotalPaymentBalance float64 `json:"total_payment_balance"`
}
