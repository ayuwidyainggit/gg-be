package model

import (
	"time"

	"gorm.io/gorm"
)

type BankTransfer struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	BankTransferNo   int        `gorm:"column:bank_transfer_no;default:nextval('acf.bank_transfer_no_seq'::regclass);not null" json:"bank_transfer_no"`
	DocNoBank        string     `gorm:"column:doc_no_bank" json:"doc_no_bank"`
	OwnerID          int        `gorm:"column:owner_id" json:"owner_id"`
	SalesmanID       *int       `gorm:"column:salesman_id" json:"salesman_id"`
	OutletID         *int       `gorm:"column:outlet_id" json:"outlet_id"`
	SupplierID       *int       `gorm:"column:sup_id" json:"sup_id"`
	BankID           int        `gorm:"column:bank_id" json:"bank_id"`
	BankIDCollecting int        `gorm:"column:bank_id_collecting" json:"bank_id_collecting"`
	AccountNo        *string    `gorm:"column:account_no" json:"account_no"`
	AccountName      string     `gorm:"column:account_name" json:"account_name"`
	OutletBankID     *int       `gorm:"column:outlet_bank_id" json:"outlet_bank_id"`
	TransferDate     *time.Time `gorm:"column:transfer_date" json:"transfer_date"`
	Amount           float64    `gorm:"column:amount" json:"amount"`
	// UsedAmount       *float64   `gorm:"column:used_amount" json:"used_amount"`
	StatusBank int        `gorm:"column:status_bank_transfer" json:"status_bank_transfer"`
	CreatedBy  int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  int        `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *BankTransfer) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (BankTransfer) TableName() string {
	return "acf.bank_transfer"
}

type BankTransferList struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	BankTransferNo   int        `gorm:"column:bank_transfer_no" json:"bank_transfer_no"`
	DocNoBank        string     `gorm:"column:doc_no_bank" json:"doc_no_bank"`
	OwnerID          int        `gorm:"column:owner_id" json:"owner_id"`
	OwnerName        string     `gorm:"column:owner_name" json:"owner_name"`
	SalesmanID       *int       `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName     *string    `gorm:"column:sales_name" json:"sales_name"`
	SalesmanCode     *string    `gorm:"column:salesman_code" json:"salesman_code"`
	SupplierID       *int       `gorm:"column:sup_id" json:"sup_id"`
	SupplierName     *string    `gorm:"column:sup_name" json:"sup_name"`
	SupplierCode     *string    `gorm:"column:sup_code" json:"sup_code"`
	OutletID         *int       `gorm:"column:outlet_id" json:"outlet_id"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	OutletCode       *string    `gorm:"column:outlet_code" json:"outlet_code"`
	BankID           int        `gorm:"column:bank_id" json:"bank_id"`
	BankName         string     `gorm:"column:bank_name" json:"bank_name"`
	BankIDCollecting int        `gorm:"column:bank_id_collecting" json:"bank_id_collecting"`
	AccountNo        *string    `gorm:"column:account_no" json:"account_no"`
	AccountName      string     `gorm:"column:account_name" json:"account_name"`
	OutletBankID     *int       `gorm:"column:outlet_bank_id" json:"outlet_bank_id"`
	TransferDate     *time.Time `gorm:"column:transfer_date" json:"transfer_date"`
	Amount           float64    `gorm:"column:amount" json:"amount"`
	UsedAmount       float64    `gorm:"column:used_amount" json:"used_amount"`
	UsedAmountOutlet float64    `gorm:"column:used_amount_outlet" json:"used_amount_outlet"`
	StatusBank       int        `gorm:"column:status_bank_transfer" json:"status_bank_transfer"`
	CreatedBy        int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy        int        `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	IsDel            bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy        *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (BankTransferList) TableName() string {
	return "acf.bank_transfer"
}

// BankTransferDepositDataRow: one row from used_amount query (deposit + deposit_detail + deposit_payment, pay_type=3)
type BankTransferDepositDataRow struct {
	DepositNo   string     `gorm:"column:deposit_no" json:"deposit_no"`
	DepositDate *time.Time `gorm:"column:deposit_date" json:"deposit_date"`
	InvoiceNo   string     `gorm:"column:invoice_no" json:"invoice_no"`
	UsedAmount  float64    `gorm:"column:used_amount" json:"used_amount"`
}

type BankTransferFile struct {
	BankTransferFileID int       `gorm:"column:bank_transfer_file_id;primaryKey;autoIncrement" json:"bank_transfer_file_id"`
	CustID             string    `gorm:"column:cust_id" json:"cust_id"`
	BankTransferNo     string    `gorm:"column:bank_transfer_no" json:"bank_transfer_no"`
	FileName           string    `gorm:"column:file_name" json:"file_name"`
	FileURL            string    `gorm:"column:file_url" json:"file_url"`
	FileKey            string    `gorm:"column:file_key" json:"file_key"`
	MediaCategory      string    `gorm:"column:media_category" json:"media_category"`
	FileSize           int64     `gorm:"column:file_size" json:"file_size"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
}

func (BankTransferFile) TableName() string {
	return "acf.bank_transfer_files"
}

type BankLookupBankTransfer struct {
	BankId   int    `db:"bank_id" json:"bank_id"`
	BankCode string `db:"bank_code" json:"bank_code"`
	BankName string `db:"bank_name" json:"bank_name"`
}

func (BankLookupBankTransfer) TableName() string {
	return "acf.bank_transfer"
}

type BankAccountLookupBankTransfer struct {
	BankId    int    `db:"bank_id" json:"bank_id"`
	AccountNo string `db:"account_no" json:"account_no"`
}

func (BankAccountLookupBankTransfer) TableName() string {
	return "acf.bank_transfer"
}
