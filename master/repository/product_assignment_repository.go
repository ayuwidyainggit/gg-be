package repository

import (
	"context"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type ProductAssignmentRepository interface {
	FindAll(dataFilter entity.ProductAssignmentQueryFilter) ([]model.ProductAssignmentHistory, int, int, error)
	FindAllHistory(dataFilter entity.ProductAssignmentQueryFilter) ([]model.ProductAssignmentHistory, error)
	Store(ctx context.Context, assignment model.ProductAssignmentMutation) (int64, error)
}

func NewProductAssignmentRepository(db *sqlx.DB) ProductAssignmentRepository {
	return &productAssignmentRepositoryImpl{db}
}

type productAssignmentRepositoryImpl struct {
	*sqlx.DB
}

func (r *productAssignmentRepositoryImpl) FindAll(dataFilter entity.ProductAssignmentQueryFilter) ([]model.ProductAssignmentHistory, int, int, error) {
	selectQuery, countQuery, args := r.buildHistoryQueries(dataFilter)

	if dataFilter.Limit > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	total, err := r.executeCountQuery(countQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	assignments, err := r.executeHistorySelectQuery(selectQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)
	return assignments, total, lastPage, nil
}

func (r *productAssignmentRepositoryImpl) FindAllHistory(dataFilter entity.ProductAssignmentQueryFilter) ([]model.ProductAssignmentHistory, error) {
	query, _, args := r.buildHistoryQueries(dataFilter)

	var rows []model.ProductAssignmentHistory
	if err := r.Select(&rows, query, args...); err != nil {
		log.Error(err)
		return nil, err
	}

	return rows, nil
}

func (r *productAssignmentRepositoryImpl) Store(ctx context.Context, assignment model.ProductAssignmentMutation) (int64, error) {
	query := `
		INSERT INTO mst.m_product_assignment (
			cust_id,
			action_date,
			pro_id,
			distributor_id,
			assignment_type,
			created_by,
			created_at
		) VALUES (
			:cust_id,
			:action_date,
			:pro_id,
			:distributor_id,
			:assignment_type,
			:created_by,
			:created_at
		) RETURNING id;
	`

	tx := GetTxFromContext(ctx)
	var (
		rows *sqlx.Rows
		err  error
		id   int64
	)

	if tx != nil {
		rows, err = tx.NamedQuery(query, assignment)
	} else {
		rows, err = r.NamedQuery(query, assignment)
	}
	if err != nil {
		log.Error("ProductAssignmentRepository, Store, err:", err.Error())
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *productAssignmentRepositoryImpl) buildHistoryQueries(dataFilter entity.ProductAssignmentQueryFilter) (string, string, []interface{}) {
	baseQuery := `
		FROM mst.m_product_assignment pa
		INNER JOIN mst.m_product p ON p.pro_id = pa.pro_id
		INNER JOIN mst.m_distributor d ON d.distributor_id = pa.distributor_id
		LEFT JOIN sys.m_user u ON u.user_id = pa.created_by
		WHERE pa.cust_id = $1
	`

	var conditions strings.Builder
	args := []interface{}{dataFilter.CustId}
	paramIndex := 2

	if len(dataFilter.QDistributorId) > 0 {
		placeholders := r.buildPlaceholders(len(dataFilter.QDistributorId), paramIndex)
		conditions.WriteString(fmt.Sprintf(" AND pa.distributor_id IN (%s)", placeholders))
		for _, id := range dataFilter.QDistributorId {
			args = append(args, id)
		}
		paramIndex += len(dataFilter.QDistributorId)
	}

	if len(dataFilter.ProId) > 0 {
		placeholders := r.buildPlaceholders(len(dataFilter.ProId), paramIndex)
		conditions.WriteString(fmt.Sprintf(" AND pa.pro_id IN (%s)", placeholders))
		for _, id := range dataFilter.ProId {
			args = append(args, id)
		}
		paramIndex += len(dataFilter.ProId)
	}

	if len(dataFilter.AssignmentType) > 0 {
		placeholders := r.buildPlaceholders(len(dataFilter.AssignmentType), paramIndex)
		conditions.WriteString(fmt.Sprintf(" AND LOWER(pa.assignment_type) IN (%s)", placeholders))
		for _, assignmentType := range dataFilter.AssignmentType {
			args = append(args, strings.ToLower(assignmentType))
		}
		paramIndex += len(dataFilter.AssignmentType)
	}

	if dataFilter.Query != "" {
		searchTerm := "%" + dataFilter.Query + "%"
		conditions.WriteString(fmt.Sprintf(
			" AND (p.pro_code ILIKE $%d OR p.pro_name ILIKE $%d OR d.distributor_code ILIKE $%d OR d.distributor_name ILIKE $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3,
		))
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	selectQuery := `
		SELECT
			pa.id,
			pa.cust_id,
			pa.action_date,
			pa.pro_id,
			p.pro_code,
			p.pro_name,
			pa.distributor_id,
			d.distributor_code,
			d.distributor_name,
			pa.assignment_type,
			pa.created_by,
			COALESCE(NULLIF(u.user_fullname, ''), u.user_name, '') AS created_by_name,
			pa.created_at
	` + baseQuery + conditions.String() + " ORDER BY pa.created_at DESC, pa.id DESC"

	countQuery := "SELECT COUNT(*) " + baseQuery + conditions.String()
	return selectQuery, countQuery, args
}

func (r *productAssignmentRepositoryImpl) buildPlaceholders(count int, startIndex int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = fmt.Sprintf("$%d", startIndex+i)
	}
	return strings.Join(placeholders, ", ")
}

func (r *productAssignmentRepositoryImpl) buildFindAllOrderBy(sort string) string {
	sortColumn := r.getSortColumn(sort)
	return " ORDER BY " + sortColumn + " DESC"
}

func (r *productAssignmentRepositoryImpl) getSortColumn(sort string) string {
	columnMap := map[string]string{
		"pro_code":         "p.pro_code",
		"distributor_code": "d.distributor_code",
		"created_at":       "p.created_at",
		"updated_at":       "p.updated_at",
		"pro_name":         "p.pro_name",
		"distributor_name": "d.distributor_name",
		"created_by_name":  "u.user_name",
	}

	if col, exists := columnMap[sort]; exists {
		return col
	}
	return "p.created_at"
}

func (r *productAssignmentRepositoryImpl) executeCountQuery(query string, args []interface{}) (int, error) {
	var total int
	if err := r.Get(&total, query, args...); err != nil {
		log.Error(err)
		return 0, err
	}
	return total, nil
}

func (r *productAssignmentRepositoryImpl) executeHistorySelectQuery(query string, args []interface{}) ([]model.ProductAssignmentHistory, error) {
	var assignments []model.ProductAssignmentHistory
	if err := r.Select(&assignments, query, args...); err != nil {
		log.Error(err)
		return nil, err
	}
	return assignments, nil
}
