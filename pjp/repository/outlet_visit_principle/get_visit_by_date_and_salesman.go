package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) GetVisitsByDateAndSalesman(ctx context.Context, tx *gorm.DB, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitListPrinciple, error) {
	var visits []model.OutletVisitListPrinciple

	query := tx.Table("pjp_principles.outlet_visit_list").
		Select("id, year, week, date, day, route_code, outlet_id, outlet_code, pjp_id, pjp_code, start, finish, skip_at, leave_at, arrive_at, on_hold, resume_at, skip_reason")

	subquery := tx.Table("pjp_principles.permanent_journey_plans").
		Select("id").
		Where("salesman_id = ?", dataFilter.SalesmanId)
	query = query.Where("pjp_principles.outlet_visit_list.pjp_id IN (?)", subquery)

	if dataFilter.Date != "" {
		query = query.Where("pjp_principles.outlet_visit_list.date = ?", dataFilter.Date)
	}
	if dataFilter.Year != "" {
		query = query.Where("EXTRACT(YEAR FROM pjp_principles.outlet_visit_list.date) = ?", dataFilter.Year)
	}
	if dataFilter.Month != "" {
		query = query.Where("EXTRACT(MONTH FROM pjp_principles.outlet_visit_list.date) = ?", dataFilter.Month)
	}

	err := query.Find(&visits).Error
	return visits, err
}
