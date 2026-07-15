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

type BrandRepository interface {
	ProductLineIsUsed(pLId int, custId string) (bool, error)
	FindOneByBrandIdAndCustId(brandId int, custId string) (model.Brand, error)
	FindOneByBrandCodeAndCustId(brandCode, custId string) (model.Brand, error)
	FindAllByCustId(dataFilter entity.BrandQueryFilter, custId string) (consPro []model.Brand, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.BrandQueryFilter, custId string) (consPro []model.Brand, total int, lastPage int, err error)
	Store(brand model.Brand) (int, error)
	Update(brandId int, request entity.UpdateBrandRequest) error
	Delete(custId string, brandId int, deletedBy int64) error
}

func NewBrandRepository(db *sqlx.DB) BrandRepository {
	return &BrandRepositoryImpl{db}
}

type BrandRepositoryImpl struct {
	*sqlx.DB
}

func (repository *BrandRepositoryImpl) ProductLineIsUsed(pLId int, custId string) (bool, error) {
	isExists := model.IsExists{}
	query := `SELECT EXISTS(SELECT brand_id 
							FROM mst.m_brand 
							WHERE is_del = false 
							AND pl_id = $1 
							AND cust_id = $2);`
	err := repository.Get(&isExists, query, pLId, custId)
	if err != nil {
		log.Println("MProductRepository, IsExists, err:", err.Error())
		return false, err
	}

	return isExists.Exists, err
}

func (repository *BrandRepositoryImpl) FindOneByBrandIdAndCustId(brandId int, custId string) (model.Brand, error) {
	brand := model.Brand{}
	query := `SELECT 
				b.brand_id, b.brand_code, 
				b.pl_id, pl.pl_code, pl.pl_name,
				b.eff_call, b.min_item, 
				b.brand_name, b.is_active, b.created_by,
				b.created_at, b.updated_by, b.updated_at,
				b.is_del, b.deleted_by, b.deleted_at
			  FROM mst.m_brand b
			  LEFT JOIN mst.m_product_line pl ON (b.pl_id=pl.pl_id) 
			  WHERE b.is_del = false 
				AND b.brand_id = $1 
				AND b.cust_id = $2`
	err := repository.Get(&brand, query, brandId, custId)
	if err != nil {
		log.Println("brandRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return brand, err
	}

	return brand, nil
}

func (repository *BrandRepositoryImpl) FindOneByBrandCodeAndCustId(brandCode, custId string) (model.Brand, error) {
	brand := model.Brand{}
	query := `SELECT 
				brand_id, brand_code, pl_id, eff_call, min_item, 
				brand_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_brand 
			  WHERE is_del = false 
			  AND brand_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&brand, query, brandCode, custId)
	if err != nil {
		log.Println("brandRepository, FindOneByBrandCode, err:", err.Error())
		return brand, err
	}

	return brand, nil
}

func (repository *BrandRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.BrandQueryFilter, custId string) ([]model.Brand, int, int, error) {

	brands := []model.Brand{}
	selectCount := ` COUNT(b.*) AS total `
	selectField := ` b.brand_id, b.brand_code, 
					b.pl_id, pl.pl_code, pl.pl_name,
					b.eff_call, b.min_item, 
					b.brand_name  `

	qWhere := ` WHERE b.cust_id = '` + custId + `' 
				AND b.is_del = false 
				AND b.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (b.brand_code ILIKE '%` + dataFilter.Query + `%' 
					OR b.brand_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.BrandIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandIds, ",")
		qWhere += ` AND b.brand_id IN (` + intArrStr + `)`
	}

	if len(dataFilter.PlIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PlIds, ",")
		qWhere += ` AND b.pl_id IN (` + intArrStr + `)`
	}

	if dataFilter.PlId != 0 {
		qWhere += `AND b.pl_id = ` + strconv.Itoa(dataFilter.PlId) + ` `
	}

	if dataFilter.BrandId != 0 {
		qWhere += `AND b.brand_id = ` + strconv.Itoa(dataFilter.BrandId) + ` `
	}

	qFrom := ` FROM mst.m_brand b
			   LEFT JOIN mst.m_product_line pl ON b.pl_id = pl.pl_id AND pl.cust_id = '` + custId + `'`

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("brandRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("brandRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`b.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `b.brand_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("brandRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("brandRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return brands, total, 1, err
	}

	return brands, total, 1, nil
}

func (repository *BrandRepositoryImpl) FindAllByCustId(dataFilter entity.BrandQueryFilter, custId string) ([]model.Brand, int, int, error) {

	brands := []model.Brand{}
	selectCount := ` COUNT(b.*) AS total `
	selectField := ` b.brand_id, b.brand_code, 
					b.pl_id, pl.pl_code, pl.pl_name,
					b.eff_call, b.min_item, 
					b.brand_name, b.is_active, b.created_by,
					b.created_at, b.updated_by, b.updated_at,
					b.is_del, b.deleted_by, b.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE b.cust_id = '` + custId + `' 
				AND b.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND (b.brand_code ILIKE '%` + dataFilter.Query + `%' 
					OR b.brand_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.BrandIds) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandIds, ",")
		qWhere += ` AND b.brand_id IN (` + intArrStr + `) `
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND b.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND b.is_active = false `
		}
	}

	if dataFilter.PlId != 0 {
		qWhere += `AND b.pl_id = ` + strconv.Itoa(dataFilter.PlId) + ` `
	}

	if dataFilter.BrandId != 0 {
		qWhere += `AND b.brand_id = ` + strconv.Itoa(dataFilter.BrandId) + ` `
	}

	qFrom := ` FROM mst.m_brand b
			   LEFT JOIN mst.m_product_line pl ON b.pl_id = pl.pl_id AND pl.cust_id = '` + custId + `'
			   LEFT JOIN sys.m_user u ON u.user_id = b.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("brandRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("brandRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`b.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `b.brand_id`
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

	// log.Println("brandRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("brandRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *BrandRepositoryImpl) Store(brand model.Brand) (int, error) {
	query :=
		`INSERT INTO mst.m_brand(
			cust_id, brand_code, brand_name, 
			pl_id, eff_call, min_item, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11, 
			$12, $13, $14 
		) RETURNING brand_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		brand.CustId, brand.BrandCode, brand.BrandName,
		brand.PlId, brand.EffCall, brand.MinItem, brand.IsActive,
		brand.CreatedBy, brand.CreatedAt, brand.UpdatedBy, brand.UpdatedAt,
		brand.IsDel, brand.DeletedBy, brand.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("brandRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *BrandRepositoryImpl) Update(brandId int, request entity.UpdateBrandRequest) error {
	var (
		r            model.BrandUpdate
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

	query := `UPDATE mst.m_brand
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND brand_id = :brand_id;`

	sqlPatch.Args["brand_id"] = brandId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("brandRepository, Update, err:", err.Error())
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

func (repository *BrandRepositoryImpl) Delete(custId string, brandId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_brand
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND brand_id = :brand_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"brand_id":   brandId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("BrandRepository, Delete, err:", err.Error())
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
