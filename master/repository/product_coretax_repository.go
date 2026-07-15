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
	"net/url"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductCoreTaxRepository interface {
	FindOneByProCodeCoreTaxAndCustId(productCoreTaxCode string, custId string) (model.ProductCoreTax, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ProductCoreTax, total int, lastPage int, err error)
	Store(productCoreTax model.ProductCoreTax) (string, error)
	Update(productCoreTaxCode string, request entity.UpdateProductCoreTaxRequest) error
	Delete(custId string, productCoreTaxCode string, deletedBy int64) error
}

func NewProductCoreTaxRepository(db *sqlx.DB) ProductCoreTaxRepository {
	return &productCoreTaxRepositoryImpl{db}
}

type productCoreTaxRepositoryImpl struct {
	*sqlx.DB
}

func (repository *productCoreTaxRepositoryImpl) FindOneByProCodeCoreTaxAndCustId(productCoreTaxCode string, custId string) (model.ProductCoreTax, error) {
	productCoreTax := model.ProductCoreTax{}
	query := `SELECT 
				cust_id, pro_code_coretax, cat_coretax,
				pro_name_coretax, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_product_coretax 
			  WHERE pro_code_coretax = $1 
			  AND cust_id = $2`
	decodedStr, _ := url.QueryUnescape(productCoreTaxCode)
	//   if err != nil {
	// 	  fmt.Println("Error:", err)
	// 	  productCoreTax, err
	//   }

	err := repository.Get(&productCoreTax, query, decodedStr, custId)
	if err != nil {
		log.Println("productCoreTaxRepository, FindOneByProCodeCoreTaxAndCustId, err:", err.Error())
		return productCoreTax, err
	}

	return productCoreTax, nil
}

func (repository *productCoreTaxRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ProductCoreTax, int, int, error) {

	productCoreTaxs := []model.ProductCoreTax{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` un.cust_id, un.pro_code_coretax, cat_coretax,
					un.pro_name_coretax, un.is_active, un.created_by,
					un.created_at, un.updated_by, un.updated_at,
					un.is_del, un.deleted_by, un.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE un.is_del = false 
				AND un.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (pro_code_coretax ILIKE '%` + dataFilter.Query + `%' 
					OR pro_name_coretax ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND un.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND un.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_product_coretax un
	LEFT JOIN sys.m_user u ON u.user_id = un.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("productCoreTaxRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productCoreTaxRepository, count total, err:", err.Error())
		return productCoreTaxs, 0, 0, err
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
		sortBy := `pro_code_coretax`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	fmt.Println("filteran", dataFilter.Limit, dataFilter.Page)

	all := dataFilter.Limit > 0 || dataFilter.Page > 0

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))
	if all {
		querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))
	}

	// log.Println("productCoreTaxRepository, querySelect:", querySelect)
	err = repository.Select(&productCoreTaxs, querySelect)
	if err != nil {
		log.Println("productCoreTaxRepository, FindAllByCustId, err:", err.Error())
		return productCoreTaxs, total, lastPage, err
	}

	return productCoreTaxs, total, lastPage, nil
}

func (repository *productCoreTaxRepositoryImpl) Store(productCoreTax model.ProductCoreTax) (string, error) {
	query :=
		`INSERT INTO mst.m_product_coretax(
			cust_id, cat_coretax, pro_code_coretax, pro_name_coretax, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12
		) RETURNING pro_code_coretax;`
	lastInsertId := productCoreTax.ProCodeCoreTax
	err := repository.QueryRow(query,
		productCoreTax.CustId, productCoreTax.CatCoreTax, productCoreTax.ProCodeCoreTax, productCoreTax.ProNameCoreTax, productCoreTax.IsActive,
		productCoreTax.CreatedBy, productCoreTax.CreatedAt, productCoreTax.UpdatedBy, productCoreTax.UpdatedAt,
		productCoreTax.IsDel, productCoreTax.DeletedBy, productCoreTax.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("productCoreTaxRepository, Store, err:", err.Error())
		return productCoreTax.ProCodeCoreTax, err
	}
	return productCoreTax.ProCodeCoreTax, nil
}

func (repository *productCoreTaxRepositoryImpl) Update(productCoreTaxCode string, request entity.UpdateProductCoreTaxRequest) error {
	var (
		r            model.ProductCoreTaxUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("productCoreTaxRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_product_coretax
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND pro_code_coretax = :pro_code_coretax_old;`

	// log.Println("productCoreTaxRepository, Update, query:", query)

	sqlPatch.Args["pro_code_coretax_old"] = productCoreTaxCode
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("productCoreTaxRepository, Update, err:", err.Error())
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

func (repository *productCoreTaxRepositoryImpl) Delete(custId string, productCoreTaxCode string, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product_coretax
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND pro_code_coretax = :pro_code_coretax;`

	wMap := map[string]interface{}{
		"cust_id":          custId,
		"pro_code_coretax": productCoreTaxCode,
		"deleted_by":       deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ProductCoreTaxRepository, Delete, err:", err.Error())
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
