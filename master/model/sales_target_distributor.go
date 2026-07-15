package model

import "time"

type SalesTargetDistributorYearly struct {
	CustId                         string     `db:"cust_id" json:"cust_id"`
	SalesTargetDistributorYearlyId int        `db:"sales_target_distributor_yearly_id" json:"sales_target_distributor_yearly_id"`
	AreaId                         int        `db:"area_id" json:"area_id"`
	RegionId                       int        `db:"region_id" json:"region_id"`
	DistributorId                  int        `db:"distributor_id" json:"distributor_id"`
	Year                           int        `db:"year" json:"year"`
	YearlyTarget                   int        `db:"yearly_target" json:"yearly_target"`
	Status                         int        `db:"status" json:"status"`
	IsActive                       bool       `db:"is_active" json:"is_active"`
	UserInactive                   *int64     `db:"user_inactive" json:"user_inactive"`
	InactiveAt                     *time.Time `db:"inactive_at" json:"inactive_at"`
	CreatedBy                      int64      `db:"created_by" json:"created_by"`
	CreatedAt                      time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy                      *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt                      *time.Time `db:"updated_at" json:"updated_at"`
	DeletedBy                      *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt                      *time.Time `db:"deleted_at" json:"deleted_at"`
	IsDel                          bool       `db:"is_del" json:"is_del"`

	// Joins for list/detail
	DistributorCode string `db:"distributor_code" json:"distributor_code"`
	DistributorName string `db:"distributor_name" json:"distributor_name"`
	AreaCode        string `db:"area_code" json:"area_code"`
	AreaName        string `db:"area_name" json:"area_name"`
	RegionCode      string `db:"region_code" json:"region_code"`
	RegionName      string `db:"region_name" json:"region_name"`
	UpdatedByName   string `db:"updated_by_name" json:"updated_by_name"`
}

type SalesTargetDistributorMonthly struct {
	CustId                          string     `db:"cust_id" json:"cust_id"`
	SalesTargetDistributorMonthlyId int        `db:"sales_target_distributor_monthly_id" json:"sales_target_distributor_monthly_id"`
	SalesTargetDistributorYearlyId  int        `db:"sales_target_distributor_yearly_id" json:"sales_target_distributor_yearly_id"`
	Month                           int        `db:"month" json:"month"`
	MonthlyTarget                   int        `db:"monthly_target" json:"monthly_target"`
	IsActive                        bool       `db:"is_active" json:"is_active"`
	CreatedBy                       int64      `db:"created_by" json:"created_by"`
	CreatedAt                       time.Time  `db:"created_at" json:"created_at"`
	UpdatedBy                       *int64     `db:"updated_by" json:"updated_by"`
	UpdatedAt                       *time.Time `db:"updated_at" json:"updated_at"`
	DeletedBy                       *int64     `db:"deleted_by" json:"deleted_by"`
	DeletedAt                       *time.Time `db:"deleted_at" json:"deleted_at"`
	IsDel                           bool       `db:"is_del" json:"is_del"`
}

type SalesTargetDistributorYearlyUpdate struct {
	Year          *int       `json:"year" sql:"year"`
	AreaId        *int       `json:"area_id" sql:"area_id"`
	RegionId      *int       `json:"region_id" sql:"region_id"`
	DistributorId *int       `json:"distributor_id" sql:"distributor_id"`
	YearlyTarget  *int       `json:"yearly_target" sql:"yearly_target"`
	Status        *int       `json:"status" sql:"status"`
	IsActive      *bool      `json:"is_active" sql:"is_active"`
	UserInactive  *int64     `json:"user_inactive" sql:"user_inactive"`
	InactiveAt    *time.Time `json:"inactive_at" sql:"inactive_at"`
	UpdatedBy     *int64     `json:"updated_by" sql:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`
}

// SalesTargetAllocationSummary represents aggregated allocation information for a yearly distributor target
type SalesTargetAllocationSummary struct {
	AllocatedTotal int64 `db:"allocated_total" json:"allocated_total"`
	IsAllocated    bool  `db:"is_allocated" json:"is_allocated"`
}

type SalesTargetMonthlyAllocation struct {
	Month          int   `db:"month" json:"month"`
	AllocatedTotal int64 `db:"allocated_total" json:"allocated_total"`
	TargetCount    int   `db:"target_count" json:"target_count"`
}
