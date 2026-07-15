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

type ReturnReasonDistributorRepository interface {
	FindOneByReturnReasonDistributorIdAndCustId(returnReasonDistributorId int, custId string) (model.ReturnReasonDistributor, error)
	FindOneByReturnReasonDistributorCodeAndCustId(returnReasonDistributorCode string, custId string) (model.ReturnReasonDistributor, error)
	FindAllByCustIdLookupMode(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) (returnReasonDistributor []model.ReturnReasonDistributor, total int, lastPage int, err error)
	FindAllByCustId(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) (returnReasonDistributor []model.ReturnReasonDistributor, total int, lastPage int, err error)
	Store(returnReasonDistributor model.ReturnReasonDistributor) (int, error)
	Update(returnReasonDistributorId int, request entity.UpdateReturnReasonDistributorRequest) error
	Delete(custId string, returnReasonDistributorId int, deletedBy int64) error
}

func NewReturnReasonDistributorRepository(db *sqlx.DB) ReturnReasonDistributorRepository {
	return &returnReasonDistributorRepositoryImpl{db}
}

type returnReasonDistributorRepositoryImpl struct {
	*sqlx.DB
}

func (repository *returnReasonDistributorRepositoryImpl) FindOneByReturnReasonDistributorIdAndCustId(returnReasonDistributorId int, custId string) (model.ReturnReasonDistributor, error) {
	returnReasonDistributor := model.ReturnReasonDistributor{}
	query := `SELECT 
				cust_id, return_reason_id, return_reason_code,
				return_reason_name, return_reason_type, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_return_reason_distributor 
			  WHERE return_reason_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&returnReasonDistributor, query, returnReasonDistributorId, custId)
	if err != nil {
		log.Println("returnReasonDistributorRepository, FindOneByReturnReasonDistributorCodeAndCustId, err:", err.Error())
		return returnReasonDistributor, err
	}

	return returnReasonDistributor, nil
}

func (repository *returnReasonDistributorRepositoryImpl) FindOneByReturnReasonDistributorCodeAndCustId(returnReasonDistributorCode string, custId string) (model.ReturnReasonDistributor, error) {
	returnReasonDistributor := model.ReturnReasonDistributor{}
	query := `SELECT 
				cust_id, return_reason_id, return_reason_code,
				return_reason_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_return_reason_distributor 
			  WHERE return_reason_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&returnReasonDistributor, query, returnReasonDistributorCode, custId)
	if err != nil {
		log.Println("returnReasonDistributorRepository, FindOneByReturnReasonDistributorCodeAndCustId, err:", err.Error())
		return returnReasonDistributor, err
	}

	return returnReasonDistributor, nil
}

func (repository *returnReasonDistributorRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) ([]model.ReturnReasonDistributor, int, int, error) {

	returnReasonDistributors := []model.ReturnReasonDistributor{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` cust_id, 
					return_reason_id, 
					return_reason_code,
					return_reason_name, return_reason_type `
	qWhere := ` WHERE is_active = true AND is_del = false 
				AND cust_id = '` + custId + `' `

	if dataFilter.ReturnReasonDistributorType != "" {
		qWhere += ` AND return_reason_type = '` + dataFilter.ReturnReasonDistributorType + `'`
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (return_reason_code ILIKE '%` + dataFilter.Query + `%' 
					OR return_reason_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := `FROM mst.m_return_reason_distributor`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("returnReasonDistributorRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("returnReasonDistributorRepository, count total, err:", err.Error())
		return returnReasonDistributors, 0, 0, err
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

	// log.Println("returnReasonDistributorRepository, querySelect:", querySelect)
	err = repository.Select(&returnReasonDistributors, querySelect)
	if err != nil {
		log.Println("returnReasonDistributorRepository, FindAllByCustId, err:", err.Error())
		return returnReasonDistributors, total, 1, err
	}

	return returnReasonDistributors, total, 1, nil
}

func (repository *returnReasonDistributorRepositoryImpl) FindAllByCustId(dataFilter entity.ReturnReasonDistributorQueryFilter, custId string) ([]model.ReturnReasonDistributor, int, int, error) {

	returnReasonDistributors := []model.ReturnReasonDistributor{}
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

	if dataFilter.ReturnReasonDistributorType != "" {
		qWhere += ` AND return_reason_type = '` + dataFilter.ReturnReasonDistributorType + `'`
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
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_return_reason_distributor t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_return_reason_distributor t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere

	// log.Println("returnReasonDistributorRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("returnReasonDistributorRepository, count total, err:", err.Error())
		return returnReasonDistributors, 0, 0, err
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

	// log.Println("returnReasonDistributorRepository, querySelect:", querySelect)
	err = repository.Select(&returnReasonDistributors, querySelect)
	if err != nil {
		log.Println("returnReasonDistributorRepository, FindAllByCustId, err:", err.Error())
		return returnReasonDistributors, total, lastPage, err
	}

	return returnReasonDistributors, total, lastPage, nil
}

func (repository *returnReasonDistributorRepositoryImpl) Store(returnReasonDistributor model.ReturnReasonDistributor) (int, error) {
	query :=
		`INSERT INTO mst.m_return_reason_distributor(
			cust_id, return_reason_code, return_reason_name, return_reason_type, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12
		) RETURNING return_reason_id;`
	lastInsertId := returnReasonDistributor.ReturnReasonDistributorId
	err := repository.QueryRow(query,
		returnReasonDistributor.CustId, returnReasonDistributor.ReturnReasonDistributorCode, returnReasonDistributor.ReturnReasonDistributorName, returnReasonDistributor.ReturnReasonDistributorType,
		returnReasonDistributor.IsActive, returnReasonDistributor.CreatedBy, returnReasonDistributor.CreatedAt, returnReasonDistributor.UpdatedBy,
		returnReasonDistributor.UpdatedAt, returnReasonDistributor.IsDel, returnReasonDistributor.DeletedBy, returnReasonDistributor.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("returnReasonDistributorRepository, Store, err:", err.Error())
		return returnReasonDistributor.ReturnReasonDistributorId, err
	}
	return returnReasonDistributor.ReturnReasonDistributorId, nil
}

func (repository *returnReasonDistributorRepositoryImpl) Update(returnReasonDistributorId int, request entity.UpdateReturnReasonDistributorRequest) error {
	var (
		r            model.ReturnReasonDistributorUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("returnReasonDistributorRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_return_reason_distributor
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND return_reason_id = :return_reason_id_old;`

	// log.Println("returnReasonDistributorRepository, Update, query:", query)

	sqlPatch.Args["return_reason_id_old"] = returnReasonDistributorId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("returnReasonDistributorRepository, Update, err:", err.Error())
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

func (repository *returnReasonDistributorRepositoryImpl) Delete(custId string, returnReasonDistributorId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_return_reason_distributor
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND return_reason_id = :return_reason_id;`

	wMap := map[string]interface{}{
		"cust_id":          custId,
		"return_reason_id": returnReasonDistributorId,
		"deleted_by":       deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ReturnReasonDistributorRepository, Delete, err:", err.Error())
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
