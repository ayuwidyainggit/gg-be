package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strings"

	"gorm.io/gorm"
)

type RoutePopPermanentRepositoryImpl struct {
	Db *gorm.DB
}

func NewRoutePopPermanentRepositoryImpl(Db *gorm.DB) RoutePopPermanentRepository {
	return &RoutePopPermanentRepositoryImpl{Db: Db}
}

func (repo *RoutePopPermanentRepositoryImpl) Save(ctx context.Context, routePopPermanent model.RoutePopPermanent) {
	result := repo.Db.WithContext(ctx).Create(&routePopPermanent)
	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopPermanentRepositoryImpl) FindByRouteCodeAndWeek(ctx context.Context, routeCode int, week int) (model.RoutePopPermanent, error) {
	var routePopDaily model.RoutePopPermanent

	result := repo.Db.WithContext(ctx).Where("route_code = ? AND week = ?", routeCode, week).First(&routePopDaily)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return routePopDaily, result.Error
	}

	return routePopDaily, nil
}

func (repo *RoutePopPermanentRepositoryImpl) FindByWeek(ctx context.Context, week int) ([]model.RoutePopPermanent, error) {
	var data []model.RoutePopPermanent
	result := repo.Db.WithContext(ctx).Where("week = ?", week).Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	if len(data) == 0 {
		return nil, errors.New("record not found")
	}

	return data, nil
}

func (repo *RoutePopPermanentRepositoryImpl) FindByPjpCodes(ctx context.Context, pjpCode []int) ([]model.RoutePopPermanent, error) {
	var data []model.RoutePopPermanent
	result := repo.Db.WithContext(ctx).Where("pjp_code IN (?)", pjpCode).Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	if len(data) == 0 {
		return nil, errors.New("pjp code not found")
	}

	return data, nil
}

func (repo *RoutePopPermanentRepositoryImpl) FindByPjpCode(ctx context.Context, pjpCode int) (model.RoutePopPermanent, error) {
	var data model.RoutePopPermanent
	result := repo.Db.WithContext(ctx).Where("pjp_code = ?", pjpCode).First(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}

func (repo *RoutePopPermanentRepositoryImpl) UpdateOrCreate(ctx context.Context, data model.RoutePopPermanent) {
	condition := model.RoutePopPermanent{
		RouteCode: data.RouteCode,
		PjpID:     data.PjpID,
		PjpCode:   data.PjpCode,
		Year:      data.Year,
		Week:      data.Week,
	}

	log.Printf("Condition: %+v", condition)

	result := repo.Db.WithContext(ctx).Where(condition).Assign(&data).FirstOrCreate(&data)
	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopPermanentRepositoryImpl) FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.RoutePopPermanent {
	var data []model.RoutePopPermanent

	query := repo.Db.Model(&data).Preload("Route").Preload("Pjp")
	fmt.Println(query.Statement.SQL.String())

	for field, value := range filters {
		switch field {
		case "salesman_name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Joins("JOIN pjp.permanent_journey_plans AS pjp_plans ON pjp.route_pop_permanent.pjp_id = pjp_plans.id").
					Where("pjp_plans.salesman_name = ?", v)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					query = query.Where(fmt.Sprintf("route_pop_permanent.%s = ?", field), v)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("route_pop_permanent.%s = ?", field), v)
				}
			}
		}
	}

	// subQuery := repo.Db.Table("route_oulet").
	// 	Select("route_code, COUNT(*) as total_outlet").
	// 	Where("cust_id", currentCustomerId).
	// 	Group("route_code")

	// query = query.Joins("LEFT JOIN (?) as ro on ro.route_outlet = route_pop_permanent.route_code", subQuery).
	// 	Select("route_pop_permanent.*, ro.total_outlet")

	// Custom join for Route to get route_name
	// query = query.Joins("JOIN pjp.routes ON pjp.route_pop_permanent.route_code = pjp.routes.route_code").
	// 	Select("pjp.route_pop_permanent.*, pjp.routes.route_name")

	result := query.WithContext(ctx).Where("route_pop_permanent.cust_id = ?", currentCustomerId).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *RoutePopPermanentRepositoryImpl) GetAllVisitDayMap(ctx context.Context, dataFilter entity.VisitDayMapQueryFilter, currentCustomerId string) []model.VisitDayMap {
	var data []model.VisitDayMap

	query := repo.Db.Table("pjp.route_pop_permanent").
		Select("pjp.route_pop_permanent.id, pjp.route_outlet.route_code,  pjp.route_outlet.route_name, pjp.route_pop_permanent.date, pjp.route_pop_permanent.week, COUNT(pjp.route_outlet.outlet_id) AS total_outlet").
		Joins("LEFT JOIN pjp.route_outlet ON pjp.route_pop_permanent.route_code = pjp.route_outlet.route_code").
		Group("pjp.route_pop_permanent.id, pjp.route_outlet.route_code, pjp.route_outlet.route_name, pjp.route_pop_permanent.date, pjp.route_pop_permanent.week")

	if dataFilter.PjpCode != 0 {
		query = query.Where("pjp.route_pop_permanent.pjp_code = ?", dataFilter.PjpCode)
	}

	if dataFilter.Sort != "" {
		sortBy := ""
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf("%s %s, ", colSort[0], colSort[1])
			}
		}

		sortBy = strings.TrimSuffix(sortBy, ", ")
		query = query.Order(sortBy)
	} else {
		query = query.Order("pjp.route_pop_permanent.id DESC")
	}

	result := query.WithContext(ctx).Where("route_pop_permanent.cust_id = ?", currentCustomerId).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

// func (repo *RoutePopPermanentRepositoryImpl) CountOutletByRoute(ctx context.Context, currentCustomerId string) map[int]any {
// 	var results []struct {
// 		RouteCode   int
// 		TotalOutlet int
// 		RouteName   string
// 	}

// 	repo.Db.Table("pjp.route_outlet").
// 		Select("pjp.route_outlet.route_code, COUNT(*) as total_outlet").
// 		Joins("JOIN pjp.routes ON pjp.route_outlet.route_code = pjp.routes.route_code").
// 		Where("cust_id = ?", currentCustomerId).
// 		Group("pjp.route_outlet.route_code, pjp.routes.route_name").
// 		Scan(&results)

// 	outletCount := make(map[int]any)
// 	for _, result := range results {
// 		outletCount[result.RouteCode] = result.TotalOutlet
// 	}

// 	return outletCount

// }

func (repo *RoutePopPermanentRepositoryImpl) CountOutletByRoute(ctx context.Context, currentCustomerId, startDate, endDate string) map[int]struct {
	TotalOutlet int
	RouteName   string
} {
	var results []struct {
		RouteCode   int    `gorm:"column:route_code"`
		TotalOutlet int    `gorm:"column:total_outlet"`
		RouteName   string `gorm:"column:route_name"`
	}

	// route_outlet --> permanent
	mainQuery := `
		SELECT pjp.route_outlet.route_code, COUNT(*) AS total_outlet, pjp.routes.route_name
		FROM pjp.route_outlet
		JOIN pjp.routes ON pjp.route_outlet.route_code = pjp.routes.route_code
		WHERE pjp.route_outlet.cust_id = ? AND (pjp.route_outlet.status = 'Approved' OR pjp.route_outlet.status = 'Approved With Propose')
		GROUP BY pjp.route_outlet.route_code, pjp.routes.route_name
	`

	// route_outlet_additional --> additional
	additionalQuery := `
		SELECT pjp.route_outlet_additional.route_code, COUNT(*) AS total_outlet, pjp.routes.route_name
		FROM pjp.route_outlet_additional
		JOIN pjp.routes ON pjp.route_outlet_additional.route_code = pjp.routes.route_code
		JOIN pjp.route_pop_daily ON pjp.route_outlet_additional.route_code = pjp.route_pop_daily.parent_route
		WHERE pjp.route_outlet_additional.cust_id = ? 
		  AND pjp.route_pop_daily.parent_route IS NOT NULL 
		  AND (pjp.route_outlet_additional.status = 'Approved' OR pjp.route_outlet_additional.status = 'Approved With Propose')
		  AND pjp.route_outlet_additional.date BETWEEN ? AND ?
		GROUP BY pjp.route_outlet_additional.route_code, pjp.routes.route_name
	`

	unionQuery := mainQuery + " UNION ALL " + additionalQuery

	err := repo.Db.Raw(unionQuery, currentCustomerId, currentCustomerId, startDate, endDate).Scan(&results).Error
	if err != nil {
		helper.ErrorPanic(err)
	}

	outletCount := make(map[int]struct {
		TotalOutlet int
		RouteName   string
	})
	for _, result := range results {
		if existing, found := outletCount[result.RouteCode]; found {
			// sum the total_outlet
			outletCount[result.RouteCode] = struct {
				TotalOutlet int
				RouteName   string
			}{
				TotalOutlet: existing.TotalOutlet + result.TotalOutlet,
				RouteName:   result.RouteName,
			}
		} else {
			// add a new entry to the map
			outletCount[result.RouteCode] = struct {
				TotalOutlet int
				RouteName   string
			}{
				TotalOutlet: result.TotalOutlet,
				RouteName:   result.RouteName,
			}
		}
	}

	return outletCount
}
func (repo *RoutePopPermanentRepositoryImpl) DeleteByRouteCode(ctx context.Context, code int) error {
	result := repo.Db.WithContext(ctx).Where("route_code = ?", code).Delete(&model.RoutePopPermanent{})

	if result.Error != nil {
		return errors.New("failed to delete route with the given code")
	}

	if result.RowsAffected == 0 {
		return errors.New("no route found with the given code")
	}

	return nil
}

func (repo *RoutePopPermanentRepositoryImpl) InitTransaction(callback func() error) error {
	err := callback()
	if err != nil {
		repo.Db.Rollback()
	}
	return err
}

func (repo *RoutePopPermanentRepositoryImpl) DeleteByParams(
	ctx context.Context,
	routeCode int,
	pjpID int,
	pjpCode int,
	year int,
	week int,
	custId string,
) {
	condition := map[string]interface{}{
		"route_code": routeCode,
		"pjp_id":     pjpID,
		"pjp_code":   pjpCode,
		"year":       year,
		"week":       week,
		"cust_id":    custId,
	}

	log.Printf("Delete condition: %+v", condition)

	result := repo.Db.WithContext(ctx).Where(condition).Delete(&model.RoutePopPermanent{})
	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopPermanentRepositoryImpl) SaveBulk(ctx context.Context, routePopPermanents []model.RoutePopPermanent) error {
	result := repo.Db.WithContext(ctx).Create(&routePopPermanents)
	return result.Error
}

func (repo *RoutePopPermanentRepositoryImpl) FindByPjpID(ctx context.Context, pjpID int) ([]model.RoutePopPermanent, error) {
	var data []model.RoutePopPermanent

	result := repo.Db.WithContext(ctx).Where("pjp_id = ?", pjpID).Find(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}
