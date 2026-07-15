package live_monitoring

import (
	"context"
	"scyllax-pjp/constant"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
)

const visitSummaryStatusNone = "none"
const visitSummaryStatusCompleted = "completed"

// GetMonitoringDetail retrieves detailed monitoring information for a specific employee
func (s *liveMonitoringService) GetMonitoringDetail(
	ctx context.Context,
	req request.LiveMonitoringDetailRequest,
	custID string,
	_ int64,
) (*response.LiveMonitoringDetailData, error) {
	// Determine if this is principal or distributor based on distributor_id
	isPrincipal := req.DistributorID == nil

	var visitInfo *response.VisitInformationData
	var err error

	// Get visit information based on user type
	if isPrincipal {
		visitInfo, err = s.getPrincipalVisitInfo(ctx, custID, req.Date, req.EmpID)
	} else {
		visitInfo, err = s.getDistributorVisitInfo(ctx, req.Date, req.EmpID, *req.DistributorID, custID)
	}
	if err != nil {
		return nil, err
	}

	// If no visit info found, return nil
	if visitInfo == nil {
		return nil, nil
	}

	// Get Salesman's CustID to filter Expense correctly (and other metrics)
	salesmanCustID, err := s.repository.GetSalesmanCustID(ctx, s.db, req.EmpID)
	if err != nil {
		return nil, err
	}
	// Use salesman's custID for queries
	targetCustIDs := []string{salesmanCustID}

	// Get sales data
	salesRows, err := s.repository.GetSales(ctx, s.db, targetCustIDs, req.Date, req.EmpID)
	if err != nil {
		return nil, err
	}
	sales := make([]response.SalesData, 0, len(salesRows))
	for _, row := range salesRows {
		sales = append(sales, response.SalesData{
			OutletID:   row.OutletID,
			OutletCode: row.OutletCode,
			OutletName: row.OutletName,
			SalesOrder: row.SalesOrder,
		})
	}

	// Get return data
	returnRows, err := s.repository.GetReturns(ctx, s.db, targetCustIDs, req.Date, req.EmpID)
	if err != nil {
		return nil, err
	}
	returns := make([]response.ReturnData, 0, len(returnRows))
	for _, row := range returnRows {
		returns = append(returns, response.ReturnData{
			OutletID:    row.OutletID,
			OutletCode:  row.OutletCode,
			OutletName:  row.OutletName,
			ReturnTotal: row.ReturnTotal,
		})
	}

	// Get expense data using salesman collector mapping.
	expenseRows, err := s.repository.GetExpenses(ctx, s.db, salesmanCustID, req.EmpID, req.Date)
	if err != nil {
		return nil, err
	}
	expenses := make([]response.ExpenseData, 0, len(expenseRows))
	for _, row := range expenseRows {
		expenses = append(expenses, response.ExpenseData{
			ExpenseTypeID: row.ExpenseTypeID,
			ExpenseType:   row.ExpenseTypeName,
			Note:          row.Note,
			ExpenseTotal:  row.Amount,
		})
	}

	// Get shipment data
	shipmentRows, err := s.repository.GetShipments(ctx, s.db, targetCustIDs, req.Date, req.EmpID)
	if err != nil {
		return nil, err
	}
	shipments := transformShipmentRows(shipmentRows)

	surveyRows, err := s.repository.GetSubmittedSurveyData(ctx, s.db, targetCustIDs, req.Date, req.EmpID)
	if err != nil {
		return nil, err
	}
	surveyData := make([]response.SurveyData, 0, len(surveyRows))
	for _, row := range surveyRows {
		surveyData = append(surveyData, response.SurveyData{
			Submission:  row.Submission,
			SurveyTitle: row.SurveyTitle,
			OutletCode:  row.OutletCode,
			OutletName:  row.OutletName,
		})
	}

	returnSummary := buildVisitSummary(len(returns))

	collectionRows, err := s.repository.GetCollections(ctx, s.db, targetCustIDs, req.Date, req.EmpID)
	if err != nil {
		return nil, err
	}
	collection := make([]response.CollectionData, 0, len(collectionRows))
	for _, row := range collectionRows {
		outletID := row.OutletID
		outletCode := row.OutletCode
		outletName := row.OutletName
		collectionTotal := row.CollectionTotal
		collection = append(collection, response.CollectionData{
			OutletID:        &outletID,
			OutletCode:      &outletCode,
			OutletName:      &outletName,
			CollectionTotal: &collectionTotal,
		})
	}
	collectionSummary := buildVisitSummary(len(collection))

	result := &response.LiveMonitoringDetailData{
		VisitInformation: *visitInfo,
		Sales:            sales,
		Return:           returns,
		Collection:       collection,
		Expense:          expenses,
		Shipment:         shipments,
		SurveyData:       surveyData,
	}
	result.VisitInformation.ReturnSummary = returnSummary
	result.VisitInformation.CollectionSummary = collectionSummary

	return result, nil
}

// getPrincipalVisitInfo retrieves visit information for principal users
func (s *liveMonitoringService) getPrincipalVisitInfo(
	ctx context.Context,
	custID, date string,
	empID int,
) (*response.VisitInformationData, error) {
	// Get child custIDs for principal
	childCustIDs, err := s.repository.GetChildCustIDs(ctx, s.db, custID)
	if err != nil {
		return nil, err
	}
	if len(childCustIDs) == 0 {
		childCustIDs = []string{custID}
	}

	// Get visit counts from principal schema history
	visitRow, err := s.repository.GetVisitInformationPrincipalFromHistory(ctx, s.db, childCustIDs, date, empID)
	if err != nil {
		return nil, err
	}
	if visitRow == nil {
		return nil, nil
	}

	// Get activity time (check-in time)
	activityTime, err := s.repository.GetActivityTime(ctx, s.db, date, empID)
	if err != nil {
		return nil, err
	}

	// Get user fullname for principal (using parent custID/login custID)
	userFullname, err := s.repository.GetUserFullname(ctx, s.db, custID)
	if err != nil {
		return nil, err
	}

	companyName := ""
	if userFullname != nil {
		companyName = *userFullname
	}

	visitInfo := &response.VisitInformationData{
		ActivityDate:      date,
		CompanyName:       companyName,
		CompanyCode:       constant.DefaultCompanyCode,
		Level:             constant.LevelPrincipal,
		EmpID:             visitRow.EmpID,
		EmpCode:           visitRow.EmpCode,
		EmpName:           visitRow.EmpName,
		ActivityTime:      activityTime,
		Planned:           visitRow.Plan,
		OnGoing:           visitRow.OnGoing,
		ExtraCall:         visitRow.ExtraCall,
		Visited:           visitRow.Visited,
		Skipped:           visitRow.TotalSkip,
		ReturnSummary:     defaultVisitSummary(),
		CollectionSummary: defaultVisitSummary(),
	}

	return visitInfo, nil
}

// getDistributorVisitInfo retrieves visit information for distributor users
func (s *liveMonitoringService) getDistributorVisitInfo(
	ctx context.Context,
	date string,
	empID, distributorID int,
	custID string,
) (*response.VisitInformationData, error) {
	// Get visit counts from distributor schema
	visitRow, err := s.repository.GetVisitInformationDistributor(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}
	if visitRow == nil {
		return nil, nil
	}

	plannedCount, err := s.repository.CountDistributorPlannedVisits(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}

	extraCallCount, err := s.repository.CountDistributorExtraCalls(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}

	onGoingCount, err := s.repository.CountDistributorOnGoingVisits(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}

	visitedCount, err := s.repository.CountDistributorVisitedVisits(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}

	skippedCount, err := s.repository.CountDistributorSkippedVisits(ctx, s.db, date, empID, distributorID)
	if err != nil {
		return nil, err
	}

	// Get activity time (check-in time)
	activityTime, err := s.repository.GetActivityTime(ctx, s.db, date, empID)
	if err != nil {
		return nil, err
	}

	// Get distributor info
	distInfo, err := s.repository.GetDistributorInfo(ctx, s.db, distributorID)
	if err != nil {
		return nil, err
	}

	companyName := ""
	companyCode := ""
	if distInfo != nil {
		companyName = distInfo.DistributorName
		companyCode = distInfo.DistributorCode
	}

	visitInfo := &response.VisitInformationData{
		ActivityDate:      date,
		CompanyName:       companyName,
		CompanyCode:       companyCode,
		Level:             constant.LevelDistributor,
		EmpID:             visitRow.EmpID,
		EmpCode:           visitRow.EmpCode,
		EmpName:           visitRow.EmpName,
		ActivityTime:      activityTime,
		Planned:           int(plannedCount),
		OnGoing:           int(onGoingCount),
		ExtraCall:         int(extraCallCount),
		Visited:           int(visitedCount),
		Skipped:           int(skippedCount),
		ReturnSummary:     defaultVisitSummary(),
		CollectionSummary: defaultVisitSummary(),
	}

	return visitInfo, nil
}

func defaultVisitSummary() response.VisitSummaryStatus {
	return response.VisitSummaryStatus{
		Count:  0,
		Status: visitSummaryStatusNone,
	}
}

func buildVisitSummary(count int) response.VisitSummaryStatus {
	if count > 0 {
		return response.VisitSummaryStatus{
			Count:  count,
			Status: visitSummaryStatusCompleted,
		}
	}

	return defaultVisitSummary()
}

// transformShipmentRows groups shipment rows by shipment_no
func transformShipmentRows(rows []model.ShipmentRow) []response.ShipmentData {
	shipmentMap := make(map[string]*response.ShipmentData)
	var shipmentOrder []string

	for _, row := range rows {
		if _, exists := shipmentMap[row.ShipmentNo]; !exists {
			shipmentMap[row.ShipmentNo] = &response.ShipmentData{
				ShipmentNo:   row.ShipmentNo,
				Status:       row.Status,
				ShipmentData: []response.ShipmentItem{},
			}
			shipmentOrder = append(shipmentOrder, row.ShipmentNo)
		}

		item := response.ShipmentItem{
			OutletID:   row.OutletID,
			OutletName: row.OutletName,
			OutletCode: row.OutletCode,
			TotalNetto: row.TotalNetto,
		}
		shipmentMap[row.ShipmentNo].ShipmentData = append(shipmentMap[row.ShipmentNo].ShipmentData, item)
	}

	result := make([]response.ShipmentData, 0, len(shipmentOrder))
	for _, no := range shipmentOrder {
		result = append(result, *shipmentMap[no])
	}

	return result
}
