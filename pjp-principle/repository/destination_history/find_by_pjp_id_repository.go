package destinationhistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationHistoryRepository) FindByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.DestinationHistory {
	var routesHistory []model.DestinationHistory
	result := tx.WithContext(ctx).Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Find(&routesHistory)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return routesHistory
}
