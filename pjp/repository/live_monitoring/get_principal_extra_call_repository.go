package live_monitoring

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

// GetPrincipalExtraCalls retrieves principal extra call monitoring rows from pjp_principles.destinations_history.
// Join shape mirrors GetPrincipalMonitoring to stay Citus-friendly (no joins to mst.m_outlet/m_distributor).
// Destination fields are taken directly from destinations_history; outlet_visit_list join uses outlet_code
// to align with the regular query and to keep the planner happy when destination_id is legacy-NULL.
func (r *liveMonitoringRepository) GetPrincipalExtraCalls(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) ([]model.LiveMonitoringPrincipalRow, error) {
	var results []model.LiveMonitoringPrincipalRow

	visitStartAt, visitEndAt, err := buildLiveMonitoringDayRange(date)
	if err != nil {
		return nil, err
	}

	query := tx.WithContext(ctx).Table("pjp_principles.destinations_history dh").
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
			COALESCE(r.route_code, dh.route_code) AS route_code,
			COALESCE(r.route_name, dh.route_name) AS route_name,
			COALESCE(dh.destination_id, 0) AS destination_id,
			COALESCE(dh.destination_code, '') AS destination_code,
			COALESCE(NULLIF(dh.destination_type, ''), 'outlet') AS destination_type,
			COALESCE(dh.destination_name, '') AS destination_name,
			COALESCE(dh.destination_address, '') AS destination_address,
			COALESCE(CAST(NULLIF(dh.longitude, '') AS DOUBLE PRECISION), 0) AS longitude,
			COALESCE(CAST(NULLIF(dh.latitude, '') AS DOUBLE PRECISION), 0) AS latitude,
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
				  AND v.outlet_code = dh.destination_code
				  AND v.created_at >= ? AND v.created_at < ?
				ORDER BY v.created_at DESC, COALESCE(v.visit_id, 0) DESC
				LIMIT 1
			), '') AS file_url,
			ovl.skip_at,
			ovl.skip_reason,
			ovl."start",
			ovl.finish,
			TRUE AS is_extra_call
		`, visitStartAt, visitEndAt).
		Joins("JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = dh.pjp_id AND pjp.cust_id = dh.cust_id").
		Joins("LEFT JOIN pjp_principles.routes r ON r.route_code = dh.route_code").
		Joins("JOIN mst.m_salesman ms2 ON pjp.salesman_id = ms2.emp_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms2.emp_id").
		Joins("JOIN smc.m_customer mc ON ms2.cust_id = mc.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Joins("LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.pjp_id = dh.pjp_id AND ovl.date = DATE(dh.date) AND ovl.is_extra_call = TRUE AND ovl.outlet_code = dh.destination_code").
		Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms WHERE ms.cust_id IN ?)", custIDs).
		Where("DATE(dh.date) = ?", date).
		Where("dh.is_extra_call = ?", true).
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

	query = query.Order("me.emp_id, pjp.id, COALESCE(r.route_code, dh.route_code), dh.id")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
