package model

import "time"

type MPeriods struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	PerYear       int        `json:"per_year" db:"per_year"`
	PerId         int        `json:"per_id" db:"per_id"`
	WeekCount     *int       `json:"week_count" db:"week_count"`
	IsActive      *bool      `json:"is_active" db:"is_active"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	IsClosed      bool       `json:"is_closed" db:"is_closed"`
	ClosedAt      *time.Time `json:"closed_at" db:"closed_at"`
	ClosedBy      *int64     `json:"closed_by" db:"closed_by"`
	ClosedByName  *string    `json:"closed_by_name" db:"closed_by_name"`
}

type MPeriodsUpdate struct {
	PerYear   *int       `json:"per_year" sql:"per_year"`
	PerId     *int       `json:"per_id" sql:"per_id"`
	WeekCount *int       `json:"week_count" sql:"week_count"`
	IsActive  *bool      `json:"is_active" sql:"is_active"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
}
