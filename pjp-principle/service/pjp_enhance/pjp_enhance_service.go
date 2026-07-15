package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	routeOutlet "scyllax-pjp/repository/destination"
	routeOutletHistory "scyllax-pjp/repository/destination_history"
	"scyllax-pjp/repository/pjp"
	"scyllax-pjp/repository/route"
	routePopPermanent "scyllax-pjp/repository/route_pop_permanent"
	"scyllax-pjp/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PjpEnhanceService interface {
	Create(ctx context.Context, request request.CreatePjpEnhanceRequest, currentCustomerId string)
	GetById(ctx context.Context, id int, currentCustomerId string) *response.PjpEnhanceResponse
	UpdatePjp(ctx context.Context, id int, request request.CreatePjpEnhanceRequest, currentCustomerId string)
	UpdateStatusPjp(ctx context.Context, id int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string)
	UpdateStatusByEmpId(ctx context.Context, EmpId int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string)
}

type pjpEnhanceService struct {
	pjpRepository                pjp.PjpRepository
	routeOutletRepository        routeOutlet.DestinationRepository
	routeOutletHistoryRepository routeOutletHistory.DestinationHistoryRepository
	routeRepository              route.RouteRepository
	routePopRepository           routePopPermanent.RoutePopPermanentRepository
	validate                     *validator.Validate
	db                           *gorm.DB
}

func NewPjpEnhanceService(pjpRepo pjp.PjpRepository, routeOutletRepository routeOutlet.DestinationRepository, routeOutletHistoryRepository routeOutletHistory.DestinationHistoryRepository, routeRepository route.RouteRepository, routePopRepository routePopPermanent.RoutePopPermanentRepository, validate *validator.Validate, db *gorm.DB) PjpEnhanceService {
	return &pjpEnhanceService{
		pjpRepository:                pjpRepo,
		routeOutletRepository:        routeOutletRepository,
		routeOutletHistoryRepository: routeOutletHistoryRepository,
		routeRepository:              routeRepository,
		routePopRepository:           routePopRepository,
		validate:                     validate,
		db:                           db,
	}
}

func buildPjpModel(request request.CreatePjpEnhanceRequest, currentCustomerId string) model.Pjp {
	return model.Pjp{
		PjpCode:        request.PjpCode,
		TeamSalesMan:   request.TeamSalesMan,
		OperationType:  request.OperationType,
		SalesManID:     request.SalesManID,
		SalesmanName:   request.SalesmanName,
		SalesmanCode:   request.SalesmanCode,
		WarehouseID:    request.WarehouseID,
		WarehouseName:  request.WarehouseName,
		Status:         "false",
		ApprovalStatus: request.ApprovalStatus,
		PjpMode:        "manual",
		CustID:         currentCustomerId,
	}
}

func buildRoute(
	routeReq request.RoutesCreatePjp,
	pjpID int,
	custID string,
) model.Route {
	return model.Route{
		RouteCode: utils.GenerateCode(4),
		RouteName: routeReq.RouteName,
		CustID:    custID,
		PjpID:     pjpID,
	}
}

func buildDestination(
	destinations []request.Destination,
	pjp model.Pjp,
	route model.Route,
	custID string,
) []model.Destination {
	var outlets []model.Destination

	for _, dest := range destinations {
		outlets = append(outlets, model.Destination{
			PjpID:              &pjp.ID,
			PjpCode:            &pjp.PjpCode,
			CustID:             custID,
			RouteCode:          route.RouteCode,
			RouteName:          route.RouteName,
			DestinationID:      dest.ID,
			DestinationCode:    dest.Code,
			DestinationName:    dest.Name,
			DestinationStatus:  dest.Status,
			DestinationAddress: dest.Address,
			DestinationType:    dest.Type,
			Longitude:          dest.Longitude,
			Latitude:           dest.Latitude,
			AvgSalesWeek:       dest.AvgSalesWeek,
		})
	}

	return outlets
}

func getRouteCodeByName(routes []model.Route, routeName string) int {
	for _, route := range routes {
		if route.RouteName == routeName {
			return route.RouteCode
		}
	}
	return 0 // fallback, bisa diganti error panic jika dianggap fatal
}
