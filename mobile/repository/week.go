package repository

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryWeekImpl struct {
		*gorm.DB
	}
)

type WeekRepository interface {
	FindAll(dataFilter entity.WeekListQueryFilter) (weeks []model.WeekListDetail, total int64, lastPage int, err error)
	FindWorkDaysByWeekId(weekId int, weekStart, weekEnd, custId string, isDistributor bool) (workDays []model.WorkDayDetail, err error)
}

func NewWeekRepository(db *gorm.DB) *RepositoryWeekImpl {
	return &RepositoryWeekImpl{db}
}

func (repo *RepositoryWeekImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryWeekImpl) FindAll(dataFilter entity.WeekListQueryFilter) (weeks []model.WeekListDetail, total int64, lastPage int, err error) {
	// Build query for data retrieval
	query := repository.DB.
		Select(`
			mw.per_year,
			mw.per_id,
			mw.week_id,
			mw.week_start,
			mw.week_end,
		  case
				when CURRENT_DATE BETWEEN mw.week_start and mw.week_end then true
   	  else false
      end as is_active,
			mw.is_closed,
			mw.closed_at,
			COALESCE(mw.closed_by, 0) as closed_by,
			COALESCE(u.user_name, '') as closed_by_name
		`).
		Table("mst.m_week AS mw").
		Joins("LEFT JOIN sys.m_user AS u ON u.user_id = mw.closed_by").
		Where("mw.working_day_calendar_id = ?", dataFilter.WDCID)

	// Apply filters
	// if dataFilter.CustID != "" {
	// 	query = query.Where("mw.cust_id = ?", dataFilter.CustID)
	// }

	if dataFilter.Month != nil && dataFilter.Year != nil {
		query = query.Where(`
						(EXTRACT(MONTH FROM mw.week_start) = ? and EXTRACT(YEAR FROM mw.week_start) = ?)
						or
						(EXTRACT(MONTH FROM mw.week_end) = ? and EXTRACT(YEAR FROM mw.week_end) = ?)
						`, *dataFilter.Month, *dataFilter.Year, *dataFilter.Month, *dataFilter.Year)
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			query = query.Where("mw.is_active = true")
		} else if *dataFilter.IsActive == 0 {
			query = query.Where("mw.is_active = false")
		}
	}

	// Count total using the same query
	err = query.Count(&total).Error
	if err != nil {
		return weeks, 0, 0, err
	}

	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			sortField := sortParts[0]
			sortOrder := strings.ToUpper(sortParts[1])
			if sortOrder == "ASC" || sortOrder == "DESC" {
				// Map field names to table columns
				if sortField == "created_date" {
					// Use week_start as default sort if created_date not available
					sortField = "mw.week_start"
				} else {
					// Add table prefix for other fields
					sortField = fmt.Sprintf("mw.%s", sortField)
				}
				query = query.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
			}
		}
	} else {
		// Default sort by week_start descending
		query = query.Order("mw.week_start DESC")
	}

	offset := (dataFilter.Page - 1) * dataFilter.Limit
	err = query.Offset(offset).Limit(dataFilter.Limit).Find(&weeks).Error
	if err != nil {
		return weeks, 0, 0, err
	}

	lastPage = int((total + int64(dataFilter.Limit) - 1) / int64(dataFilter.Limit))
	if lastPage <= 0 {
		lastPage = 1
	}

	return weeks, total, lastPage, nil
}

func (repository *RepositoryWeekImpl) FindWorkDaysByWeekId(weekId int, weekStart, weekEnd, parentCustID string, isDistributor bool) (workDays []model.WorkDayDetail, err error) {
	query := repository.DB.
		Select(`
			mwd.per_year,
			mwd.per_id,
			mwd.week_id,
			mwd.work_date,
			mwd.is_work,
			case
			 when CURRENT_DATE = mwd.work_date then true
			 else false
			end as is_active,
			mwd.is_closed,
			mwd.closed_at,
			COALESCE(mwd.closed_by, 0) as closed_by,
			COALESCE(u.user_name, '') as closed_by_name,
			COALESCE(COUNT(DISTINCT ovl.outlet_id), 0) as number_of_outlet
		`).
		Table("mst.m_work_day AS mwd").
		Joins("LEFT JOIN sys.m_user AS u ON u.user_id = mwd.closed_by").
		Joins("LEFT JOIN pjp.outlet_visit_list AS ovl ON ovl.date = mwd.work_date").
		Where("mwd.week_id = ?", weekId)

	if parentCustID != "" {
		query = query.Where("mwd.cust_id = ?", parentCustID)
	}

	query.Where("mwd.work_date >= ?", weekStart)
	query.Where("mwd.work_date <= ?", weekEnd)

	err = query.
		Group("mwd.per_year, mwd.per_id, mwd.week_id, mwd.work_date, mwd.is_work, mwd.is_active, mwd.is_closed, mwd.closed_at, mwd.closed_by, u.user_name").
		Order("mwd.work_date ASC").
		Find(&workDays).Error

	return workDays, err
}
