package model

import (
	"time"

	"gorm.io/gorm"
)

type WorkDay struct {
	PerYear  int        `gorm:"column:per_year" json:"per_year"`
	PerId    int        `gorm:"column:per_id" json:"per_id"`
	WeekId   int        `gorm:"column:week_id" json:"week_id"`
	WorkDate *time.Time `gorm:"column:work_date" json:"work_date"`
	IsWork   *bool      `gorm:"column:is_work" json:"is_work"`
	IsActive *bool      `gorm:"column:is_active" json:"is_active"`
	IsClosed *bool      `gorm:"column:is_closed" json:"is_closed" `
}

func (WorkDay) TableName() string {
	return "mst.m_work_day"
}

func (m *WorkDay) BeforeCreate(trx *gorm.DB) (err error) {

	return nil
}
