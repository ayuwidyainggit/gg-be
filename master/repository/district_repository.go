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

type DistrictRepository interface {
	FindOneByDistrictIdAndCustId(districtId int, custId string) (model.District, error)
	FindOneByDistrictCodeAndCustId(districtCode string, custId string) (model.District, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.District, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.District, total int, lastPage int, err error)
	Store(district model.District) (int, error)
	Update(districtId int, request entity.UpdateDistrictRequest) error
	Delete(custId string, districtId int, deletedBy int64) error
}

func NewDistrictRepository(db *sqlx.DB) DistrictRepository {
	return &districtRepositoryImpl{db}
}

type districtRepositoryImpl struct {
	*sqlx.DB
}

func (repository *districtRepositoryImpl) FindOneByDistrictIdAndCustId(districtId int, custId string) (model.District, error) {
	district := model.District{}
	query := `SELECT 
				cust_id, district_id, district_code,
				district_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_district 
			  WHERE district_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&district, query, districtId, custId)
	if err != nil {
		log.Println("districtRepository, FindOneByDistrictCodeAndCustId, err:", err.Error())
		return district, err
	}

	return district, nil
}

func (repository *districtRepositoryImpl) FindOneByDistrictCodeAndCustId(districtCode string, custId string) (model.District, error) {
	district := model.District{}
	query := `SELECT 
				cust_id, district_id, district_code,
				district_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_district 
			  WHERE district_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&district, query, districtCode, custId)
	if err != nil {
		log.Println("districtRepository, FindOneByDistrictCodeAndCustId, err:", err.Error())
		return district, err
	}

	return district, nil
}

func (repository *districtRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.District, int, int, error) {

	districts := []model.District{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.district_id, a.district_code, a.district_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.district_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.district_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_district a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("districtRepository, count total, err:", err.Error())
		return districts, 0, 0, err
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
		sortBy := `a.district_id`
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

	// log.Println("districtRepository, querySelect:", querySelect)
	err = repository.Select(&districts, querySelect)
	if err != nil {
		log.Println("districtRepository, FindAllByCustId, err:", err.Error())
		return districts, total, lastPage, err
	}

	return districts, total, lastPage, nil
}

func (repository *districtRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.District, int, int, error) {

	districts := []model.District{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.district_id, a.district_code, a.district_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.district_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.district_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_district a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("districtRepository, count total, err:", err.Error())
		return districts, 0, 0, err
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
		sortBy := `a.district_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&districts, querySelect)
	if err != nil {
		log.Println("districtRepository, FindAllByCustId, err:", err.Error())
		return districts, total, 1, err
	}

	return districts, total, 1, nil
}

func (repository *districtRepositoryImpl) Store(district model.District) (int, error) {
	query :=
		`INSERT INTO mst.m_district(
			cust_id, district_code, district_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING district_id;`
	lastInsertId := district.DistrictId
	err := repository.QueryRow(query,
		district.CustId, district.DistrictCode, district.DistrictName,
		district.IsActive, district.CreatedBy, district.CreatedAt, district.UpdatedBy,
		district.UpdatedAt, district.IsDel, district.DeletedBy, district.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("districtRepository, Store, err:", err.Error())
		return district.DistrictId, err
	}
	return district.DistrictId, nil
}

func (repository *districtRepositoryImpl) Update(districtId int, request entity.UpdateDistrictRequest) error {
	var (
		r            model.DistrictUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("districtRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_district
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND district_id = :district_id_old;`

	log.Println("districtRepository, Update, query:", query)

	sqlPatch.Args["district_id_old"] = districtId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("districtRepository, Update, err:", err.Error())
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

func (repository *districtRepositoryImpl) Delete(custId string, districtId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_district
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND district_id = :district_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"district_id": districtId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("DistrictRepository, Delete, err:", err.Error())
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
