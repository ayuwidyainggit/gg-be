package model

import (
	"time"

	"gorm.io/gorm"
)

type Visit struct {
	CustID        string    `gorm:"column:cust_id" json:"cust_id"`
	VisitId       *int64    `gorm:"column:visit_id;primaryKey" json:"visit_id"`
	EmpCode       *string   `gorm:"column:emp_code" json:"emp_code"`
	Type          *int      `gorm:"column:type" json:"type"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	Latitude      *string   `gorm:"column:latitude" json:"latitude"`
	Longitude     *string   `gorm:"column:longitude" json:"longitude"`
	OutletCode    *string   `gorm:"column:outlet_code" json:"outlet_code"`
	IsInOutlet    *bool     `gorm:"column:is_in_outlet" json:"is_in_outlet"`
	Reason        *string   `gorm:"column:reason" json:"reason"`
	UpComingVisit *string   `gorm:"column:upcoming_visit" json:"upcoming_visit"`
}

func (m *Visit) BeforeCreate(trx *gorm.DB) (err error) {
	return nil
}
func (Visit) TableName() string {
	return "mobile.visits"
}

type VisitRead struct {
	VisitId   *int64    `gorm:"column:visit_id;primaryKey" json:"visit_id"`
	EmpCode   *string   `gorm:"column:emp_code" json:"emp_code"`
	Latitude  *string   `gorm:"column:latitude" json:"latitude"`
	Longitude *string   `gorm:"column:longitude" json:"longitude"`
	Type      *int      `gorm:"column:type" json:"type"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (VisitRead) TableName() string {
	return "mobile.visits"
}
