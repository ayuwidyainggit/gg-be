package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MWorkingDayRepository interface {
	FindOneActive(custId string, generatedOnly bool) (model.MWorkingDayActive, error)
	FindOneByPerYearAndPerIdAndWorkingDayIdAndCustId(PerYear int, PerId int, WeekId int, WorkDate string, custId string, generatedOnly bool) (model.MWorkingDay, error)
	// FindOneBymWorkingDayCodeAndCustId(mWorkingDayCode string, custId string) (model.mWorkingDay, error)
	FindAllByCustId(dataFilter entity.MWorkingDayQueryFilter, custId string, generatedOnly bool) (consPro []model.MWorkingDay, total int, lastPage int, err error)
	HasGeneratedCalendarRows(custId string, perYear string) (bool, error)
	HasLegacyRows(custId string, perYear string) (bool, error)
	Store(mWorkingDay model.MWorkingDay) (int, error)
	Update(PerYear int, PerId int, WeekId int, WorkDate string, request entity.UpdateMWorkingDayRequest) error
	Delete(custId string, PerYear int, PerId int, WeekId int, WorkDate string, closedBy int64, closedByName string) error
}

func NewMWorkingDayRepository(db *sqlx.DB) MWorkingDayRepository {
	return &mWorkingDayRepositoryImpl{db}
}

type mWorkingDayRepositoryImpl struct {
	*sqlx.DB
}

func (repository *mWorkingDayRepositoryImpl) FindOneActive(custId string, generatedOnly bool) (model.MWorkingDayActive, error) {
	mWorkingDay := model.MWorkingDayActive{}
	query := `	SELECT per_year, per_id, week_id, work_date, is_active, is_closed
				FROM mst.m_work_day 
				WHERE cust_id = $1 AND is_active = true`
	if generatedOnly {
		query += ` AND working_day_calendar_id IS NOT NULL`
	}
	query += `;`
	err := repository.Get(&mWorkingDay, query, custId)
	if err != nil {
		log.Println("mWorkingDayRepository, FindOneActive, err:", err.Error())
		return mWorkingDay, err
	}

	return mWorkingDay, nil
}

func (repository *mWorkingDayRepositoryImpl) FindOneByPerYearAndPerIdAndWorkingDayIdAndCustId(PerYear int, PerId int, WeekId int, WorkDate string, custId string, generatedOnly bool) (model.MWorkingDay, error) {
	mWorkingDay := model.MWorkingDay{}
	query := `	SELECT * 
				FROM mst.m_work_day 
				WHERE per_year = $1 
				AND per_id = $2 
				AND week_id = $3 
				AND cust_id = $4 
				AND work_date = $5 `
	if generatedOnly {
		query += ` AND working_day_calendar_id IS NOT NULL`
	}
	// AND is_closed = FALSE`
	err := repository.Get(&mWorkingDay, query, PerYear, PerId, WeekId, custId, WorkDate)
	if err != nil {
		log.Println("mWorkingDayRepository, FindOneByPerYearAndPerIdAndWeekIdAndCustId, err:", err.Error())
		return mWorkingDay, err
	}

	return mWorkingDay, nil
}

func (repository *mWorkingDayRepositoryImpl) FindAllByCustId(dataFilter entity.MWorkingDayQueryFilter, custId string, generatedOnly bool) ([]model.MWorkingDay, int, int, error) {

	mWorkingDays := []model.MWorkingDay{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` *`
	qWhere := ` WHERE cust_id = '` + custId + `' `
	if generatedOnly {
		qWhere += ` AND working_day_calendar_id IS NOT NULL `
	}

	// if dataFilter.Query != "" {
	// qWhere += ` AND (week_id = %` + dataFilter.Query + `% )`
	// }

	if dataFilter.IsActive != nil {
		// fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND is_active = false `
		}
	}

	if dataFilter.PerYear != "" {
		qWhere += ` AND per_year = '` + dataFilter.PerYear + `'`
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_work_day ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_work_day ` + qWhere

	// log.Println("mWorkingDaysRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("mWorkingDaysRepository, count total, err:", err.Error())
		return mWorkingDays, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `per_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("mWorkingDaysRepository, querySelect:", querySelect)
	err = repository.Select(&mWorkingDays, querySelect)
	if err != nil {
		log.Println("mWorkingDaysRepository, FindAllByCustId, err:", err.Error())
		return mWorkingDays, total, lastPage, err
	}

	return mWorkingDays, total, lastPage, nil
}

func (repository *mWorkingDayRepositoryImpl) HasGeneratedCalendarRows(custId string, perYear string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM mst.m_work_day
		WHERE cust_id = $1
			AND working_day_calendar_id IS NOT NULL`
	args := []interface{}{custId}
	if perYear != "" {
		query += ` AND per_year = $2::int`
		args = append(args, perYear)
	}
	query += `)`
	err := repository.Get(&exists, query, args...)
	return exists, err
}

func (repository *mWorkingDayRepositoryImpl) HasLegacyRows(custId string, perYear string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM mst.m_work_day
		WHERE cust_id = $1
			AND working_day_calendar_id IS NULL`
	args := []interface{}{custId}
	if perYear != "" {
		query += ` AND per_year = $2::int`
		args = append(args, perYear)
	}
	query += `)`
	err := repository.Get(&exists, query, args...)
	return exists, err
}

func (repository *mWorkingDayRepositoryImpl) Store(mWorkingDay model.MWorkingDay) (int, error) {
	query :=
		`INSERT INTO mst.m_work_day(cust_id, per_year, per_id, week_id, work_date, work_day_id)
		VALUES ( 
			$1, $2, $3, $4, $5, $6
		) RETURNING per_id;`
	lastInsertId := mWorkingDay.PerId
	err := repository.QueryRow(query,
		mWorkingDay.CustId, mWorkingDay.PerYear, mWorkingDay.PerId, mWorkingDay.WeekId, mWorkingDay.WorkDate, mWorkingDay.WorkDayId).Scan(&lastInsertId)
	if err != nil {
		log.Println("mWorkingDayRepository, Store, err:", err.Error())
		return mWorkingDay.PerId, err
	}
	return mWorkingDay.PerId, nil
}

func (repository *mWorkingDayRepositoryImpl) Update(PerYear int, PerId int, WeekId int, WorkDate string, request entity.UpdateMWorkingDayRequest) error {
	var (
		r            model.MWorkingDayUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("mWorkingDayRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_work_day
			  SET ` + sqlSetFields + `
			  WHERE is_closed = false AND working_day_calendar_id IS NULL AND cust_id = :cust_id AND per_year = :per_year_old and per_id = :per_id_old 
			  and week_id = :week_id_old and work_date = :work_date_old ;`

	log.Println("mWorkingDayRepository, Update, query:", query)

	sqlPatch.Args["per_year_old"] = PerYear
	sqlPatch.Args["per_id_old"] = PerId
	sqlPatch.Args["week_id_old"] = WeekId
	sqlPatch.Args["work_date_old"] = WorkDate
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("mWorkingDayRepository, Update, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *mWorkingDayRepositoryImpl) Delete(custId string, PerYear int, PerId int, WeekId int, WorkDate string, closedBy int64, closedByName string) error {
	var nRows int64
	query := `UPDATE mst.m_work_day SET is_closed = true, closed_at = CURRENT_TIMESTAMP, closed_by = :closed_by, closed_by_name =:closed_by_name  
			WHERE is_closed = false AND working_day_calendar_id IS NULL AND cust_id = :cust_id AND per_year = :per_year and per_id = :per_id and week_id = :week_id and work_date = :work_date;`

	wMap := map[string]interface{}{
		"cust_id":        custId,
		"per_year":       PerYear,
		"per_id":         PerId,
		"week_id":        WeekId,
		"work_date":      WorkDate,
		"closed_by":      closedBy,
		"closed_by_name": closedByName,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("mWorkingDayRepository, Delete, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
