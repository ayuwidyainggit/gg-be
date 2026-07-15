package entity

// BusinessUnitQueryFilter for business-unit endpoint
type BusinessUnitQueryFilter struct {
	CustId        string // from JWT context
	ParentCustId  string // from JWT context
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	Query         string `query:"q"`
	Sort          string `query:"sort"`
	UserName      string `query:"user_name"`
	DistributorId *int   `query:"distributor_id"` // NULL/0 = principal, NOT NULL = distributor
	RegionId      []int  `query:"region_id"`
	AreaId        []int  `query:"area_id"`
	IsActive      []int  `query:"is_active"`
	EmployeeId    int
	Scope         EmployeeDropdownScope
}

// BusinessUnitDistributorData for distributor_data array in principal response
type BusinessUnitDistributorData struct {
	CustId          string `json:"cust_id"`
	DistributorId   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	AreaId          int    `json:"area_id"`
	AreaCode        string `json:"area_code"`
	AreaName        string `json:"area_name"`
	RegionId        int    `json:"region_id"`
	RegionCode      string `json:"region_code"`
	RegionName      string `json:"region_name"`
}

// BusinessUnitPrincipalResponse for principal user (distributor_id = NULL)
type BusinessUnitPrincipalResponse struct {
	CustId          string                        `json:"cust_id"`
	UserId          int                           `json:"user_id"`
	UserFullname    string                        `json:"user_fullname"`
	CustName        string                        `json:"cust_name"`
	DistributorId   string                        `json:"distributor_id"` // empty string for principal
	DistributorData []BusinessUnitDistributorData `json:"distributor_data"`
}

// BusinessUnitDistributorResponse for distributor user (distributor_id = NOT NULL)
type BusinessUnitDistributorResponse struct {
	CustId          string `json:"cust_id"`
	UserId          int    `json:"user_id"`
	UserFullname    string `json:"user_fullname"`
	DistributorId   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	AreaId          int    `json:"area_id"`
	AreaCode        string `json:"area_code"`
	AreaName        string `json:"area_name"`
	RegionId        int    `json:"region_id"`
	RegionCode      string `json:"region_code"`
	RegionName      string `json:"region_name"`
}
