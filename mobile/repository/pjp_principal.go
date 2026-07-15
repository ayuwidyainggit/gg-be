package repository

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryPjpPrincipalImpl struct {
		*gorm.DB
	}
)

type PjpPrincipalRepository interface {
	UpdatePJPStatus(ctx context.Context, PJPID int64, status string) error
	StorePermanentJourneyPlan(ctx context.Context, data *model.PermanentJourneyPlanPrincipal) error
	StoreRoute(ctx context.Context, data *model.RoutePrincipal) error
	StoreDestination(ctx context.Context, data *model.DestinationPrincipal) error
	StoreDestinationHistory(ctx context.Context, data *model.DestinationHistoryPrincipal) error
	StoreRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanentPrincipal) error
	GetWarehouseName(ctx context.Context, warehouseId int64, custId string) (string, error)
	GetPJPInfo(ctx context.Context, empID int64) (pjp model.PermanentJourneyPlan, err error)
	FindAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery, destinationType string) ([]model.PJPPrincipleDestinationList, error)
	CountOutletsByDate(ctx context.Context, salesmanID int64, date string) (model.DestinationCount, error)
	FindOnePJPByID(ctx context.Context, id int64) (model.PermanentJourneyPlanPrincipal, error)
	FindOnePJPByCode(ctx context.Context, code int64) (model.PermanentJourneyPlanPrincipal, error)
	FindOnePJPByCodeAndCustId(ctx context.Context, code int64, custId string) (model.PermanentJourneyPlanPrincipal, error)
	UpdatePJP(ctx context.Context, data *model.PermanentJourneyPlanPrincipal) error
	FindAllRoutesByPJPID(ctx context.Context, pjpID int64) ([]model.RoutePrincipal, error)
	UpdateRoute(ctx context.Context, data *model.RoutePrincipal) error
	BulkDeleteRoutes(ctx context.Context, ids []int64) error
	FindAllDestinationsByParams(ctx context.Context, pjpID, routeCode int64, custID string) ([]model.DestinationPrincipal, error)
	UpdateDestination(ctx context.Context, data *model.DestinationPrincipal) error
	BulkDeleteDestinations(ctx context.Context, ids []int64) error
	FindAllDestinationHistoriesByPJPID(ctx context.Context, pjpID, routeCode int64, custID, date string) ([]model.DestinationHistoryPrincipal, error)
	UpdateDestinationHistory(ctx context.Context, data *model.DestinationHistoryPrincipal) error
	BulkDeleteDestinationHistories(ctx context.Context, ids []int64) error
	FindAllRoutePopPermanentsByParams(ctx context.Context, pjpID int64, custID string) ([]model.RoutePopPermanentPrincipal, error)
	UpdateRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanentPrincipal) error
	BulkDeleteRoutePopPermanents(ctx context.Context, ids []int64) error
	GetRouteNameByIndexDay(ctx context.Context, custID string, salesmanID int64, sequence int64) (int64, string, error)
	GetLastRouteCode(ctx context.Context) (int64, error)
	GetSalesmanAndTeam(ctx context.Context, empID int64) (entity.SalesmanAndTeamResult, error)
	GetRouteByRouteCode(ctx context.Context, routeCode int) (model.Route, error)
	FindPJPByEmpIDAndCustID(ctx context.Context, empID int64, custID string) (entity.RoutePopDailyPrincipalResult, error)
	CheckRouteOutlet(ctx context.Context, pjpID int64, routeCode, week, year []int) (map[string]int, error)
	FindOnePJPByCodeAndSalesmanID(ctx context.Context, code, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error)
	CountAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery) (int64, error)
	FindOneBySalesmanAndCustID(ctx context.Context, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error)
	GetVisitOverview(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) (entity.VisitOverview, error)
	GetNotBuyReasons(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) ([]entity.NotBuyReasonItem, error)
	GetSkipReasons(ctx context.Context, pjpID int64, startDate string, endDate string) ([]entity.SkipReasonItem, error)
}

func NewPjpPrincipalRepository(db *gorm.DB) *RepositoryPjpPrincipalImpl {
	return &RepositoryPjpPrincipalImpl{db}
}

func (repo *RepositoryPjpPrincipalImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryPjpPrincipalImpl) StorePermanentJourneyPlan(ctx context.Context, data *model.PermanentJourneyPlanPrincipal) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpPrincipalImpl) UpdatePJPStatus(ctx context.Context, PJPID int64, status string) error {
	return repository.model(ctx).Table("pjp_principles.permanent_journey_plans").Where("id = ?", PJPID).Update("approval_status", status).Error
}

func (repository *RepositoryPjpPrincipalImpl) StoreRoute(ctx context.Context, data *model.RoutePrincipal) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpPrincipalImpl) StoreDestination(ctx context.Context, data *model.DestinationPrincipal) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpPrincipalImpl) StoreDestinationHistory(ctx context.Context, data *model.DestinationHistoryPrincipal) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpPrincipalImpl) StoreRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanentPrincipal) error {
	return repository.model(ctx).Create(data).Error
}

func (repository *RepositoryPjpPrincipalImpl) GetWarehouseName(ctx context.Context, warehouseId int64, custId string) (string, error) {
	var warehouseName string
	err := repository.model(ctx).
		Table("mst.m_warehouse").
		Select("wh_name").
		Where("wh_id = ? AND cust_id = ?", warehouseId, custId).
		Take(&warehouseName).Error
	return warehouseName, err
}

func (repository *RepositoryPjpPrincipalImpl) GetPJPInfo(ctx context.Context, empID int64) (pjp model.PermanentJourneyPlan, err error) {
	err = repository.WithContext(ctx).
		Select("pjp_principles.permanent_journey_plans.*").
		Table("pjp_principles.permanent_journey_plans").
		Where("pjp_principles.permanent_journey_plans.salesman_id = ? ", empID).
		Order("id DESC").
		Take(&pjp).Error
	if err != nil {
		return model.PermanentJourneyPlan{}, err
	}

	return pjp, nil
}

func (repository *RepositoryPjpPrincipalImpl) FindAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery, destinationType string) ([]model.PJPPrincipleDestinationList, error) {
	var results []model.PJPPrincipleDestinationList

	// Initialize the base query
	db := repository.WithContext(ctx).Table("pjp_principles.permanent_journey_plans AS pj")

	// 2. Set defaults and apply pagination/sorting
	if f.Limit <= 0 || f.Limit > 999 {
		f.Limit = 10
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	// Validate sorting column to prevent SQL injection
	allowedSortColumns := map[string]bool{"route_code": true, "date": true, "outlet_name": true}
	if !allowedSortColumns[f.Sort] {
		f.Sort = "date"
	}

	orderClause := fmt.Sprintf("%s %s", f.Sort, f.SortOrder)

	switch destinationType {
	case "outlet":
		db.Select(`
			DISTINCT
			COALESCE(
				dh.route_code,
				(SELECT (array_agg(r.route_code))[EXTRACT(ISODOW FROM ?::date)]
					FROM pjp_principles.routes r where r.pjp_id = pj.id)
			) as route_code,
			COALESCE(
				dh.route_name,
				(SELECT (array_agg(r.route_name))[EXTRACT(ISODOW FROM ?::date)]
					FROM pjp_principles.routes r where r.pjp_id = pj.id)
	    	) as route_name,
		  dh.year,
		  dh.week,
		  COALESCE(dh.date, ?) as date,
		  mo.outlet_id AS destination_id,
		  mo.outlet_name AS destination_name,
		  mo.outlet_code AS destination_code,
		  mo.outlet_status AS destination_status,
		  mo.address1 AS destination_address,
		  mo.longitude,
		  mo.latitude
		`, f.Date, f.Date, f.Date).
			Joins("LEFT JOIN pjp_principles.destinations_history AS dh ON pj.id = dh.pjp_id AND dh.destination_type = ? AND dh.date = ?", destinationType, f.Date).
			Joins("LEFT JOIN mst.m_outlet AS mo ON mo.outlet_id = dh.destination_id AND mo.is_del = false AND mo.verification_status = 1 AND mo.outlet_status != 4")

	case "distributor":
		db.Select(`
			DISTINCT
			COALESCE(
					dh.route_code,
					(SELECT (array_agg(r.route_code))[EXTRACT(ISODOW FROM ?::date)]
						FROM pjp_principles.routes r where r.pjp_id = pj.id)
					) as route_code,
			COALESCE(
					dh.route_name,
					(SELECT (array_agg(r.route_name))[EXTRACT(ISODOW FROM ?::date)]
						FROM pjp_principles.routes r where r.pjp_id = pj.id)
					) as route_name,
					dh.year,
					dh.week,
					COALESCE(dh.date, ?) as date,
					md.distributor_id AS destination_id,
					md.distributor_name AS destination_name,
					md.distributor_code AS destination_code,
					md.address AS destination_address,
					md.longitude,
					md.latitude
		`, f.Date, f.Date, f.Date).
			Joins("LEFT JOIN pjp_principles.destinations_history AS dh ON pj.id = dh.pjp_id AND dh.destination_type = ? AND dh.date = ?", destinationType, f.Date).
			Joins("LEFT JOIN mst.m_distributor AS md ON md.distributor_id = dh.destination_id AND md.distributor_id IS NOT NULL AND md.is_del = false AND md.is_active = true")
	default:
	}
	err := db.Where("pj.salesman_id = ?", f.EmpID).
		Order(orderClause).
		Limit(f.Limit).
		Offset(offset).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repository *RepositoryPjpPrincipalImpl) CountAllOutletByDate(ctx context.Context, f entity.OutletPJPListQuery) (int64, error) {
	var totalRecords int64

	errCount := repository.WithContext(ctx).Table("pjp_principles.permanent_journey_plans AS pj").
		Select("(COUNT(mo.outlet_id) + COUNT(md.distributor_id))").
		Joins("LEFT JOIN pjp_principles.destinations_history AS dh ON pj.id = dh.pjp_id AND dh.date = ?", f.Date).
		Joins("LEFT JOIN mst.m_outlet AS mo ON mo.outlet_id = dh.destination_id AND mo.is_del = false AND mo.verification_status = 1 AND mo.outlet_status != 4").
		Joins("LEFT JOIN mst.m_distributor AS md ON md.distributor_id = dh.destination_id AND md.distributor_id IS NOT NULL and md.is_del = false AND md.is_active = true").
		Where("pj.salesman_id = ?", f.EmpID).Scan(&totalRecords).Error

	if errCount != nil {
		return 0, errCount
	}
	return totalRecords, nil
}

// CountOutletsByDate counts the number of outlets and distributors per date for a given salesman (principal user)
func (repository *RepositoryPjpPrincipalImpl) CountOutletsByDate(ctx context.Context, salesmanID int64, date string) (model.DestinationCount, error) {
	var result model.DestinationCount

	err := repository.WithContext(ctx).
		Table("pjp_principles.destinations_history AS dh").
		Select(`
			COUNT(CASE 
				WHEN dh.destination_type = 'outlet' 
					AND mo.outlet_id IS NOT NULL 
				THEN 1 
			END) AS total_outlet,
			COUNT(CASE 
				WHEN dh.destination_type = 'distributor' 
					AND md.distributor_id IS NOT NULL 
				THEN 1 
			END) AS total_distributor
		`).
		Joins("JOIN pjp_principles.permanent_journey_plans AS pjp ON pjp.id = dh.pjp_id").
		Joins("LEFT JOIN mst.m_outlet AS mo ON mo.outlet_id = dh.destination_id AND dh.destination_type = 'outlet'").
		Joins("LEFT JOIN mst.m_distributor AS md ON md.distributor_id = dh.destination_id AND dh.destination_type = 'distributor'").
		Where("dh.date = ?", date).
		Where("pjp.salesman_id = ?", salesmanID).
		Where(`(
			(dh.destination_type = 'outlet' AND mo.outlet_id IS NOT NULL)
			OR 
			(dh.destination_type = 'distributor' AND md.distributor_id IS NOT NULL)
		)`).
		Scan(&result).Error

	if err != nil {
		return model.DestinationCount{}, err
	}

	return result, nil
}

func (repository *RepositoryPjpPrincipalImpl) FindOnePJPByID(ctx context.Context, id int64) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.permanent_journey_plans").
		Where("id = ?", id).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) FindOnePJPByCode(ctx context.Context, code int64) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.permanent_journey_plans").
		Where("pjp_code = ?", code).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) FindOnePJPByCodeAndCustId(ctx context.Context, code int64, custId string) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.permanent_journey_plans").
		Where("pjp_code = ? AND cust_id = ?", code, custId).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) UpdatePJP(ctx context.Context, data *model.PermanentJourneyPlanPrincipal) error {
	return repository.model(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpPrincipalImpl) FindAllRoutesByPJPID(ctx context.Context, pjpID int64) ([]model.RoutePrincipal, error) {
	var routes []model.RoutePrincipal
	err := repository.model(ctx).
		Table("pjp_principles.routes").
		Where("pjp_id = ?", pjpID).
		Find(&routes).Error
	return routes, err
}

func (repository *RepositoryPjpPrincipalImpl) UpdateRoute(ctx context.Context, data *model.RoutePrincipal) error {
	return repository.model(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpPrincipalImpl) BulkDeleteRoutes(ctx context.Context, ids []int64) error {
	return repository.model(ctx).Where("id IN (?)", ids).Delete(&model.RoutePrincipal{}).Error
}

func (repository *RepositoryPjpPrincipalImpl) FindAllDestinationsByParams(ctx context.Context, pjpID, routeCode int64, custID string) ([]model.DestinationPrincipal, error) {
	var destinations []model.DestinationPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.destinations").
		Joins("LEFT JOIN pjp_principles.destinations_history dh on dh.id = destinations.id").
		Where("pjp_principles.destinations.pjp_id = ? AND pjp_principles.destinations.route_code = ? AND pjp_principles.destinations.cust_id = ?", pjpID, routeCode, custID).
		Find(&destinations).Error
	return destinations, err
}

func (repository *RepositoryPjpPrincipalImpl) UpdateDestination(ctx context.Context, data *model.DestinationPrincipal) error {
	return repository.model(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpPrincipalImpl) BulkDeleteDestinations(ctx context.Context, ids []int64) error {
	return repository.model(ctx).Where("id IN (?)", ids).Delete(&model.DestinationPrincipal{}).Error
}

func (repository *RepositoryPjpPrincipalImpl) FindAllDestinationHistoriesByPJPID(ctx context.Context, pjpID, routeCode int64, custID, date string) ([]model.DestinationHistoryPrincipal, error) {
	var histories []model.DestinationHistoryPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.destinations_history").
		Where("pjp_id = ? AND route_code = ? AND cust_id = ?", pjpID, routeCode, custID).
		Where("date = ?", date).
		Find(&histories).Error
	return histories, err
}

func (repository *RepositoryPjpPrincipalImpl) UpdateDestinationHistory(ctx context.Context, data *model.DestinationHistoryPrincipal) error {
	return repository.model(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpPrincipalImpl) BulkDeleteDestinationHistories(ctx context.Context, ids []int64) error {
	return repository.model(ctx).Where("id IN (?)", ids).Delete(&model.DestinationHistoryPrincipal{}).Error
}

func (repository *RepositoryPjpPrincipalImpl) FindAllRoutePopPermanentsByParams(ctx context.Context, pjpID int64, custID string) ([]model.RoutePopPermanentPrincipal, error) {
	var pops []model.RoutePopPermanentPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.route_pop_permanent").
		Where("pjp_id = ? AND cust_id = ?", pjpID, custID).
		Find(&pops).Error
	return pops, err
}

func (repository *RepositoryPjpPrincipalImpl) UpdateRoutePopPermanent(ctx context.Context, data *model.RoutePopPermanentPrincipal) error {
	return repository.model(ctx).Save(data).Where("id = ?", data.ID).Error
}

func (repository *RepositoryPjpPrincipalImpl) BulkDeleteRoutePopPermanents(ctx context.Context, ids []int64) error {
	return repository.model(ctx).Where("id IN (?)", ids).Delete(&model.RoutePopPermanentPrincipal{}).Error
}

func (repository *RepositoryPjpPrincipalImpl) GetRouteNameByIndexDay(ctx context.Context, custID string, salesmanID int64, sequence int64) (int64, string, error) {
	var result struct {
		RouteCode int64
		RouteName string
	}
	err := repository.WithContext(ctx).
		Table("pjp_principles.routes r").
		Select("r.route_code, r.route_name").
		Joins("INNER JOIN pjp_principles.permanent_journey_plans pj on pj.id = r.pjp_id").
		Where("pj.salesman_id = ?", salesmanID).
		Where("r.cust_id= ? AND r.sequence= ?", custID, sequence).
		Take(&result).Error
	if err != nil {
		return 0, "", err
	}
	return result.RouteCode, result.RouteName, nil
}

func (repository *RepositoryPjpPrincipalImpl) GetLastRouteCode(ctx context.Context) (int64, error) {
	var routeCode int64
	err := repository.WithContext(ctx).
		Table("pjp_principles.routes AS r").
		Select("r.route_code").
		Order("r.route_code DESC").
		Limit(1).
		Scan(&routeCode).Error
	return routeCode, err
}

func (repository *RepositoryPjpPrincipalImpl) GetRouteByRouteCode(ctx context.Context, routeCode int) (model.Route, error) {
	var route model.Route
	err := repository.WithContext(ctx).
		Table("pjp_principles.routes AS r").
		Where("r.route_code = ?", routeCode).
		Take(&route).Error
	return route, err
}

func (repository *RepositoryPjpPrincipalImpl) GetSalesmanAndTeam(ctx context.Context, empID int64) (entity.SalesmanAndTeamResult, error) {
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

func (repository *RepositoryPjpPrincipalImpl) FindPJPByEmpIDAndCustID(ctx context.Context, empID int64, custID string) (pjp entity.RoutePopDailyPrincipalResult, err error) {
	err = repository.WithContext(ctx).
		Raw(`
			WITH getPjpIds AS (
				SELECT id
				FROM pjp_principles.permanent_journey_plans
				WHERE salesman_id = ? AND cust_id = ?
			)
			select
				route_code,
				pjp_id,
				pjp_code
			from pjp_principles.route_pop_dailies
			where day = TO_CHAR(NOW(), 'Dy') and pjp_id = (SELECT id FROM getPjpIds)
			order by created_at desc limit 1;
		`, empID, custID).
		Scan(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) CheckRouteOutlet(ctx context.Context, pjpID int64, routeCode, week, year []int) (map[string]int, error) {
	var result []model.RouteOutletInUse
	err := repository.WithContext(ctx).
		Table("pjp_principles.destinations_history AS roh").
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

func (repository *RepositoryPjpPrincipalImpl) FindOnePJPByCodeAndSalesmanID(ctx context.Context, code, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.permanent_journey_plans").
		Where("pjp_code = ? AND salesman_id = ? AND cust_id = ?", code, salesman_id, cust_id).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) FindOneBySalesmanAndCustID(ctx context.Context, salesman_id int64, cust_id string) (model.PermanentJourneyPlanPrincipal, error) {
	var pjp model.PermanentJourneyPlanPrincipal
	err := repository.model(ctx).
		Table("pjp_principles.permanent_journey_plans").
		Where("salesman_id = ? AND cust_id = ?", salesman_id, cust_id).
		Take(&pjp).Error
	return pjp, err
}

func (repository *RepositoryPjpPrincipalImpl) GetVisitOverview(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) (entity.VisitOverview, error) {
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
	FROM pjp_principles.outlet_visit_list ovl
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

func (repository *RepositoryPjpPrincipalImpl) GetNotBuyReasons(ctx context.Context, pjpID int64, salesmanID int64, startDate string, endDate string) ([]entity.NotBuyReasonItem, error) {
	var results []entity.NotBuyReasonItem

	q := `
	SELECT
		t.reason,
		COUNT(t.no_order_id) AS total
	FROM sls.no_order t
	JOIN pjp_principles.outlet_visit_list ovl
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

func (repository *RepositoryPjpPrincipalImpl) GetSkipReasons(ctx context.Context, pjpID int64, startDate string, endDate string) ([]entity.SkipReasonItem, error) {
	var results []entity.SkipReasonItem

	q := `
	SELECT
		CASE WHEN ovl.skip_at IS NOT NULL
			THEN ovl.skip_reason
			ELSE 'Unvisited'
		END AS skip_reason,
		COUNT(ovl.id) AS total
	FROM pjp_principles.outlet_visit_list ovl
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
