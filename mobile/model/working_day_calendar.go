package model

import (
	"time"
)

type WorkingDayCalendarDetail struct {
	WorkingDayCalendarID int       `gorm:"column:working_day_calendar_id"`
	CustID               string    `gorm:"column:cust_id"`
	Title                string    `gorm:"column:title"`
	StartDate            time.Time `gorm:"column:start_date"`
	NumberOfWeeks        int       `gorm:"column:number_of_weeks"`
	EndDate              time.Time `gorm:"column:end_date"`
	DefaultHolidays      string    `gorm:"column:default_holidays"`
	IsClosed             bool      `gorm:"column:is_closed"`
	IsActive             bool      `gorm:"column:is_active"`
}

func (WorkingDayCalendarDetail) TableName() string {
	return "mst.working_day_calendar"
}

type WorkingDayCalendarMonthDetail struct {
	IsActive  bool   `gorm:"column:is_active"`
	Month     int    `gorm:"column:month"`
	Year      int    `gorm:"column:year"`
	TextMonth string `gorm:"column:text_month"`
}
