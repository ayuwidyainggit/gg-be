package routeoutlet

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.RouteOutlet) {
	if len(outlets) == 0 {
		return
	}

	result := tx.WithContext(ctx).Create(&outlets)
	helper.ErrorPanic(result.Error)
}
