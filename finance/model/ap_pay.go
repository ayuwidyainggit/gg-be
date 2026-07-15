package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ApPay struct {
	CustID      string         `gorm:"column:cust_id" json:"cust_id"`
	ApPayNo     string         `gorm:"column:ap_pay_no;primaryKey" json:"ap_pay_no"`
	ApPayDate   *time.Time     `gorm:"column:ap_pay_date" json:"ap_pay_date"`
	TrCode      *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID       *int64         `gorm:"column:sup_id" json:"sup_id"`
	CashAmt     *float64       `gorm:"column:cash_amt" json:"cash_amt"`
	CndnAmt     *float64       `gorm:"column:cndn_amt" json:"cndn_amt"`
	ReturnAmt   *float64       `gorm:"column:return_amt" json:"return_amt"`
	ChequeAmt   *float64       `gorm:"column:cheque_amt" json:"cheque_amt"`
	TransferAmt *float64       `gorm:"column:transfer_amt" json:"transfer_amt"`
	TotalAmt    *float64       `gorm:"column:total_amt" json:"total_amt"`
	DataStatus  *int64         `gorm:"column:data_status" json:"data_status"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	IsDel       bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy   *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	IsPosted    *bool          `gorm:"column:is_posted" json:"is_posted"`
	PostedAt    *time.Time     `gorm:"column:posted_at" json:"posted_at"`
}

func (ApPay) TableName() string {
	return "acf.ap_pay"
}
func (m *ApPay) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()
	if m.IsPosted != nil {
		if *m.IsPosted {
			m.PostedAt = &now
		}
	}
	intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	m.ApPayNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

type ApPayList struct {
	CustID        string         `gorm:"column:cust_id" json:"cust_id"`
	ApPayNo       string         `gorm:"column:ap_pay_no;primaryKey" json:"ap_pay_no"`
	ApPayDate     *time.Time     `gorm:"column:ap_pay_date" json:"ap_pay_date"`
	TrCode        *string        `gorm:"column:tr_code" json:"tr_code"`
	SupID         *int64         `gorm:"column:sup_id" json:"sup_id"`
	SupCode       *string        `gorm:"column:sup_code" json:"sup_code"`
	SupName       *string        `gorm:"column:sup_name" json:"sup_name"`
	CashAmt       *float64       `gorm:"column:cash_amt" json:"cash_amt"`
	CndnAmt       *float64       `gorm:"column:cndn_amt" json:"cndn_amt"`
	ReturnAmt     *float64       `gorm:"column:return_amt" json:"return_amt"`
	ChequeAmt     *float64       `gorm:"column:cheque_amt" json:"cheque_amt"`
	TransferAmt   *float64       `gorm:"column:transfer_amt" json:"transfer_amt"`
	TotalAmt      *float64       `gorm:"column:total_amt" json:"total_amt"`
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

func (ApPayList) TableName() string {
	return "acf.ap_pay"
}

type ApPayJoinDet struct {
	CustID  string `gorm:"column:cust_id" json:"cust_id"`
	ApPayNo string `gorm:"column:ap_pay_no;primaryKey" json:"ap_pay_no"`
}

func (ApPayJoinDet) TableName() string {
	return "acf.ap_pay"
}
