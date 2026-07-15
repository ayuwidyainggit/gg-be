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

type ConsProductRepository interface {
	FindOneByCProIdAndCustId(cProId int, custId string) (model.ConsProduct, error)
	FindOneByCProCodeAndCustId(cProCode, custId string) (model.ConsProduct, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ConsProduct, total int, lastPage int, err error)
	FindAllLookupByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ConsProduct, total int, lastPage int, err error)
	Store(consProduct model.ConsProduct) (int, error)
	Update(cProId int, request entity.UpdateConsProductRequest) error
	Delete(custId string, cProId int, deletedBy int64) error
}

func NewConsProductRepository(db *sqlx.DB) ConsProductRepository {
	return &consProductRepositoryImpl{db}
}

type consProductRepositoryImpl struct {
	*sqlx.DB
}

func (repository *consProductRepositoryImpl) FindOneByCProIdAndCustId(cProId int, custId string) (model.ConsProduct, error) {
	consProduct := model.ConsProduct{}
	query := `SELECT 
				cust_id, c_pro_id, c_pro_code,
				c_pro_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_cons_product 
			  WHERE is_del = false 
			  AND c_pro_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&consProduct, query, cProId, custId)
	if err != nil {
		log.Println("consProductRepository, FindOneByCProIdAndCustId, err:", err.Error())
		return consProduct, err
	}

	return consProduct, nil
}

func (repository *consProductRepositoryImpl) FindOneByCProCodeAndCustId(cProCode, custId string) (model.ConsProduct, error) {
	consProduct := model.ConsProduct{}
	query := `SELECT 
				cust_id, c_pro_id, c_pro_code,
				c_pro_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_cons_product 
			  WHERE is_del = false 
			  AND c_pro_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&consProduct, query, cProCode, custId)
	if err != nil {
		log.Println("consProductRepository, FindOneByCProCode, err:", err.Error())
		return consProduct, err
	}

	return consProduct, nil
}

func (repository *consProductRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ConsProduct, int, int, error) {

	consProducts := []model.ConsProduct{}
	selectCount := `COUNT(*) AS total `
	selectField := `
	cp.cust_id, cp.c_pro_id, cp.c_pro_code,
	cp.c_pro_name, cp.is_active, cp.created_by,
	cp.created_at, cp.updated_by, cp.updated_at,
	cp.is_del, cp.deleted_by, cp.deleted_at,
	u.user_fullname AS updated_by_name `
	
	qWhere := `
	WHERE cp.is_del = false 
	AND cp.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += `
		AND (cp.c_pro_code ILIKE '%` + dataFilter.Query + `%' 
		OR cp.c_pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += `
			AND cp.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += `
			AND cp.is_active = false `
		}
	}

	qFrom := `
	FROM mst.m_cons_product cp
	LEFT JOIN sys.m_user u ON u.user_id = cp.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("consProductRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("consProductRepository, count total, err:", err.Error())
		return consProducts, 0, 0, err
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
		querySelect += fmt.Sprintf(`
		ORDER BY %s`, sortBy)
	} else {
		sortBy := `c_pro_id`
		querySelect += fmt.Sprintf(`
		ORDER BY %s DESC`, sortBy)
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

	querySelect += fmt.Sprintf(`
	LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("consProductRepository, querySelect:", querySelect)
	err = repository.Select(&consProducts, querySelect)
	if err != nil {
		log.Println("consProductRepository, FindAllByCustId, err:", err.Error())
		return consProducts, total, lastPage, err
	}

	return consProducts, total, lastPage, nil
}

func (repository *consProductRepositoryImpl) FindAllLookupByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ConsProduct, int, int, error) {

	consProducts := []model.ConsProduct{}
	selectCount := `COUNT(*) AS total `
	selectField := `cp.cust_id, cp.c_pro_id, cp.c_pro_code, cp.c_pro_name, cp.is_active `
	qWhere := `
	WHERE cp.is_del = false 
	AND cp.is_active = true 
	AND cp.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += `
		AND (cp.c_pro_code ILIKE '%` + dataFilter.Query + `%' 
		OR cp.c_pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := `
	FROM mst.m_cons_product cp`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("consProductRepository, count total, err:", err.Error())
		return consProducts, 0, 0, err
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
		querySelect += fmt.Sprintf(`
		ORDER BY %s`, sortBy)
	} else {
		sortBy := `cp.c_pro_id`
		querySelect += fmt.Sprintf(`
		ORDER BY %s DESC`, sortBy)
	}

	// log.Println("consProductRepository, querySelect:", querySelect)
	err = repository.Select(&consProducts, querySelect)
	if err != nil {
		log.Println("consProductRepository, FindAllByCustId, err:", err.Error())
		return consProducts, total, 1, err
	}

	return consProducts, total, 1, nil
}

func (repository *consProductRepositoryImpl) Store(consProduct model.ConsProduct) (int, error) {
	query :=
		`INSERT INTO mst.m_cons_product(
			cust_id, c_pro_code, c_pro_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING c_pro_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		consProduct.CustId, consProduct.CProCode, consProduct.CProName, consProduct.IsActive,
		consProduct.CreatedBy, consProduct.CreatedAt, consProduct.UpdatedBy, consProduct.UpdatedAt,
		consProduct.IsDel, consProduct.DeletedBy, consProduct.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("consProductRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *consProductRepositoryImpl) Update(cProId int, request entity.UpdateConsProductRequest) error {
	var (
		r            model.ConsProductUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("consProductRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_cons_product
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND c_pro_id = :c_pro_id;`

	// log.Println("consProductRepository, Update, query:", query)

	sqlPatch.Args["c_pro_id"] = cProId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("consProductRepository, Update, err:", err.Error())
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

func (repository *consProductRepositoryImpl) Delete(custId string, cProId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_cons_product
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND c_pro_id = :c_pro_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"c_pro_id":   cProId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ConsProductRepository, Delete, err:", err.Error())
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
