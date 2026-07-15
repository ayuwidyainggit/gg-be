package destinationhistory

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type DestinationHistoryRepository interface {
	CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.DestinationHistory)
	FindByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string) []model.DestinationHistory
	FindByPjpCodeRouteCodeAndDate(ctx context.Context, tx *gorm.DB, pjpCode, routeCode int, date, custId string) []model.DestinationHistory
	DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custId string)
	GetAll(ctx context.Context, tx *gorm.DB, filters map[string]interface{}) []model.DestinationHistory
}

type destinationHistoryRepository struct {
}

func NewDestinationHistoryRepository() DestinationHistoryRepository {
	return &destinationHistoryRepository{}
}
