package outlet_visit_principle

import (
	"context"
	"gorm.io/gorm"
	"scyllax-pjp/helper"
)

func (repo *outletVisitPrincipleRepository) UpdateOutletVisitListColumnAt(ctx context.Context, tx *gorm.DB, column string, currentTime *int64, date string, id int64) {

	result := tx.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET `+column+` = ?
		WHERE date = ? AND id = ?
	`, currentTime, date, id)

	helper.ErrorPanic(result.Error)
}
