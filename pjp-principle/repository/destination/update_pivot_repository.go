package destination

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationRepository) UpdatePivot(ctx context.Context, tx *gorm.DB, route model.Destination) {
	dataset := model.Destination{
		ID: route.ID,
		// RouteCode:    route.RouteCode,
		// DestinationCode:   route.DestinationCode,
		Status:       route.Status,
		VerifiedDate: route.VerifiedDate,
	}
	result := tx.Model(&route).WithContext(ctx).
		Where("id = ?", dataset.ID).
		// Where("route_code = ?", dataset.RouteCode).
		// Where("outlet_code = ?", dataset.DestinationCode).
		Updates(dataset)
	helper.ErrorPanic(result.Error)
}
