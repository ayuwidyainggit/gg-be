package destination

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationRepository) FindAllOutletsByPjpCode(ctx context.Context, tx *gorm.DB, pjpCode int, custId string) []model.Destination {
	var routeOutlets []model.Destination
	db := tx.WithContext(ctx)
	if err := db.Where("pjp_code = ? AND cust_id = ?", pjpCode, custId).Find(&routeOutlets).Error; err != nil {
		helper.ErrorPanic(err)
	}

	return routeOutlets
}
