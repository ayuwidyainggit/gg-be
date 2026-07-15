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

type PluProductRepository interface {
	FindOneByPluProductIdAndCustId(pluProId int, custId string) (model.PluProduct, error)
	FindOneByPluProductGrpAndCustId(pluGrpId int, custId string) (model.PluProduct, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.PluProduct, total int, lastPage int, err error)
	Store(pluProduct model.PluProduct) (int, error)
	Update(pluProductId int, request entity.UpdatePluProductRequest) error
	Delete(custId string, pluProId int, deletedBy int64) error
}

func NewPluProductRepository(db *sqlx.DB) PluProductRepository {
	return &pluProductRepositoryImpl{db}
}

type pluProductRepositoryImpl struct {
	*sqlx.DB
}

func (repository *pluProductRepositoryImpl) FindOneByPluProductIdAndCustId(pluProId int, custId string) (model.PluProduct, error) {
	pluProduct := model.PluProduct{}
	query := `SELECT 
				*
			  FROM mst.m_plu_product 
			  WHERE plu_pro_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&pluProduct, query, pluProId, custId)
	if err != nil {
		log.Println("pluProductRepository, FindOneByPluProductIdAndCustId, err:", err.Error())
		return pluProduct, err
	}

	return pluProduct, nil
}

func (repository *pluProductRepositoryImpl) FindOneByPluProductGrpAndCustId(pluGrpId int, custId string) (model.PluProduct, error) {
	pluProduct := model.PluProduct{}
	query := `SELECT 
				*
			  FROM mst.m_plu_product 
			  WHERE plu_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&pluProduct, query, pluGrpId, custId)
	if err != nil {
		log.Println("pluProductRepository, FindOneByPluProductGrpAndCustId, err:", err.Error())
		return pluProduct, err
	}

	return pluProduct, nil
}

func (repository *pluProductRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PluProduct, int, int, error) {

	pluProducts := []model.PluProduct{}
	selectCount := ` COUNT(a.*) AS total `
	selectField := ` a.cust_id, a.plu_pro_id, a.plu_grp_id, a.pro_id, a.plu_no, a.created_at, a.created_by, a.updated_at, a.updated_by, a.deleted_at, a.deleted_by, b.pro_code, b.pro_name,
	u.user_name AS updated_by_name  `
	qWhere := ` WHERE a.is_del = false 
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (b.pro_code ILIKE '%` + dataFilter.Query + `%' OR b.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_plu_product a ` + qWhere
	querySelect :=
		`SELECT ` + selectField +
			` FROM mst.m_plu_product a ` +
			` JOIN mst.m_product b ON b.pro_id = a.pro_id 
			LEFT JOIN sys.m_user u ON u.user_id = a.updated_by
			` +
			qWhere

	// log.Println("pluProductRepository, queryCount:", queryCount)

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("pluProductRepository, count total, err:", err.Error())
		return pluProducts, 0, 0, err
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
		sortBy := `a.plu_grp_id`
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

	// log.Println("pluProductRepository, querySelect:", querySelect)
	err = repository.Select(&pluProducts, querySelect)
	if err != nil {
		log.Println("pluProductRepository, FindAllByCustId, err:", err.Error())
		return pluProducts, total, lastPage, err
	}

	return pluProducts, total, lastPage, nil
}

func (repository *pluProductRepositoryImpl) Store(pluProduct model.PluProduct) (int, error) {

	query :=
		`INSERT INTO mst.m_plu_product(
			cust_id, plu_grp_id, pro_id, plu_no, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING plu_pro_id;`
	lastInsertId := pluProduct.PluProId
	err := repository.QueryRow(query,
		pluProduct.CustId, pluProduct.PluGrpId, pluProduct.ProId,
		pluProduct.PluNo, pluProduct.CreatedBy, pluProduct.CreatedAt, pluProduct.UpdatedBy,
		pluProduct.UpdatedAt, pluProduct.IsDel, pluProduct.DeletedBy, pluProduct.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("pluProductRepository, Store, err:", err.Error())
		return pluProduct.PluProId, err
	}
	return pluProduct.PluProId, nil
}

func (repository *pluProductRepositoryImpl) Update(pluProId int, request entity.UpdatePluProductRequest) error {
	var (
		r            model.PluProductUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("pluProductRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_plu_product
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND plu_pro_id = :plu_pro_id_old;`

	// log.Println("pluProductRepository, Update, query:", query)

	sqlPatch.Args["plu_pro_id_old"] = pluProId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("pluProductRepository, Update, err:", err.Error())
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

func (repository *pluProductRepositoryImpl) Delete(custId string, pluProId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_plu_product
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND plu_pro_id = :plu_pro_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"plu_pro_id": pluProId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PluProductRepository, Delete, err:", err.Error())
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
