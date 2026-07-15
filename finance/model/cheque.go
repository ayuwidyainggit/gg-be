package model

import (
	"time"

	"gorm.io/gorm"
)

type Cheque struct {
	CustID       string         `gorm:"column:cust_id" json:"cust_id"`
	ChqTrNo      string         `gorm:"column:chq_tr_no" json:"chq_tr_no"`
	TrCode       *string        `gorm:"column:tr_code" json:"tr_code"`
	ChqNo        *int64         `gorm:"column:chq_no;primaryKey" json:"chq_no"`
	ChqTrType    *int           `gorm:"column:chq_tr_type" json:"chq_tr_type"`
	ChqDate      *time.Time     `gorm:"column:chq_date" json:"chq_date"`
	ChqDueDate   *time.Time     `gorm:"column:chq_due_date" json:"chq_due_date"`
	BankId       *int64         `gorm:"column:bank_id" json:"bank_id"`
	AccountNo    *string        `gorm:"column:account_no" json:"account_no"`
	ChqMat       *float64       `gorm:"column:chq_amt" json:"chq_amt"`
	ChqUsedAmt   *float64       `gorm:"column:chq_used_amt" json:"chq_used_amt"`
	SalesmanId   *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	OutletId     *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	Notes        *string        `gorm:"column:notes" json:"notes"`
	ClearingDate *time.Time     `gorm:"column:clearing_date" json:"clearing_date"`
	ChqStatus    *int64         `gorm:"column:chq_status" json:"chq_status"`
	StatusDate   *time.Time     `gorm:"column:status_date" json:"status_date"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel        bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy    *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted     *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt     *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (m *Cheque) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	// intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	// m.ChqTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (Cheque) TableName() string {
	return "acf.cheque"
}

type ChequeList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ChqTrNo       string         `gorm:"column:chq_tr_no" json:"chq_tr_no"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	ChqNo         *int64         `gorm:"column:chq_no;primaryKey" json:"chq_no"`
	ChqTrType     *int           `gorm:"column:chq_tr_type" json:"chq_tr_type"`
	ChqDate       *time.Time     `gorm:"column:chq_date" json:"chq_date"`
	ChqDueDate    *time.Time     `gorm:"column:chq_due_date" json:"chq_due_date"`
	BankId        *int64         `gorm:"column:bank_id" json:"bank_id"`
	BankCode      *string        `gorm:"column:bank_code" json:"bank_code"`
	Bankname      *string        `gorm:"column:bank_name" json:"bank_name"`
	AccountNo     *string        `gorm:"column:account_no" json:"account_no"`
	ChqMat        *float64       `gorm:"column:chq_amt" json:"chq_amt"`
	ChqUsedAmt    *float64       `gorm:"column:chq_used_amt" json:"chq_used_amt"`
	SalesmanId    *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  *string        `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string        `gorm:"column:salesman_name" json:"salesman_name"`
	OutletId      *int64         `gorm:"column:outlet_id" json:"outlet_id"`
	OutletCode    *string        `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    *string        `gorm:"column:outlet_name" json:"outlet_name"`
	Notes         *string        `gorm:"column:notes" json:"notes"`
	ClearingDate  *time.Time     `gorm:"column:clearing_date" json:"clearing_date"`
	ChqStatus     *int64         `gorm:"column:chq_status" json:"chq_status"`
	StatusDate    *time.Time     `gorm:"column:status_date" json:"status_date"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted      *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt      *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (ChequeList) TableName() string {
	return "acf.cheque"
}
