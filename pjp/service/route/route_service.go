package route

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/repository/pjp"
	routeoutlet "scyllax-pjp/repository/route_outlet"
	routepopdaily "scyllax-pjp/repository/route_pop_daily"
	routepoppermanent "scyllax-pjp/repository/route_pop_permanent"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RouteService interface {
	UpdateStatus(ctx context.Context, request request.UpdateStatusRequest, custId string)
	UpdateStatusEnhance(ctx context.Context, request request.UpdateStatusEnhanceRequest, custId string)
	GetAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingResponse, response.Meta, error)
	GetAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingEnhanceResponse, response.Meta, error)
}
type routeService struct {
	validate              *validator.Validate
	pjpRepo               pjp.PjpRepository
	routeOutletRepo       routeoutlet.RouteOutletRepository
	routePopPermanentRepo routepoppermanent.RoutePopPermanentRepository
	routePopDailyRepo     routepopdaily.RoutePopDailyRepository
	db                    *gorm.DB
}

func NewRouteService(validate *validator.Validate, pjpRepo pjp.PjpRepository, routeOutletRepo routeoutlet.RouteOutletRepository, routePopPermanentRepo routepoppermanent.RoutePopPermanentRepository, routePopDailyRepo routepopdaily.RoutePopDailyRepository, db *gorm.DB) RouteService {
	return &routeService{
		validate:              validate,
		pjpRepo:               pjpRepo,
		routeOutletRepo:       routeOutletRepo,
		routePopPermanentRepo: routePopPermanentRepo,
		routePopDailyRepo:     routePopDailyRepo,
		db:                    db,
	}
}
