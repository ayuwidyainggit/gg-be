package model

// LiveMonitoringPrincipalRow represents the raw result from principal monitoring query
type LiveMonitoringPrincipalRow struct {
	EmpID              int     `gorm:"column:emp_id"`
	EmpCode            string  `gorm:"column:emp_code"`
	EmpName            string  `gorm:"column:emp_name"`
	DistributorID      int     `gorm:"column:distributor_id"`
	AreaID             int     `gorm:"column:area_id"`
	RegionID           int     `gorm:"column:region_id"`
	PjpID              int     `gorm:"column:pjp_id"`
	PjpCode            int     `gorm:"column:pjp_code"`
	ApprovalStatus     string  `gorm:"column:approval_status"`
	RouteCode          int64   `gorm:"column:route_code"`
	RouteName          string  `gorm:"column:route_name"`
	DestinationID      int     `gorm:"column:destination_id"`
	DestinationCode    string  `gorm:"column:destination_code"`
	DestinationType    string  `gorm:"column:destination_type"`
	DestinationName    string  `gorm:"column:destination_name"`
	DestinationAddress string  `gorm:"column:destination_address"`
	Longitude          float64 `gorm:"column:longitude"`
	Latitude           float64 `gorm:"column:latitude"`
	ArriveAt           *int64  `gorm:"column:arrive_at"`
	LeaveAt            *int64  `gorm:"column:leave_at"`
	ArriveLongitude    float64 `gorm:"column:arrive_longitude"`
	ArriveLatitude     float64 `gorm:"column:arrive_latitude"`
	LeaveLongitude     *string `gorm:"column:leave_longitude"`
	LeaveLatitude      *string `gorm:"column:leave_latitude"`
	FileURL            *string `gorm:"column:file_url"`
	Start              *int64  `gorm:"column:start"`
	Finish             *int64  `gorm:"column:finish"`
	SkipAt             *int64  `gorm:"column:skip_at"`
	SkipReason         *string `gorm:"column:skip_reason"`
	IsExtraCall        bool    `gorm:"column:is_extra_call"`
}

// LiveMonitoringDistributorRow represents the raw result from distributor monitoring query
type LiveMonitoringDistributorRow struct {
	CustID              string  `gorm:"column:cust_id"`
	EmpID               int     `gorm:"column:emp_id"`
	EmpCode             string  `gorm:"column:emp_code"`
	EmpName             string  `gorm:"column:emp_name"`
	SalesmanCode        string  `gorm:"column:salesman_code"`
	DistributorID       int     `gorm:"column:distributor_id"`
	AreaID              int     `gorm:"column:area_id"`
	RegionID            int     `gorm:"column:region_id"`
	AttendanceID        *int64  `gorm:"column:attendance_id"`
	AttendanceLongitude float64 `gorm:"column:attendance_longitude"`
	AttendanceLatitude  float64 `gorm:"column:attendance_latitude"`
	AttendanceAt        *int64  `gorm:"column:attendance_at"`
	CurrentLongitude    float64 `gorm:"column:current_longitude"`
	CurrentLatitude     float64 `gorm:"column:current_latitude"`
	CurrentCoordinateAt *int64  `gorm:"column:current_coordinate_at"`
	CurrentSource       string  `gorm:"column:current_coordinate_source"`
	PjpID               int     `gorm:"column:pjp_id"`
	ApprovalStatus      string  `gorm:"column:approval_status"`
	RouteCode           int64   `gorm:"column:route_code"`
	RouteName           string  `gorm:"column:route_name"`
	OutletID            int     `gorm:"column:outlet_id"`
	OutletCode          string  `gorm:"column:outlet_code"`
	OutletName          string  `gorm:"column:outlet_name"`
	DestinationType     string  `gorm:"column:destination_type"`
	Longitude           float64 `gorm:"column:longitude"`
	Latitude            float64 `gorm:"column:latitude"`
	ArriveAt            *int64  `gorm:"column:arrive_at"`
	LeaveAt             *int64  `gorm:"column:leave_at"`
	ArriveLongitude     float64 `gorm:"column:arrive_longitude"`
	ArriveLatitude      float64 `gorm:"column:arrive_latitude"`
	LeaveLongitude      *string `gorm:"column:leave_longitude"`
	LeaveLatitude       *string `gorm:"column:leave_latitude"`
	FileURL             *string `gorm:"column:file_url"`
	Start               *int64  `gorm:"column:start"`
	Finish              *int64  `gorm:"column:finish"`
	SkipAt              *int64  `gorm:"column:skip_at"`
	SkipReason          *string `gorm:"column:skip_reason"`
	IsExtraCall         bool    `gorm:"column:is_extra_call"`
}

type LatestVisitCoordinateRow struct {
	CustID          string  `gorm:"column:cust_id"`
	EmpCode         string  `gorm:"column:emp_code"`
	OutletCode      string  `gorm:"column:outlet_code"`
	ArriveLongitude float64 `gorm:"column:arrive_longitude"`
	ArriveLatitude  float64 `gorm:"column:arrive_latitude"`
	FileURL         *string `gorm:"column:file_url"`
}

type DistributorEmployeeMetaRow struct {
	EmpID         int    `gorm:"column:emp_id"`
	EmpCode       string `gorm:"column:emp_code"`
	EmpName       string `gorm:"column:emp_name"`
	DistributorID int    `gorm:"column:distributor_id"`
	AreaID        int    `gorm:"column:area_id"`
	RegionID      int    `gorm:"column:region_id"`
}

type DistributorRouteMetaRow struct {
	CustID    string `gorm:"column:cust_id"`
	RouteCode int64  `gorm:"column:route_code"`
	RouteName string `gorm:"column:route_name"`
}

type DistributorOutletMetaRow struct {
	CustID     string `gorm:"column:cust_id"`
	OutletID   int    `gorm:"column:outlet_id"`
	OutletCode string `gorm:"column:outlet_code"`
	OutletName string `gorm:"column:outlet_name"`
}

// VisitInformationRow represents the raw result from visit information query
type VisitInformationRow struct {
	EmpID     int    `gorm:"column:emp_id"`
	EmpCode   string `gorm:"column:emp_code"`
	EmpName   string `gorm:"column:emp_name"`
	Plan      int    `gorm:"column:plan"`
	ExtraCall int    `gorm:"column:extra_call"`
	OnGoing   int    `gorm:"column:on_going"`
	Visited   int    `gorm:"column:visited"`
	TotalSkip int    `gorm:"column:total_skip"`
	Matched   int    `gorm:"column:matched"`
}

// SalesRow represents the raw result from sales query
type SalesRow struct {
	OutletID   int     `gorm:"column:outlet_id"`
	OutletCode string  `gorm:"column:outlet_code"`
	OutletName string  `gorm:"column:outlet_name"`
	SalesOrder float64 `gorm:"column:sales_order"`
}

// ReturnRow represents the raw result from return query
type ReturnRow struct {
	OutletID    int     `gorm:"column:outlet_id"`
	OutletCode  string  `gorm:"column:outlet_code"`
	OutletName  string  `gorm:"column:outlet_name"`
	ReturnTotal float64 `gorm:"column:return_total"`
}

// CollectionRow represents collection data per outlet
type CollectionRow struct {
	OutletID        int     `gorm:"column:outlet_id"`
	OutletCode      string  `gorm:"column:outlet_code"`
	OutletName      string  `gorm:"column:outlet_name"`
	CollectionTotal float64 `gorm:"column:collection_total"`
}

// ExpenseRow represents the raw result from expense query
type ExpenseRow struct {
	ExpenseTypeID   int     `gorm:"column:expense_type_id"`
	ExpenseTypeName string  `gorm:"column:expense_type_name"`
	Note            string  `gorm:"column:note"`
	Amount          float64 `gorm:"column:amount"`
}

// ShipmentRow represents the raw result from shipment query
type ShipmentRow struct {
	ShipmentNo string  `gorm:"column:shipment_no"`
	Status     string  `gorm:"column:status"`
	OutletID   int     `gorm:"column:outlet_id"`
	OutletCode string  `gorm:"column:outlet_code"`
	OutletName string  `gorm:"column:outlet_name"`
	TotalNetto float64 `gorm:"column:total_netto"`
}

// SurveyDataRow represents grouped submitted survey data per outlet
type SurveyDataRow struct {
	Submission  int64  `gorm:"column:submission"`
	SurveyTitle string `gorm:"column:survey_title"`
	OutletCode  string `gorm:"column:outlet_code"`
	OutletName  string `gorm:"column:outlet_name"`
}

// AttendanceRow represents the raw result from attendance query
type AttendanceRow struct {
	EmpID        int     `gorm:"column:emp_id"`
	AttendanceID *int64  `gorm:"column:attendance_id"`
	CreatedAt    string  `gorm:"column:created_at"`
	Timestamp    *int64  `gorm:"column:attendance_at"`
	Longitude    float64 `gorm:"column:attendance_longitude"`
	Latitude     float64 `gorm:"column:attendance_latitude"`
	ClockOutID   *int64  `gorm:"column:clock_out"`
	ClockOutAt   *int64  `gorm:"column:clock_out_at"`
	ClockOutLong float64 `gorm:"column:clock_out_longitude"`
	ClockOutLat  float64 `gorm:"column:clock_out_latitude"`
}

// CurrentCoordinateRow represents resolved current coordinate data for an employee.
type CurrentCoordinateRow struct {
	EmpID        int     `gorm:"column:emp_id"`
	Longitude    float64 `gorm:"column:current_longitude"`
	Latitude     float64 `gorm:"column:current_latitude"`
	Timestamp    *int64  `gorm:"column:current_coordinate_at"`
	Source       string  `gorm:"column:current_coordinate_source"`
	SourceRank   int     `gorm:"column:source_rank"`
	SourceRecord int64   `gorm:"column:source_record"`
}

// DistributorInfoRow represents distributor information
type DistributorInfoRow struct {
	DistributorID   int    `gorm:"column:distributor_id"`
	DistributorCode string `gorm:"column:distributor_code"`
	DistributorName string `gorm:"column:distributor_name"`
}

// UserFullnameRow represents user fullname for principal
type UserFullnameRow struct {
	UserFullname string `gorm:"column:user_fullname"`
}

// UpdateLocationRow represents one chronological location timeline event from union query.
type UpdateLocationRow struct {
	RecordedAt      string  `gorm:"column:recorded_at"`
	Type            string  `gorm:"column:type"`
	Latitude        float64 `gorm:"column:latitude"`
	Longitude       float64 `gorm:"column:longitude"`
	DestinationID   *int64  `gorm:"column:destination_id"`
	DestinationType *string `gorm:"column:destination_type"`
	DestinationName *string `gorm:"column:destination_name"`
}
