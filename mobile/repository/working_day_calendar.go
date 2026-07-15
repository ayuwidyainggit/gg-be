package repository

import (
	"mobile/entity"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryWorkingDayCalendarImpl struct {
		*gorm.DB
	}
)

type WorkingDayCalendarRepository interface {
	FindAll(dataFilter entity.WorkingDayCalendarQueryFilter) (data []model.WorkingDayCalendarDetail, err error)
	FindMonthsByWDCID(wdcID int) (data []model.WorkingDayCalendarMonthDetail, err error)
}

func NewWorkingDayCalendarRepository(db *gorm.DB) *RepositoryWorkingDayCalendarImpl {
	return &RepositoryWorkingDayCalendarImpl{db}
}

func (repo *RepositoryWorkingDayCalendarImpl) FindAll(dataFilter entity.WorkingDayCalendarQueryFilter) (data []model.WorkingDayCalendarDetail, err error) {
	err = repo.DB.
		Select(`
			wdc.working_day_calendar_id,
			wdc.cust_id,
			wdc.title,
			wdc.start_date,
			wdc.number_of_weeks,
			wdc.end_date,
			wdc.default_holidays,
			wdc.is_closed,
			CASE
				WHEN CURRENT_DATE BETWEEN wdc.start_date AND wdc.end_date THEN true
				ELSE false
			END AS is_active
		`).
		Table("mst.working_day_calendar AS wdc").
		Where("wdc.cust_id = ?", dataFilter.ParentCustID).
		Where("wdc.is_closed = false").
		Order("wdc.start_date ASC").
		Find(&data).Error

	return data, err
}

func (repo *RepositoryWorkingDayCalendarImpl) FindMonthsByWDCID(wdcID int) (data []model.WorkingDayCalendarMonthDetail, err error) {
	err = repo.DB.
		Select(`
			(b.first_date = DATE_TRUNC('month', CURRENT_DATE)::date) AS is_active,
			EXTRACT(MONTH FROM b.first_date)::int AS month,
			EXTRACT(YEAR FROM b.first_date)::int AS year,
			TRIM(TO_CHAR(b.first_date, 'Month')) || ' (' ||
				TO_CHAR(b.first_date, 'DD Mon YYYY') || ' – ' ||
				TO_CHAR((b.first_date + INTERVAL '1 month' - INTERVAL '1 day')::date, 'DD Mon YYYY') || ')'
				AS text_month
		`).
		Table("mst.working_day_calendar wdc").
		Joins(`CROSS JOIN LATERAL (
			SELECT generate_series(
				DATE_TRUNC('month', wdc.start_date),
				DATE_TRUNC('month', wdc.end_date),
				INTERVAL '1 month'
			)::date AS first_date
		) b`).
		Where("wdc.working_day_calendar_id = ?", wdcID).
		Order("b.first_date").
		Find(&data).Error

	return data, err
}
