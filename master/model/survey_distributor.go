package model

import "time"

// SurveyDistributor represents mst.m_survey_distributor table.
type SurveyDistributor struct {
	MSurveyDistributorId int        `db:"m_survey_distributor_id" json:"m_survey_distributor_id"`
	CustId               string     `db:"cust_id" json:"cust_id"`
	SurveyId             int        `db:"survey_id" json:"survey_id"`
	DistributorId        int        `db:"distributor_id" json:"distributor_id"`
	IsDel                bool       `db:"is_del" json:"is_del"`
	CreatedAt            *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedBy            *int64     `db:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt            *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	UpdatedBy            *int64     `db:"updated_by,omitempty" json:"updated_by,omitempty"`
	DeletedAt            *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	DeletedBy            *int64     `db:"deleted_by,omitempty" json:"deleted_by,omitempty"`
	DistributorCode      *string    `db:"distributor_code" json:"distributor_code,omitempty"`
	DistributorName      *string    `db:"distributor_name" json:"distributor_name,omitempty"`
}
