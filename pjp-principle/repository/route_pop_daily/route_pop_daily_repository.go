package routepopdaily

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RoutePopDailyRepository interface {
	Create(ctx context.Context, tx *gorm.DB, routePopDaily model.RoutePopDaily)
}
type routePopDailyRepository struct{}

func NewRoutePopDailyRepositoryImpl() RoutePopDailyRepository {
	return &routePopDailyRepository{}
}
