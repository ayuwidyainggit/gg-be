package repository

import (
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const workingDayCalendarInsertBatchSize = 500

type WorkingDayCalendarRepository interface {
	FindLatestCalendarByCustID(custID string) (model.WorkingDayCalendar, error)
	FindLatestWeekIDByCustIDs(custIDs []string, perYear int) (int, error)
	StoreCalendarWithDetails(calendar model.WorkingDayCalendar, holidays []model.WorkingDayCalendarHoliday, weeks []model.MWeek, days []model.MWorkingDay) (int64, error)
	FindAll(filter entity.WorkingDayCalendarQueryFilter, ownerCustID string) ([]model.WorkingDayCalendar, int, int, error)
	FindByID(id int64, ownerCustID string) (model.WorkingDayCalendar, error)
	FindImportedHolidays(id int64) ([]model.WorkingDayCalendarHoliday, error)
	FindCalendarDays(id int64, ownerCustID string, dateFrom, dateTo time.Time) ([]model.WorkingDayCalendarDay, error)
	ReplaceImportedHolidaysAndWorkDays(calendarID int64, custID string, userID int64, holidays []model.WorkingDayCalendarHoliday, days []model.MWorkingDay) error
}

type workingDayCalendarRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkingDayCalendarRepository(db *sqlx.DB) WorkingDayCalendarRepository {
	return &workingDayCalendarRepositoryImpl{DB: db}
}

func (r *workingDayCalendarRepositoryImpl) FindLatestCalendarByCustID(custID string) (model.WorkingDayCalendar, error) {
	row := model.WorkingDayCalendar{}
	query := `
		SELECT working_day_calendar_id, cust_id, title, start_date, number_of_weeks, end_date,
			default_holidays, is_closed, created_at, created_by, updated_at, updated_by,
			closed_at, closed_by, closed_by_name
		FROM mst.working_day_calendar
		WHERE cust_id = $1
			AND is_closed = false
		ORDER BY end_date DESC, working_day_calendar_id DESC
		LIMIT 1`
	err := r.Get(&row, query, custID)
	return row, err
}

func (r *workingDayCalendarRepositoryImpl) FindLatestWeekIDByCustIDs(custIDs []string, perYear int) (int, error) {
	if len(custIDs) == 0 {
		return 0, nil
	}
	query, args, err := sqlx.In(`
		SELECT COALESCE(MAX(week_id), 0)
		FROM mst.m_week
		WHERE cust_id IN (?)
			AND per_year = ?`, custIDs, perYear)
	if err != nil {
		return 0, err
	}
	query = r.Rebind(query)
	var latestWeekID int
	err = r.Get(&latestWeekID, query, args...)
	return latestWeekID, err
}

func (r *workingDayCalendarRepositoryImpl) StoreCalendarWithDetails(calendar model.WorkingDayCalendar, holidays []model.WorkingDayCalendarHoliday, weeks []model.MWeek, days []model.MWorkingDay) (int64, error) {
	tx, err := r.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	calendarID, err := storeWorkingDayCalendar(tx, calendar)
	if err != nil {
		return 0, err
	}
	if err := storeWorkingDayCalendarHolidays(tx, calendarID, holidays); err != nil {
		return 0, err
	}
	if err := storeWorkingDayCalendarWeeks(tx, calendarID, weeks); err != nil {
		return 0, err
	}
	if err := storeWorkingDayCalendarWorkDays(tx, calendarID, days); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return calendarID, nil
}

func (r *workingDayCalendarRepositoryImpl) FindAll(filter entity.WorkingDayCalendarQueryFilter, ownerCustID string) ([]model.WorkingDayCalendar, int, int, error) {
	rows := []model.WorkingDayCalendar{}
	args := []interface{}{ownerCustID}
	where := ` WHERE cust_id = $1 AND is_closed = false`
	if query := strings.TrimSpace(filter.Query); query != "" {
		args = append(args, "%"+query+"%")
		where += fmt.Sprintf(` AND title ILIKE $%d`, len(args))
	}

	countQuery := `SELECT COUNT(*) FROM mst.working_day_calendar` + where
	var total int
	if err := r.Get(&total, countQuery, args...); err != nil {
		return rows, 0, 0, err
	}

	selectQuery := `
		SELECT working_day_calendar_id, cust_id, title, start_date, number_of_weeks, end_date,
			default_holidays, is_closed, created_at, created_by, updated_at, updated_by,
			closed_at, closed_by, closed_by_name
		FROM mst.working_day_calendar` + where + buildWorkingDayCalendarOrderBy(filter.Sort)

	limit := normalizeWorkingDayCalendarLimit(filter.Limit)
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	selectQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	if err := r.Select(&rows, selectQuery, args...); err != nil {
		return rows, total, sql_helper.CalculateLastPage(total, limit), err
	}
	return rows, total, sql_helper.CalculateLastPage(total, limit), nil
}

func (r *workingDayCalendarRepositoryImpl) FindByID(id int64, ownerCustID string) (model.WorkingDayCalendar, error) {
	row := model.WorkingDayCalendar{}
	query := `
		SELECT working_day_calendar_id, cust_id, title, start_date, number_of_weeks, end_date,
			default_holidays, is_closed, created_at, created_by, updated_at, updated_by,
			closed_at, closed_by, closed_by_name
		FROM mst.working_day_calendar
		WHERE working_day_calendar_id = $1
			AND cust_id = $2
			AND is_closed = false`
	err := r.Get(&row, query, id, ownerCustID)
	return row, err
}

func (r *workingDayCalendarRepositoryImpl) FindImportedHolidays(id int64) ([]model.WorkingDayCalendarHoliday, error) {
	rows := []model.WorkingDayCalendarHoliday{}
	query := `
		SELECT working_day_calendar_holiday_id, working_day_calendar_id, distributor_cust_id, holiday_date, notes, created_at, created_by
		FROM mst.working_day_calendar_holiday
		WHERE working_day_calendar_id = $1
			AND distributor_cust_id IS NULL
		ORDER BY holiday_date ASC`
	err := r.Select(&rows, query, id)
	return rows, err
}

func (r *workingDayCalendarRepositoryImpl) FindCalendarDays(id int64, ownerCustID string, dateFrom, dateTo time.Time) ([]model.WorkingDayCalendarDay, error) {
	rows := []model.WorkingDayCalendarDay{}
	query := `
		SELECT
			wd.work_date,
			wd.week_id,
			COALESCE(mw.calendar_week_no, 0) AS calendar_week_no,
			COALESCE(wd.is_work, true) AS is_work,
			wd.holiday_source,
			wd.holiday_note,
			COALESCE(wd.holiday_source IN ('default', 'default_imported'), false) AS is_default_holiday,
			COALESCE(wd.holiday_source IN ('imported', 'default_imported'), false) AS is_imported_holiday
		FROM mst.m_work_day wd
		INNER JOIN mst.working_day_calendar c
			ON c.working_day_calendar_id = wd.working_day_calendar_id
			AND c.cust_id = $2
			AND c.is_closed = false
		LEFT JOIN mst.m_week mw
			ON mw.working_day_calendar_id = wd.working_day_calendar_id
			AND mw.cust_id = wd.cust_id
			AND mw.per_year = wd.per_year
			AND mw.week_id = wd.week_id
		WHERE wd.working_day_calendar_id = $1
			AND wd.cust_id = $2
			AND wd.work_date::date >= $3::date
			AND wd.work_date::date <= $4::date
		ORDER BY wd.work_date ASC`
	err := r.Select(&rows, query, id, ownerCustID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))
	return rows, err
}

func (r *workingDayCalendarRepositoryImpl) ReplaceImportedHolidaysAndWorkDays(calendarID int64, custID string, userID int64, holidays []model.WorkingDayCalendarHoliday, days []model.MWorkingDay) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		DELETE FROM mst.working_day_calendar_holiday
		WHERE working_day_calendar_id = $1
			AND distributor_cust_id IS NULL`, calendarID); err != nil {
		return err
	}

	if err := storeWorkingDayCalendarHolidays(tx, calendarID, holidays); err != nil {
		return err
	}
	if err := updateWorkingDayCalendarWorkDays(tx, calendarID, custID, days); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE mst.working_day_calendar
		SET updated_at = CURRENT_TIMESTAMP,
			updated_by = $2
		WHERE working_day_calendar_id = $1
			AND cust_id = $3`, calendarID, userID, custID); err != nil {
		return err
	}

	return tx.Commit()
}

func storeWorkingDayCalendar(tx *sqlx.Tx, calendar model.WorkingDayCalendar) (int64, error) {
	query := `
		INSERT INTO mst.working_day_calendar (
			cust_id, title, start_date, number_of_weeks, end_date, default_holidays, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING working_day_calendar_id`
	var calendarID int64
	err := tx.Get(&calendarID, query,
		calendar.CustID,
		calendar.Title,
		calendar.StartDate.Format("2006-01-02"),
		calendar.NumberOfWeeks,
		calendar.EndDate.Format("2006-01-02"),
		pq.Array(calendar.DefaultHolidays),
		calendar.CreatedBy,
	)
	return calendarID, err
}

func storeWorkingDayCalendarHolidays(tx *sqlx.Tx, calendarID int64, holidays []model.WorkingDayCalendarHoliday) error {
	if len(holidays) == 0 {
		return nil
	}
	for start := 0; start < len(holidays); start += workingDayCalendarInsertBatchSize {
		end := minWorkingDayCalendarBatchEnd(start, len(holidays))
		values := make([]string, 0, end-start)
		args := make([]interface{}, 0, (end-start)*5)
		argPos := 1
		for _, holiday := range holidays[start:end] {
			values = append(values, workingDayCalendarValuePlaceholders(argPos, 5))
			args = append(args,
				calendarID,
				holiday.DistributorCustID,
				holiday.HolidayDate.Format("2006-01-02"),
				holiday.Notes,
				holiday.CreatedBy,
			)
			argPos += 5
		}
		query := `
			INSERT INTO mst.working_day_calendar_holiday (
				working_day_calendar_id, distributor_cust_id, holiday_date, notes, created_by
			) VALUES ` + strings.Join(values, ", ")
		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func storeWorkingDayCalendarWeeks(tx *sqlx.Tx, calendarID int64, weeks []model.MWeek) error {
	if len(weeks) == 0 {
		return nil
	}
	for start := 0; start < len(weeks); start += workingDayCalendarInsertBatchSize {
		end := minWorkingDayCalendarBatchEnd(start, len(weeks))
		values := make([]string, 0, end-start)
		args := make([]interface{}, 0, (end-start)*12)
		argPos := 1
		for _, week := range weeks[start:end] {
			values = append(values, workingDayCalendarValuePlaceholders(argPos, 12))
			args = append(args,
				week.CustId,
				week.PerYear,
				week.PerId,
				week.WeekId,
				week.WeekStart,
				week.WeekEnd,
				week.IsActive,
				week.IsClosed,
				week.ClosedAt,
				week.ClosedBy,
				calendarID,
				week.CalendarWeekNo,
			)
			argPos += 12
		}
		query := `
			INSERT INTO mst.m_week (
			cust_id, per_year, per_id, week_id, week_start, week_end, is_active,
			is_closed, closed_at, closed_by, working_day_calendar_id, calendar_week_no
		) VALUES ` + strings.Join(values, ", ")
		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func storeWorkingDayCalendarWorkDays(tx *sqlx.Tx, calendarID int64, days []model.MWorkingDay) error {
	if len(days) == 0 {
		return nil
	}
	for start := 0; start < len(days); start += workingDayCalendarInsertBatchSize {
		end := minWorkingDayCalendarBatchEnd(start, len(days))
		values := make([]string, 0, end-start)
		args := make([]interface{}, 0, (end-start)*10)
		argPos := 1
		for _, day := range days[start:end] {
			values = append(values, workingDayCalendarValuePlaceholders(argPos, 10))
			args = append(args,
				day.CustId,
				day.PerYear,
				day.PerId,
				day.WeekId,
				day.WorkDate,
				day.IsActive,
				day.IsWork,
				calendarID,
				day.HolidaySource,
				day.HolidayNote,
			)
			argPos += 10
		}
		query := `
			INSERT INTO mst.m_work_day (
			cust_id, per_year, per_id, week_id, work_date, is_active, is_work,
			working_day_calendar_id, holiday_source, holiday_note
		) VALUES ` + strings.Join(values, ", ")
		if _, err := tx.Exec(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func updateWorkingDayCalendarWorkDays(tx *sqlx.Tx, calendarID int64, custID string, days []model.MWorkingDay) error {
	for _, day := range days {
		if day.WorkDate == nil || day.IsWork == nil {
			continue
		}
		if _, err := tx.Exec(`
			UPDATE mst.m_work_day
			SET is_work = $1,
				holiday_source = $2,
				holiday_note = $3
			WHERE working_day_calendar_id = $4
				AND cust_id = $5
				AND work_date::date = $6::date`,
			day.IsWork,
			day.HolidaySource,
			day.HolidayNote,
			calendarID,
			custID,
			day.WorkDate,
		); err != nil {
			return err
		}
	}
	return nil
}

func buildWorkingDayCalendarOrderBy(sort string) string {
	columnMap := map[string]string{
		"title":           "title",
		"start_date":      "start_date",
		"end_date":        "end_date",
		"number_of_weeks": "number_of_weeks",
		"created_at":      "created_at",
	}
	if sort == "" {
		return " ORDER BY start_date DESC, working_day_calendar_id DESC"
	}
	parts := strings.Split(sort, ",")
	orders := make([]string, 0, len(parts))
	for _, part := range parts {
		segs := strings.SplitN(strings.TrimSpace(part), ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "DESC"
		if len(segs) > 1 {
			upper := strings.ToUpper(strings.TrimSpace(segs[1]))
			if upper == "ASC" || upper == "DESC" {
				dir = upper
			}
		}
		if column, ok := columnMap[key]; ok {
			orders = append(orders, column+" "+dir)
		}
	}
	if len(orders) == 0 {
		return " ORDER BY start_date DESC, working_day_calendar_id DESC"
	}
	return " ORDER BY " + strings.Join(orders, ", ")
}

func normalizeWorkingDayCalendarLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func minWorkingDayCalendarBatchEnd(start, total int) int {
	end := start + workingDayCalendarInsertBatchSize
	if end > total {
		return total
	}
	return end
}

func workingDayCalendarValuePlaceholders(start, count int) string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		placeholders = append(placeholders, "$"+strconv.Itoa(start+i))
	}
	return "(" + strings.Join(placeholders, ", ") + ")"
}
