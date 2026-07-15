package routeoutlethistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletHistoryRepository) DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) {
	result := tx.WithContext(ctx).Where("pjp_id = ? AND cust_id = ?", pjpId, custId).Delete(&model.RouteOutletHistory{})
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
}

func (repo *routeOutletHistoryRepository) DeleteByVisitDay(ctx context.Context, tx *gorm.DB, history model.RouteOutletHistory) {
	if history.PjpID == nil {
		return
	}
	result := tx.WithContext(ctx).
		Where("pjp_id = ? AND cust_id = ? AND year = ? AND week = ? AND date = ? AND route_code = ?",
			*history.PjpID,
			history.CustID,
			history.Year,
			history.Week,
			history.Date,
			history.RouteCode,
		).
		Delete(&model.RouteOutletHistory{})
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
}
