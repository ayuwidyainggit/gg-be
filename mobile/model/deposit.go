package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	DepositStatusNeedReview = 1 // need review
	DepositStatusApproved   = 2 // approved
)

type Deposit struct {
	CustID              string     `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	DepositNo           string     `gorm:"column:deposit_no;type:varchar(30);not null" json:"deposit_no"`
	DepositDate         *time.Time `gorm:"column:deposit_date;type:date" json:"deposit_date"`
	CollectionNo        *string    `gorm:"column:collection_no;type:varchar(30)" json:"collection_no"`
	EmpGrpID            *int       `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpID               *int       `gorm:"column:emp_id" json:"emp_id"`
	SalesmanID          *int       `gorm:"column:salesman_id" json:"salesman_id"`
	InvoiceDateFrom     *time.Time `gorm:"column:invoice_date_from;type:date" json:"invoice_date_from"`
	InvoiceDateTo       *time.Time `gorm:"column:invoice_date_to;type:date" json:"invoice_date_to"`
	DueDateFrom         *time.Time `gorm:"column:due_date_from;type:date" json:"due_date_from"`
	DueDateTo           *time.Time `gorm:"column:due_date_to;type:date" json:"due_date_to"`
	DepositStatus       int        `gorm:"column:deposit_status;type:int;not null" json:"deposit_status"`
	CreatedBy           *int       `gorm:"column:created_by" json:"created_by"`
	CreatedAt           *time.Time `gorm:"column:created_at;type:timestamptz" json:"created_at"`
	UpdatedBy           *int       `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
	DeletedBy           *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           *time.Time `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
	RemainingAmount     float64    `gorm:"column:remaining_amount;type:numeric(20,4);default:0" json:"remaining_amount"`
	TotalDiscount       float64    `gorm:"column:total_discount;type:numeric(20,4);default:0" json:"total_discount"`
	TotalMaterai        float64    `gorm:"column:total_materai;type:numeric(20,4);default:0" json:"total_materai"`
	TotalPaymentBalance float64    `gorm:"column:total_payment_balance;type:numeric(20,4);default:0" json:"total_payment_balance"`
	TotalPayment        float64    `gorm:"column:total_payment;type:numeric(20,4);default:0" json:"total_payment"`
	IsApproved          bool       `gorm:"column:is_approved;type:bool;default:false;not null" json:"is_approved"`
	ApprovedBy          *int       `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt          *time.Time `gorm:"column:approved_at;type:timestamp" json:"approved_at"`
}

func (m *Deposit) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = &now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (Deposit) TableName() string {
	return "acf.deposit"
}

type DepositList struct {
	CustID              string     `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	DepositNo           string     `gorm:"column:deposit_no;type:varchar(30);not null" json:"deposit_no"`
	DepositDate         *time.Time `gorm:"column:deposit_date;type:date" json:"deposit_date"`
	CollectionNo        *string    `gorm:"column:collection_no;type:varchar(30)" json:"collection_no"`
	CollectionDate      *time.Time `gorm:"column:collection_date;type:varchar(30)" json:"collection_date"`
	EmpGrpID            *int       `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpName          *string    `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	OutletGroupID       *int       `gorm:"column:ot_grp_id" json:"outlet_group_id"`
	OutletGroupCode     *string    `gorm:"column:ot_grp_code" json:"outlet_group_code"`
	OutletGroupName     *string    `gorm:"column:ot_grp_name" json:"outlet_group_name"`
	EmpID               *int       `gorm:"column:emp_id" json:"emp_id"`
	EmpName             *string    `gorm:"column:emp_name" json:"emp_name"`
	EmpCode             *string    `gorm:"column:emp_code" json:"emp_code"`
	SalesmanID          *int       `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName        *string    `gorm:"column:salesman_name" json:"salesman_name"`
	InvoiceDateFrom     *time.Time `gorm:"column:invoice_date_from;type:date" json:"invoice_date_from"`
	InvoiceDateTo       *time.Time `gorm:"column:invoice_date_to;type:date" json:"invoice_date_to"`
	DueDateFrom         *time.Time `gorm:"column:due_date_from;type:date" json:"due_date_from"`
	DueDateTo           *time.Time `gorm:"column:due_date_to;type:date" json:"due_date_to"`
	DepositStatus       int        `gorm:"column:deposit_status;type:int;not null" json:"deposit_status"`
	RemainingAmount     float64    `gorm:"column:remaining_amount;type:numeric(20,4);default:0" json:"remaining_amount"`
	TotalDiscount       float64    `gorm:"column:total_discount;type:numeric(20,4);default:0" json:"total_discount"`
	TotalMaterai        float64    `gorm:"column:total_materai;type:numeric(20,4);default:0" json:"total_materai"`
	TotalPaymentBalance float64    `gorm:"column:total_payment_balance;type:numeric(20,4);default:0" json:"total_payment_balance"`
	TotalPayment        float64    `gorm:"column:total_payment;type:numeric(20,4);default:0" json:"total_payment"`
	IsApproved          bool       `gorm:"column:is_approved;type:bool;default:false;not null" json:"is_approved"`
	ApprovedBy          *int       `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt          *time.Time `gorm:"column:approved_at;type:timestamp" json:"approved_at"`
	CreatedBy           *int       `gorm:"column:created_by" json:"created_by"`
	CreatedAt           *time.Time `gorm:"column:created_at;type:timestamptz" json:"created_at"`
	UpdatedBy           *int       `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt           *time.Time `gorm:"column:updated_at;type:timestamptz" json:"updated_at"`
	DeletedBy           *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt           *time.Time `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (DepositList) TableName() string {
	return "acf.deposit"
}

type DepositDetail struct {
	CustID           string  `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	DepositDetailID  int     `gorm:"column:deposit_detail_id;type:int8;default:nextval('acf.deposit_detail_id_seq'::regclass);not null" json:"deposit_detail_id"`
	DepositNo        string  `gorm:"column:deposit_no;type:varchar(30);not null" json:"deposit_no"`
	InvoiceNo        string  `gorm:"column:invoice_no;type:varchar(30);not null" json:"invoice_no"`
	Discount         float64 `gorm:"column:discount;type:numeric(20,4);default:0" json:"discount"`
	PaymentBalance   float64 `gorm:"column:payment_balance;type:numeric(20,4);default:0" json:"payment_balance"`
	Materai          float64 `gorm:"column:materai;type:numeric(20,4);default:0" json:"materai"`
	InvoiceAmount    float64 `gorm:"column:invoice_amount;type:numeric(20,4);default:0" json:"invoice_amount"`
	TotalPayment     float64 `gorm:"column:total_payment;type:numeric(20,4);default:0" json:"total_payment"`
	RemainingPayment float64 `gorm:"column:remaining_payment;type:numeric(20,4);default:0" json:"remaining_payment"`
	IsCollection     bool    `gorm:"column:is_collection;type:bool;default:true" json:"is_collection"`
}

func (DepositDetail) TableName() string {
	return "acf.deposit_detail" // Ganti dengan nama tabel yang sesuai
}

type DepositDetailList struct {
	DepositDetail
	RoNo        *string    `gorm:"ro_no" json:"ro_no"`
	OutletID    *int       `gorm:"outlet_id" json:"outlet_id"`
	OutletCode  *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName  *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceDate *time.Time `gorm:"invoice_date" json:"invoice_date"`
	DueDate     *time.Time `gorm:"due_date" json:"due_date"`
	SalesmanId  *int       `gorm:"salesman_id" json:"salesman_id"`
	SalesName   *string    `gorm:"sales_name" json:"sales_name"`
}

func (DepositDetailList) TableName() string {
	return "acf.deposit_detail" // Ganti dengan nama tabel yang sesuai
}

type DepositPayment struct {
	CustID           string   `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	DepositPaymentID int      `gorm:"column:deposit_payment_id;type:int8;default:nextval('acf.deposit_payment_id_seq'::regclass);not null" json:"deposit_payment_id"`
	DepositNo        string   `gorm:"column:deposit_no;type:varchar(30);not null" json:"deposit_no"`
	InvoiceNo        string   `gorm:"column:invoice_no;type:varchar(30);not null" json:"invoice_no"`
	PayType          int16    `gorm:"column:pay_type;type:int2" json:"pay_type"`
	DocumentNo       string   `gorm:"column:document_no;type:varchar(30)" json:"document_no"`
	Balance          float64  `gorm:"column:balance;type:numeric(20,4);default:0" json:"balance"`
	PaymentAmount    float64  `gorm:"column:payment_amount;type:numeric(20,4);default:0" json:"payment_amount"`
	Discount         *float64 `gorm:"column:discount;type:numeric(20,4);default:0" json:"discount"`
	Materai          *float64 `gorm:"column:materai;type:numeric(20,4);default:0" json:"materai"`
}

// TableName sets the insert table name for this struct type
func (DepositPayment) TableName() string {
	return "acf.deposit_payment"
}

type DepositPaymentImage struct {
	DepositImageID int    `gorm:"column:deposit_image_id;type:int8;not null" json:"deposit_image_id"` // Change to int
	DepositNo      string `gorm:"column:deposit_no;type:varchar(30);not null" json:"deposit_no"`
	InvoiceNo      string `gorm:"column:invoice_no;type:varchar(30);not null" json:"invoice_no"`
	ImageUrl       string `gorm:"column:image_url;type:text;not null" json:"image_url"`
}

// TableName sets the insert table name for this struct type
func (DepositPaymentImage) TableName() string {
	return "acf.deposit_image"
}

type DepositPaymentInvoice struct {
	DepositPayment
	SalesmanId   *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID     *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode   *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName   *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceDate  *time.Time `gorm:"invoice_date" json:"invoice_date"`
}

// TableName sets the insert table name for this struct type
func (DepositPaymentInvoice) TableName() string {
	return "acf.deposit_payment"
}

type CollectionNoPayment struct {
	CustID                 string     `gorm:"column:cust_id;type:varchar(10);not null" json:"cust_id"`
	CollectionNoPaymentID  int        `gorm:"column:collection_no_payment_id;type:serial;primaryKey" json:"collection_no_payment_id"`
	SalesmanID             *int64     `gorm:"column:salesman_id;type:bigint" json:"salesman_id,omitempty"`
	OutletID               *int64     `gorm:"column:outlet_id;type:bigint" json:"outlet_id,omitempty"`
	CollectionNo           string     `gorm:"column:collection_no;type:varchar(30);not null" json:"collection_no"`
	InvoiceNo              string     `gorm:"column:invoice_no;type:varchar(30);not null" json:"invoice_no"`
	MissedPaymentReasonsID *int64     `gorm:"column:missed_payment_reasons_id;type:bigint" json:"missed_payment_reasons_id,omitempty"`
	Reason                 *string    `gorm:"column:reason;type:varchar(255)" json:"reason,omitempty"`
	PaymentDate            *time.Time `gorm:"column:payment_date;type:timestamp with time zone" json:"payment_date,omitempty"`
	CreatedBy              *int       `gorm:"column:created_by;type:int" json:"created_by,omitempty"`
	CreatedAt              *time.Time `gorm:"column:created_at;type:timestamp without time zone" json:"created_at,omitempty"`
}

// TableName sets the table name for the struct
func (CollectionNoPayment) TableName() string {
	return "acf.collection_no_payment"
}
