package live_monitoring

import (
	"context"
	"fmt"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
)

// GetPrincipalMonitoring retrieves principal monitoring data and transforms it to response format.
// Pagination is applied at the employee level (not raw row level) to prevent destination truncation.
func (s *liveMonitoringService) GetPrincipalMonitoring(
	ctx context.Context,
	req request.LiveMonitoringRequest,
	custID string,
) ([]response.LiveMonitoringData, response.LiveMonitoringPaging, error) {
	// Set default pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 9999 {
		limit = 9999
	}

	offset := (page - 1) * limit

	// Convert epoch to date string
	date := epochToDateString(req.Date)
	requestedEmpIDs := req.GetEmpIDs()

	// Resolve child cust_ids with fallback (mirrors distributor pattern)
	custIDs, err := s.repository.GetChildCustIDs(ctx, s.db, custID)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}
	if len(custIDs) == 0 {
		custIDs = []string{custID}
	}

	// Page at employee scope level — not raw destination rows
	scopedEmpIDs, err := s.repository.GetPrincipalEmployeeIDs(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, requestedEmpIDs, req.Status,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	if len(scopedEmpIDs) == 0 {
		return []response.LiveMonitoringData{}, calculatePaging(0, page, limit), nil
	}

	totalCount := int64(len(scopedEmpIDs))

	// Slice the employee page
	var pagedEmpIDs []int
	if offset < len(scopedEmpIDs) {
		end := offset + limit
		if end > len(scopedEmpIDs) {
			end = len(scopedEmpIDs)
		}
		pagedEmpIDs = scopedEmpIDs[offset:end]
	} else {
		pagedEmpIDs = []int{}
	}

	if len(pagedEmpIDs) == 0 {
		return []response.LiveMonitoringData{}, calculatePaging(totalCount, page, limit), nil
	}

	// Fetch all destination rows for the paged employees — no SQL row limit
	rows, err := s.repository.GetPrincipalMonitoring(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, pagedEmpIDs, req.Status, 0, 0,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	extraCallRows, err := s.repository.GetPrincipalExtraCalls(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, pagedEmpIDs, req.Status,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	allRows := make([]model.LiveMonitoringPrincipalRow, 0, len(rows)+len(extraCallRows))
	allRows = append(allRows, rows...)
	allRows = append(allRows, extraCallRows...)

	// Transform raw data to response format
	result := transformPrincipalRows(allRows)

	// Enrich top-level attendance, clock-out, and current-coordinate fields
	// using the same repo calls and enrichment logic as the distributor path.
	attendanceMap, err := s.repository.GetDistributorAttendance(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, pagedEmpIDs, req.Status,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	currentCoordinateMap, err := s.repository.GetDistributorCurrentCoordinates(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, pagedEmpIDs, req.Status,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	enrichPrincipalMonitoringData(date, result, attendanceMap, currentCoordinateMap)

	return result, calculatePaging(totalCount, page, limit), nil
}

// enrichPrincipalMonitoringData populates top-level attendance, clock-out, and
// current-coordinate fields on principal monitoring results. It mirrors the
// distributor enrichment but skips the destination-level daily-tracking filter
// (principal destinations are already date-scoped by the query join on ovl.date).
func enrichPrincipalMonitoringData(
	date string,
	result []response.LiveMonitoringData,
	attendanceMap map[int]model.AttendanceRow,
	currentCoordinateMap map[int]model.CurrentCoordinateRow,
) {
	dayStartEpoch, dayEndEpoch, err := buildJakartaBusinessDayEpochRange(date)
	if err != nil {
		for i := range result {
			resetPrincipalAttendanceState(&result[i])
		}
		return
	}

	for i := range result {
		attendance, hasAttendance := attendanceMap[result[i].EmpID]
		hasAttendanceOnRequestedDate := hasAttendance && attendance.AttendanceID != nil && isBusinessDayEpoch(attendance.Timestamp, dayStartEpoch, dayEndEpoch)

		if !hasAttendanceOnRequestedDate {
			resetPrincipalAttendanceState(&result[i])
			continue
		}

		result[i].AttendanceID = attendance.AttendanceID
		result[i].AttendanceLongitude = attendance.Longitude
		result[i].AttendanceLatitude = attendance.Latitude
		result[i].AttendanceAt = attendance.Timestamp
		result[i].ClockOut = attendance.ClockOutID
		result[i].ClockOutLongitude = attendance.ClockOutLong
		result[i].ClockOutLatitude = attendance.ClockOutLat
		result[i].ClockOutAt = attendance.ClockOutAt

		if currentCoordinate, exists := currentCoordinateMap[result[i].EmpID]; exists && isBusinessDayEpoch(currentCoordinate.Timestamp, dayStartEpoch, dayEndEpoch) {
			result[i].CurrentLongitude = currentCoordinate.Longitude
			result[i].CurrentLatitude = currentCoordinate.Latitude
			result[i].CurrentCoordinateAt = currentCoordinate.Timestamp
			result[i].CurrentCoordinateSource = currentCoordinate.Source
		} else {
			result[i].CurrentLongitude = 0
			result[i].CurrentLatitude = 0
			result[i].CurrentCoordinateAt = nil
			result[i].CurrentCoordinateSource = ""
		}
	}
}

func resetPrincipalAttendanceState(data *response.LiveMonitoringData) {
	data.AttendanceID = nil
	data.AttendanceLongitude = 0
	data.AttendanceLatitude = 0
	data.AttendanceAt = nil
	data.ClockOut = nil
	data.ClockOutLongitude = 0
	data.ClockOutLatitude = 0
	data.ClockOutAt = nil
	data.CurrentLongitude = 0
	data.CurrentLatitude = 0
	data.CurrentCoordinateAt = nil
	data.CurrentCoordinateSource = ""
}

// transformPrincipalRows transforms raw database rows to hierarchical response structure.
// Rows with IsExtraCall=true are split into ExtraCallData; others go into RouteData.
func transformPrincipalRows(rows []model.LiveMonitoringPrincipalRow) []response.LiveMonitoringData {
	empMap := make(map[int]*response.LiveMonitoringData)
	pjpMap := make(map[string]*response.LiveMonitoringPjpData)
	regularRouteMap := make(map[string]*response.LiveMonitoringRouteData)
	extraRouteMap := make(map[string]*response.LiveMonitoringRouteData)

	var empOrder []int
	pjpOrder := make(map[int][]string)
	regularRouteOrder := make(map[string][]string)
	extraRouteOrder := make(map[string][]string)

	for _, row := range rows {
		empKey := row.EmpID
		if _, exists := empMap[empKey]; !exists {
			empMap[empKey] = &response.LiveMonitoringData{
				EmpID:         row.EmpID,
				EmpCode:       row.EmpCode,
				EmpName:       row.EmpName,
				DistributorID: row.DistributorID,
				AreaID:        row.AreaID,
				RegionID:      row.RegionID,
				PjpData:       []response.LiveMonitoringPjpData{},
			}
			empOrder = append(empOrder, empKey)
		}

		pjpKey := fmt.Sprintf("%d_%d", row.EmpID, row.PjpID)
		if _, exists := pjpMap[pjpKey]; !exists {
			pjpCode := row.PjpCode
			pjpMap[pjpKey] = &response.LiveMonitoringPjpData{
				PjpID:          row.PjpID,
				PjpCode:        &pjpCode,
				ApprovalStatus: row.ApprovalStatus,
				RouteData:      []response.LiveMonitoringRouteData{},
				ExtraCallData:  []response.LiveMonitoringRouteData{},
			}
			pjpOrder[row.EmpID] = append(pjpOrder[row.EmpID], pjpKey)
		}

		routeKey := fmt.Sprintf("%d_%d_%d", row.EmpID, row.PjpID, row.RouteCode)
		targetMap, targetOrder := splitPrincipalRouteTarget(row.IsExtraCall, regularRouteMap, extraRouteMap, regularRouteOrder, extraRouteOrder)

		if _, exists := targetMap[routeKey]; !exists {
			targetMap[routeKey] = &response.LiveMonitoringRouteData{
				RouteCode:       fmt.Sprintf("%d", row.RouteCode),
				RouteName:       row.RouteName,
				DestinationData: []response.LiveMonitoringDestinationData{},
			}
			targetOrder[pjpKey] = append(targetOrder[pjpKey], routeKey)
		}

		targetMap[routeKey].DestinationData = append(targetMap[routeKey].DestinationData, response.LiveMonitoringDestinationData{
			DestinationID:      row.DestinationID,
			DestinationCode:    row.DestinationCode,
			DestinationType:    row.DestinationType,
			DestinationName:    row.DestinationName,
			DestinationAddress: row.DestinationAddress,
			Longitude:          row.Longitude,
			Latitude:           row.Latitude,
			ArriveAt:           row.ArriveAt,
			LeaveAt:            row.LeaveAt,
			ArriveLongitude:    row.ArriveLongitude,
			ArriveLatitude:     row.ArriveLatitude,
			LeaveLongitude:     row.LeaveLongitude,
			LeaveLatitude:      row.LeaveLatitude,
			FileURL:            row.FileURL,
			Start:              row.Start,
			Finish:             row.Finish,
			SkipAt:             row.SkipAt,
			SkipReason:         row.SkipReason,
		})
	}

	result := make([]response.LiveMonitoringData, 0, len(empOrder))
	for _, empID := range empOrder {
		emp := empMap[empID]
		emp.PjpData = []response.LiveMonitoringPjpData{}

		for _, pjpKey := range pjpOrder[empID] {
			pjp := pjpMap[pjpKey]
			pjp.RouteData = []response.LiveMonitoringRouteData{}
			pjp.ExtraCallData = []response.LiveMonitoringRouteData{}

			for _, routeKey := range regularRouteOrder[pjpKey] {
				pjp.RouteData = append(pjp.RouteData, *regularRouteMap[routeKey])
			}
			for _, routeKey := range extraRouteOrder[pjpKey] {
				pjp.ExtraCallData = append(pjp.ExtraCallData, *extraRouteMap[routeKey])
			}

			emp.PjpData = append(emp.PjpData, *pjp)
		}

		result = append(result, *emp)
	}

	return result
}

// splitPrincipalRouteTarget returns the correct route map and order map based on IsExtraCall flag.
func splitPrincipalRouteTarget(
	isExtraCall bool,
	regularRouteMap, extraRouteMap map[string]*response.LiveMonitoringRouteData,
	regularRouteOrder, extraRouteOrder map[string][]string,
) (map[string]*response.LiveMonitoringRouteData, map[string][]string) {
	if isExtraCall {
		return extraRouteMap, extraRouteOrder
	}
	return regularRouteMap, regularRouteOrder
}
