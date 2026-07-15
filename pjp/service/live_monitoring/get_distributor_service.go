package live_monitoring

import (
	"context"
	"fmt"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
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

// GetDistributorMonitoring retrieves distributor monitoring data and transforms it to response format
func (s *liveMonitoringService) GetDistributorMonitoring(
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

	custIDs, err := s.repository.GetChildCustIDs(ctx, s.db, custID)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}
	if len(custIDs) == 0 {
		custIDs = []string{custID}
	}

	scopedEmpIDs, err := s.repository.GetDistributorEmployeeIDs(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, requestedEmpIDs, req.Status,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	if len(scopedEmpIDs) == 0 {
		return []response.LiveMonitoringData{}, calculatePaging(0, page, limit), nil
	}

	totalCount := int64(len(scopedEmpIDs))
	pagedEmpIDs := scopedEmpIDs
	if offset < len(scopedEmpIDs) {
		end := offset + limit
		if end > len(scopedEmpIDs) {
			end = len(scopedEmpIDs)
		}
		pagedEmpIDs = scopedEmpIDs[offset:end]
	} else {
		pagedEmpIDs = []int{}
	}

	// If no data, return empty result
	if totalCount == 0 {
		return []response.LiveMonitoringData{}, calculatePaging(0, page, limit), nil
	}

	if len(pagedEmpIDs) == 0 {
		return []response.LiveMonitoringData{}, calculatePaging(totalCount, page, limit), nil
	}

	rows, err := s.repository.GetDistributorMonitoring(
		ctx, s.db, custIDs, date, req.RegionID, req.AreaID, req.DistributorID, pagedEmpIDs, req.Status, 0, 0,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	routeCodes, outletIDs := collectDistributorRouteAndOutletKeys(rows)

	employeeMetaMap, err := s.repository.GetDistributorEmployeeMeta(
		ctx, s.db, custIDs, pagedEmpIDs,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	routeMetaMap, err := s.repository.GetDistributorRouteMeta(
		ctx, s.db, custIDs, routeCodes,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	outletMetaMap, err := s.repository.GetDistributorOutletMeta(
		ctx, s.db, custIDs, outletIDs,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	enrichDistributorRowsWithMetadata(rows, employeeMetaMap, routeMetaMap, outletMetaMap)

	latestVisitCoordinateMap, err := s.repository.GetDistributorLatestVisitCoordinates(
		ctx, s.db, custIDs, date, pagedEmpIDs,
	)
	if err != nil {
		return nil, response.LiveMonitoringPaging{}, err
	}

	enrichDistributorRowsWithLatestVisits(rows, latestVisitCoordinateMap)

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

	// Transform raw data to response format
	result := transformDistributorRows(rows)
	enrichDistributorMonitoringData(date, result, attendanceMap, currentCoordinateMap)

	return result, calculatePaging(totalCount, page, limit), nil
}

func enrichDistributorRowsWithLatestVisits(rows []model.LiveMonitoringDistributorRow, latestVisitCoordinateMap map[string]model.LatestVisitCoordinateRow) {
	for index := range rows {
		visitCoordinate, exists := latestVisitCoordinateMap[buildDistributorVisitCoordinateKey(rows[index].CustID, rows[index].SalesmanCode, rows[index].OutletCode)]
		if !exists {
			continue
		}

		rows[index].ArriveLongitude = visitCoordinate.ArriveLongitude
		rows[index].ArriveLatitude = visitCoordinate.ArriveLatitude
		rows[index].FileURL = visitCoordinate.FileURL
	}
}

func enrichDistributorRowsWithMetadata(
	rows []model.LiveMonitoringDistributorRow,
	employeeMetaMap map[int]model.DistributorEmployeeMetaRow,
	routeMetaMap map[string]model.DistributorRouteMetaRow,
	outletMetaMap map[string]model.DistributorOutletMetaRow,
) {
	for index := range rows {
		employeeMeta, hasEmployeeMeta := employeeMetaMap[rows[index].EmpID]
		if hasEmployeeMeta {
			rows[index].EmpCode = employeeMeta.EmpCode
			rows[index].EmpName = employeeMeta.EmpName
			rows[index].DistributorID = employeeMeta.DistributorID
			rows[index].AreaID = employeeMeta.AreaID
			rows[index].RegionID = employeeMeta.RegionID
		}

		routeMeta, hasRouteMeta := routeMetaMap[buildDistributorRouteMetaKey(rows[index].CustID, rows[index].RouteCode)]
		if hasRouteMeta {
			rows[index].RouteName = routeMeta.RouteName
		}

		outletMeta, hasOutletMeta := outletMetaMap[buildDistributorOutletMetaKey(rows[index].CustID, rows[index].OutletID)]
		if hasOutletMeta {
			rows[index].OutletCode = outletMeta.OutletCode
			rows[index].OutletName = outletMeta.OutletName
		}
	}
}

func collectDistributorRouteAndOutletKeys(rows []model.LiveMonitoringDistributorRow) ([]int64, []int) {
	routeCodeMap := make(map[int64]struct{})
	outletIDMap := make(map[int]struct{})

	for _, row := range rows {
		routeCodeMap[row.RouteCode] = struct{}{}
		outletIDMap[row.OutletID] = struct{}{}
	}

	routeCodes := make([]int64, 0, len(routeCodeMap))
	for routeCode := range routeCodeMap {
		routeCodes = append(routeCodes, routeCode)
	}

	outletIDs := make([]int, 0, len(outletIDMap))
	for outletID := range outletIDMap {
		outletIDs = append(outletIDs, outletID)
	}

	return routeCodes, outletIDs
}

func enrichDistributorMonitoringData(
	date string,
	result []response.LiveMonitoringData,
	attendanceMap map[int]model.AttendanceRow,
	currentCoordinateMap map[int]model.CurrentCoordinateRow,
) {
	dayStartEpoch, dayEndEpoch, err := buildJakartaBusinessDayEpochRange(date)
	if err != nil {
		for i := range result {
			resetDistributorDailyTrackingState(&result[i])
		}
		return
	}

	for i := range result {
		attendance, hasAttendance := attendanceMap[result[i].EmpID]
		hasAttendanceOnRequestedDate := hasAttendance && attendance.AttendanceID != nil && isBusinessDayEpoch(attendance.Timestamp, dayStartEpoch, dayEndEpoch)
		if hasAttendanceOnRequestedDate {
			result[i].AttendanceID = attendance.AttendanceID
			result[i].AttendanceLongitude = attendance.Longitude
			result[i].AttendanceLatitude = attendance.Latitude
			result[i].AttendanceAt = attendance.Timestamp
			result[i].ClockOut = attendance.ClockOutID
			result[i].ClockOutLongitude = attendance.ClockOutLong
			result[i].ClockOutLatitude = attendance.ClockOutLat
			result[i].ClockOutAt = attendance.ClockOutAt
		}

		if !hasAttendanceOnRequestedDate {
			resetDistributorDailyTrackingState(&result[i])
			continue
		}

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

		for pjpIndex := range result[i].PjpData {
			filterDistributorRoutesDailyTrackingState(result[i].PjpData[pjpIndex].RouteData, dayStartEpoch, dayEndEpoch)
			filterDistributorRoutesDailyTrackingState(result[i].PjpData[pjpIndex].ExtraCallData, dayStartEpoch, dayEndEpoch)
		}
	}
}

func resetDistributorDailyTrackingState(data *response.LiveMonitoringData) {
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

	for pjpIndex := range data.PjpData {
		resetDistributorRoutesDailyTrackingState(data.PjpData[pjpIndex].RouteData)
		resetDistributorRoutesDailyTrackingState(data.PjpData[pjpIndex].ExtraCallData)
	}
}

func resetDistributorRoutesDailyTrackingState(routes []response.LiveMonitoringRouteData) {
	for routeIndex := range routes {
		for destinationIndex := range routes[routeIndex].DestinationData {
			resetDistributorDestinationDailyTrackingState(&routes[routeIndex].DestinationData[destinationIndex])
		}
	}
}

func filterDistributorRoutesDailyTrackingState(routes []response.LiveMonitoringRouteData, dayStartEpoch, dayEndEpoch int64) {
	for routeIndex := range routes {
		for destinationIndex := range routes[routeIndex].DestinationData {
			filterDistributorDestinationDailyTrackingState(&routes[routeIndex].DestinationData[destinationIndex], dayStartEpoch, dayEndEpoch)
		}
	}
}

func filterDistributorDestinationDailyTrackingState(destination *response.LiveMonitoringDestinationData, dayStartEpoch, dayEndEpoch int64) {
	hasArriveAt := isBusinessDayEpoch(destination.ArriveAt, dayStartEpoch, dayEndEpoch)
	hasLeaveAt := isBusinessDayEpoch(destination.LeaveAt, dayStartEpoch, dayEndEpoch)
	hasStart := isBusinessDayEpoch(destination.Start, dayStartEpoch, dayEndEpoch)
	hasFinish := isBusinessDayEpoch(destination.Finish, dayStartEpoch, dayEndEpoch)
	hasSkipAt := isBusinessDayEpoch(destination.SkipAt, dayStartEpoch, dayEndEpoch)

	if !hasArriveAt {
		destination.ArriveAt = nil
		destination.ArriveLongitude = 0
		destination.ArriveLatitude = 0
	}

	if !hasLeaveAt {
		destination.LeaveAt = nil
	}

	if !hasStart {
		destination.Start = nil
	}

	if !hasFinish {
		destination.Finish = nil
	}

	if !hasSkipAt {
		destination.SkipAt = nil
		destination.SkipReason = nil
	}

	if !hasArriveAt && !hasLeaveAt && !hasStart && !hasFinish && !hasSkipAt {
		resetDistributorDestinationDailyTrackingState(destination)
	}
}

func resetDistributorDestinationDailyTrackingState(destination *response.LiveMonitoringDestinationData) {
	destination.ArriveAt = nil
	destination.LeaveAt = nil
	destination.ArriveLongitude = 0
	destination.ArriveLatitude = 0
	destination.Start = nil
	destination.Finish = nil
	destination.SkipAt = nil
	destination.SkipReason = nil
}

// transformDistributorRows transforms raw database rows to hierarchical response structure
func transformDistributorRows(rows []model.LiveMonitoringDistributorRow) []response.LiveMonitoringData {
	// Use maps to group data by employee -> pjp -> route -> destinations
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
				EmpID:                   row.EmpID,
				EmpCode:                 row.EmpCode,
				EmpName:                 row.EmpName,
				DistributorID:           row.DistributorID,
				AreaID:                  row.AreaID,
				RegionID:                row.RegionID,
				AttendanceID:            row.AttendanceID,
				AttendanceLongitude:     row.AttendanceLongitude,
				AttendanceLatitude:      row.AttendanceLatitude,
				AttendanceAt:            row.AttendanceAt,
				CurrentLongitude:        row.CurrentLongitude,
				CurrentLatitude:         row.CurrentLatitude,
				CurrentCoordinateAt:     row.CurrentCoordinateAt,
				CurrentCoordinateSource: row.CurrentSource,
				PjpData:                 []response.LiveMonitoringPjpData{},
			}
			empOrder = append(empOrder, empKey)
		}

		pjpKey := fmt.Sprintf("%d_%d", row.EmpID, row.PjpID)
		if _, exists := pjpMap[pjpKey]; !exists {
			pjpMap[pjpKey] = &response.LiveMonitoringPjpData{
				PjpID:          row.PjpID,
				ApprovalStatus: row.ApprovalStatus,
				RouteData:      []response.LiveMonitoringRouteData{},
				ExtraCallData:  []response.LiveMonitoringRouteData{},
			}
			pjpOrder[row.EmpID] = append(pjpOrder[row.EmpID], pjpKey)
		}

		routeKey := fmt.Sprintf("%d_%d_%d", row.EmpID, row.PjpID, row.RouteCode)
		targetMap := regularRouteMap
		targetOrder := regularRouteOrder
		if row.IsExtraCall {
			targetMap = extraRouteMap
			targetOrder = extraRouteOrder
		}

		if _, exists := targetMap[routeKey]; !exists {
			targetMap[routeKey] = &response.LiveMonitoringRouteData{
				RouteCode:       fmt.Sprintf("%d", row.RouteCode),
				RouteName:       row.RouteName,
				DestinationData: []response.LiveMonitoringDestinationData{},
			}
			targetOrder[pjpKey] = append(targetOrder[pjpKey], routeKey)
		}

		destinationType := row.DestinationType
		if destinationType == "" {
			destinationType = "Outlet"
		}

		targetMap[routeKey].DestinationData = append(targetMap[routeKey].DestinationData, response.LiveMonitoringDestinationData{
			DestinationID:   row.OutletID,
			DestinationCode: row.OutletCode,
			DestinationName: row.OutletName,
			DestinationType: destinationType,
			Longitude:       row.Longitude,
			Latitude:        row.Latitude,
			ArriveAt:        row.ArriveAt,
			LeaveAt:         row.LeaveAt,
			ArriveLongitude: row.ArriveLongitude,
			ArriveLatitude:  row.ArriveLatitude,
			LeaveLongitude:  row.LeaveLongitude,
			LeaveLatitude:   row.LeaveLatitude,
			FileURL:         row.FileURL,
			Start:           row.Start,
			Finish:          row.Finish,
			SkipAt:          row.SkipAt,
			SkipReason:      row.SkipReason,
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
