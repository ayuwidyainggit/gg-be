package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderApprovalRequest struct {
	OrderApprovalRequestID *int64    `gorm:"column:order_approval_request_id;default:nextval('sls.order_approval_request_seq'::regclass);not null" json:"order_approval_request_id"`
	CustID                 string    `gorm:"column:cust_id" json:"cust_id"`
	RoNo                   string    `gorm:"column:ro_no" json:"ro_no"`
	CreatedBy              *int64    `gorm:"column:created_by" json:"created_by"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
}

func (OrderApprovalRequest) TableName() string {
	return "sls.order_approval_requests"
}

func (m *OrderApprovalRequest) BeforeCreate(trx *gorm.DB) (err error) {

	m.CreatedAt = time.Now()

	return nil
}

type OrderApprovalRequestRead struct {
	OrderApprovalRequestID int64      `gorm:"column:order_approval_request_id" json:"order_approval_request_id"`
	CustID                 string     `gorm:"column:cust_id" json:"cust_id"`
	RoNo                   string     `gorm:"column:ro_no" json:"ro_no"`
	CreatedBy              *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt              time.Time  `gorm:"column:created_at" json:"created_at"`
	FinishedAt             *time.Time `gorm:"column:finished_at" json:"finished_at"`
}

func (OrderApprovalRequestRead) TableName() string {
	return "sls.order_approval_requests"
}
