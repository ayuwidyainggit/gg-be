package model

import (
	"time"

	"gorm.io/gorm"
)

type MChequeReject struct {
	CustId        string         `gorm:"cust_id" json:"cust_id"`
	ChqRejectId   *int           `gorm:"column:chq_reject_id;primaryKey" json:"chq_reject_id"`
	ChqRejectName *string        `gorm:"chq_reject_name" json:"chq_reject_name"`
	CreatedBy     *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"updated_at" json:"updated_at"`
	IsDel         bool           `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `grom:"deleted_at" json:"deleted_at"`
}

func (m *MChequeReject) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}
func (MChequeReject) TableName() string {
	return "acf.m_cheque_reject"
}

type MChequeRejectList struct {
	CustId        string         `gorm:"cust_id" json:"cust_id"`
	ChqRejectId   *int           `gorm:"column:chq_reject_id;primaryKey" json:"chq_reject_id"`
	ChqRejectName *string        `gorm:"chq_reject_name" json:"chq_reject_name"`
	CreatedBy     *int64         `gorm:"created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"updated_by" json:"updated_by"`
	UpdatedAt     time.Time      `gorm:"updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `grom:"deleted_at" json:"deleted_at"`
}

func (MChequeRejectList) TableName() string {
	return "acf.m_cheque_reject"
}
