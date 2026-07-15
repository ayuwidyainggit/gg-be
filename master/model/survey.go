package model

import "time"

// Survey represents mst.m_survey table
type Survey struct {
	SurveyId          int        `db:"survey_id" json:"survey_id"`
	CustId            string     `db:"cust_id" json:"cust_id"`
	SurveyTitle       string     `db:"survey_title" json:"survey_title"`
	AnswerFrequency   string     `db:"answer_frequency" json:"answer_frequency"`
	ResponseType      string     `db:"response_type" json:"response_type"`
	TargetType        string     `db:"target_type" json:"target_type"`
	LevelTarget       *string    `db:"level_target" json:"level_target"`
	EmpId             *int       `db:"emp_id" json:"emp_id"`
	EfectiveDateStart *time.Time `db:"efective_date_start" json:"efective_date_start"`
	EfectiveDateEnd   *time.Time `db:"efective_date_end" json:"efective_date_end"`
	Status            int        `db:"status" json:"status"`
	IsDel             bool       `db:"is_del" json:"is_del"`
	CreatedAt         *time.Time `db:"created_at,omitempty" json:"created_at"`
	CreatedBy         *int64     `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt         *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedBy         *int64     `db:"updated_by,omitempty" json:"updated_by"`
	DeletedAt         *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	DeletedBy         *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
	// Joined fields
	SalesName       *string `db:"sales_name" json:"sales_name,omitempty"`
	DistributorId   *int    `db:"distributor_id" json:"distributor_id,omitempty"`
	DistributorCode *string `db:"distributor_code" json:"distributor_code,omitempty"`
	DistributorName *string `db:"distributor_name" json:"distributor_name,omitempty"`
}
