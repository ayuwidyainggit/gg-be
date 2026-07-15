package entity

type TravelListResponse struct {
	ArriveAt int64 `json:"arrive_at"`
	UnloadAt int64 `json:"unload_at"`
	LeaveAt  int64 `json:"leave_at"`
}
