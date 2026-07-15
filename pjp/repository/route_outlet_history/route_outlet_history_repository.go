package routeoutlethistory

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type RouteOutletHistoryRepository interface {
	CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.RouteOutletHistory)
	FindByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.RouteOutletHistory
	FindByPjpIdToday(ctx context.Context, tx *gorm.DB, pjpIds []int, custId string) []model.RouteOutletHistory
	DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string)
	DeleteByVisitDay(ctx context.Context, tx *gorm.DB, history model.RouteOutletHistory)
}

type routeOutletHistoryRepository struct {
}

func NewRouteOutletHistoryRepository() RouteOutletHistoryRepository {
	return &routeOutletHistoryRepository{}
}
