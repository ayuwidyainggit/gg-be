package entity

import "time"

type MWorkingDayQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	PerYear  string `query:"per_year"`
	IsActive *int   `query:"is_active"`
}
type MWorkingDayResponse struct {
	PerYear      int        `json:"per_year"`
	PerId        int        `json:"per_id"`
	WeekId       int        `json:"week_id"`
	WorkDate     string     `json:"work_date"`
	IsWork       bool       `json:"is_work"`
	IsActive     bool       `json:"is_active"`
	IsClosed     bool       `json:"is_closed"`
	ClosedAt     *time.Time `json:"closed_at"`
	ClosedBy     int64      `json:"closed_by"`
	ClosedByName string     `json:"closed_by_name"`
}

type CreateMWorkingDayBody struct {
	CustId   string  `json:"cust_id"`
	PerYear  int     `json:"per_year"`
	PerId    int     `json:"per_id"`
	WeekId   int     `json:"week_id"`
	WorkDate *string `json:"work_date"`
	IsWork   bool    `json:"is_work"`
	IsActive bool    `json:"is_active"`
}

type DetailCreateMWorkingDayParams struct {
	PerYear  int    `params:"per_year" validate:"required"`
	PerId    int    `params:"per_id" validate:"required"`
	WeekId   int    `params:"week_id" validate:"required"`
	WorkDate string `params:"work_date" validate:"required"`
}

type UpdateMWorkingDayParams struct {
	PerYear  int    `params:"per_year" validate:"required"`
	PerId    int    `params:"per_id" validate:"required"`
	WeekId   int    `params:"week_id" validate:"required"`
	WorkDate string `params:"work_date" validate:"required"`
}

type DeleteMWorkingDayParams struct {
	PerYear  int    `params:"per_year" validate:"required"`
	PerId    int    `params:"per_id" validate:"required"`
	WeekId   int    `params:"week_id" validate:"required"`
	WorkDate string `params:"work_date" validate:"required"`
}

type UpdateMWorkingDayRequest struct {
	CustId   string  `json:"cust_id"`
	WorkDate *string `json:"work_date"`
	IsWork   bool    `json:"is_work"`
	IsActive bool    `json:"is_active"`
}

type MWorkingDayActiveResponse struct {
	PerYear  int     `json:"per_year"`
	PerId    int     `json:"per_id"`
	WeekId   int     `json:"week_id"`
	WorkDate *string `json:"work_date"`
	IsWork   bool    `json:"is_work"`
	IsActive bool    `json:"is_active"`
	IsClosed bool    `json:"is_closed"`
}
