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

type tprLimitRepositoryImpl struct {
	*sqlx.DB
}

func NewTprLimitRepository(db *sqlx.DB) *tprLimitRepositoryImpl {
	return &tprLimitRepositoryImpl{db}
}

type TprLimitRepository interface {
	FindOneByTprLimitIdAndCustId(tprlimitID int, custId string) (model.MTprLimitList, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MTprLimitList, int, int, error)
	Store(tpr model.MTprLimit) (int, error)
	Update(tprLimitId int, request entity.UpdateTprLimitRequest) error
	Delete(custId string, tprLimitId int, deletedBy int64) error
}

func (repository *tprLimitRepositoryImpl) FindOneByTprLimitIdAndCustId(tprlimitID int, custId string) (model.MTprLimitList, error) {
	tpr := model.MTprLimitList{}
	query := `SELECT 
				tl.*, p.pro_code, p.pro_name, 
				u.user_fullname AS updated_by_name 
			  FROM mst.m_tpr_limit tl
			  LEFT JOIN sys.m_user u ON u.user_id = tl.updated_by
			  LEFT JOIN mst.m_product p ON p.pro_id = tl.pro_id 
			  WHERE tl.tpr_limit_id = $1 
			  AND tl.cust_id = $2`
	err := repository.Get(&tpr, query, tprlimitID, custId)
	if err != nil {
		log.Println("tprRepository, FindOneByTprLimitIdAndCustId, err:", err.Error())
		return tpr, err
	}

	return tpr, nil
}

func (repository *tprLimitRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MTprLimitList, int, int, error) {

	tprLimits := []model.MTprLimitList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` tl.*, p.pro_code, p.pro_name, 
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE tl.is_del = false 
				AND tl.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					  OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	// comment first, is_active field not available in database
	// if dataFilter.IsActive != nil {
	// 	fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
	// 	if *dataFilter.IsActive == 1 {
	// 		qWhere += ` AND tl.is_active = true `
	// 	}
	// 	if *dataFilter.IsActive == 2 {
	// 		qWhere += ` AND tl.is_active = false `
	// 	}
	// }
	
	qFrom := ` FROM mst.m_tpr_limit tl
			   LEFT JOIN sys.m_user u ON u.user_id = tl.updated_by
			   LEFT JOIN mst.m_product p ON p.pro_id = tl.pro_id `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("tprLimitRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("tprLimitRepository, count total, err:", err.Error())
		return tprLimits, 0, 0, err
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
		sortBy := `pro_id`
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

	err = repository.Select(&tprLimits, querySelect)
	if err != nil {
		log.Println("tprLimitRepository, FindAllByCustId, err:", err.Error())
		return tprLimits, total, lastPage, err
	}

	return tprLimits, total, lastPage, nil
}

func (repository *tprLimitRepositoryImpl) Store(tprLimit model.MTprLimit) (int, error) {
	query :=
		`INSERT INTO mst.m_tpr_limit(
			cust_id, pro_id, tpr_type, date_start, date_end, value_limit, value_used, value_used_str, vat_type, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, 
			$9, $10, $11, $12, $13, $14, $15, 
			$16
		) RETURNING tpr_limit_id;`
	lastInsertId := tprLimit.TprLimitId
	err := repository.QueryRow(query,
		tprLimit.CustId, tprLimit.ProId, tprLimit.TprType, tprLimit.DateStart, tprLimit.DateEnd,
		tprLimit.ValueLimit, tprLimit.ValueUsed, tprLimit.ValueUsedStr, tprLimit.VatType,
		tprLimit.CreatedBy, tprLimit.CreatedAt, tprLimit.UpdatedBy,
		tprLimit.UpdatedAt, tprLimit.IsDel, tprLimit.DeletedBy, tprLimit.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("tprLimitRepository, Store, err:", err.Error())
		return tprLimit.TprLimitId, err
	}
	return tprLimit.TprLimitId, nil
}

func (repository *tprLimitRepositoryImpl) Update(tprLimitId int, request entity.UpdateTprLimitRequest) error {
	var (
		r            model.MTprLimitUpdate
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

	query := `UPDATE mst.m_tpr_limit
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND tpr_limit_id = :tpr_limit_id_old;`

	sqlPatch.Args["tpr_limit_id_old"] = tprLimitId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("tprLimitRepository, Update, err:", err.Error())
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

func (repository *tprLimitRepositoryImpl) Delete(custId string, tprLimitId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_tpr_limit
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND tpr_limit_id = :tpr_limit_id;`

	wMap := map[string]interface{}{
		"cust_id":      custId,
		"tpr_limit_id": tprLimitId,
		"deleted_by":   deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ConvGroupDetRepository, Delete, err:", err.Error())
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
