package routeoutlet

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) FindAllOutletsByPjpIds(ctx context.Context, tx *gorm.DB, pjpId []int, custId string) []model.RouteOutlet {
	var routeOutlets []model.RouteOutlet
	db := tx.WithContext(ctx)
	if err := db.Where("pjp_id IN ? AND cust_id = ?", pjpId, custId).Find(&routeOutlets).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return routeOutlets
}
