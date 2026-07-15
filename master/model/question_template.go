package model

import "time"

// QuestionTemplate represents mst.question_template table
type QuestionTemplate struct {
	QuestionTemplateId int        `db:"question_template_id" json:"question_template_id"`
	SurveyTemplateId   int        `db:"survey_template_id" json:"survey_template_id"`
	Question           string     `db:"question" json:"question"`
	InputType          string     `db:"input_type" json:"input_type"`
	AnswerType         string     `db:"answer_type" json:"answer_type"`
	UseImage           bool       `db:"use_image" json:"use_image"`
	Seq                int        `db:"seq" json:"seq"`
	IsDel              bool       `db:"is_del" json:"is_del"`
	CreatedAt          *time.Time `db:"created_at,omitempty" json:"created_at"`
	CreatedBy          *int64     `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt          *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedBy          *int64     `db:"updated_by,omitempty" json:"updated_by"`
	DeletedAt          *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	DeletedBy          *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
}
