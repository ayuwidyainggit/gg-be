package destinationhistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationHistoryRepository) DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) {
	result := tx.WithContext(ctx).Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Delete(&model.DestinationHistory{})
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
}
