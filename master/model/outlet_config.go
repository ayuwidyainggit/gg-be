package model

import "time"

// OutletConfig row from mst.m_outlet_config with user fullnames
type OutletConfig struct {
	CustId             string     `db:"cust_id" json:"cust_id"`
	OutletConfigId     int        `db:"outlet_config_id" json:"outlet_config_id"`
	VerificationStatus string     `db:"verification_status" json:"verification_status"`
	RulesType          string     `db:"rules_type" json:"rules_type"`
	CreatedBy          *int64    `db:"created_by" json:"created_by"`
	CreatedByName      *string   `db:"created_by_name" json:"created_by_name"`
	CreatedAt          *time.Time `db:"created_at" json:"created_at"`
	UpdatedBy          *int64    `db:"updated_by" json:"updated_by"`
	UpdatedByName      *string   `db:"updated_by_name" json:"updated_by_name"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updated_at"`
	Status             int        `db:"status" json:"status"`
}

// OutletConfigHeader minimal header for Detail (mst.m_outlet_config)
type OutletConfigHeader struct {
	OutletConfigId     int    `db:"outlet_config_id" json:"outlet_config_id"`
	CustId             string `db:"cust_id" json:"cust_id"`
	VerificationStatus string `db:"verification_status" json:"verification_status"`
	RulesType          string `db:"rules_type" json:"rules_type"`
}

// OutletConfigDet row from mst.m_outlet_config_det (with status_description from mst.m_outlet_config_status)
type OutletConfigDet struct {
	OutletConfigDetId   int     `db:"outlet_config_det_id" json:"outlet_config_det_id"`
	OutletConfigId      int     `db:"outlet_config_id" json:"outlet_config_id"`
	Status              int     `db:"status" json:"status"`
	StatusDescription   *string `db:"status_description" json:"-"`
	ValidateTrx         bool    `db:"validate_trx" json:"validate_trx"`
	CountingPeriod      *int    `db:"counting_period" json:"counting_period"`
}

// OutletConfigStatus row from mst.m_outlet_config_status
type OutletConfigStatus struct {
	OutletConfigStatusId int    `db:"outlet_config_status_id" json:"outlet_config_status_id"`
	StatusCode           string `db:"status_code" json:"status_code"`
	StatusDescription    string `db:"status_description" json:"status_description"`
	IsTrx                bool   `db:"is_trx" json:"is_trx"`
	SortOrder            int    `db:"sort_order" json:"sort_order"`
	IsActive             bool   `db:"is_active" json:"is_active"`
}
