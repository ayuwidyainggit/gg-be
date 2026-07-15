package service

import (
	"context"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
)

type RoutePopService interface {
	SaveWeekly(ctx context.Context, request request.SaveWeeklyRequest, currentCustomerId string)
	SaveDelegateRoute(ctx context.Context, request request.SaveDelegateRequest, currentCustomerId string)
	CopyAllPermanentToDaily(ctx context.Context, request request.CopyAllRequest)
	CopyPartialPermanentToDaily(ctx context.Context, request request.CopyPartialRequest)
	CopyToSpecificDaily(ctx context.Context, request request.RoutesMapping, currentCustomerId string)
	CopyRouteDailyToDaily(ctx context.Context, request request.RoutesMapping, currentCustomerId string)
	FindAllPermanent(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopPermanentResponse
	FindAllDaily(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopDailyResponse
	FindByRouteOutletAdditional(ctx context.Context, code int, currentCustomerId string) response.RouteDetailResponse
	GetAllVisitDayMap(ctx context.Context, dataFilter entity.VisitDayMapQueryFilter, currentCustomerId string) (response []entity.VisitDayMapResponse)
	SaveOutletToRoute(ctx context.Context, request request.AddOutletToRouteRequest, currentCustomerId string, customerCode string)
	CancelOutletToRoute(ctx context.Context, request request.CancelAddOutletToRouteRequest)
}
