package model

import (
	"time"
)

type DepositReport struct {
	DepositPaymentID int        `gorm:"column:deposit_payment_id" json:"deposit_payment_id"`
	PayType          int        `gorm:"column:pay_type" json:"pay_type"`
	PayTypeName      string     `gorm:"column:pay_type_name" json:"pay_type_name"`
	DepositDate      time.Time  `gorm:"column:deposit_date" json:"deposit_date"`
	DepositNo        string     `gorm:"column:deposit_no" json:"deposit_no"`
	DocumentDate     string     `gorm:"column:document_date" json:"document_date"`
	DocumentNo       *string    `gorm:"column:document_no" json:"document_no"`
	DueDate          time.Time  `gorm:"column:due_date" json:"due_date"`
	Owner            string     `gorm:"column:owner_name" json:"owner_name"`
	AccountNo        *string    `gorm:"column:account_no" json:"account_no"`
	BankID           *int       `gorm:"column:bank_id" json:"bank_id"`
	BankName         *string    `gorm:"column:bank_name" json:"bank_name"`
	PaymentAmount    float64    `gorm:"column:payment_amount" json:"payment_amount"`
	InvoiceDate      time.Time  `gorm:"column:invoice_date" json:"invoice_date"`
	InvoiceNo        string     `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceAmount    float64    `gorm:"column:invoice_amount" json:"invoice_amount"`
	OutletID         int64      `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode       string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       string     `gorm:"column:outlet_name" json:"outlet_name"`
	EmpID            int        `gorm:"column:emp_id" json:"emp_id"`
	EmpName          string     `gorm:"column:emp_name" json:"emp_name"`
	EmpCode          string     `gorm:"column:emp_code" json:"emp_code"`
	EmpGrpID         int        `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpName       string     `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	ClearingDate     *time.Time `gorm:"column:clearing_date" json:"clearing_date"`
	StatusCheque     *int       `gorm:"column:clearing_status" json:"clearing_status"`
	Notes            *string    `gorm:"column:notes" json:"notes"`
}

// TableName sets the insert table name for this struct type
func (DepositReport) TableName() string {
	return "acf.deposit_payment"
}

type DepositPayTypeLookup struct {
	PayType     int    `db:"pay_type" json:"pay_type"`
	PayTypeName string `db:"pay_type_name" json:"pay_type_name"`
}

func (DepositPayTypeLookup) TableName() string {
	return "acf.deposit_payment"
}

type DepositPaymentLookup struct {
	DocNo   string  `db:"doc_no" json:"doc_no"`
	Amount  float64 `db:"amount" json:"amount"`
	Balance float64 `db:"balance" json:"balance"`
}

type DepositNoReportLookup struct {
	DepositNo string `db:"deposit_no" json:"deposit_no"`
}

func (DepositNoReportLookup) TableName() string {
	return "acf.deposit_payment"
}
