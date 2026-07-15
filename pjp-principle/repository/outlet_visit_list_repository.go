package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"time"

	"gorm.io/gorm"
)

type OutletVisitRepo interface {
	Create(ctx context.Context, data model.OutletVisitList)
	Updates(ctx context.Context, Ids []int, data model.OutletVisitList)
	FindByOneColumn(ctx context.Context, selectColumns []string, column string, query any) (data model.OutletVisitList, err error)
	UpsertOutletVisit(ctx context.Context, salesCode, custId string, currentTime int64, date string)
	CreateOutletVisit(ctx context.Context, salesCode, custId, date string)
	UpdateOutletVisitListMultiRowColumnAt(ctx context.Context, salesCode, custId, column string, currentTime int64, date string)
	UpdateOutletVisitListColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime *int64, date string, id int64)
	UpdateOutletVisitListSkipColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime int64, date string, id int64, skipReson string)
	GetSummary(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data model.OutletVisitList, err error)
	GetSummaryStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data response.VisitStatusResponse, err error)
	FindAll(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data []model.OutletVisitList, err error)
	FindById(ctx context.Context, outletVisitId int) (data response.TodoListResponse, err error)
	Delete(ctx context.Context, date string, week int, DestinationCode string, routeCode int) (err error)
	GetVisitsByDateAndSalesman(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitList, error)
	CheckRouteOutletAdditional(ctx context.Context, routeCode *int, DestinationID int) (bool, error)
	MobileCancelAddOutletToRoute(ctx context.Context, data model.OutletVisitList, tx ...*gorm.DB) error
	UpdateByPjpIDAndDate(ctx context.Context, pjpID int, date string, data model.OutletVisitList)
}

type OutletVisitRepoImpl struct {
	Db *gorm.DB
}

func NewOutletVisitRepoImpl(Db *gorm.DB) OutletVisitRepo {
	return &OutletVisitRepoImpl{Db: Db}
}

func (repo *OutletVisitRepoImpl) FindById(ctx context.Context, outletVisitId int) (data response.TodoListResponse, err error) {
	result := repo.Db.WithContext(ctx).
		Table("pjp_principles.outlet_visit_list").
		Select("arrive_at, leave_at", "on_hold", "resume_at", "skip_at").
		Where("id = ?", outletVisitId).
		Scan(&data)

	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) Create(ctx context.Context, data model.OutletVisitList) {
	result := repo.Db.WithContext(ctx).Create(&data)
	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) Updates(ctx context.Context, Ids []int, data model.OutletVisitList) {
	result := repo.Db.WithContext(ctx).Where("id IN ?", Ids).Updates(&data)
	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) FindByOneColumn(ctx context.Context, selectColumns []string, column string, query any) (data model.OutletVisitList, err error) {
	result := repo.Db.WithContext(ctx).Select(selectColumns).Where(column+" = ?", query).Find(&data)

	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) UpsertOutletVisit(ctx context.Context, salesCode, custId string, currentTime int64, date string) {
	// Start the transaction
	tx := repo.Db.WithContext(ctx).Begin()

	// Check if there are any records to insert or update
	var count int64
	checkResult := tx.Raw(`
		SELECT COUNT(*)
		FROM (
			WITH getPjpIds AS (
				SELECT id 
				FROM pjp_principles.permanent_journey_plans 
				WHERE salesman_code = ?
				AND cust_id = ?
			)
			SELECT DISTINCT ON (COALESCE(roa.outlet_id, ro.outlet_id)) 
				rpd.year, rpd.week, rpd.date, rpd.day, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.route_code 
					ELSE ro.route_code 
				END AS route_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_code 
					ELSE ro.outlet_code 
				END AS outlet_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id 
					ELSE ro.outlet_id 
				END AS outlet_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_id 
					ELSE ro.pjp_id 
				END AS pjp_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_code 
					ELSE ro.pjp_code 
				END AS pjp_code
			FROM pjp_principles.route_pop_dailies rpd
			LEFT JOIN pjp_principles.destinations ro 
				ON ro.route_code = rpd.route_code 
				AND ro.pjp_id = rpd.pjp_id
			LEFT JOIN pjp_principles.destinations_additional roa 
				ON roa.route_code = rpd.parent_route 
				AND roa.pjp_id = rpd.pjp_id
			WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
			AND rpd.date = ?
			AND (ro.status = 'Approved' 
				OR roa.status = 'Approved' 
				OR ro.status = 'Approved With Propose' 
				OR roa.status = 'Approved With Propose')
		) subquery
	`, salesCode, custId, date).Count(&count)

	if checkResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(checkResult.Error)
	}

	// If no data is found, rollback the transaction and panic with the error
	if count == 0 {
		tx.Rollback()
		helper.ErrorPanic(fmt.Errorf("no data found on selected date. please approve route first to process"))
	}

	// First, attempt to update existing records
	updateResult := tx.Exec(`
		WITH getPjpIds AS (
			SELECT id 
			FROM pjp_principles.permanent_journey_plans 
			WHERE salesman_code = ?
			AND cust_id = ?
		)
		UPDATE pjp_principles.outlet_visit_list
		SET start = ?
		FROM pjp_principles.route_pop_dailies rpd
		LEFT JOIN pjp_principles.destinations ro 
			ON ro.route_code = rpd.route_code 
			AND ro.pjp_id = rpd.pjp_id
		LEFT JOIN pjp_principles.destinations_additional roa 
			ON roa.route_code = rpd.parent_route 
			AND roa.pjp_id = rpd.pjp_id
		WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
		AND rpd.date = ?
		AND (ro.status = 'Approved' 
			OR roa.status = 'Approved' 
			OR ro.status = 'Approved With Propose' 
			OR roa.status = 'Approved With Propose')
		AND pjp_principles.outlet_visit_list.outlet_id = ro.outlet_id
		AND pjp_principles.outlet_visit_list.pjp_id = ro.pjp_id
		AND pjp_principles.outlet_visit_list.date = rpd.date
		AND pjp_principles.outlet_visit_list.day = rpd.day
	`, salesCode, custId, currentTime, date)

	// If update fails, rollback and panic
	if updateResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(updateResult.Error)
	}

	// Insert the records that didn't exist before (i.e., do not match the conditions of the update)
	insertResult := tx.Exec(`
		WITH getPjpIds AS (
			SELECT id 
			FROM pjp_principles.permanent_journey_plans 
			WHERE salesman_code = ?
			AND cust_id = ?
		)
		INSERT INTO pjp_principles.outlet_visit_list (year, week, date, day, route_code, outlet_code, outlet_id, pjp_id, pjp_code, start)
		SELECT DISTINCT ON (COALESCE(roa.outlet_id, ro.outlet_id)) 
				rpd.year, rpd.week, rpd.date, rpd.day, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.route_code 
					ELSE ro.route_code 
				END AS route_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_code 
					ELSE ro.outlet_code 
				END AS outlet_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id 
					ELSE ro.outlet_id 
				END AS outlet_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_id 
					ELSE ro.pjp_id 
				END AS pjp_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_code 
					ELSE ro.pjp_code 
				END AS pjp_code,
				?
		FROM pjp_principles.route_pop_dailies rpd
		LEFT JOIN pjp_principles.destinations ro 
			ON ro.route_code = rpd.route_code 
			AND ro.pjp_id = rpd.pjp_id
		LEFT JOIN pjp_principles.destinations_additional roa 
			ON roa.route_code = rpd.parent_route 
			AND roa.pjp_id = rpd.pjp_id
		WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
		AND rpd.date = ?
		AND (ro.status = 'Approved' 
			OR roa.status = 'Approved' 
			OR ro.status = 'Approved With Propose' 
			OR roa.status = 'Approved With Propose')
		AND NOT EXISTS (
			SELECT 1
			FROM pjp_principles.outlet_visit_list ovl
			WHERE ovl.outlet_id = ro.outlet_id
			AND ovl.pjp_id = ro.pjp_id
			AND ovl.date = rpd.date
			AND ovl.day = rpd.day
		)
	`, salesCode, custId, currentTime, date)

	// Handle error from the insert
	if insertResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(insertResult.Error)
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		helper.ErrorPanic(err)
	}
}

func (repo *OutletVisitRepoImpl) CreateOutletVisit(ctx context.Context, salesCode, custId, date string) {
	tx := repo.Db.WithContext(ctx).Begin()

	var count int64
	checkResult := tx.Raw(`
		SELECT COUNT(*)
		FROM (
			WITH getPjpIds AS (
				SELECT id 
				FROM pjp_principles.permanent_journey_plans 
				WHERE salesman_code = ?
				AND cust_id = ?
			)
			SELECT DISTINCT ON (COALESCE(roa.outlet_id, ro.outlet_id)) 
				rpd.year, rpd.week, rpd.date, rpd.day, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.route_code 
					ELSE ro.route_code 
				END AS route_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_code 
					ELSE ro.outlet_code 
				END AS outlet_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id 
					ELSE ro.outlet_id 
				END AS outlet_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_id 
					ELSE ro.pjp_id 
				END AS pjp_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_code 
					ELSE ro.pjp_code 
				END AS pjp_code
			FROM pjp_principles.route_pop_dailies rpd
			LEFT JOIN pjp_principles.destinations ro 
				ON ro.route_code = rpd.route_code 
				AND ro.pjp_id = rpd.pjp_id
			LEFT JOIN pjp_principles.destinations_additional roa 
				ON roa.route_code = rpd.parent_route 
				AND roa.pjp_id = rpd.pjp_id
			WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
			AND rpd.date = ?
			AND (ro.status = 'Approved' 
				OR roa.status = 'Approved' 
				OR ro.status = 'Approved With Propose' 
				OR roa.status = 'Approved With Propose')
		) subquery
	`, salesCode, custId, date).Count(&count)

	if checkResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(checkResult.Error)
	}

	if count == 0 {
		tx.Rollback()
		helper.ErrorPanic(fmt.Errorf("no data found on selected date. please approve route first to process"))
	}

	insertResult := tx.Exec(`
	WITH getPjpIds AS (
		SELECT id 
		FROM pjp_principles.permanent_journey_plans 
		WHERE salesman_code = ?
		AND cust_id = ?
	)
	INSERT INTO pjp_principles.outlet_visit_list (year, week, date, day, route_code, outlet_code, outlet_id, pjp_id, pjp_code)
	SELECT DISTINCT ON (COALESCE(roa.outlet_id, ro.outlet_id)) 
			rpd.year, rpd.week, rpd.date, rpd.day, 
			CASE 
				WHEN rpd.parent_route IS NOT NULL THEN roa.route_code 
				ELSE ro.route_code 
			END AS route_code, 
			CASE 
				WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_code 
				ELSE ro.outlet_code 
			END AS outlet_code, 
			CASE 
				WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id 
				ELSE ro.outlet_id 
			END AS outlet_id, 
			CASE 
				WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_id 
				ELSE ro.pjp_id 
			END AS pjp_id, 
			CASE 
				WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_code 
				ELSE ro.pjp_code 
			END AS pjp_code
	FROM pjp_principles.route_pop_dailies rpd
	LEFT JOIN pjp_principles.destinations ro 
		ON ro.route_code = rpd.route_code 
		AND ro.pjp_id = rpd.pjp_id
	LEFT JOIN pjp_principles.destinations_additional roa 
		ON roa.route_code = rpd.parent_route 
		AND roa.pjp_id = rpd.pjp_id
	WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
	AND rpd.date = ?
	AND (ro.status = 'Approved' 
		OR roa.status = 'Approved' 
		OR ro.status = 'Approved With Propose' 
		OR roa.status = 'Approved With Propose')
	AND NOT EXISTS (
		SELECT 1
		FROM pjp_principles.outlet_visit_list ovl
		WHERE ovl.outlet_id = ro.outlet_id
		AND ovl.pjp_id = ro.pjp_id
		AND ovl.date = rpd.date
		AND ovl.day = rpd.day
	)
	`, salesCode, custId, date)

	if insertResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(insertResult.Error)
	}

	if err := tx.Commit().Error; err != nil {
		helper.ErrorPanic(err)
	}
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListMultiRowColumnAt(ctx context.Context, salesCode, custId, column string, currentTime int64, date string) {
	result := repo.Db.WithContext(ctx).Exec(`
		WITH getPjpIds AS (
			SELECT id
			FROM pjp_principles.permanent_journey_plans
			WHERE salesman_code = ? AND cust_id = ?
		)
		UPDATE pjp_principles.outlet_visit_list
		SET `+column+` = ?
		FROM (
			SELECT DISTINCT ON (COALESCE(roa.outlet_id, ro.outlet_id)) 
				rpd.year, rpd.week, rpd.date, rpd.day, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.route_code 
					ELSE ro.route_code 
				END AS route_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_code 
					ELSE ro.outlet_code 
				END AS outlet_code, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id 
					ELSE ro.outlet_id 
				END AS outlet_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_id 
					ELSE ro.pjp_id 
				END AS pjp_id, 
				CASE 
					WHEN rpd.parent_route IS NOT NULL THEN roa.pjp_code 
					ELSE ro.pjp_code 
				END AS pjp_code
			FROM pjp_principles.route_pop_dailies rpd
			LEFT JOIN pjp_principles.destinations ro 
				ON ro.route_code = rpd.route_code 
				AND ro.pjp_id = rpd.pjp_id
			LEFT JOIN pjp_principles.destinations_additional roa 
				ON roa.route_code = rpd.parent_route 
				AND roa.pjp_id = rpd.pjp_id
			WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
			AND rpd.date = ?
			AND (ro.status = 'Approved' 
				OR roa.status = 'Approved' 
				OR ro.status = 'Approved With Propose' 
				OR roa.status = 'Approved With Propose')
		) AS updated_data
		WHERE pjp_principles.outlet_visit_list.pjp_id = updated_data.pjp_id
		AND pjp_principles.outlet_visit_list.year = updated_data.year
		AND pjp_principles.outlet_visit_list.week = updated_data.week
		AND pjp_principles.outlet_visit_list.date = updated_data.date
		AND pjp_principles.outlet_visit_list.day = updated_data.day
		AND pjp_principles.outlet_visit_list.route_code = updated_data.route_code
		AND pjp_principles.outlet_visit_list.outlet_code = updated_data.outlet_code
        AND pjp_principles.outlet_visit_list.outlet_id = updated_data.outlet_id
		AND pjp_principles.outlet_visit_list.pjp_code = updated_data.pjp_code
	`, salesCode, custId, currentTime, date)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime *int64, date string, id int64) {

	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET `+column+` = ?
		WHERE date = ? AND id = ?
	`, currentTime, date, id)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListSkipColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime int64, date string, id int64, skipReson string) {

	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp_principles.outlet_visit_list
		SET `+column+` = ?, skip_reason = ?
		WHERE date = ? AND id = ?
	`, currentTime, skipReson, date, id)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) GetSummary(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data model.OutletVisitList, err error) {
	result := repo.Db.WithContext(ctx).Raw(`
		WITH getPjpIds AS (
			SELECT id
			FROM pjp_principles.permanent_journey_plans
			WHERE salesman_code = ?
		)
		SELECT start, finish
		FROM pjp_principles.outlet_visit_list ovl
		WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
        AND ovl.date = ?
	`, dataFilter.SalesmanCode, dataFilter.Date).Scan(&data)

	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) GetSummaryStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data response.VisitStatusResponse, err error) {
	query := `
		WITH getPjpIds AS (
			SELECT id
			FROM pjp_principles.permanent_journey_plans
			WHERE salesman_code = ?
		)
		SELECT 
			COUNT(CASE 
				WHEN start IS NOT NULL AND arrive_at IS NULL AND leave_at IS NULL AND skip_at IS NULL AND on_hold IS NULL 
				THEN 1 
			END) AS planned,
			COUNT(CASE 
				WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND skip_at IS NULL AND on_hold IS NULL AND leave_at IS NULL
				THEN 1 
			END) AS on_progress,
			COUNT(CASE 
				WHEN skip_at IS NOT NULL AND leave_at IS NULL 
				THEN 1 
			END) AS skipped,
			COUNT(CASE 
				WHEN on_hold IS NOT NULL AND resume_at IS NULL AND skip_at IS NULL AND leave_at IS NULL 
				THEN 1 
			END) AS on_hold,
			COUNT(CASE 
				WHEN leave_at IS NOT NULL 
				THEN 1 
			END) AS finished
		FROM 
			pjp_principles.outlet_visit_list ovl
		WHERE 
			ovl.pjp_id IN (SELECT id FROM getPjpIds)
			AND ovl.date = ?
	`

	result := repo.Db.WithContext(ctx).Raw(query, dataFilter.SalesmanCode, dataFilter.Date).Scan(&data)
	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) FindAll(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data []model.OutletVisitList, err error) {

	result := repo.Db.WithContext(ctx).Raw(`
		WITH getPjpIds AS (
			SELECT id
			FROM pjp_principles.permanent_journey_plans
			WHERE salesman_code = ? AND cust_id = ?
		)
		SELECT 
			ovl.*,
			ro.outlet_name,
			ro.outlet_address,
			CASE 
				WHEN ovl.start IS NOT NULL AND ovl.arrive_at IS NULL AND ovl.leave_at IS NULL AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL THEN 'planned'
				WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL AND ovl.leave_at IS NULL THEN 'on_progress'
				WHEN ovl.skip_at IS NOT NULL AND ovl.leave_at IS NULL THEN 'skipped'
				WHEN ovl.on_hold IS NOT NULL AND resume_at IS NULL AND ovl.skip_at IS NULL AND ovl.leave_at IS NULL THEN 'on_hold'
				WHEN ovl.leave_at IS NOT NULL THEN 'finished'
				ELSE 'planned' 
			END AS status
		FROM pjp_principles.outlet_visit_list ovl
		LEFT JOIN (
			SELECT DISTINCT ON (outlet_id) outlet_id, outlet_name, outlet_address
			FROM pjp_principles.destinations
			ORDER BY outlet_id
		) ro ON ovl.outlet_id = ro.outlet_id
		WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
		AND ovl.date = ?
	`, dataFilter.SalesmanCode, dataFilter.CustID, dataFilter.Date).Scan(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) Delete(ctx context.Context, date string, week int, DestinationCode string, routeCode int) error {
	var outletVisitList model.OutletVisitList

	// First, try to find the record
	findResult := repo.Db.WithContext(ctx).
		Where("date = ? AND week = ? AND route_code = ? AND outlet_code = ?", date, week, routeCode, DestinationCode).
		First(&outletVisitList)

	log.Printf("Data: %+v", findResult)

	// If no records are found, return nil or an error if needed
	if findResult.Error != nil {
		if findResult.Error == gorm.ErrRecordNotFound {
			// Return nil if you don't want to throw an error when no record is found
			return nil
		}
		return findResult.Error
	}

	// Proceed to delete the record if it exists
	deleteResult := repo.Db.WithContext(ctx).
		Where("date = ? AND week = ? AND route_code = ? AND outlet_code = ?", date, week, routeCode, DestinationCode).
		Delete(&outletVisitList)

	// Check for any errors during deletion
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil
}

func (repo *OutletVisitRepoImpl) GetVisitsByDateAndSalesman(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitList, error) {
	var visits []model.OutletVisitList

	query := repo.Db.Table("pjp_principles.outlet_visit_list").
		Select("id, year, week, date, day, route_code, outlet_id, outlet_code, pjp_id, pjp_code, start, finish, skip_at, leave_at, arrive_at, on_hold, resume_at, skip_reason")

	subquery := repo.Db.Table("pjp_principles.permanent_journey_plans").
		Select("id").
		Where("salesman_id = ?", dataFilter.SalesmanId)
	query = query.Where("pjp_principles.outlet_visit_list.pjp_id IN (?)", subquery)

	if dataFilter.Date != "" {
		query = query.Where("pjp_principles.outlet_visit_list.date = ?", dataFilter.Date)
	}
	if dataFilter.Year != "" {
		query = query.Where("EXTRACT(YEAR FROM pjp_principles.outlet_visit_list.date) = ?", dataFilter.Year)
	}
	if dataFilter.Month != "" {
		query = query.Where("EXTRACT(MONTH FROM pjp_principles.outlet_visit_list.date) = ?", dataFilter.Month)
	}

	err := query.Find(&visits).Error
	return visits, err
}

func (repo *OutletVisitRepoImpl) CheckRouteOutletAdditional(ctx context.Context, routeCode *int, DestinationID int) (bool, error) {
	if routeCode == nil {
		return false, nil
	}

	var count int64
	err := repo.Db.Table("pjp_principles.destinations_additional").
		Where("route_code = ? AND outlet_id = ?", *routeCode, DestinationID).
		Count(&count).Error

	return count > 0, err
}

func (repo *OutletVisitRepoImpl) MobileCancelAddOutletToRoute(ctx context.Context, data model.OutletVisitList, tx ...*gorm.DB) error {
	var db *gorm.DB
	if len(tx) > 0 && tx[0] != nil { // Jika tx diberikan dan tidak nil, gunakan tx
		db = tx[0]
	} else { // Jika tidak, gunakan koneksi default repo.Db
		db = repo.Db
	}
	today := time.Now().Format("2006-01-02")

	results := db.WithContext(ctx).
		Model(&data).
		Where("route_code = ? AND outlet_code = ? AND pjp_id = ? AND DATE(date) = ? AND is_planned = ?",
			data.RouteCode, data.DestinationCode, data.PjpID, today, false).
		Delete(&data)

	if results.Error != nil {
		return results.Error
	}
	if results.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (repo *OutletVisitRepoImpl) UpdateByPjpIDAndDate(ctx context.Context, pjpID int, date string, data model.OutletVisitList) {
	result := repo.Db.WithContext(ctx).
		Where("pjp_id = ? AND date = ?", pjpID, date).
		Updates(&data)

	helper.ErrorPanic(result.Error)
}
