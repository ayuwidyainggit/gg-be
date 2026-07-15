package model

import "time"

type SurveyAnswerListItem struct {
	SurveyID           int64  `json:"survey_id"`
	SurveyTitle        string `json:"survey_title"`
	AnswerFrequency    string `json:"answer_frequency"`
	ResponseType       string `json:"response_type"`
	EffectiveDateStart string `json:"effective_date_start"`
	EffectiveDateEnd   string `json:"effective_date_end"`
	QuestionTotal      int    `json:"question_total"`
	SurveyDetailID     int64  `json:"survey_detail_id"`
	SurveyTemplateID   int64  `json:"survey_template_id"`
	TemplateCode       string `json:"template_code"`
	TemplateTitle      string `json:"template_title"`
	UseImage           bool   `json:"use_image"`
	SurveyAnswerID     int64  `json:"survey_answer_id"`
}

type SurveyResponse struct {
	SurveyID           int64   `json:"survey_id"`
	SurveyTitle        string  `json:"survey_title"`
	AnswerFrequency    string  `json:"answer_frequency"`
	ResponseType       string  `json:"response_type"`
	EffectiveDateStart string  `json:"effective_date_start"`
	EffectiveDateEnd   string  `json:"effective_date_end"`
	QuestionTotal      int     `json:"question_total"`
	SurveyAnswerID     *int64  `json:"survey_answer_id"`
	ReSurvey           bool    `json:"re_survey"`
	TakeSurvey         bool    `json:"take_survey"`
}

type SurveyDetailResponse struct {
	SurveyID           int64               `json:"survey_id"`
	AreaID             int64               `json:"area_id"`
	SurveyTitle        string              `json:"survey_title"`
	AnswerFrequency    string              `json:"answer_frequency"`
	ResponseType       string              `json:"response_type"`
	EffectiveDateStart string              `json:"effective_date_start"`
	EffectiveDateEnd   string              `json:"effective_date_end"`
	QuestionTotal      int                 `json:"question_total"`
	TakeSurvey         bool                `json:"take_survey"`
	TemplateData       *SurveyTemplateData `json:"template_data"`
}

type SurveyTemplateData struct {
	MSurveyDetailID  int64            `json:"m_survey_detail_id"`
	SurveyTemplateID int64            `json:"survey_template_id"`
	TemplateCode     string           `json:"template_code"`
	TemplateTitle    string           `json:"template_title"`
	QuestionTotal    int              `json:"question_total"`
	UseImage         bool             `json:"use_image"`
	Questions        []SurveyQuestion `json:"questions"`
}

type SurveyQuestion struct {
	QuestionTemplateID int64                  `json:"question_template_id"`
	Question           string                 `json:"question"`
	InputType          string                 `json:"input_type"`
	AnswerType         string                 `json:"answer_type"`
	Seq                int                    `json:"seq"`
	UseImage           bool                   `json:"use_image"`
	IsAnswered         bool                   `json:"is_answered"`
	FreeTextAnswer     *string                `json:"free_text_answer"`
	PhotoPath          *string                `json:"photo_path"`
	Options            []SurveyQuestionOption `json:"options" gorm:"foreignKey:QOptionTemplateID;references:QuestionTemplateID"`
	AnswerFiles        []SurveyAnswerFile     `json:"answer_files" gorm:"foreignKey:SurveyAnswerDetailID;references:QuestionTemplateID"`
	AnswerOptions      []SurveyAnswerOption   `json:"answer_options" gorm:"foreignKey:SurveyAnswerDetailID;references:QuestionTemplateID"`
}

type SurveyQuestionOption struct {
	QOptionTemplateID int64  `json:"q_option_template_id"`
	Option            string `json:"option"`
}

type SurveyAnswer struct {
	CustID           string               `gorm:"column:cust_id" json:"cust_id"`
	SurveyAnswerID   int64                `gorm:"column:survey_answer_id;primaryKey;autoIncrement" json:"survey_answer_id"`
	SurveyTemplateID int64                `gorm:"column:survey_template_id" json:"survey_template_id"`
	SurveyID         int64                `gorm:"column:survey_id" json:"survey_id"`
	EmpID            int64                `gorm:"column:emp_id" json:"emp_id"`
	OutletID         int64                `gorm:"column:outlet_id" json:"outlet_id"`
	AreaID           *int64               `gorm:"column:area_id" json:"area_id"`
	AnswerDate       *time.Time           `gorm:"column:answer_date" json:"answer_date"`
	Status           string               `gorm:"column:status" json:"status"`
	CreatedBy        int64                `gorm:"column:created_by" json:"created_by"`
	CreatedAt        *time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedBy        int64                `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt        *time.Time           `gorm:"column:updated_at" json:"updated_at"`
	IsDel            bool                 `gorm:"column:is_del" json:"is_del"`
	DeletedBy        *int64               `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt        *time.Time           `gorm:"column:deleted_at" json:"deleted_at"`
	Details          []SurveyAnswerDetail `gorm:"foreignKey:SurveyAnswerID;references:SurveyAnswerID" json:"details"`
}

func (SurveyAnswer) TableName() string {
	return "mst.survey_answer"
}

type SurveyAnswerDetail struct {
	CustID               string               `gorm:"column:cust_id" json:"cust_id"`
	SurveyAnswerDetailID int64                `gorm:"column:survey_answer_detail_id;primaryKey;autoIncrement" json:"survey_answer_detail_id"`
	SurveyAnswerID       int64                `gorm:"column:survey_answer_id" json:"survey_answer_id"`
	QuestionTemplateID   int64                `gorm:"column:question_template_id" json:"question_template_id"`
	InputType            string               `gorm:"column:input_type" json:"input_type"`
	AnswerType           string               `gorm:"column:answer_type" json:"answer_type"`
	Seq                  int                  `gorm:"column:seq" json:"seq"`
	IsAnswered           bool                 `gorm:"column:is_answered" json:"is_answered"`
	FreeTextAnswer       *string              `gorm:"column:free_text_answer" json:"free_text_answer"`
	PhotoPath            *string              `gorm:"column:photo_path" json:"photo_path"`
	CreatedBy            int64                `gorm:"column:created_by" json:"created_by"`
	CreatedAt            *time.Time           `gorm:"column:created_at" json:"created_at"`
	UpdatedBy            int64                `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt            *time.Time           `gorm:"column:updated_at" json:"updated_at"`
	IsDel                bool                 `gorm:"column:is_del" json:"is_del"`
	DeletedBy            *int64               `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt            *time.Time           `gorm:"column:deleted_at" json:"deleted_at"`
	Files                []SurveyAnswerFile   `gorm:"foreignKey:SurveyAnswerDetailID;references:SurveyAnswerDetailID" json:"files"`
	Options              []SurveyAnswerOption `gorm:"foreignKey:SurveyAnswerDetailID;references:SurveyAnswerDetailID" json:"options"`
}

func (SurveyAnswerDetail) TableName() string {
	return "mst.survey_answer_detail"
}

type SurveyAnswerFile struct {
	CustID               string `gorm:"column:cust_id" json:"cust_id"`
	SurveyAnswerFilesID  int64  `gorm:"column:survey_answer_files;primaryKey;autoIncrement" json:"survey_answer_files_id"`
	SurveyAnswerDetailID int64  `gorm:"column:survey_answer_detail_id" json:"survey_answer_detail_id"`
	FileName             string `gorm:"column:file_name" json:"file_name"`
	FileData             []byte `gorm:"column:file_data" json:"file_data,omitempty"`
	FileKey              string `gorm:"column:file_key" json:"file_key"`
	MediaCategory        string `gorm:"column:media_category" json:"media_category"`
	FileSize             *int64 `gorm:"column:file_size" json:"file_size"`
}

func (SurveyAnswerFile) TableName() string {
	return "mst.survey_answer_files"
}

type SurveyAnswerOption struct {
	CustID               string     `gorm:"column:cust_id" json:"cust_id"`
	SurveyAnswerOptionID int64      `gorm:"column:survey_answer_option_id;primaryKey;autoIncrement" json:"survey_answer_option_id"`
	SurveyAnswerDetailID int64      `gorm:"column:survey_answer_detail_id" json:"survey_answer_detail_id"`
	QOptionTemplateID    int64      `gorm:"column:q_option_template_id" json:"q_option_template_id"`
	OptionLabel          string     `gorm:"column:option_label" json:"option_label"`
	CreatedBy            int64      `gorm:"column:created_by" json:"created_by"`
	CreatedAt            *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy            int64      `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt            *time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsDel                bool       `gorm:"column:is_del" json:"is_del"`
	DeletedBy            *int64     `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt            *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (SurveyAnswerOption) TableName() string {
	return "mst.survey_answer_option"
}

type GetSurveySubmitted struct {
	SurveyAnswerID     int64               `json:"survey_answer_id"`
	SurveyID           int64               `json:"survey_id"`
	SurveyTitle        string              `json:"survey_title"`
	AnswerFrequency    string              `json:"answer_frequency"`
	ResponseType       string              `json:"response_type"`
	EffectiveDateStart time.Time           `json:"effective_date_start"`
	EffectiveDateEnd   time.Time           `json:"effective_date_end"`
	QuestionTotal      int                 `json:"question_total"`
	TakeSurvey         bool                `json:"take_survey"`
	TemplateData       *SurveyTemplateData `json:"template_data"`
}
