package routeoutlethistory

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *routeOutletHistoryRepository) CreateBulk(ctx context.Context, tx *gorm.DB, outlets []model.RouteOutletHistory) {
	if len(outlets) == 0 {
		return
	}
	tx.WithContext(ctx).Create(&outlets)
}
