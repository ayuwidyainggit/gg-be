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

type TakingOrderRepository interface {
	FindOneByReasonAndCustId(reason, custId string) (model.TakingOrder, error)
	FindOneByTakingOrderIdAndCustId(TakingOrderId int, custId string) (model.TakingOrder, error)
	FindAllByCustId(dataFilter entity.TakingOrderQueryFilter, custId string) (consPro []model.TakingOrder, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.TakingOrderQueryFilter, custId string) (consPro []model.TakingOrder, total int, lastPage int, err error)
	Store(brand model.TakingOrder) (int, error)
	Update(TakingOrderId int, request entity.UpdateTakingOrderRequest) error
	Delete(custId string, TakingOrderId int, deletedBy int64) error
}

func NewTakingOrderRepository(db *sqlx.DB) TakingOrderRepository {
	return &TakingOrderRepositoryImpl{db}
}

type TakingOrderRepositoryImpl struct {
	*sqlx.DB
}

func (repository *TakingOrderRepositoryImpl) FindOneByReasonAndCustId(reason string, custId string) (model.TakingOrder, error) {
	TakingOrder := model.TakingOrder{}
	query := `SELECT 
				tko.taking_order_id, 
				tko.taking_order_name,
				COALESCE(tko.image_url, '') as image_url,
				tko.is_active, tko.created_by,
				tko.created_at, tko.updated_by, tko.updated_at,
				tko.is_del, tko.deleted_by, tko.deleted_at
			  FROM mst.m_taking_order tko
			  WHERE tko.is_del = false 
				AND tko.taking_order_name = $1 
				AND tko.cust_id = $2`
	err := repository.Get(&TakingOrder, query, reason, custId)
	if err != nil {
		log.Println("takingOrderRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return TakingOrder, err
	}

	return TakingOrder, nil
}

func (repository *TakingOrderRepositoryImpl) FindOneByTakingOrderIdAndCustId(TakingOrderId int, custId string) (model.TakingOrder, error) {
	TakingOrder := model.TakingOrder{}
	query := `SELECT 
				tko.taking_order_id, 
				tko.taking_order_name,
				COALESCE(tko.image_url, '') as image_url,
				tko.is_active, tko.created_by,
				tko.created_at, tko.updated_by, tko.updated_at,
				tko.is_del, tko.deleted_by, tko.deleted_at
			  FROM mst.m_taking_order tko
			  WHERE tko.is_del = false 
				AND tko.taking_order_id = $1 
				AND tko.cust_id = $2`
	err := repository.Get(&TakingOrder, query, TakingOrderId, custId)
	if err != nil {
		log.Println("takingOrderRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return TakingOrder, err
	}

	return TakingOrder, nil
}

func (repository *TakingOrderRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.TakingOrderQueryFilter, custId string) ([]model.TakingOrder, int, int, error) {

	TakingOrder := []model.TakingOrder{}
	selectCount := ` COUNT(tko.*) AS total `
	selectField := `tko.taking_order_id,
					tko.taking_order_name  `
	qWhere := ` WHERE tko.is_del = false AND tko.is_active = true 
				AND tko.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (tko.taking_order_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.TakingOrderId != 0 {
		qWhere += `AND tko.taking_order_id = ` + strconv.Itoa(dataFilter.TakingOrderId) + ` `
	}

	qFrom := ` FROM mst.m_taking_order tko `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("TakingOrderRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("TakingOrderRepository, count total, err:", err.Error())
		return TakingOrder, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`tko.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `tko.taking_order_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("TakingOrderRepository, querySelect:", querySelect)
	err = repository.Select(&TakingOrder, querySelect)
	if err != nil {
		log.Println("TakingOrderRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return TakingOrder, total, 1, err
	}

	return TakingOrder, total, 1, nil
}

func (repository *TakingOrderRepositoryImpl) FindAllByCustId(dataFilter entity.TakingOrderQueryFilter, custId string) ([]model.TakingOrder, int, int, error) {

	brands := []model.TakingOrder{}
	selectCount := ` COUNT(tko.*) AS total `
	selectField := `tko.taking_order_id,
					tko.taking_order_name,
					COALESCE(tko.image_url, '') as image_url,
					tko.is_active, tko.created_by,
					tko.created_at, tko.updated_by, tko.updated_at,
					tko.is_del, tko.deleted_by, tko.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE tko.is_del = false 
				AND tko.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (tko.taking_order_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND tko.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND tko.is_active = false `
		}
	}

	if dataFilter.TakingOrderId != 0 {
		qWhere += `AND tko.taking_order_id = ` + strconv.Itoa(dataFilter.TakingOrderId) + ` `
	}

	qFrom := ` FROM mst.m_taking_order tko
			   LEFT JOIN sys.m_user u ON u.user_id = tko.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("TakingOrderRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("TakingOrderRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`tko.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `tko.taking_order_id`
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

	// log.Println("TakingOrderRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("TakingOrderRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *TakingOrderRepositoryImpl) Store(TakingOrder model.TakingOrder) (int, error) {
	query :=
		`INSERT INTO mst.m_taking_order(
			cust_id, taking_order_name, image_url, 
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11 
		) RETURNING taking_order_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		TakingOrder.CustId, TakingOrder.TakingOrderName, TakingOrder.ImageUrl, TakingOrder.IsActive,
		TakingOrder.CreatedBy, TakingOrder.CreatedAt, TakingOrder.UpdatedBy, TakingOrder.UpdatedAt,
		TakingOrder.IsDel, TakingOrder.DeletedBy, TakingOrder.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("TakingOrderRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *TakingOrderRepositoryImpl) Update(TakingOrderId int, request entity.UpdateTakingOrderRequest) error {
	var (
		r            model.TakingOrderUpdate
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

	query := `UPDATE mst.m_taking_order
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND taking_order_id = :taking_order_id;`

	sqlPatch.Args["taking_order_id"] = TakingOrderId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("TakingOrderRepository, Update, err:", err.Error())
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

func (repository *TakingOrderRepositoryImpl) Delete(custId string, TakingOrderId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_taking_order
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND taking_order_id = :taking_order_id;`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"taking_order_id": TakingOrderId,
		"deleted_by":      deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("TakingOrderRepository, Delete, err:", err.Error())
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
