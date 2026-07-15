package service

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
)

type RouteService interface {
	Create(ctx context.Context, request request.CreateRouteRequest, currentCustomerId string)
	SaveOutlet(ctx context.Context, request request.SaveOutletRequest)
	SavePjp(ctx context.Context, request request.SavePjpRequest)
	DeleteOutlet(ctx context.Context, request request.DeleteOutletRequest)
	DeleteOutletAdditional(ctx context.Context, request request.DeleteOutletAdditionalRequest)
	UpdatePjp(ctx context.Context, request request.UpdatePjpInRouteRequest)
	FindAll(ctx context.Context, limit int, page int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingResponse, response.Meta, error)
	FindAllEnhance(ctx context.Context, limit int, page int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingEnhanceResponse, response.Meta, error)
	FindAllRoute(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RouteResponse
	FindByRouteOutlet(ctx context.Context, routeCode, pjpCode int) []response.RouteOutletsResponse
	Delete(ctx context.Context, routeId int)
	FindRouteByPjpCode(ctx context.Context, pjpCode, routeCode int) []response.RouteResponse
	FindDailyRouteByPjpCode(ctx context.Context, pjpCode, routeCode int, date string) []response.RouteDailyResponse
	UpdateRoute(ctx context.Context, request request.UpdateRoutesRequest, currentCustomerId string)
	SaveRouteConfirmation(ctx context.Context, request request.SaveRouteConfirmationRequest)
	DeletePjp(ctx context.Context, request request.DeletePjpRequest)
	UpdateNewRoute(ctx context.Context, request request.NewRouteRequest)
	FindByRouteCode(ctx context.Context, routeCode int) []response.RouteResponse
	DuplicateRoute(ctx context.Context, request request.DuplicateRoute, currentCustomerId string) error
}
