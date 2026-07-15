package entity

type EventsRequest struct {
}

type EventsResponse struct {
	Id        string   `json:"id"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Location  EventDet `json:"location"`
}

type EventDet struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
}
