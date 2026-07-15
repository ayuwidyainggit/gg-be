package model

import "time"

type MWeek struct {
	CustId               string     `json:"cust_id" db:"cust_id"`
	PerYear              int        `json:"per_year" db:"per_year"`
	PerId                int        `json:"per_id" db:"per_id"`
	WeekId               int        `json:"week_id" db:"week_id"`
	WeekStart            *string    `json:"week_start" db:"week_start"`
	WeekEnd              *string    `json:"week_end" db:"week_end"`
	WorkingDayCalendarId *int64     `json:"working_day_calendar_id" db:"working_day_calendar_id"`
	CalendarWeekNo       *int       `json:"calendar_week_no" db:"calendar_week_no"`
	IsActive             *bool      `json:"is_active" db:"is_active"`
	IsClosed             bool       `json:"is_closed" db:"is_closed"`
	ClosedAt             *time.Time `json:"closed_at" db:"closed_at"`
	ClosedBy             *int64     `json:"closed_by" db:"closed_by"`
	ClosedByName         *string    `json:"closed_by_name" db:"closed_by_name"`
}

type MWeekUpdate struct {
	PerYear   *int    `json:"per_year" sql:"per_year"`
	PerId     *int    `json:"per_id" sql:"per_id"`
	WeekId    *int    `json:"week_id" sql:"week_id"`
	WeekStart *string `json:"week_start" sql:"week_start"`
	WeekEnd   *string `json:"week_end" sql:"week_end"`
	IsActive  *bool   `json:"is_active" sql:"is_active"`
}
