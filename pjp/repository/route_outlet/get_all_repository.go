package routeoutlet

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) GetAll(ctx context.Context, tx *gorm.DB, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]model.RouteOutlet, int) {
	var routes []model.RouteOutlet
	var totalData int64

	baseQuery := tx.Model(&model.RouteOutlet{}).
		Where("pjp.route_outlet.cust_id = ?", currentCustomerId)

	for field, value := range filters {
		baseQuery = applyFilter(baseQuery, field, value)
	}

	// Hitung total data sebelum paginasi
	countQuery := baseQuery.Session(&gorm.Session{}) // clone tanpa Preload & Scopes
	result := countQuery.Count(&totalData)
	helper.ErrorPanic(result.Error)

	// Ambil data dengan preload dan paginasi
	query := baseQuery.Preload("Pjp").Preload("PjpOld").
		WithContext(ctx).
		Scopes(response.Scopes(page, limit)).
		Find(&routes)

	helper.ErrorPanic(query.Error)

	return routes, int(totalData)
}
