package model

import "time"

type ShipmentOrderStatus struct {
	ID          int       `json:"id"`
	OrderNo     *string   `json:"order_no"`
	StatusOrder string    `json:"status_order"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ShipmentOrderStatus) TableName() string {
	return "tms.shipment_order_status"
}
