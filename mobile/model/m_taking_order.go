package model

import (
	"time"

	"gorm.io/gorm"
)

type MTakingOrder struct {
	CustId          string         `gorm:"column:cust_id" json:"cust_id" `
	TakingOrderId   int64          `gorm:"column:taking_order_id" json:"taking_order_id"`
	TakingOrderName string         `gorm:"column:taking_order_name" json:"taking_order_name"`
	ImageUrl        string         `gorm:"column:image_url" json:"image_url"`
	IsActive        bool           `gorm:"column:is_active" json:"is_active"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt       *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel           bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy       *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MTakingOrder) TableName() string {
	return "mst.m_taking_order"
}
