package entity

type SurveyQueryFilter struct {
	OutletID      int64  `query:"outlet_id" validate:"required"`
	EmpID         int64  `query:"-"`
	Q             string `query:"q"`
	Page          int64  `query:"page"`
	Limit         int64  `query:"limit"`
	Sort          string `query:"sort"`
	DistributorID int64  `query:"-"`
}

type SurveyAnswerListFilter struct {
	SurveyID int64 `query:"survey_id" validate:"required"`
	OutletID int64 `query:"outlet_id" validate:"required"`
	Page     int   `query:"page"`
	Limit    int   `query:"limit"`
}

type SubmitSurveyRequest struct {
	SurveyID         int64            `json:"survey_id" validate:"required"`
	EmpID            int64            `json:"emp_id" validate:"required"`
	OutletID         int64            `json:"outlet_id" validate:"required"`
	SurveyTemplateID int64            `json:"survey_template_id" validate:"required"`
	Questions        []SurveyQuestion `json:"questions" validate:"required,dive"`
	CustID           string           `json:"-"`
	UserID           int64            `json:"-"`
	DistributorID    int64            `json:"-"`
}

type SurveyQuestion struct {
	QuestionTemplateID int64                  `json:"question_template_id" validate:"required"`
	InputType          string                 `json:"input_type" validate:"required"`
	AnswerType         string                 `json:"answer_type" validate:"required"`
	Seq                int                    `json:"seq" validate:"required"`
	IsAnswered         bool                   `json:"is_answered"`
	FreeTextAnswer     *string                `json:"free_text_answer"`
	Files              []SurveyAnswerFile     `json:"files,omitempty"`
	Options            []SurveyQuestionOption `json:"options,omitempty"`
}

type SurveyAnswerFile struct {
	ClientFileName string `json:"client_file_name"`
}

type SurveyQuestionOption struct {
	QOptionTemplateID int64  `json:"q_option_template_id"`
	Option            string `json:"option"`
}
