package destination

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationRepository) GetAll(ctx context.Context, tx *gorm.DB, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]model.Destination, int) {
	var routes []model.Destination
	var totalData int64

	baseQuery := tx.Model(&model.Destination{}).
		Where("pjp_principles.destinations.cust_id = ?", currentCustomerId)

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
