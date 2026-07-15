package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	JobID            uuid.UUID `gorm:"column:job_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"job_id"`
	JobName          string    `gorm:"column:job_name" json:"job_name"`
	JobDesc          string    `gorm:"column:job_desc" json:"job_desc"`
	JobType          string    `gorm:"column:job_type;not null" json:"job_type"`
	CronExpression   string    `gorm:"column:cron_expression" json:"cron_expression"`
	DayOfWeekOrMonth int       `gorm:"column:day_of_week_or_month;type:integer" json:"day_of_week_or_month"`
	TimeOfDay        *string   `gorm:"column:time_of_day;type:text" json:"time_of_day"`
	RunAt            *string   `gorm:"column:run_at;type:text" json:"run_at"` // or pq.StringArray, or pq.GenericArray, depending on needs.
	Task             string    `gorm:"column:task;not null" json:"task"`
	Url              string    `gorm:"column:url;" json:"url"`
	Payload          string    `gorm:"column:payload;" json:"payload"`
	Active           bool      `gorm:"column:active;default:true" json:"active"`
	CreatedBy        string    `gorm:"column:created_by" json:"created_by"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (m *Job) BeforeCreate(trx *gorm.DB) (err error) {

	// m.Jobpass = &tempHashPass
	// m.CreatedAt = time.Now()

	return nil
}

func (Job) TableName() string {
	return "public.jobs"
}
