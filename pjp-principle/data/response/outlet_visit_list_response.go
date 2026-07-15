package response

import "time"

type OutletVisitListResponse struct {
	ID                 int       `json:"id"`
	Year               int       `json:"year"`
	Week               int       `json:"week"`
	Date               time.Time `json:"date"`
	Day                string    `json:"day"`
	RouteCode          *int      `json:"route_code"`
	DestinationID      int       `json:"outlet_id"`
	DestinationCode    string    `json:"outlet_code"`
	DestinationName    string    `json:"outlet_name"`
	DestinationAddress string    `json:"outlet_address"`
	DueDate            string    `json:"due_date"`
	PjpID              *int      `json:"pjp_id"`
	PjpCode            *int      `json:"pjp_code"`
	Start              *int64    `json:"start"`
	Finish             *int64    `json:"finish"`
	SkipAt             *int64    `json:"skip_at"`
	LeaveAt            *int64    `json:"leave_at"`
	ArriveAt           *int64    `json:"arrive_at"`
	OnHold             *int64    `json:"on_hold"`
	ResumeAt           *int64    `json:"resume_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Status             string    `json:"status"`
	IsPlanned          bool      `json:"is_planned"`
}
