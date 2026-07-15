package entity

import "time"

type WorkingDayCalendarQueryFilter struct {
	ParentCustID string `query:"parent_cust_id" validate:"required"`
}

type WorkingDayCalendarResponse struct {
	WorkingDayCalendarID int       `json:"working_day_calendar_id"`
	CustID               string    `json:"cust_id"`
	Title                string    `json:"title"`
	StartDate            time.Time `json:"start_date"`
	NumberOfWeeks        int       `json:"number_of_weeks"`
	EndDate              time.Time `json:"end_date"`
	DefaultHolidays      string    `json:"default_holidays"`
	IsClosed             bool      `json:"is_closed"`
	IsActive             bool      `json:"is_active"`
}

type WorkingDayCalendarMonthQueryFilter struct {
	WDCID int `query:"wdc_id" validate:"required"`
}

type WorkingDayCalendarMonthResponse struct {
	IsActive  bool   `json:"is_active"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	TextMonth string `json:"text_month"`
}
