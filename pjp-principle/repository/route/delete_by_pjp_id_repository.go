package route

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeRepository) DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) {
	db := tx.WithContext(ctx)
	if err := db.Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Delete(&model.Route{}).Error; err != nil {
		helper.ErrorPanic(err)
	}
}
