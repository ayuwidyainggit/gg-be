package pjp

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpRepository) Patch(ctx context.Context, tx *gorm.DB, pjpId int, pjpMode, custId string) {
	result := tx.Model(&model.Pjp{}).WithContext(ctx).
		Where("id = ? AND cust_id", pjpId, custId).
		Update("pjp_mode", pjpMode)
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
	if result.RowsAffected == 0 {
		helper.ErrorPanic(errors.New("no record found"))
	}
}
