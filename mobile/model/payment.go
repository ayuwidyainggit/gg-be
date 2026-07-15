package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PaymentType represents acf.payment_type table
type PaymentType struct {
	PaymentTypeID   int        `gorm:"column:payment_type_id;primaryKey;autoIncrement" json:"payment_type_id"`
	PaymentTypeCode string     `gorm:"column:payment_type_code;type:varchar(20);not null" json:"payment_type_code"`
	PaymentTypeName string     `gorm:"column:payment_type_name;type:varchar(50);not null" json:"payment_type_name"`
	CreatedBy       int        `gorm:"column:created_by;type:int4;not null" json:"created_by"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy       *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	DeletedBy       *int64     `gorm:"column:deleted_by;type:int8" json:"deleted_by"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	IsDel           bool       `gorm:"column:is_del;default:false" json:"is_del"`
}

func (PaymentType) TableName() string {
	return "acf.payment_type"
}

// PaymentTrx represents acf.payment_trx table
type PaymentTrx struct {
	CustID           string          `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	InvoiceNo        string          `gorm:"column:invoice_no;type:varchar(255);not null" json:"invoice_no"`
	CollectionNo     string          `gorm:"column:collection_no;type:varchar(255);not null" json:"collection_no"`
	PaymentTrxID     int64           `gorm:"column:payment_trx_id;primaryKey;autoIncrement" json:"payment_trx_id"`
	OutletID         int64           `gorm:"column:outlet_id;type:int4;not null" json:"outlet_id"`
	EmpID            int64           `gorm:"column:emp_id;type:int4;not null" json:"emp_id"`
	PONumber         string          `gorm:"column:po_number;type:varchar(50);not null" json:"po_number"`
	DocumentNo       string          `gorm:"column:document_no;type:varchar(50);not null" json:"document_no"`
	TrxSource        string          `gorm:"column:trx_source;type:varchar(1);not null" json:"trx_source"`
	TrxRefNo         *int64          `gorm:"column:trx_ref_no;type:int8" json:"trx_ref_no"`
	TotalTransaction float64         `gorm:"column:total_transaction;type:numeric(20,4);not null;default:0" json:"total_transaction"`
	PaymentAmount    float64         `gorm:"column:payment_amount;type:numeric(20,4);not null;default:0" json:"payment_amount"`
	RemainingAmount  float64         `gorm:"column:remaining_amount;type:numeric(20,4);not null;default:0" json:"remaining_amount"`
	Notes            *string         `gorm:"column:notes;type:varchar(100)" json:"notes"`
	Files            any             `gorm:"column:files;type:jsonb" json:"files"`
	Date             time.Time       `gorm:"column:date;type:date;not null" json:"date"`
	CreatedBy        int64           `gorm:"column:created_by;type:int4;not null" json:"created_by"`
	CreatedAt        time.Time       `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy        *int64          `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt        *time.Time      `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	DeletedBy        *int            `gorm:"column:deleted_by;type:int4" json:"deleted_by"`
	DeletedAt        *time.Time      `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	IsDel            bool            `gorm:"column:is_del;default:false" json:"is_del"`
	Details          []PaymentTrxDet `gorm:"foreignKey:CustID,PaymentTrxID;references:CustID,PaymentTrxID" json:"details,omitempty"`
}

func (PaymentTrx) TableName() string {
	return "acf.payment_trx"
}

// PaymentTrxDet represents acf.payment_trx_detail table
type PaymentTrxDet struct {
	CustID          string     `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	PaymentTrxDetID int64      `gorm:"column:payment_trx_det_id;primaryKey;autoIncrement" json:"payment_trx_det_id"`
	PaymentTrxID    int64      `gorm:"column:payment_trx_id;type:int8;not null" json:"payment_trx_id"`
	CNDNNo          *string    `gorm:"column:cndn_no;type:varchar(30)" json:"cndn_no"`
	PayType         int16      `gorm:"column:pay_type;type:int2;not null" json:"pay_type"`
	BankTransferNo  *int       `gorm:"column:bank_transfer_no;type:int4" json:"bank_transfer_no"`
	ChequeGiroNo    *int       `gorm:"column:cheque_giro_no;type:int4" json:"cheque_giro_no"`
	Amount          float64    `gorm:"column:amount;type:numeric(20,4);not null;default:0" json:"amount"`
	CreatedBy       int64      `gorm:"column:created_by;type:int4;not null" json:"created_by"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy       *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	DeletedBy       *int       `gorm:"column:deleted_by;type:int4" json:"deleted_by"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	IsDel           bool       `gorm:"column:is_del;default:false" json:"is_del"`
}

func (PaymentTrxDet) TableName() string {
	return "acf.payment_trx_detail"
}

// BeforeCreate hook for PaymentTrx
func (p *PaymentTrx) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	// Generate document_no if empty
	if p.DocumentNo == "" {
		now := time.Now()
		prefix := fmt.Sprintf("TRX%s", now.Format("060102")) // YYMMDD

		var lastDocNo string
		err := tx.Table("acf.payment_trx").
			Select("document_no").
			Where("cust_id = ? AND document_no LIKE ?", p.CustID, prefix+"%").
			Order("document_no DESC").
			Limit(1).
			Scan(&lastDocNo).Error

		if err != nil {
			return err
		}

		runningNumber := 1
		if lastDocNo != "" && len(lastDocNo) >= 13 {
			// Extract last 4 digits
			fmt.Sscanf(lastDocNo[len(lastDocNo)-4:], "%d", &runningNumber)
			runningNumber++
		}

		p.DocumentNo = fmt.Sprintf("%s%04d", prefix, runningNumber)
	}

	return nil
}

// BeforeUpdate hook for PaymentTrx
func (p *PaymentTrx) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	p.UpdatedAt = &now
	return nil
}

// BankTransfer represents acf.bank_transfer table
type BankTransfer struct {
	CustID             string     `gorm:"column:cust_id;type:varchar(10);primaryKey" json:"cust_id"`
	BankTransferNo     int        `gorm:"column:bank_transfer_no;primaryKey;autoIncrement:false;default:nextval('acf.cheque_giro_no_seq')" json:"bank_transfer_no"`
	DocNoBank          string     `gorm:"column:doc_no_bank;type:varchar(25);not null" json:"doc_no_bank"`
	OwnerID            int        `gorm:"column:owner_id;type:int4;not null" json:"owner_id"`
	SalesmanID         *int64     `gorm:"column:salesman_id;type:int8" json:"salesman_id"`
	OutletID           *int64     `gorm:"column:outlet_id;type:int8" json:"outlet_id"`
	BankID             *int64     `gorm:"column:bank_id;type:int8" json:"bank_id"`
	BankIDCollecting   *int64     `gorm:"column:bank_id_collecting;type:int8" json:"bank_id_collecting"`
	AccountNo          *string    `gorm:"column:account_no;type:varchar(50)" json:"account_no"`
	TransferDate       *time.Time `gorm:"column:transfer_date;type:date" json:"transfer_date"`
	Amount             float64    `gorm:"column:amount;type:numeric(20,4);default:0" json:"amount"`
	StatusBankTransfer int        `gorm:"column:status_bank_transfer;type:int4;default:2" json:"status_bank_transfer"`
	CreatedBy          int64      `gorm:"column:created_by;type:int8" json:"created_by"`
	CreatedAt          time.Time  `gorm:"column:created_at;type:timestamptz(6);default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy          *int64     `gorm:"column:updated_by;type:int8" json:"updated_by"`
	UpdatedAt          *time.Time `gorm:"column:updated_at;type:timestamptz(6)" json:"updated_at"`
	IsDel              bool       `gorm:"column:is_del;default:false" json:"is_del"`
	DeletedBy          *int64     `gorm:"column:deleted_by;type:int8" json:"deleted_by"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;type:timestamptz(6)" json:"deleted_at"`
	OutletBankID       *int64     `gorm:"column:outlet_bank_id;type:int8" json:"outlet_bank_id"`
	SupID              *int64     `gorm:"column:sup_id;type:int8" json:"sup_id"`
	RemainingAmount    float64    `gorm:"column:remaining_amount;type:numeric(20,4);not null;default:0" json:"remaining_amount"`
	PaidAmount         float64    `gorm:"column:paid_amount;type:numeric(20,4);not null;default:0" json:"paid_amount"`
	AccountName        string     `gorm:"column:account_name;type:varchar(255);not null;default:''" json:"account_name"`
}

func (BankTransfer) TableName() string {
	return "acf.bank_transfer"
}

// BeforeCreate hook for BankTransfer
func (b *BankTransfer) BeforeCreate(tx *gorm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}

	// Set Default Values
	if b.OwnerID == 0 {
		b.OwnerID = 1
	}
	if b.StatusBankTransfer == 0 {
		b.StatusBankTransfer = 2
	}
	if b.PaidAmount == 0 {
		b.PaidAmount = 0
	}
	if b.RemainingAmount == 0 {
		b.RemainingAmount = b.Amount
	}

	// Generate doc_no_bank if empty
	if b.DocNoBank == "" {
		now := time.Now()
		prefix := fmt.Sprintf("TF%s", now.Format("060102")) // YYMMDD

		var lastDocNo string
		err := tx.Table("acf.bank_transfer").
			Select("doc_no_bank").
			Where("cust_id = ? AND doc_no_bank LIKE ?", b.CustID, prefix+"%").
			Order("doc_no_bank DESC").
			Limit(1).
			Scan(&lastDocNo).Error

		if err != nil {
			return err
		}

		runningNumber := 1
		if lastDocNo != "" && len(lastDocNo) >= 12 {
			// TFYYMMDDNNNN is 12 characters
			fmt.Sscanf(lastDocNo[len(lastDocNo)-4:], "%d", &runningNumber)
			runningNumber++
		}

		b.DocNoBank = fmt.Sprintf("%s%04d", prefix, runningNumber)
	}

	return nil
}

// PaymentTrxSalesData represents query result for sales data finding
type PaymentTrxSalesData struct {
	PaymentAmount   float64 `gorm:"column:payment_amount" json:"payment_amount"`
	PaymentTypeCode string  `gorm:"column:payment_type_code" json:"payment_type_code"`
}

// InvoiceListItem represents a single invoice with its payment amount
type InvoiceListItem struct {
	InvoiceNo        string               `json:"invoice_number"`
	InvoiceDate      string               `json:"invoice_date"`
	RONo             string               `json:"ro_no"`
	OrderNo          string               `json:"order_no"`
	DueDate          string               `json:"due_date"`
	OutletID         int                  `json:"outlet_id"`
	OutletCode       string               `json:"outlet_code"`
	OutletName       string               `json:"outlet_name"`
	SalesmanID       int                  `json:"salesman_id"`
	SalesmanCode     string               `json:"salesman_code"`
	SalesmanName     string               `json:"salesman_name"`
	InvoiceAmount    float64              `json:"invoice_amount"`
	RemainingAmount  float64              `json:"remaining_amount"`
	PaidAmount       float64              `json:"paid_amount"`
	TotalPayment     float64              `json:"total_payment"`
	Discount         float64              `json:"discount"`
	Materai          float64              `json:"materai"`
	PaymentBalance   float64              `json:"payment_balance"`
	RemainingPayment float64              `json:"remaining_payment"`
	IsCollection     bool                 `json:"is_collection"`
	Notes            string               `json:"notes"`
	Payments         []PaymentInvoiceList `json:"payments" gorm:"foreignKey:InvoiceNo;references:InvoiceNo"`
}

type PaymentInvoiceList struct {
	InvoiceNo     string  `json:"invoice_no"`
	PayType       int     `json:"pay_type"`
	DocumentNo    string  `json:"document_no"`
	PaymentAmount float64 `json:"payment_amount"`
}

type InvoicePayment struct {
	PaymentInvoiceList
	InvoiceAmount    float64 `json:"invoice_amount"`
	RemainingAmount  float64 `json:"remaining_amount"`
	TotalPayment     float64 `json:"total_payment"`
	RemainingPayment float64 `json:"remaining_payment"`
}
