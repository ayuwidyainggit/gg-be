package entity

type SalesmanReport struct {
	TotalOutlets    int              `json:"total_outlets"`
	TotalVisit      int              `json:"total_visit"`
	TotalNotVisit   int              `json:"total_not_visit"`
	Visit           VisitDetails     `json:"visit"`
	NotVisit        NotVisitDetails  `json:"not_visit"`
	NotVisitReasons []NotVisitReason `json:"not_visit_reasons"`
}

type VisitDetails struct {
	Planned    int `json:"planned"`     // Total planned visits
	NotPlanned int `json:"not_planned"` // Total not-planned visits
}

type NotVisitDetails struct {
	Planned    int `json:"planned"`     // Total planned but not visited
	NotPlanned int `json:"not_planned"` // Total not-planned and not visited
}

type NotVisitReason struct {
	Reason     string `json:"reason"`
	Count      int    `json:"count"`
	Percentage int    `json:"percentage"`
}

type SalesmanReportQueryFilter struct {
	SalesmanId string `validate:"required" json:"salesman_id"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Date       string `json:"date"`
}
