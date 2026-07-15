package model

import (
	"time"
)

// SalesTarget database model
type SalesTarget struct {
	CustId                          string     `db:"cust_id" json:"cust_id"`
	SalesTargetId                   int64      `db:"sales_target_id" json:"sales_target_id"`
	SalesTargetDistributorYearlyId  int64      `db:"sales_target_distributor_yearly_id" json:"sales_target_distributor_yearly_id"`
	SalesTargetDistributorMonthlyId int64      `db:"sales_target_distributor_monthly_id" json:"sales_target_distributor_monthly_id"`
	Month                           int        `db:"month" json:"month"`
	Year                            int        `db:"year" json:"year"`
	AllocatedTotal                  int64      `db:"allocated_total" json:"allocated_total"`
	MonthlyTarget                   int64      `db:"monthly_target" json:"monthly_target"`
	Remaining                       int64      `db:"remaining" json:"remaining"`
	Status                          int        `db:"status" json:"status"`
	CreatedBy                       int64      `db:"created_by" json:"created_by"`
	CreatedAt                       time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy                       *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt                       *time.Time `db:"updated_at" json:"updated_at"`
	DeletedBy                       *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt                       *time.Time `db:"deleted_at" json:"deleted_at"`
	IsDel                           bool       `db:"is_del" json:"is_del"`
}

// SalesTargetList for list query with joins
type SalesTargetList struct {
	SalesTargetId  int64      `db:"sales_target_id" json:"sales_target_id"`
	Month          int        `db:"month" json:"month"`
	Year           int        `db:"year" json:"year"`
	AllocatedTotal int64      `db:"allocated_total" json:"allocated_total"`
	MonthlyTarget  int64      `db:"monthly_target" json:"monthly_target"`
	Remaining      int64      `db:"remaining" json:"remaining"`
	Status         int        `db:"status" json:"status"`
	UpdatedBy      *string    `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy      *string    `db:"created_by_name" json:"created_by_name"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`
	CreatedAtRaw   time.Time  `db:"created_at_raw" json:"-"`
}

// SalesAllocated database model
type SalesAllocated struct {
	CustId           string     `db:"cust_id" json:"cust_id"`
	SalesAllocatedId int64      `db:"sales_allocated_id" json:"sales_allocated_id"`
	SalesTargetId    int64      `db:"sales_target_id" json:"sales_target_id"`
	SalesmanId       int64      `db:"salesman_id" json:"salesman_id"`
	SalesTeamId      *int64     `db:"sales_team_id" json:"sales_team_id"`
	Allocated        int64      `db:"allocated" json:"allocated"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	CreatedBy        int64      `db:"created_by" json:"created_by"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy        *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updated_at"`
	DeletedBy        *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deleted_at"`
	IsDel            bool       `db:"is_del" json:"is_del"`
}

// SalesAllocatedDetail for detail query with joins
type SalesAllocatedDetail struct {
	SalesAllocatedId int64   `db:"sales_allocated_id" json:"sales_allocated_id"`
	SalesTargetId    int64   `db:"sales_target_id" json:"sales_target_id"`
	SalesmanId       int64   `db:"salesman_id" json:"salesman_id"`
	SalesName        string  `db:"sales_name" json:"sales_name"`
	OprType          string  `db:"opr_type" json:"opr_type"`
	DistributorId    *int64  `db:"distributor_id" json:"distributor_id"`
	DistributorCode  *string `db:"distributor_code" json:"distributor_code"`
	DistributorName  *string `db:"distributor_name" json:"distributor_name"`
	ChannelId        *int64  `db:"channel_id" json:"channel_id"`
	ChannelCode      *string `db:"channel_code" json:"channel_code"`
	ChannelName      *string `db:"channel_name" json:"channel_name"`
	SalesTeamId      *int64  `db:"sales_team_id" json:"sales_team_id"`
	SalesTeamCode    *string `db:"sales_team_code" json:"sales_team_code"`
	SalesTeamName    *string `db:"sales_team_name" json:"sales_team_name"`
	Allocated        int64   `db:"allocated" json:"allocated"`
	IsActive         bool    `db:"is_active" json:"is_active"`
}

// SalesTargetMonthlyDist for monthly distributor query
type SalesTargetMonthlyDist struct {
	DistributorId int64 `db:"distributor_id" json:"distributor_id"`
	Year          int   `db:"year" json:"year"`
	Month         int   `db:"month" json:"month"`
	MonthlyTarget int64 `db:"monthly_target" json:"monthly_target"`
}
