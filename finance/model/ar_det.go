package model

import "time"

type ArPaymentRead struct {
	DepositPaymentID   int64      `gorm:"column:deposit_payment_id;primaryKey" json:"deposit_payment_id"`
	VisitDate          time.Time  `gorm:"column:visit_date" json:"visit_date"`
	DepositDate        time.Time  `gorm:"column:deposit_date" json:"deposit_date"`
	DepositNo          string     `gorm:"column:deposit_no" json:"deposit_no"`
	EmpId              int64      `gorm:"column:emp_id" json:"emp_id"`
	EmpCode            *string    `gorm:"column:emp_code" json:"emp_code"`
	EmpName            *string    `gorm:"column:emp_name" json:"emp_name"`
	EmpGrpId           int64      `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode         *string    `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGrpName         *string    `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	VerificationStatus int64      `gorm:"column:verification_status" json:"verification_status"`
	TotalPayment       float64    `gorm:"column:total_payment" json:"total_payment"`
	RemainingPayment   float64    `gorm:"column:remaining_payment" json:"remaining_payment"`
	PaymentMethod      int64      `gorm:"column:payment_method" json:"payment_method"`
	Amount             float64    `gorm:"column:amount" json:"amount"`
	VerifiedBy         *int64     `gorm:"column:verified_by" json:"verified_by"`
	VerifiedByName     *string    `gorm:"column:verified_by_name" json:"verified_by_name"`
	VerifiedDate       *time.Time `gorm:"column:verified_date" json:"verified_date"`
	Reason             *string    `gorm:"column:reason" json:"reason"`
	AdditionalInfo     *string    `gorm:"column:additional_info" json:"additional_info"`
}

func (ArPaymentRead) TableName() string {
	return "acf.deposit_payment"
}

type CollectionDet struct {
	CustID          string    `gorm:"column:cust_id" json:"cust_id"`
	CollectionDetID *int64    `gorm:"column:collection_det_id;primaryKey" json:"collection_det_id"`
	CollectionNo    string    `gorm:"column:collection_no" json:"collection_no"`
	InvoiceNo       string    `gorm:"column:invoice_no" json:"invoice_no"`
	SalesmanID      *int64    `gorm:"column:salesman_id" json:"salesman_id"`
	InvoiceAmount   float64   `gorm:"column:invoice_amount" json:"invoice_amount" default:"0"`
	RemainingAmount float64   `gorm:"column:remaining_amount" json:"remaining_amount" default:"0"`
	PaidAmount      float64   `gorm:"column:paid_amount" json:"paid_amount" default:"0"`
	CreatedBy       *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

type CollectionDetList struct {
	CustID              string     `gorm:"column:cust_id" json:"cust_id"`
	CollectionDetID     *int64     `gorm:"column:collection_det_id;primaryKey" json:"collection_det_id"`
	CollectionNo        string     `gorm:"column:collection_no" json:"collection_no"`
	SalesOrder          string     `gorm:"column:sales_order" json:"sales_order"`
	InvoiceDate         *time.Time `gorm:"column:invoice_date" json:"invoice_date"`
	DueDate             *time.Time `gorm:"column:due_date" json:"due_date"`
	SalesmanId          int64      `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName        string     `gorm:"column:salesman_name" json:"salesman_name"`
	SalesmanCode        string     `gorm:"column:salesman_code" json:"salesman_code"`
	OutletId            int64      `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode          string     `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName          string     `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceNo           string     `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceAmount       *float64   `gorm:"column:invoice_amount" json:"invoice_amount"`
	RemainingAmount     *float64   `gorm:"column:remaining_amount" json:"remaining_amount"`
	PaidAmount          *float64   `gorm:"column:paid_amount" json:"paid_amount"`
	TotalInvoicePayment *float64   `gorm:"column:total_invoice_amount" json:"total_invoice_amount"`
	CreatedBy           *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedByName       *string    `gorm:"column:created_by_name" json:"created_by_name"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (CollectionDet) TableName() string {
	return "acf.collection_det"
}

func (CollectionDetList) TableName() string {
	return "acf.collection_det"
}
