package live_monitoring

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"scyllax-pjp/model"

	"gorm.io/gorm"
)

const numericCoordPattern = `^[-+]{0,1}([0-9]+([.][0-9]*){0,1}|[.][0-9]+)([eE][-+]{0,1}[0-9]+){0,1}$`

func (r *liveMonitoringRepository) GetEmployeeRole(ctx context.Context, tx *gorm.DB, empID int, jwtCust string) (string, error) {
	var custID string
	err := tx.WithContext(ctx).Raw(`SELECT cust_id FROM mst.m_employee WHERE emp_id = ? AND cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = ? OR parent_cust_id = ?) LIMIT 1`, empID, jwtCust, jwtCust).Scan(&custID).Error
	if err != nil {
		return "", err
	}
	if custID == "" {
		return "", gorm.ErrRecordNotFound
	}
	return custID, nil
}

func (r *liveMonitoringRepository) GetUpdateLocations(ctx context.Context, tx *gorm.DB, empID int, date string, jwtCust string, branch string) ([]model.UpdateLocationRow, error) {
	if branch != "pjp" && branch != "pjp_principles" {
		return nil, errors.New("invalid update locations branch")
	}
	if strings.ContainsAny(branch, "'\"`; ") {
		return nil, errors.New("invalid update locations branch")
	}
	pat := numericCoordPattern
	branchTable := branch + ".outlet_visit_list"
	branchPlans := branch + ".permanent_journey_plans"
	query := fmt.Sprintf(`
SELECT type, latitude, longitude, destination_id, destination_type, destination_name, recorded_at
FROM (
 SELECT CASE WHEN a.type = 1 THEN 'clock_in' ELSE 'clock_out' END AS type,
        COALESCE(CASE WHEN btrim(a.latitude) ~ '%s' THEN btrim(a.latitude)::double precision END, 0::double precision) AS latitude,
        COALESCE(CASE WHEN btrim(a.longitude) ~ '%s' THEN btrim(a.longitude)::double precision END, 0::double precision) AS longitude,
        NULL::bigint AS destination_id, NULL::text AS destination_type,
        NULL::text AS destination_name, a.created_at AS recorded_at
 FROM mobile.attendances a
 JOIN mst.m_employee e ON e.emp_code = a.emp_code
 WHERE e.emp_id = ? AND a.created_at >= (?::date AT TIME ZONE 'Asia/Jakarta') AND a.created_at < ((?::date + INTERVAL '1 day') AT TIME ZONE 'Asia/Jakarta')
   AND e.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = ? OR parent_cust_id = ?)
 UNION ALL
 SELECT 'gps',
        COALESCE(CASE WHEN btrim(ul.latitude) ~ '%s' THEN btrim(ul.latitude)::double precision END, 0::double precision),
        COALESCE(CASE WHEN btrim(ul.longitude) ~ '%s' THEN btrim(ul.longitude)::double precision END, 0::double precision),
        NULL::bigint, NULL::text, NULL::text, ul.created_at
 FROM sys.user_location ul
 WHERE ul.emp_id = ? AND ul.created_at >= (?::date AT TIME ZONE 'Asia/Jakarta') AND ul.created_at < ((?::date + INTERVAL '1 day') AT TIME ZONE 'Asia/Jakarta')
   AND ul.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = ? OR parent_cust_id = ?)
 UNION ALL
  SELECT v.type,
         COALESCE(CASE WHEN btrim(v.latitude) ~ '%s' THEN btrim(v.latitude)::double precision END, 0::double precision),
         COALESCE(CASE WHEN btrim(v.longitude) ~ '%s' THEN btrim(v.longitude)::double precision END, 0::double precision),
         v.outlet_id, v.destination_type, o.outlet_name, v.recorded_at
 FROM (
  SELECT 'arrive'::text AS type, ov.latitude, ov.longitude, ov.outlet_id, NULL::text AS destination_type, to_timestamp(ov.arrive_at) AS recorded_at, p.salesman_id AS emp_id, p.cust_id
    FROM %s ov JOIN %s p ON p.id = ov.pjp_id
    WHERE ov.arrive_at IS NOT NULL AND ov.arrive_at > 0
    UNION ALL
    SELECT 'leave'::text, ov.leave_latitude, ov.leave_longitude, ov.outlet_id, NULL::text, to_timestamp(ov.leave_at), p.salesman_id, p.cust_id
    FROM %s ov JOIN %s p ON p.id = ov.pjp_id
    WHERE ov.leave_at IS NOT NULL AND ov.leave_at > 0
 ) v LEFT JOIN mst.m_outlet o ON o.outlet_id = v.outlet_id
 WHERE v.emp_id = ? AND v.recorded_at >= (?::date AT TIME ZONE 'Asia/Jakarta') AND v.recorded_at < ((?::date + INTERVAL '1 day') AT TIME ZONE 'Asia/Jakarta')
   AND v.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE cust_id = ? OR parent_cust_id = ?)
) timeline ORDER BY recorded_at ASC, type ASC, destination_id ASC`,
		pat, pat, pat, pat, pat, pat,
		branchTable, branchPlans, branchTable, branchPlans,
	)
	var rows []model.UpdateLocationRow
	err := tx.WithContext(ctx).Raw(query, empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust, empID, date, date, jwtCust, jwtCust).Scan(&rows).Error
	return rows, err
}
