package repository

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryPjpDistributorImpl struct {
		*gorm.DB
	}
)

type PjpDistributorRepository interface {
	StorePermanentJourneyPlan(ctx context.Context, data *model.PermanentJourneyPlan) error
	UpdatePJPStatus(ctx context.Context, PJPID int64, status string) error

	// Route
	StoreRoute(ctx context.Context, data *model.Route) error
	UpdateRoute(ctx context.Context, route *model.Route) error
	BulkDeleteRoutes(ctx context.Context, pjpID []int64) error

	// RouteOutlet
	StoreRouteOutlet(ctx context.Context, data *model.RouteOutlet) error
	UpdateRouteOutlet(ctx context.Context, data *model.RouteOutlet) error
	FindAllRouteOutletsByParams(ctx context.Context, pjpID, routeCode int64, CustID string) (routes []model.RouteOutlet, err error)
	BulkDeleteRouteOutlets(ctx context.Context, routeOutletID []int64) error

	// RouteOutletHistory
	StoreRouteOutletHistory(ctx context.Context, data *model.RouteOutletHistory) error
	UpdateRouteOutletHistory(ctx context.Context, data *model.RouteOutletHistory) error
	BulkDeleteRouteOutletsHistories(ctx context.Context, routeOutletID []int64) error

	StoreRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanent) error
	FindAllRoutePopPermanentsByParams(ctx context.Context, pjpID int64, custID string) (routes []model.RoutePopPermanent, err error)
	UpdateRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanent) error
	BulkDeleteRoutePopPermanents(ctx context.Context, routePopPermanentID []int64) error

	GetPJPInfo(ctx context.Context, empID int64) (pjp model.PermanentJourneyPlan, err error)
	FindAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery) ([]model.PJPDistributorOutletList, int64, error)
	FindOnePJPByID(ctx context.Context, pjpID int64) (pjp model.PermanentJourneyPlan, err error)
	FindOnePJPByCode(ctx context.Context, pjpCode int64) (pjp model.PermanentJourneyPlan, err error)
	FindOnePJPByCodeAndCustId(ctx context.Context, pjpCode int64, custId string) (pjp model.PermanentJourneyPlan, err error)
	UpdatePJP(ctx context.Context, pjp *model.PermanentJourneyPlan) error
	FindAllRoutesByPJPID(ctx context.Context, pjpID int64) (routes []model.Route, err error)
	FindPJPByEmpIDAndCustID(ctx context.Context, empID int64, custID string) (pjp entity.RoutePopDailyResult, err error)

	FindAllRouteOutletHistoriesByPJPID(ctx context.Context, pjpID, routeCode int64, custID, date string) (routes []model.RouteOutletHistory, err error)
	CountOutletsByDate(ctx context.Context, salesmanID int64, dates string) (int, error)
	GetRouteNameByIndexDay(ctx context.Context, custID string, salesmanID int64, sequence int64) (int64, string, error)
	GetLastRouteCode(ctx context.Context) (int64, error)
	GetSalesmanAndTeam(ctx context.Context, empID int64) (entity.SalesmanAndTeamResult, error)
	GetWarehouseName(ctx context.Context, warehouseId int64, custId string) (string, error)
	GetRouteByRouteCode(ctx context.Context, routeCode int) (model.Route, error)
	CheckRouteOutlet(ctx context.Context, pjpID int64, routeCode, week, year []int) (map[string]int, error)
	FindOnePJPByCodeAndSalesman(ctx context.Context, pjpCode, salesman_id int64, cust_id string) (pjp model.PermanentJourneyPlan, err error)
	FindOneBySalesmanAndCustID(ctx context.Context, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error)
	GetVisitOverview(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) (entity.VisitOverview, error)
	GetNotBuyReasons(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) ([]entity.NotBuyReasonItem, error)
	GetSkipReasons(ctx context.Context, pjpID int64, startDate string, endDate string) ([]entity.SkipReasonItem, error)
}

func NewPjpDistributorRepository(db *gorm.DB) *RepositoryPjpDistributorImpl {
	return &RepositoryPjpDistributorImpl{db}
}

func (repository *RepositoryPjpDistributorImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repository.WithContext(ctx)
}

func (repository *RepositoryPjpDistributorImpl) StorePermanentJourneyPlan(ctx context.Context, data *model.PermanentJourneyPlan) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpDistributorImpl) UpdatePJPStatus(ctx context.Context, PJPID int64, status string) error {
	return repository.model(ctx).Table("pjp.permanent_journey_plans").Where("id = ?", PJPID).Update("approval_status", status).Error
}

func (repository *RepositoryPjpDistributorImpl) StoreRoute(ctx context.Context, data *model.Route) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpDistributorImpl) StoreRouteOutlet(ctx context.Context, data *model.RouteOutlet) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpDistributorImpl) StoreRouteOutletHistory(ctx context.Context, data *model.RouteOutletHistory) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpDistributorImpl) StoreRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanent) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpDistributorImpl) FindAllRoutesByPJPID(ctx context.Context, pjpID int64) (routes []model.Route, err error) {
	err = repository.WithContext(ctx).
		Select([]string{
			"routes.id",
			"routes.route_code",
			"routes.route_name",
			"routes.is_assign",
			"routes.cust_id",
			"routes.created_at",
			"routes.updated_at",
			"routes.pjp_id",
			"routes.is_pjp_old",
		}).
		Table("pjp.routes").
		Where("pjp.routes.pjp_id = ?", pjpID).
		Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (repository *RepositoryPjpDistributorImpl) FindAllRouteOutletsByParams(ctx context.Context, pjpID, routeCode int64, CustID string) (routes []model.RouteOutlet, err error) {
	err = repository.WithContext(ctx).
		Select([]string{
			"route_outlet.id",
			"route_outlet.route_code",
			"route_outlet.route_name",
			"route_outlet.outlet_id",
			"route_outlet.outlet_code",
			"route_outlet.outlet_name",
			"route_outlet.longitude",
			"route_outlet.latitude",
			"route_outlet.outlet_status",
			"route_outlet.outlet_address",
			"route_outlet.pjp_id",
			"route_outlet.pjp_code",
			"route_outlet.cust_id",
			"route_outlet.status",
			"route_outlet.created_at",
			"route_outlet.updated_at",
			"route_outlet.verified_date",
			"route_outlet.old_pjp_id",
			"route_outlet.old_pjp_code",
			"route_outlet.photo",
			"route_outlet.signature",
			"route_outlet.avg_sales_week",
		}).
		Table("pjp.route_outlet").
		Joins("INNER JOIN pjp.route_outlet_history ON pjp.route_outlet_history.pjp_id = pjp.route_outlet.pjp_id AND pjp.route_outlet_history.route_code = pjp.route_outlet.route_code AND pjp.route_outlet_history.cust_id = pjp.route_outlet.cust_id").
		Where("pjp.route_outlet.pjp_id= ?", pjpID).
		Where("pjp.route_outlet.route_code= ?", routeCode).
		Where("pjp.route_outlet.cust_id= ?", CustID).
		Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (repository *RepositoryPjpDistributorImpl) FindAllRouteOutletHistoriesByPJPID(ctx context.Context, pjpID, routeCode int64, custID, date string) (routes []model.RouteOutletHistory, err error) {
	err = repository.WithContext(ctx).
		Select([]string{
			"route_outlet_history.id",
			"route_outlet_history.route_code",
			"route_outlet_history.route_name",
			"route_outlet_history.outlet_id",
			"route_outlet_history.outlet_code",
			"route_outlet_history.outlet_name",
			"route_outlet_history.longitude",
			"route_outlet_history.latitude",
			"route_outlet_history.outlet_status",
			"route_outlet_history.outlet_address",
			"route_outlet_history.pjp_id",
			"route_outlet_history.pjp_code",
			"route_outlet_history.cust_id",
			"route_outlet_history.status",
			"route_outlet_history.created_at",
			"route_outlet_history.updated_at",
			"route_outlet_history.verified_date",
			"route_outlet_history.old_pjp_id",
			"route_outlet_history.old_pjp_code",
			"route_outlet_history.photo",
			"route_outlet_history.signature",
			"route_outlet_history.avg_sales_week",
			"route_outlet_history.index_day",
			"route_outlet_history.start_week",
			"route_outlet_history.is_in_current_year",
			"route_outlet_history.week",
			"route_outlet_history.year",
			"route_outlet_history.date",
			"route_outlet_history.is_additional",
		}).
		Table("pjp.route_outlet_history").
		Where("pjp.route_outlet_history.pjp_id = ?", pjpID).
		Where("pjp.route_outlet_history.route_code= ?", routeCode).
		Where("pjp.route_outlet_history.cust_id= ?", custID).
		Where("pjp.route_outlet_history.date= ?", date).
		Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (repository *RepositoryPjpDistributorImpl) FindAllRoutePopPermanentsByParams(ctx context.Context, pjpID int64, custID string) (routes []model.RoutePopPermanent, err error) {
	err = repository.WithContext(ctx).
		Select([]string{
			"route_pop_permanent.id",
			"route_pop_permanent.year",
			"route_pop_permanent.week",
			"route_pop_permanent.date",
			"route_pop_permanent.day",
			"route_pop_permanent.route_code",
			"route_pop_permanent.pjp_id",
			"route_pop_permanent.pjp_code",
			"route_pop_permanent.cust_id",
			"route_pop_permanent.created_at",
			"route_pop_permanent.updated_at",
		}).
		Table("pjp.route_pop_permanent").
		Where("pjp.route_pop_permanent.pjp_id= ?", pjpID).
		Where("pjp.route_pop_permanent.cust_id= ?", custID).
		Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (repository *RepositoryPjpDistributorImpl) GetPJPInfo(ctx context.Context, empID int64) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp.permanent_journey_plans.*").
		Table("pjp.permanent_journey_plans").
		Where("pjp.permanent_journey_plans.salesman_id = ? ", empID).
		Order("id DESC").
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpDistributorImpl) FindAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery) ([]model.PJPDistributorOutletList, int64, error) {
	var results []model.PJPDistributorOutletList
	var totalRecords int64

	// Initialize the base query
	db := repository.WithContext(ctx).Table("pjp.permanent_journey_plans AS pj")
	dbCount := repository.WithContext(ctx).Table("pjp.permanent_journey_plans AS pj")

	// 2. Set defaults and apply pagination/sorting
	if f.Limit <= 0 || f.Limit > 999 {
		f.Limit = 10
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	// Validate sorting column to prevent SQL injection
	allowedSortColumns := map[string]bool{
		"route_code":       true,
		"route_name":       true,
		"date":             true,
		"outlet_name":      true,
		"destination_name": true}
	if !allowedSortColumns[f.Sort] {
		f.Sort = "date"
	}

	orderClause := fmt.Sprintf("%s %s", f.Sort, f.SortOrder)

	err := db.Select(`
			DISTINCT
			COALESCE(
				roh.route_code,
				(SELECT (array_agg(r.route_code))[EXTRACT(ISODOW FROM ?::date)]
					FROM pjp.routes r where r.pjp_id = pj.id)
			) as route_code,
			COALESCE(
					roh.route_name,
					(SELECT (array_agg(r.route_name))[EXTRACT(ISODOW FROM ?::date)]
						FROM pjp.routes r where r.pjp_id = pj.id)
			) as route_name,
			roh.year,
			roh.week,
			COALESCE(roh.date, ?) as date,
			mo.outlet_id,
			mo.outlet_code,
			mo.outlet_name,
			mo.outlet_status,
			mo.address1 AS outlet_address,
			mo.longitude,
			mo.latitude
		`, f.Date, f.Date, f.Date).
		Joins(`LEFT JOIN pjp.route_outlet_history roh ON pj.id = roh.pjp_id AND roh.date = ?`, f.Date).
		Joins("LEFT JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id AND mo.is_del = false AND mo.verification_status = 1 AND mo.outlet_status != 4").
		Where("pj.salesman_id = ?", f.EmpID).
		Order(orderClause).
		Limit(f.Limit).
		Offset(offset).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	errCount := dbCount.Select("COUNT(roh.outlet_id)").
		Joins(`LEFT JOIN pjp.route_outlet_history roh ON pj.id = roh.pjp_id AND roh.date = ?`, f.Date).
		Joins("LEFT JOIN mst.m_outlet mo ON mo.outlet_id = roh.outlet_id AND mo.is_del = false AND mo.verification_status = 1 AND mo.outlet_status != 4").
		Where("pj.salesman_id = ?", f.EmpID).Scan(&totalRecords).Error
	if errCount != nil {
		return nil, 0, errCount
	}
	return results, totalRecords, err
}

func (repository *RepositoryPjpDistributorImpl) FindOnePJPByID(ctx context.Context, pjpID int64) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp.permanent_journey_plans.*").
		Table("pjp.permanent_journey_plans").
		Where("pjp.permanent_journey_plans.id = ? ", pjpID).
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpDistributorImpl) FindOnePJPByCode(ctx context.Context, pjpCode int64) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp.permanent_journey_plans.*").
		Table("pjp.permanent_journey_plans").
		Where("pjp.permanent_journey_plans.pjp_code = ? ", pjpCode).
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpDistributorImpl) FindOnePJPByCodeAndCustId(ctx context.Context, pjpCode int64, custId string) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp.permanent_journey_plans.*").
		Table("pjp.permanent_journey_plans").
		Where("pjp.permanent_journey_plans.pjp_code = ? AND pjp.permanent_journey_plans.cust_id = ?", pjpCode, custId).
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpDistributorImpl) UpdatePJP(ctx context.Context, pjp *model.PermanentJourneyPlan) error {
	return repository.WithContext(ctx).Save(pjp).Where("id = ?", pjp.ID).Error
}

func (repository *RepositoryPjpDistributorImpl) BulkDeleteRoutes(ctx context.Context, pjpID []int64) error {
	return repository.WithContext(ctx).Where("id in (?)", pjpID).Delete(&model.Route{}).Error
}

func (repository *RepositoryPjpDistributorImpl) UpdateRoute(ctx context.Context, route *model.Route) error {
	return repository.WithContext(ctx).Save(route).Where("id = ?", route.ID).Error
}

func (repository *RepositoryPjpDistributorImpl) UpdateRouteOutlet(ctx context.Context, data *model.RouteOutlet) error {
	return repository.WithContext(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpDistributorImpl) BulkDeleteRouteOutlets(ctx context.Context, routeOutletID []int64) error {
	return repository.WithContext(ctx).Where("id in (?)", routeOutletID).Delete(&model.RouteOutlet{}).Error
}

func (repository *RepositoryPjpDistributorImpl) BulkDeleteRouteOutletsHistories(ctx context.Context, routeOutletID []int64) error {
	return repository.WithContext(ctx).Where("id in (?)", routeOutletID).Delete(&model.RouteOutletHistory{}).Error
}

func (repository *RepositoryPjpDistributorImpl) UpdateRouteOutletHistory(ctx context.Context, data *model.RouteOutletHistory) error {
	return repository.WithContext(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpDistributorImpl) UpdateRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanent) error {
	return repository.WithContext(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpDistributorImpl) BulkDeleteRoutePopPermanents(ctx context.Context, routePopPermanentID []int64) error {
	return repository.WithContext(ctx).Where("id in (?)", routePopPermanentID).Delete(&model.RoutePopPermanent{}).Error
}

// CountOutletsByDate counts the number of outlets per date for a given salesman
func (repository *RepositoryPjpDistributorImpl) CountOutletsByDate(ctx context.Context, salesmanID int64, dates string) (int, error) {
	type DateCount struct {
		Date  string
		Count int
	}

	var results []DateCount

	err := repository.WithContext(ctx).
		Table("pjp.route_outlet_history AS ovl").
		Select("ovl.date::text as date, COUNT(ovl.outlet_id)::int as count").
		Joins("JOIN pjp.permanent_journey_plans AS pjp ON ovl.pjp_id = pjp.id").
		Joins("JOIN mst.m_outlet AS mo ON ovl.outlet_id = mo.outlet_id").
		Where("pjp.salesman_id = ?", salesmanID).
		Where("ovl.date = ?", dates).
		Group("ovl.date").
		Scan(&results).Error

	if err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	return results[0].Count, nil
}

func (repository *RepositoryPjpDistributorImpl) GetRouteNameByIndexDay(ctx context.Context, custID string, salesmanID int64, sequence int64) (int64, string, error) {
	var result struct {
		RouteCode int64
		RouteName string
	}
	err := repository.WithContext(ctx).
		Table("pjp.routes r").
		Select("r.route_code, r.route_name").
		Joins("INNER JOIN pjp.permanent_journey_plans pj ON pj.id = r.pjp_id").
		Where("pj.salesman_id = ?", salesmanID).
		Where("r.cust_id= ? AND r.sequence= ?", custID, sequence).
		Take(&result).Error
	if err != nil {
		return 0, "", err
	}
	return result.RouteCode, result.RouteName, nil
}

func (repository *RepositoryPjpDistributorImpl) GetLastRouteCode(ctx context.Context) (int64, error) {
	var routeCode int64
	err := repository.WithContext(ctx).
		Table("pjp.routes AS r").
		Select("r.route_code").
		Order("r.route_code DESC").
		Limit(1).
		Scan(&routeCode).Error
	return routeCode, err
}

func (repository *RepositoryPjpDistributorImpl) GetRouteByRouteCode(ctx context.Context, routeCode int) (model.Route, error) {
	var route model.Route
	err := repository.WithContext(ctx).
		Table("pjp.routes AS r").
		Where("r.route_code = ?", routeCode).
		Take(&route).Error
	return route, err
}

func (repository *RepositoryPjpDistributorImpl) GetSalesmanAndTeam(ctx context.Context, empID int64) (entity.SalesmanAndTeamResult, error) {
	var result entity.SalesmanAndTeamResult
	err := repository.WithContext(ctx).
		Table("mst.m_employee me").
		Select("me.emp_id, me.emp_code, ms.sales_name, st.sales_team_name").
		Joins("inner join mst.m_salesman ms on ms.emp_id = me.emp_id").
		Joins("inner join mst.m_sales_team st on st.sales_team_id = ms.sales_team_id").
		Where("me.emp_id = ?", empID).
		Scan(&result).Error
	return result, err
}

func (repository *RepositoryPjpDistributorImpl) GetWarehouseName(ctx context.Context, warehouseId int64, custId string) (string, error) {
	var warehouseName string
	err := repository.model(ctx).
		Table("mst.m_warehouse").
		Select("wh_name").
		Where("wh_id = ? AND cust_id = ?", warehouseId, custId).
		Take(&warehouseName).Error
	return warehouseName, err
}

func (repository *RepositoryPjpDistributorImpl) FindPJPByEmpIDAndCustID(ctx context.Context, empID int64, custID string) (pjp entity.RoutePopDailyResult, err error) {
	err = repository.WithContext(ctx).
		Raw(`
			WITH getPjpIds AS (
				SELECT id
				FROM pjp.permanent_journey_plans
				WHERE salesman_id = ? AND cust_id = ?
			)
			select
				route_code,
				pjp_id,
				pjp_code
			from pjp.route_pop_daily
			where day = TO_CHAR(NOW(), 'Dy') and pjp_id = (SELECT id FROM getPjpIds)
			order by created_at desc limit 1;
		`, empID, custID).
		Scan(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpDistributorImpl) CheckRouteOutlet(ctx context.Context, pjpID int64, routeCode, week, year []int) (map[string]int, error) {
	var result []model.RouteOutletInUse
	err := repository.WithContext(ctx).
		Table("pjp.route_outlet_history AS roh").
		Select("COUNT(id) as total, route_code").
		Where("roh.pjp_id = ?", pjpID).
		Where("roh.route_code IN ?", routeCode).
		Where("roh.week NOT IN ?", week).
		Where("roh.year IN ?", year).
		Group("route_code").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		resultMap := make(map[string]int)
		for _, r := range result {
			resultMap[r.RouteCode] = r.Total
		}
		return resultMap, nil
	}

	return nil, nil
}

func (repository *RepositoryPjpDistributorImpl) FindOnePJPByCodeAndSalesman(ctx context.Context, pjpCode, salesman_id int64, cust_id string) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp.permanent_journey_plans.*").
		Table("pjp.permanent_journey_plans").
		Where("pjp.permanent_journey_plans.pjp_code = ? AND salesman_id = ? AND cust_id = ?", pjpCode, salesman_id, cust_id).
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpDistributorImpl) FindOneBySalesmanAndCustID(ctx context.Context, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp.permanent_journey_plans").
		Where("salesman_id = ? AND cust_id = ?", salesman_id, cust_id).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpDistributorImpl) GetVisitOverview(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) (entity.VisitOverview, error) {
	var result entity.VisitOverview

	q := `
	SELECT
		COUNT(o.outlet_id) AS total_buy,
		COUNT(t.outlet_id) AS total_not_buy,
		COUNT(ovl.id) AS total_visit_list,
		COUNT(CASE WHEN ovl.arrive_at IS NOT NULL THEN 1 END) AS total_visit,
		COUNT(CASE WHEN ovl.arrive_at IS NULL THEN 1 END) AS total_not_visit,
		COUNT(CASE WHEN ovl.arrive_at IS NOT NULL AND is_planned = TRUE THEN 1 END) AS visit_planned,
		COUNT(CASE WHEN ovl.arrive_at IS NOT NULL AND is_planned = FALSE THEN 1 END) AS visit_not_planned,
		COUNT(CASE WHEN ovl.arrive_at IS NULL AND is_planned = TRUE THEN 1 END) AS not_visit_planned,
		COUNT(CASE WHEN ovl.arrive_at IS NULL AND is_planned = FALSE THEN 1 END) AS not_visit_not_planned
	FROM pjp.outlet_visit_list ovl
	LEFT JOIN sls.no_order t
		ON t.salesman_id = ?
		AND t.outlet_id = ovl.outlet_id
		AND t.created_at::date = ovl.date
	LEFT JOIN sls.order o
		ON o.salesman_id = ?
		AND o.outlet_id = ovl.outlet_id
		AND o.ro_date = ovl.date
	WHERE ovl.pjp_id = ?
		AND ovl.date BETWEEN ? AND ?`

	err := repository.WithContext(ctx).
		Raw(q, salesmanID, salesmanID, pjpID, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return entity.VisitOverview{}, err
	}

	return result, nil
}

func (repository *RepositoryPjpDistributorImpl) GetNotBuyReasons(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) ([]entity.NotBuyReasonItem, error) {
	var results []entity.NotBuyReasonItem

	q := `
	SELECT
		t.reason,
		COUNT(t.no_order_id) AS total
	FROM sls.no_order t
	JOIN pjp.outlet_visit_list ovl
		ON t.outlet_id = ovl.outlet_id
		AND ovl.pjp_id = ?
		AND ovl.date BETWEEN ? AND ?
	WHERE t.salesman_id = ?
		AND t.created_at::date = ovl.date
	GROUP BY t.reason`

	err := repository.WithContext(ctx).
		Raw(q, pjpID, startDate, endDate, salesmanID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repository *RepositoryPjpDistributorImpl) GetSkipReasons(ctx context.Context, pjpID int64, startDate string, endDate string) ([]entity.SkipReasonItem, error) {
	var results []entity.SkipReasonItem

	q := `
	SELECT
		CASE WHEN ovl.skip_at IS NOT NULL
			THEN ovl.skip_reason
			ELSE 'Unvisited'
		END AS skip_reason,
		COUNT(ovl.id) AS total
	FROM pjp.outlet_visit_list ovl
	WHERE ovl.pjp_id = ?
		AND ovl.date BETWEEN ? AND ?
		AND ovl.arrive_at IS NULL
	GROUP BY
		CASE WHEN ovl.skip_at IS NOT NULL
			THEN ovl.skip_reason
			ELSE 'Unvisited'
		END`

	err := repository.WithContext(ctx).
		Raw(q, pjpID, startDate, endDate).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
