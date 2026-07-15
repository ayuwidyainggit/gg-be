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

type EmployeeRepository interface {
	FindOneByEmployeeIdAndCustId(params entity.DetailEmployeeParams) (model.Employee, error)
	FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error)
	FindEmployeeTerritoryDetail(params entity.DetailEmployeeParams) (model.EmployeeTerritoryDetail, error)
	FindOneByEmployeeCodeAndCustId(params entity.DetailEmployeeParams) (model.Employee, error)
	FindAllByCustId(dataFilter entity.EmployeeQueryFilter) (consPro []model.Employee, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.EmployeeQueryFilter) (consPro []model.Employee, total int, lastPage int, err error)
	FindAllByCustIdLookupModeWithoutSalesman(dataFilter entity.EmployeeQueryFilter) (consPro []model.Employee, total int, lastPage int, err error)
	FindLookupMinimal(dataFilter entity.EmployeeLookupAPIFilter) (rows []model.EmployeeLookupMinimal, total int, lastPage int, err error)
	FindAllExport(dataFilter entity.EmployeeQueryFilter) ([]model.EmployeeExport, int, error)
	Store(employee model.Employee) (int, error)
	ValidateEmployeeTerritoryMapping(custId string, territory model.EmployeeTerritoryMapping) error
	Update(employeeId int, request entity.UpdateEmployeeRequest) error
	Delete(custId string, employeeId int, deletedBy int64) error
	UpdateSalesman(employeeId int, isActive bool) error
	CreateImportHistory(uploadType, fileName, custId string, uploadedBy int64, totalData int) (int64, error)
	UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error
	ImportUpdate(employeeId int, request entity.UpdateEmployeeRequest) error
	InsertEmployeeTemp(temp entity.ImportEmployeeTemp) error
	InsertEmployeeUpdateTemp(temp entity.ImportEmployeeUpdateTemp) error
	GetImportTotalData(historyId int64) (int, error)
	GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error)
	GetEmployeeDataForTemplateUpdate(custId string, fields []string) (map[string][][]string, error)
	FindEmpGroupIdByName(custId string, name string) (int, error)
	FindDivisionIdByName(custId string, name string) (int, error)
	FindProvinceIdByName(custId string, name string) (string, error)
	FindCityIdByName(custId string, name string) (string, error)
	FindSubDistrictIdByName(custId string, name string) (string, error)
	FindWardIdByName(custId string, name string) (string, error)
	CheckEmployeeExists(custId, empCode string) (bool, error)
	FindAllForPJP(dataFilter entity.EmployeePJPQueryFilter) (employees []model.EmployeePJP, total int, lastPage int, err error)
}

func NewEmployeeRepository(db *sqlx.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{db}
}

type EmployeeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *EmployeeRepositoryImpl) FindOneByEmployeeIdAndCustId(params entity.DetailEmployeeParams) (model.Employee, error) {
	employee := model.Employee{}
	query := `SELECT 
				ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
				ep.address, ep.emp_type_id, ep.emp_grp_id,
				et.emp_type_name, eg.emp_grp_code, eg.emp_grp_name, 
				ep.work_date, ep.last_education, ep.dob, ep.phone_no,
				ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
				ep.created_at, ep.updated_by, ep.updated_at,
				ep.is_del, ep.deleted_by, ep.deleted_at, ep.province_id, ep.city_id, ep.sub_district_id, ep.ward_id,
				p.province, r.regency as city, sd.sub_district, w.ward,
				ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
				d.division_name,
				ep.region_scope, ep.area_scope, ep.distributor_scope
			  FROM mst.m_employee ep
			  LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + params.ParentCustId + `'
			  WHERE ep.emp_id = $1 
			  AND ep.cust_id = $2`
	err := repository.Get(&employee, query, params.EmployeeId, params.CustId)
	if err != nil {
		log.Println("EmployeeRepository, FindOneByEmployeeIdAndCustId, err:", err.Error())
		return employee, err
	}

	return employee, nil
}

func (repository *EmployeeRepositoryImpl) FindEmployeeTerritoryDetail(params entity.DetailEmployeeParams) (model.EmployeeTerritoryDetail, error) {
	detail := model.EmployeeTerritoryDetail{
		Regions:      []model.EmployeeRegionMappingDetail{},
		Areas:        []model.EmployeeAreaMappingDetail{},
		Distributors: []model.EmployeeDistributorMappingDetail{},
	}

	regionQuery := `SELECT
			COALESCE(r.region_id, 0) AS region_id,
			COALESCE(r.region_code, '') AS region_code,
			COALESCE(r.region_name, '') AS region_name
		FROM mst.m_employee_region_mapping erm
		INNER JOIN mst.m_region r
			ON r.region_id = erm.region_id
			AND r.cust_id = $3
			AND r.is_del = false
		WHERE erm.cust_id = $1
			AND erm.emp_id = $2
			AND erm.is_del = false
		ORDER BY r.region_name, r.region_id`
	if err := repository.Select(&detail.Regions, regionQuery, params.CustId, params.EmployeeId, params.ParentCustId); err != nil {
		log.Println("EmployeeRepository, FindEmployeeTerritoryDetail regions, err:", err.Error())
		return detail, err
	}

	areaQuery := `SELECT
			a.area_id,
			a.area_code,
			a.area_name,
			COALESCE(r.region_id, 0) AS region_id,
			COALESCE(r.region_code, '') AS region_code,
			COALESCE(r.region_name, '') AS region_name
		FROM mst.m_employee_area_mapping eam
		INNER JOIN mst.m_area a
			ON a.area_id = eam.area_id
			AND a.cust_id = $3
			AND a.is_del = false
		LEFT JOIN mst.m_region r
			ON r.region_id = a.region_id
			AND r.cust_id = $3
			AND r.is_del = false
		WHERE eam.cust_id = $1
			AND eam.emp_id = $2
			AND eam.is_del = false
		ORDER BY r.region_name, a.area_name, a.area_id`
	if err := repository.Select(&detail.Areas, areaQuery, params.CustId, params.EmployeeId, params.ParentCustId); err != nil {
		log.Println("EmployeeRepository, FindEmployeeTerritoryDetail areas, err:", err.Error())
		return detail, err
	}

	distributorQuery := `SELECT
			d.distributor_id,
			d.distributor_code,
			d.distributor_name,
			COALESCE(a.area_id, 0) AS area_id,
			COALESCE(a.area_code, '') AS area_code,
			COALESCE(a.area_name, '') AS area_name,
			COALESCE(r.region_id, 0) AS region_id,
			COALESCE(r.region_code, '') AS region_code,
			COALESCE(r.region_name, '') AS region_name
		FROM mst.m_employee_distributor_mapping edm
		INNER JOIN mst.m_distributor d
			ON d.distributor_id = edm.distributor_id
			AND d.is_del = false
			AND (d.cust_id = $3 OR d.cust_id LIKE $4 OR d.parent_cust_id = $3)
		LEFT JOIN mst.m_area a
			ON a.area_id = d.area_id
			AND a.cust_id = $3
			AND a.is_del = false
		LEFT JOIN mst.m_region r
			ON r.region_id = a.region_id
			AND r.cust_id = $3
			AND r.is_del = false
		WHERE edm.cust_id = $1
			AND edm.emp_id = $2
			AND edm.is_del = false
		ORDER BY r.region_name, a.area_name, d.distributor_name, d.distributor_id`
	if err := repository.Select(&detail.Distributors, distributorQuery, params.CustId, params.EmployeeId, params.ParentCustId, params.ParentCustId+"%"); err != nil {
		log.Println("EmployeeRepository, FindEmployeeTerritoryDetail distributors, err:", err.Error())
		return detail, err
	}

	return detail, nil
}

func (repository *EmployeeRepositoryImpl) FindEmployeeDropdownScope(empID int, custID string) (model.Employee, error) {
	employee := model.Employee{}
	query := `SELECT cust_id, emp_id, COALESCE(region_scope,'') AS region_scope, COALESCE(area_scope,'') AS area_scope, COALESCE(distributor_scope,'') AS distributor_scope
		FROM mst.m_employee
		WHERE emp_id = $1 AND cust_id = $2 AND is_del = false`
	err := repository.Get(&employee, query, empID, custID)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

func (repository *EmployeeRepositoryImpl) FindOneByEmployeeCodeAndCustId(params entity.DetailEmployeeParams) (model.Employee, error) {
	employee := model.Employee{}
	query := `SELECT 
				ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
				ep.address, ep.emp_type_id, ep.emp_grp_id,
				ep.work_date, ep.last_education, ep.dob, ep.phone_no,
				ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
				ep.created_at, ep.updated_by, ep.updated_at,
				ep.is_del, ep.deleted_by, ep.deleted_at, 
				ep.province_id, ep.city_id, ep.sub_district_id, ep.ward_id, 
				p.province, r.regency as city, sd.sub_district, w.ward,
				ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
				d.division_name
			  FROM mst.m_employee ep
			  LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + params.ParentCustId + `'
			  LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + params.ParentCustId + `'
			  WHERE ep.cust_id = $2
			  AND ep.emp_code = $1 
			  AND ep.is_del = false`
	err := repository.Get(&employee, query, params.EmployeeCode, params.CustId)
	if err != nil {
		log.Println("EmployeeRepository, FindOneByEmployeeCodeAndCustId, err:", err.Error())
		return employee, err
	}

	return employee, nil
}

func (repository *EmployeeRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.EmployeeQueryFilter) ([]model.Employee, int, int, error) {

	employees := []model.Employee{}
	selectCount := ` COUNT(DISTINCT ep.emp_id) AS total `
	selectField := `ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
					ep.address, r.regency as city, ep.emp_type_id, ep.emp_grp_id,
					et.emp_type_name, eg.emp_grp_code, eg.emp_grp_name, 
					ep.work_date, ep.last_education, ep.dob, ep.phone_no,
					ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
					ep.created_at, ep.updated_by, ep.updated_at,
					ep.is_del, ep.deleted_by, ep.deleted_at, ep.province_id, ep.city_id,
					ep.sub_district_id, ep.ward_id, p.province, r.regency as city, sd.sub_district, w.ward,
					ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
					d.division_name `
	qWhere := ` WHERE ep.is_del = false AND ep.is_active = true 
				AND ep.cust_id = '` + dataFilter.CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (ep.emp_code ILIKE '%` + dataFilter.Query + `%' 
					OR ep.emp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.EmpTypeId != "" {
		qWhere += ` AND ep.emp_type_id = '` + dataFilter.EmpTypeId + `' `
	}

	if dataFilter.EmpGrpId > 0 {
		qWhere += ` AND ep.emp_grp_id = ` + strconv.Itoa(dataFilter.EmpGrpId) + ` `
	}

	if dataFilter.EmpGrpName != "" {
		qWhere += ` AND eg.emp_grp_name ILIKE '%` + dataFilter.EmpGrpName + `%' `
	}

	if dataFilter.DivisionId > 0 {
		qWhere += ` AND ep.division_id = ` + strconv.Itoa(dataFilter.DivisionId) + ` `
	}

	qFrom := ` FROM mst.m_employee ep
			   LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + dataFilter.ParentCustId + `' 
			   LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + dataFilter.ParentCustId + `' `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("EmployeeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, count total, err:", err.Error())
		return employees, 0, 0, err
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
		sortBy := `ep.emp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("EmployeeRepository, querySelect:", querySelect)
	err = repository.Select(&employees, querySelect)
	if err != nil {
		log.Println("EmployeeRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return employees, total, 1, err
	}

	return employees, total, 1, nil
}

func (repository *EmployeeRepositoryImpl) FindAllByCustIdLookupModeWithoutSalesman(dataFilter entity.EmployeeQueryFilter) ([]model.Employee, int, int, error) {
	employees := []model.Employee{}
	selectCount := ` COUNT(*) AS total `
	selectField := `ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
					ep.address, r.regency as city, ep.emp_type_id, ep.emp_grp_id,
					et.emp_type_name, eg.emp_grp_code, eg.emp_grp_name, 
					ep.work_date, ep.last_education, ep.dob, ep.phone_no,
					ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
					ep.created_at, ep.updated_by, ep.updated_at,
					ep.is_del, ep.deleted_by, ep.deleted_at, ep.province_id, ep.city_id,
					ep.sub_district_id, ep.ward_id, p.province, r.regency as city, sd.sub_district, w.ward,
					ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
					d.division_name, case when ms.emp_id is null then 0 else ms.emp_id end as emp_id_salesman`
	qWhere := ` WHERE ep.is_del = false AND ep.is_active = true 
				AND ep.cust_id = '` + dataFilter.CustId + `' 
				AND ep.emp_id NOT IN (
					select ep.emp_id
					FROM
						mst.m_employee ep
					LEFT JOIN mst.m_salesman ms on ms.emp_id = ep.emp_id
						and ms.cust_id = '` + dataFilter.CustId + `'
					WHERE
						ep.is_del = false
						AND ep.is_active = true
						AND ms.is_del = false
						AND ep.cust_id = '` + dataFilter.CustId + `'
				) `

	if dataFilter.Query != "" {
		qWhere += ` AND (ep.emp_code ILIKE '%` + dataFilter.Query + `%' 
					OR ep.emp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.EmpTypeId != "" {
		qWhere += ` AND ep.emp_type_id = '` + dataFilter.EmpTypeId + `' `
	}

	// Temporary: hide emp_grp_id filter for without-salesman lookup.
	// if dataFilter.EmpGrpId > 0 {
	// 	qWhere += ` AND ep.emp_grp_id = ` + strconv.Itoa(dataFilter.EmpGrpId) + ` `
	// }

	if dataFilter.DivisionId > 0 {
		qWhere += ` AND ep.division_id = ` + strconv.Itoa(dataFilter.DivisionId) + ` `
	}

	groupBy := `GROUP BY ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
	ep.address, r.regency, ep.emp_type_id, ep.emp_grp_id,
	et.emp_type_name, eg.emp_grp_code, eg.emp_grp_name, 
	ep.work_date, ep.last_education, ep.dob, ep.phone_no,
	ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
	ep.created_at, ep.updated_by, ep.updated_at,
	ep.is_del, ep.deleted_by, ep.deleted_at, ep.province_id, ep.city_id,
	ep.sub_district_id, ep.ward_id, p.province, r.regency, sd.sub_district, w.ward,
	ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
	d.division_name, case when ms.emp_id is null then 0 else ms.emp_id end `

	qFrom := ` FROM mst.m_employee ep
			   LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + dataFilter.ParentCustId + `' 
			   LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_salesman ms on ms.emp_id = ep.emp_id and ms.cust_id = '` + dataFilter.CustId + `' `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere + groupBy

	// log.Println("EmployeeRepository, queryCount:", querySelect)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, count total, err:", err.Error())
		return employees, 0, 0, err
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
		sortBy := `ep.emp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("EmployeeRepository, querySelect:", querySelect)
	err = repository.Select(&employees, querySelect)
	if err != nil {
		log.Println("EmployeeRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return employees, total, 1, err
	}

	return employees, total, 1, nil
}

func buildEmployeeLookupSortSQL(sort string) string {
	if strings.TrimSpace(sort) == "" {
		return "ep.created_at DESC"
	}
	parts := strings.Split(sort, ",")
	var out []string
	colMap := map[string]string{
		"created_date": "ep.created_at",
		"emp_id":       "ep.emp_id",
		"emp_code":     "ep.emp_code",
		"emp_name":     "ep.emp_name",
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		col := strings.ToLower(strings.TrimSpace(kv[0]))
		dir := "DESC"
		if len(kv) > 1 {
			d := strings.ToUpper(strings.TrimSpace(kv[1]))
			if d == "ASC" || d == "DESC" {
				dir = d
			}
		}
		if sqlCol, ok := colMap[col]; ok {
			out = append(out, sqlCol+" "+dir)
		}
	}
	if len(out) == 0 {
		return "ep.created_at DESC"
	}
	return strings.Join(out, ", ")
}

func (repository *EmployeeRepositoryImpl) FindLookupMinimal(dataFilter entity.EmployeeLookupAPIFilter) ([]model.EmployeeLookupMinimal, int, int, error) {
	custIds := dataFilter.FilterCustIds
	if len(custIds) == 0 {
		custIds = []string{dataFilter.CustId}
	}

	inClause, inArgs, err := sqlx.In("ep.cust_id IN (?)", custIds)
	if err != nil {
		return nil, 0, 0, err
	}

	whereParts := []string{
		"ep.is_del = false",
		"ep.is_active = true",
		inClause,
	}
	args := append([]interface{}{}, inArgs...)

	if dataFilter.Query != "" {
		whereParts = append(whereParts, "(ep.emp_code ILIKE ? OR ep.emp_name ILIKE ?)")
		like := "%" + dataFilter.Query + "%"
		args = append(args, like, like)
	}

	whereSQL := strings.Join(whereParts, " AND ")
	fromSQL := ` FROM mst.m_employee ep `

	countQuery := `SELECT COUNT(*) ` + fromSQL + ` WHERE ` + whereSQL
	countQuery = repository.Rebind(countQuery)

	var total int
	err = repository.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, FindLookupMinimal count, err:", err.Error())
		return nil, 0, 0, err
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 5
	}
	offset := (page - 1) * limit

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if total == 0 {
		lastPage = 0
	}

	sortSQL := buildEmployeeLookupSortSQL(dataFilter.Sort)

	selectQuery := `SELECT ep.emp_id, ep.emp_code, ep.emp_name ` + fromSQL + ` WHERE ` + whereSQL + ` ORDER BY ` + sortSQL + ` LIMIT ? OFFSET ?`
	selectQuery = repository.Rebind(selectQuery)
	selArgs := append(append([]interface{}{}, args...), limit, offset)

	rows := []model.EmployeeLookupMinimal{}
	err = repository.Select(&rows, selectQuery, selArgs...)
	if err != nil {
		log.Println("EmployeeRepository, FindLookupMinimal, err:", err.Error())
		return nil, total, lastPage, err
	}

	return rows, total, lastPage, nil
}

func (repository *EmployeeRepositoryImpl) FindAllExport(dataFilter entity.EmployeeQueryFilter) ([]model.EmployeeExport, int, error) {
	employees := []model.EmployeeExport{}
	parentCustId := dataFilter.ParentCustId
	if parentCustId == "" {
		parentCustId = dataFilter.CustId
	}
	whereParts := []string{"ep.is_del = false", "ep.cust_id = $1"}
	args := []interface{}{dataFilter.CustId}
	argIdx := len(args)

	if q := strings.TrimSpace(dataFilter.Query); q != "" {
		argIdx++
		pattern := "%" + q + "%"
		whereParts = append(whereParts, fmt.Sprintf("(ep.emp_code ILIKE $%d OR ep.emp_name ILIKE $%d)", argIdx, argIdx))
		args = append(args, pattern)
	}

	if dataFilter.EmpTypeId != "" {
		argIdx++
		whereParts = append(whereParts, fmt.Sprintf("ep.emp_type_id = $%d", argIdx))
		args = append(args, dataFilter.EmpTypeId)
	}

	if dataFilter.EmpGrpId > 0 {
		argIdx++
		whereParts = append(whereParts, fmt.Sprintf("ep.emp_grp_id = $%d", argIdx))
		args = append(args, dataFilter.EmpGrpId)
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 || *dataFilter.IsActive == 2 {
			argIdx++
			isActive := *dataFilter.IsActive == 1
			whereParts = append(whereParts, fmt.Sprintf("ep.is_active = $%d", argIdx))
			args = append(args, isActive)
		}
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = " WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(1) FROM mst.m_employee ep" + whereClause
	var total int
	if err := repository.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Println("EmployeeRepository, FindAllExport, count err:", err.Error())
		return employees, 0, err
	}

	selectFields := `
		ep.cust_id,
		ep.emp_id,
		ep.emp_code,
		ep.emp_name,
		ep.address,
		ep.emp_type_id,
		et.emp_type_name,
		ep.emp_grp_id,
		eg.emp_grp_code,
		eg.emp_grp_name,
		ep.work_date,
		ep.last_education,
		ep.dob,
		ep.phone_no,
		ep.wa_no,
		ep.email,
		ep.is_active,
		ep.is_del,
		ep.device_id,
		ep.mac_address,
		ep.image_url,
		ep.identity_no,
		ep.is_wa_no,
		ep.province_id,
		p.province,
		ep.city_id,
		r.regency AS city,
		ep.sub_district_id,
		sd.sub_district,
		ep.ward_id,
		w.ward,
		ep.post_code,
		ep.division_id,
		d.division_code,
		d.division_name
	`

	fromClause := `
		FROM mst.m_employee ep
		LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + parentCustId + `'
	`

	selectQuery := "SELECT " + selectFields + " " + fromClause + whereClause + " ORDER BY ep.emp_id DESC"
	if err := repository.Select(&employees, selectQuery, args...); err != nil {
		log.Println("EmployeeRepository, FindAllExport, select err:", err.Error())
		return employees, total, err
	}

	return employees, total, nil
}

func (repository *EmployeeRepositoryImpl) FindAllByCustId(dataFilter entity.EmployeeQueryFilter) ([]model.Employee, int, int, error) {

	employees := []model.Employee{}
	selectCount := ` COUNT(*) AS total `
	selectField := `ep.cust_id, ep.emp_id, ep.emp_code, ep.emp_name, 
					ep.address, ep.emp_type_id, ep.emp_grp_id,
					et.emp_type_name, eg.emp_grp_code, eg.emp_grp_name, 
					ep.province_id, ep.city_id, ep.sub_district_id, ep.ward_id,
					p.province, r.regency as city, sd.sub_district, w.ward,
					ep.work_date, ep.last_education, ep.dob, ep.phone_no,
					ep.wa_no, ep.email, ep.image_url, ep.is_active, ep.created_by,
					ep.created_at, ep.updated_by, ep.updated_at,
					ep.is_del, ep.deleted_by, ep.deleted_at, u.user_fullname AS updated_by_name,
					ep.identity_no, ep.is_wa_no, ep.post_code, ep.division_id,
					d.division_name `
	qWhere := ` WHERE ep.is_del = false 
				AND ep.cust_id = '` + dataFilter.CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (ep.emp_code ILIKE '%` + dataFilter.Query + `%' 
					OR ep.emp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND ep.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND ep.is_active = false `
		}
	}

	if dataFilter.EmpTypeId != "" {
		qWhere += ` AND ep.emp_type_id = '` + dataFilter.EmpTypeId + `' `
	}

	if dataFilter.EmpGrpId > 0 {
		qWhere += ` AND ep.emp_grp_id = ` + strconv.Itoa(dataFilter.EmpGrpId) + ` `
	}

	if dataFilter.DivisionId > 0 {
		qWhere += ` AND ep.division_id = ` + strconv.Itoa(dataFilter.DivisionId) + ` `
	}

	qFrom := ` FROM mst.m_employee ep
			   LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + dataFilter.ParentCustId + `' 
			   LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = ep.updated_by
			   LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND d.cust_id = '` + dataFilter.ParentCustId + `' `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("EmployeeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("EmployeeRepository, count total, err:", err.Error())
		return employees, 0, 0, err
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
		sortBy := `ep.emp_id`
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
	err = repository.Select(&employees, querySelect)
	if err != nil {
		log.Println("EmployeeRepository, FindAllByCustId, err:", err.Error())
		return employees, total, lastPage, err
	}

	return employees, total, lastPage, nil
}

func (repository *EmployeeRepositoryImpl) ValidateEmployeeTerritoryMapping(custId string, territory model.EmployeeTerritoryMapping) error {
	if len(territory.RegionIds) > 0 {
		if err := repository.validateEmployeeRegionIds(custId, territory.RegionIds); err != nil {
			return err
		}
	}
	if len(territory.AreaIds) > 0 {
		if err := repository.validateEmployeeAreaIds(custId, territory.AreaIds, territory.RegionIds); err != nil {
			return err
		}
	}
	if len(territory.DistributorIds) > 0 {
		if err := repository.validateEmployeeDistributorIds(custId, territory.DistributorIds, territory.AreaIds); err != nil {
			return err
		}
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) validateEmployeeRegionIds(custId string, regionIds []int) error {
	total, err := repository.countEmployeeTerritoryRows(
		`SELECT COUNT(DISTINCT region_id)
		 FROM mst.m_region
		 WHERE cust_id = ?
		 AND is_del = false
		 AND region_id IN (?)`,
		custId,
		regionIds,
	)
	if err != nil {
		return err
	}
	if total != len(regionIds) {
		return errors.New("region_ids contains invalid region_id")
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) validateEmployeeAreaIds(custId string, areaIds []int, regionIds []int) error {
	total, err := repository.countEmployeeTerritoryRows(
		`SELECT COUNT(DISTINCT area_id)
		 FROM mst.m_area
		 WHERE cust_id = ?
		 AND is_del = false
		 AND area_id IN (?)
		 AND region_id IN (?)`,
		custId,
		areaIds,
		regionIds,
	)
	if err != nil {
		return err
	}
	if total != len(areaIds) {
		return errors.New("area_ids must belong to selected region_ids")
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) validateEmployeeDistributorIds(custId string, distributorIds []int, areaIds []int) error {
	total, err := repository.countEmployeeTerritoryRows(
		`SELECT COUNT(DISTINCT distributor_id)
		 FROM mst.m_distributor
		 WHERE is_del = false
		 AND distributor_id IN (?)
		 AND area_id IN (?)
		 AND (cust_id = ? OR cust_id LIKE ? OR parent_cust_id = ?)`,
		distributorIds,
		areaIds,
		custId,
		custId+"%",
		custId,
	)
	if err != nil {
		return err
	}
	if total != len(distributorIds) {
		return errors.New("distributor_ids must belong to selected area_ids")
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) countEmployeeTerritoryRows(query string, args ...interface{}) (int, error) {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = repository.Rebind(query)

	total := 0
	if err = repository.Get(&total, query, args...); err != nil {
		return 0, err
	}

	return total, nil
}

func (repository *EmployeeRepositoryImpl) Store(employee model.Employee) (int, error) {
	query :=
		`INSERT INTO mst.m_employee(
			cust_id, emp_code, emp_name, address, 
			emp_type_id, emp_grp_id, work_date,
			last_education, dob, phone_no, wa_no,
			email, image_url, is_active, created_by, created_at,
			updated_by, updated_at, is_del, deleted_by,
			deleted_at, province_id, city_id, sub_district_id, ward_id,identity_no,is_wa_no,post_code,division_id,
			region_scope, area_scope, distributor_scope)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28,
			$29, $30, $31, $32
		) RETURNING emp_id;`
	tx, err := repository.Beginx()
	if err != nil {
		log.Println("EmployeeRepository, Store, begin trx err:", err.Error())
		return employee.EmployeeId, err
	}
	defer tx.Rollback()

	lastInsertId := 0
	err = tx.QueryRow(query,
		employee.CustId, employee.EmployeeCode, employee.EmployeeName, employee.Address,
		employee.EmpTypeId, employee.EmpGrpId, employee.WorkDate,
		employee.LastEducation, employee.Dob, employee.PhoneNo, employee.WaNo,
		employee.Email, employee.ImageUrl, employee.IsActive, employee.CreatedBy, employee.CreatedAt,
		employee.UpdatedBy, employee.UpdatedAt, employee.IsDel, employee.DeletedBy,
		employee.DeletedAt, employee.ProvinceId, employee.CityId, employee.SubDistrictId, employee.WardId, employee.IdentityNo, employee.IsWaNo, employee.PostCode, employee.DivisionId,
		employeeTerritoryScopeOrDefault(employee.RegionScope), employeeTerritoryScopeOrDefault(employee.AreaScope), employeeTerritoryScopeOrDefault(employee.DistributorScope)).Scan(&lastInsertId)
	if err != nil {
		log.Println("EmployeeRepository, Store, err:", err.Error())
		return employee.EmployeeId, err
	}

	if err = repository.insertEmployeeTerritoryMappings(tx, employee.CustId, lastInsertId, employee.CreatedBy, employee.RegionIds, employee.AreaIds, employee.DistributorIds); err != nil {
		return employee.EmployeeId, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("EmployeeRepository, Store, commit trx err:", err.Error())
		return employee.EmployeeId, err
	}

	return lastInsertId, nil
}

func employeeTerritoryScopeOrDefault(scope string) string {
	scope = strings.ToUpper(strings.TrimSpace(scope))
	if scope == "SELECTED" {
		return scope
	}
	return "ALL"
}

func (repository *EmployeeRepositoryImpl) insertEmployeeTerritoryMappings(tx *sqlx.Tx, custId string, employeeId int, createdBy *int64, regionIds []int, areaIds []int, distributorIds []int) error {
	for _, regionId := range regionIds {
		if _, err := tx.Exec(
			`INSERT INTO mst.m_employee_region_mapping (cust_id, emp_id, region_id, created_by)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING`,
			custId,
			employeeId,
			regionId,
			createdBy,
		); err != nil {
			log.Println("EmployeeRepository, insertEmployeeTerritoryMappings region, err:", err.Error())
			return err
		}
	}

	for _, areaId := range areaIds {
		if _, err := tx.Exec(
			`INSERT INTO mst.m_employee_area_mapping (cust_id, emp_id, area_id, created_by)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING`,
			custId,
			employeeId,
			areaId,
			createdBy,
		); err != nil {
			log.Println("EmployeeRepository, insertEmployeeTerritoryMappings area, err:", err.Error())
			return err
		}
	}

	for _, distributorId := range distributorIds {
		if _, err := tx.Exec(
			`INSERT INTO mst.m_employee_distributor_mapping (cust_id, emp_id, distributor_id, created_by)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT DO NOTHING`,
			custId,
			employeeId,
			distributorId,
			createdBy,
		); err != nil {
			log.Println("EmployeeRepository, insertEmployeeTerritoryMappings distributor, err:", err.Error())
			return err
		}
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) Update(employeeId int, request entity.UpdateEmployeeRequest) error {
	var (
		r            model.EmployeeUpdate
		sqlSetFields string
		nRows        int64
	)

	// requestFormat, _ := json.Marshal(request)
	// fmt.Printf("EmployeeRepository, Update, request: %s\n", requestFormat)
	// panic("test")
	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("EmployeeRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_employee
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND emp_id = :emp_id_old;`

	// log.Println("EmployeeRepository, Update, query:", query)

	sqlPatch.Args["emp_id_old"] = employeeId
	sqlPatch.Args["cust_id"] = request.CustId

	// s.hooks.reset()

	if request.TerritoryMappingProvided {
		tx, err := repository.Beginx()
		if err != nil {
			log.Println("EmployeeRepository, Update, begin trx err:", err.Error())
			return err
		}
		defer tx.Rollback()

		result, err := tx.NamedExec(query, sqlPatch.Args)
		if err != nil {
			log.Println("EmployeeRepository, Update, err:", err.Error())
			return err
		}

		if nRows, err = result.RowsAffected(); err != nil {
			return errors.New("no rows affected")
		}
		if nRows == 0 {
			return errors.New("no rows affected")
		}

		if err = repository.replaceEmployeeTerritoryMappings(tx, request.CustId, employeeId, request.UpdatedBy, request.RegionIds, request.AreaIds, request.DistributorIds); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			log.Println("EmployeeRepository, Update, commit trx err:", err.Error())
			return err
		}

		return nil
	}

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("EmployeeRepository, Update, err:", err.Error())
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

func (repository *EmployeeRepositoryImpl) replaceEmployeeTerritoryMappings(tx *sqlx.Tx, custId string, employeeId int, updatedBy int64, regionIds []int, areaIds []int, distributorIds []int) error {
	if err := repository.softDeleteEmployeeTerritoryMappings(tx, custId, employeeId, updatedBy); err != nil {
		return err
	}
	return repository.insertEmployeeTerritoryMappings(tx, custId, employeeId, &updatedBy, regionIds, areaIds, distributorIds)
}

func (repository *EmployeeRepositoryImpl) softDeleteEmployeeTerritoryMappings(tx *sqlx.Tx, custId string, employeeId int, updatedBy int64) error {
	statements := []string{
		`UPDATE mst.m_employee_region_mapping
		 SET is_del = true,
		     deleted_by = $3,
		     deleted_at = CURRENT_TIMESTAMP,
		     updated_by = $3,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE cust_id = $1
		   AND emp_id = $2
		   AND is_del = false`,
		`UPDATE mst.m_employee_area_mapping
		 SET is_del = true,
		     deleted_by = $3,
		     deleted_at = CURRENT_TIMESTAMP,
		     updated_by = $3,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE cust_id = $1
		   AND emp_id = $2
		   AND is_del = false`,
		`UPDATE mst.m_employee_distributor_mapping
		 SET is_del = true,
		     deleted_by = $3,
		     deleted_at = CURRENT_TIMESTAMP,
		     updated_by = $3,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE cust_id = $1
		   AND emp_id = $2
		   AND is_del = false`,
	}

	for _, statement := range statements {
		if _, err := tx.Exec(statement, custId, employeeId, updatedBy); err != nil {
			log.Println("EmployeeRepository, softDeleteEmployeeTerritoryMappings, err:", err.Error())
			return err
		}
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) ImportUpdate(employeeId int, request entity.UpdateEmployeeRequest) error {
	query := `
		UPDATE mst.m_employee
		SET 
			emp_name       = COALESCE(:emp_name, emp_name),
			address        = COALESCE(:address, address),
			last_education = COALESCE(:last_education, last_education),
			phone_no       = COALESCE(:phone_no, phone_no),
			wa_no          = COALESCE(:wa_no, wa_no),
			email          = COALESCE(:email, email),
			emp_type_id    = COALESCE(:emp_type_id, emp_type_id),
			identity_no    = COALESCE(:identity_no, identity_no),
			post_code      = COALESCE(:post_code, post_code),
			image_url      = COALESCE(:image_url, image_url),
			province_id    = COALESCE(:province_id, province_id),
			city_id        = COALESCE(:city_id, city_id),
			sub_district_id= COALESCE(:sub_district_id, sub_district_id),
			ward_id        = COALESCE(:ward_id, ward_id),
			work_date      = COALESCE(:work_date, work_date),
			dob            = COALESCE(:dob, dob),
			emp_grp_id     = COALESCE(:emp_grp_id, emp_grp_id),
			division_id    = COALESCE(:division_id, division_id),
			is_active      = COALESCE(:is_active, is_active),
			is_wa_no       = COALESCE(:is_wa_no, is_wa_no),
			updated_by     = :updated_by,
			updated_at     = CURRENT_TIMESTAMP
		WHERE emp_id = :emp_id
		AND cust_id = :cust_id
		AND is_del = false
	`

	args := map[string]interface{}{
		"emp_id":          employeeId,
		"cust_id":         request.CustId,
		"emp_name":        repository.nullIfEmpty(request.EmployeeName),
		"address":         repository.nullIfEmpty(request.Address),
		"last_education":  repository.nullIfEmpty(request.LastEducation),
		"phone_no":        repository.nullIfEmpty(request.PhoneNo),
		"wa_no":           repository.nullIfEmpty(request.WaNo),
		"email":           repository.nullIfEmpty(request.Email),
		"emp_type_id":     repository.nullIfEmpty(request.EmpTypeId),
		"identity_no":     repository.nullIfEmpty(request.IdentityNo),
		"post_code":       repository.nullIfEmpty(request.PostCode),
		"image_url":       repository.nullIfEmpty(request.ImageUrl),
		"province_id":     repository.nullIfEmpty(request.ProvinceId),
		"city_id":         repository.nullIfEmpty(request.CityId),
		"sub_district_id": repository.nullIfEmpty(request.SubDistrictId),
		"ward_id":         repository.nullIfEmpty(request.WardId),
		"work_date":       repository.nullIfEmpty(request.WorkDate),
		"dob":             repository.nullIfEmpty(request.Dob),
		"emp_grp_id":      request.EmpGrpId,
		"division_id":     request.DivisionId,
		"is_active":       request.IsActive,
		"is_wa_no":        request.IsWaNo,
		"updated_by":      request.UpdatedBy,
	}

	result, err := repository.NamedExec(query, args)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (repository *EmployeeRepositoryImpl) nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func (repository *EmployeeRepositoryImpl) UpdateSalesman(employeeId int, isActive bool) error {

	query := `UPDATE mst.m_salesman
			  SET is_active = :is_active,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE emp_id = :emp_id;`

	// log.Println("EmployeeRepository, Update, query:", query)

	wMap := map[string]interface{}{
		"emp_id":    employeeId,
		"is_active": isActive,
	}
	_, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("EmployeeRepository, Update Salesman, err:", err.Error())
		return err
	}

	return nil
}

func (repository *EmployeeRepositoryImpl) Delete(custId string, employeeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_employee
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND emp_id = :emp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"emp_id":     employeeId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("EmployeeRepository, Delete, err:", err.Error())
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

func (repository *EmployeeRepositoryImpl) CreateImportHistory(uploadType, fileName, custId string, uploadedBy int64, totalData int) (int64, error) {
	var historyId int64
	query := `
		INSERT INTO import.import_history (file_name, uploaded_by, total_data, upload_type, cust_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING history_id
	`
	if err := repository.DB.QueryRow(query, fileName, uploadedBy, totalData, uploadType, custId).Scan(&historyId); err != nil {
		return 0, err
	}
	return historyId, nil
}

func (repository *EmployeeRepositoryImpl) UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error {
	query := `
		UPDATE import.import_history
		SET successful_data = $1,
			failed_data = $2,
			status_reupload = $3
		WHERE history_id = $4
	`
	_, err := repository.DB.Exec(query, success, failed, statusReupload, historyId)
	return err
}

func (repository *EmployeeRepositoryImpl) InsertEmployeeTemp(temp entity.ImportEmployeeTemp) error {
	query := `
		INSERT INTO import.employee_temp (
			history_id, cust_id,
			emp_code, emp_name, address,
			emp_type_name,
			emp_grp_name,
			division_name,
			work_date, last_education, dob,
			phone_no, wa_no, email,
			is_active,
			image_url,
			identity_no, is_wa_no,
			province,
			city,
			sub_district,
			ward,
			post_code,
			status_insert, error_message,
			created_at
		) VALUES (
			:history_id, :cust_id,
			:emp_code, :emp_name, :address,
			:emp_type_name,
			:emp_grp_name,
			:division_name,
			CAST(NULLIF(:work_date, '') AS DATE), :last_education, CAST(NULLIF(:dob, '') AS DATE),
			:phone_no, :wa_no, :email,
			:is_active,
			:image_url,
			:identity_no, :is_wa_no,
			:province,
			:city,
			:sub_district,
			:ward,
			:post_code,
			:status_insert, :error_message,
			NOW()
		)
	`
	_, err := repository.DB.NamedExec(query, temp)
	return err
}

func (repository *EmployeeRepositoryImpl) InsertEmployeeUpdateTemp(temp entity.ImportEmployeeUpdateTemp) error {
	query := `
		INSERT INTO import.employee_update_temp (
			history_id, cust_id, emp_id,
			emp_code, emp_name, address,
			emp_type_name,
			emp_grp_name,
			division_name,
			work_date, last_education, dob,
			phone_no, wa_no, email,
			is_active, image_url,
			identity_no, is_wa_no,
			province,
			city,
			sub_district,
			ward,
			post_code,
			status_insert, error_message,
			created_at
		) VALUES (
			:history_id, :cust_id, :emp_id,
			:emp_code, :emp_name, :address,
			:emp_type_name,
			:emp_grp_name,
			:division_name,
			CAST(NULLIF(:work_date, '') AS DATE), :last_education, CAST(NULLIF(:dob, '') AS DATE),
			:phone_no, :wa_no, :email,
			:is_active,
			:image_url,
			:identity_no, :is_wa_no,
			:province,
			:city,
			:sub_district,
			:ward,
			:post_code,
			:status_insert, :error_message,
			NOW()
		)
	`
	_, err := repository.DB.NamedExec(query, temp)
	return err
}

func (repository *EmployeeRepositoryImpl) GetImportTotalData(historyId int64) (int, error) {
	var total int
	if err := repository.DB.Get(&total, `SELECT total_data FROM import.import_history WHERE history_id = $1`, historyId); err != nil {
		return 0, err
	}
	return total, nil
}

func (repository *EmployeeRepositoryImpl) GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error) {
	rows := []entity.ImportInstruction{}
	query := `SELECT instruction_id, instruction_type, kolom, mandatory, keterangan, step
			FROM import.import_instructions WHERE instruction_type = $1 
			ORDER BY 
			CASE 
				WHEN step ILIKE 'Step 1%' THEN 1
				WHEN step ILIKE 'Step 2%' THEN 2
				WHEN step ILIKE 'Step 3%' THEN 3
				ELSE 4
			END, instruction_id;`
	if err := repository.DB.Select(&rows, query, instructionType); err != nil {
		return nil, err
	}
	return rows, nil
}

func (repository *EmployeeRepositoryImpl) CheckEmployeeExists(custId, empCode string) (bool, error) {
	if strings.TrimSpace(empCode) == "" {
		return false, nil
	}

	query := `
		SELECT COUNT(1)
		FROM mst.m_employee
		WHERE cust_id = $1 AND emp_code = $2 AND deleted_at IS NULL
	`
	var count int
	err := repository.DB.QueryRow(query, custId, empCode).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa duplikasi karyawan: %v", err)
	}
	return count > 0, nil
}

func buildEmployeePJPSortClause(sort string) string {
	allowedColumns := map[string]string{
		"emp_id":   "me.emp_id",
		"emp_code": "me.emp_code",
		"emp_name": "me.emp_name",
	}

	sortParts := []string{}
	for _, row := range strings.Split(sort, ",") {
		colSort := strings.SplitN(strings.TrimSpace(row), ":", 2)
		if len(colSort) != 2 {
			continue
		}

		column, ok := allowedColumns[strings.ToLower(strings.TrimSpace(colSort[0]))]
		if !ok {
			continue
		}

		direction := strings.ToLower(strings.TrimSpace(colSort[1]))
		if direction != "asc" && direction != "desc" {
			continue
		}

		sortParts = append(sortParts, fmt.Sprintf(`%s %s`, column, direction))
	}

	if len(sortParts) == 0 {
		return `me.emp_id DESC`
	}

	return strings.Join(sortParts, ", ")
}

func buildEmployeePJPQuery(dataFilter entity.EmployeePJPQueryFilter) (string, []interface{}, string, []interface{}) {
	selectCount := ` COUNT(*) AS total `
	selectField := ` me.emp_id, me.emp_code, me.emp_name `
	qFrom := ` FROM mst.m_employee me `

	whereClauses := []string{
		`me.is_del = false`,
		`EXISTS (
			SELECT 1
			FROM sys.m_user mu
			JOIN sys.user_roles ur
				ON ur.user_id = mu.user_id
				AND ur.cust_id = mu.cust_id
			JOIN sys.m_role mr
				ON mr.role_id = ur.role_id
				AND mr.cust_id = ur.cust_id
			WHERE mu.emp_id = me.emp_id
				AND mu.cust_id = me.cust_id
				AND LOWER(mr.role_name) = 'salesman'
		)`,
	}
	queryArgs := []interface{}{}

	if dataFilter.DistributorId != nil && *dataFilter.DistributorId > 0 {
		scopeParent := strings.TrimSpace(dataFilter.ParentCustId)
		if scopeParent == "" {
			scopeParent = dataFilter.CustId
		}

		whereClauses = append(whereClauses, `me.cust_id = (
			SELECT mc.cust_id
			FROM smc.m_customer mc
			WHERE mc.distributor_id = ?
				AND mc.parent_cust_id = ?
			LIMIT 1
		)`)
		queryArgs = append(queryArgs, *dataFilter.DistributorId, scopeParent)
	} else if dataFilter.FilterCustId != nil && strings.TrimSpace(*dataFilter.FilterCustId) != "" {
		whereClauses = append(whereClauses, `me.cust_id = ?`)
		queryArgs = append(queryArgs, strings.TrimSpace(*dataFilter.FilterCustId))
	} else {
		whereClauses = append(whereClauses, `me.cust_id = ?`)
		queryArgs = append(queryArgs, dataFilter.CustId)
	}

	if len(dataFilter.IsActive) > 0 {
		for _, status := range dataFilter.IsActive {
			if status == 1 {
				whereClauses = append(whereClauses, `me.is_active = ?`)
				queryArgs = append(queryArgs, true)
				break
			}
			if status == 0 {
				whereClauses = append(whereClauses, `me.is_active = ?`)
				queryArgs = append(queryArgs, false)
				break
			}
		}
	}

	if query := strings.TrimSpace(dataFilter.Query); query != "" {
		whereClauses = append(whereClauses, `(me.emp_code ILIKE ? OR me.emp_name ILIKE ?)`)
		queryLike := "%" + query + "%"
		queryArgs = append(queryArgs, queryLike, queryLike)
	}

	qWhere := ` WHERE ` + strings.Join(whereClauses, ` AND `)
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere + ` ORDER BY ` + buildEmployeePJPSortClause(dataFilter.Sort)

	limit := dataFilter.Limit
	if limit == 0 {
		limit = 9999
	}
	if limit > 9999 {
		limit = 9999
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	countArgs := append([]interface{}{}, queryArgs...)
	selectArgs := append([]interface{}{}, queryArgs...)

	return queryCount, countArgs, querySelect, selectArgs
}

func (repository *EmployeeRepositoryImpl) FindAllForPJP(dataFilter entity.EmployeePJPQueryFilter) ([]model.EmployeePJP, int, int, error) {
	employees := []model.EmployeePJP{}
	queryCount, countArgs, querySelect, selectArgs := buildEmployeePJPQuery(dataFilter)
	queryCount = repository.Rebind(queryCount)
	querySelect = repository.Rebind(querySelect)

	var total int
	err := repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		return employees, 0, 0, err
	}

	limit := dataFilter.Limit
	if limit == 0 {
		limit = 9999
	}
	if limit > 9999 {
		limit = 9999
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	err = repository.Select(&employees, querySelect, selectArgs...)
	if err != nil {
		return employees, total, lastPage, err
	}

	return employees, total, lastPage, nil
}

func (repository *EmployeeRepositoryImpl) GetEmployeeDataForTemplateUpdate(custId string, fields []string) (map[string][][]string, error) {
	result := make(map[string][][]string)
	if len(fields) == 0 {
		return result, nil
	}

	fieldSet := map[string]struct{}{}
	for _, f := range fields {
		fieldSet[strings.ToLower(strings.TrimSpace(f))] = struct{}{}
	}

	// Base employee data
	if _, ok := fieldSet["employee"]; ok {
		rows, err := repository.fetchData(`
			SELECT
				COALESCE(emp.emp_code, '') AS emp_code,
				COALESCE(emp.emp_name, '') AS emp_name,
				COALESCE(emp.identity_no, '') AS identity_no,
				COALESCE(emp.email, '') AS email,
				COALESCE(emp.phone_no, '') AS phone_no,
				COALESCE(p.province, '') AS province,
				COALESCE(r.regency, '') AS city,
				COALESCE(sd.sub_district, '') AS sub_district,
				COALESCE(w.ward, '') AS ward,
				COALESCE(emp.post_code, '') AS post_code,
				COALESCE(emp.address, '') AS address,
				COALESCE(eg.emp_grp_name, '') AS emp_grp_name,
				COALESCE(d.division_name, '') AS division_name,
				COALESCE(emp.image_url, '') AS image_url,
				COALESCE(TO_CHAR(emp.dob, 'YYYY-MM-DD'), '') AS dob,
				COALESCE(emp.last_education, '') AS last_education,
				COALESCE(TO_CHAR(emp.work_date, 'YYYY-MM-DD'), '') AS work_date,
				COALESCE(emp.wa_no, '') AS wa_no
			FROM mst.m_employee emp
			LEFT JOIN mst.m_emp_group eg ON eg.cust_id = emp.cust_id AND eg.emp_grp_id = emp.emp_grp_id
			LEFT JOIN mst.m_division d ON d.cust_id = emp.cust_id AND d.division_id = emp.division_id
			LEFT JOIN mst.m_province p ON p.province_id = emp.province_id AND p.cust_id = emp.cust_id
			LEFT JOIN mst.m_regency r ON r.regency_id = emp.city_id AND r.cust_id = emp.cust_id
			LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = emp.sub_district_id AND sd.cust_id = emp.cust_id
			LEFT JOIN mst.m_ward w ON w.ward_id = emp.ward_id AND w.cust_id = emp.cust_id
			WHERE emp.cust_id = $1 AND emp.is_del = false
			ORDER BY emp.emp_code
		`, custId)
		if err != nil {
			return nil, err
		}
		result["employee"] = rows
	}

	if _, ok := fieldSet["emp_group"]; ok {
		rows, err := repository.fetchData(`
			SELECT emp_grp_id, emp_grp_code, emp_grp_name
			FROM mst.m_emp_group
			WHERE cust_id = $1
			ORDER BY emp_grp_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["emp_group"] = rows
	}

	if _, ok := fieldSet["emp_type"]; ok {
		rows, err := repository.fetchData(`
			SELECT emp_type_id, emp_type_name
			FROM mst.m_emp_type
			WHERE cust_id = $1
			ORDER BY emp_type_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["emp_type"] = rows
	}

	if _, ok := fieldSet["division"]; ok {
		rows, err := repository.fetchData(`
			SELECT division_id, division_code, division_name
			FROM mst.m_division
			WHERE cust_id = $1
			ORDER BY division_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["division"] = rows
	}

	if _, ok := fieldSet["province"]; ok {
		rows, err := repository.fetchData(`
			SELECT province_id, province
			FROM mst.m_province
			WHERE cust_id = $1
			ORDER BY province_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["province"] = rows
	}

	if _, ok := fieldSet["city"]; ok {
		rows, err := repository.fetchData(`
			SELECT regency_id AS city_id, regency AS city, province_id
			FROM mst.m_regency
			WHERE cust_id = $1
			ORDER BY regency_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["city"] = rows
	}

	if _, ok := fieldSet["sub_district"]; ok {
		rows, err := repository.fetchData(`
			SELECT sub_district_id, sub_district, regency_id AS city_id
			FROM mst.m_sub_district
			WHERE cust_id = $1
			ORDER BY sub_district_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["sub_district"] = rows
	}

	if _, ok := fieldSet["ward"]; ok {
		rows, err := repository.fetchData(`
			SELECT ward_id, ward, sub_district_id
			FROM mst.m_ward
			WHERE cust_id = $1
			ORDER BY ward_id
		`, custId)
		if err != nil {
			return nil, err
		}
		result["ward"] = rows
	}

	return result, nil
}

func (repository *EmployeeRepositoryImpl) fetchData(query string, args ...interface{}) ([][]string, error) {
	rows, err := repository.DB.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result [][]string
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		rowStrings := make([]string, len(columns))
		for i, v := range values {
			if v == nil {
				rowStrings[i] = ""
				continue
			}
			switch val := v.(type) {
			case []byte:
				rowStrings[i] = string(val)
			default:
				rowStrings[i] = fmt.Sprint(val)
			}
		}
		result = append(result, rowStrings)
	}

	return result, rows.Err()
}

func (repository *EmployeeRepositoryImpl) FindEmpGroupIdByName(custId string, name string) (int, error) {
	var id int
	query := `SELECT emp_grp_id FROM mst.m_emp_group WHERE cust_id = $1 AND LOWER(emp_grp_name) = LOWER($2) LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return 0, err
	}
	return id, nil
}

func (repository *EmployeeRepositoryImpl) FindDivisionIdByName(custId string, name string) (int, error) {
	var id int
	query := `SELECT division_id FROM mst.m_division WHERE cust_id = $1 AND LOWER(division_name) = LOWER($2) LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return 0, err
	}
	return id, nil
}

func (repository *EmployeeRepositoryImpl) FindProvinceIdByName(custId string, name string) (string, error) {
	var id string
	query := `
		SELECT province_id
		FROM mst.m_province
		WHERE cust_id = $1
			AND LOWER(TRIM(province)) = LOWER(TRIM($2))
		LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return "", err
	}
	return id, nil
}

func (repository *EmployeeRepositoryImpl) FindCityIdByName(custId string, name string) (string, error) {
	var id string
	query := `
		SELECT regency_id
		FROM mst.m_regency
		WHERE cust_id = $1
			AND LOWER(TRIM(regency)) = LOWER(TRIM($2))
		LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return "", err
	}
	return id, nil
}

func (repository *EmployeeRepositoryImpl) FindSubDistrictIdByName(custId string, name string) (string, error) {
	var id string
	query := `
		SELECT sub_district_id
		FROM mst.m_sub_district
		WHERE cust_id = $1
			AND LOWER(TRIM(sub_district)) = LOWER(TRIM($2))
		LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return "", err
	}
	return id, nil
}

func (repository *EmployeeRepositoryImpl) FindWardIdByName(custId string, name string) (string, error) {
	var id string
	query := `
		SELECT ward_id
		FROM mst.m_ward
		WHERE cust_id = $1
			AND LOWER(TRIM(ward)) = LOWER(TRIM($2))
		LIMIT 1`
	if err := repository.DB.Get(&id, query, custId, strings.TrimSpace(name)); err != nil {
		return "", err
	}
	return id, nil
}
