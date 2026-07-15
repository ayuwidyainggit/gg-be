package repository

import (
	"master/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type QuestionTemplateRepository interface {
	FindAllBySurveyTemplateId(surveyTemplateId int) ([]model.QuestionTemplate, error)
	Store(tx *sqlx.Tx, question model.QuestionTemplate) (int, error)
	Update(tx *sqlx.Tx, questionTemplateId int, question model.QuestionTemplate) error
	DeleteBySurveyTemplateId(tx *sqlx.Tx, surveyTemplateId int, deletedBy int64) error
	DeleteByIds(tx *sqlx.Tx, ids []int, deletedBy int64) error
}

func NewQuestionTemplateRepository(db *sqlx.DB) QuestionTemplateRepository {
	return &questionTemplateRepositoryImpl{DB: db}
}

type questionTemplateRepositoryImpl struct {
	*sqlx.DB
}

func (r *questionTemplateRepositoryImpl) FindAllBySurveyTemplateId(surveyTemplateId int) ([]model.QuestionTemplate, error) {
	var questions []model.QuestionTemplate
	query := `SELECT question_template_id, survey_template_id, question, input_type, answer_type, use_image, seq
		FROM mst.question_template
		WHERE survey_template_id = $1 AND is_del = false
		ORDER BY seq ASC`
	err := r.DB.Select(&questions, query, surveyTemplateId)
	return questions, err
}

func (r *questionTemplateRepositoryImpl) Store(tx *sqlx.Tx, question model.QuestionTemplate) (int, error) {
	var id int
	query := `INSERT INTO mst.question_template 
		(survey_template_id, question, input_type, answer_type, use_image, seq, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING question_template_id`
	err := tx.Get(&id, query,
		question.SurveyTemplateId, question.Question, question.InputType, question.AnswerType, question.UseImage, question.Seq,
		question.CreatedAt, question.CreatedBy, question.UpdatedAt, question.UpdatedBy)
	return id, err
}

func (r *questionTemplateRepositoryImpl) Update(tx *sqlx.Tx, questionTemplateId int, question model.QuestionTemplate) error {
	query := `UPDATE mst.question_template 
		SET question = $1, input_type = $2, answer_type = $3, seq = $4, updated_at = $5, updated_by = $6
		WHERE question_template_id = $7 AND is_del = false`
	_, err := tx.Exec(query, question.Question, question.InputType, question.AnswerType, question.Seq,
		question.UpdatedAt, question.UpdatedBy, questionTemplateId)
	return err
}

func (r *questionTemplateRepositoryImpl) DeleteBySurveyTemplateId(tx *sqlx.Tx, surveyTemplateId int, deletedBy int64) error {
	now := time.Now().UTC()
	query := `UPDATE mst.question_template 
		SET is_del = true, deleted_at = $1, deleted_by = $2
		WHERE survey_template_id = $3 AND is_del = false`
	_, err := tx.Exec(query, now, deletedBy, surveyTemplateId)
	return err
}

func (r *questionTemplateRepositoryImpl) DeleteByIds(tx *sqlx.Tx, ids []int, deletedBy int64) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now().UTC()
	query, args, err := sqlx.In(`UPDATE mst.question_template 
		SET is_del = true, deleted_at = ?, deleted_by = ?
		WHERE question_template_id IN (?) AND is_del = false`, now, deletedBy, ids)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	_, err = tx.Exec(query, args...)
	return err
}
