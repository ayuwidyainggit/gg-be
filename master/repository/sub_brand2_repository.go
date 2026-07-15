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

type SubBrand2Repository interface {
	FindOneBySBrand2IdAndCustId(cProId int, custId string) (model.SubBrand2, error)
	FindOneBySBrand2CodeAndCustId(cProCode, custId string) (model.SubBrand2, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.SubBrand2, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.SubBrand2, total int, lastPage int, err error)
	Store(subBrand2 model.SubBrand2) (int, error)
	Update(cProId int, request entity.UpdateSubBrand2Request) error
	Delete(custId string, cProId int, deletedBy int64) error
}

func NewSubBrand2Repository(db *sqlx.DB) SubBrand2Repository {
	return &subBrand2RepositoryImpl{db}
}

type subBrand2RepositoryImpl struct {
	*sqlx.DB
}

func (repository *subBrand2RepositoryImpl) FindOneBySBrand2IdAndCustId(cProId int, custId string) (model.SubBrand2, error) {
	subBrand2 := model.SubBrand2{}
	query := `SELECT 
				cust_id, sbrand2_id, sbrand2_code,
				sbrand2_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sub_brand2 
			  WHERE is_del = false 
			  AND sbrand2_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&subBrand2, query, cProId, custId)
	if err != nil {
		log.Println("subBrand2Repository, FindOneBySBrand2IdAndCustId, err:", err.Error())
		return subBrand2, err
	}

	return subBrand2, nil
}

func (repository *subBrand2RepositoryImpl) FindOneBySBrand2CodeAndCustId(cProCode, custId string) (model.SubBrand2, error) {
	subBrand2 := model.SubBrand2{}
	query := `SELECT 
				cust_id, sbrand2_id, sbrand2_code,
				sbrand2_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sub_brand2 
			  WHERE is_del = false 
			  AND sbrand2_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&subBrand2, query, cProCode, custId)
	if err != nil {
		log.Println("subBrand2Repository, FindOneBySBrand2Code, err:", err.Error())
		return subBrand2, err
	}

	return subBrand2, nil
}

func (repository *subBrand2RepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SubBrand2, int, int, error) {

	subBrand2s := []model.SubBrand2{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb2.cust_id, sb2.sbrand2_id, sb2.sbrand2_code,
	sb2.sbrand2_name, sb2.is_active, sb2.created_by,
	sb2.created_at, sb2.updated_by, sb2.updated_at,
	sb2.is_del, sb2.deleted_by, sb2.deleted_at,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sb2.is_del = false 
				AND sb2.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb2.sbrand2_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb2.sbrand2_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND sb2.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND sb2.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_sub_brand2 sb2
	LEFT JOIN sys.m_user u ON u.user_id = sb2.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("subBrand2Repository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand2Repository, count total, err:", err.Error())
		return subBrand2s, 0, 0, err
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
		sortBy := `sbrand2_id`
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

	// log.Println("subBrand2Repository, querySelect:", querySelect)
	err = repository.Select(&subBrand2s, querySelect)
	if err != nil {
		log.Println("subBrand2Repository, FindAllByCustId, err:", err.Error())
		return subBrand2s, total, lastPage, err
	}

	return subBrand2s, total, lastPage, nil
}

func (repository *subBrand2RepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SubBrand2, int, int, error) {

	subBrand2s := []model.SubBrand2{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` sb2.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
	u.user_name AS updated_by_name`
	qWhere := ` WHERE sb2.is_del = false AND sb2.is_active = true 
				AND sb2.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sb2.sbrand2_code ILIKE '%` + dataFilter.Query + `%' 
					OR sb2.sbrand2_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` FROM mst.m_sub_brand2 sb2
	LEFT JOIN sys.m_user u ON u.user_id = sb2.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBrand2Repository, FindAllByCustIdLookupMode, count total, err:", err.Error())
		return subBrand2s, 0, 0, err
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
		sortBy := `sbrand2_code`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&subBrand2s, querySelect)
	if err != nil {
		log.Println("subBrand2Repository, FindAllByCustIdLookupMode, err:", err.Error())
		return subBrand2s, total, 1, err
	}

	return subBrand2s, total, 1, nil
}

func (repository *subBrand2RepositoryImpl) Store(subBrand2 model.SubBrand2) (int, error) {
	query :=
		`INSERT INTO mst.m_sub_brand2(
			cust_id, sbrand2_code, sbrand2_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING sbrand2_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		subBrand2.CustId, subBrand2.SBrand2Code, subBrand2.SBrand2Name, subBrand2.IsActive,
		subBrand2.CreatedBy, subBrand2.CreatedAt, subBrand2.UpdatedBy, subBrand2.UpdatedAt,
		subBrand2.IsDel, subBrand2.DeletedBy, subBrand2.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("subBrand2Repository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *subBrand2RepositoryImpl) Update(cProId int, request entity.UpdateSubBrand2Request) error {
	var (
		r            model.SubBrand2Update
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("subBrand2Repository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sub_brand2
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sbrand2_id = :sbrand2_id;`

	// log.Println("subBrand2Repository, Update, query:", query)

	sqlPatch.Args["sbrand2_id"] = cProId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("subBrand2Repository, Update, err:", err.Error())
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

func (repository *subBrand2RepositoryImpl) Delete(custId string, cProId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sub_brand2
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND sbrand2_id = :sbrand2_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"sbrand2_id": cProId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SubBrand2Repository, Delete, err:", err.Error())
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
