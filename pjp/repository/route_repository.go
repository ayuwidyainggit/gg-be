package repository

import (
	"context"
	"scyllax-pjp/model"
)

type RouteRepository interface {
	FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.Route
	Insert(ctx context.Context, route model.Route) (model.Route, error)
	FindByRouteCode(ctx context.Context, code int) (model.Route, error)
	Update(ctx context.Context, route model.Route)
	FindById(ctx context.Context, routeId int, currentCustomerId string) (model.Route, error)
	Delete(ctx context.Context, routeId int) error
	DeleteByPjpId(ctx context.Context, pjpId int) error
	DeleteByRouteCode(ctx context.Context, routeCode int) error
	FindByPjpCode(ctx context.Context, pjpCode, routeCode int) (data []model.Route)
	FindByPjpCodeRouteAdditional(ctx context.Context, pjpCode, routeCode int, date string) (data []model.Route)
	QueryByRouteCode(ctx context.Context, routeCode int) (data []model.Route)
	FindAllByRouteCode(ctx context.Context, code int) ([]model.Route, error)
	FindRouteOutletByRouteCode(ctx context.Context, code int) (model.Route, error)
	FindByRouteCodes(ctx context.Context, routeCodes []int, custID string) []model.Route
	FindByPjpID(ctx context.Context, pjpID int, custID string) []model.Route
}
