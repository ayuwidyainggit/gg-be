package entity

import "time"

type MPeriodsQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	PerYear  string `query:"per_year"`
	IsActive *int   `query:"is_active"`
}

type MPeriodsResponse struct {
	PerYear       int        `json:"per_year"`
	PerId         int        `json:"per_id"`
	WeekCount     int        `json:"week_count"`
	IsActive      bool       `json:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
	IsClosed      bool       `json:"is_closed"`
	ClosedAt      *time.Time `json:"closed_at"`
	ClosedByName  string     `json:"closed_by_name"`
}

type CreateMPeriodsBody struct {
	CustId       string `json:"cust_id"`
	ParentCustId string `json:"parent_cust_id"`
	PerYear      int    `json:"per_year"`
	PerId        int    `json:"per_id"`
	WeekCount    int    `json:"week_count"`
	IsActive     bool   `json:"is_active"`
	UpdatedBy    int64  `json:"updated_by"`
}

type DetailCreateMPeriodsParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required"`
}

type UpdateMPeriodsParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required,min=1,max=20"`
}

type DeleteMPeriodsParams struct {
	PerYear int `params:"per_year" validate:"required"`
	PerId   int `params:"per_id" validate:"required"`
}

type UpdateMPeriodsRequest struct {
	ParentCustId string
	CustId       string `json:"cust_id"`
	WeekCount    int    `json:"week_count" validate:"max=6"`
	IsActive     *bool  `json:"is_active"`
	UpdatedBy    int64  `json:"updated_by"`
}

type MPeriodsListYear struct {
	PerYear int `json:"per_year"`
}
