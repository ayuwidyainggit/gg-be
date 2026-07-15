package pjp

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetAll(
	ctx context.Context,
	tx *gorm.DB,
	limit int,
	page int,
	filters map[string]interface{},
	currentCustomerId string,
) ([]model.Pjp, int64) {
	var (
		pjp   []model.Pjp
		total int64
		db    = tx.WithContext(ctx)
	)

	// Base query
	baseQuery := db.
		Table("pjp.permanent_journey_plans AS pjp").
		Select(`
			pjp.id, 
			pjp.pjp_code, 
			pjp.salesman_id, 
			pjp.salesman_code, 
			pjp.salesman_name, 
			pjp.operation_type, 
			pjp.team_salesman, 
			pjp.status, 
			pjp.approval_status, 
			pjp.cust_id,
			pjp.pjp_mode,
			pjp.warehouse_id,
			pjp.warehouse_name,
			COUNT(DISTINCT route_outlet.route_code) AS total_route,
			COUNT(route_outlet.id) AS total_destinations
		`).
		Joins("LEFT JOIN pjp.route_outlet AS route_outlet ON route_outlet.pjp_id = pjp.id").
		Where("pjp.cust_id = ?", currentCustomerId).
		Group(`
			pjp.id, 
			pjp.pjp_code, 
			pjp.salesman_id, 
			pjp.salesman_code, 
			pjp.salesman_name, 
			pjp.operation_type, 
			pjp.team_salesman, 
			pjp.status,
			pjp.approval_status,
			pjp.cust_id,
			pjp.pjp_mode,
			pjp.warehouse_id,
			pjp.warehouse_name
		`).
		Order("pjp.pjp_code ASC")

	// Apply filters
	baseQuery = applyFilters(db, baseQuery, filters)

	// Count query (distinct pjp.id)
	countQuery := db.
		Table("pjp.permanent_journey_plans AS pjp").
		Joins("LEFT JOIN pjp.route_outlet AS route_outlet ON route_outlet.pjp_id = pjp.id").
		Where("pjp.cust_id = ?", currentCustomerId).
		Distinct("pjp.id")

	countQuery = applyFilters(db, countQuery, filters)

	if err := countQuery.Count(&total).Error; err != nil {
		helper.ErrorPanic(err)
	}

	// Pagination
	if err := baseQuery.Scopes(response.Scopes(page, limit)).Find(&pjp).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return pjp, total
}
