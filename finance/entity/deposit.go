package entity

import "time"

type DepositQueryFilter struct {
	GeneralQueryFilter
	CollectionNo []string `query:"collection_no"`
	Status       []int    `query:"status_no"`
	DepositNo    []string `query:"deposit_no"`
	InvoiceNo    []string `query:"invoice_no"`
	DocumentNo   []string `query:"document_no"`
	Type         []int    `query:"type"`
}

type DepositNumberListQueryFilter struct {
	CustId       string `json:"-"`
	Query        string `query:"q"`
	Page         int    `query:"page" validate:"required,min=1"`
	Limit        int    `query:"limit" validate:"required,min=1,max=9999"`
	Sort         string `query:"sort" validate:"required"`
	CollectorIDs []int  `query:"collector_id" validate:"required,min=1,dive,gt=0"`
}

type DepositNumberListItemResponse struct {
	DepositNo   string `json:"deposit_no"`
	CollectorID int    `json:"collector_id"`
	DepositDate string `json:"deposit_date"`
}

type DepositNumberListPagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

type DepositResponse struct {
	CustID              string     `json:"cust_id"`
	DepositNo           string     `json:"deposit_no"`
	DepositDate         *string    `json:"deposit_date"`
	CollectionNo        *string    `json:"collection_no"`
	CollectionDate      *string    `json:"collection_date"`
	EmpGrpID            *int       `json:"emp_grp_id"`
	EmpGrpName          *string    `json:"emp_grp_name"`
	EmpID               *int       `json:"emp_id"`
	EmpName             *string    `json:"emp_name"`
	EmpCode             *string    `json:"emp_code"`
	OutletGroupID       *int       `json:"outlet_group_id"`
	OutletGroupCode     *string    `json:"outlet_group_code"`
	OutletGroupName     *string    `json:"outlet_group_name"`
	SalesmanID          *int       `json:"salesman_id"`
	SalesmanName        *string    `json:"salesman_name"`
	InvoiceDateFrom     *string    `json:"invoice_date_from"`
	InvoiceDateTo       *string    `json:"invoice_date_to"`
	DueDateFrom         *string    `json:"due_date_from"`
	DueDateTo           *string    `json:"due_date_to"`
	DepositStatus       int        `json:"deposit_status"`
	DepositStatusName   string     `json:"deposit_status_name"`
	RemainingAmount     float64    `json:"remaining_amount"`
	TotalDiscount       float64    `json:"total_discount"`
	TotalMaterai        float64    `json:"total_materai"`
	TotalPaymentBalance float64    `json:"total_payment_balance"`
	ExpenseTotal        float64    `json:"expense_total"`
	TotalPayment        float64    `json:"total_payment"`
	IsApproved          bool       `json:"is_approved"`
	ApprovedBy          *int       `json:"approved_by"`
	ApprovedAt          *time.Time `json:"approved_at"`
	CreatedBy           *int       `json:"created_by"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedBy           *int       `json:"updated_by"`
	UpdatedByName       *int       `json:"updated_by_name"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedBy           *int       `json:"deleted_by"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

type DepositPaymentInvoice struct {
	DepositPayment
	SalesmanId     *int64  `json:"salesman_id"`
	SalesmanCode   *string `json:"salesman_code"`
	SalesmanName   *string `json:"salesman_name"`
	OutletID       *int64  `json:"outlet_id"`
	OutletCode     *string `json:"outlet_code"`
	OutletName     *string `json:"outlet_name"`
	InvoiceDate    string  `json:"invoice_date"`
	Materai        float64 `json:"materai"`
	Discount       float64 `json:"discount"`
	PaymentBalance float64 `json:"payment_balance"`
	TotalPayment   float64 `json:"total_payment"`
	Notes          string  `json:"notes"`
}

type DepositDetailResponse struct {
	DepositResponse
	Details []DepositDetail         `json:"details"`
	Cash    []DepositPaymentInvoice `json:"cash"`
	Cek     []DepositPaymentInvoice `json:"cek"`
	Trasfer []DepositPaymentInvoice `json:"transfer"`
	Return  []DepositPaymentInvoice `json:"return"`
	CNDN    []DepositPaymentInvoice `json:"cndn"`
	Expense []DepositExpense        `json:"expense"`
}

type DepositDetailReportResponse struct {
	DepositResponse
	Details []DepositDetailReport `json:"details"`
}

type DepositPayment struct {
	DepositPaymentID int64   `json:"deposit_payment_id"`
	DepositNo        string  `json:"deposit_no"`
	InvoiceNo        string  `json:"invoice_no"`
	SalesmanId       int64   `json:"salesman_id"`
	PayType          int16   `json:"pay_type"`
	DocumentNo       string  `json:"document_no"`
	Balance          float64 `json:"balance"`
	PaymentAmount    float64 `json:"payment_amount"`
}

type DepositDetail struct {
	DepositDetailID  int64            `json:"deposit_detail_id"`
	DepositNo        string           `json:"deposit_no"`
	InvoiceNo        string           `json:"invoice_no"`
	Discount         float64          `json:"discount"`
	PaymentBalance   float64          `json:"payment_balance"`
	Materai          float64          `json:"materai"`
	InvoiceAmount    float64          `json:"invoice_amount"`
	TotalPayment     float64          `json:"total_payment"`
	RemainingPayment float64          `json:"remaining_payment"`
	IsCollection     bool             `json:"is_collection"`
	PayType          int16            `json:"pay_type"`
	InvoiceDate      string           `json:"invoice_date"`
	DocumentNo       string           `json:"document_no"`
	Notes            string           `json:"notes"`
	RoNo             string           `json:"ro_no"`
	OutletID         int              `json:"outlet_id"`
	OutletCode       string           `json:"outlet_code"`
	OutletName       string           `json:"outlet_name"`
	DueDate          string           `json:"due_date"`
	SalesmanId       int              `json:"salesman_id"`
	SalesmanName     string           `json:"salesman_name"`
	PaidAmount       float64          `json:"paid_amount"`
	Payment          []DepositPayment `json:"payment" validate:"dive,required"`
}

type DepositDetailReport struct {
	DepositDetailID  int64   `json:"deposit_detail_id"`
	DepositNo        string  `json:"deposit_no"`
	InvoiceNo        string  `json:"invoice_no"`
	Discount         float64 `json:"discount"`
	PaymentBalance   float64 `json:"payment_balance"`
	Materai          float64 `json:"materai"`
	InvoiceAmount    float64 `json:"invoice_amount"`
	TotalPayment     float64 `json:"total_payment"`
	RemainingPayment float64 `json:"remaining_payment"`
	IsCollection     bool    `json:"is_collection"`
	PayType          int16   `json:"pay_type"`
	InvoiceDate      string  `json:"invoice_date"`
	DocumentNo       string  `json:"document_no"`
	OutletID         *int64  `json:"outlet_id"`
	OutletCode       *string `json:"outlet_code"`
	OutletName       *string `json:"outlet_name"`

	NilaiTunai    *float64 `json:"nilai_tunai"`
	NoTunai       *string  `json:"no_tunai"`
	NilaiRetur    *float64 `json:"nilai_retur"`
	NoRetur       *string  `json:"no_retur"`
	NilaiTransfer *float64 `json:"nilai_transfer"`
	NoTransfer    *string  `json:"no_transfer"`
	NilaiCek      *float64 `json:"nilai_cek"`
	NoCek         *string  `json:"no_cek"`
	NilaiCndn     *float64 `json:"nilai_cndn"`
	NoCndn        *string  `json:"no_cndn"`
	Notes         string   `json:"notes"`

	// Payment          []DepositPayment `json:"payment" validate:"required,dive,required"`
}

type CreateDepositBodyByCollection struct {
	CustID              string           `json:"cust_id"`
	DepositNo           string           `json:"deposit_no"`
	DepositDate         string           `json:"deposit_date"`
	CollectionNo        *string          `json:"collection_no"`
	EmpGrpID            *int             `json:"emp_grp_id"`
	EmpID               *int             `json:"emp_id"`
	DepositStatus       int              `json:"deposit_status"`
	RemainingAmount     float64          `json:"remaining_amount"`
	TotalDiscount       float64          `json:"total_discount"`
	TotalMaterai        float64          `json:"total_materai"`
	TotalPaymentBalance float64          `json:"total_payment_balance"`
	TotalPayment        float64          `json:"total_payment"`
	Details             []DepositDetail  `json:"detail" validate:"required,dive,required"`
	Expense             []DepositExpense `json:"expense"`
	CreatedBy           *int64           `json:"created_by"`
	CreatedAt           *time.Time       `json:"created_at"`
	UpdatedBy           *int64           `json:"updated_by"`
	UpdatedAt           *time.Time       `json:"updated_at"`
}

type CreateDepositBodyByInvoice struct {
	CustID              string           `json:"cust_id"`
	DepositNo           string           `json:"deposit_no"`
	DepositDate         string           `json:"deposit_date"`
	EmpGrpID            *int             `json:"emp_grp_id"`
	EmpID               *int             `json:"emp_id"`
	SalesmanID          *int             `json:"salesman_id"`
	InvoiceDateFrom     string           `json:"invoice_date_from"`
	InvoiceDateTo       string           `json:"invoice_date_to"`
	DueDateFrom         string           `json:"due_date_from"`
	DueDateTo           string           `json:"due_date_to"`
	DepositStatus       int              `json:"deposit_status"`
	RemainingAmount     float64          `json:"remaining_amount"`
	TotalDiscount       float64          `json:"total_discount"`
	TotalMaterai        float64          `json:"total_materai"`
	TotalPaymentBalance float64          `json:"total_payment_balance"`
	TotalPayment        float64          `json:"total_payment"`
	Details             []DepositDetail  `json:"detail" validate:"required,dive,required"`
	Expense             []DepositExpense `json:"expense"`
	CreatedBy           *int64           `json:"created_by"`
	CreatedAt           *time.Time       `json:"created_at"`
	UpdatedBy           *int64           `json:"updated_by"`
	UpdatedAt           *time.Time       `json:"updated_at"`
}

type UpdateDepositBodyCollection struct {
	CustID              string           `json:"cust_id"`
	DepositNo           string           `json:"deposit_no"`
	DepositDate         string           `json:"deposit_date"`
	CollectionNo        *string          `json:"collection_no"`
	EmpGrpID            *int             `json:"emp_grp_id"`
	EmpID               *int             `json:"emp_id"`
	DepositStatus       int              `json:"deposit_status"`
	RemainingAmount     float64          `json:"remaining_amount"`
	TotalDiscount       float64          `json:"total_discount"`
	TotalMaterai        float64          `json:"total_materai"`
	TotalPaymentBalance float64          `json:"total_payment_balance"`
	TotalPayment        float64          `json:"total_payment"`
	Detail              []DepositDetail  `json:"detail" validate:"required,dive,required"`
	Expense             []DepositExpense `json:"expense"`
	UpdatedBy           *int64           `json:"updated_by"`
	UpdatedAt           *time.Time       `json:"updated_at"`
}

type UpdateDepositBodyInvoice struct {
	CustID              string           `json:"cust_id"`
	DepositNo           string           `json:"deposit_no"`
	DepositDate         string           `json:"deposit_date"`
	EmpGrpID            *int             `json:"emp_grp_id"`
	EmpID               *int             `json:"emp_id"`
	SalesmanID          *int             `json:"salesman_id"`
	InvoiceDateFrom     string           `json:"invoice_date_from"`
	InvoiceDateTo       string           `json:"invoice_date_to"`
	DueDateFrom         string           `json:"due_date_from"`
	DueDateTo           string           `json:"due_date_to"`
	DepositStatus       int              `json:"deposit_status"`
	RemainingAmount     float64          `json:"remaining_amount"`
	TotalDiscount       float64          `json:"total_discount"`
	TotalMaterai        float64          `json:"total_materai"`
	TotalPaymentBalance float64          `json:"total_payment_balance"`
	TotalPayment        float64          `json:"total_payment"`
	Detail              []DepositDetail  `json:"detail" validate:"required,dive,required"`
	Expense             []DepositExpense `json:"expense"`
	UpdatedBy           *int64           `json:"updated_by"`
	UpdatedAt           *time.Time       `json:"updated_at"`
}

type DetailDepositParams struct {
	DepositNo string `params:"deposit_no" validate:"required"`
}
type DeleteDepositParams struct {
	DepositNo string `params:"deposit_no" validate:"required"`
}
type UpdateDepositParams struct {
	DepositNo string `params:"deposit_no" validate:"required"`
}
