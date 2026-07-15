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

type ReturnReasonRepository interface {
	FindOneByReturnReasonIdAndCustId(returnReasonId int, custId string) (model.ReturnReason, error)
	FindOneByReturnReasonCodeAndCustId(returnReasonCode string, custId string) (model.ReturnReason, error)
	FindAllByCustIdLookupMode(dataFilter entity.ReturnReasonQueryFilter, custId string) (returnReason []model.ReturnReason, total int, lastPage int, err error)
	FindAllByCustId(dataFilter entity.ReturnReasonQueryFilter, custId string) (returnReason []model.ReturnReason, total int, lastPage int, err error)
	Store(returnReason model.ReturnReason) (int, error)
	Update(returnReasonId int, request entity.UpdateReturnReasonRequest) error
	Delete(custId string, returnReasonId int, deletedBy int64) error
}

func NewReturnReasonRepository(db *sqlx.DB) ReturnReasonRepository {
	return &returnReasonRepositoryImpl{db}
}

type returnReasonRepositoryImpl struct {
	*sqlx.DB
}

func (repository *returnReasonRepositoryImpl) FindOneByReturnReasonIdAndCustId(returnReasonId int, custId string) (model.ReturnReason, error) {
	returnReason := model.ReturnReason{}
	query := `SELECT 
				cust_id, return_reason_id, return_reason_code,
				return_reason_name, return_reason_type, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_return_reason 
			  WHERE return_reason_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&returnReason, query, returnReasonId, custId)
	if err != nil {
		log.Println("returnReasonRepository, FindOneByReturnReasonCodeAndCustId, err:", err.Error())
		return returnReason, err
	}

	return returnReason, nil
}

func (repository *returnReasonRepositoryImpl) FindOneByReturnReasonCodeAndCustId(returnReasonCode string, custId string) (model.ReturnReason, error) {
	returnReason := model.ReturnReason{}
	query := `SELECT 
				cust_id, return_reason_id, return_reason_code,
				return_reason_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_return_reason 
			  WHERE return_reason_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&returnReason, query, returnReasonCode, custId)
	if err != nil {
		log.Println("returnReasonRepository, FindOneByReturnReasonCodeAndCustId, err:", err.Error())
		return returnReason, err
	}

	return returnReason, nil
}

func (repository *returnReasonRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.ReturnReasonQueryFilter, custId string) ([]model.ReturnReason, int, int, error) {

	returnReasons := []model.ReturnReason{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` cust_id, 
					return_reason_id, 
					return_reason_code,
					return_reason_name, return_reason_type `
	qWhere := ` WHERE is_active = true AND is_del = false 
				AND cust_id = '` + custId + `' `

	if dataFilter.ReturnReasonType != "" {
		qWhere += ` AND return_reason_type = '` + dataFilter.ReturnReasonType + `'`
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (return_reason_code ILIKE '%` + dataFilter.Query + `%' 
					OR return_reason_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := `FROM mst.m_return_reason`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("returnReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("returnReasonRepository, count total, err:", err.Error())
		return returnReasons, 0, 0, err
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
		sortBy := `return_reason_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("returnReasonRepository, querySelect:", querySelect)
	err = repository.Select(&returnReasons, querySelect)
	if err != nil {
		log.Println("returnReasonRepository, FindAllByCustId, err:", err.Error())
		return returnReasons, total, 1, err
	}

	return returnReasons, total, 1, nil
}

func (repository *returnReasonRepositoryImpl) FindAllByCustId(dataFilter entity.ReturnReasonQueryFilter, custId string) ([]model.ReturnReason, int, int, error) {

	returnReasons := []model.ReturnReason{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` t.cust_id, t.return_reason_id, t.return_reason_code,
					t.return_reason_name, t.is_active, t.created_by,
					t.created_at, t.updated_by, t.updated_at,
					t.is_del, t.deleted_by, t.deleted_at, t.return_reason_type, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE t.is_del = false 
				AND t.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (return_reason_code ILIKE '%` + dataFilter.Query + `%' 
					OR return_reason_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.ReturnReasonType != "" {
		qWhere += ` AND return_reason_type = '` + dataFilter.ReturnReasonType + `'`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND t.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND t.is_active = false `
		}
	}
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_return_reason t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_return_reason t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere

	// log.Println("returnReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("returnReasonRepository, count total, err:", err.Error())
		return returnReasons, 0, 0, err
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
		sortBy := `return_reason_id`
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

	// log.Println("returnReasonRepository, querySelect:", querySelect)
	err = repository.Select(&returnReasons, querySelect)
	if err != nil {
		log.Println("returnReasonRepository, FindAllByCustId, err:", err.Error())
		return returnReasons, total, lastPage, err
	}

	return returnReasons, total, lastPage, nil
}

func (repository *returnReasonRepositoryImpl) Store(returnReason model.ReturnReason) (int, error) {
	query :=
		`INSERT INTO mst.m_return_reason(
			cust_id, return_reason_code, return_reason_name, return_reason_type, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12
		) RETURNING return_reason_id;`
	lastInsertId := returnReason.ReturnReasonId
	err := repository.QueryRow(query,
		returnReason.CustId, returnReason.ReturnReasonCode, returnReason.ReturnReasonName, returnReason.ReturnReasonType,
		returnReason.IsActive, returnReason.CreatedBy, returnReason.CreatedAt, returnReason.UpdatedBy,
		returnReason.UpdatedAt, returnReason.IsDel, returnReason.DeletedBy, returnReason.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("returnReasonRepository, Store, err:", err.Error())
		return returnReason.ReturnReasonId, err
	}
	return returnReason.ReturnReasonId, nil
}

func (repository *returnReasonRepositoryImpl) Update(returnReasonId int, request entity.UpdateReturnReasonRequest) error {
	var (
		r            model.ReturnReasonUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("returnReasonRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_return_reason
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND return_reason_id = :return_reason_id_old;`

	// log.Println("returnReasonRepository, Update, query:", query)

	sqlPatch.Args["return_reason_id_old"] = returnReasonId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("returnReasonRepository, Update, err:", err.Error())
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

func (repository *returnReasonRepositoryImpl) Delete(custId string, returnReasonId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_return_reason
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND return_reason_id = :return_reason_id;`

	wMap := map[string]interface{}{
		"cust_id":          custId,
		"return_reason_id": returnReasonId,
		"deleted_by":       deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ReturnReasonRepository, Delete, err:", err.Error())
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
