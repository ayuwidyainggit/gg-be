package repository

import (
	"fmt"
	"master/entity"
	"master/model"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SurveyTemplateRepository interface {
	FindAllByCustId(filter entity.SurveyTemplateQueryFilter, custId string) ([]model.SurveyTemplate, int, int, error)
	FindOneById(surveyTemplateId int, custId string) (model.SurveyTemplate, error)
	Store(tx *sqlx.Tx, template model.SurveyTemplate) (int, error)
	Update(tx *sqlx.Tx, surveyTemplateId int, custId string, template model.SurveyTemplate) error
	Delete(tx *sqlx.Tx, custId string, surveyTemplateId int, deletedBy int64) error
	GenerateTemplateCode(custId string) (string, error)
	BeginTx() (*sqlx.Tx, error)
}

func NewSurveyTemplateRepository(db *sqlx.DB) SurveyTemplateRepository {
	return &surveyTemplateRepositoryImpl{DB: db}
}

type surveyTemplateRepositoryImpl struct {
	*sqlx.DB
}

func (r *surveyTemplateRepositoryImpl) BeginTx() (*sqlx.Tx, error) {
	return r.DB.Beginx()
}

func (r *surveyTemplateRepositoryImpl) GenerateTemplateCode(custId string) (string, error) {
	var lastCode string
	query := `SELECT COALESCE(MAX(template_code), 'TMP000') FROM mst.m_survey_template WHERE cust_id = $1`
	err := r.DB.Get(&lastCode, query, custId)
	if err != nil {
		return "", err
	}

	// Extract number from last code and increment
	var num int
	fmt.Sscanf(lastCode, "TMP%d", &num)
	num++

	return fmt.Sprintf("TMP%03d", num), nil
}

func (r *surveyTemplateRepositoryImpl) FindAllByCustId(filter entity.SurveyTemplateQueryFilter, custId string) ([]model.SurveyTemplate, int, int, error) {
	var templates []model.SurveyTemplate
	var total int

	// Base query
	baseQuery := `FROM mst.m_survey_template WHERE cust_id = $1 AND is_del = false`
	args := []interface{}{custId}
	argIdx := 2

	// Filter by status
	if filter.Status != nil {
		if *filter.Status == 1 {
			baseQuery += fmt.Sprintf(" AND is_active = $%d", argIdx)
			args = append(args, true)
		} else {
			baseQuery += fmt.Sprintf(" AND is_active = $%d", argIdx)
			args = append(args, false)
		}
		argIdx++
	}

	// Filter by search query
	if filter.Query != "" {
		baseQuery += fmt.Sprintf(" AND (template_title ILIKE $%d OR template_code ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.DB.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, 0, err
	}

	// Calculate pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(filter.Limit)))

	// Sort
	sortField := "created_at"
	sortOrder := "DESC"
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ":")
		if len(parts) == 2 {
			sortField = parts[0]
			sortOrder = strings.ToUpper(parts[1])
		}
	}

	// Data query
	dataQuery := fmt.Sprintf(`SELECT survey_template_id, cust_id, template_code, template_title, 
		question_total, use_image, is_active, created_at %s ORDER BY %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, sortField, sortOrder, argIdx, argIdx+1)
	args = append(args, filter.Limit, offset)

	err = r.DB.Select(&templates, dataQuery, args...)
	if err != nil {
		return nil, 0, 0, err
	}

	return templates, total, lastPage, nil
}

func (r *surveyTemplateRepositoryImpl) FindOneById(surveyTemplateId int, custId string) (model.SurveyTemplate, error) {
	var template model.SurveyTemplate
	query := `SELECT survey_template_id, cust_id, template_code, template_title, 
		question_total, use_image, is_active, created_at 
		FROM mst.m_survey_template 
		WHERE survey_template_id = $1 AND cust_id = $2 AND is_del = false`
	err := r.DB.Get(&template, query, surveyTemplateId, custId)
	return template, err
}

func (r *surveyTemplateRepositoryImpl) Store(tx *sqlx.Tx, template model.SurveyTemplate) (int, error) {
	var id int
	query := `INSERT INTO mst.m_survey_template 
		(cust_id, template_code, template_title, question_total, use_image, is_active, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING survey_template_id`
	err := tx.Get(&id, query,
		template.CustId, template.TemplateCode, template.TemplateTitle,
		template.QuestionTotal, template.UseImage, template.IsActive,
		template.CreatedAt, template.CreatedBy, template.UpdatedAt, template.UpdatedBy)
	return id, err
}

func (r *surveyTemplateRepositoryImpl) Update(tx *sqlx.Tx, surveyTemplateId int, custId string, template model.SurveyTemplate) error {
	query := `UPDATE mst.m_survey_template 
		SET template_title = $1, question_total = $2, use_image = $3, is_active = $4, updated_at = $5, updated_by = $6
		WHERE survey_template_id = $7 AND cust_id = $8 AND is_del = false`
	_, err := tx.Exec(query,
		template.TemplateTitle, template.QuestionTotal, template.UseImage, template.IsActive,
		template.UpdatedAt, template.UpdatedBy, surveyTemplateId, custId)
	return err
}

func (r *surveyTemplateRepositoryImpl) Delete(tx *sqlx.Tx, custId string, surveyTemplateId int, deletedBy int64) error {
	now := time.Now().UTC()
	query := `UPDATE mst.m_survey_template 
		SET is_del = true, deleted_at = $1, deleted_by = $2
		WHERE survey_template_id = $3 AND cust_id = $4 AND is_del = false`
	result, err := tx.Exec(query, now, deletedBy, surveyTemplateId, custId)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("survey template not found")
	}
	return nil
}
