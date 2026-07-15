package entity

import "time"

// ============ Response Structs ============

// SurveyTemplateListResponse for List endpoint response
type SurveyTemplateListResponse struct {
	SurveyTemplateId int        `json:"survey_template_id"`
	TemplateCode     string     `json:"template_code"`
	TemplateTitle    string     `json:"template_title"`
	QuestionTotal    int        `json:"question_total"`
	UseImage         bool       `json:"use_image"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        *time.Time `json:"created_at"`
}

// SurveyTemplateDetailResponse for Detail endpoint response
type SurveyTemplateDetailResponse struct {
	SurveyTemplateId int                        `json:"survey_template_id"`
	TemplateCode     string                     `json:"template_code"`
	TemplateTitle    string                     `json:"template_title"`
	QuestionTotal    int                        `json:"question_total"`
	UseImage         bool                       `json:"use_image"`
	IsActive         bool                       `json:"is_active"`
	CreatedAt        *time.Time                 `json:"created_at"`
	QuestionTemplate []QuestionTemplateResponse `json:"question_template"`
}

// QuestionTemplateResponse for nested question in detail response
type QuestionTemplateResponse struct {
	QuestionTemplateId int                       `json:"question_template_id"`
	SurveyTemplateId   int                       `json:"survey_template_id"`
	Question           string                    `json:"question"`
	InputType          string                    `json:"input_type"`
	AnswerType         string                    `json:"answer_type"`
	UseImage           bool                      `json:"use_image"`
	MQOptionTemplate   []QOptionTemplateResponse `json:"m_q_option_template"`
}

// QOptionTemplateResponse for nested option in question response
type QOptionTemplateResponse struct {
	QOptionTemplateId int    `json:"q_option_template_id"`
	Option            string `json:"option"`
}

// ============ Request Structs ============

// CreateSurveyTemplateBody for POST request
type CreateSurveyTemplateBody struct {
	TemplateTitle string                    `json:"template_title" validate:"required,max=150"`
	QuestionTotal int                       `json:"question_total" validate:"required,min=0"`
	UseImage      bool                      `json:"use_image"`
	IsActive      bool                      `json:"is_active"`
	Question      []QuestionTemplateRequest `json:"question" validate:"required,dive"`
	CustId        string                    `json:"-"`
	CreatedBy     int64                     `json:"-"`
}

// UpdateSurveyTemplateBody for PUT request
type UpdateSurveyTemplateBody struct {
	TemplateTitle string                    `json:"template_title" validate:"required,max=150"`
	QuestionTotal int                       `json:"question_total" validate:"required,min=0"`
	UseImage      bool                      `json:"use_image"`
	IsActive      bool                      `json:"is_active"`
	Question      []QuestionTemplateRequest `json:"question" validate:"required,dive"`
	CustId        string                    `json:"-"`
	UpdatedBy     int64                     `json:"-"`
}

// QuestionTemplateRequest for nested question in create/update request
type QuestionTemplateRequest struct {
	QuestionTemplateId int                      `json:"question_template_id,omitempty"` // Optional for update
	Question           string                   `json:"question" validate:"required,max=225"`
	InputType          string                   `json:"input_type" validate:"required,oneof=textfield dropdown radiobutton toggle checkbox"`
	AnswerType         string                   `json:"answer_type" validate:"required,oneof=Single Multiple 'Free Text'"`
	UseImage           bool                     `json:"use_image"`
	QOption            []QOptionTemplateRequest `json:"q_option"`
}

// QOptionTemplateRequest for nested option in question request
type QOptionTemplateRequest struct {
	QOptionTemplateId int    `json:"q_option_template_id,omitempty"` // Optional for update
	Option            string `json:"option" validate:"max=225"`
}

// ============ Params & Filter Structs ============

// SurveyTemplateParams for path parameter
type SurveyTemplateParams struct {
	SurveyTemplateId int `params:"survey_template_id" validate:"required"`
}

// SurveyTemplateQueryFilter for List endpoint query parameters
type SurveyTemplateQueryFilter struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit" validate:"required"`
	Query  string `query:"q"`
	Sort   string `query:"sort"`
	Status *int   `query:"status"` // 1=active, 0=inactive, nil=all
}
