package model

import (
	"time"
)

type PromoStatusLog struct {
	CustID        string    `gorm:"column:cust_id" json:"cust_id"`
	PromoID       string    `gorm:"column:promo_id" json:"promo_id"`
	PromoStatusID int       `gorm:"column:promo_status_id" json:"promo_status_id"`
	Remarks       string    `gorm:"column:remarks" json:"remarks,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PromoStatusLog) TableName() string {
	return "sls.promo_status_logs"
}
