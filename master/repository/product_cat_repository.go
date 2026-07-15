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

type ProductCatRepository interface {
	FindOneByPCatIdAndCustId(pCatId int, custId string) (model.ProductCat, error)
	FindOneByPCatCodeAndCustId(pCatCode, custId string) (model.ProductCat, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ProductCat, total int, lastPage int, err error)
	Store(productCat model.ProductCat) (int, error)
	Update(pCatId int, request entity.UpdateProductCatRequest) error
	Delete(custId string, pCatId int, deletedBy int64) error
}

func NewProductCatRepository(db *sqlx.DB) ProductCatRepository {
	return &productCatRepositoryImpl{db}
}

type productCatRepositoryImpl struct {
	*sqlx.DB
}

func (repository *productCatRepositoryImpl) FindOneByPCatIdAndCustId(pCatId int, custId string) (model.ProductCat, error) {
	productCat := model.ProductCat{}
	query := `SELECT 
				cust_id, pcat_id, pcat_code,
				pcat_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_product_cat 
			  WHERE is_del = false 
			  AND pcat_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&productCat, query, pCatId, custId)
	if err != nil {
		log.Println("productCatRepository, FindOneByPCatIdAndCustId, err:", err.Error())
		return productCat, err
	}

	return productCat, nil
}

func (repository *productCatRepositoryImpl) FindOneByPCatCodeAndCustId(pCatCode, custId string) (model.ProductCat, error) {
	productCat := model.ProductCat{}
	query := `SELECT 
				cust_id, pcat_id, pcat_code,
				pcat_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_product_cat 
			  WHERE is_del = false 
			  AND pcat_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&productCat, query, pCatCode, custId)
	if err != nil {
		log.Println("productCatRepository, FindOneByPCatCode, err:", err.Error())
		return productCat, err
	}

	return productCat, nil
}

func (repository *productCatRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ProductCat, int, int, error) {

	productCats := []model.ProductCat{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pc.cust_id, pc.pcat_id, pc.pcat_code,
				pc.pcat_name, pc.is_active, pc.created_by,
				pc.created_at, pc.updated_by, pc.updated_at,
				pc.is_del, pc.deleted_by, pc.deleted_at,
				u.user_fullname AS updated_by_name `
	qWhere := ` WHERE pc.is_del = false 
				AND pc.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (pc.pcat_code ILIKE '%` + dataFilter.Query + `%' 
					OR pc.pcat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND pc.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND pc.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_product_cat pc
	LEFT JOIN sys.m_user u ON u.user_id = pc.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("productCatRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productCatRepository, count total, err:", err.Error())
		return productCats, 0, 0, err
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
		sortBy := `pc.pcat_id`
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

	// log.Println("productCatRepository, querySelect:", querySelect)
	err = repository.Select(&productCats, querySelect)
	if err != nil {
		log.Println("productCatRepository, FindAllByCustId, err:", err.Error())
		return productCats, total, lastPage, err
	}

	return productCats, total, lastPage, nil
}

func (repository *productCatRepositoryImpl) Store(productCat model.ProductCat) (int, error) {
	query :=
		`INSERT INTO mst.m_product_cat(
			cust_id, pcat_code, pcat_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING pcat_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		productCat.CustId, productCat.PCatCode, productCat.PCatName, productCat.IsActive,
		productCat.CreatedBy, productCat.CreatedAt, productCat.UpdatedBy, productCat.UpdatedAt,
		productCat.IsDel, productCat.DeletedBy, productCat.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("productCatRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *productCatRepositoryImpl) Update(pCatId int, request entity.UpdateProductCatRequest) error {
	var (
		r            model.ProductCatUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("productCatRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_product_cat
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND pcat_id = :pcat_id;`

	// log.Println("productCatRepository, Update, query:", query)

	sqlPatch.Args["pcat_id"] = pCatId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("productCatRepository, Update, err:", err.Error())
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

func (repository *productCatRepositoryImpl) Delete(custId string, pCatId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product_cat
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND pcat_id = :pcat_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"pcat_id":    pCatId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ProductCatRepository, Delete, err:", err.Error())
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
