package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CashTr struct {
	CustID      string         `gorm:"column:cust_id" json:"cust_id"`
	CashTrNo    string         `gorm:"column:cash_tr_no;primaryKey" json:"cash_tr_no"`
	TrCode      *string        `gorm:"column:tr_code" json:"tr_code"`
	CashTrDate  *time.Time     `gorm:"column:cash_tr_date" json:"cash_tr_date"`
	CoaIdTo     *int64         `gorm:"column:coa_id_to" json:"coa_id_to"`
	Notes       *string        `gorm:"column:notes" json:"notes"`
	AccountNo   *string        `gorm:"column:account_no" json:"account_no"`
	AccountName *string        `gorm:"column:account_name" json:"account_name"`
	Amount      *float64       `gorm:"column:amount" json:"amount"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel       bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy   *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CashTr) TableName() string {
	return "acf.cash_tr"
}
func (m *CashTr) BeforeCreate(trx *gorm.DB) (err error) {
	intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	m.CashTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type CashTrList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	CashTrNo      string         `gorm:"column:cash_tr_no;primaryKey" json:"cash_tr_no"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	CashTrDate    *time.Time     `gorm:"column:cash_tr_date" json:"cash_tr_date"`
	CoaIdTo       *int64         `gorm:"column:coa_id_to" json:"coa_id_to"`
	CoaCodeTo     *string        `gorm:"column:coa_code_to" json:"coa_code_to"`
	CoaNameTo     *string        `gorm:"column:coa_name_to" json:"coa_name_to"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	AccountNo     *string        `gorm:"column:account_no" json:"account_no"`
	AccountName   *string        `gorm:"column:account_name" json:"account_name"`
	Amount        *float64       `gorm:"column:amount" json:"amount"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CashTrList) TableName() string {
	return "acf.cash_tr"
}
