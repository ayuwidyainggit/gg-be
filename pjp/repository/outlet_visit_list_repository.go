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
	CreateOutletVisit(ctx context.Context, salesmanID int64, custId, date string)
	UpdateOutletVisitListMultiRowColumnAt(ctx context.Context, salesCode, custId, column string, currentTime int64, date string)
	UpdateOutletVisitListColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime *int64, date string, id int64)
	UpdateOutletVisitListSkipColumnAt(ctx context.Context, column string, currentTime int64, date string, id int64, skipReson string, InOutlet bool, fileInfo model.OutletVisitList)
	GetSummary(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data model.OutletVisitList, err error)
	GetSummaryStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data response.VisitStatusResponse, err error)
	FindAll(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data []model.OutletVisitList, err error)
	FindById(ctx context.Context, outletVisitId int) (data response.TodoListResponse, err error)
	Delete(ctx context.Context, date string, week int, outletCode string, routeCode int) (err error)
	GetVisitsByDateAndSalesman(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitList, error)
	CheckRouteOutletAdditional(ctx context.Context, routeCode *int, outletID int) (bool, error)
	MobileCancelAddOutletToRoute(ctx context.Context, data model.OutletVisitList, tx ...*gorm.DB) error
	UpdateByPjpIDAndDate(ctx context.Context, pjpID int, date string, data model.OutletVisitList)
	UpdateOutletVisitListWithFile(ctx context.Context, id int64, date string, fileInfo model.OutletVisitList) error
	InsertOrUpdateMobileVisit(ctx context.Context, custID, empCode string, outletID int, latitude, longitude float64, fileURL string, arriveAt int64, reason string)
	UpdateOutletVisitListByID(ctx context.Context, tx *gorm.DB, id int64, data model.OutletVisitList) error
}

type OutletVisitRepoImpl struct {
	Db *gorm.DB
}

func NewOutletVisitRepoImpl(Db *gorm.DB) OutletVisitRepo {
	return &OutletVisitRepoImpl{Db: Db}
}

func (repo *OutletVisitRepoImpl) FindById(ctx context.Context, outletVisitId int) (data response.TodoListResponse, err error) {
	result := repo.Db.WithContext(ctx).
		Table("pjp.outlet_visit_list").
		Select("arrive_at, leave_at", "on_hold", "resume_at", "skip_at", "skip_reason", "skip_in_outlet as in_outlet").
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
				FROM pjp.permanent_journey_plans 
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
			FROM pjp.route_pop_daily rpd
			LEFT JOIN pjp.route_outlet ro 
				ON ro.route_code = rpd.route_code 
				AND ro.pjp_id = rpd.pjp_id
			LEFT JOIN pjp.route_outlet_additional roa 
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
			FROM pjp.permanent_journey_plans 
			WHERE salesman_code = ?
			AND cust_id = ?
		)
		UPDATE pjp.outlet_visit_list
		SET start = ?
		FROM pjp.route_pop_daily rpd
		LEFT JOIN pjp.route_outlet ro 
			ON ro.route_code = rpd.route_code 
			AND ro.pjp_id = rpd.pjp_id
		LEFT JOIN pjp.route_outlet_additional roa 
			ON roa.route_code = rpd.parent_route 
			AND roa.pjp_id = rpd.pjp_id
		WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
		AND rpd.date = ?
		AND (ro.status = 'Approved' 
			OR roa.status = 'Approved' 
			OR ro.status = 'Approved With Propose' 
			OR roa.status = 'Approved With Propose')
		AND pjp.outlet_visit_list.outlet_id = ro.outlet_id
		AND pjp.outlet_visit_list.pjp_id = ro.pjp_id
		AND pjp.outlet_visit_list.date = rpd.date
		AND pjp.outlet_visit_list.day = rpd.day
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
			FROM pjp.permanent_journey_plans 
			WHERE salesman_code = ?
			AND cust_id = ?
		)
		INSERT INTO pjp.outlet_visit_list (year, week, date, day, route_code, outlet_code, outlet_id, pjp_id, pjp_code, start)
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
		FROM pjp.route_pop_daily rpd
		LEFT JOIN pjp.route_outlet ro 
			ON ro.route_code = rpd.route_code 
			AND ro.pjp_id = rpd.pjp_id
		LEFT JOIN pjp.route_outlet_additional roa 
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
			FROM pjp.outlet_visit_list ovl
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

func (repo *OutletVisitRepoImpl) CreateOutletVisit(ctx context.Context, salesmanID int64, custId, date string) {
	tx := repo.Db.WithContext(ctx).Begin()

	var distributorID *string

	err := tx.Table("smc.m_customer").
		Select("distributor_id").
		Where("cust_id = ?", custId).
		Scan(&distributorID).Error
	if err != nil {
		tx.Rollback()
		helper.ErrorPanic(fmt.Errorf("failed to check customer status: %v", err))
	}

	isPrincipal := len(custId) == 6 && distributorID == nil
	var insertQuery string
	var checkResultQuery string

	if isPrincipal {
		checkResultQuery = `
            SELECT COUNT(*)
                FROM (
                    WITH getPjpIds AS (
                        SELECT id
                        FROM pjp_principles.permanent_journey_plans
                        WHERE salesman_id = ?
                        AND cust_id = ?
						AND approval_status = 'Approved'
                    )
                    SELECT
                        roh.year,
						roh.week,
						roh.date,
						to_char(roh.date, 'Dy') as day,
						roh.route_code,
						roh.destination_code,
						roh.destination_id,
						roh.destination_type,
						roh.pjp_id,
						roh.pjp_code
                    FROM pjp_principles.destinations_history roh
                    WHERE roh.pjp_id IN (SELECT id FROM getPjpIds)
					AND roh.date = ?
                    AND NOT EXISTS (
						SELECT 1 
						FROM pjp_principles.outlet_visit_list ovl
						WHERE ovl.outlet_id = roh.destination_id
						AND ovl.pjp_id = roh.pjp_id
						AND ovl.date = roh.date
					)
                ) subquery`

		insertQuery = `
            WITH getPjpIds AS (
                SELECT id
                FROM pjp_principles.permanent_journey_plans
                WHERE salesman_id = ?
                AND cust_id = ?
				AND approval_status = 'Approved'
            )
            INSERT INTO pjp_principles.outlet_visit_list (year, week, date, day, route_code, outlet_code, outlet_id, destination_type, pjp_id, pjp_code)
            SELECT
				roh.year,
				roh.week,
				roh.date,
				to_char(roh.date, 'Dy') as day,
				roh.route_code,
				roh.destination_code,
				roh.destination_id,
				roh.destination_type,
				roh.pjp_id,
				roh.pjp_code
			FROM pjp_principles.destinations_history roh
			WHERE roh.pjp_id IN (SELECT id FROM getPjpIds)
			AND roh.date = ?
			AND NOT EXISTS (
				SELECT 1 
				FROM pjp_principles.outlet_visit_list ovl
				WHERE ovl.outlet_id = roh.destination_id
				AND ovl.pjp_id = roh.pjp_id
				AND ovl.date = roh.date
			)`
		var count int64
		checkResult := tx.Raw(checkResultQuery, salesmanID, custId, date).Count(&count)

		if checkResult.Error != nil {
			tx.Rollback()
			helper.ErrorPanic(checkResult.Error)
		}

		insertResult := tx.Exec(insertQuery, salesmanID, custId, date)

		if insertResult.Error != nil {
			tx.Rollback()
			helper.ErrorPanic(insertResult.Error)
		}

		if err := tx.Commit().Error; err != nil {
			helper.ErrorPanic(err)
		}
		return
	}
	// distributor query
	checkResultQuery = `SELECT COUNT(*)
	FROM (
		WITH getPjpIds AS (
			SELECT id 
			FROM pjp.permanent_journey_plans 
			WHERE salesman_id = ?
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
		FROM pjp.route_pop_daily rpd
		LEFT JOIN pjp.route_outlet ro 
			ON ro.route_code = rpd.route_code 
			AND ro.pjp_id = rpd.pjp_id
		LEFT JOIN pjp.route_outlet_additional roa 
			ON roa.route_code = rpd.parent_route 
			AND roa.pjp_id = rpd.pjp_id
		JOIN pjp.route_outlet_history roh 
			ON roh.date = ?
			AND roh.route_code = rpd.route_code
			AND roh.pjp_id = rpd.pjp_id
			AND roh.outlet_id = CASE WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id ELSE ro.outlet_id END
		WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
		AND rpd.date = ?
		AND (ro.status = 'Approved' 
			OR roa.status = 'Approved' 
			OR ro.status = 'Approved With Propose' 
			OR roa.status = 'Approved With Propose')
	) subquery`

	insertQuery = `WITH getPjpIds AS (
		SELECT id 
		FROM pjp.permanent_journey_plans 
		WHERE salesman_id = ?
		AND cust_id = ?
	)
	INSERT INTO pjp.outlet_visit_list (year, week, date, day, route_code, outlet_code, outlet_id, pjp_id, pjp_code)
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
	FROM pjp.route_pop_daily rpd
	LEFT JOIN pjp.route_outlet ro 
		ON ro.route_code = rpd.route_code 
		AND ro.pjp_id = rpd.pjp_id
	LEFT JOIN pjp.route_outlet_additional roa 
		ON roa.route_code = rpd.parent_route 
		AND roa.pjp_id = rpd.pjp_id
	JOIN pjp.route_outlet_history roh 
		ON roh.date = ?
		AND roh.route_code = rpd.route_code
		AND roh.pjp_id = rpd.pjp_id
		AND roh.outlet_id = CASE WHEN rpd.parent_route IS NOT NULL THEN roa.outlet_id ELSE ro.outlet_id END
	WHERE rpd.pjp_id IN (SELECT id FROM getPjpIds)
	AND rpd.date = ?
	AND (ro.status = 'Approved' 
		OR roa.status = 'Approved' 
		OR ro.status = 'Approved With Propose' 
		OR roa.status = 'Approved With Propose')
	AND NOT EXISTS (
		SELECT 1
		FROM pjp.outlet_visit_list ovl
		WHERE ovl.outlet_id = ro.outlet_id
		AND ovl.pjp_id = ro.pjp_id
		AND ovl.date = rpd.date
		AND ovl.day = rpd.day
	)`

	var count int64
	checkResult := tx.Raw(checkResultQuery, salesmanID, custId, date, date).Count(&count)

	if checkResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(checkResult.Error)
	}

	insertResult := tx.Exec(insertQuery, salesmanID, custId, date, date)

	if insertResult.Error != nil {
		tx.Rollback()
		helper.ErrorPanic(insertResult.Error)
	}

	if err := tx.Commit().Error; err != nil {
		helper.ErrorPanic(err)
	}
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListMultiRowColumnAt(ctx context.Context, salesCode, custId, column string, currentTime int64, date string) {
	var distributorID *string
	var ids []int64

	err := repo.Db.Table("smc.m_customer").
		Select("distributor_id").
		Where("cust_id = ?", custId).
		Scan(&distributorID).Error
	if err != nil {
		helper.ErrorPanic(fmt.Errorf("failed to check customer status: %v", err))
	}

	isPrincipal := len(custId) == 6 && distributorID == nil
	var query string

	if isPrincipal {
		errFind := repo.Db.Table("pjp_principles.permanent_journey_plans").
			Select("id").
			Where("salesman_code = ? AND cust_id = ?", salesCode, custId).
			Find(&ids).Error

		if errFind != nil {
			helper.ErrorPanic(fmt.Errorf("failed when get pjp id by salescode and cust id: %v", err))
		}
		if len(ids) == 0 {
			helper.ErrorPanic(fmt.Errorf("pjp not found"))
		}

		query = `UPDATE pjp_principles.outlet_visit_list ovl
				SET ` + column + ` = ?
				WHERE ovl.pjp_id IN (?)
				AND ovl.date = ?`
	} else {
		errFind := repo.Db.Table("pjp.permanent_journey_plans").
			Select("id").
			Where("salesman_code = ? AND cust_id = ?", salesCode, custId).
			Find(&ids).Error

		if errFind != nil {
			helper.ErrorPanic(fmt.Errorf("failed when get pjp id by salescode and cust id: %v", err))
		}
		if len(ids) == 0 {
			helper.ErrorPanic(fmt.Errorf("pjp not found"))
		}

		query = `UPDATE pjp.outlet_visit_list ovl
				SET ` + column + ` = ?
				WHERE ovl.pjp_id IN (?)
				AND ovl.date = ?`
	}
	result := repo.Db.WithContext(ctx).Exec(query, currentTime, ids, date)
	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListColumnAt(ctx context.Context, salesmanCode, custId, column string, currentTime *int64, date string, id int64) {

	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp.outlet_visit_list
		SET `+column+` = ?
		WHERE date = ? AND id = ?
	`, currentTime, date, id)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListSkipColumnAt(ctx context.Context, column string, currentTime int64, date string, id int64, skipReson string, inOutlet bool, fileInfo model.OutletVisitList) {
	if fileInfo.FileUrl == "" {
		result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp.outlet_visit_list
		SET `+column+` = ?, skip_reason = ?, skip_in_outlet = ?
		WHERE date = ? AND id = ?
		`, currentTime, skipReson, inOutlet, date, id)
		helper.ErrorPanic(result.Error)
		return
	}

	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp.outlet_visit_list
		SET 
		    `+column+` = ?,
			skip_reason = ?,
			is_update_location = ?,
			file_name = ?,
			file_type = ?,
			media_category = ?,
			file_url = ?,
			file_size = ?,
			file_base64 = ?,
			photo_path = ?,
			folder = ?,
			latitude = ?,
			longitude = ?,
			allowed_radius = ?,
			distance_meter = ?,
			location_status = ?,
			skip_in_outlet = ?
		WHERE date = ? AND id = ?
	`,
		currentTime,
		skipReson,
		fileInfo.IsUpdateLocation,
		fileInfo.FileName,
		fileInfo.FileType,
		fileInfo.MediaCategory,
		fileInfo.FileUrl,
		fileInfo.FileSize,
		fileInfo.FileBase64,
		fileInfo.PhotoPath,
		fileInfo.Folder,
		fileInfo.Latitude,
		fileInfo.Longitude,
		fileInfo.AllowedRadius,
		fileInfo.DistanceMeter,
		fileInfo.LocationStatus,
		fileInfo.SkipInOutlet,
		date,
		id,
	)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) GetSummary(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data model.OutletVisitList, err error) {
	var query string
	if dataFilter.IsPrinciple {
		query = `WITH getPjpIds AS (
            SELECT id
            FROM pjp_principles.permanent_journey_plans
            WHERE salesman_code = ?
        )
        SELECT start, finish
        FROM pjp_principles.outlet_visit_list ovl
        WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
        AND ovl.date = ?`
	} else {
		query = `WITH getPjpIds AS (
            SELECT id
            FROM pjp.permanent_journey_plans
            WHERE salesman_code = ?
        )
        SELECT start, finish
        FROM pjp.outlet_visit_list ovl
        WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
        AND ovl.date = ?
		ORDER BY finish ASC`
	}
	result := repo.Db.WithContext(ctx).Raw(query, dataFilter.SalesmanCode, dataFilter.Date).Scan(&data)

	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) GetSummaryStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data response.VisitStatusResponse, err error) {
	var query string

	if dataFilter.IsPrinciple {
		query = `WITH getPjpIds AS (
            SELECT id
            FROM pjp_principles.permanent_journey_plans
            WHERE salesman_code = ?
        )
        SELECT
            COUNT(CASE
                WHEN is_planned = true
                THEN 1
                ELSE NULL
            END) AS planned,
            COUNT(CASE
                WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND skip_at IS NULL AND on_hold IS NULL AND leave_at IS NULL
                THEN 1
                ELSE NULL
            END) AS on_progress,
            COUNT(CASE
                WHEN skip_at IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS skipped,
            COUNT(CASE
                WHEN on_hold IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS on_hold,
            COUNT(CASE
                WHEN leave_at IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS finished,
            COUNT(CASE
                WHEN is_extra_call = true
                THEN 1
                ELSE NULL
            END) AS extra_call
        FROM
            pjp_principles.outlet_visit_list ovl
        WHERE
            ovl.pjp_id IN (SELECT id FROM getPjpIds)
            AND ovl.date = ?
        `
	} else {
		query = `WITH getPjpIds AS (
            SELECT id
            FROM pjp.permanent_journey_plans
            WHERE salesman_code = ?
        )
        SELECT
            COUNT(CASE
                WHEN is_planned = true
                THEN 1
                ELSE NULL
            END) AS planned,
            COUNT(CASE
                WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND skip_at IS NULL AND on_hold IS NULL AND leave_at IS NULL
                THEN 1
                ELSE NULL
            END) AS on_progress,
            COUNT(CASE
                WHEN skip_at IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS skipped,
            COUNT(CASE
                WHEN on_hold IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS on_hold,
            COUNT(CASE
                WHEN leave_at IS NOT NULL
                THEN 1
                ELSE NULL
            END) AS finished,
            COUNT(CASE
                WHEN is_extra_call = true
                THEN 1
                ELSE NULL
            END) AS extra_call
        FROM
            pjp.outlet_visit_list ovl
        WHERE
            ovl.pjp_id IN (SELECT id FROM getPjpIds)
            AND ovl.date = ?
        `
	}

	result := repo.Db.WithContext(ctx).Raw(query, dataFilter.SalesmanCode, dataFilter.Date).Scan(&data)
	if result.Error != nil {
		return data, errors.New("record not found")
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) FindAll(ctx context.Context, dataFilter entity.SummaryQueryFilter) (data []model.OutletVisitList, err error) {
	var query string
	if dataFilter.IsPrinciple {
		query = `
                    WITH getPjpIds AS (
                        SELECT id
                        FROM pjp_principles.permanent_journey_plans
                        WHERE salesman_id = ? AND cust_id = ?
                    )
                    SELECT
                        ovl.*,
						CASE WHEN destination_type = 'distributor'
							THEN md.distributor_name
							ELSE mo.outlet_name
						END as outlet_name,
						CASE WHEN destination_type = 'distributor'
							THEN md.address
							ELSE mo.address1 
						END as outlet_address,
						CASE WHEN destination_type = 'distributor'
							THEN md.longitude
							ELSE mo.longitude 
						END as outlet_longitude,
						CASE WHEN destination_type = 'distributor'
							THEN md.longitude
							ELSE mo.latitude 
						END AS outlet_latitude,
						COALESCE(mo.top, 0) as top,
                        CASE
                            WHEN ovl.start IS NOT NULL AND ovl.arrive_at IS NULL AND ovl.leave_at IS NULL AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL THEN 'planned'
                            WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL AND ovl.leave_at IS NULL THEN 'on_progress'
							WHEN ovl.skip_at IS NOT NULL AND ovl.arrive_at IS NOT NULL THEN 'visited'
                            WHEN ovl.skip_at IS NOT NULL AND ovl.leave_at IS NULL THEN 'skipped'
                            WHEN ovl.on_hold IS NOT NULL AND resume_at IS NULL AND ovl.skip_at IS NULL AND ovl.leave_at IS NULL THEN 'on_hold'
                            WHEN ovl.leave_at IS NOT NULL THEN 'finished'
                            ELSE 'planned'
                        END AS status
                    FROM pjp_principles.outlet_visit_list ovl
                    LEFT JOIN mst.m_outlet mo ON mo.outlet_id = ovl.outlet_id 
						AND mo.is_del = false 
						AND mo.verification_status = 1 
						AND mo.outlet_status != 4 
					LEFT JOIN mst.m_distributor md ON md.distributor_id = ovl.outlet_id 
						AND md.is_del = false 
						AND md.is_active = true
                    WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
                    AND ovl.date = ? ORDER BY ovl.start DESC NULLS LAST`
	} else {
		query = `
                WITH getPjpIds AS (
                    SELECT id
                    FROM pjp.permanent_journey_plans
                    WHERE salesman_id = ? AND cust_id = ?
                )
                SELECT
                    ovl.*,
                    mo.outlet_name,
					mo.address1 as outlet_address,
					mo.longitude as outlet_longitude,
					mo.latitude as outlet_latitude,
					mo.top,
                    CASE
                        WHEN ovl.start IS NOT NULL AND ovl.arrive_at IS NULL AND ovl.leave_at IS NULL AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL THEN 'planned'
                        WHEN (arrive_at IS NOT NULL OR resume_at IS NOT NULL) AND ovl.skip_at IS NULL AND ovl.on_hold IS NULL AND ovl.leave_at IS NULL THEN 'on_progress'
						WHEN ovl.skip_at IS NOT NULL AND ovl.arrive_at IS NOT NULL THEN 'visited'
                        WHEN ovl.skip_at IS NOT NULL AND ovl.leave_at IS NULL THEN 'skipped'
                        WHEN ovl.on_hold IS NOT NULL AND resume_at IS NULL AND ovl.skip_at IS NULL AND ovl.leave_at IS NULL THEN 'on_hold'
                        WHEN ovl.leave_at IS NOT NULL THEN 'finished'
                        ELSE 'planned'
                    END AS status
                FROM pjp.outlet_visit_list ovl
                LEFT JOIN (
                    SELECT DISTINCT ON (outlet_id) outlet_id, outlet_name, outlet_address
                    FROM pjp.route_outlet
                    ORDER BY outlet_id
                ) ro ON ovl.outlet_id = ro.outlet_id
				JOIN mst.m_outlet mo ON mo.outlet_id = ovl.outlet_id AND mo.is_del = false AND mo.verification_status = 1 AND mo.outlet_status != 4
                WHERE ovl.pjp_id IN (SELECT id FROM getPjpIds)
                AND ovl.date = ? ORDER BY ovl.start DESC NULLS LAST`

	}
	result := repo.Db.WithContext(ctx).Raw(query, dataFilter.EmpID, dataFilter.CustID, dataFilter.Date).Scan(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}

func (repo *OutletVisitRepoImpl) Delete(ctx context.Context, date string, week int, outletCode string, routeCode int) error {
	var outletVisitList model.OutletVisitList

	// First, try to find the record
	findResult := repo.Db.WithContext(ctx).
		Where("date = ? AND week = ? AND route_code = ? AND outlet_code = ?", date, week, routeCode, outletCode).
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
		Where("date = ? AND week = ? AND route_code = ? AND outlet_code = ?", date, week, routeCode, outletCode).
		Delete(&outletVisitList)

	// Check for any errors during deletion
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil
}

func (repo *OutletVisitRepoImpl) GetVisitsByDateAndSalesman(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitList, error) {
	var visits []model.OutletVisitList

	query := repo.Db.Table("pjp.outlet_visit_list").
		Select("id, year, week, date, day, route_code, outlet_id, outlet_code, pjp_id, pjp_code, start, finish, skip_at, leave_at, arrive_at, on_hold, resume_at, skip_reason")

	subquery := repo.Db.Table("pjp.permanent_journey_plans").
		Select("id").
		Where("salesman_id = ?", dataFilter.SalesmanId)
	query = query.Where("pjp.outlet_visit_list.pjp_id IN (?)", subquery)

	if dataFilter.Date != "" {
		query = query.Where("pjp.outlet_visit_list.date = ?", dataFilter.Date)
	}
	if dataFilter.Year != "" {
		query = query.Where("EXTRACT(YEAR FROM pjp.outlet_visit_list.date) = ?", dataFilter.Year)
	}
	if dataFilter.Month != "" {
		query = query.Where("EXTRACT(MONTH FROM pjp.outlet_visit_list.date) = ?", dataFilter.Month)
	}

	err := query.Find(&visits).Error
	return visits, err
}

func (repo *OutletVisitRepoImpl) CheckRouteOutletAdditional(ctx context.Context, routeCode *int, outletId int) (bool, error) {
	if routeCode == nil {
		return false, nil
	}

	var count int64
	err := repo.Db.Table("pjp.route_outlet_additional").
		Where("route_code = ? AND outlet_id = ?", *routeCode, outletId).
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
			data.RouteCode, data.OutletCode, data.PjpID, today, false).
		Delete(&data)

	if results.Error != nil {
		return results.Error
	}
	if results.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

// UpdateOutletVisitListWithFile updates outlet visit list with file information and location data
// Updates fields: is_update_location, file_name, file_type, media_category, file_url, file_size,
// file_base64, photo_path, folder, latitude, longitude
func (repo *OutletVisitRepoImpl) UpdateOutletVisitListWithFile(ctx context.Context, id int64, date string, fileInfo model.OutletVisitList) error {
	result := repo.Db.WithContext(ctx).Exec(`
		UPDATE pjp.outlet_visit_list
		SET 
			is_update_location = ?,
			file_name = ?,
			file_type = ?,
			media_category = ?,
			file_url = ?,
			file_size = ?,
			file_base64 = ?,
			photo_path = ?,
			folder = ?,
			latitude = ?,
			longitude = ?,
			allowed_radius = ?,
			distance_meter = ?,
			location_status = ?
		WHERE date = ? AND id = ?
	`,
		fileInfo.IsUpdateLocation,
		fileInfo.FileName,
		fileInfo.FileType,
		fileInfo.MediaCategory,
		fileInfo.FileUrl,
		fileInfo.FileSize,
		fileInfo.FileBase64,
		fileInfo.PhotoPath,
		fileInfo.Folder,
		fileInfo.Latitude,
		fileInfo.Longitude,
		fileInfo.AllowedRadius,
		fileInfo.DistanceMeter,
		fileInfo.LocationStatus,
		date,
		id,
	)

	return result.Error
}

func (repo *OutletVisitRepoImpl) UpdateByPjpIDAndDate(ctx context.Context, pjpID int, date string, data model.OutletVisitList) {
	result := repo.Db.WithContext(ctx).
		Where("pjp_id = ? AND date = ?", pjpID, date).
		Updates(&data)

	helper.ErrorPanic(result.Error)
}

func (repo *OutletVisitRepoImpl) InsertOrUpdateMobileVisit(ctx context.Context, custID, empCode string, outletID int, latitude, longitude float64, fileURL string, arriveAt int64, reason string) {
	var outletInfo struct {
		OutletCode string
		OutletName string
	}

	// Get Outlet Code and Name
	if err := repo.Db.WithContext(ctx).Table("mst.m_outlet").
		Select("outlet_code, outlet_name").
		Where("outlet_id = ?", outletID).
		Scan(&outletInfo).Error; err != nil {
		helper.ErrorPanic(err)
	}

	// Convert float to string for lat/long
	latStr := fmt.Sprintf("%f", latitude)
	longStr := fmt.Sprintf("%f", longitude)

	// Convert arriveAt (epoch millis) to time.Time
	createdAt := time.Unix(arriveAt/1000, (arriveAt%1000)*1000000)

	// Insert into mobile.visits
	query := `
		INSERT INTO mobile.visits (
			cust_id, emp_code, outlet_code, 
			type, created_at, 
			latitude, longitude, file_url, 
			is_in_outlet, reason
		) VALUES (
			?, ?, ?, 
			10, ?, 
			?, ?, ?, 
			?, ?
		)
	`
	isInOutlet := true
	if reason == "Skip" {
		isInOutlet = false
	}

	err := repo.Db.WithContext(ctx).Exec(query,
		custID, empCode, outletInfo.OutletCode,
		createdAt,
		latStr, longStr, fileURL, isInOutlet, reason,
	).Error

	if err != nil {
		helper.ErrorPanic(err)
	}
}

func (repo *OutletVisitRepoImpl) UpdateOutletVisitListByID(ctx context.Context, tx *gorm.DB, id int64, data model.OutletVisitList) error {
	err := tx.WithContext(ctx).Model(&model.OutletVisitList{}).
		Where("id = ?", id).
		Updates(&data).Error

	if err != nil {
		return err
	}
	return nil
}
