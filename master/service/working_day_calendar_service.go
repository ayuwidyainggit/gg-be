package service

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"master/entity"
	"master/model"
	"master/pkg/generator"
	"master/repository"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

type WorkingDayCalendarService interface {
	List(entity.WorkingDayCalendarQueryFilter, string, string) ([]entity.WorkingDayCalendarListItem, int, int, error)
	Create(entity.CreateWorkingDayCalendarBody, string, string, int64) (entity.WorkingDayCalendarDetailResponse, error)
	Detail(int64, string, string) (entity.WorkingDayCalendarDetailResponse, error)
	Calendar(int64, entity.WorkingDayCalendarViewFilter, string, string) (entity.WorkingDayCalendarViewResponse, error)
	DownloadHolidayTemplate(int64, string, string, string) (*bytes.Buffer, string, string, error)
	ImportHolidays(int64, entity.WorkingDayCalendarImportHolidayRequest, string, string, int64) (entity.WorkingDayCalendarImportHolidayResponse, error)
}

type workingDayCalendarServiceImpl struct {
	repository repository.WorkingDayCalendarRepository
	now        func() time.Time
}

func NewWorkingDayCalendarService(repository repository.WorkingDayCalendarRepository) WorkingDayCalendarService {
	return &workingDayCalendarServiceImpl{
		repository: repository,
		now:        time.Now,
	}
}

func (s *workingDayCalendarServiceImpl) List(filter entity.WorkingDayCalendarQueryFilter, custID, parentCustID string) ([]entity.WorkingDayCalendarListItem, int, int, error) {
	filter.CustID = custID
	filter.ParentCustID = parentCustID
	ownerCustID := workingDayCalendarOwnerCustID(custID, parentCustID)
	rows, total, lastPage, err := s.repository.FindAll(filter, ownerCustID)
	if err != nil {
		return nil, 0, 0, err
	}

	items := make([]entity.WorkingDayCalendarListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.WorkingDayCalendarListItem{
			WorkingDayCalendarID: row.WorkingDayCalendarID,
			Title:                row.Title,
			StartDate:            formatCalendarDate(row.StartDate),
			EndDate:              formatCalendarDate(row.EndDate),
			NumberOfWeeks:        row.NumberOfWeeks,
			DefaultHolidays:      int64ArrayToInts(row.DefaultHolidays),
		})
	}
	return items, total, lastPage, nil
}

func (s *workingDayCalendarServiceImpl) Create(request entity.CreateWorkingDayCalendarBody, custID, parentCustID string, userID int64) (entity.WorkingDayCalendarDetailResponse, error) {
	if custID != parentCustID {
		return entity.WorkingDayCalendarDetailResponse{}, errors.New("working day calendar can only be created by principal")
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		return entity.WorkingDayCalendarDetailResponse{}, errors.New("title is required")
	}
	if len([]rune(title)) > 100 {
		return entity.WorkingDayCalendarDetailResponse{}, errors.New("title maximum is 100 characters")
	}
	if request.NumberOfWeeks < 1 || request.NumberOfWeeks > 99 {
		return entity.WorkingDayCalendarDetailResponse{}, errors.New("number_of_weeks must be between 1 and 99")
	}

	startDate, err := parseWorkingDayCalendarDate(request.StartDate)
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}

	latestCalendar, err := s.repository.FindLatestCalendarByCustID(custID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}
	if err == nil && !startDate.After(normalizeServiceCalendarDate(latestCalendar.EndDate)) {
		return entity.WorkingDayCalendarDetailResponse{}, fmt.Errorf("start_date must be after latest calendar end date %s", formatCalendarDate(latestCalendar.EndDate))
	}

	perYear := workingDayCalendarPerYear(startDate)

	defaultWeekdays, err := intsToWeekdays(request.DefaultHolidays)
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}
	generated, err := generator.GenerateWorkingDayCalendar(generator.WorkingDayCalendarInput{
		StartDate:              startDate,
		NumberOfWeeks:          request.NumberOfWeeks,
		FirstWeekID:            1,
		DefaultHolidayWeekdays: defaultWeekdays,
	})
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}

	calendar := model.WorkingDayCalendar{
		CustID:          custID,
		Title:           title,
		StartDate:       generated.StartDate,
		NumberOfWeeks:   request.NumberOfWeeks,
		EndDate:         generated.EndDate,
		DefaultHolidays: intsToInt64Array(request.DefaultHolidays),
		CreatedBy:       &userID,
	}
	weeks, days := materializeWorkingDayCalendarRows(custID, perYear, generated)
	calendarID, err := s.repository.StoreCalendarWithDetails(calendar, nil, weeks, days)
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}
	calendar.WorkingDayCalendarID = calendarID

	return calendarDetailResponse(calendar, nil), nil
}

func (s *workingDayCalendarServiceImpl) Detail(id int64, custID, parentCustID string) (entity.WorkingDayCalendarDetailResponse, error) {
	ownerCustID := workingDayCalendarOwnerCustID(custID, parentCustID)
	calendar, err := s.repository.FindByID(id, ownerCustID)
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}
	holidays, err := s.repository.FindImportedHolidays(id)
	if err != nil {
		return entity.WorkingDayCalendarDetailResponse{}, err
	}
	return calendarDetailResponse(calendar, holidays), nil
}

func (s *workingDayCalendarServiceImpl) Calendar(id int64, filter entity.WorkingDayCalendarViewFilter, custID, parentCustID string) (entity.WorkingDayCalendarViewResponse, error) {
	ownerCustID := workingDayCalendarOwnerCustID(custID, parentCustID)
	calendar, err := s.repository.FindByID(id, ownerCustID)
	if err != nil {
		return entity.WorkingDayCalendarViewResponse{}, err
	}

	view := strings.ToLower(strings.TrimSpace(filter.View))
	if view == "" {
		view = "month"
	}
	dateFrom, dateTo, month, year, err := s.calendarDateWindow(view, filter)
	if err != nil {
		return entity.WorkingDayCalendarViewResponse{}, err
	}

	days, err := s.repository.FindCalendarDays(id, ownerCustID, dateFrom, dateTo)
	if err != nil {
		return entity.WorkingDayCalendarViewResponse{}, err
	}

	items := make([]entity.WorkingDayCalendarDateItem, 0, len(days))
	for _, day := range days {
		items = append(items, entity.WorkingDayCalendarDateItem{
			Date:              formatCalendarDate(day.WorkDate),
			WeekID:            day.WeekID,
			CalendarWeekNo:    day.CalendarWeekNo,
			WeekLabel:         fmt.Sprintf("Week %d", day.WeekID),
			IsWork:            day.IsWork,
			IsDefaultHoliday:  day.IsDefaultHoliday,
			IsImportedHoliday: day.IsImportedHoliday,
			Notes:             day.HolidayNote,
		})
	}

	return entity.WorkingDayCalendarViewResponse{
		WorkingDayCalendarID: calendar.WorkingDayCalendarID,
		Title:                calendar.Title,
		View:                 view,
		Month:                month,
		Year:                 year,
		Dates:                items,
	}, nil
}

func (s *workingDayCalendarServiceImpl) DownloadHolidayTemplate(id int64, format, custID, parentCustID string) (*bytes.Buffer, string, string, error) {
	ownerCustID := workingDayCalendarOwnerCustID(custID, parentCustID)
	holidays := []model.WorkingDayCalendarHoliday{}
	if id > 0 {
		if _, err := s.repository.FindByID(id, ownerCustID); err != nil {
			return nil, "", "", err
		}
		rows, err := s.repository.FindImportedHolidays(id)
		if err != nil {
			return nil, "", "", err
		}
		holidays = rows
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "csv":
		buf, err := createWorkingDayCalendarHolidayTemplateCSV(holidays)
		return buf, "text/csv", "working_day_calendar_holidays.csv", err
	case "xls":
		buf, err := createWorkingDayCalendarHolidayTemplateWorkbook(holidays)
		return buf, "application/vnd.ms-excel", "working_day_calendar_holidays.xls", err
	default:
		buf, err := createWorkingDayCalendarHolidayTemplateWorkbook(holidays)
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "working_day_calendar_holidays.xlsx", err
	}
}

func (s *workingDayCalendarServiceImpl) ImportHolidays(id int64, request entity.WorkingDayCalendarImportHolidayRequest, custID, parentCustID string, userID int64) (entity.WorkingDayCalendarImportHolidayResponse, error) {
	resp := entity.WorkingDayCalendarImportHolidayResponse{
		FileURL:     request.FileURL,
		ProcessedAt: s.now().Format(time.RFC3339),
	}
	if custID != parentCustID {
		return resp, errors.New("working day calendar holidays can only be imported by principal")
	}

	calendar, err := s.repository.FindByID(id, custID)
	if err != nil {
		return resp, err
	}

	rows, fileName, err := downloadWorkingDayCalendarHolidayRows(request.FileURL)
	if err != nil {
		return resp, err
	}
	resp.FileName = fileName

	headerIndex, err := detectWorkingDayCalendarHolidayHeader(rows)
	if err != nil {
		return resp, err
	}

	seen := map[string]struct{}{}
	failed := make([]string, 0)
	imported := make([]generator.WorkingDayImportedHoliday, 0)
	holidayRows := make([]model.WorkingDayCalendarHoliday, 0)
	for idx, row := range rows[headerIndex+1:] {
		actualRow := idx + headerIndex + 2
		if isWorkingDayCalendarHolidayRowEmpty(row) {
			continue
		}
		resp.TotalRow++

		dateText := workingDayCalendarHolidayCell(row, 0)
		notes := strings.TrimSpace(workingDayCalendarHolidayCell(row, 1))
		if dateText == "" {
			failed = append(failed, fmt.Sprintf("row %d: date is required", actualRow))
			continue
		}
		if notes == "" {
			failed = append(failed, fmt.Sprintf("row %d: notes is required", actualRow))
			continue
		}

		date, err := parseWorkingDayCalendarDate(dateText)
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}
		dateKey := formatCalendarDate(date)
		if date.Before(normalizeServiceCalendarDate(calendar.StartDate)) || date.After(normalizeServiceCalendarDate(calendar.EndDate)) {
			failed = append(failed, fmt.Sprintf("row %d: imported holiday %s is outside calendar range", actualRow, dateKey))
			continue
		}
		if _, exists := seen[dateKey]; exists {
			failed = append(failed, fmt.Sprintf("row %d: duplicate imported holiday %s", actualRow, dateKey))
			continue
		}
		seen[dateKey] = struct{}{}

		imported = append(imported, generator.WorkingDayImportedHoliday{Date: date, Notes: notes})
		createdBy := userID
		holidayRows = append(holidayRows, model.WorkingDayCalendarHoliday{
			WorkingDayCalendarID: id,
			HolidayDate:          date,
			Notes:                notes,
			CreatedBy:            &createdBy,
		})
	}

	resp.FailedReasons = failed
	resp.FailedRow = len(failed)
	if len(failed) > 0 {
		return resp, errors.New("import validation failed")
	}

	defaultWeekdays, err := intsToWeekdays(int64ArrayToInts(calendar.DefaultHolidays))
	if err != nil {
		return resp, err
	}
	generated, err := generator.GenerateWorkingDayCalendar(generator.WorkingDayCalendarInput{
		StartDate:              calendar.StartDate,
		NumberOfWeeks:          calendar.NumberOfWeeks,
		FirstWeekID:            1,
		DefaultHolidayWeekdays: defaultWeekdays,
		ImportedHolidays:       imported,
	})
	if err != nil {
		return resp, err
	}

	_, days := materializeWorkingDayCalendarRows(custID, workingDayCalendarPerYear(calendar.StartDate), generated)
	if err := s.repository.ReplaceImportedHolidaysAndWorkDays(id, custID, userID, holidayRows, days); err != nil {
		return resp, err
	}

	resp.SuccessRow = resp.TotalRow
	return resp, nil
}

func (s *workingDayCalendarServiceImpl) calendarDateWindow(view string, filter entity.WorkingDayCalendarViewFilter) (time.Time, time.Time, *int, int, error) {
	now := s.now()
	year := filter.Year
	if year == 0 {
		year = now.Year()
	}

	switch view {
	case "month":
		month := filter.Month
		if month == 0 {
			month = int(now.Month())
		}
		if month < 1 || month > 12 {
			return time.Time{}, time.Time{}, nil, 0, errors.New("month must be between 1 and 12")
		}
		dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		dateTo := dateFrom.AddDate(0, 1, -1)
		return dateFrom, dateTo, &month, year, nil
	case "year":
		dateFrom := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
		return dateFrom, dateTo, nil, year, nil
	default:
		return time.Time{}, time.Time{}, nil, 0, errors.New("view must be month or year")
	}
}

func materializeWorkingDayCalendarRows(principalCustID string, perYear int, generated generator.WorkingDayCalendarResult) ([]model.MWeek, []model.MWorkingDay) {
	isActive := false
	weeks := make([]model.MWeek, 0, len(generated.Weeks))
	days := make([]model.MWorkingDay, 0, len(generated.Days))
	for _, week := range generated.Weeks {
		weekStart := formatCalendarDate(week.WeekStart)
		weekEnd := formatCalendarDate(week.WeekEnd)
		calendarWeekNo := week.CalendarWeekNo
		weeks = append(weeks, model.MWeek{
			CustId:         principalCustID,
			PerYear:        perYear,
			PerId:          week.WeekID,
			WeekId:         week.WeekID,
			WeekStart:      &weekStart,
			WeekEnd:        &weekEnd,
			IsActive:       &isActive,
			CalendarWeekNo: &calendarWeekNo,
		})
	}
	for _, day := range generated.Days {
		workDate := formatCalendarDate(day.WorkDate)
		isWork := day.IsWork
		days = append(days, model.MWorkingDay{
			CustId:        principalCustID,
			PerYear:       perYear,
			PerId:         day.WeekID,
			WeekId:        day.WeekID,
			WorkDate:      &workDate,
			IsActive:      &isActive,
			IsWork:        &isWork,
			HolidaySource: day.HolidaySource,
			HolidayNote:   day.HolidayNote,
		})
	}
	return weeks, days
}

func calendarDetailResponse(calendar model.WorkingDayCalendar, holidays []model.WorkingDayCalendarHoliday) entity.WorkingDayCalendarDetailResponse {
	importedHolidays := make([]entity.WorkingDayCalendarHoliday, 0, len(holidays))
	for _, holiday := range holidays {
		importedHolidays = append(importedHolidays, entity.WorkingDayCalendarHoliday{
			DistributorCustID: holiday.DistributorCustID,
			Date:              formatCalendarDate(holiday.HolidayDate),
			Notes:             holiday.Notes,
		})
	}
	return entity.WorkingDayCalendarDetailResponse{
		WorkingDayCalendarID: calendar.WorkingDayCalendarID,
		Title:                calendar.Title,
		StartDate:            formatCalendarDate(calendar.StartDate),
		EndDate:              formatCalendarDate(calendar.EndDate),
		NumberOfWeeks:        calendar.NumberOfWeeks,
		DefaultHolidays:      int64ArrayToInts(calendar.DefaultHolidays),
		ImportedHolidays:     importedHolidays,
	}
}

func workingDayCalendarOwnerCustID(custID, parentCustID string) string {
	if parentCustID != "" {
		return parentCustID
	}
	return custID
}

func parseWorkingDayCalendarDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{
		"2006-01-02",
		"2006-1-2",
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, raw, time.UTC)
		if err == nil {
			return normalizeServiceCalendarDate(parsed), nil
		}
	}
	return time.Time{}, errors.New("date must use YYYY-MM-DD or DD/MM/YYYY format")
}

func normalizeServiceCalendarDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func workingDayCalendarPerYear(startDate time.Time) int {
	return startDate.AddDate(0, 0, 6).Year()
}

func intsToWeekdays(values []int) ([]time.Weekday, error) {
	result := make([]time.Weekday, 0, len(values))
	seen := map[int]bool{}
	for _, value := range values {
		if value < 0 || value > 6 {
			return nil, errors.New("default_holidays must contain weekday values between 0 and 6")
		}
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, time.Weekday(value))
	}
	return result, nil
}

func intsToInt64Array(values []int) pq.Int64Array {
	if len(values) == 0 {
		return nil
	}
	result := make(pq.Int64Array, 0, len(values))
	seen := map[int]bool{}
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, int64(value))
	}
	return result
}

func int64ArrayToInts(values pq.Int64Array) []int {
	result := make([]int, 0, len(values))
	for _, value := range values {
		result = append(result, int(value))
	}
	return result
}

func formatCalendarDate(t time.Time) string {
	return normalizeServiceCalendarDate(t).Format("2006-01-02")
}

func createWorkingDayCalendarHolidayTemplateWorkbook(holidays []model.WorkingDayCalendarHoliday) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Template"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}})
	headers := []string{"date", "notes"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	for i, holiday := range holidays {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), formatCalendarDate(holiday.HolidayDate))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), holiday.Notes)
	}
	if _, err := f.NewSheet("Instructions"); err == nil {
		f.SetCellValue("Instructions", "A1", "Column")
		f.SetCellValue("Instructions", "B1", "Instruction")
		f.SetCellValue("Instructions", "A2", "date")
		f.SetCellValue("Instructions", "B2", "Required. Use DD/MM/YYYY or YYYY-MM-DD. Date must be inside the selected calendar range.")
		f.SetCellValue("Instructions", "A3", "notes")
		f.SetCellValue("Instructions", "B3", "Required. Holiday description.")
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createWorkingDayCalendarHolidayTemplateCSV(holidays []model.WorkingDayCalendarHoliday) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	if err := writer.Write([]string{"date", "notes"}); err != nil {
		return nil, err
	}
	for _, holiday := range holidays {
		if err := writer.Write([]string{formatCalendarDate(holiday.HolidayDate), holiday.Notes}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf, writer.Error()
}

func downloadWorkingDayCalendarHolidayRows(fileURL string) ([][]string, string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	fileName := fileURL
	if parts := strings.Split(fileURL, "/"); len(parts) > 0 {
		fileName = parts[len(parts)-1]
	}
	if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		reader := csv.NewReader(bytes.NewReader(content))
		rows, err := reader.ReadAll()
		return rows, fileName, err
	}

	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = f.Close() }()
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, "", errors.New("excel file has no sheet")
	}
	rows, err := f.GetRows(sheetName)
	return rows, fileName, err
}

func detectWorkingDayCalendarHolidayHeader(rows [][]string) (int, error) {
	if len(rows) == 0 {
		return 0, errors.New("template does not contain header")
	}
	for idx, row := range rows {
		if normalizeWorkingDayCalendarHolidayHeader(workingDayCalendarHolidayCell(row, 0)) == "DATE" &&
			normalizeWorkingDayCalendarHolidayHeader(workingDayCalendarHolidayCell(row, 1)) == "NOTES" {
			return idx, nil
		}
	}
	return 0, errors.New("invalid template header")
}

func normalizeWorkingDayCalendarHolidayHeader(value string) string {
	return strings.TrimSpace(strings.ToUpper(value))
}

func workingDayCalendarHolidayCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func isWorkingDayCalendarHolidayRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
