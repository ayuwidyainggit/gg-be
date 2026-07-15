package repository

import (
	"context"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type SurveyRepository interface {
	Count(ctx context.Context, filter entity.SurveyQueryFilter) (int64, error)
	List(ctx context.Context, filter entity.SurveyQueryFilter) ([]model.SurveyResponse, error)
	Detail(ctx context.Context, surveyId int64) (*model.SurveyDetailResponse, error)
	SubmitSurvey(ctx context.Context, answer *model.SurveyAnswer) error
	GetSubmittedSurvey(ctx context.Context, surveyID, surveyAnswerID int64) (*model.SurveyDetailResponse, error)
	CheckExistingSubmission(ctx context.Context, surveyID, empID, outletID int64) (bool, error)
	CountSurveyAnswer(ctx context.Context, surveyID, outletID int64) (int64, error)
	ListSurveyAnswer(ctx context.Context, surveyID, outletID int64, limit int, offset int) ([]model.SurveyAnswerListItem, error)
}

type surveyRepository struct {
	db *gorm.DB
}

func NewSurveyRepository(db *gorm.DB) SurveyRepository {
	return &surveyRepository{
		db: db,
	}
}

func (r *surveyRepository) Count(ctx context.Context, filter entity.SurveyQueryFilter) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM mst.m_survey ms
INNER JOIN mst.m_survey_detail msd ON msd.survey_id = ms.survey_id AND msd.is_del = false
INNER JOIN mst.m_survey_template mst_tpl ON msd.survey_template_id = mst_tpl.survey_template_id AND mst_tpl.is_del = false
LEFT JOIN (
	SELECT survey_id, COUNT(CASE WHEN outlet_id = ? THEN 1 END) AS is_matched
	FROM mst.m_survey_outlet WHERE is_del = false GROUP BY survey_id
) mso ON mso.survey_id = ms.survey_id
LEFT JOIN (
	SELECT survey_id, COUNT(CASE WHEN salesman_id = ? THEN 1 END) AS is_matched
	FROM mst.m_survey_salesman WHERE is_del = false GROUP BY survey_id
) mss ON mss.survey_id = ms.survey_id
LEFT JOIN (
	SELECT survey_id, COUNT(*) AS total_answers,
		COUNT(CASE WHEN created_at::date = CURRENT_DATE THEN 1 END) AS today_answers
	FROM mst.survey_answer WHERE outlet_id = ? AND is_del = false GROUP BY survey_id
) msa ON msa.survey_id = ms.survey_id
WHERE ms.is_del = false
	AND ms.status = 1
	AND ms.efective_date_start <= CURRENT_DATE
	AND ms.efective_date_end >= CURRENT_DATE
	AND (
		(mso.survey_id IS NULL AND mss.survey_id IS NULL)
		OR (mso.is_matched > 0 AND mss.survey_id IS NULL)
		OR (mso.survey_id IS NULL AND mss.is_matched > 0)
		OR (mso.is_matched > 0 AND mss.is_matched > 0)
	)`

	args := []interface{}{filter.OutletID, filter.EmpID, filter.OutletID}

	if filter.Q != "" {
		query += " AND ms.survey_title ILIKE ?"
		args = append(args, "%"+filter.Q+"%")
	}

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *surveyRepository) List(ctx context.Context, filter entity.SurveyQueryFilter) ([]model.SurveyResponse, error) {
	var records []model.SurveyResponse

	query := `SELECT
	ms.survey_id,
	ms.survey_title,
	ms.answer_frequency,
	ms.response_type,
	ms.efective_date_start,
	ms.efective_date_end,
	mst_tpl.question_total,
	CASE
		WHEN ms.answer_frequency = 'One Time' THEN msa.first_answer_id
		WHEN ms.answer_frequency = 'Multiple Times, Different Day' THEN msa.today_answer_id
		ELSE NULL
	END AS survey_answer_id,
	CASE
		WHEN ms.answer_frequency = 'One Time'
			AND msa.total_answers > 0 THEN false
		WHEN ms.answer_frequency = 'Multiple Times, Different Day'
			AND msa.today_answers > 0 THEN false
		ELSE true
	END AS re_survey,
  CASE
       WHEN msa.today_answers > 0 THEN true
       WHEN ms.answer_frequency = 'One Time'
            AND msa.total_answers > 0 THEN true
       ELSE false
   END AS take_survey
FROM mst.m_survey ms
INNER JOIN mst.m_survey_detail msd ON msd.survey_id = ms.survey_id AND msd.is_del = false
INNER JOIN mst.m_survey_template mst_tpl ON msd.survey_template_id = mst_tpl.survey_template_id AND mst_tpl.is_del = false
LEFT JOIN (
	SELECT survey_id, COUNT(CASE WHEN outlet_id = ? THEN 1 END) AS is_matched
	FROM mst.m_survey_outlet WHERE is_del = false GROUP BY survey_id
) mso ON mso.survey_id = ms.survey_id
LEFT JOIN (
	SELECT survey_id, COUNT(CASE WHEN salesman_id = ? THEN 1 END) AS is_matched
	FROM mst.m_survey_salesman WHERE is_del = false GROUP BY survey_id
) mss ON mss.survey_id = ms.survey_id
LEFT JOIN (
	SELECT
		survey_id,
		COUNT(*) AS total_answers,
		COUNT(CASE WHEN created_at::date = CURRENT_DATE THEN 1 END) AS today_answers,
		MIN(survey_answer_id) AS first_answer_id,
		MAX(CASE WHEN created_at::date = CURRENT_DATE THEN survey_answer_id END) AS today_answer_id
	FROM mst.survey_answer
	WHERE outlet_id = ? AND is_del = false
	GROUP BY survey_id
) msa ON msa.survey_id = ms.survey_id
WHERE ms.is_del = false
	AND ms.status = 1
	AND ms.efective_date_start <= CURRENT_DATE
	AND ms.efective_date_end >= CURRENT_DATE
	AND (
		(mso.survey_id IS NULL AND mss.survey_id IS NULL)
		OR (mso.is_matched > 0 AND mss.survey_id IS NULL)
		OR (mso.survey_id IS NULL AND mss.is_matched > 0)
		OR (mso.is_matched > 0 AND mss.is_matched > 0)
	)`

	args := []interface{}{filter.OutletID, filter.EmpID, filter.OutletID}

	// Apply filter
	if filter.Q != "" {
		query += " AND ms.survey_title ILIKE ?"
		args = append(args, "%"+filter.Q+"%")
	}

	// Apply sorting
	if filter.Sort != "" {
		// example: "created_date:asc" -> "created_date asc"
		sortConfig := strings.Split(filter.Sort, ":")
		if len(sortConfig) == 2 {
			query += " ORDER BY " + sortConfig[0] + " " + sortConfig[1]
		} else {
			query += " ORDER BY " + filter.Sort // fallback
		}
	} else {
		// Default sorting
		query += " ORDER BY ms.survey_id DESC"
	}

	// Apply pagination
	if filter.Limit > 0 && filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += " LIMIT ? OFFSET ?"
		args = append(args, filter.Limit, offset)
	} else if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *surveyRepository) Detail(ctx context.Context, surveyId int64) (*model.SurveyDetailResponse, error) {
	var detail model.SurveyDetailResponse
	var templateData model.SurveyTemplateData

	err := r.db.WithContext(ctx).Table("mst.m_survey ms").
		Select(`
			ms.survey_id,
			ms.survey_title,
			ms.answer_frequency,
			ms.response_type,
			ms.efective_date_start,
			ms.efective_date_end,
			mst.question_total,
			CASE WHEN msa.survey_answer_id IS NOT NULL THEN true ELSE false END as take_survey,
			msd.survey_detail_id,
			mst.survey_template_id,
			mst.template_code,
			mst.template_title,
			mst.use_image
		`).
		Joins("inner join mst.m_survey_detail msd on msd.survey_id = ms.survey_id").
		Joins("inner join mst.m_survey_template mst on msd.survey_template_id = mst.survey_template_id").
		Joins("left join mst.survey_answer msa on msa.survey_id = ms.survey_id").
		Where("ms.survey_id = ?", surveyId).
		Row().Scan(
		&detail.SurveyID,
		&detail.SurveyTitle,
		&detail.AnswerFrequency,
		&detail.ResponseType,
		&detail.EffectiveDateStart,
		&detail.EffectiveDateEnd,
		&detail.QuestionTotal,
		&detail.TakeSurvey,
		&templateData.MSurveyDetailID,
		&templateData.SurveyTemplateID,
		&templateData.TemplateCode,
		&templateData.TemplateTitle,
		&templateData.UseImage,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, err
		}
		return nil, err
	}

	templateData.QuestionTotal = detail.QuestionTotal
	detail.TemplateData = &templateData
	detail.TemplateData.Questions = make([]model.SurveyQuestion, 0)

	var questions []model.SurveyQuestion
	err = r.db.WithContext(ctx).Table("mst.question_template mqt").
		Select(`
			mqt.question_template_id,
			mqt.question,
			mqt.answer_type,
			mqt.seq,
			mqt.use_image,
		  mqt.input_type
		`).
		Where("mqt.survey_template_id = ?", templateData.SurveyTemplateID).
		Where("mqt.is_del = ?", false).
		Order("mqt.seq asc").
		Find(&questions).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(questions) > 0 {
		var questionIds []int64
		for _, q := range questions {
			questionIds = append(questionIds, q.QuestionTemplateID)
		}

		type questionOption struct {
			QuestionTemplateID int64
			QOptionTemplateID  int64
			Option             string
		}
		var options []questionOption

		err = r.db.WithContext(ctx).Table("mst.m_q_option_template mqo").
			Select(`
				mqo.question_template_id,
				mqo.q_option_template_id,
				mqo.option
			`).
			Where("mqo.question_template_id IN ?", questionIds).
			Find(&options).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		optionsMap := make(map[int64][]model.SurveyQuestionOption)
		for _, opt := range options {
			optionsMap[opt.QuestionTemplateID] = append(optionsMap[opt.QuestionTemplateID], model.SurveyQuestionOption{
				QOptionTemplateID: opt.QOptionTemplateID,
				Option:            opt.Option,
			})
		}

		for i := range questions {
			if opts, ok := optionsMap[questions[i].QuestionTemplateID]; ok {
				questions[i].Options = opts
			} else {
				questions[i].Options = make([]model.SurveyQuestionOption, 0)
			}
		}

		detail.TemplateData.Questions = questions
	}

	return &detail, nil
}

func (r *surveyRepository) SubmitSurvey(ctx context.Context, answer *model.SurveyAnswer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

func (r *surveyRepository) CheckExistingSubmission(ctx context.Context, surveyID, empID, outletID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("mst.survey_answer").
		Where("survey_id = ? AND emp_id = ? AND outlet_id = ?", surveyID, empID, outletID).
		Count(&count).Error
	return count > 0, err
}

func (r *surveyRepository) GetSubmittedSurvey(ctx context.Context, surveyID, surveyAnswerID int64) (*model.SurveyDetailResponse, error) {
	var detail model.SurveyDetailResponse
	var templateData model.SurveyTemplateData

	err := r.db.WithContext(ctx).Raw(`
SELECT
	ms.survey_id,
	ms.survey_title,
	ms.answer_frequency,
	ms.response_type,
	ms.efective_date_start,
	ms.efective_date_end,
	mst_tpl.question_total,
	CASE
		WHEN msa.survey_answer_id IS NOT NULL THEN true
		ELSE false
	END AS take_survey,
	msd.survey_detail_id,
	mst_tpl.survey_template_id,
	mst_tpl.template_code,
	mst_tpl.template_title,
	mst_tpl.use_image
FROM mst.m_survey ms
INNER JOIN mst.m_survey_detail msd
	ON msd.survey_id = ms.survey_id
INNER JOIN mst.m_survey_template mst_tpl
	ON msd.survey_template_id = mst_tpl.survey_template_id
LEFT JOIN mst.survey_answer msa
	ON msa.survey_id = ms.survey_id
WHERE ms.survey_id = ?
AND msa.survey_answer_id = ?
`, surveyID, surveyAnswerID).Row().Scan(
		&detail.SurveyID,
		&detail.SurveyTitle,
		&detail.AnswerFrequency,
		&detail.ResponseType,
		&detail.EffectiveDateStart,
		&detail.EffectiveDateEnd,
		&detail.QuestionTotal,
		&detail.TakeSurvey,
		&templateData.MSurveyDetailID,
		&templateData.SurveyTemplateID,
		&templateData.TemplateCode,
		&templateData.TemplateTitle,
		&templateData.UseImage,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, err
		}
		return nil, err
	}

	templateData.QuestionTotal = detail.QuestionTotal
	detail.TemplateData = &templateData
	detail.TemplateData.Questions = make([]model.SurveyQuestion, 0)

	var questions []model.SurveyQuestion
	err = r.db.WithContext(ctx).Table("mst.question_template mqt").
		Select(`
			mqt.question_template_id,
			mqt.question,
			mqt.answer_type,
			mqt.seq,
			mqt.use_image,
		  mqt.input_type
		`).
		Where("mqt.survey_template_id = ?", templateData.SurveyTemplateID).
		Where("mqt.is_del = ?", false).
		Order("mqt.seq asc").
		Find(&questions).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(questions) > 0 {
		questionTemplateIDs := make([]int64, 0)
		for _, sq := range questions {
			if sq.QuestionTemplateID != 0 {
				questionTemplateIDs = append(questionTemplateIDs, sq.QuestionTemplateID)
			}
		}

		questionIds := make([]int64, 0)
		for _, q := range questions {
			questionIds = append(questionIds, q.QuestionTemplateID)
		}

		// Fetch options for each question template
		type questionOption struct {
			QuestionTemplateID int64
			QOptionTemplateID  int64
			Option             string
		}
		var options []questionOption

		err = r.db.WithContext(ctx).Table("mst.m_q_option_template mqo").
			Select(`
				mqo.question_template_id,
				mqo.q_option_template_id,
				mqo.option
			`).
			Where("mqo.question_template_id IN ?", questionIds).
			Where("mqo.is_del = false").
			Find(&options).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		optionsMap := make(map[int64][]model.SurveyQuestionOption)
		for _, opt := range options {
			optionsMap[opt.QuestionTemplateID] = append(optionsMap[opt.QuestionTemplateID], model.SurveyQuestionOption{
				QOptionTemplateID: opt.QOptionTemplateID,
				Option:            opt.Option,
			})
		}

		for i := range questions {
			if opts, ok := optionsMap[questions[i].QuestionTemplateID]; ok {
				questions[i].Options = opts
			} else {
				questions[i].Options = make([]model.SurveyQuestionOption, 0)
			}
		}

		detail.TemplateData.Questions = questions

		// Fetch answer details filtered by survey_answer_id
		var (
			answersMap    = make(map[int64]model.SurveyAnswerDetail)
			surveyAnswers []model.SurveyAnswer
		)

		err = r.db.
			Where("survey_answer_id = ?", surveyAnswerID).
			Preload("Details", "is_del = ? AND question_template_id in (?)", false, questionTemplateIDs).
			Find(&surveyAnswers).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		var (
			surveyAnswersDetailID []int64
			detailIDToQuestionID  = make(map[int64]int64)
		)
		for _, sa := range surveyAnswers {
			for _, sd := range sa.Details {
				answersMap[sd.QuestionTemplateID] = sd
				surveyAnswersDetailID = append(surveyAnswersDetailID, sd.SurveyAnswerDetailID)
				detailIDToQuestionID[sd.SurveyAnswerDetailID] = sd.QuestionTemplateID
			}
		}

		var answerFiles = make([]model.SurveyAnswerFile, 0)
		if len(surveyAnswersDetailID) != 0 {
			err = r.db.
				Where("survey_answer_detail_id in (?)", surveyAnswersDetailID).
				Find(&answerFiles).Error

			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			filesByDetailID := make(map[int64][]model.SurveyAnswerFile)
			for _, f := range answerFiles {
				filesByDetailID[f.SurveyAnswerDetailID] = append(filesByDetailID[f.SurveyAnswerDetailID], f)
			}
			for detailID, files := range filesByDetailID {
				if qID, ok := detailIDToQuestionID[detailID]; ok {
					if ans, ok := answersMap[qID]; ok {
						ans.Files = files
						answersMap[qID] = ans
					}
				}
			}

			var answerOptions = make([]model.SurveyAnswerOption, 0)
			err = r.db.
				Where("is_del = ? AND survey_answer_detail_id in (?)", false, surveyAnswersDetailID).
				Find(&answerOptions).Error

			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			if len(answerOptions) > 0 {
				optionsByDetailID := make(map[int64][]model.SurveyAnswerOption)
				for _, opt := range answerOptions {
					optionsByDetailID[opt.SurveyAnswerDetailID] = append(optionsByDetailID[opt.SurveyAnswerDetailID], opt)
				}
				for detailID, opts := range optionsByDetailID {
					if qID, ok := detailIDToQuestionID[detailID]; ok {
						if ans, ok := answersMap[qID]; ok {
							ans.Options = opts
							answersMap[qID] = ans
						}
					}
				}
			}
		}

		for i := range detail.TemplateData.Questions {
			// Map answers
			if ans, ok := answersMap[detail.TemplateData.Questions[i].QuestionTemplateID]; ok {
				detail.TemplateData.Questions[i].IsAnswered = ans.IsAnswered
				detail.TemplateData.Questions[i].FreeTextAnswer = ans.FreeTextAnswer
				detail.TemplateData.Questions[i].PhotoPath = ans.PhotoPath
				detail.TemplateData.Questions[i].AnswerFiles = ans.Files
				detail.TemplateData.Questions[i].AnswerOptions = ans.Options
			}

			// Map template options
			if opts, ok := optionsMap[detail.TemplateData.Questions[i].QuestionTemplateID]; ok {
				detail.TemplateData.Questions[i].Options = opts
			} else {
				detail.TemplateData.Questions[i].Options = make([]model.SurveyQuestionOption, 0)
			}
		}
	}

	return &detail, nil
}

func (r *surveyRepository) CountSurveyAnswer(ctx context.Context, surveyID, outletID int64) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Raw(`
SELECT
	COUNT(msa.survey_answer_id)
FROM mst.m_survey ms
INNER JOIN mst.m_survey_detail msd
	ON msd.survey_id = ms.survey_id
INNER JOIN mst.m_survey_template mst_tpl
	ON msd.survey_template_id = mst_tpl.survey_template_id
JOIN mst.survey_answer msa
	ON msa.survey_id = ms.survey_id
WHERE
	ms.survey_id = ?
	AND msa.outlet_id = ?
	AND msa.created_at::date = CURRENT_DATE
`, surveyID, outletID).Scan(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *surveyRepository) ListSurveyAnswer(ctx context.Context, surveyID, outletID int64, limit int, offset int) ([]model.SurveyAnswerListItem, error) {
	var records []model.SurveyAnswerListItem

	err := r.db.WithContext(ctx).Raw(`
SELECT
	ms.survey_id,
	ms.survey_title,
	ms.answer_frequency,
	ms.response_type,
	ms.efective_date_start,
	ms.efective_date_end,
	mst_tpl.question_total,
	msd.survey_detail_id,
	mst_tpl.survey_template_id,
	mst_tpl.template_code,
	mst_tpl.template_title,
	mst_tpl.use_image,
	msa.survey_answer_id
FROM mst.m_survey ms
INNER JOIN mst.m_survey_detail msd
	ON msd.survey_id = ms.survey_id
INNER JOIN mst.m_survey_template mst_tpl
	ON msd.survey_template_id = mst_tpl.survey_template_id
JOIN mst.survey_answer msa
	ON msa.survey_id = ms.survey_id
WHERE
	ms.survey_id = ?
	AND msa.outlet_id = ?
	AND msa.created_at::date = CURRENT_DATE
ORDER BY msa.created_at DESC
LIMIT ?
OFFSET ?
`, surveyID, outletID, limit, offset).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}
