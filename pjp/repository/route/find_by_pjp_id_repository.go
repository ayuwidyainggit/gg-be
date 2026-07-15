package route

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeRepository) FindByPjpID(ctx context.Context, tx *gorm.DB, pjpID int, custID string) []model.Route {
	var routes []model.Route

	result := tx.WithContext(ctx).
		Where("pjp_id = ? AND cust_id = ?", pjpID, custID).
		Find(&routes)

	helper.ErrorPanic(result.Error)
	return routes
}
