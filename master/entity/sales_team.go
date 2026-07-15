package entity

import (
	"time"
)

type SalesTeamResponse struct {
	SalesTeamId   int        `json:"sales_team_id"`
	SalesTeamCode string     `json:"sales_team_code"`
	SalesTeamName string     `json:"sales_team_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type SalesTeamLookupResponse struct {
	SalesTeamId   int    `json:"sales_team_id"`
	SalesTeamCode string `json:"sales_team_code"`
	SalesTeamName string `json:"sales_team_name"`
}

type CreateSalesTeamBody struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	CreatedBy     int64  `json:"created_by" validate:"required"`
	SalesTeamCode string `json:"sales_team_code" validate:"required,max=3,numeric"`
	SalesTeamName string `json:"sales_team_name" validate:"required,max=20,alphanumericSpace"`
	IsActive      bool   `json:"is_active"`
}

type DetailSalesTeamParams struct {
	SalesTeamId int `params:"sales_team_id" validate:"required"`
}

type UpdateSalesTeamParams struct {
	SalesTeamId int `params:"sales_team_id" validate:"required"`
}

type DeleteSalesTeamParams struct {
	SalesTeamId int `params:"sales_team_id" validate:"required"`
}

type UpdateSalesTeamRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	SalesTeamCode string `json:"sales_team_code,omitempty" validate:"required,max=3,omitempty,numeric"`
	SalesTeamName string `json:"sales_team_name,omitempty" validate:"max=20,omitempty,alphanumericSpace"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
