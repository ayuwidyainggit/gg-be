package repository

import (
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OfficialRepository interface {
	TrxBegin() (*officialTransaction, error)
	FindOneEmployeeByEmpIdAndCustId(empId int, custId string) (model.OfficialEmployee, error)
	FindOneByOfficialIdAndCustId(officialId int, custId string) (model.OfficialList, error)
	FindOneByOfficialTypeAndCustId(officialType int, custId string) (model.OfficialList, error)
	FindAllByCustId(dataFilter entity.OfficialQueryFilter, custId string) (official []model.OfficialList, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.OfficialQueryFilter, custId string) (official []model.OfficialList, total int, lastPage int, err error)
	DeleteAllByCustId(custId string) error
	HierarchyByCustId(dataFilter entity.OfficialQueryFilter, custId string) (official []model.AllOfficialHierarchy, err error)
}

func NewOfficialRepository(db *sqlx.DB) OfficialRepository {
	return &officialRepositoryImpl{db}
}

type officialRepositoryImpl struct {
	*sqlx.DB
}

type officialTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewTransactionOfficial(db *sqlx.DB) (trxObj *officialTransaction, err error) {
	trx := db.MustBegin()

	return &officialTransaction{tx: trx, db: db}, nil
}

func (repo *officialRepositoryImpl) TrxBegin() (*officialTransaction, error) {
	trxObj, err := NewTransactionOfficial(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *officialTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *officialTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *officialRepositoryImpl) FindOneEmployeeByEmpIdAndCustId(empId int, custId string) (model.OfficialEmployee, error) {
	log.Println("empId:", empId)
	employee := model.OfficialEmployee{}
	query := `SELECT 
				cust_id, emp_id, emp_code, emp_name, is_active
			FROM mst.m_employee emp
			WHERE emp_id = $1 
			AND cust_id = $2 `
	err := repository.DB.Get(&employee, query, empId, custId)
	if err != nil {
		log.Println("officialRepository, FindOneEmployeeByEmpIdAndCustId, err:", err.Error())
		return employee, err
	}

	return employee, nil
}

func (repository *officialRepositoryImpl) FindOneByOfficialIdAndCustId(officialId int, custId string) (model.OfficialList, error) {
	official := model.OfficialList{}
	query := `SELECT 
				of.cust_id, of.official_id, of.official_type, of.official_name,
				of.emp_id, of.created_at, of.created_by,
				
				emp.emp_code, emp.emp_name, 
				ofh.hierarchy_code,
				of.supervisor_id as supervisor_id2,
				ofemp_spv2.emp_name as supervisor_name2,
				ofemp_spv2.emp_code as supervisor_code2,
				ofh_spv2.hierarchy_code as hierarchy_code2,
				spv2.supervisor_id as supervisor_id1,
				ofemp_spv1.emp_code as supervisor_code1,
				ofemp_spv1.emp_name as supervisor_name1,
				ofh_spv1.hierarchy_code as hierarchy_code1

			FROM mst.m_official of
			LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type 
			LEFT JOIN mst.m_employee emp ON emp.emp_id = of.emp_id 
			LEFT JOIN mst.m_employee ofemp_spv2 ON ofemp_spv2.emp_id = of.supervisor_id 
			LEFT JOIN mst.m_official spv2 ON spv2.emp_id = of.supervisor_id 
			LEFT JOIN mst.m_official_hierarchy ofh_spv2 ON ofh_spv2.official_type = spv2.official_type
			LEFT JOIN mst.m_employee ofemp_spv1 ON ofemp_spv1.emp_id = spv2.supervisor_id 
			LEFT JOIN mst.m_official spv1 ON spv1.emp_id = spv2.supervisor_id 
			LEFT JOIN mst.m_official_hierarchy ofh_spv1 ON ofh_spv1.official_type = spv1.official_type

			WHERE of.official_id = $1 
			AND of.cust_id = $2 `
	err := repository.Get(&official, query, officialId, custId)
	if err != nil {
		log.Println("officialRepository, FindOneByOfficialIdAndCustId, err:", err.Error())
		return official, err
	}

	return official, nil
}

func (repository *officialRepositoryImpl) FindOneByOfficialTypeAndCustId(officialType int, custId string) (model.OfficialList, error) {
	official := model.OfficialList{}
	query := `SELECT 
				of.cust_id, of.emp_id,
				of.official_id, of.official_type, of.official_name,
				of.created_by, of.created_at, 
				emp.emp_code, emp.emp_name, ofh.hierarchy_code,
				of.supervisor_id as supervisor_id2,
				ofemp_spv2.emp_name as supervisor_name2,
				ofemp_spv2.emp_code as supervisor_code2,
				ofh_spv2.hierarchy_code as hierarchy_code2,
				spv2.supervisor_id as supervisor_id1,
				ofemp_spv1.emp_code as supervisor_code1,
				ofemp_spv1.emp_name as supervisor_name1,
				ofh_spv1.hierarchy_code as hierarchy_code1
			  FROM mst.m_official of
			  LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type 
			  LEFT JOIN mst.m_employee emp ON emp.emp_id = of.emp_id 
			  LEFT JOIN mst.m_employee ofemp_spv2 ON ofemp_spv2.emp_id = of.supervisor_id 
			  LEFT JOIN mst.m_official spv2 ON spv2.emp_id = of.supervisor_id 
			  LEFT JOIN mst.m_official_hierarchy ofh_spv2 ON ofh_spv2.official_type = spv2.official_type
			  LEFT JOIN mst.m_employee ofemp_spv1 ON ofemp_spv1.emp_id = spv2.supervisor_id 
			  LEFT JOIN mst.m_official spv1 ON spv1.emp_id = spv2.supervisor_id 
			  LEFT JOIN mst.m_official_hierarchy ofh_spv1 ON ofh_spv1.official_type = spv1.official_type
			  WHERE of.official_type = $1 
			  AND of.cust_id = $2`
	err := repository.Get(&official, query, officialType, custId)
	if err != nil {
		log.Println("officialRepository, FindOneByOfficialTypeAndCustId, err:", err.Error())
		return official, err
	}

	return official, nil
}

func (repository *officialRepositoryImpl) FindAllByCustId(dataFilter entity.OfficialQueryFilter, custId string) ([]model.OfficialList, int, int, error) {

	officials := []model.OfficialList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` 	of.cust_id, of.official_id, of.official_type, of.official_name,
						of.emp_id, emp.emp_code, emp.emp_name, ofh.hierarchy_code,
						of.supervisor_id as supervisor_id2,
						ofemp_spv2.emp_name as supervisor_name2,
						ofemp_spv2.emp_code as supervisor_code2,
						ofh_spv2.hierarchy_code as hierarchy_code2,
						spv2.supervisor_id as supervisor_id1,
						ofemp_spv1.emp_code as supervisor_code1,
						ofemp_spv1.emp_name as supervisor_name1,
						ofh_spv1.hierarchy_code as hierarchy_code1,
						of.created_at, of.created_by `
	qWhere := ` WHERE of.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (emp.emp_name = '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OfficialType != 0 {
		qWhere += ` AND of.official_type = ` + strconv.Itoa(dataFilter.OfficialType) + ` `
	}

	qFrom := `	FROM mst.m_official of 
				LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type AND ofh.official_type = 3 AND ofh.cust_id = '` + custId + `'
				LEFT JOIN mst.m_employee emp ON emp.emp_id = of.emp_id 
				LEFT JOIN mst.m_employee ofemp_spv2 ON ofemp_spv2.emp_id = of.supervisor_id 
				LEFT JOIN mst.m_official spv2 ON spv2.emp_id = of.supervisor_id 
				LEFT JOIN mst.m_official_hierarchy ofh_spv2 ON ofh_spv2.official_type = spv2.official_type AND ofh_spv2.official_type = 2 AND ofh_spv2.cust_id = '` + custId + `'
				LEFT JOIN mst.m_employee ofemp_spv1 ON ofemp_spv1.emp_id = spv2.supervisor_id 
				LEFT JOIN mst.m_official spv1 ON spv1.emp_id = spv2.supervisor_id 
				LEFT JOIN mst.m_official_hierarchy ofh_spv1 ON ofh_spv1.official_type = spv1.official_type AND ofh_spv1.official_type = 1 AND ofh_spv1.cust_id = '` + custId + `'
				LEFT JOIN sys.m_user u ON u.user_id = of.created_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("officialRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("officialRepository, count total, err:", err.Error())
		return officials, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`of.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `of.official_id`
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

	// log.Println("officialRepository, querySelect:", querySelect)
	err = repository.Select(&officials, querySelect)
	if err != nil {
		log.Println("officialRepository, FindAllByCustId, err:", err.Error())
		return officials, total, lastPage, err
	}

	return officials, total, lastPage, nil
}

func (repository *officialRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.OfficialQueryFilter, custId string) ([]model.OfficialList, int, int, error) {

	officials := []model.OfficialList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` 	of.cust_id, of.official_id, of.official_type, of.official_name,
						of.created_at, of.created_by,
						of.emp_id, emp.emp_code, emp.emp_name, ofh.hierarchy_code,
						of.supervisor_id as supervisor_id2,
						ofemp_spv2.emp_name as supervisor_name2,
						ofemp_spv2.emp_code as supervisor_code2,
						ofh_spv2.hierarchy_code as hierarchy_code2,
						spv2.supervisor_id as supervisor_id1,
						ofemp_spv1.emp_code as supervisor_code1,
						ofemp_spv1.emp_name as supervisor_name1,
						ofh_spv1.hierarchy_code as hierarchy_code1  `

	qWhere := ` WHERE of.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (emp.emp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OfficialType != 0 {
		qWhere += ` AND of.official_type = ` + strconv.Itoa(dataFilter.OfficialType) + ` `
	}

	qFrom := ` FROM mst.m_official of 
			   LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type AND ofh.official_type = 3 AND ofh.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_employee emp ON emp.emp_id = of.emp_id 
			   LEFT JOIN mst.m_employee ofemp_spv2 ON ofemp_spv2.emp_id = of.supervisor_id 
			   LEFT JOIN mst.m_official spv2 ON spv2.emp_id = of.supervisor_id 
			   LEFT JOIN mst.m_official_hierarchy ofh_spv2 ON ofh_spv2.official_type = spv2.official_type AND ofh_spv2.official_type = 2 AND ofh_spv2.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_employee ofemp_spv1 ON ofemp_spv1.emp_id = spv2.supervisor_id 
			   LEFT JOIN mst.m_official spv1 ON spv1.emp_id = spv2.supervisor_id 
			   LEFT JOIN mst.m_official_hierarchy ofh_spv1 ON ofh_spv1.official_type = spv1.official_type AND ofh_spv1.official_type = 1 AND ofh_spv1.cust_id = '` + custId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = of.created_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("officialRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("officialRepository, count total, err:", err.Error())
		return officials, 0, 0, err
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
		sortBy := `of.official_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	lastPage := 1

	// log.Println("officialRepository, querySelect:", querySelect)
	err = repository.Select(&officials, querySelect)
	if err != nil {
		log.Println("officialRepository, FindAllByCustIdLookup, err:", err.Error())
		return officials, total, lastPage, err
	}

	return officials, total, lastPage, nil
}

func (repository *officialRepositoryImpl) DeleteAllByCustId(custId string) error {
	var nRows int64
	query := `DELETE FROM mst.m_official
			  WHERE cust_id = :cust_id;`

	wMap := map[string]interface{}{
		"cust_id": custId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OfficialRepository, DeleteAllByCustId, err:", err.Error())
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

func (repository *officialRepositoryImpl) HierarchyByCustId(dataFilter entity.OfficialQueryFilter, custId string) ([]model.AllOfficialHierarchy, error) {

	officials := []model.AllOfficialHierarchy{}
	selectField := ` 	of.cust_id, moh.hierarchy_code,
						of.official_id, of.official_type, of.official_name,
						of.emp_id, emp.emp_name,
						of.supervisor_id,
						spv_emp.emp_name AS supervisor_name,
						of.created_by, of.created_at`
	qWhere := ` WHERE of.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (emp.emp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OfficialType != 0 {
		qWhere += ` AND of.official_type = ` + strconv.Itoa(dataFilter.OfficialType) + ` `
	}

	qFrom := ` 	FROM mst.m_official of 
				LEFT JOIN mst.m_employee emp ON emp.emp_id = of.emp_id AND emp.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_official spv ON spv.official_id = of.supervisor_id AND spv.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_employee spv_emp ON spv_emp.emp_id = of.supervisor_id AND spv_emp.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_official_hierarchy moh ON moh.official_type = of.official_type AND moh.cust_id = '` + custId + `' 
				LEFT JOIN sys.m_user u ON u.user_id = of.created_by`
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`of.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `of.official_type`
		querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)
	}

	// log.Println("officialRepository, querySelect:", querySelect)
	err := repository.Select(&officials, querySelect)
	if err != nil {
		log.Println("officialRepository, FindAllByCustId, err:", err.Error())
		return officials, err
	}

	return officials, nil
}

func (repository *officialTransaction) StoreWithTrx(official model.Official) (int, error) {
	query :=
		`INSERT INTO mst.m_official(
			cust_id, official_type, official_name ,emp_id, supervisor_id, 
			created_by, created_at)
		VALUES ( 
			$1, $2, $3, $4, $5, 
			$6, $7
		) RETURNING official_id;`
	lastInsertId := official.OfficialId
	err := repository.tx.QueryRow(query,
		official.CustId, official.OfficialType, official.OfficialName, official.EmpId, official.SupervisorId,
		official.CreatedBy, official.CreatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("officialRepository, StoreWithTrx, err:", err.Error())
		return official.OfficialId, err
	}
	return official.OfficialId, nil
}

func (repository *officialTransaction) DeleteAllByCustIdWithTrx(custId string) error {

	query := `DELETE FROM mst.m_official
			  WHERE cust_id = :cust_id;`

	wMap := map[string]interface{}{
		"cust_id": custId,
	}

	_, err := repository.tx.NamedExec(query, wMap)
	if err != nil {
		log.Println("OfficialRepository, DeleteAllByCustIdWithTrx, err:", err.Error())
		return err
	}

	return nil
}
