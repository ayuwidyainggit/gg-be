package repository

import (
	"context"
	"scyllax-pjp/model"
)

type RoutePopDailyRepository interface {
	InitTransaction(callback func() error) error
	Insert(ctx context.Context, routePopDaily model.RoutePopDaily)
	Save(ctx context.Context, routePopDaily model.RoutePopDaily)
	FindByRouteCode(ctx context.Context, code int) (model.RoutePopDaily, error)
	//FindAllTes(ctx context.Context, filters map[string]interface{}) []model.RoutePopDailyWithOutlet
	FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.RoutePopDaily
	FindByParentRoute(ctx context.Context, code int, currentCustomerId string) ([]model.RoutePopDaily, error)
	UpdateOrCreate(ctx context.Context, data model.RoutePopDaily)
	UpdateOrCreateDaily(ctx context.Context, data model.RoutePopDaily)
	DeleteByRouteCode(ctx context.Context, code int) error
	DeleteByParams(
		ctx context.Context,
		routeCode int,
		pjpID int,
		pjpCode int,
		year int,
		week int,
		custId string,
	) error
}
