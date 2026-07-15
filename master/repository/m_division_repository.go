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

type MDivisionRepository interface {
	FindOneByMDivisionIdAndCustId(MDivisionId int64, custId string) (model.MDivision, error)
	FindOneByMDivisionCodeAndCustId(divisionCode, custId string) (model.MDivision, error)
	Store(MDivision model.MDivision) (int64, error)
	Update(divisionID int64, request entity.UpdateDivisionBody) error
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MDivision, int, int, error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.MDivision, int, int, error)
	Delete(custId string, MDivisionId int64, deletedBy int64) error
}

func NewMDivisionRepository(db *sqlx.DB) MDivisionRepository {
	return &MDivisionRepositoryImpl{db}
}

type MDivisionRepositoryImpl struct {
	*sqlx.DB
}

func (repository *MDivisionRepositoryImpl) FindOneByMDivisionIdAndCustId(MDivisionId int64, custId string) (model.MDivision, error) {
	MDivision := model.MDivision{}
	query := `SELECT 
				*
			  FROM mst.m_division 
			  WHERE division_id = $1 AND is_del = false
			  AND cust_id = $2`
	err := repository.Get(&MDivision, query, MDivisionId, custId)
	if err != nil {
		log.Println("MDivisionRepository, FindOneByMDivisionIdAndCustId, err:", err.Error())
		return MDivision, err
	}

	return MDivision, nil
}

func (repository *MDivisionRepositoryImpl) FindOneByMDivisionCodeAndCustId(divisionCode, custId string) (model.MDivision, error) {
	MDivision := model.MDivision{}
	query := `SELECT 
				*
			  FROM mst.m_division 
			  WHERE division_code = $1 AND is_del = false
			  AND cust_id = $2`
	err := repository.Get(&MDivision, query, divisionCode, custId)
	if err != nil {
		log.Println("MDivisionRepository, FindOneByMDivisionIdAndCustId, err:", err.Error())
		return MDivision, err
	}

	return MDivision, nil
}

func (repository *MDivisionRepositoryImpl) Store(MDivision model.MDivision) (int64, error) {
	query :=
		`INSERT INTO mst.m_division(
			cust_id, division_code, division_name, is_active, created_by, created_at, 
			updated_by, updated_at, is_del, deleted_by, 
			deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING division_id;`
	lastInsertId := MDivision.DivisionID
	err := repository.QueryRow(query,
		MDivision.CustId, MDivision.DivisionCode, MDivision.DivisionName, MDivision.IsActive, MDivision.CreatedBy, MDivision.CreatedAt,
		MDivision.UpdatedBy, MDivision.UpdatedAt, MDivision.IsDel, MDivision.DeletedBy,
		MDivision.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("MDivisionRepository, Store, err:", err.Error())
		return MDivision.DivisionID, err
	}
	return MDivision.DivisionID, nil
}

func (repository *MDivisionRepositoryImpl) Update(divisionID int64, request entity.UpdateDivisionBody) error {
	var (
		r            model.MDivisionUpdate
		sqlSetFields string
		nRows        int64
	)

	// requestFormat, _ := json.Marshal(request)
	// fmt.Printf("divisionRepository, Update, request: %s\n", requestFormat)
	// panic("test")
	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("divisionRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_division
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND division_id = :division_id_old;`

	// log.Println("divisionRepository, Update, query:", query)

	sqlPatch.Args["division_id_old"] = divisionID
	sqlPatch.Args["cust_id"] = request.CustId

	// s.hooks.reset()

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("divisionRepository, Update, err:", err.Error())
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

func (repository *MDivisionRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MDivision, int, int, error) {

	divisions := []model.MDivision{}
	selectCount := ` COUNT(*) AS total `
	selectField := `d.*, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE d.is_del = false 
				AND d.cust_id = '` + dataFilter.CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (d.division_code ILIKE '%` + dataFilter.Query + `%' 
					OR d.division_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND d.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND d.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_division d
			   LEFT JOIN sys.m_user u ON u.user_id = d.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("EmployeeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, count total, err:", err.Error())
		return divisions, 0, 0, err
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
		sortBy := `d.division_id`
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

	// log.Println("EmployeeRepository, querySelect:", querySelect)
	err = repository.Select(&divisions, querySelect)
	if err != nil {
		log.Println("DivisionRepository, FindAllByCustId, err:", err.Error())
		return divisions, total, lastPage, err
	}

	return divisions, total, lastPage, nil
}

func (repository *MDivisionRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.MDivision, int, int, error) {

	divisions := []model.MDivision{}
	selectCount := ` COUNT(*) AS total `
	selectField := `d.*, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE d.is_del = false 
				AND d.cust_id = '` + dataFilter.ParentCustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (d.division_code ILIKE '%` + dataFilter.Query + `%' 
					OR d.division_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		// fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND d.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND d.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_division d
			   LEFT JOIN sys.m_user u ON u.user_id = d.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("EmployeeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, count total, err:", err.Error())
		return divisions, 0, 0, err
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
		sortBy := `d.division_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}

	if dataFilter.Limit > 1000 {
		dataFilter.Limit = 1000
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("EmployeeRepository, querySelect:", querySelect)
	err = repository.Select(&divisions, querySelect)
	if err != nil {
		log.Println("DivisionRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return divisions, total, lastPage, err
	}

	return divisions, total, lastPage, nil
}

func (repository *MDivisionRepositoryImpl) Delete(custId string, MDivisionId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_division
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND division_id = :division_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"division_id": MDivisionId,
		"deleted_by":  deletedBy,
	}
	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("mDivisionRepository, Delete, err:", err.Error())
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
