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

type SpecialPriceGroupRepository interface {
	FindOneBySpecialPriceGroupIdAndCustId(specialPriceGroupId int, custId string) (model.SpecialPriceGroup, error)
	FindOneBySpecialPriceGroupCodeAndCustId(specialPriceGroupCode string, custId string) (model.SpecialPriceGroup, error)
	FindAllByCustId(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) (consPro []model.SpecialPriceGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) (consPro []model.SpecialPriceGroup, total int, lastPage int, err error)
	Store(brand model.SpecialPriceGroup) (int, error)
	Update(specialPriceGroupId int, request entity.UpdateSpecialPriceGroupRequest) error
	Delete(custId string, specialPriceGroupId int, deletedBy int64) error
}

func NewSpecialPriceGroupRepository(db *sqlx.DB) SpecialPriceGroupRepository {
	return &SpecialPriceGroupRepositoryImpl{db}
}

type SpecialPriceGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *SpecialPriceGroupRepositoryImpl) FindOneBySpecialPriceGroupIdAndCustId(specialPriceGroupId int, custId string) (model.SpecialPriceGroup, error) {
	specialPriceGroup := model.SpecialPriceGroup{}
	query := `SELECT 
				sr.sp_price_grp_id, 
				sr.sp_price_grp_code,
				sr.sp_price_grp_name, 
				sr.is_active, sr.created_by,
				sr.created_at, sr.updated_by, sr.updated_at,
				sr.is_del, sr.deleted_by, sr.deleted_at
			  FROM mst.m_sp_price_group sr
			  WHERE sr.is_del = false 
				AND sr.sp_price_grp_id = $1 
				AND sr.cust_id = $2`
	err := repository.Get(&specialPriceGroup, query, specialPriceGroupId, custId)
	if err != nil {
		log.Println("specialPriceGroupRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return specialPriceGroup, err
	}

	return specialPriceGroup, nil
}

func (repository *SpecialPriceGroupRepositoryImpl) FindOneBySpecialPriceGroupCodeAndCustId(specialPriceGroupCode string, custId string) (model.SpecialPriceGroup, error) {
	specialPriceGroup := model.SpecialPriceGroup{}
	query := `SELECT 
				sr.sp_price_grp_id, 
				sr.sp_price_grp_code,
				sr.sp_price_grp_name, 
				sr.is_active, sr.created_by,
				sr.created_at, sr.updated_by, sr.updated_at,
				sr.is_del, sr.deleted_by, sr.deleted_at
			  FROM mst.m_sp_price_group sr
			  WHERE sr.is_del = false 
				AND sr.sp_price_grp_code = $1 
				AND sr.cust_id = $2`
	err := repository.Get(&specialPriceGroup, query, specialPriceGroupCode, custId)
	if err != nil {
		log.Println("specialPriceGroupRepository, FindOneByBrandCodeAndCustId, err:", err.Error())
		return specialPriceGroup, err
	}

	return specialPriceGroup, nil
}

func (repository *SpecialPriceGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) ([]model.SpecialPriceGroup, int, int, error) {

	specialPriceGroup := []model.SpecialPriceGroup{}
	selectCount := ` COUNT(sr.*) AS total `
	selectField := `sr.sp_price_grp_id,
					sr.sp_price_grp_name,
					sr.sp_price_grp_code  `
	qWhere := ` WHERE sr.is_del = false AND sr.is_active = true 
				AND sr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sr.sp_price_grp_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.SpecialPriceGroupId != 0 {
		qWhere += `AND sr.sp_price_grp_id = ` + strconv.Itoa(dataFilter.SpecialPriceGroupId) + ` `
	}

	qFrom := ` FROM mst.m_sp_price_group sr `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("specialPriceGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("specialPriceGroupRepository, count total, err:", err.Error())
		return specialPriceGroup, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`sr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `sr.sp_price_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("specialPriceGroupRepository, querySelect:", querySelect)
	err = repository.Select(&specialPriceGroup, querySelect)
	if err != nil {
		log.Println("specialPriceGroupRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return specialPriceGroup, total, 1, err
	}

	return specialPriceGroup, total, 1, nil
}

func (repository *SpecialPriceGroupRepositoryImpl) FindAllByCustId(dataFilter entity.SpecialPriceGroupQueryFilter, custId string) ([]model.SpecialPriceGroup, int, int, error) {

	brands := []model.SpecialPriceGroup{}
	selectCount := ` COUNT(sr.*) AS total `
	selectField := `sr.sp_price_grp_id,
					sr.sp_price_grp_code,
					sr.sp_price_grp_name, 
					sr.is_active, sr.created_by,
					sr.created_at, sr.updated_by, sr.updated_at,
					sr.is_del, sr.deleted_by, sr.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE sr.is_del = false 
				AND sr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (sr.sp_price_grp_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND sr.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND sr.is_active = false `
		}
	}

	if dataFilter.SpecialPriceGroupId != 0 {
		qWhere += `AND sr.sp_price_grp_id = ` + strconv.Itoa(dataFilter.SpecialPriceGroupId) + ` `
	}

	qFrom := ` FROM mst.m_sp_price_group sr
			   LEFT JOIN sys.m_user u ON u.user_id = sr.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("specialPriceGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("specialPriceGroupRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`sr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `sr.sp_price_grp_id`
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

	// log.Println("specialPriceGroupRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("specialPriceGroupRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *SpecialPriceGroupRepositoryImpl) Store(specialPriceGroup model.SpecialPriceGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_sp_price_group(
			cust_id, sp_price_grp_name, sp_price_grp_code,
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11
		) RETURNING sp_price_grp_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		specialPriceGroup.CustId, specialPriceGroup.SpecialPriceGroupName, specialPriceGroup.SpecialPriceGroupCode, specialPriceGroup.IsActive,
		specialPriceGroup.CreatedBy, specialPriceGroup.CreatedAt, specialPriceGroup.UpdatedBy, specialPriceGroup.UpdatedAt,
		specialPriceGroup.IsDel, specialPriceGroup.DeletedBy, specialPriceGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("specialPriceGroupRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *SpecialPriceGroupRepositoryImpl) Update(specialPriceGroupId int, request entity.UpdateSpecialPriceGroupRequest) error {
	var (
		r            model.SpecialPriceGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sp_price_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sp_price_grp_id = :sp_price_grp_id;`

	sqlPatch.Args["sp_price_grp_id"] = specialPriceGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("specialPriceGroupRepository, Update, err:", err.Error())
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

func (repository *SpecialPriceGroupRepositoryImpl) Delete(custId string, specialPriceGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sp_price_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND sp_price_grp_id = :sp_price_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"sp_price_grp_id": specialPriceGroupId,
		"deleted_by":      deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SpecialPriceGroupRepository, Delete, err:", err.Error())
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
