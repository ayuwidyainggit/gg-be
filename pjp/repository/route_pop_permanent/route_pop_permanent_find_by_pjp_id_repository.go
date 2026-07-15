package routepoppermanent

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routePopPermanentRepository) FindByPjpID(ctx context.Context, tx *gorm.DB, pjpID int, custId, parentCustID string) []model.RoutePopPermanent {
	var data []model.RoutePopPermanent

	result := tx.WithContext(ctx).
		Table("pjp.route_pop_permanent").
		Select(`pjp.route_pop_permanent.*,
			COALESCE(
				pjp.route_pop_permanent.working_day_calendar_id,
				(
				SELECT mw.working_day_calendar_id
				FROM mst.m_week mw
				WHERE
					(mw.cust_id = pjp.route_pop_permanent.cust_id OR (? <> '' AND mw.cust_id = ?))
					AND mw.per_year = pjp.route_pop_permanent.year
					AND mw.week_id = pjp.route_pop_permanent.week
				ORDER BY CASE WHEN mw.cust_id = pjp.route_pop_permanent.cust_id THEN 0 ELSE 1 END
				LIMIT 1
				)
			) AS working_day_calendar_id`, parentCustID, parentCustID).
		Where("pjp.route_pop_permanent.pjp_id = ? AND pjp.route_pop_permanent.cust_id = ?", pjpID, custId).
		Find(&data)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return data
}
