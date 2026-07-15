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

type UnitCoreTaxRepository interface {
	FindOneByUnitIdCoreTaxAndCustId(unitId string, custId string) (model.UnitCoreTax, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.UnitCoreTax, total int, lastPage int, err error)
	Store(unit model.UnitCoreTax) (string, error)
	Update(unitId string, request entity.UpdateUnitCoreTaxRequest) error
	Delete(custId string, unitId string, deletedBy int64) error
}

func NewUnitCoreTaxRepository(db *sqlx.DB) UnitCoreTaxRepository {
	return &unitCoreTaxRepositoryImpl{db}
}

type unitCoreTaxRepositoryImpl struct {
	*sqlx.DB
}

func (repository *unitCoreTaxRepositoryImpl) FindOneByUnitIdCoreTaxAndCustId(unitId string, custId string) (model.UnitCoreTax, error) {
	unit := model.UnitCoreTax{}
	query := `SELECT 
				cust_id, unit_id_coretax, 
				unit_name_coretax, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_unit_coretax 
			  WHERE unit_id_coretax = $1 
			  AND cust_id = $2`
	err := repository.Get(&unit, query, unitId, custId)
	if err != nil {
		log.Println("unitRepository, FindOneByUnitIdCoreTaxAndCustId, err:", err.Error())
		return unit, err
	}

	return unit, nil
}

func (repository *unitCoreTaxRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.UnitCoreTax, int, int, error) {

	units := []model.UnitCoreTax{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` un.cust_id, un.unit_id_coretax, 
					un.unit_name_coretax, un.is_active, un.created_by,
					un.created_at, un.updated_by, un.updated_at,
					un.is_del, un.deleted_by, un.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE un.is_del = false 
				AND un.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (unit_id_coretax ILIKE '%` + dataFilter.Query + `%' 
					OR unit_name_coretax ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND un.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND un.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_unit_coretax un
	LEFT JOIN sys.m_user u ON u.user_id = un.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("unitRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("unitRepository, count total, err:", err.Error())
		return units, 0, 0, err
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
		sortBy := `unit_id_coretax`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	fmt.Println("filteran", dataFilter.Limit, dataFilter.Page)

	all := dataFilter.Limit > 0 || dataFilter.Page > 0

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))
	if all {
		querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))
	}

	// log.Println("unitRepository, querySelect:", querySelect)
	err = repository.Select(&units, querySelect)
	if err != nil {
		log.Println("unitRepository, FindAllByCustId, err:", err.Error())
		return units, total, lastPage, err
	}

	return units, total, lastPage, nil
}

func (repository *unitCoreTaxRepositoryImpl) Store(unit model.UnitCoreTax) (string, error) {
	query :=
		`INSERT INTO mst.m_unit_coretax(
			cust_id, unit_id_coretax, unit_name_coretax, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING unit_id_coretax;`
	lastInsertId := unit.UnitIdCoreTax
	err := repository.QueryRow(query,
		unit.CustId, unit.UnitIdCoreTax, unit.UnitNameCoreTax, unit.IsActive,
		unit.CreatedBy, unit.CreatedAt, unit.UpdatedBy, unit.UpdatedAt,
		unit.IsDel, unit.DeletedBy, unit.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("unitRepository, Store, err:", err.Error())
		return unit.UnitIdCoreTax, err
	}
	return unit.UnitIdCoreTax, nil
}

func (repository *unitCoreTaxRepositoryImpl) Update(unitId string, request entity.UpdateUnitCoreTaxRequest) error {
	var (
		r            model.UnitCoreTaxUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("unitRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_unit_coretax
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND unit_id_coretax = :unit_id_coretax_old;`

	// log.Println("unitRepository, Update, query:", query)

	sqlPatch.Args["unit_id_coretax_old"] = unitId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("unitRepository, Update, err:", err.Error())
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

func (repository *unitCoreTaxRepositoryImpl) Delete(custId string, unitId string, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_unit_coretax
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND unit_id_coretax = :unit_id_coretax;`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"unit_id_coretax": unitId,
		"deleted_by":      deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("UnitCoreTaxRepository, Delete, err:", err.Error())
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
