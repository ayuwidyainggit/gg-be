package model

// SurveyDetail represents mst.m_survey_detail table (survey to template mapping)
type SurveyDetail struct {
	SurveyDetailId   int  `db:"survey_detail_id" json:"survey_detail_id"`
	SurveyId         int  `db:"survey_id" json:"survey_id"`
	SurveyTemplateId int  `db:"survey_template_id" json:"survey_template_id"`
	IsDel            bool `db:"is_del" json:"is_del"`
}
