package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	"scyllax-pjp/repository/pjp"
	"scyllax-pjp/repository/route"
	routeOutlet "scyllax-pjp/repository/route_outlet"
	routeOutletHistory "scyllax-pjp/repository/route_outlet_history"
	routePopPermanent "scyllax-pjp/repository/route_pop_permanent"
	"scyllax-pjp/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PjpEnhanceService interface {
	Create(ctx context.Context, request request.CreatePjpEnhanceRequest, currentCustomerId string)
	GetById(ctx context.Context, id int, currentCustomerId, parentCustomerId string) *response.PjpEnhanceResponse
	UpdatePjp(ctx context.Context, id int, request request.CreatePjpEnhanceRequest, currentCustomerId string)
	UpdateStatusPjp(ctx context.Context, id int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string)
	UpdateStatusByEmpId(ctx context.Context, EmpId int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string)
}

type pjpEnhanceService struct {
	pjpRepository                pjp.PjpRepository
	routeOutletRepository        routeOutlet.RouteOutletRepository
	routeOutletHistoryRepository routeOutletHistory.RouteOutletHistoryRepository
	routeRepository              route.RouteRepository
	routePopRepository           routePopPermanent.RoutePopPermanentRepository
	validate                     *validator.Validate
	db                           *gorm.DB
}

func NewPjpEnhanceService(pjpRepo pjp.PjpRepository, routeOutletRepository routeOutlet.RouteOutletRepository, routeOutletHistoryRepository routeOutletHistory.RouteOutletHistoryRepository, routeRepository route.RouteRepository, routePopRepository routePopPermanent.RoutePopPermanentRepository, validate *validator.Validate, db *gorm.DB) PjpEnhanceService {
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
	existingRoutes []model.Route,
	sequence int,
) model.Route {
	return model.Route{
		RouteCode: resolveRouteCode(routeReq, existingRoutes, sequence),
		RouteName: routeReq.RouteName,
		CustID:    custID,
		PjpID:     pjpID,
	}
}

func resolveRouteCode(routeReq request.RoutesCreatePjp, existingRoutes []model.Route, sequence int) int {
	if routeReq.RouteCode != nil && *routeReq.RouteCode != 0 {
		return *routeReq.RouteCode
	}

	requestedID := routeReq.ID
	if routeReq.RouteID != nil {
		requestedID = routeReq.RouteID
	}

	if requestedID != nil {
		for _, route := range existingRoutes {
			if route.ID == *requestedID {
				return route.RouteCode
			}
		}
	}

	for _, route := range existingRoutes {
		if route.Sequence == sequence {
			return route.RouteCode
		}
	}

	for _, route := range existingRoutes {
		if route.RouteName == routeReq.RouteName {
			return route.RouteCode
		}
	}

	return utils.GenerateCode(4)
}

func buildRouteOutlets(
	destinations []request.Destination,
	pjp model.Pjp,
	route model.Route,
	custID string,
) []model.RouteOutlet {
	var outlets []model.RouteOutlet

	for _, dest := range destinations {
		if dest.Type != "outlet" {
			continue
		}

		outlets = append(outlets, model.RouteOutlet{
			PjpID:         &pjp.ID,
			PjpCode:       &pjp.PjpCode,
			CustID:        custID,
			RouteCode:     route.RouteCode,
			RouteName:     route.RouteName,
			OutletID:      dest.ID,
			OutletCode:    dest.Code,
			OutletName:    dest.Name,
			Longitude:     dest.Longitude,
			Latitude:      dest.Latitude,
			AvgSalesWeek:  dest.AvgSalesWeek,
			OutletStatus:  dest.Status,
			OutletAddress: dest.Address,
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
