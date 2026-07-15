package routeoutlethistory

import (
	"context"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"time"

	"gorm.io/gorm"
)

func (repo *routeOutletHistoryRepository) FindByPjpIdToday(ctx context.Context, tx *gorm.DB, pjpIds []int, custId string) []model.RouteOutletHistory {
	var routesHistory []model.RouteOutletHistory

	// Ambil tanggal hari ini dalam format YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	result := tx.WithContext(ctx).
		Where("pjp_id IN ? AND cust_id = ? AND date = ?", pjpIds, custId, today).
		Find(&routesHistory)

	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return routesHistory
}
