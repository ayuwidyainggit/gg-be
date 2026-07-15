package model

import (
	"time"
)

type ProductAssignment struct {
	ID              int64     `json:"id" db:"id"`
	CustID          string    `json:"cust_id" db:"cust_id"`
	ProID           int64     `json:"pro_id" db:"pro_id"`
	ProCode         string    `json:"pro_code" db:"pro_code"`
	ProName         string    `json:"pro_name" db:"pro_name"`
	DistributorID   int64     `json:"distributor_id" db:"distributor_id"`
	DistributorCode string    `json:"distributor_code" db:"distributor_code"`
	DistributorName string    `json:"distributor_name" db:"distributor_name"`
	CreatedBy       int64     `json:"created_by" db:"created_by"`
	CreatedByName   string    `json:"created_by_name" db:"created_by_name"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy       int64     `json:"updated_by" db:"updated_by"`
}

type ProductAssignmentHistory struct {
	ID              int64     `db:"id"`
	CustID          string    `db:"cust_id"`
	ActionDate      time.Time `db:"action_date"`
	ProID           int64     `db:"pro_id"`
	ProCode         string    `db:"pro_code"`
	ProName         string    `db:"pro_name"`
	DistributorID   int64     `db:"distributor_id"`
	DistributorCode string    `db:"distributor_code"`
	DistributorName string    `db:"distributor_name"`
	AssignmentType  string    `db:"assignment_type"`
	CreatedBy       int64     `db:"created_by"`
	CreatedByName   string    `db:"created_by_name"`
	CreatedAt       time.Time `db:"created_at"`
}

type ProductAssignmentMutation struct {
	CustID         string    `db:"cust_id"`
	ActionDate     time.Time `db:"action_date"`
	ProID          int64     `db:"pro_id"`
	DistributorID  int64     `db:"distributor_id"`
	AssignmentType string    `db:"assignment_type"`
	CreatedBy      int64     `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
}
