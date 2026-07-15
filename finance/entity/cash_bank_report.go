package entity

import "time"

type DepositPayTypeLookup struct {
	PayType     int    `json:"pay_type"`
	PayTypeName string `json:"pay_type_name"`
}

type DepositReportResponse struct {
	DepositPaymentID int        `json:"deposit_payment_id"`
	PayType          int        `json:"pay_type"`
	PayTypeName      string     `json:"pay_type_name"`
	DepositDate      time.Time  `json:"deposit_date"`
	DepositNo        string     `json:"deposit_no"`
	DocumentDate     string     `json:"document_date"`
	DocumentNo       *string    `json:"document_no"`
	DueDate          time.Time  `json:"due_date"`
	Owner            string     `json:"owner_name"`
	AccountNo        *string    `json:"account_no"`
	BankName         *string    `json:"bank_name"`
	PaymentAmount    float64    `json:"payment_amount"`
	InvoiceDate      time.Time  `json:"invoice_date"`
	InvoiceNo        string     `json:"invoice_no"`
	InvoiceAmount    float64    `json:"invoice_amount"`
	OutletID         int64      `json:"outlet_id"`
	OutletCode       string     `json:"outlet_code"`
	OutletName       string     `json:"outlet_name"`
	EmpID            int        `json:"emp_id"`
	EmpName          string     `json:"emp_name"`
	EmpCode          string     `json:"emp_code"`
	EmpGrpID         int        `json:"emp_grp_id"`
	EmpGrpName       string     `json:"emp_grp_name"`
	ClearingDate     *time.Time `json:"clearing_date"`
	StatusClearing   *string    `json:"status_clearing"`
	Notes            *string    `json:"notes"`
}

var PayType = map[int]string{
	1: "Cash",
	2: "Cheque",
	3: "Transfer",
	4: "Return",
	5: "Credit",
}
