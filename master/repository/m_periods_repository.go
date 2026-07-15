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

type MPeriodsRepository interface {
	TrxBegin() (*mPeriodsTransaction, error)
	FindOneByPerYearAndPerIdAndCustId(PerYear int, PerId int, custId string) (model.MPeriods, error)
	FindAllBymPeriodsCodeAndCustId(PerYear int, PerId int, custId string) ([]model.MPeriods, error)
	FindOneByWeekStartMinimum(PerYear int, PerId int, custId string) (model.MWeek, error)
	FindAllByCustId(dataFilter entity.MPeriodsQueryFilter, custId string) (consPro []model.MPeriods, total int, lastPage int, err error)
	Store(mPeriods model.MPeriods) (int, error)
	Delete(custId string, PerYear int, PerId int, closedBy int64, closedByName string) error
	FindAllYearByCustId(dataFilter entity.MPeriodsQueryFilter, custId string) (consPro []model.MPeriods, err error)
	FindOneByYearAndCustId(PerYear int, custId string) (model.MPeriods, error)
	FindOneLastDataCustId(perYear int, custId string) (model.MWeek, error)
	FindConfigByCustId(custId, configId string) (model.MConfig, error)
	FindConfigsByCustId(custId, configIds string) ([]model.MConfig, error)
	FindDistributorByCustId(custId string) ([]model.MCustomer, error)
	HasGeneratedCalendarRows(custId string, perYear int) (bool, error)
}

func NewMPeriodsRepository(db *sqlx.DB) MPeriodsRepository {
	return &mPeriodsRepositoryImpl{db}
}

type mPeriodsRepositoryImpl struct {
	*sqlx.DB
}
type mPeriodsTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewTransactionmPeriods(db *sqlx.DB) (trxObj *mPeriodsTransaction, err error) {
	trx := db.MustBegin()

	return &mPeriodsTransaction{tx: trx, db: db}, nil
}
func (repo *mPeriodsRepositoryImpl) TrxBegin() (*mPeriodsTransaction, error) {
	trxObj, err := NewTransactionmPeriods(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}
func (repo *mPeriodsTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *mPeriodsTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *mPeriodsTransaction) Update(PerYear int, PerId int, request entity.UpdateMPeriodsRequest) error {
	var (
		r            model.MPeriodsUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("mPeriodsRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_periods
			  SET ` + sqlSetFields + `, updated_at = CURRENT_TIMESTAMP
			  WHERE is_closed = false AND cust_id = :cust_id AND per_year = :per_year_old and per_id = :per_id_old ;`

	log.Println("mPeriodsRepository, Update, query:", query)

	sqlPatch.Args["per_year_old"] = PerYear
	sqlPatch.Args["per_id_old"] = PerId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("mPeriodsRepository, Update, err:", err.Error())
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

func (repository *mPeriodsTransaction) FindAllBymPeriodsCodeAndCustId(PerYear int, PerId int, custId string) ([]model.MPeriods, error) {
	mPeriods := []model.MPeriods{}
	query := `SELECT * FROM mst.m_periods 
			  WHERE 
			  	per_year = $1 AND 
			  	per_id >= $2 AND 
				cust_id = $3
				order by per_id`
	err := repository.tx.Select(&mPeriods, query, PerYear, PerId, custId)
	if err != nil {
		log.Println("mPeriodsTransaction, FindAllBymPeriodsCodeAndCustId, err:", err.Error())
		return mPeriods, err
	}

	return mPeriods, nil
}

func (repository *mPeriodsTransaction) FindTotalWeekCountExclByPerIdAndCustId(perYear int, perId int, custId string) (model.MPeriods, error) {
	mPeriod := model.MPeriods{}
	query := `SELECT sum(week_count) AS week_count 
			  FROM mst.m_periods 
			  WHERE per_year = $1 
			  AND per_id != $2 
			  AND cust_id = $3`
	err := repository.tx.Get(&mPeriod, query, perYear, perId, custId)
	if err != nil {
		log.Println("mPeriodsTransaction, FindAllBymPeriodsCodeAndCustId, err:", err.Error())
		return mPeriod, err
	}

	return mPeriod, nil
}

func (repository *mPeriodsTransaction) DeleteMWeekGreaterThanPeriod(perYear int, perId int, groupCustId string) error {
	query := `delete FROM mst.m_week 
			  WHERE 
			  	per_year = :per_year AND 
			  	per_id >= :per_id AND 
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year": perYear,
		"per_id":   perId,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mPeriodsRepository, DeleteMWeekGreaterThanPeriod, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mPeriodsTransaction) DeleteMWorkDayGreaterThanPeriod(perYear int, perId int, groupCustId string) error {
	query := `delete FROM mst.m_work_day 
			  WHERE 
			  	per_year = :per_year AND 
			  	per_id >= :per_id AND 
				working_day_calendar_id IS NULL AND
				cust_id IN (` + groupCustId + `)`
	param := map[string]interface{}{
		"per_year": perYear,
		"per_id":   perId,
	}
	_, err := repository.tx.NamedExec(query, param)
	if err != nil {
		log.Println("mPeriodsRepository, DeleteMWorkDayGreaterThanPeriod, err:", err.Error())
		return err
	}

	return nil
}

func (repository *mPeriodsTransaction) StoreMweek(mWeek model.MWeek) (int, error) {
	query :=
		`INSERT INTO mst.m_week(cust_id, per_year, per_id, week_id, week_start, week_end, is_active, is_closed, closed_at, closed_by)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING per_id;`
	lastInsertId := mWeek.PerId
	err := repository.tx.QueryRow(query,
		mWeek.CustId, mWeek.PerYear, mWeek.PerId, mWeek.WeekId, mWeek.WeekStart, mWeek.WeekEnd, mWeek.IsActive, mWeek.IsClosed, mWeek.ClosedAt, mWeek.ClosedBy).Scan(&lastInsertId)
	if err != nil {
		log.Println("mWeekRepository, Store, err:", err.Error())
		return mWeek.PerId, err
	}
	return mWeek.PerId, nil
}
func (repository *mPeriodsTransaction) StoreMWorkingDay(mWorkingDay model.MWorkingDay) (string, error) {
	query :=
		`INSERT INTO mst.m_work_day(cust_id, per_year, per_id, week_id, work_date, is_active, is_work)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING work_date;`
	lastInsertWorkDate := ""
	err := repository.tx.QueryRow(query,
		mWorkingDay.CustId, mWorkingDay.PerYear, mWorkingDay.PerId, mWorkingDay.WeekId, mWorkingDay.WorkDate, mWorkingDay.IsActive, mWorkingDay.IsWork).Scan(&lastInsertWorkDate)
	if err != nil {
		log.Println("mWorkingDayRepository, Store, err:", err.Error())
		return *mWorkingDay.WorkDate, err
	}
	return *mWorkingDay.WorkDate, nil
}
func (repository *mPeriodsTransaction) FindOneByWeekStartMinimum(PerYear int, PerId int, custId string) (model.MWeek, error) {
	mWeek := model.MWeek{}
	query := `SELECT * FROM mst.m_week WHERE per_year = $1 AND per_id = $2 AND cust_id = $3 order by week_start asc limit 1`
	err := repository.tx.Get(&mWeek, query, PerYear, PerId, custId)
	if err != nil {
		log.Println("mWeekRepository, FindOneByWeekStartMinimum, err:", err.Error())
		return mWeek, err
	}

	return mWeek, nil
}

func (repository *mPeriodsTransaction) StorePeriod(mPeriods model.MPeriods) (err error) {
	query := `INSERT INTO mst.m_periods(
				cust_id, per_year, per_id ,week_count, is_active, 
				updated_at, updated_by, is_closed, closed_at, closed_by)
			  VALUES (
				:cust_id, :per_year, :per_id, :week_count, :is_active, 
				:updated_at, :updated_by, :is_closed, :closed_at, :closed_by
			  );`

	_, err = repository.tx.NamedExec(query, mPeriods)
	if err != nil {
		log.Println("MPeriodsRepository, StoreBulkPeriods, err:", err.Error())
		return err
	}

	return err
}

func (repository *mPeriodsRepositoryImpl) FindOneByPerYearAndPerIdAndCustId(PerYear int, PerId int, custId string) (model.MPeriods, error) {
	mPeriods := model.MPeriods{}
	query := `SELECT * FROM mst.m_periods 
			  WHERE 
			  	per_year = $1 AND 
			  	per_id = $2 AND 
				cust_id = $3`
	err := repository.Get(&mPeriods, query, PerYear, PerId, custId)
	if err != nil {
		log.Println("FindOneByPerYearAndPerIdAndCustId, FindOneByPerYearAndPerIdAndCustId, err:", err.Error())
		return mPeriods, err
	}

	return mPeriods, nil
}

func (repository *mPeriodsRepositoryImpl) FindAllBymPeriodsCodeAndCustId(PerYear int, PerId int, custId string) ([]model.MPeriods, error) {
	mPeriods := []model.MPeriods{}
	query := `SELECT * FROM mst.m_periods 
			  WHERE 
			  	per_year = $1 AND 
			  	per_id >= $2 AND 
				cust_id = $3
				order by per_id`
	err := repository.Select(&mPeriods, query, PerYear, PerId, custId)
	if err != nil {
		log.Println("remarkPromoRepository, FindAllBymPeriodsCodeAndCustId, err:", err.Error())
		return mPeriods, err
	}

	return mPeriods, nil
}

func (repository *mPeriodsRepositoryImpl) FindOneByWeekStartMinimum(PerYear int, PerId int, custId string) (model.MWeek, error) {
	mWeek := model.MWeek{}
	query := `SELECT * FROM mst.m_week WHERE per_year = $1 AND per_id = $2 AND cust_id = $3 order by week_start asc limit 1`
	err := repository.Get(&mWeek, query, PerYear, PerId, custId)
	if err != nil {
		log.Println("mWeekRepository, FindOneByWeekStartMinimum, err:", err.Error())
		return mWeek, err
	}

	return mWeek, nil
}

func (repository *mPeriodsRepositoryImpl) FindAllByCustId(dataFilter entity.MPeriodsQueryFilter, custId string) ([]model.MPeriods, int, int, error) {

	mPeriods := []model.MPeriods{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.*,
					u.user_fullname AS updated_by_name,
					cl.user_fullname AS closed_by_name	`
	qWhere := ` WHERE p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.per_id = ` + dataFilter.Query + `' 
					OR p.week_count = '%` + dataFilter.Query + `' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	if dataFilter.PerYear != "" {
		qWhere += ` AND p.per_year = '` + dataFilter.PerYear + `'`
	}

	qFrom := ` FROM mst.m_periods p
			   LEFT JOIN sys.m_user u ON u.user_id = p.updated_by 
			   LEFT JOIN sys.m_user cl ON cl.user_id = p.closed_by `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("mPeriodsRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("mPeriodsRepository, count total, err:", err.Error())
		return mPeriods, 0, 0, err
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
		sortBy := `p.per_id`
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

	// log.Println("mPeriodsRepository, querySelect:", querySelect)
	err = repository.Select(&mPeriods, querySelect)
	if err != nil {
		log.Println("mPeriodsRepository, FindAllByCustId, err:", err.Error())
		return mPeriods, total, lastPage, err
	}

	return mPeriods, total, lastPage, nil
}

func (repository *mPeriodsRepositoryImpl) Store(mPeriods model.MPeriods) (int, error) {
	query :=
		`INSERT INTO mst.m_periods(
			cust_id, per_year, per_id ,week_count, is_active, 
			updated_at, updated_by, is_closed, closed_at, closed_by)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING per_id;`
	lastInsertId := mPeriods.PerId
	err := repository.QueryRow(query,
		mPeriods.CustId, mPeriods.PerYear, mPeriods.PerId, mPeriods.WeekCount, mPeriods.IsActive,
		mPeriods.UpdatedAt, mPeriods.UpdatedBy, mPeriods.IsClosed, mPeriods.ClosedAt, mPeriods.ClosedBy).
		Scan(&lastInsertId)
	if err != nil {
		log.Println("mPeriodsRepository, Store, err:", err.Error())
		return mPeriods.PerId, err
	}
	return mPeriods.PerId, nil
}

func (repository *mPeriodsRepositoryImpl) Delete(custId string, PerYear int, PerId int, closedBy int64, closedByName string) error {
	var nRows int64
	query := `UPDATE mst.m_periods SET is_closed = true, closed_at = CURRENT_TIMESTAMP, closed_by = :closed_by, closed_by_name =:closed_by_name  
			WHERE is_closed = false AND cust_id = :cust_id AND per_year = :per_year and per_id = :per_id;`

	wMap := map[string]interface{}{
		"cust_id":        custId,
		"per_year":       PerYear,
		"per_id":         PerId,
		"closed_by":      closedBy,
		"closed_by_name": closedByName,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("MPeriodsRepository, Delete, err:", err.Error())
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

func (repository *mPeriodsRepositoryImpl) FindAllYearByCustId(dataFilter entity.MPeriodsQueryFilter, custId string) ([]model.MPeriods, error) {

	mPeriods := []model.MPeriods{}

	selectField := ` p.per_year`
	qWhere := ` WHERE p.cust_id = '` + custId + `' 
				GROUP BY p.per_year
				ORDER BY p.per_year ASC`

	qFrom := ` FROM mst.m_periods p `

	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	err := repository.Select(&mPeriods, querySelect)
	if err != nil {
		log.Println("mPeriodsRepository, FindAllByCustId, err:", err.Error())
		return mPeriods, err
	}

	return mPeriods, nil
}

func (repository *mPeriodsRepositoryImpl) FindOneByYearAndCustId(PerYear int, custId string) (model.MPeriods, error) {
	mPeriods := model.MPeriods{}
	query := `SELECT * FROM mst.m_periods 
			  WHERE 
			  	per_year = $1 AND 
				cust_id = $2`
	err := repository.Get(&mPeriods, query, PerYear, custId)
	if err != nil {
		log.Println("m_periods_repository, FindOneByYearAndCustId, err:", err.Error())
		return mPeriods, err
	}

	return mPeriods, nil
}

func (repository *mPeriodsRepositoryImpl) FindOneLastDataCustId(perYear int, custId string) (model.MWeek, error) {
	mWeeks := model.MWeek{}
	query := `SELECT * FROM mst.m_week 
			  WHERE cust_id = $1 AND per_year = $2
			  ORDER BY week_id DESC
			  LIMIT 1`
	err := repository.Get(&mWeeks, query, custId, perYear)
	if err != nil {
		log.Println("m_periods_repository, FindOneLastDataCustId, err:", err.Error())
		return mWeeks, err
	}

	return mWeeks, nil
}

func (repository *mPeriodsRepositoryImpl) FindConfigByCustId(custId, configId string) (model.MConfig, error) {
	mConfig := model.MConfig{}
	query := `SELECT * 
			  FROM sys.m_config 
			  WHERE 
			  	config_id = $1 
				AND cust_id = $2`
	err := repository.Get(&mConfig, query, configId, custId)
	if err != nil {
		log.Println("m_periods_repository, FindConfigByCustId, err:", err.Error())
		return mConfig, err
	}

	return mConfig, nil
}

func (repository *mPeriodsRepositoryImpl) FindConfigsByCustId(custId, configIds string) ([]model.MConfig, error) {
	mConfigs := []model.MConfig{}
	query := `SELECT config_id, config_value, data_type, config_desc, module, created_date, config_group 
			  FROM sys.m_config 
			  WHERE 
			  	config_id IN (` + configIds + `) 
				AND cust_id = $1`
	err := repository.Select(&mConfigs, query, custId)
	if err != nil {
		log.Println("m_periods_repository, FindConfigsByCustId, err:", err.Error())
		return mConfigs, err
	}

	return mConfigs, nil
}

func (repository *mPeriodsRepositoryImpl) FindDistributorByCustId(custId string) ([]model.MCustomer, error) {
	customers := []model.MCustomer{}
	querySelect := `SELECT cust_id, cust_name from smc.m_customer WHERE parent_cust_id ='` + custId + `' AND cust_id != parent_cust_id `
	err := repository.Select(&customers, querySelect)
	if err != nil {
		log.Println("m_periods_repository, FindDistributorByCustId, err:", err.Error())
		return customers, err
	}

	return customers, nil
}

func (repository *mPeriodsRepositoryImpl) HasGeneratedCalendarRows(custId string, perYear int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1
		FROM mst.m_week
		WHERE cust_id = $1
			AND per_year = $2
			AND is_closed = false
			AND working_day_calendar_id IS NOT NULL
	)`
	err := repository.Get(&exists, query, custId, perYear)
	return exists, err
}
