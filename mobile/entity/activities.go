package entity

type SummaryDailyRequest struct {
	CustId       string
	ParentCustId string
	EmployeeCode string
	EmployeeId   int64
}

type SummaryDailyResponse struct {
	LastUpdate    string `json:"last_update"`
	Plan          int    `json:"plan"`
	Visit         int    `json:"visit"`
	ExtraCall     int    `json:"extra_call"`
	EffectiveCall int    `json:"effective_call"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	DriveTime     string `json:"drive_time"`
	EstTime       string `json:"est_time"`
}
