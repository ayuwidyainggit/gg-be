package route

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeRepository) Create(ctx context.Context, tx *gorm.DB, route model.Route) model.Route {
	result := tx.WithContext(ctx).Create(&route)
	helper.ErrorPanic(result.Error)

	return route
}
