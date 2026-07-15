package generator

import (
	"errors"
	"fmt"
	"time"
)

const (
	WorkingDayHolidaySourceDefault         = "default"
	WorkingDayHolidaySourceImported        = "imported"
	WorkingDayHolidaySourceDefaultImported = "default_imported"
)

type WorkingDayCalendarInput struct {
	StartDate              time.Time
	NumberOfWeeks          int
	FirstWeekID            int
	DefaultHolidayWeekdays []time.Weekday
	ImportedHolidays       []WorkingDayImportedHoliday
}

type WorkingDayImportedHoliday struct {
	Date  time.Time
	Notes string
}

type WorkingDayCalendarResult struct {
	StartDate time.Time
	EndDate   time.Time
	Weeks     []WorkingDayCalendarWeek
	Days      []WorkingDayCalendarDay
}

type WorkingDayCalendarWeek struct {
	WeekID         int
	CalendarWeekNo int
	WeekStart      time.Time
	WeekEnd        time.Time
}

type WorkingDayCalendarDay struct {
	WorkDate       time.Time
	WeekID         int
	CalendarWeekNo int
	IsWork         bool
	HolidaySource  *string
	HolidayNote    *string
}

func GenerateWorkingDayCalendar(input WorkingDayCalendarInput) (WorkingDayCalendarResult, error) {
	if input.StartDate.IsZero() {
		return WorkingDayCalendarResult{}, errors.New("start date is required")
	}
	if input.NumberOfWeeks < 1 || input.NumberOfWeeks > 99 {
		return WorkingDayCalendarResult{}, errors.New("number of weeks must be between 1 and 99")
	}
	if input.FirstWeekID < 1 {
		return WorkingDayCalendarResult{}, errors.New("first week id must be greater than zero")
	}

	startDate := normalizeCalendarDate(input.StartDate)
	endDate := startDate.AddDate(0, 0, input.NumberOfWeeks*7-1)
	defaultHolidaySet := buildDefaultHolidaySet(input.DefaultHolidayWeekdays)
	importedHolidayMap, err := buildImportedHolidayMap(input.ImportedHolidays, startDate, endDate)
	if err != nil {
		return WorkingDayCalendarResult{}, err
	}

	result := WorkingDayCalendarResult{
		StartDate: startDate,
		EndDate:   endDate,
	}

	for weekOffset := 0; weekOffset < input.NumberOfWeeks; weekOffset++ {
		weekStart := startDate.AddDate(0, 0, weekOffset*7)
		weekEnd := weekStart.AddDate(0, 0, 6)
		weekID := input.FirstWeekID + weekOffset
		calendarWeekNo := weekOffset + 1
		result.Weeks = append(result.Weeks, WorkingDayCalendarWeek{
			WeekID:         weekID,
			CalendarWeekNo: calendarWeekNo,
			WeekStart:      weekStart,
			WeekEnd:        weekEnd,
		})

		for dayOffset := 0; dayOffset < 7; dayOffset++ {
			workDate := weekStart.AddDate(0, 0, dayOffset)
			dateKey := calendarDateKey(workDate)
			importedNote, imported := importedHolidayMap[dateKey]
			isDefault := defaultHolidaySet[workDate.Weekday()]

			isWork := !(isDefault || imported)
			holidaySource, holidayNote := calendarHolidayFields(isDefault, imported, importedNote)
			result.Days = append(result.Days, WorkingDayCalendarDay{
				WorkDate:       workDate,
				WeekID:         weekID,
				CalendarWeekNo: calendarWeekNo,
				IsWork:         isWork,
				HolidaySource:  holidaySource,
				HolidayNote:    holidayNote,
			})
		}
	}

	return result, nil
}

func buildDefaultHolidaySet(weekdays []time.Weekday) map[time.Weekday]bool {
	result := map[time.Weekday]bool{}
	for _, weekday := range weekdays {
		if weekday >= time.Sunday && weekday <= time.Saturday {
			result[weekday] = true
		}
	}
	return result
}

func buildImportedHolidayMap(holidays []WorkingDayImportedHoliday, startDate, endDate time.Time) (map[string]string, error) {
	result := map[string]string{}
	for _, holiday := range holidays {
		date := normalizeCalendarDate(holiday.Date)
		if date.Before(startDate) || date.After(endDate) {
			return nil, fmt.Errorf("imported holiday %s is outside calendar range", calendarDateKey(date))
		}
		key := calendarDateKey(date)
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("duplicate imported holiday %s", key)
		}
		result[key] = holiday.Notes
	}
	return result, nil
}

func calendarHolidayFields(isDefault, imported bool, importedNote string) (*string, *string) {
	var source string
	switch {
	case isDefault && imported:
		source = WorkingDayHolidaySourceDefaultImported
	case isDefault:
		source = WorkingDayHolidaySourceDefault
	case imported:
		source = WorkingDayHolidaySourceImported
	default:
		return nil, nil
	}

	var note *string
	if imported {
		note = &importedNote
	}
	return &source, note
}

func normalizeCalendarDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func calendarDateKey(t time.Time) string {
	return normalizeCalendarDate(t).Format("2006-01-02")
}
