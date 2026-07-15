package live_monitoring

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (r *liveMonitoringRepository) buildPrincipalEmployeeScopeQuery(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) *gorm.DB {
	query := tx.WithContext(ctx).Table("pjp_principles.permanent_journey_plans pjp").
		Joins("JOIN pjp_principles.route_pop_permanent rpp ON rpp.pjp_id = pjp.id").
		Joins("JOIN mst.m_salesman ms2 ON pjp.salesman_id = ms2.emp_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms2.emp_id").
		Joins("JOIN smc.m_customer mc ON ms2.cust_id = mc.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms WHERE ms.cust_id IN ?)", custIDs).
		Where("DATE(rpp.date) = ?", date).
		Where("pjp.approval_status IN ?", statuses)

	if areaID > 0 {
		query = query.Where("md.area_id = ?", areaID)
	}
	if regionID > 0 {
		query = query.Where("md.region_id = ?", regionID)
	}
	if distributorID > 0 {
		query = query.Where("md.distributor_id = ?", distributorID)
	}
	if len(empIDs) > 0 {
		query = query.Where("me.emp_id IN ?", empIDs)
	}

	return query
}

// GetPrincipalEmployeeIDs retrieves scoped principal employee IDs for pagination.
func (r *liveMonitoringRepository) GetPrincipalEmployeeIDs(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) ([]int, error) {
	var employeeIDs []int

	query := r.buildPrincipalEmployeeScopeQuery(ctx, tx, custIDs, date, regionID, areaID, distributorID, empIDs, statuses).
		Distinct("me.emp_id").
		Order("me.emp_id")

	if err := query.Pluck("me.emp_id", &employeeIDs).Error; err != nil {
		return nil, err
	}

	return employeeIDs, nil
}

// GetPrincipalMonitoring retrieves principal monitoring data from pjp_principles schema
func (r *liveMonitoringRepository) GetPrincipalMonitoring(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
	limit, offset int,
) ([]model.LiveMonitoringPrincipalRow, error) {
	var results []model.LiveMonitoringPrincipalRow

	visitStartAt, visitEndAt, err := buildLiveMonitoringDayRange(date)
	if err != nil {
		return nil, err
	}

	query := tx.WithContext(ctx).Table("pjp_principles.permanent_journey_plans pjp").
		Select(`
			me.emp_id,
			me.emp_code,
			me.emp_name,
			md.distributor_id,
			md.area_id,
			md.region_id,
			pjp.id AS pjp_id,
			pjp.pjp_code,
			pjp.approval_status,
			r.route_code,
			r.route_name,
			d.id AS destination_id,
			d.destination_code,
			d.destination_type,
			d.destination_name,
			COALESCE(d.destination_address, '') AS destination_address,
			COALESCE(CAST(d.longitude AS FLOAT), 0) AS longitude,
			COALESCE(CAST(d.latitude AS FLOAT), 0) AS latitude,
			ovl.arrive_at,
			ovl.leave_at,
			COALESCE(CAST(NULLIF(ovl.longitude, '') AS DOUBLE PRECISION), 0) AS arrive_longitude,
			COALESCE(CAST(NULLIF(ovl.latitude, '') AS DOUBLE PRECISION), 0) AS arrive_latitude,
			NULLIF(ovl.leave_longitude, '') AS leave_longitude,
			NULLIF(ovl.leave_latitude, '') AS leave_latitude,
			NULLIF((
				SELECT v.file_url FROM mobile.visits v
				WHERE v.cust_id = me.cust_id
				  AND v.emp_code = me.emp_code
				  AND v.outlet_code = d.destination_code
				  AND v.created_at >= ? AND v.created_at < ?
				ORDER BY v.created_at DESC, COALESCE(v.visit_id, 0) DESC
				LIMIT 1
			), '') AS file_url,
			ovl.skip_at,
			ovl.skip_reason,
			ovl."start",
			ovl.finish,
			FALSE AS is_extra_call
		`, visitStartAt, visitEndAt).
		Joins("JOIN pjp_principles.route_pop_permanent rpp ON rpp.pjp_id = pjp.id").
		Joins("JOIN pjp_principles.routes r ON r.route_code = rpp.route_code").
		Joins("JOIN pjp_principles.destinations d ON d.route_code = r.route_code").
		Joins("JOIN mst.m_salesman ms2 ON pjp.salesman_id = ms2.emp_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms2.emp_id").
		Joins("JOIN smc.m_customer mc ON ms2.cust_id = mc.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Joins("LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.pjp_id = pjp.id AND ovl.date = DATE(rpp.date) AND ovl.outlet_code = d.destination_code").
		Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms WHERE ms.cust_id IN ?)", custIDs).
		Where("DATE(rpp.date) = ?", date).
		Where("pjp.approval_status IN ?", statuses)

	if areaID > 0 {
		query = query.Where("md.area_id = ?", areaID)
	}
	if regionID > 0 {
		query = query.Where("md.region_id = ?", regionID)
	}
	if distributorID > 0 {
		query = query.Where("md.distributor_id = ?", distributorID)
	}
	if len(empIDs) > 0 {
		query = query.Where("me.emp_id IN ?", empIDs)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	query = query.Order("me.emp_id, pjp.id, r.route_code, d.id")

	err = query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CountPrincipalMonitoring counts unique employees for pagination
func (r *liveMonitoringRepository) CountPrincipalMonitoring(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) (int64, error) {
	var count int64

	query := r.buildPrincipalEmployeeScopeQuery(ctx, tx, custIDs, date, regionID, areaID, distributorID, empIDs, statuses).
		Distinct("me.emp_id")

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
