package entity

import "time"

type RegionResponse struct {
	RegionId      int        `json:"region_id"`
	RegionCode    string     `json:"region_code"`
	RegionName    string     `json:"region_name"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type RegionLookupResponse struct {
	RegionId   int    `json:"region_id"`
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
}

type CreateRegionBody struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	CreatedBy  int64  `json:"created_by" validate:"required"`
	RegionCode string `json:"region_code" validate:"required,max=4,alphanumericSpace"`
	RegionName string `json:"region_name" validate:"required,max=25,alphanumericSpace"`
	IsActive   bool   `json:"is_active"`
}

type DetailRegionParams struct {
	RegionId int `params:"region_id" validate:"required"`
}

type UpdateRegionParams struct {
	RegionId int `params:"region_id" validate:"required"`
}

type DeleteRegionParams struct {
	RegionId int `params:"region_id" validate:"required"`
}

type UpdateRegionRequest struct {
	CustId     string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy  int64  `json:"updated_by" validate:"required"`
	RegionCode string `json:"region_code,omitempty" validate:"required,max=4,alphanumericSpace"`
	RegionName string `json:"region_name,omitempty" validate:"max=25,omitempty,alphanumericSpace"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type RegionIdCodeNameResp struct {
	RegionId   int    `json:"region_id"`
	RegionCode string `json:"region_code"`
	RegionName string `json:"region_name"`
}

type RegionQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	RegionId     []int  `query:"region_id"`
	EmployeeId   int
	DistributorId int
	Scope        EmployeeDropdownScope
}
