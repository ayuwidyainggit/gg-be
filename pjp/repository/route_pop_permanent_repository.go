package repository

import (
	"context"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/model"
)

type RoutePopPermanentRepository interface {
	InitTransaction(callback func() error) error
	Save(ctx context.Context, routePopPermanent model.RoutePopPermanent)
	FindByRouteCodeAndWeek(ctx context.Context, routeCode int, week int) (model.RoutePopPermanent, error)
	FindByWeek(ctx context.Context, week int) ([]model.RoutePopPermanent, error)
	FindByPjpCodes(ctx context.Context, pjpCode []int) ([]model.RoutePopPermanent, error)
	FindByPjpCode(ctx context.Context, pjpCode int) (model.RoutePopPermanent, error)
	FindAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []model.RoutePopPermanent
	UpdateOrCreate(ctx context.Context, data model.RoutePopPermanent)
	GetAllVisitDayMap(ctx context.Context, dataFilter entity.VisitDayMapQueryFilter, currentCustomerId string) []model.VisitDayMap
	CountOutletByRoute(ctx context.Context, currentCustomerId, startDate, endDate string) map[int]struct {
		TotalOutlet int
		RouteName   string
	}
	DeleteByRouteCode(ctx context.Context, code int) error
	DeleteByParams(
		ctx context.Context,
		routeCode int,
		pjpID int,
		pjpCode int,
		year int,
		week int,
		custId string,
	)

	SaveBulk(ctx context.Context, routePopPermanents []model.RoutePopPermanent) error
	FindByPjpID(ctx context.Context, pjpID int) ([]model.RoutePopPermanent, error)
}
