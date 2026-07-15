package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/sql_helper"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EmployeeRepository interface {
	FindOneByEmployeeIdAndCustId(params entity.DetailEmployeeParams) (model.Employee, error)
	FindOneByEmployeeCodeAndCustId(params entity.DetailEmployeeParams) (model.Employee, error)
	FindAllByCustId(dataFilter entity.EmployeeQueryFilter) (consPro []model.Employee, total int64, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.EmployeeQueryFilter) (consPro []model.Employee, total int64, lastPage int, err error)
	Store(employee model.Employee) (int, error)
	Update(employeeId int, request entity.UpdateEmployeeRequest) error
	Delete(custId string, employeeId int, deletedBy int64) error
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
				d.division_name
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

func (repository *EmployeeRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.EmployeeQueryFilter) ([]model.Employee, int64, int, error) {

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
	var total int64
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

func (repository *EmployeeRepositoryImpl) FindAllByCustId(dataFilter entity.EmployeeQueryFilter) ([]model.Employee, int64, int, error) {

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

	qFrom := ` FROM mst.m_employee ep
			   LEFT JOIN mst.m_emp_type et ON et.emp_type_id = ep.emp_type_id AND et.cust_id = '` + dataFilter.ParentCustId + `' 
			   LEFT JOIN mst.m_province p ON p.province_id = ep.province_id AND p.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_regency r ON r.regency_id = ep.city_id AND r.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = ep.sub_district_id AND sd.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_ward w ON w.ward_id = ep.ward_id AND w.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN mst.m_emp_group eg ON eg.emp_grp_id = ep.emp_grp_id AND eg.cust_id = '` + dataFilter.ParentCustId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = ep.updated_by
			   LEFT JOIN mst.m_division d ON d.division_id = ep.division_id AND ep.cust_id = '` + dataFilter.ParentCustId + `' `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("EmployeeRepository, queryCount:", queryCount)
	var total int64
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

func (repository *EmployeeRepositoryImpl) Store(employee model.Employee) (int, error) {
	query :=
		`INSERT INTO mst.m_employee(
			cust_id, emp_code, emp_name, address, 
			emp_type_id, emp_grp_id, work_date,
			last_education, dob, phone_no, wa_no,
			email, image_url, is_active, created_by, created_at,
			updated_by, updated_at, is_del, deleted_by,
			deleted_at, province_id, city_id, sub_district_id, ward_id,identity_no,is_wa_no,post_code,division_id)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28,
			$29
		) RETURNING emp_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		employee.CustId, employee.EmployeeCode, employee.EmployeeName, employee.Address,
		employee.EmpTypeId, employee.EmpGrpId, employee.WorkDate,
		employee.LastEducation, employee.Dob, employee.PhoneNo, employee.WaNo,
		employee.Email, employee.ImageUrl, employee.IsActive, employee.CreatedBy, employee.CreatedAt,
		employee.UpdatedBy, employee.UpdatedAt, employee.IsDel, employee.DeletedBy,
		employee.DeletedAt, employee.ProvinceId, employee.CityId, employee.SubDistrictId, employee.WardId, employee.IdentityNo, employee.IsWaNo, employee.PostCode, employee.DivisionId).Scan(&lastInsertId)
	if err != nil {
		log.Println("EmployeeRepository, Store, err:", err.Error())
		return employee.EmployeeId, err
	}
	return employee.EmployeeId, nil
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
