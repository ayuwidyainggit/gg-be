package routepop

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/repository"
	destinationhistory "scyllax-pjp/repository/destination_history"
	"scyllax-pjp/repository/pjp"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RoutePopService interface {
	SaveDailyRouteMap(ctx context.Context, request request.SaveDailyRouteMap, currentCustomerId string)
	GetByPjpAndRouteCode(ctx context.Context, pjpCode, routeCode int, date, currentCustomerId string) response.DailyRouteMap
	GetAll(ctx context.Context, filters map[string]interface{}, currentCustomerId string) []response.RoutePopPermanentResponse
}

type routePopService struct {
	pjpRepo                 pjp.PjpRepository
	routeRepository         repository.RouteRepository
	routePopDailyRepository repository.RoutePopDailyRepository
	routeOutletRepository   repository.RouteOutletRepository
	destinationHistoryRepo  destinationhistory.DestinationHistoryRepository
	validate                *validator.Validate
	db                      *gorm.DB
}

func NewRoutePopService(pjpRepo pjp.PjpRepository, destinationhistoryRepo destinationhistory.DestinationHistoryRepository, routePopDailyRepository repository.RoutePopDailyRepository, routeRepository repository.RouteRepository, routeOutletRepository repository.RouteOutletRepository, validate *validator.Validate, db *gorm.DB) RoutePopService {
	return &routePopService{
		pjpRepo:                 pjpRepo,
		destinationHistoryRepo:  destinationhistoryRepo,
		routePopDailyRepository: routePopDailyRepository,
		routeRepository:         routeRepository,
		routeOutletRepository:   routeOutletRepository,
		validate:                validate,
		db:                      db,
	}
}
