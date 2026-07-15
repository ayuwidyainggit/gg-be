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

type SalesTypeRepository interface {
	FindOneBySalesTypeIdAndCustId(salesTypeId int, custId string) (model.SalesType, error)
	FindOneBySalesTypeCodeAndCustId(salesTypeCode string, custId string) (model.SalesType, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (salesType []model.SalesType, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (salesType []model.SalesType, total int, lastPage int, err error)
	Store(salesType model.SalesType) (int, error)
	Update(salesTypeId int, request entity.UpdateSalesTypeRequest) error
	Delete(custId string, salesTypeId int, deletedBy int64) error
}

func NewSalesTypeRepository(db *sqlx.DB) SalesTypeRepository {
	return &salesTypeRepositoryImpl{db}
}

type salesTypeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *salesTypeRepositoryImpl) FindOneBySalesTypeIdAndCustId(salesTypeId int, custId string) (model.SalesType, error) {
	salesType := model.SalesType{}
	query := `SELECT 
				cust_id, sales_type_id, sales_type_code,
				sales_type_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sales_type 
			  WHERE sales_type_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&salesType, query, salesTypeId, custId)
	if err != nil {
		log.Println("salesTypeRepository, FindOneBySalesTypeCodeAndCustId, err:", err.Error())
		return salesType, err
	}

	return salesType, nil
}

func (repository *salesTypeRepositoryImpl) FindOneBySalesTypeCodeAndCustId(salesTypeCode string, custId string) (model.SalesType, error) {
	salesType := model.SalesType{}
	query := `SELECT 
				cust_id, sales_type_id, sales_type_code,
				sales_type_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sales_type 
			  WHERE sales_type_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&salesType, query, salesTypeCode, custId)
	if err != nil {
		log.Println("salesTypeRepository, FindOneBySalesTypeCodeAndCustId, err:", err.Error())
		return salesType, err
	}

	return salesType, nil
}

func (repository *salesTypeRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SalesType, int, int, error) {

	salesTypes := []model.SalesType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sales_type_id, a.sales_type_code,
					a.sales_type_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sales_type_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sales_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_sales_type a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesTypeRepository, count total, err:", err.Error())
		return salesTypes, 0, 0, err
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
		sortBy := `a.sales_type_id`
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

	// log.Println("salesTypeRepository, querySelect:", querySelect)
	err = repository.Select(&salesTypes, querySelect)
	if err != nil {
		log.Println("salesTypeRepository, FindAllByCustId, err:", err.Error())
		return salesTypes, total, lastPage, err
	}

	return salesTypes, total, lastPage, nil
}

func (repository *salesTypeRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SalesType, int, int, error) {

	salesTypes := []model.SalesType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sales_type_id, a.sales_type_code,
					a.sales_type_name, a.is_active `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sales_type_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sales_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_sales_type a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesTypeRepository, count total, err:", err.Error())
		return salesTypes, 0, 0, err
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
		sortBy := `a.sales_type_id`
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

	// log.Println("salesTypeRepository, querySelect:", querySelect)
	err = repository.Select(&salesTypes, querySelect)
	if err != nil {
		log.Println("salesTypeRepository, FindAllByCustId, err:", err.Error())
		return salesTypes, total, lastPage, err
	}

	return salesTypes, total, lastPage, nil
}

func (repository *salesTypeRepositoryImpl) Store(salesType model.SalesType) (int, error) {
	query :=
		`INSERT INTO mst.m_sales_type(
			cust_id, sales_type_code, sales_type_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING sales_type_id;`
	lastInsertId := salesType.SalesTypeId
	err := repository.QueryRow(query,
		salesType.CustId, salesType.SalesTypeCode, salesType.SalesTypeName,
		salesType.IsActive, salesType.CreatedBy, salesType.CreatedAt, salesType.UpdatedBy,
		salesType.UpdatedAt, salesType.IsDel, salesType.DeletedBy, salesType.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("salesTypeRepository, Store, err:", err.Error())
		return salesType.SalesTypeId, err
	}
	return salesType.SalesTypeId, nil
}

func (repository *salesTypeRepositoryImpl) Update(salesTypeId int, request entity.UpdateSalesTypeRequest) error {
	var (
		r            model.SalesTypeUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("salesTypeRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sales_type
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sales_type_id = :sales_type_id_old;`

	log.Println("salesTypeRepository, Update, query:", query)

	sqlPatch.Args["sales_type_id_old"] = salesTypeId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("salesTypeRepository, Update, err:", err.Error())
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

func (repository *salesTypeRepositoryImpl) Delete(custId string, salesTypeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sales_type
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND sales_type_id = :sales_type_id;`

	wMap := map[string]interface{}{
		"cust_id":      custId,
		"sales_type_id": salesTypeId,
		"deleted_by":   deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesTypeRepository, Delete, err:", err.Error())
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
