package model

import "time"

// SurveyTemplate represents mst.m_survey_template table
type SurveyTemplate struct {
	SurveyTemplateId int        `db:"survey_template_id" json:"survey_template_id"`
	CustId           string     `db:"cust_id" json:"cust_id"`
	TemplateCode     string     `db:"template_code" json:"template_code"`
	TemplateTitle    string     `db:"template_title" json:"template_title"`
	QuestionTotal    int        `db:"question_total" json:"question_total"`
	UseImage         bool       `db:"use_image" json:"use_image"`
	IsActive         bool       `db:"is_active" json:"is_active"`
	IsDel            bool       `db:"is_del" json:"is_del"`
	CreatedAt        *time.Time `db:"created_at,omitempty" json:"created_at"`
	CreatedBy        *int64     `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt        *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedBy        *int64     `db:"updated_by,omitempty" json:"updated_by"`
	DeletedAt        *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	DeletedBy        *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
}
