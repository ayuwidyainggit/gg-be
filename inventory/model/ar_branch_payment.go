package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArBranchPaymentList struct {
	GrBranchPaymentId  int        `gorm:"column:gr_branch_payment_id;primaryKey" json:"gr_branch_payment_id"`
	CustID             string     `gorm:"column:cust_id" json:"cust_id"`
	InvoiceNoBranch    string     `gorm:"column:invoice_no_branch" json:"invoice_no_branch"`
	PaymentOption      int        `gorm:"column:payment_option" json:"payment_option"`
	PaymentType        int        `gorm:"column:payment_type" json:"payment_type"`
	PaymentAmount      float64    `gorm:"column:payment_amount" json:"payment_amount"`
	PaymentBalance     float64    `gorm:"column:payment_balance" json:"payment_balance"`
	Discount           float64    `gorm:"column:discount" json:"discount"`
	TotalPayment       float64    `gorm:"column:total_payment" json:"total_payment"`
	DepositNo          string     `gorm:"column:deposit_no" json:"deposit_no"`
	DepositDate        *time.Time `gorm:"column:deposit_date" json:"deposit_date"`
	VerificationStatus int        `gorm:"column:verification_status" json:"verification_status"`
	VerifiedBy         *int64     `gorm:"column:verified_by" json:"verified_by"`
	VerifiedByName     *string    `gorm:"column:verified_by_name" json:"verified_by_name"`
	VerifiedAt         *time.Time `gorm:"column:verified_at" json:"verified_at"`
	Notes              *string    `gorm:"column:notes" json:"notes"`
}

func (ArBranchPaymentList) TableName() string {
	return "inv.gr_branch_payment"
}

type ArBranchPaymentCreate struct {
	GrBranchPaymentId  int64      `gorm:"column:gr_branch_payment_id;primaryKey" json:"gr_branch_payment_id"`
	CustID             string     `gorm:"column:cust_id" json:"cust_id"`
	InvoiceNoBranch    *string    `gorm:"invoice_no_branch" json:"invoice_no_branch"`
	PaymentOption      *int       `gorm:"column:payment_option" json:"payment_option"`
	PaymentType        *int       `gorm:"column:payment_type" json:"payment_type"`
	PaymentAmount      *float64   `gorm:"column:payment_amount" json:"payment_amount"`
	PaymentBalance     *float64   `gorm:"column:payment_balance" json:"payment_balance"`
	Discount           *float64   `gorm:"column:discount" json:"discount"`
	TotalPayment       *float64   `gorm:"column:total_payment" json:"total_payment"`
	VerificationStatus *int       `gorm:"column:verification_status" json:"verification_status"`
	DepositNo          *string    `gorm:"deposit_no" json:"deposit_no"`
	DepositDate        *time.Time `gorm:"deposit_date" json:"deposit_date"`
	Notes              *string    `gorm:"notes" json:"notes"`
}

func (ArBranchPaymentCreate) TableName() string {
	return "inv.gr_branch_payment"
}

func (m *ArBranchPaymentCreate) BeforeCreate(trx *gorm.DB) (err error) {
	var invoiceNoBranch InvoiceNoBranch
	trCode := "DPB"
	depositDateStr := m.DepositDate.Format("2006-01-02")
	depositDateSubtr := depositDateStr[2:4] + depositDateStr[5:7] + depositDateStr[8:10]
	// log.Println("grBranchDateStr:", grBranchDateStr)
	// log.Println("grBranchDateSubtr:", grBranchDateSubtr)

	queryStr := fmt.Sprintf(`SELECT
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(deposit_no,10,3),'999')),0)+1, '000')) AS get_no_fn
	FROM inv.gr_branch_payment
	WHERE substr(deposit_no,4,6) = '%v' AND cust_id = '%v'`, depositDateSubtr, strings.ToUpper(m.CustID))

	err = trx.Raw(queryStr).Scan(&invoiceNoBranch).Error
	if err != nil {
		return err
	}

	// log.Println("grBranchNo:", grBranchNo.GrBranchNo)
	invoiceNo := trCode + depositDateSubtr + invoiceNoBranch.InvoiceNoBranch
	m.DepositNo = &invoiceNo
	return nil
}

type InvoiceNoBranch struct {
	InvoiceNoBranch string `gorm:"column:get_no_fn"`
}
