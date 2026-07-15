package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/str"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SubBrand1Repository interface {
	FindOneBySubBrand1IdAndCustId(subBrand1Id int, custId string) (model.SubBrand1, error)
	FindOneBySubBrand1CodeAndCustId(subBrand1Code, custId string) (model.SubBrand1, error)
	FindAllByCustId(dataFilter entity.SubBrand1QueryFilter, custId string) (subBrand1 []model.SubBrand1, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.SubBrand1QueryFilter, custId string) (subBrand1 []model.SubBrand1, total int, lastPage int, err error)
	FindAllByCustIdMatGroupMode(dataFilter entity.SubBrand1QueryFilter, custId string) (subBrand1 []model.SubBrand1, total int, lastPage int, err error)
	FindAllByCustIdSubBrand(dataFilter entity.SubBrandQueryFilter, custId string) (subBrand1 []model.SubBrand1, total int, lastPage int, err error)
	Store(subBrand1 model.SubBrand1) (int, error)
	Update(subBrand1Id int, request entity.UpdateSubBrand1Request) error
	Delete(custId string, subBrand1Id int, deletedBy int64) error
}

func NewSubBrand1Repository(db *sqlx.DB) SubBrand1Repository {
	return &subBrand1RepositoryImpl{db}
}

type subBrand1RepositoryImpl struct {
	*sqlx.DB
}

func (repository *subBrand1RepositoryImpl) FindOneBySubBrand1IdAndCustId(subBrand1Id int, custId string) (model.SubBrand1, error) {
	subBrand1 := model.SubBrand1{}
	query := `SELECT 
				sb1.cust_id,
				sb1.brand_id, 
				sb1.sbrand1_id,
				sb1.sbrand1_code,
				sb1.sbrand1_name,
				sb1.eff_call,
				sb1.min_item,
				sb1.is_active, 
				sb1.created_by,
				sb1.created_at, 
				sb1.updated_by, sb1.updated_at,
				sb1.is_del, sb1.deleted_by, sb1.deleted_at,
				b.brand_code, b.brand_name,
				pl.pl_id, pl.pl_code, pl.pl_name
			  FROM mst.m_sub_brand1 AS sb1
			  LEFT JOIN mst.m_brand AS b ON b.brand_id = sb1.brand_id AND b.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_product_line AS pl ON pl.pl_id = b.pl_id AND pl.cust_id = '` + custId + `'
			  WHERE sb1.is_del = false 
			  AND sb1.sbrand1_id = $1 
			  AND sb1.cust_id = $2`
	err := repository.Get(&subBrand1, query, subBrand1Id, custId)
	if err != nil {
		log.Println("subBrand1Repository, FindOneBySBrand1IdAndCustId, err:", err.Error())
		return subBrand1, err
	}

	return subBrand1, nil
}

func (repository *subBrand1RepositoryImpl) FindOneBySubBrand1CodeAndCustId(subBrand1Code, custId string) (model.SubBrand1, error) {
	subBrand1 := model.SubBrand1{}
	query := `SELECT 
				sb1.cust_id,
				sb1.brand_id,
				sb1.sbrand1_id,
				sb1.sbrand1_code,
				sb1.sbrand1_name,
				sb1.eff_call,
				sb1.min_item,
				sb1.is_active, 
				sb1.created_by,
				sb1.created_at, 
				sb1.updated_by, sb1.updated_at,
				sb1.is_del, sb1.deleted_by, sb1.deleted_at,
				b.brand_code, b.brand_name,
				pl.pl_id, pl.pl_code, pl.pl_name
			  FROM mst.m_sub_brand1 AS sb1
			  LEFT JOIN mst.m_brand AS b ON b.brand_id = sb1.brand_id AND b.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_product_line AS pl ON pl.pl_id = b.pl_id AND pl.cust_id = '` + custId + `'
			  WHERE sb1.is_del = false 
			  AND sb1.sbrand1_code = $1 
			  AND sb1.cust_id = $2`
	err := repository.Get(&subBrand1, query, subBrand1Code, custId)
	if err != nil {
		log.Println("subBrand1Repository, FindOneBySBrand1Code, err:", err.Error())
		return subBrand1, err
	}

	return subBrand1, nil
}

func (repository *subBrand1RepositoryImpl) FindAllByCustId(dataFilter entity.SubBrand1QueryFilter, custId string) ([]model.SubBrand1, int, int, error) {

	subBrand1s := []model.SubBrand1{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb1.cust_id,
					sb1.brand_id, 
					sb1.sbrand1_id,
					sb1.sbrand1_code,
					sb1.sbrand1_name,
					sb1.eff_call,
					sb1.min_item,
					sb1.is_active, 
					sb1.created_by,
					sb1.created_at, 
					sb1.updated_by, sb1.updated_at,
					sb1.is_del, sb1.deleted_by, sb1.deleted_at,
					b.brand_code, b.brand_name,
					pl.pl_id, pl.pl_code, pl.pl_name,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sb1.cust_id = '` + custId + `' AND sb1.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb1.sbrand1_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb1.sbrand1_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.SBrand1Ids) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SBrand1Ids, ",")
		qWhere += ` AND sb1.sbrand1_id IN (` + intArrStr + `)`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND sb1.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND sb1.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_sub_brand1 sb1 
			   LEFT JOIN mst.m_brand AS b ON b.brand_id = sb1.brand_id AND b.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_product_line AS pl ON pl.pl_id = b.pl_id AND pl.cust_id = '` + custId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = sb1.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("subBrand1Repository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand1Repository, count total, err:", err.Error())
		return subBrand1s, 0, 0, err
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
		sortBy := `sbrand1_id`
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

	// log.Println("subBrand1Repository, querySelect:", querySelect)
	err = repository.Select(&subBrand1s, querySelect)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustId, err:", err.Error())
		return subBrand1s, total, lastPage, err
	}

	return subBrand1s, total, lastPage, nil
}

func (repository *subBrand1RepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.SubBrand1QueryFilter, custId string) ([]model.SubBrand1, int, int, error) {

	subBrand1s := []model.SubBrand1{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb1.brand_id, sb1.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name, 
					 b.brand_code, b.brand_name,
					 pl.pl_id, pl.pl_code, pl.pl_name,
					 u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sb1.cust_id = '` + custId + `' AND sb1.is_del = false AND sb1.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb1.sbrand1_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb1.sbrand1_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.SBrand1Ids) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SBrand1Ids, ",")
		qWhere += ` AND sb1.sbrand1_id IN (` + intArrStr + `)`
	}

	if len(dataFilter.BrandIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandIds, ",")
		qWhere += ` AND b.brand_id IN (` + intArrStr + `)`
	}

	qFrom := ` FROM mst.m_sub_brand1 sb1 
			   LEFT JOIN mst.m_brand AS b ON b.brand_id = sb1.brand_id AND b.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_product_line AS pl ON pl.pl_id = b.pl_id AND pl.cust_id = '` + custId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = pl.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdLookupMode, count total, err:", err.Error())
		return subBrand1s, 0, 0, err
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
		sortBy := `sbrand1_code`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&subBrand1s, querySelect)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdLookupMode, err:", err.Error())
		return subBrand1s, total, 1, err
	}

	return subBrand1s, total, 1, nil
}

func (repository *subBrand1RepositoryImpl) FindAllByCustIdMatGroupMode(dataFilter entity.SubBrand1QueryFilter, custId string) ([]model.SubBrand1, int, int, error) {

	subBrand1s := []model.SubBrand1{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb1.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					b.brand_id, b.brand_code, b.brand_name,
					b.pl_id, pl.pl_code, pl.pl_name, 
					CONCAT(pl.pl_code, '-', b.brand_code, '-', sb1.sbrand1_code) AS mat_group_code,
					CONCAT(pl.pl_name, '-', b.brand_name, '-', sb1.sbrand1_name) AS mat_group_name,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sb1.cust_id = '` + custId + `' AND sb1.is_del = false AND sb1.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb1.sbrand1_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb1.sbrand1_name ILIKE '%` + dataFilter.Query + `%'
					OR b.brand_name ILIKE '%` + dataFilter.Query + `%'
					OR pl.pl_name ILIKE '%` + dataFilter.Query + `%' 
					OR CONCAT(pl.pl_code, '-', b.brand_code, '-', sb1.sbrand1_code) ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SBrand1Id != nil {
		if *dataFilter.SBrand1Id != 0 {
			qWhere += ` AND sb1.sbrand1_id = ` + strconv.Itoa(*dataFilter.SBrand1Id)
		}
	}

	if len(dataFilter.SBrand1Ids) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SBrand1Ids, ",")
		qWhere += ` AND sb1.sbrand1_id IN (` + intArrStr + `)`
	}

	qFrom := ` FROM mst.m_sub_brand1 AS sb1 
			   LEFT JOIN mst.m_brand AS b ON b.brand_id = sb1.brand_id AND b.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_product_line AS pl ON pl.pl_id = b.pl_id AND pl.cust_id = '` + custId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = sb1.updated_by  `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdMatGroupMode, count total, err:", err.Error())
		return subBrand1s, 0, 0, err
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
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		sortBy := `sb1.sbrand1_code`
		querySelect += fmt.Sprintf(` ORDER BY %s DESC`, sortBy)
	}

	fmt.Println("dataFilter.Limit >>>", dataFilter.Limit)

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
	err = repository.Select(&subBrand1s, querySelect)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdMatGroupMode, err:", err.Error())
		return subBrand1s, total, lastPage, err
	}

	return subBrand1s, total, lastPage, nil
}

func (repository *subBrand1RepositoryImpl) Store(subBrand1 model.SubBrand1) (int, error) {
	query :=
		`INSERT INTO mst.m_sub_brand1(
			cust_id, brand_id, sbrand1_code, 
			sbrand1_name, eff_call, min_item, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8, $9, $10, $11, 
			$12, $13, $14
		) RETURNING sbrand1_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		subBrand1.CustId, subBrand1.BrandId, subBrand1.SBrand1Code,
		subBrand1.SBrand1Name, subBrand1.EffCall, subBrand1.MinItem, subBrand1.IsActive,
		subBrand1.CreatedBy, subBrand1.CreatedAt, subBrand1.UpdatedBy, subBrand1.UpdatedAt,
		subBrand1.IsDel, subBrand1.DeletedBy, subBrand1.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("subBrand1Repository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *subBrand1RepositoryImpl) Update(subBrand1Id int, request entity.UpdateSubBrand1Request) error {
	var (
		r            model.SubBrand1Update
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("subBrand1Repository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sub_brand1
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sbrand1_id = :sbrand1_id;`

	// log.Println("subBrand1Repository, Update, query:", query)

	sqlPatch.Args["sbrand1_id"] = subBrand1Id
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("subBrand1Repository, Update, err:", err.Error())
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

func (repository *subBrand1RepositoryImpl) Delete(custId string, subBrand1Id int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sub_brand1
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND sbrand1_id = :sbrand1_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"sbrand1_id": subBrand1Id,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SubBrand1Repository, Delete, err:", err.Error())
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

func (repository *subBrand1RepositoryImpl) FindAllByCustIdSubBrand(dataFilter entity.SubBrandQueryFilter, custId string) ([]model.SubBrand1, int, int, error) {
	subBrand1s := []model.SubBrand1{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb1.brand_id, 
					sb1.sbrand1_id,
					sb1.sbrand1_code,
					sb1.sbrand1_name,
					sb1.eff_call,
					sb1.min_item,
					sb1.is_active `
	qWhere := ` WHERE sb1.cust_id = '` + custId + `' AND sb1.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb1.sbrand1_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb1.sbrand1_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.Status) > 0 {
		statusConditions := []string{}
		for _, status := range dataFilter.Status {
			if status == 1 {
				statusConditions = append(statusConditions, "sb1.is_active = true")
			} else if status == 0 {
				statusConditions = append(statusConditions, "sb1.is_active = false")
			}
		}
		if len(statusConditions) > 0 {
			qWhere += ` AND (` + strings.Join(statusConditions, " OR ") + `) `
		}
	}

	qFrom := ` FROM mst.m_sub_brand1 sb1 `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdSubBrand, count total, err:", err.Error())
		return subBrand1s, 0, 0, err
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				colName := colSort[0]
				order := colSort[1]
				// Map created_date to created_at
				if colName == "created_date" {
					colName = "created_at"
				}
				sortBy += fmt.Sprintf(`sb1.%s %s, `, colName, strings.ToUpper(order))
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		// Default sort: created_date:desc
		querySelect += ` ORDER BY sb1.created_at DESC`
	}

	// Handle pagination
	if dataFilter.Limit == 0 {
		dataFilter.Limit = 9999
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))
	if lastPage == 0 && total > 0 {
		lastPage = 1
	}

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	err = repository.Select(&subBrand1s, querySelect)
	if err != nil {
		log.Println("subBrand1Repository, FindAllByCustIdSubBrand, err:", err.Error())
		return subBrand1s, total, lastPage, err
	}

	return subBrand1s, total, lastPage, nil
}
