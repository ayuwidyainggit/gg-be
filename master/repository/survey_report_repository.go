package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SurveyReportRepository interface {
	FindList(filter entity.SurveyReportQueryFilter, withPagination bool) ([]model.SurveyReportListRow, int, int, error)
	FindDetail(surveyAnswerID int64, custID string) (model.SurveyReportDetailRow, error)
	FindQuestionDetails(surveyAnswerID int64, custID string) ([]model.SurveyReportQuestionRow, error)
	FindQuestionOptions(questionTemplateIDs []int64) ([]model.SurveyReportQuestionOptionRow, error)
	FindSelectedOptions(detailIDs []int64, custID string) ([]model.SurveyReportSelectedOptionRow, error)
	FindAnswerFiles(detailIDs []int64, custID string) ([]model.SurveyReportAnswerFileRow, error)
	FindExportRows(filter entity.SurveyReportQueryFilter) ([]model.SurveyReportExportRow, error)
}

type surveyReportRepositoryImpl struct {
	*sqlx.DB
}

func NewSurveyReportRepository(db *sqlx.DB) SurveyReportRepository {
	return &surveyReportRepositoryImpl{DB: db}
}

func (r *surveyReportRepositoryImpl) FindList(filter entity.SurveyReportQueryFilter, withPagination bool) ([]model.SurveyReportListRow, int, int, error) {
	var rows []model.SurveyReportListRow

	baseQuery, args, err := r.buildListBaseQuery(filter)
	if err != nil {
		return nil, 0, 0, err
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	var total int
	if err := r.DB.Get(&total, countQuery, args...); err != nil {
		return nil, 0, 0, err
	}

	lastPage := 1
	if filter.Limit > 0 {
		lastPage = (total + filter.Limit - 1) / filter.Limit
		if lastPage == 0 {
			lastPage = 1
		}
	}

	selectQuery := `SELECT
		sa.survey_answer_id,
		sa.survey_id,
		s.survey_title,
		s.answer_frequency,
		s.response_type,
		sa.answer_date,
		COALESCE(sa.created_at, sa.answer_date) AS created_date,
		s.efective_date_start AS effective_date_start,
		s.efective_date_end AS effective_date_end,
		COALESCE(md.area_id, sa.area_id) AS area_id,
		COALESCE(a.area_code, '') AS area_code,
		COALESCE(a.area_name, '') AS area_name,
		md.distributor_id,
		COALESCE(md.distributor_code, '') AS distributor_code,
		COALESCE(md.distributor_name, '') AS distributor_name,
		sa.outlet_id,
		COALESCE(o.outlet_principal_code, o.outlet_code, '') AS outlet_code,
		COALESCE(o.outlet_name, '') AS outlet_name,
		sa.emp_id,
		COALESCE(emp.emp_code, '') AS emp_code,
		COALESCE(emp.emp_name, '') AS emp_name,
		COALESCE(sm.sales_name, emp.emp_name, '') AS salesman_name,
		COALESCE(sa.status, '') AS status ` + baseQuery + ` ORDER BY ` + buildSurveyReportSortClause(filter.Sort)

	dataArgs := append([]interface{}{}, args...)
	if withPagination {
		offset := (filter.Page - 1) * filter.Limit
		selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(dataArgs)+1, len(dataArgs)+2)
		dataArgs = append(dataArgs, filter.Limit, offset)
	}

	if err := r.DB.Select(&rows, selectQuery, dataArgs...); err != nil {
		return nil, 0, 0, err
	}

	return rows, total, lastPage, nil
}

func (r *surveyReportRepositoryImpl) buildListBaseQuery(filter entity.SurveyReportQueryFilter) (string, []interface{}, error) {
	return r.buildListBaseQueryWithExtraJoins(filter, "")
}

func (r *surveyReportRepositoryImpl) buildListBaseQueryWithExtraJoins(filter entity.SurveyReportQueryFilter, extraJoins string) (string, []interface{}, error) {
	baseQuery := `FROM mst.survey_answer sa
		JOIN mst.m_survey s ON s.survey_id = sa.survey_id AND s.is_del = false
		LEFT JOIN smc.m_customer mc ON mc.cust_id = sa.cust_id
		LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.is_del = false
		LEFT JOIN mst.m_area a ON a.area_id = COALESCE(md.area_id, sa.area_id)
		LEFT JOIN mst.m_outlet o ON o.outlet_id = sa.outlet_id AND o.cust_id = sa.cust_id AND o.is_del = false
		LEFT JOIN mst.m_salesman sm ON sm.emp_id = sa.emp_id AND sm.cust_id = sa.cust_id AND sm.is_del = false
		LEFT JOIN mst.m_employee emp ON emp.emp_id = sa.emp_id AND emp.cust_id = sa.cust_id`
	if extraJoins != "" {
		baseQuery += extraJoins
	}
	baseQuery += `
		WHERE s.cust_id = $1 AND sa.is_del = false`
	args := []interface{}{filter.CustID}
	argIdx := 2

	if filter.Query != "" {
		baseQuery += fmt.Sprintf(" AND (s.survey_title ILIKE $%d OR COALESCE(o.outlet_name, '') ILIKE $%d OR COALESCE(o.outlet_principal_code, o.outlet_code, '') ILIKE $%d OR COALESCE(a.area_name, '') ILIKE $%d OR COALESCE(a.area_code, '') ILIKE $%d OR COALESCE(md.distributor_name, '') ILIKE $%d OR COALESCE(md.distributor_code, '') ILIKE $%d OR COALESCE(emp.emp_name, sm.sales_name, '') ILIKE $%d OR COALESCE(emp.emp_code, '') ILIKE $%d)", argIdx, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND sa.answer_date >= $%d::date", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}

	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND sa.answer_date < ($%d::date + INTERVAL '1 day')", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}

	if len(filter.SurveyID) > 0 {
		query, inArgs, err := sqlx.In(" AND s.survey_id IN (?)", filter.SurveyID)
		if err != nil {
			return "", nil, err
		}
		query = r.DB.Rebind(query)
		query = shiftPositionalPlaceholders(query, argIdx-1)
		baseQuery += query
		args = append(args, inArgs...)
		argIdx += len(inArgs)
	}

	if len(filter.SurveyTitle) > 0 {
		query, inArgs, err := sqlx.In(" AND s.survey_title IN (?)", filter.SurveyTitle)
		if err != nil {
			return "", nil, err
		}
		query = r.DB.Rebind(query)
		query = shiftPositionalPlaceholders(query, argIdx-1)
		baseQuery += query
		args = append(args, inArgs...)
		argIdx += len(inArgs)
	}

	if len(filter.AreaID) > 0 {
		query, inArgs, err := sqlx.In(" AND COALESCE(md.area_id, sa.area_id) IN (?)", filter.AreaID)
		if err != nil {
			return "", nil, err
		}
		query = r.DB.Rebind(query)
		query = shiftPositionalPlaceholders(query, argIdx-1)
		baseQuery += query
		args = append(args, inArgs...)
		argIdx += len(inArgs)
	}

	return baseQuery, args, nil
}

func (r *surveyReportRepositoryImpl) FindDetail(surveyAnswerID int64, custID string) (model.SurveyReportDetailRow, error) {
	var row model.SurveyReportDetailRow
	query := `SELECT
		sa.survey_answer_id,
		sa.cust_id AS answer_cust_id,
		sa.survey_id,
		s.survey_title,
		s.answer_frequency,
		s.response_type,
		sa.answer_date,
		COALESCE(sa.created_at, sa.answer_date) AS created_date,
		s.efective_date_start AS effective_date_start,
		s.efective_date_end AS effective_date_end,
		COALESCE(md.area_id, sa.area_id) AS area_id,
		COALESCE(a.area_code, '') AS area_code,
		COALESCE(a.area_name, '') AS area_name,
		md.distributor_id,
		COALESCE(md.distributor_code, '') AS distributor_code,
		COALESCE(md.distributor_name, '') AS distributor_name,
		sa.outlet_id,
		COALESCE(o.outlet_principal_code, o.outlet_code, '') AS outlet_code,
		COALESCE(o.outlet_name, '') AS outlet_name,
		sa.emp_id,
		COALESCE(emp.emp_code, '') AS emp_code,
		COALESCE(emp.emp_name, '') AS emp_name,
		COALESCE(sm.sales_name, emp.emp_name, '') AS salesman_name,
		COALESCE(sa.status, '') AS status
		FROM mst.survey_answer sa
		JOIN mst.m_survey s ON s.survey_id = sa.survey_id AND s.is_del = false
		LEFT JOIN smc.m_customer mc ON mc.cust_id = sa.cust_id
		LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id AND md.is_del = false
		LEFT JOIN mst.m_area a ON a.area_id = COALESCE(md.area_id, sa.area_id)
		LEFT JOIN mst.m_outlet o ON o.outlet_id = sa.outlet_id AND o.cust_id = sa.cust_id AND o.is_del = false
		LEFT JOIN mst.m_salesman sm ON sm.emp_id = sa.emp_id AND sm.cust_id = sa.cust_id AND sm.is_del = false
		LEFT JOIN mst.m_employee emp ON emp.emp_id = sa.emp_id AND emp.cust_id = sa.cust_id
		WHERE sa.survey_answer_id = $1 AND s.cust_id = $2 AND sa.is_del = false`
	err := r.DB.Get(&row, query, surveyAnswerID, custID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.SurveyReportDetailRow{}, err
		}
		return model.SurveyReportDetailRow{}, err
	}
	return row, nil
}

func (r *surveyReportRepositoryImpl) FindQuestionDetails(surveyAnswerID int64, custID string) ([]model.SurveyReportQuestionRow, error) {
	var rows []model.SurveyReportQuestionRow
	query := `SELECT
			sad.survey_answer_detail_id,
			sad.question_template_id,
			COALESCE(qt.survey_template_id, 0) AS survey_template_id,
			COALESCE(qt.question, '') AS question,
			COALESCE(qt.input_type, sad.input_type, '') AS input_type,
			COALESCE(qt.answer_type, sad.answer_type, '') AS answer_type,
			COALESCE(sad.seq, 0) AS seq,
			COALESCE(sad.is_answered, false) AS is_answered,
			sad.free_text_answer,
			sad.photo_path
		FROM mst.survey_answer_detail sad
		LEFT JOIN mst.question_template qt ON qt.question_template_id = sad.question_template_id
		WHERE sad.survey_answer_id = $1 AND sad.cust_id = $2 AND sad.is_del = false
		ORDER BY sad.seq ASC, sad.survey_answer_detail_id ASC`
	err := r.DB.Select(&rows, query, surveyAnswerID, custID)
	return rows, err
}

func (r *surveyReportRepositoryImpl) FindQuestionOptions(questionTemplateIDs []int64) ([]model.SurveyReportQuestionOptionRow, error) {
	if len(questionTemplateIDs) == 0 {
		return []model.SurveyReportQuestionOptionRow{}, nil
	}

	query, args, err := sqlx.In(`SELECT question_template_id, q_option_template_id, option
		FROM mst.m_q_option_template
		WHERE is_del = false AND question_template_id IN (?)
		ORDER BY seq ASC, q_option_template_id ASC`, questionTemplateIDs)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)

	var rows []model.SurveyReportQuestionOptionRow
	err = r.DB.Select(&rows, query, args...)
	return rows, err
}

func (r *surveyReportRepositoryImpl) FindSelectedOptions(detailIDs []int64, custID string) ([]model.SurveyReportSelectedOptionRow, error) {
	if len(detailIDs) == 0 {
		return []model.SurveyReportSelectedOptionRow{}, nil
	}

	query, args, err := sqlx.In(`SELECT survey_answer_detail_id, survey_answer_option_id, q_option_template_id, option_label
		FROM mst.survey_answer_option
		WHERE is_del = false AND cust_id = ? AND survey_answer_detail_id IN (?)
		ORDER BY survey_answer_option_id ASC`, custID, detailIDs)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)

	var rows []model.SurveyReportSelectedOptionRow
	err = r.DB.Select(&rows, query, args...)
	return rows, err
}

func (r *surveyReportRepositoryImpl) FindAnswerFiles(detailIDs []int64, custID string) ([]model.SurveyReportAnswerFileRow, error) {
	if len(detailIDs) == 0 {
		return []model.SurveyReportAnswerFileRow{}, nil
	}

	query, args, err := sqlx.In(`SELECT survey_answer_detail_id, survey_answer_files AS survey_answer_files_id, file_name, file_key, media_category, file_size
		FROM mst.survey_answer_files
		WHERE cust_id = ? AND survey_answer_detail_id IN (?)
		ORDER BY survey_answer_files ASC`, custID, detailIDs)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)

	var rows []model.SurveyReportAnswerFileRow
	err = r.DB.Select(&rows, query, args...)
	return rows, err
}

func (r *surveyReportRepositoryImpl) FindExportRows(filter entity.SurveyReportQueryFilter) ([]model.SurveyReportExportRow, error) {
	query, args, err := r.buildExportRowsQuery(filter)
	if err != nil {
		return nil, err
	}

	var rows []model.SurveyReportExportRow
	err = r.DB.Select(&rows, query, args...)
	return rows, err
}

func (r *surveyReportRepositoryImpl) buildExportRowsQuery(filter entity.SurveyReportQueryFilter) (string, []interface{}, error) {
	extraJoins := `
		JOIN mst.survey_answer_detail sad ON sad.survey_answer_id = sa.survey_answer_id AND sad.cust_id = sa.cust_id AND sad.is_del = false
		JOIN mst.question_template qt ON qt.question_template_id = sad.question_template_id`
	baseQuery, args, err := r.buildListBaseQueryWithExtraJoins(filter, extraJoins)
	if err != nil {
		return "", nil, err
	}

	query := `SELECT
		sa.answer_date AS survey_date,
		s.survey_title,
		COALESCE(a.area_code, '') AS area_code,
		COALESCE(a.area_name, '') AS area_name,
		COALESCE(md.distributor_code, '') AS distributor_code,
		COALESCE(md.distributor_name, '') AS distributor_name,
		COALESCE(emp.emp_code, '') AS emp_code,
		COALESCE(emp.emp_name, sm.sales_name, '') AS emp_name,
		COALESCE(o.outlet_principal_code, o.outlet_code, '') AS outlet_code,
		COALESCE(o.outlet_name, '') AS outlet_name,
		COALESCE(qt.question, '') AS question,
		CASE
			WHEN COALESCE(qt.answer_type, sad.answer_type, '') = 'Free Text' THEN COALESCE(sad.free_text_answer, '')
			ELSE COALESCE((
				SELECT string_agg(sao.option_label, ', ' ORDER BY sao.survey_answer_option_id)
				FROM mst.survey_answer_option sao
				WHERE sao.survey_answer_detail_id = sad.survey_answer_detail_id
					AND sao.cust_id = sa.cust_id
					AND sao.is_del = false
			), '')
		END AS answer,
		COALESCE((SELECT COALESCE(NULLIF(saf.file_key, ''), NULLIF(saf.file_name, ''), '') FROM mst.survey_answer_files saf WHERE saf.survey_answer_detail_id = sad.survey_answer_detail_id AND saf.cust_id = sa.cust_id ORDER BY saf.survey_answer_files LIMIT 1 OFFSET 0), '') AS attachment_1,
		COALESCE((SELECT COALESCE(NULLIF(saf.file_key, ''), NULLIF(saf.file_name, ''), '') FROM mst.survey_answer_files saf WHERE saf.survey_answer_detail_id = sad.survey_answer_detail_id AND saf.cust_id = sa.cust_id ORDER BY saf.survey_answer_files LIMIT 1 OFFSET 1), '') AS attachment_2,
		COALESCE((SELECT COALESCE(NULLIF(saf.file_key, ''), NULLIF(saf.file_name, ''), '') FROM mst.survey_answer_files saf WHERE saf.survey_answer_detail_id = sad.survey_answer_detail_id AND saf.cust_id = sa.cust_id ORDER BY saf.survey_answer_files LIMIT 1 OFFSET 2), '') AS attachment_3
		` + baseQuery + `
		ORDER BY sa.answer_date DESC, s.survey_title ASC, emp.emp_code ASC, o.outlet_name ASC, sad.seq ASC`

	return query, args, nil
}

func buildSurveyReportSortClause(rawSort string) string {
	allowedFields := map[string]string{
		"created_date":     "COALESCE(sa.created_at, sa.answer_date)",
		"answer_date":      "sa.answer_date",
		"survey_title":     "s.survey_title",
		"area_name":        "a.area_name",
		"outlet_name":      "o.outlet_name",
		"salesman_name":    "COALESCE(sm.sales_name, emp.emp_name, '')",
		"emp_name":         "emp.emp_name",
		"distributor_name": "md.distributor_name",
	}

	field := "COALESCE(sa.created_at, sa.answer_date)"
	order := "DESC"
	parts := strings.Split(rawSort, ":")
	if len(parts) == 2 {
		if mappedField, ok := allowedFields[strings.TrimSpace(parts[0])]; ok {
			field = mappedField
		}
		if strings.EqualFold(strings.TrimSpace(parts[1]), "asc") {
			order = "ASC"
		}
	}

	return field + " " + order + ", sa.survey_answer_id DESC"
}

func shiftPositionalPlaceholders(query string, offset int) string {
	if offset == 0 {
		return query
	}

	var builder strings.Builder
	builder.Grow(len(query))
	for index := 0; index < len(query); index++ {
		if query[index] != '$' || index+1 >= len(query) || query[index+1] < '0' || query[index+1] > '9' {
			builder.WriteByte(query[index])
			continue
		}

		placeholderStart := index + 1
		placeholderEnd := placeholderStart
		placeholderNumber := 0
		for placeholderEnd < len(query) && query[placeholderEnd] >= '0' && query[placeholderEnd] <= '9' {
			placeholderNumber = placeholderNumber*10 + int(query[placeholderEnd]-'0')
			placeholderEnd++
		}

		builder.WriteString(fmt.Sprintf("$%d", placeholderNumber+offset))
		index = placeholderEnd - 1
	}
	return builder.String()
}
