package model

import (
	"time"
)

type MWeek struct {
	PerYear      int        `gorm:"column:per_year" json:"per_year"`
	PerId        int        `gorm:"column:per_id" json:"per_id"`
	WeekId       int        `gorm:"column:week_id" json:"week_id"`
	WeekStart    time.Time  `gorm:"column:week_start" json:"week_start"`
	WeekEnd      time.Time  `gorm:"column:week_end" json:"week_end"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	IsClosed     bool       `gorm:"column:is_closed" json:"is_closed"`
	ClosedAt     *time.Time `gorm:"column:closed_at" json:"closed_at"`
	ClosedBy     *int       `gorm:"column:closed_by" json:"closed_by"`
	ClosedByName *string    `gorm:"column:closed_by_name" json:"closed_by_name"`
	CreatedAt    *time.Time `gorm:"column:created_date" json:"created_date"`
}

func (MWeek) TableName() string {
	return "mst.m_week"
}

type MWorkDay struct {
	PerYear      int        `gorm:"column:per_year" json:"per_year"`
	PerId        int        `gorm:"column:per_id" json:"per_id"`
	WeekId       int        `gorm:"column:week_id" json:"week_id"`
	WorkDate     time.Time  `gorm:"column:work_date" json:"work_date"`
	IsWork       bool       `gorm:"column:is_work" json:"is_work"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	IsClosed     bool       `gorm:"column:is_closed" json:"is_closed"`
	ClosedAt     *time.Time `gorm:"column:closed_at" json:"closed_at"`
	ClosedBy     *int       `gorm:"column:closed_by" json:"closed_by"`
	ClosedByName *string    `gorm:"column:closed_by_name" json:"closed_by_name"`
}

func (MWorkDay) TableName() string {
	return "mst.m_work_day"
}

type WeekListDetail struct {
	PerYear      int        `gorm:"column:per_year" json:"per_year"`
	PerId        int        `gorm:"column:per_id" json:"per_id"`
	WeekId       int        `gorm:"column:week_id" json:"week_id"`
	WeekStart    time.Time  `gorm:"column:week_start" json:"week_start"`
	WeekEnd      time.Time  `gorm:"column:week_end" json:"week_end"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	IsClosed     bool       `gorm:"column:is_closed" json:"is_closed"`
	ClosedAt     *time.Time `gorm:"column:closed_at" json:"closed_at"`
	ClosedBy     int        `gorm:"column:closed_by" json:"closed_by"`
	ClosedByName string     `gorm:"column:closed_by_name" json:"closed_by_name"`
}

type WorkDayDetail struct {
	PerYear      int        `gorm:"column:per_year" json:"per_year"`
	PerId        int        `gorm:"column:per_id" json:"per_id"`
	WeekId       int        `gorm:"column:week_id" json:"week_id"`
	WorkDate     time.Time  `gorm:"column:work_date" json:"work_date"`
	IsWork       bool       `gorm:"column:is_work" json:"is_work"`
	IsActive     bool       `gorm:"column:is_active" json:"is_active"`
	IsClosed     bool       `gorm:"column:is_closed" json:"is_closed"`
	ClosedAt     *time.Time `gorm:"column:closed_at" json:"closed_at"`
	ClosedBy     int        `gorm:"column:closed_by" json:"closed_by"`
	ClosedByName string     `gorm:"column:closed_by_name" json:"closed_by_name"`
}

// DestinationCount holds the count of outlets and distributors for principal users
type DestinationCount struct {
	TotalOutlet      int `gorm:"column:total_outlet" json:"total_outlet"`
	TotalDistributor int `gorm:"column:total_distributor" json:"total_distributor"`
}
