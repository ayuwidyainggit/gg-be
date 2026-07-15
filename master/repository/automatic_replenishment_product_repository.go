package repository

import (
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type AutomaticReplenishmentProductRepository interface {
	FindAll(dataFilter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]*model.AutomaticReplenishmentProduct, int, int, error)
	FindAllExport(dataFilter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]*model.AutomaticReplenishmentProduct, error)
	IsExistsByProductAndDistributor(custId string, proId, distributorId int64) (bool, error)
	FindOne(params entity.DetailAutomaticReplenishmentProductParams) (*model.AutomaticReplenishmentProduct, error)
	Create(data *model.AutomaticReplenishmentProduct) (*model.AutomaticReplenishmentProduct, error)
	Update(id int64, data *model.AutomaticReplenishmentProduct) error
	Delete(custId string, id int64, deletedBy int64) error
	IsExists(id int64, custId string) (bool, error)
}

type automaticReplenishmentProductRepository struct {
	db *sqlx.DB
}

func NewAutomaticReplenishmentProductRepository(db *sqlx.DB) AutomaticReplenishmentProductRepository {
	return &automaticReplenishmentProductRepository{
		db: db,
	}
}

func (r *automaticReplenishmentProductRepository) FindAll(dataFilter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]*model.AutomaticReplenishmentProduct, int, int, error) {
	return r.findAll(dataFilter, custId, true)
}

func (r *automaticReplenishmentProductRepository) FindAllExport(dataFilter entity.AutomaticReplenishmentProductQueryFilter, custId string) ([]*model.AutomaticReplenishmentProduct, error) {
	products, _, _, err := r.findAll(dataFilter, custId, false)
	return products, err
}

func (r *automaticReplenishmentProductRepository) findAll(dataFilter entity.AutomaticReplenishmentProductQueryFilter, custId string, withPagination bool) ([]*model.AutomaticReplenishmentProduct, int, int, error) {
	var products []*model.AutomaticReplenishmentProduct
	var total int

	filterQuery := ``
	filterArgs := []interface{}{}

	if len(dataFilter.DistributorID) > 0 {
		filterQuery += ` AND arp.distributor_id IN (?)`
		filterArgs = append(filterArgs, dataFilter.DistributorID)
	}

	if len(dataFilter.ProID) > 0 {
		filterQuery += ` AND arp.pro_id IN (?)`
		filterArgs = append(filterArgs, dataFilter.ProID)
	}

	if dataFilter.Query != "" {
		filterQuery += ` AND (p.pro_code ILIKE ? OR p.pro_name ILIKE ? OR d.distributor_name ILIKE ?)`
		query := "%" + dataFilter.Query + "%"
		filterArgs = append(filterArgs, query, query, query)
	}

	countQuery := `
		SELECT COUNT(*)
		FROM mst.auto_replenishment_product arp
		JOIN mst.m_product p ON arp.pro_id = p.pro_id
		JOIN mst.m_distributor d ON arp.distributor_id = d.distributor_id
		WHERE arp.cust_id = ? AND arp.is_del IS NOT TRUE
	`
	countQuery += filterQuery
	countArgs := append([]interface{}{custId}, filterArgs...)

	query, args, err := sqlx.In(countQuery, countArgs...)
	if err != nil {
		return products, 0, 0, err
	}

	query = r.db.Rebind(query)

	err = r.db.Get(&total, query, args...)
	if err != nil {
		return products, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)

	selectQuery := `
		SELECT arp.*, p.pro_code, p.pro_name, d.distributor_code, d.distributor_name, uc.user_name AS created_by_name, ud.user_name AS updated_by_name
		FROM mst.auto_replenishment_product arp
		JOIN mst.m_product p ON arp.pro_id = p.pro_id
		JOIN mst.m_distributor d ON arp.distributor_id = d.distributor_id
		LEFT JOIN sys.m_user uc ON uc.user_id = arp.created_by
		LEFT JOIN sys.m_user ud ON ud.user_id = arp.updated_by
		WHERE arp.cust_id = ? AND arp.is_del IS NOT TRUE
	`
	selectQuery += filterQuery
	selectArgs := append([]interface{}{custId}, filterArgs...)

	// Sorting
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			sortField := sortParts[0]
			sortOrder := sortParts[1]
			if sortOrder != "asc" && sortOrder != "desc" {
				sortOrder = "desc"
			}
			selectQuery += ` ORDER BY ` + sortField + ` ` + sortOrder
		}
	} else {
		selectQuery += ` ORDER BY arp.created_at DESC`
	}

	if withPagination {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += ` LIMIT ? OFFSET ?`
		selectArgs = append(selectArgs, dataFilter.Limit, offset)
	}

	query, args, err = sqlx.In(selectQuery, selectArgs...)
	if err != nil {
		return products, 0, 0, err
	}

	query = r.db.Rebind(query)

	err = r.db.Select(&products, query, args...)

	return products, total, lastPage, nil
}

func (r *automaticReplenishmentProductRepository) FindOne(params entity.DetailAutomaticReplenishmentProductParams) (*model.AutomaticReplenishmentProduct, error) {
	product := &model.AutomaticReplenishmentProduct{}

	query := `
		SELECT arp.*, p.pro_code, p.pro_name, d.distributor_code, d.distributor_name, uc.user_name AS created_by_name, ud.user_name AS updated_by_name
		FROM mst.auto_replenishment_product arp
		JOIN mst.m_product p ON arp.pro_id = p.pro_id
		JOIN mst.m_distributor d ON arp.distributor_id = d.distributor_id
		LEFT JOIN sys.m_user uc ON uc.user_id = arp.created_by
		LEFT JOIN sys.m_user ud ON ud.user_id = arp.updated_by
		WHERE arp.id = :id AND arp.cust_id = :cust_id AND arp.is_del IS NOT TRUE
	`

	arg := map[string]interface{}{
		"id":      params.Id,
		"cust_id": params.CustId,
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	err = r.db.Get(product, query, args...)

	return product, nil
}

func (r *automaticReplenishmentProductRepository) IsExistsByProductAndDistributor(custId string, proId, distributorId int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM mst.auto_replenishment_product
			WHERE cust_id = $1
				AND pro_id = $2
				AND distributor_id = $3
				AND is_del IS NOT TRUE
		)
	`

	if err := r.db.Get(&exists, query, custId, proId, distributorId); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *automaticReplenishmentProductRepository) Create(data *model.AutomaticReplenishmentProduct) (*model.AutomaticReplenishmentProduct, error) {
	query := `
		INSERT INTO mst.auto_replenishment_product (
			cust_id, pro_id, distributor_id, limit_action, max_order_qty, max_order_type,
			min_stock_qty, min_stock_type, safety_stock_qty, safety_stock_type,
			min_order_qty, min_order_type, is_active, created_by, created_at, updated_by, updated_at
		) VALUES (
			:cust_id, :pro_id, :distributor_id, :limit_action, :max_order_qty, :max_order_type,
			:min_stock_qty, :min_stock_type, :safety_stock_qty, :safety_stock_type,
			:min_order_qty, :min_order_type, :is_active, :created_by, :created_at, :updated_by, :updated_at
		) RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return data, err
	}

	var id int64
	err = stmt.Get(&id, data)
	if err != nil {
		return data, err
	}

	data.Id = id
	return data, nil
}

func (r *automaticReplenishmentProductRepository) Update(id int64, data *model.AutomaticReplenishmentProduct) error {
	query := `
		UPDATE mst.auto_replenishment_product SET
			pro_id = :pro_id, distributor_id = :distributor_id, limit_action = :limit_action, max_order_qty = :max_order_qty, max_order_type = :max_order_type,
			min_stock_qty = :min_stock_qty, min_stock_type = :min_stock_type, safety_stock_qty = :safety_stock_qty, safety_stock_type = :safety_stock_type,
			min_order_qty = :min_order_qty, min_order_type = :min_order_type, updated_by = :updated_by, updated_at = :updated_at
		WHERE id = :id AND cust_id = :cust_id AND is_del IS NOT TRUE
	`

	arg := map[string]interface{}{
		"pro_id":            data.ProId,
		"distributor_id":    data.DistributorId,
		"limit_action":      data.LimitAction,
		"max_order_qty":     data.MaxOrderQty,
		"max_order_type":    data.MaxOrderType,
		"min_stock_qty":     data.MinStockQty,
		"min_stock_type":    data.MinStockType,
		"safety_stock_qty":  data.SafetyStockQty,
		"safety_stock_type": data.SafetyStockType,
		"min_order_qty":     data.MinOrderQty,
		"min_order_type":    data.MinOrderType,
		"updated_by":        data.UpdatedBy,
		"updated_at":        data.UpdatedAt,
		"id":                id,
		"cust_id":           data.CustId,
	}

	result, err := r.db.NamedExec(query, arg)
	if err != nil {
		log.Error("AutomaticReplenishmentProductRepository, Update error:", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("AutomaticReplenishmentProductRepository, Update RowsAffected error:", err.Error())
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("automatic replenishment product not updated")
	}

	return nil
}

func (r *automaticReplenishmentProductRepository) Delete(custId string, id int64, deletedBy int64) error {
	query := `
		UPDATE mst.auto_replenishment_product SET
			is_del = true, deleted_by = :deleted_by, deleted_at = NOW()
		WHERE id = :id AND cust_id = :cust_id AND is_del IS NOT TRUE
	`

	arg := map[string]interface{}{
		"deleted_by": deletedBy,
		"id":         id,
		"cust_id":    custId,
	}

	_, err := r.db.NamedExec(query, arg)
	if err != nil {
		log.Error("AutomaticReplenishmentProductRepository, Delete error:", err.Error())
		return err
	}

	return nil
}

func (r *automaticReplenishmentProductRepository) IsExists(id int64, custId string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM mst.auto_replenishment_product WHERE id = :id AND cust_id = :cust_id AND is_del IS NOT TRUE`

	params := map[string]interface{}{
		"id":      id,
		"cust_id": custId,
	}

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		log.Error("AutomaticReplenishmentProductRepository, IsExists error:", err.Error())
		return false, err
	}

	query = r.db.Rebind(query)
	err = r.db.Get(&count, query, args...)
	if err != nil {
		log.Error("AutomaticReplenishmentProductRepository, IsExists error:", err.Error())
		return false, err
	}

	return count > 0, nil
}
