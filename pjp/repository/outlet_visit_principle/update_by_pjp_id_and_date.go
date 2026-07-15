package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) UpdateByPjpIDandDate(ctx context.Context, tx *gorm.DB, pjpID int, date string, data model.OutletVisitListPrinciple) {
	result := tx.WithContext(ctx).
		Where("pjp_id = ? AND date = ?", pjpID, date).
		Updates(&data)

	helper.ErrorPanic(result.Error)
}
