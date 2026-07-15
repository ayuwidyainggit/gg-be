package repository

import (
	"context"
	"errors"
	"fmt"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RouteRepositoryImpl struct {
	Db *gorm.DB
}

func NewRouteRepositoryImpl(Db *gorm.DB) RouteRepository {
	return &RouteRepositoryImpl{Db: Db}
}

func (repo *RouteRepositoryImpl) Insert(ctx context.Context, route model.Route) (model.Route, error) {
	result := repo.Db.WithContext(ctx).Create(&route)
	if result.Error != nil {
		return model.Route{}, result.Error
	}

	return route, nil
}

func (repo *RouteRepositoryImpl) Update(ctx context.Context, route model.Route) {
	result := repo.Db.WithContext(ctx).Save(&route)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteRepositoryImpl) FindByRouteCode(ctx context.Context, code int) (model.Route, error) {
	var route model.Route

	result := repo.Db.WithContext(ctx).First(&route, "route_code = ?", code)

	if result.Error != nil {
		return route, errors.New("route code not found")
	}
	return route, nil
}

func (repo *RouteRepositoryImpl) FindRouteOutletByRouteCode(ctx context.Context, code int) (model.Route, error) {
	var route model.Route

	// Use Preload to load the related outlets
	result := repo.Db.WithContext(ctx).
		Preload("RouteOutlets").
		First(&route, "route_code = ?", code)

	if result.Error != nil {
		return route, errors.New("route code not found")
	}
	return route, nil
}

func (repo *RouteRepositoryImpl) FindById(ctx context.Context, routeId int, currentCustomerId string) (model.Route, error) {
	var route model.Route
	result := repo.Db.WithContext(ctx).Where("cust_id = ?", currentCustomerId).Find(&route, routeId)

	if result.Error != nil {
		return route, result.Error
	}

	if result.RowsAffected == 0 {
		return route, errors.New("route id is not found")
	}

	return route, nil
}

func (repo *RouteRepositoryImpl) Delete(ctx context.Context, routeId int) error {
	var route model.Route
	result := repo.Db.WithContext(ctx).Where("id = ?", routeId).Delete(&route)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("route is not found")
	}

	return nil
}

func (repo *RouteRepositoryImpl) DeleteByPjpId(ctx context.Context, pjpId int) error {
	var route model.Route
	result := repo.Db.WithContext(ctx).Where("pjp_id = ?", pjpId).Delete(&route)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("no route found for the given pjp_id")
	}

	return nil
}

func (repo *RouteRepositoryImpl) DeleteByRouteCode(ctx context.Context, routeCode int) error {
	var route model.Route
	result := repo.Db.WithContext(ctx).Where("route_code = ?", routeCode).Delete(&route)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("route is not found")
	}

	return nil
}

func (repo *RouteRepositoryImpl) FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.Route {
	var routes []model.Route

	query := repo.Db.Model(&routes).
		Preload("RouteOutlets")

	for field, value := range filters {
		switch v := value.(type) {
		case string:
			if v != "" {
				query = query.Where(fmt.Sprintf("%s = ?", field), v)
			}
		case int:
			if v != 0 {
				query = query.Where(fmt.Sprintf("%s = ?", field), v)
			}
		case bool:
			query = query.Where(fmt.Sprintf("%s = ?", field), v)
		}
	}

	result := query.WithContext(ctx).Where("pjp.routes.cust_id = ?", currentCustomerId).Find(&routes)
	helper.ErrorPanic(result.Error)

	return routes

}

func (repo *RouteRepositoryImpl) FindByPjpCode(ctx context.Context, pjpCode, routeCode int) (data []model.Route) {
	result := repo.Db.WithContext(ctx).
		Select(`
			pjp.routes.*, 
			pjp.route_outlet.outlet_id, 
			pjp.route_outlet.outlet_code, 
			pjp.route_outlet.outlet_name, 
			pjp.route_outlet.outlet_address, 
			pjp.route_outlet.outlet_status, 
			pjp.route_outlet.longitude, 
			pjp.route_outlet.latitude, 
			pjp.route_outlet.status,
			pjp.route_outlet.avg_sales_week,
			CASE
				WHEN EXISTS (
					SELECT 1 
					FROM pjp.route_outlet 
					WHERE pjp.route_outlet.route_code = ? AND pjp.route_outlet.pjp_code IS NOT NULL
				) THEN true
				ELSE false
			END AS is_assign_pjp
		`, routeCode).
		Joins("JOIN pjp.route_outlet ON pjp.routes.route_code = pjp.route_outlet.route_code").
		Where("pjp.route_outlet.pjp_code = ?", pjpCode).
		Where("pjp.route_outlet.route_code = ?", routeCode).
		Group("pjp.routes.id, pjp.route_outlet.outlet_id, pjp.route_outlet.outlet_code, pjp.route_outlet.outlet_name, pjp.route_outlet.outlet_address, pjp.route_outlet.outlet_status, pjp.route_outlet.longitude, pjp.route_outlet.latitude, pjp.route_outlet.status, pjp.route_outlet.avg_sales_week").
		Find(&data)

	helper.ErrorPanic(result.Error)

	return data
}

// func (repo *RouteRepositoryImpl) FindByPjpCodeRouteAdditional(ctx context.Context, pjpCode, routeCode int) (data []model.Route) {
// 	result := repo.Db.WithContext(ctx).
// 		Select(`
// 			pjp.routes.*,
// 			pjp.route_outlet.outlet_id,
// 			pjp.route_outlet.outlet_code,
// 			pjp.route_outlet.outlet_name,
// 			pjp.route_outlet.outlet_address,
// 			pjp.route_outlet.outlet_status,
// 			pjp.route_outlet.longitude,
// 			pjp.route_outlet.latitude,
// 			pjp.route_outlet.status,
// 			pjp.route_outlet_additional.outlet_id AS additional_outlet_id,
// 			pjp.route_outlet_additional.outlet_code AS additional_outlet_code,
// 			pjp.route_outlet_additional.outlet_name AS additional_outlet_name,
// 			pjp.route_outlet_additional.outlet_address AS additional_outlet_address,
// 			pjp.route_outlet_additional.outlet_status AS additional_outlet_status,
// 			pjp.route_outlet_additional.longitude AS additional_longitude,
// 			pjp.route_outlet_additional.latitude AS additional_latitude,
// 			pjp.route_outlet_additional.status AS additional_status,
// 			pjp.route_pop_daily.status AS route_pop_status
// 		`).
// 		Joins("JOIN pjp.route_outlet ON pjp.routes.route_code = pjp.route_outlet.route_code").
// 		Joins("LEFT JOIN pjp.route_outlet_additional ON pjp.routes.route_code = pjp.route_outlet_additional.route_code").
// 		Joins("LEFT JOIN pjp.route_pop_daily ON pjp.route_pop_daily.parent_route = pjp.routes.route_code").
// 		Where("pjp.route_pop_daily.pjp_code = ? OR pjp.route_outlet.pjp_code = ?", pjpCode, pjpCode).
// 		Where("pjp.routes.route_code = ?", routeCode).
// 		Group(`
// 			pjp.routes.id,
// 			pjp.route_outlet.outlet_id,
// 			pjp.route_outlet.outlet_code,
// 			pjp.route_outlet.outlet_name,
// 			pjp.route_outlet.outlet_address,
// 			pjp.route_outlet.outlet_status,
// 			pjp.route_outlet.longitude,
// 			pjp.route_outlet.latitude,
// 			pjp.route_outlet.status,
// 			pjp.route_outlet_additional.outlet_id,
// 			pjp.route_outlet_additional.outlet_code,
// 			pjp.route_outlet_additional.outlet_name,
// 			pjp.route_outlet_additional.outlet_address,
// 			pjp.route_outlet_additional.outlet_status,
// 			pjp.route_outlet_additional.longitude,
// 			pjp.route_outlet_additional.latitude,
// 			pjp.route_outlet_additional.status,
// 			pjp.route_pop_daily.status
// 		`).
// 		Find(&data)

// 	helper.ErrorPanic(result.Error)

// 	return data
// }

// func (repo *RouteRepositoryImpl) FindByPjpCodeRouteAdditional(ctx context.Context, pjpCode, routeCode int) (data []model.Route) {
// 	result := repo.Db.WithContext(ctx).
// 		Select(`
// 			pjp.routes.*,
// 			pjp.route_outlet.outlet_id,
// 			pjp.route_outlet.outlet_code,
// 			pjp.route_outlet.outlet_name,
// 			pjp.route_outlet.outlet_address,
// 			pjp.route_outlet.outlet_status,
// 			pjp.route_outlet.longitude,
// 			pjp.route_outlet.latitude,
// 			pjp.route_outlet.status AS outlet_status,
// 			pjp.route_outlet_additional.outlet_id AS additional_outlet_id,
// 			pjp.route_outlet_additional.outlet_code AS additional_outlet_code,
// 			pjp.route_outlet_additional.outlet_name AS additional_outlet_name,
// 			pjp.route_outlet_additional.outlet_address AS additional_outlet_address,
// 			pjp.route_outlet_additional.outlet_status AS additional_outlet_status,
// 			pjp.route_outlet_additional.longitude AS additional_longitude,
// 			pjp.route_outlet_additional.latitude AS additional_latitude,
// 			pjp.route_outlet_additional.status AS additional_status,
// 			COALESCE(pjp.route_pop_daily.status, 'permanent') AS route_pop_status
// 		`).
// 		Joins("JOIN pjp.route_outlet ON pjp.routes.route_code = pjp.route_outlet.route_code").
// 		Joins("LEFT JOIN pjp.route_outlet_additional ON pjp.routes.route_code = pjp.route_outlet_additional.route_code").
// 		Joins("LEFT JOIN pjp.route_pop_daily ON pjp.route_pop_daily.parent_route = pjp.routes.route_code AND pjp.route_pop_daily.pjp_code = ?", pjpCode).
// 		Where("pjp.route_outlet.pjp_code = ? OR pjp.route_outlet_additional.pjp_code = ?", pjpCode, pjpCode).
// 		Where("pjp.routes.route_code = ?", routeCode).
// 		Group(`
// 			pjp.routes.id,
// 			pjp.route_outlet.outlet_id,
// 			pjp.route_outlet.outlet_code,
// 			pjp.route_outlet.outlet_name,
// 			pjp.route_outlet.outlet_address,
// 			pjp.route_outlet.outlet_status,
// 			pjp.route_outlet.longitude,
// 			pjp.route_outlet.latitude,
// 			pjp.route_outlet.status,
// 			pjp.route_outlet_additional.outlet_id,
// 			pjp.route_outlet_additional.outlet_code,
// 			pjp.route_outlet_additional.outlet_name,
// 			pjp.route_outlet_additional.outlet_address,
// 			pjp.route_outlet_additional.outlet_status,
// 			pjp.route_outlet_additional.longitude,
// 			pjp.route_outlet_additional.latitude,
// 			pjp.route_outlet_additional.status,
// 			route_pop_status
// 		`).
// 		Find(&data)

// 	helper.ErrorPanic(result.Error)

// 	return data
// }

func (repo *RouteRepositoryImpl) FindByRouteCodes(ctx context.Context, routeCodes []int, custID string) []model.Route {
	var routes []model.Route

	result := repo.Db.WithContext(ctx).
		Where("route_code IN ? AND cust_id = ?", routeCodes, custID).
		Find(&routes)

	helper.ErrorPanic(result.Error)
	return routes
}

func (repo *RouteRepositoryImpl) FindByPjpCodeRouteAdditional(ctx context.Context, pjpCode, routeCode int, date string) []model.Route {
	var data []model.Route

	// Define the raw SQL query
	query := `
        -- Permanent Data
        SELECT
            r.route_code,
            ro.outlet_id,
            ro.outlet_code,
            ro.outlet_name,
            ro.outlet_address,
            ro.outlet_status,
            ro.longitude,
            ro.latitude,
			ro.avg_sales_week,
            'permanent' AS route_pop_status
        FROM
            pjp.routes r
        LEFT JOIN
            pjp.route_outlet ro ON r.route_code = ro.route_code
        WHERE
            r.route_code = ? AND ro.pjp_code = ? 
			AND (ro.status = 'Approved' OR ro.status = 'Approved With Propose')

        UNION ALL

        -- Additional Data (using parent_route)
        SELECT
            r.route_code,
            roa.outlet_id,
            roa.outlet_code,
            roa.outlet_name,
            roa.outlet_address,
            roa.outlet_status,
            roa.longitude,
            roa.latitude,
			roa.avg_sales_week,
            'additional' AS route_pop_status
        FROM
            pjp.routes r
        LEFT JOIN
            pjp.route_pop_daily rpd ON r.route_code = rpd.parent_route AND rpd.pjp_code = ?
        LEFT JOIN
            pjp.route_outlet_additional roa ON rpd.parent_route = roa.route_code AND roa.pjp_code = ?
        WHERE
            r.route_code = ? AND (roa.status = 'Approved' OR roa.status = 'Approved With Propose')
			AND roa.date = ?

    `

	// Execute the query with the parameters
	repo.Db.WithContext(ctx).
		Raw(query, routeCode, pjpCode, pjpCode, pjpCode, routeCode, date).
		Scan(&data)

	return data
}

func (repo *RouteRepositoryImpl) QueryByRouteCode(ctx context.Context, routeCode int) (data []model.Route) {
	result := repo.Db.WithContext(ctx).Raw(`
		WITH distinct_routes AS (
			SELECT 
				pjp.routes.route_code,
				COUNT(pjp.route_outlet.outlet_id) AS total_outlet
			FROM pjp.routes
			JOIN pjp.route_outlet ON pjp.routes.route_code = pjp.route_outlet.route_code
			WHERE pjp.route_outlet.route_code = ? AND pjp.route_outlet.pjp_code IS NULL
			GROUP BY pjp.routes.route_code
		)
		SELECT 
			pjp.routes.*,
			CASE WHEN EXISTS (
					SELECT 1 
					FROM pjp.route_outlet 
					WHERE route_code = ? AND pjp_code IS NOT NULL
				) THEN true
				ELSE false
			END AS is_assign_pjp,
			pjp.route_outlet.outlet_id,
			pjp.route_outlet.outlet_code,
			pjp.route_outlet.outlet_name,
			pjp.route_outlet.outlet_address,
			pjp.route_outlet.outlet_status,
			pjp.route_outlet.longitude,
			pjp.route_outlet.latitude,
			pjp.route_outlet.status,
			pjp.route_outlet.avg_sales_week,
			distinct_routes.total_outlet
		FROM distinct_routes
		JOIN pjp.routes ON distinct_routes.route_code = pjp.routes.route_code
		JOIN pjp.route_outlet ON pjp.routes.route_code = pjp.route_outlet.route_code
        WHERE pjp.route_outlet.pjp_code IS NULL
    `, routeCode, routeCode).Scan(&data)

	helper.ErrorPanic(result.Error)

	return data
}

func (repo *RouteRepositoryImpl) FindAllByRouteCode(ctx context.Context, code int) ([]model.Route, error) {
	var routeOutlets []model.Route

	query := `
		SELECT 
			outlet_id, 
			outlet_code, 
			outlet_name, 
			outlet_address, 
			outlet_status, 
			longitude, 
			latitude 
		FROM 
			pjp.route_outlet
		WHERE 
			route_code = ?
	`

	result := repo.Db.WithContext(ctx).Raw(query, code).Scan(&routeOutlets)

	helper.ErrorPanic(result.Error)

	return routeOutlets, nil
}

func (repo *RouteRepositoryImpl) FindByPjpID(ctx context.Context, pjpID int, custID string) []model.Route {
	var routes []model.Route

	result := repo.Db.WithContext(ctx).
		Where("pjp_id = ? AND cust_id = ?", pjpID, custID).
		Find(&routes)

	helper.ErrorPanic(result.Error)
	return routes
}
