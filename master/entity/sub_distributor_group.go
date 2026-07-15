package entity

import "time"

type CreateSubDistributorGroupBody struct {
	SubDistributorGroupID   int    `json:"sub_distributor_group_id"`
	CustID                  string `json:"cust_id" validate:"required,max=10"`
	SubDistributorGroupCode string `json:"sub_distributor_group_code" validate:"required,max=4,alphanumericSpace"`
	SubDistributorGroupName string `json:"sub_distributor_group_name" validate:"required,max=25,alphanumericSpace"`
	IsActive                bool   `json:"is_active"`
	CreatedBy               int64  `json:"created_by" validate:"required"`
}

type UpdateSubDistributorGroupBody struct {
	CustID                  string `json:"cust_id" validate:"required,max=10"`
	SubDistributorGroupCode string `json:"sub_distributor_group_code" validate:"required,max=4,alphanumericSpace"`
	SubDistributorGroupName string `json:"sub_distributor_group_name" validate:"required,max=25,alphanumericSpace"`
	IsActive                bool   `json:"is_active"`
	UpdatedBy               int64  `json:"updated_by" validate:"required"`
}
type SubDistributorGroupListResponse struct {
	SubDistributorGroupID   int        `json:"sub_distributor_group_id"`
	CustID                  string     `json:"cust_id"`
	SubDistributorGroupCode string     `json:"sub_distributor_group_code" `
	SubDistributorGroupName string     `json:"sub_distributor_group_name"`
	IsActive                bool       `json:"is_active"`
	UpdatedBy               *int64     `json:"updated_by"`
	UpdatedAt               *time.Time `json:"updated_at"`
	UpdatedByName           string     `json:"updated_by_name"`
}

type DetailSubDistributorGroupParams struct {
	SubDistributorGroupId int `params:"sub_distributor_group_id" validate:"required"`
}

type UpdateSubDistributorGroupParams struct {
	SubDistributorGroupId int `params:"sub_distributor_group_id" validate:"required"`
}

type DeleteSubDistributorGroupParams struct {
	SubDistributorGroupId int `params:"sub_distributor_group_id" validate:"required"`
}

type SubDistributorGroupResponse struct {
	SubDistributorGroupID   int        `json:"sub_distributor_group_id"`
	CustID                  string     `json:"cust_id"`
	SubDistributorGroupCode string     `json:"sub_distributor_group_code" `
	SubDistributorGroupName string     `json:"sub_distributor_group_name"`
	IsActive                bool       `json:"is_active"`
	UpdatedBy               *int64     `json:"updated_by"`
	UpdatedAt               *time.Time `json:"updated_at"`
	UpdatedByName           string     `json:"updated_by_name"`
}

type SubDistributorGroupLookupResponse struct {
	SubDistributorGroupID   int    `json:"sub_distributor_group_id"`
	SubDistributorGroupCode string `json:"sub_distributor_group_code" `
	SubDistributorGroupName string `json:"sub_distributor_group_name"`
}

type SubDistributorGroupQueryFilter struct {
	Page                  int    `query:"page"`
	Limit                 int    `query:"limit" validate:"required"`
	Query                 string `query:"q"`
	Mode                  string `query:"mode"`
	Sort                  string `query:"sort"`
	SubDistributorGroupID int    `query:"sub_distributor_group_id"`
	IsActive              *int   `query:"is_active"`
}

type SubDistributorGroupUpdateRequest struct {
	CustID                  string `json:"cust_id" validate:"required,max=10"`
	SubDistributorGroupCode string `json:"sub_distributor_group_code" validate:"required,max=4,alphanumericSpace"`
	SubDistributorGroupName string `json:"sub_distributor_group_name" validate:"required,max=25,alphanumericSpace"`
	UpdatedBy               int64  `json:"updated_by"`
	IsActive                *bool  `json:"is_active"`
}
