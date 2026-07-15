package response

import "time"

type OutletVisitListResponse struct {
	ID              int       `json:"id"`
	Year            int       `json:"year"`
	Week            int       `json:"week"`
	Date            time.Time `json:"date"`
	Day             string    `json:"day"`
	RouteCode       *int      `json:"route_code"`
	OutletID        int       `json:"outlet_id"`
	OutletCode      string    `json:"outlet_code"`
	OutletName      string    `json:"outlet_name"`
	OutletAddress   string    `json:"outlet_address"`
	OutletLongitude string    `json:"outlet_longitude"`
	OutletLatitude  string    `json:"outlet_latitude"`
	DueDate         string    `json:"due_date"`
	PjpID           *int      `json:"pjp_id"`
	PjpCode         *int      `json:"pjp_code"`
	Start           *int64    `json:"start"`
	Finish          *int64    `json:"finish"`
	SkipAt          *int64    `json:"skip_at"`
	LeaveAt         *int64    `json:"leave_at"`
	ArriveAt        *int64    `json:"arrive_at"`
	OnHold          *int64    `json:"on_hold"`
	ResumeAt        *int64    `json:"resume_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Status          string    `json:"status"`
	IsPlanned       bool      `json:"is_planned"`
	DestinationType *string   `json:"destination_type"`
	SkipReason      *string   `json:"skip_reason"`
	InOutlet        bool      `json:"in_outlet"`
}
