package route

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RouteRepository interface {
	DeleteByRouteCodes(ctx context.Context, tx *gorm.DB, routeCodes []int, custId string)
	Create(ctx context.Context, tx *gorm.DB, route model.Route) model.Route
	FindByPjpID(ctx context.Context, tx *gorm.DB, pjpID int, custID string) []model.Route
	DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string)
}

type routeRepository struct{}

func NewRouteRepository() RouteRepository {
	return &routeRepository{}
}
