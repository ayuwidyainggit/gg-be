package model

import (
	"time"

	"github.com/lib/pq"
)

type WorkingDayCalendar struct {
	WorkingDayCalendarID int64         `db:"working_day_calendar_id" json:"working_day_calendar_id"`
	CustID               string        `db:"cust_id" json:"cust_id"`
	Title                string        `db:"title" json:"title"`
	StartDate            time.Time     `db:"start_date" json:"start_date"`
	NumberOfWeeks        int           `db:"number_of_weeks" json:"number_of_weeks"`
	EndDate              time.Time     `db:"end_date" json:"end_date"`
	DefaultHolidays      pq.Int64Array `db:"default_holidays" json:"default_holidays"`
	IsClosed             bool          `db:"is_closed" json:"is_closed"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
	CreatedBy            *int64        `db:"created_by" json:"created_by"`
	UpdatedAt            *time.Time    `db:"updated_at" json:"updated_at"`
	UpdatedBy            *int64        `db:"updated_by" json:"updated_by"`
	ClosedAt             *time.Time    `db:"closed_at" json:"closed_at"`
	ClosedBy             *int64        `db:"closed_by" json:"closed_by"`
	ClosedByName         *string       `db:"closed_by_name" json:"closed_by_name"`
}

type WorkingDayCalendarDay struct {
	WorkDate          time.Time `db:"work_date" json:"work_date"`
	WeekID            int       `db:"week_id" json:"week_id"`
	CalendarWeekNo    int       `db:"calendar_week_no" json:"calendar_week_no"`
	IsWork            bool      `db:"is_work" json:"is_work"`
	HolidaySource     *string   `db:"holiday_source" json:"holiday_source"`
	HolidayNote       *string   `db:"holiday_note" json:"holiday_note"`
	IsDefaultHoliday  bool      `db:"is_default_holiday" json:"is_default_holiday"`
	IsImportedHoliday bool      `db:"is_imported_holiday" json:"is_imported_holiday"`
}

type WorkingDayCalendarHoliday struct {
	WorkingDayCalendarHolidayID int64     `db:"working_day_calendar_holiday_id" json:"working_day_calendar_holiday_id"`
	WorkingDayCalendarID        int64     `db:"working_day_calendar_id" json:"working_day_calendar_id"`
	DistributorCustID           *string   `db:"distributor_cust_id" json:"distributor_cust_id,omitempty"`
	HolidayDate                 time.Time `db:"holiday_date" json:"holiday_date"`
	Notes                       string    `db:"notes" json:"notes"`
	CreatedAt                   time.Time `db:"created_at" json:"created_at"`
	CreatedBy                   *int64    `db:"created_by" json:"created_by"`
}
