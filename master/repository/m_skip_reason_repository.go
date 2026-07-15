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

type SkipReasonRepository interface {
	FindOneBySkipReasonIdAndCustId(skipReasonId int, custId string) (model.SkipReason, error)
	FindOneBySkipReasonCodeAndCustId(skipReasonCode, custId string) (model.SkipReason, error)
	FindAllByCustId(dataFilter entity.SkipReasonQueryFilter, custId string) (consPro []model.SkipReason, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.SkipReasonQueryFilter, custId string) (consPro []model.SkipReason, total int, lastPage int, err error)
	Store(brand model.SkipReason) (int, error)
	Update(skipReasonId int, request entity.UpdateSkipReasonRequest) error
	Delete(custId string, skipReasonId int, deletedBy int64) error
}

func NewSkipReasonRepository(db *sqlx.DB) SkipReasonRepository {
	return &SkipReasonRepositoryImpl{db}
}

type SkipReasonRepositoryImpl struct {
	*sqlx.DB
}

func (repository *SkipReasonRepositoryImpl) FindOneBySkipReasonIdAndCustId(skipReasonId int, custId string) (model.SkipReason, error) {
	skipReason := model.SkipReason{}
	query := `SELECT 
				sr.skip_reason_id, 
				sr.skip_reason_code,
				sr.skip_reason_name, 
				sr.is_active, sr.created_by,
				sr.created_at, sr.updated_by, sr.updated_at,
				sr.is_del, sr.deleted_by, sr.deleted_at
			  FROM mst.m_skip_reason sr
			  WHERE sr.is_del = false 
				AND sr.skip_reason_id = $1 
				AND sr.cust_id = $2`
	err := repository.Get(&skipReason, query, skipReasonId, custId)
	if err != nil {
		log.Println("skipReasonRepository, FindOneBySkipReasonIdAndCustId, err:", err.Error())
		return skipReason, err
	}

	return skipReason, nil
}

func (repository *SkipReasonRepositoryImpl) FindOneBySkipReasonCodeAndCustId(skipReasonCode, custId string) (model.SkipReason, error) {
	skipReason := model.SkipReason{}
	query := `SELECT 
				sr.skip_reason_id, 
				sr.skip_reason_code,
				sr.skip_reason_name, 
				sr.is_active, sr.created_by,
				sr.created_at, sr.updated_by, sr.updated_at,
				sr.is_del, sr.deleted_by, sr.deleted_at
			  FROM mst.m_skip_reason sr
			  WHERE sr.is_del = false 
				AND sr.skip_reason_code = $1 
				AND sr.cust_id = $2`
	err := repository.Get(&skipReason, query, skipReasonCode, custId)
	if err != nil {
		log.Println("skipReasonRepository, FindOneBySkipReasonCodeAndCustId, err:", err.Error())
		return skipReason, err
	}

	return skipReason, nil
}

func (repository *SkipReasonRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.SkipReasonQueryFilter, custId string) ([]model.SkipReason, int, int, error) {

	skipReason := []model.SkipReason{}
	selectCount := ` COUNT(sr.*) AS total `
	selectField := `sr.skip_reason_id,
					sr.skip_reason_name,
					sr.skip_reason_code  `
	qWhere := ` WHERE sr.is_del = false AND sr.is_active = true 
				AND sr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sr.skip_reason_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.SkipReasonId != 0 {
		qWhere += `AND sr.skip_reason_id = ` + strconv.Itoa(dataFilter.SkipReasonId) + ` `
	}

	qFrom := ` FROM mst.m_skip_reason sr `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("skipReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("skipReasonRepository, count total, err:", err.Error())
		return skipReason, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`sr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `sr.skip_reason_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("skipReasonRepository, querySelect:", querySelect)
	err = repository.Select(&skipReason, querySelect)
	if err != nil {
		log.Println("skipReasonRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return skipReason, total, 1, err
	}

	return skipReason, total, 1, nil
}

func (repository *SkipReasonRepositoryImpl) FindAllByCustId(dataFilter entity.SkipReasonQueryFilter, custId string) ([]model.SkipReason, int, int, error) {

	brands := []model.SkipReason{}
	selectCount := ` COUNT(sr.*) AS total `
	selectField := `sr.skip_reason_id,
					sr.skip_reason_code,
					sr.skip_reason_name, 
					sr.is_active, sr.created_by,
					sr.created_at, sr.updated_by, sr.updated_at,
					sr.is_del, sr.deleted_by, sr.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sr.is_del = false 
				AND sr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sr.skip_reason_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND sr.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND sr.is_active = false `
		}
	}

	if dataFilter.SkipReasonId != 0 {
		qWhere += `AND sr.skip_reason_id = ` + strconv.Itoa(dataFilter.SkipReasonId) + ` `
	}

	qFrom := ` FROM mst.m_skip_reason sr
			   LEFT JOIN sys.m_user u ON u.user_id = sr.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("skipReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("skipReasonRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`sr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `sr.skip_reason_id`
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

	// log.Println("skipReasonRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("skipReasonRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *SkipReasonRepositoryImpl) Store(skipReason model.SkipReason) (int, error) {
	query :=
		`INSERT INTO mst.m_skip_reason(
			cust_id, skip_reason_name, skip_reason_code,
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11
		) RETURNING skip_reason_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		skipReason.CustId, skipReason.SkipReasonName, skipReason.SkipReasonCode, skipReason.IsActive,
		skipReason.CreatedBy, skipReason.CreatedAt, skipReason.UpdatedBy, skipReason.UpdatedAt,
		skipReason.IsDel, skipReason.DeletedBy, skipReason.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("skipReasonRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *SkipReasonRepositoryImpl) Update(skipReasonId int, request entity.UpdateSkipReasonRequest) error {
	var (
		r            model.SkipReasonUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_skip_reason
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND skip_reason_id = :skip_reason_id;`

	sqlPatch.Args["skip_reason_id"] = skipReasonId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("skipReasonRepository, Update, err:", err.Error())
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

func (repository *SkipReasonRepositoryImpl) Delete(custId string, skipReasonId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_skip_reason
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND skip_reason_id = :skip_reason_id;`

	wMap := map[string]interface{}{
		"cust_id":        custId,
		"skip_reason_id": skipReasonId,
		"deleted_by":     deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SkipReasonRepository, Delete, err:", err.Error())
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
