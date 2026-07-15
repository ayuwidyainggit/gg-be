package entity

import "time"

type SalesTargetDistributorQueryFilter struct {
	GeneralQueryFilter
	Year   *int   `query:"year"`
	Status *[]int `query:"status"`
}

type SalesTargetDistributorMonthlyDetail struct {
	SalesTargetDistributorMonthlyId int    `json:"sales_target_distributor_monthly_id"`
	Month                           int    `json:"month"`
	MonthlyTarget                   int    `json:"monthly_target"`
	IsActive                        bool   `json:"is_active"`
	AllocatedTotal                  int64  `json:"allocated_total"`
	Remaining                       int64  `json:"remaining"`
	IsAllocated                     bool   `json:"is_allocated"`
	IsPastMonth                     bool   `json:"is_past_month"`
	IsEditable                      bool   `json:"is_editable"`
	DisableReason                   string `json:"disable_reason"`
}

type SalesTargetDistributorListResponse struct {
	SalesTargetDistributorYearlyId int        `json:"sales_target_distributor_yearly_id"`
	DistributorId                  int        `json:"distributor_id"`
	DistributorCode                string     `json:"distributor_code"`
	DistributorName                string     `json:"distributor_name"`
	Year                           int        `json:"year"`
	YearlyTarget                   int        `json:"yearly_target"`
	UpdatedBy                      int64      `json:"updated_by"`
	UpdatedByName                  string     `json:"updated_by_name"`
	UpdatedAt                      *time.Time `json:"updated_at"`
	Status                         string     `json:"status"`
	UserInactive                   *int64     `json:"user_inactive"`
	InactiveAt                     *time.Time `json:"inactive_at"`
}

type SalesTargetDistributorDetailResponse struct {
	SalesTargetDistributorYearlyId int                                   `json:"sales_target_distributor_yearly_id"`
	Year                           int                                   `json:"year"`
	YearlyTarget                   int                                   `json:"yearly_target"`
	UpdatedBy                      int64                                 `json:"updated_by"`
	UpdatedByName                  string                                `json:"updated_by_name"`
	UpdatedAt                      *time.Time                            `json:"updated_at"`
	AreaId                         int                                   `json:"area_id"`
	AreaCode                       string                                `json:"area_code"`
	AreaName                       string                                `json:"area_name"`
	RegionId                       int                                   `json:"region_id"`
	RegionCode                     string                                `json:"region_code"`
	RegionName                     string                                `json:"region_name"`
	DistributorId                  int                                   `json:"distributor_id"`
	DistributorCode                string                                `json:"distributor_code"`
	DistributorName                string                                `json:"distributor_name"`
	Status                         string                                `json:"status"`
	UserInactive                   *int64                                `json:"user_inactive"`
	InactiveAt                     *time.Time                            `json:"inactive_at"`
	IsAllocated                    bool                                  `json:"is_allocated"`
	AllocationTotal                int64                                 `json:"allocation_total"`
	Details                        []SalesTargetDistributorMonthlyDetail `json:"details"`
}

type CreateSalesTargetDistributorMonthly struct {
	Month         int `json:"month" validate:"required"`
	MonthlyTarget int `json:"monthly_target" validate:"required"`
}

type UpdateSalesTargetDistributorMonthly struct {
	Month         int `json:"month"`
	MonthlyTarget int `json:"monthly_target"`
}

type CreateSalesTargetDistributorBody struct {
	Year          int                                   `json:"year" validate:"required"`
	AreaId        int                                   `json:"area_id" validate:"required"`
	RegionId      int                                   `json:"region_id" validate:"required"`
	DistributorId int                                   `json:"distributor_id" validate:"required"`
	YearlyTarget  int                                   `json:"yearly_target" validate:"required"`
	Status        *int                                  `json:"status" validate:"required"`
	Data          []CreateSalesTargetDistributorMonthly `json:"data" validate:"required"`
	// Internal fields
	CustId    string `json:"-"`
	CreatedBy int64  `json:"-"`
}

type UpdateSalesTargetDistributorRequest struct {
	Year          *int                                  `json:"year"`
	AreaId        *int                                  `json:"area_id"`
	RegionId      *int                                  `json:"region_id"`
	DistributorId *int                                  `json:"distributor_id"`
	YearlyTarget  *int                                  `json:"yearly_target"`
	Status        *int                                  `json:"status"`
	IsActive      *bool                                 `json:"is_active"`
	Data          []UpdateSalesTargetDistributorMonthly `json:"data"`
	// Internal fields
	CustId       string     `json:"-"`
	UpdatedBy    int64      `json:"-"`
	UserInactive *int64     `json:"-"`
	InactiveAt   *time.Time `json:"-"`
}

type DetailSalesTargetDistributorParams struct {
	SalesTargetDistributorYearlyId int `params:"sales_target_distributor_yearly_id" validate:"required"`
}

type UpdateSalesTargetDistributorParams struct {
	SalesTargetDistributorYearlyId int `params:"sales_target_distributor_yearly_id" validate:"required"`
}
