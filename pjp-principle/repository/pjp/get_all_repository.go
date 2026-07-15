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
		Table("pjp_principles.permanent_journey_plans AS pjp_principles").
		Select(`
		pjp_principles.id, 
		pjp_principles.pjp_code, 
		pjp_principles.salesman_id, 
		pjp_principles.salesman_code, 
		pjp_principles.salesman_name, 
		pjp_principles.operation_type, 
		pjp_principles.team_salesman, 
		pjp_principles.status, 
		pjp_principles.approval_status, 
		pjp_principles.cust_id,
		pjp_principles.pjp_mode,
		pjp_principles.warehouse_id,
		pjp_principles.warehouse_name,
		COUNT(DISTINCT destinations.route_code) AS total_route,
		COUNT(destinations.id) AS total_destinations,
		COUNT(CASE WHEN destinations.destination_type = 'outlet' THEN 1 END) AS total_outlet,
		COUNT(CASE WHEN destinations.destination_type = 'distributor' THEN 1 END) AS total_distributor
	`).
		Joins("LEFT JOIN pjp_principles.destinations AS destinations ON destinations.pjp_id = pjp_principles.id").
		Where("pjp_principles.cust_id = ?", currentCustomerId).
		Group(`
		pjp_principles.id, 
		pjp_principles.pjp_code, 
		pjp_principles.salesman_id, 
		pjp_principles.salesman_code, 
		pjp_principles.salesman_name, 
		pjp_principles.operation_type, 
		pjp_principles.team_salesman, 
		pjp_principles.status,
		pjp_principles.approval_status,
		pjp_principles.cust_id,
		pjp_principles.pjp_mode,
		pjp_principles.warehouse_id,
		pjp_principles.warehouse_name
	`).
		Order("pjp_principles.pjp_code ASC")

	// Apply filters
	baseQuery = applyFilters(db, baseQuery, filters)

	// Count query (distinct pjp_principles.id)
	countQuery := db.
		Table("pjp_principles.permanent_journey_plans AS pjp_principles").
		Joins("LEFT JOIN pjp_principles.destinations AS destinations ON destinations.pjp_id = pjp_principles.id").
		Where("pjp_principles.cust_id = ?", currentCustomerId).
		Distinct("pjp_principles.id")

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
