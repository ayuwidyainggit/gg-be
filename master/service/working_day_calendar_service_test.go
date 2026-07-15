package service

import (
	"database/sql"
	"master/entity"
	"master/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
)

type fakeWorkingDayCalendarRepository struct {
	latestCalendar    model.WorkingDayCalendar
	latestCalendarErr error
	latestWeekID      int
	latestWeekIDErr   error
	storeCalendarID   int64
	storeErr          error
	listRows          []model.WorkingDayCalendar
	listTotal         int
	listLastPage      int
	listErr           error
	detail            model.WorkingDayCalendar
	detailErr         error
	holidays          []model.WorkingDayCalendarHoliday
	holidaysErr       error
	calendarDays      []model.WorkingDayCalendarDay
	calendarDaysErr   error

	findLatestCustID  string
	findLatestWeekID  []string
	findLatestYear    int
	storeCalendar     model.WorkingDayCalendar
	storeHolidays     []model.WorkingDayCalendarHoliday
	storeWeeks        []model.MWeek
	storeDays         []model.MWorkingDay
	findAllOwnerID    string
	findByID          int64
	findByOwnerID     string
	findDaysOwnerID   string
	findDaysFrom      time.Time
	findDaysTo        time.Time
	replaceCalendarID int64
	replaceCustID     string
	replaceUserID     int64
	replaceHolidays   []model.WorkingDayCalendarHoliday
	replaceDays       []model.MWorkingDay
	replaceErr        error
}

func (r *fakeWorkingDayCalendarRepository) FindLatestCalendarByCustID(custID string) (model.WorkingDayCalendar, error) {
	r.findLatestCustID = custID
	return r.latestCalendar, r.latestCalendarErr
}

func (r *fakeWorkingDayCalendarRepository) FindLatestWeekIDByCustIDs(custIDs []string, perYear int) (int, error) {
	r.findLatestWeekID = append([]string{}, custIDs...)
	r.findLatestYear = perYear
	return r.latestWeekID, r.latestWeekIDErr
}

func (r *fakeWorkingDayCalendarRepository) StoreCalendarWithDetails(calendar model.WorkingDayCalendar, holidays []model.WorkingDayCalendarHoliday, weeks []model.MWeek, days []model.MWorkingDay) (int64, error) {
	r.storeCalendar = calendar
	r.storeHolidays = append([]model.WorkingDayCalendarHoliday{}, holidays...)
	r.storeWeeks = append([]model.MWeek{}, weeks...)
	r.storeDays = append([]model.MWorkingDay{}, days...)
	return r.storeCalendarID, r.storeErr
}

func (r *fakeWorkingDayCalendarRepository) FindAll(filter entity.WorkingDayCalendarQueryFilter, ownerCustID string) ([]model.WorkingDayCalendar, int, int, error) {
	r.findAllOwnerID = ownerCustID
	return r.listRows, r.listTotal, r.listLastPage, r.listErr
}

func (r *fakeWorkingDayCalendarRepository) FindByID(id int64, ownerCustID string) (model.WorkingDayCalendar, error) {
	r.findByID = id
	r.findByOwnerID = ownerCustID
	return r.detail, r.detailErr
}

func (r *fakeWorkingDayCalendarRepository) FindImportedHolidays(id int64) ([]model.WorkingDayCalendarHoliday, error) {
	return r.holidays, r.holidaysErr
}

func (r *fakeWorkingDayCalendarRepository) FindCalendarDays(id int64, ownerCustID string, dateFrom, dateTo time.Time) ([]model.WorkingDayCalendarDay, error) {
	r.findDaysOwnerID = ownerCustID
	r.findDaysFrom = dateFrom
	r.findDaysTo = dateTo
	return r.calendarDays, r.calendarDaysErr
}

func (r *fakeWorkingDayCalendarRepository) ReplaceImportedHolidaysAndWorkDays(calendarID int64, custID string, userID int64, holidays []model.WorkingDayCalendarHoliday, days []model.MWorkingDay) error {
	r.replaceCalendarID = calendarID
	r.replaceCustID = custID
	r.replaceUserID = userID
	r.replaceHolidays = append([]model.WorkingDayCalendarHoliday{}, holidays...)
	r.replaceDays = append([]model.MWorkingDay{}, days...)
	return r.replaceErr
}

func TestWorkingDayCalendarServiceCreateMaterializesWeeksAndDays(t *testing.T) {
	repo := &fakeWorkingDayCalendarRepository{
		latestCalendarErr: sql.ErrNoRows,
		latestWeekID:      5,
		storeCalendarID:   42,
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	response, err := svc.Create(entity.CreateWorkingDayCalendarBody{
		Title:           "Ramadan 2026",
		StartDate:       "2026-01-01",
		NumberOfWeeks:   2,
		DefaultHolidays: []int{0},
	}, "P001", "P001", 99)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if response.WorkingDayCalendarID != 42 {
		t.Fatalf("response id = %d, want 42", response.WorkingDayCalendarID)
	}
	if repo.findLatestCustID != "P001" {
		t.Fatalf("latest calendar cust = %q, want P001", repo.findLatestCustID)
	}
	if repo.findLatestYear != 0 {
		t.Fatalf("week lookup year = %d, want no lookup", repo.findLatestYear)
	}
	if len(repo.findLatestWeekID) != 0 {
		t.Fatalf("week lookup cust ids = %v, want no lookup", repo.findLatestWeekID)
	}

	if repo.storeCalendar.Title != "Ramadan 2026" || formatCalendarDate(repo.storeCalendar.EndDate) != "2026-01-14" {
		t.Fatalf("stored calendar = %+v", repo.storeCalendar)
	}
	if got := int64ArrayToInts(repo.storeCalendar.DefaultHolidays); len(got) != 1 || got[0] != 0 {
		t.Fatalf("stored default holidays = %v, want [0]", got)
	}
	if len(repo.storeWeeks) != 2 {
		t.Fatalf("stored weeks = %d, want 2", len(repo.storeWeeks))
	}
	if len(repo.storeDays) != 14 {
		t.Fatalf("stored days = %d, want 14", len(repo.storeDays))
	}

	firstWeek := repo.storeWeeks[0]
	if firstWeek.CustId != "P001" || firstWeek.PerYear != 2026 || firstWeek.PerId != 1 || firstWeek.WeekId != 1 {
		t.Fatalf("first week ids = %+v", firstWeek)
	}
	if firstWeek.CalendarWeekNo == nil || *firstWeek.CalendarWeekNo != 1 {
		t.Fatalf("calendar_week_no = %v, want 1", firstWeek.CalendarWeekNo)
	}
	if firstWeek.WeekStart == nil || *firstWeek.WeekStart != "2026-01-01" || firstWeek.WeekEnd == nil || *firstWeek.WeekEnd != "2026-01-07" {
		t.Fatalf("first week dates = %+v", firstWeek)
	}

	sunday := findWorkDay(repo.storeDays, "P001", "2026-01-04")
	if sunday == nil {
		t.Fatal("expected Sunday work day to be materialized")
	}
	if sunday.IsWork == nil || *sunday.IsWork {
		t.Fatalf("Sunday is_work = %v, want false", sunday.IsWork)
	}
	if sunday.HolidaySource == nil || *sunday.HolidaySource != "default" {
		t.Fatalf("Sunday holiday_source = %v, want default", sunday.HolidaySource)
	}
}

func TestWorkingDayCalendarServiceCreateStoresNoDefaultHolidaysAsNull(t *testing.T) {
	repo := &fakeWorkingDayCalendarRepository{
		latestCalendarErr: sql.ErrNoRows,
		storeCalendarID:   42,
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	_, err := svc.Create(entity.CreateWorkingDayCalendarBody{
		Title:           "No Holidays",
		StartDate:       "2026-01-01",
		NumberOfWeeks:   1,
		DefaultHolidays: []int{},
	}, "P001", "P001", 99)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repo.storeCalendar.DefaultHolidays != nil {
		t.Fatalf("stored default holidays = %v, want nil", repo.storeCalendar.DefaultHolidays)
	}
}

func TestWorkingDayCalendarServiceCreateRejectsDistributorAndOverlappingCalendar(t *testing.T) {
	repo := &fakeWorkingDayCalendarRepository{}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	_, err := svc.Create(entity.CreateWorkingDayCalendarBody{
		Title:         "Title",
		StartDate:     "2026-01-15",
		NumberOfWeeks: 1,
	}, "D001", "P001", 99)
	if err == nil || !strings.Contains(err.Error(), "principal") {
		t.Fatalf("Create() distributor error = %v, want principal error", err)
	}

	repo.latestCalendar = model.WorkingDayCalendar{
		EndDate: mustDate(t, "2026-01-14"),
	}
	repo.latestCalendarErr = nil
	_, err = svc.Create(entity.CreateWorkingDayCalendarBody{
		Title:         "Title",
		StartDate:     "2026-01-14",
		NumberOfWeeks: 1,
	}, "P001", "P001", 99)
	if err == nil || !strings.Contains(err.Error(), "after latest calendar end date") {
		t.Fatalf("Create() overlap error = %v, want append-only error", err)
	}
	if len(repo.storeWeeks) != 0 || len(repo.storeDays) != 0 {
		t.Fatalf("stored rows after rejected create: weeks=%d days=%d", len(repo.storeWeeks), len(repo.storeDays))
	}
}

func TestWorkingDayCalendarServiceListDetailAndCalendar(t *testing.T) {
	note := "National holiday"
	repo := &fakeWorkingDayCalendarRepository{
		listRows: []model.WorkingDayCalendar{{
			WorkingDayCalendarID: 7,
			Title:                "Cycle A",
			StartDate:            mustDate(t, "2026-01-01"),
			EndDate:              mustDate(t, "2026-01-14"),
			NumberOfWeeks:        2,
			DefaultHolidays:      pq.Int64Array{0},
		}},
		listTotal:    1,
		listLastPage: 1,
		detail: model.WorkingDayCalendar{
			WorkingDayCalendarID: 7,
			Title:                "Cycle A",
			StartDate:            mustDate(t, "2026-01-01"),
			EndDate:              mustDate(t, "2026-01-14"),
			NumberOfWeeks:        2,
			DefaultHolidays:      pq.Int64Array{0},
		},
		holidays: []model.WorkingDayCalendarHoliday{{
			HolidayDate: mustDate(t, "2026-01-02"),
			Notes:       note,
		}},
		calendarDays: []model.WorkingDayCalendarDay{{
			WorkDate:          mustDate(t, "2026-01-02"),
			WeekID:            1,
			CalendarWeekNo:    1,
			IsWork:            false,
			HolidayNote:       &note,
			IsImportedHoliday: true,
		}},
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	list, total, lastPage, err := svc.List(entity.WorkingDayCalendarQueryFilter{}, "D001", "P001")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || lastPage != 1 || len(list) != 1 || repo.findAllOwnerID != "P001" {
		t.Fatalf("List() total=%d lastPage=%d list=%v owner=%q", total, lastPage, list, repo.findAllOwnerID)
	}

	detail, err := svc.Detail(7, "D001", "P001")
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}
	if detail.Title != "Cycle A" || len(detail.ImportedHolidays) != 1 || detail.ImportedHolidays[0].Date != "2026-01-02" {
		t.Fatalf("Detail() = %+v", detail)
	}

	calendar, err := svc.Calendar(7, entity.WorkingDayCalendarViewFilter{View: "month", Month: 1, Year: 2026}, "D001", "P001")
	if err != nil {
		t.Fatalf("Calendar() error = %v", err)
	}
	if calendar.View != "month" || calendar.Month == nil || *calendar.Month != 1 || calendar.Year != 2026 {
		t.Fatalf("Calendar() window = %+v", calendar)
	}
	if repo.findByID != 7 || repo.findByOwnerID != "P001" || repo.findDaysOwnerID != "P001" {
		t.Fatalf("owner lookup id=%d detailOwner=%q daysOwner=%q", repo.findByID, repo.findByOwnerID, repo.findDaysOwnerID)
	}
	if repo.findDaysFrom.Format("2006-01-02") != "2026-01-01" || repo.findDaysTo.Format("2006-01-02") != "2026-01-31" {
		t.Fatalf("days window = %s to %s", repo.findDaysFrom, repo.findDaysTo)
	}
	if len(calendar.Dates) != 1 || calendar.Dates[0].WeekLabel != "Week 1" || !calendar.Dates[0].IsImportedHoliday {
		t.Fatalf("Calendar dates = %+v", calendar.Dates)
	}
}

func TestWorkingDayCalendarServiceDownloadHolidayTemplateIncludesExistingRows(t *testing.T) {
	repo := &fakeWorkingDayCalendarRepository{
		detail: model.WorkingDayCalendar{
			WorkingDayCalendarID: 7,
			CustID:               "P001",
			Title:                "Cycle A",
		},
		holidays: []model.WorkingDayCalendarHoliday{{
			HolidayDate: mustDate(t, "2026-01-02"),
			Notes:       "National holiday",
		}},
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	buf, contentType, filename, err := svc.DownloadHolidayTemplate(7, "csv", "D001", "P001")
	if err != nil {
		t.Fatalf("DownloadHolidayTemplate() error = %v", err)
	}
	if contentType != "text/csv" || filename != "working_day_calendar_holidays.csv" {
		t.Fatalf("template metadata contentType=%q filename=%q", contentType, filename)
	}
	if !strings.Contains(buf.String(), "2026-01-02") || !strings.Contains(buf.String(), "National holiday") {
		t.Fatalf("template body = %q", buf.String())
	}
	if repo.findByOwnerID != "P001" {
		t.Fatalf("template owner = %q, want P001", repo.findByOwnerID)
	}
}

func TestWorkingDayCalendarServiceImportHolidaysReplacesAndRegeneratesDays(t *testing.T) {
	csvBody := "date,notes\n2026-01-04,Imported Sunday\n2026-01-05,National holiday\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(csvBody))
	}))
	defer server.Close()

	repo := &fakeWorkingDayCalendarRepository{
		detail: model.WorkingDayCalendar{
			WorkingDayCalendarID: 7,
			CustID:               "P001",
			Title:                "Cycle A",
			StartDate:            mustDate(t, "2026-01-01"),
			EndDate:              mustDate(t, "2026-01-07"),
			NumberOfWeeks:        1,
			DefaultHolidays:      pq.Int64Array{0},
		},
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	resp, err := svc.ImportHolidays(7, entity.WorkingDayCalendarImportHolidayRequest{FileURL: server.URL + "/holidays.csv"}, "P001", "P001", 99)
	if err != nil {
		t.Fatalf("ImportHolidays() error = %v resp=%+v", err, resp)
	}
	if resp.TotalRow != 2 || resp.SuccessRow != 2 || resp.FailedRow != 0 {
		t.Fatalf("ImportHolidays() resp = %+v", resp)
	}
	if repo.replaceCalendarID != 7 || repo.replaceCustID != "P001" || repo.replaceUserID != 99 {
		t.Fatalf("replace call id=%d cust=%q user=%d", repo.replaceCalendarID, repo.replaceCustID, repo.replaceUserID)
	}
	if len(repo.replaceHolidays) != 2 {
		t.Fatalf("replace holidays = %d, want 2", len(repo.replaceHolidays))
	}
	jan4 := findWorkDay(repo.replaceDays, "P001", "2026-01-04")
	if jan4 == nil || jan4.IsWork == nil || *jan4.IsWork || jan4.HolidaySource == nil || *jan4.HolidaySource != "default_imported" || jan4.HolidayNote == nil || *jan4.HolidayNote != "Imported Sunday" {
		t.Fatalf("jan4 regenerated day = %+v", jan4)
	}
	jan5 := findWorkDay(repo.replaceDays, "P001", "2026-01-05")
	if jan5 == nil || jan5.IsWork == nil || *jan5.IsWork || jan5.HolidaySource == nil || *jan5.HolidaySource != "imported" || jan5.HolidayNote == nil || *jan5.HolidayNote != "National holiday" {
		t.Fatalf("jan5 regenerated day = %+v", jan5)
	}
}

func TestWorkingDayCalendarServiceImportHolidaysRejectsDistributorAndInvalidRows(t *testing.T) {
	repo := &fakeWorkingDayCalendarRepository{
		detail: model.WorkingDayCalendar{
			WorkingDayCalendarID: 7,
			CustID:               "P001",
			StartDate:            mustDate(t, "2026-01-01"),
			EndDate:              mustDate(t, "2026-01-07"),
			NumberOfWeeks:        1,
		},
	}
	svc := &workingDayCalendarServiceImpl{repository: repo, now: fixedNow}

	_, err := svc.ImportHolidays(7, entity.WorkingDayCalendarImportHolidayRequest{FileURL: "https://example.com/holidays.csv"}, "D001", "P001", 99)
	if err == nil || !strings.Contains(err.Error(), "principal") {
		t.Fatalf("ImportHolidays() distributor error = %v, want principal error", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("date,notes\n2026-01-08,Outside\n2026-01-05,\n"))
	}))
	defer server.Close()

	resp, err := svc.ImportHolidays(7, entity.WorkingDayCalendarImportHolidayRequest{FileURL: server.URL + "/holidays.csv"}, "P001", "P001", 99)
	if err == nil || !strings.Contains(err.Error(), "validation") {
		t.Fatalf("ImportHolidays() invalid rows error = %v resp=%+v", err, resp)
	}
	if resp.FailedRow != 2 || len(repo.replaceHolidays) != 0 || len(repo.replaceDays) != 0 {
		t.Fatalf("invalid import resp=%+v replaceHolidays=%d replaceDays=%d", resp, len(repo.replaceHolidays), len(repo.replaceDays))
	}
}

func fixedNow() time.Time {
	return time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC)
}

func mustDate(t *testing.T, raw string) time.Time {
	t.Helper()
	date, err := time.ParseInLocation("2006-01-02", raw, time.UTC)
	if err != nil {
		t.Fatalf("parse date %q: %v", raw, err)
	}
	return date
}

func findWorkDay(days []model.MWorkingDay, custID, workDate string) *model.MWorkingDay {
	for i := range days {
		if days[i].CustId == custID && days[i].WorkDate != nil && *days[i].WorkDate == workDate {
			return &days[i]
		}
	}
	return nil
}

func assertStrings(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("strings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("strings = %v, want %v", got, want)
		}
	}
}
