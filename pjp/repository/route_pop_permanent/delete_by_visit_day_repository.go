package routepoppermanent

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routePopPermanentRepository) DeleteByVisitDay(ctx context.Context, tx *gorm.DB, routePopPermanent model.RoutePopPermanent) {
	if routePopPermanent.PjpID == nil || routePopPermanent.RouteCode == nil {
		return
	}
	result := tx.WithContext(ctx).
		Where("pjp_id = ? AND cust_id = ? AND year = ? AND week = ? AND date = ? AND route_code = ?",
			*routePopPermanent.PjpID,
			routePopPermanent.CustID,
			routePopPermanent.Year,
			routePopPermanent.Week,
			routePopPermanent.Date,
			*routePopPermanent.RouteCode,
		).
		Delete(&model.RoutePopPermanent{})
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}
}
