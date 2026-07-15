package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ArPay struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ArPayNo       string         `gorm:"column:ar_pay_no;primaryKey" json:"ar_pay_no"`
	ArPayDate     *time.Time     `gorm:"column:ar_pay_date" json:"ar_pay_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	ArNo          *string        `gorm:"column:ar_no" json:"ar_no"`
	SalesmanID    *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	CashAmt       *float64       `gorm:"column:cash_amt" json:"cash_amt"`
	ChequeAmt     *float64       `gorm:"column:cheque_amt" json:"cheque_amt"`
	TransferAmt   *float64       `gorm:"column:transfer_amt" json:"transfer_amt"`
	ReturnAmt     *float64       `gorm:"column:return_amt" json:"return_amt"`
	CndnAmt       *float64       `gorm:"column:cndn_amt" json:"cndn_amt"`
	DiscAmt       *float64       `gorm:"column:disc_amt" json:"disc_amt"`
	DutyStampAmt  *float64       `gorm:"column:duty_stamp_amt" json:"duty_stamp_amt"`
	TotalAmt      *float64       `gorm:"column:total_amt" json:"total_amt"`
	TotalDiff     *float64       `gorm:"column:total_diff" json:"total_diff"`
	TotalAmtRound *float64       `gorm:"column:total_amt_round" json:"total_amt_round"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted      *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt      *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (ArPay) TableName() string {
	return "acf.ar_pay"
}
func (m *ArPay) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	m.ArPayNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

type ArPayList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ArPayNo       string         `gorm:"column:ar_pay_no;primaryKey" json:"ar_pay_no"`
	ArPayDate     *time.Time     `gorm:"column:ar_pay_date" json:"ar_pay_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	ArNo          *string        `gorm:"column:ar_no" json:"ar_no"`
	SalesmanID    *int64         `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanCode  *string        `gorm:"column:salesman_code" json:"salesman_code"`
	SalesmanName  *string        `gorm:"column:salesman_name" json:"salesman_name"`
	CashAmt       *float64       `gorm:"column:cash_amt" json:"cash_amt"`
	ChequeAmt     *float64       `gorm:"column:cheque_amt" json:"cheque_amt"`
	TransferAmt   *float64       `gorm:"column:transfer_amt" json:"transfer_amt"`
	ReturnAmt     *float64       `gorm:"column:return_amt" json:"return_amt"`
	CndnAmt       *float64       `gorm:"column:cndn_amt" json:"cndn_amt"`
	DiscAmt       *float64       `gorm:"column:disc_amt" json:"disc_amt"`
	DutyStampAmt  *float64       `gorm:"column:duty_stamp_amt" json:"duty_stamp_amt"`
	TotalAmt      *float64       `gorm:"column:total_amt" json:"total_amt"`
	TotalDiff     *float64       `gorm:"column:total_diff" json:"total_diff"`
	TotalAmtRound *float64       `gorm:"column:total_amt_round" json:"total_amt_round"`
	DataStatus    *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted      *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt      *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (ArPayList) TableName() string {
	return "acf.ar_pay"
}
