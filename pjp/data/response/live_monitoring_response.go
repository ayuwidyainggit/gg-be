package response

// LiveMonitoringResponse is the response structure for location monitoring (Principal & Distributor)
type LiveMonitoringResponse struct {
	Message   string               `json:"message"`
	Data      []LiveMonitoringData `json:"data"`
	Paging    LiveMonitoringPaging `json:"paging"`
	RequestID string               `json:"request_id"`
}

// LiveMonitoringData represents employee data with their PJP information
type LiveMonitoringData struct {
	EmpID                   int                     `json:"emp_id"`
	EmpCode                 string                  `json:"emp_code"`
	EmpName                 string                  `json:"emp_name"`
	DistributorID           int                     `json:"distributor_id"`
	AreaID                  int                     `json:"area_id"`
	RegionID                int                     `json:"region_id"`
	AttendanceID            *int64                  `json:"attendance_id"`
	AttendanceLongitude     float64                 `json:"attendance_longitude"`
	AttendanceLatitude      float64                 `json:"attendance_latitude"`
	AttendanceAt            *int64                  `json:"attendance_at"`
	ClockOut                *int64                  `json:"clock_out"`
	ClockOutLongitude       float64                 `json:"clock_out_longitude"`
	ClockOutLatitude        float64                 `json:"clock_out_latitude"`
	ClockOutAt              *int64                  `json:"clock_out_at"`
	CurrentLongitude        float64                 `json:"current_longitude"`
	CurrentLatitude         float64                 `json:"current_latitude"`
	CurrentCoordinateAt     *int64                  `json:"current_coordinate_at"`
	CurrentCoordinateSource string                  `json:"current_coordinate_source"`
	PjpData                 []LiveMonitoringPjpData `json:"pjp_data"`
}

// LiveMonitoringPjpData represents PJP data with route information
type LiveMonitoringPjpData struct {
	PjpID          int                       `json:"pjp_id"`
	PjpCode        *int                      `json:"pjp_code,omitempty"`
	ApprovalStatus string                    `json:"approval_status"`
	RouteData      []LiveMonitoringRouteData `json:"route_data"`
	ExtraCallData  []LiveMonitoringRouteData `json:"extra_call_data"`
}

// LiveMonitoringRouteData represents route data with destination information
type LiveMonitoringRouteData struct {
	RouteCode       string                          `json:"route_code"`
	RouteName       string                          `json:"route_name"`
	DestinationData []LiveMonitoringDestinationData `json:"destination_data"`
}

// LiveMonitoringDestinationData represents destination/outlet data for Principal
type LiveMonitoringDestinationData struct {
	DestinationID      int     `json:"destination_id,omitempty"`
	DestinationCode    string  `json:"destination_code,omitempty"`
	DestinationType    string  `json:"destination_type,omitempty"`
	DestinationName    string  `json:"destination_name,omitempty"`
	DestinationAddress string  `json:"destination_address,omitempty"`
	Longitude          float64 `json:"longitude"`
	Latitude           float64 `json:"latitude"`
	ArriveAt           *int64  `json:"arrive_at"`
	LeaveAt            *int64  `json:"leave_at"`
	ArriveLongitude    float64 `json:"arrive_longitude"`
	ArriveLatitude     float64 `json:"arrive_latitude"`
	LeaveLongitude     *string `json:"leave_longitude"`
	LeaveLatitude      *string `json:"leave_latitude"`
	FileURL            *string `json:"file_url"`
	Start              *int64  `json:"start"`
	Finish             *int64  `json:"finish"`
	SkipAt             *int64  `json:"skip_at"`
	SkipReason         *string `json:"skip_reason,omitempty"`
}

// LiveMonitoringPaging represents pagination information
type LiveMonitoringPaging struct {
	TotalRecord int `json:"total_record"`
	PageCurrent int `json:"page_current"`
	PageLimit   int `json:"page_limit"`
	PageTotal   int `json:"page_total"`
}

// ============== Detail Response Structures ==============

// LiveMonitoringDetailResponse is the response structure for location monitoring detail
type LiveMonitoringDetailResponse struct {
	Message   string                     `json:"message"`
	Data      []LiveMonitoringDetailData `json:"data"`
	RequestID string                     `json:"request_id"`
}

// UpdateLocationsResponse contains one employee's daily location timeline.
type UpdateLocationsResponse struct {
	Timeline []TimelineItem `json:"timeline"`
}

// TimelineItem represents one chronological employee location event.
type TimelineItem struct {
	Sequence        int     `json:"sequence"`
	Type            string  `json:"type"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	DestinationID   *int64  `json:"destination_id"`
	DestinationType *string `json:"destination_type"`
	DestinationName *string `json:"destination_name"`
	RecordedAt      string  `json:"recorded_at"`
}

// LiveMonitoringDetailData represents the complete detail data
type LiveMonitoringDetailData struct {
	VisitInformation VisitInformationData `json:"visit_information"`
	Sales            []SalesData          `json:"sales"`
	Return           []ReturnData         `json:"return"`
	Collection       []CollectionData     `json:"collection"`
	Expense          []ExpenseData        `json:"expense"`
	Shipment         []ShipmentData       `json:"shipment"`
	SurveyData       []SurveyData         `json:"survey_data"`
}

// VisitInformationData represents visit information for an employee
type VisitInformationData struct {
	ActivityDate      string             `json:"activity_date"`
	CompanyName       string             `json:"company_name"`
	CompanyCode       string             `json:"company_code"`
	Level             string             `json:"level"` // "Principal" or "Distributor"
	EmpID             int                `json:"emp_id"`
	EmpCode           string             `json:"emp_code"`
	EmpName           string             `json:"emp_name"`
	ActivityTime      *string            `json:"activity_time"`
	Planned           int                `json:"planned"`
	OnGoing           int                `json:"on_going"`
	ExtraCall         int                `json:"extra_call"`
	Visited           int                `json:"visited"`
	Skipped           int                `json:"skipped"`
	ReturnSummary     VisitSummaryStatus `json:"return_summary"`
	CollectionSummary VisitSummaryStatus `json:"collection_summary"`
}

// VisitSummaryStatus represents additive count and completion status details.
type VisitSummaryStatus struct {
	Count  int    `json:"count"`
	Status string `json:"status"`
}

// SalesData represents sales order data per outlet
type SalesData struct {
	OutletID   int     `json:"outlet_id"`
	OutletCode string  `json:"outlet_code"`
	OutletName string  `json:"outlet_name"`
	SalesOrder float64 `json:"sales_order"`
}

// ReturnData represents return data per outlet
type ReturnData struct {
	OutletID    int     `json:"outlet_id"`
	OutletCode  string  `json:"outlet_code"`
	OutletName  string  `json:"outlet_name"`
	ReturnTotal float64 `json:"return_total"`
}

// CollectionData represents collection data per outlet (currently empty - waiting for dev sby)
type CollectionData struct {
	OutletID        *int     `json:"outlet_id"`
	OutletCode      *string  `json:"outlet_code"`
	OutletName      *string  `json:"outlet_name"`
	CollectionTotal *float64 `json:"collection_total"`
}

// ExpenseData represents expense data
type ExpenseData struct {
	ExpenseTypeID int     `json:"expense_type_id"`
	ExpenseType   string  `json:"expense_type"`
	Note          string  `json:"note"`
	ExpenseTotal  float64 `json:"expense_total"`
}

// ShipmentData represents shipment data with outlets
type ShipmentData struct {
	ShipmentNo   string         `json:"shipment_no"`
	Status       string         `json:"status"`
	ShipmentData []ShipmentItem `json:"shipment_data"`
}

// ShipmentItem represents individual shipment item per outlet
type ShipmentItem struct {
	OutletID   int     `json:"outlet_id"`
	OutletName string  `json:"outlet_name"`
	OutletCode string  `json:"outlet_code"`
	TotalNetto float64 `json:"total_netto"`
}

// SurveyData represents survey submission aggregation per survey and outlet
type SurveyData struct {
	Submission  int64  `json:"submission"`
	SurveyTitle string `json:"survey_title"`
	OutletCode  string `json:"outlet_code"`
	OutletName  string `json:"outlet_name"`
}
