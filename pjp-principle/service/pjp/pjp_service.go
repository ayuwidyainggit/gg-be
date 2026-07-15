package pjp

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	routeoutlet "scyllax-pjp/repository/destination"
	"scyllax-pjp/repository/pjp"
	"scyllax-pjp/repository/route"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PjpService interface {
	Create(ctx context.Context, request request.PjpRequest, currentCustomerId string)
	GetAll(ctx context.Context, limit int, page int, filters map[string]interface{}, currentCustomerId string) ([]response.PjpResponse, response.Meta, error)
	GetPjpWithRoute(ctx context.Context, q string, custId string) []response.PjpResponse
	Update(ctx context.Context, request request.PjpRequest, currentCustomerId string)
	Delete(ctx context.Context, pjpId int, custId string)
	ListPjpApprove(ctx context.Context, q string, custId string) []response.PjpResponse
	GetById(ctx context.Context, pjpId int, currentCustomerId string) response.PjpResponse
}

type pjpService struct {
	pjpRepository         pjp.PjpRepository
	routeOutletRepository routeoutlet.DestinationRepository
	routeRepository       route.RouteRepository
	validate              *validator.Validate
	db                    *gorm.DB
}

func NewPjpService(pjpRepo pjp.PjpRepository, routeOutletRepository routeoutlet.DestinationRepository, routeRepository route.RouteRepository, validate *validator.Validate, db *gorm.DB) PjpService {
	return &pjpService{
		pjpRepository:         pjpRepo,
		routeOutletRepository: routeOutletRepository,
		routeRepository:       routeRepository,
		validate:              validate,
		db:                    db,
	}
}

func toPjpResponse(value model.Pjp) response.PjpResponse {
	statusBool, _ := strconv.ParseBool(value.Status)

	res := response.PjpResponse{
		PjpCode: helper.FormatPjpCode(value.PjpCode),
		Status:  statusBool,
	}
	helper.Automapper(value, &res)

	if value.RouteCode != 0 {
		res.Route = []response.RoutesEntity{{
			RouteCode:   value.RouteCode,
			TotalOutlet: value.TotalOutlet,
		}}
	}

	return res
}

func mapRequestToModel(req request.PjpRequest, customerId string) model.Pjp {
	return model.Pjp{
		ID:            req.ID,
		PjpCode:       req.PjpCode,
		TeamSalesMan:  req.TeamSalesMan,
		OperationType: req.OperationType,
		SalesManID:    req.SalesManID,
		SalesmanName:  req.SalesmanName,
		SalesmanCode:  req.SalesmanCode,
		WarehouseID:   req.WarehouseID,
		WarehouseName: req.WarehouseName,
		Status:        "false",
		PjpMode:       "manual",
		CustID:        customerId,
	}
}
