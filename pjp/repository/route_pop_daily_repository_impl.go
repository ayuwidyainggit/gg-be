package repository

import (
	"context"
	"errors"
	"fmt"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoutePopDailyRepositoryImpl struct {
	Db *gorm.DB
}

func NewRoutePopDailyRepositoryImpl(Db *gorm.DB) RoutePopDailyRepository {
	return &RoutePopDailyRepositoryImpl{Db: Db}
}

func (repo *RoutePopDailyRepositoryImpl) Insert(ctx context.Context, routePopDaily model.RoutePopDaily) {
	result := repo.Db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "year"},
			{Name: "week"},
			{Name: "date"},
			{Name: "day"},
			{Name: "route_code"},
			{Name: "pjp_id"},
			{Name: "pjp_code"},
			{Name: "cust_id"},
			{Name: "status"},
		},
		UpdateAll: true,
	}).Create(&routePopDaily)

	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopDailyRepositoryImpl) Save(ctx context.Context, routePopDaily model.RoutePopDaily) {
	result := repo.Db.WithContext(ctx).Save(&routePopDaily)
	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopDailyRepositoryImpl) FindByRouteCode(ctx context.Context, code int) (model.RoutePopDaily, error) {
	var route model.RoutePopDaily

	result := repo.Db.WithContext(ctx).First(&route, "route_code = ?", code)

	if result.Error != nil {
		return route, errors.New("route code not found")
	}
	return route, nil
}

func (repo *RoutePopDailyRepositoryImpl) InitTransaction(callback func() error) error {
	err := callback()
	if err != nil {
		repo.Db.Rollback()
	}
	return err
}

func (repo *RoutePopDailyRepositoryImpl) FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.RoutePopDaily {
	var data []model.RoutePopDaily

	query := repo.Db.Model(&data).
		Preload("Pjp").
		Preload("Route.RouteOutlets")

	for field, value := range filters {
		switch field {
		case "salesman_name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Joins("JOIN pjp.permanent_journey_plans ON pjp.route_pop_daily.pjp_id = pjp.permanent_journey_plans.id").
					Where("permanent_journey_plans.salesman_name = ?", v)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					query = query.Where(fmt.Sprintf("%s = ?", field), v)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("%s = ?", field), v)
				}
			}
		}
	}

	result := query.WithContext(ctx).Where("cust_id = ?", currentCustomerId).Find(&data)
	helper.ErrorPanic(result.Error)
	return data
}

func (repo *RoutePopDailyRepositoryImpl) FindByParentRoute(ctx context.Context, code int, currentCustomerId string) ([]model.RoutePopDaily, error) {
	var route []model.RoutePopDaily
	result := repo.Db.WithContext(ctx).Preload("Route.RouteOutlets").Where("cust_id = ?", currentCustomerId).Where("parent_route = ?", code).Find(&route)

	if result.Error != nil {
		return route, result.Error
	}

	if result.RowsAffected == 0 {
		return route, errors.New("route pop daily not found")
	}

	return route, nil
}

func (repo *RoutePopDailyRepositoryImpl) UpdateOrCreate(ctx context.Context, data model.RoutePopDaily) {
	condition := model.RoutePopPermanent{
		RouteCode: data.RouteCode,
		PjpID:     data.PjpID,
		PjpCode:   data.PjpCode,
		Year:      data.Year,
		Week:      data.Week,
	}

	// result := repo.Db.WithContext(ctx).Where(condition).Assign(&data).FirstOrCreate(&data)

	var existingData model.RoutePopDaily
	var existingPermanent model.RoutePopPermanent
	result := repo.Db.WithContext(ctx).Where(condition).First(&existingData)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = repo.Db.WithContext(ctx).Create(&data)
		}
	} else {
		result = repo.Db.WithContext(ctx).Model(&existingData).Updates(&data)
	}
	helper.ErrorPanic(result.Error)

	result = repo.Db.WithContext(ctx).Where(condition).First(&existingPermanent)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			permanentData := model.RoutePopPermanent{
				RouteCode: data.RouteCode,
				PjpID:     data.PjpID,
				PjpCode:   data.PjpCode,
				Year:      data.Year,
				Week:      data.Week,
				Day:       data.Day,
				Date:      data.Date,
				CustID:    data.CustID,
			}
			result = repo.Db.WithContext(ctx).Create(&permanentData)
		}
	} else {
		permanentData := model.RoutePopPermanent{
			RouteCode: data.RouteCode,
			PjpID:     data.PjpID,
			PjpCode:   data.PjpCode,
			Year:      data.Year,
			Week:      data.Week,
			Day:       data.Day,
			Date:      data.Date,
			CustID:    data.CustID,
		}
		result = repo.Db.WithContext(ctx).Model(&existingPermanent).Updates(&permanentData)
	}

	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopDailyRepositoryImpl) UpdateOrCreateDaily(ctx context.Context, data model.RoutePopDaily) {
	condition := model.RoutePopDaily{
		RouteCode:   data.RouteCode,
		PjpID:       data.PjpID,
		PjpCode:     data.PjpCode,
		Year:        data.Year,
		Week:        data.Week,
		Status:      data.Status,
		ParentRoute: data.ParentRoute,
	}

	result := repo.Db.WithContext(ctx).Where(condition).Assign(&data).FirstOrCreate(&data)
	// result := repo.Db.WithContext(ctx).Create(&data)
	helper.ErrorPanic(result.Error)
}

func (repo *RoutePopDailyRepositoryImpl) DeleteByRouteCode(ctx context.Context, code int) error {
	result := repo.Db.WithContext(ctx).Where("route_code = ?", code).Delete(&model.RoutePopDaily{})

	if result.Error != nil {
		return errors.New("failed to delete route with the given code")
	}

	if result.RowsAffected == 0 {
		return errors.New("no route found with the given code")
	}

	return nil
}

func (repo *RoutePopDailyRepositoryImpl) DeleteByParams(
	ctx context.Context,
	routeCode int,
	pjpID int,
	pjpCode int,
	year int,
	week int,
	custId string,
) error {
	condition := map[string]interface{}{
		"route_code": routeCode,
		"pjp_id":     pjpID,
		"pjp_code":   pjpCode,
		"year":       year,
		"week":       week,
		"cust_id":    custId,
	}

	result := repo.Db.WithContext(ctx).Where(condition).Delete(&model.RoutePopDaily{})

	if result.Error != nil {
		return errors.New("failed to delete route with the given parameters")
	}

	if result.RowsAffected == 0 {
		return errors.New("no route found with the given parameters")
	}

	return nil
}
