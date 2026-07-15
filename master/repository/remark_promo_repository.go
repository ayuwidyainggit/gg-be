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

type RemarkPromoRepository interface {
	FindOneByRemarkPromoIdAndCustId(remarkPromoId int, custId string) (model.RemarkPromo, error)
	FindOneByRemarkPromoCodeAndCustId(remarkPromoCode string, custId string) (model.RemarkPromo, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.RemarkPromo, total int, lastPage int, err error)
	Store(remarkPromo model.RemarkPromo) (int, error)
	Update(remarkPromoId int, request entity.UpdateRemarkPromoRequest) error
	Delete(custId string, remPromoId int, deletedBy int64) error
}

func NewDiscRemarkPromoRepository(db *sqlx.DB) RemarkPromoRepository {
	return &remarkPromoRepositoryImpl{db}
}

type remarkPromoRepositoryImpl struct {
	*sqlx.DB
}

func (repository *remarkPromoRepositoryImpl) FindOneByRemarkPromoIdAndCustId(remarkPromoId int, custId string) (model.RemarkPromo, error) {
	remarkPromos := model.RemarkPromo{}
	query := `SELECT * FROM mst.m_remark_promo WHERE rem_promo_id = $1 AND cust_id = $2`
	err := repository.Get(&remarkPromos, query, remarkPromoId, custId)
	if err != nil {
		log.Println("remarkPromoRepository, FindOneByRemarkPromoIdAndCustId, err:", err.Error())
		return remarkPromos, err
	}

	return remarkPromos, nil
}

func (repository *remarkPromoRepositoryImpl) FindOneByRemarkPromoCodeAndCustId(remarkPromoCode string, custId string) (model.RemarkPromo, error) {
	remarkPromo := model.RemarkPromo{}
	query := `SELECT * FROM mst.m_remark_promo WHERE rem_promo_code = $1 AND cust_id = $2`
	err := repository.Get(&remarkPromo, query, remarkPromoCode, custId)
	if err != nil {
		log.Println("remarkPromoRepository, FindOneByRemarkPromoCodeAndCustId, err:", err.Error())
		return remarkPromo, err
	}

	return remarkPromo, nil
}

func (repository *remarkPromoRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.RemarkPromo, int, int, error) {

	remarkPromos := []model.RemarkPromo{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` rp.*,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE rp.is_del = false and rp.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (rp.rem_promo_code ILIKE '%` + dataFilter.Query + `%' OR rp.rem_promo_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND rp.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND rp.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_remark_promo rp
	LEFT JOIN sys.m_user u ON u.user_id = rp.updated_by `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("remarkPromoRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("remarkPromoRepository, count total, err:", err.Error())
		return remarkPromos, 0, 0, err
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
		sortBy := `rem_promo_id`
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

	// log.Println("remarkPromoRepository, querySelect:", querySelect)
	err = repository.Select(&remarkPromos, querySelect)
	if err != nil {
		log.Println("remarkPromoRepository, FindAllByCustId, err:", err.Error())
		return remarkPromos, total, lastPage, err
	}

	return remarkPromos, total, lastPage, nil
}

func (repository *remarkPromoRepositoryImpl) Store(remarkPromo model.RemarkPromo) (int, error) {
	query :=
		`INSERT INTO mst.m_remark_promo(cust_id, rem_promo_code, rem_promo_name, is_active, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING rem_promo_id;`
	lastInsertId := remarkPromo.RemPromoId
	err := repository.QueryRow(query,
		remarkPromo.CustId, remarkPromo.RemPromoCode, remarkPromo.RemPromoName,
		remarkPromo.IsActive, remarkPromo.CreatedBy, remarkPromo.CreatedAt, remarkPromo.UpdatedBy, remarkPromo.UpdatedAt, remarkPromo.IsDel, remarkPromo.DeletedBy, remarkPromo.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("remarkPromoRepository, Store, err:", err.Error())
		return remarkPromo.RemPromoId, err
	}
	return remarkPromo.RemPromoId, nil
}

func (repository *remarkPromoRepositoryImpl) Update(remarkPromoId int, request entity.UpdateRemarkPromoRequest) error {
	var (
		r            model.RemarkPromoUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("remarkPromoRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_remark_promo
			  SET ` + sqlSetFields + `, updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false AND cust_id = :cust_id AND rem_promo_id = :rem_promo_id_old;`

	log.Println("remarkPromoRepository, Update, query:", query)

	sqlPatch.Args["rem_promo_id_old"] = remarkPromoId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("remarkPromoRepository, Update, err:", err.Error())
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

func (repository *remarkPromoRepositoryImpl) Delete(custId string, remPromoId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_remark_promo SET is_del = true, deleted_at = CURRENT_TIMESTAMP, deleted_by = :deleted_by 
			WHERE is_del = false AND cust_id = :cust_id AND rem_promo_id = :rem_promo_id;`

	wMap := map[string]interface{}{
		"cust_id":      custId,
		"rem_promo_id": remPromoId,
		"deleted_by":   deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("RemarkPromoRepository, Delete, err:", err.Error())
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
