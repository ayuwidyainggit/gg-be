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

type WardRepository interface {
	FindOneByWardId(wardId string, custId string) (model.Ward, error)
	Store(ward model.Ward) (string, error)
	Update(wardId string, request entity.UpdateWardRequest) error
	Delete(wardId string, deletedBy int64, custId string) error
	FindAllByCustId(dataFilter entity.WardQueryFilter, custId string) (province []model.Ward, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.WardQueryFilter, custId string) (province []model.Ward, total int, lastPage int, err error)
}

func NewWardRepository(db *sqlx.DB) WardRepository {
	return &WardRepositoryImpl{db}
}

type WardRepositoryImpl struct {
	*sqlx.DB
}

func (repository *WardRepositoryImpl) FindOneByWardId(wardId string, custId string) (model.Ward, error) {
	ward := model.Ward{}
	query := `SELECT 
				ward_id, ward, province_id, regency_id, sub_district_id, is_active, created_by,
				created_at, updated_by, updated_at
			  FROM mst.m_ward 
			  WHERE ward_id = $1 and cust_id = $2`
	err := repository.Get(&ward, query, wardId, custId)
	if err != nil {
		log.Println("WardRepository, FindOneByWardId, err:", err.Error())
		return ward, err
	}

	return ward, nil
}

func (repository *WardRepositoryImpl) FindAllByCustId(dataFilter entity.WardQueryFilter, custId string) ([]model.Ward, int, int, error) {

	wards := []model.Ward{}
	selectCount := ` COUNT(*) AS total `
	selectField := `   mw.ward_id, mw.ward, mw.province_id, mw.regency_id, mw.sub_district_id,
					msd.sub_district, 
					mp.province, mr.regency, mw.is_active, 
					mw.is_active, mw.created_by, mw.created_at, mw.updated_by, mw.updated_at,
					u.user_fullname AS updated_by_name `

	qFrom := ` 	FROM mst.m_ward mw
	LEFT JOIN mst.m_sub_district msd on msd.sub_district_id = mw.sub_district_id AND msd.cust_id = '` + custId + `'
	LEFT JOIN mst.m_province mp on mp.province_id = mw.province_id AND mp.cust_id = '` + custId + `'
	LEFT JOIN mst.m_regency mr on mr.regency_id = mw.regency_id AND mr.cust_id = '` + custId + `'
	LEFT JOIN sys.m_user u ON u.user_id = mw.updated_by   `

	qWhere := ` WHERE mw.cust_id = '` + custId + `'  `

	if dataFilter.Query != "" {
		qWhere += ` AND (mw.ward_id ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.ProvinceId != "" {
		qWhere += `AND mw.province_id = '` + dataFilter.ProvinceId + `' `
	}

	if dataFilter.RegencyId != "" {
		qWhere += `AND mw.regency_id = '` + dataFilter.RegencyId + `' `
	}

	if dataFilter.SubDistrictId != "" {
		qWhere += `AND mw.sub_district_id = '` + dataFilter.SubDistrictId + `' `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND mw.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND mw.is_active = false `
		}
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("WardRepository, count total, err:", err.Error())
		return wards, 0, 0, err
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
		sortBy := `ward_id`
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

	// log.Println("WardRepository, querySelect:", querySelect)
	err = repository.Select(&wards, querySelect)
	if err != nil {
		log.Println("WardRepository, FindAllByCustId, err:", err.Error())
		return wards, total, lastPage, err
	}

	return wards, total, lastPage, nil
}

func (repository *WardRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.WardQueryFilter, custId string) ([]model.Ward, int, int, error) {

	wards := []model.Ward{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` mw.ward_id, mw.ward, mw.province_id, mw.regency_id, mw.sub_district_id,
					msd.sub_district, mp.province, mr.regency, mw.is_active `

	qWhere := ` WHERE mw.cust_id = '` + custId + `' AND mw.is_active = true  `

	qFrom := ` 	FROM mst.m_ward mw
				LEFT JOIN mst.m_sub_district msd on msd.sub_district_id = mw.sub_district_id AND msd.cust_id = '` + custId + `'
				LEFT JOIN mst.m_province mp on mp.province_id = mw.province_id AND mp.cust_id = '` + custId + `'
				LEFT JOIN mst.m_regency mr on mr.regency_id = mw.regency_id AND mr.cust_id = '` + custId + `'
				LEFT JOIN sys.m_user u ON u.user_id = mw.updated_by   `

	if dataFilter.Query != "" {
		qWhere += ` AND (mw.ward_id ILIKE '%` + dataFilter.Query + `%' OR mw.ward ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.ProvinceId != "" {
		qWhere += `AND mw.province_id = '` + dataFilter.ProvinceId + `' `
	}

	if dataFilter.RegencyId != "" {
		qWhere += `AND mw.regency_id = '` + dataFilter.RegencyId + `' `
	}

	if dataFilter.SubDistrictId != "" {
		qWhere += `AND mw.sub_district_id = '` + dataFilter.SubDistrictId + `' `
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("WardRepository, count total, err:", err.Error())
		return wards, 0, 0, err
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
		sortBy := `ward_id`
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

	// log.Println("WardRepository, querySelect:", querySelect)
	err = repository.Select(&wards, querySelect)
	if err != nil {
		log.Println("WardRepository, FindAllByCustId, err:", err.Error())
		return wards, total, lastPage, err
	}

	return wards, total, lastPage, nil
}

func (repository *WardRepositoryImpl) Store(ward model.Ward) (string, error) {
	query :=
		`INSERT INTO mst.m_ward(
			cust_id, ward_id, ward, province_id, regency_id, sub_district_id, is_active, 
			created_by, created_at, updated_by, updated_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING ward_id;`
	lastInsertId := ward.WardId
	err := repository.QueryRow(query,
		ward.CustId, ward.WardId, ward.Ward, ward.ProvinceId, ward.RegencyId, ward.SubDistrictId,
		ward.IsActive, ward.CreatedBy, ward.CreatedAt, ward.UpdatedBy,
		ward.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("WardRepository, Store, err:", err.Error())
		return ward.WardId, err
	}
	return ward.WardId, nil
}

func (repository *WardRepositoryImpl) Update(wardId string, request entity.UpdateWardRequest) error {
	var (
		r            model.WardUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("WardRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_ward
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id_old
			  AND ward_id = :ward_id_old;`

	// log.Println("WardRepository, Update, query:", query)

	sqlPatch.Args["ward_id_old"] = wardId
	sqlPatch.Args["cust_id_old"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("WardRepository, Update, err:", err.Error())
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

func (repository *WardRepositoryImpl) Delete(wardId string, deletedBy int64, custId string) error {

	var nRows int64
	query := `DELETE FROM mst.m_ward
			  WHERE cust_id = :cust_id
			  AND ward_id = :ward_id`

	wMap := map[string]interface{}{
		"cust_id": custId,
		"ward_id": wardId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("WardRepositoryImpl, Delete, err:", err.Error())
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
