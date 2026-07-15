package routepoppermanent

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RoutePopPermanentRepository interface {
	CreateBulk(ctx context.Context, tx *gorm.DB, routePopPermanents []model.RoutePopPermanent)
	FindByPjpID(ctx context.Context, tx *gorm.DB, pjpID int, custId string) []model.RoutePopPermanent
}

type routePopPermanentRepository struct {
}

func NewRoutePopPermanentRepository() RoutePopPermanentRepository {
	return &routePopPermanentRepository{}
}
