package generator

import (
	"testing"
	"time"
)

func TestGenerateWorkingDayCalendarContinuesWeekIDAndStartsCalendarWeekNo(t *testing.T) {
	result, err := GenerateWorkingDayCalendar(WorkingDayCalendarInput{
		StartDate:     date(2026, 1, 1),
		NumberOfWeeks: 2,
		FirstWeekID:   6,
	})
	if err != nil {
		t.Fatalf("GenerateWorkingDayCalendar returned error: %v", err)
	}

	if got := len(result.Weeks); got != 2 {
		t.Fatalf("expected 2 weeks, got %d", got)
	}
	if got := len(result.Days); got != 14 {
		t.Fatalf("expected 14 days, got %d", got)
	}

	assertWeek(t, result.Weeks[0], 6, 1, "2026-01-01", "2026-01-07")
	assertWeek(t, result.Weeks[1], 7, 2, "2026-01-08", "2026-01-14")

	if got := result.EndDate.Format("2006-01-02"); got != "2026-01-14" {
		t.Fatalf("expected end date 2026-01-14, got %s", got)
	}
}

func TestGenerateWorkingDayCalendarAppliesDefaultAndImportedHolidays(t *testing.T) {
	result, err := GenerateWorkingDayCalendar(WorkingDayCalendarInput{
		StartDate:              date(2026, 1, 1),
		NumberOfWeeks:          1,
		FirstWeekID:            1,
		DefaultHolidayWeekdays: []time.Weekday{time.Sunday},
		ImportedHolidays: []WorkingDayImportedHoliday{
			{Date: date(2026, 1, 4), Notes: "Imported Sunday"},
			{Date: date(2026, 1, 5), Notes: "Special holiday"},
		},
	})
	if err != nil {
		t.Fatalf("GenerateWorkingDayCalendar returned error: %v", err)
	}

	jan4 := findDay(t, result.Days, "2026-01-04")
	if jan4.IsWork {
		t.Fatalf("expected 2026-01-04 to be non-work")
	}
	if jan4.HolidaySource == nil || *jan4.HolidaySource != WorkingDayHolidaySourceDefaultImported {
		t.Fatalf("expected default_imported source for 2026-01-04, got %v", jan4.HolidaySource)
	}
	if jan4.HolidayNote == nil || *jan4.HolidayNote != "Imported Sunday" {
		t.Fatalf("expected imported note for 2026-01-04, got %v", jan4.HolidayNote)
	}

	jan5 := findDay(t, result.Days, "2026-01-05")
	if jan5.IsWork {
		t.Fatalf("expected 2026-01-05 to be non-work")
	}
	if jan5.HolidaySource == nil || *jan5.HolidaySource != WorkingDayHolidaySourceImported {
		t.Fatalf("expected imported source for 2026-01-05, got %v", jan5.HolidaySource)
	}
	if jan5.HolidayNote == nil || *jan5.HolidayNote != "Special holiday" {
		t.Fatalf("expected imported note for 2026-01-05, got %v", jan5.HolidayNote)
	}

	jan6 := findDay(t, result.Days, "2026-01-06")
	if !jan6.IsWork {
		t.Fatalf("expected 2026-01-06 to be work day")
	}
	if jan6.HolidaySource != nil || jan6.HolidayNote != nil {
		t.Fatalf("expected no holiday fields for 2026-01-06, got source=%v note=%v", jan6.HolidaySource, jan6.HolidayNote)
	}
}

func TestGenerateWorkingDayCalendarRejectsImportedHolidayOutsideRange(t *testing.T) {
	_, err := GenerateWorkingDayCalendar(WorkingDayCalendarInput{
		StartDate:     date(2026, 1, 1),
		NumberOfWeeks: 1,
		FirstWeekID:   1,
		ImportedHolidays: []WorkingDayImportedHoliday{
			{Date: date(2026, 1, 8), Notes: "Outside"},
		},
	})
	if err == nil {
		t.Fatal("expected outside range imported holiday to fail")
	}
}

func TestGenerateWorkingDayCalendarCrossesYearBoundary(t *testing.T) {
	result, err := GenerateWorkingDayCalendar(WorkingDayCalendarInput{
		StartDate:     date(2025, 12, 29),
		NumberOfWeeks: 2,
		FirstWeekID:   52,
	})
	if err != nil {
		t.Fatalf("GenerateWorkingDayCalendar returned error: %v", err)
	}

	assertWeek(t, result.Weeks[0], 52, 1, "2025-12-29", "2026-01-04")
	assertWeek(t, result.Weeks[1], 53, 2, "2026-01-05", "2026-01-11")
}

func assertWeek(t *testing.T, week WorkingDayCalendarWeek, weekID, calendarWeekNo int, start, end string) {
	t.Helper()
	if week.WeekID != weekID {
		t.Fatalf("expected week_id %d, got %d", weekID, week.WeekID)
	}
	if week.CalendarWeekNo != calendarWeekNo {
		t.Fatalf("expected calendar_week_no %d, got %d", calendarWeekNo, week.CalendarWeekNo)
	}
	if got := week.WeekStart.Format("2006-01-02"); got != start {
		t.Fatalf("expected week start %s, got %s", start, got)
	}
	if got := week.WeekEnd.Format("2006-01-02"); got != end {
		t.Fatalf("expected week end %s, got %s", end, got)
	}
}

func findDay(t *testing.T, days []WorkingDayCalendarDay, date string) WorkingDayCalendarDay {
	t.Helper()
	for _, day := range days {
		if day.WorkDate.Format("2006-01-02") == date {
			return day
		}
	}
	t.Fatalf("day %s not found", date)
	return WorkingDayCalendarDay{}
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
