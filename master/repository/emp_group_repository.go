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

type EmpGroupRepository interface {
	FindOneByEmpGroupIdAndCustId(empGroupId int, custId string) (model.EmpGroup, error)
	FindOneByEmpTypeIdAndCustId(empTypeId string, custId string) (model.EmpType, error)
	FindOneByEmpGroupCodeAndCustId(empGroupCode string, custId string) (model.EmpGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) (empGroup []model.EmpGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) (empGroup []model.EmpGroup, total int, lastPage int, err error)
	EmpTypeFindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (empGroup []model.EmpType, total int, lastPage int, err error)
	Store(empGroup model.EmpGroup) (int, error)
	Update(empGroupId int, request entity.UpdateEmpGroupRequest) error
	Delete(custId string, empGroupId int, deletedBy int64) error
}

func NewEmpGroupRepository(db *sqlx.DB) EmpGroupRepository {
	return &EmpGroupRepositoryImpl{db}
}

type EmpGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *EmpGroupRepositoryImpl) FindOneByEmpGroupIdAndCustId(empGroupId int, custId string) (model.EmpGroup, error) {
	empGroup := model.EmpGroup{}
	query := `SELECT 
				cust_id, emp_grp_id, emp_grp_code,
				emp_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_emp_group 
			  WHERE emp_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&empGroup, query, empGroupId, custId)
	if err != nil {
		log.Println("EmpGroupRepository, FindOneByEmpGroupCodeAndCustId, err:", err.Error())
		return empGroup, err
	}

	return empGroup, nil
}

func (repository *EmpGroupRepositoryImpl) FindOneByEmpTypeIdAndCustId(empTypeId string, custId string) (model.EmpType, error) {
	empType := model.EmpType{}
	query := `SELECT 
				cust_id, emp_type_id, emp_type_name, 
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_emp_type 
			  WHERE emp_type_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&empType, query, empTypeId, custId)
	if err != nil {
		log.Println("EmpTypeRepository, FindOneByEmpTypeIdAndCustId, err:", err.Error())
		return empType, err
	}

	return empType, nil
}

func (repository *EmpGroupRepositoryImpl) FindOneByEmpGroupCodeAndCustId(empGroupCode string, custId string) (model.EmpGroup, error) {
	empGroup := model.EmpGroup{}
	query := `SELECT 
				cust_id, emp_grp_id, emp_grp_code,
				emp_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_emp_group 
			  WHERE emp_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&empGroup, query, empGroupCode, custId)
	if err != nil {
		log.Println("EmpGroupRepository, FindOneByEmpGroupCodeAndCustId, err:", err.Error())
		return empGroup, err
	}

	return empGroup, nil
}

func (repository *EmpGroupRepositoryImpl) EmpTypeFindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.EmpType, int, int, error) {

	empTypes := []model.EmpType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` emp_type_id, emp_type_name `
	qWhere := ` WHERE cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (emp_type_id ILIKE '%` + dataFilter.Query + `%' 
					OR emp_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_emp_type ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_emp_type ` + qWhere

	// log.Println("EmpGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmpGroupRepository, count total, err:", err.Error())
		return empTypes, 0, 0, err
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
		sortBy := `emp_type_id`
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

	// log.Println("EmpGroupRepository, querySelect:", querySelect)
	err = repository.Select(&empTypes, querySelect)
	if err != nil {
		log.Println("EmpGroupRepository, EmpTypeFindAllByCustId, err:", err.Error())
		return empTypes, total, lastPage, err
	}

	return empTypes, total, lastPage, nil
}

func (repository *EmpGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.EmpGroup, int, int, error) {

	empGroups := []model.EmpGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` emp_grp_id, emp_grp_code, emp_grp_name `
	qWhere := ` WHERE is_del = false AND is_active = true
				AND cust_id = '` + dataFilter.ParentCustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (emp_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR emp_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_emp_group ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_emp_group ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmpGroupRepository, FindAllByCustIdLookupMode, count total, err:", err.Error())
		return empGroups, 0, 0, err
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
		sortBy := `emp_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)
	}

	err = repository.Select(&empGroups, querySelect)
	if err != nil {
		log.Println("EmpGroupRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return empGroups, total, 1, err
	}

	return empGroups, total, 1, nil
}

func (repository *EmpGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.EmpGroup, int, int, error) {

	empGroups := []model.EmpGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.emp_grp_id, a.emp_grp_code, a.emp_grp_name, a.is_active, a.created_by, a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + dataFilter.CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.emp_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.emp_grp_name ILIKE '%` + dataFilter.Query + `%' )`
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
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_emp_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_emp_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by ` + qWhere

	// log.Println("EmpGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmpGroupRepository, count total, err:", err.Error())
		return empGroups, 0, 0, err
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
		sortBy := `a.emp_grp_id`
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

	// log.Println("EmpGroupRepository, querySelect:", querySelect)
	err = repository.Select(&empGroups, querySelect)
	if err != nil {
		log.Println("EmpGroupRepository, FindAllByCustId, err:", err.Error())
		return empGroups, total, lastPage, err
	}

	return empGroups, total, lastPage, nil
}

func (repository *EmpGroupRepositoryImpl) Store(empGroup model.EmpGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_emp_group(
			cust_id, emp_grp_code, emp_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING emp_grp_id;`
	lastInsertId := empGroup.EmpGroupId
	err := repository.QueryRow(query,
		empGroup.CustId, empGroup.EmpGroupCode, empGroup.EmpGroupName,
		empGroup.IsActive, empGroup.CreatedBy, empGroup.CreatedAt, empGroup.UpdatedBy,
		empGroup.UpdatedAt, empGroup.IsDel, empGroup.DeletedBy, empGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("EmpGroupRepository, Store, err:", err.Error())
		return empGroup.EmpGroupId, err
	}
	return empGroup.EmpGroupId, nil
}

func (repository *EmpGroupRepositoryImpl) Update(empGroupId int, request entity.UpdateEmpGroupRequest) error {
	var (
		r            model.EmpGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("EmpGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_emp_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND emp_grp_id = :emp_grp_id_old;`

	// log.Println("EmpGroupRepository, Update, query:", query)

	sqlPatch.Args["emp_grp_id_old"] = empGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("EmpGroupRepository, Update, err:", err.Error())
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

func (repository *EmpGroupRepositoryImpl) Delete(custId string, empGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_emp_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND emp_grp_id = :emp_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"emp_grp_id": empGroupId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("EmpGroupRepository, Delete, err:", err.Error())
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
