package response

type SummaryResponse struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type VisitStatusResponse struct {
	Planned    int64 `json:"planned"`
	Finished   int64 `json:"finished"`
	Skipped    int64 `json:"skipped"`
	OnProgress int64 `json:"on_progress"`
	OnHold     int64 `json:"on_hold"`
	ExtraCall  int64 `json:"extra_call"`
}

type TodoListResponse struct {
	ArriveAt   int64   `json:"arrive_at"`
	LeaveAt    int64   `json:"leave_at"`
	OnHold     int64   `json:"on_hold"`
	ResumeAt   int64   `json:"resume_at"`
	SkipAt     int64   `json:"skip_at"`
	SkipReason *string `json:"skip_reason"`
	InOutlet   bool    `json:"in_outlet"`
}
