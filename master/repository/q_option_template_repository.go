package repository

import (
	"master/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type QOptionTemplateRepository interface {
	FindAllByQuestionTemplateId(questionTemplateId int) ([]model.QOptionTemplate, error)
	FindAllByQuestionTemplateIds(questionTemplateIds []int) ([]model.QOptionTemplate, error)
	Store(tx *sqlx.Tx, option model.QOptionTemplate) (int, error)
	DeleteByQuestionTemplateId(tx *sqlx.Tx, questionTemplateId int, deletedBy int64) error
	DeleteByQuestionTemplateIds(tx *sqlx.Tx, questionTemplateIds []int, deletedBy int64) error
}

func NewQOptionTemplateRepository(db *sqlx.DB) QOptionTemplateRepository {
	return &qOptionTemplateRepositoryImpl{DB: db}
}

type qOptionTemplateRepositoryImpl struct {
	*sqlx.DB
}

func (r *qOptionTemplateRepositoryImpl) FindAllByQuestionTemplateId(questionTemplateId int) ([]model.QOptionTemplate, error) {
	var options []model.QOptionTemplate
	query := `SELECT q_option_template_id, question_template_id, option, seq 
		FROM mst.m_q_option_template 
		WHERE question_template_id = $1 AND is_del = false
		ORDER BY seq ASC`
	err := r.DB.Select(&options, query, questionTemplateId)
	return options, err
}

func (r *qOptionTemplateRepositoryImpl) FindAllByQuestionTemplateIds(questionTemplateIds []int) ([]model.QOptionTemplate, error) {
	if len(questionTemplateIds) == 0 {
		return []model.QOptionTemplate{}, nil
	}
	var options []model.QOptionTemplate
	query, args, err := sqlx.In(`SELECT q_option_template_id, question_template_id, option, seq 
		FROM mst.m_q_option_template 
		WHERE question_template_id IN (?) AND is_del = false
		ORDER BY question_template_id, seq ASC`, questionTemplateIds)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)
	err = r.DB.Select(&options, query, args...)
	return options, err
}

func (r *qOptionTemplateRepositoryImpl) Store(tx *sqlx.Tx, option model.QOptionTemplate) (int, error) {
	var id int
	query := `INSERT INTO mst.m_q_option_template 
		(question_template_id, option, seq, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING q_option_template_id`
	err := tx.Get(&id, query,
		option.QuestionTemplateId, option.Option, option.Seq,
		option.CreatedAt, option.CreatedBy, option.UpdatedAt, option.UpdatedBy)
	return id, err
}

func (r *qOptionTemplateRepositoryImpl) DeleteByQuestionTemplateId(tx *sqlx.Tx, questionTemplateId int, deletedBy int64) error {
	now := time.Now().UTC()
	query := `UPDATE mst.m_q_option_template 
		SET is_del = true, deleted_at = $1, deleted_by = $2
		WHERE question_template_id = $3 AND is_del = false`
	_, err := tx.Exec(query, now, deletedBy, questionTemplateId)
	return err
}

func (r *qOptionTemplateRepositoryImpl) DeleteByQuestionTemplateIds(tx *sqlx.Tx, questionTemplateIds []int, deletedBy int64) error {
	if len(questionTemplateIds) == 0 {
		return nil
	}
	now := time.Now().UTC()
	query, args, err := sqlx.In(`UPDATE mst.m_q_option_template 
		SET is_del = true, deleted_at = ?, deleted_by = ?
		WHERE question_template_id IN (?) AND is_del = false`, now, deletedBy, questionTemplateIds)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	return err
}
