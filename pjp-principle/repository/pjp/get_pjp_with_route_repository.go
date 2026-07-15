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
		Table("pjp_principles.permanent_journey_plans").
		Select(`
			pjp_principles.permanent_journey_plans.*,
			pjp_distinct.route_code,
			COUNT(*) AS total_outlet
		`).
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT pjp_id, route_code, destination_id, status, verified_date
				FROM pjp_principles.destinations
			) AS pjp_distinct ON pjp_distinct.pjp_id = pjp_principles.permanent_journey_plans.id
		`).
		Where("pjp_principles.permanent_journey_plans.cust_id = ?", custId).
		Group("pjp_principles.permanent_journey_plans.id, pjp_distinct.route_code")

	if q != "" {
		pjpCode, err := strconv.Atoi(q)
		helper.ErrorPanic(err)
		query = query.Where("pjp_principles.permanent_journey_plans.pjp_code = ?", pjpCode)
	}

	result := query.Find(&rows)
	helper.ErrorPanic(result.Error)
	return rows
}
