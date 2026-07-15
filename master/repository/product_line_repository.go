package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/str"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductLineRepository interface {
	FindOneByPLIdAndCustId(pLId int, custId string) (model.ProductLine, error)
	FindOneByPLCodeAndCustId(pLCode, custId string) (model.ProductLine, error)
	FindAllByCustId(dataFilter entity.ProductLineQueryFilter, custId string) (pl []model.ProductLine, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.ProductLineQueryFilter, custId string) (pl []model.ProductLine, total int, lastPage int, err error)
	Store(productLine model.ProductLine) (int, error)
	Update(pLId int, request entity.UpdateProductLineRequest) error
	Delete(custId string, pLId int, deletedBy int64) error
}

func NewProductLineRepository(db *sqlx.DB) ProductLineRepository {
	return &productLineRepositoryImpl{db}
}

type productLineRepositoryImpl struct {
	*sqlx.DB
}

func (repository *productLineRepositoryImpl) FindOneByPLIdAndCustId(pLId int, custId string) (model.ProductLine, error) {
	productLine := model.ProductLine{}
	query := `SELECT 
				cust_id, pl_id, pl_code, 
				eff_call, min_item, 
				pl_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_product_line 
			  WHERE is_del = false 
			  AND pl_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&productLine, query, pLId, custId)
	if err != nil {
		log.Println("productLineRepository, FindOneByPLIdAndCustId, err:", err.Error())
		return productLine, err
	}

	return productLine, nil
}

func (repository *productLineRepositoryImpl) FindOneByPLCodeAndCustId(pLCode, custId string) (model.ProductLine, error) {
	productLine := model.ProductLine{}
	query := `SELECT 
				cust_id, pl_id, pl_code, eff_call, min_item, 
				pl_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_product_line 
			  WHERE is_del = false 
			  AND pl_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&productLine, query, pLCode, custId)
	if err != nil {
		log.Println("productLineRepository, FindOneByPLCode, err:", err.Error())
		return productLine, err
	}

	return productLine, nil
}

func (repository *productLineRepositoryImpl) FindAllByCustId(dataFilter entity.ProductLineQueryFilter, custId string) ([]model.ProductLine, int, int, error) {

	productLines := []model.ProductLine{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pl.cust_id, pl.pl_id, pl.pl_code, pl.eff_call, pl.min_item, 
	pl.pl_name, pl.is_active, pl.created_by,
	pl.created_at, pl.updated_by, pl.updated_at,
	pl.is_del, pl.deleted_by, pl.deleted_at,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE pl.is_del = false 
				AND pl.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (pl.pl_code ILIKE '%` + dataFilter.Query + `%' 
					OR pl.pl_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.PlIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PlIds, ",")
		qWhere += ` AND pl.pl_id IN (` + intArrStr + `)`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND pl.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND pl.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_product_line pl
	LEFT JOIN sys.m_user u ON u.user_id = pl.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("productLineRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productLineRepository, count total, err:", err.Error())
		return productLines, 0, 0, err
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
		sortBy := `pl_code`
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

	// log.Println("productLineRepository, querySelect:", querySelect)
	err = repository.Select(&productLines, querySelect)
	if err != nil {
		log.Println("productLineRepository, FindAllByCustId, err:", err.Error())
		return productLines, total, lastPage, err
	}

	return productLines, total, lastPage, nil
}

func (repository *productLineRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.ProductLineQueryFilter, custId string) ([]model.ProductLine, int, int, error) {

	productLines := []model.ProductLine{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pl.cust_id, pl.pl_id, pl.pl_code, pl.eff_call, pl.min_item, 
					pl.pl_name, pl.is_active, pl.created_by,
					pl.created_at, pl.updated_by, pl.updated_at,
					pl.is_del, pl.deleted_by, pl.deleted_at,
					u.user_fullname AS updated_by_name`
	qWhere := ` WHERE pl.cust_id = '` + custId + `' AND pl.is_del = false AND pl.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (pl_code ILIKE '%` + dataFilter.Query + `%' 
					OR pl_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.PlIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PlIds, ",")
		qWhere += ` AND pl.pl_id IN (` + intArrStr + `)`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND pl.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND pl.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_product_line pl
	LEFT JOIN sys.m_user u ON u.user_id = pl.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("productLineRepository, FindAllByCustIdLookupMode, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productLineRepository, FindAllByCustIdLookupMode, count total, err:", err.Error())
		return productLines, 0, 0, err
	}

	sortBy := `pl.pl_code` // default sort by
	if sortBy != "" {
		querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)
	}

	// log.Println("productLineRepository, FindAllByCustIdLookupMode, querySelect:", querySelect)
	err = repository.Select(&productLines, querySelect)
	if err != nil {
		log.Println("productLineRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return productLines, total, 1, err
	}

	return productLines, total, 1, nil
}

func (repository *productLineRepositoryImpl) Store(productLine model.ProductLine) (int, error) {
	query :=
		`INSERT INTO mst.m_product_line(
			cust_id, pl_code, pl_name, 
			eff_call, min_item, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8,	$9, $10, 
			$11, $12, $13 
		) RETURNING pl_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		productLine.CustId, productLine.PlCode, productLine.PlName,
		productLine.EffCall, productLine.MinItem, productLine.IsActive,
		productLine.CreatedBy, productLine.CreatedAt, productLine.UpdatedBy, productLine.UpdatedAt,
		productLine.IsDel, productLine.DeletedBy, productLine.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("productLineRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *productLineRepositoryImpl) Update(pLId int, request entity.UpdateProductLineRequest) error {
	var (
		r            model.ProductLineUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("productLineRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_product_line
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND pl_id = :pl_id;`

	// log.Println("productLineRepository, Update, query:", query)

	sqlPatch.Args["pl_id"] = pLId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("productLineRepository, Update, err:", err.Error())
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

func (repository *productLineRepositoryImpl) Delete(custId string, pLId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product_line
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND pl_id = :pl_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"pl_id":      pLId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ProductLineRepository, Delete, err:", err.Error())
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
