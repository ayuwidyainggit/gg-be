package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Ar struct {
	CustID     string         `gorm:"column:cust_id" json:"cust_id"`
	ArNo       string         `gorm:"column:ar_no;primaryKey" json:"ar_no"`
	ArDate     *time.Time     `gorm:"column:ar_date" json:"ar_date"`
	TrCode     *string        `gorm:"column:tr_code" json:"tr_code"`
	SalesmanID *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	RefNo      *string        `gorm:"column:ref_no" json:"ref_no"`
	DataStatus *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted   *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt   *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (Ar) TableName() string {
	return "acf.ar"
}

func (m *Ar) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	m.ArNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

type ArList struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	SalesmanId    int64      `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode  *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName  *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID      int64      `gorm:"outlet_id" json:"outlet_id"`
	OutletCode    *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceAmount float64    `gorm:"invoice_amount" json:"invoice_amount"`
	InvoiceNo     string     `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time `gorm:"invoice_date" json:"invoice_date"`
	DueDate       *time.Time `gorm:"due_date" json:"due_date"`
}

func (ArList) TableName() string {
	return "sls.order"
}

type ArRead struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	SalesmanId    int64      `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode  *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName  *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID      int64      `gorm:"outlet_id" json:"outlet_id"`
	OutletCode    *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string    `gorm:"column:outlet_name" json:"outlet_name"`
	OrderNo       *string    `gorm:"order_no" json:"order_no"`
	InvoiceAmount float64    `gorm:"invoice_amount" json:"invoice_amount"`
	InvoiceNo     string     `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate   *time.Time `gorm:"invoice_date" json:"invoice_date"`
	DueDate       *time.Time `gorm:"due_date" json:"due_date"`
}

func (ArRead) TableName() string {
	return "sls.order"
}

type Collection struct {
	CustID          string         `gorm:"column:cust_id" json:"cust_id"`
	CollectionNo    string         `gorm:"column:collection_no" json:"collection_no"`
	CollectionDate  *time.Time     `gorm:"column:collection_date" json:"collection_date"`
	EmpID           *int64         `gorm:"column:emp_id" json:"emp_id"`
	OtGrpID         *int64         `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	Notes           *string        `gorm:"column:notes" json:"notes"`
	TotalAmount     float64        `gorm:"column:total_amount" json:"total_amount" default:"0"`
	RemainingAmount float64        `gorm:"column:remaining_amount" json:"remaining_amount" default:"0"`
	InvoiceDateFrom *time.Time     `gorm:"column:invoice_date_from" json:"invoice_date_from"`
	InvoiceDateTo   *time.Time     `gorm:"column:invoice_date_to" json:"invoice_date_to"`
	DueDateFrom     *time.Time     `gorm:"column:due_date_from" json:"due_date_from"`
	DueDateTo       *time.Time     `gorm:"column:due_date_to" json:"due_date_to"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel           bool           `gorm:"column:is_del" json:"is_del" default:"false"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPrinted       bool           `gorm:"column:is_printed" json:"is_printed" default:"false"`
	PrintedBy       *int64         `gorm:"column:printed_by" json:"printed_by"`
	PrintedAt       *time.Time     `gorm:"column:printed_at" json:"printed_at"`
}

func (Collection) TableName() string {
	return "acf.collection"
}

// func (m *Collection) BeforeCreate(trx *gorm.DB) (err error) {
// 	now := time.Now()
// 	intTmpsStr := now.UnixNano() / int64(time.Millisecond)
// 	m.CollectionNo = strconv.Itoa(int(intTmpsStr))
// 	m.CreatedAt = now
// 	return nil
// }

type CollectionNo struct {
	CollectionNo string `gorm:"column:get_no_fn"`
}

func (m *Collection) BeforeCreate(trx *gorm.DB) (err error) {
	var collectionNo CollectionNo
	trCode := "CL"
	collectionDateStr := m.CollectionDate.Format("2006-01-02")
	collectionDateSubtr := collectionDateStr[2:4] + collectionDateStr[5:7] + collectionDateStr[8:10]

	queryStr := fmt.Sprintf(`SELECT
	TRIM(to_char(COALESCE(MAX(TO_NUMBER(SUBSTR(collection_no,9,4),'9999')),0)+1, '0000')) AS get_no_fn
	FROM acf.collection
	WHERE substr(collection_no,3,6) = '%v' AND cust_id = '%v'`, collectionDateSubtr, strings.ToUpper(m.CustID))
	err = trx.Raw(queryStr).Scan(&collectionNo).Error
	if err != nil {
		return err
	}
	// log.Println("collectionNo:", collectionNo.CollectionNo)

	m.CollectionNo = trCode + collectionDateSubtr + collectionNo.CollectionNo
	// log.Println("m.CollectionNo:", m.CollectionNo)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy

	return nil
}

type CollectionList struct {
	CustID          string         `gorm:"column:cust_id" json:"cust_id"`
	CollectionNo    string         `gorm:"column:collection_no;primaryKey" json:"collection_no"`
	CollectionDate  *time.Time     `gorm:"column:collection_date" json:"collection_date"`
	EmpID           *int64         `gorm:"column:emp_id" json:"emp_id"`
	EmpCode         *string        `gorm:"column:emp_code" json:"emp_code"`
	EmpName         *string        `gorm:"column:emp_name" json:"emp_name"`
	OtGrpID         *int64         `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	OtGrpCode       *string        `gorm:"column:ot_grp_code" json:"ot_grp_code"`
	OtGrpName       *string        `gorm:"column:ot_grp_name" json:"ot_grp_name"`
	EmpGrpID        *int64         `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGrpCode      *string        `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGrpName      *string        `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	Notes           *string        `gorm:"column:notes" json:"notes"`
	TotalAmount     *float64       `gorm:"column:total_amount" json:"total_amount"`
	RemainingAmount *float64       `gorm:"column:remaining_amount" json:"remaining_amount"`
	InvoiceDateFrom *time.Time     `gorm:"column:invoice_date_from" json:"invoice_date_from"`
	InvoiceDateTo   *time.Time     `gorm:"column:invoice_date_to" json:"invoice_date_to"`
	DueDateFrom     *time.Time     `gorm:"column:due_date_from" json:"due_date_from"`
	DueDateTo       *time.Time     `gorm:"column:due_date_to" json:"due_date_to"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedByName   *string        `gorm:"column:created_by_name" json:"created_by_name"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName   *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel           bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedByName   *string        `gorm:"column:deleted_by_name" json:"deleted_by_name"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPrinted       *bool          `gorm:"column:is_printed" json:"is_printed"`
	PrintedBy       *int64         `gorm:"column:printed_by" json:"printed_by"`
	PrintedByName   *string        `gorm:"column:printed_by_name" json:"printed_by_name"`
	PrintedAt       *time.Time     `gorm:"column:printed_at" json:"printed_at"`
}

func (CollectionList) TableName() string {
	return "acf.collection"
}

type EmployeeGroup struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	EmpGroupId    int             `gorm:"column:emp_grp_id" json:"emp_grp_id"`
	EmpGroupCode  string          `gorm:"column:emp_grp_code" json:"emp_grp_code"`
	EmpGroupName  string          `gorm:"column:emp_grp_name" json:"emp_grp_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (EmployeeGroup) TableName() string {
	return "mst.m_emp_group"
}

type OutletGroupFilter struct {
	CustId          string          `gorm:"column:cust_id" json:"cust_id"`
	OutletGroupId   int             `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	OutletGroupCode string          `gorm:"column:ot_grp_code" json:"ot_grp_code"`
	OutletGroupName string          `gorm:"column:ot_grp_name" json:"ot_grp_name"`
	IsActive        bool            `gorm:"column:is_active" json:"is_active"`
	IsDel           bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy       *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt       *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy       *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName   *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy       *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletGroupFilter) TableName() string {
	return "sls.order"
}

type Employee struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	EmpId         int             `gorm:"column:emp_id" json:"emp_id"`
	EmpCode       string          `gorm:"column:emp_code" json:"emp_code"`
	EmpName       string          `gorm:"column:emp_name" json:"emp_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Employee) TableName() string {
	return "mst.m_employee"
}

type SalesmanFilter struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	SalesmanId    int             `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  string          `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  string          `gorm:"column:salesman_name" json:"salesman_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SalesmanFilter) TableName() string {
	return "sls.order"
}

type InvoiceList struct {
	CustID          string     `gorm:"cust_id" json:"cust_id"`
	RoNo            string     `gorm:"ro_no" json:"ro_no"`
	SalesmanId      *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesmanCode    *string    `gorm:"salesman_code" json:"salesman_code"`
	SalesmanName    *string    `gorm:"salesman_name" json:"salesman_name"`
	OutletID        *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode      *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName      *string    `gorm:"column:outlet_name" json:"outlet_name"`
	InvoiceAmount   *float64   `gorm:"invoice_amount" json:"invoice_amount"`
	RemainingAmount *float64   `gorm:"remaining_amount" json:"remaining_amount"`
	PaidAmount      *float64   `gorm:"paid_amount" json:"paid_amount"`
	InvoiceNo       *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate     *time.Time `gorm:"invoice_date" json:"invoice_date"`
	DueDate         *time.Time `gorm:"due_date" json:"due_date"`
}

func (InvoiceList) TableName() string {
	return "sls.order"
}

type OutletFilter struct {
	CustId        string          `gorm:"column:cust_id" json:"cust_id"`
	OutletId      int             `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    string          `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string          `gorm:"column:outlet_name" json:"outlet_name"`
	IsActive      bool            `gorm:"column:is_active" json:"is_active"`
	IsDel         bool            `gorm:"column:is_del" json:"is_del"`
	CreatedBy     *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt     *time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" db:"updated_by_name"`
	DeletedBy     *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OutletFilter) TableName() string {
	return "sls.order"
}

type Collector struct {
	CustId   string `gorm:"column:cust_id" json:"cust_id"`
	EmpId    int    `gorm:"column:emp_id" json:"emp_id"`
	EmpCode  string `gorm:"column:emp_code" json:"emp_code"`
	EmpName  string `gorm:"column:emp_name" json:"emp_name"`
	IsActive bool   `gorm:"column:is_active" json:"is_active"`
	IsDel    bool   `gorm:"column:is_del" json:"is_del"`
}

func (Collector) TableName() string {
	return "acf.collection"
}

type InvoicePaidAmount struct {
	PaidAmount float64 `gorm:"column:paid_amount" json:"paid_amount"`
}

func (InvoicePaidAmount) TableName() string {
	return "acf.deposit_detail"
}

type LastApprovedDeposit struct {
	DepositNo  string     `gorm:"column:deposit_no;primaryKey" json:"deposit_no"`
	InvoiceNo  string     `gorm:"column:invoice_no" json:"invoice_no"`
	ApprovedAt *time.Time `gorm:"column:approved_at" json:"approved_at"`
}

func (LastApprovedDeposit) TableName() string {
	return "acf.deposit_detail"
}
