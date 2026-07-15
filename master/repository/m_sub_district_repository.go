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

type SubDistrictRepository interface {
	FindOneBySubDistrictId(subDistrictId string, custId string) (model.SubDistrict, error)
	Store(subDistrict model.SubDistrict) (string, error)
	Update(subDistrictId string, request entity.UpdateSubDistrictRequest) error
	Delete(subDistrictId string, deletedBy int64, custId string) error
	FindAllByCustId(dataFilter entity.SubDistrictQueryFilter, custId string) (province []model.SubDistrict, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.SubDistrictQueryFilter, custId string) (province []model.SubDistrict, total int, lastPage int, err error)
}

func NewSubDistrictRepository(db *sqlx.DB) SubDistrictRepository {
	return &SubDistrictRepositoryImpl{db}
}

type SubDistrictRepositoryImpl struct {
	*sqlx.DB
}

func (repository *SubDistrictRepositoryImpl) FindOneBySubDistrictId(subDistrictId string, custId string) (model.SubDistrict, error) {
	subDistrict := model.SubDistrict{}
	query := `SELECT 
				sub_district_id, sub_district, province_id, regency_id, is_active, created_by,
				created_at, updated_by, updated_at
			  FROM mst.m_sub_district 
			  WHERE sub_district_id = $1 and cust_id = $2`
	err := repository.Get(&subDistrict, query, subDistrictId, custId)
	if err != nil {
		log.Println("SubDistrictRepository, FindOneBySubDistrictId, err:", err.Error())
		return subDistrict, err
	}

	return subDistrict, nil
}

func (repository *SubDistrictRepositoryImpl) FindAllByCustId(dataFilter entity.SubDistrictQueryFilter, custId string) ([]model.SubDistrict, int, int, error) {

	subDistricts := []model.SubDistrict{}
	selectCount := ` COUNT(*) AS total `
	selectField := `   	msd.sub_district_id, 
						msd.sub_district, 
						msd.province_id, 
						msd.regency_id, 
						msd.is_active, 
						mp.province, 
						mr.regency, msd.created_by,
					msd.created_at, msd.updated_by, msd.updated_at,
					u.user_fullname AS updated_by_name `

	qFrom := ` FROM mst.m_sub_district msd 
	LEFT JOIN mst.m_province mp on mp.province_id = msd.province_id AND mp.cust_id = '` + custId + `'
	LEFT JOIN mst.m_regency mr on mr.regency_id = msd.regency_id AND mr.cust_id = '` + custId + `'
	LEFT JOIN sys.m_user u ON u.user_id = msd.updated_by   `

	qWhere := ` WHERE msd.cust_id = '` + custId + `'  `

	if dataFilter.Query != "" {
		qWhere += ` AND (msd.sub_district_id ILIKE '%` + dataFilter.Query + `%' OR msd.sub_district ILIKE '%` + dataFilter.Query + `%' OR mr.regency_id ILIKE '%` + dataFilter.Query + `%' OR mr.regency ILIKE '%` + dataFilter.Query + `%')`
	}
	if dataFilter.RegencyId != "" {
		qWhere += `AND msd.regency_id = '` + dataFilter.RegencyId + `' `
	}
	if dataFilter.ProvinceId != "" {
		qWhere += `AND msd.province_id = '` + dataFilter.ProvinceId + `' `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND msd.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND msd.is_active = false `
		}
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("SubDistrictRepository, count total, err:", err.Error())
		return subDistricts, 0, 0, err
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
		sortBy := `sub_district_id`
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

	// log.Println("SubDistrictRepository, querySelect:", querySelect)
	err = repository.Select(&subDistricts, querySelect)
	if err != nil {
		log.Println("SubDistrictRepository, FindAllByCustId, err:", err.Error())
		return subDistricts, total, lastPage, err
	}

	return subDistricts, total, lastPage, nil
}

func (repository *SubDistrictRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.SubDistrictQueryFilter, custId string) ([]model.SubDistrict, int, int, error) {

	subDistricts := []model.SubDistrict{}
	selectCount := ` COUNT(*) AS total `
	selectField := `  
		msd.sub_district_id, 
		msd.sub_district, 
		msd.province_id, 
		msd.regency_id, 
		msd.is_active, 
		mp.province, 
		mr.regency `

	qFrom := ` 	FROM mst.m_sub_district msd 
	LEFT JOIN mst.m_province mp on mp.province_id = msd.province_id AND mp.cust_id = '` + custId + `'
	LEFT JOIN mst.m_regency mr on mr.regency_id = msd.regency_id AND mr.cust_id = '` + custId + `'
	LEFT JOIN sys.m_user u ON u.user_id = msd.updated_by   `

	qWhere := ` WHERE msd.cust_id = '` + custId + `' AND msd.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (msd.sub_district_id ILIKE '%` + dataFilter.Query + `%' OR msd.sub_district ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.RegencyId != "" {
		qWhere += `AND msd.regency_id = '` + dataFilter.RegencyId + `' `
	}
	if dataFilter.ProvinceId != "" {
		qWhere += `AND msd.province_id = '` + dataFilter.ProvinceId + `' `
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("SubDistrictRepository, count total, err:", err.Error())
		return subDistricts, 0, 0, err
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
		sortBy := `sub_district_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// if dataFilter.Limit == 0 {
	// 	dataFilter.Limit = 10
	// }

	lastPage := int(math.Ceil(float64(float64(total))))

	if dataFilter.Limit > 0 {

		page := dataFilter.Page
		if page-1 < 1 {
			page = 1
		}
		offset := (page - 1) * dataFilter.Limit

		lastPage = int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

		querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))
	}
	// log.Println("SubDistrictRepository, querySelect:", querySelect)
	err = repository.Select(&subDistricts, querySelect)
	if err != nil {
		log.Println("SubDistrictRepository, FindAllByCustId, err:", err.Error())
		return subDistricts, total, lastPage, err
	}

	return subDistricts, total, lastPage, nil
}

func (repository *SubDistrictRepositoryImpl) Store(subDistrict model.SubDistrict) (string, error) {
	query :=
		`INSERT INTO mst.m_sub_district(
			cust_id, sub_district_id, sub_district, province_id, 
			regency_id, is_active, created_by, created_at, 
			updated_by, updated_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10
		) RETURNING sub_district_id;`
	lastInsertId := subDistrict.SubDistrictId
	err := repository.QueryRow(query,
		subDistrict.CustId, subDistrict.SubDistrictId, subDistrict.SubDistrict, subDistrict.ProvinceId,
		subDistrict.RegencyId, subDistrict.IsActive, subDistrict.CreatedBy, subDistrict.CreatedAt, subDistrict.UpdatedBy,
		subDistrict.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("SubDistrictRepository, Store, err:", err.Error())
		return subDistrict.SubDistrictId, err
	}
	return subDistrict.SubDistrictId, nil
}

func (repository *SubDistrictRepositoryImpl) Update(subDistrictId string, request entity.UpdateSubDistrictRequest) error {
	var (
		r            model.SubDistrictUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("SubDistrictRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sub_district
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id_old
			  AND sub_district_id = :sub_district_id_old;`

	log.Println("SubDistrictRepository, Update, query:", query)

	sqlPatch.Args["sub_district_id_old"] = subDistrictId
	sqlPatch.Args["cust_id_old"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("SubDistrictRepository, Update, err:", err.Error())
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

func (repository *SubDistrictRepositoryImpl) Delete(subDistrictId string, deletedBy int64, custId string) error {

	var nRows int64
	query := `DELETE FROM mst.m_sub_district
			  WHERE cust_id = :cust_id
			  AND sub_district_id = :sub_district_id`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"sub_district_id": subDistrictId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SubDistrictRepository, Delete, err:", err.Error())
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
