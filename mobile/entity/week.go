package entity

import "time"

type WeekListQueryFilter struct {
	Page          int    `query:"page" validate:"required"`
	Limit         int    `query:"limit" validate:"required"`
	Sort          string `query:"sort" validate:"required"`
	Month         *int   `query:"month"`
	Year          *int   `query:"year"`
	IsActive      *int   `query:"is_active"` // null = semua, 1 = active, 0 = nonaktif
	CustID        string `query:"-" json:"-"`
	IsDistributor bool   `query:"-" json:"-"` // true if distributor user, false if principal
	EmpID         int64  `query:"-" json:"-"`
	WDCID         int    `query:"wdc_id" validate:"required"`
	ParentCustID  string `query:"parent_cust_id" validate:"required"`
}

type WorkDayData struct {
	PerYear             int        `json:"per_year"`
	PerId               int        `json:"per_id"`
	WeekId              int        `json:"week_id"`
	WorkDate            time.Time  `json:"work_date"`
	IsWork              bool       `json:"is_work"`
	IsActive            bool       `json:"is_active"`
	IsClosed            bool       `json:"is_closed"`
	ClosedAt            *time.Time `json:"closed_at"`
	ClosedBy            int        `json:"closed_by"`
	ClosedByName        string     `json:"closed_by_name"`
	RouteCode           int64      `json:"route_code"`
	RouteName           string     `json:"route_name"`       // From BE: Route 1, Route 2, etc
	NumberOfOutlet      int        `json:"number_of_outlet"` // Count from pjp.outlet_visit_list
	NumberOfDistributor int        `json:"number_of_distributor"`
}

type WeekListResponse struct {
	PerYear      int           `json:"per_year"`
	PerId        int           `json:"per_id"`
	WeekId       int           `json:"week_id"`
	WeekStart    time.Time     `json:"week_start"`
	WeekEnd      time.Time     `json:"week_end"`
	IsActive     bool          `json:"is_active"`
	IsClosed     bool          `json:"is_closed"`
	ClosedAt     *time.Time    `json:"closed_at"`
	ClosedBy     int           `json:"closed_by"`
	ClosedByName string        `json:"closed_by_name"`
	WorkDays     []WorkDayData `json:"work_days"`
}
