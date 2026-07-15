package destination

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationRepository) FindAllOutletsByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.Destination {
	var routeOutlets []model.Destination
	db := tx.WithContext(ctx)
	if err := db.Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Find(&routeOutlets).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return routeOutlets
}
