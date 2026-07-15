package entity

import (
	"encoding/base64"
	"time"
)

const (
	FILE_STATUS_READY      = 1
	FILE_STATUS_PROCESSING = 2
	FILE_STATUS_FAILED     = 3
	FILE_STATUS_EXPIRED    = 4
)

var SurveyReportFileStatusName = map[int]string{
	FILE_STATUS_READY:      "Ready",
	FILE_STATUS_PROCESSING: "Processing",
	FILE_STATUS_FAILED:     "Failed",
	FILE_STATUS_EXPIRED:    "Expired",
}

type SurveyReportQueryFilter struct {
	Page        int        `query:"page"`
	Limit       int        `query:"limit"`
	Query       string     `query:"q"`
	Sort        string     `query:"sort"`
	StartDate   *time.Time `query:"-"`
	EndDate     *time.Time `query:"-"`
	SurveyID    []int64    `query:"survey_id"`
	SurveyTitle []string   `query:"survey_title"`
	AreaID      []int64    `query:"area_id"`

	CustID string `query:"-"`
}

type SurveyReportParams struct {
	SurveyAnswerID int64 `params:"survey_answer_id" validate:"required,min=1"`
}

type SurveyReportListResponse struct {
	SurveyAnswerID     int64      `json:"survey_answer_id"`
	SurveyID           int64      `json:"survey_id"`
	SurveyTitle        string     `json:"survey_title"`
	AnswerFrequency    string     `json:"answer_frequency"`
	ResponseType       string     `json:"response_type"`
	AnswerDate         *time.Time `json:"answer_date"`
	CreatedDate        *time.Time `json:"created_date"`
	EffectiveDateStart *time.Time `json:"effective_date_start"`
	EffectiveDateEnd   *time.Time `json:"effective_date_end"`
	AreaID             *int64     `json:"area_id"`
	AreaCode           string     `json:"area_code"`
	AreaName           string     `json:"area_name"`
	DistributorID      *int64     `json:"distributor_id"`
	DistributorCode    string     `json:"distributor_code"`
	DistributorName    string     `json:"distributor_name"`
	OutletID           int64      `json:"outlet_id"`
	OutletCode         string     `json:"outlet_code"`
	OutletName         string     `json:"outlet_name"`
	EmpID              int64      `json:"emp_id"`
	EmpCode            string     `json:"emp_code"`
	EmpName            string     `json:"emp_name"`
	SalesmanName       string     `json:"salesman_name"`
	Status             string     `json:"status"`
}

type SurveyReportDetailResponse struct {
	SurveyAnswerID     int64                           `json:"survey_answer_id"`
	SurveyID           int64                           `json:"survey_id"`
	SurveyTitle        string                          `json:"survey_title"`
	AnswerFrequency    string                          `json:"answer_frequency"`
	ResponseType       string                          `json:"response_type"`
	AnswerDate         *time.Time                      `json:"answer_date"`
	CreatedDate        *time.Time                      `json:"created_date"`
	EffectiveDateStart *time.Time                      `json:"effective_date_start"`
	EffectiveDateEnd   *time.Time                      `json:"effective_date_end"`
	Area               SurveyReportAreaResponse        `json:"area"`
	Distributor        SurveyReportDistributorResponse `json:"distributor"`
	Outlet             SurveyReportOutletResponse      `json:"outlet"`
	Salesman           SurveyReportSalesmanResponse    `json:"salesman"`
	Status             string                          `json:"status"`
	Details            []SurveyReportQuestion          `json:"details"`
}

type SurveyReportAreaResponse struct {
	AreaID   *int64 `json:"area_id"`
	AreaCode string `json:"area_code"`
	AreaName string `json:"area_name"`
}

type SurveyReportDistributorResponse struct {
	DistributorID   *int64 `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
}

type SurveyReportOutletResponse struct {
	OutletID   int64  `json:"outlet_id"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
}

type SurveyReportSalesmanResponse struct {
	EmpID        int64  `json:"emp_id"`
	EmpCode      string `json:"emp_code"`
	EmpName      string `json:"emp_name"`
	SalesmanName string `json:"salesman_name"`
}

type SurveyReportQuestion struct {
	SurveyAnswerDetailID int64                        `json:"survey_answer_detail_id"`
	QuestionTemplateID   int64                        `json:"question_template_id"`
	SurveyTemplateID     int64                        `json:"survey_template_id"`
	Question             string                       `json:"question"`
	InputType            string                       `json:"input_type"`
	AnswerType           string                       `json:"answer_type"`
	Answer               string                       `json:"answer"`
	Seq                  int                          `json:"seq"`
	IsAnswered           bool                         `json:"is_answered"`
	FreeTextAnswer       *string                      `json:"free_text_answer"`
	PhotoPath            *string                      `json:"photo_path"`
	Options              []SurveyReportQuestionOption `json:"options"`
	SelectedOptions      []SurveyReportSelectedOption `json:"selected_options"`
	Files                []SurveyReportAnswerFile     `json:"files"`
}

type SurveyReportQuestionOption struct {
	QOptionTemplateID int64  `json:"q_option_template_id"`
	Option            string `json:"option"`
}

type SurveyReportSelectedOption struct {
	SurveyAnswerOptionID int64  `json:"survey_answer_option_id"`
	QOptionTemplateID    int64  `json:"q_option_template_id"`
	OptionLabel          string `json:"option_label"`
}

type SurveyReportAnswerFile struct {
	SurveyAnswerFilesID int64  `json:"survey_answer_files_id"`
	FileName            string `json:"file_name"`
	FileKey             string `json:"file_key"`
	MediaCategory       string `json:"media_category"`
	FileSize            *int64 `json:"file_size"`
}

type SurveyReportExportResponse struct {
	ReportID       string `json:"report_id"`
	ReportName     string `json:"report_name"`
	FileStatus     int    `json:"file_status"`
	FileStatusName string `json:"file_status_name"`
	FileBase64     string `json:"file_base64,omitempty"`
	CreatedBy      string `json:"created_by"`
}

func (report SurveyReportExportResponse) GetFileStatusName() string {
	return SurveyReportFileStatusName[report.FileStatus]
}

func (report SurveyReportExportResponse) DecodeFileBase64() ([]byte, error) {
	return base64.StdEncoding.DecodeString(report.FileBase64)
}
