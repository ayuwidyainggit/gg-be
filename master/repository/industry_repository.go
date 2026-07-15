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

type IndustryRepository interface {
	FindOneByIndustryIdAndCustId(industryId int, custId string) (model.Industry, error)
	FindOneByIndustryCodeAndCustId(industryCode string, custId string) (model.Industry, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (industry []model.Industry, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (industry []model.Industry, total int, lastPage int, err error)
	Store(industry model.Industry) (int, error)
	Update(industryId int, request entity.UpdateIndustryRequest) error
	Delete(custId string, industryId int, deletedBy int64) error
}

func NewIndustryRepository(db *sqlx.DB) IndustryRepository {
	return &industryRepositoryImpl{db}
}

type industryRepositoryImpl struct {
	*sqlx.DB
}

func (repository *industryRepositoryImpl) FindOneByIndustryIdAndCustId(industryId int, custId string) (model.Industry, error) {
	industry := model.Industry{}
	query := `SELECT 
				cust_id, industry_id, industry_code,
				industry_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_industry 
			  WHERE industry_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&industry, query, industryId, custId)
	if err != nil {
		log.Println("industryRepository, FindOneByIndustryCodeAndCustId, err:", err.Error())
		return industry, err
	}

	return industry, nil
}

func (repository *industryRepositoryImpl) FindOneByIndustryCodeAndCustId(industryCode string, custId string) (model.Industry, error) {
	industry := model.Industry{}
	query := `SELECT 
				cust_id, industry_id, industry_code,
				industry_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_industry 
			  WHERE industry_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&industry, query, industryCode, custId)
	if err != nil {
		log.Println("industryRepository, FindOneByIndustryCodeAndCustId, err:", err.Error())
		return industry, err
	}

	return industry, nil
}

func (repository *industryRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Industry, int, int, error) {

	industrys := []model.Industry{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.industry_id, a.industry_code, a.industry_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.industry_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.industry_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_industry a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("industryRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("industryRepository, count total, err:", err.Error())
		return industrys, 0, 0, err
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
		sortBy := `a.industry_id`
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

	// log.Println("industryRepository, querySelect:", querySelect)
	err = repository.Select(&industrys, querySelect)
	if err != nil {
		log.Println("industryRepository, FindAllByCustId, err:", err.Error())
		return industrys, total, lastPage, err
	}

	return industrys, total, lastPage, nil
}

func (repository *industryRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Industry, int, int, error) {

	industrys := []model.Industry{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.industry_id, a.industry_code, a.industry_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.industry_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.industry_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_industry a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("industryRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("industryRepository, count total, err:", err.Error())
		return industrys, 0, 0, err
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
		sortBy := `a.industry_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&industrys, querySelect)
	if err != nil {
		log.Println("industryRepository, FindAllByCustId, err:", err.Error())
		return industrys, total, 1, err
	}

	return industrys, total, 1, nil
}

func (repository *industryRepositoryImpl) Store(industry model.Industry) (int, error) {
	query :=
		`INSERT INTO mst.m_industry(
			cust_id, industry_code, industry_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING industry_id;`
	lastInsertId := industry.IndustryId
	err := repository.QueryRow(query,
		industry.CustId, industry.IndustryCode, industry.IndustryName,
		industry.IsActive, industry.CreatedBy, industry.CreatedAt, industry.UpdatedBy,
		industry.UpdatedAt, industry.IsDel, industry.DeletedBy, industry.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("industryRepository, Store, err:", err.Error())
		return industry.IndustryId, err
	}
	return industry.IndustryId, nil
}

func (repository *industryRepositoryImpl) Update(industryId int, request entity.UpdateIndustryRequest) error {
	var (
		r            model.IndustryUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("industryRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_industry
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND industry_id = :industry_id_old;`

	// log.Println("industryRepository, Update, query:", query)

	sqlPatch.Args["industry_id_old"] = industryId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("industryRepository, Update, err:", err.Error())
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

func (repository *industryRepositoryImpl) Delete(custId string, industryId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_industry
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND industry_id = :industry_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"industry_id": industryId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("IndustryRepository, Delete, err:", err.Error())
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
