package entity

type AnnouncementsRequest struct {
}

type AnnouncementsResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"start_time"`
	IsNew     bool   `json:"is_new"`
	Image     string `json:"image"`
}
