package pjp

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"strconv"

	"gorm.io/gorm"
)

func (repo *pjpRepository) GetPjpWithRoute(ctx context.Context, tx *gorm.DB, q string, custId string) []response.PjpWithRouteRow {
	var rows []response.PjpWithRouteRow

	query := tx.WithContext(ctx).
		Table("pjp.permanent_journey_plans").
		Select(`
			pjp.permanent_journey_plans.*,
			pjp_distinct.route_code,
			pjp_distinct.outlet_id,
			COUNT(*) FILTER (
				WHERE pjp_distinct.status != 'Reject'
				AND (pjp_distinct.status = 'Approved' OR pjp_distinct.status = 'Approved With Propose')
			) AS total_outlet
		`).
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT pjp_id, route_code, outlet_id, status, verified_date
				FROM pjp.route_outlet
			) AS pjp_distinct ON pjp_distinct.pjp_id = pjp.permanent_journey_plans.id
		`).
		Where("pjp.permanent_journey_plans.cust_id = ?", custId).
		Group("pjp.permanent_journey_plans.id, pjp_distinct.route_code, pjp_distinct.outlet_id")

	if q != "" {
		pjpCode, err := strconv.Atoi(q)
		helper.ErrorPanic(err)
		query = query.Where("pjp.permanent_journey_plans.pjp_code = ?", pjpCode)
	}

	result := query.Find(&rows)
	helper.ErrorPanic(result.Error)
	return rows
}
