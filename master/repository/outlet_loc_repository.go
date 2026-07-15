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

type OutletLocRepository interface {
	FindOneByOutletLocIdAndCustId(outletLocId int, custId string) (model.OutletLoc, error)
	FindOneByOutletLocCodeAndCustId(outletLocCode string, custId string) (model.OutletLoc, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (outletLoc []model.OutletLoc, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (outletLoc []model.OutletLoc, total int, lastPage int, err error)
	Store(outletLoc model.OutletLoc) (int, error)
	Update(outletLocId int, request entity.UpdateOutletLocRequest) error
	Delete(custId string, outletLocId int, deletedBy int64) error
}

func NewOutletLocRepository(db *sqlx.DB) OutletLocRepository {
	return &outletLocRepositoryImpl{db}
}

type outletLocRepositoryImpl struct {
	*sqlx.DB
}

func (repository *outletLocRepositoryImpl) FindOneByOutletLocIdAndCustId(outletLocId int, custId string) (model.OutletLoc, error) {
	outletLoc := model.OutletLoc{}
	query := `SELECT 
				cust_id, ot_loc_id, ot_loc_code,
				ot_loc_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_loc 
			  WHERE ot_loc_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletLoc, query, outletLocId, custId)
	if err != nil {
		log.Println("outletLocRepository, FindOneByOutletLocCodeAndCustId, err:", err.Error())
		return outletLoc, err
	}

	return outletLoc, nil
}

func (repository *outletLocRepositoryImpl) FindOneByOutletLocCodeAndCustId(outletLocCode string, custId string) (model.OutletLoc, error) {
	outletLoc := model.OutletLoc{}
	query := `SELECT 
				cust_id, ot_loc_id, ot_loc_code,
				ot_loc_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_loc 
			  WHERE ot_loc_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletLoc, query, outletLocCode, custId)
	if err != nil {
		log.Println("outletLocRepository, FindOneByOutletLocCodeAndCustId, err:", err.Error())
		return outletLoc, err
	}

	return outletLoc, nil
}

func (repository *outletLocRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletLoc, int, int, error) {

	outletLocs := []model.OutletLoc{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_loc_id, a.ot_loc_code, a.ot_loc_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_loc_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_loc_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_outlet_loc a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("outletLocRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletLocRepository, count total, err:", err.Error())
		return outletLocs, 0, 0, err
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
		sortBy := `a.ot_loc_id`
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

	// log.Println("outletLocRepository, querySelect:", querySelect)
	err = repository.Select(&outletLocs, querySelect)
	if err != nil {
		log.Println("outletLocRepository, FindAllByCustId, err:", err.Error())
		return outletLocs, total, lastPage, err
	}

	return outletLocs, total, lastPage, nil
}

func (repository *outletLocRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletLoc, int, int, error) {

	outletLocs := []model.OutletLoc{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_loc_id, a.ot_loc_code, a.ot_loc_name  `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_loc_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_loc_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_outlet_loc a  `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("outletLocRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletLocRepository, count total, err:", err.Error())
		return outletLocs, 0, 0, err
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
		sortBy := `a.ot_loc_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("outletLocRepository, querySelect:", querySelect)
	err = repository.Select(&outletLocs, querySelect)
	if err != nil {
		log.Println("outletLocRepository, FindAllByCustId, err:", err.Error())
		return outletLocs, total, 1, err
	}

	return outletLocs, total, 1, nil
}

func (repository *outletLocRepositoryImpl) Store(outletLoc model.OutletLoc) (int, error) {
	query :=
		`INSERT INTO mst.m_outlet_loc(
			cust_id, ot_loc_code, ot_loc_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING ot_loc_id;`
	lastInsertId := outletLoc.OtLocId
	err := repository.QueryRow(query,
		outletLoc.CustId, outletLoc.OtLocCode, outletLoc.OtLocName,
		outletLoc.IsActive, outletLoc.CreatedBy, outletLoc.CreatedAt, outletLoc.UpdatedBy,
		outletLoc.UpdatedAt, outletLoc.IsDel, outletLoc.DeletedBy, outletLoc.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletLocRepository, Store, err:", err.Error())
		return outletLoc.OtLocId, err
	}
	return outletLoc.OtLocId, nil
}

func (repository *outletLocRepositoryImpl) Update(outletLocId int, request entity.UpdateOutletLocRequest) error {
	var (
		r            model.OutletLocUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("outletLocRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_outlet_loc
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND ot_loc_id = :ot_loc_id_old;`

	// log.Println("outletLocRepository, Update, query:", query)

	sqlPatch.Args["ot_loc_id_old"] = outletLocId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletLocRepository, Update, err:", err.Error())
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

func (repository *outletLocRepositoryImpl) Delete(custId string, outletLocId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_outlet_loc
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND ot_loc_id = :ot_loc_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"ot_loc_id":  outletLocId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OutletLocRepository, Delete, err:", err.Error())
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
