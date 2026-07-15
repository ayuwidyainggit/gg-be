package model

import "time"

type SurveyReportListRow struct {
	SurveyAnswerID     int64      `db:"survey_answer_id"`
	SurveyID           int64      `db:"survey_id"`
	SurveyTitle        string     `db:"survey_title"`
	AnswerFrequency    string     `db:"answer_frequency"`
	ResponseType       string     `db:"response_type"`
	AnswerDate         *time.Time `db:"answer_date"`
	CreatedDate        *time.Time `db:"created_date"`
	EffectiveDateStart *time.Time `db:"effective_date_start"`
	EffectiveDateEnd   *time.Time `db:"effective_date_end"`
	AreaID             *int64     `db:"area_id"`
	AreaCode           string     `db:"area_code"`
	AreaName           string     `db:"area_name"`
	DistributorID      *int64     `db:"distributor_id"`
	DistributorCode    string     `db:"distributor_code"`
	DistributorName    string     `db:"distributor_name"`
	OutletID           int64      `db:"outlet_id"`
	OutletCode         string     `db:"outlet_code"`
	OutletName         string     `db:"outlet_name"`
	EmpID              int64      `db:"emp_id"`
	EmpCode            string     `db:"emp_code"`
	EmpName            string     `db:"emp_name"`
	SalesmanName       string     `db:"salesman_name"`
	Status             string     `db:"status"`
}

type SurveyReportDetailRow struct {
	SurveyAnswerID     int64      `db:"survey_answer_id"`
	AnswerCustID       string     `db:"answer_cust_id"`
	SurveyID           int64      `db:"survey_id"`
	SurveyTitle        string     `db:"survey_title"`
	AnswerFrequency    string     `db:"answer_frequency"`
	ResponseType       string     `db:"response_type"`
	AnswerDate         *time.Time `db:"answer_date"`
	CreatedDate        *time.Time `db:"created_date"`
	EffectiveDateStart *time.Time `db:"effective_date_start"`
	EffectiveDateEnd   *time.Time `db:"effective_date_end"`
	AreaID             *int64     `db:"area_id"`
	AreaCode           string     `db:"area_code"`
	AreaName           string     `db:"area_name"`
	DistributorID      *int64     `db:"distributor_id"`
	DistributorCode    string     `db:"distributor_code"`
	DistributorName    string     `db:"distributor_name"`
	OutletID           int64      `db:"outlet_id"`
	OutletCode         string     `db:"outlet_code"`
	OutletName         string     `db:"outlet_name"`
	EmpID              int64      `db:"emp_id"`
	EmpCode            string     `db:"emp_code"`
	EmpName            string     `db:"emp_name"`
	SalesmanName       string     `db:"salesman_name"`
	Status             string     `db:"status"`
}

type SurveyReportQuestionRow struct {
	SurveyAnswerDetailID int64   `db:"survey_answer_detail_id"`
	QuestionTemplateID   int64   `db:"question_template_id"`
	SurveyTemplateID     int64   `db:"survey_template_id"`
	Question             string  `db:"question"`
	InputType            string  `db:"input_type"`
	AnswerType           string  `db:"answer_type"`
	Seq                  int     `db:"seq"`
	IsAnswered           bool    `db:"is_answered"`
	FreeTextAnswer       *string `db:"free_text_answer"`
	PhotoPath            *string `db:"photo_path"`
}

type SurveyReportQuestionOptionRow struct {
	QuestionTemplateID int64  `db:"question_template_id"`
	QOptionTemplateID  int64  `db:"q_option_template_id"`
	Option             string `db:"option"`
}

type SurveyReportSelectedOptionRow struct {
	SurveyAnswerDetailID int64  `db:"survey_answer_detail_id"`
	SurveyAnswerOptionID int64  `db:"survey_answer_option_id"`
	QOptionTemplateID    int64  `db:"q_option_template_id"`
	OptionLabel          string `db:"option_label"`
}

type SurveyReportAnswerFileRow struct {
	SurveyAnswerDetailID int64  `db:"survey_answer_detail_id"`
	SurveyAnswerFilesID  int64  `db:"survey_answer_files_id"`
	FileName             string `db:"file_name"`
	FileKey              string `db:"file_key"`
	MediaCategory        string `db:"media_category"`
	FileSize             *int64 `db:"file_size"`
}

type SurveyReportExportRow struct {
	SurveyDate      *time.Time `db:"survey_date"`
	SurveyTitle     string     `db:"survey_title"`
	AreaCode        string     `db:"area_code"`
	AreaName        string     `db:"area_name"`
	DistributorCode string     `db:"distributor_code"`
	DistributorName string     `db:"distributor_name"`
	EmpCode         string     `db:"emp_code"`
	EmpName         string     `db:"emp_name"`
	OutletCode      string     `db:"outlet_code"`
	OutletName      string     `db:"outlet_name"`
	Question        string     `db:"question"`
	Answer          string     `db:"answer"`
	Attachment1     string     `db:"attachment_1"`
	Attachment2     string     `db:"attachment_2"`
	Attachment3     string     `db:"attachment_3"`
}
