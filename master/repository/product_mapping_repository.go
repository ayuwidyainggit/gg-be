package repository

import (
	"context"
	"database/sql"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type ProductMappingRepository interface {
	FindDistributorSummary(dataFilter entity.ProductMappingListQueryFilter) ([]model.ProductMappingListRow, int, int, error)
	FindDetailByDistributor(dataFilter entity.ProductMappingDetailQueryFilter) ([]model.ProductMappingDetailRow, int, int, error)
	FindOneByProIDAndPrincipal(proID int64, principalCustID string) (model.ProductMappingProductRow, error)
	ExistsProCode(custID, proCode string, excludeProID int64) (bool, error)
	ExistsProName(custID, proName string, excludeProID int64) (bool, error)
	ExistsParentMappingByDistributor(distributorID int64, parentProID int64) (bool, error)
	UpdateMapping(ctx context.Context, proID int64, proCode, proName, unitID1, unitID2, unitID3 string, updatedBy int64) error
	SoftDelete(ctx context.Context, proID int64, deletedBy int64) error
	CountActiveByDistributor(distributorID int64) (int, error)
}

func NewProductMappingRepository(db *sqlx.DB) ProductMappingRepository {
	return &productMappingRepositoryImpl{DB: db}
}

type productMappingRepositoryImpl struct {
	*sqlx.DB
}

func (r *productMappingRepositoryImpl) FindDistributorSummary(dataFilter entity.ProductMappingListQueryFilter) ([]model.ProductMappingListRow, int, int, error) {
	baseFrom := `
		FROM mst.m_product p
		INNER JOIN mst.m_distributor d ON d.distributor_id = p.distributor_id
		LEFT JOIN LATERAL (
			SELECT p2.created_by, p2.created_at, p2.updated_by, p2.updated_at
			FROM mst.m_product p2
			WHERE p2.distributor_id = d.distributor_id
				AND COALESCE(p2.is_product_mapping, false) = true
				AND COALESCE(p2.is_del, false) = false
				AND p2.is_active = true
			ORDER BY p2.updated_at DESC NULLS LAST, p2.created_at DESC NULLS LAST, p2.pro_id DESC
			LIMIT 1
		) latest ON true
		LEFT JOIN sys.m_user cu ON cu.user_id = latest.created_by
		LEFT JOIN sys.m_user uu ON uu.user_id = latest.updated_by
		WHERE COALESCE(p.is_product_mapping, false) = true
			AND COALESCE(p.is_del, false) = false
			AND p.is_active = true
			AND p.cust_id LIKE $1
			AND d.parent_cust_id = $2
	`

	args := []interface{}{dataFilter.CustId + "%", dataFilter.CustId}
	conditions := strings.Builder{}
	paramIndex := 3

	if strings.TrimSpace(dataFilter.Search) != "" {
		conditions.WriteString(fmt.Sprintf(" AND (d.distributor_name ILIKE $%d OR d.distributor_code ILIKE $%d)", paramIndex, paramIndex))
		args = append(args, "%"+strings.TrimSpace(dataFilter.Search)+"%")
		paramIndex++
	}

	groupBy := `
		GROUP BY d.distributor_id, d.distributor_code, d.distributor_name,
			latest.created_by, cu.user_fullname, latest.created_at,
			latest.updated_by, uu.user_fullname, latest.updated_at
	`

	countQuery := `SELECT COUNT(*) FROM (SELECT d.distributor_id ` + baseFrom + conditions.String() + groupBy + `) grouped`
	selectQuery := `
		SELECT
			d.distributor_id,
			d.distributor_code,
			d.distributor_name,
			COUNT(p.pro_id) AS total_product,
			latest.created_by,
			COALESCE(cu.user_fullname, '') AS created_by_name,
			latest.created_at,
			latest.updated_by,
			COALESCE(uu.user_fullname, '') AS updated_by_name,
			latest.updated_at
	` + baseFrom + conditions.String() + groupBy

	sortColumn := r.resolveListSortColumn(dataFilter.SortBy)
	sortOrder := "ASC"
	if strings.EqualFold(dataFilter.SortOrder, "desc") {
		sortOrder = "DESC"
	}
	selectQuery += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	total, err := r.executeCount(countQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	if dataFilter.Limit > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	rows := []model.ProductMappingListRow{}
	if err := r.Select(&rows, selectQuery, args...); err != nil {
		log.Error("ProductMappingRepository, FindDistributorSummary, err:", err.Error())
		return nil, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)
	return rows, total, lastPage, nil
}

func (r *productMappingRepositoryImpl) FindDetailByDistributor(dataFilter entity.ProductMappingDetailQueryFilter) ([]model.ProductMappingDetailRow, int, int, error) {
	baseWhere := `
		FROM mst.m_product p
		LEFT JOIN mst.m_product parent ON parent.pro_id = p.parent_pro_id
		INNER JOIN mst.m_distributor d ON d.distributor_id = p.distributor_id
		WHERE p.distributor_id = $1
			AND d.parent_cust_id = $2
			AND COALESCE(p.is_product_mapping, false) = true
			AND p.level > 0
			AND COALESCE(p.is_del, false) = false
			AND p.is_active = true
	`
	args := []interface{}{dataFilter.DistributorId, dataFilter.CustId}

	countQuery := `SELECT COUNT(p.pro_id) ` + baseWhere
	total, err := r.executeCount(countQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	selectQuery := `
		SELECT
			p.pro_id,
			p.parent_pro_id,
			COALESCE(parent.pro_code, '') AS parent_pro_code,
			COALESCE(parent.pro_name, '') AS parent_pro_name,
			p.pro_code,
			p.pro_name,
			p.unit_id3 AS largest_uom,
			NULLIF(p.unit_id2, '') AS middle_uom,
			NULLIF(p.unit_id1, '') AS smallest_uom
	` + baseWhere + ` ORDER BY p.pro_id ASC`

	if dataFilter.Limit > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	rows := []model.ProductMappingDetailRow{}
	if err := r.Select(&rows, selectQuery, args...); err != nil {
		log.Error("ProductMappingRepository, FindDetailByDistributor, err:", err.Error())
		return nil, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)
	return rows, total, lastPage, nil
}

func (r *productMappingRepositoryImpl) FindOneByProIDAndPrincipal(proID int64, principalCustID string) (model.ProductMappingProductRow, error) {
	row := model.ProductMappingProductRow{}
	query := `
		SELECT
			p.pro_id,
			p.cust_id,
			COALESCE(p.distributor_id, 0) AS distributor_id,
			COALESCE(p.parent_pro_id, 0) AS parent_pro_id,
			p.pro_code,
			p.pro_name,
			COALESCE(p.is_product_mapping, false) AS is_product_mapping
		FROM mst.m_product p
		INNER JOIN mst.m_distributor d ON d.distributor_id = p.distributor_id
		WHERE p.pro_id = $1
			AND d.parent_cust_id = $2
			AND COALESCE(p.is_product_mapping, false) = true
			AND COALESCE(p.is_del, false) = false
			AND p.is_active = true
		LIMIT 1
	`
	err := r.Get(&row, query, proID, principalCustID)
	if err != nil {
		log.Error("ProductMappingRepository, FindOneByProIDAndPrincipal, err:", err.Error())
	}
	return row, err
}

func (r *productMappingRepositoryImpl) ExistsProCode(custID, proCode string, excludeProID int64) (bool, error) {
	if strings.TrimSpace(proCode) == "" {
		return false, nil
	}
	query := `
		SELECT COUNT(*) > 0
		FROM mst.m_product
		WHERE cust_id = $1
			AND LOWER(TRIM(pro_code)) = LOWER(TRIM($2))
			AND pro_id <> $3
			AND COALESCE(is_del, false) = false
	`
	var exists bool
	err := r.Get(&exists, query, custID, proCode, excludeProID)
	return exists, err
}

func (r *productMappingRepositoryImpl) ExistsProName(custID, proName string, excludeProID int64) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM mst.m_product
		WHERE cust_id = $1
			AND LOWER(TRIM(pro_name)) = LOWER(TRIM($2))
			AND pro_id <> $3
			AND COALESCE(is_del, false) = false
	`
	var exists bool
	err := r.Get(&exists, query, custID, proName, excludeProID)
	return exists, err
}

func (r *productMappingRepositoryImpl) ExistsParentMappingByDistributor(distributorID int64, parentProID int64) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM mst.m_product
		WHERE distributor_id = $1
			AND parent_pro_id = $2
			AND COALESCE(is_product_mapping, false) = true
			AND COALESCE(is_del, false) = false
			AND is_active = true
	`
	var exists bool
	err := r.Get(&exists, query, distributorID, parentProID)
	return exists, err
}

func (r *productMappingRepositoryImpl) UpdateMapping(ctx context.Context, proID int64, proCode, proName, unitID1, unitID2, unitID3 string, updatedBy int64) error {
	query := `
		UPDATE mst.m_product
		SET pro_code = $1,
			pro_name = $2,
			unit_id1 = $3,
			unit_id2 = $4,
			unit_id3 = $5,
			updated_by = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE pro_id = $7
			AND COALESCE(is_product_mapping, false) = true
			AND COALESCE(is_del, false) = false
	`
	tx := GetTxFromContext(ctx)
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.Exec(query, proCode, proName, unitID1, unitID2, unitID3, updatedBy, proID)
	} else {
		result, err = r.Exec(query, proCode, proName, unitID1, unitID2, unitID3, updatedBy, proID)
	}
	if err != nil {
		log.Error("ProductMappingRepository, UpdateMapping, err:", err.Error())
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (r *productMappingRepositoryImpl) SoftDelete(ctx context.Context, proID int64, deletedBy int64) error {
	query := `
		UPDATE mst.m_product
		SET is_del = true,
			is_active = false,
			deleted_by = $1,
			deleted_at = CURRENT_TIMESTAMP
		WHERE pro_id = $2
			AND COALESCE(is_product_mapping, false) = true
			AND COALESCE(is_del, false) = false
	`
	tx := GetTxFromContext(ctx)
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.Exec(query, deletedBy, proID)
	} else {
		result, err = r.Exec(query, deletedBy, proID)
	}
	if err != nil {
		log.Error("ProductMappingRepository, SoftDelete, err:", err.Error())
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (r *productMappingRepositoryImpl) CountActiveByDistributor(distributorID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM mst.m_product
		WHERE distributor_id = $1
			AND COALESCE(is_product_mapping, false) = true
			AND COALESCE(is_del, false) = false
			AND is_active = true
	`
	var total int
	err := r.Get(&total, query, distributorID)
	return total, err
}

func (r *productMappingRepositoryImpl) resolveListSortColumn(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "distributor_code":
		return "d.distributor_code"
	case "total_product":
		return "total_product"
	case "created_at":
		return "latest.created_at"
	case "updated_at":
		return "latest.updated_at"
	default:
		return "d.distributor_name"
	}
}

func (r *productMappingRepositoryImpl) executeCount(query string, args []interface{}) (int, error) {
	var total int
	if err := r.Get(&total, query, args...); err != nil {
		log.Error("ProductMappingRepository, executeCount, err:", err.Error())
		return 0, err
	}
	return total, nil
}
