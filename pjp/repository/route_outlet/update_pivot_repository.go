package routeoutlet

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) UpdatePivot(ctx context.Context, tx *gorm.DB, route model.RouteOutlet) {
	dataset := model.RouteOutlet{
		ID: route.ID,
		// RouteCode:    route.RouteCode,
		// OutletCode:   route.OutletCode,
		Status:       route.Status,
		VerifiedDate: route.VerifiedDate,
	}
	result := tx.Model(&route).WithContext(ctx).
		Where("id = ?", dataset.ID).
		// Where("route_code = ?", dataset.RouteCode).
		// Where("outlet_code = ?", dataset.OutletCode).
		Updates(dataset)
	helper.ErrorPanic(result.Error)
}
