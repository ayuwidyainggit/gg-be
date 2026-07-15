package model

import (
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	CustID        string     `gorm:"cust_id" json:"cust_id"`
	RoNo          string     `gorm:"ro_no" json:"ro_no"`
	SalesmanId    *int64     `gorm:"salesman_id" json:"salesman_id"`
	WhId          *int64     `gorm:"wh_id" json:"wh_id"`
	RoDate        *time.Time `gorm:"ro_date" json:"ro_date"`
	ValDate       *time.Time `gorm:"val_date" json:"val_date"`
	OutletID      *int64     `gorm:"outlet_id" json:"outlet_id"`
	DeliveryDate  *time.Time `gorm:"delivery_date" json:"delivery_date"`
	CollectionNo  *string    `gorm:"order_no" json:"order_no"`
	PoNo          *string    `gorm:"po_no" json:"po_no"`
	VehicleNo     *string    `gorm:"vehicle_no" json:"vehicle_no"`
	PayType       *int64     `gorm:"pay_type" json:"pay_type"`
	ReffNo        *string    `gorm:"reff_no" json:"reff_no"`
	MobileID      *int64     `gorm:"mobile_id" json:"mobile_id"`
	SubTotal      *float64   `gorm:"sub_total" json:"sub_total"`
	Disc          *float64   `gorm:"disc" json:"disc"`
	DiscValue     *float64   `gorm:"disc_value" json:"disc_value"`
	PromoValue    *float64   `gorm:"promo_value" json:"promo_value"`
	CashDiscValue *float64   `gorm:"cash_disc_value" json:"cash_disc_value"`
	TotDisc1      *float64   `gorm:"tot_disc1" json:"tot_disc1"`
	TotDisc2      *float64   `gorm:"tot_disc2" json:"tot_disc2"`
	Vat           *float64   `gorm:"vat" json:"vat"`
	VatValue      *float64   `gorm:"vat_value" json:"vat_value"`
	Total         *float64   `gorm:"total" json:"total"`
	DataStatus    *int64     `gorm:"data_status" json:"data_status"`
	DataSource    *int64     `gorm:"data_source" json:"data_source"`

	TrCode *string `gorm:"tr_code" json:"tr_code"`
	// IsClosed    bool       `gorm:"is_closed" json:"is_closed"`
	// ClosedBy    *int64     `gorm:"closed_by" json:"closed_by"`
	// ClosedAt    *time.Time `gorm:"closed_at" json:"closed_at"`
	DueDate     *time.Time `gorm:"due_date" json:"due_date"`
	Notes       *string    `gorm:"notes" json:"notes"`
	InvoiceNo   *string    `gorm:"invoice_no" json:"invoice_no"`
	InvoiceDate *time.Time `gorm:"invoice_date" json:"invoice_date"`

	CreatedBy *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel     bool            `gorm:"is_del" json:"is_del"`
	DeletedBy *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (Collection) TableName() string {
	return "acf.collection"
}

func (m *Collection) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type NoCollection struct {
	CustID             string     `gorm:"cust_id" json:"cust_id"`
	SalesmanId         *int64     `gorm:"salesman_id" json:"salesman_id"`
	NoCollectionDate   *time.Time `gorm:"no_order_date" json:"no_order_date"`
	OutletId           *int64     `gorm:"outlet_id" json:"outlet_id"`
	TakingCollectionId *int64     `gorm:"taking_order_id" json:"taking_order_id"`

	Reason *string `gorm:"reason" json:"reason"`

	CreatedBy *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel     bool            `gorm:"is_del" json:"is_del"`
	DeletedBy *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (NoCollection) TableName() string {
	return "sls.no_order"
}

func (m *NoCollection) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type CollectionList struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	CollectionNo    string     `gorm:"column:collection_no;primaryKey" json:"collection_no"`
	CollectionDate  *time.Time `gorm:"column:collection_date" json:"collection_date"`
	EmpID           *int64     `gorm:"column:emp_id" json:"emp_id"`
	SalesName       *string    `gorm:"sales_name" json:"sales_name"`
	OtGrpID         *int64     `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	Notes           *string    `gorm:"column:notes" json:"notes"`
	TotalAmount     float64    `gorm:"column:total_amount" json:"total_amount" default:"0"`
	RemainingAmount float64    `gorm:"column:remaining_amount" json:"remaining_amount" default:"0"`
	InvoiceNo       *string    `gorm:"column:invoice_no" json:"invoice_no"`
	InvoiceDateFrom *time.Time `gorm:"column:invoice_date_from" json:"invoice_date_from"`
	InvoiceDateTo   *time.Time `gorm:"column:invoice_date_to" json:"invoice_date_to"`
	DueDateFrom     *time.Time `gorm:"column:due_date_from" json:"due_date_from"`
	DueDateTo       *time.Time `gorm:"column:due_date_to" json:"due_date_to"`
	PaidAmount      float64    `gorm:"column:paid_amount" json:"paid_amount"`
	RoNo            string     `gorm:"column:ro_no" json:"ro_no"`
	OrderNo         string     `gorm:"column:order_no" json:"order_no"`

	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (CollectionList) TableName() string {
	return "acf.collection"
}

type NoCollectionList struct {
	CustID string `gorm:"cust_id" json:"cust_id"`
	// RoNo        string     `gorm:"ro_no" json:"ro_no"`
	NoCollectionId        *int32     `gorm:"no_order_id" json:"no_order_id"`
	SalesmanId            *int64     `gorm:"salesman_id" json:"salesman_id"`
	SalesName             *string    `gorm:"sales_name" json:"sales_name"`
	Reason                *string    `gorm:"reason" json:"reason"`
	NoCollectionDate      *time.Time `gorm:"no_order_date" json:"no_order_date"`
	OutletID              *int64     `gorm:"outlet_id" json:"outlet_id"`
	OutletCode            *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName            *string    `gorm:"column:outlet_name" json:"outlet_name"`
	TakingCollectionId    *int64     `gorm:"taking_order_id" json:"taking_order_id"`
	TakingCollectionName  *string    `gorm:"taking_order_name" json:"taking_order_name"`
	TakingCollectionImage *string    `gorm:"column:image_url" json:"image_url"`

	CreatedBy     *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string         `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool            `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

type MissedPaymentReason struct {
	CustId            string         `gorm:"column:cust_id" json:"cust_id" `
	MissedPaymentId   int64          `gorm:"column:missed_payment_reasons_id" json:"missed_payment_reasons_id"`
	MissedPaymentCode string         `gorm:"column:missed_payment_reasons_code" json:"missed_payment_reasons_code"`
	MissedPaymentName string         `gorm:"column:missed_payment_reasons_name" json:"missed_payment_reasons_name"`
	ImageUrl          string         `gorm:"column:image_url" json:"image_url"`
	IsActive          bool           `gorm:"column:is_active" json:"is_active"`
	CreatedBy         *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy         *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt         *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel             bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy         *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MissedPaymentReason) TableName() string {
	return "mst.missed_payment_reasons"
}

// CollectionSummaryDetail represents an invoice inside a collection
type CollectionSummaryDetail struct {
	InvoiceNumber    string               `json:"invoice_number"`
	InvoiceDate      string               `json:"invoice_date"`
	RONo             string               `json:"ro_no"`
	DueDate          string               `json:"due_date"`
	OutletID         int                  `json:"outlet_id"`
	OutletCode       string               `json:"outlet_code"`
	OutletName       string               `json:"outlet_name"`
	SalesmanID       int                  `json:"salesman_id"`
	SalesmanCode     string               `json:"salesman_code"`
	SalesmanName     string               `json:"salesman_name"`
	InvoiceAmount    float64              `json:"invoice_amount"`
	RemainintAmount  float64              `json:"remaining_amount"`
	TotalPayment     float64              `json:"total_payment"`
	Discount         float64              `json:"discount"`
	Materai          float64              `json:"materai"`
	PaymentBalance   float64              `json:"payment_balance"`
	RemainingPayment float64              `json:"remaining_payment"`
	IsCollection     bool                 `json:"is_collection"`
	Notes            string               `json:"notes"`
	Payments         []PaymentInvoiceList `json:"payments" gorm:"foreignKey:InvoiceNo;references:InvoiceNo"`
}

// CollectionListItem represents a collection with its associated invoices
type CollectionListItem struct {
	CollectionNo string                    `json:"collection_no"`
	TotalAmount  float64                   `json:"total_amount"`
	Details      []CollectionSummaryDetail `json:"details"`
}

type CollectionTotal struct {
	Total        int64   `json:"total"`
	TotalInvoice float64 `json:"total_invoice"`
}

type CollectionDetail struct {
	CustID          string     `json:"cust_id" gorm:"cust_id"`
	CollectionNo    string     `json:"collection_no" gorm:"collection_no"`
	InvoiceNo       string     `json:"invoice_no" gorm:"invoice_no"`
	SalesmanID      int64      `json:"salesman_id" gorm:"salesman_id"`
	InvoiceAmount   float64    `json:"invoice_amount" gorm:"invoice_amount"`
	RemainingAmount float64    `json:"remaining_amount" gorm:"remaining_amount"`
	PaidAmount      float64    `json:"paid_amount" gorm:"paid_amount"`
	CreatedBy       *int64     `json:"created_by" gorm:"created_by"`
	CreatedAt       *time.Time `json:"created_at" gorm:"created_at"`
	Source          string     `json:"source" gorm:"source"`
}

func (CollectionDetail) TableName() string {
	return "acf.collection_det"
}

type CollectionModel struct {
	CustID          string     `gorm:"column:cust_id" json:"cust_id"`
	CollectionNo    string     `gorm:"column:collection_no;primaryKey" json:"collection_no"`
	CollectionDate  *time.Time `gorm:"column:collection_date" json:"collection_date"`
	EmpID           *int64     `gorm:"column:emp_id" json:"emp_id"`
	OtGrpID         *int64     `gorm:"column:ot_grp_id" json:"ot_grp_id"`
	Notes           *string    `gorm:"column:notes" json:"notes"`
	TotalAmount     float64    `gorm:"column:total_amount" json:"total_amount" default:"0"`
	RemainingAmount float64    `gorm:"column:remaining_amount" json:"remaining_amount" default:"0"`
	InvoiceDateFrom *time.Time `gorm:"column:invoice_date_from" json:"invoice_date_from"`
	InvoiceDateTo   *time.Time `gorm:"column:invoice_date_to" json:"invoice_date_to"`
	DueDateFrom     *time.Time `gorm:"column:due_date_from" json:"due_date_from"`
	DueDateTo       *time.Time `gorm:"column:due_date_to" json:"due_date_to"`

	CreatedBy *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt time.Time       `gorm:"updated_at" json:"updated_at"`
	IsDel     bool            `gorm:"is_del" json:"is_del"`
	DeletedBy *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
	PrintedBy *int64          `gorm:"printed_by" json:"printed_by"`
	PrintedAt *time.Time      `gorm:"printed_at" json:"printed_at"`
	IsPrinted bool            `gorm:"is_printed" json:"is_printed"`
	Source    string          `gorm:"source" json:"source"`
}

func (CollectionModel) TableName() string {
	return "acf.collection"
}

type InvoiceList struct {
	CollectionNo string `json:"collection_no"`
	InvoiceNo    string `json:"invoice_no"`
}
