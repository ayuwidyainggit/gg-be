package routeoutlet

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) FindAllOutletsByPjpCode(ctx context.Context, tx *gorm.DB, pjpCode int, custId string) []model.RouteOutlet {
	var routeOutlets []model.RouteOutlet
	db := tx.WithContext(ctx)
	if err := db.Where("pjp_code = ? AND cust_id = ?", pjpCode, custId).Find(&routeOutlets).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return routeOutlets
}
