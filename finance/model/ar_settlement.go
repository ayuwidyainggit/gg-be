package model

import (
	"time"

	"gorm.io/gorm"
)

type ArSettlement struct {
	CustID              string         `gorm:"column:cust_id" json:"cust_id"`
	DepositNo           string         `gorm:"column:deposit_no" json:"deposit_no"`
	DepositDate         *time.Time     `gorm:"deposit_date" json:"deposit_date"`
	EmpId               *int64         `gorm:"emp_id" json:"emp_id"`
	EmpCode             *string        `gorm:"emp_code" json:"emp_code"`
	EmpName             *string        `gorm:"emp_name" json:"emp_name"`
	TotalDiscount       *float64       `gorm:"total_discount" json:"total_discount"`
	TotalMaterai        *float64       `gorm:"total_materai" json:"total_materai"`
	TotalPaymentBalance *float64       `gorm:"total_payment_balance" json:"total_payment_balance"`
	TotalPayment        *float64       `gorm:"total_payment" json:"total_payment"`
	RemainingAmount     *float64       `gorm:"remaining_amount" json:"remaining_amount"`
	DepositStatus       *int64         `gorm:"deposit_status" json:"deposit_status"`
	CreatedBy           *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy           *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt           *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel               bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsApproved          bool           `gorm:"column:is_approved" json:"is_approved"`
	ApprovedBy          *string        `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt          *time.Time     `gorm:"column:approved_at" json:"approved_at"`
}

func (ArSettlement) TableName() string {
	return "acf.deposit"
}

type ArSettlementPayment struct {
	CustId           *string    `gorm:"column:cust_id" json:"cust_id"`
	CustName         *string    `gorm:"column:cust_name" json:"cust_name"`
	DepositPaymentID int64      `gorm:"column:deposit_payment_id;primaryKey" json:"deposit_payment_id"`
	InvoiceNo        *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate      *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	PayType          *int64     `gorm:"column:pay_type" json:"pay_type"`
	DocumentNo       *string    `gorm:"column:document_no" json:"document_no"`
	Balance          *float64   `gorm:"column:balance" json:"balance"`
	PaymentAmount    *float64   `gorm:"column:payment_amount" json:"payment_amount"`
	SalesmanId       *int64     `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName     *string    `gorm:"column:salesman_name" json:"salesman_name"`
	OutletId         *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	Discount         *float64   `gorm:"column:discount" json:"discount"`
	Materai          *float64   `gorm:"column:materai" json:"materai"`
	PaymentBalance   *float64   `gorm:"column:payment_balance" json:"payment_balance"`
	TotalPayment     *float64   `gorm:"column:total_payment" json:"total_payment"`
	RemainingPayment *float64   `gorm:"column:remaining_payment" json:"remaining_payment"`
}

func (ArSettlementPayment) TableName() string {
	return "acf.deposit_payment"
}

type ArBranchSettlementPayment struct {
	CustId           *string    `gorm:"column:cust_id" json:"cust_id"`
	DepositPaymentID int64      `gorm:"column:deposit_payment_id;primaryKey" json:"deposit_payment_id"`
	InvoiceNo        *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDate      *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	PayType          *int64     `gorm:"column:pay_type" json:"pay_type"`
	// DocumentNo       *string    `gorm:"column:document_no" json:"document_no"`
	// Balance          *float64   `gorm:"column:balance" json:"balance"`
	PaymentAmount *float64 `gorm:"column:payment_amount" json:"payment_amount"`
	// SalesmanId       *int64     `gorm:"column:salesman_id" json:"salesman_id"`
	// SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	// SalesmanName     *string    `gorm:"column:salesman_name" json:"salesman_name"`
	// OutletId         *int64     `gorm:"column:outlet_id" json:"outlet_id"`
	// OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	// OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	Discount *float64 `gorm:"column:discount" json:"discount"`
	// Materai          *float64   `gorm:"column:materai" json:"materai"`
	PaymentBalance   *float64 `gorm:"column:payment_balance" json:"payment_balance"`
	TotalPayment     *float64 `gorm:"column:total_payment" json:"total_payment"`
	RemainingPayment *float64 `gorm:"column:remaining_payment" json:"remaining_payment"`
}

func (ArBranchSettlementPayment) TableName() string {
	return "inv.gr_branch_payment"
}

type ArSettlementList struct {
	CustID              *string    `gorm:"column:cust_id" json:"cust_id"`
	CustName            *string    `gorm:"column:cust_name" json:"cust_name"`
	DepositNo           *string    `gorm:"deposit_no" json:"deposit_no"`
	DepositDate         *time.Time `gorm:"deposit_date" json:"deposit_date"`
	CollectionNo        *string    `gorm:"collection_no" json:"collection_no"`
	CollectionDate      *time.Time `gorm:"collection_date" json:"collection_date"`
	EmpId               *int64     `gorm:"emp_id" json:"emp_id"`
	EmpCode             *string    `gorm:"emp_code" json:"emp_code"`
	EmpName             *string    `gorm:"emp_name" json:"emp_name"`
	OtGrpId             *int64     `gorm:"ot_grp_id" json:"ot_grp_id"`
	OtGrpCode           *string    `gorm:"ot_grp_code" json:"ot_grp_code"`
	OtGrpName           *string    `gorm:"ot_grp_name" json:"ot_grp_name"`
	TotalDiscount       *float64   `gorm:"total_discount" json:"total_discount"`
	TotalMaterai        *float64   `gorm:"total_materai" json:"total_materai"`
	TotalPaymentBalance *float64   `gorm:"total_payment_balance" json:"total_payment_balance"`
	TotalPayment        *float64   `gorm:"total_payment" json:"total_payment"`
	RemainingAmount     *float64   `gorm:"remaining_amount" json:"remaining_amount"`
	DepositStatus       *int64     `gorm:"deposit_status" json:"deposit_status"`
	ApprovedBy          *int64     `gorm:"column:approved_by" json:"approved_by"`
	ApprovedByName      *string    `gorm:"column:approved_by_name" json:"approved_by_name"`
}

func (ArSettlementList) TableName() string {
	return "acf.deposit"
}

type ArBranchSettlementList struct {
	CustID              *string    `gorm:"column:cust_id" json:"cust_id"`
	CustName            *string    `gorm:"column:cust_name" json:"cust_name"`
	DepositPaymentId    int64      `gorm:"deposit_payment_id" json:"deposit_payment_id"`
	DepositNo           *string    `gorm:"deposit_no" json:"deposit_no"`
	DepositDate         *time.Time `gorm:"deposit_date" json:"deposit_date"`
	CollectionNo        *string    `gorm:"collection_no" json:"collection_no"`
	CollectionDate      *time.Time `gorm:"collection_date" json:"collection_date"`
	InvoiceNoBranch     *string    `gorm:"invoice_no_branch" json:"invoice_no_branch"`
	InvoiceDateBranch   *time.Time `gorm:"invoice_date_branch" json:"invoice_date_branch"`
	PayType             *int64     `gorm:"pay_type" json:"pay_type"`
	EmpId               *int64     `gorm:"emp_id" json:"emp_id"`
	EmpCode             *string    `gorm:"emp_code" json:"emp_code"`
	EmpName             *string    `gorm:"emp_name" json:"emp_name"`
	OtGrpId             *int64     `gorm:"ot_grp_id" json:"ot_grp_id"`
	OtGrpCode           *string    `gorm:"ot_grp_code" json:"ot_grp_code"`
	OtGrpName           *string    `gorm:"ot_grp_name" json:"ot_grp_name"`
	TotalDiscount       *float64   `gorm:"total_discount" json:"total_discount"`
	TotalMaterai        *float64   `gorm:"total_materai" json:"total_materai"`
	TotalPaymentBalance *float64   `gorm:"total_payment_balance" json:"total_payment_balance"`
	TotalPayment        *float64   `gorm:"total_payment" json:"total_payment"`
	RemainingAmount     *float64   `gorm:"remaining_amount" json:"remaining_amount"`
	DepositStatus       *int64     `gorm:"deposit_status" json:"deposit_status"`
	ApprovedBy          *int64     `gorm:"column:approved_by" json:"approved_by"`
	ApprovedByName      *string    `gorm:"column:approved_by_name" json:"approved_by_name"`
}

func (ArBranchSettlementList) TableName() string {
	return "inv.gr_branch_payment"
}

type SettlementCollectorFilter struct {
	EmpId   int    `gorm:"column:emp_id" json:"emp_id"`
	EmpCode string `gorm:"column:emp_code" json:"emp_code"`
	EmpName string `gorm:"column:emp_name" json:"emp_name"`
}

func (SettlementCollectorFilter) TableName() string {
	return "acf.deposit"
}

type SettlementDepositStatusFilter struct {
	DepositStatusId int `gorm:"column:deposit_status" json:"deposit_status"`
}

func (SettlementDepositStatusFilter) TableName() string {
	return "acf.deposit"
}

type DepositDetailByInvoice struct {
	InvoiceAmount    float64 `gorm:"column:invoice_amount" json:"invoice_amount"`
	TotalPayment     float64 `gorm:"column:total_payment" json:"total_payment"`
	RemainingPayment float64 `gorm:"column:remaining_payment" json:"remaining_payment"`
}

func (DepositDetailByInvoice) TableName() string {
	return "acf.deposit_detail"
}

type DepositBranchDetailByInvoice struct {
	InvoiceAmount    float64 `gorm:"column:invoice_amount" json:"invoice_amount"`
	TotalPayment     float64 `gorm:"column:total_payment" json:"total_payment"`
	RemainingPayment float64 `gorm:"column:remaining_payment" json:"remaining_payment"`
}

func (DepositBranchDetailByInvoice) TableName() string {
	return "inv.gr_branch_payment"
}

type ArBranchSettlement struct {
	DepositNo   string     `gorm:"column:deposit_no" json:"deposit_no"`
	DepositDate *time.Time `gorm:"deposit_date" json:"deposit_date"`
	// EmpId               *int64         `gorm:"emp_id" json:"emp_id"`
	// EmpCode             *string        `gorm:"emp_code" json:"emp_code"`
	// EmpName             *string        `gorm:"emp_name" json:"emp_name"`
	// TotalDiscount       *float64       `gorm:"total_discount" json:"total_discount"`
	// TotalMaterai        *float64       `gorm:"total_materai" json:"total_materai"`
	// TotalPaymentBalance *float64       `gorm:"total_payment_balance" json:"total_payment_balance"`
	TotalPayment *float64 `gorm:"total_payment" json:"total_payment"`
	// RemainingAmount     *float64       `gorm:"remaining_amount" json:"remaining_amount"`
	VerificationStatus *int64 `gorm:"verification_status" json:"verification_status"`
	// CreatedBy           *int64         `gorm:"column:created_by" json:"created_by"`
	// CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	// UpdatedBy *int64     `gorm:"column:updated_by" json:"updated_by"`
	// UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	// IsDel               bool           `gorm:"column:is_del" json:"is_del"`
	// DeletedBy           *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	// DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsApproved bool       `gorm:"column:is_approved" json:"is_approved"`
	ApprovedBy *string    `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt *time.Time `gorm:"column:approved_at" json:"approved_at"`
}

func (ArBranchSettlement) TableName() string {
	return "inv.gr_branch_payment"
}
