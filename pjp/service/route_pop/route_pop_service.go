package routepop

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/repository"
	routeoutlethistory "scyllax-pjp/repository/route_outlet_history"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RoutePopService interface {
	SaveDailyRouteMap(ctx context.Context, request request.SaveDailyRouteMap, currentCustomerId string)
}

type routePopService struct {
	routeRepository         repository.RouteRepository
	routePopDailyRepository repository.RoutePopDailyRepository
	routeOutletRepository   repository.RouteOutletRepository
	routeOutletHistoryRepo  routeoutlethistory.RouteOutletHistoryRepository
	validate                *validator.Validate
	db                      *gorm.DB
}

func NewRoutePopService(routeOutletHistoryRepo routeoutlethistory.RouteOutletHistoryRepository, routePopDailyRepository repository.RoutePopDailyRepository, routeRepository repository.RouteRepository, routeOutletRepository repository.RouteOutletRepository, validate *validator.Validate, db *gorm.DB) RoutePopService {
	return &routePopService{
		routeOutletHistoryRepo:  routeOutletHistoryRepo,
		routePopDailyRepository: routePopDailyRepository,
		routeRepository:         routeRepository,
		routeOutletRepository:   routeOutletRepository,
		validate:                validate,
		db:                      db,
	}
}
