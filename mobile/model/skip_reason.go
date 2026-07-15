package model

import "time"

type SkipReason struct {
	CustId         string    `gorm:"cust_id" json:"cust_id"`
	SkipReasonId   int64     `gorm:"skip_reason_id" json:"skip_reason_id"`
	SkipReasonCode string    `gorm:"skip_reason_code" json:"skip_reason_code"`
	SkipReasonName string    `gorm:"skip_reason_name" json:"skip_reason_name"`
	IsActive       bool      `gorm:"is_active" json:"is_active"`
	CreatedBy      int64     `gorm:"created_by" json:"created_by"`
	CreatedAt      time.Time `gorm:"created_at" json:"created_at"`
	UpdatedBy      int64     `gorm:"updated_by" json:"updated_by"`
	UpdatedAt      time.Time `gorm:"updated_at" json:"updated_at"`
	IsDel          bool      `gorm:"is_del" json:"is_del"`
	DeletedBy      int64     `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt      time.Time `gorm:"deleted_at" json:"deleted_at"`
}

func (SkipReason) TableName() string {
	return "mst.m_skip_reason"
}
