package model

import (
	"time"

	"gorm.io/gorm"
)

type ChequeGiro struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	ChequeGiroNo     int        `gorm:"column:cheque_giro_no;default:nextval('acf.cheque_giro_no_seq'::regclass);not null" json:"cheque_giro_no"`
	DocNoCheque      string     `gorm:"column:doc_no_cheque" json:"doc_no_cheque"`
	OwnerID          int        `gorm:"column:owner_id" json:"owner_id"`
	SalesmanID       *int       `gorm:"column:salesman_id" json:"salesman_id"`
	OutletID         *int       `gorm:"column:outlet_id" json:"outlet_id"`
	SupplierID       *int       `gorm:"column:sup_id" json:"sup_id"`
	BankID           int        `gorm:"column:bank_id" json:"bank_id"`
	BankIDCollecting int        `gorm:"column:bank_id_collecting" json:"bank_id_collecting"`
	AccountNo        string     `gorm:"column:account_no" json:"account_no"`
	OutletBankID     *int       `gorm:"column:outlet_bank_id" json:"outlet_bank_id"`
	DocDateCheque    *time.Time `gorm:"column:doc_date_cheque" json:"doc_date_cheque"`
	DueDate          *time.Time `gorm:"column:due_date" json:"due_date"`
	Amount           float64    `gorm:"column:amount" json:"amount"`
	StatusCheque     int        `gorm:"column:status_cheque" json:"status_cheque"`
	ClearingDate     *time.Time `gorm:"column:clearing_date" json:"clearing_date"`
	CreatedBy        int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy        int        `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	IsDel            bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy        *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *ChequeGiro) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (ChequeGiro) TableName() string {
	return "acf.cheque_giro"
}

type ChequeGiroList struct {
	CustID           string     `gorm:"column:cust_id" json:"cust_id"`
	ChequeGiroNo     int        `gorm:"column:cheque_giro_no" json:"cheque_giro_no"`
	DocNoCheque      string     `gorm:"column:doc_no_cheque" json:"doc_no_cheque"`
	OwnerID          int        `gorm:"column:owner_id" json:"owner_id"`
	OwnerName        string     `gorm:"column:owner_name" json:"owner_name"`
	SalesmanID       *int       `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName     *string    `gorm:"column:sales_name" json:"sales_name"`
	SupplierID       *int       `gorm:"column:sup_id" json:"sup_id"`
	SupplierName     *string    `gorm:"column:sup_name" json:"sup_name"`
	OutletID         *int       `gorm:"column:outlet_id" json:"outlet_id"`
	OutletName       *string    `gorm:"column:outlet_name" json:"outlet_name"`
	BankID           int        `gorm:"column:bank_id" json:"bank_id"`
	BankName         string     `gorm:"column:bank_name" json:"bank_name"`
	BankIDCollecting int        `gorm:"column:bank_id_collecting" json:"bank_id_collecting"`
	AccountNo        string     `gorm:"column:account_no" json:"account_no"`
	OutletBankID     *int       `gorm:"column:outlet_bank_id" json:"outlet_bank_id"`
	DocDateCheque    *time.Time `gorm:"column:doc_date_cheque" json:"doc_date_cheque"`
	DueDate          *time.Time `gorm:"column:due_date" json:"due_date"`
	Amount           float64    `gorm:"column:amount" json:"amount"`
	UsedAmount       *float64   `gorm:"column:used_amount" json:"used_amount"`
	UsedAmountOutlet float64    `gorm:"column:used_amount_outlet" json:"used_amount_outlet"`
	StatusCheque     int        `gorm:"column:status_cheque" json:"status_cheque"`
	ClearingDate     *time.Time `gorm:"column:clearing_date" json:"clearing_date"`
	CreatedBy        int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy        int        `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	IsDel            bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy        *int       `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ChequeGiroList) TableName() string {
	return "acf.cheque_giro"
}

type BankLookup struct {
	BankId   int    `db:"bank_id" json:"bank_id"`
	BankCode string `db:"bank_code" json:"bank_code"`
	BankName string `db:"bank_name" json:"bank_name"`
}

func (BankLookup) TableName() string {
	return "acf.cheque_giro"
}

type BankAccountLookup struct {
	BankId    int    `db:"bank_id" json:"bank_id"`
	AccountNo string `db:"account_no" json:"account_no"`
}

func (BankAccountLookup) TableName() string {
	return "acf.cheque_giro"
}
