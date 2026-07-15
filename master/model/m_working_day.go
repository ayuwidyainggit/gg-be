package model

import "time"

type MWorkingDay struct {
	CustId               string     `db:"cust_id" json:"cust_id"`
	PerYear              int        `db:"per_year" json:"per_year"`
	PerId                int        `db:"per_id" json:"per_id"`
	WeekId               int        `db:"week_id" json:"week_id"`
	WorkDate             *string    `db:"work_date" json:"work_date"`
	WorkDayId            *int       `db:"work_day_id" json:"work_day_id"`
	WorkingDayCalendarId *int64     `db:"working_day_calendar_id" json:"working_day_calendar_id"`
	HolidaySource        *string    `db:"holiday_source" json:"holiday_source"`
	HolidayNote          *string    `db:"holiday_note" json:"holiday_note"`
	IsActive             *bool      `db:"is_active" json:"is_active"`
	IsWork               *bool      `db:"is_work" json:"is_work"`
	IsClosed             *bool      `db:"is_closed" json:"is_closed"`
	ClosedAt             *time.Time `db:"closed_at" json:"closed_at"`
	ClosedBy             *string    `db:"closed_by" json:"closed_by"`
}

type MWorkingDayUpdate struct {
	PerYear  *int    `json:"per_year" sql:"per_year"`
	PerId    *int    `json:"per_id" sql:"per_id"`
	WeekId   *int    `json:"week_id" db:"week_id"`
	WorkDate *string `json:"work_date" db:"work_date"`
	IsWork   *bool   `json:"is_work" db:"is_work"`
	IsClosed *bool   `json:"is_closed" db:"is_closed"`
}

type MWorkingDayActive struct {
	PerYear  int        `json:"per_year" db:"per_year"`
	PerId    int        `json:"per_id" db:"per_id"`
	WeekId   int        `json:"week_id" db:"week_id"`
	WorkDate *time.Time `json:"work_date" db:"work_date"`
	IsWork   *bool      `json:"is_work" db:"is_work"`
	IsActive *bool      `json:"is_active" db:"is_active"`
	IsClosed *bool      `json:"is_closed" db:"is_closed" `
}
