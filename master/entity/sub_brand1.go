package entity

import (
	"time"
)

type SubBrand1QueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
	SBrand1Id    *int   `query:"sbrand1_id"`
	SBrand1Ids   []int  `query:"sbrand1_ids"`
	BrandIds     []int  `query:"brand_ids"`
}

type SubBrand1Response struct {
	BrandId     int        `json:"brand_id"`
	BrandCode   string     `json:"brand_code"`
	BrandName   string     `json:"brand_name"`
	Sbrand1Id   int        `json:"sbrand1_id"`
	Sbrand1Code string     `json:"sbrand1_code"`
	Sbrand1Name string     `json:"sbrand1_name"`
	PlId        int        `json:"pl_id"`
	PlCode      string     `json:"pl_code"`
	PlName      string     `json:"pl_name"`
	EffCall     int        `json:"eff_call"`
	MinItem     int        `json:"min_item"`
	IsActive    bool       `json:"is_active"`
	UpdatedBy   *int64     `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
type SubBrand1ListResponse struct {
	BrandId       int        `json:"brand_id"`
	BrandCode     string     `json:"brand_code"`
	BrandName     string     `json:"brand_name"`
	Sbrand1Id     int        `json:"sbrand1_id"`
	Sbrand1Code   string     `json:"sbrand1_code"`
	Sbrand1Name   string     `json:"sbrand1_name"`
	PlId          int        `json:"pl_id"`
	PlCode        string     `json:"pl_code"`
	PlName        string     `json:"pl_name"`
	EffCall       int        `json:"eff_call"`
	MinItem       int        `json:"min_item"`
	IsActive      bool       `json:"is_active"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName string     `json:"updated_by_name"`
}

type SubBrand1LookupResponse struct {
	BrandId     int    `json:"brand_id"`
	BrandCode   string `json:"brand_code"`
	BrandName   string `json:"brand_name"`
	Sbrand1Id   int    `json:"sbrand1_id"`
	Sbrand1Code string `json:"sbrand1_code"`
	Sbrand1Name string `json:"sbrand1_name"`
	PlId        int    `json:"pl_id"`
	PlCode      string `json:"pl_code"`
	PlName      string `json:"pl_name"`
	EffCall     int    `json:"eff_call"`
	MinItem     int    `json:"min_item"`
}

type SubBrand1LookupListResponse struct {
	BrandId       int    `json:"brand_id"`
	BrandCode     string `json:"brand_code"`
	BrandName     string `json:"brand_name"`
	Sbrand1Id     int    `json:"sbrand1_id"`
	Sbrand1Code   string `json:"sbrand1_code"`
	Sbrand1Name   string `json:"sbrand1_name"`
	EffCall       int    `json:"eff_call"`
	MinItem       int    `json:"min_item"`
	PlId          int    `json:"pl_id"`
	PlCode        string `json:"pl_code"`
	PlName        string `json:"pl_name"`
	UpdatedByName string `json:"updated_by_name"`
}

type SubBrand1MatGroupResponse struct {
	BrandId      int    `json:"brand_id"`
	BrandCode    string `json:"brand_code"`
	BrandName    string `json:"brand_name"`
	Sbrand1Id    int    `json:"sbrand1_id"`
	Sbrand1Code  string `json:"sbrand1_code"`
	Sbrand1Name  string `json:"sbrand1_name"`
	PlId         int    `json:"pl_id"`
	PlCode       string `json:"pl_code"`
	PlName       string `json:"pl_name"`
	MatGroupCode string `json:"mat_group_code"`
	MatGroupName string `json:"mat_group_name"`
}
type SubBrand1MatGroupListResponse struct {
	BrandId       int    `json:"brand_id"`
	BrandCode     string `json:"brand_code"`
	Sbrand1Id     int    `json:"sbrand1_id"`
	Sbrand1Code   string `json:"sbrand1_code"`
	PlId          int    `json:"pl_id"`
	PlCode        string `json:"pl_code"`
	MatGroupCode  string `json:"mat_group_code"`
	MatGroupName  string `json:"mat_group_name"`
	UpdatedByName string `json:"updated_by_name"`
}

type CreateSubBrand1Body struct {
	CustId      string  `json:"cust_id" validate:"required,max=10"`
	CreatedBy   int64   `json:"created_by" validate:"required"`
	BrandId     int     `json:"brand_id" validate:"required"`
	Sbrand1Code string  `json:"sbrand1_code" validate:"required,max=5,alphanumericSpace"`
	Sbrand1Name string  `json:"sbrand1_name" validate:"required,max=40,alphanumericSpace"`
	EffCall     float32 `json:"eff_call" validate:"min=0"`
	MinItem     float32 `json:"min_item" validate:"min=0"`
	IsActive    bool    `json:"is_active"`
}

type DetailSubBrand1Params struct {
	Sbrand1Id int `params:"sbrand1_id" validate:"required"`
}

type UpdateSubBrand1Params struct {
	Sbrand1Id int `params:"sbrand1_id" validate:"required"`
}

type DeleteSubBrand1Params struct {
	Sbrand1Id int `params:"sbrand1_id" validate:"required"`
}

type UpdateSubBrand1Request struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	UpdatedBy   int64  `json:"updated_by" validate:"required"`
	BrandId     int    `json:"brand_id" validate:"required"`
	Sbrand1Code string `json:"sbrand1_code,omitempty" validate:"max=5,omitempty,alphanumericSpace"`
	Sbrand1Name string `json:"sbrand1_name,omitempty" validate:"max=40,omitempty,alphanumericSpace"`
	EffCall     int    `json:"eff_call,omitempty" validate:""`
	MinItem     int    `json:"min_item,omitempty" validate:""`
	PgrpCode    string `json:"pgrp_code,omitempty" validate:""`
	PgroupName  string `json:"pgroup_name,omitempty" validate:""`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type SubBrandQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Query        string `query:"q"`
	Sort         string `query:"sort"`
	Status       []int  `query:"status"`
}

type SubBrandResponse struct {
	BrandId     int    `json:"brand_id"`
	Sbrand1Id   int    `json:"sbrand1_id"`
	Sbrand1Code string `json:"sbrand1_code"`
	Sbrand1Name string `json:"sbrand1_name"`
	EffCall     int    `json:"eff_call"`
	MinItem     int    `json:"min_item"`
	IsActive    bool   `json:"is_active"`
}
