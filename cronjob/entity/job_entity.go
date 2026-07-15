package entity

import "github.com/google/uuid"

type CreateJobBody struct {
	JobName          string `json:"job_name" validate:"required,max=100"`
	JobDesc          string `json:"job_desc" validate:"required,max=255"`
	JobType          string `json:"job_type" validate:"required,oneof=duration random_duration cron daily weekly monthly one_time"`
	CronExpression   string `json:"cron_expression,omitempty"`
	DayOfWeekOrMonth int    `json:"day_of_week,omitempty"`
	TimeOfDay        string `json:"time_of_day,omitempty"`
	RunAt            string `json:"run_at,omitempty"`
	Task             string `json:"task" validate:"required"`
	Url              string `json:"url" validate:"omitempty"`
	Payload          string `json:"payload" validate:"omitempty"`
	CreatedBy        string `json:"created_by" validate:"required,max=100"`
}

type UpdateJobBody struct {
	JobName          string `json:"job_name" validate:"required,max=100"`
	JobDesc          string `json:"job_desc" validate:"required,max=255"`
	JobType          string `json:"job_type" validate:"required,oneof=duration random_duration cron daily weekly monthly one_time"`
	CronExpression   string `json:"cron_expression,omitempty"`
	DayOfWeekOrMonth int    `json:"day_of_week,omitempty"`
	TimeOfDay        string `json:"time_of_day,omitempty"`
	RunAt            string `json:"run_at,omitempty"`
	Task             string `json:"task" validate:"required"`
	Url              string `json:"url" validate:"omitempty"`
	Payload          string `json:"payload" validate:"omitempty"`
	UpdatedBy        string `json:"updated_by" validate:"required,max=100"`
	Active           bool   `json:"active"`
}

type UpdateJobBodyParam struct {
	JobID uuid.UUID `params:"job_id" validate:"required"`
}

type DetailJobBodyParam struct {
	JobID uuid.UUID `params:"job_id" validate:"required"`
}

type DeleteJobBodyParams struct {
	JobID uuid.UUID `params:"job_id" validate:"required"`
}

type JobList struct {
	ID               uuid.UUID `json:"job_id"`
	JobName          string    `json:"job_name"`
	JobDesc          string    `json:"job_desc"`
	JobType          string    `json:"job_type"`
	CronExpression   string    `json:"cron_expression"`
	DayOfWeekOrMonth int       `json:"day_of_week_or_month"`
	TimeOfDay        string    `json:"time_of_day,omitempty"`
	RunAt            string    `json:"run_at"`
	Task             string    `json:"task"`
	Url              string    `json:"url"`
	Payload          string    `json:"payload"`
	Active           bool      `json:"active"`
	CreatedBy        string    `json:"created_by"`
}
