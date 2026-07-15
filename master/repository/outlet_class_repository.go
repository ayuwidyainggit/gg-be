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

type OutletClassRepository interface {
	FindOneByOutletClassIdAndCustId(outletClassId int, custId string) (model.OutletClass, error)
	FindOneByOutletClassCodeAndCustId(outletClassCode string, custId string) (model.OutletClass, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (outletClass []model.OutletClass, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (outletClass []model.OutletClass, total int, lastPage int, err error)
	Store(outletClass model.OutletClass) (int, error)
	Update(outletClassId int, request entity.UpdateOutletClassRequest) error
	Delete(custId string, outletClassId int, deletedBy int64) error
}

func NewOutletClassRepository(db *sqlx.DB) OutletClassRepository {
	return &outletClassRepositoryImpl{db}
}

type outletClassRepositoryImpl struct {
	*sqlx.DB
}

func (repository *outletClassRepositoryImpl) FindOneByOutletClassIdAndCustId(outletClassId int, custId string) (model.OutletClass, error) {
	outletClass := model.OutletClass{}
	query := `SELECT 
				cust_id, ot_class_id, ot_class_code,
				ot_class_name, ot_class_limit, is_active, 
				created_by, created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_class 
			  WHERE ot_class_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletClass, query, outletClassId, custId)
	if err != nil {
		log.Println("outletClassRepository, FindOneByOutletClassCodeAndCustId, err:", err.Error())
		return outletClass, err
	}

	return outletClass, nil
}

func (repository *outletClassRepositoryImpl) FindOneByOutletClassCodeAndCustId(outletClassCode string, custId string) (model.OutletClass, error) {
	outletClass := model.OutletClass{}
	query := `SELECT 
				cust_id, ot_class_id, ot_class_code,
				ot_class_name, ot_class_limit, is_active, 
				created_by, created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_class 
			  WHERE ot_class_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletClass, query, outletClassCode, custId)
	if err != nil {
		log.Println("outletClassRepository, FindOneByOutletClassCodeAndCustId, err:", err.Error())
		return outletClass, err
	}

	return outletClass, nil
}

func (repository *outletClassRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletClass, int, int, error) {

	outletClasss := []model.OutletClass{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_class_id, a.ot_class_code, a.ot_class_name, a.ot_class_limit, a.is_active, 
					 a.created_by, a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_class_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_class_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_outlet_class a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("outletClassRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletClassRepository, count total, err:", err.Error())
		return outletClasss, 0, 0, err
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
		sortBy := `a.ot_class_id`
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

	// log.Println("outletClassRepository, querySelect:", querySelect)
	err = repository.Select(&outletClasss, querySelect)
	if err != nil {
		log.Println("outletClassRepository, FindAllByCustId, err:", err.Error())
		return outletClasss, total, lastPage, err
	}

	return outletClasss, total, lastPage, nil
}

func (repository *outletClassRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletClass, int, int, error) {

	outletClasss := []model.OutletClass{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_class_id, a.ot_class_code, a.ot_class_name, a.ot_class_limit `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_class_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_class_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_outlet_class a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletClassRepository, count total, err:", err.Error())
		return outletClasss, 0, 0, err
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
		sortBy := `a.ot_class_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}
	
	err = repository.Select(&outletClasss, querySelect)
	if err != nil {
		log.Println("outletClassRepository, FindAllByCustId, err:", err.Error())
		return outletClasss, total, 1, err
	}

	return outletClasss, total, 1, nil
}

func (repository *outletClassRepositoryImpl) Store(outletClass model.OutletClass) (int, error) {
	query :=
		`INSERT INTO mst.m_outlet_class(
			cust_id, ot_class_code, ot_class_name, ot_class_limit, 
			is_active, created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12
		) RETURNING ot_class_id;`
	lastInsertId := outletClass.OtClassId
	err := repository.QueryRow(query,
		outletClass.CustId, outletClass.OtClassCode, outletClass.OtClassName, outletClass.OtClassLimit,
		outletClass.IsActive, outletClass.CreatedBy, outletClass.CreatedAt, outletClass.UpdatedBy,
		outletClass.UpdatedAt, outletClass.IsDel, outletClass.DeletedBy, outletClass.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletClassRepository, Store, err:", err.Error())
		return outletClass.OtClassId, err
	}
	return outletClass.OtClassId, nil
}

func (repository *outletClassRepositoryImpl) Update(outletClassId int, request entity.UpdateOutletClassRequest) error {
	var (
		r            model.OutletClassUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("outletClassRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_outlet_class
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND ot_class_id = :ot_class_id_old;`

	// log.Println("outletClassRepository, Update, query:", query)

	sqlPatch.Args["ot_class_id_old"] = outletClassId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletClassRepository, Update, err:", err.Error())
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

func (repository *outletClassRepositoryImpl) Delete(custId string, outletClassId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_outlet_class
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND ot_class_id = :ot_class_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"ot_class_id": outletClassId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OutletClassRepository, Delete, err:", err.Error())
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
