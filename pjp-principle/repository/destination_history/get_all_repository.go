package destinationhistory

import (
	"context"
	"fmt"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *destinationHistoryRepository) GetAll(ctx context.Context, tx *gorm.DB, filters map[string]interface{}) []model.DestinationHistory {
	var routesHistory []model.DestinationHistory

	query := tx.WithContext(ctx).Model(&model.DestinationHistory{})

	for key, value := range filters {
		if value != "" && value != nil {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	result := query.Find(&routesHistory)
	if result.Error != nil {
		helper.ErrorPanic(result.Error)
	}

	return routesHistory
}
