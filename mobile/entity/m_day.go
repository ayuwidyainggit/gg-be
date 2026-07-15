package entity

type MDayListResponse struct {
	DayId   int64  `json:"day_id"`
	DayName string `json:"day_name"`
	LangId  string `json:"lang_id"`
}

type DetailMDayBodyParam struct {
	DayId int64 `params:"day_id"`
}
