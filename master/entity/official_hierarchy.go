package entity

import (
	"time"
)

type OfficialHierarchyQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	OfficialType int    `query:"official_type"`
}

type OfficialHierarchyResponse struct {
	OfficialType  int    `json:"official_type"`
	HierarchyCode string `json:"hierarchy_code"`
	IsActive      bool   `json:"is_active"`
}

type OfficialHierarchyListResponse struct {
	OfficialType  int    `json:"official_type"`
	HierarchyCode string `json:"hierarchy_code"`
	IsActive      bool   `json:"is_active"`
}

type OfficialHierarchyLookupResponse struct {
	OfficialType  int    `json:"official_type"`
	HierarchyCode string `json:"hierarchy_code"`
}

type CreateOfficialHierarchyBody struct {
	CustId        string    `json:"cust_id" validate:""`
	CreatedBy     int64     `json:"created_by" validate:""`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedBy     int64     `json:"updated_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	OfficialType  int       `json:"official_type" validate:"required,min=1,max=3"`
	HierarchyCode string    `json:"hierarchy_code" validate:"required,max=5"`
	IsActive      *bool     `json:"is_active"`
}

type DetailOfficialHierarchyParams struct {
	OfficialType int `params:"official_type" validate:"required"`
}

type UpdateOfficialHierarchyParams struct {
	OfficialType int `params:"official_type" validate:"required"`
}

type DeleteOfficialHierarchyParams struct {
	OfficialType int `params:"official_type" validate:"required"`
}

type UpdateOfficialHierarchyRequest struct {
	CustId        string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy     int64  `json:"updated_by" validate:"required"`
	OfficialType  int    `json:"official_type,omitempty" validate:"max=3,omitempty"`
	HierarchyCode string `json:"hierarchy_code" validate:"max=5,omitempty"`
	IsActive      *bool  `json:"is_active,omitempty"`
}

type BulkUpsertOfficialHierarchyBody struct {
	CustID                  string                        `json:"cust_id" validate:"required,max=10"`
	CreatedBy               int64                         `json:"created_by" validate:"required"`
	CreatedAt               time.Time                     `json:"created_at"`
	UpdatedBy               int64                         `json:"updated_by"`
	UpdatedAt               time.Time                     `json:"updated_at"`
	UpsertOfficialHierarchy []CreateOfficialHierarchyBody `json:"data" validate:"required,min=3,dive"`
}
