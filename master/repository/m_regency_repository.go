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

type RegencyRepository interface {
	FindOneByRegencyId(regencyId string, custId string) (model.Regency, error)
	Store(regency model.Regency) (string, error)
	Update(regencyId string, request entity.UpdateRegencyRequest) error
	Delete(regencyId string, deletedBy int64, custId string) error
	FindAllByCustId(dataFilter entity.RegencyQueryFilter, custId string) (province []model.Regency, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.RegencyQueryFilter, custId string) (province []model.Regency, total int, lastPage int, err error)
}

func NewRegencyRepository(db *sqlx.DB) RegencyRepository {
	return &RegencyRepositoryImpl{db}
}

type RegencyRepositoryImpl struct {
	*sqlx.DB
}

func (repository *RegencyRepositoryImpl) FindOneByRegencyId(regencyId string, custId string) (model.Regency, error) {
	regency := model.Regency{}
	query := `SELECT 
				regency_id,regency, province_id, is_active, created_by,
				created_at, updated_by, updated_at
			  FROM mst.m_regency 
			  WHERE regency_id = $1 and cust_id = $2`
	err := repository.Get(&regency, query, regencyId, custId)
	if err != nil {
		log.Println("RegencyRepository, FindOneByRegencyId, err:", err.Error())
		return regency, err
	}

	return regency, nil
}

func (repository *RegencyRepositoryImpl) FindAllByCustId(dataFilter entity.RegencyQueryFilter, custId string) ([]model.Regency, int, int, error) {

	regencys := []model.Regency{}
	selectCount := ` COUNT(*) AS total `
	selectField := `mr.regency_id, mr.regency, mr.province_id, 
					mp.province, 
					mr.is_active, mr.created_by,
					mr.created_at, mr.updated_by, mr.updated_at,
					u.user_fullname AS updated_by_name `
	
	qFrom := ` 	FROM mst.m_regency mr 
				LEFT JOIN mst.m_province mp on mp.province_id = mr.province_id AND mp.cust_id = '` + custId + `'
				LEFT JOIN sys.m_user u ON u.user_id = mr.updated_by   `

	qWhere := ` WHERE mr.cust_id = '` + custId + `' `
	
	if dataFilter.Query != "" {
		qWhere += ` AND (mr.regency_id ILIKE '%` + dataFilter.Query + `%' OR mr.regency ILIKE '%` + dataFilter.Query + `%' OR mr.province_id ILIKE '%` + dataFilter.Query + `%' )`
	}
	
	if dataFilter.ProvinceId != "" {
		qWhere += `AND mr.province_id = '` + dataFilter.ProvinceId + `' `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND mr.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND mr.is_active = false `
		}
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("RegencyRepository, count total, err:", err.Error())
		return regencys, 0, 0, err
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
		sortBy := `regency_id`
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

	// log.Println("RegencyRepository, querySelect:", querySelect)
	err = repository.Select(&regencys, querySelect)
	if err != nil {
		log.Println("RegencyRepository, FindAllByCustId, err:", err.Error())
		return regencys, total, lastPage, err
	}

	return regencys, total, lastPage, nil
}

func (repository *RegencyRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.RegencyQueryFilter, custId string) ([]model.Regency, int, int, error) {
	regencys := []model.Regency{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` mr.regency_id, mr.regency, mr.province_id, 
					mp.province, mr.is_active `

	qFrom := ` FROM mst.m_regency mr 
	LEFT JOIN mst.m_province mp on mp.province_id = mr.province_id AND mp.cust_id = '` + custId + `'
	LEFT JOIN sys.m_user u ON u.user_id = mr.updated_by   `

	qWhere := ` WHERE mr.cust_id = '` + custId + `' AND mr.is_active = true  `

	if dataFilter.Query != "" {
		qWhere += ` AND (mr.regency_id ILIKE '%` + dataFilter.Query + `%' OR mr.regency ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.ProvinceId != "" {
		qWhere += `AND mr.province_id = '` + dataFilter.ProvinceId + `' `
	}

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("RegencyRepository, count total, err:", err.Error())
		return regencys, 0, 0, err
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
		sortBy := `regency_id`
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

	// log.Println("RegencyRepository, querySelect:", querySelect)
	err = repository.Select(&regencys, querySelect)
	if err != nil {
		log.Println("RegencyRepository, FindAllByCustId, err:", err.Error())
		return regencys, total, lastPage, err
	}

	return regencys, total, lastPage, nil
}

func (repository *RegencyRepositoryImpl) Store(regency model.Regency) (string, error) {
	query :=
		`INSERT INTO mst.m_regency(
			cust_id, regency_id, regency, province_id, 
			is_active, created_by, created_at, updated_by, 
			updated_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9
		) RETURNING regency_id;`
	lastInsertId := regency.RegencyId
	err := repository.QueryRow(query,
		regency.CustId, regency.RegencyId, regency.Regency, regency.ProvinceId,
		regency.IsActive, regency.CreatedBy, regency.CreatedAt, regency.UpdatedBy,
		regency.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("RegencyRepository, Store, err:", err.Error())
		return regency.RegencyId, err
	}
	return regency.RegencyId, nil
}

func (repository *RegencyRepositoryImpl) Update(regencyId string, request entity.UpdateRegencyRequest) error {
	var (
		r            model.RegencyUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("RegencyRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_regency
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id_old
			  AND regency_id = :regency_id_old;`

	// log.Println("RegencyRepository, Update, query:", query)

	sqlPatch.Args["regency_id_old"] = regencyId
	sqlPatch.Args["cust_id_old"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("RegencyRepository, Update, err:", err.Error())
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

func (repository *RegencyRepositoryImpl) Delete(regencyId string, deletedBy int64, custId string) error {

	var nRows int64
	query := `DELETE FROM mst.m_regency
			  WHERE cust_id = :cust_id
			  AND regency_id = :regency_id`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"regency_id": regencyId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("RegencyRepositoryImpl, Delete, err:", err.Error())
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
