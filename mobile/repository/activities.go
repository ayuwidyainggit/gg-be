package repository

import (
	"context"
	"gorm.io/gorm"
	"log"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/times"
	"time"
)

type (
	RepositoryActivitiesImpl struct {
		*gorm.DB
	}
)

type ActivitiesRepository interface {
	FindSummaryActivity(params entity.SummaryDailyRequest) (model.ActivitiesSummary, error)
}

func NewActivitiesRepository(db *gorm.DB) *RepositoryActivitiesImpl {
	return &RepositoryActivitiesImpl{db}
}

func (repo *RepositoryActivitiesImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

// FindSummaryActivity retrieves daily summary for a salesman.
// Enhanced: All queries now filter by CURRENT_DATE for accurate daily reporting.
// - Plan: Count of outlets to visit today (not yet arrived/skipped/left)
// - Visit: Count of outlets visited today (arrived but not skipped/left)
// - Effective Call: Count of orders placed today at visited outlets
// - StartTime: Clock in time (type=1) from mobile.attendances for current date
// - EndTime: Clock out time (type=2) from mobile.attendances for current date
func (repository *RepositoryActivitiesImpl) FindSummaryActivity(params entity.SummaryDailyRequest) (model.ActivitiesSummary, error) {
	activitiesSummary := model.ActivitiesSummary{}
	var distributorID *string
	var isPrincipal bool
	var query string

	err := repository.WithContext(context.Background()).Table("smc.m_customer").
		Select("distributor_id").
		Where("cust_id = ?", params.CustId).
		Scan(&distributorID).Error

	if err != nil {
		return activitiesSummary, err
	}
	timeNow, err := times.GetCurrentTime()
	if err != nil {
		return activitiesSummary, err
	}

	isPrincipal = len(params.CustId) == 6 && distributorID == nil

	if isPrincipal {
		query = `
        SELECT
            (
                    SELECT COUNT(ovl.id) AS plan
                    FROM pjp_principles.outlet_visit_list AS ovl
                    JOIN pjp_principles.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND is_planned = true
                    AND ovl.date = $2
            ) AS plan,
            (
                    SELECT COUNT(ovl.id) AS visit
                    FROM pjp_principles.outlet_visit_list AS ovl
                    JOIN pjp_principles.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.arrive_at IS NOT NULL
                    AND ovl.date = $2
            ) AS visit,
            (
                    SELECT COUNT(ovl.id) AS extra_call
                    FROM pjp_principles.outlet_visit_list AS ovl
                    JOIN pjp_principles.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.is_extra_call = true
                    AND ovl.date = $2
            ) AS extra_call,
            (
                    SELECT COUNT(ovl.id) AS skip
                    FROM pjp_principles.outlet_visit_list AS ovl
                    JOIN pjp_principles.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.skip_at IS NOT NULL
                    AND ovl.date = $2
            ) AS skip,
            (
                    SELECT COUNT(DISTINCT orders.outlet_id) AS effective_call
                    FROM sls.order AS orders
                    JOIN mst.m_outlet AS ovl ON ovl.outlet_id = orders.outlet_id
                    WHERE orders.salesman_id = $1
                    AND DATE(orders.ro_date) = $2
            ) AS effective_call,
            (
                    SELECT TO_CHAR(created_at, 'HH24:MI:SS')
                    FROM mobile.attendances
                    WHERE emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = $1 LIMIT 1)
                    AND type = 1
                    AND DATE(created_at) = $2
                    ORDER BY created_at ASC
                    LIMIT 1
            ) AS start_time,
            (
                    SELECT TO_CHAR(created_at, 'HH24:MI:SS')
                    FROM mobile.attendances
                    WHERE emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = $1 LIMIT 1)
                    AND type = 2
                    AND DATE(created_at) = $2
                    ORDER BY created_at DESC
                    LIMIT 1
            ) AS end_time`
	} else {
		query = `
        SELECT
            (
                    SELECT COUNT(ovl.id) AS plan
                    FROM pjp.outlet_visit_list AS ovl
                    JOIN pjp.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND is_planned = true
                    AND ovl.date = $2
            ) AS plan,
            (
                    SELECT COUNT(ovl.id) AS visit
                    FROM pjp.outlet_visit_list AS ovl
                    JOIN pjp.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.arrive_at IS NOT NULL
                    AND ovl.date = $2
            ) AS visit,
            (
                    SELECT COUNT(ovl.id) AS extra_call
                    FROM pjp.outlet_visit_list AS ovl
                    JOIN pjp.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.is_extra_call = true
                    AND ovl.date = $2
            ) AS extra_call,
            (
                    SELECT COUNT(ovl.id) AS skip
                    FROM pjp.outlet_visit_list AS ovl
                    JOIN pjp.permanent_journey_plans AS perpjp ON perpjp.pjp_code = ovl.pjp_code
                    WHERE perpjp.salesman_id = $1
                    AND ovl.skip_at IS NOT NULL
                    AND ovl.date = $2
            ) AS skip,
            (
                    SELECT COUNT(DISTINCT orders.outlet_id) AS effective_call
                    FROM sls.order AS orders
                    JOIN mst.m_outlet AS ovl ON ovl.outlet_id = orders.outlet_id
                    WHERE orders.salesman_id = $1
                    AND DATE(orders.ro_date) = $2
            ) AS effective_call,
            (
                    SELECT TO_CHAR(created_at, 'HH24:MI:SS')
                    FROM mobile.attendances
                    WHERE emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = $1 LIMIT 1)
                    AND type = 1
                    AND DATE(created_at) = $2
                    ORDER BY created_at ASC
                    LIMIT 1
            ) AS start_time,
            (
                    SELECT TO_CHAR(created_at, 'HH24:MI:SS')
                    FROM mobile.attendances
                    WHERE emp_code = (SELECT emp_code FROM mst.m_employee WHERE emp_id = $1 LIMIT 1)
                    AND type = 2
                    AND DATE(created_at) = $2
                    ORDER BY created_at DESC
                    LIMIT 1
            ) AS end_time`
	}

	err = repository.model(context.Background()).Raw(query, params.EmployeeId, timeNow.Format(time.DateOnly)).Scan(&activitiesSummary).Error
	if err != nil {
		log.Println("EmployeeRepository, FindOneByEmployeeIdAndCustId, err:", err.Error())
		return activitiesSummary, err
	}

	return activitiesSummary, nil
}
