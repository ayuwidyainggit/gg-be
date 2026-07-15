package entity

import "time"

// OutletConfigListFilter query params for GET /v1/outlet_config
type OutletConfigListFilter struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Sort       string `query:"sort"`
	Q          string `query:"q"`
	RulesType  string `query:"rules_type"`
	Status     *int   `query:"status"` // 1: is_active true, 0: is_active false, null: all
}

// OutletConfigListResponse item for list API
type OutletConfigListResponse struct {
	CustId             string     `json:"cust_id"`   
	OutletConfigId     int        `json:"outlet_config_id"`
	VerificationStatus string     `json:"verification_status"`
	RulesType          string     `json:"rules_type"`
	CreatedBy          string     `json:"created_by"`   // sys.m_user.user_fullname
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedBy          string     `json:"updated_by"`   // sys.m_user.user_fullname
	UpdatedAt          *time.Time `json:"updated_at"`
	Status             int        `json:"status"`
	IsActive           bool       `json:"is_active"`
	CanEdit            bool       `json:"can_edit"`  
	CanDelete          bool       `json:"can_delete"` 
}

// OutletConfigDetailParams path param for GET /v1/outlet_config/:outlet_config_id
type OutletConfigDetailParams struct {
	OutletConfigId int `params:"outlet_config_id" validate:"required"`
}

// OutletConfigDetItem item in outlet_status array (mst.m_outlet_config_det)
type OutletConfigDetItem struct {
	OutletConfigDetId int   `json:"outlet_config_det_id"`
	Status            int   `json:"status"`
	StatusDesc        string `json:"status_desc"`   // from appendix
	ValidateTrx       bool   `json:"validate_trx"`
	CountingPeriod    *int   `json:"counting_period"` // nullable
}

// OutletConfigDetailResponse response for Detail API
type OutletConfigDetailResponse struct {
	CustId             string                 `json:"cust_id"` 
	OutletConfigId     int                    `json:"outlet_config_id"`
	VerificationStatus string                 `json:"verification_status"`
	RulesType          string                 `json:"rules_type"`
	OutletStatus       []OutletConfigDetItem  `json:"outlet_status"`
	CanEdit            bool                   `json:"can_edit"`
	CanDelete          bool                   `json:"can_delete"`
}

// CreateOutletConfigDetItem item in outlet_status array for Create
type CreateOutletConfigDetItem struct {
	Status         int  `json:"status" validate:"required"`
	ValidateTrx    bool `json:"validate_trx"`
	CountingPeriod *int `json:"counting_period"`
}

// CreateOutletConfigBody body for POST /v1/outlet_config
type CreateOutletConfigBody struct {
	VerificationStatus string                   `json:"verification_status" validate:"required"`
	RulesType          string                   `json:"rules_type" validate:"required"`
	OutletStatus       []CreateOutletConfigDetItem `json:"outlet_status" validate:"required,dive"`
}

// OutletConfigStatusListFilter query params for GET /v1/outlet_config_status
type OutletConfigStatusListFilter struct {
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
	Sort     string  `query:"sort"`
	Q        string  `query:"q"`
	IsActive *bool   `query:"is_active"` // true = active only, false = inactive only, nil = all
}

// OutletConfigStatusListResponse item for outlet_config_status list API
type OutletConfigStatusListResponse struct {
	OutletConfigStatusId int    `json:"outlet_config_status_id"`
	StatusCode           string `json:"status_code"`
	StatusDescription    string `json:"status_description"`
	IsTrx                bool   `json:"is_trx"`
	SortOrder            int    `json:"sort_order"`
	IsActive             bool   `json:"is_active"`
}
