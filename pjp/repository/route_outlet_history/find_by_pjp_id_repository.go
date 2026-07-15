package routeoutlethistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletHistoryRepository) FindByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.RouteOutletHistory {
	var routesHistory []model.RouteOutletHistory
	result := tx.WithContext(ctx).Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Find(&routesHistory)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return routesHistory
}
