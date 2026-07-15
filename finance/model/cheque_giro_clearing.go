package model

import (
	"time"
)

type ChequeGiroClearingList struct {
	CustID             string     `gorm:"column:cust_id" json:"cust_id"`
	ChequeGiroNo       int        `gorm:"column:cheque_giro_no" json:"cheque_giro_no"`
	DocNoCheque        string     `gorm:"column:doc_no_cheque" json:"doc_no_cheque"`
	OwnerID            int        `gorm:"column:owner_id" json:"owner_id"`
	OwnerName          string     `gorm:"column:owner_name" json:"owner_name"`
	OutletID           *int       `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode         *string    `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName         *string    `gorm:"column:outlet_name" json:"outlet_name"`
	BankID             int        `gorm:"column:bank_id" json:"bank_id"`
	BankName           string     `gorm:"column:bank_name" json:"bank_name"`
	BankIDCollecting   int        `gorm:"column:bank_id_collecting" json:"bank_id_collecting"`
	BankNameCollecting string     `gorm:"column:bank_id_collecting_name" json:"bank_id_collecting_name"`
	AccountNo          string     `gorm:"column:account_no" json:"account_no"`
	DocDateCheque      *time.Time `gorm:"column:doc_date_cheque" json:"doc_date_cheque"`
	DueDate            *time.Time `gorm:"column:due_date" json:"due_date"`
	Amount             float64    `gorm:"column:amount" json:"amount"`
	UsedAmount         *float64   `gorm:"column:used_amount" json:"used_amount"`
	StatusCheque       int        `gorm:"column:status_cheque" json:"status_cheque"`
	ClearingDate       *time.Time `gorm:"column:clearing_date" json:"clearing_date"`
	Reason             *string    `gorm:"column:reason" json:"reason"`
	CashNo             *string    `gorm:"column:cash_no" json:"cash_no"`
	CashAmount         *float64   `gorm:"column:cash_amount" json:"cash_amount"`
	TransferNo         *string    `gorm:"column:transfer_no" json:"transfer_no"`
	TransferAmount     *float64   `gorm:"column:transfer_amount" json:"transfer_amount"`
	ChequeNo           *string    `gorm:"column:cheque_no" json:"cheque_no"`
	ChequeAmount       *float64   `gorm:"column:cheque_amount" json:"cheque_amount"`
	CreatedBy          int        `gorm:"column:created_by" json:"created_by"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedBy          int        `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (ChequeGiroClearingList) TableName() string {
	return "acf.cheque_giro"
}
