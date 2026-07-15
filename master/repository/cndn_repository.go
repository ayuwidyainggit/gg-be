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

type CndnRepository interface {
	FindOneByCndnIdAndCustId(cndnId int, custId string) (model.Cndn, error)
	FindOneByCndnCodeAndCustId(cndnCode string, custId string) (model.Cndn, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (cndn []model.Cndn, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (cndn []model.Cndn, total int, lastPage int, err error)
	Store(cndn model.Cndn) (int, error)
	Update(cndnId int, request entity.UpdateCndnRequest) error
	Delete(custId string, cndnId int, deletedBy int64) error
}

func NewCndnRepository(db *sqlx.DB) CndnRepository {
	return &cndnRepositoryImpl{db}
}

type cndnRepositoryImpl struct {
	*sqlx.DB
}

func (repository *cndnRepositoryImpl) FindOneByCndnIdAndCustId(cndnId int, custId string) (model.Cndn, error) {
	cndn := model.Cndn{}
	query := `SELECT 
				cust_id, cndn_id, cndn_code,
				cndn_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_cndn 
			  WHERE cndn_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&cndn, query, cndnId, custId)
	if err != nil {
		log.Println("cndnRepository, FindOneByCndnCodeAndCustId, err:", err.Error())
		return cndn, err
	}

	return cndn, nil
}

func (repository *cndnRepositoryImpl) FindOneByCndnCodeAndCustId(cndnCode string, custId string) (model.Cndn, error) {
	cndn := model.Cndn{}
	query := `SELECT 
				cust_id, cndn_id, cndn_code,
				cndn_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_cndn 
			  WHERE cndn_code = $1 
			  AND cust_id = $2 and is_del = false`
	err := repository.Get(&cndn, query, cndnCode, custId)
	if err != nil {
		log.Println("cndnRepository, FindOneByCndnCodeAndCustId, err:", err.Error())
		return cndn, err
	}

	return cndn, nil
}

func (repository *cndnRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Cndn, int, int, error) {

	cndns := []model.Cndn{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.cndn_id, a.cndn_code,
					a.cndn_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `

	qWhere := ` WHERE a.is_del = false 
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.cndn_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.cndn_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_cndn a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by  `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("cndnRepository, count total, err:", err.Error())
		return cndns, 0, 0, err
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
		sortBy := `a.cndn_id`
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

	err = repository.Select(&cndns, querySelect)
	if err != nil {
		log.Println("cndnRepository, FindAllByCustId, err:", err.Error())
		return cndns, total, lastPage, err
	}

	return cndns, total, lastPage, nil
}

func (repository *cndnRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Cndn, int, int, error) {

	cndns := []model.Cndn{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.cndn_id, a.cndn_code,
					a.cndn_name  `

	qWhere := ` WHERE a.is_del = false AND a.is_active = true
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.cndn_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.cndn_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_cndn a   `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("cndnRepository, count total, err:", err.Error())
		return cndns, 0, 0, err
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
		sortBy := `a.cndn_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&cndns, querySelect)
	if err != nil {
		log.Println("cndnRepository, FindAllByCustId, err:", err.Error())
		return cndns, total, 1, err
	}

	return cndns, total, 1, nil
}

func (repository *cndnRepositoryImpl) Store(cndn model.Cndn) (int, error) {
	query :=
		`INSERT INTO mst.m_cndn(
			cust_id, cndn_code, cndn_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING cndn_id;`
	lastInsertId := cndn.CndnId
	err := repository.QueryRow(query,
		cndn.CustId, cndn.CndnCode, cndn.CndnName,
		cndn.IsActive, cndn.CreatedBy, cndn.CreatedAt, cndn.UpdatedBy,
		cndn.UpdatedAt, cndn.IsDel, cndn.DeletedBy, cndn.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("cndnRepository, Store, err:", err.Error())
		return cndn.CndnId, err
	}
	return cndn.CndnId, nil
}

func (repository *cndnRepositoryImpl) Update(cndnId int, request entity.UpdateCndnRequest) error {
	var (
		r            model.CndnUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("cndnRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_cndn
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND cndn_id = :cndn_id_old;`

	// log.Println("cndnRepository, Update, query:", query)

	sqlPatch.Args["cndn_id_old"] = cndnId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("cndnRepository, Update, err:", err.Error())
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

func (repository *cndnRepositoryImpl) Delete(custId string, cndnId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_cndn
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND cndn_id = :cndn_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"cndn_id":    cndnId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("CndnRepository, Delete, err:", err.Error())
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
