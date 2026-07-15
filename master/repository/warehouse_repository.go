package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WarehouseRepository interface {
	FindOneByWarehouseIdAndCustId(warehouseId int, custId string) (model.Warehouse, error)
	FindOneByWarehouseCodeAndCustId(warehouseCode string, custId string) (model.Warehouse, error)
	FindAllByCustId(dataFilter entity.WarehouseQueryFilter, custId string) (consPro []model.Warehouse, total int, lastPage int, err error)
	FindAllLookupByCustId(dataFilter entity.WarehouseQueryFilter, custId string) (consPro []model.Warehouse, total int, lastPage int, err error)
	FindLookupByDistributorUnion(dataFilter entity.WarehouseQueryFilter, distributorIDs []int) (consPro []model.Warehouse, total int, lastPage int, err error)
	Store(warehouse model.Warehouse) (int, error)
	Update(warehouseId int, request entity.UpdateWarehouseRequest) error
	Delete(custId string, warehouseId int, deletedBy int64) error
	CountWarehouseInStockByCustId(warehouseId int, custId string) (model.TotalWarehouse, error)
}

func NewWarehouseRepository(db *sqlx.DB) WarehouseRepository {
	return &warehouseRepositoryImpl{db}
}

type warehouseRepositoryImpl struct {
	*sqlx.DB
}

func (repository *warehouseRepositoryImpl) FindOneByWarehouseIdAndCustId(warehouseId int, custId string) (model.Warehouse, error) {
	warehouse := model.Warehouse{}
	query := `SELECT 
				cust_id, wh_id, wh_code, stock_type, latitude, longitude,
				wh_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_warehouse 
			  WHERE wh_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&warehouse, query, warehouseId, custId)
	if err != nil {
		log.Println("warehouseRepository, FindOneByWarehouseCodeAndCustId, err:", err.Error())
		return warehouse, err
	}

	return warehouse, nil
}

func (repository *warehouseRepositoryImpl) FindOneByWarehouseCodeAndCustId(warehouseCode string, custId string) (model.Warehouse, error) {
	warehouse := model.Warehouse{}
	query := `SELECT 
				cust_id, wh_id, wh_code, stock_type,
				wh_name, is_active, 
				created_by, created_at, updated_by, 
				updated_at, is_del, deleted_by, deleted_at
			  FROM mst.m_warehouse 
			  WHERE cust_id = $2  
			  AND wh_code = $1
			  AND is_del = false`
	err := repository.Get(&warehouse, query, warehouseCode, custId)
	if err != nil {
		log.Println("warehouseRepository, FindOneByWarehouseCodeAndCustId, err:", err.Error())
		return warehouse, err
	}

	return warehouse, nil
}

func (repository *warehouseRepositoryImpl) FindAllByCustId(dataFilter entity.WarehouseQueryFilter, custId string) ([]model.Warehouse, int, int, error) {

	warehouses := []model.Warehouse{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.wh_id, a.wh_code, a.wh_name, a.stock_type, a.latitude, a.longitude,
					a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false `
	if dataFilter.IsReplenishment == nil || !*dataFilter.IsReplenishment {
		qWhere += ` AND a.cust_id = '` + custId + `' `
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (a.wh_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.wh_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.StockType != "" {
		qWhere += ` AND a.stock_type = '` + dataFilter.StockType + `'`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	if dataFilter.Gudang != nil {
		fmt.Println("dataFilter.Gudang:", *dataFilter.Gudang)
		if *dataFilter.Gudang == 1 {
			qWhere += ` AND a.latitude IS NOT NULL AND a.longitude IS NOT NULL `
		}
	}

	if len(dataFilter.DistributorIDs) > 0 {
		ids := make([]string, 0, len(dataFilter.DistributorIDs))
		for _, id := range dataFilter.DistributorIDs {
			ids = append(ids, fmt.Sprintf("%d", id))
		}
		distributorIn := strings.Join(ids, ",")
		qWhere += ` AND (
			EXISTS (
				SELECT 1
				FROM sys.m_user mu
				INNER JOIN mst.m_salesman ms ON ms.emp_id = mu.emp_id
				INNER JOIN smc.m_customer mc ON mc.cust_id = mu.cust_id
				WHERE ms.wh_id = a.wh_id
					AND mc.distributor_id IN (` + distributorIn + `)
			)
			OR EXISTS (
				SELECT 1
				FROM sys.m_user mu
				INNER JOIN mst.m_salesman_canvas msc ON msc.emp_id = mu.emp_id
				INNER JOIN smc.m_customer mc ON mc.cust_id = mu.cust_id
				WHERE msc.wh_id = a.wh_id
					AND mc.distributor_id IN (` + distributorIn + `)
			)
			OR EXISTS (
				SELECT 1
				FROM smc.m_customer mcw
				WHERE mcw.cust_id = a.cust_id
					AND mcw.distributor_id IN (` + distributorIn + `)
			)
		) `
	}

	qFrom := ` FROM mst.m_warehouse a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	log.Println("warehouseRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("warehouseRepository, count total, err:", err.Error())
		return warehouses, 0, 0, err
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
		sortBy := `wh_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	// log.Println("warehouseRepository, querySelect:", querySelect)
	err = repository.Select(&warehouses, querySelect)
	if err != nil {
		log.Println("warehouseRepository, FindAllByCustId, err:", err.Error())
		return warehouses, total, 0, err
	}

	return warehouses, total, 0, nil
}

func (repository *warehouseRepositoryImpl) FindAllLookupByCustId(dataFilter entity.WarehouseQueryFilter, custId string) ([]model.Warehouse, int, int, error) {

	warehouses := []model.Warehouse{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` cust_id, wh_id, wh_code, wh_name, stock_type,
					is_active, created_by,
					created_at, updated_by, updated_at,
					is_del, deleted_by, deleted_at `
	qWhere := ` WHERE is_del = false AND is_active = true `
	if dataFilter.IsReplenishment == nil || !*dataFilter.IsReplenishment {
		qWhere += ` AND cust_id = '` + custId + `' `
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (wh_code ILIKE '%` + dataFilter.Query + `%' 
					OR wh_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.StockType != "" {
		qWhere += ` AND stock_type = '` + dataFilter.StockType + `'`
	}

	qFrom := ` FROM mst.m_warehouse `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("warehouseRepository, count total, err:", err.Error())
		return warehouses, 0, 0, err
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
		sortBy := `wh_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	lastPage := sql_helper.CalculateLastPage(total, limit)

	err = repository.Select(&warehouses, querySelect)
	if err != nil {
		log.Println("warehouseRepository, FindAllLookupByCustId, err:", err.Error())
		return warehouses, total, 0, err
	}

	return warehouses, total, lastPage, nil
}

func warehouseLookupStockTypeArg(st string) string {
	st = strings.TrimSpace(st)
	if st == "G" || st == "E" || st == "BS" {
		return st
	}
	return ""
}

func buildWarehouseLookupOrderSQL(sort string) string {
	if sort == "" {
		return "ORDER BY u.wh_name ASC"
	}
	parts := strings.Split(sort, ",")
	var orders []string
	columnMap := map[string]string{
		"wh_id":      "u.wh_id",
		"wh_name":    "u.wh_name",
		"wh_code":    "u.wh_code",
		"stock_type": "u.stock_type",
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		segs := strings.SplitN(p, ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "ASC"
		if len(segs) > 1 {
			d := strings.ToUpper(strings.TrimSpace(segs[1]))
			if d == "ASC" || d == "DESC" {
				dir = d
			}
		}
		col, ok := columnMap[key]
		if !ok {
			continue
		}
		orders = append(orders, col+" "+dir)
	}
	if len(orders) == 0 {
		return "ORDER BY u.wh_name ASC"
	}
	return "ORDER BY " + strings.Join(orders, ", ")
}

func (repository *warehouseRepositoryImpl) FindLookupByDistributorUnion(dataFilter entity.WarehouseQueryFilter, distributorIDs []int) ([]model.Warehouse, int, int, error) {
	warehouses := []model.Warehouse{}
	if len(distributorIDs) == 0 {
		return warehouses, 0, 1, nil
	}

	stockArg := warehouseLookupStockTypeArg(dataFilter.StockType)
	baseUnion := `
(
	SELECT DISTINCT mw.wh_id, mw.wh_code, mw.wh_name, mw.stock_type, mw.latitude, mw.longitude
	FROM sys.m_user mu
	INNER JOIN mst.m_salesman ms ON ms.emp_id = mu.emp_id
	INNER JOIN mst.m_warehouse mw ON mw.wh_id = ms.wh_id
	INNER JOIN smc.m_customer mc ON mc.cust_id = mu.cust_id
	WHERE mc.distributor_id = ANY($1::int[])
		AND COALESCE(mw.is_del, false) = false
		AND ($2::text = '' OR mw.stock_type = $2)
)
UNION
(
	SELECT DISTINCT mw.wh_id, mw.wh_code, mw.wh_name, mw.stock_type, mw.latitude, mw.longitude
	FROM sys.m_user mu
	INNER JOIN mst.m_salesman_canvas msc ON msc.emp_id = mu.emp_id
	INNER JOIN mst.m_warehouse mw ON mw.wh_id = msc.wh_id
	INNER JOIN smc.m_customer mc ON mc.cust_id = mu.cust_id
	WHERE mc.distributor_id = ANY($1::int[])
		AND COALESCE(mw.is_del, false) = false
		AND ($2::text = '' OR mw.stock_type = $2)
)
UNION
(
	SELECT DISTINCT mw.wh_id, mw.wh_code, mw.wh_name, mw.stock_type, mw.latitude, mw.longitude
	FROM mst.m_warehouse mw
	INNER JOIN smc.m_customer mc ON mc.cust_id = mw.cust_id
	WHERE mc.distributor_id = ANY($1::int[])
		AND COALESCE(mw.is_del, false) = false
		AND ($2::text = '' OR mw.stock_type = $2)
)`

	args := []interface{}{pq.Array(distributorIDs), stockArg}

	countQuery := `SELECT COUNT(*) FROM (` + baseUnion + `) AS u`
	var total int
	if err := repository.Get(&total, countQuery, args...); err != nil {
		log.Println("warehouseRepository, FindLookupByDistributorUnion count, err:", err.Error())
		return warehouses, 0, 0, err
	}

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	orderSQL := buildWarehouseLookupOrderSQL(dataFilter.Sort)
	selectQuery := `
SELECT u.wh_id, u.wh_code, u.wh_name, u.stock_type, u.latitude, u.longitude
FROM (` + baseUnion + `) AS u
` + orderSQL + fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	if err := repository.Select(&warehouses, selectQuery, args...); err != nil {
		log.Println("warehouseRepository, FindLookupByDistributorUnion select, err:", err.Error())
		return warehouses, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, limit)
	return warehouses, total, lastPage, nil
}

func (repository *warehouseRepositoryImpl) Store(warehouse model.Warehouse) (int, error) {
	query :=
		`INSERT INTO mst.m_warehouse(
			cust_id, wh_code, wh_name, stock_type, latitude, longitude,
			is_active, created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12, $13, $14
		) RETURNING wh_id;`
	lastInsertId := warehouse.WarehouseId
	err := repository.QueryRow(query,
		warehouse.CustId, warehouse.WarehouseCode, warehouse.WarehouseName, warehouse.StockType, warehouse.Latitude, warehouse.Longitude,
		warehouse.IsActive, warehouse.CreatedBy, warehouse.CreatedAt, warehouse.UpdatedBy, warehouse.UpdatedAt,
		warehouse.IsDel, warehouse.DeletedBy, warehouse.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("warehouseRepository, Store, err:", err.Error())
		return warehouse.WarehouseId, err
	}
	return warehouse.WarehouseId, nil
}

func (repository *warehouseRepositoryImpl) Update(warehouseId int, request entity.UpdateWarehouseRequest) error {
	var (
		r            model.WarehouseUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("warehouseRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_warehouse
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND wh_id = :wh_id_old;`

	// log.Println("warehouseRepository, Update, query:", query)

	sqlPatch.Args["wh_id_old"] = warehouseId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("warehouseRepository, Update, err:", err.Error())
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

func (repository *warehouseRepositoryImpl) Delete(custId string, warehouseId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_warehouse
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND wh_id = :wh_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"wh_id":      warehouseId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("WarehouseRepository, Delete, err:", err.Error())
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

func (repository *warehouseRepositoryImpl) CountWarehouseInStockByCustId(warehouseId int, custId string) (model.TotalWarehouse, error) {
	warehouse := model.TotalWarehouse{}
	query := `SELECT COUNT(wh_id) AS total 
			  FROM inv.warehouse_stock 
			  WHERE cust_id = $2 AND wh_id = $1`
	err := repository.Get(&warehouse, query, warehouseId, custId)
	if err != nil {
		log.Println("warehouseRepository, CountWarehouseInStockByCustId, err:", err.Error())
		return warehouse, err
	}

	return warehouse, nil
}
