package live_monitoring

import (
	"context"
	"fmt"
	"scyllax-pjp/model"
	"time"

	"gorm.io/gorm"
)

func buildDistributorVisitCoordinateKey(custID, empCode, outletCode string) string {
	return fmt.Sprintf("%s|%s|%s", custID, empCode, outletCode)
}

func buildDistributorRouteMetaKey(custID string, routeCode int64) string {
	return fmt.Sprintf("%s|%d", custID, routeCode)
}

func buildDistributorOutletMetaKey(custID string, outletID int) string {
	return fmt.Sprintf("%s|%d", custID, outletID)
}

var liveMonitoringJakartaLocation = loadLiveMonitoringJakartaLocation()

func loadLiveMonitoringJakartaLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("WIB", 7*60*60)
	}

	return location
}

func buildLiveMonitoringDayRange(date string) (string, string, error) {
	startAt, err := time.ParseInLocation("2006-01-02", date, liveMonitoringJakartaLocation)
	if err != nil {
		return "", "", err
	}

	endAt := startAt.AddDate(0, 0, 1)

	return startAt.Format("2006-01-02 15:04:05"), endAt.Format("2006-01-02 15:04:05"), nil
}

func (r *liveMonitoringRepository) buildDistributorEmployeeScopeQuery(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) *gorm.DB {
	query := tx.WithContext(ctx).Table("pjp.route_outlet_history roh").
		Select("DISTINCT pjp.salesman_id").
		Joins("JOIN pjp.permanent_journey_plans pjp ON roh.pjp_id = pjp.id AND roh.cust_id = pjp.cust_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id").
		Where("roh.date = ?", date).
		Where("pjp.approval_status IN ?", statuses)

	if len(custIDs) > 0 {
		query = query.Where("pjp.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		query = query.Where("pjp.salesman_id IN ?", empIDs)
	}

	if distributorID > 0 {
		query = query.Where("md.distributor_id = ?", distributorID)
	}
	if areaID > 0 {
		query = query.Where("md.area_id = ?", areaID)
	}
	if regionID > 0 {
		query = query.Where("md.region_id = ?", regionID)
	}

	return query
}

func (r *liveMonitoringRepository) GetDistributorEmployeeIDs(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) ([]int, error) {
	var employeeIDs []int

	query := r.buildDistributorEmployeeScopeQuery(ctx, tx, custIDs, date, regionID, areaID, distributorID, empIDs, statuses).
		Order("pjp.salesman_id")

	if err := query.Pluck("pjp.salesman_id", &employeeIDs).Error; err != nil {
		return nil, err
	}

	return employeeIDs, nil
}

// GetDistributorMonitoring retrieves distributor monitoring data from pjp schema
func (r *liveMonitoringRepository) GetDistributorMonitoring(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
	limit, offset int,
) ([]model.LiveMonitoringDistributorRow, error) {
	var results []model.LiveMonitoringDistributorRow

	query := tx.WithContext(ctx).Table("pjp.route_outlet_history roh").
		Select(`
			roh.cust_id,
			pjp.salesman_id AS emp_id,
			pjp.salesman_code,
			pjp.id AS pjp_id,
			pjp.approval_status,
			roh.route_code,
			roh.outlet_id,
			CASE WHEN roh.is_extra_call THEN 'Distributor' ELSE 'Outlet' END AS destination_type,
			COALESCE(CAST(roh.longitude AS FLOAT), 0) AS longitude,
			COALESCE(CAST(roh.latitude AS FLOAT), 0) AS latitude,
			ovl.arrive_at,
			ovl.leave_at,
			0 AS arrive_longitude,
			0 AS arrive_latitude,
			NULLIF(ovl.leave_longitude, '') AS leave_longitude,
			NULLIF(ovl.leave_latitude, '') AS leave_latitude,
			ovl."start",
			ovl.finish,
			ovl.skip_at,
			ovl.skip_reason,
			roh.is_extra_call
		`).
		Joins("JOIN pjp.permanent_journey_plans pjp ON roh.pjp_id = pjp.id AND roh.cust_id = pjp.cust_id").
		Joins("LEFT JOIN pjp.outlet_visit_list ovl ON ovl.outlet_id = roh.outlet_id AND ovl.pjp_id = roh.pjp_id AND ovl.date = ? AND ovl.is_extra_call = roh.is_extra_call", date).
		Where("roh.date = ?", date).
		Where("pjp.approval_status IN ?", statuses)

	if len(custIDs) > 0 {
		query = query.Where("pjp.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		query = query.Where("pjp.salesman_id IN ?", empIDs)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	query = query.Order("pjp.salesman_id, pjp.id, roh.is_extra_call, roh.route_code, roh.outlet_id")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *liveMonitoringRepository) GetDistributorLatestVisitCoordinates(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empIDs []int,
) (map[string]model.LatestVisitCoordinateRow, error) {
	var rows []model.LatestVisitCoordinateRow

	visitStartAt, visitEndAt, err := buildLiveMonitoringDayRange(date)
	if err != nil {
		return nil, err
	}

	baseQuery := tx.WithContext(ctx).Table("mobile.visits mv").
		Select(`
			mv.cust_id,
			mv.emp_code,
			mv.outlet_code,
			COALESCE(CAST(NULLIF(mv.longitude, '') AS DOUBLE PRECISION), 0) AS arrive_longitude,
			COALESCE(CAST(NULLIF(mv.latitude, '') AS DOUBLE PRECISION), 0) AS arrive_latitude,
			NULLIF(mv.file_url, '') AS file_url,
			ROW_NUMBER() OVER (
				PARTITION BY mv.cust_id, mv.emp_code, mv.outlet_code
				ORDER BY mv.created_at DESC, COALESCE(mv.visit_id, 0) DESC
			) AS row_number
		`).
		Where("mv.created_at >= ? AND mv.created_at < ?", visitStartAt, visitEndAt)

	if len(custIDs) > 0 {
		baseQuery = baseQuery.Where("mv.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		baseQuery = baseQuery.Joins("JOIN mst.m_employee me ON me.emp_code = mv.emp_code AND me.cust_id = mv.cust_id").Where("me.emp_id IN ?", empIDs)
	}

	query := tx.WithContext(ctx).Table("(?) AS latest_visits", baseQuery).
		Select(`
			cust_id,
			emp_code,
			outlet_code,
			arrive_longitude,
			arrive_latitude,
			file_url
		`).
		Where("row_number = 1")

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]model.LatestVisitCoordinateRow, len(rows))
	for _, row := range rows {
		result[buildDistributorVisitCoordinateKey(row.CustID, row.EmpCode, row.OutletCode)] = row
	}

	return result, nil
}

func (r *liveMonitoringRepository) GetDistributorEmployeeMeta(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	empIDs []int,
) (map[int]model.DistributorEmployeeMetaRow, error) {
	var rows []model.DistributorEmployeeMetaRow

	query := tx.WithContext(ctx).Table("mst.m_salesman ms").
		Select(`
			me.emp_id,
			me.emp_code,
			me.emp_name,
			md.distributor_id,
			md.area_id,
			md.region_id
		`).
		Joins("JOIN mst.m_employee me ON me.emp_id = ms.emp_id AND me.cust_id = ms.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id")

	if len(custIDs) > 0 {
		query = query.Where("ms.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		query = query.Where("ms.emp_id IN ?", empIDs)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[int]model.DistributorEmployeeMetaRow, len(rows))
	for _, row := range rows {
		result[row.EmpID] = row
	}

	return result, nil
}

func (r *liveMonitoringRepository) GetDistributorRouteMeta(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	routeCodes []int64,
) (map[string]model.DistributorRouteMetaRow, error) {
	if len(routeCodes) == 0 {
		return map[string]model.DistributorRouteMetaRow{}, nil
	}

	var rows []model.DistributorRouteMetaRow

	query := tx.WithContext(ctx).Table("pjp.routes r").
		Select(`
			r.cust_id,
			r.route_code,
			r.route_name
		`).
		Where("r.route_code IN ?", routeCodes)

	if len(custIDs) > 0 {
		query = query.Where("r.cust_id IN ?", custIDs)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]model.DistributorRouteMetaRow, len(rows))
	for _, row := range rows {
		result[buildDistributorRouteMetaKey(row.CustID, row.RouteCode)] = row
	}

	return result, nil
}

func (r *liveMonitoringRepository) GetDistributorOutletMeta(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	outletIDs []int,
) (map[string]model.DistributorOutletMetaRow, error) {
	if len(outletIDs) == 0 {
		return map[string]model.DistributorOutletMetaRow{}, nil
	}

	var rows []model.DistributorOutletMetaRow

	query := tx.WithContext(ctx).Table("mst.m_outlet mo").
		Select(`
			mo.cust_id,
			mo.outlet_id,
			mo.outlet_code,
			mo.outlet_name
		`).
		Where("mo.outlet_id IN ?", outletIDs)

	if len(custIDs) > 0 {
		query = query.Where("mo.cust_id IN ?", custIDs)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]model.DistributorOutletMetaRow, len(rows))
	for _, row := range rows {
		result[buildDistributorOutletMetaKey(row.CustID, row.OutletID)] = row
	}

	return result, nil
}

func (r *liveMonitoringRepository) getDistributorAttendanceByType(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empIDs []int,
	attendanceType int,
	selectClause string,
) ([]model.AttendanceRow, error) {
	var rows []model.AttendanceRow

	dayStartAt, dayEndAt, err := buildLiveMonitoringDayRange(date)
	if err != nil {
		return nil, err
	}

	baseQuery := tx.WithContext(ctx).Table("mobile.attendances a").
		Select(selectClause).
		Joins("JOIN mst.m_employee me ON me.emp_code = a.emp_code AND me.cust_id = a.cust_id").
		Where("a.created_at >= ? AND a.created_at < ?", dayStartAt, dayEndAt).
		Where("a.type = ?", attendanceType)

	if len(custIDs) > 0 {
		baseQuery = baseQuery.Where("a.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		baseQuery = baseQuery.Where("me.emp_id IN ?", empIDs)
	}

	query := tx.WithContext(ctx).Table("(?) AS distributor_attendance", baseQuery).
		Where("row_number = 1")

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

// GetDistributorAttendance retrieves first daily check-in and latest daily check-out per employee for distributor monitoring.
func (r *liveMonitoringRepository) GetDistributorAttendance(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) (map[int]model.AttendanceRow, error) {
	checkInRows, err := r.getDistributorAttendanceByType(
		ctx,
		tx,
		custIDs,
		date,
		empIDs,
		1,
		`
			me.emp_id,
			a.attendance_id,
			a.created_at::text AS created_at,
			EXTRACT(EPOCH FROM a.created_at)::bigint AS attendance_at,
			COALESCE(CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION), 0) AS attendance_longitude,
			COALESCE(CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION), 0) AS attendance_latitude,
			ROW_NUMBER() OVER (
				PARTITION BY me.emp_id
				ORDER BY a.created_at ASC, a.attendance_id ASC
			) AS row_number
		`,
	)
	if err != nil {
		return nil, err
	}

	checkOutRows, err := r.getDistributorAttendanceByType(
		ctx,
		tx,
		custIDs,
		date,
		empIDs,
		2,
		`
			me.emp_id,
			a.attendance_id AS clock_out,
			EXTRACT(EPOCH FROM a.created_at)::bigint AS clock_out_at,
			COALESCE(CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION), 0) AS clock_out_longitude,
			COALESCE(CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION), 0) AS clock_out_latitude,
			ROW_NUMBER() OVER (
				PARTITION BY me.emp_id
				ORDER BY a.created_at DESC, a.attendance_id DESC
			) AS row_number
		`,
	)
	if err != nil {
		return nil, err
	}

	result := make(map[int]model.AttendanceRow, len(checkInRows)+len(checkOutRows))
	for _, row := range checkInRows {
		result[row.EmpID] = row
	}

	for _, row := range checkOutRows {
		current := result[row.EmpID]
		current.EmpID = row.EmpID
		current.ClockOutID = row.ClockOutID
		current.ClockOutAt = row.ClockOutAt
		current.ClockOutLong = row.ClockOutLong
		current.ClockOutLat = row.ClockOutLat
		result[row.EmpID] = current
	}

	return result, nil
}

// GetDistributorCurrentCoordinates resolves the latest valid coordinate per employee.
// Priority order:
// 1. Latest valid attendance checkout
// 2. Latest valid attendance any type
// 3. Latest valid mobile.visits record
// 4. Latest valid pjp.outlet_visit_list record
func (r *liveMonitoringRepository) GetDistributorCurrentCoordinates(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) (map[int]model.CurrentCoordinateRow, error) {
	var rows []model.CurrentCoordinateRow

	dayStartAt, dayEndAt, err := buildLiveMonitoringDayRange(date)
	if err != nil {
		return nil, err
	}

	attendanceCheckoutQuery := tx.WithContext(ctx).Table("mobile.attendances a").
		Select(`
			me.emp_id,
			CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) AS current_longitude,
			CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) AS current_latitude,
			EXTRACT(EPOCH FROM a.created_at)::bigint AS current_coordinate_at,
			'attendance_checkout' AS current_coordinate_source,
			1 AS source_rank,
			a.attendance_id AS source_record
		`).
		Joins("JOIN mst.m_employee me ON me.emp_code = a.emp_code AND me.cust_id = a.cust_id").
		Where("a.created_at >= ? AND a.created_at < ?", dayStartAt, dayEndAt).
		Where("a.type = 2").
		Where("CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) BETWEEN -90 AND 90").
		Where("CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) BETWEEN -180 AND 180").
		Where("CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) <> 0").
		Where("CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) <> 0").
		Where("NULLIF(a.latitude, '') IS NOT NULL").
		Where("NULLIF(a.longitude, '') IS NOT NULL")

	attendanceAnyQuery := tx.WithContext(ctx).Table("mobile.attendances a").
		Select(`
			me.emp_id,
			CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) AS current_longitude,
			CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) AS current_latitude,
			EXTRACT(EPOCH FROM a.created_at)::bigint AS current_coordinate_at,
			'attendance' AS current_coordinate_source,
			2 AS source_rank,
			a.attendance_id AS source_record
		`).
		Joins("JOIN mst.m_employee me ON me.emp_code = a.emp_code AND me.cust_id = a.cust_id").
		Where("a.created_at >= ? AND a.created_at < ?", dayStartAt, dayEndAt).
		Where("CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) BETWEEN -90 AND 90").
		Where("CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) BETWEEN -180 AND 180").
		Where("CAST(NULLIF(a.latitude, '') AS DOUBLE PRECISION) <> 0").
		Where("CAST(NULLIF(a.longitude, '') AS DOUBLE PRECISION) <> 0").
		Where("NULLIF(a.latitude, '') IS NOT NULL").
		Where("NULLIF(a.longitude, '') IS NOT NULL")

	visitsQuery := tx.WithContext(ctx).Table("mobile.visits v").
		Select(`
			me.emp_id,
			CAST(NULLIF(v.longitude, '') AS DOUBLE PRECISION) AS current_longitude,
			CAST(NULLIF(v.latitude, '') AS DOUBLE PRECISION) AS current_latitude,
			EXTRACT(EPOCH FROM v.created_at)::bigint AS current_coordinate_at,
			'mobile_visit' AS current_coordinate_source,
			3 AS source_rank,
			COALESCE(v.visit_id, EXTRACT(EPOCH FROM v.created_at)::bigint) AS source_record
		`).
		Joins("JOIN mst.m_employee me ON me.emp_code = v.emp_code AND me.cust_id = v.cust_id").
		Where("v.created_at >= ? AND v.created_at < ?", dayStartAt, dayEndAt).
		Where("CAST(NULLIF(v.latitude, '') AS DOUBLE PRECISION) BETWEEN -90 AND 90").
		Where("CAST(NULLIF(v.longitude, '') AS DOUBLE PRECISION) BETWEEN -180 AND 180").
		Where("CAST(NULLIF(v.latitude, '') AS DOUBLE PRECISION) <> 0").
		Where("CAST(NULLIF(v.longitude, '') AS DOUBLE PRECISION) <> 0").
		Where("NULLIF(v.latitude, '') IS NOT NULL").
		Where("NULLIF(v.longitude, '') IS NOT NULL")

	outletVisitQuery := tx.WithContext(ctx).Table("pjp.outlet_visit_list ovl").
		Select(`
			pjp.salesman_id AS emp_id,
			CAST(NULLIF(ovl.longitude, '') AS DOUBLE PRECISION) AS current_longitude,
			CAST(NULLIF(ovl.latitude, '') AS DOUBLE PRECISION) AS current_latitude,
			COALESCE(ovl.leave_at, ovl.skip_at, ovl.finish, ovl.arrive_at, ovl.start) AS current_coordinate_at,
			'outlet_visit_list' AS current_coordinate_source,
			4 AS source_rank,
			ovl.id AS source_record
		`).
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = ovl.pjp_id").
		Where("ovl.date = ?", date).
		Where("CAST(NULLIF(ovl.latitude, '') AS DOUBLE PRECISION) BETWEEN -90 AND 90").
		Where("CAST(NULLIF(ovl.longitude, '') AS DOUBLE PRECISION) BETWEEN -180 AND 180").
		Where("CAST(NULLIF(ovl.latitude, '') AS DOUBLE PRECISION) <> 0").
		Where("CAST(NULLIF(ovl.longitude, '') AS DOUBLE PRECISION) <> 0").
		Where("NULLIF(ovl.latitude, '') IS NOT NULL").
		Where("NULLIF(ovl.longitude, '') IS NOT NULL").
		Where("COALESCE(ovl.leave_at, ovl.skip_at, ovl.finish, ovl.arrive_at, ovl.start) IS NOT NULL")

	if len(custIDs) > 0 {
		attendanceCheckoutQuery = attendanceCheckoutQuery.Where("a.cust_id IN ?", custIDs)
		attendanceAnyQuery = attendanceAnyQuery.Where("a.cust_id IN ?", custIDs)
		visitsQuery = visitsQuery.Where("v.cust_id IN ?", custIDs)
		outletVisitQuery = outletVisitQuery.Where("pjp.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		attendanceCheckoutQuery = attendanceCheckoutQuery.Where("me.emp_id IN ?", empIDs)
		attendanceAnyQuery = attendanceAnyQuery.Where("me.emp_id IN ?", empIDs)
		visitsQuery = visitsQuery.Where("me.emp_id IN ?", empIDs)
		outletVisitQuery = outletVisitQuery.Where("pjp.salesman_id IN ?", empIDs)
	}

	unionQuery := tx.WithContext(ctx).Table("(?) AS current_coordinate_candidates",
		tx.Raw("(?) UNION ALL (?) UNION ALL (?) UNION ALL (?)", attendanceCheckoutQuery, attendanceAnyQuery, visitsQuery, outletVisitQuery),
	).Select(`
		emp_id,
		current_longitude,
		current_latitude,
		current_coordinate_at,
		current_coordinate_source,
		source_rank,
		source_record
	`)

	if err := unionQuery.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[int]model.CurrentCoordinateRow, len(rows))
	for _, row := range rows {
		current, exists := result[row.EmpID]
		if !exists || shouldReplaceCurrentCoordinate(current, row) {
			result[row.EmpID] = row
		}
	}

	return result, nil
}

// CountDistributorMonitoring counts unique employees for pagination
func (r *liveMonitoringRepository) CountDistributorMonitoring(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	regionID, areaID, distributorID int,
	empIDs []int,
	statuses []string,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.route_outlet_history roh").
		Select("COUNT(DISTINCT me.emp_id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON roh.pjp_id = pjp.id AND roh.cust_id = pjp.cust_id").
		Joins("JOIN mst.m_salesman ms ON pjp.salesman_id = ms.emp_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms.emp_id AND me.cust_id = ms.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.cust_id = mc.cust_id").
		Where("roh.date = ?", date).
		Where("pjp.approval_status IN ?", statuses)

	if len(custIDs) > 0 {
		query = query.Where("pjp.cust_id IN ?", custIDs)
	}

	if len(empIDs) > 0 {
		query = query.Where("me.emp_id IN ?", empIDs)
	}

	if distributorID > 0 {
		query = query.Where("md.distributor_id = ?", distributorID)
	}
	if areaID > 0 {
		query = query.Where("md.area_id = ?", areaID)
	}
	if regionID > 0 {
		query = query.Where("md.region_id = ?", regionID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
