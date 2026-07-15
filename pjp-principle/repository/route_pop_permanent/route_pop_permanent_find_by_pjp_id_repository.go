package routepoppermanent

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routePopPermanentRepository) FindByPjpID(ctx context.Context, tx *gorm.DB, pjpID int, custId string) []model.RoutePopPermanent {
	var data []model.RoutePopPermanent

	result := tx.WithContext(ctx).Where("pjp_id = ? AND cust_id = ?", pjpID, custId).Find(&data)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return data
}
