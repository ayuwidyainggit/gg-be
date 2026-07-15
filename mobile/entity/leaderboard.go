package entity

type LeaderboardsRequest struct {
}

type LeaderboardsResponse struct {
	Rank     int     `json:"rank"`
	Name     string  `json:"name"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
	Image    string  `json:"image"`
	IsMe     bool    `json:"is_me"`
}
