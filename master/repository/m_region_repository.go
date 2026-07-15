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
	"strings"

	"github.com/jmoiron/sqlx"
)

type RegionRepository interface {
	FindOneByRegionIdAndCustId(regionId int, custId string) (model.Region, error)
	FindOneByRegionCodeAndCustId(regionCode string, custId string) (model.Region, error)
	FindAllByCustId(dataFilter entity.RegionQueryFilter) (region []model.Region, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.RegionQueryFilter) (region []model.Region, total int, lastPage int, err error)
	Store(Region model.Region) (int, error)
	Update(regionId int, request entity.UpdateRegionRequest) error
	Delete(custId string, regionId int, deletedBy int64) error
}

func NewRegionRepository(db *sqlx.DB) RegionRepository {
	return &regionRepositoryImpl{db}
}

type regionRepositoryImpl struct {
	*sqlx.DB
}

func (repository *regionRepositoryImpl) FindOneByRegionIdAndCustId(regionId int, custId string) (model.Region, error) {
	region := model.Region{}
	query := `SELECT 
				cust_id, region_id, region_code,
				region_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_region 
			  WHERE region_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&region, query, regionId, custId)
	if err != nil {
		log.Println("regionRepository, FindOneByBankCodeAndCustId, err:", err.Error())
		return region, err
	}

	return region, nil
}

func (repository *regionRepositoryImpl) FindOneByRegionCodeAndCustId(regionCode string, custId string) (model.Region, error) {
	region := model.Region{}
	query := `SELECT 
				cust_id, region_id, region_code,
				region_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_region 
			  WHERE region_code = $1 
			  AND cust_id = $2 and is_del = false`
	err := repository.Get(&region, query, regionCode, custId)
	if err != nil {
		log.Println("regionRepository, FindOneByRegionCodeAndCustId, err:", err.Error())
		return region, err
	}

	return region, nil
}

func (repository *regionRepositoryImpl) Store(region model.Region) (int, error) {
	query :=
		`INSERT INTO mst.m_region(
			cust_id, region_code, region_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING region_id;`
	lastInsertId := region.RegionId
	err := repository.QueryRow(query,
		region.CustId, region.RegionCode, region.RegionName,
		region.IsActive, region.CreatedBy, region.CreatedAt, region.UpdatedBy,
		region.UpdatedAt, region.IsDel, region.DeletedBy, region.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("regionRepository, Store, err:", err.Error())
		return region.RegionId, err
	}
	return region.RegionId, nil
}

func (repository *regionRepositoryImpl) FindAllByCustId(dataFilter entity.RegionQueryFilter) ([]model.Region, int, int, error) {
	regions := []model.Region{}
	queryCount, countArgs, querySelect, selectArgs := buildRegionListQuery(dataFilter, false)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return regions, 0, 0, err
	}
	queryCount = repository.Rebind(queryCount)

	var total int
	err = repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Println("regionRepository, count total, err:", err.Error())
		return regions, 0, 0, err
	}

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return regions, total, 0, err
	}
	querySelect = repository.Rebind(querySelect)

	lastPage := calculateLastPage(total, dataFilter.Limit)
	err = repository.Select(&regions, querySelect, selectArgs...)
	if err != nil {
		log.Println("regionRepository, FindAllByCustId, err:", err.Error())
		return regions, total, lastPage, err
	}

	return regions, total, lastPage, nil
}

func (repository *regionRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.RegionQueryFilter) ([]model.Region, int, int, error) {
	regions := []model.Region{}
	queryCount, countArgs, querySelect, selectArgs := buildRegionListQuery(dataFilter, true)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return regions, 0, 0, err
	}
	queryCount = repository.Rebind(queryCount)

	var total int
	err = repository.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Println("regionRepository, count total, err:", err.Error())
		return regions, 0, 0, err
	}

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return regions, total, 0, err
	}
	querySelect = repository.Rebind(querySelect)

	lastPage := calculateLastPage(total, dataFilter.Limit)
	err = repository.Select(&regions, querySelect, selectArgs...)
	if err != nil {
		log.Println("regionRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return regions, total, lastPage, err
	}

	return regions, total, lastPage, nil
}

func buildRegionListQuery(dataFilter entity.RegionQueryFilter, lookup bool) (string, []interface{}, string, []interface{}) {
	selectCount := `COUNT(DISTINCT a.region_id) AS total`
	selectField := `DISTINCT a.cust_id, a.region_id, a.region_code, a.region_name, a.is_active`
	if !lookup {
		selectField = `DISTINCT a.cust_id, a.region_id, a.region_code,
	a.region_name, a.is_active, a.created_by,
	a.created_at, a.updated_by, a.updated_at,
	a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name`
	}

	qFrom := ` FROM mst.m_region a
	LEFT JOIN sys.m_user u ON u.user_id = a.updated_by`
	whereClauses := []string{"a.is_del = false"}
	args := make([]interface{}, 0)

	isPrincipal := dataFilter.DistributorId == 0 && dataFilter.EmployeeId > 0
	if isPrincipal && dataFilter.Scope.RegionScope == "specific" {
		qFrom += ` INNER JOIN mst.m_employee_region_mapping erm
			ON erm.region_id = a.region_id
			AND erm.cust_id = ?
			AND erm.emp_id = ?
			AND erm.is_del = false`
		args = append(args, dataFilter.CustId, dataFilter.EmployeeId)
		whereClauses = append(whereClauses, "a.cust_id = ?")
		args = append(args, dataFilter.CustId)
	} else {
		scopeCustID := dataFilter.ParentCustId
		if isPrincipal {
			scopeCustID = dataFilter.CustId
		}
		whereClauses = append(whereClauses, "a.cust_id = ?")
		args = append(args, scopeCustID)
	}

	if lookup {
		whereClauses = append(whereClauses, "a.is_active = true")
	}
	if dataFilter.Query != "" {
		queryLike := "%" + dataFilter.Query + "%"
		whereClauses = append(whereClauses, "(a.region_code ILIKE ? OR a.region_name ILIKE ?)")
		args = append(args, queryLike, queryLike)
	}
	if len(dataFilter.RegionId) > 0 {
		whereClauses = append(whereClauses, "a.region_id IN (?)")
		args = append(args, dataFilter.RegionId)
	}
	if !lookup && dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			whereClauses = append(whereClauses, "a.is_active = ?")
			args = append(args, true)
		}
		if *dataFilter.IsActive == 2 {
			whereClauses = append(whereClauses, "a.is_active = ?")
			args = append(args, false)
		}
	}
	if lookup && dataFilter.IsActive != nil && *dataFilter.IsActive == 2 {
		whereClauses = append(whereClauses, "a.is_active = ?")
		args = append(args, false)
	}

	qWhere := ` WHERE ` + strings.Join(whereClauses, ` AND `)
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere
	querySelect += buildRegionOrderAndPagination(dataFilter)

	return queryCount, args, querySelect, args
}

func buildRegionOrderAndPagination(dataFilter entity.RegionQueryFilter) string {
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
		sortBy = `a.region_id DESC`
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

func calculateLastPage(total int, limit int) int {
	if limit <= 0 {
		limit = 10
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func (repository *regionRepositoryImpl) Update(regionId int, request entity.UpdateRegionRequest) error {
	var (
		r            model.RegionUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("regionRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_region
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND region_id = :region_id_old;`

	log.Println("regionRepository, Update, query:", query)

	sqlPatch.Args["region_id_old"] = regionId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("regionRepository, Update, err:", err.Error())
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

func (repository *regionRepositoryImpl) Delete(custId string, regionId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_region
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND region_id = :region_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"region_id":  regionId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("RegionRepository, Delete, err:", err.Error())
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
