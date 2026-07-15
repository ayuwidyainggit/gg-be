package destinationhistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationHistoryRepository) FindByPjpCodeRouteCodeAndDate(ctx context.Context, tx *gorm.DB, pjpCode, routeCode int, date, custId string) []model.DestinationHistory {
	var routesHistory []model.DestinationHistory
	result := tx.WithContext(ctx).Where("pjp_code = ? AND route_code = ? AND date = ? AND cust_id = ?", pjpCode, routeCode, date, custId).Find(&routesHistory)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return routesHistory
}
