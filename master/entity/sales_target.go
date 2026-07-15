package entity

import (
	"time"
)

const (
	StatusDraft    = 0
	StatusActive   = 1
	StatusInactive = 2
)

// Sales Target List Query Filter
type SalesTargetQueryFilter struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit" validate:"required"`
	Sort  string `query:"sort"`
	Year  int    `query:"year"`
}

// Sales Target List Response
type SalesTargetListResponse struct {
	SalesTargetId  int64      `json:"sales_target_id"`
	Year           int        `json:"year"`
	Month          int        `json:"month"`
	AllocatedTotal int64      `json:"allocated_total"`
	MonthlyTarget  int64      `json:"monthly_target"`
	Remaining      int64      `json:"remaining"`
	Status         string     `json:"status"`
	UpdatedBy      string     `json:"updated_by"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

// Sales Target Detail Response
type SalesTargetDetailResponse struct {
	SalesTargetId  int64                      `json:"sales_target_id"`
	Year           int                        `json:"year"`
	Month          int                        `json:"month"`
	AllocatedTotal int64                      `json:"allocated_total"`
	MonthlyTarget  int64                      `json:"monthly_target"`
	Remaining      int64                      `json:"remaining"`
	Status         string                     `json:"status"`
	UpdatedBy      string                     `json:"updated_by"`
	UpdatedAt      *time.Time                 `json:"updated_at"`
	Details        []SalesAllocatedDetailResp `json:"details"`
}

// Sales Allocated Detail Response
type SalesAllocatedDetailResp struct {
	SalesAllocatedId int64  `json:"sales_allocated_id"`
	SalesTargetId    int64  `json:"sales_target_id"`
	SalesmanId       int64  `json:"salesman_id"`
	SalesName        string `json:"sales_name"`
	OprType          string `json:"opr_type"`
	DistributorId    *int64 `json:"distributor_id"`
	DistributorCode  string `json:"distributor_code"`
	DistributorName  string `json:"distributor_name"`
	ChannelId        *int64 `json:"channel_id"`
	ChannelCode      string `json:"channel_code"`
	ChannelName      string `json:"channel_name"`
	SalesTeamId      *int64 `json:"sales_team_id"`
	SalesTeamCode    string `json:"sales_team_code"`
	SalesTeamName    string `json:"sales_team_name"`
	Allocated        int64  `json:"allocated"`
	IsActive         bool   `json:"is_active"`
}

// Sales Target Monthly Distributor Response
type SalesTargetMonthlyDistResp struct {
	DistributorId int64 `json:"distributor_id"`
	Year          int   `json:"year"`
	Month         int   `json:"month"`
	MonthlyTarget int64 `json:"monthly_target"`
}

// Sales Target Monthly Distributor Query
type SalesTargetMonthlyDistQuery struct {
	Year          int   `query:"year" validate:"required"`
	Month         int   `query:"month" validate:"required"`
	DistributorId int64 `query:"distributor_id" validate:"required"`
}

// Create Sales Target Request
type CreateSalesTargetRequest struct {
	CustId                          string                      `json:"-"`
	ParentCustId                    string                      `json:"-"`
	CreatedBy                       int64                       `json:"-"`
	SalesTargetDistributorYearlyId  int64                       `json:"sales_target_distributor_yearly_id" validate:"required"`
	SalesTargetDistributorMonthlyId int64                       `json:"sales_target_distributor_monthly_id" validate:"required"`
	Month                           int                         `json:"month" validate:"required"`
	Year                            int                         `json:"year" validate:"required"`
	AllocatedTotal                  int64                       `json:"allocated_total" validate:"required"`
	MonthlyTarget                   int64                       `json:"monthly_target" validate:"required"`
	Remaining                       *int64                      `json:"remaining" validate:"required,gte=0"`
	Status                          *int                        `json:"status" validate:"required,oneof=0 1 2"`
	Data                            []SalesAllocatedItemRequest `json:"data" validate:"required,dive"`
}

// Sales Allocated Item Request
type SalesAllocatedItemRequest struct {
	SalesmanId  int64 `json:"salesman_id" validate:"required"`
	SalesTeamId int64 `json:"sales_team_id" validate:"required"`
	Allocated   int64 `json:"allocated" validate:"required"`
}

// Update Sales Target Request
type UpdateSalesTargetRequest struct {
	CustId                          string                      `json:"-"`
	ParentCustId                    string                      `json:"-"`
	UpdatedBy                       int64                       `json:"-"`
	SalesTargetDistributorYearlyId  *int64                      `json:"sales_target_distributor_yearly_id,omitempty"`
	SalesTargetDistributorMonthlyId *int64                      `json:"sales_target_distributor_monthly_id,omitempty"`
	Month                           *int                        `json:"month,omitempty"`
	Year                            *int                        `json:"year,omitempty"`
	AllocatedTotal                  *int64                      `json:"allocated_total,omitempty"`
	MonthlyTarget                   *int64                      `json:"monthly_target,omitempty"`
	Remaining                       *int64                      `json:"remaining,omitempty"`
	Status                          *int                        `json:"status,omitempty" validate:"omitempty,oneof=0 1 2"`
	Data                            []SalesAllocatedItemRequest `json:"data,omitempty" validate:"omitempty,dive"`
}

// Sales Target Params
type SalesTargetParams struct {
	SalesTargetId int64 `params:"sales_target_id" validate:"required"`
}
