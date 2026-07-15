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

type ProvinceRepository interface {
	FindOneByProvinceId(provinceId string, custId string) (model.Province, error)
	Store(province model.Province) (string, error)
	Update(provinceId string, request entity.UpdateProvinceRequest) error
	Delete(provinceId string, deletedBy int64, custId string) error
	FindAll(dataFilter entity.GeneralQueryFilter, custId string) (province []model.Province, total int, lastPage int, err error)
	FindAllLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (province []model.Province, total int, lastPage int, err error)
}

func NewProvinceRepository(db *sqlx.DB) ProvinceRepository {
	return &provinceRepositoryImpl{db}
}

type provinceRepositoryImpl struct {
	*sqlx.DB
}

func (repository *provinceRepositoryImpl) FindOneByProvinceId(provinceId string, custId string) (model.Province, error) {
	province := model.Province{}
	query := `SELECT 
				province_id, province, is_active, created_by,
				created_at, updated_by, updated_at
			  FROM mst.m_province 
			  WHERE province_id = $1 and cust_id = $2`
	err := repository.Get(&province, query, provinceId, custId)
	if err != nil {
		log.Println("provinceRepository, FindOneByProvinceId, err:", err.Error())
		return province, err
	}

	return province, nil
}

func (repository *provinceRepositoryImpl) FindAll(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Province, int, int, error) {

	provinces := []model.Province{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.province_id,
					a.province, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.province_id ILIKE '%` + dataFilter.Query + `%' 
					OR a.province ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_province a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("provinceRepository, count total, err:", err.Error())
		return provinces, 0, 0, err
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
		sortBy := `a.province_id`
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

	// log.Println("provinceRepository, querySelect:", querySelect)
	err = repository.Select(&provinces, querySelect)
	if err != nil {
		log.Println("provinceRepository, FindAllByCustId, err:", err.Error())
		return provinces, total, lastPage, err
	}

	return provinces, total, lastPage, nil
}

func (repository *provinceRepositoryImpl) FindAllLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Province, int, int, error) {

	provinces := []model.Province{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.province_id, a.province, a.is_active `
	qWhere := ` WHERE a.cust_id = '` + custId + `' AND a.is_active = true  `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.province_id ILIKE '%` + dataFilter.Query + `%' 
					OR a.province ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_province a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("provinceRepository, count total, err:", err.Error())
		return provinces, 0, 0, err
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
		sortBy := `a.province_id`
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

	// log.Println("provinceRepository, querySelect:", querySelect)
	err = repository.Select(&provinces, querySelect)
	if err != nil {
		log.Println("provinceRepository, FindAllByCustId, err:", err.Error())
		return provinces, total, lastPage, err
	}

	return provinces, total, lastPage, nil
}

func (repository *provinceRepositoryImpl) Store(province model.Province) (string, error) {
	query :=
		`INSERT INTO mst.m_province(
			cust_id, province_id, province, is_active, 
			created_by, created_at, updated_by, updated_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8
		) RETURNING province_id;`
	lastInsertId := province.ProvinceId
	err := repository.QueryRow(query,
		province.CustId, province.ProvinceId, province.Province,
		province.IsActive, province.CreatedBy, province.CreatedAt, province.UpdatedBy,
		province.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("provinceRepository, Store, err:", err.Error())
		return province.ProvinceId, err
	}
	return province.ProvinceId, nil
}

func (repository *provinceRepositoryImpl) Update(provinceId string, request entity.UpdateProvinceRequest) error {
	var (
		r            model.ProvinceUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("provinceRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_province
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id_old
			  AND province_id = :province_id_old;`

	log.Println("provinceRepository, Update, query:", query)

	sqlPatch.Args["province_id_old"] = provinceId
	sqlPatch.Args["cust_id_old"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("provinceRepository, Update, err:", err.Error())
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

func (repository *provinceRepositoryImpl) Delete(provinceId string, deletedBy int64, custId string) error {

	var nRows int64
	query := `DELETE FROM mst.m_province
			  WHERE cust_id = :cust_id
			  AND province_id = :province_id`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"province_id": provinceId,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("provinceRepositoryImpl, Delete, err:", err.Error())
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
