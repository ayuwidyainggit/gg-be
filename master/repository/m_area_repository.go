package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AreaRepository interface {
	FindOneByAreaIdAndCustId(areaId int, custId string) (model.AreaList, error)
	FindOneByAreaCodeAndCustId(areaCode string, custId string) (model.Area, error)
	FindAllByCustId(dataFilter entity.AreaQueryFilter) (area []model.AreaList, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.AreaQueryFilter) (area []model.AreaList, total int, lastPage int, err error)
	Store(area model.Area) (int, error)
	Update(areaId int, request entity.UpdateAreaRequest) error
	Delete(custId string, areaId int, deletedBy int64) error
}

func NewAreaRepository(db *sqlx.DB) AreaRepository {
	return &areaRepositoryImpl{db}
}

type areaRepositoryImpl struct {
	*sqlx.DB
}

func (repository *areaRepositoryImpl) FindOneByAreaIdAndCustId(areaId int, custId string) (model.AreaList, error) {
	area := model.AreaList{}
	query := `SELECT
				ma.cust_id, ma.area_id, ma.area_code, ma.area_name, ma.is_active, ma.updated_by, ma.updated_at,  
				mr.region_id, mr.region_code, mr.region_name,  
				mo.official_id, mo.official_type,
				moh.hierarchy_code as official_hierarchy_code, me.emp_name as official_emp_name,
				u.user_fullname AS updated_by_name, ma.created_by,
				ma.created_at, ma.updated_by, ma.updated_at,
				ma.is_del, ma.deleted_by, ma.deleted_at
			from mst.m_area ma
				left join mst.m_official mo on mo.official_id = ma.official_id 
				left join mst.m_employee me on me.emp_id = mo.emp_id 
				left join mst.m_official_hierarchy moh on moh.official_type = mo.official_type 
				left join mst.m_region mr on mr.region_id = ma.region_id 
				left join sys.m_user u ON u.user_id = ma.updated_by 
			WHERE ma.area_id = $1
			AND ma.cust_id = $2`
	err := repository.Get(&area, query, areaId, custId)
	if err != nil {
		log.Println("areaRepository, FindOneByBankCodeAndCustId, err:", err.Error())
		return area, err
	}

	return area, nil
}

func (repository *areaRepositoryImpl) FindOneByAreaCodeAndCustId(areaCode string, custId string) (model.Area, error) {
	area := model.Area{}
	query := `SELECT
				cust_id, area_id, area_code,
				area_name, region_id, official_id, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_area
			  WHERE area_code = $1
			  AND cust_id = $2 and is_del = false`
	err := repository.Get(&area, query, areaCode, custId)
	if err != nil {
		log.Println("areaRepository, FindOneByAreaCodeAndCustId, err:", err.Error())
		return area, err
	}

	return area, nil
}

func (repository *areaRepositoryImpl) Store(area model.Area) (int, error) {
	query :=
		`INSERT INTO mst.m_area(
			cust_id, area_code, area_name, region_id, official_id, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12, $13
		) RETURNING area_id;`
	lastInsertId := area.AreaId
	err := repository.QueryRow(query,
		area.CustId, area.AreaCode, area.AreaName, area.RegionId, area.OfficialId,
		area.IsActive, area.CreatedBy, area.CreatedAt, area.UpdatedBy,
		area.UpdatedAt, area.IsDel, area.DeletedBy, area.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("areaRepository, Store, err:", err.Error())
		return area.AreaId, err
	}
	return area.AreaId, nil
}

func (repository *areaRepositoryImpl) FindAllByCustId(dataFilter entity.AreaQueryFilter) ([]model.AreaList, int, int, error) {
	areas := []model.AreaList{}
	queryCount, countArgs, querySelect, selectArgs := buildAreaListQuery(dataFilter, false)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return areas, 0, 0, err
	}
	queryCount = repository.Rebind(queryCount)

	var total int
	err = repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Println("areaRepository, count total, err:", err.Error())
		return areas, 0, 0, err
	}

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return areas, total, 0, err
	}
	querySelect = repository.Rebind(querySelect)

	lastPage := calculateLastPage(total, dataFilter.Limit)
	err = repository.Select(&areas, querySelect, selectArgs...)
	if err != nil {
		log.Println("areaRepository, FindAllByCustId, err:", err.Error())
		return areas, total, lastPage, err
	}

	return areas, total, lastPage, nil
}

func (repository *areaRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.AreaQueryFilter) ([]model.AreaList, int, int, error) {
	areas := []model.AreaList{}
	queryCount, countArgs, querySelect, selectArgs := buildAreaListQuery(dataFilter, true)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return areas, 0, 0, err
	}
	queryCount = repository.Rebind(queryCount)

	var total int
	err = repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Println("areaRepository, count total, err:", err.Error())
		return areas, 0, 0, err
	}

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return areas, total, 0, err
	}
	querySelect = repository.Rebind(querySelect)

	lastPage := calculateLastPage(total, dataFilter.Limit)
	err = repository.Select(&areas, querySelect, selectArgs...)
	if err != nil {
		log.Println("areaRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return areas, total, lastPage, err
	}

	return areas, total, lastPage, nil
}

func buildAreaListQuery(dataFilter entity.AreaQueryFilter, lookup bool) (string, []interface{}, string, []interface{}) {
	selectCount := `COUNT(DISTINCT ma.area_id) AS total`
	selectField := `DISTINCT ma.cust_id, ma.area_id, ma.area_code, ma.area_name, ma.is_active, ma.updated_by, ma.updated_at,
	mr.region_id, mr.region_code, mr.region_name,
	mo.official_id, mo.official_type,
	moh.hierarchy_code as official_hierarchy_code, me.emp_name as official_emp_name,
	u.user_fullname AS updated_by_name`

	qFrom := ` FROM mst.m_area ma
	LEFT JOIN mst.m_official mo on mo.official_id = ma.official_id
	LEFT JOIN mst.m_employee me on me.emp_id = mo.emp_id
	LEFT JOIN mst.m_official_hierarchy moh on moh.official_type = mo.official_type AND moh.cust_id = ?
	LEFT JOIN mst.m_region mr on mr.region_id = ma.region_id AND mr.cust_id = ?
	LEFT JOIN sys.m_user u ON u.user_id = ma.updated_by`
	args := []interface{}{dataFilter.CustId, dataFilter.ParentCustId}
	whereClauses := []string{"ma.is_del = false"}

	isPrincipal := dataFilter.DistributorId == 0 && dataFilter.EmployeeId > 0
	if isPrincipal && dataFilter.Scope.AreaScope == "specific" {
		qFrom += ` INNER JOIN mst.m_employee_area_mapping eam
			ON eam.area_id = ma.area_id
			AND eam.cust_id = ?
			AND eam.emp_id = ?
			AND eam.is_del = false`
		args = append(args, dataFilter.CustId, dataFilter.EmployeeId)
		whereClauses = append(whereClauses, "ma.cust_id = ?")
		args = append(args, dataFilter.CustId)
	} else if isPrincipal && dataFilter.Scope.RegionScope == "specific" {
		qFrom += ` INNER JOIN mst.m_employee_region_mapping erm
			ON erm.region_id = ma.region_id
			AND erm.cust_id = ?
			AND erm.emp_id = ?
			AND erm.is_del = false`
		args = append(args, dataFilter.CustId, dataFilter.EmployeeId)
		whereClauses = append(whereClauses, "ma.cust_id = ?")
		args = append(args, dataFilter.CustId)
	} else {
		scopeCustID := dataFilter.CustId
		if lookup && !isPrincipal {
			scopeCustID = dataFilter.ParentCustId
		}
		whereClauses = append(whereClauses, "ma.cust_id = ?")
		args = append(args, scopeCustID)
	}

	if lookup {
		whereClauses = append(whereClauses, "ma.is_active = true")
	}
	if dataFilter.Query != "" {
		queryLike := "%" + dataFilter.Query + "%"
		whereClauses = append(whereClauses, "(ma.area_code ILIKE ? OR ma.area_name ILIKE ?)")
		args = append(args, queryLike, queryLike)
	}
	if len(dataFilter.RegionId) > 0 {
		whereClauses = append(whereClauses, "ma.region_id IN (?)")
		args = append(args, dataFilter.RegionId)
	}
	if len(dataFilter.AreaId) > 0 {
		whereClauses = append(whereClauses, "ma.area_id IN (?)")
		args = append(args, dataFilter.AreaId)
	}
	if !lookup && dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			whereClauses = append(whereClauses, "ma.is_active = ?")
			args = append(args, true)
		}
		if *dataFilter.IsActive == 2 {
			whereClauses = append(whereClauses, "ma.is_active = ?")
			args = append(args, false)
		}
	}
	if lookup && dataFilter.IsActive != nil && *dataFilter.IsActive == 2 {
		whereClauses = append(whereClauses, "ma.is_active = ?")
		args = append(args, false)
	}

	qWhere := ` WHERE ` + strings.Join(whereClauses, ` AND `)
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere
	querySelect += buildAreaOrderAndPagination(dataFilter)
	return queryCount, args, querySelect, args
}

func buildAreaOrderAndPagination(dataFilter entity.AreaQueryFilter) string {
	sortBy := ""
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
	}
	if sortBy == "" {
		sortBy = `ma.area_id DESC`
	}

	limit := dataFilter.Limit
	if limit == 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	return fmt.Sprintf(` ORDER BY %s LIMIT %d OFFSET %d`, sortBy, limit, offset)
}

func (repository *areaRepositoryImpl) Update(areaId int, request entity.UpdateAreaRequest) error {
	var (
		r            model.AreaUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("areaRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_area
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false
			  AND cust_id = :cust_id
			  AND area_id = :area_id_old;`

	log.Println("areaRepository, Update, query:", query)

	sqlPatch.Args["area_id_old"] = areaId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("areaRepository, Update, err:", err.Error())
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

func (repository *areaRepositoryImpl) Delete(custId string, areaId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_area
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by
			WHERE is_del = false
			AND cust_id = :cust_id
			AND area_id = :area_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"area_id":    areaId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("areaRepository, Delete, err:", err.Error())
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
