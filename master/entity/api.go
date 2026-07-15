package entity

type GeneralQueryFilter struct {
	CustId         string
	ParentCustId   string
	Page           int    `query:"page"`
	Limit          int    `query:"limit" validate:"required"`
	Query          string `query:"q"`
	Mode           string `query:"mode"`
	Sort           string `query:"sort"`
	IsActive       *int   `query:"is_active"`
	DistributorID  int64  `query:"-"`
	DistributorIDs []int  `query:"-"`
}

type Pagination struct {
	TotalRecord int    `json:"total_record"`
	PageCurrent int    `json:"page_current"`
	PageLimit   int    `json:"page_limit"`
	PageTotal   int    `json:"page_total"`
	RequestID   string `json:"request_id,omitempty"`
}

type ApiResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Paging    interface{} `json:"paging,omitempty"`
	RequestId string      `json:"request_id"`
}

type PublishJob struct {
	JobName          string `json:"job_name" validate:"required,max=100"`
	JobDesc          string `json:"job_desc" validate:"required,max=255"`
	JobType          string `json:"job_type" validate:"required,oneof=duration random_duration cron daily weekly monthly one_time"`
	CronExpression   string `json:"cron_expression,omitempty"`
	DayOfWeekOrMonth int    `json:"day_of_week_or_month,omitempty"`
	TimeOfDay        string `json:"time_of_day,omitempty"`
	RunAt            string `json:"run_at,omitempty"`
	Task             string `json:"task" validate:"required"`
	Url              string `json:"url" validate:"omitempty"`
	Payload          string `json:"payload" validate:"omitempty"`
	CreatedBy        string `json:"created_by" validate:"required,max=100"`
}
