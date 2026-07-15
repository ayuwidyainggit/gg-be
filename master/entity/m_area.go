package entity

import "time"

type AreaResponse struct {
	AreaId        int        `json:"area_id"`
	AreaCode      string     `json:"area_code"`
	AreaName      string     `json:"area_name"`
	RegionId      int        `json:"region_id"`
	OfficialId    int        `json:"official_id"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type AreaListResponse struct {
	AreaId                int        `json:"area_id"`
	AreaCode              string     `json:"area_code"`
	AreaName              string     `json:"area_name"`
	RegionId              int        `json:"region_id"`
	RegionCode            string     `json:"region_code"`
	RegionName            string     `json:"region_name"`
	OfficialId            int        `json:"official_id"`
	OfficialType          int        `json:"official_type"`
	OfficialHierarchyCode string     `json:"official_hierarchy_code"`
	OfficialEmpName       string     `json:"official_emp_name"`
	IsActive              bool       `json:"is_active"`
	UpdatedBy             *int64     `json:"updated_by"`
	UpdatedByName         *string    `json:"updated_by_name"`
	UpdatedAt             *time.Time `json:"updated_at"`
}

type AreaLookupResponse struct {
	AreaId     int    `json:"area_id"`
	AreaCode   string `json:"area_code"`
	AreaName   string `json:"area_name"`
	RegionId   int    `json:"region_id"`
	OfficialId int    `json:"official_id"`
}

type CreateAreaBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	AreaCode   string `json:"area_code" validate:"required,max=4,alphanumericSpace"`
	AreaName   string `json:"area_name" validate:"required,max=25,alphanumericSpace"`
	RegionId   int    `json:"region_id" validate:"required"`
	OfficialId int    `json:"official_id" validate:""`
	IsActive   bool   `json:"is_active"`
}

type DetailAreaParams struct {
	AreaId int `params:"area_id" validate:"required"`
}

type UpdateAreaParams struct {
	AreaId int `params:"area_id" validate:"required"`
}

type DeleteAreaParams struct {
	AreaId int `params:"area_id" validate:"required"`
}

type UpdateAreaRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	AreaCode   string `json:"area_code" validate:"required,max=4,alphanumericSpace"`
	AreaName   string `json:"area_name" validate:"required,max=25,alphanumericSpace"`
	RegionId   int    `json:"region_id" validate:"required"`
	OfficialId int    `json:"official_id" validate:""`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type AreaQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive      *int   `query:"is_active"`
	RegionId      []int  `query:"region_id"`
	AreaId        []int  `query:"area_id"`
	EmployeeId    int
	DistributorId int
	Scope         EmployeeDropdownScope
}

type AreaIdCodeNameResp struct {
	AreaId   int    `json:"area_id"`
	AreaCode string `json:"area_code"`
	AreaName string `json:"area_name"`
}
