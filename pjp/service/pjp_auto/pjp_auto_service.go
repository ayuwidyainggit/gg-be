package pjpauto

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/repository"
	"scyllax-pjp/repository/pjp"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PjpAutoService interface {
	Create(ctx context.Context, request request.CreatePjpAuto, currentCustomerID string)
}

type pjpAutoService struct {
	pjpRepository         pjp.PjpRepository
	routeRepository       repository.RouteRepository
	routeOutletRepository repository.RouteOutletRepository
	validate              *validator.Validate
	db                    *gorm.DB
}

func NewPjpAutoService(pjpRepository pjp.PjpRepository, routeRepository repository.RouteRepository, routeOutletRepository repository.RouteOutletRepository, validate *validator.Validate, db *gorm.DB) PjpAutoService {
	return &pjpAutoService{
		pjpRepository:         pjpRepository,
		routeRepository:       routeRepository,
		routeOutletRepository: routeOutletRepository,
		validate:              validate,
		db:                    db,
	}
}
