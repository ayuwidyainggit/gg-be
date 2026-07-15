package model

import "time"

// QOptionTemplate represents mst.m_q_option_template table
type QOptionTemplate struct {
	QOptionTemplateId  int        `db:"q_option_template_id" json:"q_option_template_id"`
	QuestionTemplateId int        `db:"question_template_id" json:"question_template_id"`
	Option             string     `db:"option" json:"option"`
	Seq                int        `db:"seq" json:"seq"`
	IsDel              bool       `db:"is_del" json:"is_del"`
	CreatedAt          *time.Time `db:"created_at,omitempty" json:"created_at"`
	CreatedBy          *int64     `db:"created_by,omitempty" json:"created_by"`
	UpdatedAt          *time.Time `db:"updated_at,omitempty" json:"updated_at"`
	UpdatedBy          *int64     `db:"updated_by,omitempty" json:"updated_by"`
	DeletedAt          *time.Time `db:"deleted_at,omitempty" json:"deleted_at"`
	DeletedBy          *int64     `db:"deleted_by,omitempty" json:"deleted_by"`
}
