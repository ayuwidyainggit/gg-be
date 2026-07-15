package entity

import "time"

type MWeekQueryFilter struct {
	Page                 int    `query:"page"`
	Limit                int    `query:"limit" validate:"required"`
	Query                string `query:"q"`
	Mode                 string `query:"mode"`
	Sort                 string `query:"sort"`
	PerYear              string `query:"per_year"`
	ParentCustId         string `query:"-"`
	IsActive             *int   `query:"is_active"`
	WorkingDayCalendarID []int  `query:"-"`
}
type MWeekResponse struct {
	CustId       string     `json:"cust_id"`
	PerYear      int        `json:"per_year"`
	PerId        int        `json:"per_id"`
	WeekId       int        `json:"week_id"`
	WeekStart    string     `json:"week_start"`
	WeekEnd      string     `json:"week_end"`
	IsActive     bool       `json:"is_active"`
	IsClosed     bool       `json:"is_closed"`
	ClosedAt     *time.Time `json:"closed_at"`
	ClosedBy     int64      `json:"closed_by"`
	ClosedByName string     `json:"closed_by_name"`
}

type CreateMWeekBody struct {
	CustId    string  `json:"cust_id"`
	PerYear   int     `json:"per_year"`
	PerId     int     `json:"per_id"`
	WeekId    int     `json:"week_id"`
	WeekStart *string `json:"week_start"`
	WeekEnd   *string `json:"week_end"`
	IsActive  bool    `json:"is_active"`
}

type DetailCreateMWeekParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required"`
	WeekId  int `params:"week_id" validate:"required"`
}

type UpdateMWeekParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required"`
	WeekId  int `params:"week_id" validate:"required"`
}

type DeleteMWeekParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required"`
	WeekId  int `params:"week_id" validate:"required"`
}

type UpdateMWeekRequest struct {
	CustId    string  `json:"cust_id"`
	WeekStart *string `json:"week_start"`
	WeekEnd   *string `json:"week_end"`
	IsActive  bool    `json:"is_active"`
}
