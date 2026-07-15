package outletvisit

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitRepository) UpdateByPjpIDandDate(ctx context.Context, tx *gorm.DB, pjpID int, date string, data model.OutletVisitList) {
	result := tx.WithContext(ctx).
		Where("pjp_id = ? AND date = ?", pjpID, date).
		Updates(&data)

	helper.ErrorPanic(result.Error)
}
