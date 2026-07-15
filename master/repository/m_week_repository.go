package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/structs"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MWeekRepository interface {
	TrxBegin() (*mWeekTransaction, error)
	FindOneByPerYearAndPerIdAndWeekIdAndCustId(perYear int, perId int, weekId int, custId string, generatedOnly bool) (model.MWeek, error)
	FindAllByGreaterThanWeekIdAndCustId(perYear int, weekId int, custId string) ([]model.MWeek, error)
	// FindOneBymWeekCodeAndCustId(mWeekCode string, custId string) (model.mWeek, error)
	FindAllByCustId(dataFilter entity.MWeekQueryFilter, custId string, generatedOnly bool) (consPro []model.MWeek, total int, lastPage int, err error)
	HasGeneratedCalendarRows(custId string, perYear string) (bool, error)
	HasLegacyRows(custId string, perYear string) (bool, error)
	Store(mWeek model.MWeek) (int, error)
	Update(PerYear int, PerId int, WeekId int, request entity.UpdateMWeekRequest) error
	Delete(custId string, PerYear int, PerId int, WeekId int, closedBy int64, closedByName string) error
	FindDistributorByCustId(custId string) ([]model.MCustomer, error)
}

func NewMWeekRepository(db *sqlx.DB) MWeekRepository {
	return &mWeekRepositoryImpl{db}
}

type mWeekRepositoryImpl struct {
	*sqlx.DB
}

type mWeekTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewTransactionmWeek(db *sqlx.DB) (trxObj *mWeekTransaction, err error) {
	trx := db.MustBegin()

	return &mWeekTransaction{tx: trx, db: db}, nil
}

func (repo *mWeekRepositoryImpl) TrxBegin() (*mWeekTransaction, error) {
	trxObj, err := NewTransactionmWeek(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}
func (repo *mWeekTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *mWeekTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *mWeekRepositoryImpl) FindOneByPerYearAndPerIdAndWeekIdAndCustId(PerYear int, PerId int, WeekId int, custId string, generatedOnly bool) (model.MWeek, error) {
	mWeek := model.MWeek{}
	query := `SELECT * FROM mst.m_week WHERE per_year = $1 AND per_id = $2 AND week_id = $3 and cust_id = $4 and is_closed = FALSE`
	if generatedOnly {
		query += ` AND working_day_calendar_id IS NOT NULL`
	}
	err := repository.Get(&mWeek, query, PerYear, PerId, WeekId, custId)
	if err != nil {
		log.Println("mWeekRepository, FindOneByPerYearAndPerIdAndWeekIdAndCustId, err:", err.Error())
		return mWeek, err
	}

	return mWeek, nil
}

func (repository *mWeekRepositoryImpl) FindAllByGreaterThanWeekIdAndCustId(perYear int, weekId int, custId string) ([]model.MWeek, error) {
	mWeeks := []model.MWeek{}
	query := `SELECT * 
			  FROM mst.m_week 
			  WHERE 
			  	per_year = $1 AND 
				week_id >= $2 AND 
				cust_id = $3
			  ORDER BY week_id ASC`
	err := repository.Select(&mWeeks, query, perYear, weekId, custId)
	if err != nil {
		log.Println("mWeekRepository, FindAllByGreaterThanWeekIdAndCustId, err:", err.Error())
		return mWeeks, err
	}

	return mWeeks, nil
}

func (repository *mWeekRepositoryImpl) FindAllByCustId(dataFilter entity.MWeekQueryFilter, custId string, generatedOnly bool) ([]model.MWeek, int, int, error) {

	mWeeks := []model.MWeek{}
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}
	queryCount, countArgs, querySelect, selectArgs := buildMWeekListQuery(dataFilter, custId, generatedOnly)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return mWeeks, 0, 0, err
	}
	queryCount = repository.Rebind(queryCount)

	var total int
	err = repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Println("mWeeksRepository, count total, err:", err.Error())
		return mWeeks, 0, 0, err
	}

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return mWeeks, total, 0, err
	}
	querySelect = repository.Rebind(querySelect)

	// log.Println("mWeeksRepository, querySelect:", querySelect)
	err = repository.Select(&mWeeks, querySelect, selectArgs...)
	if err != nil {
		log.Println("mWeeksRepository, FindAllByCustId, err:", err.Error())
		return mWeeks, total, calculateMWeekLastPage(total, dataFilter.Limit), err
	}

	return mWeeks, total, calculateMWeekLastPage(total, dataFilter.Limit), nil
}

func buildMWeekListQuery(dataFilter entity.MWeekQueryFilter, custId string, generatedOnly bool) (string, []interface{}, string, []interface{}) {
	selectCount := ` COUNT(*) AS total `
	selectField := ` *`
	queryCustID := custId
	if generatedOnly && dataFilter.ParentCustId != "" && dataFilter.ParentCustId != custId {
		queryCustID = dataFilter.ParentCustId
	}
	qWhere := ` WHERE is_closed = false and cust_id = '` + queryCustID + `' `
	args := make([]interface{}, 0)
	if generatedOnly {
		qWhere += ` AND working_day_calendar_id IS NOT NULL `
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (week_id ILIKE '%` + dataFilter.Query + `%' )`
	}

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

	if len(dataFilter.WorkingDayCalendarID) > 0 && generatedOnly {
		qWhere += ` AND working_day_calendar_id IN (?)`
		args = append(args, dataFilter.WorkingDayCalendarID)
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_week ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_week ` + qWhere

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

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	return queryCount, args, querySelect, args
}

func calculateMWeekLastPage(total int, limit int) int {
	if limit == 0 {
		limit = 10
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func (repository *mWeekRepositoryImpl) HasGeneratedCalendarRows(custId string, perYear string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM mst.m_week
		WHERE cust_id = $1
			AND is_closed = false
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

func (repository *mWeekRepositoryImpl) HasLegacyRows(custId string, perYear string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM mst.m_week
		WHERE cust_id = $1
			AND is_closed = false
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

func (repository *mWeekRepositoryImpl) Store(mWeek model.MWeek) (int, error) {
	query :=
		`INSERT INTO mst.m_week(cust_id, per_year, per_id, week_id, week_start, week_end, is_active, is_closed, closed_at, closed_by)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING per_id;`
	lastInsertId := mWeek.PerId
	err := repository.QueryRow(query,
		mWeek.CustId, mWeek.PerYear, mWeek.PerId, mWeek.WeekId, mWeek.WeekStart, mWeek.WeekEnd, mWeek.IsActive, mWeek.IsClosed, mWeek.ClosedAt, mWeek.ClosedBy).Scan(&lastInsertId)
	if err != nil {
		log.Println("mWeekRepository, Store, err:", err.Error())
		return mWeek.PerId, err
	}
	return mWeek.PerId, nil
}

func (repository *mWeekRepositoryImpl) Update(PerYear int, PerId int, WeekId int, request entity.UpdateMWeekRequest) error {
	var (
		r            model.MWeekUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("mWeekRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_week
			  SET ` + sqlSetFields + `
			  WHERE is_closed = false AND working_day_calendar_id IS NULL AND cust_id = :cust_id AND per_year = :per_year_old and per_id = :per_id_old and week_id = :week_id_old ;`

	log.Println("mWeekRepository, Update, query:", query)

	sqlPatch.Args["per_year_old"] = PerYear
	sqlPatch.Args["per_id_old"] = PerId
	sqlPatch.Args["week_id_old"] = WeekId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("mWeekRepository, Update, err:", err.Error())
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

func (repository *mWeekRepositoryImpl) Delete(custId string, PerYear int, PerId int, WeekId int, closedBy int64, closedByName string) error {
	var nRows int64
	query := `UPDATE mst.m_week SET is_closed = true, closed_at = CURRENT_TIMESTAMP, closed_by = :closed_by, closed_by_name =:closed_by_name  
			WHERE is_closed = false AND working_day_calendar_id IS NULL AND cust_id = :cust_id AND per_year = :per_year and per_id = :per_id and week_id = :week_id;`

	wMap := map[string]interface{}{
		"cust_id":        custId,
		"per_year":       PerYear,
		"per_id":         PerId,
		"week_id":        WeekId,
		"closed_by":      closedBy,
		"closed_by_name": closedByName,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("mWeekRepository, Delete, err:", err.Error())
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

func (repository *mWeekRepositoryImpl) FindDistributorByCustId(custId string) ([]model.MCustomer, error) {
	customers := []model.MCustomer{}
	querySelect := `SELECT cust_id, cust_name from smc.m_customer WHERE parent_cust_id ='` + custId + `' AND cust_id != parent_cust_id `
	err := repository.Select(&customers, querySelect)
	if err != nil {
		log.Println("m_periods_repository, FindDistributorByCustId, err:", err.Error())
		return customers, err
	}

	return customers, nil
}

func (repository *mWeekTransaction) UpdateTrx(perYear, perId, weekId int, weekStart, weekEnd, groupCustId string) error {
	query := `UPDATE mst.m_week
			  SET 
				week_start = :week_start,
				week_end = :week_end
			  WHERE 
			  	is_closed = false AND 
				working_day_calendar_id IS NULL AND
				per_year = :per_year AND 
				per_id = :per_id AND 
				week_id = :week_id AND
				cust_id IN (` + groupCustId + `);`

	args := map[string]interface{}{
		"per_year":   perYear,
		"per_id":     perId,
		"week_id":    weekId,
		"week_start": weekStart,
		"week_end":   weekEnd,
	}

	result, err := repository.tx.NamedExec(query, args)
	if err != nil {
		log.Println("mWeekTransaction, UpdateTrx, err:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		log.Println("mWeekTransaction, UpdateTrx, err:", err.Error())
	}

	if nRows == 0 {
		log.Println("mWeekTransaction, UpdateTrx, nRows:", nRows)
	}

	return nil
}

func (repository *mWeekTransaction) DeleteMWeekGreaterThanWeekId(perYear int, weekId int, groupCustId string) error {
	query := `delete FROM mst.m_week 
			  WHERE 
			  	per_year = :per_year AND 
			  	week_id >= :week_id AND 
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year": perYear,
		"week_id":  weekId,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, DeleteMWeekGreaterThanWeekId, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) DeleteMWorkDayGreaterThanWeekId(perYear int, weekId int, groupCustId string) error {
	query := `delete FROM mst.m_work_day 
			  WHERE 
			  	per_year = :per_year AND 
				week_id >= :week_id AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year": perYear,
		"week_id":  weekId,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, DeleteMWorkDayGreaterThanWeekId, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) DeleteMWorkDayByWeekIdExceptFirstDate(perYear, perId, weekId int, exceptDateStr, groupCustId string) error {
	query := `DELETE FROM mst.m_work_day 
			  WHERE 
			  	per_year = :per_year AND 
				per_id = :per_id AND
				week_id = :week_id AND
				work_date != :work_date AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year":  perYear,
		"per_id":    perId,
		"week_id":   weekId,
		"work_date": exceptDateStr,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, DeleteMWorkDayByWeekIdExceptFirstDate, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) UpdateMWorkDayByWeekIdGreaterThanDate(perYear, perId, weekId int, date, groupCustId string, isNeedUpdatePerId bool) error {
	newWeekId := weekId + 1
	newPerId := perId
	if isNeedUpdatePerId {
		newPerId += 1
	}
	query := `UPDATE mst.m_work_day 
			  SET 
				per_id = :next_per_id,
				week_id = :next_week_id
			  WHERE 
			  	per_year = :per_year AND 
				per_id = :per_id AND
				week_id = :week_id AND
				work_date > :work_date AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year":     perYear,
		"per_id":       perId,
		"next_per_id":  newPerId,
		"next_week_id": newWeekId,
		"week_id":      weekId,
		"work_date":    date,
	}
	result, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, UpdateMWorkDayByWeekIdGreaterThanDate, err:", err.Error())
		return err
	}
	log.Println("UpdateMWorkDayByWeekIdGreaterThanDate, result:", structs.StructToJson(result))
	return nil
}

func (repository *mWeekTransaction) UpdateMWorkDayByWeekIdLowerThanDate(perYear, perId, weekId int, date, groupCustId string, isNeedUpdatePerId bool) error {
	newWeekId := weekId - 1
	newPerId := perId
	if isNeedUpdatePerId {
		newPerId -= 1
	}
	query := `UPDATE mst.m_work_day 
			  SET 
				per_id = :next_per_id,
				week_id = :next_week_id
			  WHERE 
			  	per_year = :per_year AND 
				per_id = :per_id AND
				week_id = :week_id AND
				work_date < :work_date AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year":     perYear,
		"per_id":       perId,
		"next_per_id":  newPerId,
		"next_week_id": newWeekId,
		"week_id":      weekId,
		"work_date":    date,
	}
	result, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, UpdateMWorkDayByWeekIdGreaterThanDate, err:", err.Error())
		return err
	}
	log.Println("UpdateMWorkDayByWeekIdGreaterThanDate, result:", structs.StructToJson(result))
	return nil
}

func (repository *mWeekTransaction) DeleteMWorkDayByGreaterThanDate(perYear int, date, groupCustId string) error {
	query := `DELETE FROM mst.m_work_day 
			  WHERE 
			  	per_year = :per_year AND 
				work_date > :work_date AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year":  perYear,
		"work_date": date,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mWeekRepository, DeleteMWorkDayByGreaterThanDate, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) UpdateMWorkDayByWeekIdAndGreaterThanDate(perYear, perId, weekId int, endDate, addDate, groupCustId string) error {
	query := `UPDATE mst.m_work_day 
			  SET 
			  	per_id = :per_id,
				week_id = :week_id
			  WHERE 
			  	per_year = :per_year AND
				work_date > :end_date AND
				work_date <= :add_date AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	args := map[string]interface{}{
		"per_year": perYear,
		"per_id":   perId,
		"week_id":  weekId,
		"end_date": endDate,
		"add_date": addDate,
	}
	_, err := repository.tx.NamedExec(query, args)
	if err != nil {
		log.Println("mWeekRepository, UpdateMWorkDayByWeekIdAndGreaterThanDate, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) UpdateMWeekStartEndDate(perYear, weekId int, weekStart, weekEnd, groupCustId string) error {
	query := `UPDATE mst.m_week 
			  SET 
			  	week_start = :week_start,
				week_end = :week_end
			  WHERE 
			  	per_year = :per_year AND
				week_id = :week_id AND
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	args := map[string]interface{}{
		"per_year":   perYear,
		"week_start": weekStart,
		"week_end":   weekEnd,
		"week_id":    weekId,
	}
	_, err := repository.tx.NamedExec(query, args)
	if err != nil {
		log.Println("mWeekRepository, UpdateMWeekStartEndDate, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mWeekTransaction) StoreMWorkDay(mWorkingDay model.MWorkingDay) error {
	query :=
		`INSERT INTO mst.m_work_day(cust_id, per_year, per_id, week_id, work_date)
		VALUES (:cust_id, :per_year, :per_id, :week_id, :work_date);`
	result, err := repository.tx.NamedExec(query, mWorkingDay)
	if err != nil {
		log.Println("mWeekTransaction, StoreMWorkDay, err:", err.Error())
		return err
	}
	log.Println("result:", result)
	return nil
}
